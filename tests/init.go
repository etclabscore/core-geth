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

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

func newUint64(n uint64) *uint64 {
	return &n
}

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
	"ETC_Atlantis": &paramtypes.MultiGethChainConfig{
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
	"ETC_Agharta": &paramtypes.MultiGethChainConfig{
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
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
