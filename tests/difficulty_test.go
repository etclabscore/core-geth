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
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

var (
	mainnetChainConfig = params.ChainConfig{
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
			uint64(0x0): new(big.Int).SetUint64(uint64(0x4563918244f40000)),
			uint64(4370000):                 new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(4370000): new(big.Int).SetUint64(uint64(0x2dc6c0)),
		},
	}
)

func TestDifficulty(t *testing.T) {
	t.Parallel()

	dt := new(testMatcher)
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

	dt.config("Ropsten", *params.TestnetChainConfig)
	dt.config("Morden", *params.TestnetChainConfig)
	dt.config("Frontier", params.ChainConfig{})

	dt.config("Homestead", params.ChainConfig{
		HomesteadBlock: big.NewInt(0),
	})

	dt.config("Byzantium", params.ChainConfig{
		ByzantiumBlock: big.NewInt(0),
		BlockRewardSchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x29a2241af62c0000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x2dc6c0)),
		},
	})

	dt.config("Frontier", *params.TestnetChainConfig)
	dt.config("MainNetwork", mainnetChainConfig)
	dt.config("CustomMainNetwork", mainnetChainConfig)
	dt.config("Constantinople", params.ChainConfig{
		ConstantinopleBlock: big.NewInt(0),
		BlockRewardSchedule: hexutil.Uint64BigMapEncodesHex{
			uint64(0): new(big.Int).SetUint64(uint64(0x1bc16d674ec80000)),
		},
		DifficultyBombDelaySchedule: hexutil.Uint64BigMapEncodesHex{
			//uint64(0): new(big.Int).SetUint64(uint64(0x2dc6c0)), // 3000000
			//uint64(0): new(big.Int).SetUint64(uint64(0x1e8480)), // 2000000
			0: big.NewInt(5000000), // Because the algo wants compounding or sum.
		},
	})
	dt.config("difficulty.json", mainnetChainConfig)

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		cfg := dt.findConfig(name)
		if test.ParentDifficulty.Cmp(params.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
