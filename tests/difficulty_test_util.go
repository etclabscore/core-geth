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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

//go:generate [gencodec -type DifficultyTest -field-override difficultyTestMarshaling -out gen_difficultytest.go]

var (
	mainnetChainConfig = &goethereum.ChainConfig{
		Ethash:         new(goethereum.EthashConfig),
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

var difficultyChainConfigurations = map[string]common2.ChainConfigurator{
	"Ropsten":  params.TestnetChainConfig,
	"Morden":   params.TestnetChainConfig,
	"Frontier": &goethereum.ChainConfig{},
	"Homestead": &goethereum.ChainConfig{
		Ethash:         new(goethereum.EthashConfig),
		HomesteadBlock: big.NewInt(0),
	},
	"Byzantium": &goethereum.ChainConfig{
		Ethash:         new(goethereum.EthashConfig),
		ByzantiumBlock: big.NewInt(0),
	},
	"MainNetwork":       mainnetChainConfig,
	"CustomMainNetwork": mainnetChainConfig,
	"Constantinople": &goethereum.ChainConfig{
		Ethash:              new(goethereum.EthashConfig),
		HomesteadBlock:      big.NewInt(0),
		ByzantiumBlock: big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
	},
	"difficulty.json": mainnetChainConfig,
	"ETC_Atlantis": &paramtypes.MultiGethChainConfig{
		Ethash:        new(goethereum.EthashConfig),
		EIP100FBlock:  big.NewInt(0),
		EIP140FBlock:  big.NewInt(0),
		EIP198FBlock:  big.NewInt(0),
		EIP211FBlock:  big.NewInt(0),
		EIP212FBlock:  big.NewInt(0),
		EIP213FBlock:  big.NewInt(0),
		EIP214FBlock:  big.NewInt(0),
		EIP658FBlock:  big.NewInt(0),
		DisposalBlock: big.NewInt(0),
	},
	"ETC_Agharta": &paramtypes.MultiGethChainConfig{
		Ethash:        new(goethereum.EthashConfig),
		EIP100FBlock:  big.NewInt(0),
		EIP140FBlock:  big.NewInt(0),
		EIP198FBlock:  big.NewInt(0),
		EIP211FBlock:  big.NewInt(0),
		EIP212FBlock:  big.NewInt(0),
		EIP213FBlock:  big.NewInt(0),
		EIP214FBlock:  big.NewInt(0),
		EIP658FBlock:  big.NewInt(0),
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		EIP1283FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0), // Petersburg
		DisposalBlock: big.NewInt(0),
	},
}

type DifficultyTest struct {
	ParentTimestamp    uint64       `json:"parentTimestamp"`
	ParentDifficulty   *big.Int     `json:"parentDifficulty"`
	UncleHash          common.Hash  `json:"parentUncles"`
	CurrentTimestamp   uint64       `json:"currentTimestamp"`
	CurrentBlockNumber uint64       `json:"currentBlockNumber"`
	CurrentDifficulty  *big.Int     `json:"currentDifficulty"`
	Chainspec          chainspecRef `json:"chainspec"`
	Name               string       `json:"name"`
}

type difficultyTestMarshaling struct {
	ParentTimestamp    math.HexOrDecimal64
	ParentDifficulty   *math.HexOrDecimal256
	CurrentTimestamp   math.HexOrDecimal64
	CurrentDifficulty  *math.HexOrDecimal256
	UncleHash          common.Hash
	CurrentBlockNumber math.HexOrDecimal64
	Chainspec          chainspecRef `json:"chainspec"`
	Name               string
}

func (t *DifficultyTest) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (test *DifficultyTest) Run(config common2.ChainConfigurator) error {
	parentNumber := big.NewInt(int64(test.CurrentBlockNumber - 1))
	parent := &types.Header{
		Difficulty: test.ParentDifficulty,
		Time:       test.ParentTimestamp,
		Number:     parentNumber,
		UncleHash:  test.UncleHash,
	}

	actual := ethash.CalcDifficulty(config, test.CurrentTimestamp, parent)
	exp := test.CurrentDifficulty

	if actual.Cmp(exp) != 0 {
		return fmt.Errorf(`%s got: %v, want: %v
test: %v
config: %v`, test.Name, actual, exp, test, config)
	}
	return nil

}
