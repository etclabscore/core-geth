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
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
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
	config, hash, err := SetupGenesisBlock(db, genA)
	if err != nil {
		t.Fatal(err)
	}

	headHash := common.HexToHash("0xe618c1b2d738dfa09052e199e5870274f09eb83c684a8a2c194b82dedc00a977")
	rawdb.WriteHeadHeaderHash(db, headHash)
	rawdb.WriteHeaderNumber(db, headHash, 9700559)

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
}
