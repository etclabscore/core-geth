// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/big"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
	"github.com/imdario/mergo"
)

func (s *RPCService) Describe() (*goopenrpcT.OpenRPCSpec1, error) {

	if s.doc == nil {
		s.doc = NewOpenRPCDescription(s.server)
	}

	for module, list := range s.methods() {
		if module == "rpc" {
			continue
		}

	methodListLoop:
		for _, methodName := range list {
			fullName := strings.Join([]string{module, methodName}, serviceMethodSeparators[0])
			method := s.server.services.services[module].callbacks[methodName]

			// FIXME: Development only.
			// There is a bug with the isPubSub method, it's not picking up #PublicEthAPI.eth_subscribeSyncStatus
			// because the isPubSub conditionals are wrong or the method is wrong.
			if method.isSubscribe || strings.Contains(fullName, subscribeMethodSuffix) {
				continue
			}

			// Dedupe. Not sure how `admin_datadir` got in there twice.
			for _, m := range s.doc.Doc.Methods {
				if m.Name == fullName {
					continue methodListLoop
				}
			}
			if err := s.doc.RegisterMethod(fullName, method); err != nil {
				return nil, err
			}
		}
	}

	//for _, m := range s.doc.Doc.Methods {
	//	for _, p := range m.Params {
	//		//sch := derefSchemaRecurse(&s.doc.Doc.Components, p.Schema)
	//		//s.doc.Doc.Components.Schemas[sch.Ref]
	//		//p.Schema = sch
	//	}
	//	//s := derefSchemaRecurse(&s.doc.Doc.Components, m.Result.Schema)
	//	//m.Result.Schema = s
	//}

	//s.doc.Doc.Components.

	//parse.GetTypes(s.doc.Doc, s.doc.Doc.Objects)
	//
	//for _, k := range s.doc.Doc.Objects.GetKeys() {
	//	kk := s.doc.Doc.Objects.Get(k)
	//	kkk := kk.GetKeys()
	//	log.Println(kk, kkk)
	//}
	////b, _ := json.MarshalIndent(s.doc.Doc.Objects, "", "    ")
	////log.Println("objects", string(b), len(s.doc.Doc.Objects.GetKeys()))

	return s.doc.Doc, nil
}

// ---

type OpenRPCDescription struct {
	Doc *goopenrpcT.OpenRPCSpec1
}

func NewOpenRPCDescription(server *Server) *OpenRPCDescription {
	doc := &goopenrpcT.OpenRPCSpec1{
		OpenRPC: "1.2.4",
		Info: goopenrpcT.Info{
			Title:          "Ethereum JSON-RPC",
			Description:    "This API lets you interact with an EVM-based client via JSON-RPC",
			TermsOfService: "https://github.com/etclabscore/core-geth/blob/master/COPYING",
			Contact: goopenrpcT.Contact{
				Name:  "",
				URL:   "",
				Email: "",
			},
			License: goopenrpcT.License{
				Name: "Apache-2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
			Version: "1.0.10",
		},
		Servers: []goopenrpcT.Server{},
		Methods: []goopenrpcT.Method{},
		Components: goopenrpcT.Components{
			ContentDescriptors:    make(map[string]*goopenrpcT.ContentDescriptor),
			Schemas:               make(map[string]spec.Schema),
			Examples:              make(map[string]goopenrpcT.Example),
			Links:                 make(map[string]goopenrpcT.Link),
			Errors:                make(map[string]goopenrpcT.Error),
			ExamplePairingObjects: make(map[string]goopenrpcT.ExamplePairing),
			Tags:                  make(map[string]goopenrpcT.Tag),
		},
		ExternalDocs: goopenrpcT.ExternalDocs{
			Description: "Source",
			URL:         "https://github.com/etclabscore/core-geth",
		},
		Objects: goopenrpcT.NewObjectMap(),
	}

	return &OpenRPCDescription{Doc: doc}
}

func (d *OpenRPCDescription) RegisterMethod(name string, cb *callback) error {

	cb.makeArgTypes()
	cb.makeRetTypes()

	rtFunc := runtime.FuncForPC(cb.fn.Pointer())
	cbFile, _ := rtFunc.FileLine(rtFunc.Entry())

	tokenset := token.NewFileSet()
	astFile, err := parser.ParseFile(tokenset, cbFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astFuncDel := getAstFunc(cb, astFile, rtFunc)

	if astFuncDel == nil {
		return fmt.Errorf("nil ast func: method name: %s", name)
	}

	method, err := makeMethod(name, cb, rtFunc, astFuncDel)
	if err != nil {
		return fmt.Errorf("make method error method=%s cb=%s error=%v", name, spew.Sdump(cb), err)
	}

	d.Doc.Methods = append(d.Doc.Methods, method)
	sort.Slice(d.Doc.Methods, func(i, j int) bool {
		return d.Doc.Methods[i].Name < d.Doc.Methods[j].Name
	})

	return nil
}

type argIdent struct {
	ident *ast.Ident
	name  string
}

func (a argIdent) Name() string {
	if a.ident != nil {
		return a.ident.Name
	}
	return a.name
}

type AnalysisT struct {
	OpenMetaDescription string
	//ErrorHandling AnalysisTErrorHandling
	//Errors []error
	schemaTitles map[string]string
	//TitleKeyer   func(schema spec.Schema) string
	lastSch *spec.Schema
}

func (a *AnalysisT) schemaAsReferenceSchema(sch spec.Schema) (refSchema spec.Schema, err error) {
	b, _ := json.Marshal(sch)
	titleKey, ok := a.schemaTitles[string(b)]
	if !ok {
		bb, _ := json.Marshal(sch)
		return refSchema, fmt.Errorf("schema not available as reference: %s @ %v", string(b), string(bb))
	}
	refSchema.Ref = spec.Ref{
		Ref: jsonreference.MustCreateRef("#/components/schemas/" + titleKey),
	}
	return
}

//func (a *AnalysisT) getRegisteredSchema()
func (a *AnalysisT) registerSchema(sch spec.Schema, titleKeyer func(schema spec.Schema) string) {
	b, _ := json.Marshal(sch)
	a.schemaTitles[string(b)] = titleKeyer(sch)
}

func (a *AnalysisT) schemaFromRef(psch spec.Schema, ref spec.Ref) (schema spec.Schema, err error) {
	v, _, err := ref.GetPointer().Get(psch)
	if err != nil {
		return
	}
	return v.(spec.Schema), nil
}

//
//func (a *AnalysisT) dereferenceSchema(sch *spec.Schema) {
//
//	// Slices.
//	for i := 0; i < len(sch.OneOf); i++ {
//		it := sch.OneOf[i]
//		a.analysisOnNode(sch, &it, onNode)
//	}
//	for i := 0; i < len(sch.AnyOf); i++ {
//		it := sch.AnyOf[i]
//		a.analysisOnNode(sch, &it, onNode)
//	}
//	for i := 0; i < len(sch.AllOf); i++ {
//		it := sch.AllOf[i]
//		a.analysisOnNode(sch, &it, onNode)
//	}
//	// Maps.
//	for k, defSch := range sch.Definitions {
//		defSch.Title = k
//		a.analysisOnNode(sch, &defSch, onNode)
//	}
//	for k := range sch.Properties {
//		v := sch.Properties[k]
//		a.analysisOnNode(sch, &v, onNode)
//	}
//	for k := range sch.PatternProperties {
//		v := sch.PatternProperties[k]
//		a.analysisOnNode(sch, &v, onNode)
//	}
//	if sch.Items == nil {
//		//onNode(prnt, sch)
//		return nil
//	}
//	if sch.Items.Len() > 1 {
//		for i := range sch.Items.Schemas {
//			a.analysisOnNode(sch, &sch.Items.Schemas[i], onNode) // PTAL: Is this right?
//		}
//	} else {
//		// Is schema
//		a.analysisOnNode(sch, sch.Items.Schema, onNode)
//	}
//	//onNode(prnt, sch)
//	return nil
//}

// analysisOnNode runs a callback function on each leaf of a the JSON schema tree.
// It will return the first error it encounters.
func (a *AnalysisT) analysisOnNode(sch *spec.Schema, onNode func(node *spec.Schema) error) error {

	// Slices.
	for i := 0; i < len(sch.OneOf); i++ {
		it := sch.OneOf[i]
		a.analysisOnNode(&it, onNode)
		sch.OneOf[i] = it
	}
	for i := 0; i < len(sch.AnyOf); i++ {
		it := sch.AnyOf[i]
		a.analysisOnNode(&it, onNode)
		sch.AnyOf[i] = it
	}
	for i := 0; i < len(sch.AllOf); i++ {
		it := sch.AllOf[i]
		a.analysisOnNode(&it, onNode)
		sch.AllOf[i] = it
	}
	// Maps.
	for k := range sch.Definitions {
		v := sch.Definitions[k]
		v.Title = k
		a.analysisOnNode(&v, onNode)
		sch.Definitions[k] = v
	}
	for k := range sch.Properties {
		v := sch.Properties[k]
		//v.Title = k
		a.analysisOnNode(&v, onNode)
		sch.Properties[k] = v
	}
	for k := range sch.PatternProperties {
		v := sch.PatternProperties[k]
		//v.Title = k
		a.analysisOnNode(&v, onNode)
		sch.PatternProperties[k] = v
	}
	if sch.Items == nil {
		//onNode(sch)
		return onNode(sch)
		//return nil
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			a.analysisOnNode(&sch.Items.Schemas[i], onNode) // PTAL: Is this right?
		}
	} else {
		a.analysisOnNode(sch.Items.Schema, onNode)
	}
	return onNode(sch)
}

// analysisOnLeaf runs a callback function on each leaf of a the JSON schema tree.
// It will return the first error it encounters.
func (a *AnalysisT) analysisOnLeaf(sch spec.Schema, onLeaf func(leaf spec.Schema) error) error {
	for i := range sch.Definitions {
		a.analysisOnLeaf(sch.Definitions[i], onLeaf)
	}
	for i := range sch.OneOf {
		a.analysisOnLeaf(sch.OneOf[i], onLeaf)
	}
	for i := range sch.AnyOf {
		a.analysisOnLeaf(sch.AnyOf[i], onLeaf)
	}
	for i := range sch.AllOf {
		a.analysisOnLeaf(sch.AllOf[i], onLeaf)
	}
	for k, _ := range sch.Properties {
		a.analysisOnLeaf(sch.Properties[k], onLeaf)
	}
	for k, _ := range sch.PatternProperties {
		a.analysisOnLeaf(sch.PatternProperties[k], onLeaf)
	}
	if sch.Items == nil {
		return onLeaf(sch)
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			a.analysisOnLeaf(sch.Items.Schemas[i], onLeaf) // PTAL: Is this right?
		}
	} else {
		// Is schema
		a.analysisOnLeaf(*sch.Items.Schema, onLeaf)
	}
	return onLeaf(sch)
}

func schemaHasRef(sch spec.Schema) bool {
	return sch.Ref.String() != ""
}

func getComponentsSchemaFromRef(cmpnts *goopenrpcT.Components, ref spec.Ref) (sch spec.Schema) {
	if cmpnts == nil || ref.String() == "" {
		return
	}
	r := filepath.Base(ref.String())
	sch = cmpnts.Schemas[r] // Trust parser
	return
}

func derefSchemaRecurse(cts *goopenrpcT.Components, sch spec.Schema) spec.Schema {
	if schemaHasRef(sch) {
		sch = getComponentsSchemaFromRef(cts, sch.Ref)
		sch = derefSchemaRecurse(cts, sch)
	}
	for i := range sch.Definitions {
		got := derefSchemaRecurse(cts, sch.Definitions[i])
		if err := mergo.Merge(&got, sch.Definitions[i]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.Definitions[i] = got
	}
	for i := range sch.OneOf {
		got := derefSchemaRecurse(cts, sch.OneOf[i])
		if err := mergo.Merge(&got, sch.OneOf[i]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.OneOf[i] = got
	}
	for i := range sch.AnyOf {
		got := derefSchemaRecurse(cts, sch.AnyOf[i])
		if err := mergo.Merge(&got, sch.AnyOf[i]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.AnyOf[i] = got
	}
	for i := range sch.AllOf {
		got := derefSchemaRecurse(cts, sch.AllOf[i])
		if err := mergo.Merge(&got, sch.AllOf[i]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.AllOf[i] = got
	}
	for k, _ := range sch.Properties {
		got := derefSchemaRecurse(cts, sch.Properties[k])
		if err := mergo.Merge(&got, sch.Properties[k]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.Properties[k] = got
	}
	for k, _ := range sch.PatternProperties {
		got := derefSchemaRecurse(cts, sch.PatternProperties[k])
		if err := mergo.Merge(&got, sch.PatternProperties[k]); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.PatternProperties[k] = got
	}
	if sch.Items == nil {
		return sch
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			got := derefSchemaRecurse(cts, sch.Items.Schemas[i])
			if err := mergo.Merge(&got, sch.Items.Schemas[i]); err != nil {
				panic(err.Error())
			}
			got.Schema = ""
			sch.Ref = spec.Ref{}
			sch.Items.Schemas[i] = got
		}
	} else {
		// Is schema
		got := derefSchemaRecurse(cts, *sch.Items.Schema)
		if err := mergo.Merge(&got, sch.Items.Schema); err != nil {
			panic(err.Error())
		}
		got.Schema = ""
		sch.Ref = spec.Ref{}
		sch.Items.Schema = &got
	}

	return sch
}

func makeMethod(name string, cb *callback, rt *runtime.Func, fn *ast.FuncDecl) (goopenrpcT.Method, error) {
	file, line := rt.FileLine(rt.Entry())

	//packageName := strings.Split(rt.Name(), ".")[0]

	m := goopenrpcT.Method{
		Name:        name,
		Tags:        []goopenrpcT.Tag{},
		Summary:     fn.Doc.Text(),
		Description: "", // fmt.Sprintf(`%s@%s:%d'`, rt.Name(), file, line),
		ExternalDocs: goopenrpcT.ExternalDocs{
			Description: fmt.Sprintf(`%s`, rt.Name()),
			URL:         fmt.Sprintf("file://%s:%d", file, line),
		},
		Params:         []*goopenrpcT.ContentDescriptor{},
		Result:         &goopenrpcT.ContentDescriptor{},
		Deprecated:     false,
		Servers:        []goopenrpcT.Server{},
		Errors:         []goopenrpcT.Error{},
		Links:          []goopenrpcT.Link{},
		ParamStructure: "by-position",
		Examples:       []goopenrpcT.ExamplePairing{},
	}

	defer func() {
		if m.Result.Name == "" {
			m.Result.Name = "null"
			m.Result.Schema.Type = []string{"null"}
			m.Result.Schema.Description = "Null"
		}
	}()

	//derefContentDescriptor := func(descriptor *goopenrpcT.ContentDescriptor) {
	//	for _, d := range descriptor.Schema.Definitions {
	//
	//	}
	//}

	if fn.Type.Params != nil {
		j := 0
		for _, field := range fn.Type.Params.List {
			if field == nil {
				continue
			}
			if cb.hasCtx && strings.Contains(fmt.Sprintf("%s", field.Type), "context") {
				continue
			}
			if len(field.Names) > 0 {
				for _, ident := range field.Names {
					if ident == nil {
						continue
					}
					if j > len(cb.argTypes)-1 {
						log.Println(name, cb.argTypes, field.Names, j)
						continue
					}
					cd, err := makeContentDescriptor(cb.argTypes[j], field, argIdent{ident, fmt.Sprintf("%sParameter%d", name, j)})
					if err != nil {
						return m, err
					}
					j++
					m.Params = append(m.Params, &cd)
				}
			} else {
				cd, err := makeContentDescriptor(cb.argTypes[j], field, argIdent{nil, fmt.Sprintf("%sParameter%d", name, j)})
				if err != nil {
					return m, err
				}
				j++
				m.Params = append(m.Params, &cd)
			}

		}
	}
	if fn.Type.Results != nil {
		j := 0
		for _, field := range fn.Type.Results.List {
			if field == nil {
				continue
			}
			if strings.Contains(fmt.Sprintf("%s", field.Type), "error") {
				continue
			}
			if len(field.Names) > 0 {
				// This really should never ever happen I don't think.
				// JSON-RPC returns _an_ result. So there can't be > 1 return value.
				// But just in case.
				for _, ident := range field.Names {
					cd, err := makeContentDescriptor(cb.retTypes[j], field, argIdent{ident, fmt.Sprintf("%sResult%d", name, j)})
					if err != nil {
						return m, err
					}
					j++
					m.Result = &cd
				}
			} else {
				cd, err := makeContentDescriptor(cb.retTypes[j], field, argIdent{nil, fmt.Sprintf("%sResult", name)})
				if err != nil {
					return m, err
				}
				j++
				m.Result = &cd
			}
		}
	}

	return m, nil
}

func makeContentDescriptor(ty reflect.Type, field *ast.Field, ident argIdent) (goopenrpcT.ContentDescriptor, error) {
	cd := goopenrpcT.ContentDescriptor{}
	if !jsonschemaPkgSupport(ty) {
		return cd, fmt.Errorf("unsupported iface: %v %v %v", spew.Sdump(ty), spew.Sdump(field), spew.Sdump(ident))
	}

	schemaType := fmt.Sprintf("%s:%s", ty.PkgPath(), ty.Name())
	switch tt := field.Type.(type) {
	case *ast.SelectorExpr:
		schemaType = fmt.Sprintf("%v.%v", tt.X, tt.Sel)
		schemaType = fmt.Sprintf("%s:%s", ty.PkgPath(), schemaType)
	case *ast.StarExpr:
		schemaType = fmt.Sprintf("%v", tt.X)
		schemaType = fmt.Sprintf("*%s:%s", ty.PkgPath(), schemaType)
		if reflect.ValueOf(ty).Type().Kind() == reflect.Ptr {
			schemaType = fmt.Sprintf("%v", ty.Elem().Name())
			schemaType = fmt.Sprintf("*%s:%s", ty.Elem().PkgPath(), schemaType)
		}
		//ty = ty.Elem() // FIXME: wart warn
	}
	//schemaType = fmt.Sprintf("%s:%s", ty.PkgPath(), schemaType)

	//cd.Name = schemaType
	cd.Name = ident.Name()

	cd.Summary = field.Doc.Text()
	cd.Description = field.Comment.Text()

	rflctr := jsonschema.Reflector{
		AllowAdditionalProperties:  false, // false,
		RequiredFromJSONSchemaTags: false,
		ExpandedStruct:             true, // false, // false,
		//IgnoredTypes:               []interface{}{chaninterface},
		TypeMapper: OpenRPCJSONSchemaTypeMapper,
	}

	jsch := rflctr.ReflectFromType(ty)

	// Poor man's type cast.
	m, err := json.Marshal(jsch)
	if err != nil {
		log.Fatal(err)
	}
	sch := spec.Schema{}
	err = json.Unmarshal(m, &sch)
	if err != nil {
		log.Fatal(err)
	}

	if schemaType != ":" && (cd.Schema.Description == "" || cd.Schema.Description == ":") {
		sch.Description = schemaType
	}

	cd.Schema = sch

	return cd, nil
}

func jsonschemaPkgSupport(r reflect.Type) bool {
	switch r.Kind() {
	case reflect.Struct,
		reflect.Map,
		reflect.Slice, reflect.Array,
		reflect.Interface,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.String,
		reflect.Ptr:
		return true
	default:
		return false
	}
}

type schemaDictEntry struct {
	t interface{}
	j string
}

func OpenRPCJSONSchemaTypeMapper(r reflect.Type) *jsonschema.Type {
	unmarshalJSONToJSONSchemaType := func(input string) *jsonschema.Type {
		var js jsonschema.Type
		err := json.Unmarshal([]byte(input), &js)
		if err != nil {
			return nil
		}
		return &js
	}

	//unmarshalJSONToJSONSchemaTypeGlom := func(input string, ref *jsonschema.Type) *jsonschema.Type {
	//	err := json.Unmarshal([]byte(input), ref)
	//	if err != nil {
	//		return nil
	//	}
	//	return ref
	//}

	//handleNil := func(r reflect.Type, ret *jsonschema.Type) *jsonschema.Type {
	//	if ret == nil {
	//		return nil
	//	}
	//	if reflect.ValueOf(r).Kind() == reflect.Ptr {
	//		if ret.OneOf == nil || len(ret.OneOf) == 0 {
	//			ttt := &jsonschema.Type{}
	//			ttt.Title = ret.Title // Only keep the title.
	//			ttt.OneOf = []*jsonschema.Type{}
	//			ttt.OneOf = append(ttt.OneOf, ret)
	//			ttt.OneOf = append(ttt.OneOf, unmarshalJSONToJSONSchemaType(`
	//	{
	//      "title": "null",
	//      "type": "null",
	//      "description": "Null"
	//    }`))
	//			return ttt
	//		}
	//	}
	//
	//	return ret
	//}

	integerD := `{
          "title": "integer",
          "type": "string",
          "pattern": "^0x[a-fA-F0-9]+$",
          "description": "Hex representation of the integer"
        }`
	commonHashD := `{
          "title": "keccak",
          "type": "string",
          "description": "Hex representation of a Keccak 256 hash",
          "pattern": "^0x[a-fA-F\\d]{64}$"
        }`
	blockNumberTagD := `{
          "title": "blockNumberTag",
          "type": "string",
          "description": "The optional block height description",
          "enum": [
            "earliest",
            "latest",
            "pending"
          ]
        }`

	// Second, handle other types.
	// Use a slice instead of a map because it preserves order, as a logic safeguard/fallback.
	dict := []schemaDictEntry{

		{new(big.Int), integerD},
		{big.Int{}, integerD},
		{new(hexutil.Big), integerD},
		{hexutil.Big{}, integerD},

		{new(common.Address), `{
          "title": "keccak",
          "type": "string",
          "description": "Hex representation of a Keccak 256 hash POINTER",
          "pattern": "^0x[a-fA-F\\d]{64}$"
        }`},

		{common.Address{}, `{
          "title": "address",
          "type": "string",
          "pattern": "^0x[a-fA-F\\d]{40}$"
        }`},

		{new(common.Hash), `{
          "title": "keccak",
          "type": "string",
          "description": "Hex representation of a Keccak 256 hash POINTER",
          "pattern": "^0x[a-fA-F\\d]{64}$"
        }`},

		{common.Hash{}, commonHashD},

		{
			hexutil.Bytes{}, `{
          "title": "dataWord",
          "type": "string",
          "description": "Hex representation of a 256 bit unit of data",
          "pattern": "^0x([a-fA-F\\d]{64})?$"
        }`},
		{
			new(hexutil.Bytes), `{
          "title": "dataWord",
          "type": "string",
          "description": "Hex representation of a 256 bit unit of data",
          "pattern": "^0x([a-fA-F\\d]{64})?$"
        }`},

		{[]byte{}, `{
          "title": "bytes",
          "type": "string",
          "description": "Hex representation of a variable length byte array",
          "pattern": "^0x([a-fA-F0-9]?)+$"
        }`},

		{BlockNumberOrHash{}, fmt.Sprintf(`{
		  "title": "blockNumberOrHash",
		  "description": "Hex representation of a block number or hash",
		  "oneOf": [%s, %s]
		}`, commonHashD, integerD)},

		{BlockNumber(0), fmt.Sprintf(`{
		  "title": "blockNumberOrTag",
		  "description": "Block tag or hex representation of a block number",
		  "oneOf": [%s, %s]
		}`, commonHashD, blockNumberTagD)},
	}

	for _, d := range dict {
		d := d
		if reflect.TypeOf(d.t) == r {
			tt := unmarshalJSONToJSONSchemaType(d.j)

			//	// If the value is a pointer, then it can be nil.
			//	// Which means it should be oneOf <the value> or <Null>.
			//	//if strings.Contains(reflect.ValueOf(d.t).Type().Name(), "*") { // .Kind() == reflect.Ptr
			//	//if d.t == nil {
			//	if reflect.ValueOf(d.t).Type().Kind() == reflect.Ptr {
			//		if tt.OneOf == nil || len(tt.OneOf) == 0 {
			//			ttt := &jsonschema.Type{}
			//			ttt.Title = tt.Title // Only keep the title.
			//			ttt.OneOf = []*jsonschema.Type{}
			//			ttt.OneOf = append(ttt.OneOf, tt)
			//			ttt.OneOf = append(ttt.OneOf, unmarshalJSONToJSONSchemaType(`
			//{
			//  "title": "null",
			//  "type": "null",
			//  "description": "Null"
			//}`))
			//			return ttt
			//		}
			//	}

			return tt // handleNil(r, tt)
		}
	}

	// First, handle primitives.
	switch r.Kind() {
	case reflect.Struct:

	case reflect.Map,
		reflect.Interface:
	case reflect.Slice, reflect.Array:

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret := unmarshalJSONToJSONSchemaType(integerD)

		//if reflect.ValueOf(r).Type().Kind() == reflect.Ptr {
		//	if ret.OneOf == nil || len(ret.OneOf) == 0 {
		//		ttt := &jsonschema.Type{}
		//		ttt.Title = ret.Title // Only keep the title.
		//		ttt.OneOf = []*jsonschema.Type{}
		//		ttt.OneOf = append(ttt.OneOf, ret)
		//		ttt.OneOf = append(ttt.OneOf, unmarshalJSONToJSONSchemaType(`
		//{
		//  "title": "null",
		//  "type": "null",
		//  "description": "Null"
		//}`))
		//		return ttt
		//	}
		//}
		return ret

	case reflect.Float32, reflect.Float64:

	case reflect.Bool:

	case reflect.String:

	case reflect.Ptr:

	default:
		panic("prevent me somewhere else please")
	}

	return nil
}

func getAstFunc(cb *callback, astFile *ast.File, rf *runtime.Func) *ast.FuncDecl {

	rfName := runtimeFuncName(rf)
	for _, decl := range astFile.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name == nil || fn.Name.Name != rfName {
			continue
		}
		//log.Println("getAstFunc", spew.Sdump(cb), spew.Sdump(fn))
		fnRecName := ""
		for _, l := range fn.Recv.List {
			if fnRecName != "" {
				break
			}
			i, ok := l.Type.(*ast.Ident)
			if ok {
				fnRecName = i.Name
				continue
			}
			s, ok := l.Type.(*ast.StarExpr)
			if ok {
				fnRecName = fmt.Sprintf("%v", s.X)
			}
		}
		// Ensure that this is the one true function.
		// Have to match receiver AND method names.
		/*
		 => recvr= <*ethapi.PublicBlockChainAPI Value> fn= PublicBlockChainAPI
		 => recvr= <*ethash.API Value> fn= API
		 => recvr= <*ethapi.PublicTxPoolAPI Value> fn= PublicTxPoolAPI
		 => recvr= <*ethapi.PublicTxPoolAPI Value> fn= PublicTxPoolAPI
		 => recvr= <*ethapi.PublicTxPoolAPI Value> fn= PublicTxPoolAPI
		 => recvr= <*ethapi.PublicNetAPI Value> fn= PublicNetAPI
		 => recvr= <*ethapi.PublicNetAPI Value> fn= PublicNetAPI
		 => recvr= <*ethapi.PublicNetAPI Value> fn= PublicNetAPI
		 => recvr= <*node.PrivateAdminAPI Value> fn= PrivateAdminAPI
		 => recvr= <*node.PublicAdminAPI Value> fn= PublicAdminAPI
		 => recvr= <*node.PublicAdminAPI Value> fn= PublicAdminAPI
		 => recvr= <*eth.PrivateAdminAPI Value> fn= PrivateAdminAPI
		*/
		reRec := regexp.MustCompile(fnRecName + `\s`)
		if !reRec.MatchString(cb.rcvr.String()) {
			continue
		}
		return fn
	}
	return nil
}

func getAstType(astFile *ast.File, t reflect.Type) *ast.TypeSpec {
	log.Println("getAstType", t.Name(), t.String())
	for _, decl := range astFile.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if d.Tok != token.TYPE {
			continue
		}
		for _, s := range d.Specs {
			sp, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if sp.Name != nil && sp.Name.Name == t.Name() {
				return sp
			} else if sp.Name != nil {
				log.Println("nomatch", sp.Name.Name)
			}
		}

	}
	return nil
}

func runtimeFuncName(rf *runtime.Func) string {
	spl := strings.Split(rf.Name(), ".")
	return spl[len(spl)-1]
}

func (d *OpenRPCDescription) findMethodByName(name string) (ok bool, method goopenrpcT.Method) {
	for _, m := range d.Doc.Methods {
		if m.Name == name {
			return true, m
		}
	}
	return false, goopenrpcT.Method{}
}

func runtimeFuncPackageName(rf *runtime.Func) string {
	re := regexp.MustCompile(`(?im)^(?P<pkgdir>.*/)(?P<pkgbase>[a-zA-Z0-9\-_]*)`)
	match := re.FindStringSubmatch(rf.Name())
	pmap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			pmap[name] = match[i]
		}
	}
	return pmap["pkgdir"] + pmap["pkgbase"]
}
