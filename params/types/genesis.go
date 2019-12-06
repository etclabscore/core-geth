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

// MultiGeth: Manually override this file generation; it requires manual editing. The
// git patch for the edited difference is supplied at gen_genesis_patch.diff.
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

func (g *Genesis) GetEIP161abcTransition() *uint64 {
	return g.Config.GetEIP161abcTransition()
}

func (g *Genesis) SetEIP161abcTransition(n *uint64) error {
	return g.Config.SetEIP161abcTransition(n)
}

func (g *Genesis) GetEIP161dTransition() *uint64 {
	return g.Config.GetEIP161dTransition()
}

func (g *Genesis) SetEIP161dTransition(n *uint64) error {
	return g.Config.SetEIP161dTransition(n)
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

func (g *Genesis) GetEIP1283ReenableTransition() *uint64 {
	return g.Config.GetEIP1283ReenableTransition()
}

func (g *Genesis) SetEIP1283ReenableTransition(n *uint64) error {
	return g.Config.SetEIP1283ReenableTransition(n)
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

func (g *Genesis) IsForked(fn func() *uint64, n *big.Int) bool {
	return g.Config.IsForked(fn, n)
}

func (g *Genesis) ForkCanonHash(n uint64) common.Hash {
	return g.Config.ForkCanonHash(n)
}

func (g *Genesis) GetConsensusEngineType() common2.ConsensusEngineT {
	return g.Config.GetConsensusEngineType()
}

func (g *Genesis) MustSetConsensusEngineType(t common2.ConsensusEngineT) error {
	return g.Config.MustSetConsensusEngineType(t)
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

func (g *Genesis) GetEthashEIP2Transition() *uint64 {
	return g.Config.GetEthashEIP2Transition()
}

func (g *Genesis) SetEthashEIP2Transition(n *uint64) error {
	return g.Config.SetEthashEIP2Transition(n)
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

func (g *Genesis) GetEthashDifficultyBombDelaySchedule() common2.Uint64BigMapEncodesHex {
	return g.Config.GetEthashDifficultyBombDelaySchedule()
}

func (g *Genesis) SetEthashDifficultyBombDelaySchedule(m common2.Uint64BigMapEncodesHex) error {
	return g.Config.SetEthashDifficultyBombDelaySchedule(m)
}

func (g *Genesis) GetEthashBlockRewardSchedule() common2.Uint64BigMapEncodesHex {
	return g.Config.GetEthashBlockRewardSchedule()
}

func (g *Genesis) SetEthashBlockRewardSchedule(m common2.Uint64BigMapEncodesHex) error {
	return g.Config.SetEthashBlockRewardSchedule(m)
}

func (g *Genesis) GetCliquePeriod() *uint64 {
	return g.Config.GetCliquePeriod()
}

func (g *Genesis) SetCliquePeriod(n uint64) error {
	return g.Config.SetCliquePeriod(n)
}

func (g *Genesis) GetCliqueEpoch() *uint64 {
	return g.Config.GetCliqueEpoch()
}

func (g *Genesis) SetCliqueEpoch(n uint64) error {
	return g.Config.SetCliqueEpoch(n)
}