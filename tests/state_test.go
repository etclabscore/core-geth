// Copyright 2015 The go-ethereum Authors
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

package tests

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params/vars"
)

func TestState(t *testing.T) {
	//t.Parallel()

	st := new(testMatcher)
	// Long tests:
	// st.whitelist(`^stAttackTest/ContractCreationSpam`)
	// st.whitelist(`^stBadOpcode/badOpcodes`)
	// st.whitelist(`^stPreCompiledContracts/modexp`)
	// st.whitelist(`^stQuadraticComplexityTest/`)
	// st.whitelist(`^stStaticCall/static_Call50000`)
	// st.whitelist(`^stStaticCall/static_Return50000`)
	// st.whitelist(`^stSystemOperationsTest/CallRecursiveBomb`)
	// st.whitelist(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	// st.skipLoad(`^stTimeConsuming/`)
	// st.skipLoad(`.*vmPerformance/loop.*`)

	// Uses 1GB RAM per tested fork
	// st.whitelist(`^stStaticCall/static_Call1MB`)

	if *testEWASM == "" {
		st.skipLoad(`^stEWASM`)
	}
	if *testEVM != "" {
		// These interpreters fail Constantinople (but pass ConstantinopleFix).
		if strings.Contains(*testEVM, "aleth") || strings.Contains(*testEVM, "evmone") {
			st.skipFork("^Constantinople$")
		}

		// These tests are noted as SLOW above, and they fail against the EVMOne.so
		// So (get it?), skip 'em.
		st.skipLoad(`^stQuadraticComplexityTest/`)
		st.skipLoad(`^stStaticCall/static_Call50000`)
	}
	if *testEVM != "" || *testEWASM != "" {
		// Berlin tests are not expected to pass for external EVMs, yet.
		//
		st.skipFork("^Berlin$") // ETH
		st.skipFork("Magneto")  // ETC
		st.skipFork("London")   // ETH
		st.skipFork("Mystique") // ETC
	}
	// The multigeth data type (like the Ethereum Foundation data type) doesn't support
	// the ETC_Mystique fork/feature configuration, which omits EIP1559 and the associated BASEFEE
	// opcode stuff. This configuration cannot be represented in their struct.
	if os.Getenv(CG_CHAINCONFIG_FEATURE_EQ_MULTIGETHV0_KEY) != "" {
		st.skipFork("ETC_Mystique")
	}

	// Un-skip this when https://github.com/ethereum/tests/issues/908 is closed
	st.skipLoad(`^stQuadraticComplexityTest/QuadraticComplexitySolidity_CallDataCopy`)

	// Broken tests:
	// Expected failures:
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	// st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

	heads := make(chan *types.Header)
	//var sub ethereum.Subscription
	var err error
	quit := make(chan bool)
	if os.Getenv("AM") != "" {
		MyTransmitter = NewTransmitter()
		_, err = MyTransmitter.client.SubscribeNewHead(context.Background(), heads)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Transmitter OKGO")
		t.Log("log: Transmitter OKGO")
		//defer sub.Unsubscribe()
	}

	// For Istanbul, older tests were moved into LegacyTests
	for _, dir := range []string{
		stateTestDir,
		//legacyStateTestDir,
	} {
		if os.Getenv("AM") != "" {
			go func() {
				for {
					select {
					//case err := <-sub.Err():
					//	if err != nil {
					//		t.Fatal("subscription error", err)
					//	}
					case head := <-heads:
						bl, err := MyTransmitter.client.BlockByHash(MyTransmitter.ctx, head.Hash())
						if err != nil {
							t.Fatal(err)
						}
						fmt.Println("New block head", "num", bl.NumberU64(), "txlen", bl.Transactions().Len(), "hash", bl.Hash().Hex())
						for _, tr := range bl.Transactions() {

							MyTransmitter.mu.Lock()
							v, ok := MyTransmitter.txPendingContracts[tr.Hash()]
							MyTransmitter.mu.Unlock()
							if !ok {
								continue
							}

							fmt.Println("Matched pending transaction", tr.Hash().Hex())

							receipt, err := MyTransmitter.client.TransactionReceipt(MyTransmitter.ctx, tr.Hash())
							if err != nil {
								panic(fmt.Sprintf("receipt err=%v tx=%x", err, tr.Hash()))
							}

							gas := v.Gas()
							lim := bl.GasLimit() - (bl.GasLimit() / vars.GasLimitBoundDivisor)
							if gas >= lim {
								gas = lim - 1
							}
							next := types.NewMessage(MyTransmitter.sender, &receipt.ContractAddress, 0, v.Value(), gas, v.GasPrice(), v.Data(), false)

							sentTxHash, err := MyTransmitter.SendMessage(next)
							if err != nil {
								panic(fmt.Sprintf("send pend err=%v tx=%v", err, next))
							}

							MyTransmitter.mu.Lock()
							delete(MyTransmitter.txPendingContracts, tr.Hash())
							MyTransmitter.mu.Unlock()

							fmt.Println("Sent was-pending tx", "hash", sentTxHash.Hex())
						}
					case <-quit:
						return

					}
				}
			}()
		}
		st.walk(t, dir, func(t *testing.T, name string, test *StateTest) {
			for _, subtest := range test.Subtests(st.skipforkpat) {
				subtest := subtest
				key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)

				fmt.Println("running test", name)
				t.Run(key+"/trie", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
						_, _, err := test.Run(subtest, vmconfig, false)
						if err != nil && *testEWASM != "" {
							err = fmt.Errorf("%v ewasm=%s", err, *testEWASM)
						}
						if err != nil && len(test.json.Post[subtest.Fork][subtest.Index].ExpectException) > 0 {
							// Ignore expected errors (TODO MariusVanDerWijden check error string)
							return nil
						}
						return st.checkFailure(t, err)
					})
				})
				// t.Run(key+"/snap", func(t *testing.T) {
				// 	withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
				// 		snaps, statedb, err := test.Run(subtest, vmconfig, true)
				// 		if snaps != nil && statedb != nil {
				// 			if _, err := snaps.Journal(statedb.IntermediateRoot(false)); err != nil {
				// 				return err
				// 			}
				// 		}
				// 		if err != nil && len(test.json.Post[subtest.Fork][subtest.Index].ExpectException) > 0 {
				// 			// Ignore expected errors (TODO MariusVanDerWijden check error string)
				// 			return nil
				// 		}
				// 		if err != nil && *testEWASM != "" {
				// 			err = fmt.Errorf("%v ewasm=%s", err, *testEWASM)
				// 		}
				// 		return st.checkFailure(t, err)
				// 	})
		})
		MyTransmitter.wg.Wait()
		//t.Run("wait for pending txs", func(t *testing.T) {
		//	time.Sleep(5*time.Second)
		//	for len(MyTransmitter.txPendingContracts) > 0 {
		//		fmt.Println("Sleeping")
		//		time.Sleep(time.Minute)
		//		//MyTransmitter.mu.Lock()
		//		//for k, v := range MyTransmitter.txPendingContracts {
		//		//
		//		//}
		//		//
		//		//MyTransmitter.mu.Unlock()
		//		//
		//	}
		//	quit <- true
		//	close(quit)
		//})
	}
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	// Use config from command line arguments.
	config := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
	err := test(config)
	if err == nil {
		return
	}

	// Test failed, re-run with tracing enabled.
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	tracer := vm.NewJSONLogger(&vm.LogConfig{}, w)
	config.Debug, config.Tracer = true, tracer
	err2 := test(config)
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	w.Flush()
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	// t.Logf("EVM output: 0x%x", tracer.Output())
	// t.Logf("EVM error: %v", tracer.Error())
}
