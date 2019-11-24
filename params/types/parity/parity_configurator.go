package parity

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var zero = uint64(0)

func (spec *ParityChainSpec) GetAccountStartNonce() *uint64 {
	return spec.Params.AccountStartNonce.Uint64P()
}

func (spec *ParityChainSpec) SetAccountStartNonce(i *uint64) error {
	if i == nil {
		return common2.ErrUnsupportedConfigFatal
	}
	spec.Params.AccountStartNonce = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetMaximumExtraDataSize() *uint64 {
	return spec.Params.MaximumExtraDataSize.Uint64P()
}

func (spec *ParityChainSpec) SetMaximumExtraDataSize(i *uint64) error {
	spec.Params.MaximumExtraDataSize = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetMinGasLimit() *uint64 {
	return spec.Params.MinGasLimit.Uint64P()
}

func (spec *ParityChainSpec) SetMinGasLimit(i *uint64) error {
	spec.Params.MinGasLimit = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetGasLimitBoundDivisor() *uint64 {
	return spec.Params.GasLimitBoundDivisor.Uint64P()
}

func (spec *ParityChainSpec) SetGasLimitBoundDivisor(i *uint64) error {
	spec.Params.GasLimitBoundDivisor = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetNetworkID() *uint64 {
	return spec.Params.NetworkID.Uint64P()
}

func (spec *ParityChainSpec) SetNetworkID(i *uint64) error {
	spec.Params.NetworkID = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetChainID() *uint64 {
	return spec.Params.ChainID.Uint64P()
}

func (spec *ParityChainSpec) SetChainID(i *uint64) error {
	spec.Params.ChainID = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetMaxCodeSize() *uint64 {
	return spec.Params.MaxCodeSize.Uint64P()
}

func (spec *ParityChainSpec) SetMaxCodeSize(i *uint64) error {
	spec.Params.MaxCodeSize = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetMaxCodeSizeTransition() *uint64 {
	return spec.Params.MaxCodeSizeTransition.Uint64P()
}

func (spec *ParityChainSpec) SetMaxCodeSizeTransition(i *uint64) error {
	spec.Params.MaxCodeSizeTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP7Transition() *uint64 {
	return spec.Engine.Ethash.Params.HomesteadTransition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP7Transition(i *uint64) error {
	spec.Engine.Ethash.Params.HomesteadTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP98Transition() *uint64 {
	return spec.Params.EIP98Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP98Transition(i *uint64) error {
	spec.Params.EIP98Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP150Transition() *uint64 {
	return spec.Params.EIP150Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP150Transition(i *uint64) error {
	spec.Params.EIP150Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP160Transition() *uint64 {
	return spec.Params.EIP160Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP160Transition(i *uint64) error {
	spec.Params.EIP160Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP161abcTransition() *uint64 {
	return spec.Params.EIP161abcTransition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP161abcTransition(i *uint64) error {
	spec.Params.EIP161abcTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP161dTransition() *uint64 {
	return spec.Params.EIP161dTransition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP161dTransition(i *uint64) error {
	spec.Params.EIP161dTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP155Transition() *uint64 {
	return spec.Params.EIP155Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP155Transition(i *uint64) error {
	spec.Params.EIP155Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP140Transition() *uint64 {
	return spec.Params.EIP140Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP140Transition(i *uint64) error {
	spec.Params.EIP140Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP211Transition() *uint64 {
	return spec.Params.EIP211Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP211Transition(i *uint64) error {
	spec.Params.EIP211Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP214Transition() *uint64 {
	return spec.Params.EIP214Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP214Transition(i *uint64) error {
	spec.Params.EIP214Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP658Transition() *uint64 {
	return spec.Params.EIP658Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP658Transition(i *uint64) error {
	spec.Params.EIP658Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP145Transition() *uint64 {
	return spec.Params.EIP145Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP145Transition(i *uint64) error {
	spec.Params.EIP145Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1014Transition() *uint64 {
	return spec.Params.EIP1014Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1014Transition(i *uint64) error {
	spec.Params.EIP1014Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1052Transition() *uint64 {
	return spec.Params.EIP1052Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1052Transition(i *uint64) error {
	spec.Params.EIP1052Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1283Transition() *uint64 {
	return spec.Params.EIP1283Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1283Transition(i *uint64) error {
	spec.Params.EIP1283Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1283DisableTransition() *uint64 {
	return spec.Params.EIP1283DisableTransition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1283DisableTransition(i *uint64) error {
	spec.Params.EIP1283DisableTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1283ReenableTransition() *uint64 {
	return spec.Params.EIP1283ReenableTransition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1283ReenableTransition(i *uint64) error {
	spec.Params.EIP1283ReenableTransition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1344Transition() *uint64 {
	return spec.Params.EIP1344Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1344Transition(i *uint64) error {
	spec.Params.EIP1344Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP1884Transition() *uint64 {
	return spec.Params.EIP1884Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP1884Transition(i *uint64) error {
	spec.Params.EIP1884Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) GetEIP2028Transition() *uint64 {
	return spec.Params.EIP2028Transition.Uint64P()
}

func (spec *ParityChainSpec) SetEIP2028Transition(i *uint64) error {
	spec.Params.EIP2028Transition = new(ParityU64).SetUint64(i)
	return nil
}

func (spec *ParityChainSpec) IsForked(fn func(*big.Int) bool, n *big.Int) bool {
	if n == nil || fn == nil {
		return false
	}
	return fn(n)
}

func (spec *ParityChainSpec) GetConsensusEngineType() paramtypes.ConsensusEngineT {
	if !reflect.DeepEqual(spec.Engine.Ethash, reflect.Zero(reflect.TypeOf(spec.Engine.Ethash)).Interface()) {
		return paramtypes.ConsensusEngineT_Ethash
	}
	if !reflect.DeepEqual(spec.Engine.Clique, reflect.Zero(reflect.TypeOf(spec.Engine.Clique)).Interface()) {
		return paramtypes.ConsensusEngineT_Clique
	}
	return paramtypes.ConsensusEngineT_Unknown
}

func (spec *ParityChainSpec) MustSetConsensusEngineType(engine paramtypes.ConsensusEngineT) error {
	switch engine {
	case paramtypes.ConsensusEngineT_Ethash, paramtypes.ConsensusEngineT_Clique:
		return nil
	default:
		return common2.ErrUnsupportedConfigFatal
	}
}

func (spec *ParityChainSpec) GetEthashMinimumDifficulty() *big.Int {
	return spec.Engine.Ethash.Params.MinimumDifficulty.ToInt()
}

func (spec *ParityChainSpec) SetEthashMinimumDifficulty(n *big.Int) error {
	if n == nil {
		return nil
	}
	spec.Engine.Ethash.Params.MinimumDifficulty = math.NewHexOrDecimal256(n.Int64())
	return nil
}

func (spec *ParityChainSpec) GetEthashDifficultyBoundDivisor() *big.Int {
	return spec.Engine.Ethash.Params.DifficultyBoundDivisor.ToInt()
}

func (spec *ParityChainSpec) SetEthashDifficultyBoundDivisor(n *big.Int) error {
	if n == nil {
		return nil
	}
	spec.Engine.Ethash.Params.DifficultyBoundDivisor = math.NewHexOrDecimal256(n.Int64())
	return nil
}

func (spec *ParityChainSpec) GetEthashDurationLimit() *big.Int {
	return spec.Engine.Ethash.Params.DurationLimit.ToInt()
}

func (spec *ParityChainSpec) SetEthashDurationLimit(n *big.Int) error {
	if n == nil {
		return nil
	}
	spec.Engine.Ethash.Params.DurationLimit = math.NewHexOrDecimal256(n.Int64())
	return nil
}

func (spec *ParityChainSpec) GetEthashHomesteadTransition() *big.Int {

	return spec.Engine.Ethash.Params.HomesteadTransition.Big()
}

func (spec *ParityChainSpec) SetEthashHomesteadTransition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.HomesteadTransition = new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashEIP2Transition() *big.Int {
	return spec.Engine.Ethash.Params.HomesteadTransition.Big()
}

func (spec *ParityChainSpec) SetEthashEIP2Transition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.HomesteadTransition = new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashECIP1010PauseTransition() *big.Int {
	return spec.Engine.Ethash.Params.ECIP1010PauseTransition.Big()
}

func (spec *ParityChainSpec) SetEthashECIP1010PauseTransition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.ECIP1010PauseTransition = new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashECIP1010ContinueTransition() *big.Int {
	return spec.Engine.Ethash.Params.ECIP1010ContinueTransition.Big()
}

func (spec *ParityChainSpec) SetEthashECIP1010ContinueTransition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.ECIP1010ContinueTransition= new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashECIP1017Transition() *big.Int {
	return spec.Engine.Ethash.Params.ECIP1017EraRounds.Big()
}

func (spec *ParityChainSpec) SetEthashECIP1017Transition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.ECIP1017Transition= new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashECIP1017EraRounds() *big.Int {
	return spec.Engine.Ethash.Params.ECIP1017EraRounds.Big()
}

func (spec *ParityChainSpec) SetEthashECIP1017EraRounds(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.ECIP1017EraRounds= new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashEIP100BTransition() *big.Int {
	return spec.Engine.Ethash.Params.EIP100bTransition.Big()
}

func (spec *ParityChainSpec) SetEthashEIP100BTransition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.EIP100bTransition= new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashECIP1041Transition() *big.Int {
	return spec.Engine.Ethash.Params.BombDefuseTransition.Big()
}

func (spec *ParityChainSpec) SetEthashECIP1041Transition(n *big.Int) error {
	if n == nil {
		return nil
	}
	nn := n.Uint64()
	spec.Engine.Ethash.Params.BombDefuseTransition= new(ParityU64).SetUint64(&nn)
	return nil
}

func (spec *ParityChainSpec) GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex {
	panic("implement me")
}

func (spec *ParityChainSpec) SetEthashDifficultyBombDelaySchedule(common2.Uint64BigMapEncodesHex) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex {
	panic("implement me")
}

func (spec *ParityChainSpec) SetEthashBlockRewardSchedule(common2.Uint64BigMapEncodesHex) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetCliquePeriod() *uint64 {
	panic("implement me")
}

func (spec *ParityChainSpec) SetCliquePeriod(uint64) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetCliqueEpoch() *uint64 {
	panic("implement me")
}

func (spec *ParityChainSpec) SetCliqueEpoch(uint64) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetSealingType() paramtypes.BlockSealingT {
	panic("implement me")
}

func (spec *ParityChainSpec) SetSealingType(paramtypes.BlockSealer) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisSealerNonce() uint64 {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisSealerNonce(uint64) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisSealerMixHash() common.Hash {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisSealerMixHash(common.Hash) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisDifficulty() *big.Int {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisDifficulty(*big.Int) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisAuthor() common.Address {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisAuthor(common.Address) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisTimestamp() uint64 {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisTimestamp(uint64) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisParentHash() common.Hash {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisParentHash(common.Hash) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisExtraData() common.Hash {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisExtraData(common.Hash) error {
	panic("implement me")
}

func (spec *ParityChainSpec) GetGenesisGasLimit() uint64 {
	panic("implement me")
}

func (spec *ParityChainSpec) SetGenesisGasLimit(uint64) error {
	panic("implement me")
}

func (spec *ParityChainSpec) ForEachAccount(fn paramtypes.AccountIteratorFn) error {
	panic("implement me")
}

func (spec *ParityChainSpec) SetPlainAccount(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error {
	panic("implement me")
}
