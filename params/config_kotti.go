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
	KottiGenesisHash = common.HexToHash("0x14c2283285a88fe5fce9bf5c573ab03d6616695d717b12a127188bcacfc743c4")

	KottiChainConfig = &multigeth.MultiGethChainConfig{
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

		// Istanbul eq, aka Aztlan
		// ECIP-1061
		EIP152FBlock:  big.NewInt(2058191),
		EIP1108FBlock: big.NewInt(2058191),
		EIP1344FBlock: big.NewInt(2058191),
		EIP1884FBlock: nil,
		EIP2028FBlock: big.NewInt(2058191),
		EIP2200FBlock: big.NewInt(2058191), // RePetersburg (== re-1283)

		// ECIP-1078, aka Phoenix Fix
		EIP2200DisableFBlock: big.NewInt(2_208_203),
		EIP1283FBlock:        big.NewInt(2_208_203),
		EIP1706FBlock:        big.NewInt(2_208_203),
		ECIP1080FBlock:       big.NewInt(2_208_203),

		ECIP1017FBlock:    big.NewInt(5000000),
		ECIP1017EraRounds: big.NewInt(5000000),

		DisposalBlock:      big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(2000000),

		RequireBlockHashes: map[uint64]common.Hash{
			0: KottiGenesisHash,
		},
	}
)
