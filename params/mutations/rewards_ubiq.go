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
package mutations

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

// Some weird constants to avoid constant memory allocs for them.
var (
	big2 = big.NewInt(2)
)

// Ubiq monetary policy reward step
type UbiqMPStep struct {
	Block  *big.Int `json:"block"`
	Reward *big.Int `json:"reward"`
}

var UbiqMonetaryPolicy = []UbiqMPStep{
	{
		Block:  big.NewInt(0),
		Reward: big.NewInt(8e+18),
	},
	{
		Block:  big.NewInt(358363),
		Reward: big.NewInt(7e+18),
	},
	{
		Block:  big.NewInt(716727),
		Reward: big.NewInt(6e+18),
	},
	{
		Block:  big.NewInt(1075090),
		Reward: big.NewInt(5e+18),
	},
	{
		Block:  big.NewInt(1433454),
		Reward: big.NewInt(4e+18),
	},
	{
		Block:  big.NewInt(1791818),
		Reward: big.NewInt(3e+18),
	},
	{
		Block:  big.NewInt(2150181),
		Reward: big.NewInt(2e+18),
	},
	{
		Block:  big.NewInt(2508545),
		Reward: big.NewInt(1e+18),
	},
}

// CalcBaseBlockReward calculates the base block reward as per the ubiq monetary policy.
func CalcBaseBlockReward(height *big.Int) (*big.Int, *big.Int) {
	reward := new(big.Int)

	for _, step := range UbiqMonetaryPolicy {
		if height.Cmp(step.Block) > 0 {
			reward = new(big.Int).Set(step.Reward)
		} else {
			break
		}
	}

	return new(big.Int).Set(UbiqMonetaryPolicy[0].Reward), reward
}

// CalcUncleBlockReward calculates the uncle miner reward based on depth.
func CalcUncleBlockReward(config ctypes.ChainConfigurator, blockHeight *big.Int, uncleHeight *big.Int, blockReward *big.Int) *big.Int {
	reward := new(big.Int)
	// calculate reward based on depth
	reward.Add(uncleHeight, big2)
	reward.Sub(reward, blockHeight)
	reward.Mul(reward, blockReward)
	reward.Div(reward, big2)

	// negative uncle reward fix. (activates along-side EIP158/160)
	if config.IsEnabled(config.GetEIP160Transition, blockHeight) && reward.Cmp(big.NewInt(0)) < 0 {
		reward = big.NewInt(0)
	}
	return reward
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
func accumulateUbiqRewards(config ctypes.ChainConfigurator, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	// block reward (miner)
	initialReward, currentReward := CalcBaseBlockReward(header.Number)

	// Uncle reward step down fix. (activates along-side andromeda)
	ufixReward := initialReward
	if config.IsEnabled(config.GetEIP1052Transition, header.Number) {
		ufixReward = currentReward
	}

	for _, uncle := range uncles {
		// uncle block miner reward (depth === 1 ? baseBlockReward * 0.5 : 0)
		uncleReward := CalcUncleBlockReward(config, header.Number, uncle.Number, ufixReward)
		// update uncle miner balance
		state.AddBalance(uncle.Coinbase, uncleReward)
		// include uncle bonus reward (baseBlockReward/32)
		uncleReward.Div(ufixReward, big32)
		currentReward.Add(currentReward, uncleReward)
	}
	// update block miner balance
	state.AddBalance(header.Coinbase, currentReward)
}
