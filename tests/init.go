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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/chainspec"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

func init() {
	for _, config := range Forks {
		if config.DifficultyBombDelaySchedule == nil {
			config.DifficultyBombDelaySchedule = hexutil.Uint64BigMapEncodesHex{}
		}
		if config.BlockRewardSchedule == nil {
			config.BlockRewardSchedule = hexutil.Uint64BigMapEncodesHex{
				uint64(0x0):      new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			}
		}
		if config.ByzantiumBlock != nil {
			config.DifficultyBombDelaySchedule[config.ByzantiumBlock.Uint64()] = big.NewInt(3000000)
			config.BlockRewardSchedule[config.ByzantiumBlock.Uint64()] = big.NewInt(3000000000000000000)
		}
		if config.ConstantinopleBlock != nil {
			config.DifficultyBombDelaySchedule[config.ConstantinopleBlock.Uint64()] = big.NewInt(2000000)
			config.BlockRewardSchedule[config.ConstantinopleBlock.Uint64()] = big.NewInt(2000000000000000000)

		}
	}

	if os.Getenv("MULTIGETH_TESTS_CHAINCONFIG_FEATURE_EQUIVALANCE") != "" {
		log.Println("Setting equivalent fork feature chain configurations")
		for _, config := range Forks {
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
		spew.Config.DisableMethods = true // Turn of ChainConfig Stringer method
		log.Println(spew.Sdump(Forks))
	}
	if os.Getenv("MULTIGETH_TESTS_CHAINCONFIG_PARITY_SPECS") != "" {
		paritySpecsDir := filepath.Join("..", "chainspecs")
		paritySpecPath := func(name string) string {
			return filepath.Join(paritySpecsDir, name)
		}
		matchups := map[string]string{
			"Frontier": paritySpecPath("frontier_test.json"),
			"Homestead": paritySpecPath("homestead_test.json"),
			"EIP150": paritySpecPath("eip150_test.json"),
			"EIP158": paritySpecPath("eip161_test.json"),
			"Byzantium": paritySpecPath("byzantium_test.json"),
			"Constantinople": paritySpecPath("constantinople_test.json"),
			"ConstantinopleFix": paritySpecPath("st_peters_test.json"),
			"EIP158ToByzantiumAt5": paritySpecPath("transition_test.json"),
			//"Istanbul": paritySpecPath("istanbul_test.json"),
		}
		for k, v := range matchups {
			spec := chainspec.ParityChainSpec{}
			b, err := ioutil.ReadFile(v)
			if err != nil {
				panic(fmt.Sprintf("%s/%s err: %s\n%s", k, v, err, b))
			}
			err = json.Unmarshal(b, &spec)
			if err != nil {
				panic(fmt.Sprintf("%s/%s err: %s\n%s", k, v, err, b))
			}
			genesis, err := chainspec.ParityConfigToMultiGethGenesis(&spec)
			if err != nil {
				panic(fmt.Sprintf("%s/%s err: %s\n%s", k, v, err, b))
			}
			Forks[k] = genesis.Config
		}
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
