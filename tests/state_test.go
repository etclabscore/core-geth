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
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.slow(`^stAttackTest/ContractCreationSpam`)
	st.slow(`^stBadOpcode/badOpcodes`)
	st.slow(`^stPreCompiledContracts/modexp`)
	st.slow(`^stQuadraticComplexityTest/`)
	st.slow(`^stStaticCall/static_Call50000`)
	st.slow(`^stStaticCall/static_Return50000`)
	st.slow(`^stSystemOperationsTest/CallRecursiveBomb`)
	st.slow(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	st.skipLoad(`^stTimeConsuming/`)

	// Uses 1GB RAM per tested fork
	st.skipLoad(`^stStaticCall/static_Call1MB`)

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
		st.skipFork("^Berlin$")
		st.skipFork("Magneto")
	}

	// Broken tests:
	// Expected failures:
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

	// For Istanbul, older tests were moved into LegacyTests
	for _, dir := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		st.walk(t, dir, func(t *testing.T, name string, test *StateTest) {
			for _, subtest := range test.Subtests(st.skipforkpat) {
				subtest := subtest
				key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
				name := name + "/" + key

				t.Run(key+"/trie", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
						_, _, err := test.Run(subtest, vmconfig, false)
						if err != nil && *testEWASM != "" {
							err = fmt.Errorf("%v ewasm=%s", err, *testEWASM)
						}
						return st.checkFailure(t, name+"/trie", err)
					})
				})
				t.Run(key+"/snap", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
						snaps, statedb, err := test.Run(subtest, vmconfig, true)
						if _, err := snaps.Journal(statedb.IntermediateRoot(false)); err != nil {
							return err
						}
						if err != nil && *testEWASM != "" {
							err = fmt.Errorf("%v ewasm=%s", err, *testEWASM)
						}
						return st.checkFailure(t, name+"/snap", err)
					})
				})
			}
		})
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
