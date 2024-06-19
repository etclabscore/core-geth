package tracetest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	// Force-load the native, to trigger registration
)

type stateDiffAccount struct {
	Balance interface{}                            `json:"balance"` // Can be either string "=" or mapping "*" => {"from": "hex", "to": "hex"}
	Code    interface{}                            `json:"code"`
	Nonce   interface{}                            `json:"nonce"`
	Storage map[common.Hash]map[string]interface{} `json:"storage"`
}

type stateDiffTest struct {
	Genesis        *genesisT.Genesis       `json:"genesis"`
	Context        *callContext            `json:"context"`
	Input          *ethapi.TransactionArgs `json:"input"`
	StateOverrides *ethapi.StateOverride
	TracerConfig   json.RawMessage                       `json:"tracerConfig"`
	Result         *map[common.Address]*stateDiffAccount `json:"result"`
}

func stateDiffTracerTestRunner(tracerName string, filename string, dirPath string, t testing.TB) error {
	// // Call tracer test found, read if from disk
	// blob, err := os.ReadFile(filepath.Join("testdata", dirPath, filename))
	// if err != nil {
	// 	return fmt.Errorf("failed to read testcase: %v", err)
	// }
	// test := new(stateDiffTest)
	// if err := json.Unmarshal(blob, test); err != nil {
	// 	return fmt.Errorf("failed to parse testcase: %v", err)
	// }

	// // Configure a blockchain with the given prestate
	// msg, err := test.Input.ToMessage(uint64(test.Context.GasLimit), nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to create transaction: %v", err)
	// }

	// // This is needed for trace_call (debug mode),
	// // as the Transaction is being run on top of the block transactions,
	// // which might lead into ErrInsufficientFundsForTransfer error

	// txContext := vm.TxContext{
	// 	Origin:   msg.From,
	// 	GasPrice: msg.GasPrice,
	// }
	// context := vm.BlockContext{
	// 	CanTransfer: core.CanTransfer,
	// 	Transfer:    core.Transfer,
	// 	Coinbase:    test.Context.Miner,
	// 	BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
	// 	Time:        uint64(test.Context.Time),
	// 	Difficulty:  (*big.Int)(test.Context.Difficulty),
	// 	GasLimit:    uint64(test.Context.GasLimit),
	// }
	// state := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false, rawdb.HashScheme)

	// if err := test.StateOverrides.Apply(state.StateDB); err != nil {
	// 	return fmt.Errorf("failed to apply test stateOverrides: %v", err)
	// }

	// // Create the tracer, the EVM environment and run it
	// tracer, err := tracers.DefaultDirectory.New(tracerName, new(tracers.Context), test.TracerConfig)
	// if err != nil {
	// 	return fmt.Errorf("failed to create state diff tracer: %v", err)
	// }
	// evm := vm.NewEVM(context, txContext, state.StateDB, test.Genesis.Config, vm.Config{Tracer: tracer})

	// if traceStateCapturer, ok := tracer.(vm.EVMLogger_StateCapturer); ok {
	// 	traceStateCapturer.CapturePreEVM(evm)
	// }

	// st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(msg.GasLimit))
	// if _, err = st.TransitionDb(); err != nil {
	// 	return fmt.Errorf("failed to execute transaction: %v", err)
	// }

	// // Retrieve the trace result and compare against the etalon
	// res, err := tracer.GetResult()
	// if err != nil {
	// 	return fmt.Errorf("failed to retrieve trace result: %v", err)
	// }
	// ret := new(map[common.Address]*stateDiffAccount)
	// if err := json.Unmarshal(res, ret); err != nil {
	// 	return fmt.Errorf("failed to unmarshal trace result: %v", err)
	// }

	// if !jsonEqualStateDiff(ret, test.Result) {
	// 	t.Logf("tracer name: %s", tracerName)

	// 	// uncomment this for easier debugging
	// 	have, _ := json.MarshalIndent(ret, "", " ")
	// 	want, _ := json.MarshalIndent(test.Result, "", " ")
	// 	t.Logf("trace mismatch: \nhave %+v\nwant %+v", string(have), string(want))

	// 	// uncomment this for harder debugging <3 meowsbits
	// 	lines := deep.Equal(ret, test.Result)
	// 	for _, l := range lines {
	// 		t.Logf("%s", l)
	// 	}

	// 	t.Fatalf("trace mismatch: \nhave %+v\nwant %+v", ret, test.Result)
	// }
	return nil
}

// Iterates over all the input-output datasets in the tracer test harness and
// runs the Native tracer against them.
func TestStateDiffTracerNative(t *testing.T) {
	testStateDiffTracer("stateDiffTracer", "state_diff", t)
}

// TestStateDiffTracer Iterates over all the input-output datasets in the state diff tracer test harness and
// runs the JavaScript tracers against them.
func testStateDiffTracer(tracerName string, dirPath string, t *testing.T) {
	files, err := os.ReadDir(filepath.Join("testdata", dirPath))
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		// if !strings.HasSuffix(file.Name(), "errInsufficientFundsForTransfer_and_gas_cost.json") {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(file.Name(), ".json")), func(t *testing.T) {
			t.Parallel()

			err := stateDiffTracerTestRunner(tracerName, file.Name(), dirPath, t)
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
