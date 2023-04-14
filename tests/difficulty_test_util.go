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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

//go:generate go run github.com/fjl/gencodec -type DifficultyTest -field-override difficultyTestMarshaling -out gen_difficultytest.go

type DifficultyTest struct {
	ParentTimestamp    uint64   `json:"parentTimestamp"`
	ParentDifficulty   *big.Int `json:"parentDifficulty"`
	ParentUncles       uint64   `json:"parentUncles"`
	CurrentTimestamp   uint64   `json:"currentTimestamp"`
	CurrentBlockNumber uint64   `json:"currentBlockNumber"`
	CurrentDifficulty  *big.Int `json:"currentDifficulty"`
}

type difficultyTestMarshaling struct {
	ParentTimestamp    math.HexOrDecimal64
	ParentDifficulty   *math.HexOrDecimal256
	CurrentTimestamp   math.HexOrDecimal64
	CurrentDifficulty  *math.HexOrDecimal256
	ParentUncles       math.HexOrDecimal64
	CurrentBlockNumber math.HexOrDecimal64
}

var uncleHashNonEmpty = types.CalcUncleHash([]*types.Header{{Number: common.Big1}})

func (test *DifficultyTest) Run(config ctypes.ChainConfigurator) error {
	parentNumber := big.NewInt(int64(test.CurrentBlockNumber - 1))
	uncleHash := types.EmptyUncleHash
	if test.ParentUncles != 0 {
		uncleHash = uncleHashNonEmpty
	}
	parent := &types.Header{
		Difficulty: test.ParentDifficulty,
		Time:       test.ParentTimestamp,
		Number:     parentNumber,
		UncleHash:  uncleHash,
	}

	actual := ethash.CalcDifficulty(config, test.CurrentTimestamp, parent)
	exp := test.CurrentDifficulty

	if actual.Cmp(exp) != 0 {
		return fmt.Errorf("parent[time %v diff %v unclehash:%x] child[time %v number %v] diff %v != expected %v",
			test.ParentTimestamp, test.ParentDifficulty, test.ParentUncles,
			test.CurrentTimestamp, test.CurrentBlockNumber, actual, exp)
	}
	return nil
}

var (
	mainnetChainConfig = &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(1150000),
		DAOForkBlock:   big.NewInt(1920000),
		DAOForkSupport: true,
		EIP150Block:    big.NewInt(2463000),
		EIP150Hash:     common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:    big.NewInt(2675000),
		EIP158Block:    big.NewInt(2675000),
		ByzantiumBlock: big.NewInt(4370000),
	}
)

var difficultyChainConfigurations = map[string]ctypes.ChainConfigurator{
	"Frontier": &goethereum.ChainConfig{},
	"Homestead": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		HomesteadBlock: big.NewInt(0),
	},
	"Byzantium": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ByzantiumBlock: big.NewInt(0),
	},
	"GrayGlacier": &goethereum.ChainConfig{
		Ethash:           new(ctypes.EthashConfig),
		GrayGlacierBlock: big.NewInt(0),
	},
	"MainNetwork":       mainnetChainConfig,
	"CustomMainNetwork": mainnetChainConfig,
	"Constantinople": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		HomesteadBlock:      big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
	},
	"difficulty.json": mainnetChainConfig,
	"ETC_Atlantis": &coregeth.CoreGethChainConfig{
		Ethash:     new(ctypes.EthashConfig),
		NetworkID:  1,
		ChainID:    big.NewInt(1),
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// // Constantinople eq, aka Agharta
		// EIP145FBlock:  big.NewInt(0),
		// EIP1014FBlock: big.NewInt(0),
		// EIP1052FBlock: big.NewInt(0),
		// // EIP1283FBlock:   big.NewInt(0),
		// // PetersburgBlock: big.NewInt(0),
		//
		// // Istanbul eq, aka Phoenix
		// // ECIP-1088
		// EIP152FBlock:  big.NewInt(0),
		// EIP1108FBlock: big.NewInt(0),
		// EIP1344FBlock: big.NewInt(0),
		// EIP1884FBlock: big.NewInt(0),
		// EIP2028FBlock: big.NewInt(0),
		// EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)
		//
		// // Berlin eq, aka Magneto
		// EIP2565FBlock: big.NewInt(0),
		// EIP2718FBlock: big.NewInt(0),
		// EIP2929FBlock: big.NewInt(0),
		// EIP2930FBlock: big.NewInt(0),
		//
		// ECIP1099FBlock: big.NewInt(0), // Etchash (DAG size limit)
		//
		// // London (partially), aka Mystique
		// EIP3529FBlock: big.NewInt(0),
		// EIP3541FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),
		ECBP1100FBlock:     big.NewInt(0), // ETA 09 Oct 2020
	},
	"ETC_Agharta": &coregeth.CoreGethChainConfig{
		Ethash:     new(ctypes.EthashConfig),
		NetworkID:  1,
		ChainID:    big.NewInt(1),
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		// EIP1283FBlock:   big.NewInt(0),
		// PetersburgBlock: big.NewInt(0),

		// // Istanbul eq, aka Phoenix
		// // ECIP-1088
		// EIP152FBlock:  big.NewInt(0),
		// EIP1108FBlock: big.NewInt(0),
		// EIP1344FBlock: big.NewInt(0),
		// EIP1884FBlock: big.NewInt(0),
		// EIP2028FBlock: big.NewInt(0),
		// EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)
		//
		// // Berlin eq, aka Magneto
		// EIP2565FBlock: big.NewInt(0),
		// EIP2718FBlock: big.NewInt(0),
		// EIP2929FBlock: big.NewInt(0),
		// EIP2930FBlock: big.NewInt(0),
		//
		// ECIP1099FBlock: big.NewInt(0), // Etchash (DAG size limit)
		//
		// // London (partially), aka Mystique
		// EIP3529FBlock: big.NewInt(0),
		// EIP3541FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),
		ECBP1100FBlock:     big.NewInt(0), // ETA 09 Oct 2020
	},
	"EIP2384": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		HomesteadBlock:      big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
	},
	"ETC_Phoenix": &coregeth.CoreGethChainConfig{
		Ethash:     new(ctypes.EthashConfig),
		NetworkID:  1,
		ChainID:    big.NewInt(1),
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		// EIP1283FBlock:   big.NewInt(0),
		// PetersburgBlock: big.NewInt(0),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// // Berlin eq, aka Magneto
		// EIP2565FBlock: big.NewInt(0),
		// EIP2718FBlock: big.NewInt(0),
		// EIP2929FBlock: big.NewInt(0),
		// EIP2930FBlock: big.NewInt(0),
		//
		// ECIP1099FBlock: big.NewInt(0), // Etchash (DAG size limit)
		//
		// // London (partially), aka Mystique
		// EIP3529FBlock: big.NewInt(0),
		// EIP3541FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),
		ECBP1100FBlock:     big.NewInt(0), // ETA 09 Oct 2020
	},
	"ETC_Magneto": &coregeth.CoreGethChainConfig{
		Ethash:     new(ctypes.EthashConfig),
		NetworkID:  1,
		ChainID:    big.NewInt(1),
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		// EIP1283FBlock:   big.NewInt(0),
		// PetersburgBlock: big.NewInt(0),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		// ECIP1099FBlock: big.NewInt(0), // Etchash (DAG size limit)
		//
		// // London (partially), aka Mystique
		// EIP3529FBlock: big.NewInt(0),
		// EIP3541FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),
		ECBP1100FBlock:     big.NewInt(0), // ETA 09 Oct 2020
	},
	"ETC_Mystique": &coregeth.CoreGethChainConfig{
		Ethash:     new(ctypes.EthashConfig),
		NetworkID:  1,
		ChainID:    big.NewInt(1),
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(0),
		EIP160FBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		// EIP1283FBlock:   big.NewInt(0),
		// PetersburgBlock: big.NewInt(0),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		ECIP1099FBlock: big.NewInt(0), // Etchash (DAG size limit)

		// London (partially), aka Mystique
		EIP3529FBlock: big.NewInt(0),
		EIP3541FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),
		ECBP1100FBlock:     big.NewInt(0), // ETA 09 Oct 2020
	},
}
