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
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/iancoleman/strcase"
)

// coregethSpecsDir is where core-geth-style configuration files for testing are stored.
var coregethSpecsDir = filepath.Join("..", "params", "coregeth.json.d")

// MapForkNameChainspecFileState is a dictionary pairing Fork names with respective
// file base names.
// These are used for StateTests, BlockchainTests, but not Difficulty tests.
// These files are expected to be found in coregethSpecsDir.
var MapForkNameChainspecFileState = map[string]string{
	"Frontier":             "frontier_test.json",
	"Homestead":            "homestead_test.json",
	"EIP150":               "eip150_test.json",
	"EIP158":               "eip161_test.json",
	"Byzantium":            "byzantium_test.json",
	"Constantinople":       "constantinople_test.json",
	"ConstantinopleFix":    "constantinople_fix_test.json",
	"EIP158ToByzantiumAt5": "eip_158_to_byzantium_at_5_test.json",
	"Istanbul":             "istanbul_test.json",
	"Berlin":               "berlin_test.json",
	"ETC_Atlantis":         "etc_atlantis_test.json",
	"ETC_Agharta":          "etc_agharta_test.json",
	"ETC_Phoenix":          "etc_phoenix_test.json",
	"ETC_Magneto":          "etc_magneto_test.json",
}

// mapForkNameChainspecFileDifficulty is a dictionary pairing fork names with respective
// file base name.
// These configurations are used exclusively for Difficulty tests.
// These files are expected to be found in coregethSpecsDir.
var mapForkNameChainspecFileDifficulty = map[string]string{
	"Ropsten":           "ropsten_difficulty_test.json",
	"Morden":            "morden_difficulty_test.json",
	"Frontier":          "frontier_difficulty_test.json",
	"Homestead":         "homestead_difficulty_test.json",
	"Byzantium":         "byzantium_difficulty_test.json",
	"MainNetwork":       "mainnetwork_difficulty_test.json",
	"CustomMainNetwork": "custom_mainnetwork_difficulty_test.json",
	"Constantinople":    "constantinople_difficulty_test.json",
	"difficulty.json":   "difficulty_json_difficulty_test.json",
	"ETC_Atlantis":      "classic_atlantis_difficulty_test.json",
	"ETC_Agharta":       "classic_agharta_difficulty_test.json",
	"EIP2384":           "eip2384_difficulty_test.json",
	"ETC_Phoenix":       "classic_phoenix_difficulty_test.json",
}

// readJSONFromFile is a utility function to read (unmarshaling) a value from a JSON file,
// which tries to return helpful errors if it is unable to, which can be useful for debugging.
// Additionally, it returns the SHA1 sum of the file if it does not error otherwise.
// This is floozy logic, but I don't really care right now.
func readJSONFromFile(name string, value interface{}) (sha1sum []byte, err error) {
	if fi, err := os.Open(name); os.IsNotExist(err) {
		return nil, err
	} else {
		fi.Close()
	}
	b, err := os.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	err = json.Unmarshal(b, value)
	if err != nil {
		if jsonError, ok := err.(*json.SyntaxError); ok {
			line, character, lcErr := lineAndCharacter(string(b), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "test failed with error: Cannot parse JSON schema due to a syntax error at line %d, character %d: %v\n", line, character, jsonError.Error())
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}
		if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
			line, character, lcErr := lineAndCharacter(string(b), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "test failed with error: The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "test failed with error: Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	bb := sha1.Sum(b)
	return bb[:], nil
}

func init() {
	if os.Getenv(CG_CHAINCONFIG_FEATURE_EQ_COREGETH_KEY) != "" {
		log.Println("converting to CoreGeth Chain Config data type.")

		for i, config := range Forks {
			mgc := &coregeth.CoreGethChainConfig{}
			if err := confp.Crush(mgc, config, true); ctypes.IsFatalUnsupportedErr(err) {
				panic(err)
			}
			Forks[i] = mgc
		}

		for k, v := range difficultyChainConfigurations {
			mgc := &coregeth.CoreGethChainConfig{}
			if err := confp.Crush(mgc, v, true); ctypes.IsFatalUnsupportedErr(err) {
				panic(err)
			}
			difficultyChainConfigurations[k] = mgc
		}
	} else if os.Getenv(CG_CHAINCONFIG_CHAINSPECS_COREGETH_KEY) != "" {
		// This logic reads Forks (used by [General]StateTests) and Difficulty configurations
		// from their respective coregeth.json.d/<file>.json files.
		// This implementation differs from that of this scope's predecessor CG_CHAINCONFIG_CHAINSPECS_OPENETHEREUM_KEY
		// because it only replaces Go values when it finds a corresponding configuration file
		// (it does not demand to replace all available configurations).
		// This avoids some unnecessary overhead for establishing configurations
		// that aren't really relevant, like Morden testnets.
		log.Println("Setting chain configurations from core-geth chainspecs")

		// newForks avoid write+iterate on Forks map.
		// All key:values in newForks will be written back to Forks.
		newForks := map[string]ctypes.ChainConfigurator{}
		for name := range Forks {
			gen := &genesisT.Genesis{
				Config: &coregeth.CoreGethChainConfig{},
			}
			specPath := filepath.Join(coregethSpecsDir, strcase.ToSnake(name)+"_test.json")

			sha1sum, err := readJSONFromFile(specPath, gen)
			if err != nil {
				log.Printf("Failed to read core-geth state config file for %s: %s", name, specPath)
				continue
			}
			chainspecRefsState[name] = chainspecRef{filepath.Base(specPath), sha1sum}
			newForks[name] = gen.Config
		}
		for name, conf := range newForks {
			Forks[name] = conf
		}

		for k, v := range mapForkNameChainspecFileDifficulty {
			conf := &coregeth.CoreGethChainConfig{}
			specPath := filepath.Join(coregethSpecsDir, v)
			sha1sum, err := readJSONFromFile(specPath, conf)
			if err != nil {
				log.Printf("Failed to read core-geth difficulty file for %s: %s", k, specPath)
				continue
			}
			if len(sha1sum) == 0 {
				panic("empty sha1 sum")
			}
			chainspecRefsDifficulty[k] = chainspecRef{filepath.Base(v), sha1sum}
			difficultyChainConfigurations[k] = conf
		}
	} else if os.Getenv(CG_CHAINCONFIG_CONSENSUS_EQ_CLIQUE) != "" {
		log.Println("converting Istanbul config to Clique consensus engine")

		for _, c := range Forks {
			if c.GetConsensusEngineType().IsEthash() {
				err := c.MustSetConsensusEngineType(ctypes.ConsensusEngineT_Clique)
				if err != nil {
					log.Fatal(err)
				}
				err = c.SetCliqueEpoch(30000)
				if err != nil {
					log.Fatal(err)
				}
				err = c.SetCliquePeriod(15)
				if err != nil {
					log.Fatal(err)
				}
			} else if c.GetConsensusEngineType().IsClique() {
				err := c.MustSetConsensusEngineType(ctypes.ConsensusEngineT_Ethash)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors/
func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1
	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}
	return line, character, nil
}
