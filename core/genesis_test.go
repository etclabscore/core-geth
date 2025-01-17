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
	"bytes"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/vars"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/pathdb"
)

func TestSetupGenesisBlock(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	defaultGenesisBlock := params.DefaultGenesisBlock()

	config, hash, err := SetupGenesisBlock(db, triedb.NewDatabase(db, nil), defaultGenesisBlock)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if wantHash := GenesisToBlock(defaultGenesisBlock, nil).Hash(); wantHash != hash {
		t.Errorf("mismatch block hash, want: %x, got: %x", wantHash, hash)
	}
	if diffs := confp.Equal(reflect.TypeOf((*ctypes.ChainConfigurator)(nil)), defaultGenesisBlock.Config, config); len(diffs) != 0 {
		for _, diff := range diffs {
			t.Error("mismatch", "diff=", diff, "in", defaultGenesisBlock.Config, "out", config)
		}
	}

	classicGenesisBlock := params.DefaultClassicGenesisBlock()

	clConfig, clHash, clErr := SetupGenesisBlock(db, triedb.NewDatabase(db, nil), classicGenesisBlock)
	if clErr != nil {
		t.Errorf("err: %v", clErr)
	}
	if wantHash := GenesisToBlock(classicGenesisBlock, nil).Hash(); wantHash != clHash {
		t.Errorf("mismatch block hash, want: %x, got: %x", wantHash, clHash)
	}
	if diffs := confp.Equal(reflect.TypeOf((*ctypes.ChainConfigurator)(nil)), classicGenesisBlock.Config, clConfig); len(diffs) != 0 {
		for _, diff := range diffs {
			t.Error("mismatch", "diff=", diff, "in", classicGenesisBlock.Config, "out", clConfig)
		}
	}
}

func TestSetupGenesis(t *testing.T) {
	testSetupGenesis(t, rawdb.HashScheme)
	testSetupGenesis(t, rawdb.PathScheme)
}

func testSetupGenesis(t *testing.T, scheme string) {
	var (
		customghash = common.HexToHash("0x89c99d90b79719238d2645c7642f2c9295246e80775b38cfd162b696817fbd50")
		customg     = genesisT.Genesis{
			Config: &goethereum.ChainConfig{HomesteadBlock: big.NewInt(3)},
			Alloc: genesisT.GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			},
		}
		oldcustomg = customg
	)
	oldcustomg.Config = &goethereum.ChainConfig{HomesteadBlock: big.NewInt(2)}
	tests := []struct {
		name       string
		fn         func(ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error)
		wantConfig ctypes.ChainConfigurator
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				return SetupGenesisBlock(db, triedb.NewDatabase(db, newDbConfig(scheme)), new(genesisT.Genesis))
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.AllEthashProtocolChanges,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				return SetupGenesisBlock(db, triedb.NewDatabase(db, newDbConfig(scheme)), nil)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "mainnet block in DB, genesis == nil",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				MustCommitGenesis(db, triedb.NewDatabase(db, nil), params.DefaultGenesisBlock())
				return SetupGenesisBlock(db, triedb.NewDatabase(db, newDbConfig(scheme)), nil)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				MustCommitGenesis(db, tdb, &customg)
				return SetupGenesisBlock(db, tdb, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "custom block in DB, genesis == mordor",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				MustCommitGenesis(db, tdb, &customg)
				return SetupGenesisBlock(db, tdb, params.DefaultMordorGenesisBlock())
			},
			wantErr:    &genesisT.GenesisMismatchError{Stored: customghash, New: params.MordorGenesisHash},
			wantHash:   params.MordorGenesisHash,
			wantConfig: params.MordorChainConfig,
		},
		{
			name: "compatible config in DB",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				MustCommitGenesis(db, triedb.NewDatabase(db, nil), &oldcustomg)
				return SetupGenesisBlock(db, tdb, &customg)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "incompatible config in DB",
			fn: func(db ethdb.Database) (ctypes.ChainConfigurator, common.Hash, error) {
				// Commit the 'old' genesis block with Homestead transition at #2.
				// Advance to block #4, past the homestead transition block of customg.
				tdb := triedb.NewDatabase(db, newDbConfig(scheme))
				MustCommitGenesis(db, tdb, &oldcustomg)

				bc, _ := NewBlockChain(db, DefaultCacheConfigWithScheme(scheme), &oldcustomg, nil, ethash.NewFullFaker(), vm.Config{}, nil, nil)
				defer bc.Stop()

				_, blocks, _ := GenerateChainWithGenesis(&oldcustomg, ethash.NewFaker(), 4, nil)
				bc.InsertChain(blocks)

				// This should return a compatibility error.
				return SetupGenesisBlock(db, tdb, &customg)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
			wantErr: &confp.ConfigCompatError{
				What:          "incompatible fork value: GetEIP2Transition",
				StoredBlock:   big.NewInt(2),
				NewBlock:      big.NewInt(3),
				RewindToBlock: 1,
			},
		},
	}

	for _, test := range tests {
		db := rawdb.NewMemoryDatabase()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := rawdb.ReadBlock(db, test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

// TestGenesisHashes checks the congruity of default genesis data to
// corresponding hardcoded genesis hash values.
func TestGenesisHashes(t *testing.T) {
	for i, c := range []struct {
		genesis *genesisT.Genesis
		want    common.Hash
	}{
		{params.DefaultGenesisBlock(), params.MainnetGenesisHash},
		{params.DefaultMordorGenesisBlock(), params.MordorGenesisHash},
		{params.DefaultSepoliaGenesisBlock(), params.SepoliaGenesisHash},
	} {
		// Test via MustCommit
		db := rawdb.NewMemoryDatabase()
		if have := MustCommitGenesis(rawdb.NewMemoryDatabase(), triedb.NewDatabase(db, triedb.HashDefaults), c.genesis).Hash(); have != c.want {
			t.Errorf("case: %d a), want: %s, got: %s", i, c.want.Hex(), have.Hex())
		}
		// TODO(meowsbits): go-ethereum has an additional Test via ToBlock. Is there a comparable method that we should also test here?
	}
}

func TestGenesis_Commit(t *testing.T) {
	genesis := &genesisT.Genesis{
		BaseFee: big.NewInt(vars.InitialBaseFee),
		Config:  params.TestChainConfig,
		// difficulty is nil
	}

	db := rawdb.NewMemoryDatabase()
	genesisBlock := MustCommitGenesis(db, triedb.NewDatabase(db, triedb.HashDefaults), genesis)

	if genesis.Difficulty != nil {
		t.Fatalf("assumption wrong")
	}

	// This value should have been set as default in the ToBlock method.
	if genesisBlock.Difficulty().Cmp(vars.GenesisDifficulty) != 0 {
		t.Errorf("assumption wrong: want: %d, got: %v", vars.GenesisDifficulty, genesisBlock.Difficulty())
	}

	// Expect the stored total difficulty to be the difficulty of the genesis block.
	stored := rawdb.ReadTd(db, genesisBlock.Hash(), genesisBlock.NumberU64())

	if stored.Cmp(genesisBlock.Difficulty()) != 0 {
		t.Errorf("inequal difficulty; stored: %v, genesisBlock: %v", stored, genesisBlock.Difficulty())
	}
}

func TestReadWriteGenesisAlloc(t *testing.T) {
	var (
		db    = rawdb.NewMemoryDatabase()
		alloc = &genesisT.GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			{2}: {Balance: big.NewInt(2), Storage: map[common.Hash]common.Hash{{2}: {2}}},
		}
		hash, _ = gaHash(alloc, false)
	)
	blob, _ := json.Marshal(alloc)
	rawdb.WriteGenesisStateSpec(db, hash, blob)

	var reload genesisT.GenesisAlloc
	err := reload.UnmarshalJSON(rawdb.ReadGenesisStateSpec(db, hash))
	if err != nil {
		t.Fatalf("Failed to load genesis state %v", err)
	}
	if len(reload) != len(*alloc) {
		t.Fatal("Unexpected genesis allocation")
	}
	for addr, account := range reload {
		want, ok := (*alloc)[addr]
		if !ok {
			t.Fatal("Account is not found")
		}
		if !reflect.DeepEqual(want, account) {
			t.Fatal("Unexpected account")
		}
	}
}

func newDbConfig(scheme string) *triedb.Config {
	if scheme == rawdb.HashScheme {
		return triedb.HashDefaults
	}
	return &triedb.Config{PathDB: pathdb.Defaults}
}

func TestVerkleGenesisCommit(t *testing.T) {
	var verkleTime uint64 = 0
	verkleConfig := &goethereum.ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  &verkleTime,
		CancunTime:                    &verkleTime,
		PragueTime:                    &verkleTime,
		VerkleTime:                    &verkleTime,
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
		Ethash:                        nil,
		Clique:                        nil,
	}

	genesis := &genesisT.Genesis{
		BaseFee:    big.NewInt(vars.InitialBaseFee),
		Config:     verkleConfig,
		Timestamp:  verkleTime,
		Difficulty: big.NewInt(0),
		Alloc: genesisT.GenesisAlloc{
			{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
		},
	}

	db := rawdb.NewMemoryDatabase()

	expected := common.Hex2Bytes("14398d42be3394ff8d50681816a4b7bf8d8283306f577faba2d5bc57498de23b")
	genesisBlock := MustCommitGenesis(db, triedb.NewDatabase(db, triedb.HashDefaults), genesis)
	got := genesisBlock.Root().Bytes()
	if !bytes.Equal(got, expected) {
		t.Fatalf("invalid genesis state root, expected %x, got %x", expected, got)
	}

	triedb := triedb.NewDatabase(db, &triedb.Config{IsVerkle: true, PathDB: pathdb.Defaults})
	block := MustCommitGenesis(db, triedb, genesis)
	if !bytes.Equal(block.Root().Bytes(), expected) {
		t.Fatalf("invalid genesis state root, expected %x, got %x", expected, got)
	}

	// Test that the trie is verkle
	if !triedb.IsVerkle() {
		t.Fatalf("expected trie to be verkle")
	}

	if !rawdb.ExistsAccountTrieNode(db, nil) {
		t.Fatal("could not find node")
	}
}
