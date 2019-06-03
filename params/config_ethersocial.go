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
	EthersocialGenesisHash = common.HexToHash("0x310dd3c4ae84dd89f1b46cfdd5e26c8f904dfddddc73f323b468127272e20e9f")

	// EthersocialChainConfig is the chain parameters to run a node on the Ethersocial main network.
	EthersocialChainConfig = &ChainConfig{
		ChainID:             big.NewInt(31102),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x310dd3c4ae84dd89f1b46cfdd5e26c8f904dfddddc73f323b468127272e20e9f"),
		EIP155Block:         big.NewInt(845000),
		EIP158Block:         big.NewInt(845000),
		ByzantiumBlock:      big.NewInt(600000),
		DisposalBlock:       nil,
		SocialBlock:         nil,
		EthersocialBlock:    big.NewInt(0),
		ConstantinopleBlock: nil,
		Ethash:              new(EthashConfig),
	}

	EthersocialBlockReward = big.NewInt(5e+18) // Block reward in wei for successfully mining a block upward for Ethersocial Network
)

func (c *ChainConfig) IsEthersocial(num *big.Int) bool {
	return isForked(c.EthersocialBlock, num)
}
