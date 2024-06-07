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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/go-test/deep"
	"github.com/iancoleman/strcase"
)

// testMatcherGen embeds the common testMatcher struct, extending it to include information
// necessary for test generation.
type testMatcherGen struct {
	*testMatcher
	references []*regexp.Regexp
	targets    []string

	gitHead     string
	errorPanics bool

	// allConfigs tracks each unique configuration used in a test-generating suite.
	// This configuration map will be written to a file included with the tests as a reference
	// so that cross-client tests can define what eg. "Berlin" means.
	allConfigsMu sync.Mutex // safety first
	allConfigs   map[string]*coregeth.CoreGethChainConfig
}

// generateFromReference assigns reference:target pairs for test generation by fork name.
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

func (tg *testMatcherGen) getGenerationTarget(ref string) string {
	// Lookup the reference:target pairing, if any.
	for i, r := range tg.references {
		if r.MatchString(ref) {
			return tg.targets[i]
		}
	}
	return ""
}

// TestGenStateAll generates tests for all State tests.
func TestGenStateAll(t *testing.T) {
	if os.Getenv(CG_GENERATE_STATE_TESTS_KEY) == "" {
		t.Skip()
	}
	if os.Getenv(CG_CHAINCONFIG_FEATURE_EQ_COREGETH_KEY) == "" {
		t.Fatal("Must use core-geth chain configs for test generation, converting if necessary.")
	}

	// There is no need to run this git command for every test, but
	// speed is not really a big deal here, and it's nice to keep as much logic out
	// out the global scope as possible.
	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	tm := new(testMatcherGen)
	tm.allConfigs = make(map[string]*coregeth.CoreGethChainConfig)
	tm.testMatcher = new(testMatcher)
	tm.noParallel = true
	tm.errorPanics = true
	tm.gitHead = head

	tm.generateFromReference("Byzantium", "ETC_Atlantis")
	tm.generateFromReference("ConstantinopleFix", "ETC_Agharta")
	tm.generateFromReference("Berlin", "ETC_Magneto")
	tm.generateFromReference("Istanbul", "ETC_Phoenix")
	tm.generateFromReference("London", "ETC_Mystique")

	for _, dir := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		tm.walkFullName(t, dir, tm.testWriteTest)
	}
}

// TestGenStateSingles generates tests for a one or a few test files.
// This should only be used as a debugging tool.
// For production use, use TestGenStateAll.
func TestGenStateSingles(t *testing.T) {
	if os.Getenv(CG_GENERATE_STATE_TESTS_KEY) == "" {
		t.Skip("Use CG_GENERATE_STATE_TESTS=1 to generate state tests.")
	}
	if os.Getenv(CG_CHAINCONFIG_FEATURE_EQ_COREGETH_KEY) == "" {
		t.Fatal("Must use core-geth chain configs for test generation, converting if necessary.")
	}

	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	filesEnv := os.Getenv("TEST_FILES")
	files := []string{}
	for _, f := range strings.Split(filesEnv, ",") {
		files = append(files, filepath.Join(stateTestDir, f))
	}

	tm := new(testMatcherGen)
	tm.allConfigs = make(map[string]*coregeth.CoreGethChainConfig)
	tm.testMatcher = new(testMatcher)
	tm.noParallel = true // disable parallelism
	tm.errorPanics = true
	tm.gitHead = head

	tm.generateFromReference("Byzantium", "ETC_Atlantis")
	tm.generateFromReference("ConstantinopleFix", "ETC_Agharta")
	tm.generateFromReference("Berlin", "ETC_Magneto")
	tm.generateFromReference("Istanbul", "ETC_Phoenix")
	tm.generateFromReference("London", "ETC_Mystique")

	for _, f := range files {
		tm.runTestFile(t, f, f, tm.testWriteTest)
	}
}

var generatedBasedir = filepath.Join(".", "testdata_generated")

func (tm *testMatcherGen) testWriteTest(t *testing.T, name string, test *StateTest) {
	// testWriteTest generates a file-based tests (writing a new file), and re-runs (testing) the generated test file.

	// Set up a temporary file to write the generated test(s) to.
	// We have to use an actual file because the on-write callback re-runs the
	// test using the newly-generated file to ensure consistency.
	// If we didn't rerun the tests after writing, we could use a buffer as the WriteCloser.
	// Note that parallelism can cause greasy bugs around file during read/write which is why
	// we use a temporary file instead of immediately overwriting the canonical file in the first place;
	// for example, I saw regular encoding errors without this pattern.
	tmpFile, err := os.CreateTemp(os.TempDir(), "geth-state-test-generation")
	if err != nil {
		t.Fatal(err)
	}
	tmpFileName := tmpFile.Name()

	tm.runTestFile(t, name, name, tm.stateTestsGen(tmpFile,
		// On-Write:
		// After generating the tests, rerun the test using the new test file, ensuring consistency.
		// If the test passes (and it damn well should), then move the file to the canonical location.
		func() {
			tm.runTestFile(t, tmpFileName, tmpFileName, tm.stateTestRunner)

			rel, err := filepath.Rel(baseDir, name)
			if err != nil {
				t.Fatal(err)
			}

			targetPath := filepath.Join(generatedBasedir, rel)

			if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
				t.Fatal(err)
			}

			if err := os.Rename(tmpFileName, targetPath); err != nil {
				t.Fatal(err)
			}

			// Write all relevant configs.

			relSpl := strings.Split(rel, string(filepath.Separator))
			targetDirCommon := filepath.Join(generatedBasedir, relSpl[0]) // e.g. "testdata_generated/GeneralStateTests"

			// FIXME This is obviously super redundant.
			sts := test.Subtests(nil)
			for _, s := range sts {
				target := tm.getGenerationTarget(s.Fork)
				if target == "" {
					continue
				}
				conf, _, err := GetChainConfig(target)
				if err != nil {
					t.Fatal(err)
				}
				b, _ := json.MarshalIndent(conf, "", "    ")
				configPathTarget := filepath.Join(targetDirCommon, "configs", fmt.Sprintf("%s_config.json", target)) // e.g. "testdata_generated/GeneralStateTests/ETC_Atlantis_config.json"
				os.MkdirAll(filepath.Dir(configPathTarget), os.ModePerm)
				if _, statErr := os.Stat(configPathTarget); os.IsNotExist(statErr) {
					if err := os.WriteFile(configPathTarget, b, os.ModePerm); err != nil {
						t.Fatal(err)
					}
				}
			}
		},
		// On-Skip:
		// This test is skipped by some condition in the test matcher, probably by skip-fork.
		// In this case we just clean up the inoperative temp file.
		func() {
			if err := os.RemoveAll(tmpFileName); err != nil {
				t.Fatal(err)
			}
		}))
}

// stateTestsGen generates state tests using a reference fork and targeting a target fork.
// The reference fork is used to the pre-state and tested transaction(s) (which are schematized as indexes),
// replacing the reference fork's chain config and genesis with that of the target fork.
// The resulting post-state is assigned to the test's post.Root and post.Logs hashes.
func (tm *testMatcherGen) stateTestsGen(w io.WriteCloser, writeCallback, skipCallback func()) func(t *testing.T, name string, test *StateTest) {
	return func(t *testing.T, name string, test *StateTest) {
		subtests := test.Subtests(nil)

		// targets are the subtests that will be generated
		targets := map[string][]stPostState{}

		// eip1559ExistsAnySubtest will be used to omit `currentBaseFee` field
		// from the test env if EIP1559 not present in any of the subtest configs.
		eip1559ExistsAnySubtest := false

		for _, s := range subtests {
			// Prior to test-generation logic, record the genesis+chain config at the testmatcher level.
			// This will allow us to generate a complete map of chain configurations for the test suite,
			// whether the subtest's fork was used to generate any tests or not.
			// Record the genesis+chain config at the testmatcher level.
			config, _, err := GetChainConfig(s.Fork)
			if err != nil {
				t.Fatal(UnsupportedForkError{s.Fork})
			}

			tm.allConfigsMu.Lock()
			// This will panic if type isn't expected.
			tm.allConfigs[s.Fork] = config.(*coregeth.CoreGethChainConfig)
			tm.allConfigsMu.Unlock()

			// Lookup the reference:target pairing, if any.
			targetFork := tm.getGenerationTarget(s.Fork)

			// If target fork is empty we know that this subtest is not intended for as a reference
			// for any target generated test.
			if targetFork == "" {
				continue
			}
			referenceFork := s.Fork

			if _, ok := Forks[targetFork]; !ok {
				t.Fatalf("missing target fork config: %s, reference: %s", targetFork, referenceFork)
			}

			// Check to see if EIP1559 is enabled in the target subtest config.
			tconfig, _, err := GetChainConfig(targetFork)
			if err != nil {
				t.Fatal(UnsupportedForkError{s.Fork})
			}
			if config.IsEnabled(tconfig.GetEIP1559Transition, big.NewInt(int64(test.json.Env.Number))) {
				eip1559ExistsAnySubtest = true
			}

			if _, ok := targets[targetFork]; !ok {
				subtestsLen := len(test.json.Post[referenceFork])
				targets[targetFork] = make([]stPostState, subtestsLen)
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
				TxBytes:         refPostState[s.Index].TxBytes,
				ExpectException: refPostState[s.Index].ExpectException,
				// Root: <This is set pending the results of the Run under the target subtests.>
				// Logs: <This is set pending the results of the Run under the target subtests.>
			}

			// vmConfig is constructed using global variables for possible EVM and EWASM interpreters.
			// These interpreters are configured with environment variables and are assigned in an init() function.
			vmConfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}

			targetSubtest := StateSubtest{
				Fork:  targetFork,
				Index: s.Index,
			}

			state, root, err := test.RunNoVerifyWithPost(targetSubtest, vmConfig, false, rawdb.HashScheme, stPost)
			if err != nil {
				// Our runner has returned an error.
				// This can either be an intentional error (testing for the error), or an "unexpected" error,
				// which, since we're transposing tests, could still be a desirable error.
				// For example, London adopts EIP1559, while its counterpart, Mystique, does not.
				// In this case, the London test might return no error, while Mystique would return an error
				// (eg "unsupported transaction type").

				// An error was returned, but none defined for this test.
				// We write the error to the generated test.
				if refPostState[s.Index].ExpectException == "" {
					// TODO: Turn this error into the kind of error constants that upstream uses, eg. TR_TypeNotSupported.
					stPost.ExpectException = err.Error()
					if errors.Is(err, types.ErrTxTypeNotSupported) {
						stPost.ExpectException = "TR_TypeNotSupported"
					}
				}

				// Either way, we maintain the incumbent post state values,
				// although we do not expect our runner to return these for us.
				stPost.Root = refPostState[s.Index].Root
				stPost.Logs = refPostState[s.Index].Logs
			}

			if err == nil {
				if refPostState[s.Index].ExpectException != "" {
					// An error was expected, but none returned.
					// We overwrite the expected error to a zero value, because it didn't fail under our target configuration.
					stPost.ExpectException = ""
				}
				// If no error was returned, we can safely expect the root and state value to exist for us.
				stPost.Root = common.UnprefixedHash(root)
				stPost.Logs = common.UnprefixedHash(rlpHash(state.StateDB.Logs()))
			}

			targets[targetFork][s.Index] = stPost
		}

		if len(targets) == 0 {
			t.Skip("No targets found for this test")
			skipCallback()
			return
		}

		if !eip1559ExistsAnySubtest {
			test.json.Env.BaseFee = nil
		}

		// Install the generated cases to the test.
		test.json.Post = make(map[string][]stPostState) // Reset the post state. We'll only add the generated post states.
		for k, v := range targets {
			test.json.Post[k] = v
		}

		// Assign provenance metadata to the test.
		test.Name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
		test.json.Info.FilledWith = fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, tm.gitHead)

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

		t.Run(key+"/trie", func(t *testing.T) {
			// vmConfig is constructed using global variables for possible EVM and EWASM interpreters.
			// These interpreters are configured with environment variables and are assigned in an init() function.
			vmConfig := vm.Config{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
			test.Run(st, vmConfig, false, rawdb.HashScheme, func(err error, state *StateTestState) {
				if err != nil && len(test.json.Post[st.Fork][st.Index].ExpectException) > 0 {
					// Ignore expected errors (TODO MariusVanDerWijden check error string)
					return
				}
				checkedErr := tm.checkFailure(t, err)
				if checkedErr != nil && *testEWASM != "" {
					checkedErr = fmt.Errorf("%w ewasm=%s", checkedErr, *testEWASM)
				}
				if checkedErr != nil {
					if tm.errorPanics {
						panic(checkedErr)
					} else {
						t.Fatal(err)
					}
				}
			})
		})
	}
}

// TestGenStateCoreGethConfigs generates core-geth-style configurations.
// This isn't a test. It generates configs.
// Skip should be installed so this function will only be run by developers as needed.
func TestGenStateCoreGethConfigs(t *testing.T) {
	t.Skip()
	st := new(testMatcher)

	// FYI: Possibly the slowest and stupidest way to write 14 files: read 42189 test to do it
	// and write each file 1486 times.
	for _, d := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		st.walkFullName(t, d, func(t *testing.T, name string, test *StateTest) {
			subtests := test.Subtests(nil)
			for _, subtest := range subtests {
				subtest := subtest

				genesis := test.genesis(Forks[subtest.Fork])

				cgConfig := &coregeth.CoreGethChainConfig{}
				err := confp.Crush(cgConfig, genesis.Config, true)
				if err != nil {
					t.Fatal(err)
				}
				genesis.Config = cgConfig

				b, err := json.MarshalIndent(genesis, "", "    ")
				if err != nil {
					t.Fatal(err)
				}
				filename := filepath.Join(
					coregethSpecsDir,
					strcase.ToSnake(subtest.Fork)+"_test.json",
				)
				err = os.WriteFile(filename, b, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestGeneratedConfigsEq tests that the CoreGeth configuration described in
// ../params/coregeth.json.d/berlin_test.json (a configuration which is relevant to
// state tests) is equivalent to the configuration coded at Forks["Berlin"].
// Note that when run with COREGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALENCE_COREGETH, for example,
// the coded configuration will have been init'd with a CoreGeth-struct value (as opposed
// to the default go-ethereum value).
func TestGeneratedConfigsEq(t *testing.T) {
	specPath := filepath.Join(coregethSpecsDir, "berlin_test.json")
	gen := &genesisT.Genesis{
		Config: &coregeth.CoreGethChainConfig{},
	}
	_, err := readJSONFromFile(specPath, gen)
	if err != nil {
		t.Fatal(err)
	}

	coded := Forks["Berlin"]

	// Special case handling for EIP1283.
	if coded.GetEIP1283Transition() == nil && coded.GetEIP1283DisableTransition() == nil {
		e1283 := gen.Config.GetEIP1283Transition()
		d1283 := gen.Config.GetEIP1283DisableTransition()

		if (e1283 == nil && d1283 == nil) ||
			(*e1283 == *d1283) {
			gen.Config.SetEIP1283Transition(nil)
			gen.Config.SetEIP1283DisableTransition(nil)
		}
	}

	err = confp.Equivalent(gen.Config, coded)
	if err != nil {
		t.Error(err)
	}
}

// TestConvertDefaultsBounce tests that for a few default Forks configuration values (which are used in tests),
// that the convert-marshal-unmarshal cycle (here, a "bounce") results in equivalent chain configurations.
// EIP1234 is used as a canary test since I have seen that fail most lately, and it involved some complicated
// encoding and interdependencies (is Ethash, rel EIP649, map encoding vs. inferrable boolean field).
func TestConvertDefaultsBounce(t *testing.T) {
	safePrint := func(n *uint64) string {
		if n == nil {
			return "nil"
		}
		return fmt.Sprintf("%d", *n)
	}
	for _, forkName := range []string{"Constantinople", "Istanbul", "Berlin", "London", "ArrowGlacier"} {
		t.Run(forkName, func(t *testing.T) {
			forkConfig := Forks[forkName]
			eip1234 := forkConfig.GetEthashEIP1234Transition()

			// Crush the original config (type) into the CoreGeth chain config data type.
			coreGethConfig := &coregeth.CoreGethChainConfig{}
			err := confp.Crush(coreGethConfig, forkConfig, true)
			if err != nil {
				t.Fatal(err)
			}

			// Read the EIP1234 value from the converted version.
			cg1234 := coreGethConfig.GetEthashEIP1234Transition()

			t.Logf("DEBUG1/%s eip=%s cg=%s", forkName, safePrint(eip1234), safePrint(cg1234))

			// Assert that the converted chain configuration data type is equivalent to the original.
			err = confp.Equivalent(forkConfig, coreGethConfig)
			if err != nil {
				t.Fatal(err)
			}

			b, _ := json.MarshalIndent(coreGethConfig, "", "    ")
			t.Log(string(b))

			coreGethConfig2 := &coregeth.CoreGethChainConfig{}
			err = json.Unmarshal(b, coreGethConfig2)
			if err != nil {
				t.Fatal(err)
			}

			cg2_1234 := coreGethConfig2.GetEthashEIP1234Transition()

			t.Logf("DEBUG2/%s eip=%s cg=%s cg2=%s", forkName, safePrint(eip1234), safePrint(cg1234), safePrint(cg2_1234))

			err = confp.Equivalent(forkConfig, coreGethConfig2)
			if err != nil {
				jb, err := json.MarshalIndent(coreGethConfig2, "", "    ")
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("%s", string(jb))

				for _, d := range deep.Equal(coreGethConfig, coreGethConfig2) {
					t.Logf("diff: %s", d)
				}

				t.Fatal(forkName, err)
			}
		})
	}
}
