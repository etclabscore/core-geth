package core

import (
	"reflect"
	"testing"

		"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/convert"
	"github.com/ethereum/go-ethereum/params/types/common"
)

func TestSetupGenesisBlock(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	defGenBl := params.DefaultGenesisBlock()

	config, hash, err := SetupGenesisBlock(db, defGenBl)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if wantHash := GenesisToBlock(defGenBl, nil).Hash(); wantHash != hash  {
		t.Errorf("mismatch block hash, want: %x, got: %x", wantHash, hash)
	}
	if diffs := convert.Equal(reflect.TypeOf((*common.ChainConfigurator)(nil)), defGenBl.Config, config); len(diffs) != 0 {
		for _, diff := range diffs {
			t.Error("mismatch", "diff=", diff, "in", defGenBl.Config, "out", config)
		}
	}

	clGenBl := params.DefaultClassicGenesisBlock()

	clConfig, clHash, clErr := SetupGenesisBlock(db, clGenBl)
	if clErr != nil {
		t.Errorf("err: %v", clErr)
	}
	if wantHash := GenesisToBlock(clGenBl, nil).Hash(); wantHash != clHash  {
		t.Errorf("mismatch block hash, want: %x, got: %x", wantHash, clHash)
	}
	if diffs := convert.Equal(reflect.TypeOf((*common.ChainConfigurator)(nil)), clGenBl.Config, clConfig); len(diffs) != 0 {
		for _, diff := range diffs {
			t.Error("mismatch", "diff=", diff, "in", clGenBl.Config, "out", clConfig)
		}
	}
}
