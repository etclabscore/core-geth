package paramtypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
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

func setBig(i *big.Int, u *uint64) *big.Int {
	if u == nil {
		return nil
	}
	i = big.NewInt(int64(*u))
	return i
}

// upstream is used as a way to share common interface methods
// This pattern should only be used where the receiver value of the method
// is not used, ie when accessing/setting global default parameters, eg. vars/ pkg values.
var upstream = goethereum.ChainConfig{}

func (c *MultiGethChainConfig) GetAccountStartNonce() *uint64 { return upstream.GetAccountStartNonce() }
func (c *MultiGethChainConfig) SetAccountStartNonce(n *uint64) error {
	return upstream.SetAccountStartNonce(n)
}
func (c *MultiGethChainConfig) GetMaximumExtraDataSize() *uint64 {
	return upstream.GetMaximumExtraDataSize()
}
func (c *MultiGethChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	return upstream.SetMaximumExtraDataSize(n)
}
func (c *MultiGethChainConfig) GetMinGasLimit() *uint64        { return upstream.GetMinGasLimit() }
func (c *MultiGethChainConfig) SetMinGasLimit(n *uint64) error { return upstream.SetMinGasLimit(n) }
func (c *MultiGethChainConfig) GetGasLimitBoundDivisor() *uint64 {
	return upstream.GetGasLimitBoundDivisor()
}
func (c *MultiGethChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	return upstream.SetGasLimitBoundDivisor(n)
}

func (c *MultiGethChainConfig) GetNetworkID() *uint64 {
	return newU64(c.NetworkID)
}

func (c *MultiGethChainConfig) SetNetworkID(n *uint64) error {
	if n == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.NetworkID = *n
	return nil
}

func (c *MultiGethChainConfig) GetChainID() *big.Int {
	return c.ChainID
}

func (c *MultiGethChainConfig) SetChainID(n *big.Int) error {
	c.ChainID = n
	return nil
}

func (c *MultiGethChainConfig) GetMaxCodeSize() *uint64        { return upstream.GetMaxCodeSize() }
func (c *MultiGethChainConfig) SetMaxCodeSize(n *uint64) error { return upstream.SetMaxCodeSize(n) }

func (c *MultiGethChainConfig) GetEIP7Transition() *uint64 {
	return bigNewU64(c.EIP7FBlock)
}

func (c *MultiGethChainConfig) SetEIP7Transition(n *uint64) error {
	c.EIP7FBlock = setBig(c.EIP7FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP150Transition() *uint64 {
	return bigNewU64(c.EIP150Block)
}

func (c *MultiGethChainConfig) SetEIP150Transition(n *uint64) error {
	c.EIP150Block = setBig(c.EIP150Block, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP152Transition() *uint64 {
	return bigNewU64(c.EIP152FBlock)
}

func (c *MultiGethChainConfig) SetEIP152Transition(n *uint64) error {
	c.EIP152FBlock = setBig(c.EIP152FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP160Transition() *uint64 {
	return bigNewU64(c.EIP160FBlock)
}

func (c *MultiGethChainConfig) SetEIP160Transition(n *uint64) error {
	c.EIP160FBlock = setBig(c.EIP160FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP161abcTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *MultiGethChainConfig) SetEIP161abcTransition(n *uint64) error {
	c.EIP161FBlock = setBig(c.EIP161FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP161dTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *MultiGethChainConfig) SetEIP161dTransition(n *uint64) error {
	c.EIP161FBlock = setBig(c.EIP161FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP170Transition() *uint64 {
	return bigNewU64(c.EIP170FBlock)
}

func (c *MultiGethChainConfig) SetEIP170Transition(n *uint64) error {
	c.EIP170FBlock = setBig(c.EIP170FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP155Transition() *uint64 {
	return bigNewU64(c.EIP155Block)
}

func (c *MultiGethChainConfig) SetEIP155Transition(n *uint64) error {
	c.EIP155Block = setBig(c.EIP155Block, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP140Transition() *uint64 {
	return bigNewU64(c.EIP140FBlock)
}

func (c *MultiGethChainConfig) SetEIP140Transition(n *uint64) error {
	c.EIP140FBlock = setBig(c.EIP140FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP198Transition() *uint64 {
	return bigNewU64(c.EIP198FBlock)
}

func (c *MultiGethChainConfig) SetEIP198Transition(n *uint64) error {
	c.EIP198FBlock = setBig(c.EIP198FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP211Transition() *uint64 {
	return bigNewU64(c.EIP211FBlock)
}

func (c *MultiGethChainConfig) SetEIP211Transition(n *uint64) error {
	c.EIP211FBlock = setBig(c.EIP211FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP212Transition() *uint64 {
	return bigNewU64(c.EIP212FBlock)
}

func (c *MultiGethChainConfig) SetEIP212Transition(n *uint64) error {
	c.EIP212FBlock = setBig(c.EIP212FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP213Transition() *uint64 {
	return bigNewU64(c.EIP213FBlock)
}

func (c *MultiGethChainConfig) SetEIP213Transition(n *uint64) error {
	c.EIP213FBlock = setBig(c.EIP213FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP214Transition() *uint64 {
	return bigNewU64(c.EIP214FBlock)
}

func (c *MultiGethChainConfig) SetEIP214Transition(n *uint64) error {
	c.EIP214FBlock = setBig(c.EIP214FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP658Transition() *uint64 {
	return bigNewU64(c.EIP658FBlock)
}

func (c *MultiGethChainConfig) SetEIP658Transition(n *uint64) error {
	c.EIP658FBlock = setBig(c.EIP658FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP145Transition() *uint64 {
	return bigNewU64(c.EIP145FBlock)
}

func (c *MultiGethChainConfig) SetEIP145Transition(n *uint64) error {
	c.EIP145FBlock = setBig(c.EIP145FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1014Transition() *uint64 {
	return bigNewU64(c.EIP1014FBlock)
}

func (c *MultiGethChainConfig) SetEIP1014Transition(n *uint64) error {
	c.EIP1014FBlock = setBig(c.EIP1014FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1052Transition() *uint64 {
	return bigNewU64(c.EIP1052FBlock)
}

func (c *MultiGethChainConfig) SetEIP1052Transition(n *uint64) error {
	c.EIP1052FBlock = setBig(c.EIP1052FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1283Transition() *uint64 {
	return bigNewU64(c.EIP1283FBlock)
}

func (c *MultiGethChainConfig) SetEIP1283Transition(n *uint64) error {
	c.EIP1283FBlock = setBig(c.EIP1283FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1283DisableTransition() *uint64 {
	return bigNewU64(c.PetersburgBlock)
}

func (c *MultiGethChainConfig) SetEIP1283DisableTransition(n *uint64) error {
	c.PetersburgBlock = setBig(c.PetersburgBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1108Transition() *uint64 {
	return bigNewU64(c.EIP1108FBlock)
}

func (c *MultiGethChainConfig) SetEIP1108Transition(n *uint64) error {
	c.EIP1108FBlock = setBig(c.EIP1108FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1283ReenableTransition() *uint64 {
	return bigNewU64(c.EIP2200FBlock)
}

func (c *MultiGethChainConfig) SetEIP1283ReenableTransition(n *uint64) error {
	c.EIP2200FBlock = setBig(c.EIP2200FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1344Transition() *uint64 {
	return bigNewU64(c.EIP1344FBlock)
}

func (c *MultiGethChainConfig) SetEIP1344Transition(n *uint64) error {
	c.EIP1344FBlock = setBig(c.EIP1344FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP1884Transition() *uint64 {
	return bigNewU64(c.EIP1884FBlock)
}

func (c *MultiGethChainConfig) SetEIP1884Transition(n *uint64) error {
	c.EIP1884FBlock = setBig(c.EIP1884FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEIP2028Transition() *uint64 {
	return bigNewU64(c.EIP2028FBlock)
}

func (c *MultiGethChainConfig) SetEIP2028Transition(n *uint64) error {
	c.EIP2028FBlock = setBig(c.EIP2028FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) IsForked(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *MultiGethChainConfig) ForkCanonHash(n uint64) common.Hash {
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

func (c *MultiGethChainConfig) GetConsensusEngineType() common2.ConsensusEngineT {
	if c.Ethash != nil {
		return common2.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return common2.ConsensusEngineT_Clique
	}
	return common2.ConsensusEngineT_Unknown
}

func (c *MultiGethChainConfig) MustSetConsensusEngineType(t common2.ConsensusEngineT) error {
	switch t {
	case common2.ConsensusEngineT_Ethash:
		c.Ethash = new(goethereum.EthashConfig)
		return nil
	case common2.ConsensusEngineT_Clique:
		c.Clique = new(goethereum.CliqueConfig)
		return nil
	default:
		return common2.ErrUnsupportedConfigFatal
	}
}

func (c *MultiGethChainConfig) GetEthashMinimumDifficulty() *big.Int {
	return upstream.GetEthashMinimumDifficulty()
}
func (c *MultiGethChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	return upstream.SetEthashMinimumDifficulty(i)
}

func (c *MultiGethChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	return upstream.GetEthashDifficultyBoundDivisor()
}

func (c *MultiGethChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	return upstream.SetEthashDifficultyBoundDivisor(i)
}

func (c *MultiGethChainConfig) GetEthashDurationLimit() *big.Int {
	return upstream.GetEthashDurationLimit()
}

func (c *MultiGethChainConfig) SetEthashDurationLimit(i *big.Int) error {
	return upstream.SetEthashDurationLimit(i)
}

func (c *MultiGethChainConfig) GetEthashHomesteadTransition() *uint64 {
	if c.EIP2FBlock == nil || c.EIP7FBlock == nil {
		return nil
	}
	return bigNewU64(math.BigMax(c.EIP2FBlock, c.EIP7FBlock))
}

func (c *MultiGethChainConfig) SetEthashHomesteadTransition(n *uint64) error {
	c.EIP2FBlock = setBig(c.EIP2FBlock, n)
	c.EIP7FBlock = setBig(c.EIP7FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashEIP2Transition() *uint64 {
	return bigNewU64(c.EIP2FBlock)
}

func (c *MultiGethChainConfig) SetEthashEIP2Transition(n *uint64) error {
	c.EIP2FBlock = setBig(c.EIP2FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashEIP779Transition() *uint64 {
	return bigNewU64(c.DAOForkBlock)
}

func (c *MultiGethChainConfig) SetEthashEIP779Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.DAOForkBlock = setBig(c.DAOForkBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashEIP649Transition() *uint64 {
	if c.eip649FInferred {
		return bigNewU64(c.EIP649FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP649FBlock = setBig(c.EIP649FBlock, diffN)
		c.eip649FInferred = true
	}()

	diffN = common2.ExtractHostageSituationN(
		c.DifficultyBombDelaySchedule,
		common2.Uint64BigMapEncodesHex(c.BlockRewardSchedule),
		vars.EIP649DifficultyBombDelay,
		vars.EIP649FBlockReward,
	)
	return diffN
}

func (c *MultiGethChainConfig) SetEthashEIP649Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}

	c.EIP649FBlock = setBig(c.EIP649FBlock, n)
	c.eip649FInferred = true

	if n == nil {
		return nil
	}

	if c.BlockRewardSchedule == nil {
		c.BlockRewardSchedule = common2.Uint64BigMapEncodesHex{}
	}
	if c.DifficultyBombDelaySchedule == nil {
		c.DifficultyBombDelaySchedule = common2.Uint64BigMapEncodesHex{}
	}

	c.BlockRewardSchedule[*n] = vars.EIP649FBlockReward

	eip1234N := c.EIP1234FBlock
	if eip1234N == nil || eip1234N.Uint64() != *n {
		c.DifficultyBombDelaySchedule[*n] = vars.EIP649DifficultyBombDelay
	}
	// Else EIP1234 has been set to equal activation value, which means the map contains a sum value (eg 5m),
	// so the EIP649 difficulty adjustment is already accounted for.
	return nil
}

func (c *MultiGethChainConfig) GetEthashEIP1234Transition() *uint64 {
	if c.eip1234FInferred {
		return bigNewU64(c.EIP1234FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP1234FBlock = setBig(c.EIP1234FBlock, diffN)
		c.eip1234FInferred = true
	}()

	diffN = common2.ExtractHostageSituationN(
		c.DifficultyBombDelaySchedule,
		c.BlockRewardSchedule,
		vars.EIP1234DifficultyBombDelay,
		vars.EIP1234FBlockReward,
	)
	return diffN
}

func (c *MultiGethChainConfig) SetEthashEIP1234Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}

	c.EIP1234FBlock = setBig(c.EIP1234FBlock, n)
	c.eip1234FInferred = true

	if n == nil {
		return nil
	}

	if c.BlockRewardSchedule == nil {
		c.BlockRewardSchedule = common2.Uint64BigMapEncodesHex{}
	}
	if c.DifficultyBombDelaySchedule == nil {
		c.DifficultyBombDelaySchedule = common2.Uint64BigMapEncodesHex{}
	}

	// Block reward is a simple lookup; doesn't matter if overwrite or not.
	c.BlockRewardSchedule[*n] = vars.EIP1234FBlockReward

	eip649N := c.EIP649FBlock
	if eip649N == nil || eip649N.Uint64() == *n {
		// EIP649 has NOT been set, OR has been set to identical block, eg. 0 for testing
		// Overwrite key with total delay (5m)
		c.DifficultyBombDelaySchedule[*n] = vars.EIP1234DifficultyBombDelay
		return nil
	}

	c.DifficultyBombDelaySchedule[*n] = new(big.Int).Sub(vars.EIP1234DifficultyBombDelay, vars.EIP649DifficultyBombDelay)

	return nil
}

func (c *MultiGethChainConfig) GetEthashECIP1010PauseTransition() *uint64 {
	return bigNewU64(c.ECIP1010PauseBlock)
}

func (c *MultiGethChainConfig) SetEthashECIP1010PauseTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	if c.ECIP1010PauseBlock == nil && c.ECIP1010Length != nil {
		c.ECIP1010PauseBlock = setBig(c.ECIP1010PauseBlock, n)
		c.ECIP1010Length = c.ECIP1010Length.Sub(c.ECIP1010Length, c.ECIP1010PauseBlock)
		return nil
	}
	c.ECIP1010PauseBlock = setBig(c.ECIP1010PauseBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashECIP1010ContinueTransition() *uint64 {
	if c.ECIP1010PauseBlock == nil {
		return nil
	}
	if c.ECIP1010Length == nil {
		return nil
	}
	// transition = pause + length
	return bigNewU64(new(big.Int).Add(c.ECIP1010PauseBlock, c.ECIP1010Length))
}

func (c *MultiGethChainConfig) SetEthashECIP1010ContinueTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	// length = continue - pause
	if n == nil {
		return common2.ErrUnsupportedConfigNoop
	}
	if c.ECIP1010PauseBlock == nil {
		c.ECIP1010Length = new(big.Int).SetUint64(*n)
		return nil
	}
	c.ECIP1010Length = new(big.Int).Sub(big.NewInt(int64(*n)), c.ECIP1010PauseBlock)
	return nil
}

func (c *MultiGethChainConfig) GetEthashECIP1017Transition() *uint64 {
	return bigNewU64(c.ECIP1017FBlock)
}

func (c *MultiGethChainConfig) SetEthashECIP1017Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.ECIP1017FBlock = setBig(c.ECIP1017FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashECIP1017EraRounds() *uint64 {
	return bigNewU64(c.ECIP1017EraRounds)
}

func (c *MultiGethChainConfig) SetEthashECIP1017EraRounds(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.ECIP1017EraRounds = setBig(c.ECIP1017EraRounds, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashEIP100BTransition() *uint64 {
	return bigNewU64(c.EIP100FBlock)
}

func (c *MultiGethChainConfig) SetEthashEIP100BTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.EIP100FBlock = setBig(c.EIP100FBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashECIP1041Transition() *uint64 {
	return bigNewU64(c.DisposalBlock)
}

func (c *MultiGethChainConfig) SetEthashECIP1041Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.DisposalBlock = setBig(c.DisposalBlock, n)
	return nil
}

func (c *MultiGethChainConfig) GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex {
	return c.DifficultyBombDelaySchedule
}

func (c *MultiGethChainConfig) SetEthashDifficultyBombDelaySchedule(m common2.Uint64BigMapEncodesHex) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.DifficultyBombDelaySchedule = m
	return nil
}

func (c *MultiGethChainConfig) GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex {
	return c.BlockRewardSchedule
}

func (c *MultiGethChainConfig) SetEthashBlockRewardSchedule(m common2.Uint64BigMapEncodesHex) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.BlockRewardSchedule = m
	return nil
}

func (c *MultiGethChainConfig) GetCliquePeriod() *uint64 {
	if c.Clique == nil {
		return nil
	}
	return newU64(c.Clique.Period)
}

func (c *MultiGethChainConfig) SetCliquePeriod(n uint64) error {
	if c.Clique == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.Clique.Period = n
	return nil
}

func (c *MultiGethChainConfig) GetCliqueEpoch() *uint64 {
	if c.Clique == nil {
		return nil
	}
	return newU64(c.Clique.Epoch)
}

func (c *MultiGethChainConfig) SetCliqueEpoch(n uint64) error {
	if c.Clique == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.Clique.Epoch = n
	return nil
}
