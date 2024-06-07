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

package genesisT

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
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:generate go run github.com/fjl/gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate go run github.com/fjl/gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var ErrGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     ctypes.ChainConfigurator `json:"config"`
	Nonce      uint64                   `json:"nonce"`
	Timestamp  uint64                   `json:"timestamp"`
	ExtraData  []byte                   `json:"extraData"`
	GasLimit   uint64                   `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int                 `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash              `json:"mixHash"`
	Coinbase   common.Address           `json:"coinbase"`
	Alloc      GenesisAlloc             `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number        uint64      `json:"number"`
	GasUsed       uint64      `json:"gasUsed"`
	ParentHash    common.Hash `json:"parentHash"`
	BaseFee       *big.Int    `json:"baseFeePerGas,omitempty"` // EIP-1559
	ExcessBlobGas *uint64     `json:"excessBlobGas,omitempty"` // EIP-4844
	BlobGasUsed   *uint64     `json:"blobGasUsed,omitempty"`   // EIP-4844
}

func (g *Genesis) GetElasticityMultiplier() uint64 {
	return g.Config.GetElasticityMultiplier()
}

func (g *Genesis) SetElasticityMultiplier(n uint64) error {
	return g.Config.SetElasticityMultiplier(n)
}

func (g *Genesis) GetBaseFeeChangeDenominator() uint64 {
	return g.Config.GetBaseFeeChangeDenominator()
}

func (g *Genesis) SetBaseFeeChangeDenominator(n uint64) error {
	return g.Config.SetBaseFeeChangeDenominator(n)
}

func (g *Genesis) GetEIP3651TransitionTime() *uint64 {
	return g.Config.GetEIP3651TransitionTime()
}

func (g *Genesis) SetEIP3651TransitionTime(n *uint64) error {
	return g.Config.SetEIP3651TransitionTime(n)
}

func (g *Genesis) GetEIP3855TransitionTime() *uint64 {
	return g.Config.GetEIP3855TransitionTime()
}

func (g *Genesis) SetEIP3855TransitionTime(n *uint64) error {
	return g.Config.SetEIP3855TransitionTime(n)
}

func (g *Genesis) GetEIP3860TransitionTime() *uint64 {
	return g.Config.GetEIP3860TransitionTime()
}

func (g *Genesis) SetEIP3860TransitionTime(n *uint64) error {
	return g.Config.SetEIP3860TransitionTime(n)
}

func (g *Genesis) GetEIP4895TransitionTime() *uint64 {
	return g.Config.GetEIP4895TransitionTime()
}

func (g *Genesis) SetEIP4895TransitionTime(n *uint64) error {
	return g.Config.SetEIP4895TransitionTime(n)
}

func (g *Genesis) GetEIP6049TransitionTime() *uint64 {
	return g.Config.GetEIP6049TransitionTime()
}

func (g *Genesis) SetEIP6049TransitionTime(n *uint64) error {
	return g.Config.SetEIP6049TransitionTime(n)
}

func (g *Genesis) GetEIP3651Transition() *uint64 {
	return g.Config.GetEIP3651Transition()
}

func (g *Genesis) SetEIP3651Transition(n *uint64) error {
	return g.Config.SetEIP3651Transition(n)
}

func (g *Genesis) GetEIP3855Transition() *uint64 {
	return g.Config.GetEIP3855Transition()
}

func (g *Genesis) SetEIP3855Transition(n *uint64) error {
	return g.Config.SetEIP3855Transition(n)
}

func (g *Genesis) GetEIP3860Transition() *uint64 {
	return g.Config.GetEIP3860Transition()
}

func (g *Genesis) SetEIP3860Transition(n *uint64) error {
	return g.Config.SetEIP3860Transition(n)
}

func (g *Genesis) GetEIP4895Transition() *uint64 {
	return g.Config.GetEIP4895Transition()
}

func (g *Genesis) SetEIP4895Transition(n *uint64) error {
	return g.Config.SetEIP4895Transition(n)
}

func (g *Genesis) GetEIP6049Transition() *uint64 {
	return g.Config.GetEIP6049Transition()
}

func (g *Genesis) SetEIP6049Transition(n *uint64) error {
	return g.Config.SetEIP6049Transition(n)
}

func (g *Genesis) GetEIP4844TransitionTime() *uint64 {
	return g.Config.GetEIP4844TransitionTime()
}

func (g *Genesis) SetEIP4844TransitionTime(n *uint64) error {
	return g.Config.SetEIP4844TransitionTime(n)
}

func (g *Genesis) GetEIP7516TransitionTime() *uint64 {
	return g.Config.GetEIP7516TransitionTime()
}

func (g *Genesis) SetEIP7516TransitionTime(n *uint64) error {
	return g.Config.SetEIP7516TransitionTime(n)
}

func (g *Genesis) GetEIP1153TransitionTime() *uint64 {
	return g.Config.GetEIP1153TransitionTime()
}

func (g *Genesis) SetEIP1153TransitionTime(n *uint64) error {
	return g.Config.SetEIP1153TransitionTime(n)
}

func (g *Genesis) GetEIP5656TransitionTime() *uint64 {
	return g.Config.GetEIP5656TransitionTime()
}

func (g *Genesis) SetEIP5656TransitionTime(n *uint64) error {
	return g.Config.SetEIP5656TransitionTime(n)
}

func (g *Genesis) GetEIP6780TransitionTime() *uint64 {
	return g.Config.GetEIP6780TransitionTime()
}

func (g *Genesis) SetEIP6780TransitionTime(n *uint64) error {
	return g.Config.SetEIP6780TransitionTime(n)
}

func (g *Genesis) GetEIP4788TransitionTime() *uint64 {
	return g.Config.GetEIP4788TransitionTime()
}

func (g *Genesis) SetEIP4788TransitionTime(n *uint64) error {
	return g.Config.SetEIP4788TransitionTime(n)
}

// Cancun by block number
func (g *Genesis) GetEIP4844Transition() *uint64 {
	return g.Config.GetEIP4844Transition()
}

func (g *Genesis) SetEIP4844Transition(n *uint64) error {
	return g.Config.SetEIP4844Transition(n)
}

func (g *Genesis) GetEIP7516Transition() *uint64 {
	return g.Config.GetEIP7516Transition()
}

func (g *Genesis) SetEIP7516Transition(n *uint64) error {
	return g.Config.SetEIP7516Transition(n)
}

func (g *Genesis) GetEIP1153Transition() *uint64 {
	return g.Config.GetEIP1153Transition()
}

func (g *Genesis) SetEIP1153Transition(n *uint64) error {
	return g.Config.SetEIP1153Transition(n)
}

func (g *Genesis) GetEIP5656Transition() *uint64 {
	return g.Config.GetEIP5656Transition()
}

func (g *Genesis) SetEIP5656Transition(n *uint64) error {
	return g.Config.SetEIP5656Transition(n)
}

func (g *Genesis) GetEIP6780Transition() *uint64 {
	return g.Config.GetEIP6780Transition()
}

func (g *Genesis) SetEIP6780Transition(n *uint64) error {
	return g.Config.SetEIP6780Transition(n)
}

func (g *Genesis) GetEIP4788Transition() *uint64 {
	return g.Config.GetEIP4788Transition()
}

func (g *Genesis) SetEIP4788Transition(n *uint64) error {
	return g.Config.SetEIP4788Transition(n)
}

// Verkle Trie
func (g *Genesis) GetVerkleTransitionTime() *uint64 {
	return g.Config.GetVerkleTransitionTime()
}

func (g *Genesis) SetVerkleTransitionTime(n *uint64) error {
	return g.Config.SetVerkleTransitionTime(n)
}

func (g *Genesis) GetVerkleTransition() *uint64 {
	return g.Config.GetVerkleTransition()
}

func (g *Genesis) SetVerkleTransition(n *uint64) error {
	return g.Config.SetVerkleTransition(n)
}

func (g *Genesis) IsVerkle() bool {
	return g.IsEnabledByTime(g.GetVerkleTransitionTime, &g.Timestamp) || g.IsEnabled(g.GetVerkleTransition, new(big.Int).SetUint64(g.Number))
}

func (g *Genesis) IsEnabledByTime(fn func() *uint64, n *uint64) bool {
	return g.Config.IsEnabledByTime(fn, n)
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
	Nonce         math.HexOrDecimal64
	Timestamp     math.HexOrDecimal64
	ExtraData     hexutil.Bytes
	GasLimit      math.HexOrDecimal64
	GasUsed       math.HexOrDecimal64
	Number        math.HexOrDecimal64
	Difficulty    *math.HexOrDecimal256
	Alloc         map[common.UnprefixedAddress]GenesisAccount
	BaseFee       *math.HexOrDecimal256
	ExcessBlobGas *math.HexOrDecimal64
	BlobGasUsed   *math.HexOrDecimal64
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
	var p []struct {
		Addr    *big.Int
		Balance *big.Int
		Misc    *struct {
			Nonce uint64
			Code  []byte
			Slots []struct {
				Key common.Hash
				Val common.Hash
			}
		} `rlp:"optional"`
	}
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		acc := GenesisAccount{Balance: account.Balance}
		if account.Misc != nil {
			acc.Nonce = account.Misc.Nonce
			acc.Code = account.Misc.Code

			acc.Storage = make(map[common.Hash]common.Hash)
			for _, slot := range account.Misc.Slots {
				acc.Storage[slot.Key] = slot.Val
			}
		}
		ga[common.BigToAddress(account.Addr)] = acc
	}
	return ga
}

// Following methods implement the ctypes.GenesisBlocker interface.

func (g *Genesis) GetSealingType() ctypes.BlockSealingT {
	return ctypes.BlockSealing_Ethereum
}

func (g *Genesis) SetSealingType(t ctypes.BlockSealingT) error {
	if t != ctypes.BlockSealing_Ethereum {
		return ctypes.ErrUnsupportedConfigFatal
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

func (g *Genesis) GetGenesisExtraData() []byte {
	return g.ExtraData
}

func (g *Genesis) SetGenesisExtraData(b []byte) error {
	g.ExtraData = b
	return nil
}

func (g *Genesis) GetGenesisGasLimit() uint64 {
	return g.GasLimit
}

func (g *Genesis) SetGenesisGasLimit(u uint64) error {
	g.GasLimit = u
	return nil
}

// Implement methods to satisfy Configurator interface.

func (g *Genesis) GetAccountStartNonce() *uint64 {
	return g.Config.GetAccountStartNonce()
}

func (g *Genesis) SetAccountStartNonce(n *uint64) error {
	return g.Config.SetAccountStartNonce(n)
}

func (g *Genesis) GetMaximumExtraDataSize() *uint64 {
	return g.Config.GetMaximumExtraDataSize()
}

func (g *Genesis) SetMaximumExtraDataSize(n *uint64) error {
	return g.Config.SetMaximumExtraDataSize(n)
}

func (g *Genesis) GetMinGasLimit() *uint64 {
	return g.Config.GetMinGasLimit()
}

func (g *Genesis) SetMinGasLimit(n *uint64) error {
	return g.Config.SetMinGasLimit(n)
}

func (g *Genesis) GetGasLimitBoundDivisor() *uint64 {
	return g.Config.GetGasLimitBoundDivisor()
}

func (g *Genesis) SetGasLimitBoundDivisor(n *uint64) error {
	return g.Config.SetGasLimitBoundDivisor(n)
}

func (g *Genesis) GetNetworkID() *uint64 {
	return g.Config.GetNetworkID()
}

func (g *Genesis) SetNetworkID(n *uint64) error {
	return g.Config.SetNetworkID(n)
}

func (g *Genesis) GetChainID() *big.Int {
	return g.Config.GetChainID()
}

func (g *Genesis) SetChainID(i *big.Int) error {
	return g.Config.SetChainID(i)
}

func (g *Genesis) GetSupportedProtocolVersions() []uint {
	return g.Config.GetSupportedProtocolVersions()
}

func (g *Genesis) SetSupportedProtocolVersions(p []uint) error {
	return g.Config.SetSupportedProtocolVersions(p)
}

func (g *Genesis) GetMaxCodeSize() *uint64 {
	return g.Config.GetMaxCodeSize()
}

func (g *Genesis) SetMaxCodeSize(n *uint64) error {
	return g.Config.SetMaxCodeSize(n)
}

func (g *Genesis) GetEIP7Transition() *uint64 {
	return g.Config.GetEIP7Transition()
}

func (g *Genesis) SetEIP7Transition(n *uint64) error {
	return g.Config.SetEIP7Transition(n)
}

func (g *Genesis) GetEIP150Transition() *uint64 {
	return g.Config.GetEIP150Transition()
}

func (g *Genesis) SetEIP150Transition(n *uint64) error {
	return g.Config.SetEIP150Transition(n)
}

func (g *Genesis) GetEIP152Transition() *uint64 {
	return g.Config.GetEIP152Transition()
}

func (g *Genesis) SetEIP152Transition(n *uint64) error {
	return g.Config.SetEIP152Transition(n)
}

func (g *Genesis) GetEIP160Transition() *uint64 {
	return g.Config.GetEIP160Transition()
}

func (g *Genesis) SetEIP160Transition(n *uint64) error {
	return g.Config.SetEIP160Transition(n)
}

func (g *Genesis) GetEIP161dTransition() *uint64 {
	return g.Config.GetEIP161dTransition()
}

func (g *Genesis) SetEIP161dTransition(n *uint64) error {
	return g.Config.SetEIP161dTransition(n)
}

func (g *Genesis) GetEIP161abcTransition() *uint64 {
	return g.Config.GetEIP161abcTransition()
}

func (g *Genesis) SetEIP161abcTransition(n *uint64) error {
	return g.Config.SetEIP161abcTransition(n)
}

func (g *Genesis) GetEIP170Transition() *uint64 {
	return g.Config.GetEIP170Transition()
}

func (g *Genesis) SetEIP170Transition(n *uint64) error {
	return g.Config.SetEIP170Transition(n)
}

func (g *Genesis) GetEIP155Transition() *uint64 {
	return g.Config.GetEIP155Transition()
}

func (g *Genesis) SetEIP155Transition(n *uint64) error {
	return g.Config.SetEIP155Transition(n)
}

func (g *Genesis) GetEIP140Transition() *uint64 {
	return g.Config.GetEIP140Transition()
}

func (g *Genesis) SetEIP140Transition(n *uint64) error {
	return g.Config.SetEIP140Transition(n)
}

func (g *Genesis) GetEIP198Transition() *uint64 {
	return g.Config.GetEIP198Transition()
}

func (g *Genesis) SetEIP198Transition(n *uint64) error {
	return g.Config.SetEIP198Transition(n)
}

func (g *Genesis) GetEIP211Transition() *uint64 {
	return g.Config.GetEIP211Transition()
}

func (g *Genesis) SetEIP211Transition(n *uint64) error {
	return g.Config.SetEIP211Transition(n)
}

func (g *Genesis) GetEIP212Transition() *uint64 {
	return g.Config.GetEIP212Transition()
}

func (g *Genesis) SetEIP212Transition(n *uint64) error {
	return g.Config.SetEIP212Transition(n)
}

func (g *Genesis) GetEIP213Transition() *uint64 {
	return g.Config.GetEIP213Transition()
}

func (g *Genesis) SetEIP213Transition(n *uint64) error {
	return g.Config.SetEIP213Transition(n)
}

func (g *Genesis) GetEIP214Transition() *uint64 {
	return g.Config.GetEIP214Transition()
}

func (g *Genesis) SetEIP214Transition(n *uint64) error {
	return g.Config.SetEIP214Transition(n)
}

func (g *Genesis) GetEIP658Transition() *uint64 {
	return g.Config.GetEIP658Transition()
}

func (g *Genesis) SetEIP658Transition(n *uint64) error {
	return g.Config.SetEIP658Transition(n)
}

func (g *Genesis) GetEIP145Transition() *uint64 {
	return g.Config.GetEIP145Transition()
}

func (g *Genesis) SetEIP145Transition(n *uint64) error {
	return g.Config.SetEIP145Transition(n)
}

func (g *Genesis) GetEIP1014Transition() *uint64 {
	return g.Config.GetEIP1014Transition()
}

func (g *Genesis) SetEIP1014Transition(n *uint64) error {
	return g.Config.SetEIP1014Transition(n)
}

func (g *Genesis) GetEIP1052Transition() *uint64 {
	return g.Config.GetEIP1052Transition()
}

func (g *Genesis) SetEIP1052Transition(n *uint64) error {
	return g.Config.SetEIP1052Transition(n)
}

func (g *Genesis) GetEIP1283Transition() *uint64 {
	return g.Config.GetEIP1283Transition()
}

func (g *Genesis) SetEIP1283Transition(n *uint64) error {
	return g.Config.SetEIP1283Transition(n)
}

func (g *Genesis) GetEIP1283DisableTransition() *uint64 {
	return g.Config.GetEIP1283DisableTransition()
}

func (g *Genesis) SetEIP1283DisableTransition(n *uint64) error {
	return g.Config.SetEIP1283DisableTransition(n)
}

func (g *Genesis) GetEIP1108Transition() *uint64 {
	return g.Config.GetEIP1108Transition()
}

func (g *Genesis) SetEIP1108Transition(n *uint64) error {
	return g.Config.SetEIP1108Transition(n)
}

func (g *Genesis) GetEIP2200Transition() *uint64 {
	return g.Config.GetEIP2200Transition()
}

func (g *Genesis) SetEIP2200Transition(n *uint64) error {
	return g.Config.SetEIP2200Transition(n)
}

func (g *Genesis) GetEIP2200DisableTransition() *uint64 {
	return g.Config.GetEIP2200DisableTransition()
}

func (g *Genesis) SetEIP2200DisableTransition(n *uint64) error {
	return g.Config.SetEIP2200DisableTransition(n)
}

func (g *Genesis) GetEIP1344Transition() *uint64 {
	return g.Config.GetEIP1344Transition()
}

func (g *Genesis) SetEIP1344Transition(n *uint64) error {
	return g.Config.SetEIP1344Transition(n)
}

func (g *Genesis) GetEIP1884Transition() *uint64 {
	return g.Config.GetEIP1884Transition()
}

func (g *Genesis) SetEIP1884Transition(n *uint64) error {
	return g.Config.SetEIP1884Transition(n)
}

func (g *Genesis) GetEIP2028Transition() *uint64 {
	return g.Config.GetEIP2028Transition()
}

func (g *Genesis) SetEIP2028Transition(n *uint64) error {
	return g.Config.SetEIP2028Transition(n)
}

func (g *Genesis) GetECIP1080Transition() *uint64 {
	return g.Config.GetECIP1080Transition()
}

func (g *Genesis) SetECIP1080Transition(n *uint64) error {
	return g.Config.SetECIP1080Transition(n)
}

func (g *Genesis) GetEIP1706Transition() *uint64 {
	return g.Config.GetEIP1706Transition()
}

func (g *Genesis) SetEIP1706Transition(n *uint64) error {
	return g.Config.SetEIP1706Transition(n)
}

func (g *Genesis) GetEIP2537Transition() *uint64 {
	return g.Config.GetEIP2537Transition()
}

func (g *Genesis) SetEIP2537Transition(n *uint64) error {
	return g.Config.SetEIP2537Transition(n)
}

func (g *Genesis) GetEIP2315Transition() *uint64 {
	return g.Config.GetEIP2315Transition()
}

func (g *Genesis) SetEIP2315Transition(n *uint64) error {
	return g.Config.SetEIP2315Transition(n)
}

func (g *Genesis) GetEIP2929Transition() *uint64 {
	return g.Config.GetEIP2929Transition()
}

func (g *Genesis) SetEIP2929Transition(n *uint64) error {
	return g.Config.SetEIP2929Transition(n)
}

func (g *Genesis) GetEIP2930Transition() *uint64 {
	return g.Config.GetEIP2930Transition()
}

func (g *Genesis) SetEIP2930Transition(n *uint64) error {
	return g.Config.SetEIP2930Transition(n)
}

func (g *Genesis) GetEIP1559Transition() *uint64 {
	return g.Config.GetEIP1559Transition()
}

func (g *Genesis) SetEIP1559Transition(n *uint64) error {
	return g.Config.SetEIP1559Transition(n)
}

func (g *Genesis) GetEIP3541Transition() *uint64 {
	return g.Config.GetEIP3541Transition()
}

func (g *Genesis) SetEIP3541Transition(n *uint64) error {
	return g.Config.SetEIP3541Transition(n)
}

func (g *Genesis) GetEIP3529Transition() *uint64 {
	return g.Config.GetEIP3529Transition()
}

func (g *Genesis) SetEIP3529Transition(n *uint64) error {
	return g.Config.SetEIP3529Transition(n)
}

func (g *Genesis) GetEIP3198Transition() *uint64 {
	return g.Config.GetEIP3198Transition()
}

func (g *Genesis) SetEIP3198Transition(n *uint64) error {
	return g.Config.SetEIP3198Transition(n)
}

func (g *Genesis) GetEIP2565Transition() *uint64 {
	return g.Config.GetEIP2565Transition()
}

func (g *Genesis) SetEIP2565Transition(n *uint64) error {
	return g.Config.SetEIP2565Transition(n)
}

func (g *Genesis) GetEIP2718Transition() *uint64 {
	return g.Config.GetEIP2718Transition()
}

func (g *Genesis) SetEIP2718Transition(n *uint64) error {
	return g.Config.SetEIP2718Transition(n)
}

func (g *Genesis) GetEIP4399Transition() *uint64 {
	return g.Config.GetEIP4399Transition()
}

func (g *Genesis) SetEIP4399Transition(n *uint64) error {
	return g.Config.SetEIP4399Transition(n)
}

func (g *Genesis) GetMergeVirtualTransition() *uint64 {
	return g.Config.GetMergeVirtualTransition()
}

func (g *Genesis) SetMergeVirtualTransition(n *uint64) error {
	return g.Config.SetMergeVirtualTransition(n)
}

func (g *Genesis) GetECBP1100Transition() *uint64 {
	return g.Config.GetECBP1100Transition()
}

func (g *Genesis) SetECBP1100Transition(n *uint64) error {
	return g.Config.SetECBP1100Transition(n)
}

func (g *Genesis) GetECBP1100DeactivateTransition() *uint64 {
	return g.Config.GetECBP1100DeactivateTransition()
}

func (g *Genesis) SetECBP1100DeactivateTransition(n *uint64) error {
	return g.Config.SetECBP1100DeactivateTransition(n)
}

func (g *Genesis) IsEnabled(fn func() *uint64, n *big.Int) bool {
	return g.Config.IsEnabled(fn, n)
}

func (g *Genesis) GetForkCanonHash(n uint64) common.Hash {
	return g.Config.GetForkCanonHash(n)
}

func (g *Genesis) SetForkCanonHash(n uint64, h common.Hash) error {
	return g.Config.SetForkCanonHash(n, h)
}

func (g *Genesis) GetForkCanonHashes() map[uint64]common.Hash {
	return g.Config.GetForkCanonHashes()
}

func (g *Genesis) GetConsensusEngineType() ctypes.ConsensusEngineT {
	return g.Config.GetConsensusEngineType()
}

func (g *Genesis) MustSetConsensusEngineType(t ctypes.ConsensusEngineT) error {
	return g.Config.MustSetConsensusEngineType(t)
}

func (g *Genesis) GetIsDevMode() bool {
	return g.Config.GetIsDevMode()
}

func (g *Genesis) SetDevMode(devMode bool) error {
	return g.Config.SetDevMode(devMode)
}

func (g *Genesis) GetEthashTerminalTotalDifficulty() *big.Int {
	return g.Config.GetEthashTerminalTotalDifficulty()
}

func (g *Genesis) SetEthashTerminalTotalDifficulty(n *big.Int) error {
	return g.Config.SetEthashTerminalTotalDifficulty(n)
}

func (g *Genesis) GetEthashTerminalTotalDifficultyPassed() bool {
	return g.Config.GetEthashTerminalTotalDifficultyPassed()
}

func (g *Genesis) SetEthashTerminalTotalDifficultyPassed(t bool) error {
	return g.Config.SetEthashTerminalTotalDifficultyPassed(t)
}

// IsTerminalPoWBlock returns whether the given block is the last block of PoW stage.
func (g *Genesis) IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool {
	terminalTotalDifficulty := g.Config.GetEthashTerminalTotalDifficulty()
	if terminalTotalDifficulty == nil {
		return false
	}
	return parentTotalDiff.Cmp(terminalTotalDifficulty) < 0 && totalDiff.Cmp(terminalTotalDifficulty) >= 0
}

func (g *Genesis) GetEthashMinimumDifficulty() *big.Int {
	return g.Config.GetEthashMinimumDifficulty()
}

func (g *Genesis) SetEthashMinimumDifficulty(i *big.Int) error {
	return g.Config.SetEthashMinimumDifficulty(i)
}

func (g *Genesis) GetEthashDifficultyBoundDivisor() *big.Int {
	return g.Config.GetEthashDifficultyBoundDivisor()
}

func (g *Genesis) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	return g.Config.SetEthashDifficultyBoundDivisor(i)
}

func (g *Genesis) GetEthashDurationLimit() *big.Int {
	return g.Config.GetEthashDurationLimit()
}

func (g *Genesis) SetEthashDurationLimit(i *big.Int) error {
	return g.Config.SetEthashDurationLimit(i)
}

func (g *Genesis) GetEthashHomesteadTransition() *uint64 {
	return g.Config.GetEthashHomesteadTransition()
}

func (g *Genesis) SetEthashHomesteadTransition(n *uint64) error {
	return g.Config.SetEthashHomesteadTransition(n)
}

func (g *Genesis) GetEIP2Transition() *uint64 {
	return g.Config.GetEIP2Transition()
}

func (g *Genesis) SetEIP2Transition(n *uint64) error {
	return g.Config.SetEIP2Transition(n)
}

func (g *Genesis) GetEthashEIP779Transition() *uint64 {
	return g.Config.GetEthashEIP779Transition()
}

func (g *Genesis) SetEthashEIP779Transition(n *uint64) error {
	return g.Config.SetEthashEIP779Transition(n)
}

func (g *Genesis) GetEthashEIP649Transition() *uint64 {
	return g.Config.GetEthashEIP649Transition()
}

func (g *Genesis) SetEthashEIP649Transition(n *uint64) error {
	return g.Config.SetEthashEIP649Transition(n)
}

func (g *Genesis) GetEthashEIP1234Transition() *uint64 {
	return g.Config.GetEthashEIP1234Transition()
}

func (g *Genesis) SetEthashEIP1234Transition(n *uint64) error {
	return g.Config.SetEthashEIP1234Transition(n)
}

func (g *Genesis) GetEthashEIP2384Transition() *uint64 {
	return g.Config.GetEthashEIP2384Transition()
}

func (g *Genesis) SetEthashEIP2384Transition(n *uint64) error {
	return g.Config.SetEthashEIP2384Transition(n)
}

func (g *Genesis) GetEthashEIP3554Transition() *uint64 {
	return g.Config.GetEthashEIP3554Transition()
}

func (g *Genesis) SetEthashEIP3554Transition(n *uint64) error {
	return g.Config.SetEthashEIP3554Transition(n)
}

func (g *Genesis) GetEthashEIP4345Transition() *uint64 {
	return g.Config.GetEthashEIP4345Transition()
}

func (g *Genesis) SetEthashEIP4345Transition(n *uint64) error {
	return g.Config.SetEthashEIP4345Transition(n)
}

func (g *Genesis) GetEthashEIP5133Transition() *uint64 {
	return g.Config.GetEthashEIP5133Transition()
}

func (g *Genesis) SetEthashEIP5133Transition(n *uint64) error {
	return g.Config.SetEthashEIP5133Transition(n)
}

func (g *Genesis) GetEthashECIP1010PauseTransition() *uint64 {
	return g.Config.GetEthashECIP1010PauseTransition()
}

func (g *Genesis) SetEthashECIP1010PauseTransition(n *uint64) error {
	return g.Config.SetEthashECIP1010PauseTransition(n)
}

func (g *Genesis) GetEthashECIP1010ContinueTransition() *uint64 {
	return g.Config.GetEthashECIP1010ContinueTransition()
}

func (g *Genesis) SetEthashECIP1010ContinueTransition(n *uint64) error {
	return g.Config.SetEthashECIP1010ContinueTransition(n)
}

func (g *Genesis) GetEthashECIP1017Transition() *uint64 {
	return g.Config.GetEthashECIP1017Transition()
}

func (g *Genesis) SetEthashECIP1017Transition(n *uint64) error {
	return g.Config.SetEthashECIP1017Transition(n)
}

func (g *Genesis) GetEthashECIP1017EraRounds() *uint64 {
	return g.Config.GetEthashECIP1017EraRounds()
}

func (g *Genesis) SetEthashECIP1017EraRounds(n *uint64) error {
	return g.Config.SetEthashECIP1017EraRounds(n)
}

func (g *Genesis) GetEthashEIP100BTransition() *uint64 {
	return g.Config.GetEthashEIP100BTransition()
}

func (g *Genesis) SetEthashEIP100BTransition(n *uint64) error {
	return g.Config.SetEthashEIP100BTransition(n)
}

func (g *Genesis) GetEthashECIP1041Transition() *uint64 {
	return g.Config.GetEthashECIP1041Transition()
}

func (g *Genesis) SetEthashECIP1041Transition(n *uint64) error {
	return g.Config.SetEthashECIP1041Transition(n)
}

func (g *Genesis) GetEthashECIP1099Transition() *uint64 {
	return g.Config.GetEthashECIP1099Transition()
}

func (g *Genesis) SetEthashECIP1099Transition(n *uint64) error {
	return g.Config.SetEthashECIP1099Transition(n)
}

func (g *Genesis) GetEthashDifficultyBombDelaySchedule() ctypes.Uint64Uint256MapEncodesHex {
	return g.Config.GetEthashDifficultyBombDelaySchedule()
}

func (g *Genesis) SetEthashDifficultyBombDelaySchedule(m ctypes.Uint64Uint256MapEncodesHex) error {
	return g.Config.SetEthashDifficultyBombDelaySchedule(m)
}

func (g *Genesis) GetEthashBlockRewardSchedule() ctypes.Uint64Uint256MapEncodesHex {
	return g.Config.GetEthashBlockRewardSchedule()
}

func (g *Genesis) SetEthashBlockRewardSchedule(m ctypes.Uint64Uint256MapEncodesHex) error {
	return g.Config.SetEthashBlockRewardSchedule(m)
}

func (g *Genesis) GetCliquePeriod() uint64 {
	return g.Config.GetCliquePeriod()
}

func (g *Genesis) SetCliquePeriod(n uint64) error {
	return g.Config.SetCliquePeriod(n)
}

func (g *Genesis) GetCliqueEpoch() uint64 {
	return g.Config.GetCliqueEpoch()
}

func (g *Genesis) SetCliqueEpoch(n uint64) error {
	return g.Config.SetCliqueEpoch(n)
}

func (g *Genesis) GetLyra2NonceTransition() *uint64 {
	return g.Config.GetLyra2NonceTransition()
}

func (g *Genesis) SetLyra2NonceTransition(n *uint64) error {
	return g.Config.SetLyra2NonceTransition(n)
}

func (g *Genesis) String() string {
	j, _ := json.MarshalIndent(g, "", "    ")
	return "Genesis: " + string(j)
}
