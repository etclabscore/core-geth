// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Implements interaction with EVMC-based VMs.
// https://github.com/ethereum/evmc

package vm

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/evmc/v10/bindings/go/evmc"
	"github.com/ethereum/go-ethereum/params/vars"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

// EVMC represents the reference to a common EVMC-based VM instance and
// the current execution context as required by go-ethereum design.
type EVMC struct {
	instance *evmc.VM        // The reference to the EVMC VM instance.
	env      *EVM            // The execution context.
	cap      evmc.Capability // The supported EVMC capability (EVM or Ewasm)
	readOnly bool            // The readOnly flag (TODO: Try to get rid of it).
}

var (
	evmModule   *evmc.VM
	ewasmModule *evmc.VM
	evmcMux     sync.Mutex
)

func InitEVMCEVM(config string) {
	evmcMux.Lock()
	defer evmcMux.Unlock()
	if evmModule != nil {
		return
	}

	evmModule = initEVMC(evmc.CapabilityEVM1, config)
	log.Info("initialized EVMC interpreter", "path", config)
}

func InitEVMCEwasm(config string) {
	evmcMux.Lock()
	defer evmcMux.Unlock()
	if ewasmModule != nil {
		return
	}
	ewasmModule = initEVMC(evmc.CapabilityEWASM, config)
}

func initEVMC(cap evmc.Capability, config string) *evmc.VM {
	options := strings.Split(config, ",")
	path := options[0]

	if path == "" {
		panic("EVMC VM path not provided, set --vm.(evm|ewasm)=/path/to/vm")
	}

	instance, err := evmc.Load(path)
	if err != nil {
		panic(err.Error())
	}
	log.Info("EVMC VM loaded", "name", instance.Name(), "version", instance.Version(), "path", path)

	// Set options before checking capabilities.
	for _, option := range options[1:] {
		if idx := strings.Index(option, "="); idx >= 0 {
			name := option[:idx]
			value := option[idx+1:]
			err := instance.SetOption(name, value)
			if err == nil {
				log.Info("EVMC VM option set", "name", name, "value", value)
			} else {
				log.Warn("EVMC VM option setting failed", "name", name, "error", err)
			}
		}
	}

	if !instance.HasCapability(cap) {
		panic(fmt.Errorf("the EVMC module %s does not have requested capability %d", path, cap))
	}
	return instance
}

// hostContext implements evmc.HostContext interface.
type hostContext struct {
	env      *EVM      // The reference to the EVM execution context.
	contract *Contract // The reference to the current contract, needed by Call-like methods.
}

func (host *hostContext) AccessAccount(addr evmc.Address) evmc.AccessStatus {
	conf := host.env.ChainConfig()
	if !conf.IsEnabled(conf.GetEIP2929Transition, host.env.Context.BlockNumber) {
		return evmc.ColdAccess
	}

	if host.env.StateDB.AddressInAccessList(common.Address(addr)) {
		return evmc.WarmAccess
	}
	host.env.StateDB.AddAddressToAccessList(common.Address(addr))
	return evmc.ColdAccess
}

func (host *hostContext) AccessStorage(addr evmc.Address, key evmc.Hash) evmc.AccessStatus {
	conf := host.env.ChainConfig()
	if !conf.IsEnabled(conf.GetEIP2929Transition, host.env.Context.BlockNumber) {
		return evmc.ColdAccess
	}

	if addrOK, slotOK := host.env.StateDB.SlotInAccessList(common.Address(addr), common.Hash(key)); addrOK && slotOK {
		return evmc.WarmAccess
	}
	host.env.StateDB.AddSlotToAccessList(common.Address(addr), common.Hash(key))
	return evmc.ColdAccess
}

func (host *hostContext) AccountExists(evmcAddr evmc.Address) bool {
	addr := common.Address(evmcAddr)
	if host.env.ChainConfig().IsEnabled(host.env.ChainConfig().GetEIP161dTransition, host.env.Context.BlockNumber) /* EIP-161 */ {
		return !host.env.StateDB.Empty(addr)
	}
	return host.env.StateDB.Exist(addr)
}

func (host *hostContext) GetStorage(addr evmc.Address, evmcKey evmc.Hash) evmc.Hash {
	var value uint256.Int
	key := common.Hash(evmcKey)
	value.SetBytes(host.env.StateDB.GetState(common.Address(addr), key).Bytes())
	return evmc.Hash(value.Bytes32())
}

// setStorageLegacy implements the legacy (pre-EIP-1283 and pre-EIP-2200) storage setting.
// See core/vm/gas_table.go#gasSStore for reference.
func (host *hostContext) setStorageLegacy(original, current, value *uint256.Int) (status evmc.StorageStatus) {
	// This checks for 3 scenario's and calculates gas accordingly:
	//
	// 1. From a zero-value address to a non-zero value         (NEW VALUE)
	// 2. From a non-zero value address to a zero-value address (DELETE)
	// 3. From a non-zero to a non-zero
	switch {
	case current.IsZero() && !value.IsZero():
		return evmc.StorageAdded
	case !current.IsZero() && value.IsZero():
		host.env.StateDB.AddRefund(vars.SstoreRefundGas)
		return evmc.StorageDeleted
	default: // non 0 => non 0 (or 0 => 0)
		return evmc.StorageAssigned
	}
}

func (host *hostContext) setStorageEIP1283(original, current, value *uint256.Int) (status evmc.StorageStatus) {

	// conf := host.env.ChainConfig()
	// if host.contract.Gas <= vars.CallStipend && conf.IsEnabled(conf.GetEIP1706Transition, host.env.Context.BlockNumber) {
	// 	return evmc.StorageAssigned
	// }

	if current.Eq(value) {
		return evmc.StorageAssigned
	}

	// (X -> X) -> Y
	if original.Eq(current) {
		// 0 -> 0 -> Y
		if original.IsZero() {
			return evmc.StorageAdded
		}

		status = evmc.StorageModified

		// X -> X -> 0
		if value.IsZero() {
			host.env.StateDB.AddRefund(vars.NetSstoreClearRefund)
			return evmc.StorageDeleted
		}
		return status
	}
	if !original.IsZero() {
		// X -> 0 -> Z
		if current.IsZero() {
			host.env.StateDB.SubRefund(vars.NetSstoreClearRefund)
			// X -> 0 -> 0
		} else if value.IsZero() {
			host.env.StateDB.AddRefund(vars.NetSstoreClearRefund)
		}
	}
	if original.Eq(value) {
		if original.IsZero() {
			host.env.StateDB.AddRefund(vars.NetSstoreResetClearRefund)
		} else {
			host.env.StateDB.AddRefund(vars.NetSstoreResetRefund)
		}
	}
	return evmc.StorageAssigned
}

func (host *hostContext) setStorageEIP2200(original, current, value *uint256.Int) (status evmc.StorageStatus) {
	// Clause 1 is irrelevant:
	// 1. "If gasleft is less than or equal to gas stipend,
	//    fail the current call frame with ‘out of gas’ exception"
	// if host.contract.Gas <= vars.SstoreSentryGasEIP2200 {
	// 	return evmc.StorageAssigned
	// }
	//
	// // 2. "If current value equals new value (this is a no-op)"
	if current.Eq(value) {
		return evmc.StorageAssigned
	}

	// 3. "If current value does not equal new value"
	//
	// 3.1. "If original value equals current value
	//      (this storage slot has not been changed by the current execution context)"
	if original.Eq(current) {

		// 3.1.1 "If original value is 0"
		if original.IsZero() {
			return evmc.StorageAdded
		}

		// 3.1.2 "Otherwise"
		//   "SSTORE_RESET_GAS gas is deducted"
		status = evmc.StorageModified

		// "If new value is 0"
		if value.IsZero() {
			// "add SSTORE_CLEARS_SCHEDULE gas to refund counter"
			host.env.StateDB.AddRefund(vars.SstoreClearsScheduleRefundEIP2200)
			status = evmc.StorageDeleted
		}

		return status
	}

	// 3.2. "If original value does not equal current value
	//      (this storage slot is dirty),
	//      SLOAD_GAS gas is deducted.
	//      Apply both of the following clauses."

	// Because we need to apply "both following clauses"
	// we first collect information which clause is triggered
	// then assign status code to combination of these clauses.
	const (
		none                 = 0
		removeClearsSchedule = 1 << 0
		addClearsSchedule    = 1 << 1
		restoredBySet        = 1 << 2
		restoredByReset      = 1 << 3
	)
	var triggeredClauses = none

	// 3.2.1. "If original value is not 0"
	if !original.IsZero() {

		// 3.2.1.1. "If current value is 0"
		if current.IsZero() {

			// "(also means that new value is not 0)"
			//    assert(!is_zero(value));
			// "remove SSTORE_CLEARS_SCHEDULE gas from refund counter"
			host.env.StateDB.SubRefund(vars.SstoreClearsScheduleRefundEIP2200)
			triggeredClauses |= removeClearsSchedule
		}

		// 3.2.1.2. "If new value is 0"
		if value.IsZero() {
			// "(also means that current value is not 0)"
			//     assert(!is_zero(current));
			// "add SSTORE_CLEARS_SCHEDULE gas to refund counter"
			host.env.StateDB.AddRefund(vars.SstoreClearsScheduleRefundEIP2200)
			triggeredClauses |= addClearsSchedule
		}
	}

	// 3.2.2. "If original value equals new value (this storage slot is reset)"
	// Except: we use term 'storage slot restored'.
	if original.Eq(value) {

		// 3.2.2.1. "If original value is 0"
		if original.IsZero() {

			// "add SSTORE_SET_GAS - SLOAD_GAS to refund counter"
			// 20000 - 800
			host.env.StateDB.AddRefund(vars.SstoreSetGasEIP2200 - vars.SloadGasEIP2200)
			triggeredClauses |= restoredBySet
		} else

		// 3.2.2.2. "Otherwise"
		{
			// "add SSTORE_RESET_GAS - SLOAD_GAS gas to refund counter"
			host.env.StateDB.AddRefund(vars.SstoreResetGasEIP2200 - vars.SloadGasEIP2200)
			triggeredClauses |= restoredByReset
		}
	}

	switch triggeredClauses {
	case removeClearsSchedule:
		return evmc.StorageDeletedAdded
	case addClearsSchedule:
		return evmc.StorageModifiedDeleted
	case removeClearsSchedule | restoredByReset:
		return evmc.StorageDeletedRestored
	case restoredBySet:
		return evmc.StorageAddedDeleted
	case restoredByReset:
		return evmc.StorageModifiedRestored
	case none:
		return evmc.StorageAssigned
	default:
		panic("other combinations are impossible")
	}
}

func (host *hostContext) SetStorage(evmcAddr evmc.Address, evmcKey evmc.Hash, evmcValue evmc.Hash) (status evmc.StorageStatus) {
	addr := common.Address(evmcAddr)
	key := common.Hash(evmcKey)
	value := new(uint256.Int).SetBytes(evmcValue[:])

	var current, original = new(uint256.Int), new(uint256.Int)
	current.SetBytes(host.env.StateDB.GetState(addr, key).Bytes())
	original.SetBytes(host.env.StateDB.GetCommittedState(addr, key).Bytes())

	host.env.StateDB.SetState(addr, key, value.Bytes32())

	conf := host.env.ChainConfig()

	if conf.IsEnabled(conf.GetEIP2200Transition, host.env.Context.BlockNumber) {
		// >= Istanbul
		status = host.setStorageEIP2200(original, current, value)
	} else if conf.IsEnabled(conf.GetEIP1283Transition, host.env.Context.BlockNumber) &&
		!conf.IsEnabled(conf.GetEIP1283DisableTransition, host.env.Context.BlockNumber) {
		// == Constantinople
		status = host.setStorageEIP1283(original, current, value)
	} else {
		// == Legacy (< Istanbul && !Constantinople)
		status = host.setStorageLegacy(original, current, value)
	}

	return status
}

func (host *hostContext) GetBalance(addr evmc.Address) evmc.Hash {
	return evmc.Hash(common.BigToHash(host.env.StateDB.GetBalance(common.Address(addr))))
}

func (host *hostContext) GetCodeSize(addr evmc.Address) int {
	return host.env.StateDB.GetCodeSize(common.Address(addr))
}

func (host *hostContext) GetCodeHash(evmcAddr evmc.Address) evmc.Hash {
	addr := common.Address(evmcAddr)
	if host.env.StateDB.Empty(addr) {
		return evmc.Hash{}
	}
	return evmc.Hash(host.env.StateDB.GetCodeHash(addr))
}

func (host *hostContext) GetCode(addr evmc.Address) []byte {
	return host.env.StateDB.GetCode(common.Address(addr))
}

func (host *hostContext) Selfdestruct(evmcAddr evmc.Address, evmcBeneficiary evmc.Address) bool {
	addr := common.Address(evmcAddr)
	beneficiary := common.Address(evmcBeneficiary)
	db := host.env.StateDB
	if !db.HasSuicided(addr) {
		db.AddRefund(vars.SelfdestructRefundGas)
	}
	db.AddBalance(beneficiary, db.GetBalance(addr))
	return db.Suicide(addr)
}

func (host *hostContext) GetTxContext() evmc.TxContext {
	txCtx := evmc.TxContext{
		GasPrice:   evmc.Hash(common.BigToHash(host.env.GasPrice)),
		Origin:     evmc.Address(host.env.TxContext.Origin),
		Coinbase:   evmc.Address(host.env.Context.Coinbase),
		Number:     host.env.Context.BlockNumber.Int64(),
		Timestamp:  host.env.Context.Time.Int64(),
		GasLimit:   int64(host.env.Context.GasLimit),
		ChainID:    evmc.Hash(common.BigToHash(host.env.chainConfig.GetChainID())),
		BaseFee:    evmc.Hash{},
		PrevRandao: evmc.Hash(common.BigToHash(host.env.Context.Difficulty)),
	}
	conf := host.env.ChainConfig()
	if conf.IsEnabled(conf.GetEIP1559Transition, host.env.Context.BlockNumber) {
		txCtx.BaseFee = evmc.Hash(common.BigToHash(host.env.Context.BaseFee))
		if host.env.Context.Random != nil {
			txCtx.PrevRandao = evmc.Hash(*host.env.Context.Random)
		}
	}
	return txCtx
}

func (host *hostContext) GetBlockHash(number int64) evmc.Hash {
	b := host.env.Context.BlockNumber.Int64()
	if number >= (b-256) && number < b {
		return evmc.Hash(host.env.Context.GetHash(uint64(number)))
	}
	return evmc.Hash{}
}

func (host *hostContext) EmitLog(addr evmc.Address, evmcTopics []evmc.Hash, data []byte) {
	topics := make([]common.Hash, len(evmcTopics))
	for i, t := range evmcTopics {
		topics[i] = common.Hash(t)
	}
	host.env.StateDB.AddLog(&types.Log{
		Address:     common.Address(addr),
		Topics:      topics,
		Data:        data,
		BlockNumber: host.env.Context.BlockNumber.Uint64(),
	})
}

// Call executes a message call transaction.
// - evmcCodeAddress: https://github.com/ethereum/evmc/commit/8314761222c837d41f573b0af1152ce1c9895a32#diff-4c54ef259154e3cd3bd18b1963e027655525f909a887700cc2d1386e075c3628R155
func (host *hostContext) Call(kind evmc.CallKind,
	evmcRecipient evmc.Address, evmcSender evmc.Address, valueBytes evmc.Hash, input []byte, gas int64, depth int,
	static bool, saltBytes evmc.Hash, evmcCodeAddress evmc.Address) (output []byte, gasLeft int64, gasRefund int64, createAddrEvmc evmc.Address, err error) {

	recipient := common.Address(evmcRecipient)
	codeAddress := common.Address(evmcCodeAddress)

	var createAddr common.Address

	gasU := uint64(gas)
	var gasLeftU uint64

	value := new(uint256.Int).SetBytes(valueBytes[:])
	salt := big.NewInt(0).SetBytes(saltBytes[:])

	switch kind {
	case evmc.Call:
		if static {
			output, gasLeftU, err = host.env.StaticCall(host.contract, recipient, input, gasU)
		} else {
			output, gasLeftU, err = host.env.Call(host.contract, recipient, input, gasU, value.ToBig())
		}
	case evmc.DelegateCall:
		output, gasLeftU, err = host.env.DelegateCall(host.contract, codeAddress, input, gasU)
	case evmc.CallCode:
		output, gasLeftU, err = host.env.CallCode(host.contract, codeAddress, input, gasU, value.ToBig())
	case evmc.Create:
		var createOutput []byte
		createOutput, createAddr, gasLeftU, err = host.env.Create(host.contract, input, gasU, value.ToBig())
		createAddrEvmc = evmc.Address(createAddr)

		conf := host.env.ChainConfig()
		isHomestead := conf.IsEnabled(conf.GetEIP7Transition, host.env.Context.BlockNumber)
		if !isHomestead && err == ErrCodeStoreOutOfGas {
			err = nil
		}
		if err == ErrExecutionReverted {
			// Assign return buffer from REVERT.
			// TODO: Bad API design: return data buffer and the code is returned in the same place. In worst case
			//       the code is returned also when there is not enough funds to deploy the code.
			output = createOutput
		}
	case evmc.Create2:
		var createOutput []byte

		saltInt256 := new(uint256.Int)
		saltInt256.SetBytes(salt.Bytes())

		createOutput, createAddr, gasLeftU, err = host.env.Create2(host.contract, input, gasU, value.ToBig(), saltInt256)
		createAddrEvmc = evmc.Address(createAddr)
		if err == ErrExecutionReverted {
			// Assign return buffer from REVERT.
			// TODO: Bad API design: return data buffer and the code is returned in the same place. In worst case
			//       the code is returned also when there is not enough funds to deploy the code.
			output = createOutput
		}
	default:
		panic(fmt.Errorf("EVMC: Unknown call kind %d", kind))
	}

	// Map errors.
	if err == ErrExecutionReverted {
		err = evmc.Revert
	} else if err != nil {
		err = evmc.Failure
	}

	gasLeft = int64(gasLeftU)
	gasRefund = int64(host.env.StateDB.GetRefund())

	return output, gasLeft, gasRefund, createAddrEvmc, err
}

// getRevision translates ChainConfig's HF block information into EVMC revision.
func getRevision(env *EVM) evmc.Revision {
	n := env.Context.BlockNumber
	conf := env.ChainConfig()
	switch {
	// This is an example of choosing to use an "abstracted" idea
	// about chain config, where I'm choosing to prioritize "indicative" features
	// as identifiers for Fork-Feature-Groups. Note that this is very different
	// from using Feature-complete sets to assert "did Forkage."
	case conf.IsEnabled(conf.GetEIP1559Transition, n):
		return evmc.London
	case conf.IsEnabled(conf.GetEIP2565Transition, n):
		return evmc.Berlin
	case conf.IsEnabled(conf.GetEIP2200Transition, n):
		return evmc.Istanbul
	case conf.IsEnabled(conf.GetEIP145Transition, n) && (conf.IsEnabled(conf.GetEIP1283DisableTransition, n) || !conf.IsEnabled(conf.GetEIP1283Transition, n)):
		return evmc.Petersburg
	case conf.IsEnabled(conf.GetEIP145Transition, n) && (conf.IsEnabled(conf.GetEIP1283Transition, n) && !conf.IsEnabled(conf.GetEIP1283DisableTransition, n)):
		return evmc.Constantinople
	case conf.IsEnabled(conf.GetEIP198Transition, n):
		return evmc.Byzantium
	case conf.IsEnabled(conf.GetEIP155Transition, n):
		return evmc.SpuriousDragon
	case conf.IsEnabled(conf.GetEIP150Transition, n):
		return evmc.TangerineWhistle
	case conf.IsEnabled(conf.GetEIP7Transition, n):
		return evmc.Homestead
	default:
		return evmc.Frontier
	}
}

// Run implements Interpreter.Run().
func (evm *EVMC) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	evm.env.depth++
	defer func() { evm.env.depth-- }()

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}

	kind := evmc.Call
	if evm.env.StateDB.GetCodeSize(contract.Address()) == 0 {
		// Guess if this is a CREATE.
		kind = evmc.Create
	}

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This makes also sure that the readOnly flag isn't removed for child calls.
	if readOnly && !evm.readOnly {
		evm.readOnly = true
		defer func() { evm.readOnly = false }()
	}

	output, gasLeft, err := evm.instance.Execute(
		&hostContext{evm.env, contract},
		getRevision(evm.env),
		kind,
		evm.readOnly,
		evm.env.depth-1,
		int64(contract.Gas),
		evmc.Address(contract.Address()),
		evmc.Address(contract.Caller()),
		input,
		evmc.Hash(common.BigToHash(contract.value)),
		contract.Code,
	)

	contract.Gas = uint64(gasLeft)

	if err == evmc.Revert {
		err = ErrExecutionReverted
	} else if evmcError, ok := err.(evmc.Error); ok && evmcError.IsInternalError() {
		panic(fmt.Sprintf("EVMC VM internal error: %s", evmcError.Error()))
	}

	return output, err
}

// CanRun implements Interpreter.CanRun().
func (evm *EVMC) CanRun(code []byte) bool {
	required := evmc.CapabilityEVM1
	wasmPreamble := []byte("\x00asm")
	if bytes.HasPrefix(code, wasmPreamble) {
		required = evmc.CapabilityEWASM
	}
	return evm.cap == required
}
