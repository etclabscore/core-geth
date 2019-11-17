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
package ethash

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func ecip1017BlockReward(config *params.ChainConfig, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	blockReward := params.FrontierBlockReward

	// Ensure value 'era' is configured.
	eraLen := config.ECIP1017EraRounds
	era := GetBlockEra(header.Number, eraLen)
	wr := GetBlockWinnerRewardByEra(era, blockReward)                    // wr "winner reward". 5, 4, 3.2, 2.56, ...
	wurs := GetBlockWinnerRewardForUnclesByEra(era, uncles, blockReward) // wurs "winner uncle rewards"
	wr.Add(wr, wurs)
	state.AddBalance(header.Coinbase, wr) // $$

	// Reward uncle miners.
	for _, uncle := range uncles {
		ur := GetBlockUncleRewardByEra(era, header, uncle, blockReward)
		state.AddBalance(uncle.Coinbase, ur) // $$
	}
}

func ecip1010Explosion(config *params.ChainConfig, next *big.Int, exPeriodRef *big.Int) {
	// https://github.com/ethereumproject/ECIPs/blob/master/ECIPs/ECIP-1010.md

	explosionBlock := new(big.Int).Add(config.ECIP1010PauseBlock, config.ECIP1010Length)
	if next.Cmp(explosionBlock) < 0 {
		exPeriodRef.Set(config.ECIP1010PauseBlock)
	} else {
		exPeriodRef.Sub(exPeriodRef, config.ECIP1010Length)
	}
}

// GetBlockEra gets which "Era" a given block is within, given an era length (ecip-1017 has era=5,000,000 blocks)
// Returns a zero-index era number, so "Era 1": 0, "Era 2": 1, "Era 3": 2 ...
func GetBlockEra(blockNum, eraLength *big.Int) *big.Int {
	// If genesis block or impossible negative-numbered block, return zero-val.
	if blockNum.Sign() < 1 {
		return new(big.Int)
	}

	remainder := big.NewInt(0).Mod(big.NewInt(0).Sub(blockNum, big.NewInt(1)), eraLength)
	base := big.NewInt(0).Sub(blockNum, remainder)

	d := big.NewInt(0).Div(base, eraLength)
	dremainder := big.NewInt(0).Mod(d, big.NewInt(1))

	return new(big.Int).Sub(d, dremainder)
}
