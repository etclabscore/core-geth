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

	common2 "github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

var (
	MordorNetworkID          uint64 = 7
	MordorDisposalBlock             = uint64(0)
	MordorECIP1017FBlock            = uint64(2000000)
	MordorECIP1017EraRounds         = uint64(2000000)
	MordorEIP160FBlock              = uint64(0)
	MordorECIP1010PauseBlock        = uint64(0)
	MordorECIP1010Length            = uint64(2000000)

	// MordorChainConfig is the chain parameters to run a node on the Ethereum Classic Mordor test network (PoW).
	MordorChainConfig = func() common2.ChainConfigurator {
		c := &goethereum.ChainConfig{
			ChainID:             big.NewInt(63),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(301243),
			PetersburgBlock:     big.NewInt(301243),
			Ethash:              new(goethereum.EthashConfig),
		}
		c.SetNetworkID(&MordorNetworkID)
		c.SetEthashECIP1041Transition(&MordorDisposalBlock)
		c.SetEthashECIP1041Transition(&MordorECIP1017FBlock)
		c.SetEthashECIP1017EraRounds(&MordorECIP1017EraRounds)
		c.SetEthashECIP1010PauseTransition(&MordorECIP1010PauseBlock)
		c.SetEthashECIP1010ContinueTransition(&MordorECIP1010Length)
		return c
	}()
)
