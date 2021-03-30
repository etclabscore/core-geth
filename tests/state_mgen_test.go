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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/confp/tconvert"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
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

	st.walkFullName(t, stateTestDir, st.withWritingTests)

	// For Istanbul, older tests were moved into LegacyTests
	// st.walkFullName(t, legacyStateTestDir, st.withWritingTests)
}

type testMatcherGen struct {
	*testMatcher
	references []*regexp.Regexp
	targets    []string

	gitHead     string
	errorPanics bool
}

func (tg *testMatcherGen) generateFromReference(ref, target string) {
	if tg.references == nil {
		tg.references = []*regexp.Regexp{}
	}
	tg.references = append(tg.references, regexp.MustCompile(ref))
	if tg.targets == nil {
		tg.targets = []string{}
	}
	tg.targets = append(tg.targets, target)
}

func (tm *testMatcherGen) stateTestsGen(w io.WriteCloser, writeCallback, skipCallback func()) func(t *testing.T, name string, test *StateTest) {
	return func(t *testing.T, name string, test *StateTest) {

		subtests := test.Subtests(nil)

		targets := map[string][]stPostState{}

		for _, s := range subtests {
			// Lookup the reference:target pairing, if any.
			var referenceFork, targetFork string
			for i, r := range tm.references {
				if r.MatchString(s.Fork) {
					referenceFork = s.Fork
					targetFork = tm.targets[i]
					break
				}
			}
			if referenceFork == "" {
				continue
			}
			if _, ok := Forks[targetFork]; !ok {
				t.Fatalf("missing target fork config: %s, reference: %s", targetFork, referenceFork)
			}

			if _, ok := targets[targetFork]; !ok {
				subtestsLen := len(test.json.Post[referenceFork])
				targets[targetFork] = make([]stPostState, subtestsLen)
			}

			targetSubtest := StateSubtest{
				Fork:  targetFork,
				Index: s.Index,
			}

			refPostState := test.json.Post[referenceFork]

			// Initialize the post state with reference indexes.
			// These indexes (containing gas, value, and data) are used to construct the message (tx) run
			// by the test.
			stPost := stPostState{
				Indexes: stPostStateIndexes{
					Data:  refPostState[s.Index].Indexes.Data,
					Gas:   refPostState[s.Index].Indexes.Gas,
					Value: refPostState[s.Index].Indexes.Value,
				},
			}

			// vmConfig is constructed using global variables for possible EVM and EWASM interpreters.
			// These interpreters are configured with environment variables and are assigned in an init() function.
			vmConfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}

			// Since we know that tests run with and without the snapshotter features are equivalent, either boolean state
			// is valid and 'false' is arbitrary.
			_, statedb, root, err := test.RunNoVerifyWithPost(targetSubtest, vmConfig, false, stPost)
			if err != nil {
				t.Fatalf("Error encountered at RunSetPost: %v", err)
			}

			// Assign the generated testable values.
			stPost.Root = common.UnprefixedHash(root)
			stPost.Logs = common.UnprefixedHash(rlpHash(statedb.Logs()))

			targets[targetFork][s.Index] = stPost
		}

		if len(targets) == 0 {
			t.Skip()
			skipCallback()
			return
		}

		// Install the generated cases to the test.
		for k, v := range targets {
			test.json.Post[k] = v
		}

		// Assign provenance metadata to the test.
		test.json.Info.WrittenWith = fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, tm.gitHead)

		// Write the augmented test to the writer.
		generatedJSON, err := json.MarshalIndent(test, "", "    ")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(generatedJSON)
		if err != nil {
			panic(err)
		}
		if err := w.Close(); err != nil {
			panic(err)
		}
		writeCallback()
	}
}

func (tm *testMatcherGen) stateTestRunner(t *testing.T, name string, test *StateTest) {
	subtests := test.Subtests(nil)
	for _, subtest := range subtests {
		st := subtest

		key := fmt.Sprintf("%s/%d", st.Fork, st.Index)
		name := name + "/" + key

		t.Run(key+"/trie", func(t *testing.T) {
			// vmConfig is constructed using global variables for possible EVM and EWASM interpreters.
			// These interpreters are configured with environment variables and are assigned in an init() function.
			vmConfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
			_, _, err := test.Run(st, vmConfig, false)
			checkedErr := tm.checkFailure(t, name+"/trie", err)
			if checkedErr != nil && *testEWASM != "" {
				checkedErr = fmt.Errorf("%w ewasm=%s", checkedErr, *testEWASM)
			}
			if checkedErr != nil {
				if tm.errorPanics {
					panic(err)
				} else {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestGenStateParityConfigs generates parity-style configurations.
// This isn't a test. It generates configs.
// Skip should be installed so this function will only be run by developers as needed.
func TestGenStateParityConfigs(t *testing.T) {
	t.Skip()
	st := new(testMatcher)

	// FYI: Possibly the slowest and stupidest way to write 14 files: read 42189 test to do it
	// and write each file 1486 times.
	for _, d := range []string{stateTestDir, legacyStateTestDir} {
		st.walkFullName(t, d, func(t *testing.T, name string, test *StateTest) {

			subtests := test.Subtests(nil)
			for _, subtest := range subtests {
				subtest := subtest

				genesis := test.genesis(Forks[subtest.Fork])

				// Write the genesis+config in the Parity config format.
				pspec, err := tconvert.NewParityChainSpec(subtest.Fork, genesis, []string{})
				if err != nil {
					t.Fatal(err)
				}
				b, err := json.MarshalIndent(pspec, "", "    ")
				if err != nil {
					t.Fatal(err)
				}
				filename := filepath.Join(
					"..",
					"params",
					"parity.json.d",
					strcase.ToSnake(subtest.Fork)+".json",
				)
				err = ioutil.WriteFile(filename, b, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}

				// Also write it in coregeth format.
				cgspec := &coregeth.CoreGethChainConfig{}
				err = confp.Convert(Forks[subtest.Fork], cgspec)
				if err != nil {
					t.Fatal(err)
				}

				cgGenesis := test.genesis(cgspec)
				b, err = json.MarshalIndent(cgGenesis, "", "    ")
				if err != nil {
					t.Fatal(err)
				}

				filename = filepath.Join(
					"..",
					"params",
					"coregeth.json.d",
					strcase.ToSnake(subtest.Fork)+".json",
				)
				err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(filename, b, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}

				// We cannot write the config in any other formats because
				// go-ethereum and multi-geth are unable to describe some of the
				// features configured, eg. ECIPs and possibly others (eg. EIP2537).
			}
		})
	}
}

func TestGenStateAll(t *testing.T) {
	if os.Getenv(CG_GENERATE_STATE_TESTS_KEY) == "" {
		t.Skip()
	}
	// if os.Getenv(CG_CHAINCONFIG_CHAINSPECS_OPENETHEREUM_KEY) == "" {
	// 	t.Fatal("Must use chainspec files for fork configurations.")
	// }

	// There is no need to run this git command for every test, but
	// speed is not really a big deal here, and it's nice to keep as much logic out
	// out the global scope as possible.
	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	tm := new(testMatcherGen)
	tm.testMatcher = new(testMatcher)
	tm.noParallel = true
	tm.errorPanics = true
	tm.gitHead = head

	tm.generateFromReference("Byzantium", "ETC_Atlantis")
	tm.generateFromReference("ConstantinopleFix", "ETC_Agharta")
	tm.generateFromReference("Berlin", "ETC_Magneto")
	tm.generateFromReference("Istanbul", "ETC_Phoenix")

	for _, dir := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		tm.walkFullName(t, dir, tm.testWriteTest)
	}
}

func TestGenStateSingles(t *testing.T) {
	if os.Getenv(CG_GENERATE_STATE_TESTS_KEY) == "" {
		t.Skip()
	}
	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	files := []string{
		filepath.Join(stateTestDir, "stStaticFlagEnabled/DelegatecallToPrecompileFromContractInitialization.json"),
		filepath.Join(stateTestDir, "stStaticCall/StaticcallToPrecompileFromCalledContract.json"),
	}

	tm := new(testMatcherGen)
	tm.testMatcher = new(testMatcher)
	tm.noParallel = true // disable parallelism
	tm.errorPanics = true
	tm.gitHead = head

	tm.generateFromReference("Byzantium", "ETC_Atlantis")
	tm.generateFromReference("ConstantinopleFix", "ETC_Agharta")
	tm.generateFromReference("Berlin", "ETC_Magneto")
	tm.generateFromReference("Istanbul", "ETC_Phoenix")

	for _, f := range files {
		tm.runTestFile(t, f, f, tm.testWriteTest)
	}
}

func mustTemporaryGenerationPath(originalName string) (tempPath string) {
	targetName := strings.Replace(originalName, "testdata", "testdata_gen", 1)
	err := os.MkdirAll(filepath.Dir(targetName), os.ModePerm)
	if err != nil {
		panic(err)
	}
	testOutFullName, err := filepath.Abs(targetName)
	if err != nil {
		panic(err)
	}
	testOutFullName = filepath.Clean(testOutFullName)
	if err := os.Truncate(testOutFullName, 0); os.IsNotExist(err) {
		if _, err := os.Create(testOutFullName); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	return testOutFullName
}

func (tm *testMatcherGen) testWriteTest(t *testing.T, name string, test *StateTest) {
	testOutFullName := mustTemporaryGenerationPath(name)
	testOut, err := os.OpenFile(testOutFullName, os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// tm.runTestFile(t, name, name, tm.stateTestRunner)
	tm.runTestFile(t, name, name, tm.stateTestsGen(testOut, func() {
		tm.runTestFile(t, testOutFullName, testOutFullName, tm.stateTestRunner)
		os.Rename(testOutFullName, name)
	}, func() {
		os.RemoveAll(testOutFullName)
	}))
}

func (tm *testMatcher) withWritingTests(t *testing.T, name string, test *StateTest) {

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
	subtests := test.Subtests(nil)

	// for _, subtest := range subtests {
	// 	subtest := subtest
	// 	if _, ok := MapForkNameChainspecFileState[subtest.Fork]; !ok {
	// 		genesis := test.genesis(Forks[subtest.Fork])
	// 		pspec, err := tconvert.NewParityChainSpec(subtest.Fork, genesis, []string{})
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		b, err := json.MarshalIndent(pspec, "", "    ")
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		filename := paritySpecPath(strcase.ToSnake(subtest.Fork) + ".json")
	// 		err = ioutil.WriteFile(filename, b, os.ModePerm)
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		sum := sha1.Sum(b)
	// 		chainspecRefsState[subtest.Fork] = chainspecRef{filepath.Base(filename), sum[:]}
	// 		t.Logf("Created new fork chainspec file: %v", chainspecRefsState[subtest.Fork])
	// 	}
	// }

	for _, subtest := range subtests {
		subtest := subtest
		runTestGenerating(t, tm, fpath, test, subtest, head, subtests)
	}
}

func runTestGenerating(t *testing.T, tm *testMatcher, fpath string, test *StateTest, subtest StateSubtest, head string, subtests []StateSubtest) {

	originalJSON, _ := json.MarshalIndent(test, "", "    ")

	// Only proceed with test forks which are destined for writing.
	// Note that using this function implies that you trust the test runner
	// to give valid output, ie. only generate tests after you're sure the
	// reference tests themselves are passing.
	// eg. Istanbul => ETC_Phoenix; Berlin => ETC_Magneto
	referenceFork := subtest.Fork
	targetFork, ok := writeStateTestsReferencePairs[referenceFork]
	// if !ok {
	// 	// t.Logf("Skipping test (non-writing): %s", subtest.Fork)
	// 	return
	// }

	var crossReference stPostState

	if ok {
		for _, s := range subtests {
			if s.Fork == targetFork {
				// This will/would regenerate an existing test subtest.
				crossReference = test.json.Post[s.Fork][s.Index]
				break
			}
		}
	} else {
		return
	}

	// if _, ok := test.json.Post[targetFork]; !ok {
	test.json.Post[targetFork] = make([]stPostState, len(test.json.Post[referenceFork]))
	copy(test.json.Post[targetFork], test.json.Post[referenceFork])
	// }

	// Initialize the subtest/index data by copy from reference.
	// test.json.Post[targetFork][subtest.Index] = test.json.Post[referenceFork][subtest.Index]

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

		// Check against cross reference, if any.
		if crossReference != (stPostState{}) {
			if crossReference.Root != test.json.Post[targetFork][subtest.Index].Root {
				panic(fmt.Sprintf(`cross reference failed
referenceFork: %s
targetFork: %s
cross.root: %s
gend.root: %s
`, referenceFork,
					targetFork,
					crossReference.Root,
					test.json.Post[targetFork][subtest.Index].Root))
			}
		}

		// Only write the test once, after all subtests have been written.
		writeFile := filledPostStates(test.json.Post[targetFork])
		generatedJSON := []byte{}
		if writeFile {
			fi, err := ioutil.ReadFile(fpath)
			if err != nil {
				t.Fatal("Error reading file, and will not write:", fpath, "test", key)
			}
			test.json.Info.WrittenWith = fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, head)
			test.json.Info.Parent = submoduleParentRef
			test.json.Info.ParentSha1Sum = fmt.Sprintf("%x", sha1.Sum(fi))
			test.json.Info.Chainspecs = chainspecRefsState

			generatedJSON, err = json.MarshalIndent(test, "", "    ")
			if err != nil {
				t.Fatalf("Error marshaling JSON: %v", err)
			}

			err = ioutil.WriteFile(fpath, generatedJSON, os.ModePerm)
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
%s -> %s`, fpath, referenceFork, targetFork)

			// // Re-run the test we just wrote
			tm.runTestFile(t, fpath, fpath, func(t *testing.T, name2 string, test2 *StateTest) {
				ssubtests := test2.Subtests(nil)
				for _, subtest2 := range ssubtests {
					subtest2 := subtest2
					key := fmt.Sprintf("%s/%d", subtest2.Fork, subtest2.Index)
					name3 := name2 + "/" + key

					t.Run(key+"/trie", func(t *testing.T) {
						wrapFatal(t, test2.gasLimit(subtest2), func(vmconfig vm.Config) error {
							_, _, err := test2.Run(subtest2, vmconfig, false)
							checkedErr := tm.checkFailure(t, name3+"/trie", err)
							if checkedErr != nil && *testEWASM != "" {
								checkedErr = fmt.Errorf("%w ewasm=%s", checkedErr, *testEWASM)
							}
							if checkedErr != nil {
								ioutil.WriteFile("original.json", originalJSON, os.ModePerm)
								ioutil.WriteFile("generated.json", generatedJSON, os.ModePerm)
								panic(checkedErr)
							}
							return checkedErr
						})
					})
				}
			})

		}
	})
}
