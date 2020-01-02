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
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
)

var (
	// Genesis hashes to enforce below configs on.
	SocialGenesisHash = common.HexToHash("0xba8314d5c2ebddaf58eb882b364b27cbfa4d3402dacd32b60986754ac25cfe8d")

	// SocialChainConfig is the chain parameters to run a node on the Ethereum Social main network.
	SocialChainConfig = &multigeth.MultiGethChainConfig{
		ChainID:           big.NewInt(28),
		EIP2FBlock:        big.NewInt(0),
		EIP7FBlock:        big.NewInt(0),
		EIP150Block:       big.NewInt(0),
		EIP155Block:       big.NewInt(0),
		Ethash:            new(ctypes.EthashConfig),
		NetworkID:         28,
		DisposalBlock:     big.NewInt(0),
		SocialBlock:       big.NewInt(0),
		EthersocialBlock:  nil,
		ECIP1017FBlock:    big.NewInt(5000000),
		ECIP1017EraRounds: big.NewInt(5000000),
		EIP160FBlock:      big.NewInt(0),
		BlockRewardSchedule: ctypes.Uint64BigMapEncodesHex{
			0: new(big.Int).Mul(big.NewInt(50), big.NewInt(1e+18)),
		},
		RequireBlockHashes: map[uint64]common.Hash{
			0: common.HexToHash("0xba8314d5c2ebddaf58eb882b364b27cbfa4d3402dacd32b60986754ac25cfe8d"),
		},
	}
)
