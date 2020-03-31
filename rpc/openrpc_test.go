package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/jst"
	"github.com/go-openapi/spec"
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

// https://stackoverflow.com/questions/46904588/efficient-way-to-to-generate-a-random-hex-string-of-a-fixed-length-in-golang
var src = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandStringBytesMaskImprSrc returns a random hexadecimal string of length n.
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := src.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}

func TestOpenRPC_Analysis(t *testing.T) {
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

	if err := Clean(doc); err != nil {
		t.Fatal(err)
	}

	docbb, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	//schemasbb, err := json.MarshalIndent(doc.Components.Schemas, "", "    ")
	//if err != nil {
	//	t.Fatal(err)
	//}

	fmt.Println(string(docbb))

	err = ioutil.WriteFile(filepath.Join("..", ".develop", "spec2.json"), docbb, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func testOnNode(node *spec.Schema) error {
	b, err := json.MarshalIndent(node, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func TestAnalysisOnNode(t *testing.T) {
	schemaJSON := `
{
	"type": "object",
	"properties": {
		"foo": {}
	}
}`

	schema := spec.Schema{}
	err := json.Unmarshal([]byte(schemaJSON), &schema)
	if err != nil {
		t.Fatal(err)
	}

	aa := jst.NewAnalysisT()
	err = aa.Traverse(&schema, testOnNode)
	if err != nil {
		t.Error(err)
	}

	schema.Properties["foo"] = schema
	err = aa.Traverse(&schema, testOnNode)
	if err != nil {
		t.Error(err)
	}

}
