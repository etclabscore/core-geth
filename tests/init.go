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
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func u64(val uint64) *uint64 { return &val }

// Forks table defines supported forks and their chain config.
var Forks = map[string]ctypes.ChainConfigurator{
	"Frontier": &goethereum.ChainConfig{
		Ethash:  new(ctypes.EthashConfig),
		ChainID: big.NewInt(1),
	},
	"Homestead": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
	},
	"EIP150": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
	},
	"EIP158": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
	},
	"Byzantium": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
		ByzantiumBlock: big.NewInt(0),
	},
	"ETC_Atlantis": &coregeth.CoreGethChainConfig{
		NetworkID:          1,
		Ethash:             new(ctypes.EthashConfig),
		ChainID:            big.NewInt(61),
		EIP2FBlock:         big.NewInt(0),
		EIP7FBlock:         big.NewInt(0),
		EIP150Block:        big.NewInt(0),
		EIP155Block:        big.NewInt(0),
		EIP160FBlock:       big.NewInt(0),
		EIP161FBlock:       big.NewInt(0),
		EIP170FBlock:       big.NewInt(0),
		EIP100FBlock:       big.NewInt(0),
		EIP140FBlock:       big.NewInt(0),
		EIP198FBlock:       big.NewInt(0),
		EIP211FBlock:       big.NewInt(0),
		EIP212FBlock:       big.NewInt(0),
		EIP213FBlock:       big.NewInt(0),
		EIP214FBlock:       big.NewInt(0),
		EIP658FBlock:       big.NewInt(0),
		EIP145FBlock:       nil,
		EIP1014FBlock:      nil,
		EIP1052FBlock:      nil,
		EIP1283FBlock:      nil,
		EIP2200FBlock:      nil, // RePetersburg
		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"Constantinople": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     nil,
	},
	"ConstantinopleFix": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
	},
	"ETC_Agharta": &coregeth.CoreGethChainConfig{
		NetworkID:          1,
		Ethash:             new(ctypes.EthashConfig),
		ChainID:            big.NewInt(61),
		EIP2FBlock:         big.NewInt(0),
		EIP7FBlock:         big.NewInt(0),
		EIP150Block:        big.NewInt(0),
		EIP155Block:        big.NewInt(0),
		EIP160FBlock:       big.NewInt(0),
		EIP161FBlock:       big.NewInt(0),
		EIP170FBlock:       big.NewInt(0),
		EIP100FBlock:       big.NewInt(0),
		EIP140FBlock:       big.NewInt(0),
		EIP198FBlock:       big.NewInt(0),
		EIP211FBlock:       big.NewInt(0),
		EIP212FBlock:       big.NewInt(0),
		EIP213FBlock:       big.NewInt(0),
		EIP214FBlock:       big.NewInt(0),
		EIP658FBlock:       big.NewInt(0),
		EIP145FBlock:       big.NewInt(0),
		EIP1014FBlock:      big.NewInt(0),
		EIP1052FBlock:      big.NewInt(0),
		EIP1283FBlock:      big.NewInt(0),
		PetersburgBlock:    big.NewInt(0),
		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"Istanbul": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
	},
	"ETC_Phoenix": &coregeth.CoreGethChainConfig{
		NetworkID:       1,
		Ethash:          new(ctypes.EthashConfig),
		ChainID:         big.NewInt(61),
		EIP2FBlock:      big.NewInt(0),
		EIP7FBlock:      big.NewInt(0),
		EIP150Block:     big.NewInt(0),
		EIP155Block:     big.NewInt(0),
		EIP160FBlock:    big.NewInt(0),
		EIP161FBlock:    big.NewInt(0),
		EIP170FBlock:    big.NewInt(0),
		EIP100FBlock:    big.NewInt(0),
		EIP140FBlock:    big.NewInt(0),
		EIP198FBlock:    big.NewInt(0),
		EIP211FBlock:    big.NewInt(0),
		EIP212FBlock:    big.NewInt(0),
		EIP213FBlock:    big.NewInt(0),
		EIP214FBlock:    big.NewInt(0),
		EIP658FBlock:    big.NewInt(0),
		EIP145FBlock:    big.NewInt(0),
		EIP1014FBlock:   big.NewInt(0),
		EIP1052FBlock:   big.NewInt(0),
		EIP1283FBlock:   big.NewInt(0),
		PetersburgBlock: big.NewInt(0),
		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"FrontierToHomesteadAt5": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(5),
	},
	"HomesteadToEIP150At5": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(5),
	},
	"HomesteadToDaoAt5": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		DAOForkBlock:   big.NewInt(5),
		DAOForkSupport: true,
	},
	"EIP158ToByzantiumAt5": &goethereum.ChainConfig{
		Ethash:         new(ctypes.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(0),
		EIP150Block:    big.NewInt(0),
		EIP155Block:    big.NewInt(0),
		EIP158Block:    big.NewInt(0),
		ByzantiumBlock: big.NewInt(5),
	},
	"ByzantiumToConstantinopleAt5": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(5),
	},
	"ByzantiumToConstantinopleFixAt5": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(5),
		PetersburgBlock:     big.NewInt(5),
	},
	"ConstantinopleFixToIstanbulAt5": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
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
	"MuirGlacier": &goethereum.ChainConfig{
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
		MuirGlacierBlock:    big.NewInt(0),
	},
	"Berlin": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
	},
	"BerlinToLondonAt5": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(5),
	},
	"London": &goethereum.ChainConfig{
		Ethash:              new(ctypes.EthashConfig),
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
	},
	"ArrowGlacier": &goethereum.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		ArrowGlacierBlock:   big.NewInt(0),
	},
	"ArrowGlacierToMergeAtDiffC0000": &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0xC0000),
	},
	"GrayGlacier": &goethereum.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		ArrowGlacierBlock:   big.NewInt(0),
		GrayGlacierBlock:    big.NewInt(0),
	},
	"Merge": &goethereum.ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(0),
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		MergeNetsplitBlock:            big.NewInt(0),
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
	},
	"ETC_Magneto": &coregeth.CoreGethChainConfig{
		NetworkID:       1,
		Ethash:          new(ctypes.EthashConfig),
		ChainID:         big.NewInt(61),
		EIP2FBlock:      big.NewInt(0),
		EIP7FBlock:      big.NewInt(0),
		EIP150Block:     big.NewInt(0),
		EIP155Block:     big.NewInt(0),
		EIP160FBlock:    big.NewInt(0),
		EIP161FBlock:    big.NewInt(0),
		EIP170FBlock:    big.NewInt(0),
		EIP100FBlock:    big.NewInt(0),
		EIP140FBlock:    big.NewInt(0),
		EIP198FBlock:    big.NewInt(0),
		EIP211FBlock:    big.NewInt(0),
		EIP212FBlock:    big.NewInt(0),
		EIP213FBlock:    big.NewInt(0),
		EIP214FBlock:    big.NewInt(0),
		EIP658FBlock:    big.NewInt(0),
		EIP145FBlock:    big.NewInt(0),
		EIP1014FBlock:   big.NewInt(0),
		EIP1052FBlock:   big.NewInt(0),
		EIP1283FBlock:   big.NewInt(0),
		PetersburgBlock: big.NewInt(0),
		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// Berlin
		EIP2565FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"ETC_Mystique": &coregeth.CoreGethChainConfig{
		NetworkID:       1,
		Ethash:          new(ctypes.EthashConfig),
		ChainID:         big.NewInt(61),
		EIP2FBlock:      big.NewInt(0),
		EIP7FBlock:      big.NewInt(0),
		EIP150Block:     big.NewInt(0),
		EIP155Block:     big.NewInt(0),
		EIP160FBlock:    big.NewInt(0),
		EIP161FBlock:    big.NewInt(0),
		EIP170FBlock:    big.NewInt(0),
		EIP100FBlock:    big.NewInt(0),
		EIP140FBlock:    big.NewInt(0),
		EIP198FBlock:    big.NewInt(0),
		EIP211FBlock:    big.NewInt(0),
		EIP212FBlock:    big.NewInt(0),
		EIP213FBlock:    big.NewInt(0),
		EIP214FBlock:    big.NewInt(0),
		EIP658FBlock:    big.NewInt(0),
		EIP145FBlock:    big.NewInt(0),
		EIP1014FBlock:   big.NewInt(0),
		EIP1052FBlock:   big.NewInt(0),
		EIP1283FBlock:   big.NewInt(0),
		PetersburgBlock: big.NewInt(0),
		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// Berlin
		EIP2565FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		// London
		/*
			https://github.com/ethereumclassic/ECIPs/blob/master/_specs/ecip-1104.md

			3529 (Alternative refund reduction) 	#22733 	Include
			3541 (Reject new contracts starting with the 0xEF byte) 	#22809 	Include
			1559 (Fee market change) 	#22837 #22896 	Omit
			3198 (BASEFEE opcode) 	#22837 	Omit
			3228 (bomb delay) 	#22840 and #22870 	Omit
		*/
		EIP3529FBlock: big.NewInt(0),
		EIP3541FBlock: big.NewInt(0),
		EIP1559FBlock: nil,
		EIP3198FBlock: nil,
		EIP3554FBlock: nil,

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"Shanghai": &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		ShanghaiTime:            u64(0),
	},
	"MergeToShanghaiAtTime15k": &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		ShanghaiTime:            u64(15_000),
	},
	"ETC_Spiral": &coregeth.CoreGethChainConfig{
		NetworkID:       1,
		Ethash:          new(ctypes.EthashConfig),
		ChainID:         big.NewInt(61),
		EIP2FBlock:      big.NewInt(0),
		EIP7FBlock:      big.NewInt(0),
		EIP150Block:     big.NewInt(0),
		EIP155Block:     big.NewInt(0),
		EIP160FBlock:    big.NewInt(0),
		EIP161FBlock:    big.NewInt(0),
		EIP170FBlock:    big.NewInt(0),
		EIP100FBlock:    big.NewInt(0),
		EIP140FBlock:    big.NewInt(0),
		EIP198FBlock:    big.NewInt(0),
		EIP211FBlock:    big.NewInt(0),
		EIP212FBlock:    big.NewInt(0),
		EIP213FBlock:    big.NewInt(0),
		EIP214FBlock:    big.NewInt(0),
		EIP658FBlock:    big.NewInt(0),
		EIP145FBlock:    big.NewInt(0),
		EIP1014FBlock:   big.NewInt(0),
		EIP1052FBlock:   big.NewInt(0),
		EIP1283FBlock:   big.NewInt(0),
		PetersburgBlock: big.NewInt(0),
		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // RePetersburg (=~ re-1283)

		// Berlin
		EIP2565FBlock: big.NewInt(0),
		EIP2929FBlock: big.NewInt(0),
		EIP2718FBlock: big.NewInt(0),
		EIP2930FBlock: big.NewInt(0),

		// London
		/*
			https://github.com/ethereumclassic/ECIPs/blob/master/_specs/ecip-1104.md

			3529 (Alternative refund reduction) 	#22733 	Include
			3541 (Reject new contracts starting with the 0xEF byte) 	#22809 	Include
			1559 (Fee market change) 	#22837 #22896 	Omit
			3198 (BASEFEE opcode) 	#22837 	Omit
			3228 (bomb delay) 	#22840 and #22870 	Omit
		*/
		EIP3529FBlock: big.NewInt(0),
		EIP3541FBlock: big.NewInt(0),
		EIP1559FBlock: nil,
		EIP3198FBlock: nil,
		EIP3554FBlock: nil,

		// Shanghai == Spiral
		EIP4399FBlock: nil,           // Supplant DIFFICULTY with PREVRANDAO. ETC does not spec 4399 because it's still PoW, and 4399 is only applicable for the PoS system.
		EIP3651FBlock: big.NewInt(0), // Warm COINBASE (gas reprice)
		EIP3855FBlock: big.NewInt(0), // PUSH0 instruction
		EIP3860FBlock: big.NewInt(0), // Limit and meter initcode
		EIP4895FBlock: nil,           // Beacon chain push withdrawals as operations
		EIP6049FBlock: big.NewInt(0), // Deprecate SELFDESTRUCT (noop)

		// ETC specifics
		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(5000000), // FIXME(meows) maybe
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
	},
	"Cancun": &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		ShanghaiTime:            u64(0),
		CancunTime:              u64(0),
	},
	"ShanghaiToCancunAtTime15k": &goethereum.ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(0),
		ShanghaiTime:            u64(0),
		CancunTime:              u64(15_000),
	},
}

// AvailableForks returns the set of defined fork names
func AvailableForks() []string {
	var availableForks []string
	for k := range Forks {
		availableForks = append(availableForks, k)
	}
	sort.Strings(availableForks)
	return availableForks
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
