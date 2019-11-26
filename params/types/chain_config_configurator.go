package paramtypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/common"
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

func setBig(i *big.Int, u *uint64) {
	if u == nil {
		return
	}
	i = big.NewInt(int64(*u))
}

// upstream is used as a way to share common interface methods
// This pattern should only be used where the receiver value of the method
// is not used, ie when accessing/setting global default parameters, eg. vars/ pkg values.
var upstream = goethereum.ChainConfig{}

func (c *ChainConfig) GetAccountStartNonce() *uint64        { return upstream.GetAccountStartNonce() }
func (c *ChainConfig) SetAccountStartNonce(n *uint64) error { return upstream.SetAccountStartNonce(n) }
func (c *ChainConfig) GetMaximumExtraDataSize() *uint64     { return upstream.GetMaximumExtraDataSize() }
func (c *ChainConfig) SetMaximumExtraDataSize(n *uint64) error {
	return upstream.SetMaximumExtraDataSize(n)
}
func (c *ChainConfig) GetMinGasLimit() *uint64          { return upstream.GetMinGasLimit() }
func (c *ChainConfig) SetMinGasLimit(n *uint64) error   { return upstream.SetMinGasLimit(n) }
func (c *ChainConfig) GetGasLimitBoundDivisor() *uint64 { return upstream.GetGasLimitBoundDivisor() }
func (c *ChainConfig) SetGasLimitBoundDivisor(n *uint64) error {
	return upstream.SetGasLimitBoundDivisor(n)
}

func (c *ChainConfig) GetNetworkID() *uint64 {
	return newU64(c.NetworkID)
}

func (c *ChainConfig) SetNetworkID(n *uint64) error {
	if n == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.NetworkID = *n
	return nil
}

func (c *ChainConfig) GetChainID() *big.Int {
	return c.ChainID
}

func (c *ChainConfig) SetChainID(n *big.Int) error {
	c.ChainID = n
	return nil
}

func (c *ChainConfig) GetMaxCodeSize() *uint64        { return upstream.GetMaxCodeSize() }
func (c *ChainConfig) SetMaxCodeSize(n *uint64) error { return upstream.SetMaxCodeSize(n) }

func (c *ChainConfig) GetEIP7Transition() *uint64 {
	return bigNewU64(c.EIP7FBlock)
}

func (c *ChainConfig) SetEIP7Transition(n *uint64) error {
	setBig(c.EIP7FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP98Transition() *uint64        { return upstream.GetEIP98Transition() }
func (c *ChainConfig) SetEIP98Transition(n *uint64) error { return upstream.SetEIP98Transition(n) }

func (c *ChainConfig) GetEIP150Transition() *uint64 {
	return bigNewU64(c.EIP150Block)
}

func (c *ChainConfig) SetEIP150Transition(n *uint64) error {
	setBig(c.EIP150Block, n)
	return nil
}

func (c *ChainConfig) GetEIP152Transition() *uint64 {
	return bigNewU64(c.EIP152FBlock)
}

func (c *ChainConfig) SetEIP152Transition(n *uint64) error {
	setBig(c.EIP152FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP160Transition() *uint64 {
	return bigNewU64(c.EIP160FBlock)
}

func (c *ChainConfig) SetEIP160Transition(n *uint64) error {
	setBig(c.EIP160FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP161abcTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *ChainConfig) SetEIP161abcTransition(n *uint64) error {
	setBig(c.EIP161FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP161dTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *ChainConfig) SetEIP161dTransition(n *uint64) error {
	setBig(c.EIP161FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP170Transition() *uint64 {
	return bigNewU64(c.EIP170FBlock)
}

func (c *ChainConfig) SetEIP170Transition(n *uint64) error {
	setBig(c.EIP170FBlock, n)
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
	return bigNewU64(c.EIP140FBlock)
}

func (c *ChainConfig) SetEIP140Transition(n *uint64) error {
	setBig(c.EIP140FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP198Transition() *uint64 {
	return bigNewU64(c.EIP198FBlock)
}

func (c *ChainConfig) SetEIP198Transition(n *uint64) error {
	setBig(c.EIP198FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP211Transition() *uint64 {
	return bigNewU64(c.EIP211FBlock)
}

func (c *ChainConfig) SetEIP211Transition(n *uint64) error {
	setBig(c.EIP211FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP212Transition() *uint64 {
	return bigNewU64(c.EIP212FBlock)
}

func (c *ChainConfig) SetEIP212Transition(n *uint64) error {
	setBig(c.EIP212FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP213Transition() *uint64 {
	return bigNewU64(c.EIP213FBlock)
}

func (c *ChainConfig) SetEIP213Transition(n *uint64) error {
	setBig(c.EIP213FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP214Transition() *uint64 {
	return bigNewU64(c.EIP214FBlock)
}

func (c *ChainConfig) SetEIP214Transition(n *uint64) error {
	setBig(c.EIP214FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP658Transition() *uint64 {
	return bigNewU64(c.EIP658FBlock)
}

func (c *ChainConfig) SetEIP658Transition(n *uint64) error {
	setBig(c.EIP658FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP145Transition() *uint64 {
	return bigNewU64(c.EIP145FBlock)
}

func (c *ChainConfig) SetEIP145Transition(n *uint64) error {
	setBig(c.EIP145FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1014Transition() *uint64 {
	return bigNewU64(c.EIP1014FBlock)
}

func (c *ChainConfig) SetEIP1014Transition(n *uint64) error {
	setBig(c.EIP1014FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1052Transition() *uint64 {
	return bigNewU64(c.EIP1052FBlock)
}

func (c *ChainConfig) SetEIP1052Transition(n *uint64) error {
	setBig(c.EIP1052FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283Transition() *uint64 {
	return bigNewU64(c.EIP1283FBlock)
}

func (c *ChainConfig) SetEIP1283Transition(n *uint64) error {
	setBig(c.EIP1283FBlock, n)
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
	return bigNewU64(c.EIP1108FBlock)
}

func (c *ChainConfig) SetEIP1108Transition(n *uint64) error {
	setBig(c.EIP1108FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1283ReenableTransition() *uint64 {
	return bigNewU64(c.EIP2200FBlock)
}

func (c *ChainConfig) SetEIP1283ReenableTransition(n *uint64) error {
	setBig(c.EIP2200FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1344Transition() *uint64 {
	return bigNewU64(c.EIP1344FBlock)
}

func (c *ChainConfig) SetEIP1344Transition(n *uint64) error {
	setBig(c.EIP1344FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP1884Transition() *uint64 {
	return bigNewU64(c.EIP1884FBlock)
}

func (c *ChainConfig) SetEIP1884Transition(n *uint64) error {
	setBig(c.EIP1884FBlock, n)
	return nil
}

func (c *ChainConfig) GetEIP2028Transition() *uint64 {
	return bigNewU64(c.EIP2028FBlock)
}

func (c *ChainConfig) SetEIP2028Transition(n *uint64) error {
	setBig(c.EIP2028FBlock, n)
	return nil
}

func (c *ChainConfig) IsForked(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *ChainConfig) GetConsensusEngineType() common.ConsensusEngineT {
	if c.Ethash != nil {
		return common2.ConsensusEngineT_Ethash
	}
	if c.Clique != nil {
		return common2.ConsensusEngineT_Clique
	}
	return common2.ConsensusEngineT_Unknown
}

func (c *ChainConfig) MustSetConsensusEngineType(t common.ConsensusEngineT) error {
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

func (c *ChainConfig) GetEthashMinimumDifficulty() *big.Int {
	return upstream.GetEthashMinimumDifficulty()
}
func (c *ChainConfig) SetEthashMinimumDifficulty(i *big.Int) error {
	return upstream.SetEthashMinimumDifficulty(i)
}

func (c *ChainConfig) GetEthashDifficultyBoundDivisor() *big.Int {
	return upstream.GetEthashDifficultyBoundDivisor()
}

func (c *ChainConfig) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	return upstream.SetEthashDifficultyBoundDivisor(i)
}

func (c *ChainConfig) GetEthashDurationLimit() *big.Int {
	return upstream.GetEthashDurationLimit()
}

func (c *ChainConfig) SetEthashDurationLimit(i *big.Int) error {
	return upstream.SetEthashDurationLimit(i)
}

func (c *ChainConfig) GetEthashHomesteadTransition() *uint64 {
	if c.EIP2FBlock == nil || c.EIP7FBlock == nil {
		return nil
	}
	return bigNewU64(math.BigMax(c.EIP2FBlock, c.EIP7FBlock))
}

func (c *ChainConfig) SetEthashHomesteadTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.EIP2FBlock, n)
	setBig(c.EIP7FBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashEIP2Transition() *uint64 {
	return bigNewU64(c.EIP2FBlock)
}

func (c *ChainConfig) SetEthashEIP2Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.EIP2FBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashEIP779Transition() *uint64 {
	return bigNewU64(c.DAOForkBlock)
}

func (c *ChainConfig) SetEthashEIP779Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.DAOForkBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashEIP649TransitionV() *uint64 {
	if c.eip649FInferred {
		return bigNewU64(c.EIP649FBlock)
	}

	var diffN *uint64
	defer func() {
		setBig(c.EIP649FBlock, diffN)
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

func (c *ChainConfig) SetEthashEIP649Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.EIP649FBlock, n)
	c.eip649FInferred = true
	if n == nil {
		return nil
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

func (c *ChainConfig) GetEthashEIP1234TransitionV() *uint64 {
	if c.eip1234FInferred {
		return bigNewU64(c.EIP1234FBlock)
	}

	var diffN *uint64
	defer func() {
		setBig(c.EIP1234FBlock, diffN)
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

func (c *ChainConfig) SetEthashEIP1234Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.EIP1234FBlock, n)
	c.eip1234FInferred = true
	if n == nil {
		return nil
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

func (c *ChainConfig) GetEthashECIP1010PauseTransition() *uint64 {
	return bigNewU64(c.ECIP1010PauseBlock)
}

func (c *ChainConfig) SetEthashECIP1010PauseTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.ECIP1010PauseBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashECIP1010ContinueTransition() *uint64 {
	if c.ECIP1010PauseBlock == nil {
		return nil
	}
	if c.ECIP1010Length == nil {
		return nil
	}
	// transition = pause + length
	return bigNewU64(new(big.Int).Add(c.ECIP1010PauseBlock, c.ECIP1010Length))
}

func (c *ChainConfig) SetEthashECIP1010ContinueTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	// length = continue - pause
	if n == nil {
		return common2.ErrUnsupportedConfigNoop
	}
	if c.ECIP1010PauseBlock == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.ECIP1010Length = new(big.Int).Sub(big.NewInt(int64(*n)), c.ECIP1010PauseBlock)
	return nil
}

func (c *ChainConfig) GetEthashECIP1017Transition() *uint64 {
	return bigNewU64(c.ECIP1017FBlock)
}

func (c *ChainConfig) SetEthashECIP1017Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.ECIP1017FBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashECIP1017EraRounds() *uint64 {
	return bigNewU64(c.ECIP1017EraRounds)
}

func (c *ChainConfig) SetEthashECIP1017EraRounds(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.ECIP1017EraRounds, n)
	return nil
}

func (c *ChainConfig) GetEthashEIP100BTransition() *uint64 {
	return bigNewU64(c.EIP100FBlock)
}

func (c *ChainConfig) SetEthashEIP100BTransition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.EIP100FBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashECIP1041Transition() *uint64 {
	return bigNewU64(c.DisposalBlock)
}

func (c *ChainConfig) SetEthashECIP1041Transition(n *uint64) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	setBig(c.DisposalBlock, n)
	return nil
}

func (c *ChainConfig) GetEthashDifficultyBombDelaySchedule() common.Uint64BigMapEncodesHex {
	return c.DifficultyBombDelaySchedule
}

func (c *ChainConfig) SetEthashDifficultyBombDelaySchedule(m common.Uint64BigMapEncodesHex) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.DifficultyBombDelaySchedule = m
	return nil
}

func (c *ChainConfig) GetEthashBlockRewardSchedule() common.Uint64BigMapEncodesHex {
	return c.BlockRewardSchedule
}

func (c *ChainConfig) SetEthashBlockRewardSchedule(m common.Uint64BigMapEncodesHex) error {
	if c.Ethash == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	c.BlockRewardSchedule = m
	return nil
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
