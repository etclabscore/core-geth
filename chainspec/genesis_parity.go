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

// ParityChainSpec is the chain specification format used by Parity.
type ParityChainSpec struct {
	Name    string `json:"name"`
	Datadir string `json:"dataDir"`
	Engine  struct {
		Ethash struct {
			Params struct {
				MinimumDifficulty      *hexutil.Big                   `json:"minimumDifficulty"`
				DifficultyBoundDivisor *hexutil.Big                   `json:"difficultyBoundDivisor"`
				DurationLimit          *hexutil.Big                   `json:"durationLimit"`
				BlockReward            hexutil.Uint64BigValOrMapHex   `json:"blockReward"`
				DifficultyBombDelays   hexutil.Uint64BigMapEncodesHex `json:"difficultyBombDelays"`
				HomesteadTransition    *hexutil.Uint64                `json:"homesteadTransition"`
				EIP100bTransition      *hexutil.Uint64                `json:"eip100bTransition"`

				// Note: DAO fields will NOT be written to Parity configs from multi-geth.
				// The chains with DAO settings are already canonical and have existing chainspecs.
				// There is no need to replicate this information.
				DaoHardforkTransition  *hexutil.Uint64  `json:"daoHardforkTransition,omitempty"`
				DaoHardforkBeneficiary *common.Address  `json:"daoHardforkBeneficiary,omitempty"`
				DaoHardforkAccounts    []common.Address `json:"daoHardforkAccounts,omitempty"`

				BombDefuseTransition       *hexutil.Uint64 `json:"bombDefuseTransition"`
				ECIP1010PauseTransition    *hexutil.Uint64 `json:"ecip1010PauseTransition,omitempty"`
				ECIP1010ContinueTransition *hexutil.Uint64 `json:"ecip1010ContinueTransition,omitempty"`
				ECIP1017EraRounds          *hexutil.Uint64 `json:"ecip1017EraRounds,omitempty"`
			} `json:"params"`
		} `json:"Ethash,omitempty"`
		Clique struct {
			Params struct {
				Period *hexutil.Uint64 `json:"period"`
				Epoch  *hexutil.Uint64 `json:"epoch"`
			} `json:"params"`
		} `json:"Clique,omitempty"`
	} `json:"engine"`

	Params struct {
		AccountStartNonce         *hexutil.Uint64      `json:"accountStartNonce,omitempty"`
		MaximumExtraDataSize      *hexutil.Uint64      `json:"maximumExtraDataSize,omitempty"`
		MinGasLimit               *hexutil.Uint64      `json:"minGasLimit,omitempty"`
		GasLimitBoundDivisor      math2.HexOrDecimal64 `json:"gasLimitBoundDivisor,omitempty"`
		NetworkID                 *hexutil.Uint64      `json:"networkID,omitempty"`
		ChainID                   *hexutil.Uint64      `json:"chainID,omitempty"`
		MaxCodeSize               *hexutil.Uint64      `json:"maxCodeSize,omitempty"`
		MaxCodeSizeTransition     *hexutil.Uint64      `json:"maxCodeSizeTransition,omitempty"`
		EIP98Transition           *hexutil.Uint64      `json:"eip98Transition,omitempty"`
		EIP150Transition          *hexutil.Uint64      `json:"eip150Transition,omitempty"`
		EIP160Transition          *hexutil.Uint64      `json:"eip160Transition,omitempty"`
		EIP161abcTransition       *hexutil.Uint64      `json:"eip161abcTransition,omitempty"`
		EIP161dTransition         *hexutil.Uint64      `json:"eip161dTransition,omitempty"`
		EIP155Transition          *hexutil.Uint64      `json:"eip155Transition,omitempty"`
		EIP140Transition          *hexutil.Uint64      `json:"eip140Transition,omitempty"`
		EIP211Transition          *hexutil.Uint64      `json:"eip211Transition,omitempty"`
		EIP214Transition          *hexutil.Uint64      `json:"eip214Transition,omitempty"`
		EIP658Transition          *hexutil.Uint64      `json:"eip658Transition,omitempty"`
		EIP145Transition          *hexutil.Uint64      `json:"eip145Transition,omitempty"`
		EIP1014Transition         *hexutil.Uint64      `json:"eip1014Transition,omitempty"`
		EIP1052Transition         *hexutil.Uint64      `json:"eip1052Transition,omitempty"`
		EIP1283Transition         *hexutil.Uint64      `json:"eip1283Transition,omitempty"`
		EIP1283DisableTransition  *hexutil.Uint64      `json:"eip1283DisableTransition,omitempty"`
		EIP1283ReenableTransition *hexutil.Uint64      `json:"eip1283ReenableTransition,omitempty"`
		EIP1344Transition         *hexutil.Uint64      `json:"eip1344Transition,omitempty"`
		EIP1884Transition         *hexutil.Uint64      `json:"eip1884Transition,omitempty"`
		EIP2028Transition         *hexutil.Uint64      `json:"eip2028Transition,omitempty"`

		ForkBlock     *hexutil.Uint64 `json:"forkBlock,omitempty"`
		ForkCanonHash *common.Hash    `json:"forkCanonHash,omitempty"`
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
	Balance math2.HexOrDecimal256       `json:"balance"`
	Nonce   math2.HexOrDecimal64        `json:"nonce,omitempty"`
	Code    hexutil.Bytes               `json:"code,omitempty"`
	Storage map[common.Hash]common.Hash `json:"storage,omitempty"`
	Builtin *parityChainSpecBuiltin     `json:"builtin,omitempty"`
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

func hexutilUint64(i uint64) *hexutil.Uint64 {
	p := hexutil.Uint64(i)
	return &p
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
		spec.Engine.Ethash.Params.DifficultyBombDelays = hexutil.Uint64BigMapEncodesHex{}
		spec.Engine.Ethash.Params.BlockReward = hexutil.Uint64BigValOrMapHex{}
		spec.Engine.Ethash.Params.BlockReward[0] = params.FrontierBlockReward

		spec.Engine.Ethash.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
		spec.Engine.Ethash.Params.DifficultyBoundDivisor = (*hexutil.Big)(params.DifficultyBoundDivisor)
		spec.Engine.Ethash.Params.DurationLimit = (*hexutil.Big)(params.DurationLimit)

		if b := params.FeatureOrMetaBlock(genesis.Config.EIP100FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.EIP100bTransition = hexutilUint64(b.Uint64())
		}

		if b := params.FeatureOrMetaBlock(genesis.Config.EIP649FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP649FBlockReward
			spec.Engine.Ethash.Params.DifficultyBombDelays[b.Uint64()] = big.NewInt(3000000)
		}
		if b := params.FeatureOrMetaBlock(genesis.Config.EIP1234FBlock, genesis.Config.ConstantinopleBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP1234FBlockReward
			spec.Engine.Ethash.Params.DifficultyBombDelays[b.Uint64()] = big.NewInt(2000000)
		}

		if b := genesis.Config.DisposalBlock; b != nil {
			spec.Engine.Ethash.Params.BombDefuseTransition = hexutilUint64(b.Uint64())
		}

		if b := genesis.Config.ECIP1010PauseBlock; b != nil {
			spec.Engine.Ethash.Params.ECIP1010PauseTransition = hexutilUint64(b.Uint64())
			if c := genesis.Config.ECIP1010Length; c != nil {
				spec.Engine.Ethash.Params.ECIP1010ContinueTransition = hexutilUint64(b.Uint64())
			}
		}
		// FIXME
		if b := params.FeatureOrMetaBlock(genesis.Config.ECIP1017EraRounds, genesis.Config.ECIP1017FBlock); b != nil {
			spec.Engine.Ethash.Params.ECIP1017EraRounds = hexutilUint64(genesis.Config.ECIP1017EraRounds.Uint64())
		}
	}
	if genesis.Config.Clique != nil {
		spec.Engine.Clique.Params.Period = hexutilUint64(genesis.Config.Clique.Period)
		spec.Engine.Clique.Params.Epoch = hexutilUint64(genesis.Config.Clique.Epoch)
	}

	// Homestead
	if b := params.OneOrAllEqOfBlocks(
		genesis.Config.HomesteadBlock,
		genesis.Config.EIP2FBlock,
		genesis.Config.EIP7FBlock,
	); b != nil {
		spec.Engine.Ethash.Params.HomesteadTransition = hexutilUint64(b.Uint64())
	}

	// Tangerine Whistle : 150
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-608.md
	if b := genesis.Config.EIP150Block; b != nil {
	spec.Params.EIP150Transition = hexutilUint64(b.Uint64())
	}

	// Spurious Dragon: 155, 160, 161, 170
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-607.md
	if b := genesis.Config.EIP155Block; b != nil {
		spec.Params.EIP155Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP160FBlock, genesis.Config.EIP158Block); b != nil {
		spec.Params.EIP160Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP161FBlock, genesis.Config.EIP158Block); b != nil {
		spec.Params.EIP161abcTransition = hexutilUint64(b.Uint64())
		spec.Params.EIP161dTransition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP170FBlock, genesis.Config.EIP158Block); b != nil {
		spec.Params.MaxCodeSizeTransition = hexutilUint64(b.Uint64())
		spec.Params.MaxCodeSize = hexutilUint64(params.MaxCodeSize)
	}

	if b := params.FeatureOrMetaBlock(genesis.Config.EIP140FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP140Transition = hexutilUint64(b.Uint64())
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
		spec.Params.EIP211Transition = hexutilUint64(b.Uint64())
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
		spec.Params.EIP214Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP658FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP658Transition = hexutilUint64(b.Uint64())
	}

	if b := params.FeatureOrMetaBlock(genesis.Config.EIP145FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP145Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1014FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1014Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1052FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1052Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1283FBlock, genesis.Config.ConstantinopleBlock); b != nil {
		spec.Params.EIP1283Transition = hexutilUint64(b.Uint64())
	}

	// ConstantinopleFix (remove eip-1283)
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.Params.EIP1283DisableTransition = hexutilUint64(num.Uint64())
	}

	// EIP-152: Add Blake2 compression function F precompile
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP152FBlock, genesis.Config.IstanbulBlock); b != nil {
		//spec.Params.EIP152Transition = hexutilUint64(b.Uint64())
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
			//spec.Params.EIP1108Transition = hexutilUint64(b.Uint64())
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
		spec.Params.EIP1344Transition = hexutilUint64(b.Uint64())
	}
	// EIP-1884: Repricing for trie-size-dependent opcodes
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1884FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP1884Transition = hexutilUint64(b.Uint64())
	}
	// EIP-2028: Calldata gas cost reduction
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP2028FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP2028Transition = hexutilUint64(b.Uint64())
	}
	// EIP-2200: Rebalance net-metered SSTORE gas cost with consideration of SLOAD gas cost change
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP2200FBlock, genesis.Config.IstanbulBlock); b != nil {
		spec.Params.EIP1283ReenableTransition = hexutilUint64(b.Uint64())
	}

	spec.Params.AccountStartNonce = hexutilUint64(0)
	spec.Params.MaximumExtraDataSize = hexutilUint64(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = hexutilUint64(params.MinGasLimit)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.NetworkID = hexutilUint64(genesis.Config.NetworkID)
	spec.Params.ChainID = hexutilUint64(genesis.Config.ChainID.Uint64())

	// Disable this one
	spec.Params.EIP98Transition = hexutilUint64(math.MaxInt64)

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
