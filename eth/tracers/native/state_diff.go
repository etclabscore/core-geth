// Copyright 2022 The go-ethereum Authors
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

package native

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

func init() {
	register("stateDiffTracer", newStateDiffTracer)
}

const (
	// diffMarkerMemory= "_"	// temp state used while running the tracer, will never be returned to the user
	diffMarkerBorn    = "+"
	diffMarkerDied    = "-"
	diffMarkerChanged = "*"
	diffMarkerSame    = "="
)

type stateDiff = map[common.Address]*stateDiffAccount
type stateDiffAccount struct {
	empty  bool               `json:"-"`
	exists StateDiffExistance `json:"-"`
	// Balance string                      `json:"balance"`
	// Nonce   uint64                      `json:"nonce"`
	// Code    string                      `json:"code"`
	// Storage map[common.Hash]common.Hash `json:"storage"`
	Balance interface{}                            `json:"balance"` // Can be either string "=" or mapping "*" => {"from": "hex", "to": "hex"}
	Nonce   interface{}                            `json:"nonce"`
	Code    interface{}                            `json:"code"`
	Storage map[common.Hash]map[string]interface{} `json:"storage"`
}

type StateDiffExistance struct {
	From bool
	To   bool
}

type StateDiffBalance struct {
	From *hexutil.Big `json:"from"`
	To   *hexutil.Big `json:"to"`
}

type StateDiffCode struct {
	From hexutil.Bytes `json:"from"`
	To   hexutil.Bytes `json:"to"`
}

type StateDiffNonce struct {
	From hexutil.Uint64 `json:"from"`
	To   hexutil.Uint64 `json:"to"`
}

type StateDiffStorage struct {
	From common.Hash `json:"from"`
	To   common.Hash `json:"to"`
}

type stateDiffTracer struct {
	env          *vm.EVM
	ctx          *tracers.Context // Holds tracer context data
	stateDiff    stateDiff
	initialState *int
	create       bool
	to           common.Address
	interrupt    uint32 // Atomic flag to signal execution interruption
	reason       error  // Textual reason for the interruption
}

func newStateDiffTracer(ctx *tracers.Context) tracers.Tracer {
	// First callframe contains tx context info
	// and is populated on start and end.
	return &stateDiffTracer{stateDiff: stateDiff{}, ctx: ctx}
}

func (t *stateDiffTracer) CapturePreEVM(env *vm.EVM, inputs map[string]interface{}) {
}

func (t *stateDiffTracer) CapturePreEVM2(env *vm.EVM, inputs map[string]interface{}) {
	t.env = env
	if t.initialState == nil {
		snapshot := t.env.StateDB.Snapshot()
		t.initialState = &snapshot
	}
}

// CaptureStart implements the EVMLogger interface to initialize the tracing operation.
func (t *stateDiffTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	// t.env = env

	// if t.initialState == nil {
	// 	snapshot := t.env.StateDB.Snapshot()
	// 	t.initialState = &snapshot
	// }

	// t.create = create
	// t.to = to

	// Compute intrinsic gas
	// eip2f := env.ChainConfig().IsEnabled(env.ChainConfig().GetEIP2Transition, new(big.Int))
	// eip2028f := env.ChainConfig().IsEnabled(env.ChainConfig().GetEIP2028Transition, new(big.Int))
	// intrinsicGas, err := core.IntrinsicGas(input, nil, create, eip2f, eip2028f)
	// if err != nil {
	// 	return
	// }

	t.prepareAccount(from)
	t.prepareAccount(to)

	// t.lookupAccount(from)
	// t.lookupAccount(to)

	// The recipient balance includes the value transferred.
	// toBal := hexutil.MustDecodeBig(t.stateDiff[to].Balance)
	// toBal = new(big.Int).Sub(toBal, value)
	// t.stateDiff[to].Balance = hexutil.EncodeBig(toBal)

	// // The sender balance is after reducing: value, gasLimit, intrinsicGas.
	// // We need to re-add them to get the pre-tx balance.
	// fromBal := hexutil.MustDecodeBig(t.stateDiff[from].Balance)
	// gasPrice := env.TxContext.GasPrice
	// consumedGas := new(big.Int).Mul(
	// 	gasPrice,
	// 	new(big.Int).Add(
	// 		new(big.Int).SetUint64(intrinsicGas),
	// 		new(big.Int).SetUint64(gas),
	// 	),
	// )
	// fromBal.Add(fromBal, new(big.Int).Add(value, consumedGas))
	// t.stateDiff[from].Balance = hexutil.EncodeBig(fromBal)
	// t.stateDiff[from].Nonce--
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (t *stateDiffTracer) CaptureEnd(output []byte, gasUsed uint64, _ time.Duration, err error) {
	// if t.create {
	// 	// Exclude created contract.
	// 	delete(t.stateDiff, t.to)
	// }
}

// CaptureState implements the EVMLogger interface to trace a single step of VM execution.
func (t *stateDiffTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	stack := scope.Stack
	stackData := stack.Data()
	stackLen := len(stackData)
	switch {
	case stackLen >= 1 && (op == vm.SLOAD || op == vm.SSTORE):
		slot := common.Hash(stackData[stackLen-1].Bytes32())
		t.lookupStorage(scope.Contract.Address(), slot)
	case stackLen >= 1 && (op == vm.EXTCODECOPY || op == vm.EXTCODEHASH || op == vm.EXTCODESIZE || op == vm.BALANCE || op == vm.SELFDESTRUCT):
		addr := common.Address(stackData[stackLen-1].Bytes20())
		// t.lookupAccount(addr)
		t.prepareAccount(addr)
	case stackLen >= 5 && (op == vm.DELEGATECALL || op == vm.CALL || op == vm.STATICCALL || op == vm.CALLCODE):
		addr := common.Address(stackData[stackLen-2].Bytes20())
		// t.lookupAccount(addr)
		t.prepareAccount(addr)
	case op == vm.CREATE:
		addr := scope.Contract.Address()
		nonce := t.env.StateDB.GetNonce(addr)
		// t.lookupAccount(crypto.CreateAddress(addr, nonce))
		t.prepareAccount(crypto.CreateAddress(addr, nonce))
	case stackLen >= 4 && op == vm.CREATE2:
		offset := stackData[stackLen-2]
		size := stackData[stackLen-3]
		init := scope.Memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		inithash := crypto.Keccak256(init)
		salt := stackData[stackLen-4]
		// t.lookupAccount(crypto.CreateAddress2(scope.Contract.Address(), salt.Bytes32(), inithash))
		t.prepareAccount(crypto.CreateAddress2(scope.Contract.Address(), salt.Bytes32(), inithash))
	}
}

// CaptureFault implements the EVMLogger interface to trace an execution fault.
func (t *stateDiffTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, _ *vm.ScopeContext, depth int, err error) {
}

// CaptureEnter is called when EVM enters a new scope (via call, create or selfdestruct).
func (t *stateDiffTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

// CaptureExit is called when EVM exits a scope, even if the scope didn't
// execute any code.
func (t *stateDiffTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
}

// GetResult returns the json-encoded nested list of call traces, and any
// error arising from the encoding or forceful termination (via `Stop`).
func (t *stateDiffTracer) GetResult() (json.RawMessage, error) {
	// TODO: consider releasing snapshot
	// t.env.StateDB.DiscardSnapshot(t.initialState)
	var accountsToRemove []common.Address

	t.prepareAccount(t.env.Context.Coinbase)

	// read from latest state
	for addr, accountDiff := range t.stateDiff {
		accountDiff.empty = t.env.StateDB.Empty(addr)
		if accountDiff.empty {
			accountsToRemove = append(accountsToRemove, addr)
			continue
		}

		accountDiff.exists = StateDiffExistance{
			To: t.env.StateDB.Exist(addr),
		}

		accountDiff.Nonce = StateDiffNonce{
			To: (hexutil.Uint64)(t.env.StateDB.GetNonce(addr)),
		}

		accountDiff.Balance = StateDiffBalance{
			To: (*hexutil.Big)(t.env.StateDB.GetBalance(addr)),
		}

		accountDiff.Code = StateDiffCode{
			To: t.env.StateDB.GetCode(addr),
		}

		for key := range accountDiff.Storage {
			accountDiff.Storage[key]["*"] = StateDiffStorage{
				To: t.env.StateDB.GetState(addr, key),
			}
		}
	}

	// read from initial state
	t.env.StateDB.RevertToSnapshot(*t.initialState)
	for addr, accountDiff := range t.stateDiff {
		accountDiff.exists.From = t.env.StateDB.Exist(addr)

		if nonce, ok := accountDiff.Nonce.(StateDiffNonce); ok {
			nonce.From = (hexutil.Uint64)(t.env.StateDB.GetNonce(addr))
			accountDiff.Nonce = nonce
		}

		if bal, ok := accountDiff.Balance.(StateDiffBalance); ok {
			bal.From = (*hexutil.Big)(t.env.StateDB.GetBalance(addr))
			accountDiff.Balance = bal
		}

		if code, ok := accountDiff.Code.(StateDiffCode); ok {
			code.From = t.env.StateDB.GetCode(addr)
			accountDiff.Code = code
		}

		for key := range accountDiff.Storage {
			if storage, ok := accountDiff.Storage[key]["*"].(StateDiffStorage); ok {
				storage.From = t.env.StateDB.GetState(addr, key)
				accountDiff.Storage[key]["*"] = storage
			}
		}
	}

	// compare states
	// TODO: merge loop with the above one
	for addr, accountDiff := range t.stateDiff {
		// TODO: rename to isEmpty
		if accountDiff.empty {
			continue
		}

		// delete empty or unchanged storage items
		for key, accountStorage := range accountDiff.Storage {
			storage := accountStorage["*"].(StateDiffStorage)
			if storage.To == (common.Hash{}) || storage.From == storage.To {
				delete(accountDiff.Storage, key)
			}
		}

		// account removed
		if accountDiff.exists.From && !accountDiff.exists.To {
			{
				m := make(map[string]hexutil.Uint64)
				if nonce, ok := accountDiff.Nonce.(StateDiffNonce); ok {
					m["-"] = nonce.From
					accountDiff.Nonce = m
				}
			}
			{
				m := make(map[string]*hexutil.Big)
				if bal, ok := accountDiff.Balance.(StateDiffBalance); ok {
					m["-"] = bal.From
					accountDiff.Balance = m
				}
			}
			{
				m := make(map[string]hexutil.Bytes)
				if code, ok := accountDiff.Code.(StateDiffCode); ok {
					m["-"] = code.From
					accountDiff.Code = m
				}
			}

			// new account added
		} else if !accountDiff.exists.From && accountDiff.exists.To {
			{
				m := make(map[string]hexutil.Uint64)
				if nonce, ok := accountDiff.Nonce.(StateDiffNonce); ok {
					m["+"] = nonce.To
					accountDiff.Nonce = m
				}
			}
			{
				m := make(map[string]*hexutil.Big)
				if bal, ok := accountDiff.Balance.(StateDiffBalance); ok {
					m["+"] = bal.To
					accountDiff.Balance = m
				}
			}
			{
				m := make(map[string]hexutil.Bytes)
				if code, ok := accountDiff.Code.(StateDiffCode); ok {
					m["+"] = code.To
					accountDiff.Code = m
				}
			}

			for _, accountStorage := range accountDiff.Storage {
				storage := accountStorage["*"].(StateDiffStorage)
				delete(accountStorage, "*")
				accountStorage["+"] = &storage.To
			}

			// account changed
		} else if accountDiff.exists.From && accountDiff.exists.To {
			allEqual := len(accountDiff.Storage) == 0

			if nonce, ok := accountDiff.Nonce.(StateDiffNonce); ok {
				if nonce.From == nonce.To {
					accountDiff.Nonce = "="
				} else {
					m := make(map[string]StateDiffNonce)
					m["*"] = nonce
					accountDiff.Nonce = m
					allEqual = false
				}
			}
			if bal, ok := accountDiff.Balance.(StateDiffBalance); ok {
				if bal.To.ToInt().Cmp(bal.From.ToInt()) == 0 {
					accountDiff.Balance = "="
				} else {
					m := make(map[string]StateDiffBalance)
					m["*"] = bal
					accountDiff.Balance = m
					allEqual = false
				}
			}
			if code, ok := accountDiff.Code.(StateDiffCode); ok {
				if bytes.Equal(code.From, code.To) {
					accountDiff.Code = "="
				} else {
					m := make(map[string]StateDiffCode)
					m["*"] = code
					accountDiff.Code = m
					allEqual = false
				}
			}

			if allEqual {
				accountsToRemove = append(accountsToRemove, addr)
			}
		} else {
			accountsToRemove = append(accountsToRemove, addr)
		}
	}

	// remove marked accounts
	for _, addr := range accountsToRemove {
		delete(t.stateDiff, addr)
	}

	res, err := json.Marshal(t.stateDiff)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(res), t.reason
}

// Stop terminates execution of the tracer at the first opportune moment.
func (t *stateDiffTracer) Stop(err error) {
	t.reason = err
	atomic.StoreUint32(&t.interrupt, 1)
}

// func (t *stateDiffTracer) WriteAccountStorage(address common.Address, incarnation uint64, key *common.Hash, original, value *uint256.Int) error {
// 	if *original == *value {
// 		return nil
// 	}
// 	accountDiff := sd.sdMap[address]
// 	if accountDiff == nil {
// 		accountDiff = &StateDiffAccount{Storage: make(map[common.Hash]map[string]interface{})}
// 		sd.sdMap[address] = accountDiff
// 	}
// 	m := make(map[string]interface{})
// 	m["*"] = &StateDiffStorage{From: common.BytesToHash(original.Bytes()), To: common.BytesToHash(value.Bytes())}
// 	accountDiff.Storage[*key] = m
// 	return nil
// }

func (t *stateDiffTracer) prepareAccount(address common.Address) error {
	if _, ok := t.stateDiff[address]; !ok {
		t.stateDiff[address] = &stateDiffAccount{Storage: make(map[common.Hash]map[string]interface{})}
	}
	return nil
}

// lookupAccount fetches details of an account and adds it to the stateDiff
// if it doesn't exist there.
func (t *stateDiffTracer) lookupAccount(addr common.Address) {
	if _, ok := t.stateDiff[addr]; ok {
		return
	}
	t.stateDiff[addr] = &stateDiffAccount{
		Balance: bigToHex(t.env.StateDB.GetBalance(addr)),
		Nonce:   t.env.StateDB.GetNonce(addr),
		Code:    bytesToHex(t.env.StateDB.GetCode(addr)),
		Storage: make(map[common.Hash]map[string]interface{}),
	}
}

// lookupStorage fetches the requested storage slot and adds
// it to the stateDiff of the given contract. It assumes `lookupAccount`
// has been performed on the contract before.
func (t *stateDiffTracer) lookupStorage(addr common.Address, key common.Hash) {
	// if _, ok := t.stateDiff[addr].Storage[key]; ok {
	// 	return
	// }
	t.stateDiff[addr].Storage[key] = make(map[string]interface{})
}
