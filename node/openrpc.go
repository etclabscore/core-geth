package node

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/alecthomas/jsonschema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

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
	     "description": "The optional block height description",
	     "enum": [
	       "earliest",
	       "latest",
	       "pending"
	     ]
	   }`

var blockNumberOrHashD = fmt.Sprintf(`{
          "oneOf": [
            %s,
            %s
          ]
        }`, blockNumberTagD, commonHashD)

var emptyInterface interface{}

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

	// Second, handle other types.
	// Use a slice instead of a map because it preserves order, as a logic safeguard/fallback.
	dict := []schemaDictEntry{
		{emptyInterface, fmt.Sprintf(`{
			"oneOf": [{"additionalProperties": true}, {"type": "null"}]
		}`)},
		{[]byte{}, bytesD},
		{big.Int{}, integerD},
		{types.BlockNonce{}, integerD},
		{common.Address{}, commonAddressD},
		{common.Hash{}, commonHashD},
		{hexutil.Big{}, integerD},
		{hexutil.Bytes{}, hexutilBytesD},
		{hexutil.Uint(0), hexutilUintD},
		{hexutil.Uint64(0), hexutilUint64D},
		{rpc.BlockNumber(0), blockNumberOrHashD},
		{rpc.BlockNumberOrHash{}, fmt.Sprintf(`{
		  "title": "blockNumberOrHash",
		  "oneOf": [
			%s,
			{
			  "allOf": [
				%s,
				{
				  "type": "object",
				  "properties": {
					"requireCanonical": {
					  "type": "boolean"
					}
				  },
				  "additionalProperties": false
				}
			  ]
			}
		  ]
		}`, blockNumberOrHashD, blockNumberOrHashD)},
		{rpc.Subscription{}, fmt.Sprintf(`{
			"type": "object",
			"title": "Subscription",
			"summary": ""
		}`)},
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
