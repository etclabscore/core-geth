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

package pyethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types"
)

// PyEthereumGenesisSpec represents the genesis specification format used by the
// Python Ethereum implementation.
type PyEthereumGenesisSpec struct {
	Nonce      hexutil.Bytes           `json:"nonce"`
	Timestamp  hexutil.Uint64          `json:"timestamp"`
	ExtraData  hexutil.Bytes           `json:"extraData"`
	GasLimit   hexutil.Uint64          `json:"gasLimit"`
	Difficulty *hexutil.Big            `json:"difficulty"`
	Mixhash    common.Hash             `json:"mixhash"`
	Coinbase   common.Address          `json:"coinbase"`
	Alloc      paramtypes.GenesisAlloc `json:"alloc"`
	ParentHash common.Hash             `json:"parentHash"`
}
