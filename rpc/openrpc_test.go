package rpc

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

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

	a := &AnalysisT{
		OpenMetaDescription: "Analysis test",
		schemaTitles:        make(map[string]string),
	}

	inspect := func(leaf spec.Schema) error {
		a.registerSchema(leaf, func(sch spec.Schema) string {

			//if sch.Title != "" {
			//	if len(sch.Type) > 0 {
			//		return sch.Type[0] + "@" + sch.Title
			//	}
			//	return strings.Join(append(sch.Type, sch.Title), ",")
			//}

			b, _ := json.Marshal(sch)
			sum := sha1.Sum(b)
			return fmt.Sprintf("%s%v%v%x", sch.Title, sch.Description, sch.AdditionalProperties, sum)

			//return RandStringBytesMaskImprSrc(8)

			//if sch.Description != "" {
			//	spl := strings.Split(sch.Description, ":")
			//	return spl[len(spl)-1]
			//}
			//if len(sch.Type) == 1 {
			//	switch sch.Type[0] {
			//	case "array":
			//		out := "array"
			//		for _, s := range sch.Items.Schemas {
			//			out += "+" + s.Type[0]
			//		}
			//		return out
			//	case "object":
			//		//return "object:" + sch.Description + sch.Pattern
			//		b, _ := json.Marshal(sch)
			//		sum := sha1.Sum(b)
			//		return fmt.Sprintf("object%x", sum)
			//	default:
			//		return sch.Type[0]
			//	}
			//}
			//
			//return strings.Join(sch.Type, "+")
		})
		//l, err := json.Marshal(leaf)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//fmt.Println(string(l))
		return nil
	}
	for _, m := range doc.Methods {
		for _, param := range m.Params {
			a.analysisOnNode(param.Schema, func(sch spec.Schema) error {
				inspect(sch)
				return nil
			})
		}
		a.analysisOnNode(m.Result.Schema, func(sch spec.Schema) error {
			inspect(sch)
			return nil
		})
	}

	for _, m := range doc.Methods {
		fmt.Println(m.Name)
		for _, param := range m.Params {
			param := param
			a.analysisOnNode(param.Schema, func(sch spec.Schema) error {
				ns := a.schemaReferenced(param.Schema)
				param.Schema = ns
				b, _ := json.Marshal(param)
				fmt.Println(" < " + param.Name, string(b))
				return nil
			})
		}
		a.analysisOnNode(m.Result.Schema, func(sch spec.Schema) error {
			ns := a.schemaReferenced(m.Result.Schema)
			m.Result.Schema = ns
			b, _ := json.Marshal(m.Result)
			fmt.Println(" > " + m.Result.Name, string(b))
			return nil
		})
	}

	for sch, tit := range a.schemaTitles {
		fmt.Println(tit, sch)
	}

	//bb, err := json.MarshalIndent(doc, "", "    ")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(string(bb))

	// Extract schemas. Put
}
