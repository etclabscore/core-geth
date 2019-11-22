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
package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
)

var (
	// Genesis hashes to enforce below configs on.
	MusicoinGenesisHash = common.HexToHash("0x4eba28a4ce8dc0701f94c936a223a8429129b38ca9974ec0e92bf9234ac952e9")

	// MusicoinChainConfig is the chain parameters to run a node on the main network.
	MusicoinChainConfig = &paramtypes.ChainConfig{
		NetworkID:      776959,
		MCIP0Block:     big.NewInt(0),
		ChainID:        big.NewInt(7762959),
		HomesteadBlock: big.NewInt(1150000),
		DAOForkBlock:   big.NewInt(36028797018963967),
		DAOForkSupport: false,
		EIP150Block:    big.NewInt(2222222),
		EIP150Hash:     common.HexToHash("0x"),
		EIP155Block:    big.NewInt(2222222),
		EIP158Block:    big.NewInt(2222222),
		ByzantiumBlock: big.NewInt(2222222),
		MCIP3Block:     big.NewInt(1200001),
		MCIP8Block:     big.NewInt(5200001),
		Ethash:         new(paramtypes.EthashConfig),
		BlockRewardSchedule: common2.Uint64BigMapEncodesHex{
			uint64(0):       new(big.Int).Mul(big.NewInt(314), big.NewInt(1e+18)),
			uint64(1200001): new(big.Int).Mul(big.NewInt(250), big.NewInt(1e+18)),
			uint64(5200001): new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18)),
		},
		DifficultyBombDelaySchedule: common2.Uint64BigMapEncodesHex{
			uint64(2222222): new(big.Int).SetUint64(uint64(0x2dc6c0)),
		},
	}

	MusicoinTimeCapsuleBlock  = int64(4200000)
	MusicoinTimeCapsuleLength = uint64(50) // Threshold of blocks that can be delayed and the value is in Blocks

	Mcip0BlockReward       = new(big.Int).Mul(big.NewInt(314), big.NewInt(1e+18)) // In musicoin code as 'FrontierBlockReward'
	Mcip3BlockReward       = new(big.Int).Mul(big.NewInt(250), big.NewInt(1e+18))
	Mcip8BlockReward       = new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18))
	MusicoinUbiBlockReward = new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18))
	MusicoinDevBlockReward = new(big.Int).Mul(big.NewInt(14), big.NewInt(1e+18))
)

