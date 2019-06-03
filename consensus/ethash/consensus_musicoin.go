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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func musicoinBlockReward(config *params.ChainConfig, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	// Select the correct block reward based on chain progression
	blockReward := params.Mcip0BlockReward
	mcip3Reward := params.Mcip3BlockReward
	mcip8Reward := params.Mcip8BlockReward
	ubiReservoir := params.MusicoinUbiBlockReward
	devReservoir := params.MusicoinDevBlockReward

	reward := new(big.Int).Set(blockReward)

	if config.IsMCIP8(header.Number) {
		state.AddBalance(header.Coinbase, mcip8Reward)
		state.AddBalance(common.HexToAddress("0x00eFdd5883eC628983E9063c7d969fE268BBf310"), ubiReservoir)
		state.AddBalance(common.HexToAddress("0x00756cF8159095948496617F5FB17ED95059f536"), devReservoir)
		blockReward := mcip8Reward
		reward := new(big.Int).Set(blockReward)
		_ = reward
	} else if config.IsMCIP3(header.Number) {
		state.AddBalance(header.Coinbase, mcip3Reward)
		state.AddBalance(common.HexToAddress("0x00eFdd5883eC628983E9063c7d969fE268BBf310"), ubiReservoir)
		state.AddBalance(common.HexToAddress("0x00756cF8159095948496617F5FB17ED95059f536"), devReservoir)
		// no change to uncle reward during UBI fork, a mistake but now a legacy
	} else {
		state.AddBalance(header.Coinbase, reward)
	}

	// Accumulate the rewards for the miner and any included uncles
	r := new(big.Int)
	for _, uncle := range uncles {
		r.Add(uncle.Number, big8)
		r.Sub(r, header.Number)
		r.Mul(r, blockReward)
		r.Div(r, big8)
		state.AddBalance(uncle.Coinbase, r)

		r.Div(blockReward, big32)
		reward.Add(reward, r)
	}
}
