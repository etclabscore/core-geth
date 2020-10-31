package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"math/big"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/alecthomas/jsonschema"
	go_openrpc_reflect "github.com/etclabscore/go-openrpc-reflect"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	meta_schema "github.com/open-rpc/meta-schema"
)

// RPCDiscoveryService defines a receiver type used for RPC discovery by reflection.
type RPCDiscoveryService struct {
	d *go_openrpc_reflect.Document
}

// Discover exposes a Discover method to the RPC receiver registration.
func (r *RPCDiscoveryService) Discover() (*meta_schema.OpenrpcDocument, error) {
	return r.d.Discover()
}

// newOpenRPCDocument returns a Document configured with application-specific logic.
func newOpenRPCDocument() *go_openrpc_reflect.Document {
	d := &go_openrpc_reflect.Document{}

	// Register "Meta" document fields.
	// These include getters for
	// - Servers object
	// - Info object
	// - ExternalDocs object
	//
	// These objects represent server-specific data that cannot be
	// reflected.
	d.WithMeta(&go_openrpc_reflect.MetaT{
		GetServersFn: func() func(listeners []net.Listener) (*meta_schema.Servers, error) {
			return func(listeners []net.Listener) (*meta_schema.Servers, error) {
				servers := []meta_schema.ServerObject{}
				for _, listener := range listeners {
					addr := "http://" + listener.Addr().String()
					network := listener.Addr().Network()
					servers = append(servers, meta_schema.ServerObject{
						Url:  (*meta_schema.ServerObjectUrl)(&addr),
						Name: (*meta_schema.ServerObjectName)(&network),
					})
				}
				return (*meta_schema.Servers)(&servers), nil
			}
		},
		GetInfoFn: func() (info *meta_schema.InfoObject) {
			info = &meta_schema.InfoObject{}
			title := "Core-Geth RPC API"
			info.Title = (*meta_schema.InfoObjectProperties)(&title)

			version := params.VersionWithMeta
			info.Version = (*meta_schema.InfoObjectVersion)(&version)
			return info
		},
		GetExternalDocsFn: func() (exdocs *meta_schema.ExternalDocumentationObject) {
			return nil // FIXME
		},
	})

	// Use a provided Ethereum default configuration as a base.
	appReflector := &go_openrpc_reflect.EthereumReflectorT{}

	// Install overrides for the json schema->type map fn used by the jsonschema reflect package.
	appReflector.FnSchemaTypeMap = func() func(ty reflect.Type) *jsonschema.Type {
		return OpenRPCJSONSchemaTypeMapper
	}

	// Install an override for method eligibility to exclude subscription methods.
	// The majority of this logic is taken from the go_openrpc_reflect package configuration default,
	// with the clause commented 'Custom' noting the custom logic.
	var errType = reflect.TypeOf((*error)(nil)).Elem()
	var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	appReflector.FnIsMethodEligible = func(method reflect.Method) bool {

		// Method must be exported.
		if method.PkgPath != "" {
			return false
		}

		// Custom: skip methods with 1 argument of context.Context.
		// We'll consider these methods subscription methods,
		// and they are not well supported by OpenRPC.
		if method.Type.NumIn() == 2 {
			if method.Type.In(1) == contextType {
				return false
			}
		}

		// Verify return types. The function must return at most one error
		// and/or one other non-error value.
		outs := make([]reflect.Type, method.Func.Type().NumOut())
		for i := 0; i < method.Func.Type().NumOut(); i++ {
			outs[i] = method.Func.Type().Out(i)
		}
		isErrorType := func(ty reflect.Type) bool {
			return ty == errType
		}

		// If an error is returned, it must be the last returned value.
		switch {
		case len(outs) > 2:
			return false
		case len(outs) == 1 && isErrorType(outs[0]):
			return true
		case len(outs) == 2:
			if isErrorType(outs[0]) || !isErrorType(outs[1]) {
				return false
			}
		}
		return true
	}

	appReflector.FnGetContentDescriptorName = func(r reflect.Value, m reflect.Method, field *ast.Field) (string, error) {
		fs := expandedFieldNamesFromList([]*ast.Field{field})
		name := fs[0].Names[0].Name
		// removeChars are characters that look like code.
		// Shane doesn't like these because they might be weird for generated clients to use
		// as variable/field names (eg for params-by-name stuff).
		removeChars := ".*[]{}-"
		for _, c := range strings.Split(removeChars, "") {
			name = strings.ReplaceAll(name, c, "")
		}
		if regexp.MustCompile(`(?m)^\d`).MatchString(name) {
			name = "num" + name
		}
		return name, nil
	}

	// Finally, register the configured reflector to the document.
	d.WithReflector(appReflector)
	return d
}

// registerOpenRPCAPIs provides a convenience logic that is reused
// congruent to the rpc package receiver registrations.
func registerOpenRPCAPIs(doc *go_openrpc_reflect.Document, apis []rpc.API) {
	for _, api := range apis {
		doc.RegisterReceiverName(api.Namespace, api.Service)
	}
}

const bytesD = `{
"title": "bytes",
"type": "string",
"description": "Hex representation of a variable length byte array",
"pattern": "^0x([a-fA-F0-9]?)+$"
}`
const integerD = `{
          "title": "integer",
          "type": "string",
          "pattern": "^0x[a-fA-F0-9]+$",
          "description": "Hex representation of the integer"
        }`
const commonAddressD = `{
          "title": "keccak",
          "type": "string",
          "description": "Hex representation of a Keccak 256 hash POINTER",
          "pattern": "^0x[a-fA-F\\d]{64}$"
        }`
const commonHashD = `{
          "title": "keccak",
          "type": "string",
          "description": "Hex representation of a Keccak 256 hash",
          "pattern": "^0x[a-fA-F\\d]{64}$"
        }`
const hexutilBytesD = `{
          "title": "dataWord",
          "type": "string",
          "description": "Hex representation of some bytes",
          "pattern": "^0x([a-fA-F\\d])+$"
        }`
const hexutilUintD = `{
		"title": "uint",
			"type": "string",
			"description": "Hex representation of a uint",
			"pattern": "^0x([a-fA-F\\d])+$"
	}`
const hexutilUint64D = `{
          "title": "uint64",
          "type": "string",
          "description": "Hex representation of a uint64",
          "pattern": "^0x([a-fA-F\\d])+$"
        }`
const blockNumberTagD = `{
	     "title": "blockNumberTag",
	     "type": "string",
	     "description": "The block height description",
	     "enum": [
	       "earliest",
	       "latest",
	       "pending"
	     ]
		}`

var blockNumberD = fmt.Sprintf(`{
		"title": "blockNumberIdentifier",
		"oneOf": [%s, %s]
		}`, blockNumberTagD, hexutilUint64D)

const requireCanonicalD = `{
		  "type": "object",
		  "properties": {
			"requireCanonical": {
			  "type": "boolean"
			}
		  },
		  "additionalProperties": false
		}`

var blockNumberOrHashD = fmt.Sprintf(`{
          "oneOf": [
            %s,
            {
				"allOf": [%s, %s]
			}
          ]
        }`, blockNumberD, commonHashD, requireCanonicalD)

// schemaDictEntry represents a type association passed to the jsonschema reflector.
type schemaDictEntry struct {
	example interface{}
	rawJson string
}

// OpenRPCJSONSchemaTypeMapper contains the application-specific type mapping
// passed to the jsonschema reflector, used in generating a runtime representation
// of specific API objects.
func OpenRPCJSONSchemaTypeMapper(ty reflect.Type) *jsonschema.Type {
	unmarshalJSONToJSONSchemaType := func(input string) *jsonschema.Type {
		var js jsonschema.Type
		err := json.Unmarshal([]byte(input), &js)
		if err != nil {
			panic(err)
		}
		return &js
	}

	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	if ty == reflect.TypeOf((*interface{})(nil)).Elem() {
		return &jsonschema.Type{Type: "object", AdditionalProperties: []byte("true")}
	}

	// Second, handle other types.
	// Use a slice instead of a map because it preserves order, as a logic safeguard/fallback.
	dict := []schemaDictEntry{
		//{interface{}{}, fmt.Sprintf(`{
		//	"oneOf": [{"additionalProperties": true}, {"type": "null"}]
		//}`)},
		{[]byte{}, bytesD},
		{big.Int{}, integerD},
		{types.BlockNonce{}, integerD},
		{common.Address{}, commonAddressD},
		{common.Hash{}, commonHashD},
		{hexutil.Big{}, integerD},
		{hexutil.Bytes{}, hexutilBytesD},
		{hexutil.Uint(0), hexutilUintD},
		{hexutil.Uint64(0), hexutilUint64D},
		{rpc.BlockNumber(0), blockNumberD},
		{rpc.BlockNumberOrHash{}, blockNumberOrHashD},

		{rpc.Subscription{}, `{
			"type": "object",
			"title": "Subscription",
			"summary": ""
		}`},
	}

	for _, d := range dict {
		if reflect.TypeOf(d.example) == ty {
			tt := unmarshalJSONToJSONSchemaType(d.rawJson)

			return tt
		}
	}

	// Handle primitive types in case there are generic cases
	// specific to our services.
	switch ty.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Return all integer types as the hex representation integer schemea.
		ret := unmarshalJSONToJSONSchemaType(integerD)
		return ret
	case reflect.Struct:
	case reflect.Map:
	case reflect.Slice, reflect.Array:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Ptr, reflect.Interface:
	default:
	}

	return nil
}

func expandedFieldNamesFromList(in []*ast.Field) (out []*ast.Field) {
	expandedFields := []*ast.Field{}
	for _, f := range in {
		expandedFields = append(expandedFields, fieldsWithNames(f)...)
	}
	return expandedFields
}

// fieldsWithNames expands a field (either parameter or result, in this case) to
// fields which all have at least one name, or at least one field with one name.
// This handles unnamed fields, and fields declared using multiple names with one type.
// Unnamed fields are assigned a name that is the 'printed' identity of the field Type,
// eg. int -> int, bool -> bool
func fieldsWithNames(f *ast.Field) (fields []*ast.Field) {
	if f == nil {
		return nil
	}

	if len(f.Names) == 0 {
		fields = append(fields, &ast.Field{
			Doc:     f.Doc,
			Names:   []*ast.Ident{{Name: printIdentField(f)}},
			Type:    f.Type,
			Tag:     f.Tag,
			Comment: f.Comment,
		})
		return
	}
	for _, ident := range f.Names {
		fields = append(fields, &ast.Field{
			Doc:     f.Doc,
			Names:   []*ast.Ident{ident},
			Type:    f.Type,
			Tag:     f.Tag,
			Comment: f.Comment,
		})
	}
	return
}

func printIdentField(f *ast.Field) string {
	b := []byte{}
	buf := bytes.NewBuffer(b)
	fs := token.NewFileSet()
	printer.Fprint(buf, fs, f.Type.(ast.Node))
	return buf.String()
}
