package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	goopenrpcT "github.com/gregdhill/go-openrpc/types"
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

func TestOpenRPC_Flatten(t *testing.T) {
	testSpecFile := filepath.Join("..", ".develop", "spec.json")
	b, err := ioutil.ReadFile(testSpecFile)
	if err != nil {
		t.Fatal(err)
	}
	doc := &goopenrpcT.OpenRPCSpec1{}
	err = json.Unmarshal(b, doc)
	if err != nil {
		t.Fatal(err)
	}

	bb, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bb))

	// Extract schemas. Put
}
