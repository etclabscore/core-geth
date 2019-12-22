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

package tconvert

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/aleth"
	"github.com/ethereum/go-ethereum/params/types/genesis"
	"github.com/ethereum/go-ethereum/params/vars"
)

// NewAlethGenesisSpec converts a go-ethereum genesis block into a Aleth-specific
// chain specification format.
func NewAlethGenesisSpec(network string, genesis *genesis.Genesis) (*aleth.AlethGenesisSpec, error) {
	// Only ethash is currently supported between go-ethereum and aleth
	if !genesis.Config.GetConsensusEngineType().IsEthash() {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Aleth format
	spec := &aleth.AlethGenesisSpec{
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

	if num := genesis.Config.GetEthashHomesteadTransition(); num != nil {
		spec.Params.HomesteadForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big()) // Fuckin beautiful
	}
	if num := genesis.Config.GetEIP150Transition(); num != nil {
		spec.Params.EIP150ForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	// Just choosing random signal blocks for any given fork bundle.
	if num := genesis.Config.GetEIP161dTransition(); num != nil {
		spec.Params.EIP158ForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	if num := genesis.Config.GetEthashEIP649Transition(); num != nil {
		spec.Params.ByzantiumForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	if num := genesis.Config.GetEthashEIP1234Transition(); num != nil {
		spec.Params.ConstantinopleForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	if num := genesis.Config.GetEIP1283DisableTransition(); num != nil {
		spec.Params.ConstantinopleFixForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	if num := genesis.Config.GetEIP145Transition(); num != nil {
		spec.Params.IstanbulForkBlock = (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())
	}
	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.GetChainID().Uint64())
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.GetChainID().Uint64())
	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(vars.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(vars.MinGasLimit)
	spec.Params.MaxGasLimit = (hexutil.Uint64)(math.MaxInt64)
	spec.Params.MinimumDifficulty = (*hexutil.Big)(vars.MinimumDifficulty)
	spec.Params.DifficultyBoundDivisor = (*math.HexOrDecimal256)(vars.DifficultyBoundDivisor)
	spec.Params.GasLimitBoundDivisor = (math.HexOrDecimal64)(vars.GasLimitBoundDivisor)
	spec.Params.DurationLimit = (*math.HexOrDecimal256)(vars.DurationLimit)
	spec.Params.BlockReward = (*hexutil.Big)(vars.FrontierBlockReward)

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
		spec.SetAccount(address, account)
	}

	spec.SetPrecompile(1, &aleth.AlethGenesisSpecBuiltin{Name: "ecrecover",
		Linear: &aleth.AlethGenesisSpecLinearPricing{Base: 3000}})
	spec.SetPrecompile(2, &aleth.AlethGenesisSpecBuiltin{Name: "sha256",
		Linear: &aleth.AlethGenesisSpecLinearPricing{Base: 60, Word: 12}})
	spec.SetPrecompile(3, &aleth.AlethGenesisSpecBuiltin{Name: "ripemd160",
		Linear: &aleth.AlethGenesisSpecLinearPricing{Base: 600, Word: 120}})
	spec.SetPrecompile(4, &aleth.AlethGenesisSpecBuiltin{Name: "identity",
		Linear: &aleth.AlethGenesisSpecLinearPricing{Base: 15, Word: 3}})
	if num := genesis.Config.GetEIP212Transition(); num != nil {
		spec.SetPrecompile(5, &aleth.AlethGenesisSpecBuiltin{Name: "modexp",
			StartingBlock: (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())})
		spec.SetPrecompile(6, &aleth.AlethGenesisSpecBuiltin{Name: "alt_bn128_G1_add",
			StartingBlock: (*hexutil.Big)(((*hexutil.Uint64)(num)).Big()),
			Linear:        &aleth.AlethGenesisSpecLinearPricing{Base: 500}})
		spec.SetPrecompile(7, &aleth.AlethGenesisSpecBuiltin{Name: "alt_bn128_G1_mul",
			StartingBlock: (*hexutil.Big)(((*hexutil.Uint64)(num)).Big()),
			Linear:        &aleth.AlethGenesisSpecLinearPricing{Base: 40000}})
		spec.SetPrecompile(8, &aleth.AlethGenesisSpecBuiltin{Name: "alt_bn128_pairing_product",
			StartingBlock: (*hexutil.Big)(((*hexutil.Uint64)(num)).Big())})
	}
	if num := genesis.Config.GetEIP1108Transition(); num != nil {
		byz := genesis.Config.GetEIP212Transition()
		if byz == nil {
			return nil, errors.New("invalid genesis, istanbul fork is enabled while byzantium is not")
		}
		spec.SetPrecompile(6, &aleth.AlethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_add",
			StartingBlock: (spec.Params.ByzantiumForkBlock),
		}) // Aleth hardcoded the gas policy
		spec.SetPrecompile(7, &aleth.AlethGenesisSpecBuiltin{
			Name:          "alt_bn128_G1_mul",
			StartingBlock: (spec.Params.ByzantiumForkBlock),
		}) // Aleth hardcoded the gas policy
		spec.SetPrecompile(9, &aleth.AlethGenesisSpecBuiltin{
			Name:          "blake2_compression",
			StartingBlock: (*hexutil.Big)(((*hexutil.Uint64)(num)).Big()),
		})
	}
	return spec, nil
}
