package testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/tests"

	// Force-load the native, to trigger registration
	"github.com/ethereum/go-ethereum/eth/tracers"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

// callTraceParity is the result of a callTracerParity run.
type callTraceParity struct {
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
	Time                string                `json:"-"`
}

type callTraceParityAction struct {
	Author         *common.Address `json:"author,omitempty"`
	RewardType     *string         `json:"rewardType,omitempty"`
	SelfDestructed *common.Address `json:"address,omitempty"`
	Balance        *hexutil.Big    `json:"balance,omitempty"`
	CallType       string          `json:"callType,omitempty"`
	CreationMethod string          `json:"creationMethod,omitempty"`
	From           common.Address  `json:"from,omitempty"`
	Gas            hexutil.Uint64  `json:"gas,omitempty"`
	Init           *hexutil.Bytes  `json:"init,omitempty"`
	Input          *hexutil.Bytes  `json:"input,omitempty"`
	RefundAddress  *common.Address `json:"refundAddress,omitempty"`
	To             common.Address  `json:"to,omitempty"`
	Value          hexutil.Big     `json:"value,omitempty"`
}

type callTraceParityResult struct {
	Address *common.Address `json:"address,omitempty"`
	Code    *hexutil.Bytes  `json:"code,omitempty"`
	GasUsed hexutil.Uint64  `json:"gasUsed,omitempty"`
	Output  hexutil.Bytes   `json:"output,omitempty"`
}

// callTracerParityTest defines a single test to check the call tracer against.
type callTracerParityTest struct {
	Genesis *genesisT.Genesis  `json:"genesis"`
	Context *callContext       `json:"context"`
	Input   string             `json:"input"`
	Result  *[]callTraceParity `json:"result"`
}

func callTracerParityTestRunner(filename string, dirPath string) error {
	// Call tracer test found, read if from disk
	blob, err := ioutil.ReadFile(filepath.Join("..", "testdata", dirPath, filename))
	if err != nil {
		return fmt.Errorf("failed to read testcase: %v", err)
	}
	test := new(callTracerParityTest)
	if err := json.Unmarshal(blob, test); err != nil {
		return fmt.Errorf("failed to parse testcase: %v", err)
	}
	// Configure a blockchain with the given prestate
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(common.FromHex(test.Input), tx); err != nil {
		return fmt.Errorf("failed to parse testcase input: %v", err)
	}
	signer := types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)))
	origin, _ := signer.Sender(tx)
	txContext := vm.TxContext{
		Origin:   origin,
		GasPrice: tx.GasPrice(),
	}
	context := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    test.Context.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
		Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
		Difficulty:  (*big.Int)(test.Context.Difficulty),
		GasLimit:    uint64(test.Context.GasLimit),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)

	// Create the tracer, the EVM environment and run it
	tracer, err := tracers.New("callTracerParity", new(tracers.Context))
	if err != nil {
		return fmt.Errorf("failed to create call tracer: %v", err)
	}
	evm := vm.NewEVM(context, txContext, statedb, test.Genesis.Config, vm.Config{Debug: true, Tracer: tracer})

	msg, err := tx.AsMessage(signer, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare transaction for tracing: %v", err)
	}
	st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(tx.Gas()))

	if _, err = st.TransitionDb(); err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}

	// Retrieve the trace result and compare against the etalon
	res, err := tracer.GetResult()
	if err != nil {
		return fmt.Errorf("failed to retrieve trace result: %v", err)
	}
	ret := new([]callTraceParity)
	if err := json.Unmarshal(res, ret); err != nil {
		return fmt.Errorf("failed to unmarshal trace result: %v", err)
	}

	if !jsonEqualParity(ret, test.Result) {
		// uncomment this for easier debugging
		// have, _ := json.MarshalIndent(ret, "", " ")
		// want, _ := json.MarshalIndent(test.Result, "", " ")
		// return fmt.Errorf("trace mismatch: \nhave %+v\nwant %+v", string(have), string(want))
		return fmt.Errorf("trace mismatch: \nhave %+v\nwant %+v", ret, test.Result)
	}
	return nil
}

// Iterates over all the input-output datasets in the tracer parity test harness and
// runs the JavaScript tracers against them.
func TestCallTracerParity(t *testing.T) {
	folderName := "call_tracer_parity"
	files, err := ioutil.ReadDir(filepath.Join("..", "testdata", folderName))
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(file.Name(), ".json")), func(t *testing.T) {
			t.Parallel()

			err := callTracerParityTestRunner(file.Name(), folderName)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// jsonEqual is similar to reflect.DeepEqual, but does a 'bounce' via json prior to
// comparison
func jsonEqualParity(x, y interface{}) bool {
	xTrace := new([]callTraceParity)
	yTrace := new([]callTraceParity)
	if xj, err := json.Marshal(x); err == nil {
		json.Unmarshal(xj, xTrace)
	} else {
		return false
	}
	if yj, err := json.Marshal(y); err == nil {
		json.Unmarshal(yj, yTrace)
	} else {
		return false
	}
	return reflect.DeepEqual(xTrace, yTrace)
}

func BenchmarkCallTracerParity(b *testing.B) {
	files, err := filepath.Glob("testdata/call_tracer_parity/*.json")
	if err != nil {
		b.Fatalf("failed to read testdata: %v", err)
	}

	for _, file := range files {
		filename := strings.TrimPrefix(file, "testdata/call_tracer_parity/")
		b.Run(camel(strings.TrimSuffix(filename, ".json")), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				err := callTracerParityTestRunner(filename, "call_tracer_parity")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type stateDiffAccount struct {
	Balance interface{}                            `json:"balance"` // Can be either string "=" or mapping "*" => {"from": "hex", "to": "hex"}
	Code    interface{}                            `json:"code"`
	Nonce   interface{}                            `json:"nonce"`
	Storage map[common.Hash]map[string]interface{} `json:"storage"`
}

type stateDiffTest struct {
	Genesis *genesisT.Genesis                    `json:"genesis"`
	Context *callContext                         `json:"context"`
	Input   *ethapi.TransactionArgs              `json:"input"`
	Result  map[common.Address]*stateDiffAccount `json:"result"`
}

func stateDiffTracerTestRunner(filename string, dirPath string) error {
	// Call tracer test found, read if from disk
	blob, err := ioutil.ReadFile(filepath.Join("..", "testdata", dirPath, filename))
	if err != nil {
		return fmt.Errorf("failed to read testcase: %v", err)
	}
	test := new(stateDiffTest)
	if err := json.Unmarshal(blob, test); err != nil {
		return fmt.Errorf("failed to parse testcase: %v", err)
	}

	// Configure a blockchain with the given prestate
	msg, err := test.Input.ToMessage(uint64(test.Context.GasLimit), nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// This is needed for trace_call (debug mode),
	// as the Transaction is being run on top of the block transactions,
	// which might lead into ErrInsufficientFundsForTransfer error
	canTransfer := func(db vm.StateDB, sender common.Address, amount *big.Int) bool {
		if msg.From() == sender {
			return true
		}
		return core.CanTransfer(db, sender, amount)
	}

	// If the actual transaction would fail, then there is no reason to actually transfer any balance at all
	transfer := func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
		toAmount := new(big.Int).Set(amount)
		senderBalance := db.GetBalance(sender)
		if senderBalance.Cmp(toAmount) < 0 {
			toAmount.Set(big.NewInt(0))
		}
		core.Transfer(db, sender, recipient, toAmount)
	}

	txContext := vm.TxContext{
		Origin:   msg.From(),
		GasPrice: msg.GasPrice(),
	}
	context := vm.BlockContext{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		Coinbase:    test.Context.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
		Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
		Difficulty:  (*big.Int)(test.Context.Difficulty),
		GasLimit:    uint64(test.Context.GasLimit),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)

	// Store the truth on whether from account has enough balance for context usage
	gasCost := new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.GasPrice())
	totalCost := new(big.Int).Add(gasCost, msg.Value())

	// It is important to use core.CanTransfer on the two following lines
	hasFromSufficientBalanceForValueAndGasCost := core.CanTransfer(statedb, msg.From(), totalCost)
	hasFromSufficientBalanceForGasCost := core.CanTransfer(statedb, msg.From(), gasCost)

	// Add extra context needed for state_diff
	taskExtraContext := map[string]interface{}{
		"hasFromSufficientBalanceForValueAndGasCost": hasFromSufficientBalanceForValueAndGasCost,
		"hasFromSufficientBalanceForGasCost":         hasFromSufficientBalanceForGasCost,
		"from":                                       msg.From(),
		"coinbase":                                   context.Coinbase,
		"gasLimit":                                   msg.Gas(),
		"gasPrice":                                   msg.GasPrice(),
	}

	if msg.To() != nil {
		taskExtraContext["msgTo"] = *msg.To()
	}

	// Create the tracer, the EVM environment and run it
	tracer, err := tracers.New("stateDiffTracer", new(tracers.Context))
	if err != nil {
		return fmt.Errorf("failed to create state diff tracer: %v", err)
	}
	evm := vm.NewEVM(context, txContext, statedb, test.Genesis.Config, vm.Config{Debug: true, Tracer: tracer})

	tracer.CapturePreEVM(evm, taskExtraContext)

	st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(msg.Gas()))
	if _, err = st.TransitionDb(); err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}

	// Retrieve the trace result and compare against the etalon
	res, err := tracer.GetResult()
	if err != nil {
		return fmt.Errorf("failed to retrieve trace result: %v", err)
	}
	ret := new(map[common.Address]*stateDiffAccount)
	if err := json.Unmarshal(res, ret); err != nil {
		return fmt.Errorf("failed to unmarshal trace result: %v", err)
	}

	if !jsonEqualStateDiff(ret, test.Result) {
		// uncomment this for easier debugging
		// have, _ := json.MarshalIndent(ret, "", " ")
		// want, _ := json.MarshalIndent(test.Result, "", " ")
		// return fmt.Errorf("trace mismatch: \nhave %+v\nwant %+v", string(have), string(want))
		return fmt.Errorf("trace mismatch: \nhave %+v\nwant %+v", ret, test.Result)
	}
	return nil
}

// TestStateDiffTracer Iterates over all the input-output datasets in the state diff tracer test harness and
// runs the JavaScript tracers against them.
func TestStateDiffTracer(t *testing.T) {
	folderName := "state_diff"
	files, err := ioutil.ReadDir(filepath.Join("..", "testdata", folderName))
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(file.Name(), ".json")), func(t *testing.T) {
			t.Parallel()

			err := stateDiffTracerTestRunner(file.Name(), folderName)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// jsonEqual is similar to reflect.DeepEqual, but does a 'bounce' via json prior to
// comparison
func jsonEqualStateDiff(x, y interface{}) bool {
	xTrace := new(map[common.Address]*stateDiffAccount)
	yTrace := new(map[common.Address]*stateDiffAccount)
	if xj, err := json.Marshal(x); err == nil {
		json.Unmarshal(xj, xTrace)
	} else {
		return false
	}
	if yj, err := json.Marshal(y); err == nil {
		json.Unmarshal(yj, yTrace)
	} else {
		return false
	}
	return reflect.DeepEqual(xTrace, yTrace)
}
