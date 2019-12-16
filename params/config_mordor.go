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

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
)

var (
	//MordorNetworkID         uint64 = 7
	//MordorDisposalBlock            = uint64(0)
	//MordorECIP1017FBlock           = uint64(2000000)
	//MordorECIP1017EraRounds        = uint64(2000000)
	//MordorEIP160FBlock             = uint64(0)

	// MordorChainConfig is the chain parameters to run a node on the Ethereum Classic Mordor test network (PoW).
	MordorChainConfig = func() ctypes.ChainConfigurator {
		//c := &goethereum.ChainConfig{
		//	ChainID:             big.NewInt(63),
		//	HomesteadBlock:      big.NewInt(0),
		//	EIP150Block:         big.NewInt(0),
		//	EIP155Block:         big.NewInt(0),
		//	EIP158Block:         big.NewInt(0),
		//	ByzantiumBlock:      big.NewInt(0),
		//	ConstantinopleBlock: big.NewInt(301243),
		//	PetersburgBlock:     big.NewInt(301243),
		//	Ethash:              new(goethereum.EthashConfig),
		//}
		//cc := &paramtypes.MultiGethChainConfig{}
		//err := convert.Convert(c, cc)
		//if err != nil {
		//	panic(err)
		//}
		//
		//cc.SetNetworkID(&MordorNetworkID)
		//cc.SetEthashECIP1041Transition(&MordorDisposalBlock)
		//cc.SetEthashECIP1017Transition(&MordorECIP1017FBlock)
		//cc.SetEthashECIP1017EraRounds(&MordorECIP1017EraRounds)

		return &paramtypes.MultiGethChainConfig{
			NetworkID:  7,
			ChainID:    big.NewInt(63),
			Ethash:     new(goethereum.EthashConfig),
			EIP2FBlock: big.NewInt(0),
			EIP7FBlock: big.NewInt(0),
			//DAOForkBlock:        big.NewInt(1920000),
			EIP150Block:        big.NewInt(0),
			EIP155Block:        big.NewInt(0),
			EIP160FBlock:       big.NewInt(0),
			EIP161FBlock:       big.NewInt(0),
			EIP170FBlock:       big.NewInt(0),
			EIP100FBlock:       big.NewInt(0),
			EIP140FBlock:       big.NewInt(0),
			EIP198FBlock:       big.NewInt(0),
			EIP211FBlock:       big.NewInt(0),
			EIP212FBlock:       big.NewInt(0),
			EIP213FBlock:       big.NewInt(0),
			EIP214FBlock:       big.NewInt(0),
			EIP658FBlock:       big.NewInt(0),
			EIP145FBlock:       big.NewInt(301243),
			EIP1014FBlock:      big.NewInt(301243),
			EIP1052FBlock:      big.NewInt(301243),
			EIP1283FBlock:      nil,
			PetersburgBlock:    nil, // Disable 1283
			EIP2200FBlock:      nil, // RePetersburg
			DisposalBlock:      big.NewInt(0),
			ECIP1017FBlock:     big.NewInt(0),
			ECIP1017EraRounds:  big.NewInt(2000000),
			ECIP1010PauseBlock: nil,
			ECIP1010Length:     nil,
		}
	}()
)
