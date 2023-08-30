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
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

// TestDifficultyGen2 generates JSON tests from scratch.
// The test case matrix can be deduced from the for-loop iterations.
// The cases are written to files PER CHAIN CONFIG, following the upstream convention,
// eg. tests/testdata_generated/BasicTests/difficultyETC_Agharta.json
func TestDifficultyGen2(t *testing.T) {
	if os.Getenv(CG_GENERATE_DIFFICULTY_TESTS_KEY) == "" {
		t.Skip()
	}

	configs := map[string]ctypes.ChainConfigurator{
		"ETC_Atlantis": Forks["ETC_Atlantis"],
		"ETC_Agharta":  Forks["ETC_Agharta"],
		"ETC_Phoenix":  Forks["ETC_Phoenix"],
		"ETC_Magneto":  Forks["ETC_Magneto"],
		"ETC_Mystique": Forks["ETC_Mystique"],
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
	blockTimeDefault := uint64(13)

	filledTests := map[string]*DifficultyTest{}
	for configName, config := range configs {
		testsName := fmt.Sprintf("difficulty%s", configName)
		targetFileBaseName := fmt.Sprintf("difficulty%s.json", configName)
		targetFilePath := filepath.Join(targetDir, targetFileBaseName)
		os.Truncate(targetFilePath, 0)
		heights := difficultyTestCaseHeights(config)
		for _, blockNumber := range heights {
			for _, timestampOffset := range timestampOffsets {
				for _, uncle := range []bool{false, true} {
					uncleCount := 0
					uncleHash := types.EmptyUncleHash
					if uncle {
						uncleCount = 1
						uncleHash = types.CalcUncleHash([]*types.Header{{Number: common.Big1}})
					}

					// Establish test parameters.
					newTest := &DifficultyTest{
						ParentTimestamp:    blockNumber*blockTimeDefault + blockTimeDefault,
						ParentDifficulty:   new(big.Int).Set(parentDifficulty),
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
				}
			}
		}
		writeDifficultyTestFileJSON(t, targetFilePath, filledTests, testsName, configName)
	}
}

// difficultyTestCaseHeights returns a list of block numbers for which we want to generate
// difficulty tests. The list is sorted in ascending order.
// Values are provided below (if possible), at (equal), and above the configuration's fork blocks.
func difficultyTestCaseHeights(config ctypes.ChainConfigurator) []uint64 {
	blockHeights := []uint64{}

	// Add the block-fork blocks.
	// We do not handle time-forks, since it is assumed that under time-based chain configuration contexts,
	// difficulty will be inoperative or otherwise disused.
	forks := confp.BlockForks(config)
	copy(blockHeights, forks)
	for _, forkBlock := range forks {
		if forkBlock > 0 {
			blockHeights = append(blockHeights, forkBlock-1)
		}
		blockHeights = append(blockHeights, forkBlock+1)
	}

	// Add random blocks, spaced out. We want to have them adequately spaced
	// so that, for example, we'd catch an accidental difficulty bomb explosion next year.
	for i := 0; i < 1000; i++ {
		blockHeights = append(blockHeights, uint64(10_000*1_000)) // 10_000 * 1_000 = 10_000_000
	}
	for i := 0; i < 100; i++ {
		blockHeights = append(blockHeights, uint64(rand.Int63n(100_000_000)))
	}

	//
	sort.Slice(blockHeights, func(i, j int) bool { return blockHeights[i] < blockHeights[j] })
	return blockHeights
}

func writeDifficultyTestFileJSON(t *testing.T, filePath string, tests map[string]*DifficultyTest, testName string, configName string) {
	head := strings.TrimSpace(build.RunGit("rev-parse", "HEAD"))
	filledWith := fmt.Sprintf("%s-%s-%s", params.VersionName, params.VersionWithMeta, head)

	enc := make(map[string]json.RawMessage)
	for k, v := range tests {
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		enc[k] = b
	}

	enc2 := map[string]interface{}{
		testName: map[string]interface{}{
			"_info": stInfo{
				Comment:            "Generated by etclabscore/core-geth tests/difficulty_mgen_test.go",
				FillingRPCServer:   filledWith,
				FillingToolVersion: filledWith,
				FilledWith:         "",
				LLLCVersion:        "",
				Source:             "",
				SourceHash:         "",
				Labels:             nil,
			},
			configName: enc,
		},
	}
	b, err := json.MarshalIndent(enc2, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filePath, b, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
