package openrpc_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/internal/openrpc"
	"github.com/ethereum/go-ethereum/rpc"
)

func TestDefaultSchema(t *testing.T) {
	if err := rpc.SetDefaultOpenRPCSchemaRaw(openrpc.OpenRPCSchema); err != nil {
		t.Fatal(err)
	}
}
