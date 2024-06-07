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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

func init() {
	tracers.DefaultDirectory.Register("stateDiffTracer", newStateDiffTracer, false)
}

type stateDiffMarker string

const (
	markerBorn    stateDiffMarker = "+"
	markerDied    stateDiffMarker = "-"
	markerChanged stateDiffMarker = "*"
	markerSame    stateDiffMarker = "="
)

type stateDiff = map[common.Address]*stateDiffAccount
type stateDiffAccount struct {
	marker  *stateDiffMarker                                `json:"-"`
	err     error                                           `json:"-"`
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
	env                *vm.EVM
	ctx                *tracers.Context // Holds tracer context data
	stateDiff          stateDiff
	initialState       *state.StateDB
	create             bool
	to                 common.Address
	accountsToRemove   []common.Address
	changedStorageKeys map[common.Address]map[common.Hash]bool
	interrupt          uint32 // Atomic flag to signal execution interruption
	reason             error  // Textual reason for the interruption
}

func (t *stateDiffTracer) CaptureTxStart(gasLimit uint64) {}

func (t *stateDiffTracer) CaptureTxEnd(restGas uint64) {}

func newStateDiffTracer(ctx *tracers.Context, j json.RawMessage) (tracers.Tracer, error) {
	// First callframe contains tx context info
	// and is populated on start and end.
	return &stateDiffTracer{stateDiff: stateDiff{}, ctx: ctx,
		changedStorageKeys: make(map[common.Address]map[common.Hash]bool)}, nil
}
func (t *stateDiffTracer) CapturePreEVM(env *vm.EVM) {
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
		marker = markerBorn
	}

	t.initAccount(from, nil)
	t.initAccount(to, &marker)
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (t *stateDiffTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {}

// CaptureState implements the EVMLogger interface to trace a single step of VM execution.
func (t *stateDiffTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	stack := scope.Stack
	stackData := stack.Data()
	stackLen := len(stackData)
	switch {
	case stackLen >= 1 && (op == vm.SLOAD || op == vm.SSTORE):
		addr := scope.Contract.Address()
		slot := common.Hash(stackData[stackLen-1].Bytes32())
		t.initStorageKey(addr, slot)

		// check if storage set/changed at least once
		if op == vm.SSTORE {
			if _, ok := t.changedStorageKeys[addr]; !ok {
				t.changedStorageKeys[addr] = make(map[common.Hash]bool)
			}

			isValueChanged, found := t.changedStorageKeys[addr][slot]
			if !found {
				t.changedStorageKeys[addr][slot] = false
			}

			if !isValueChanged {
				val := common.Hash(stackData[stackLen-2].Bytes32())
				if val != (common.Hash{}) {
					t.changedStorageKeys[addr][slot] = true
				}
			}
		}
	case stackLen >= 1 && (op == vm.EXTCODECOPY || op == vm.EXTCODEHASH || op == vm.EXTCODESIZE || op == vm.BALANCE):
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.initAccount(addr, nil)
	case stackLen >= 5 && (op == vm.DELEGATECALL || op == vm.CALL || op == vm.STATICCALL || op == vm.CALLCODE):
		addr := common.Address(stackData[stackLen-2].Bytes20())
		t.initAccount(addr, nil)
	case op == vm.CREATE:
		addr := scope.Contract.Address()
		nonce := t.env.StateDB.GetNonce(addr)
		marker := markerBorn
		t.initAccount(crypto.CreateAddress(addr, nonce), &marker)
	case stackLen >= 4 && op == vm.CREATE2:
		offset := stackData[stackLen-2]
		size := stackData[stackLen-3]
		init := scope.Memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		inithash := crypto.Keccak256(init)
		salt := stackData[stackLen-4]
		marker := markerBorn
		t.initAccount(crypto.CreateAddress2(scope.Contract.Address(), salt.Bytes32(), inithash), &marker)
	case stackLen >= 1 && op == vm.SELFDESTRUCT:
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.initAccount(addr, nil)

		// on SELFDESTRUCT mark the contract address as died
		marker := markerDied

		// account won't be SELFDESTRUCTed if out of gas happens on same instruction
		if err != nil && err.Error() == "out of gas" {
			marker = ""
		}
		t.initAccount(scope.Contract.Address(), &marker)
	}

	// log any account errors, in order we decide removal of accounts later
	if err != nil {
		if account, ok := t.stateDiff[scope.Contract.Address()]; ok {
			account.err = err
		}
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

		// read any special predefined marker set
		var marker stateDiffMarker
		if accountDiff.marker != nil {
			marker = *accountDiff.marker
		}

		hasDied := marker == markerDied

		// if an account has been Born within this run and also Died,
		// this means it will never be persisted to the state
		if t.create && addr == t.to && hasDied {
			t.accountsToRemove = append(t.accountsToRemove, addr)
			continue
		}

		// remove accounts with errors, except "out of gas"
		// though, when "out of gas", happens on new account creation, then we remove it as well
		if accountDiff.err != nil &&
			(accountDiff.err.Error() != "out of gas" || marker == markerBorn) {
			t.accountsToRemove = append(t.accountsToRemove, addr)
			continue
		}

		initialExist := t.initialState.Exist(addr)
		exist := t.env.StateDB.Exist(addr)

		// if initialState doesn't have the account (new account creation),
		// and hasDied, then account will be removed from state
		if !initialExist && hasDied {
			t.accountsToRemove = append(t.accountsToRemove, addr)
			continue
		}

		// handle storage keys
		var storageKeysToRemove []common.Hash

		// fill storage
		for key := range accountDiff.Storage {
			hasChanged := false
			if changedKeys, ok := t.changedStorageKeys[addr]; ok {
				if changed, ok := changedKeys[key]; ok && changed {
					hasChanged = true
				}
			}

			fromStorage := t.initialState.GetState(addr, key)
			toStorage := t.env.StateDB.GetState(addr, key)

			if initialExist && exist {
				// mark unchanged storage items for deletion
				if fromStorage == toStorage || (fromStorage == (common.Hash{}) && toStorage == (common.Hash{})) {
					storageKeysToRemove = append(storageKeysToRemove, key)
				} else {
					accountDiff.Storage[key][markerChanged] = &StateDiffStorage{
						From: fromStorage,
						To:   toStorage,
					}
				}
			} else if !initialExist && exist {
				if !hasChanged {
					storageKeysToRemove = append(storageKeysToRemove, key)
					continue
				}
				accountDiff.Storage[key][markerBorn] = toStorage
			} else if initialExist && !exist {
				accountDiff.Storage[key][markerDied] = fromStorage
			}
		}

		// remove marked storage keys
		for _, key := range storageKeysToRemove {
			delete(accountDiff.Storage, key)
		}

		allEqual := len(accountDiff.Storage) == 0

		// account creation
		if !initialExist && exist && !hasDied {
			accountDiff.Nonce = map[stateDiffMarker]hexutil.Uint64{
				markerBorn: hexutil.Uint64(t.env.StateDB.GetNonce(addr)),
			}
			accountDiff.Balance = map[stateDiffMarker]*hexutil.Big{
				markerBorn: (*hexutil.Big)(t.env.StateDB.GetBalance(addr).ToBig()),
			}
			accountDiff.Code = map[stateDiffMarker]hexutil.Bytes{
				markerBorn: t.env.StateDB.GetCode(addr),
			}

			// account has been removed
		} else if initialExist && !exist || hasDied {
			fromNonce := t.initialState.GetNonce(addr)
			accountDiff.Nonce = map[stateDiffMarker]hexutil.Uint64{
				markerDied: hexutil.Uint64(fromNonce),
			}
			accountDiff.Balance = map[stateDiffMarker]*hexutil.Big{
				markerDied: (*hexutil.Big)(t.initialState.GetBalance(addr).ToBig()),
			}
			accountDiff.Code = map[stateDiffMarker]hexutil.Bytes{
				markerDied: t.initialState.GetCode(addr),
			}

			// account changed
		} else if initialExist && exist {
			fromNonce := t.initialState.GetNonce(addr)
			toNonce := t.env.StateDB.GetNonce(addr)
			if fromNonce == toNonce {
				accountDiff.Nonce = markerSame
			} else {
				diff := make(map[stateDiffMarker]*StateDiffNonce)
				diff[markerChanged] = &StateDiffNonce{
					From: hexutil.Uint64(fromNonce),
					To:   hexutil.Uint64(toNonce),
				}
				accountDiff.Nonce = diff
				allEqual = false
			}

			fromBalance := t.initialState.GetBalance(addr)
			toBalance := t.env.StateDB.GetBalance(addr)
			if fromBalance.Cmp(toBalance) == 0 {
				accountDiff.Balance = markerSame
			} else {
				diff := make(map[stateDiffMarker]*StateDiffBalance)
				diff[markerChanged] = &StateDiffBalance{From: (*hexutil.Big)(fromBalance.ToBig()), To: (*hexutil.Big)(toBalance.ToBig())}
				accountDiff.Balance = diff
				allEqual = false
			}

			fromCode := t.initialState.GetCode(addr)
			toCode := t.env.StateDB.GetCode(addr)
			if bytes.Equal(fromCode, toCode) {
				accountDiff.Code = markerSame
			} else {
				diff := make(map[stateDiffMarker]*StateDiffCode)
				diff[markerChanged] = &StateDiffCode{From: fromCode, To: toCode}
				accountDiff.Code = diff
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
		// update the marker if account already inited
		if marker != nil && *marker != "" {
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
