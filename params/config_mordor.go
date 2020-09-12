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
	// MordorChainConfig is the chain parameters to run a node on the Ethereum Classic Mordor test network (PoW).
	MordorChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID: 7,
		ChainID:   big.NewInt(63),
		Ethash:    new(ctypes.EthashConfig),

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),

		EIP155Block: big.NewInt(0),

		// EIP158 eq
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(301243),
		EIP1014FBlock: big.NewInt(301243),
		EIP1052FBlock: big.NewInt(301243),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(999_983),
		EIP1108FBlock: big.NewInt(999_983),
		EIP1344FBlock: big.NewInt(999_983),
		EIP1884FBlock: big.NewInt(999_983),
		EIP2028FBlock: big.NewInt(999_983),
		EIP2200FBlock: big.NewInt(999_983), // RePetersburg (== re-1283)

		DisposalBlock:      big.NewInt(0),
		ECIP1017FBlock:     big.NewInt(0),
		ECIP1017EraRounds:  big.NewInt(2000000),
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,
		ECBP1100FBlock:     big.NewInt(2290740), // ETA 15 Sept 2020, ~1500 UTC
		RequireBlockHashes: map[uint64]common.Hash{
			840013: common.HexToHash("0x2ceada2b191879b71a5bcf2241dd9bc50d6d953f1640e62f9c2cee941dc61c9d"),
			840014: common.HexToHash("0x8ec29dd692c8985b82410817bac232fc82805b746538d17bc924624fe74a0fcf"),
		},
	}
)
