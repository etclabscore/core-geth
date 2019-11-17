// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package chainspec

import (
	"encoding/binary"
	"errors"
	"math"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	math2 "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

// AlethGenesisSpec represents the genesis specification format used by the
// C++ Ethereum implementation.
type AlethGenesisSpec struct {
	SealEngine string `json:"sealEngine"`
	Params     struct {
		AccountStartNonce          math2.HexOrDecimal64   `json:"accountStartNonce"`
		MaximumExtraDataSize       hexutil.Uint64         `json:"maximumExtraDataSize"`
		HomesteadForkBlock         *hexutil.Big           `json:"homesteadForkBlock,omitempty"`
		DaoHardforkBlock           math2.HexOrDecimal64   `json:"daoHardforkBlock"`
		EIP150ForkBlock            *hexutil.Big           `json:"EIP150ForkBlock,omitempty"`
		EIP158ForkBlock            *hexutil.Big           `json:"EIP158ForkBlock,omitempty"`
		ByzantiumForkBlock         *hexutil.Big           `json:"byzantiumForkBlock,omitempty"`
		ConstantinopleForkBlock    *hexutil.Big           `json:"constantinopleForkBlock,omitempty"`
		ConstantinopleFixForkBlock *hexutil.Big           `json:"constantinopleFixForkBlock,omitempty"`
		IstanbulForkBlock          *hexutil.Big           `json:"istanbulForkBlock,omitempty"`
		MinGasLimit                hexutil.Uint64         `json:"minGasLimit"`
		MaxGasLimit                hexutil.Uint64         `json:"maxGasLimit"`
		TieBreakingGas             bool                   `json:"tieBreakingGas"`
		GasLimitBoundDivisor       math2.HexOrDecimal64   `json:"gasLimitBoundDivisor"`
		MinimumDifficulty          *hexutil.Big           `json:"minimumDifficulty"`
		DifficultyBoundDivisor     *math2.HexOrDecimal256 `json:"difficultyBoundDivisor"`
		DurationLimit              *math2.HexOrDecimal256 `json:"durationLimit"`
		BlockReward                *hexutil.Big           `json:"blockReward"`
		NetworkID                  hexutil.Uint64         `json:"networkID"`
		ChainID                    hexutil.Uint64         `json:"chainID"`
		AllowFutureBlocks          bool                   `json:"allowFutureBlocks"`
	} `json:"params"`
	Genesis struct {
		Nonce      hexutil.Bytes  `json:"nonce"`
		Difficulty *hexutil.Big   `json:"difficulty"`
		MixHash    common.Hash    `json:"mixHash"`
		Author     common.Address `json:"author"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  hexutil.Bytes  `json:"extraData"`
		GasLimit   hexutil.Uint64 `json:"gasLimit"`
	} `json:"genesis"`

	Accounts map[common.UnprefixedAddress]*alethGenesisSpecAccount `json:"accounts"`
}

// alethGenesisSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type alethGenesisSpecAccount struct {
	Balance     *math2.HexOrDecimal256   `json:"balance,omitempty"`
	Nonce       uint64                   `json:"nonce,omitempty"`
	Precompiled *alethGenesisSpecBuiltin `json:"precompiled,omitempty"`
}

// alethGenesisSpecBuiltin is the precompiled contract definition.
type alethGenesisSpecBuiltin struct {
	Name          string                         `json:"name,omitempty"`
	StartingBlock *hexutil.Big                   `json:"startingBlock,omitempty"`
	Linear        *alethGenesisSpecLinearPricing `json:"linear,omitempty"`
}

type alethGenesisSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

// NewAlethGenesisSpec converts a go-ethereum genesis block into a Aleth-specific
// chain specification format.
func NewAlethGenesisSpec(network string, genesis *core.Genesis) (*AlethGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and aleth
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Aleth format
	spec := &AlethGenesisSpec{
		SealEngine: "Ethash",
	}
	// Some defaults
	spec.Params.AccountStartNonce = 0
	spec.Params.TieBreakingGas = false
	spec.Params.AllowFutureBlocks = false

	// Dao hardfork block is a special one. The fork block is listed as 0 in the
	// config but aleth will sync with ETC clients up until the actual dao hard
	// fork block.
	spec.Params.DaoHardforkBlock = 0

	if num := genesis.Config.HomesteadBlock; num != nil {
		spec.Params.HomesteadForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.EIP150Block; num != nil {
		spec.Params.EIP150ForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.EIP158Block; num != nil {
		spec.Params.EIP158ForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.Params.ByzantiumForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.Params.ConstantinopleForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.Params.ConstantinopleFixForkBlock = (*hexutil.Big)(num)
	}
	if num := genesis.Config.IstanbulBlock; num != nil {
		spec.Params.IstanbulForkBlock = (*hexutil.Big)(num)
	}
	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(params.MinGasLimit)
	spec.Params.MaxGasLimit = (hexutil.Uint64)(math.MaxInt64)
	spec.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
	spec.Params.DifficultyBoundDivisor = (*math2.HexOrDecimal256)(params.DifficultyBoundDivisor)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.DurationLimit = (*math2.HexOrDecimal256)(params.DurationLimit)
	spec.Params.BlockReward = (*hexutil.Big)(params.FrontierBlockReward)

	spec.Genesis.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Genesis.Nonce[:], genesis.Nonce)

	spec.Genesis.MixHash = genesis.Mixhash
	spec.Genesis.Difficulty = (*hexutil.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (hexutil.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = (hexutil.Bytes)(genesis.ExtraData)
	spec.Genesis.GasLimit = (hexutil.Uint64)(genesis.GasLimit)

	for address, account := range genesis.Alloc {
		spec.setAccount(address, account)
	}

	spec.setPrecompile(1, &alethGenesisSpecBuiltin{Name: "ecrecover",
		Linear: &alethGenesisSpecLinearPricing{Base: 3000}})
	spec.setPrecompile(2, &alethGenesisSpecBuiltin{Name: "sha256",
		Linear: &alethGenesisSpecLinearPricing{Base: 60, Word: 12}})
	spec.setPrecompile(3, &alethGenesisSpecBuiltin{Name: "ripemd160",
		Linear: &alethGenesisSpecLinearPricing{Base: 600, Word: 120}})
	spec.setPrecompile(4, &alethGenesisSpecBuiltin{Name: "identity",
		Linear: &alethGenesisSpecLinearPricing{Base: 15, Word: 3}})
	if genesis.Config.ByzantiumBlock != nil {
		spec.setPrecompile(5, &alethGenesisSpecBuiltin{Name: "modexp",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock)})
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_add",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 500}})
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{Name: "alt_bn128_G1_mul",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
			Linear:        &alethGenesisSpecLinearPricing{Base: 40000}})
		spec.setPrecompile(8, &alethGenesisSpecBuiltin{Name: "alt_bn128_pairing_product",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock)})
	}
	if genesis.Config.IstanbulBlock != nil {
		if genesis.Config.ByzantiumBlock == nil {
			return nil, errors.New("invalid genesis, istanbul fork is enabled while byzantium is not")
		}
		spec.setPrecompile(6, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_add",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(7, &alethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_mul",
			StartingBlock: (*hexutil.Big)(genesis.Config.ByzantiumBlock),
		}) // Aleth hardcoded the gas policy
		spec.setPrecompile(9, &alethGenesisSpecBuiltin{
			Name:          "blake2_compression",
			StartingBlock: (*hexutil.Big)(genesis.Config.IstanbulBlock),
		})
	}
	return spec, nil
}

func (spec *AlethGenesisSpec) setPrecompile(address byte, data *alethGenesisSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alethGenesisSpecAccount)
	}
	addr := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[addr]; !exist {
		spec.Accounts[addr] = &alethGenesisSpecAccount{}
	}
	spec.Accounts[addr].Precompiled = data
}

func (spec *AlethGenesisSpec) setAccount(address common.Address, account core.GenesisAccount) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alethGenesisSpecAccount)
	}

	a, exist := spec.Accounts[common.UnprefixedAddress(address)]
	if !exist {
		a = &alethGenesisSpecAccount{}
		spec.Accounts[common.UnprefixedAddress(address)] = a
	}
	a.Balance = (*math2.HexOrDecimal256)(account.Balance)
	a.Nonce = account.Nonce

}

// ParityChainSpec is the chain specification format used by Parity.
type ParityChainSpec struct {
	Name    string `json:"name"`
	Datadir string `json:"dataDir"`
	Engine  struct {
		Ethash struct {
			Params struct {
				MinimumDifficulty      *hexutil.Big                 `json:"minimumDifficulty"`
				DifficultyBoundDivisor *hexutil.Big                 `json:"difficultyBoundDivisor"`
				DurationLimit          *hexutil.Big                 `json:"durationLimit"`
				BlockReward            hexutil.Uint64BigValOrMapHex `json:"blockReward"`
				DifficultyBombDelays   map[string]string            `json:"difficultyBombDelays"`
				HomesteadTransition    hexutil.Uint64               `json:"homesteadTransition"`
				EIP100bTransition      hexutil.Uint64               `json:"eip100bTransition"`

				// Note: DAO fields will NOT be written to Parity configs from multi-geth.
				// The chains with DAO settings are already canonical and have existing chainspecs.
				// There is no need to replicate this information.
				DaoHardforkTransition  hexutil.Uint64   `json:"daoHardforkTransition,omitempty"`
				DaoHardforkBeneficiary *common.Address  `json:"daoHardforkBeneficiary,omitempty"`
				DaoHardforkAccounts    []common.Address `json:"daoHardforkAccounts,omitempty"`

				BombDefuseTransition       hexutil.Uint64 `json:"bombDefuseTransition"`
				ECIP1010PauseTransition    hexutil.Uint64 `json:"ecip1010PauseTransition,omitempty"`
				ECIP1010ContinueTransition hexutil.Uint64 `json:"ecip1010ContinueTransition,omitempty"`
				ECIP1017EraRounds          hexutil.Uint64 `json:"ecip1017EraRounds,omitempty"`
			} `json:"params"`
		} `json:"Ethash,omitempty"`
		Clique struct {
			Params struct {
				Period hexutil.Uint64 `json:"period"`
				Epoch  hexutil.Uint64 `json:"epoch"`
			} `json:"params"`
		} `json:"Clique,omitempty"`
	} `json:"engine"`

	Params struct {
		AccountStartNonce         hexutil.Uint64       `json:"accountStartNonce"`
		MaximumExtraDataSize      hexutil.Uint64       `json:"maximumExtraDataSize"`
		MinGasLimit               hexutil.Uint64       `json:"minGasLimit"`
		GasLimitBoundDivisor      math2.HexOrDecimal64 `json:"gasLimitBoundDivisor"`
		NetworkID                 hexutil.Uint64       `json:"networkID"`
		ChainID                   hexutil.Uint64       `json:"chainID"`
		MaxCodeSize               hexutil.Uint64       `json:"maxCodeSize"`
		MaxCodeSizeTransition     hexutil.Uint64       `json:"maxCodeSizeTransition"`
		EIP98Transition           hexutil.Uint64       `json:"eip98Transition"`
		EIP150Transition          hexutil.Uint64       `json:"eip150Transition"`
		EIP160Transition          hexutil.Uint64       `json:"eip160Transition"`
		EIP161abcTransition       hexutil.Uint64       `json:"eip161abcTransition"`
		EIP161dTransition         hexutil.Uint64       `json:"eip161dTransition"`
		EIP155Transition          hexutil.Uint64       `json:"eip155Transition"`
		EIP140Transition          hexutil.Uint64       `json:"eip140Transition"`
		EIP211Transition          hexutil.Uint64       `json:"eip211Transition"`
		EIP214Transition          hexutil.Uint64       `json:"eip214Transition"`
		EIP658Transition          hexutil.Uint64       `json:"eip658Transition"`
		EIP145Transition          hexutil.Uint64       `json:"eip145Transition"`
		EIP1014Transition         hexutil.Uint64       `json:"eip1014Transition"`
		EIP1052Transition         hexutil.Uint64       `json:"eip1052Transition"`
		EIP1283Transition         hexutil.Uint64       `json:"eip1283Transition"`
		EIP1283DisableTransition  hexutil.Uint64       `json:"eip1283DisableTransition"`
		EIP1283ReenableTransition hexutil.Uint64       `json:"eip1283ReenableTransition"`
		EIP1344Transition         hexutil.Uint64       `json:"eip1344Transition"`
		EIP1884Transition         hexutil.Uint64       `json:"eip1884Transition"`
		EIP2028Transition         hexutil.Uint64       `json:"eip2028Transition"`
	} `json:"params"`

	Genesis struct {
		Seal struct {
			Ethereum struct {
				Nonce   hexutil.Bytes `json:"nonce"`
				MixHash hexutil.Bytes `json:"mixHash"`
			} `json:"ethereum"`
		} `json:"seal"`

		Difficulty *hexutil.Big   `json:"difficulty"`
		Author     common.Address `json:"author"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  hexutil.Bytes  `json:"extraData"`
		GasLimit   hexutil.Uint64 `json:"gasLimit"`
	} `json:"genesis"`

	Nodes    []string                                             `json:"nodes"`
	Accounts map[common.UnprefixedAddress]*parityChainSpecAccount `json:"accounts"`
}

// parityChainSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type parityChainSpecAccount struct {
	Balance math2.HexOrDecimal256   `json:"balance"`
	Nonce   math2.HexOrDecimal64    `json:"nonce,omitempty"`
	Builtin *parityChainSpecBuiltin `json:"builtin,omitempty"`
}

// parityChainSpecBuiltin is the precompiled contract definition.
type parityChainSpecBuiltin struct {
	Name              string                  `json:"name"`                         // Each builtin should has it own name
	Pricing           *parityChainSpecPricing `json:"pricing"`                      // Each builtin should has it own price strategy
	ActivateAt        *hexutil.Big            `json:"activate_at,omitempty"`        // ActivateAt can't be omitted if empty, default means no fork
	EIP1108Transition *hexutil.Big            `json:"eip1108_transition,omitempty"` // EIP1108Transition can't be omitted if empty, default means no fork
}

// parityChainSpecPricing represents the different pricing models that builtin
// contracts might advertise using.
type parityChainSpecPricing struct {
	Linear              *parityChainSpecLinearPricing              `json:"linear,omitempty"`
	ModExp              *parityChainSpecModExpPricing              `json:"modexp,omitempty"`
	AltBnPairing        *parityChainSpecAltBnPairingPricing        `json:"alt_bn128_pairing,omitempty"`
	AltBnConstOperation *parityChainSpecAltBnConstOperationPricing `json:"alt_bn128_const_operations,omitempty"`

	// Blake2F is the price per round of Blake2 compression
	Blake2F *parityChainSpecBlakePricing `json:"blake2_f,omitempty"`
}

type parityChainSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

type parityChainSpecModExpPricing struct {
	Divisor uint64 `json:"divisor"`
}

type parityChainSpecAltBnConstOperationPricing struct {
	Price                  uint64 `json:"price"`
	EIP1108TransitionPrice uint64 `json:"eip1108_transition_price,omitempty"` // Before Istanbul fork, this field is nil
}

type parityChainSpecAltBnPairingPricing struct {
	Base                  uint64 `json:"base"`
	Pair                  uint64 `json:"pair"`
	EIP1108TransitionBase uint64 `json:"eip1108_transition_base,omitempty"` // Before Istanbul fork, this field is nil
	EIP1108TransitionPair uint64 `json:"eip1108_transition_pair,omitempty"` // Before Istanbul fork, this field is nil
}

type parityChainSpecBlakePricing struct {
	GasPerRound uint64 `json:"gas_per_round"`
}

// NewParityChainSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewParityChainSpec(network string, genesis *core.Genesis, bootnodes []string) (*ParityChainSpec, error) {
	// Only ethash and clique are currently supported between go-ethereum and Parity
	if genesis.Config.Ethash == nil && genesis.Config.Clique == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Parity's format
	spec := &ParityChainSpec{
		Name:    network,
		Nodes:   bootnodes,
		Datadir: strings.ToLower(network),
	}
	if genesis.Config.Ethash != nil {
		spec.Engine.Ethash.Params.DifficultyBombDelays = make(map[string]string)
		spec.Engine.Ethash.Params.BlockReward = hexutil.Uint64BigValOrMapHex{}
		spec.Engine.Ethash.Params.BlockReward[0] = params.FrontierBlockReward

		spec.Engine.Ethash.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
		spec.Engine.Ethash.Params.DifficultyBoundDivisor = (*hexutil.Big)(params.DifficultyBoundDivisor)
		spec.Engine.Ethash.Params.DurationLimit = (*hexutil.Big)(params.DurationLimit)

		if b := params.FeatureOrMetaBlock(genesis.Config.EIP100FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.EIP100bTransition = hexutil.Uint64(b.Uint64())
		}

		if b := params.FeatureOrMetaBlock(genesis.Config.EIP649FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP649FBlockReward
			spec.Engine.Ethash.Params.DifficultyBombDelays[hexutil.EncodeBig(b)] = hexutil.EncodeUint64(3000000)
		}
		if b := params.FeatureOrMetaBlock(genesis.Config.EIP1234FBlock, genesis.Config.ConstantinopleBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP1234FBlockReward
			spec.Engine.Ethash.Params.DifficultyBombDelays[hexutil.EncodeBig(b)] = hexutil.EncodeUint64(2000000)
		}

		if b := genesis.Config.DisposalBlock; b != nil {
			spec.Engine.Ethash.Params.BombDefuseTransition = hexutil.Uint64(b.Uint64())
		}

		if b := genesis.Config.ECIP1010PauseBlock; b != nil {
			spec.Engine.Ethash.Params.ECIP1010PauseTransition = hexutil.Uint64(b.Uint64())
			if c := genesis.Config.ECIP1010Length; c != nil {
				spec.Engine.Ethash.Params.ECIP1010ContinueTransition = hexutil.Uint64(b.Uint64())
			}
		}
		// FIXME
		if b := params.FeatureOrMetaBlock(genesis.Config.ECIP1017EraRounds, genesis.Config.ECIP1017FBlock); b != nil {
			spec.Engine.Ethash.Params.ECIP1017EraRounds = hexutil.Uint64(genesis.Config.ECIP1017EraRounds.Uint64())
		}
	}
	if genesis.Config.Clique != nil {
		spec.Engine.Clique.Params.Period = hexutil.Uint64(genesis.Config.Clique.Period)
		spec.Engine.Clique.Params.Epoch = hexutil.Uint64(genesis.Config.Clique.Epoch)
	}

	// Homestead
	spec.Engine.Ethash.Params.HomesteadTransition = hexutil.Uint64(
		params.OneOrAllEqOfBlocks(
			genesis.Config.HomesteadBlock,
			genesis.Config.EIP2FBlock,
			genesis.Config.EIP7FBlock,
		).Uint64(),
	)

	// Tangerine Whistle : 150
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-608.md
	spec.Params.EIP150Transition = hexutil.Uint64(genesis.Config.EIP150Block.Uint64())

	// Spurious Dragon: 155, 160, 161, 170
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-607.md
	spec.Params.EIP155Transition = hexutil.Uint64(genesis.Config.EIP155Block.Uint64())
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP160FBlock, genesis.Config.EIP158Block); b != nil {
		spec.Params.EIP160Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP161FBlock, genesis.Config.EIP158Block); b != nil {
		spec.Params.EIP161abcTransition = hexutil.Uint64(b.Uint64())
		spec.Params.EIP161dTransition = hexutil.Uint64(b.Uint64())
	}

	if b := params.FeatureOrMetaBlock(genesis.Config.EIP140FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP140Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP198FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.setPrecompile(5, &parityChainSpecBuiltin{
			Name:       "modexp",
			ActivateAt: (*hexutil.Big)(b),
			Pricing: &parityChainSpecPricing{
				ModExp: &parityChainSpecModExpPricing{Divisor: 20}},
		})
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP211FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP211Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP212FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.setPrecompile(8, &parityChainSpecBuiltin{
			Name:       "alt_bn128_pairing",
			ActivateAt: (*hexutil.Big)(b),
			Pricing: &parityChainSpecPricing{
				AltBnPairing: &parityChainSpecAltBnPairingPricing{Base: 100000, Pair: 80000}},
		})
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.setPrecompile(6, &parityChainSpecBuiltin{
			Name:       "alt_bn128_add",
			ActivateAt: (*hexutil.Big)(b),
			Pricing: &parityChainSpecPricing{
				AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 500}},
		})
		spec.setPrecompile(7, &parityChainSpecBuiltin{
			Name:       "alt_bn128_mul",
			ActivateAt: (*hexutil.Big)(b),
			Pricing: &parityChainSpecPricing{
				AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 40000}},
		})
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP214FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP214Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP658FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP658Transition = hexutil.Uint64(b.Uint64())
	}

	if b := params.FeatureOrMetaBlock(genesis.Config.EIP145FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP145Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1014FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1014Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1052FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1052Transition = hexutil.Uint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1283FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1283Transition = hexutil.Uint64(b.Uint64())
	}

	// ConstantinopleFix (remove eip-1283)
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.Params.EIP1283DisableTransition = hexutil.Uint64(num.Uint64())
	}

	// EIP-152: Add Blake2 compression function F precompile
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP152FBlock, genesis.Config.IstanbulBlock); b != nil {
		//spec.Params.EIP152Transition = hexutil.Uint64(b.Uint64())
		spec.setPrecompile(9, &parityChainSpecBuiltin{
			Name:       "blake2_f",
			ActivateAt: (*hexutil.Big)(b),
			Pricing: &parityChainSpecPricing{
				Blake2F: &parityChainSpecBlakePricing{GasPerRound: 1}},
		})
	}
	// EIP-1108: Reduce alt_bn128 precompile gas costs
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1108FBlock, genesis.Config.IstanbulBlock); b != nil {
		if genesis.Config.IsEIP212F(b) && genesis.Config.IsEIP213F(b) {
			//spec.Params.EIP1108Transition = hexutil.Uint64(b.Uint64())
			spec.setPrecompile(6, &parityChainSpecBuiltin{
				Name:              "alt_bn128_add",
				ActivateAt:        (*hexutil.Big)(params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock)),
				EIP1108Transition: (*hexutil.Big)(b),
				Pricing: &parityChainSpecPricing{
					AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 500, EIP1108TransitionPrice: 150}},
			})
			spec.setPrecompile(7, &parityChainSpecBuiltin{
				Name:              "alt_bn128_mul",
				ActivateAt:        (*hexutil.Big)(params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock)),
				EIP1108Transition: (*hexutil.Big)(b),
				Pricing: &parityChainSpecPricing{
					AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 40000, EIP1108TransitionPrice: 6000}},
			})
			spec.setPrecompile(8, &parityChainSpecBuiltin{
				Name:              "alt_bn128_pairing",
				ActivateAt:        (*hexutil.Big)(params.FeatureOrMetaBlock(genesis.Config.EIP212FBlock, genesis.Config.ByzantiumBlock)),
				EIP1108Transition: (*hexutil.Big)(b),
				Pricing: &parityChainSpecPricing{
					AltBnPairing: &parityChainSpecAltBnPairingPricing{Base: 100000, Pair: 80000, EIP1108TransitionBase: 45000, EIP1108TransitionPair: 34000}},
			})
		}
	}

	// EIP-1344: Add ChainID opcode
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1344FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP1344Transition = hexutil.Uint64(b.Uint64())
	}
	// EIP-1884: Repricing for trie-size-dependent opcodes
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1884FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP1884Transition = hexutil.Uint64(b.Uint64())
	}
	// EIP-2028: Calldata gas cost reduction
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP2028FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP2028Transition = hexutil.Uint64(b.Uint64())
	}
	// EIP-2200: Rebalance net-metered SSTORE gas cost with consideration of SLOAD gas cost change
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP2200FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP1283ReenableTransition = hexutil.Uint64(b.Uint64())
	}

	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(params.MinGasLimit)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.NetworkID)
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaxCodeSize = params.MaxCodeSize
	// geth has it set from zero
	spec.Params.MaxCodeSizeTransition = 0

	// Disable this one
	spec.Params.EIP98Transition = math.MaxInt64

	spec.Genesis.Seal.Ethereum.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Genesis.Seal.Ethereum.Nonce[:], genesis.Nonce)

	spec.Genesis.Seal.Ethereum.MixHash = (hexutil.Bytes)(genesis.Mixhash[:])
	spec.Genesis.Difficulty = (*hexutil.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (hexutil.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = (hexutil.Bytes)(genesis.ExtraData)
	spec.Genesis.GasLimit = (hexutil.Uint64)(genesis.GasLimit)

	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*parityChainSpecAccount)
	}
	for address, account := range genesis.Alloc {
		bal := math2.HexOrDecimal256(*account.Balance)

		a := common.UnprefixedAddress(address)
		if _, exist := spec.Accounts[a]; !exist {
			spec.Accounts[a] = &parityChainSpecAccount{}
		}
		spec.Accounts[a].Balance = bal
		spec.Accounts[a].Nonce = math2.HexOrDecimal64(account.Nonce)
	}
	spec.setPrecompile(1, &parityChainSpecBuiltin{Name: "ecrecover",
		Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 3000}}})

	spec.setPrecompile(2, &parityChainSpecBuiltin{
		Name: "sha256", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 60, Word: 12}},
	})
	spec.setPrecompile(3, &parityChainSpecBuiltin{
		Name: "ripemd160", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 600, Word: 120}},
	})
	spec.setPrecompile(4, &parityChainSpecBuiltin{
		Name: "identity", Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 15, Word: 3}},
	})
	return spec, nil
}

func (spec *ParityChainSpec) setPrecompile(address byte, data *parityChainSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*parityChainSpecAccount)
	}
	a := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[a]; !exist {
		spec.Accounts[a] = &parityChainSpecAccount{}
	}
	spec.Accounts[a].Builtin = data
}

// PyEthereumGenesisSpec represents the genesis specification format used by the
// Python Ethereum implementation.
type PyEthereumGenesisSpec struct {
	Nonce      hexutil.Bytes     `json:"nonce"`
	Timestamp  hexutil.Uint64    `json:"timestamp"`
	ExtraData  hexutil.Bytes     `json:"extraData"`
	GasLimit   hexutil.Uint64    `json:"gasLimit"`
	Difficulty *hexutil.Big      `json:"difficulty"`
	Mixhash    common.Hash       `json:"mixhash"`
	Coinbase   common.Address    `json:"coinbase"`
	Alloc      core.GenesisAlloc `json:"alloc"`
	ParentHash common.Hash       `json:"parentHash"`
}

// NewPyEthereumGenesisSpec converts a go-ethereum genesis block into a Parity specific
// chain specification format.
func NewPyEthereumGenesisSpec(network string, genesis *core.Genesis) (*PyEthereumGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and pyethereum
	if genesis.Config.Ethash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	spec := &PyEthereumGenesisSpec{
		Timestamp:  (hexutil.Uint64)(genesis.Timestamp),
		ExtraData:  genesis.ExtraData,
		GasLimit:   (hexutil.Uint64)(genesis.GasLimit),
		Difficulty: (*hexutil.Big)(genesis.Difficulty),
		Mixhash:    genesis.Mixhash,
		Coinbase:   genesis.Coinbase,
		Alloc:      genesis.Alloc,
		ParentHash: genesis.ParentHash,
	}
	spec.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Nonce[:], genesis.Nonce)

	return spec, nil
}
