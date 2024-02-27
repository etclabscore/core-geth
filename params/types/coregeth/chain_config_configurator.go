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

/*
This file contains logic implementing the Configurator interface for multi-geth.

Notes:
When setting the difficulty bomb delay map using a wanted total difficulty
value. The map, following Parity's format, uses aggregating values (summed) to yield a net difficulty delay,
while the specs use a max value (eg 3m, 5m, 9m). The model's difficulty bomb delay map data type has a
method SetValueForHeight which is used for this.
*/

package coregeth

import (
	"math/big"
	"reflect"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/internal"
	"github.com/ethereum/go-ethereum/params/vars"
)

func newU64(u uint64) *uint64 {
	return &u
}

func bigNewU64(i *big.Int) *uint64 {
	if i == nil {
		return nil
	}
	return newU64(i.Uint64())
}

// nolint: staticcheck
func setBig(i *big.Int, u *uint64) *big.Int {
	if u == nil {
		return nil
	}
	i = big.NewInt(int64(*u))
	return i
}

func (c *CoreGethChainConfig) ensureExistingRewardSchedule() {
	if c.BlockRewardSchedule == nil {
		c.BlockRewardSchedule = ctypes.Uint64Uint256MapEncodesHex{}
	}
}

func (c *CoreGethChainConfig) ensureExistingDifficultySchedule() {
	if c.DifficultyBombDelaySchedule == nil {
		c.DifficultyBombDelaySchedule = ctypes.Uint64Uint256MapEncodesHex{}
	}
}

func (c *CoreGethChainConfig) GetAccountStartNonce() *uint64 {
	return internal.GlobalConfigurator().GetAccountStartNonce()
}
func (c *CoreGethChainConfig) SetAccountStartNonce(n *uint64) error {
	return internal.GlobalConfigurator().SetAccountStartNonce(n)
}
func (c *CoreGethChainConfig) GetMaximumExtraDataSize() *uint64 {
	return internal.GlobalConfigurator().GetMaximumExtraDataSize()
}
func (c *CoreGethChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	return internal.GlobalConfigurator().SetMaximumExtraDataSize(n)
}
func (c *CoreGethChainConfig) GetMinGasLimit() *uint64 {
	return internal.GlobalConfigurator().GetMinGasLimit()
}
func (c *CoreGethChainConfig) SetMinGasLimit(n *uint64) error {
	return internal.GlobalConfigurator().SetMinGasLimit(n)
}
func (c *CoreGethChainConfig) GetGasLimitBoundDivisor() *uint64 {
	return internal.GlobalConfigurator().GetGasLimitBoundDivisor()
}
func (c *CoreGethChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	return internal.GlobalConfigurator().SetGasLimitBoundDivisor(n)
}

func (c *CoreGethChainConfig) GetNetworkID() *uint64 {
	return newU64(c.NetworkID)
}

func (c *CoreGethChainConfig) SetNetworkID(n *uint64) error {
	if n == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.NetworkID = *n
	return nil
}

func (c *CoreGethChainConfig) GetChainID() *big.Int {
	return c.ChainID
}

func (c *CoreGethChainConfig) SetChainID(n *big.Int) error {
	c.ChainID = n
	return nil
}

// GetSupportedProtocolVersions returns the protocol versions supported by this configuration value.
// When GetSupportedProtocolVersions is called, if the field containing the associated value (SupportedProtocolVersions)
// is empty, this method will assign the app-default value to that field.
// This establishes an in-data way of handling default behavior, and plays nicely with configurator equivalence
// and conversion methods.
func (c *CoreGethChainConfig) GetSupportedProtocolVersions() []uint {
	if len(c.SupportedProtocolVersions) == 0 {
		c.SupportedProtocolVersions = vars.DefaultProtocolVersions
	}
	return c.SupportedProtocolVersions
}

func (c *CoreGethChainConfig) SetSupportedProtocolVersions(p []uint) error {
	c.SupportedProtocolVersions = p
	return nil
}

func (c *CoreGethChainConfig) GetMaxCodeSize() *uint64 {
	return internal.GlobalConfigurator().GetMaxCodeSize()
}
func (c *CoreGethChainConfig) SetMaxCodeSize(n *uint64) error {
	return internal.GlobalConfigurator().SetMaxCodeSize(n)
}

func (c *CoreGethChainConfig) GetElasticityMultiplier() uint64 {
	return internal.GlobalConfigurator().GetElasticityMultiplier()
}

func (c *CoreGethChainConfig) SetElasticityMultiplier(n uint64) error {
	return internal.GlobalConfigurator().SetElasticityMultiplier(n)
}

func (c *CoreGethChainConfig) GetBaseFeeChangeDenominator() uint64 {
	return internal.GlobalConfigurator().GetBaseFeeChangeDenominator()
}

func (c *CoreGethChainConfig) SetBaseFeeChangeDenominator(n uint64) error {
	return internal.GlobalConfigurator().SetBaseFeeChangeDenominator(n)
}

func (c *CoreGethChainConfig) GetEIP7Transition() *uint64 {
	return bigNewU64(c.EIP7FBlock)
}

func (c *CoreGethChainConfig) SetEIP7Transition(n *uint64) error {
	c.EIP7FBlock = setBig(c.EIP7FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP150Transition() *uint64 {
	return bigNewU64(c.EIP150Block)
}

func (c *CoreGethChainConfig) SetEIP150Transition(n *uint64) error {
	c.EIP150Block = setBig(c.EIP150Block, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP152Transition() *uint64 {
	return bigNewU64(c.EIP152FBlock)
}

func (c *CoreGethChainConfig) SetEIP152Transition(n *uint64) error {
	c.EIP152FBlock = setBig(c.EIP152FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP160Transition() *uint64 {
	return bigNewU64(c.EIP160FBlock)
}

func (c *CoreGethChainConfig) SetEIP160Transition(n *uint64) error {
	c.EIP160FBlock = setBig(c.EIP160FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP161dTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *CoreGethChainConfig) SetEIP161dTransition(n *uint64) error {
	c.EIP161FBlock = setBig(c.EIP161FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP161abcTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *CoreGethChainConfig) SetEIP161abcTransition(n *uint64) error {
	c.EIP161FBlock = setBig(c.EIP161FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP170Transition() *uint64 {
	return bigNewU64(c.EIP170FBlock)
}

func (c *CoreGethChainConfig) SetEIP170Transition(n *uint64) error {
	c.EIP170FBlock = setBig(c.EIP170FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP155Transition() *uint64 {
	return bigNewU64(c.EIP155Block)
}

func (c *CoreGethChainConfig) SetEIP155Transition(n *uint64) error {
	c.EIP155Block = setBig(c.EIP155Block, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP140Transition() *uint64 {
	return bigNewU64(c.EIP140FBlock)
}

func (c *CoreGethChainConfig) SetEIP140Transition(n *uint64) error {
	c.EIP140FBlock = setBig(c.EIP140FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP198Transition() *uint64 {
	return bigNewU64(c.EIP198FBlock)
}

func (c *CoreGethChainConfig) SetEIP198Transition(n *uint64) error {
	c.EIP198FBlock = setBig(c.EIP198FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP211Transition() *uint64 {
	return bigNewU64(c.EIP211FBlock)
}

func (c *CoreGethChainConfig) SetEIP211Transition(n *uint64) error {
	c.EIP211FBlock = setBig(c.EIP211FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP212Transition() *uint64 {
	return bigNewU64(c.EIP212FBlock)
}

func (c *CoreGethChainConfig) SetEIP212Transition(n *uint64) error {
	c.EIP212FBlock = setBig(c.EIP212FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP213Transition() *uint64 {
	return bigNewU64(c.EIP213FBlock)
}

func (c *CoreGethChainConfig) SetEIP213Transition(n *uint64) error {
	c.EIP213FBlock = setBig(c.EIP213FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP214Transition() *uint64 {
	return bigNewU64(c.EIP214FBlock)
}

func (c *CoreGethChainConfig) SetEIP214Transition(n *uint64) error {
	c.EIP214FBlock = setBig(c.EIP214FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP658Transition() *uint64 {
	return bigNewU64(c.EIP658FBlock)
}

func (c *CoreGethChainConfig) SetEIP658Transition(n *uint64) error {
	c.EIP658FBlock = setBig(c.EIP658FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP145Transition() *uint64 {
	return bigNewU64(c.EIP145FBlock)
}

func (c *CoreGethChainConfig) SetEIP145Transition(n *uint64) error {
	c.EIP145FBlock = setBig(c.EIP145FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1014Transition() *uint64 {
	return bigNewU64(c.EIP1014FBlock)
}

func (c *CoreGethChainConfig) SetEIP1014Transition(n *uint64) error {
	c.EIP1014FBlock = setBig(c.EIP1014FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1052Transition() *uint64 {
	return bigNewU64(c.EIP1052FBlock)
}

func (c *CoreGethChainConfig) SetEIP1052Transition(n *uint64) error {
	c.EIP1052FBlock = setBig(c.EIP1052FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1283Transition() *uint64 {
	return bigNewU64(c.EIP1283FBlock)
}

func (c *CoreGethChainConfig) SetEIP1283Transition(n *uint64) error {
	c.EIP1283FBlock = setBig(c.EIP1283FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1283DisableTransition() *uint64 {
	return bigNewU64(c.PetersburgBlock)
}

func (c *CoreGethChainConfig) SetEIP1283DisableTransition(n *uint64) error {
	c.PetersburgBlock = setBig(c.PetersburgBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1108Transition() *uint64 {
	return bigNewU64(c.EIP1108FBlock)
}

func (c *CoreGethChainConfig) SetEIP1108Transition(n *uint64) error {
	c.EIP1108FBlock = setBig(c.EIP1108FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2200Transition() *uint64 {
	return bigNewU64(c.EIP2200FBlock)
}

func (c *CoreGethChainConfig) SetEIP2200Transition(n *uint64) error {
	c.EIP2200FBlock = setBig(c.EIP2200FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2200DisableTransition() *uint64 {
	return bigNewU64(c.EIP2200DisableFBlock)
}

func (c *CoreGethChainConfig) SetEIP2200DisableTransition(n *uint64) error {
	c.EIP2200DisableFBlock = setBig(c.EIP2200DisableFBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1344Transition() *uint64 {
	return bigNewU64(c.EIP1344FBlock)
}

func (c *CoreGethChainConfig) SetEIP1344Transition(n *uint64) error {
	c.EIP1344FBlock = setBig(c.EIP1344FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1884Transition() *uint64 {
	return bigNewU64(c.EIP1884FBlock)
}

func (c *CoreGethChainConfig) SetEIP1884Transition(n *uint64) error {
	c.EIP1884FBlock = setBig(c.EIP1884FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2028Transition() *uint64 {
	return bigNewU64(c.EIP2028FBlock)
}

func (c *CoreGethChainConfig) SetEIP2028Transition(n *uint64) error {
	c.EIP2028FBlock = setBig(c.EIP2028FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetECIP1080Transition() *uint64 {
	return bigNewU64(c.ECIP1080FBlock)
}

func (c *CoreGethChainConfig) SetECIP1080Transition(n *uint64) error {
	c.ECIP1080FBlock = setBig(c.ECIP1080FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1706Transition() *uint64 {
	return bigNewU64(c.EIP1706FBlock)
}

func (c *CoreGethChainConfig) SetEIP1706Transition(n *uint64) error {
	c.EIP1706FBlock = setBig(c.EIP1706FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2537Transition() *uint64 {
	return bigNewU64(c.EIP2537FBlock)
}

func (c *CoreGethChainConfig) SetEIP2537Transition(n *uint64) error {
	c.EIP2537FBlock = setBig(c.EIP2537FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetECBP1100Transition() *uint64 {
	return bigNewU64(c.ECBP1100FBlock)
}

func (c *CoreGethChainConfig) SetECBP1100Transition(n *uint64) error {
	c.ECBP1100FBlock = setBig(c.ECBP1100FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetECBP1100DeactivateTransition() *uint64 {
	return bigNewU64(c.ECBP1100DeactivateFBlock)
}

func (c *CoreGethChainConfig) SetECBP1100DeactivateTransition(n *uint64) error {
	c.ECBP1100DeactivateFBlock = setBig(c.ECBP1100DeactivateFBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2315Transition() *uint64 {
	return bigNewU64(c.EIP2315FBlock)
}

func (c *CoreGethChainConfig) SetEIP2315Transition(n *uint64) error {
	c.EIP2315FBlock = setBig(c.EIP2315FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2929Transition() *uint64 {
	return bigNewU64(c.EIP2929FBlock)
}

func (c *CoreGethChainConfig) SetEIP2929Transition(n *uint64) error {
	c.EIP2929FBlock = setBig(c.EIP2929FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2930Transition() *uint64 {
	return bigNewU64(c.EIP2930FBlock)
}

func (c *CoreGethChainConfig) SetEIP2930Transition(n *uint64) error {
	c.EIP2930FBlock = setBig(c.EIP2930FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2565Transition() *uint64 {
	return bigNewU64(c.EIP2565FBlock)
}

func (c *CoreGethChainConfig) SetEIP2565Transition(n *uint64) error {
	c.EIP2565FBlock = setBig(c.EIP2565FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2718Transition() *uint64 {
	return bigNewU64(c.EIP2718FBlock)
}

func (c *CoreGethChainConfig) SetEIP2718Transition(n *uint64) error {
	c.EIP2718FBlock = setBig(c.EIP2718FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP1559Transition() *uint64 {
	return bigNewU64(c.EIP1559FBlock)
}

func (c *CoreGethChainConfig) SetEIP1559Transition(n *uint64) error {
	c.EIP1559FBlock = setBig(c.EIP1559FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP3541Transition() *uint64 {
	return bigNewU64(c.EIP3541FBlock)
}

func (c *CoreGethChainConfig) SetEIP3541Transition(n *uint64) error {
	c.EIP3541FBlock = setBig(c.EIP3541FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP3529Transition() *uint64 {
	return bigNewU64(c.EIP3529FBlock)
}

func (c *CoreGethChainConfig) SetEIP3529Transition(n *uint64) error {
	c.EIP3529FBlock = setBig(c.EIP3529FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP3198Transition() *uint64 {
	return bigNewU64(c.EIP3198FBlock)
}

func (c *CoreGethChainConfig) SetEIP3198Transition(n *uint64) error {
	c.EIP3198FBlock = setBig(c.EIP3198FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP4399Transition() *uint64 {
	return bigNewU64(c.EIP4399FBlock)
}

func (c *CoreGethChainConfig) SetEIP4399Transition(n *uint64) error {
	c.EIP4399FBlock = setBig(c.EIP4399FBlock, n)
	return nil
}

// EIP3651: Warm COINBASE
func (c *CoreGethChainConfig) GetEIP3651TransitionTime() *uint64 {
	return c.EIP3651FTime
}

func (c *CoreGethChainConfig) SetEIP3651TransitionTime(n *uint64) error {
	c.EIP3651FTime = n
	return nil
}

// GetEIP3855TransitionTime EIP3855: PUSH0 instruction
func (c *CoreGethChainConfig) GetEIP3855TransitionTime() *uint64 {
	return c.EIP3855FTime
}

func (c *CoreGethChainConfig) SetEIP3855TransitionTime(n *uint64) error {
	c.EIP3855FTime = n
	return nil
}

// GetEIP3860TransitionTime EIP3860: Limit and meter initcode
func (c *CoreGethChainConfig) GetEIP3860TransitionTime() *uint64 {
	return c.EIP3860FTime
}

func (c *CoreGethChainConfig) SetEIP3860TransitionTime(n *uint64) error {
	c.EIP3860FTime = n
	return nil
}

// GetEIP4895TransitionTime EIP4895: Beacon chain push withdrawals as operations
func (c *CoreGethChainConfig) GetEIP4895TransitionTime() *uint64 {
	return c.EIP4895FTime
}

func (c *CoreGethChainConfig) SetEIP4895TransitionTime(n *uint64) error {
	c.EIP4895FTime = n
	return nil
}

// GetEIP6049TransitionTime EIP6049: Deprecate SELFDESTRUCT
func (c *CoreGethChainConfig) GetEIP6049TransitionTime() *uint64 {
	return c.EIP6049FTime
}

func (c *CoreGethChainConfig) SetEIP6049TransitionTime(n *uint64) error {
	c.EIP6049FTime = n
	return nil
}

// Shanghai by block
// EIP3651: Warm COINBASE
func (c *CoreGethChainConfig) GetEIP3651Transition() *uint64 {
	return bigNewU64(c.EIP3651FBlock)
}

func (c *CoreGethChainConfig) SetEIP3651Transition(n *uint64) error {
	c.EIP3651FBlock = setBig(c.EIP3651FBlock, n)
	return nil
}

// GetEIP3855Transition EIP3855: PUSH0 instruction
func (c *CoreGethChainConfig) GetEIP3855Transition() *uint64 {
	return bigNewU64(c.EIP3855FBlock)
}

func (c *CoreGethChainConfig) SetEIP3855Transition(n *uint64) error {
	c.EIP3855FBlock = setBig(c.EIP3855FBlock, n)
	return nil
}

// GetEIP3860Transition EIP3860: Limit and meter initcode
func (c *CoreGethChainConfig) GetEIP3860Transition() *uint64 {
	return bigNewU64(c.EIP3860FBlock)
}

func (c *CoreGethChainConfig) SetEIP3860Transition(n *uint64) error {
	c.EIP3860FBlock = setBig(c.EIP3860FBlock, n)
	return nil
}

// GetEIP4895Transition EIP4895: Beacon chain push withdrawals as operations
func (c *CoreGethChainConfig) GetEIP4895Transition() *uint64 {
	return bigNewU64(c.EIP4895FBlock)
}

func (c *CoreGethChainConfig) SetEIP4895Transition(n *uint64) error {
	c.EIP4895FBlock = setBig(c.EIP4895FBlock, n)
	return nil
}

// GetEIP6049Transition EIP6049: Deprecate SELFDESTRUCT
func (c *CoreGethChainConfig) GetEIP6049Transition() *uint64 {
	return bigNewU64(c.EIP6049FBlock)
}

func (c *CoreGethChainConfig) SetEIP6049Transition(n *uint64) error {
	c.EIP6049FBlock = setBig(c.EIP6049FBlock, n)
	return nil
}

// GetEIP4844TransitionTime EIP4844: Shard Blob Transactions
func (c *CoreGethChainConfig) GetEIP4844TransitionTime() *uint64 {
	return c.EIP4844FTime
}

func (c *CoreGethChainConfig) SetEIP4844TransitionTime(n *uint64) error {
	c.EIP4844FTime = n
	return nil
}

// GetEIP7516TransitionTime EIP7516: Blob Base Fee Opcode
func (c *CoreGethChainConfig) GetEIP7516TransitionTime() *uint64 {
	return c.EIP7516FTime
}

func (c *CoreGethChainConfig) SetEIP7516TransitionTime(n *uint64) error {
	c.EIP7516FTime = n
	return nil
}

// GetEIP1153TransitionTime EIP1153: Transient Storage opcodes
func (c *CoreGethChainConfig) GetEIP1153TransitionTime() *uint64 {
	return c.EIP1153FTime
}

func (c *CoreGethChainConfig) SetEIP1153TransitionTime(n *uint64) error {
	c.EIP1153FTime = n
	return nil
}

// GetEIP5656TransitionTime EIP5656: MCOPY - Memory copying instruction
func (c *CoreGethChainConfig) GetEIP5656TransitionTime() *uint64 {
	return c.EIP5656FTime
}

func (c *CoreGethChainConfig) SetEIP5656TransitionTime(n *uint64) error {
	c.EIP5656FTime = n
	return nil
}

// GetEIP6780TransitionTime EIP6780: SELFDESTRUCT only in same transaction
func (c *CoreGethChainConfig) GetEIP6780TransitionTime() *uint64 {
	return c.EIP6780FTime
}

func (c *CoreGethChainConfig) SetEIP6780TransitionTime(n *uint64) error {
	c.EIP6780FTime = n
	return nil
}

// GetEIP6780TransitionTime EIP4788: Beacon block root in the EVM
func (c *CoreGethChainConfig) GetEIP4788TransitionTime() *uint64 {
	return c.EIP4788FTime
}

func (c *CoreGethChainConfig) SetEIP4788TransitionTime(n *uint64) error {
	c.EIP4788FTime = n
	return nil
}

// Cancun by block
// GetEIP4844Transition EIP4844: Shard Blob Transactions
func (c *CoreGethChainConfig) GetEIP4844Transition() *uint64 {
	return bigNewU64(c.EIP4844FBlock)
}

func (c *CoreGethChainConfig) SetEIP4844Transition(n *uint64) error {
	c.EIP4844FBlock = setBig(c.EIP4844FBlock, n)
	return nil
}

// GetEIP7516Transition EIP7516: Blob Base Fee Opcode
func (c *CoreGethChainConfig) GetEIP7516Transition() *uint64 {
	return bigNewU64(c.EIP7516FBlock)
}

func (c *CoreGethChainConfig) SetEIP7516Transition(n *uint64) error {
	c.EIP7516FBlock = setBig(c.EIP7516FBlock, n)
	return nil
}

// GetEIP1153Transition EIP1153: Transient Storage opcodes
func (c *CoreGethChainConfig) GetEIP1153Transition() *uint64 {
	return bigNewU64(c.EIP1153FBlock)
}

func (c *CoreGethChainConfig) SetEIP1153Transition(n *uint64) error {
	c.EIP1153FBlock = setBig(c.EIP1153FBlock, n)
	return nil
}

// GetEIP5656Transition EIP5656: MCOPY - Memory copying instruction
func (c *CoreGethChainConfig) GetEIP5656Transition() *uint64 {
	return bigNewU64(c.EIP5656FBlock)
}

func (c *CoreGethChainConfig) SetEIP5656Transition(n *uint64) error {
	c.EIP5656FBlock = setBig(c.EIP5656FBlock, n)
	return nil
}

// GetEIP6780Transition EIP6780: SELFDESTRUCT only in same transaction
func (c *CoreGethChainConfig) GetEIP6780Transition() *uint64 {
	return bigNewU64(c.EIP6780FBlock)
}

func (c *CoreGethChainConfig) SetEIP6780Transition(n *uint64) error {
	c.EIP6780FBlock = setBig(c.EIP6780FBlock, n)
	return nil
}

// GetEIP6780Transition EIP4788: Beacon block root in the EVM
func (c *CoreGethChainConfig) GetEIP4788Transition() *uint64 {
	return bigNewU64(c.EIP4788FBlock)
}

func (c *CoreGethChainConfig) SetEIP4788Transition(n *uint64) error {
	c.EIP4788FBlock = setBig(c.EIP4788FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetMergeVirtualTransition() *uint64 {
	return bigNewU64(c.MergeNetsplitVBlock)
}

func (c *CoreGethChainConfig) SetMergeVirtualTransition(n *uint64) error {
	c.MergeNetsplitVBlock = setBig(c.MergeNetsplitVBlock, n)
	return nil
}

// Verkle Trie
func (c *CoreGethChainConfig) GetVerkleTransitionTime() *uint64 {
	return c.VerkleFTime
}

func (c *CoreGethChainConfig) SetVerkleTransitionTime(n *uint64) error {
	c.VerkleFTime = n
	return nil
}

func (c *CoreGethChainConfig) GetVerkleTransition() *uint64 {
	return bigNewU64(c.VerkleFBlock)
}

func (c *CoreGethChainConfig) SetVerkleTransition(n *uint64) error {
	c.VerkleFBlock = setBig(c.VerkleFBlock, n)
	return nil
}

func (c *CoreGethChainConfig) IsEnabled(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if strings.Contains(fnName, "ECBP1100Transition") {
		deactivateTransition := c.GetECBP1100DeactivateTransition()
		if deactivateTransition != nil {
			return big.NewInt(int64(*deactivateTransition)).Cmp(n) > 0 && big.NewInt(int64(*f)).Cmp(n) <= 0
		}
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *CoreGethChainConfig) IsEnabledByTime(fn func() *uint64, n *uint64) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return *f <= *n
}

func (c *CoreGethChainConfig) GetForkCanonHash(n uint64) common.Hash {
	if c.RequireBlockHashes == nil {
		return common.Hash{}
	}
	for k, v := range c.RequireBlockHashes {
		if k == n {
			return v
		}
	}
	return common.Hash{}
}

func (c *CoreGethChainConfig) SetForkCanonHash(n uint64, h common.Hash) error {
	if c.RequireBlockHashes == nil {
		c.RequireBlockHashes = make(map[uint64]common.Hash)
	}
	c.RequireBlockHashes[n] = h
	return nil
}

func (c *CoreGethChainConfig) GetForkCanonHashes() map[uint64]common.Hash {
	return c.RequireBlockHashes
}

func (c *CoreGethChainConfig) GetConsensusEngineType() ctypes.ConsensusEngineT {
	if c.Ethash != nil {
		return ctypes.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return ctypes.ConsensusEngineT_Clique
	}
	if c.Lyra2 != nil {
		return ctypes.ConsensusEngineT_Lyra2
	}
	return ctypes.ConsensusEngineT_Unknown
}

func (c *CoreGethChainConfig) MustSetConsensusEngineType(t ctypes.ConsensusEngineT) error {
	switch t {
	case ctypes.ConsensusEngineT_Ethash:
		c.Ethash = new(ctypes.EthashConfig)
		c.Clique = nil
		return nil
	case ctypes.ConsensusEngineT_Clique:
		c.Clique = new(ctypes.CliqueConfig)
		c.Ethash = nil
		return nil
	case ctypes.ConsensusEngineT_Lyra2:
		c.Lyra2 = new(ctypes.Lyra2Config)
		c.Ethash = nil
		c.Clique = nil
		return nil
	default:
		return ctypes.ErrUnsupportedConfigFatal
	}
}

func (c *CoreGethChainConfig) GetIsDevMode() bool {
	return c.IsDevMode
}

func (c *CoreGethChainConfig) SetDevMode(devMode bool) error {
	c.IsDevMode = devMode
	return nil
}

func (c *CoreGethChainConfig) GetEthashTerminalTotalDifficulty() *big.Int {
	return c.TerminalTotalDifficulty
}

func (c *CoreGethChainConfig) SetEthashTerminalTotalDifficulty(n *big.Int) error {
	c.TerminalTotalDifficulty = n
	return nil
}

func (c *CoreGethChainConfig) GetEthashTerminalTotalDifficultyPassed() bool {
	return c.TerminalTotalDifficultyPassed
}

func (c *CoreGethChainConfig) SetEthashTerminalTotalDifficultyPassed(t bool) error {
	c.TerminalTotalDifficultyPassed = t
	return nil
}

// IsTerminalPoWBlock returns whether the given block is the last block of PoW stage.
func (c *CoreGethChainConfig) IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool {
	terminalTotalDifficulty := c.GetEthashTerminalTotalDifficulty()
	if terminalTotalDifficulty == nil {
		return false
	}
	return parentTotalDiff.Cmp(terminalTotalDifficulty) < 0 && totalDiff.Cmp(terminalTotalDifficulty) >= 0
}

func (c *CoreGethChainConfig) GetEthashMinimumDifficulty() *big.Int {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return internal.GlobalConfigurator().GetEthashMinimumDifficulty()
}
func (c *CoreGethChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	return internal.GlobalConfigurator().SetEthashMinimumDifficulty(i)
}

func (c *CoreGethChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return internal.GlobalConfigurator().GetEthashDifficultyBoundDivisor()
}

func (c *CoreGethChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	return internal.GlobalConfigurator().SetEthashDifficultyBoundDivisor(i)
}

func (c *CoreGethChainConfig) GetEthashDurationLimit() *big.Int {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return internal.GlobalConfigurator().GetEthashDurationLimit()
}

func (c *CoreGethChainConfig) SetEthashDurationLimit(i *big.Int) error {
	return internal.GlobalConfigurator().SetEthashDurationLimit(i)
}

func (c *CoreGethChainConfig) GetEthashHomesteadTransition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.EIP2FBlock == nil || c.EIP7FBlock == nil {
		return nil
	}
	return bigNewU64(math.BigMax(c.EIP2FBlock, c.EIP7FBlock))
}

func (c *CoreGethChainConfig) SetEthashHomesteadTransition(n *uint64) error {
	c.EIP2FBlock = setBig(c.EIP2FBlock, n)
	c.EIP7FBlock = setBig(c.EIP7FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEIP2Transition() *uint64 {
	return bigNewU64(c.EIP2FBlock)
}

func (c *CoreGethChainConfig) SetEIP2Transition(n *uint64) error {
	c.EIP2FBlock = setBig(c.EIP2FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP779Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.DAOForkBlock)
}

func (c *CoreGethChainConfig) SetEthashEIP779Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.DAOForkBlock = setBig(c.DAOForkBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP649Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip649FInferred {
		return bigNewU64(c.EIP649FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP649FBlock = setBig(c.EIP649FBlock, diffN)
		c.eip649FInferred = true
	}()

	// Get block number (key) from maps where EIP649 criteria is met.
	diffN = ctypes.MapMeetsSpecification(
		c.DifficultyBombDelaySchedule,
		c.BlockRewardSchedule,
		vars.EIP649DifficultyBombDelay,
		vars.EIP649FBlockReward,
	)
	if diffN == nil {
		diffN = c.GetEthashEIP1234Transition()
	}
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP649Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP649FBlock = setBig(c.EIP649FBlock, n)
	c.eip649FInferred = true

	if n == nil {
		return nil
	}

	if eip1234 := c.GetEthashEIP1234Transition(); eip1234 != nil {
		if *eip1234 <= *n {
			return nil
		}
	}

	c.ensureExistingRewardSchedule()
	c.BlockRewardSchedule[*n] = vars.EIP649FBlockReward

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP649DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP1234Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip1234FInferred {
		return bigNewU64(c.EIP1234FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP1234FBlock = setBig(c.EIP1234FBlock, diffN)
		c.eip1234FInferred = true
	}()

	// Get block number (key) from maps where EIP1234 criteria is met.
	diffN = ctypes.MapMeetsSpecification(
		c.DifficultyBombDelaySchedule,
		c.BlockRewardSchedule,
		vars.EIP1234DifficultyBombDelay,
		vars.EIP1234FBlockReward,
	)
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP1234Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP1234FBlock = setBig(c.EIP1234FBlock, n)
	c.eip1234FInferred = true

	if n == nil {
		return nil
	}

	// Block reward is a simple lookup; doesn't matter if overwrite or not.
	c.ensureExistingRewardSchedule()
	c.BlockRewardSchedule[*n] = vars.EIP1234FBlockReward

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP1234DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP2384Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip2384Inferred {
		return bigNewU64(c.EIP2384FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP2384FBlock = setBig(c.EIP2384FBlock, diffN)
		c.eip2384Inferred = true
	}()

	// Get block number (key) from map where EIP2384 criteria is met.
	diffN = ctypes.MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, vars.EIP2384DifficultyBombDelay, nil)
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP2384Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP2384FBlock = setBig(c.EIP2384FBlock, n)
	c.eip2384Inferred = true

	if n == nil {
		return nil
	}

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP2384DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP3554Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip3554Inferred {
		return bigNewU64(c.EIP3554FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP3554FBlock = setBig(c.EIP3554FBlock, diffN)
		c.eip3554Inferred = true
	}()

	// Get block number (key) from map where EIP3554 criteria is met.
	diffN = ctypes.MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, vars.EIP3554DifficultyBombDelay, nil)
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP3554Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP3554FBlock = setBig(c.EIP3554FBlock, n)
	c.eip3554Inferred = true

	if n == nil {
		return nil
	}

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP3554DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP4345Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip4345Inferred {
		return bigNewU64(c.EIP4345FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP4345FBlock = setBig(c.EIP4345FBlock, diffN)
		c.eip4345Inferred = true
	}()

	// Get block number (key) from map where EIP4345 criteria is met.
	diffN = ctypes.MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, vars.EIP4345DifficultyBombDelay, nil)
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP4345Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP4345FBlock = setBig(c.EIP4345FBlock, n)
	c.eip4345Inferred = true

	if n == nil {
		return nil
	}

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP4345DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1010PauseTransition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.ECIP1010PauseBlock)
}

func (c *CoreGethChainConfig) SetEthashECIP1010PauseTransition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	if c.ECIP1010PauseBlock == nil && c.ECIP1010Length != nil {
		c.ECIP1010PauseBlock = setBig(c.ECIP1010PauseBlock, n)
		c.ECIP1010Length = c.ECIP1010Length.Sub(c.ECIP1010Length, c.ECIP1010PauseBlock)
		return nil
	}
	c.ECIP1010PauseBlock = setBig(c.ECIP1010PauseBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1010ContinueTransition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.ECIP1010PauseBlock == nil {
		return nil
	}
	if c.ECIP1010Length == nil {
		return nil
	}
	// transition = pause + length
	return bigNewU64(new(big.Int).Add(c.ECIP1010PauseBlock, c.ECIP1010Length))
}

func (c *CoreGethChainConfig) SetEthashECIP1010ContinueTransition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	// length = continue - pause
	if n == nil {
		return ctypes.ErrUnsupportedConfigNoop
	}
	if c.ECIP1010PauseBlock == nil {
		c.ECIP1010Length = new(big.Int).SetUint64(*n)
		return nil
	}
	c.ECIP1010Length = new(big.Int).Sub(big.NewInt(int64(*n)), c.ECIP1010PauseBlock)
	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1017Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.ECIP1017FBlock)
}

func (c *CoreGethChainConfig) SetEthashECIP1017Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ECIP1017FBlock = setBig(c.ECIP1017FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1017EraRounds() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.ECIP1017EraRounds)
}

func (c *CoreGethChainConfig) SetEthashECIP1017EraRounds(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ECIP1017EraRounds = setBig(c.ECIP1017EraRounds, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP100BTransition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.EIP100FBlock)
}

func (c *CoreGethChainConfig) SetEthashEIP100BTransition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.EIP100FBlock = setBig(c.EIP100FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1041Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.DisposalBlock)
}

func (c *CoreGethChainConfig) SetEthashECIP1041Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.DisposalBlock = setBig(c.DisposalBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashECIP1099Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return bigNewU64(c.ECIP1099FBlock)
}

func (c *CoreGethChainConfig) SetEthashECIP1099Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.ECIP1099FBlock = setBig(c.ECIP1099FBlock, n)
	return nil
}

func (c *CoreGethChainConfig) GetEthashEIP5133Transition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	if c.eip5133Inferred {
		return bigNewU64(c.EIP5133FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP5133FBlock = setBig(c.EIP5133FBlock, diffN)
		c.eip5133Inferred = true
	}()

	// Get block number (key) from map where EIP5133 criteria is met.
	diffN = ctypes.MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, vars.EIP5133DifficultyBombDelay, nil)
	return diffN
}

func (c *CoreGethChainConfig) SetEthashEIP5133Transition(n *uint64) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.EIP5133FBlock = setBig(c.EIP5133FBlock, n)
	c.eip5133Inferred = true

	if n == nil {
		return nil
	}

	c.ensureExistingDifficultySchedule()
	c.DifficultyBombDelaySchedule.SetValueTotalForHeight(n, vars.EIP5133DifficultyBombDelay)

	return nil
}

func (c *CoreGethChainConfig) GetEthashDifficultyBombDelaySchedule() ctypes.Uint64Uint256MapEncodesHex {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return c.DifficultyBombDelaySchedule
}

func (c *CoreGethChainConfig) SetEthashDifficultyBombDelaySchedule(m ctypes.Uint64Uint256MapEncodesHex) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.DifficultyBombDelaySchedule = m
	return nil
}

func (c *CoreGethChainConfig) GetEthashBlockRewardSchedule() ctypes.Uint64Uint256MapEncodesHex {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Ethash {
		return nil
	}
	return c.BlockRewardSchedule
}

func (c *CoreGethChainConfig) SetEthashBlockRewardSchedule(m ctypes.Uint64Uint256MapEncodesHex) error {
	if c.Ethash == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.BlockRewardSchedule = m
	return nil
}

func (c *CoreGethChainConfig) GetCliquePeriod() uint64 {
	if c.Clique == nil {
		return 0
	}
	return c.Clique.Period
}

func (c *CoreGethChainConfig) SetCliquePeriod(n uint64) error {
	if c.Clique == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.Clique.Period = n
	return nil
}

func (c *CoreGethChainConfig) GetCliqueEpoch() uint64 {
	if c.Clique == nil {
		return 0
	}
	return c.Clique.Epoch
}

func (c *CoreGethChainConfig) SetCliqueEpoch(n uint64) error {
	if c.Clique == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}
	c.Clique.Epoch = n
	return nil
}

func (c *CoreGethChainConfig) GetLyra2NonceTransition() *uint64 {
	if c.GetConsensusEngineType() != ctypes.ConsensusEngineT_Lyra2 {
		return nil
	}
	return bigNewU64(c.Lyra2NonceTransitionBlock)
}

func (c *CoreGethChainConfig) SetLyra2NonceTransition(n *uint64) error {
	if c.Lyra2 == nil {
		return ctypes.ErrUnsupportedConfigFatal
	}

	c.Lyra2NonceTransitionBlock = setBig(c.Lyra2NonceTransitionBlock, n)

	return nil
}
