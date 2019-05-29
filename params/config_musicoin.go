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
)

var (
	// Genesis hashes to enforce below configs on.
	MusicoinGenesisHash = common.HexToHash("0x4eba28a4ce8dc0701f94c936a223a8429129b38ca9974ec0e92bf9234ac952e9")

	// MusicoinChainConfig is the chain parameters to run a node on the main network.
	MusicoinChainConfig = &ChainConfig{
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
		Ethash:         new(EthashConfig),
	}

	MusicoinTimeCapsuleBlock  = int64(4200000)
	MusicoinTimeCapsuleLength = uint64(50) // Threshold of blocks that can be delayed and the value is in Blocks

	Mcip0BlockReward       = new(big.Int).Mul(big.NewInt(314), big.NewInt(1e+18)) // In musicoin code as 'FrontierBlockReward'
	Mcip3BlockReward       = new(big.Int).Mul(big.NewInt(250), big.NewInt(1e+18))
	Mcip8BlockReward       = new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18))
	MusicoinUbiBlockReward = new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18))
	MusicoinDevBlockReward = new(big.Int).Mul(big.NewInt(14), big.NewInt(1e+18))
)

// IsMCIP0 returns whether MCIP0 block is engaged; this is equivalent to 'IsMusicoin'.
// (There is no MCIP-0).
func (c *ChainConfig) IsMCIP0(num *big.Int) bool {
	return isForked(c.MCIP0Block, num)
}

// IsMCIP3 returns whether MCIP3-UBI block is engaged.
func (c *ChainConfig) IsMCIP3(num *big.Int) bool {
	return isForked(c.MCIP3Block, num)
}

// IsMCIP8 returns whether MCIP3-QT block is engaged.
func (c *ChainConfig) IsMCIP8(num *big.Int) bool {
	return isForked(c.MCIP8Block, num)
}
