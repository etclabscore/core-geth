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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	math2 "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// ParityChainSpec is the chain specification format used by Parity.
type ParityChainSpec struct {
	Name    string `json:"name"`
	Datadir string `json:"dataDir"`
	Engine  struct {
		Ethash struct {
			Params struct {
				MinimumDifficulty      *math2.HexOrDecimal256                   `json:"minimumDifficulty"`
				DifficultyBoundDivisor *math2.HexOrDecimal256                   `json:"difficultyBoundDivisor"`
				DurationLimit          *math2.HexOrDecimal256                   `json:"durationLimit"`
				BlockReward            hexutil.Uint64BigValOrMapHex   `json:"blockReward"`
				DifficultyBombDelays   hexutil.Uint64BigMapEncodesHex `json:"difficultyBombDelays,omitempty"`
				HomesteadTransition    *math2.HexOrDecimal64                `json:"homesteadTransition"`
				EIP100bTransition      *math2.HexOrDecimal64                `json:"eip100bTransition"`

				// Note: DAO fields will NOT be written to Parity configs from multi-geth.
				// The chains with DAO settings are already canonical and have existing chainspecs.
				// There is no need to replicate this information.
				DaoHardforkTransition  *math2.HexOrDecimal64  `json:"daoHardforkTransition,omitempty"`
				DaoHardforkBeneficiary *common.Address  `json:"daoHardforkBeneficiary,omitempty"`
				DaoHardforkAccounts    []common.Address `json:"daoHardforkAccounts,omitempty"`

				BombDefuseTransition       *math2.HexOrDecimal64 `json:"bombDefuseTransition"`
				ECIP1010PauseTransition    *math2.HexOrDecimal64 `json:"ecip1010PauseTransition,omitempty"`
				ECIP1010ContinueTransition *math2.HexOrDecimal64 `json:"ecip1010ContinueTransition,omitempty"`
				ECIP1017EraRounds          *math2.HexOrDecimal64 `json:"ecip1017EraRounds,omitempty"`
			} `json:"params"`
		} `json:"Ethash,omitempty"`
		Clique struct {
			Params struct {
				Period *math2.HexOrDecimal64 `json:"period"`
				Epoch  *math2.HexOrDecimal64 `json:"epoch"`
			} `json:"params"`
		} `json:"Clique,omitempty"`
	} `json:"engine"`

	Params struct {
		AccountStartNonce         *math2.HexOrDecimal64      `json:"accountStartNonce,omitempty"`
		MaximumExtraDataSize      *math2.HexOrDecimal64      `json:"maximumExtraDataSize,omitempty"`
		MinGasLimit               *math2.HexOrDecimal64      `json:"minGasLimit,omitempty"`
		GasLimitBoundDivisor      math2.HexOrDecimal64 `json:"gasLimitBoundDivisor,omitempty"`
		NetworkID                 *math2.HexOrDecimal64      `json:"networkID,omitempty"`
		ChainID                   *math2.HexOrDecimal64      `json:"chainID,omitempty"`
		MaxCodeSize               *math2.HexOrDecimal64      `json:"maxCodeSize,omitempty"`
		MaxCodeSizeTransition     *math2.HexOrDecimal64      `json:"maxCodeSizeTransition,omitempty"`
		EIP98Transition           *math2.HexOrDecimal64      `json:"eip98Transition,omitempty"`
		EIP150Transition          *math2.HexOrDecimal64      `json:"eip150Transition,omitempty"`
		EIP160Transition          *math2.HexOrDecimal64      `json:"eip160Transition,omitempty"`
		EIP161abcTransition       *math2.HexOrDecimal64      `json:"eip161abcTransition,omitempty"`
		EIP161dTransition         *math2.HexOrDecimal64      `json:"eip161dTransition,omitempty"`
		EIP155Transition          *math2.HexOrDecimal64      `json:"eip155Transition,omitempty"`
		EIP140Transition          *math2.HexOrDecimal64      `json:"eip140Transition,omitempty"`
		EIP211Transition          *math2.HexOrDecimal64      `json:"eip211Transition,omitempty"`
		EIP214Transition          *math2.HexOrDecimal64      `json:"eip214Transition,omitempty"`
		EIP658Transition          *math2.HexOrDecimal64      `json:"eip658Transition,omitempty"`
		EIP145Transition          *math2.HexOrDecimal64      `json:"eip145Transition,omitempty"`
		EIP1014Transition         *math2.HexOrDecimal64      `json:"eip1014Transition,omitempty"`
		EIP1052Transition         *math2.HexOrDecimal64      `json:"eip1052Transition,omitempty"`
		EIP1283Transition         *math2.HexOrDecimal64      `json:"eip1283Transition,omitempty"`
		EIP1283DisableTransition  *math2.HexOrDecimal64      `json:"eip1283DisableTransition,omitempty"`
		EIP1283ReenableTransition *math2.HexOrDecimal64      `json:"eip1283ReenableTransition,omitempty"`
		EIP1344Transition         *math2.HexOrDecimal64      `json:"eip1344Transition,omitempty"`
		EIP1884Transition         *math2.HexOrDecimal64      `json:"eip1884Transition,omitempty"`
		EIP2028Transition         *math2.HexOrDecimal64      `json:"eip2028Transition,omitempty"`

		ForkBlock     *math2.HexOrDecimal64 `json:"forkBlock,omitempty"`
		ForkCanonHash *common.Hash    `json:"forkCanonHash,omitempty"`
	} `json:"params"`

	Genesis struct {
		Seal struct {
			Ethereum struct {
				Nonce   types.BlockNonce `json:"nonce"`
				MixHash hexutil.Bytes    `json:"mixHash"`
			} `json:"ethereum"`
		} `json:"seal"`

		Difficulty *math2.HexOrDecimal256   `json:"difficulty"`
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
	Name              string                       `json:"name"`                         // Each builtin should has it own name
	Pricing           *parityChainSpecPricingMaybe `json:"pricing"`                      // Each builtin should has it own price strategy
	ActivateAt        *math2.HexOrDecimal256                 `json:"activate_at,omitempty"`        // ActivateAt can't be omitted if empty, default means no fork
	EIP1108Transition *math2.HexOrDecimal256                 `json:"eip1108_transition,omitempty"` // EIP1108Transition can't be omitted if empty, default means no fork
}

type parityChainSpecPricingMaybe struct {
	Map     map[hexutil.Uint64]parityChainSpecPricingPrice
	Pricing *parityChainSpecPricing
}

type parityChainSpecPricingPrice struct {
	parityChainSpecPricing `json:"price"`
}

func (p *parityChainSpecPricingMaybe) UnmarshalJSON(input []byte) error {
	pricing := parityChainSpecPricing{}
	err := json.Unmarshal(input, &pricing)
	if err == nil && !reflect.DeepEqual(pricing, parityChainSpecPricing{}) {
		p.Pricing = &pricing
		return nil
	}
	p.Map = make(map[hexutil.Uint64]parityChainSpecPricingPrice)
	err = json.Unmarshal(input, &p.Map)
	if err != nil {
		return err
	}
	if len(p.Map) == 0 {
		panic("0map")
	}
	return nil
}
func (p parityChainSpecPricingMaybe) MarshalJSON() ([]byte, error) {
	if p.Map != nil {
		return json.Marshal(p.Map)
	}
	return json.Marshal(p.Pricing)
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

func hexutilUint64(i uint64) *math2.HexOrDecimal64 {
	p := math2.HexOrDecimal64(i)
	return &p
}

func hexOrDecimal256FromBig(i *big.Int) *math2.HexOrDecimal256 {
	if i == nil {
		return nil
	}
	return math2.NewHexOrDecimal256(i.Int64())
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

		spec.Engine.Ethash.Params.MinimumDifficulty = hexOrDecimal256FromBig(params.MinimumDifficulty)
		spec.Engine.Ethash.Params.DifficultyBoundDivisor = hexOrDecimal256FromBig(params.DifficultyBoundDivisor)
		spec.Engine.Ethash.Params.DurationLimit = hexOrDecimal256FromBig(params.DurationLimit)

		if b := params.FeatureOrMetaBlock(genesis.Config.EIP100FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.EIP100bTransition = hexutilUint64(b.Uint64())
		}

		if genesis.Config.BlockRewardSchedule != nil && len(genesis.Config.BlockRewardSchedule) > 0 {
			for k, v := range genesis.Config.BlockRewardSchedule {
				spec.Engine.Ethash.Params.BlockReward[k] = v
			}
		} else if b := params.FeatureOrMetaBlock(genesis.Config.EIP1234FBlock, genesis.Config.ConstantinopleBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP1234FBlockReward
		} else if b := params.FeatureOrMetaBlock(genesis.Config.EIP649FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.BlockReward[b.Uint64()] = params.EIP649FBlockReward
		}

		if genesis.Config.DifficultyBombDelaySchedule != nil && len(genesis.Config.DifficultyBombDelaySchedule) > 0 {
			for k, v := range genesis.Config.DifficultyBombDelaySchedule {
				spec.Engine.Ethash.Params.DifficultyBombDelays[k] = v
			}
		} else if b := params.FeatureOrMetaBlock(genesis.Config.EIP1234FBlock, genesis.Config.ConstantinopleBlock); b != nil {
			spec.Engine.Ethash.Params.DifficultyBombDelays[b.Uint64()] = big.NewInt(2000000)
		} else if b := params.FeatureOrMetaBlock(genesis.Config.EIP649FBlock, genesis.Config.ByzantiumBlock); b != nil {
			spec.Engine.Ethash.Params.DifficultyBombDelays[b.Uint64()] = big.NewInt(3000000)
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
			ActivateAt: hexOrDecimal256FromBig(b),
			Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{
				ModExp: &parityChainSpecModExpPricing{Divisor: 20}}},
		})
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP211FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.Params.EIP211Transition = hexutilUint64(b.Uint64())
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP212FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.setPrecompile(8, &parityChainSpecBuiltin{
			Name: "alt_bn128_pairing",
			//ActivateAt: hexOrDecimal256FromBig(b),
			Pricing: &parityChainSpecPricingMaybe{
				Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
					hexutil.Uint64(b.Uint64()): {
						parityChainSpecPricing{AltBnPairing: &parityChainSpecAltBnPairingPricing{Base: 100000, Pair: 80000}}}},
			}})
	}
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock); b != nil {
		spec.setPrecompile(6, &parityChainSpecBuiltin{
			Name: "alt_bn128_add",
			//ActivateAt: hexOrDecimal256FromBig(b),
			Pricing: &parityChainSpecPricingMaybe{
				Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
					hexutil.Uint64(b.Uint64()): {
						parityChainSpecPricing{AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 500}}}},
			}})
		spec.setPrecompile(7, &parityChainSpecBuiltin{
			Name: "alt_bn128_mul",
			//ActivateAt: hexOrDecimal256FromBig(b),
			Pricing: &parityChainSpecPricingMaybe{
				Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
					hexutil.Uint64(b.Uint64()): {
						parityChainSpecPricing{AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 40000}}}},
			}})
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
			ActivateAt: hexOrDecimal256FromBig(b),
			Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{
				Blake2F: &parityChainSpecBlakePricing{GasPerRound: 1}}},
		})
	}
	// EIP-1108: Reduce alt_bn128 precompile gas costs
	if b := params.FeatureOrMetaBlock(genesis.Config.EIP1108FBlock, genesis.Config.IstanbulBlock); b != nil {
		if genesis.Config.IsEIP212F(b) && genesis.Config.IsEIP213F(b) {
			spec.setPrecompile(6, &parityChainSpecBuiltin{
				Name: "alt_bn128_add",
				//ActivateAt: hexOrDecimal256FromBig(b),
				Pricing: &parityChainSpecPricingMaybe{
					Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
						hexutil.Uint64(params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock).Uint64()): parityChainSpecPricingPrice{parityChainSpecPricing{
							AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 500}},
						},
						hexutil.Uint64(b.Uint64()): parityChainSpecPricingPrice{
							parityChainSpecPricing{AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 150}}},
					},
				},
			})
			spec.setPrecompile(7, &parityChainSpecBuiltin{
				Name: "alt_bn128_mul",
				//ActivateAt: hexOrDecimal256FromBig(b),
				Pricing: &parityChainSpecPricingMaybe{
					Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
						hexutil.Uint64(params.FeatureOrMetaBlock(genesis.Config.EIP213FBlock, genesis.Config.ByzantiumBlock).Uint64()): parityChainSpecPricingPrice{
							parityChainSpecPricing{AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 40000}}},
						hexutil.Uint64(b.Uint64()): parityChainSpecPricingPrice{
							parityChainSpecPricing{AltBnConstOperation: &parityChainSpecAltBnConstOperationPricing{Price: 6000}}},
					},
				}})
			spec.setPrecompile(8, &parityChainSpecBuiltin{
				Name: "alt_bn128_pairing",
				//ActivateAt: hexOrDecimal256FromBig(b),
				Pricing: &parityChainSpecPricingMaybe{
					Map: map[hexutil.Uint64]parityChainSpecPricingPrice{
						hexutil.Uint64(params.FeatureOrMetaBlock(genesis.Config.EIP212FBlock, genesis.Config.ByzantiumBlock).Uint64()): parityChainSpecPricingPrice{
							parityChainSpecPricing{AltBnPairing: &parityChainSpecAltBnPairingPricing{Base: 100000, Pair: 80000}}},
						hexutil.Uint64(b.Uint64()): parityChainSpecPricingPrice{
							parityChainSpecPricing{AltBnPairing: &parityChainSpecAltBnPairingPricing{Base: 45000, Pair: 34000}}},
					},
				}})

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
	if id := genesis.Config.ChainID; id != nil {
		spec.Params.ChainID = hexutilUint64(id.Uint64())
	} else {
		spec.Params.ChainID = spec.Params.NetworkID
	}

	// Disable this one
	spec.Params.EIP98Transition = hexutilUint64(math.MaxInt64)
	spec.Genesis.Seal.Ethereum.Nonce = types.EncodeNonce(genesis.Nonce)

	spec.Genesis.Seal.Ethereum.MixHash = (hexutil.Bytes)(genesis.Mixhash[:])
	spec.Genesis.Difficulty = hexOrDecimal256FromBig(genesis.Difficulty)
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
	spec.setPrecompile(1, &parityChainSpecBuiltin{
		Name: "ecrecover", Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 3000}}},
	})
	spec.setPrecompile(2, &parityChainSpecBuiltin{
		Name: "sha256", Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 60, Word: 12}}},
	})
	spec.setPrecompile(3, &parityChainSpecBuiltin{
		Name: "ripemd160", Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 600, Word: 120}}},
	})
	spec.setPrecompile(4, &parityChainSpecBuiltin{
		Name: "identity", Pricing: &parityChainSpecPricingMaybe{Pricing: &parityChainSpecPricing{Linear: &parityChainSpecLinearPricing{Base: 15, Word: 3}}},
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

// ToMultiGethGenesis converts a Parity chainspec to the corresponding MultiGeth datastructure.
// Note that the return value 'core.Genesis' includes the respective 'params.ChainConfig' values.
func ParityConfigToMultiGethGenesis(c *ParityChainSpec) (*core.Genesis, error) {
	mgc := &params.ChainConfig{}
	if pars := c.Params; pars.NetworkID != nil {
		if err := checkUnsupportedValsMust(c); err != nil {
			panic(err)
		}

		if *pars.AccountStartNonce != 0 {
			return nil, errors.New("nonzero account start nonce configuration unsupported")
		}
		mgc.NetworkID = (uint64)(*pars.NetworkID)
		mgc.ChainID = pars.ChainID.Big()

		// Defaults according to Parity documentation https://wiki.parity.io/Chain-specification.html
		if mgc.ChainID == nil && pars.NetworkID != nil {
			mgc.ChainID = pars.NetworkID.Big()
		}

		// DAO
		setMultiGethDAOConfigsFromParity(mgc, c)

		// Tangerine
		mgc.EIP150Block = pars.EIP150Transition.Big()
		// mgc.EIP150Hash // optional@mg

		// Spurious
		mgc.EIP155Block = pars.EIP155Transition.Big()
		mgc.EIP160FBlock = pars.EIP160Transition.Big()
		mgc.EIP161FBlock = pars.EIP161abcTransition.Big() // and/or d
		mgc.EIP170FBlock = pars.MaxCodeSizeTransition.Big()
		if mgc.EIP170FBlock != nil && uint64(*pars.MaxCodeSize) != uint64(24576) {
			panic(fmt.Sprintf("%v != %v - unsupported configuration value", pars.MaxCodeSize, 24576))
		}

		// Byzantium
		// 100
		mgc.EIP140FBlock = pars.EIP140Transition.Big()
		// 198
		mgc.EIP211FBlock = pars.EIP211Transition.Big() // FIXME this might actually be for EIP210. :-$
		// 212
		// 213
		mgc.EIP214FBlock = pars.EIP214Transition.Big()
		// 649 - metro diff bomb, block reward
		mgc.EIP658FBlock = pars.EIP658Transition.Big()

		for _, v := range c.Accounts {
			if v.Builtin != nil {
				switch v.Builtin.Name {
				case "ripemd160", "ecrecover", "sha256", "identity":
				case "modexp":
					mgc.EIP198FBlock = new(big.Int).Set(v.Builtin.ActivateAt.ToInt())

				case "blake2_f":
					if v.Builtin.Pricing.Pricing != nil {
						mgc.EIP152FBlock = new(big.Int).Set(v.Builtin.ActivateAt.ToInt())
					}

				case "alt_bn128_pairing":
					if v.Builtin.Pricing.Pricing != nil {
						mgc.EIP212FBlock = new(big.Int).Set(v.Builtin.ActivateAt.ToInt())
						if v.Builtin.EIP1108Transition != nil {
							mgc.EIP1108FBlock = new(big.Int).Set(v.Builtin.EIP1108Transition.ToInt())
						}
					} else {
						for k, vv := range v.Builtin.Pricing.Map {
							if vv.AltBnPairing.Base == 100000 && vv.AltBnPairing.Pair == 80000 {
								mgc.EIP212FBlock = k.Big()
							} else if vv.AltBnPairing.Base == 45000 && vv.AltBnPairing.Pair == 34000 {
								if mgc.EIP212FBlock == nil {
									mgc.EIP212FBlock = k.Big()
								}
								mgc.EIP1108FBlock = k.Big()
							}
						}
					}

				case "alt_bn128_add", "alt_bn128_mul":
					if v.Builtin.Pricing.Pricing != nil {
						mgc.EIP213FBlock = new(big.Int).Set(v.Builtin.ActivateAt.ToInt())
						if v.Builtin.EIP1108Transition != nil {
							mgc.EIP1108FBlock = new(big.Int).Set(v.Builtin.EIP1108Transition.ToInt())
						}
					} else {
						for k, vv := range v.Builtin.Pricing.Map {
							if v.Builtin.Name == "alt_bn128_add" {
								if vv.AltBnConstOperation.Price == 500 {
									mgc.EIP213FBlock = k.Big()
								}
								if vv.AltBnConstOperation.Price == 150 {
									if mgc.EIP213FBlock == nil {
										mgc.EIP213FBlock = k.Big()
									}
									mgc.EIP1108FBlock = k.Big()
								}
							}
							if v.Builtin.Name == "alt_bn128_mul" {
								if vv.AltBnConstOperation.Price == 40000 {
									mgc.EIP213FBlock = k.Big()
								}
								if vv.AltBnConstOperation.Price == 6000 {
									if mgc.EIP213FBlock == nil {
										mgc.EIP213FBlock = k.Big()
									}
									mgc.EIP1108FBlock = k.Big()
								}
							}
						}
					}
				default:
					panic("unsupported builtin contract: " + v.Builtin.Name)
				}
			}
		}

		// Constantinople
		mgc.EIP145FBlock = pars.EIP145Transition.Big()
		mgc.EIP1014FBlock = pars.EIP1014Transition.Big()
		mgc.EIP1052FBlock = pars.EIP1052Transition.Big()
		mgc.EIP1283FBlock = pars.EIP1283Transition.Big()

		// St. Peters aka ConstantinopleFix
		mgc.PetersburgBlock = pars.EIP1283DisableTransition.Big()

		// Istanbul
		mgc.EIP1344FBlock = pars.EIP1344Transition.Big()
		mgc.EIP1884FBlock = pars.EIP1884Transition.Big()
		mgc.EIP2028FBlock = pars.EIP2028Transition.Big()
		mgc.EIP2200FBlock = pars.EIP1283ReenableTransition.Big()
	}

	if ethash := c.Engine.Ethash; ethash.Params.MinimumDifficulty != nil {

		pars := ethash.Params

		mgc.Ethash = &params.EthashConfig{}

		mgc.HomesteadBlock = pars.HomesteadTransition.Big()
		mgc.EIP100FBlock = pars.EIP100bTransition.Big()
		mgc.DisposalBlock = pars.BombDefuseTransition.Big()
		mgc.ECIP1010PauseBlock = pars.ECIP1010PauseTransition.Big()
		mgc.ECIP1010Length = func() *big.Int {
			if pars.ECIP1010ContinueTransition != nil {
				return new(big.Int).Sub(pars.ECIP1010ContinueTransition.Big(), pars.ECIP1010PauseTransition.Big())
			} else if pars.ECIP1010PauseTransition == nil && pars.ECIP1010ContinueTransition == nil {
				return nil
			}
			return big.NewInt(0)
		}()
		mgc.ECIP1017EraRounds = pars.ECIP1017EraRounds.Big()

		mgc.DifficultyBombDelaySchedule = hexutil.Uint64BigMapEncodesHex{}
		for k, v := range pars.DifficultyBombDelays {
			mgc.DifficultyBombDelaySchedule[k] = v
		}
		mgc.BlockRewardSchedule = hexutil.Uint64BigMapEncodesHex{}
		for k, v := range pars.BlockReward {
			mgc.BlockRewardSchedule[k] = v
		}

	} else if clique := c.Engine.Clique; clique.Params.Period != nil {
		mgc.Clique = &params.CliqueConfig{
			Period: (uint64)(*clique.Params.Period),
			Epoch:  (uint64)(*clique.Params.Epoch),
		}

	} else {
		return nil, errors.New("unsupported engine")
	}
	mgg := &core.Genesis{
		Config: mgc,
	}
	if c.Genesis.Difficulty != nil {
		seal := c.Genesis.Seal.Ethereum

		mgg.Nonce = seal.Nonce.Uint64()
		mgg.Mixhash = common.BytesToHash(seal.MixHash)
		mgg.Timestamp = (uint64)(c.Genesis.Timestamp)
		mgg.GasLimit = (uint64)(c.Genesis.GasLimit)
		mgg.Difficulty = c.Genesis.Difficulty.ToInt()
		mgg.Coinbase = c.Genesis.Author
		mgg.ParentHash = c.Genesis.ParentHash
		mgg.ExtraData = c.Genesis.ExtraData
	}
	if c.Accounts != nil {
		mgg.Alloc = core.GenesisAlloc{}
		for k, v := range c.Accounts {
			addr := common.Address(k)

			bal := (big.Int)(v.Balance)
			mgg.Alloc[addr] = core.GenesisAccount{
				Nonce:   (uint64)(v.Nonce),
				Balance: &bal,
				Code:    v.Code,
				Storage: v.Storage,
			}
		}
	}
	return mgg, nil
}

func checkUnsupportedValsMust(spec *ParityChainSpec) error {
	// FIXME

	if spec.Params.EIP161abcTransition != nil && spec.Params.EIP161dTransition != nil &&
		*spec.Params.EIP161abcTransition != *spec.Params.EIP161dTransition {
		panic(spec.Name + ": eip161abc vs. eip161d transition not supported")
	}
	// TODO...
	// unsupportedValuesMust := map[interface{}]interface{}{
	// 	pars.AccountStartNonce:                       uint64(0),
	// 	pars.MaximumExtraDataSize:                    uint64(32),
	// 	pars.MinGasLimit:                             uint64(5000),
	// 	pars.SubProtocolName:                         "",
	// 	pars.ValidateChainIDTransition:               nil,
	// 	pars.ValidateChainReceiptsTransition:         nil,
	// 	pars.DustProtectionTransition:                nil,
	// 	pars.NonceCapIncrement:                       nil,
	// 	pars.RemoveDustContracts:                     false,
	// 	pars.EIP210Transition:                        nil,
	// 	pars.EIP210ContractAddress:                   nil,
	// 	pars.EIP210ContractCode:                      nil,
	// 	pars.ApplyReward:                             false,
	// 	pars.TransactionPermissionContract:           nil,
	// 	pars.TransactionPermissionContractTransition: nil,
	// 	pars.KIP4Transition:                          nil,
	// 	pars.KIP6Transition:                          nil,
	// }
	// i := -1
	// for k, v := range unsupportedValuesMust {
	// 	i++
	// 	if v == nil && k == nil {
	// 		continue
	// 	}
	// 	if v != nil && !reflect.DeepEqual(k, v) {
	// 		panic(fmt.Sprintf("%d: %v != %v - unsupported configuration value", i, k, v))
	// 	}
	// }
	return nil
}

func setMultiGethDAOConfigsFromParity(mgc *params.ChainConfig, spec *ParityChainSpec) {
	if spec.Params.ForkCanonHash != nil {
		if (*spec.Params.ForkCanonHash == common.HexToHash("0x4985f5ca3d2afbec36529aa96f74de3cc10a2a4a6c44f2157a57d2c6059a11bb")) ||
			(*spec.Params.ForkCanonHash == common.HexToHash("0x3e12d5c0f8d63fbc5831cc7f7273bd824fa4d0a9a4102d65d99a7ea5604abc00")) {

			mgc.DAOForkBlock = new(big.Int).SetUint64(uint64(*spec.Params.ForkBlock))
			mgc.DAOForkSupport = true
		}
		if *spec.Params.ForkCanonHash == common.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f") {
			mgc.DAOForkBlock = new(big.Int).SetUint64(uint64(*spec.Params.ForkBlock))
		}
	}
}
