// Copyright 2017 The go-ethereum Authors
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

package tracers

import (
	"crypto/ecdsa"
	"crypto/rand"
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
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/tests"
)

// To generate a new callTracer test, copy paste the makeTest method below into
// a Geth console and call it with a transaction hash you which to export.

/*
// makeTest generates a callTracer test by running a prestate reassembled and a
// call trace run, assembling all the gathered information into a test case.
var makeTest = function(tx, reexec) {
  // Generate the genesis block from the block, transaction and prestate data
  var block   = eth.getBlock(eth.getTransaction(tx).blockHash);
  var genesis = eth.getBlock(block.parentHash);

  delete genesis.gasUsed;
  delete genesis.logsBloom;
  delete genesis.parentHash;
  delete genesis.receiptsRoot;
  delete genesis.sha3Uncles;
  delete genesis.size;
  delete genesis.transactions;
  delete genesis.transactionsRoot;
  delete genesis.uncles;

  genesis.gasLimit  = genesis.gasLimit.toString();
  genesis.number    = genesis.number.toString();
  genesis.timestamp = genesis.timestamp.toString();

  genesis.alloc = debug.traceTransaction(tx, {tracer: "prestateTracer", reexec: reexec});
  for (var key in genesis.alloc) {
    genesis.alloc[key].nonce = genesis.alloc[key].nonce.toString();
  }
  genesis.config = admin.nodeInfo.protocols.eth.config;

  // Generate the call trace and produce the test input
  var result = debug.traceTransaction(tx, {tracer: "callTracer", reexec: reexec});
  delete result.time;

  console.log(JSON.stringify({
    genesis: genesis,
    context: {
      number:     block.number.toString(),
      difficulty: block.difficulty,
      timestamp:  block.timestamp.toString(),
      gasLimit:   block.gasLimit.toString(),
      miner:      block.miner,
    },
    input:  eth.getRawTransaction(tx),
    result: result,
  }, null, 2));
}
*/

// callTrace is the result of a callTracer run.
type callTrace struct {
	Type    string          `json:"type"`
	From    common.Address  `json:"from"`
	To      common.Address  `json:"to"`
	Input   hexutil.Bytes   `json:"input"`
	Output  hexutil.Bytes   `json:"output"`
	Gas     *hexutil.Uint64 `json:"gas,omitempty"`
	GasUsed *hexutil.Uint64 `json:"gasUsed,omitempty"`
	Value   *hexutil.Big    `json:"value,omitempty"`
	Error   string          `json:"error,omitempty"`
	Calls   []callTrace     `json:"calls,omitempty"`
}

type callContext struct {
	Number     math.HexOrDecimal64   `json:"number"`
	Difficulty *math.HexOrDecimal256 `json:"difficulty"`
	Time       math.HexOrDecimal64   `json:"timestamp"`
	GasLimit   math.HexOrDecimal64   `json:"gasLimit"`
	Miner      common.Address        `json:"miner"`
}

// callTracerTest defines a single test to check the call tracer against.
type callTracerTest struct {
	Genesis *genesisT.Genesis `json:"genesis"`
	Context *callContext      `json:"context"`
	Input   string            `json:"input"`
	Result  *callTrace        `json:"result"`
}

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

func TestPrestateTracerCreate2(t *testing.T) {
	unsignedTx := types.NewTransaction(1, common.HexToAddress("0x00000000000000000000000000000000deadbeef"),
		new(big.Int), 5000000, big.NewInt(1), []byte{})

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	signer := types.NewEIP155Signer(big.NewInt(1))
	tx, err := types.SignTx(unsignedTx, signer, privateKeyECDSA)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	/**
		This comes from one of the test-vectors on the Skinny Create2 - EIP

	    address 0x00000000000000000000000000000000deadbeef
	    salt 0x00000000000000000000000000000000000000000000000000000000cafebabe
	    init_code 0xdeadbeef
	    gas (assuming no mem expansion): 32006
	    result: 0x60f3f640a8508fC6a86d45DF051962668E1e8AC7
	*/
	origin, _ := signer.Sender(tx)
	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Origin:      origin,
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(8000000),
		Time:        new(big.Int).SetUint64(5),
		Difficulty:  big.NewInt(0x30000),
		GasLimit:    uint64(6000000),
		GasPrice:    big.NewInt(1),
	}
	alloc := genesisT.GenesisAlloc{}

	// The code pushes 'deadbeef' into memory, then the other params, and calls CREATE2, then returns
	// the address
	alloc[common.HexToAddress("0x00000000000000000000000000000000deadbeef")] = genesisT.GenesisAccount{
		Nonce:   1,
		Code:    hexutil.MustDecode("0x63deadbeef60005263cafebabe6004601c6000F560005260206000F3"),
		Balance: big.NewInt(1),
	}
	alloc[origin] = genesisT.GenesisAccount{
		Nonce:   1,
		Code:    []byte{},
		Balance: big.NewInt(500000000000000),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), alloc, false)

	// Create the tracer, the EVM environment and run it
	tracer, err := New("prestateTracer")
	if err != nil {
		t.Fatalf("failed to create call tracer: %v", err)
	}
	evm := vm.NewEVM(context, statedb, params.MainnetChainConfig, vm.Config{Debug: true, Tracer: tracer})

	msg, err := tx.AsMessage(signer)
	if err != nil {
		t.Fatalf("failed to prepare transaction for tracing: %v", err)
	}
	st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(tx.Gas()))
	if _, err = st.TransitionDb(); err != nil {
		t.Fatalf("failed to execute transaction: %v", err)
	}
	// Retrieve the trace result and compare against the etalon
	res, err := tracer.GetResult()
	if err != nil {
		t.Fatalf("failed to retrieve trace result: %v", err)
	}
	ret := make(map[string]interface{})
	if err := json.Unmarshal(res, &ret); err != nil {
		t.Fatalf("failed to unmarshal trace result: %v", err)
	}
	if _, has := ret["0x60f3f640a8508fc6a86d45df051962668e1e8ac7"]; !has {
		t.Fatalf("Expected 0x60f3f640a8508fc6a86d45df051962668e1e8ac7 in result")
	}
}

// Iterates over all the input-output datasets in the tracer test harness and
// runs the JavaScript tracers against them.
func TestCallTracer(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "call_tracer_") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(strings.TrimPrefix(file.Name(), "call_tracer_"), ".json")), func(t *testing.T) {
			t.Parallel()

			// Call tracer test found, read if from disk
			blob, err := ioutil.ReadFile(filepath.Join("testdata", file.Name()))
			if err != nil {
				t.Fatalf("failed to read testcase: %v", err)
			}
			test := new(callTracerTest)
			if err := json.Unmarshal(blob, test); err != nil {
				t.Fatalf("failed to parse testcase: %v", err)
			}
			// Configure a blockchain with the given prestate
			tx := new(types.Transaction)
			if err := rlp.DecodeBytes(common.FromHex(test.Input), tx); err != nil {
				t.Fatalf("failed to parse testcase input: %v", err)
			}
			signer := types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)))
			origin, _ := signer.Sender(tx)

			context := vm.Context{
				CanTransfer: core.CanTransfer,
				Transfer:    core.Transfer,
				Origin:      origin,
				Coinbase:    test.Context.Miner,
				BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
				Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
				Difficulty:  (*big.Int)(test.Context.Difficulty),
				GasLimit:    uint64(test.Context.GasLimit),
				GasPrice:    tx.GasPrice(),
			}
			_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)

			// Create the tracer, the EVM environment and run it
			tracer, err := New("callTracer")
			if err != nil {
				t.Fatalf("failed to create call tracer: %v", err)
			}
			evm := vm.NewEVM(context, statedb, test.Genesis.Config, vm.Config{Debug: true, Tracer: tracer})

			msg, err := tx.AsMessage(signer)
			if err != nil {
				t.Fatalf("failed to prepare transaction for tracing: %v", err)
			}
			st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(tx.Gas()))
			if _, err = st.TransitionDb(); err != nil {
				t.Fatalf("failed to execute transaction: %v", err)
			}
			// Retrieve the trace result and compare against the etalon
			res, err := tracer.GetResult()
			if err != nil {
				t.Fatalf("failed to retrieve trace result: %v", err)
			}
			ret := new(callTrace)
			if err := json.Unmarshal(res, ret); err != nil {
				t.Fatalf("failed to unmarshal trace result: %v", err)
			}

			if !jsonEqual(ret, test.Result) {
				// uncomment this for easier debugging
				//have, _ := json.MarshalIndent(ret, "", " ")
				//want, _ := json.MarshalIndent(test.Result, "", " ")
				//t.Fatalf("trace mismatch: \nhave %+v\nwant %+v", string(have), string(want))
				t.Fatalf("trace mismatch: \nhave %+v\nwant %+v", ret, test.Result)
			}
		})
	}
}

// jsonEqual is similar to reflect.DeepEqual, but does a 'bounce' via json prior to
// comparison
func jsonEqual(x, y interface{}) bool {
	xTrace := new(callTrace)
	yTrace := new(callTrace)
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

func callTracerParityTestRunner(filename string) error {
	// Call tracer test found, read if from disk
	blob, err := ioutil.ReadFile(filepath.Join("testdata", filename))
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

	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Origin:      origin,
		Coinbase:    test.Context.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
		Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
		Difficulty:  (*big.Int)(test.Context.Difficulty),
		GasLimit:    uint64(test.Context.GasLimit),
		GasPrice:    tx.GasPrice(),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)

	// Create the tracer, the EVM environment and run it
	tracer, err := New("callTracerParity")
	if err != nil {
		return fmt.Errorf("failed to create call tracer: %v", err)
	}
	evm := vm.NewEVM(context, statedb, test.Genesis.Config, vm.Config{Debug: true, Tracer: tracer})

	msg, err := tx.AsMessage(signer)
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
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "parity_call_tracer_") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(strings.TrimPrefix(file.Name(), "parity_call_tracer_"), ".json")), func(t *testing.T) {
			t.Parallel()

			err := callTracerParityTestRunner(file.Name())
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
	files, err := filepath.Glob("testdata/parity_call_tracer_*.json")
	if err != nil {
		b.Fatalf("failed to read testdata: %v", err)
	}

	for _, file := range files {
		filename := strings.TrimPrefix(file, "testdata/")
		b.Run(camel(strings.TrimSuffix(strings.TrimPrefix(filename, "parity_call_tracer_"), ".json")), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				err := callTracerParityTestRunner(filename)
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

type stateDiffBalance struct {
	From *hexutil.Big `json:"from"`
	To   *hexutil.Big `json:"to"`
}

type stateDiffCode struct {
	From hexutil.Bytes `json:"from"`
	To   hexutil.Bytes `json:"to"`
}

type stateDiffNonce struct {
	From hexutil.Uint64 `json:"from"`
	To   hexutil.Uint64 `json:"to"`
}

type stateDiffStorage struct {
	From common.Hash `json:"from"`
	To   common.Hash `json:"to"`
}

type stateDiffTest struct {
	Genesis *genesisT.Genesis                    `json:"genesis"`
	Context *callContext                         `json:"context"`
	Input   ethapi.CallArgs                      `json:"input"`
	Result  map[common.Address]*stateDiffAccount `json:"result"`
}

func stateDiffTracerTestRunner(filename string) error {
	// Call tracer test found, read if from disk
	blob, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		return fmt.Errorf("failed to read testcase: %v", err)
	}
	test := new(stateDiffTest)
	if err := json.Unmarshal(blob, test); err != nil {
		return fmt.Errorf("failed to parse testcase: %v", err)
	}

	// Configure a blockchain with the given prestate
	msg := test.Input.ToMessage(uint64(test.Context.GasLimit))

	// This is needed for trace_call (debug mode),
	// as the Transaction is being run on top of the block transactions,
	// which might lead into ErrInsufficientFundsForTransfer error
	canTransfer := func(db vm.StateDB, sender common.Address, amount *big.Int) bool {
		if msg.From() == sender {
			return true
		}
		return core.CanTransfer(db, sender, amount)
	}

	// If the actual transaction would fail, then their is no reason to actually transfer any balance at all
	transfer := func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
		toAmount := new(big.Int).Set(amount)
		senderBalance := db.GetBalance(sender)
		if senderBalance.Cmp(toAmount) < 0 {
			toAmount.Set(big.NewInt(0))
		}
		core.Transfer(db, sender, recipient, toAmount)
	}

	context := vm.Context{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		Origin:      msg.From(),
		Coinbase:    test.Context.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
		Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
		Difficulty:  (*big.Int)(test.Context.Difficulty),
		GasLimit:    uint64(test.Context.GasLimit),
		GasPrice:    msg.GasPrice(),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)

	// Store the truth on wether from acount has enough balance for context usage
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
	tracer, err := New("stateDiffTracer")
	if err != nil {
		return fmt.Errorf("failed to create state diff tracer: %v", err)
	}
	evm := vm.NewEVM(context, statedb, test.Genesis.Config, vm.Config{Debug: true, Tracer: tracer})

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
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to retrieve tracer test suite: %v", err)
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "state_diff_tracer_") {
			continue
		}
		file := file // capture range variable
		t.Run(camel(strings.TrimSuffix(strings.TrimPrefix(file.Name(), "state_diff_tracer_"), ".json")), func(t *testing.T) {
			t.Parallel()

			err := stateDiffTracerTestRunner(file.Name())
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
