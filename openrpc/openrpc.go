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

package openrpc

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

type schemaDictEntry struct {
	t interface{}
	j string
}

func OpenRPCJSONSchemaTypeMapper(r reflect.Type) *jsonschema.Type {
	unmarshalJSONToJSONSchemaType := func(input string) *jsonschema.Type {
		var js jsonschema.Type
		err := json.Unmarshal([]byte(input), &js)
		if err != nil {
			panic(err)
		}
		return &js
	}

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

	blockNumberOrHashD := fmt.Sprintf(`{
          "oneOf": [
            %s,
            %s
          ]
        }`, blockNumberTagD, commonHashD)

	//s := jsonschema.Reflect(ethapi.Progress{})
	//ethSyncingResultProgress, err := json.Marshal(s)
	//if err != nil {
	//	return nil
	//}

	// Second, handle other types.
	// Use a slice instead of a map because it preserves order, as a logic safeguard/fallback.
	dict := []schemaDictEntry{

		{new(big.Int), integerD},
		{big.Int{}, integerD},
		{new(hexutil.Big), integerD},
		{hexutil.Big{}, integerD},

		{types.BlockNonce{}, integerD},

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

		{rpc.BlockNumber(0),
			blockNumberOrHashD,
		},

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

		//{
		//	BlockNumber(0): blockNumberOrHashD,
		//},

		//{BlockNumberOrHash{}, fmt.Sprintf(`{
		//  "title": "blockNumberOrHash",
		//  "description": "Hex representation of a block number or hash",
		//  "oneOf": [%s, %s]
		//}`, commonHashD, integerD)},

		//{BlockNumber(0), fmt.Sprintf(`{
		//  "title": "blockNumberOrTag",
		//  "description": "Block tag or hex representation of a block number",
		//  "oneOf": [%s, %s]
		//}`, commonHashD, blockNumberTagD)},

		//		{ethapi.EthSyncingResult{}, fmt.Sprintf(`{
		//          "title": "ethSyncingResult",
		//		  "description": "Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
		//yet received the latest block headers from its pears. In case it is synchronizing:
		//- startingBlock: block number this node started to synchronise from
		//- currentBlock:  block number this node is currently importing
		//- highestBlock:  block number of the highest block header this node has received from peers
		//- pulledStates:  number of state entries processed until now
		//- knownStates:   number of known state entries that still need to be pulled",
		//		  "oneOf": [%s, %s]
		//		}`, `{
		//        "type": "boolean"
		//      }`, `{"type": "object"}`)},

	}

	for _, d := range dict {
		d := d
		if reflect.TypeOf(d.t) == r {
			tt := unmarshalJSONToJSONSchemaType(d.j)

			return tt
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