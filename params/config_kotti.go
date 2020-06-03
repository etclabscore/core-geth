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
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

var (
	// Genesis hashes to enforce below configs on.
	KottiGenesisHash = common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4")

	KottiChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID: 6,
		ChainID:   big.NewInt(6),
		Clique: &ctypes.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block: big.NewInt(0),

		// EIP158 eq
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(716617),
		EIP170FBlock: big.NewInt(716617),

		// Byzantium eq
		EIP100FBlock: big.NewInt(716617),
		EIP140FBlock: big.NewInt(716617),
		EIP198FBlock: big.NewInt(716617),
		EIP211FBlock: big.NewInt(716617),
		EIP212FBlock: big.NewInt(716617),
		EIP213FBlock: big.NewInt(716617),
		EIP214FBlock: big.NewInt(716617),
		EIP658FBlock: big.NewInt(716617),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(1705549),
		EIP1014FBlock: big.NewInt(1705549),
		EIP1052FBlock: big.NewInt(1705549),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(2_200_013),
		EIP1108FBlock: big.NewInt(2_200_013),
		EIP1344FBlock: big.NewInt(2_200_013),
		EIP1884FBlock: big.NewInt(2_200_013),
		EIP2028FBlock: big.NewInt(2_200_013),
		EIP2200FBlock: big.NewInt(2_200_013), // RePetersburg (== re-1283)

		RequireBlockHashes: map[uint64]common.Hash{
			0: KottiGenesisHash,
			/*
				########## BAD BLOCK #########
				Chain config: NetworkID: 6, ChainID: 6 Engine: clique EIP1014: 1705549 EIP1052: 1705549 EIP1108: 2200013 EIP1344: 2200013 EIP140: 716617 EIP145: 1705549 EIP1
				50: 0 EIP152: 2200013 EIP155: 0 EIP160: 0 EIP161abc: 716617 EIP161d: 716617 EIP170: 716617 EIP1884: 2200013 EIP198: 716617 EIP2028: 2200013 EIP211: 716617 EI
				P212: 716617 EIP213: 716617 EIP214: 716617 EIP2200: 2200013 EIP658: 716617 EIP7: 0 EthashECIP1010Continue: 2000000 EthashECIP1010Pause: 0 EthashECIP1017: 500
				0000 EthashECIP1041: 0 EthashEIP100B: 716617 EthashEIP2: 0 EthashHomestead: 0
				Number: 2058192
				Hash: 0xf56efd9b074eed67c79d903c26bd0059e701c37200dc3cb86cc6c38159406d36
				         0: cumulative: 25352 gas: 25352 contract: 0x0000000000000000000000000000000000000000 status: 1 tx: 0x0ae3504d52fed3fdc92fb4267d6b36523d022b33663fc44
				b4eeddea0bcaa51cb logs: [] bloom: 000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
				0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
				0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
				000000000000000000000000000000000000000000000000000000000000000000000000000 state:
				Error: invalid gas used (remote: 22024 local: 25352)
				##############################
			*/
			2_058_192: common.HexToHash("0x60b1f84737789f0233d8483fb0e6f460e068f690373893a878a2d12c9e59be1e"),
		},
	}
)
