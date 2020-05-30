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
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params/vars"
)

func TestState(t *testing.T) {
	//t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.whitelist(`^stAttackTest/ContractCreationSpam`)
	st.whitelist(`^stBadOpcode/badOpcodes`)
	st.whitelist(`^stPreCompiledContracts/modexp`)
	st.whitelist(`^stQuadraticComplexityTest/`)
	st.whitelist(`^stStaticCall/static_Call50000`)
	st.whitelist(`^stStaticCall/static_Return50000`)
	st.whitelist(`^stSystemOperationsTest/CallRecursiveBomb`)
	st.whitelist(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	st.whitelist(`^stTimeConsuming/`)

	// Uses 1GB RAM per tested fork
	st.whitelist(`^stStaticCall/static_Call1MB`)

	// Broken tests:
	// Expected failures:
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

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
			for _, subtest := range test.Subtests() {
				subtest := subtest
				key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
				name := name + "/" + key

				fmt.Println("running test", name)
				t.Run(key+"/trie", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
						_, _, err := test.Run(subtest, vmconfig, false)
						return st.checkFailure(t, name+"/trie", err)
					})
				})
				//time.Sleep(time.Second)
				//t.Run(key+"/snap", func(t *testing.T) {
				//	withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
				//		snaps, statedb, err := test.Run(subtest, vmconfig, true)
				//		if _, err := snaps.Journal(statedb.IntermediateRoot(false)); err != nil {
				//			return err
				//		}
				//		return st.checkFailure(t, name+"/snap", err)
				//	})
				//})
			}
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
	tracer := vm.NewJSONLogger(&vm.LogConfig{DisableMemory: true}, w)
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
	//t.Logf("EVM output: 0x%x", tracer.Output())
	//t.Logf("EVM error: %v", tracer.Error())
}
