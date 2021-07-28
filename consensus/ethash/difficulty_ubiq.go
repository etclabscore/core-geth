// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	"math/big"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	big88                 = big.NewInt(88)
	minimumDifficultyUbiq = big.NewInt(131072)
	digishieldV3Config    = &diffConfig{
		AveragingWindow: big.NewInt(21),
		MaxAdjustDown:   big.NewInt(16), // 16%
		MaxAdjustUp:     big.NewInt(8),  // 8%
		Factor:          big.NewInt(100),
	}

	digishieldV3ModConfig = &diffConfig{
		AveragingWindow: big.NewInt(88),
		MaxAdjustDown:   big.NewInt(3), // 3%
		MaxAdjustUp:     big.NewInt(2), // 2%
		Factor:          big.NewInt(100),
	}

	fluxConfig = &diffConfig{
		AveragingWindow: big.NewInt(88),
		MaxAdjustDown:   big.NewInt(5), // 0.5%
		MaxAdjustUp:     big.NewInt(3), // 0.3%
		Dampen:          big.NewInt(1), // 0.1%
		Factor:          big.NewInt(1000),
	}
)

type diffConfig struct {
	AveragingWindow *big.Int `json:"averagingWindow"`
	MaxAdjustDown   *big.Int `json:"maxAdjustDown"`
	MaxAdjustUp     *big.Int `json:"maxAdjustUp"`
	Dampen          *big.Int `json:"dampen,omitempty"`
	Factor          *big.Int `json:"factor"`
}

// Difficulty timespans
func averagingWindowTimespan(config *diffConfig) *big.Int {
	x := new(big.Int)
	return x.Mul(config.AveragingWindow, big88)
}

func minActualTimespan(config *diffConfig, dampen bool) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	if dampen {
		x.Sub(config.Factor, config.Dampen)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	} else {
		x.Sub(config.Factor, config.MaxAdjustUp)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	}
	return z
}

func maxActualTimespan(config *diffConfig, dampen bool) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	if dampen {
		x.Add(config.Factor, config.Dampen)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	} else {
		x.Add(config.Factor, config.MaxAdjustDown)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	}
	return z
}

// CalcDifficulty determines which difficulty algorithm to use for calculating a new block
func calcDifficultyUbiq(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	parentTime := parent.Time
	parentNumber := parent.Number
	parentDiff := parent.Difficulty

	config := chain.Config()
	ethashConfig := config.GetEthashConfig()

	if parentNumber.Cmp(ethashConfig.FluxFBlock) < 0 {
		if parentNumber.Cmp(ethashConfig.DigishieldV3ModFBlock) < 0 {
			// Original DigishieldV3
			return calcDifficultyDigishieldV3(chain, parentNumber, parentDiff, parent, digishieldV3Config)
		}
		// Modified DigishieldV3
		return calcDifficultyDigishieldV3(chain, parentNumber, parentDiff, parent, digishieldV3ModConfig)
	}
	// Flux
	return calcDifficultyFlux(chain, big.NewInt(int64(time)), big.NewInt(int64(parentTime)), parentNumber, parentDiff, parent)
}

// calcDifficultyDigishieldV3 is the original difficulty adjustment algorithm.
// It returns the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
// Based on Digibyte's Digishield v3 retargeting
func calcDifficultyDigishieldV3(chain consensus.ChainHeaderReader, parentNumber, parentDiff *big.Int, parent *types.Header, digishield *diffConfig) *big.Int {
	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	nFirstBlock := new(big.Int)
	nFirstBlock.Sub(parentNumber, digishield.AveragingWindow)

	// Check we have enough blocks
	if parentNumber.Cmp(digishield.AveragingWindow) < 1 {
		x.Set(parentDiff)
		return x
	}

	// Limit adjustment step
	// Use medians to prevent time-warp attacks
	nLastBlockTime := chain.CalcPastMedianTime(parentNumber.Uint64(), parent)
	nFirstBlockTime := chain.CalcPastMedianTime(nFirstBlock.Uint64(), parent)
	nActualTimespan := new(big.Int)
	nActualTimespan.Sub(nLastBlockTime, nFirstBlockTime)

	y := new(big.Int)
	y.Sub(nActualTimespan, averagingWindowTimespan(digishield))
	y.Div(y, big.NewInt(4))
	nActualTimespan.Add(y, averagingWindowTimespan(digishield))
	if nActualTimespan.Cmp(minActualTimespan(digishield, false)) < 0 {
		nActualTimespan.Set(minActualTimespan(digishield, false))
	} else if nActualTimespan.Cmp(maxActualTimespan(digishield, false)) > 0 {
		nActualTimespan.Set(maxActualTimespan(digishield, false))
	}

	// Retarget
	x.Mul(parentDiff, averagingWindowTimespan(digishield))
	x.Div(x, nActualTimespan)

	if x.Cmp(minimumDifficultyUbiq) < 0 {
		x.Set(minimumDifficultyUbiq)
	}

	return x
}

func calcDifficultyFlux(chain consensus.ChainHeaderReader, time, parentTime, parentNumber, parentDiff *big.Int, parent *types.Header) *big.Int {
	x := new(big.Int)
	nFirstBlock := new(big.Int)
	nFirstBlock.Sub(parentNumber, fluxConfig.AveragingWindow)

	// Check we have enough blocks
	if parentNumber.Cmp(fluxConfig.AveragingWindow) < 1 {
		x.Set(parentDiff)
		return x
	}

	diffTime := new(big.Int)
	diffTime.Sub(time, parentTime)

	nLastBlockTime := chain.CalcPastMedianTime(parentNumber.Uint64(), parent)
	nFirstBlockTime := chain.CalcPastMedianTime(nFirstBlock.Uint64(), parent)
	nActualTimespan := new(big.Int)
	nActualTimespan.Sub(nLastBlockTime, nFirstBlockTime)

	y := new(big.Int)
	y.Sub(nActualTimespan, averagingWindowTimespan(fluxConfig))
	y.Div(y, big.NewInt(4))
	nActualTimespan.Add(y, averagingWindowTimespan(fluxConfig))

	if nActualTimespan.Cmp(minActualTimespan(fluxConfig, false)) < 0 {
		doubleBig88 := new(big.Int)
		doubleBig88.Mul(big88, big.NewInt(2))
		if diffTime.Cmp(doubleBig88) > 0 {
			nActualTimespan.Set(minActualTimespan(fluxConfig, true))
		} else {
			nActualTimespan.Set(minActualTimespan(fluxConfig, false))
		}
	} else if nActualTimespan.Cmp(maxActualTimespan(fluxConfig, false)) > 0 {
		halfBig88 := new(big.Int)
		halfBig88.Div(big88, big.NewInt(2))
		if diffTime.Cmp(halfBig88) < 0 {
			nActualTimespan.Set(maxActualTimespan(fluxConfig, true))
		} else {
			nActualTimespan.Set(maxActualTimespan(fluxConfig, false))
		}
	}

	x.Mul(parentDiff, averagingWindowTimespan(fluxConfig))
	x.Div(x, nActualTimespan)

	if x.Cmp(minimumDifficultyUbiq) < 0 {
		x.Set(minimumDifficultyUbiq)
	}

	return x
}
