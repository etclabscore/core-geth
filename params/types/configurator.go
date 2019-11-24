package paramtypes

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/params/types/common"

)

type ChainConfigurator interface {
	CatHerder
	Forker
	ConsensusEnginator // Consensus Engine
	GenesisBlocker
	Accounter
	// CHTer
}

type CatHerder interface {
	GetAccountStartNonce() *uint64
	SetAccountStartNonce(*uint64) error
	GetMaximumExtraDataSize() *uint64
	SetMaximumExtraDataSize(*uint64) error
	GetMinGasLimit() *uint64
	SetMinGasLimit(*uint64) error
	GetGasLimitBoundDivisor() *uint64
	SetGasLimitBoundDivisor(*uint64) error
	GetNetworkID() *uint64
	SetNetworkID(*uint64) error
	GetChainID() *uint64
	SetChainID(*uint64) error
	GetMaxCodeSize() *uint64
	SetMaxCodeSize(*uint64) error
	GetMaxCodeSizeTransition() *uint64
	SetMaxCodeSizeTransition(*uint64) error
	GetEIP7Transition() *uint64
	SetEIP7Transition(*uint64) error
	GetEIP98Transition() *uint64
	SetEIP98Transition(*uint64) error
	GetEIP150Transition() *uint64
	SetEIP150Transition(*uint64) error
	GetEIP160Transition() *uint64
	SetEIP160Transition(*uint64) error
	GetEIP161abcTransition() *uint64
	SetEIP161abcTransition(*uint64) error
	GetEIP161dTransition() *uint64
	SetEIP161dTransition(*uint64) error
	GetEIP155Transition() *uint64
	SetEIP155Transition(*uint64) error
	GetEIP140Transition() *uint64
	SetEIP140Transition(*uint64) error
	GetEIP211Transition() *uint64
	SetEIP211Transition(*uint64) error
	GetEIP214Transition() *uint64
	SetEIP214Transition(*uint64) error
	GetEIP658Transition() *uint64
	SetEIP658Transition(*uint64) error
	GetEIP145Transition() *uint64
	SetEIP145Transition(*uint64) error
	GetEIP1014Transition() *uint64
	SetEIP1014Transition(*uint64) error
	GetEIP1052Transition() *uint64
	SetEIP1052Transition(*uint64) error
	GetEIP1283Transition() *uint64
	SetEIP1283Transition(*uint64) error
	GetEIP1283DisableTransition() *uint64
	SetEIP1283DisableTransition(*uint64) error
	GetEIP1283ReenableTransition() *uint64
	SetEIP1283ReenableTransition(*uint64) error
	GetEIP1344Transition() *uint64
	SetEIP1344Transition(*uint64) error
	GetEIP1884Transition() *uint64
	SetEIP1884Transition(*uint64) error
	GetEIP2028Transition() *uint64
	SetEIP2028Transition(*uint64) error
}

type Forker interface {
	IsForked(func(*big.Int) bool, *big.Int) bool
}

type ConsensusEngineT int
const (
	ConsensusEngineT_Unknown = iota
	ConsensusEngineT_Ethash
	ConsensusEngineT_Clique
)

func (c ConsensusEngineT) String() string {
	switch c {
	case ConsensusEngineT_Ethash:
		return "ethash"
	case ConsensusEngineT_Clique:
		return "clique"
	default:
		return "unknown"
	}
}

type ConsensusEnginator interface {
	GetConsensusEngineType() ConsensusEngineT
	MustSetConsensusEngineType(ConsensusEngineT)
	EthashConfigurator
	CliqueConfigurator
}

type EthashConfigurator interface {
	GetEthashMinimumDifficulty() *big.Int
	SetEthashMinimumDifficulty(*big.Int) error
	GetEthashDifficultyBoundDivisor() *big.Int
	SetEthashDifficultyBoundDivisor(*big.Int) error
	GetEthashHomesteadTransition() *big.Int
	SetEthashHomesteadTransition(*big.Int) error
	GetEthashEIP2Transition() *big.Int
	SetEthashEIP2Transition(*big.Int) error
	GetEthashECIP1010PauseTransition() *big.Int
	SetEthashECIP1010PauseTransition(*big.Int) error
	GetEthashECIP1010ContinueTransition() *big.Int
	SetEthashECIP1010ContinueTransition(*big.Int) error
	GetEthashECIP1017Transition() *big.Int
	SetEthashECIP1017Transition(*big.Int) error
	GetEthashECIP1017EraRounds() *big.Int
	SetEthashECIP1017EraRounds(*big.Int) error
	GetEthashEIP100BTransition() *big.Int
	SetEthashEIP100BTransition(*big.Int) error
	GetEthashECIP1041Transition() *big.Int
	SetEthashECIP1041Transition(*big.Int) error

	GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex
	SetEthashDifficultyBombDelaySchedule(common2.Uint64BigMapEncodesHex) error
	GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex
	SetEthashBlockRewardSchedule(common2.Uint64BigMapEncodesHex) error
}

type CliqueConfigurator interface {
	GetCliquePeriod() *uint64
	SetCliquePeriod(uint64) error
	GetCliqueEpoch() *uint64
	SetCliqueEpoch(uint64) error
}

type BlockSealingT int

const (
	BlockSealing_Unknown = iota
	BlockSealing_Ethereum
)

func (b BlockSealingT) String() string {
	switch b {
	case BlockSealing_Ethereum:
		return "ethereum"
	default:
		return "unknown"
	}
}

type BlockSealer interface {
	GetSealingType() BlockSealingT
	SetSealingType(BlockSealer) error
	BlockSealerEthereum
}

type BlockSealerEthereum interface {
	GetGenesisSealerNonce() uint64
	SetGenesisSealerNonce(uint64) error
	GetGenesisSealerMixHash() common.Hash
	SetGenesisSealerMixHash(common.Hash) error
}

type GenesisBlocker interface {
	BlockSealer
	GetGenesisDifficulty() *big.Int
	SetGenesisDifficulty(*big.Int) error
	GetGenesisAuthor() common.Address
	SetGenesisAuthor(common.Address) error
	GetGenesisTimestamp() uint64
	SetGenesisTimestamp(uint64) error
	GetGenesisParentHash() common.Hash
	SetGenesisParentHash(common.Hash) error
	GetGenesisExtraData() common.Hash
	SetGenesisExtraData(common.Hash) error
	GetGenesisGasLimit() uint64
	SetGenesisGasLimit(uint64) error
}

//type BuiltinContractT int
//
//const (
//	BuiltinContract_Unknown = iota
//	BuiltinContract_ECRecover
//	BuiltinContract_SHA256
//	BuiltinContract_RipeMD160
//	BuiltinContract_Identity
//	BuiltinContract_ModExp
//	BuiltinContract_Blake2F
//	BuiltinContract_AltBn128Add
//	BuiltinContract_AltBn128Mul
//	BuiltinContract_AltBn128Pairing
//)
//
//type BuiltinT struct {
//	Contract        BuiltinContractT
//	ActivationBlock uint64
//	SpecificationIP string
//}

type AccountIteratorFn func(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error
type Accounter interface {
	ForEachAccount(fn AccountIteratorFn) error
	SetPlainAccount(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error
	//SetBuiltin(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash, builtin *BuiltinT) error
}


