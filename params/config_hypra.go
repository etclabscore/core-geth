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
	"github.com/ethereum/go-ethereum/params/vars"
)

var (
	// HypraChainConfig is the chain parameters to run a node on the Classic main network.
	HypraChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 622277,
		EthashB3:                  new(ctypes.EthashB3Config),
		ChainID:                   big.NewInt(622277),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		// DAOForkBlock:        big.NewInt(1920000),

		EIP150Block: big.NewInt(0),

		EIP155Block:        big.NewInt(0),
		EIP160FBlock:       big.NewInt(0),
		ECIP1010PauseBlock: big.NewInt(0),
		ECIP1010Length:     big.NewInt(0),

		ECIP1017FBlock:    big.NewInt(0),
		ECIP1017EraRounds: big.NewInt(0),

		DisposalBlock: big.NewInt(0),

		// EIP158~
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(1001),
		EIP140FBlock: big.NewInt(1001),
		EIP198FBlock: big.NewInt(1001),
		EIP211FBlock: big.NewInt(1001),
		EIP212FBlock: big.NewInt(1001),
		EIP213FBlock: big.NewInt(1001),
		EIP214FBlock: big.NewInt(1001),
		EIP658FBlock: big.NewInt(1001),

		// Constantinople eq, aka Agharta
		EIP145FBlock:    big.NewInt(5503),
		EIP1014FBlock:   big.NewInt(5503),
		EIP1052FBlock:   big.NewInt(5503),
		EIP1283FBlock:   big.NewInt(5503),
		PetersburgBlock: big.NewInt(5507),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(5519),
		EIP1108FBlock: big.NewInt(5519),
		EIP1344FBlock: big.NewInt(5519),
		EIP1884FBlock: big.NewInt(5519),
		EIP2028FBlock: big.NewInt(5519),
		EIP2200FBlock: big.NewInt(5519), // RePetersburg (=~ re-1283)

		// Hypra does not use ETC Improvements
		// ECBP1100FBlock:           big.NewInt(11_380_000), // ETA 09 Oct 2020
		// ECBP1100DeactivateFBlock: big.NewInt(19_250_000), // ETA 31 Jan 2023 (== Spiral hard fork)
		// ECIP1099FBlock:           big.NewInt(11_700_000), // Etchash (DAG size limit)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(5527),
		EIP2718FBlock: big.NewInt(5527),
		EIP2929FBlock: big.NewInt(5527),
		EIP2930FBlock: big.NewInt(5527),

		// London (partially), aka Mystique
		EIP3529FBlock: big.NewInt(14_525_000),
		EIP3541FBlock: big.NewInt(14_525_000),

		// Spiral, aka Shanghai (partially)
		// EIP4399FBlock: nil, // Supplant DIFFICULTY with PREVRANDAO. ETC does not spec 4399 because it's still PoW, and 4399 is only applicable for the PoS system.
		EIP3651FBlock: big.NewInt(19_250_000), // Warm COINBASE (gas reprice)
		EIP3855FBlock: big.NewInt(19_250_000), // PUSH0 instruction
		EIP3860FBlock: big.NewInt(19_250_000), // Limit and meter initcode
		// EIP4895FBlock: nil, // Beacon chain push withdrawals as operations
		EIP6049FBlock: big.NewInt(19_250_000), // Deprecate SELFDESTRUCT (noop)

		RequireBlockHashes: map[uint64]common.Hash{
			1920000: common.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f"),
			2500000: common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		},
	}
)
