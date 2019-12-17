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

package goethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/convert"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/internal"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"github.com/ethereum/go-ethereum/params/vars"
)

// File contains the go-ethereum implementation of the Configurator interface.
// TODO: Handle 'unsupported' Feature, Fork-only cases (where unequal feature settings cause _undetermined_ behavior),
// eg. where SetEIP1052 -> Constantinople AND SetEIP145 -> Constantinople;
// If these values are different, the GetConstantinople result is undetermined.
// Maybe this should return an error which is handled by any conversion, either logging to Warn, or something similar;
// "You're configuring a chainspec from a more precise chainspec, and if not used carefully, some settings may be undetermined
// and/or unexpected."

func newU64(u uint64) *uint64 {
	return &u
}

func bigNewU64(i *big.Int) *uint64 {
	if i == nil {
		return nil
	}
	return newU64(i.Uint64())
}

func setBig(i *big.Int, u *uint64) *big.Int {
	if u == nil {
		return nil
	}
	i = big.NewInt(int64(*u))
	return i
}

func (c *ChainConfig) SwapIfAlt() ctypes.ChainConfigurator {
	switch ChainConfigUseAlt {
	case "MULTICC":
		if c.altDT != nil {
			return c.altDT
		}
		mgc := &multigeth.MultiGethChainConfig{}
		if err := convert.Convert(c, mgc); err != nil {
			panic(err)
		}
		c.altDT = mgc
		return c.altDT
	case "PARITYCC":
		if c.altDT != nil {
			return c.altDT
		}
		pc := &parity.ParityChainSpec{}
		if err := convert.Convert(c, pc); err != nil {
			panic(err)
		}
		c.altDT = pc
		return c.altDT
	default:
		panic("invalid alt chain config")
	}
}

func (c *ChainConfig) GetAccountStartNonce() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetAccountStartNonce()
		}
	}

	return internal.One().GetAccountStartNonce()
}

func (c *ChainConfig) SetAccountStartNonce(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetAccountStartNonce(n)
		}
	}

	return internal.One().SetAccountStartNonce(n)
}

func (c *ChainConfig) GetMaximumExtraDataSize() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetMaximumExtraDataSize()
		}
	}

	return internal.One().GetMaximumExtraDataSize()
}

func (c *ChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetMaximumExtraDataSize(n)
		}
	}

	return internal.One().SetMaximumExtraDataSize(n)
}

func (c *ChainConfig) GetMinGasLimit() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetMinGasLimit()
		}
	}

	return internal.One().GetMinGasLimit()
}

func (c *ChainConfig) SetMinGasLimit(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetMinGasLimit(n)
		}
	}

	return internal.One().SetMinGasLimit(n)
}

func (c *ChainConfig) GetGasLimitBoundDivisor() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetGasLimitBoundDivisor()
		}
	}

	return internal.One().GetGasLimitBoundDivisor()
}

func (c *ChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetGasLimitBoundDivisor(n)
		}
	}

	return internal.One().SetGasLimitBoundDivisor(n)
}

// GetNetworkID and the following Set/Getters for ChainID too
// are... opinionated... because of where and how currently the NetworkID
// value is designed.
// This can cause unexpected and/or counter-intuitive behavior, especially with SetNetworkID.
// In order to use these logic properly, one should call NetworkID setter before ChainID setter.
// FIXME.
func (c *ChainConfig) GetNetworkID() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetNetworkID()
		}
	}
	if c.NetworkID != 0 {
		return &c.NetworkID
	}
	if c.ChainID != nil {
		return newU64(c.ChainID.Uint64())
	}
	return newU64(vars.DefaultNetworkID)
}

func (c *ChainConfig) SetNetworkID(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetNetworkID(n)
		}
	}

	if n == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	if c.ChainID == nil {
		c.ChainID = new(big.Int).SetUint64(*n)
	}
	c.NetworkID = *n
	return nil
}

func (c *ChainConfig) GetChainID() *big.Int {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetChainID()
		}
	}

	return c.ChainID
}

func (c *ChainConfig) SetChainID(n *big.Int) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetChainID(n)
		}
	}

	c.ChainID = n
	return nil
}

func (c *ChainConfig) GetMaxCodeSize() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetMaxCodeSize()
		}
	}

	return internal.One().GetMaxCodeSize()
}

func (c *ChainConfig) SetMaxCodeSize(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetMaxCodeSize(n)
		}
	}

	return internal.One().SetMaxCodeSize(n)
}

func (c *ChainConfig) GetEIP7Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP7Transition()
		}
	}

	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEIP7Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP7Transition(n)
		}
	}

	c.HomesteadBlock = setBig(c.HomesteadBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP150Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP150Transition()
		}
	}

	return bigNewU64(c.EIP150Block)
}

func (c *ChainConfig) SetEIP150Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP150Transition(n)
		}
	}

	c.EIP150Block = setBig(c.EIP150Block, n)
	return nil
}

func (c *ChainConfig) GetEIP152Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP152Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP152Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP152Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP160Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP160Transition()
		}
	}

	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP160Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP160Transition(n)
		}
	}

	c.EIP158Block = setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP161dTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP161dTransition()
		}
	}

	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP161dTransition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP161dTransition(n)
		}
	}

	c.EIP158Block = setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP161abcTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP161abcTransition()
		}
	}

	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP161abcTransition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP161abcTransition(n)
		}
	}

	c.EIP158Block = setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP170Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP170Transition()
		}
	}

	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP170Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP170Transition(n)
		}
	}

	c.EIP158Block = setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP155Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP155Transition()
		}
	}

	return bigNewU64(c.EIP155Block)
}

func (c *ChainConfig) SetEIP155Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP155Transition(n)
		}
	}

	c.EIP155Block = setBig(c.EIP155Block, n)
	return nil
}

func (c *ChainConfig) GetEIP140Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP140Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP140Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP140Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP198Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP198Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP198Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP198Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP211Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP211Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP211Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP211Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP212Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP212Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP212Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP212Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP213Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP213Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP213Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP213Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP214Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP214Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP214Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP214Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP658Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP658Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP658Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP658Transition(n)
		}
	}

	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP145Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP145Transition()
		}
	}

	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP145Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP145Transition(n)
		}
	}

	c.ConstantinopleBlock = setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1014Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1014Transition()
		}
	}

	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1014Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1014Transition(n)
		}
	}

	c.ConstantinopleBlock = setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1052Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1052Transition()
		}
	}

	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1052Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1052Transition(n)
		}
	}

	c.ConstantinopleBlock = setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1283Transition()
		}
	}

	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1283Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1283Transition(n)
		}
	}

	c.ConstantinopleBlock = setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283DisableTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1283DisableTransition()
		}
	}

	return bigNewU64(c.PetersburgBlock)
}

func (c *ChainConfig) SetEIP1283DisableTransition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1283DisableTransition(n)
		}
	}

	c.PetersburgBlock = setBig(c.PetersburgBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1108Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1108Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1108Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1108Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP2200Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP2200Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP2200Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP2200Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1344Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1344Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1344Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1344Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1884Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP1884Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1884Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP1884Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP2028Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEIP2028Transition()
		}
	}

	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP2028Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEIP2028Transition(n)
		}
	}

	c.IstanbulBlock = setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) IsForked(fn func() *uint64, n *big.Int) bool {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.IsForked(fn, n)
		}
	}

	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *ChainConfig) GetForkCanonHash(n uint64) common.Hash {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetForkCanonHash(n)
		}
	}

	if c.EIP150Block != nil && c.EIP150Block.Uint64() == n {
		return c.EIP150Hash
	}
	return common.Hash{}
}

func (c *ChainConfig) SetForkCanonHash(n uint64, h common.Hash) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetForkCanonHash(n, h)
		}
	}

	if c.GetEIP150Transition() != nil && *c.GetEIP150Transition() == n {
		c.EIP150Hash = h
		return nil
	}
	return ctypes.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetForkCanonHashes() map[uint64]common.Hash {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetForkCanonHashes()
		}
	}

	if c.EIP150Block == nil || c.EIP150Hash == (common.Hash{}) {
		return nil
	}
	return map[uint64]common.Hash{
		c.EIP150Block.Uint64(): c.EIP150Hash,
	}
}

func (c *ChainConfig) GetConsensusEngineType() ctypes.ConsensusEngineT {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetConsensusEngineType()
		}
	}

	if c.Ethash != nil {
		return ctypes.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return ctypes.ConsensusEngineT_Clique
	}
	return ctypes.ConsensusEngineT_Ethash
}

func (c *ChainConfig) MustSetConsensusEngineType(t ctypes.ConsensusEngineT) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.MustSetConsensusEngineType(t)
		}
	}

	switch t {
	case ctypes.ConsensusEngineT_Ethash:
		c.Ethash = new(ctypes.EthashConfig)
		return nil
	case ctypes.ConsensusEngineT_Clique:
		c.Clique = new(ctypes.CliqueConfig)
		return nil
	default:
		return ctypes.ErrUnsupportedConfigFatal
	}
}

func (c *ChainConfig) GetEthashMinimumDifficulty() *big.Int {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashMinimumDifficulty()
		}
	}

	return internal.One().GetEthashMinimumDifficulty()
}

func (c *ChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashMinimumDifficulty(i)
		}
	}

	return internal.One().SetEthashMinimumDifficulty(i)
}

func (c *ChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashDifficultyBoundDivisor()
		}
	}

	return internal.One().GetEthashDifficultyBoundDivisor()
}

func (c *ChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashDifficultyBoundDivisor(i)
		}
	}

	return internal.One().SetEthashDifficultyBoundDivisor(i)
}

func (c *ChainConfig) GetEthashDurationLimit() *big.Int {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashDurationLimit()
		}
	}

	return internal.One().GetEthashDurationLimit()
}

func (c *ChainConfig) SetEthashDurationLimit(i *big.Int) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashDurationLimit(i)
		}
	}

	return internal.One().SetEthashDurationLimit(i)
}

// NOTE: Checking for if c.Ethash == nil is a consideration.
// If set, settings are strictly enforced, and can avoid misconfiguration.
// If not, settings are more lenient, and allow for more shorthand testing.
// For the current implementation I have chosen to USE the nil check
// for Set_ methods, and to abstain for Get_ methods.
// This allows for shorthand-initialized structs, eg. for testing,
// but refuses un-strict Conversion methods.

func (c *ChainConfig) GetEthashHomesteadTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashHomesteadTransition()
		}
	}

	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEthashHomesteadTransition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashHomesteadTransition(i)
		}
	}

	c.HomesteadBlock = setBig(c.HomesteadBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashEIP2Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashEIP2Transition()
		}
	}

	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEthashEIP2Transition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashEIP2Transition(i)
		}
	}

	c.HomesteadBlock = setBig(c.HomesteadBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashEIP779Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashEIP779Transition()
		}
	}

	return bigNewU64(c.DAOForkBlock)
}

func (c *ChainConfig) SetEthashEIP779Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashEIP779Transition(n)
		}
	}

	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.DAOForkBlock = setBig(c.DAOForkBlock, n)
	if c.DAOForkBlock == nil {
		c.DAOForkSupport = false
	}
	return nil
}

func (c *ChainConfig) GetEthashEIP649Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashEIP649Transition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEthashEIP649Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashEIP649Transition(n)
		}
	}

	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ByzantiumBlock = setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashEIP1234Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashEIP1234Transition()
		}
	}

	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEthashEIP1234Transition(n *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashEIP1234Transition(n)
		}
	}

	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ConstantinopleBlock = setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashECIP1010PauseTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashECIP1010PauseTransition()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashECIP1010PauseTransition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashECIP1010PauseTransition(i)
		}
	}

	if i == nil {
		return nil
	}
	return ctypes.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1010ContinueTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashECIP1010ContinueTransition()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashECIP1010ContinueTransition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashECIP1010ContinueTransition(i)
		}
	}

	if i == nil {
		return nil
	}
	return ctypes.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashECIP1017Transition()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashECIP1017Transition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashECIP1017Transition(i)
		}
	}

	if i == nil {
		return nil
	}
	return ctypes.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017EraRounds() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashECIP1017EraRounds()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashECIP1017EraRounds(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashECIP1017EraRounds(i)
		}
	}

	if i == nil {
		return nil
	}
	return ctypes.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashEIP100BTransition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashEIP100BTransition()
		}
	}

	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEthashEIP100BTransition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashEIP100BTransition(i)
		}
	}

	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ByzantiumBlock = setBig(c.ByzantiumBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashECIP1041Transition() *uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashECIP1041Transition()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashECIP1041Transition(i *uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashECIP1041Transition(i)
		}
	}

	if i == nil {
		return nil
	}
	return ctypes.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashDifficultyBombDelaySchedule() ctypes.Uint64BigMapEncodesHex {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashDifficultyBombDelaySchedule()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashDifficultyBombDelaySchedule(m ctypes.Uint64BigMapEncodesHex) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashDifficultyBombDelaySchedule(m)
		}
	}

	return ctypes.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetEthashBlockRewardSchedule() ctypes.Uint64BigMapEncodesHex {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetEthashBlockRewardSchedule()
		}
	}

	return nil
}

func (c *ChainConfig) SetEthashBlockRewardSchedule(m ctypes.Uint64BigMapEncodesHex) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetEthashBlockRewardSchedule(m)
		}
	}

	return ctypes.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetCliquePeriod() uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetCliquePeriod()
		}
	}

	if c.Clique == nil {
		return 0
	}
	return c.Clique.Period
}

func (c *ChainConfig) SetCliquePeriod(n uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetCliquePeriod(n)
		}
	}

	if c.Clique == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.Clique.Period = n
	return nil
}

func (c *ChainConfig) GetCliqueEpoch() uint64 {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.GetCliqueEpoch()
		}
	}

	if c.Clique == nil {
		return 0
	}
	return c.Clique.Epoch
}

func (c *ChainConfig) SetCliqueEpoch(n uint64) error {
	if !c.Converting {
		if a := c.SwapIfAlt(); a != nil {
			return a.SetCliqueEpoch(n)
		}
	}

	if c.Clique == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.Clique.Epoch = n
	return nil
}
