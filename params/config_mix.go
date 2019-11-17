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
	MixGenesisHash = common.HexToHash("0x4fa57903dad05875ddf78030c16b5da886f7d81714cf66946a4c02566dbb2af5")

	// MixChainConfig is the chain parameters to run a node on the MIX main network.
	MixChainConfig = &ChainConfig{
		NetworkID:           76,
		ChainID:             big.NewInt(76),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x4fa57903dad05875ddf78030c16b5da886f7d81714cf66946a4c02566dbb2af5"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      nil,
		DisposalBlock:       nil,
		ConstantinopleBlock: nil,
		EIP160FBlock:        big.NewInt(0),
	}
)
