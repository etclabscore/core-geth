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
package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

// MixGenesisBlock returns the Mix genesis block.
func DefaultMixGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.MixChainConfig,
		Nonce:      0x1391eaa92b871f91,
		ExtraData:  hexutil.MustDecode("0x77656c636f6d65746f7468656c696e6b6564776f726c64000000000000000000"),
		GasLimit:   3000000,
		Difficulty: big.NewInt(1048576),
		Alloc:      decodePrealloc(mixAllocData),
	}
}
