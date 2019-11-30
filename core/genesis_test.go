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
	if diff, eq := convert.Equal(reflect.TypeOf((*common.ChainConfigurator)(nil)), defGenBl.Config, config); !eq {
		t.Error("mismatch", "diff=", diff, "in", defGenBl.Config, "out", config)
	}

	clGenBl := params.DefaultClassicGenesisBlock()

	clConfig, clHash, clErr := SetupGenesisBlock(db, clGenBl)
	if clErr != nil {
		t.Errorf("err: %v", clErr)
	}
	if wantHash := GenesisToBlock(clGenBl, nil).Hash(); wantHash != clHash  {
		t.Errorf("mismatch block hash, want: %x, got: %x", wantHash, clHash)
	}
	if diff, eq := convert.Equal(reflect.TypeOf((*common.ChainConfigurator)(nil)), clGenBl.Config, clConfig); !eq {
		t.Error("mismatch", "diff=", diff, "in", clGenBl.Config, "out", clConfig)
	}
}
