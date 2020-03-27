package rpc

import (
	"encoding/json"
	"testing"
)

func mustMarshalJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}

func TestOpenRPCDescription(t *testing.T) {
	server := newTestServer()

	rpcService := &RPCService{server: server, doc: NewOpenRPCDescription(server)}
	err := server.RegisterName(MetadataApi, rpcService)
	if err != nil {
		t.Fatal(err)
	}

	desribed, err := rpcService.Describe()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("doc %s", mustMarshalJSON(desribed))
}
