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
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

var (
	// Genesis hashes to enforce below configs on.
	KottiGenesisHash = common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4")

	KottiNetworkID          uint64 = 6
	//KottiDisposalBlock             = uint64(0)
	//KottiECIP1017FBlock            = uint64(5000000)
	//KottiECIP1017EraRounds         = uint64(5000000)
	//KottiEIP160FBlock              = uint64(0)
	//KottiECIP1010PauseBlock        = uint64(0)
	//KottiECIP1010Length            = uint64(2000000)

	KottiChainConfig = func() ctypes.ChainConfigurator {
		c := &goethereum.ChainConfig{
			ChainID:             big.NewInt(6),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP150Hash:          common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4"),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(716617),
			ByzantiumBlock:      big.NewInt(716617),
			ConstantinopleBlock: big.NewInt(1705549),
			PetersburgBlock:     big.NewInt(1705549),
			Clique: &ctypes.CliqueConfig{
				Period: 15,
				Epoch:  30000,
			},
		}
		c.SetNetworkID(&KottiNetworkID)

		return c
	}()
)
