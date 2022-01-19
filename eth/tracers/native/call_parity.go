// Copyright 2021 The go-ethereum Authors
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
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

var parityErrorMapping = map[string]string{
	"contract creation code storage out of gas": "Out of gas",
	"out of gas":                      "Out of gas",
	"gas uint64 overflow":             "Out of gas",
	"max code size exceeded":          "Out of gas",
	"invalid jump destination":        "Bad jump destination",
	"execution reverted":              "Reverted",
	"return data out of bounds":       "Out of bounds",
	"stack limit reached 1024 (1023)": "Out of stack",
	"precompiled failed":              "Built-in failed",
	"invalid input length":            "Built-in failed",
}

var parityErrorMappingStartingWith = map[string]string{
	"invalid opcode:": "Bad instruction",
	"stack underflow": "Stack underflow",
}

func init() {
	register("callParityTracer", NewCallParityTracer)
}

// callParityFrame is the result of a callParityTracerParity run.
type callParityFrame struct {
	Action              callTraceParityAction `json:"action"`
	BlockHash           *common.Hash          `json:"-"`
	BlockNumber         uint64                `json:"-"`
	Error               string                `json:"error,omitempty"`
	Result              callTraceParityResult `json:"result"`
	Subtraces           int                   `json:"subtraces"`
	TraceAddress        []int                 `json:"traceAddress"`
	TransactionHash     *common.Hash          `json:"-"`
	TransactionPosition *uint64               `json:"-"`
	Type                string                `json:"type"`
	Calls               []callParityFrame     `json:"-"`
}

// type callTraceParityAction struct {
// 	Author         string `json:"author,omitempty"`
// 	RewardType     string `json:"rewardType,omitempty"`
// 	SelfDestructed string `json:"address,omitempty"`
// 	Balance        string `json:"balance,omitempty"`
// 	CallType       string `json:"callType,omitempty"`
// 	CreationMethod string `json:"creationMethod,omitempty"`
// 	From           string `json:"from,omitempty"`
// 	Gas            string `json:"gas,omitempty"`
// 	Init           string `json:"init,omitempty"`
// 	Input          string `json:"input,omitempty"`
// 	RefundAddress  string `json:"refundAddress,omitempty"`
// 	To             string `json:"to,omitempty"`
// 	Value          string `json:"value,omitempty"`
// }

type callTraceParityAction struct {
	Author         *common.Address `json:"author,omitempty"`
	RewardType     *string         `json:"rewardType,omitempty"`
	SelfDestructed *common.Address `json:"address,omitempty"`
	Balance        *hexutil.Big    `json:"balance,omitempty"`
	CallType       string          `json:"callType,omitempty"`
	CreationMethod string          `json:"creationMethod,omitempty"`
	From           *common.Address `json:"from,omitempty"`
	Gas            hexutil.Uint64  `json:"gas,omitempty"`
	Init           *hexutil.Bytes  `json:"init,omitempty"`
	Input          *hexutil.Bytes  `json:"input,omitempty"`
	RefundAddress  *common.Address `json:"refundAddress,omitempty"`
	To             *common.Address `json:"to,omitempty"`
	Value          hexutil.Big     `json:"value,omitempty"`
}

type callTraceParityResult struct {
	Address *common.Address `json:"address,omitempty"`
	Code    *hexutil.Bytes  `json:"code,omitempty"`
	GasUsed hexutil.Uint64  `json:"gasUsed,omitempty"`
	Output  hexutil.Bytes   `json:"output,omitempty"`
}

// type callParityFrame struct {
// 	Type    string            `json:"type"`
// 	From    string            `json:"from"`
// 	To      string            `json:"to,omitempty"`
// 	Value   string            `json:"value,omitempty"`
// 	Gas     string            `json:"gas"`
// 	GasUsed string            `json:"gasUsed"`
// 	Input   string            `json:"input"`
// 	Output  string            `json:"output,omitempty"`
// 	Error   string            `json:"error,omitempty"`
// 	Calls   []callParityFrame `json:"calls,omitempty"`
// }

type callParityTracer struct {
	callstack         []callParityFrame
	interrupt         uint32           // Atomic flag to signal execution interruption
	reason            error            // Textual reason for the interruption
	activePrecompiles []common.Address // Updated on CaptureStart based on given rules
}

// NewCallParityTracer returns a native go tracer which tracks
// call frames of a tx, and implements vm.EVMLogger.
func NewCallParityTracer() tracers.Tracer {
	// First callParityframe contains tx context info
	// and is populated on start and end.
	t := &callParityTracer{callstack: make([]callParityFrame, 1)}
	return t
}

// isPrecompiled returns whether the addr is a precompile. Logic borrowed from newJsTracer in eth/tracers/js/tracer.go
func (t *callParityTracer) isPrecompiled(addr common.Address) bool {
	for _, p := range t.activePrecompiles {
		if p == addr {
			return true
		}
	}
	return false
}

func (l *callParityTracer) CapturePreEVM(env *vm.EVM, inputs map[string]interface{}) {}

func (t *callParityTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	// Skip any pre-compile invocations, those are just fancy opcodes
	t.activePrecompiles = env.ActivePrecompiles()

	inputHex := hexutil.Bytes(common.CopyBytes(input))

	t.callstack[0] = callParityFrame{
		Type: strings.ToLower(vm.CALL.String()),
		Action: callTraceParityAction{
			From:  &from,
			To:    &to,
			Input: &inputHex,
			Gas:   hexutil.Uint64(gas),
		},
		Result: callTraceParityResult{},
	}
	if value != nil {
		t.callstack[0].Action.Value = hexutil.Big(*value)
	}
	if create {
		t.callstack[0].Type = strings.ToLower(vm.CREATE.String())
	}
}

func (t *callParityTracer) CaptureEnd(output []byte, gasUsed uint64, _ time.Duration, err error) {
	if err != nil {
		t.callstack[0].Error = err.Error()
		if err.Error() == "execution reverted" && len(output) > 0 {
			t.callstack[0].Result.Output = hexutil.Bytes(common.CopyBytes(output))
		}
	} else {
		// TODO (ziogaschr): move back outside of if, makes sense to have it always. Is addition, no API breaks
		t.callstack[0].Result.GasUsed = hexutil.Uint64(gasUsed)

		t.callstack[0].Result.Output = hexutil.Bytes(common.CopyBytes(output))
	}
}

func (t *callParityTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

func (t *callParityTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, _ *vm.ScopeContext, depth int, err error) {
}

func (t *callParityTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	// Skip if tracing was interrupted
	if atomic.LoadUint32(&t.interrupt) > 0 {
		// TODO: env.Cancel()
		return
	}

	// Skip any pre-compile invocations, those are just fancy opcodes
	if t.isPrecompiled(to) && (typ == vm.CALL || typ == vm.STATICCALL) {
		return
	}

	inputHex := hexutil.Bytes(common.CopyBytes(input))

	call := callParityFrame{
		Type: strings.ToLower(typ.String()),
		Action: callTraceParityAction{
			From:  &from,
			To:    &to,
			Input: &inputHex,
			Gas:   hexutil.Uint64(gas),
			// Value: hexutil.Big(*value),
		},
		Result: callTraceParityResult{},
	}
	if value != nil {
		call.Action.Value = hexutil.Big(*value)
	}
	t.callstack = append(t.callstack, call)
}

func (t *callParityTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	size := len(t.callstack)
	if size <= 1 {
		return
	}
	// pop call
	call := t.callstack[size-1]
	t.callstack = t.callstack[:size-1]
	size -= 1

	call.Result.GasUsed = hexutil.Uint64(gasUsed)
	if err == nil {
		call.Result.Output = hexutil.Bytes(common.CopyBytes(output))
	} else {
		call.Error = err.Error()
		typ := vm.StringToOp(strings.ToUpper(call.Type))
		if typ == vm.CREATE || typ == vm.CREATE2 {
			call.Action.To = nil
		}
	}
	t.callstack[size-1].Calls = append(t.callstack[size-1].Calls, call)
}

func (t *callParityTracer) Finalize(call callParityFrame, traceAddress []int) ([]callParityFrame, error) {
	typ := vm.StringToOp(strings.ToUpper(call.Type))
	if typ == vm.CREATE || typ == vm.CREATE2 {
		t.formatCreateResult(&call)
	} else if typ == vm.SELFDESTRUCT {
		t.formatSuicideResult(&call)
	} else {
		t.formatCallResult(&call)
	}

	// for _, errorContains := range paritySkipTracesForErrors {
	// 	if strings.Contains(call.Error, errorContains) {
	// 		return
	// 	}
	// }

	t.convertErrorToParity(&call)

	if subtraces := len(call.Calls); subtraces > 0 {
		call.Subtraces = subtraces
	}

	call.TraceAddress = traceAddress

	results := []callParityFrame{call}

	for i := 0; i < len(call.Calls); i++ {
		childCall := call.Calls[i]

		var childTraceAddress []int
		childTraceAddress = append(childTraceAddress, traceAddress...)
		childTraceAddress = append(childTraceAddress, i)

		// Delegatecall uses the value from parent
		if (childCall.Type == "DELEGATECALL" || childCall.Type == "STATICCALL") && childCall.Action.Value.ToInt().Cmp(common.Big0) == 0 {
			childCall.Action.Value = call.Action.Value
		}

		child, err := t.Finalize(childCall, childTraceAddress)
		if err != nil {
			return nil, errors.New("failed to parse trace frame")
		}

		results = append(results, child...)
	}
	// fmt.Println("results:", results)

	return results, nil
}

func (t *callParityTracer) GetResult() (json.RawMessage, error) {
	if len(t.callstack) != 1 {
		return nil, errors.New("incorrect number of top-level calls")
	}

	traceAddress := []int{}
	result, err := t.Finalize(t.callstack[0], traceAddress)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(res), t.reason
}

func (t *callParityTracer) Stop(err error) {
	t.reason = err
	atomic.StoreUint32(&t.interrupt, 1)
}

func (t *callParityTracer) formatCreateResult(call *callParityFrame) {
	call.Action.CreationMethod = call.Type
	call.Type = strings.ToLower(vm.CREATE.String())

	input := call.Action.Input
	call.Action.Init = input
	call.Action.Input = nil

	to := call.Action.To
	call.Result.Address = to
	call.Action.To = nil

	output := call.Result.Output
	call.Result.Code = &output
	call.Result.Output = nil
}

func (t *callParityTracer) formatCallResult(call *callParityFrame) {
	call.Action.CallType = call.Type

	typ := vm.StringToOp(strings.ToUpper(call.Type))

	// update after callResult so as it affects only the root type
	if typ == vm.CALLCODE || typ == vm.DELEGATECALL || typ == vm.STATICCALL {
		call.Type = strings.ToLower(vm.CALL.String())
	}
}

func (t *callParityTracer) formatSuicideResult(call *callParityFrame) {
	call.Type = "suicide"

	addrFrom := call.Action.From
	call.Action.SelfDestructed = addrFrom
	call.Action.From = nil

	addrTo := call.Action.To
	call.Action.RefundAddress = addrTo
	call.Action.Balance = &call.Action.Value

	call.Action.Input = nil

	call.Result = callTraceParityResult{}
}

func (t *callParityTracer) convertErrorToParity(call *callParityFrame) {
	if call.Error == "" {
		return
	}

	if parityError, ok := parityErrorMapping[call.Error]; ok {
		call.Error = parityError
		call.Result = callTraceParityResult{}
	} else {
		for gethError, parityError := range parityErrorMappingStartingWith {
			if strings.HasPrefix(call.Error, gethError) {
				call.Error = parityError
				call.Result = callTraceParityResult{}
			}
		}
	}
}
