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
		forkPair, ok := writeStateTestsReferencePairs[subtest.Fork]
		if !ok {
			// t.Logf("Skipping test (non-writing): %s", subtest.Fork)
			continue
		}

		if _, ok := test.json.Post[forkPair]; !ok {
			test.json.Post[forkPair] = make([]stPostState, len(test.json.Post[subtest.Fork]))
		}

		// Initialize the subtest/index data by copy from reference.
		reftestFork := subtest.Fork
		test.json.Post[forkPair][subtest.Index] = test.json.Post[reftestFork][subtest.Index]

		// Set new fork name, so new test config will be used instead.
		subtest.Fork = forkPair

		key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)

		t.Run(key, func(t *testing.T) {
			withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
				err := test.RunSetPost(subtest, vmconfig)

				// Only write the test once, after all subtests have been written.
				if err == nil && filledPostStates(test.json.Post[subtest.Fork]) {
					b, err := json.MarshalIndent(test, "", "    ")
					if err != nil {
						return err
					}
					fi, err := ioutil.ReadFile(fpath)
					if err != nil {
						t.Fatal("Error reading file, and will not write:", fpath, "test", string(b))
						return nil
					}
					test.json.Info.WrittenWith = fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, head)
					test.json.Info.Parent = submoduleParentRef
					test.json.Info.ParentSha1Sum = fmt.Sprintf("%x", sha1.Sum(fi))
					test.json.Info.Chainspecs = chainspecRefsState

					err = ioutil.WriteFile(fpath, b, os.ModePerm)
					if err != nil {
						panic(err)
					}
					t.Logf("Wrote test file: %s\n", fpath)
				} else {
					t.Errorf("Error encountered at RunSetPost: %v", err)
				}
				return nil
			})
		})
	}
}
