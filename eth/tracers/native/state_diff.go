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
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

func init() {
	register("stateDiffTracer", newStateDiffTracer)
}

type stateDiffMarker string

const (
	stateDiffMarkerBorn    stateDiffMarker = "+"
	stateDiffMarkerDied                    = "-"
	stateDiffMarkerChanged                 = "*"
	stateDiffMarkerSame                    = "="
)

type stateDiff = map[common.Address]*stateDiffAccount
type stateDiffAccount struct {
	marker  *stateDiffMarker                                `json:"-"`
	Balance interface{}                                     `json:"balance"`
	Nonce   interface{}                                     `json:"nonce"`
	Code    interface{}                                     `json:"code"`
	Storage map[common.Hash]map[stateDiffMarker]interface{} `json:"storage"`
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
	env              *vm.EVM
	ctx              *tracers.Context // Holds tracer context data
	stateDiff        stateDiff
	initialState     *state.StateDB
	create           bool
	to               common.Address
	accountsToRemove []common.Address
	interrupt        uint32 // Atomic flag to signal execution interruption
	reason           error  // Textual reason for the interruption
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
		t.initialState = t.env.StateDB.(*state.StateDB).Copy()
	}
}

// CaptureStart implements the EVMLogger interface to initialize the tracing operation.
func (t *stateDiffTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	t.create = create
	t.to = to

	var marker stateDiffMarker
	if create {
		marker = stateDiffMarkerBorn
	}

	t.initAccount(from, nil)
	t.initAccount(to, &marker)
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
	if err != nil {
		t.accountsToRemove = append(t.accountsToRemove, scope.Contract.Address())
	}

	stack := scope.Stack
	stackData := stack.Data()
	stackLen := len(stackData)
	switch {
	case stackLen >= 1 && (op == vm.SLOAD || op == vm.SSTORE):
		slot := common.Hash(stackData[stackLen-1].Bytes32())
		t.initStorageKey(scope.Contract.Address(), slot)
	case stackLen >= 1 && (op == vm.EXTCODECOPY || op == vm.EXTCODEHASH || op == vm.EXTCODESIZE || op == vm.BALANCE):
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.initAccount(addr, nil)
	case stackLen >= 5 && (op == vm.DELEGATECALL || op == vm.CALL || op == vm.STATICCALL || op == vm.CALLCODE):
		addr := common.Address(stackData[stackLen-2].Bytes20())
		t.initAccount(addr, nil)
	case op == vm.CREATE:
		addr := scope.Contract.Address()
		nonce := t.env.StateDB.GetNonce(addr)
		t.initAccount(crypto.CreateAddress(addr, nonce), nil)
	case stackLen >= 4 && op == vm.CREATE2:
		offset := stackData[stackLen-2]
		size := stackData[stackLen-3]
		init := scope.Memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		inithash := crypto.Keccak256(init)
		salt := stackData[stackLen-4]
		t.initAccount(crypto.CreateAddress2(scope.Contract.Address(), salt.Bytes32(), inithash), nil)
	case stackLen >= 1 && op == vm.SELFDESTRUCT:
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.initAccount(addr, nil)
		var marker stateDiffMarker
		marker = stateDiffMarkerDied
		t.initAccount(scope.Contract.Address(), &marker)
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
	t.initAccount(t.env.Context.Coinbase, nil)

	for addr, accountDiff := range t.stateDiff {
		// remove empty accounts
		if t.env.StateDB.Empty(addr) {
			t.accountsToRemove = append(t.accountsToRemove, addr)
			continue
		}

		// fmt.Println("initialState:", t.initialState)
		// fmt.Println("StateDB:", t.env.StateDB)
		initialExist := t.initialState.Exist(addr)
		exist := t.env.StateDB.Exist(addr)

		var marker stateDiffMarker
		if accountDiff.marker != nil {
			marker = *accountDiff.marker
		}

		if t.create && addr == t.to && marker == stateDiffMarkerDied {
			t.accountsToRemove = append(t.accountsToRemove, addr)
			continue
		}

		// handle storage keys
		var storageKeysToRemove []common.Hash

		// fill storage
		for key := range accountDiff.Storage {
			fromStorage := t.initialState.GetState(addr, key)
			toStorage := t.env.StateDB.GetState(addr, key)

			// mark unchanged storage items for deletion
			if toStorage == (common.Hash{}) || fromStorage == toStorage {
				storageKeysToRemove = append(storageKeysToRemove, key)
			} else if initialExist && exist {
				accountDiff.Storage[key][stateDiffMarkerChanged] = &StateDiffStorage{
					From: fromStorage,
					To:   toStorage,
				}
			} else if !initialExist && exist {
				accountDiff.Storage[key][stateDiffMarkerBorn] = toStorage
			} else if initialExist && !exist {
				accountDiff.Storage[key][stateDiffMarkerDied] = fromStorage
			}
		}

		// remove marked storage keys
		for _, key := range storageKeysToRemove {
			delete(accountDiff.Storage, key)
		}

		allEqual := len(accountDiff.Storage) == 0

		if !initialExist && exist {
			accountDiff.Nonce = map[stateDiffMarker]hexutil.Uint64{
				stateDiffMarkerBorn: hexutil.Uint64(t.env.StateDB.GetNonce(addr)),
			}
			accountDiff.Balance = map[stateDiffMarker]*hexutil.Big{
				stateDiffMarkerBorn: (*hexutil.Big)(t.env.StateDB.GetBalance(addr)),
			}
			accountDiff.Code = map[stateDiffMarker]hexutil.Bytes{
				stateDiffMarkerBorn: t.env.StateDB.GetCode(addr),
			}
		} else if initialExist && !exist || marker == stateDiffMarkerDied {
			accountDiff.Nonce = map[stateDiffMarker]hexutil.Uint64{
				stateDiffMarkerDied: hexutil.Uint64(t.initialState.GetNonce(addr)),
			}
			accountDiff.Balance = map[stateDiffMarker]*hexutil.Big{
				stateDiffMarkerDied: (*hexutil.Big)(t.initialState.GetBalance(addr)),
			}
			accountDiff.Code = map[stateDiffMarker]hexutil.Bytes{
				stateDiffMarkerDied: t.initialState.GetCode(addr),
			}
		} else if initialExist && exist {
			fromNonce := t.initialState.GetNonce(addr)
			toNonce := t.env.StateDB.GetNonce(addr)
			if fromNonce == toNonce {
				accountDiff.Nonce = stateDiffMarkerSame
			} else {
				m := make(map[stateDiffMarker]*StateDiffNonce)
				m[stateDiffMarkerChanged] = &StateDiffNonce{
					From: hexutil.Uint64(fromNonce),
					To:   hexutil.Uint64(toNonce),
				}
				accountDiff.Nonce = m
				allEqual = false
			}

			fromBalance := t.initialState.GetBalance(addr)
			toBalance := t.env.StateDB.GetBalance(addr)
			if fromBalance.Cmp(toBalance) == 0 {
				accountDiff.Balance = stateDiffMarkerSame
			} else {
				m := make(map[stateDiffMarker]*StateDiffBalance)
				m[stateDiffMarkerChanged] = &StateDiffBalance{From: (*hexutil.Big)(fromBalance), To: (*hexutil.Big)(toBalance)}
				accountDiff.Balance = m
				allEqual = false
			}

			fromCode := t.initialState.GetCode(addr)
			toCode := t.env.StateDB.GetCode(addr)
			if bytes.Equal(fromCode, toCode) {
				accountDiff.Code = stateDiffMarkerSame
			} else {
				m := make(map[stateDiffMarker]*StateDiffCode)
				m[stateDiffMarkerChanged] = &StateDiffCode{From: fromCode, To: toCode}
				accountDiff.Code = m
				allEqual = false
			}

			if allEqual {
				t.accountsToRemove = append(t.accountsToRemove, addr)
			}

		} else {
			t.accountsToRemove = append(t.accountsToRemove, addr)
		}
	}

	// remove marked accounts
	for _, addr := range t.accountsToRemove {
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

// initAccount stores the account address, in order we fetch the data in GetResult
func (t *stateDiffTracer) initAccount(address common.Address, marker *stateDiffMarker) error {
	if _, ok := t.stateDiff[address]; !ok {
		t.stateDiff[address] = &stateDiffAccount{
			marker:  marker,
			Storage: make(map[common.Hash]map[stateDiffMarker]interface{}),
		}
	} else {
		if marker != nil && *marker != "" {
			fmt.Println("acc address:", address, marker)
			t.stateDiff[address].marker = marker
		}
	}
	return nil
}

// initStorageKey stores the storage key in the account, in order we fetch the data in GetResult. It assumes `lookupAccount`
// has been performed on the contract before.
func (t *stateDiffTracer) initStorageKey(addr common.Address, key common.Hash) {
	t.stateDiff[addr].Storage[key] = make(map[stateDiffMarker]interface{})
}
