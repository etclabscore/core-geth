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

package tests

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/vars"
)

// outNDJSONFile is the file where difficulty tests get written to.
// The file is encoded as line-delimited JSON.
var outNDJSONFile = filepath.Join(difficultyTestDir, "mgen_difficulty.ndjson")

// TestDifficultyTestConfigGen generates the difficulty test configuration files
// for all existing tests' configuration and the configuration for any to-be generated
// test configurations (via dt.generateFromReference).
func TestDifficultyTestConfigGen(t *testing.T) {
	if os.Getenv(CG_GENERATE_DIFFICULTY_TEST_CONFIGS_KEY) == "" {
		t.Skip()
	}

	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	dt := new(testMatcherGen)
	dt.allConfigs = make(map[string]*coregeth.CoreGethChainConfig)
	dt.testMatcher = new(testMatcher)
	dt.noParallel = true // disable parallelism
	dt.errorPanics = true
	dt.gitHead = head

	// Not difficulty-tests
	dt.skipLoad("hexencodetest.*")
	dt.skipLoad("crypto.*")
	dt.skipLoad("blockgenesistest\\.json")
	dt.skipLoad("genesishashestest\\.json")
	dt.skipLoad("keyaddrtest\\.json")
	dt.skipLoad("txtest\\.json")

	// files are 2 years old, contains strange values
	dt.skipLoad("difficultyCustomHomestead\\.json")
	dt.skipLoad("difficultyMorden\\.json")
	dt.skipLoad("difficultyOlimpic\\.json")

	dt.generateFromReference("Byzantium", "ETC_Atlantis")
	dt.generateFromReference("Constantinople", "ETC_Agharta")
	dt.generateFromReference("EIP2384", "ETC_Phoenix")
	/*
		My rationale for not adding ETC_Magneto was that difficulty hasn't changed for the Foundation since EIP2384 Muir Glacier,
		and the tests haven't changed in at least that long.
		This leads to me to think that adding an ETC_Magneto case would only duplicate the set of ETC_Phoenix.
	*/

	for k, v := range difficultyChainConfigurations {
		dt.config(k, v)
	}

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		config, matchedName := dt.findConfig(name)
		// t.Logf("name: %s, matchedName: %s", name, matchedName)

		fileBasename, ok := mapForkNameChainspecFileDifficulty[matchedName]
		if !ok {
			t.Fatalf("unmatched config name: %s", matchedName)
		}
		specPath := filepath.Join(coregethSpecsDir, fileBasename)

		cgConfig := &coregeth.CoreGethChainConfig{}
		err := confp.Convert(config, cgConfig)
		if err != nil {
			t.Fatal(err)
		}

		j, err := json.MarshalIndent(cgConfig, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(specPath, j, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		targetForkName := dt.getGenerationTarget(matchedName)
		if targetForkName == "" {
			return
		}

		targetConfiguration, ok := difficultyChainConfigurations[targetForkName]
		if !ok {
			t.Fatalf("config association failed; no existing Go chain config found: %s", targetForkName)
		}

		targetConfigurationBasename := mapForkNameChainspecFileDifficulty[targetForkName]
		specPath = filepath.Join(coregethSpecsDir, targetConfigurationBasename)

		cgConfig = &coregeth.CoreGethChainConfig{}
		err = confp.Convert(targetConfiguration, cgConfig)
		if err != nil {
			t.Fatal(err)
		}

		j, err = json.MarshalIndent(cgConfig, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(specPath, j, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestDifficultyGen generated line-delimited JSON tests using the existing default
// difficulty tests suite as a base reference set.
// Reference:Target pairs are defined with dt.generateFromReference.
// The original test will be run and must pass in order for it to be copied to the
// generated set and for it to be used as a reference, if applicable, for test generation.
func TestDifficultyGen(t *testing.T) {
	if os.Getenv(CG_GENERATE_DIFFICULTY_TESTS_KEY) == "" {
		t.Skip()
	}

	head := build.RunGit("rev-parse", "HEAD")
	head = strings.TrimSpace(head)

	err := os.MkdirAll(filepath.Dir(outNDJSONFile), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// Truncate/touch output file.
	err = ioutil.WriteFile(outNDJSONFile, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	dt := new(testMatcherGen)
	dt.allConfigs = make(map[string]*coregeth.CoreGethChainConfig)
	dt.testMatcher = new(testMatcher)
	dt.noParallel = true // disable parallelism
	dt.errorPanics = true
	dt.gitHead = head

	dt.generateFromReference("Byzantium", "ETC_Atlantis")
	dt.generateFromReference("Constantinople", "ETC_Agharta")
	dt.generateFromReference("EIP2384", "ETC_Phoenix")
	/*
		My rationale for not adding ETC_Magneto was that difficulty hasn't changed for the Foundation since EIP2384 Muir Glacier,
		and the tests haven't changed in at least that long.
		This leads to me to think that adding an ETC_Magneto case would only duplicate the set of ETC_Phoenix.
	*/

	// Not difficulty-tests
	dt.skipLoad("hexencodetest.*")
	dt.skipLoad("crypto.*")
	dt.skipLoad("blockgenesistest\\.json")
	dt.skipLoad("genesishashestest\\.json")
	dt.skipLoad("keyaddrtest\\.json")
	dt.skipLoad("txtest\\.json")

	// files are 2 years old, contains strange values
	dt.skipLoad("difficultyCustomHomestead\\.json")
	dt.skipLoad("difficultyMorden\\.json")
	dt.skipLoad("difficultyOlimpic\\.json")
	dt.skipLoad("all_difficulty_tests\\.json")

	for k, v := range difficultyChainConfigurations {
		dt.config(k, v)
	}

	mustSha1SumForFile := func(filePath string) []byte {
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatal(err)
		}

		s1s := sha1.Sum(b)
		return s1s[:]
	}

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		cfg, key := dt.findConfig(name)

		if test.ParentDifficulty.Cmp(vars.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, test.Run(cfg)); err != nil {
			t.Fatalf("failed to run difficulty test, err=%v", err)
		} else {

			fileBasename, ok := mapForkNameChainspecFileDifficulty[key]
			if !ok {
				t.Fatalf("unmatched config name: %s", key)
			}
			specPath := filepath.Join(coregethSpecsDir, fileBasename)

			test.Chainspec = chainspecRef{Filename: fileBasename, Sha1Sum: mustSha1SumForFile(specPath)}
			test.Name = strings.ReplaceAll(name, ".json", "")

			mustAppendTestToFileNDJSON(t, test, outNDJSONFile)
			// -- This is as far as ALL tests go.

			targetForkName := dt.getGenerationTarget(key)
			if targetForkName == "" {
				// This configuration is not a reference for a generated target.
				return
			}

			// Lookup Target configuration from canonical coded list, eg.
			/*
				"Ropsten":  params.RopstenChainConfig,
				"Morden":   params.RopstenChainConfig,
				"Frontier": &goethereum.ChainConfig{},
			*/
			targetConfiguration, ok := difficultyChainConfigurations[targetForkName]
			if !ok {
				t.Fatalf("config association failed; no existing Go chain config found: %s", targetForkName)
			}

			// Lookup respective Target configuration FILE name, eg.
			/*
				"Ropsten":           "ropsten_difficulty_test.json",
				"Morden":            "morden_difficulty_test.json",
				"Frontier":          "frontier_difficulty_test.json",
			*/
			targetConfigurationBasename := mapForkNameChainspecFileDifficulty[targetForkName]

			// Establish the specification reference value.
			// Note that this assumes that the files sourced here are indeed consistent with the
			// the coded configurations (ie. that TestDifficultyTestConfigGen has already been
			// run and has done its job properly).
			// There is no coded relationship or test between the test file and the configuration
			// actually used here.
			targetSpecRef := chainspecRef{
				Filename: targetConfigurationBasename,
				Sha1Sum:  mustSha1SumForFile(filepath.Join(coregethSpecsDir, targetConfigurationBasename))}

			newTest := &DifficultyTest{
				ParentTimestamp:    test.ParentTimestamp,
				ParentDifficulty:   test.ParentDifficulty,
				UncleHash:          test.UncleHash,
				CurrentTimestamp:   test.CurrentTimestamp,
				CurrentBlockNumber: test.CurrentBlockNumber,
				CurrentDifficulty: ethash.CalcDifficulty(targetConfiguration, test.CurrentTimestamp, &types.Header{
					Difficulty: test.ParentDifficulty,
					Time:       test.ParentTimestamp,
					Number:     big.NewInt(int64(test.CurrentBlockNumber - 1)),
					UncleHash:  test.UncleHash,
				}),
				Chainspec: targetSpecRef,
				Name:      strings.ReplaceAll(test.Name, key, targetForkName),
			}

			// "Dogfood".
			if err := newTest.Run(targetConfiguration); err != nil {
				t.Fatal(err)
			}
			mustAppendTestToFileNDJSON(t, newTest, outNDJSONFile)
			t.Logf("OK [generated] %v", newTest)

		}
	})
}

// mustAppendTestToFileNDJSON appends a difficulty test to a file as newline-delimited JSON.
func mustAppendTestToFileNDJSON(t *testing.T, test *DifficultyTest, filep string) {
	b, _ := json.Marshal(test)
	out := []byte{}
	buf := bytes.NewBuffer(out)
	err := json.Compact(buf, b)
	if err != nil {
		t.Fatal(err)
	}
	buf.Write([]byte("\n"))

	fi, err := os.OpenFile(filep, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = fi.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	err = fi.Close()
	if err != nil {
		t.Fatal(err)
	}
}
