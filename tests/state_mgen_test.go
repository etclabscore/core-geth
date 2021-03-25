// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp/tconvert"
	"github.com/iancoleman/strcase"
)

func TestGenState(t *testing.T) {
	if os.Getenv(CG_GENERATE_STATE_TESTS_KEY) == "" {
		t.Skip()
	}
	if os.Getenv(CG_CHAINCONFIG_CHAINSPECS_OPENETHEREUM_KEY) == "" {
		t.Fatal("Must use chainspec files for fork configurations.")
	}

	st := new(testMatcher)

	// Generating tests should NOT skip slow or time consuming tests.

	// Long tests:
	//st.slow(`^stAttackTest/ContractCreationSpam`)
	//st.slow(`^stBadOpcode/badOpcodes`)
	//st.slow(`^stPreCompiledContracts/modexp`)
	//st.slow(`^stQuadraticComplexityTest/`)
	//st.slow(`^stStaticCall/static_Call50000`)
	//st.slow(`^stStaticCall/static_Return50000`)
	//st.slow(`^stStaticCall/static_Call1MB`)
	//st.slow(`^stSystemOperationsTest/CallRecursiveBomb`)
	//st.slow(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	//st.skipLoad(`^stTimeConsuming/`)

	// Broken tests:
	// Expected failures:
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

	st.walkFullName(t, stateTestDir, withWritingTests)

	// For Istanbul, older tests were moved into LegacyTests
	st.walkFullName(t, legacyStateTestDir, withWritingTests)
}

var (
	// core-geth
	baseDirCG            = filepath.Join(".", "testdata.core-geth")
	blockTestDirCG       = filepath.Join(baseDir, "BlockchainTests")
	stateTestDirCG       = filepath.Join(baseDir, "GeneralStateTests")
	legacyStateTestDirCG = filepath.Join(baseDir, "LegacyTests", "Constantinople", "GeneralStateTests")

	// foundation v1.10.1
	baseDir101            = filepath.Join(".", "testdata.ethereum")
	blockTestDir101       = filepath.Join(baseDir, "BlockchainTests")
	stateTestDir101       = filepath.Join(baseDir, "GeneralStateTests")
	legacyStateTestDir101 = filepath.Join(baseDir, "LegacyTests", "Constantinople", "GeneralStateTests")
)

func debug() {

}

func withWritingTests(t *testing.T, name string, test *StateTest) {

	// Test output is written here.
	//fpath := filepath.Join(currentTestDir, name)
	//test.Name = strings.TrimSuffix(filepath.Base(fpath), ".json")

	fpath := name
	test.Name = strings.TrimSuffix(filepath.Base(name), ".json")

	// There is no need to run this git command for every test, but
	// speed is not really a big deal here, and it's nice to keep as much logic out
	// out the global scope as possible.
	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	// For tests using a config that does not have an associated chainspec file,
	// then generate that file.
	for _, subtest := range test.Subtests(nil) {
		subtest := subtest
		if _, ok := MapForkNameChainspecFileState[subtest.Fork]; !ok {
			genesis := test.genesis(Forks[subtest.Fork])
			pspec, err := tconvert.NewParityChainSpec(subtest.Fork, genesis, []string{})
			if err != nil {
				t.Fatal(err)
			}
			b, err := json.MarshalIndent(pspec, "", "    ")
			if err != nil {
				t.Fatal(err)
			}
			filename := paritySpecPath(strcase.ToSnake(subtest.Fork) + ".json")
			err = ioutil.WriteFile(filename, b, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
			sum := sha1.Sum(b)
			chainspecRefsState[subtest.Fork] = chainspecRef{filepath.Base(filename), sum[:]}
			t.Logf("Created new fork chainspec file: %v", chainspecRefsState[subtest.Fork])
		}
	}

	for _, subtest := range test.Subtests(nil) {
		subtest := subtest

		// Only proceed with test forks which are destined for writing.
		// Note that using this function implies that you trust the test runner
		// to give valid output, ie. only generate tests after you're sure the
		// reference tests themselves are passing.
		targetFork, ok := writeStateTestsReferencePairs[subtest.Fork]
		if !ok {
			// t.Logf("Skipping test (non-writing): %s", subtest.Fork)
			continue
		}

		if _, ok := test.json.Post[targetFork]; !ok {
			test.json.Post[targetFork] = make([]stPostState, len(test.json.Post[subtest.Fork]))
		}

		// Initialize the subtest/index data by copy from reference.
		referenceFork := subtest.Fork
		test.json.Post[targetFork][subtest.Index] = test.json.Post[referenceFork][subtest.Index]

		// Set new fork name, so new test config will be used instead.
		subtest.Fork = targetFork

		key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
		t.Run(key, func(t *testing.T) {
			vmConfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}

			// This is where the magic happens.
			err := test.RunSetPost(subtest, vmConfig)
			if err != nil {
				t.Fatalf("Error encountered at RunSetPost: %v", err)
			}

			// Only write the test once, after all subtests have been written.
			writeFile := filledPostStates(test.json.Post[subtest.Fork])
			if writeFile {
				fi, err := ioutil.ReadFile(fpath)
				if err != nil {
					t.Fatal("Error reading file, and will not write:", fpath, "test", key)
				}
				test.json.Info.WrittenWith = fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, head)
				test.json.Info.Parent = submoduleParentRef
				test.json.Info.ParentSha1Sum = fmt.Sprintf("%x", sha1.Sum(fi))
				test.json.Info.Chainspecs = chainspecRefsState

				b, err := json.MarshalIndent(test, "", "    ")
				if err != nil {
					t.Fatalf("Error marshaling JSON: %v", err)
				}

				err = ioutil.WriteFile(fpath, b, os.ModePerm)
				if err != nil {
					panic(err)
				}
			}

			_, _, err = test.Run(subtest, vmConfig, false)
			if err != nil {
				t.Fatalf("FAIL snap=false %v", err)
			}
			_, _, err = test.Run(subtest, vmConfig, true)
			if err != nil {
				t.Fatalf("FAIL snap=true %v", err)
			}

			if writeFile {
				t.Logf(`Wrote test file: %s
%s -> %s`, fpath, referenceFork, subtest.Fork)
			}
		})
	}
}
