package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/vars"
)

func EthashBlockReward(c ChainConfigurator, n *big.Int) *big.Int {
	// Select the correct block reward based on chain progression
	blockReward := vars.FrontierBlockReward
	if c == nil || n == nil {
		return blockReward
	}

	if c.IsForked(c.GetEthashEIP1234Transition, n) {
		return vars.EIP1234FBlockReward
	} else if c.IsForked(c.GetEthashEIP649Transition, n) {
		return vars.EIP649FBlockReward
	} else if len(c.GetEthashBlockRewardSchedule()) > 0 {
		// Because the map is not necessarily sorted low-high, we
		// have to ensure that we're walking upwards only.
		var lastActivation uint64
		for activation, reward := range c.GetEthashBlockRewardSchedule() {
			if activation <= n.Uint64() { // Is forked
				if activation >= lastActivation {
					lastActivation = activation
					blockReward = reward
				}
			}
		}
	}

	return blockReward
}
