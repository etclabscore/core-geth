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
	KottiGenesisHash = common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4")

	// KottiChainConfig is the chain parameters to run a node on the Kotti main network.
	KottiChainConfig = &ChainConfig{
		ChainID:             big.NewInt(6),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         nil,
		ByzantiumBlock:      nil,
		DisposalBlock:       big.NewInt(0),
		SocialBlock:         nil,
		EthersocialBlock:    nil,
		ConstantinopleBlock: nil,
		ECIP1017EraRounds:   big.NewInt(5000000),
		EIP160FBlock:        big.NewInt(0),
		ECIP1010PauseBlock:  big.NewInt(0),
		ECIP1010Length:      big.NewInt(2000000),
		Clique: &CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// KottiBootnodes are the enode URLs of the P2P bootstrap nodes running on the
	// Kotti test network.
	KottiBootnodes = []string{
		"enode://06333009fc9ef3c9e174768e495722a7f98fe7afd4660542e983005f85e556028410fd03278944f44cfe5437b1750b5e6bd1738f700fe7da3626d52010d2954c@51.141.15.254:30303",
		"enode://ae8658da8d255d1992c3ec6e62e11d6e1c5899aa1566504bc1ff96a0c9c8bd44838372be643342553817f5cc7d78f1c83a8093dee13d77b3b0a583c050c81940@18.232.185.151:30303",
		"enode://67913271d14f445689e8310270c304d42f268428f2de7a4ac0275bea97690e021df6f549f462503ff4c7a81d9dd27288867bbfa2271477d0911378b8944fae55@157.230.239.163:30303",
		"enode://e8a786a894db053fe6886e283fc4385389ad034e04a692a26335f30b714059efd5cead0e410ecd783ce095888fdafcc21a685f13501594e969d6f5ac7ba0388c@86.103.236.55:63384",
	}
)
