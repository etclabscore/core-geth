// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/trie"
)

func uint64P(n uint64) *uint64 {
	return &n
}

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new   ctypes.ChainConfigurator
		headBlock     uint64
		headTimestamp uint64
		wantErr       *confp.ConfigCompatError
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 0, headTimestamp: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 0, headTimestamp: uint64(time.Now().Unix()), wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, headBlock: 100, wantErr: nil},
		{
			stored:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), EIP150Block: big.NewInt(10)},
			new:       &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), EIP150Block: big.NewInt(20)},
			headBlock: 9,
			wantErr:   nil,
		},
		// case index 4
		{
			stored:    AllEthashProtocolChanges,
			new:       &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: nil},
			headBlock: 3,
			wantErr: &confp.ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      nil,
				RewindToBlock: 0,
			},
		},
		// 5
		{
			stored:    AllEthashProtocolChanges,
			new:       &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(1)},
			headBlock: 3,
			wantErr: &confp.ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      big.NewInt(1),
				RewindToBlock: 0,
			},
		},
		// 6
		// Want the EIP150 block error because it is below the head block AND below the Homestead blocks.
		// The compat error should always be the earliest fork block incompatibility.
		{
			stored:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:       &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			headBlock: 25,
			wantErr: &confp.ConfigCompatError{
				What:          "EIP150 fork block",
				StoredBlock:   big.NewInt(10),
				NewBlock:      big.NewInt(20),
				RewindToBlock: 9,
			},
		},
		// 7
		{
			stored:    &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(30), EIP649FBlock: big.NewInt(30)},
			new:       &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(24), EIP649FBlock: big.NewInt(24)},
			headBlock: 25,
			wantErr: &confp.ConfigCompatError{
				What:          "EIP100F fork block",
				StoredBlock:   big.NewInt(30),
				NewBlock:      big.NewInt(24),
				RewindToBlock: 23,
			},
		},
		// 8
		{
			stored:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), ByzantiumBlock: big.NewInt(30)},
			new:       &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP211FBlock: big.NewInt(26)},
			headBlock: 25,
			wantErr:   nil,
		},
		// 9
		{
			stored:    &goethereum.ChainConfig{Ethash: new(ctypes.EthashConfig), ByzantiumBlock: big.NewInt(30)},
			new:       &coregeth.CoreGethChainConfig{Ethash: new(ctypes.EthashConfig), EIP100FBlock: big.NewInt(26), EIP649FBlock: big.NewInt(26)},
			headBlock: 25,
			wantErr:   nil,
		},
		// 10
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				if err := confp.Crush(c, MainnetChainConfig, true); err != nil {
					panic(err)
				}
				c.SetEthashEIP779Transition(uint64P(1900000))
				return c
			}(),
			headBlock: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:          "DAO fork support flag",
				StoredBlock:   MainnetChainConfig.DAOForkBlock,
				NewBlock:      big.NewInt(1900000),
				RewindToBlock: 1900000 - 1,
			},
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				confp.Crush(c, MainnetChainConfig, true)
				c.SetEthashEIP779Transition(nil)
				return c
			}(),
			headBlock: MainnetChainConfig.DAOForkBlock.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:          "DAO fork support flag",
				StoredBlock:   MainnetChainConfig.DAOForkBlock,
				NewBlock:      nil,
				RewindToBlock: 1920000 - 1,
			},
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{}
				*c = *MainnetChainConfig
				c.SetChainID(new(big.Int).Sub(MainnetChainConfig.EIP155Block, common.Big1))
				return c
			}(),
			headBlock: MainnetChainConfig.EIP158Block.Uint64(),
			wantErr: &confp.ConfigCompatError{
				What:          "EIP155 chain ID",
				StoredBlock:   MainnetChainConfig.EIP155Block,
				NewBlock:      MainnetChainConfig.EIP155Block,
				RewindToBlock: new(big.Int).Sub(MainnetChainConfig.EIP158Block, common.Big1).Uint64(),
			},
		},
		{
			stored: func() ctypes.ChainConfigurator {
				c := &goethereum.ChainConfig{
					Ethash:         new(ctypes.EthashConfig),
					DAOForkBlock:   big.NewInt(3),
					DAOForkSupport: false,
				}
				return c
			}(),
			new: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{
					Ethash:       new(ctypes.EthashConfig),
					DAOForkBlock: nil,
				}
				return c
			}(),
			headBlock: 5,
			wantErr:   nil,
		},
		{
			// v1.9.5 -> v1.9.7
			stored: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				*c = *ClassicChainConfig
				c.SetEIP145Transition(nil)
				c.SetEIP1014Transition(nil)
				c.SetEIP1052Transition(nil)
				c.SetEIP152Transition(nil)
				c.SetEIP1108Transition(nil)
				c.SetEIP1344Transition(nil)
				// c.SetEIP1884Transition(nil)
				c.SetEIP2028Transition(nil)
				c.SetEIP2200Transition(nil)
				return c
			}(),
			new:       ClassicChainConfig,
			headBlock: 9550000,
			wantErr:   nil,
		},
		{
			// v1.9.6 -> v1.9.7
			stored: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				*c = *ClassicChainConfig
				c.SetEIP152Transition(nil)
				c.SetEIP1108Transition(nil)
				c.SetEIP1344Transition(nil)
				// c.SetEIP1884Transition(nil)
				c.SetEIP2028Transition(nil)
				c.SetEIP2200Transition(nil)
				return c
			}(),
			new:       ClassicChainConfig,
			headBlock: 9550000,
			wantErr:   nil,
		},
		{
			stored: MainnetChainConfig,
			new: func() ctypes.ChainConfigurator {
				c := &coregeth.CoreGethChainConfig{}
				err := confp.Crush(c, MainnetChainConfig, true)
				if err != nil {
					panic(err)
				}
				return c
			}(),
		},
		// https://github.com/ethereum/go-ethereum/pull/21473
		// This is to enable private chains running on older Geth release 1.8.27 with Constantinople fork enabled (but not Petersburg) to apply Petersburg retroactively when upgrading to Geth 1.9.
		// ... but @meowsbits thinks this isn't reasonable.
		// This is allowance would presume that the private chains were unaffected by the constantinople vs. petersburg
		// changes, and were safe to retroactively set the config. I do not understand how this is allowed or even desirable...
		// Are the configs equivalent? No. Do we have any observables in place to ensure that the chain data will not be
		// retroactively corrupted? No.
		// {
		// 	stored:  &goethereum.ChainConfig{ConstantinopleBlock: big.NewInt(30)},
		// 	new:     &goethereum.ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(30)},
		// 	headBlock:    40,
		// 	wantErr: nil,
		// },
		{
			stored:    &goethereum.ChainConfig{ConstantinopleBlock: big.NewInt(30)},
			new:       &goethereum.ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(31)},
			headBlock: 40,
			wantErr: &confp.ConfigCompatError{
				What:          "Petersburg fork block",
				StoredBlock:   nil,
				NewBlock:      big.NewInt(31),
				RewindToBlock: 30,
			},
		},
	}

	for i, test := range tests {
		err := confp.Compatible(new(big.Int).SetUint64(test.headBlock), nil, test.stored, test.new)
		if (err == nil && test.wantErr != nil) || (err != nil && test.wantErr == nil) {
			t.Errorf("[case %d] nil/nonnil, error mismatch:\nstored: %v\nnew: %v\nheadBlock: %v\nerr: %v\nwant: %v", i, test.stored, test.new, test.headBlock, err, test.wantErr)
		} else if err != nil && (err.RewindToBlock != test.wantErr.RewindToBlock) {
			// if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("[case %d] error mismatch:\nstored: %v\nnew: %v\nheadBlock: %v\nerr: %v\nwant: %v", i, test.stored, test.new, test.headBlock, err, test.wantErr)
		}
	}
}

func TestFoundationIsForked(t *testing.T) {
	c := MainnetChainConfig
	if !c.IsEnabled(c.GetEthashEIP2384Transition, big.NewInt(9200001)) {
		t.Fatal("nofork muir bad")
	}
}

func TestClassicIs649(t *testing.T) {
	c := ClassicChainConfig
	got := c.GetEthashEIP649Transition()
	if got != nil {
		t.Fatal("classic config doesn't support 649; difficulty bomb was disposed of")
	}
}

func TestFoundationIsEIP779(t *testing.T) {
	blockNumbers := []*big.Int{
		big.NewInt(0),
		big.NewInt(1_920_000),
		big.NewInt(10_000_000),
	}
	for _, bn := range blockNumbers {
		if bn.Cmp(big.NewInt(0)) > 0 && !MainnetChainConfig.IsEnabled(MainnetChainConfig.GetEthashEIP779Transition, bn) {
			t.Fatal("bad")
		}
		if *MainnetChainConfig.GetEthashEIP779Transition() != 1_920_000 {
			t.Fatal("bad")
		}
	}
}

// genesisToBlock is a helper function duplicating the logic at core.MustCommitGenesis,
// which cannot be imported into this package.
func genesisToBlock(g *genesisT.Genesis, db ethdb.Database) *types.Block {
	if db == nil {
		db = rawdb.NewMemoryDatabase()
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root := statedb.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       g.Timestamp,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = vars.GenesisGasLimit
	}
	if g.Difficulty == nil {
		head.Difficulty = vars.GenesisDifficulty
	}
	statedb.Commit(head.Number.Uint64(), false)
	statedb.Database().TrieDB().Commit(root, true)

	return types.NewBlock(head, nil, nil, nil, trie.NewStackTrie(nil))
}

func TestGenesisHashes(t *testing.T) {
	cases := []struct {
		genesis *genesisT.Genesis
		hash    common.Hash
	}{
		{
			genesis: DefaultClassicGenesisBlock(),
			hash:    MainnetGenesisHash,
		},
		{
			genesis: DefaultGoerliGenesisBlock(),
			hash:    GoerliGenesisHash,
		},
	}
	for i, c := range cases {
		b := genesisToBlock(c.genesis, rawdb.NewMemoryDatabase())
		if got := b.Hash(); got != c.hash {
			t.Errorf("case: %d, want: %s, got: %s", i, c.hash.Hex(), got.Hex())
		}
	}
}
