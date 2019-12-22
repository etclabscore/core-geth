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

package aleth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params/types/genesis"
)

// AlethGenesisSpec represents the genesis specification format used by the
// C++ Ethereum implementation.
type AlethGenesisSpec struct {
	SealEngine string `json:"sealEngine"`
	Params     struct {
		AccountStartNonce          math.HexOrDecimal64   `json:"accountStartNonce"`
		MaximumExtraDataSize       hexutil.Uint64        `json:"maximumExtraDataSize"`
		HomesteadForkBlock         *hexutil.Big          `json:"homesteadForkBlock,omitempty"`
		DaoHardforkBlock           math.HexOrDecimal64   `json:"daoHardforkBlock"`
		EIP150ForkBlock            *hexutil.Big          `json:"EIP150ForkBlock,omitempty"`
		EIP158ForkBlock            *hexutil.Big          `json:"EIP158ForkBlock,omitempty"`
		ByzantiumForkBlock         *hexutil.Big          `json:"byzantiumForkBlock,omitempty"`
		ConstantinopleForkBlock    *hexutil.Big          `json:"constantinopleForkBlock,omitempty"`
		ConstantinopleFixForkBlock *hexutil.Big          `json:"constantinopleFixForkBlock,omitempty"`
		IstanbulForkBlock          *hexutil.Big          `json:"istanbulForkBlock,omitempty"`
		MinGasLimit                hexutil.Uint64        `json:"minGasLimit"`
		MaxGasLimit                hexutil.Uint64        `json:"maxGasLimit"`
		TieBreakingGas             bool                  `json:"tieBreakingGas"`
		GasLimitBoundDivisor       math.HexOrDecimal64   `json:"gasLimitBoundDivisor"`
		MinimumDifficulty          *hexutil.Big          `json:"minimumDifficulty"`
		DifficultyBoundDivisor     *math.HexOrDecimal256 `json:"difficultyBoundDivisor"`
		DurationLimit              *math.HexOrDecimal256 `json:"durationLimit"`
		BlockReward                *hexutil.Big          `json:"blockReward"`
		NetworkID                  hexutil.Uint64        `json:"networkID"`
		ChainID                    hexutil.Uint64        `json:"chainID"`
		AllowFutureBlocks          bool                  `json:"allowFutureBlocks"`
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

	Accounts map[common.UnprefixedAddress]*AlethGenesisSpecAccount `json:"accounts"`
}

// AlethGenesisSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type AlethGenesisSpecAccount struct {
	Balance     *math.HexOrDecimal256    `json:"balance,omitempty"`
	Nonce       uint64                   `json:"nonce,omitempty"`
	Precompiled *AlethGenesisSpecBuiltin `json:"precompiled,omitempty"`
}

// AlethGenesisSpecBuiltin is the precompiled contract definition.
type AlethGenesisSpecBuiltin struct {
	Name          string                         `json:"name,omitempty"`
	StartingBlock *hexutil.Big                   `json:"startingBlock,omitempty"`
	Linear        *AlethGenesisSpecLinearPricing `json:"linear,omitempty"`
}

type AlethGenesisSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

func (spec *AlethGenesisSpec) SetPrecompile(address byte, data *AlethGenesisSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*AlethGenesisSpecAccount)
	}
	addr := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[addr]; !exist {
		spec.Accounts[addr] = &AlethGenesisSpecAccount{}
	}
	spec.Accounts[addr].Precompiled = data
}

func (spec *AlethGenesisSpec) SetAccount(address common.Address, account genesis.GenesisAccount) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*AlethGenesisSpecAccount)
	}

	a, exist := spec.Accounts[common.UnprefixedAddress(address)]
	if !exist {
		a = &AlethGenesisSpecAccount{}
		spec.Accounts[common.UnprefixedAddress(address)] = a
	}
	a.Balance = (*math.HexOrDecimal256)(account.Balance)
	a.Nonce = account.Nonce

}
