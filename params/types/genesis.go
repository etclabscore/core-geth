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


package paramtypes

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	common2 "github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// MultiGeth: Manually override this file generation; it requires manual editing.
//// go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var ErrGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     common2.ChainConfigurator `json:"config"`
	Nonce      uint64                    `json:"nonce"`
	Timestamp  uint64                    `json:"timestamp"`
	ExtraData  []byte                    `json:"extraData"`
	GasLimit   uint64                    `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int                  `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash               `json:"mixHash"`
	Coinbase   common.Address            `json:"coinbase"`
	Alloc      GenesisAlloc              `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

func (g *Genesis) ForEachAccount(fn func(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error) error {
	for k, v := range g.Alloc {
		if err := fn(k, v.Balance, v.Nonce, v.Code, v.Storage); err != nil {
			return err
		}
	}
	return nil
}

func (g *Genesis) UpdateAccount(address common.Address, bal *big.Int, nonce uint64, code []byte, storage map[common.Hash]common.Hash) error {
	if g.Alloc == nil {
		g.Alloc = GenesisAlloc{}
	}

	var acc GenesisAccount
	if _, ok := g.Alloc[address]; ok {
		acc = g.Alloc[address]
	}

	acc.Nonce = nonce
	acc.Balance = bal
	acc.Code = code
	acc.Storage = storage

	g.Alloc[address] = acc
	return nil
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database contains incompatible genesis (have %x, new %x)", e.Stored, e.New)
}

func DecodePreAlloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

// Following methods implement the common2.GenesisBlocker interface.

func (g *Genesis) GetSealingType() common2.BlockSealingT {
	return common2.BlockSealing_Ethereum
}

func (g *Genesis) SetSealingType(t common2.BlockSealingT) error {
	if t != common2.BlockSealing_Ethereum {
		return common2.ErrUnsupportedConfigFatal
	}
	return nil
}

func (g *Genesis) GetGenesisSealerEthereumNonce() uint64 {
	return g.Nonce
}

func (g *Genesis) SetGenesisSealerEthereumNonce(n uint64) error {
	g.Nonce = n
	return nil
}

func (g *Genesis) GetGenesisSealerEthereumMixHash() common.Hash {
	return g.Mixhash
}

func (g *Genesis) SetGenesisSealerEthereumMixHash(h common.Hash) error {
	g.Mixhash = h
	return nil
}

func (g *Genesis) GetGenesisDifficulty() *big.Int {
	return g.Difficulty
}

func (g *Genesis) SetGenesisDifficulty(i *big.Int) error {
	g.Difficulty = i
	return nil
}

func (g *Genesis) GetGenesisAuthor() common.Address {
	return g.Coinbase
}

func (g *Genesis) SetGenesisAuthor(a common.Address) error {
	g.Coinbase = a
	return nil
}

func (g *Genesis) GetGenesisTimestamp() uint64 {
	return g.Timestamp
}

func (g *Genesis) SetGenesisTimestamp(u uint64) error {
	g.Timestamp = u
	return nil
}

func (g *Genesis) GetGenesisParentHash() common.Hash {
	return g.ParentHash
}

func (g *Genesis) SetGenesisParentHash(h common.Hash) error {
	g.ParentHash = h
	return nil
}

func (g *Genesis) GetGenesisExtraData() common.Hash {
	return common.BytesToHash(g.ExtraData)
}

func (g *Genesis) SetGenesisExtraData(h common.Hash) error {
	g.ExtraData = h[:]
	return nil
}

func (g *Genesis) GetGenesisGasLimit() uint64 {
	return g.GasLimit
}

func (g *Genesis) SetGenesisGasLimit(u uint64) error {
	g.GasLimit = u
	return nil
}
