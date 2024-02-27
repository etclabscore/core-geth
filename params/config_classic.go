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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/holiman/uint256"
)

var (
	// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
	ClassicChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 1,
		Ethash:                    new(ctypes.EthashConfig),
		ChainID:                   big.NewInt(61),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(1150000),
		EIP7FBlock: big.NewInt(1150000),

		// DAOForkBlock:        big.NewInt(1920000),

		EIP150Block: big.NewInt(2500000),

		EIP155Block:        big.NewInt(3000000),
		EIP160FBlock:       big.NewInt(3000000),
		ECIP1010PauseBlock: big.NewInt(3000000),
		ECIP1010Length:     big.NewInt(2000000),

		ECIP1017FBlock:    big.NewInt(5000000),
		ECIP1017EraRounds: big.NewInt(5000000),

		DisposalBlock: big.NewInt(5900000),

		// EIP158~
		EIP161FBlock: big.NewInt(8772000),
		EIP170FBlock: big.NewInt(8772000),

		// Byzantium eq
		EIP100FBlock: big.NewInt(8772000),
		EIP140FBlock: big.NewInt(8772000),
		EIP198FBlock: big.NewInt(8772000),
		EIP211FBlock: big.NewInt(8772000),
		EIP212FBlock: big.NewInt(8772000),
		EIP213FBlock: big.NewInt(8772000),
		EIP214FBlock: big.NewInt(8772000),
		EIP658FBlock: big.NewInt(8772000),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(9573000),
		EIP1014FBlock: big.NewInt(9573000),
		EIP1052FBlock: big.NewInt(9573000),
		// EIP1283FBlock:   big.NewInt(9573000),
		// PetersburgBlock: big.NewInt(9573000),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(10_500_839),
		EIP1108FBlock: big.NewInt(10_500_839),
		EIP1344FBlock: big.NewInt(10_500_839),
		EIP1884FBlock: big.NewInt(10_500_839),
		EIP2028FBlock: big.NewInt(10_500_839),
		EIP2200FBlock: big.NewInt(10_500_839), // RePetersburg (=~ re-1283)

		ECBP1100FBlock:           big.NewInt(11_380_000), // ETA 09 Oct 2020
		ECBP1100DeactivateFBlock: big.NewInt(19_250_000), // ETA 31 Jan 2023 (== Spiral hard fork)
		ECIP1099FBlock:           big.NewInt(11_700_000), // Etchash (DAG size limit)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(13_189_133),
		EIP2718FBlock: big.NewInt(13_189_133),
		EIP2929FBlock: big.NewInt(13_189_133),
		EIP2930FBlock: big.NewInt(13_189_133),

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

	DisinflationRateQuotient = uint256.NewInt(4)      // Disinflation rate quotient for ECIP1017
	DisinflationRateDivisor  = uint256.NewInt(5)      // Disinflation rate divisor for ECIP1017
	ExpDiffPeriod            = uint256.NewInt(100000) // Exponential diff period for diff bomb & ECIP1010

	MessNetConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 1,
		Ethash:                    new(ctypes.EthashConfig),
		ChainID:                   big.NewInt(6161),
		SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(1),
		EIP7FBlock: big.NewInt(1),

		DAOForkBlock: nil,

		EIP150Block: big.NewInt(2),

		EIP155Block:  big.NewInt(3),
		EIP160FBlock: big.NewInt(3),

		// EIP158~
		EIP161FBlock: big.NewInt(8),
		EIP170FBlock: big.NewInt(8),

		// Byzantium eq
		EIP100FBlock: big.NewInt(8),
		EIP140FBlock: big.NewInt(8),
		EIP198FBlock: big.NewInt(8),
		EIP211FBlock: big.NewInt(8),
		EIP212FBlock: big.NewInt(8),
		EIP213FBlock: big.NewInt(8),
		EIP214FBlock: big.NewInt(8),
		EIP658FBlock: big.NewInt(8),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(9),
		EIP1014FBlock: big.NewInt(9),
		EIP1052FBlock: big.NewInt(9),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(10),
		EIP1108FBlock: big.NewInt(10),
		EIP1344FBlock: big.NewInt(10),
		EIP1884FBlock: big.NewInt(10),
		EIP2028FBlock: big.NewInt(10),
		EIP2200FBlock: big.NewInt(10), // RePetersburg (=~ re-1283)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(11),
		EIP2718FBlock: big.NewInt(11),
		EIP2929FBlock: big.NewInt(11),
		EIP2930FBlock: big.NewInt(11),

		DisposalBlock:      big.NewInt(5),
		ECIP1017FBlock:     big.NewInt(5),
		ECIP1017EraRounds:  big.NewInt(5000),
		ECIP1010PauseBlock: big.NewInt(3),
		ECIP1010Length:     big.NewInt(2),
		ECBP1100FBlock:     big.NewInt(11),
	}
)

func DefaultMessNetGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     MessNetConfig,
		Timestamp:  1598650845,
		ExtraData:  hexutil.MustDecode("0x4235353535353535353535353535353535353535353535353535353535353535"),
		GasLimit:   10485760,
		Difficulty: big.NewInt(37103392657464),
		Alloc: map[common.Address]genesisT.GenesisAccount{
			common.BytesToAddress([]byte{1}): {Balance: big.NewInt(1)}, // ECRecover
			common.BytesToAddress([]byte{2}): {Balance: big.NewInt(1)}, // SHA256
			common.BytesToAddress([]byte{3}): {Balance: big.NewInt(1)}, // RIPEMD
			common.BytesToAddress([]byte{4}): {Balance: big.NewInt(1)}, // Identity
			common.BytesToAddress([]byte{5}): {Balance: big.NewInt(1)}, // ModExp
			common.BytesToAddress([]byte{6}): {Balance: big.NewInt(1)}, // ECAdd
			common.BytesToAddress([]byte{7}): {Balance: big.NewInt(1)}, // ECScalarMul
			common.BytesToAddress([]byte{8}): {Balance: big.NewInt(1)}, // ECPairing
		},
	}
}
