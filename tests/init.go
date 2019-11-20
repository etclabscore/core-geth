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
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/chainspec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/internal/build"
	"github.com/ethereum/go-ethereum/params"
)

var (
	MG_GENERATE_STATE_TESTS_KEY      = "MULTIGETH_TESTS_GENERATE_STATE_TESTS"
	MG_GENERATE_DIFFICULTY_TESTS_KEY = "MULTIGETH_TESTS_GENERATE_DIFFICULTY_TESTS"
	MG_CHAINCONFIG_FEATURE_EQ_KEY    = "MULTIGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALANCE"
	MG_CHAINCONFIG_CHAINSPEC_KEY     = "MULTIGETH_TESTS_CHAINCONFIG_PARITY_SPECS"
)

// writeStateTestsReferencePairs defines reference pairs for use when writing tests.
// The reference (key) is used to define the environment and parameters, while the
// output from these tests run against the <value> state is actually written.
var writeStateTestsReferencePairs = map[string]string{
	"Byzantium":         "ETC_Atlantis",
	"ConstantinopleFix": "ETC_Agharta",
}

var writeDifficultyTestsReferencePairs = map[string]string{
	"Byzantium":      "ETC_Atlantis",
	"Constantinople": "ETC_Agharta",
}

type chainspecRefsT map[string]chainspecRef

var chainspecRefsState = chainspecRefsT{}
var chainspecRefsDifficulty = chainspecRefsT{}

type chainspecRef struct {
	Filename string `json:"filename"`
	Sha1Sum  []byte `json:"sha1sum"`
}

func (c chainspecRef) String() string {
	return fmt.Sprintf("file: %s, file.sha1sum: %x", c.Filename, c.Sha1Sum)
}

func (c *chainspecRef) UnmarshalJSON(input []byte) error {
	type xT struct {
		F string `json:"filename"`
		S string `json:"sha1sum"`
	}
	var x = xT{}
	err := json.Unmarshal(input, &x)
	if err != nil {
		return err
	}
	c.Filename = x.F
	c.Sha1Sum = common.Hex2Bytes(x.S)
	return nil
}

func (c chainspecRef) MarshalJSON() ([]byte, error) {
	var x = struct {
		F string `json:"filename"`
		S string `json:"sha1sum"`
	}{
		F: c.Filename,
		S: common.Bytes2Hex(c.Sha1Sum[:]),
	}

	return json.MarshalIndent(x, "", "    ")
}

// submoduleParentRef captures the current git status of the tests submodule.
// This is used for reference when writing tests.
var submoduleParentRef = func() string {
	subModOut := build.RunGit("submodule", "status")
	subModOut = strings.ReplaceAll(strings.TrimSpace(subModOut), " ", "_")
	return subModOut
}()

var paritySpecsDir = filepath.Join("..", "chainspecs")

func paritySpecPath(name string) string {
	p := filepath.Join(paritySpecsDir, name)
	if fi, err := os.Open(p); err == nil {
		fi.Close()
		return p
	} else if os.IsNotExist(err) {
		p = filepath.Join("..", paritySpecsDir, name)
	}
	return p
}

var mapForkNameChainspecFileState = map[string]string{
	"Frontier":             "frontier_test.json",
	"Homestead":            "homestead_test.json",
	"EIP150":               "eip150_test.json",
	"EIP158":               "eip161_test.json",
	"Byzantium":            "byzantium_test.json",
	"Constantinople":       "constantinople_test.json",
	"ConstantinopleFix":    "st_peters_test.json",
	"EIP158ToByzantiumAt5": "transition_test.json",
	"Istanbul":             "istanbul_test.json",
	"ETC_Atlantis":         "classic_atlantis_test.json",
	"ETC_Agharta":          "classic_agharta_test.json",
}

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
}

func readConfigFromSpecFile(name string) (genesis *core.Genesis, sha1sum []byte, err error) {
	spec := chainspec.ParityChainSpec{}
	if fi, err := os.Open(name); os.IsNotExist(err) {
		return nil, nil, err
	} else {
		fi.Close()
	}
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	err = json.Unmarshal(b, &spec)
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
	genesis, err = chainspec.ParityConfigToMultiGethGenesis(&spec)
	if err != nil {
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	bb := sha1.Sum(b)
	return genesis, bb[:], nil
}

func init() {
	for _, config := range Forks {
		config := config
		convertMetaForkBlocksDifficultyAndRewardSchedules(config)
	}

	if os.Getenv(MG_CHAINCONFIG_FEATURE_EQ_KEY) != "" {
		log.Println("Setting equivalent fork feature chain configurations")

		for _, config := range Forks {
			config := config
			convertMetaForkBlocksToFeatures(config)
		}

		for k, v := range difficultyChainConfiguations {
			v := v
			convertMetaForkBlocksToFeatures(&v)
			difficultyChainConfiguations[k] = v
		}

	} else if os.Getenv(MG_CHAINCONFIG_CHAINSPEC_KEY) != "" {
		log.Println("Setting chain configurations from Parity chainspecs")

		for k, v := range mapForkNameChainspecFileState {
			genesis, sha1sum, err := readConfigFromSpecFile(paritySpecPath(v))
			if os.IsNotExist(err) {
				wd, wde := os.Getwd()
				if wde != nil {
					panic(wde)
				}
				panic(fmt.Sprintf("failed to find chainspec, wd: %s", wd))
			} else if err != nil {
				panic(err)
			}
			chainspecRefsState[k] = chainspecRef{filepath.Base(v), sha1sum}
			Forks[k] = genesis.Config
		}

		for k, v := range mapForkNameChainspecFileDifficulty {
			genesis, sha1sum, err := readConfigFromSpecFile(paritySpecPath(v))
			if os.IsNotExist(err) && os.Getenv(MG_GENERATE_DIFFICULTY_TESTS_KEY) != "" {
				log.Println("Will generate chainspec file for", k, v)
			} else if len(sha1sum) == 0 {
				panic("zero sum game")
			} else {
				chainspecRefsDifficulty[k] = chainspecRef{filepath.Base(v), sha1sum}
				difficultyChainConfiguations[k] = *genesis.Config
			}
		}
	}
}

func convertMetaForkBlocksDifficultyAndRewardSchedules(config *params.ChainConfig) {
	if config.BlockRewardSchedule == nil {
		config.BlockRewardSchedule = hexutil.Uint64BigMapEncodesHex{
			uint64(0x0): new(big.Int).SetUint64(uint64(0x4563918244f40000)),
		}
	}
	if config.DifficultyBombDelaySchedule == nil {
		config.DifficultyBombDelaySchedule = hexutil.Uint64BigMapEncodesHex{}
		if config.ByzantiumBlock != nil {
			config.DifficultyBombDelaySchedule[config.ByzantiumBlock.Uint64()] = big.NewInt(3000000)
			config.BlockRewardSchedule[config.ByzantiumBlock.Uint64()] = big.NewInt(3000000000000000000)
		}
		if config.ConstantinopleBlock != nil {
			config.DifficultyBombDelaySchedule[config.ConstantinopleBlock.Uint64()] = big.NewInt(2000000)
			config.BlockRewardSchedule[config.ConstantinopleBlock.Uint64()] = big.NewInt(2000000000000000000)
		}
	}
}

func convertMetaForkBlocksToFeatures(config *params.ChainConfig) {
	if config.HomesteadBlock != nil {
		config.EIP2FBlock = config.HomesteadBlock
		config.EIP7FBlock = config.HomesteadBlock
		config.HomesteadBlock = nil
	}
	if config.EIP158Block != nil {
		config.EIP160FBlock = config.EIP158Block
		config.EIP161FBlock = config.EIP158Block
		config.EIP170FBlock = config.EIP158Block
		config.EIP158Block = nil
	}
	if config.ByzantiumBlock != nil {
		// Difficulty adjustment to target mean block time including uncles
		// https://github.com/ethereum/EIPs/issues/100
		config.EIP100FBlock = config.ByzantiumBlock
		// Opcode REVERT
		// https://eips.ethereum.org/EIPS/eip-140
		config.EIP140FBlock = config.ByzantiumBlock
		// Precompiled contract for bigint_modexp
		// https://github.com/ethereum/EIPs/issues/198
		config.EIP198FBlock = config.ByzantiumBlock
		// Opcodes RETURNDATACOPY, RETURNDATASIZE
		// https://github.com/ethereum/EIPs/issues/211
		config.EIP211FBlock = config.ByzantiumBlock
		// Precompiled contract for pairing check
		// https://github.com/ethereum/EIPs/issues/212
		config.EIP212FBlock = config.ByzantiumBlock
		// Precompiled contracts for addition and scalar multiplication on the elliptic curve alt_bn128
		// https://github.com/ethereum/EIPs/issues/213
		config.EIP213FBlock = config.ByzantiumBlock
		// Opcode STATICCALL
		// https://github.com/ethereum/EIPs/issues/214
		config.EIP214FBlock = config.ByzantiumBlock
		// Metropolis diff bomb delay and reducing block reward
		// https://github.com/ethereum/EIPs/issues/649
		// note that this is closely related to EIP100.
		// In fact, EIP100 is bundled in
		config.EIP649FBlock = config.ByzantiumBlock
		// Transaction receipt status
		// https://github.com/ethereum/EIPs/issues/658
		config.EIP658FBlock = config.ByzantiumBlock
		// NOT CONFIGURABLE: prevent overwriting contracts
		// https://github.com/ethereum/EIPs/issues/684
		// EIP684FBlock *big.Int `json:"eip684BFlock,omitempty"`

		config.ByzantiumBlock = nil
	}
	if config.ConstantinopleBlock != nil {
		// Opcodes SHR, SHL, SAR
		// https://eips.ethereum.org/EIPS/eip-145
		config.EIP145FBlock = config.ConstantinopleBlock
		// Opcode CREATE2
		// https://eips.ethereum.org/EIPS/eip-1014
		config.EIP1014FBlock = config.ConstantinopleBlock
		// Opcode EXTCODEHASH
		// https://eips.ethereum.org/EIPS/eip-1052
		config.EIP1052FBlock = config.ConstantinopleBlock
		// Constantinople difficulty bomb delay and block reward adjustment
		// https://eips.ethereum.org/EIPS/eip-1234
		config.EIP1234FBlock = config.ConstantinopleBlock
		// Net gas metering
		// https://eips.ethereum.org/EIPS/eip-1283
		config.EIP1283FBlock = config.ConstantinopleBlock

		config.ConstantinopleBlock = nil
	}
	if config.IstanbulBlock != nil {
		config.EIP152FBlock = config.IstanbulBlock
		config.EIP1108FBlock = config.IstanbulBlock
		config.EIP1344FBlock = config.IstanbulBlock
		config.EIP1884FBlock = config.IstanbulBlock
		config.EIP2028FBlock = config.IstanbulBlock
		config.EIP2200FBlock = config.IstanbulBlock
		config.IstanbulBlock = nil
	}
}

// Forks table defines supported forks and their chain config.
var Forks = map[string]*params.ChainConfig{
	"Frontier": {
		ChainID: big.NewInt(1),
	},
	"Homestead": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
	},
	"EIP150": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
	},
	"EIP158": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
	},
	"Byzantium": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
		DAOForkBlock:   big.NewInt(0),
		ByzantiumBlock: big.NewInt(0),
	},
	"ETC_Atlantis": {
		ChainID:           big.NewInt(1),
		HomesteadBlock:    big.NewInt(0),
		EIP150Block:       big.NewInt(0),
		EIP155Block:       big.NewInt(0),
		EIP158Block:       big.NewInt(0),
		DAOForkBlock:      big.NewInt(0),
		ByzantiumBlock:    big.NewInt(0),
		ECIP1017EraRounds: big.NewInt(5000000),
		ECIP1017FBlock:    big.NewInt(0),
		DisposalBlock:     big.NewInt(0),
	},
	"Constantinople": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(10000000),
	},
	"ConstantinopleFix": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
	},
	"ETC_Agharta": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		ECIP1017EraRounds:   big.NewInt(5000000),
		ECIP1017FBlock:      big.NewInt(0),
		DisposalBlock:       big.NewInt(0),
	},
	"Istanbul": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
	},
	"FrontierToHomesteadAt5": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(5),
	},
	"HomesteadToEIP150At5": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(5),
	},
	"HomesteadToDaoAt5": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		DAOForkBlock:   big.NewInt(5),
		DAOForkSupport: true,
	},
	"EIP158ToByzantiumAt5": {
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
		ByzantiumBlock: big.NewInt(5),
	},
	"ByzantiumToConstantinopleAt5": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(5),
	},
	"ByzantiumToConstantinopleFixAt5": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(5),
		PetersburgBlock:     big.NewInt(5),
	},
	"ConstantinopleFixToIstanbulAt5": {
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(5),
	},
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
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
