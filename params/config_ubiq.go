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
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

var (
	// UbiqChainConfig is the chain parameters to run a node on the Ubiq main network.
	UbiqChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID: 8,
		Ethash: &ctypes.EthashConfig{
			DigishieldV3FBlock:    big.NewInt(0),
			DigishieldV3ModFBlock: big.NewInt(4088),
			FluxFBlock:            big.NewInt(8000),
			UIP0FBlock:            big.NewInt(0),
			UIP1FEpoch:            big.NewInt(22),
		},
		ChainID:                   big.NewInt(8),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		//DAOForkBlock:        big.NewInt(1920000),

		EIP150Block: big.NewInt(0),

		EIP155Block:  big.NewInt(10),
		EIP160FBlock: big.NewInt(10),

		// EIP158~
		EIP161FBlock: big.NewInt(10),
		EIP170FBlock: big.NewInt(10),

		// Byzantium eq, aka Andromeda
		EIP100FBlock: big.NewInt(1075090),
		EIP140FBlock: big.NewInt(1075090),
		EIP198FBlock: big.NewInt(1075090),
		EIP211FBlock: big.NewInt(1075090),
		EIP212FBlock: big.NewInt(1075090),
		EIP213FBlock: big.NewInt(1075090),
		EIP214FBlock: big.NewInt(1075090),
		EIP658FBlock: big.NewInt(1075090),

		// Constantinople eq, aka Andromeda
		EIP145FBlock:  big.NewInt(1075090),
		EIP1014FBlock: big.NewInt(1075090),
		EIP1052FBlock: big.NewInt(1075090),
		// EIP1283FBlock:   big.NewInt(9573000),
		// PetersburgBlock: big.NewInt(9573000),

		// Istanbul eq, aka Taurus
		// ECIP-1088
		EIP152FBlock:  big.NewInt(1_500_000),
		EIP1108FBlock: big.NewInt(1_500_000),
		EIP1344FBlock: big.NewInt(1_500_000),
		EIP1884FBlock: big.NewInt(1_500_000),
		EIP2028FBlock: big.NewInt(1_500_000),
		EIP2200FBlock: big.NewInt(1_500_000), // RePetersburg (=~ re-1283)

		// Berlin eq, aka Aries
		EIP2565FBlock: big.NewInt(math.MaxInt64),
		EIP2718FBlock: big.NewInt(math.MaxInt64),
		EIP2929FBlock: big.NewInt(math.MaxInt64),
		EIP2930FBlock: big.NewInt(math.MaxInt64),

		ECIP1099FBlock: nil, // Etchash

		DisposalBlock:      big.NewInt(0), // Dispose difficulty bomb
		ECIP1017FBlock:     nil,           // Ethereum Classic's disinflationary monetary policy
		ECIP1017EraRounds:  nil,
		ECIP1010PauseBlock: nil, // No need to delay difficulty bomb, is defused by default
		ECIP1010Length:     nil,
		ECBP1100FBlock:     nil, // ECBP1100 (MESS artificial finality)

		RequireBlockHashes: map[uint64]common.Hash{},
	}
)
