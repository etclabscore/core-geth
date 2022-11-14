package evmc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/evmc/v10/bindings/go/evmc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"
	"github.com/nsf/jsondiff"
)

type evmcVM struct {
	cap       evmc.Capability
	path      string
	canonical bool // should only be true for native evm
	evmLogger vm.EVMLogger
}

// subtestResult is the result of running a state subtest.
type subtestResult struct {
	fork string
	env  interface{}
	pre  interface{}
	tx   interface{}

	err      error
	dump     state.Dump
	dumpJSON []byte

	debugOutput []byte
}

var (
	// Configure the tested external EVMs
	soEVMOne = evmcVM{
		cap:       evmc.CapabilityEVM1,
		path:      "../build/_workspace/evmone/lib/libevmone.so",
		evmLogger: &logger.JSONLogger{},
	}
	soHera = evmcVM{
		cap:       evmc.CapabilityEWASM,
		path:      "../build/_workspace/hera/build/src/libhera.so",
		evmLogger: &logger.JSONLogger{},
	}
	native = evmcVM{
		canonical: true,
		// evmLogger: &logger.StructLogger{}, // <- The StructLogger gives more readable opcode step output.
		evmLogger: &logger.JSONLogger{},
	}
)

var (
	// Define fork(s) to run
	forks = []string{"Istanbul", "Homestead"}
)

var (
	// Test file targets
	myTestFiles = []string{
		// "../tests/testdata/LegacyTests/Constantinople/GeneralStateTests/VMTests/vmArithmeticTest/add.json",
		// "../tests/testdata/GeneralStateTests/VMTests/vmArithmeticTest/add.json",

		"../tests/testdata/LegacyTests/Constantinople/GeneralStateTests/stRandom/randomStatetest153.json",
		// "../tests/testdata/GeneralStateTests/stRandom/randomStatetest153.json",

		// "../tests/testdata/GeneralStateTests/stBadOpcode/operationDiffGas.json",
		// "../tests/testdata/GeneralStateTests/stCallCodes/callcallcallcode_001_SuicideMiddle.json",
		// "../tests/testdata/GeneralStateTests/stBadOpcode/opcEADiffPlaces.json",
		// "../tests/testdata/GeneralStateTests/stZeroKnowledge2/ecadd_0-0_0-0_25000_128.json",
	}
)

func TestStateEVMC(t *testing.T) {

	// Configure the go-ethereum logger
	verbose := false
	if verbose {
		glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
		glogger.Verbosity(log.Lvl(log.LvlDebug))
		log.Root().SetHandler(glogger)
	}

	for _, myTestFile := range myTestFiles {
		runEVMCStateTestFile(t, []evmcVM{soEVMOne, soHera}, myTestFile)
	}
}

func runEVMCStateTestFile(t *testing.T, myEVMs []evmcVM, testFile string) {
	// Load the test content from the input file
	src, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	var mytests map[string]tests.StateTest
	if err = json.Unmarshal(src, &mytests); err != nil {
		t.Fatal(err)
	}

	// Run the tests
	t.Log("====== Running tests from", testFile)
	for _, mytest := range mytests {

		nativeResults := runEVMStateTest(native, mytest, forks)
		evmoneResults := runEVMStateTest(soEVMOne, mytest, forks)
		heraResults := runEVMStateTest(soHera, mytest, forks)

		// Compare output
		for i, nativeResult := range nativeResults {

			if nativeResult.err != nil {
				t.Fatal(err)
			}

			evmoneResult := evmoneResults[i]
			heraResult := heraResults[i]

			compareResults(t, "EVMOne", nativeResult, evmoneResult)
			compareResults(t, "Hera", nativeResult, heraResult)
		}
	}
}

func runEVMStateTest(myvm evmcVM, test tests.StateTest, forks []string) []subtestResult {
	loggerConfig := &logger.Config{
		Debug:            true,
		EnableMemory:     true,
		DisableStack:     false,
		DisableStorage:   false,
		EnableReturnData: true,
	}
	var (
		tracer vm.EVMLogger
	)

	results := []subtestResult{}
	for _, st := range test.Subtests(nil) {

		res := subtestResult{
			fork: st.Fork,
			env:  test.GetEnv(),
			tx:   test.GetTx(),
			pre:  test.GetPre(),
		}

		if len(forks) > 0 {
			forkMatched := false
			for _, fork := range forks {
				if strings.Contains(fork, st.Fork) {
					forkMatched = true
					break
				}
			}
			if !forkMatched {
				continue
			}
		}

		// JSONLogger vars. Unused if logger type is StructLogger.
		var buf *bytes.Buffer
		var w *bufio.Writer

		if _, ok := myvm.evmLogger.(*logger.StructLogger); ok {
			tracer = logger.NewStructLogger(loggerConfig)
		} else if _, ok := myvm.evmLogger.(*logger.JSONLogger); ok {
			buf = new(bytes.Buffer)
			w = bufio.NewWriter(buf)
			tracer = logger.NewJSONLogger(loggerConfig, w)
		}

		// Iterate over all the tests, run them and aggregate the results
		cfg := vm.Config{
			Tracer:                  tracer,
			Debug:                   true,
			EnablePreimageRecording: true,
		}

		if myvm.path == "" {

		} else if myvm.cap == evmc.CapabilityEVM1 {
			cfg.EVMInterpreter = myvm.path
			vm.InitEVMCEVM(cfg.EVMInterpreter)
		} else if myvm.cap == evmc.CapabilityEWASM {
			cfg.EWASMInterpreter = myvm.path
			vm.InitEVMCEwasm(cfg.EWASMInterpreter)
		}

		_, s, err := test.Run(st, cfg, false)
		res.err = err

		if _, ok := myvm.evmLogger.(*logger.StructLogger); ok {

			res.debugOutput, _ = tracer.(*logger.StructLogger).GetResult()

			marshalIndented := false
			if marshalIndented {
				execRes := logger.ExecutionResult{}
				if e := json.Unmarshal(res.debugOutput, &execRes); e != nil {
					panic(e)
				}
				// Indent it.
				res.debugOutput, _ = json.MarshalIndent(execRes, "", "  ")
			}

		} else if _, ok = myvm.evmLogger.(*logger.JSONLogger); ok {
			w.Flush()
			res.debugOutput = buf.Bytes()
		}

		if s != nil {
			res.dump = s.RawDump(nil)
			res.dumpJSON, _ = json.MarshalIndent(res.dump, "", "  ")
		}

		results = append(results, res)
	}
	return results
}

func compareResults(t *testing.T, name string, referenceResult subtestResult, result subtestResult) {
	if referenceResult.err != result.err {
		t.Errorf("[%s] Error mismatch: %v != %v", name, referenceResult.err, result.err)
	}

	opts := jsondiff.DefaultConsoleOptions()
	diff, diffStr := jsondiff.Compare(referenceResult.dumpJSON, result.dumpJSON, &opts)

	if diff != jsondiff.FullMatch {
		t.Logf("[%s] ====== Test failed, fork: %s", name, result.fork)
		t.Log("JSON diff:", diff)
		t.Log(diffStr)
		t.Log("Reference debug output:")
		t.Log(string(referenceResult.debugOutput))
		t.Log("Target debug output:")
		t.Log(string(result.debugOutput))

		b, _ := json.MarshalIndent(result.env, "", "  ")
		t.Log("Test.env", string(b))

		b, _ = json.MarshalIndent(result.pre, "", "  ")
		t.Log("Test.pre", string(b))

		b, _ = json.MarshalIndent(result.tx, "", "  ")
		t.Log("Test.tx", string(b))
	} else {
		t.Logf("[%s] ====== Test passed, fork: %s", name, result.fork)
	}
}
