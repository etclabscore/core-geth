






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-23T04:50:40-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### trace_block

Block returns the structured logs created during the execution of
EVM and returns them as a JSON object.
The correct name will be TraceBlockByNumber, though we want to be compatible with Parity trace module.


#### Params (2)

Parameters must be given _by position_.  


__1:__ 
number <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- title: `blockNumberTag`
			- type: string


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: string


	- title: `blockNumberIdentifier`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "The block height description",
                "enum": [
                    "earliest",
                    "latest",
                    "pending"
                ],
                "title": "blockNumberTag",
                "type": [
                    "string"
                ]
            },
            {
                "description": "Hex representation of a uint64",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": [
                    "string"
                ]
            }
        ],
        "title": "blockNumberIdentifier"
    }
	```




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Debug: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- Tracer: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Debug": {
                "type": "boolean"
            },
            "DisableMemory": {
                "type": "boolean"
            },
            "DisableReturnData": {
                "type": "boolean"
            },
            "DisableStack": {
                "type": "boolean"
            },
            "DisableStorage": {
                "type": "boolean"
            },
            "Limit": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Reexec": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Timeout": {
                "type": "string"
            },
            "Tracer": {
                "type": "string"
            },
            "overrides": {
                "additionalProperties": true
            }
        },
        "type": [
            "object"
        ]
    }
	```





#### Result



interface <code>[]interface{}</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `true`


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "additionalProperties": true
            }
        ],
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "trace_block", "params": [<number>, <config>]}'
	```

=== "Javascript Console"

	``` js
	trace.block(number,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateTraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([ // Block returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
// The correct name will be TraceBlockByNumber, though we want to be compatible with Parity trace module.
]interface{}, error) {
	var block *types.Block
	switch number {
	case rpc.PendingBlockNumber:
		block = api.eth.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		block = api.eth.blockchain.CurrentBlock()
	default:
		block = api.eth.blockchain.GetBlockByNumber(uint64(number))
	}
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", number)
	}
	config = setConfigTracerToParity(config)
	traceResults, err := traceBlockByNumber(ctx, api.eth, number, config)
	if err != nil {
		return nil, err
	}
	traceReward, err := traceBlockReward(ctx, api.eth, block, config)
	if err != nil {
		return nil, err
	}
	traceUncleRewards, err := traceBlockUncleRewards(ctx, api.eth, block, config)
	if err != nil {
		return nil, err
	}
	results := [ // Fetch the block that we want to trace
	]interface{}{}
	for _, result := range traceResults {
		var tmp []interface{}
		if err := json.Unmarshal(result.Result.(json.RawMessage), &tmp); err != nil {
			return nil, err
		}
		results = append(results, tmp...)
	}
	results = append(results, traceReward)
	for _, uncleReward := range traceUncleRewards {
		results = append(results, uncleReward)
	}
	return results, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer_parity.go#L114" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_filter



#### Params (2)

Parameters must be given _by position_.  


__1:__ 
args <code>ethapi.CallArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- data: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- value: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "data": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "from": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "gas": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasPrice": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "to": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "value": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            }
        },
        "type": [
            "object"
        ]
    }
	```




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Debug: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- Tracer: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Debug": {
                "type": "boolean"
            },
            "DisableMemory": {
                "type": "boolean"
            },
            "DisableReturnData": {
                "type": "boolean"
            },
            "DisableStack": {
                "type": "boolean"
            },
            "DisableStorage": {
                "type": "boolean"
            },
            "Limit": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Reexec": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Timeout": {
                "type": "string"
            },
            "Tracer": {
                "type": "string"
            },
            "overrides": {
                "additionalProperties": true
            }
        },
        "type": [
            "object"
        ]
    }
	```





#### Result



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- error: 
					- type: `string`

				- result: 
					- additionalProperties: `true`


			- type: object


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "additionalProperties": false,
                "properties": {
                    "error": {
                        "type": "string"
                    },
                    "result": {
                        "additionalProperties": true
                    }
                },
                "type": [
                    "object"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "trace_filter", "params": [<args>, <config>]}'
	```

=== "Javascript Console"

	``` js
	trace.filter(args,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateTraceAPI) Filter(ctx context.Context, args ethapi.CallArgs, config *TraceConfig) ([]*txTraceResult, error) {
	return nil, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer_parity.go#L176" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_transaction

Transaction returns the structured logs created during the execution of EVM
and returns them as a JSON object.


#### Params (2)

Parameters must be given _by position_.  


__1:__ 
hash <code>common.Hash</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Debug: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- Tracer: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Debug": {
                "type": "boolean"
            },
            "DisableMemory": {
                "type": "boolean"
            },
            "DisableReturnData": {
                "type": "boolean"
            },
            "DisableStack": {
                "type": "boolean"
            },
            "DisableStorage": {
                "type": "boolean"
            },
            "Limit": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Reexec": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Timeout": {
                "type": "string"
            },
            "Tracer": {
                "type": "string"
            },
            "overrides": {
                "additionalProperties": true
            }
        },
        "type": [
            "object"
        ]
    }
	```





#### Result



interface <code>interface{}</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "trace_transaction", "params": [<hash>, <config>]}'
	```

=== "Javascript Console"

	``` js
	trace.transaction(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateTraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	config = setConfigTracerToParity(config)
	return traceTransaction(ctx, api.eth, hash, config)
}// Transaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer_parity.go#L169" target="_">View on GitHub →</a>
</p>
</details>

---

