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
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/rlp"
)

var UbiqGenesisHash = common.HexToHash("0x406f1b7dd39fca54d8c702141851ed8b755463ab5b560e6f19b963b4047418af")

func DefaultUbiqGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     UbiqChainConfig,
		Nonce:      2184,
		Timestamp:  1485633600,
		ExtraData:  hexutil.MustDecode("0x4a756d6275636b734545"),
		GasLimit:   134217728,
		Difficulty: big.NewInt(80000000000),
		Coinbase:   common.BytesToAddress([]byte("3333333333333333333333333333333333333333")),
		Alloc:      decodeUbiqPrealloc(UbiqAllocData),
	}
}

func decodeUbiqPrealloc(data string) genesisT.GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(genesisT.GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = genesisT.GenesisAccount{Balance: account.Balance}
	}
	return ga
}
