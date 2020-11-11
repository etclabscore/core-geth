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
	"time"

	"github.com/alecthomas/jsonschema"
	go_openrpc_reflect "github.com/etclabscore/go-openrpc-reflect"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
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

			version := params.VersionWithMeta + "/generated-at:" + time.Now().Format(time.RFC3339)
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

	appReflector.FnIsMethodEligible = func(method reflect.Method) bool {

		// Method must be exported.
		if method.PkgPath != "" {
			return false
		}

		// Exclude methods that handle subscriptions, but do so without adhering to the conventional code pattern.
		// Eg. *filters.PublicFiltersAPI.SubscribeNewHeads handles eth_subscribe("newHeads"), but there
		// isn't a method called `eth_subscribeNewHeads`. So we blacklist all these methods and use
		// the mock subscription receiver type RPCSubscription.
		// This pattern matches all strings that start with Subscribe and are suffixed with a non-zero
		// number of A-z characters.
		if regexp.MustCompile(`^Subscribe[A-Za-z]+`).MatchString(method.Name) {
			return false
		}

		// Verify return types. The function must return at most one error
		// and/or one other non-error value.
		outs := make([]reflect.Type, method.Func.Type().NumOut())
		for i := 0; i < method.Func.Type().NumOut(); i++ {
			outs[i] = method.Func.Type().Out(i)
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

	appReflector.FnGetContentDescriptorRequired = func(r reflect.Value, m reflect.Method, field *ast.Field) (bool, error) {
		// Custom handling for eth_subscribe optional second parameter (depends on channel).
		if m.Name == "Subscribe" && len(field.Names) > 0 && field.Names[0].Name == "subscriptionOptions" {
			return false, nil
		}

		// Otherwise return the default.
		return go_openrpc_reflect.EthereumReflector.GetContentDescriptorRequired(r, m, field)
	}

	// Finally, register the configured reflector to the document.
	d.WithReflector(appReflector)
	return d
}

/*
The following struct RPCSubscription and RPCSubscription.Unsubscribe method
are documentation-only mocks to represent the otherwise invisible (un-reflected)
method.
It is appended to the OpenRPC document when the eth/api/filters.PublicFilterAPI receiver
is registered.
*/
type RPCSubscription struct{}

// Unsubscribe terminates an existing subscription by ID.
func (sub *RPCSubscription) Unsubscribe(id rpc.ID) error {
	// This is a mock function, not the real one.
	return nil
}

type RPCSubscriptionParamsName string

// Subscribe creates a subscription to an event channel.
// Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.
func (sub *RPCSubscription) Subscribe(subscriptionName RPCSubscriptionParamsName, subscriptionOptions interface{}) (subscriptionID rpc.ID, err error) {
	// This is a mock function, not the real one.
	return
}

// registerOpenRPCAPIs provides a convenience logic that is reused
// congruent to the rpc package receiver registrations.
func registerOpenRPCAPIs(doc *go_openrpc_reflect.Document, apis []rpc.API) {
	for _, api := range apis {
		doc.RegisterReceiverName(api.Namespace, api.Service)

		// Append a mock interface for the eth_unsubscribe method, which
		// would otherwise not occur in the document.
		switch api.Service.(type) {
		case *filters.PublicFilterAPI:
			doc.RegisterReceiverName("eth", &RPCSubscription{})
		}
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

const rpcSubscriptionIDD = `{
		"title": "subscriptionID",
		"type": "string",
		"description": "Subscription identifier"
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

var rpcSubscriptionParamsNameD = fmt.Sprintf(`{
		"oneOf": [
			{"type": "string", "enum": ["newHeads"], "description": "Fires a notification each time a new header is appended to the chain, including chain reorganizations."},
			{"type": "string", "enum": ["logs"], "description": "Returns logs that are included in new imported blocks and match the given filter criteria."},
			{"type": "string", "enum": ["newPendingTransactions"], "description": "Returns the hash for all transactions that are added to the pending state and are signed with a key that is available in the node."},
			{"type": "string", "enum": ["syncing"], "description": "Indicates when the node starts or stops synchronizing. The result can either be a boolean indicating that the synchronization has started (true), finished (false) or an object with various progress indicators."}
		]
	}`)

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
		return &jsonschema.Type{AdditionalProperties: []byte("true")}
	}

	// Second, handle other types.
	// Use a slice instead of a map because it preserves order, as a logic safeguard/fallback.
	dict := []schemaDictEntry{
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
		{rpc.Subscription{}, rpcSubscriptionIDD},
		{RPCSubscriptionParamsName(""), rpcSubscriptionParamsNameD},
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

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	subscriptionType = reflect.TypeOf(rpc.Subscription{})
	stringType       = reflect.TypeOf("")
)

// Is t context.Context or *context.Context?
func isContextType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t == contextType
}

// Is t Subscription or *Subscription?
func isSubscriptionType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t == subscriptionType
}

// Does t satisfy the error interface?
func isErrorType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Implements(errorType)
}

// isPubSub tests whether the given method has as as first argument a context.Context and
// returns the pair (Subscription, error).
// This function is taken directly from rpc/service.go.
func isPubSub(methodType reflect.Type) bool {
	// numIn(0) is the receiver type
	if methodType.NumIn() < 2 || methodType.NumOut() != 2 {
		return false
	}
	return isContextType(methodType.In(1)) &&
		isSubscriptionType(methodType.Out(0)) &&
		isErrorType(methodType.Out(1))
}
