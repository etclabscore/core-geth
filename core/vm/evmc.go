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
	if getRevision(host.env) < evmc.Berlin /* EIP-2929 */ {
		return evmc.ColdAccess
	}

	if host.env.StateDB.AddressInAccessList(common.Address(addr)) {
		return evmc.WarmAccess
	}
	return evmc.ColdAccess
}

func (host *hostContext) AccessStorage(addr evmc.Address, key evmc.Hash) evmc.AccessStatus {
	if getRevision(host.env) < evmc.Berlin /* EIP-2929 */ {
		return evmc.ColdAccess
	}

	if addrOK, slotOK := host.env.StateDB.SlotInAccessList(common.Address(addr), common.Hash(key)); addrOK && slotOK {
		return evmc.WarmAccess
	}
	return evmc.ColdAccess
}

func (host *hostContext) AccountExists(evmcAddr evmc.Address) bool {
	addr := common.Address(evmcAddr)
	if getRevision(host.env) >= evmc.SpuriousDragon /* EIP-161 */ {
		if !host.env.StateDB.Empty(addr) {
			return true
		}
	} else if host.env.StateDB.Exist(addr) {
		return true
	}
	return false
}

func (host *hostContext) GetStorage(addr evmc.Address, evmcKey evmc.Hash) evmc.Hash {
	var value uint256.Int
	key := common.Hash(evmcKey)
	value.SetBytes(host.env.StateDB.GetState(common.Address(addr), key).Bytes())
	return evmc.Hash(value.Bytes32())
}

func (host *hostContext) SetStorage(evmcAddr evmc.Address, evmcKey evmc.Hash, evmcValue evmc.Hash) (status evmc.StorageStatus) {
	addr := common.Address(evmcAddr)
	key := common.Hash(evmcKey)
	value := new(uint256.Int).SetBytes(evmcValue[:])
	var oldValue uint256.Int
	oldValue.SetBytes(host.env.StateDB.GetState(addr, key).Bytes())
	if oldValue.Eq(value) {
		return evmc.StorageAssigned
	}

	var current, original uint256.Int
	current.SetBytes(host.env.StateDB.GetState(addr, key).Bytes())
	original.SetBytes(host.env.StateDB.GetCommittedState(addr, key).Bytes())

	host.env.StateDB.SetState(addr, key, common.BytesToHash(value.Bytes()))

	// Here's a great example of one of the limits of our (core-geth) current chainconfig interface model.
	// Should we handle the logic here about historic-featuro logic (which really is nice, because when reading the strange-incantation implemations, it's nice to see why it is),
	// or should we handle the question of where we handle the rest of the questions like this, since this logic is
	// REALLY logic that belongs to the abstract idea of a chainconfiguration (aka chainconfig), which makes sense
	// but depends on ECIPs having steadier and more predictable logic.

	// isIstanbul was originally defined as two conditions: EIP1884, EIP2200;
	// but now the code expects any configuration to always apply (or have applied) them together at (as) the Istanbul fork.
	isIstanbul := getRevision(host.env) >= evmc.Istanbul

	if !isIstanbul {
		status = evmc.StorageModified
		if oldValue.IsZero() {
			return evmc.StorageAdded
		} else if value.IsZero() {
			host.env.StateDB.AddRefund(vars.SstoreRefundGas)
			return evmc.StorageDeleted
		}
		return evmc.StorageModified
	}

	// /**
	//  * The effect of an attempt to modify a contract storage item.
	//  *
	//  * See @ref storagestatus for additional information about design of this enum
	//  * and analysis of the specification.
	//  *
	//  * For the purpose of explaining the meaning of each element, the following
	//  * notation is used:
	//  * - 0 is zero value,
	//  * - X != 0 (X is any value other than 0),
	//  * - Y != 0, Y != X,  (Y is any value other than X and 0),
	//  * - Z != 0, Z != X, Z != X (Z is any value other than Y and X and 0),
	//  * - the "o -> c -> v" triple describes the change status in the context of:
	//  *   - o: original value (cold value before a transaction started),
	//  *   - c: current storage value,
	//  *   - v: new storage value to be set.
	//  *
	//  * The order of elements follows EIPs introducing net storage gas costs:
	//  * - EIP-2200: https://eips.ethereum.org/EIPS/eip-2200,
	//  * - EIP-1283: https://eips.ethereum.org/EIPS/eip-1283.
	//  */
	// enum evmc_storage_status
	// {
	// 	/**
	// 	 * The new/same value is assigned to the storage item without affecting the cost structure.
	// 	 *
	// 	 * The storage value item is either:
	// 	 * - left unchanged (c == v) or
	// 	 * - the dirty value (o != c) is modified again (c != v).
	// 	 * This is the group of cases related to minimal gas cost of only accessing warm storage.
	// 	 * 0|X   -> 0 -> 0 (current value unchanged)
	// 	 * 0|X|Y -> Y -> Y (current value unchanged)
	// 	 * 0|X   -> Y -> Z (modified previously added/modified value)
	// 	 *
	// 	 * This is "catch all remaining" status. I.e. if all other statuses are correctly matched
	// 	 * this status should be assigned to all remaining cases.
	// 	 */
	// 	EVMC_STORAGE_ASSIGNED = 0,
	status = evmc.StorageAssigned
	//
	// 	/**
	// 	 * A new storage item is added by changing
	// 	 * the current clean zero to a nonzero value.
	// 	 * 0 -> 0 -> Z
	// 	 */
	// 		EVMC_STORAGE_ADDED = 1,
	//
	if original.IsZero() && current.IsZero() && !value.IsZero() {
		status = evmc.StorageAdded
	}
	// /**
	//  * A storage item is deleted by changing
	//  * the current clean nonzero to the zero value.
	//  * X -> X -> 0
	//  */
	// 	EVMC_STORAGE_DELETED = 2,
	//
	if !original.IsZero() && !current.IsZero() && value.IsZero() {
		status = evmc.StorageDeleted
	}
	// /**
	//  * A storage item is modified by changing
	//  * the current clean nonzero to other nonzero value.
	//  * X -> X -> Z
	//  */
	// 	EVMC_STORAGE_MODIFIED = 3,
	//
	if !original.IsZero() && !current.IsZero() && !value.IsZero() {
		if original == current {
			status = evmc.StorageModified
		}
	}
	// /**
	//  * A storage item is added by changing
	//  * the current dirty zero to a nonzero value other than the original value.
	//  * X -> 0 -> Z
	//  */
	// 	EVMC_STORAGE_DELETED_ADDED = 4,
	//
	if !original.IsZero() && current.IsZero() && !value.IsZero() {
		if !original.Eq(value) {
			status = evmc.StorageDeletedAdded
		}
	}
	// /**
	//  * A storage item is deleted by changing
	//  * the current dirty nonzero to the zero value and the original value is not zero.
	//  * X -> Y -> 0
	//  */
	// 	EVMC_STORAGE_MODIFIED_DELETED = 5,
	//
	if !original.IsZero() && !current.IsZero() && value.IsZero() {
		if original != current {
			status = evmc.StorageModifiedDeleted
		}
	}
	// /**
	//  * A storage item is added by changing
	//  * the current dirty zero to the original value.
	//  * X -> 0 -> X
	//  */
	// 	EVMC_STORAGE_DELETED_RESTORED = 6,
	//
	if !original.IsZero() && current.IsZero() && !value.IsZero() {
		if original.Eq(value) {
			status = evmc.StorageDeletedRestored
		}
	}
	// /**
	//  * A storage item is deleted by changing
	//  * the current dirty nonzero to the original zero value.
	//  * 0 -> Y -> 0
	//  */
	// 	EVMC_STORAGE_ADDED_DELETED = 7,
	//
	if original.IsZero() && !current.IsZero() && value.IsZero() {
		status = evmc.StorageAddedDeleted
	}
	// /**
	//  * A storage item is modified by changing
	//  * the current dirty nonzero to the original nonzero value other than the current value.
	//  * X -> Y -> X
	//  */
	// 	EVMC_STORAGE_MODIFIED_RESTORED = 8
	// };
	//
	if !original.IsZero() && !current.IsZero() && !value.IsZero() {
		if original != current && original.Eq(value) {
			status = evmc.StorageModifiedRestored
		}
	}

	// if original == current {
	// 	// 0 -> 0 -> Z
	// 	if original.IsZero() { // create slot (2.1.1)
	// 		return evmc.StorageAdded
	// 	}
	//
	// 	if value.IsZero() { // delete slot (2.1.2b)
	// 		host.env.StateDB.AddRefund(vars.NetSstoreClearRefund)
	// 		return evmc.StorageDeleted
	// 	}
	// 	return evmc.StorageModified
	// }
	//
	// if !original.IsZero() {
	// 	if current.IsZero() { // recreate slot (2.2.1.1)
	// 		host.env.StateDB.SubRefund(vars.NetSstoreClearRefund)
	// 	} else if value.IsZero() { // delete slot (2.2.1.2)
	// 		host.env.StateDB.AddRefund(vars.NetSstoreClearRefund)
	// 	}
	// }
	// if original.Eq(value) {
	// 	if original.IsZero() { // reset to original inexistent slot (2.2.2.1)
	// 		host.env.StateDB.AddRefund(vars.NetSstoreResetClearRefund)
	// 	} else { // reset to original existing slot (2.2.2.2)
	// 		host.env.StateDB.AddRefund(vars.NetSstoreResetRefund)
	// 	}
	// }

	if current.Eq(value) {
		return
	}

	if original == current {
		if original.IsZero() {
			return
		}
		if value.IsZero() {
			host.env.StateDB.AddRefund(vars.NetSstoreClearRefund)
		}
		return
	}
	if !original.IsZero() {
		if current.IsZero() {
			host.env.StateDB.SubRefund(vars.NetSstoreClearRefund)
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
		GasPrice:  evmc.Hash(common.BigToHash(host.env.GasPrice)),
		Origin:    evmc.Address(host.env.TxContext.Origin),
		Coinbase:  evmc.Address(host.env.Context.Coinbase),
		Number:    host.env.Context.BlockNumber.Int64(),
		Timestamp: host.env.Context.Time.Int64(),
		GasLimit:  int64(host.env.Context.GasLimit),
		ChainID:   evmc.Hash(common.BigToHash(host.env.chainConfig.GetChainID())),
	}
	if getRevision(host.env) >= evmc.London {
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
	evmcDestination evmc.Address, evmcSender evmc.Address, valueBytes evmc.Hash, input []byte, gas int64, depth int,
	static bool, saltBytes evmc.Hash, evmcCodeAddress evmc.Address) (output []byte, gasLeft int64, gasRefund int64, createAddrEvmc evmc.Address, err error) {

	destination := common.Address(evmcDestination)
	codeTarget := common.Address(evmcCodeAddress)

	var createAddr common.Address

	gasU := uint64(gas)
	var gasLeftU uint64

	value := new(uint256.Int)
	value.SetBytes(valueBytes[:])

	salt := big.NewInt(0)
	salt.SetBytes(saltBytes[:])

	switch kind {
	case evmc.Call:
		if static {
			output, gasLeftU, err = host.env.StaticCall(host.contract, destination, input, gasU)
		} else {
			output, gasLeftU, err = host.env.Call(host.contract, destination, input, gasU, value.ToBig())
		}
	case evmc.DelegateCall:
		output, gasLeftU, err = host.env.DelegateCall(host.contract, codeTarget, input, gasU)
	case evmc.CallCode:
		output, gasLeftU, err = host.env.CallCode(host.contract, codeTarget, input, gasU, value.ToBig())
	case evmc.Create:
		var createOutput []byte
		createOutput, createAddr, gasLeftU, err = host.env.Create(host.contract, input, gasU, value.ToBig())
		createAddrEvmc = evmc.Address(createAddr)
		isHomestead := getRevision(host.env) >= evmc.Homestead
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
	gasRefund = gasLeft
	if getRevision(host.env) >= evmc.London {
		gasRefund = int64(host.env.StateDB.GetRefund())
	}
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
	case conf.IsEnabled(conf.GetEIP1884Transition, n):
		return evmc.Istanbul
	case conf.IsEnabled(conf.GetEIP1283DisableTransition, n):
		return evmc.Petersburg
	case conf.IsEnabled(conf.GetEIP145Transition, n):
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
