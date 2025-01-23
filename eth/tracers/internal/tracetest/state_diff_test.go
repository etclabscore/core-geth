package tracetest

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/tests"
	// Force-load the native, to trigger registration
)

type stateDiffAccount struct {
	Balance interface{}                            `json:"balance"` // Can be either string "=" or mapping "*" => {"from": "hex", "to": "hex"}
	Code    interface{}                            `json:"code"`
	Nonce   interface{}                            `json:"nonce"`
	Storage map[common.Hash]map[string]interface{} `json:"storage"`
}

type stateDiffTest struct {
	Genesis        *genesisT.Genesis `json:"genesis"`
	Context        *callContext      `json:"context"`
	Input          string            `json:"input"`
	StateOverrides *ethapi.StateOverride
	TracerConfig   json.RawMessage                       `json:"tracerConfig"`
	Result         *map[common.Address]*stateDiffAccount `json:"result"`
}

func stateDiffTracerTestRunner(tracerName string, filename string, dirPath string, t testing.TB) error {
	// Call tracer test found, read if from disk
	blob, err := os.ReadFile(filepath.Join("testdata", dirPath, filename))
	if err != nil {
		return fmt.Errorf("failed to read testcase: %v", err)
	}
	test := new(stateDiffTest)
	if err := json.Unmarshal(blob, test); err != nil {
		return fmt.Errorf("failed to parse testcase: %v", err)
	}

	// Configure a blockchain with the given prestate
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

	// Configure a blockchain with the given prestate
	// tx := test.Input.ToTransaction()
	// fmt.Printf("tx: %+v\n", tx)

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(common.FromHex(test.Input)); err != nil {
		t.Fatalf("failed to parse testcase input: %v", err)
	}
	// if err := rlp.DecodeBytes(common.FromHex(test.Input), tx); err != nil {
	// 	return fmt.Errorf("failed to parse testcase input: %v", err)
	// }
	signer := types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)), uint64(test.Context.Time))
	context := test.Context.toBlockContext(test.Genesis)
	state := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false, rawdb.HashScheme)
	defer state.Close()

	// Create the tracer, the EVM environment and run it
	tracer, err := tracers.DefaultDirectory.New(tracerName, new(tracers.Context), test.TracerConfig)
	if err != nil {
		return fmt.Errorf("failed to create call tracer: %v", err)
	}

	state.StateDB.SetLogger(tracer.Hooks)
	msg, err := core.TransactionToMessage(tx, signer, context.BaseFee)
	if err != nil {
		return fmt.Errorf("failed to prepare transaction for tracing: %v", err)
	}
	evm := vm.NewEVM(context, core.NewEVMTxContext(msg), state.StateDB, test.Genesis.Config, vm.Config{Tracer: tracer.Hooks})
	tracer.OnTxStart(evm.GetVMContext(), tx, msg.From)
	vmRet, err := core.ApplyMessage(evm, msg, new(core.GasPool).AddGas(tx.Gas()))
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}
	tracer.OnTxEnd(&types.Receipt{GasUsed: vmRet.UsedGas}, nil)

	// Retrieve the trace result and compare against the etalon
	res, err := tracer.GetResult()
	// fmt.Println("res:", string(res))
	if err != nil {
		return fmt.Errorf("failed to retrieve trace result: %v", err)
	}
	// var ret prestateTrace
	ret := new(map[common.Address]*stateDiffAccount)
	if err := json.Unmarshal(res, ret); err != nil {
		return fmt.Errorf("failed to unmarshal trace result: %v", err)
	}
	have, _ := json.MarshalIndent(ret, "", " ")
	fmt.Printf("have: %+v", string(have))

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
	t.FailNow()
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
		if !strings.HasSuffix(file.Name(), "test.json") {
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
