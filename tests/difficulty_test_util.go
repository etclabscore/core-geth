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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

//go:generate [gencodec -type DifficultyTest -field-override difficultyTestMarshaling -out gen_difficultytest.go]

var (
	mainnetChainConfig = params.ChainConfig{
		Ethash:         new(params.EthashConfig),
		ChainID:        big.NewInt(1),
		HomesteadBlock: big.NewInt(1150000),
		DAOForkBlock:   big.NewInt(1920000),
		DAOForkSupport: true,
		EIP150Block:    big.NewInt(2463000),
		EIP150Hash:     common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:    big.NewInt(2675000),
		EIP158Block:    big.NewInt(2675000),
		ByzantiumBlock: big.NewInt(4370000),
		BlockRewardSchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0x0):     new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			uint64(4370000): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(4370000): new(big.Int).SetUint64(uint64(0x2dc6c0)),
		},
	}
)

var difficultyChainConfiguations = map[string]params.ChainConfig{
	"Ropsten":  *params.TestnetChainConfig,
	"Morden":   *params.TestnetChainConfig,
	"Frontier": {
		Ethash:                      new(params.EthashConfig),
		BlockRewardSchedule:         hexutil.Uint64BigMapEncodesHex{},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{},
	},
	"Homestead": {
		Ethash:                      new(params.EthashConfig),
		HomesteadBlock:              big.NewInt(0),
		BlockRewardSchedule:         hexutil.Uint64BigMapEncodesHex{},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{},
	},
	"Byzantium": {
		Ethash:         new(params.EthashConfig),
		ByzantiumBlock: big.NewInt(0),
		BlockRewardSchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x2dc6c0)),
		},
	},
	"MainNetwork":       mainnetChainConfig,
	"CustomMainNetwork": mainnetChainConfig,
	"Constantinople": {
		Ethash:              new(params.EthashConfig),
		HomesteadBlock:      big.NewInt(0),
		EIP100FBlock:        big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		BlockRewardSchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			//uint64(0): new(big.Int).SetUint64(uint64(0x2dc6c0)), // 3000000
			//uint64(0): new(big.Int).SetUint64(uint64(0x1e8480)), // 2000000
			0: big.NewInt(5000000), // Because the algo wants compounding or sum.
		},
	},
	"difficulty.json": mainnetChainConfig,
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

func (test *DifficultyTest) Run(config *params.ChainConfig) error {
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
		return fmt.Errorf(`parent[time %v diff %v unclehash:%x]
child[time %v number %v]
diff %v != expected %v
chainspec %v
config %v`,
			test.ParentTimestamp, test.ParentDifficulty, test.UncleHash,
			test.CurrentTimestamp, test.CurrentBlockNumber, actual, exp,
			test.Chainspec, config,
		)
	}
	return nil

}
