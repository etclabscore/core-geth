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

package params

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var HypraGenesisHash = common.HexToHash("0x0fb784d1481f0aa911d21d639641763ca09641413842b35f1d10eb5d208abdf8")

// ClassicGenesisBlock returns the Ethereum Classic genesis block.
func DefaultHypraGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     HypraChainConfig,
		Nonce:      4009,
		ExtraData:  hexutil.MustDecode("0x5465736c6120656e746572732074686520656c65637472696320747275636b20626174746c65206166746572206669727374204379626572747275636b20726f6c6c73206f6666207468652070726f64756374696f6e206c696e65"),
		GasLimit:   50000,
		Difficulty: big.NewInt(9_035_329),
		Alloc:      genesisT.DecodePreAlloc(hypraAllocData),
		Timestamp:  1689703500,
	}
}
