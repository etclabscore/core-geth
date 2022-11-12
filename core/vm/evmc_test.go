package vm

import (
	"testing"

	"github.com/ethereum/evmc/v10/bindings/go/evmc"
	"github.com/ethereum/go-ethereum/params"
)

func TestGetRevision(t *testing.T) {
	conf := params.AllEthashProtocolChanges

	env := NewEVM(BlockContext{
		BlockNumber: big0,
	}, TxContext{}, nil, conf, Config{})

	rev := getRevision(env)

	if rev != evmc.London {
		t.Errorf("expected revision London, got %v", rev)
	}
}
