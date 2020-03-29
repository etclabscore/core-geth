package rpc

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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

	uniqueKeyFn := func(sch spec.Schema) string {
		b, _ := json.Marshal(sch)
		sum := sha1.Sum(b)
		out := fmt.Sprintf("%x", sum[:4])

		if sch.Title != "" {
			out = fmt.Sprintf("%s.", sch.Title) + out
		}

		if len(sch.Type) != 0 {
			out = fmt.Sprintf("%s.", strings.Join(sch.Type, "+")) + out
		}

		spl := strings.Split(sch.Description, ":")
		splv := spl[len(spl)-1]
		if splv != "" && !strings.Contains(splv, " ") {
			out = splv + "_" + out
		}

		return out
	}

	registerSchema := func(leaf spec.Schema) error {
		a.registerSchema(leaf, uniqueKeyFn)
		return nil
	}

	//schemaIsEmpty := func(sch *spec.Schema) bool {
	//	return sch == nil || reflect.DeepEqual(*sch, spec.Schema{})
	//}

	//onSchema :=

	mustMarshalString := func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	}

	doc.Components.Schemas = make(map[string]spec.Schema)
	for im := 0; im < len(doc.Methods); im++ {

		met := doc.Methods[im]
		fmt.Println(met.Name)

		parent := &spec.Schema{}

		deferencer := func(sch *spec.Schema) error {

			if sch.Ref.String() == "" {
				return nil
			}
			rr, err := a.schemaFromRef(*parent, sch.Ref)
			if err == nil {
				*sch = rr
			}
			return nil
		}

		referencer := func(sch *spec.Schema) error {

			err = registerSchema(*sch)
			if err != nil {
				fmt.Println("!!! ", err)
				return err
			}

			fmt.Println("   *", mustMarshalString(sch))

			r, err := a.schemaAsReferenceSchema(*sch)
			if err != nil {
				fmt.Println("error getting schema as ref-only schema")
				return err
			}

			doc.Components.Schemas[uniqueKeyFn(*sch)] = *sch
			*sch = r

			fmt.Println("    @", mustMarshalString(sch))
			return nil
		}

		// Params.
		for ip := 0; ip < len(met.Params); ip++ {
			par := met.Params[ip]
			fmt.Println(" < ", par.Name)

			*parent = par.Schema
			a.analysisOnNode(&par.Schema, deferencer)
			a.analysisOnNode(&par.Schema, referencer)
			met.Params[ip] = par
			fmt.Println("   :", mustMarshalString(par))
		}

		// Result (single).
		fmt.Println(" > ", doc.Methods[im].Result.Name)
		*parent = met.Result.Schema
		a.analysisOnNode(&met.Result.Schema, deferencer)
		a.analysisOnNode(&met.Result.Schema, referencer)
		fmt.Println("   :", mustMarshalString(&met.Result))
	}

	//for schv, tit := range a.schemaTitles {
	//	sch := spec.Schema{}
	//	err := json.Unmarshal([]byte(schv), &sch)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	doc.Components.Schemas[tit] = sch
	//}

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

	// Extract schemas. Put
}

//for _, m := range doc.Methods {
//	for _, param := range m.Params {
//		a.analysisOnNode(param.Schema, func(sch spec.Schema) error {
//			if err := registerSchema(sch); err != nil {
//				return err
//			}
//			param.Schema = sch
//			return nil
//		})
//	}
//	a.analysisOnNode(m.Result.Schema, func(sch spec.Schema) error {
//		if err := registerSchema(sch); err != nil {
//			return err
//		}
//		m.Result.Schema = sch
//		return nil
//	})
//}
//
//for _, m := range doc.Methods {
//	for _, param := range m.Params {
//		a.analysisOnNode(param.Schema, func(sch spec.Schema) error {
//			ns, err := a.schemaAsReferenceSchema(sch)
//			if err != nil {
//				return err
//			}
//			param.Schema = ns
//			return nil
//		})
//
//	}
//	a.analysisOnNode(m.Result.Schema, func(sch spec.Schema) error {
//		ns, err := a.schemaAsReferenceSchema(sch)
//		if err != nil {
//			return err
//		}
//		m.Result.Schema = ns
//		return nil
//	})
//}
//
//doc.Components.Schemas = make(map[string]spec.Schema)
//
//// Add schema to component.schemas for all leaves.
//for _, m := range doc.Methods {
//	for _, param := range m.Params {
//		a.analysisOnLeaf(param.Schema, func(sch spec.Schema) error {
//			tit, err := a.getRegisteredSchemaKey(sch)
//			if err != nil {
//				return err
//			}
//			doc.Components.Schemas[tit] = sch
//			return nil
//		})
//	}
//	a.analysisOnLeaf(m.Result.Schema, func(sch spec.Schema) error {
//		tit, err := a.getRegisteredSchemaKey(sch)
//		if err != nil {
//			return err
//		}
//		doc.Components.Schemas[tit] = sch
//		doc.Components.Schemas[uniqueKeyFn(sch)] = sch
//		return nil
//	})
//}

//for schv, tit := range a.schemaTitles {
//	sch := spec.Schema{}
//	err := json.Unmarshal([]byte(schv), &sch)
//	if err != nil {
//		t.Fatal(err)
//	}
//	doc.Components.Schemas[tit] = sch
//}

//for tit, sch := range doc.Components.Schemas {
//	a.analysisOnNode(sch, func(node spec.Schema) error {
//		ns, err := a.schemaAsReferenceSchema(node)
//		if err != nil {
//			return err // NOTE
//		}
//		key := uniqueKeyFn(ns)
//		if v, ok := doc.Components.Schemas[key]; ok {
//			if v.Ref.String() == "" {
//				return nil
//			}
//			v = spec.Schema{
//				SchemaProps: spec.SchemaProps{
//					Ref: spec.Ref{
//						Ref: jsonreference.MustCreateRef("#/components/schemas/" + tit),
//					},
//				},
//			}
//		}
//		return nil
//	})
//}
