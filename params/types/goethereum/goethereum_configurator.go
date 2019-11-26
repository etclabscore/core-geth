package goethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
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

func setBig(i *big.Int, u *uint64) {
	if u == nil {
		return
	}
	i = big.NewInt(int64(*u))
}

func (c *ChainConfig) GetAccountStartNonce() *uint64 {
	return newU64(0)
}

func (c *ChainConfig) SetAccountStartNonce(n *uint64) error {
	if *n != 0 {
		return common2.ErrUnsupportedConfigFatal
	}
	return nil
}

func (c *ChainConfig) GetMaximumExtraDataSize() *uint64 {
	return newU64(vars.MaximumExtraDataSize)
}

func (c *ChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	vars.MaximumExtraDataSize = *n
	return nil
}

func (c *ChainConfig) GetMinGasLimit() *uint64 {
	return newU64(vars.MinGasLimit)
}

func (c *ChainConfig) SetMinGasLimit(n *uint64) error {
	vars.MinGasLimit = *n
	return nil
}

func (c *ChainConfig) GetGasLimitBoundDivisor() *uint64 {
	return newU64(vars.GasLimitBoundDivisor)
}

func (c *ChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	vars.GasLimitBoundDivisor = *n
	return nil
}

func (c *ChainConfig) GetNetworkID() *uint64 {
	if c.ChainID != nil {
		return newU64(c.ChainID.Uint64())
	}
	return newU64(vars.DefaultNetworkID)
}

func (c *ChainConfig) SetNetworkID(n *uint64) error {
	if n == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	vars.DefaultNetworkID = *n
	return nil
}

func (c *ChainConfig) GetChainID() *uint64 {
	return bigNewU64(c.ChainID)
}

func (c *ChainConfig) SetChainID(n *uint64) error {
	if n == nil {
		return nil
	}
	c.ChainID = big.NewInt(int64(*n))
	return nil
}

func (c *ChainConfig) GetMaxCodeSize() *uint64 {
	return newU64(vars.MaxCodeSize)
}

func (c *ChainConfig) SetMaxCodeSize(n *uint64) error {
	if n == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	vars.MaxCodeSize = *n
	return nil
}

func (c *ChainConfig) GetEIP7Transition() *uint64 {
	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEIP7Transition(n *uint64) error {
	setBig(c.HomesteadBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP98Transition() *uint64 {
	return newU64(math.MaxUint64)
}

func (c *ChainConfig) SetEIP98Transition(n *uint64) error {
	return common2.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetEIP150Transition() *uint64 {
	return bigNewU64(c.EIP150Block)
}

func (c *ChainConfig) SetEIP150Transition(n *uint64) error {
	setBig(c.EIP150Block, n)
	return nil
}

func (c *ChainConfig) GetEIP152Transition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP152Transition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP160Transition() *uint64 {
	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP160Transition(n *uint64) error {
	setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP161abcTransition() *uint64 {
	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP161abcTransition(n *uint64) error {
	setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP161dTransition() *uint64 {
	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP161dTransition(n *uint64) error {
	setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP170Transition() *uint64 {
	return bigNewU64(c.EIP158Block)
}

func (c *ChainConfig) SetEIP170Transition(n *uint64) error {
	setBig(c.EIP158Block, n)
	return nil
}

func (c *ChainConfig) GetEIP155Transition() *uint64 {
	return bigNewU64(c.EIP155Block)
}

func (c *ChainConfig) SetEIP155Transition(n *uint64) error {
	setBig(c.EIP155Block, n)
	return nil
}

func (c *ChainConfig) GetEIP140Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP140Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP198Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP198Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP211Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP211Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP212Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP212Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP213Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP213Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP214Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP214Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP658Transition() *uint64 {
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEIP658Transition(n *uint64) error {
	setBig(c.ByzantiumBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP145Transition() *uint64 {
	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP145Transition(n *uint64) error {
	setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1014Transition() *uint64 {
	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1014Transition(n *uint64) error {
	setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1052Transition() *uint64 {
	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1052Transition(n *uint64) error {
	setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283Transition() *uint64 {
	return bigNewU64(c.ConstantinopleBlock)
}

func (c *ChainConfig) SetEIP1283Transition(n *uint64) error {
	setBig(c.ConstantinopleBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283DisableTransition() *uint64 {
	return bigNewU64(c.PetersburgBlock)
}

func (c *ChainConfig) SetEIP1283DisableTransition(n *uint64) error {
	setBig(c.PetersburgBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1108Transition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1108Transition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283ReenableTransition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1283ReenableTransition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1344Transition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1344Transition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1884Transition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP1884Transition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP2028Transition() *uint64 {
	return bigNewU64(c.IstanbulBlock)
}

func (c *ChainConfig) SetEIP2028Transition(n *uint64) error {
	setBig(c.IstanbulBlock, n)
	return nil
}

func (c *ChainConfig) IsForked(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *ChainConfig) GetConsensusEngineType() common2.ConsensusEngineT {
	if c.Ethash != nil {
		return common2.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return common2.ConsensusEngineT_Clique
	}
	return common2.ConsensusEngineT_Unknown
}

func (c *ChainConfig) MustSetConsensusEngineType(t common2.ConsensusEngineT) error {
	switch t {
	case common2.ConsensusEngineT_Ethash:
		c.Ethash = new(EthashConfig)
		return nil
	case common2.ConsensusEngineT_Clique:
		c.Clique = new(CliqueConfig)
		return nil
	default:
		return common2.ErrUnsupportedConfigFatal
	}
}

func (c *ChainConfig) GetEthashMinimumDifficulty() *big.Int {
	return vars.MinimumDifficulty
}

func (c *ChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	if i == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	vars.MinimumDifficulty = i
	return nil
}

func (c *ChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	return vars.DifficultyBoundDivisor
}

func (c *ChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	if i == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	vars.DifficultyBoundDivisor = i
	return nil
}

func (c *ChainConfig) GetEthashDurationLimit() *big.Int {
	return vars.DurationLimit
}

func (c *ChainConfig) SetEthashDurationLimit(i *big.Int) error {
	if i == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	vars.DurationLimit = i
	return nil
}

func (c *ChainConfig) GetEthashHomesteadTransition() *uint64 {
	if c.Ethash == nil {
		return nil
	}
	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEthashHomesteadTransition(i *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.HomesteadBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashEIP2Transition() *uint64 {
	if c.Ethash == nil {
		return nil
	}
	return bigNewU64(c.HomesteadBlock)
}

func (c *ChainConfig) SetEthashEIP2Transition(i *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.HomesteadBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashECIP1010PauseTransition() *uint64 {
	return nil
}

func (c *ChainConfig) SetEthashECIP1010PauseTransition(i *uint64) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1010ContinueTransition() *uint64 {
	return nil
}

func (c *ChainConfig) SetEthashECIP1010ContinueTransition(i *uint64) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017Transition() *uint64 {
	return nil
}

func (c *ChainConfig) SetEthashECIP1017Transition(i *uint64) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017EraRounds() *uint64 {
	return nil
}

func (c *ChainConfig) SetEthashECIP1017EraRounds(i *uint64) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashEIP100BTransition() *uint64 {
	if c.Ethash == nil {
		return nil
	}
	return bigNewU64(c.ByzantiumBlock)
}

func (c *ChainConfig) SetEthashEIP100BTransition(i *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.ByzantiumBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashECIP1041Transition() *uint64 {
	return nil
}

func (c *ChainConfig) SetEthashECIP1041Transition(i *uint64) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex {
	return nil
}

func (c *ChainConfig) SetEthashDifficultyBombDelaySchedule(m common2.Uint64BigMapEncodesHex) error {
	return common2.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex {
	return nil
}

func (c *ChainConfig) SetEthashBlockRewardSchedule(m common2.Uint64BigMapEncodesHex) error {
	return common2.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetCliquePeriod() *uint64 {
	if c.Clique == nil {
		return nil
	}
	return newU64(c.Clique.Period)
}

func (c *ChainConfig) SetCliquePeriod(n uint64) error {
	if c.Clique == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.Clique.Period = n
	return nil
}

func (c *ChainConfig) GetCliqueEpoch() *uint64 {
	if c.Clique == nil {
		return nil
	}
	return newU64(c.Clique.Epoch)
}

func (c *ChainConfig) SetCliqueEpoch(n uint64) error {
	if c.Clique == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.Clique.Epoch = n
	return nil
}
