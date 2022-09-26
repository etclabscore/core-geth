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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

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

// TestDifficultyGen2 generates JSON tests from scratch.
// The test case matrix can be deduced from the for-loop iterations.
// The cases are written to files PER CHAIN CONFIG, following the upstream convention,
// eg. tests/testdata_generated/BasicTests/difficultyETC_Agharta.json
func TestDifficultyGen2(t *testing.T) {
	if os.Getenv(CG_GENERATE_DIFFICULTY_TESTS_KEY) == "" {
		t.Skip()
	}

	configs := map[string]ctypes.ChainConfigurator{
		"ETC": params.ClassicChainConfig,
	}

	targetDir := filepath.Join(generatedBasedir, "DifficultyTests", "dfETC")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Establish test matrix values.
	//
	// Multiples of 8 offset by 1, maxing out at 121.
	timestampOffsets := func() []uint64 {
		var r []uint64
		for i := 1; i <= 121; i += 8 {
			r = append(r, uint64(i))
		}
		return r
	}()

	parentDifficulty := new(big.Int).Mul(vars.MinimumDifficulty, big.NewInt(100))

	filledTests := map[string]*DifficultyTest{}
	for configName, config := range configs {
		testsName := fmt.Sprintf("difficulty%s", configName)
		targetFileBaseName := fmt.Sprintf("difficulty%s.json", configName)
		targetFilePath := filepath.Join(targetDir, targetFileBaseName)
		os.Truncate(targetFilePath, 0)
		forks := confp.Forks(config)
		writeForks := []uint64{}
		for _, f := range forks {
			writeForks = append(writeForks, f)
			if f >= 1 {
				writeForks = append(writeForks, f-1)
			}
			writeForks = append(writeForks, f+1)
		}
		for _, blockNumber := range writeForks {
			for _, timestampOffset := range timestampOffsets {
				for _, uncle := range []bool{false, true} {
					uncleCount := 0
					if uncle {
						uncleCount = 1
					}
					uncleHash := types.EmptyUncleHash
					if uncle {
						uncleHash = types.CalcUncleHash([]*types.Header{{Number: common.Big1}})
					}

					// Establish test parameters.
					newTest := &DifficultyTest{
						ParentTimestamp:    blockNumber*13 + 13,
						ParentDifficulty:   parentDifficulty,
						ParentUncles:       uint64(uncleCount),
						CurrentBlockNumber: blockNumber,
						// CurrentTimestamp:  This gets filled later.
						// CurrentDifficulty: This gets filled later.
					}

					newTest.CurrentTimestamp = newTest.ParentTimestamp + timestampOffset

					// Fill the expected difficulty from the test params we've just established.
					newTest.CurrentDifficulty = ethash.CalcDifficulty(config, newTest.CurrentTimestamp, &types.Header{
						Difficulty: newTest.ParentDifficulty,
						Time:       newTest.ParentTimestamp,
						Number:     big.NewInt(int64(newTest.CurrentBlockNumber - 1)),
						UncleHash:  uncleHash,
					})

					filledTests[fmt.Sprintf("difficulty_n%d_t%d_u%t", blockNumber, timestampOffset, uncle)] = newTest

					mustWriteTestFileJSON2(t, targetFilePath, filledTests, testsName, configName)
				}
			}
		}
	}
}

func mustWriteTestFileJSON2(tt *testing.T, filePath string, tests map[string]*DifficultyTest, testName string, configName string) {
	enc := make(map[string]json.RawMessage)
	for k, v := range tests {
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			tt.Fatal(err)
		}
		enc[k] = b
	}
	enc2 := map[string]interface{}{
		testName: map[string]interface{}{
			configName: enc,
		},
	}
	b, err := json.MarshalIndent(enc2, "", "    ")
	if err != nil {
		tt.Fatal(err)
	}
	err = ioutil.WriteFile(filePath, b, os.ModePerm)
	if err != nil {
		tt.Fatal(err)
	}
}

// TestDifficultyGen generates JSON tests from scratch.
// The test case matrix can be deduced from the for-loop iterations.
// The cases are written to files PER CHAIN CONFIG, following the upstream convention,
// eg. tests/testdata_generated/BasicTests/difficultyETC_Agharta.json
func TestDifficultyGen(t *testing.T) {
	if os.Getenv(CG_GENERATE_DIFFICULTY_TESTS_KEY) == "" {
		t.Skip()
	}

	configs := map[string]ctypes.ChainConfigurator{
		"ETC_Atlantis": difficultyChainConfigurations["ETC_Atlantis"],
		"ETC_Agharta":  difficultyChainConfigurations["ETC_Agharta"],
		"ETC_Phoenix":  difficultyChainConfigurations["ETC_Phoenix"],
		"ETC_Magneto":  difficultyChainConfigurations["ETC_Magneto"],
		"ETC_Mystique": difficultyChainConfigurations["ETC_Mystique"],
	}

	targetDir := filepath.Join(generatedBasedir, "BasicTests")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Write all the configs we'll use to JSON files for reference and provenance.
	if err := os.MkdirAll(filepath.Join(targetDir, "configs"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	for configName, config := range configs {
		cgConfig := &coregeth.CoreGethChainConfig{}
		err := confp.Convert(config, cgConfig)
		if err != nil {
			t.Fatal(err)
		}

		j, err := json.MarshalIndent(cgConfig, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(filepath.Join(targetDir, "configs", configName+".json"), j, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Establish test matrix values.
	//
	// Multiples of 8 offset by 1, maxing out at 121.
	timestampOffsets := func() []uint64 {
		var r []uint64
		for i := 1; i <= 121; i += 8 {
			r = append(r, uint64(i))
		}
		return r
	}()

	parentDifficulty := new(big.Int).Mul(vars.MinimumDifficulty, big.NewInt(10))

	for name, config := range configs {
		targetFileBaseName := fmt.Sprintf("difficulty%s.json", name)
		targetFilePath := filepath.Join(targetDir, targetFileBaseName)
		os.Truncate(targetFilePath, 0)
		for _, blockNumber := range []uint64{0, 1, 10_000_000} {
			for _, timestampOffset := range timestampOffsets {
				for _, uncle := range []bool{false, true} {
					uncleCount := 0
					if uncle {
						uncleCount = 1
					}
					uncleHash := types.EmptyUncleHash
					if uncle {
						uncleHash = types.CalcUncleHash([]*types.Header{{Number: common.Big1}})
					}
					// Establish test parameters.
					newTest := &DifficultyTest{
						ParentTimestamp:    blockNumber*13 + 13,
						ParentDifficulty:   parentDifficulty,
						ParentUncles:       uint64(uncleCount),
						CurrentBlockNumber: blockNumber,
						// CurrentTimestamp:  This gets filled later.
						// CurrentDifficulty: This gets filled later.
					}

					newTest.CurrentTimestamp = newTest.ParentTimestamp + timestampOffset

					// Fill the expected difficulty from the test params we've just established.
					newTest.CurrentDifficulty = ethash.CalcDifficulty(config, newTest.CurrentTimestamp, &types.Header{
						Difficulty: newTest.ParentDifficulty,
						Time:       newTest.ParentTimestamp,
						Number:     big.NewInt(int64(newTest.CurrentBlockNumber - 1)),
						UncleHash:  uncleHash,
					})

					m := mustReadTestFileJSON(t, targetFilePath)
					m[fmt.Sprintf("difficulty_n%d_t%d_u%t", blockNumber, timestampOffset, uncle)] = newTest
					mustWriteTestFileJSON(t, targetFilePath, m)
				}
			}
		}
	}
}

func mustReadTestFileJSON(tt *testing.T, filePath string) map[string]*DifficultyTest {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		tt.Log(err)
		return map[string]*DifficultyTest{}
	}
	var tests map[string]json.RawMessage
	if err := json.Unmarshal(b, &tests); err != nil {
		tt.Log(err)
		return map[string]*DifficultyTest{}
	}
	out := make(map[string]*DifficultyTest)
	for k, v := range tests {
		dt := new(DifficultyTest)
		if err := json.Unmarshal(v, dt); err != nil {
			tt.Log(err)
			continue
		}
		out[k] = dt
	}
	return out
}

func mustWriteTestFileJSON(tt *testing.T, filePath string, tests map[string]*DifficultyTest) {
	enc := make(map[string]json.RawMessage)
	for k, v := range tests {
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			tt.Fatal(err)
		}
		enc[k] = b
	}
	b, err := json.MarshalIndent(enc, "", "    ")
	if err != nil {
		tt.Fatal(err)
	}
	err = ioutil.WriteFile(filePath, b, os.ModePerm)
	if err != nil {
		tt.Fatal(err)
	}
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
