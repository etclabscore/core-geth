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
	// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
	ClassicChainConfig = &ChainConfig{
		ChainID:             big.NewInt(61),
		HomesteadBlock:      big.NewInt(1150000),
		DAOForkBlock:        big.NewInt(1920000),
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(2500000),
		EIP150Hash:          common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		EIP155Block:         big.NewInt(3000000),
		EIP158Block:         nil,
		ByzantiumBlock:      nil,
		DisposalBlock:       big.NewInt(5900000),
		SocialBlock:         nil,
		EthersocialBlock:    nil,
		ConstantinopleBlock: nil,
		ECIP1017EraRounds:   big.NewInt(5000000),
		EIP160FBlock:        big.NewInt(3000000),
		ECIP1010PauseBlock:  big.NewInt(3000000),
		ECIP1010Length:      big.NewInt(2000000),
		Ethash:              new(EthashConfig),
	}

	DisinflationRateQuotient = big.NewInt(4)      // Disinflation rate quotient for ECIP1017
	DisinflationRateDivisor  = big.NewInt(5)      // Disinflation rate divisor for ECIP1017
	ExpDiffPeriod            = big.NewInt(100000) // Exponential diff period for diff bomb & ECIP1010
)
