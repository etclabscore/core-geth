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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
	"github.com/ethereum/go-ethereum/params/types/oldmultigeth"
)

func TestSetupGenesisBlock(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	defaultGenesisBlock := params.DefaultGenesisBlock()

	config, hash, err := SetupGenesisBlock(db, defaultGenesisBlock)
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

	clConfig, clHash, clErr := SetupGenesisBlock(db, classicGenesisBlock)
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

func TestSetupGenesisBlockOldVsNewMultigeth(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	// Setup a genesis mocking <=v1.9.6, aka "old".
	genA := params.DefaultGenesisBlock()
	genA.Config = &oldmultigeth.ChainConfig{
		ChainID:              big.NewInt(61),
		HomesteadBlock:       big.NewInt(1150000),
		DAOForkBlock:         big.NewInt(1920000),
		DAOForkSupport:       false,
		EIP150Block:          big.NewInt(2500000),
		EIP150Hash:           common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		EIP155Block:          big.NewInt(3000000),
		EIP158Block:          big.NewInt(8772000),
		ByzantiumBlock:       big.NewInt(8772000),
		DisposalBlock:        big.NewInt(5900000),
		SocialBlock:          nil,
		EthersocialBlock:     nil,
		ConstantinopleBlock:  big.NewInt(9573000),
		PetersburgBlock:      big.NewInt(9573000),
		IstanbulBlock:        big.NewInt(10500839),
		EIP1884DisableFBlock: big.NewInt(10500839),
		ECIP1017EraRounds:    big.NewInt(5000000),
		EIP160FBlock:         big.NewInt(3000000),
		ECIP1010PauseBlock:   big.NewInt(3000000),
		ECIP1010Length:       big.NewInt(2000000),
		Ethash:               new(ctypes.EthashConfig),
	}

	// Set it up.
	config, hash, err := SetupGenesisBlock(db, genA)
	if err != nil {
		t.Fatal(err)
	}

	// Capture and debug-log the marshaled stored config.
	b, _ := json.MarshalIndent(config, "", "    ")
	//t.Log(string(b))

	// Read the stored config manually.
	stored := rawdb.ReadCanonicalHash(db, 0)
	storedConfig := rawdb.ReadChainConfig(db, stored)

	b2, _ := json.MarshalIndent(storedConfig, "", "    ")
	//t.Log(string(b2))

	if !bytes.Equal(b, b2) {
		t.Fatal("different chain config read vs. wrote")
	}

	headHeight := uint64(9700559)
	headHash := common.HexToHash("0xe618c1b2d738dfa09052e199e5870274f09eb83c684a8a2c194b82dedc00a977")
	rawdb.WriteHeadHeaderHash(db, headHash)
	rawdb.WriteHeaderNumber(db, headHash, headHeight)

	genB := params.DefaultClassicGenesisBlock()

	newConfig, newHash, err := SetupGenesisBlock(db, genB)
	if err != nil {
		t.Fatal("incompat conf", err)
	}
	if hash != newHash {
		t.Fatal("hash mismatch")
	}

	if !confp.Identical(config, newConfig, []string{"NetworkID", "ChainID"}) {
		t.Fatal("chain config identities not same")
	}

	// These should be redundant to the SetupGenesisBlock method, but this is
	// for double double double extra sureness.
	if compatErr := confp.Compatible(&headHeight, genA, genB); compatErr != nil {
		t.Fatal(err)
	}
	if compatErr := confp.Compatible(&headHeight, config, newConfig); compatErr != nil {
		t.Fatal(err)
	}
}

func TestSetupGenesisBlock2(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	// An example of v1.9.6 multigeth config marshaled to JSON.
	// Note the fields EIP1108FBlock; these were included accidentally because
	// of a typo in the struct field json tags, and because of that, will
	// not be omitted when empty, nor "properly" (lowercase) named.
	//
	// This should be treated as an 'oldmultigeth' data type, since it has values which are
	// not present in the contemporary multigeth data type.
	//
	// In this test we'll assume that this is the config which has been
	// written to the database, and which should be superceded by the
	// config below (cc_v197_a).
	var cc_v196_a = `{
  "chainId":61,
  "homesteadBlock":1150000,
  "daoForkBlock":1920000,
  "eip150Block":2500000,
  "eip150Hash":"0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a",
  "eip155Block":3000000,
  "eip158Block":8772000,
  "byzantiumBlock":8772000,
  "constantinopleBlock":9573000,
  "petersburgBlock":9573000,
  "ethash":{

  
},
  "trustedCheckpoint":null,
  "trustedCheckpointOracle":null,
  "networkId":1,
  "eip7FBlock":null,
  "eip160Block":3000000,
  "EIP1108FBlock":null,
  "EIP1344FBlock":null,
  "EIP1884FBlock":null,
  "EIP2028FBlock":null,
  "EIP2200FBlock":null,
  "ecip1010PauseBlock":3000000,
  "ecip1010Length":2000000,
  "ecip1017FBlock":5000000,
  "ecip1017EraRounds":5000000,
  "disposalBlock":5900000
}`

	// An example of a "healthy" multigeth configuration marshaled to JSON.
	var cc_v197_a = `{
    "networkId": 1,
    "chainId": 61,
    "eip2FBlock": 1150000,
    "eip7FBlock": 1150000,
    "eip150Block": 2500000,
    "eip155Block": 3000000,
    "eip160Block": 3000000,
    "eip161FBlock": 8772000,
    "eip170FBlock": 8772000,
    "eip100FBlock": 8772000,
    "eip140FBlock": 8772000,
    "eip198FBlock": 8772000,
    "eip211FBlock": 8772000,
    "eip212FBlock": 8772000,
    "eip213FBlock": 8772000,
    "eip214FBlock": 8772000,
    "eip658FBlock": 8772000,
    "eip145FBlock": 9573000,
    "eip1014FBlock": 9573000,
    "eip1052FBlock": 9573000,
    "eip152FBlock": 10500839,
    "eip1108FBlock": 10500839,
    "eip1344FBlock": 10500839,
    "eip2028FBlock": 10500839,
    "eip2200FBlock": 10500839,
    "ecip1010PauseBlock": 3000000,
    "ecip1010Length": 2000000,
    "ecip1017FBlock": 5000000,
    "ecip1017EraRounds": 5000000,
    "disposalBlock": 5900000,
    "ethash": {},
    "trustedCheckpoint": null,
    "trustedCheckpointOracle": null,
    "requireBlockHashes": {
        "1920000": "0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f",
        "2500000": "0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"
    }
}`
	headHeight := uint64(9700559)
	genHash := common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	headHash := common.HexToHash("0xe618c1b2d738dfa09052e199e5870274f09eb83c684a8a2c194b82dedc00a977")

	_, hash, err := SetupGenesisBlock(db, params.DefaultClassicGenesisBlock())
	if err != nil {
		t.Fatal(err)
	}
	if genHash != hash {
		t.Fatal("mismatch genesis hash")
	}
	// Simulate that the stored config is the v1.9.6 version.
	// This skips the marshaling step of the rawdb.WriteChainConfig method,
	// allowing us to just slap this value in there straight.
	err = db.Put(rawdb.ConfigKey(genHash), []byte(cc_v196_a))
	if err != nil {
		t.Fatal(err)
	}

	// First test: show that the config we've stored in the database gets unmarshaled
	// as an 'oldmultigeth' config.
	storedConf := rawdb.ReadChainConfig(db, genHash)
	if storedConf == nil {
		t.Fatal("nil stored conf")
	}
	wantType := reflect.TypeOf(&oldmultigeth.ChainConfig{})
	if reflect.TypeOf(storedConf) != wantType {
		t.Fatalf("mismatch, want: %v, got: %v", wantType, reflect.TypeOf(storedConf))
	}

	// "Fast forward" the database indicators.
	rawdb.WriteHeadHeaderHash(db, headHash)
	rawdb.WriteHeaderNumber(db, headHash, headHeight)

	// Setup genesis again, but now with contemporary chain config, ie v1.9.7+
	conf2, hash2, err := SetupGenesisBlock(db, params.DefaultClassicGenesisBlock())
	if err != nil {
		t.Fatal(err)
	}
	if hash2 != hash {
		t.Fatal("mismatch hash")
	}
	// Test that our setup config return the proper type configurator.
	wantType = reflect.TypeOf(&multigeth.MultiGethChainConfig{})
	if reflect.TypeOf(conf2) != wantType {
		t.Fatalf("mismatch, want: %v, got: %v", wantType, reflect.TypeOf(conf2))
	}

	// Nitty gritty test that the contemporary stored config, when compactly marshaled,
	// is equal to the expected "healthy" variable value set above.
	// Use compaction to remove whitespace considerations.
	outConf := rawdb.ReadChainConfig(db, genHash)
	outConfMarshal, err := json.MarshalIndent(outConf, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	bCompactB := []byte{}
	bufCompactB := bytes.NewBuffer(bCompactB)

	bCompactA := []byte{}
	bufCompactA := bytes.NewBuffer(bCompactA)

	err = json.Compact(bufCompactB, outConfMarshal)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Compact(bufCompactA, []byte(cc_v197_a))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bCompactB, bCompactA) {
		t.Fatal("different config")
	}
}
