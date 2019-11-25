package goethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
)

// File contains the go-ethereum implementation of the ChainConfigurator interface.
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
	return newU64(params.MaximumExtraDataSize)
}

func (c *ChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	params.MaximumExtraDataSize = *n
	return nil
}

func (c *ChainConfig) GetMinGasLimit() *uint64 {
	return newU64(params.MinGasLimit)
}

func (c *ChainConfig) SetMinGasLimit(n *uint64) error {
	params.MinGasLimit = *n
	return nil
}

func (c *ChainConfig) GetGasLimitBoundDivisor() *uint64 {
	return newU64(params.GasLimitBoundDivisor)
}

func (c *ChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	params.GasLimitBoundDivisor = *n
	return nil
}

func (c *ChainConfig) GetNetworkID() *uint64 {
	if c.ChainID != nil {
		return newU64(c.ChainID.Uint64())
	}
	return newU64(params.DefaultNetworkID)
}

func (c *ChainConfig) SetNetworkID(n *uint64) error {
	params.DefaultNetworkID = *n
	return nil
}

func (c *ChainConfig) GetChainID() *uint64 {
	return newU64(c.ChainID.Uint64())
}

func (c *ChainConfig) SetChainID(n *uint64) error {
	c.ChainID = big.NewInt(int64(*n))
	return nil
}

func (c *ChainConfig) GetMaxCodeSize() *uint64 {
	return newU64(params.MaxCodeSize)
}

func (c *ChainConfig) SetMaxCodeSize(n *uint64) error {
	params.MaxCodeSize = *n
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

func (c *ChainConfig) IsForked(fn func(*big.Int) bool, n *big.Int) bool {
	if n == nil || fn == nil {
		return false
	}
	return fn(n)
}

func (c *ChainConfig) GetConsensusEngineType() paramtypes.ConsensusEngineT {
	if c.Ethash != nil {
		return paramtypes.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return paramtypes.ConsensusEngineT_Clique
	}
	return paramtypes.ConsensusEngineT_Unknown
}

func (c *ChainConfig) MustSetConsensusEngineType(t paramtypes.ConsensusEngineT) error {
	switch t {
	case paramtypes.ConsensusEngineT_Ethash:
		c.Ethash = new(EthashConfig)
		return nil
	case paramtypes.ConsensusEngineT_Clique:
		c.Clique = new(CliqueConfig)
		return nil
	default:
		return common2.ErrUnsupportedConfigFatal
	}
}

func (c *ChainConfig) GetEthashMinimumDifficulty() *big.Int {
	return params.MinimumDifficulty
}

func (c *ChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	params.MinimumDifficulty = i
	return nil
}

func (c *ChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	return params.DifficultyBoundDivisor
}

func (c *ChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	params.DifficultyBoundDivisor = i
	return nil
}

func (c *ChainConfig) GetEthashDurationLimit() *big.Int {
	return params.DurationLimit
}

func (c *ChainConfig) SetEthashDurationLimit(i *big.Int) error {
	params.DurationLimit = i
	return nil
}

func (c *ChainConfig) GetEthashHomesteadTransition() *big.Int {
	return c.HomesteadBlock
}

func (c *ChainConfig) SetEthashHomesteadTransition(i *big.Int) error {
	c.HomesteadBlock = i
	return nil
}

func (c *ChainConfig) GetEthashEIP2Transition() *big.Int {
	return c.HomesteadBlock
}

func (c *ChainConfig) SetEthashEIP2Transition(i *big.Int) error {
	c.HomesteadBlock = i
	return nil
}

func (c *ChainConfig) GetEthashECIP1010PauseTransition() *big.Int {
	return nil
}

func (c *ChainConfig) SetEthashECIP1010PauseTransition(i *big.Int) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1010ContinueTransition() *big.Int {
	return nil
}

func (c *ChainConfig) SetEthashECIP1010ContinueTransition(i *big.Int) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017Transition() *big.Int {
	return nil
}

func (c *ChainConfig) SetEthashECIP1017Transition(i *big.Int) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashECIP1017EraRounds() *big.Int {
	return nil
}

func (c *ChainConfig) SetEthashECIP1017EraRounds(i *big.Int) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashEIP100BTransition() *big.Int {
	return c.ByzantiumBlock
}

func (c *ChainConfig) SetEthashEIP100BTransition(i *big.Int) error {
	setBig(c.ByzantiumBlock, i)
	return nil
}

func (c *ChainConfig) GetEthashECIP1041Transition() *big.Int {
	return nil
}

func (c *ChainConfig) SetEthashECIP1041Transition(i *big.Int) error {
	if i == nil {
		return nil
	}
	return common2.ErrUnsupportedConfigFatal
}

func (c *ChainConfig) GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex {
	m := common2.Uint64BigMapEncodesHex{
		0: params.FrontierBlockReward,
	}
	if c.ByzantiumBlock != nil {
		m[c.ByzantiumBlock.Uint64()] = params.EIP649FBlockReward
	}
	if c.ConstantinopleBlock != nil {
		m[c.ConstantinopleBlock.Uint64()] = params.EIP1234FBlockReward
	}
	return m
}

func (c *ChainConfig) SetEthashDifficultyBombDelaySchedule(m common2.Uint64BigMapEncodesHex) error {
	return common2.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex {
	m := common2.Uint64BigMapEncodesHex{
		0: params.FrontierBlockReward
	}
	if c.ByzantiumBlock != nil {
		m[c.ByzantiumBlock.Uint64()] = params.EIP649FBlockReward
	}
	if c.ConstantinopleBlock != nil {
		m[c.ConstantinopleBlock.Uint64()] = params.EIP1234FBlockReward
	}
	return m
}

func (c *ChainConfig) SetEthashBlockRewardSchedule(m common2.Uint64BigMapEncodesHex) error {
	return common2.ErrUnsupportedConfigNoop
}

func (c *ChainConfig) GetCliquePeriod() *uint64 {
	panic("implement me")
}

func (c *ChainConfig) SetCliquePeriod(n uint64) error {
	panic("implement me")
}

func (c *ChainConfig) GetCliqueEpoch() *uint64 {
	panic("implement me")
}

func (c *ChainConfig) SetCliqueEpoch(n uint64) error {
	panic("implement me")
}

func (c *ChainConfig) GetSealingType() paramtypes.BlockSealingT {
	panic("implement me")
}

func (c *ChainConfig) SetSealingType(paramtypes.BlockSealingT) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisSealerEthereumNonce() uint64 {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisSealerEthereumNonce(n uint64) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisSealerEthereumMixHash() common.Hash {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisSealerEthereumMixHash(h common.Hash) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisDifficulty() *big.Int {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisDifficulty(i *big.Int) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisAuthor() common.Address {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisAuthor(a common.Address) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisTimestamp() uint64 {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisTimestamp(u uint64) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisParentHash() common.Hash {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisParentHash(h common.Hash) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisExtraData() common.Hash {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisExtraData(h common.Hash) error {
	panic("implement me")
}

func (c *ChainConfig) GetGenesisGasLimit() uint64 {
	panic("implement me")
}

func (c *ChainConfig) SetGenesisGasLimit(u uint64) error {
	panic("implement me")
}

func (c *ChainConfig) ForEachAccount(fn func(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error) error {
	panic("implement me")
}

func (c *ChainConfig) UpdateAccount(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error {
	panic("implement me")
}
