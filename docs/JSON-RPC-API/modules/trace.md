






| Entity | Version |
| --- | --- |
| Source | <code>1.11.23-unstable/generated-at:2021-04-30T19:24:24+03:00</code> |
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

		- NestedTraceOutput: 
			- type: `boolean`

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
            "NestedTraceOutput": {
                "type": "boolean"
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "trace_block", "params": [<number>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_block", "params": [<number>, <config>]}'
	```


=== "Javascript Console"

	``` js
	trace.block(number,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *TraceAPI) Block(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([ // Block returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
// The correct name will be TraceBlockByNumber, though we want to be compatible with Parity trace module.
]interface{}, error) {
	config = setTraceConfigDefaultTracer(config)
	block, err := api.debugAPI.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	traceResults, err := api.debugAPI.traceBlock(ctx, block, config)
	if err != nil {
		return nil, err
	}
	traceReward, err := api.traceBlockReward(ctx, block, config)
	if err != nil {
		return nil, err
	}
	traceUncleRewards, err := api.traceBlockUncleRewards(ctx, block, config)
	if err != nil {
		return nil, err
	}
	results := []interface{}{}
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api_parity.go#L175" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_call

Call lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
if the given transaction was added on top of the provided block and returns them as a JSON object.
You can provide -2 as a block number to trace on top of the pending block.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
args <code>ethapi.CallArgs</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accessList: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- storageKeys: 
						- items: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- type: `array`


				- type: `object`

			- type: `array`

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
            "accessList": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "storageKeys": {
                            "items": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "type": "array"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
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
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes





__3:__ 
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

		- NestedTraceOutput: 
			- type: `boolean`

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
            "NestedTraceOutput": {
                "type": "boolean"
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "trace_call", "params": [<args>, <blockNrOrHash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_call", "params": [<args>, <blockNrOrHash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	trace.call(args,blockNrOrHash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *TraceAPI) Call(ctx context.Context, args ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	config = setTraceConfigDefaultTracer(config)
	res, err := api.debugAPI.TraceCall(ctx, args, blockNrOrHash, config)
	if err != nil {
		return nil, err
	}
	return decorateResponse(res, config)
}// Call lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api_parity.go#L240" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_callMany

CallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
if the given transaction was added on top of the provided block and returns them as a JSON object.
You can provide -2 as a block number to trace on top of the pending block.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
txs <code>[]ethapi.CallArgs</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- accessList: 
					- items: 
						- additionalProperties: `false`
						- properties: 
							- address: 
								- pattern: `^0x[a-fA-F\d]{64}$`
								- title: `keccak`
								- type: `string`

							- storageKeys: 
								- items: 
									- description: `Hex representation of a Keccak 256 hash`
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- type: `array`


						- type: `object`

					- type: `array`

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


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "additionalProperties": false,
                "properties": {
                    "accessList": {
                        "items": {
                            "additionalProperties": false,
                            "properties": {
                                "address": {
                                    "pattern": "^0x[a-fA-F\\d]{64}$",
                                    "title": "keccak",
                                    "type": "string"
                                },
                                "storageKeys": {
                                    "items": {
                                        "description": "Hex representation of a Keccak 256 hash",
                                        "pattern": "^0x[a-fA-F\\d]{64}$",
                                        "title": "keccak",
                                        "type": "string"
                                    },
                                    "type": "array"
                                }
                            },
                            "type": "object"
                        },
                        "type": "array"
                    },
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
        ],
        "type": [
            "array"
        ]
    }
	```




__2:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes





__3:__ 
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

		- NestedTraceOutput: 
			- type: `boolean`

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
            "NestedTraceOutput": {
                "type": "boolean"
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "trace_callMany", "params": [<txs>, <blockNrOrHash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_callMany", "params": [<txs>, <blockNrOrHash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	trace.callMany(txs,blockNrOrHash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *TraceAPI) CallMany(ctx context.Context, txs [ // CallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
]ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	config = setTraceConfigDefaultTracer(config)
	return api.debugAPI.TraceCallMany(ctx, txs, blockNrOrHash, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api_parity.go#L252" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_filter

Filter configures a new tracer according to the provided configuration, and
executes all the transactions contained within. The return value will be one item
per transaction, dependent on the requested tracer.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
args <code>TraceFilterArgs</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- after: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- count: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- fromAddress: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- fromBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- toAddress: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- toBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "after": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "count": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "fromAddress": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "fromBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "toAddress": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "toBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
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

		- NestedTraceOutput: 
			- type: `boolean`

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
            "NestedTraceOutput": {
                "type": "boolean"
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




<code>*rpc.Subscription</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Subscription identifier`
	- title: `subscriptionID`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Subscription identifier",
        "title": "subscriptionID",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples








=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_subscribe", "params": ["filter", <args>, <config>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *TraceAPI) Filter(ctx context.Context, args TraceFilterArgs, config *TraceConfig) (*rpc.Subscription, error) {
	config = setTraceConfigDefaultTracer(config)
	start := rpc.BlockNumber(args.FromBlock)
	end := rpc.BlockNumber(args.ToBlock)
	return api.debugAPI.TraceChain(ctx, start, end, config)
}// Filter configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requested tracer.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api_parity.go#L227" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_subscribe

Subscribe creates a subscription to an event channel.
Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
subscriptionName <code>RPCTraceSubscriptionParamsName</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `Returns transaction traces for the filtered addresses within a range of blocks.`
			- enum: filter
			- type: string


	- title: `subscriptionName`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "Returns transaction traces for the filtered addresses within a range of blocks.",
                "enum": [
                    "filter"
                ],
                "type": [
                    "string"
                ]
            }
        ],
        "title": "subscriptionName"
    }
	```




__2:__ 
subscriptionOptions <code>interface{}</code> 

  + Required: No






#### Result



subscriptionID <code>rpc.ID</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Subscription identifier`
	- title: `subscriptionID`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Subscription identifier",
        "title": "subscriptionID",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples






=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_subscribe", "params": [<subscriptionName>, <subscriptionOptions>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (sub *RPCTraceSubscription) Subscribe(subscriptionName RPCTraceSubscriptionParamsName, subscriptionOptions interface{}) (subscriptionID rpc.ID, err error) {
	return
}// Subscribe creates a subscription to an event channel.
// Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L269" target="_">View on GitHub →</a>
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

		- NestedTraceOutput: 
			- type: `boolean`

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
            "NestedTraceOutput": {
                "type": "boolean"
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "trace_transaction", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_transaction", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	trace.transaction(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *TraceAPI) Transaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	config = setTraceConfigDefaultTracer(config)
	return api.debugAPI.TraceTransaction(ctx, hash, config)
}// Transaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api_parity.go#L219" target="_">View on GitHub →</a>
</p>
</details>

---



### trace_unsubscribe

Unsubscribe terminates an existing subscription by ID.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
id <code>rpc.ID</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Subscription identifier`
	- title: `subscriptionID`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Subscription identifier",
        "title": "subscriptionID",
        "type": [
            "string"
        ]
    }
	```





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "trace_unsubscribe", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "trace_unsubscribe", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	trace.unsubscribe(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (sub *RPCTraceSubscription) Unsubscribe(id rpc.ID) error {
	return nil
}// Unsubscribe terminates an existing subscription by ID.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L260" target="_">View on GitHub →</a>
</p>
</details>

---

