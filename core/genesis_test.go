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
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/confp"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
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
