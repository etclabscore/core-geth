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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

func DefaultMordorGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.MordorChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  hexutil.MustDecode("0x70686f656e697820636869636b656e206162737572642062616e616e61"),
		GasLimit:   hexutil.MustDecodeUint64("0x2fefd8"),
		Difficulty: hexutil.MustDecodeBig("0x20000"),
		Timestamp: hexutil.MustDecodeUint64("0x5d9676db"),
		Alloc:      GenesisAlloc{},
	}
}
