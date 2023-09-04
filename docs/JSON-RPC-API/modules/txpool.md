






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### txpool_content

Content returns the transactions contained within the transaction pool.


#### Params (0)

_None_

#### Result



mapstringmapstringmapstringRPCTransaction <code>map[string]map[string]map[string]*RPCTransaction</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- patternProperties: 
				- .*: 
					- patternProperties: 
						- .*: 
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

								- blockHash: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- blockNumber: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- chainId: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
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

								- hash: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- input: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `dataWord`
									- type: `string`

								- maxFeePerGas: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- maxPriorityFeePerGas: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- nonce: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`

								- r: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- s: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- to: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- transactionIndex: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`

								- type: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`

								- v: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- value: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- yParity: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`


							- type: `object`


					- type: `object`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "patternProperties": {
            ".*": {
                "patternProperties": {
                    ".*": {
                        "patternProperties": {
                            ".*": {
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
                                    "blockHash": {
                                        "pattern": "^0x[a-fA-F\\d]{64}$",
                                        "title": "keccak",
                                        "type": "string"
                                    },
                                    "blockNumber": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "chainId": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
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
                                    "hash": {
                                        "pattern": "^0x[a-fA-F\\d]{64}$",
                                        "title": "keccak",
                                        "type": "string"
                                    },
                                    "input": {
                                        "pattern": "^0x([a-fA-F\\d])+$",
                                        "title": "dataWord",
                                        "type": "string"
                                    },
                                    "maxFeePerGas": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "maxPriorityFeePerGas": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "nonce": {
                                        "pattern": "^0x([a-fA-F\\d])+$",
                                        "title": "uint64",
                                        "type": "string"
                                    },
                                    "r": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "s": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "to": {
                                        "pattern": "^0x[a-fA-F\\d]{64}$",
                                        "title": "keccak",
                                        "type": "string"
                                    },
                                    "transactionIndex": {
                                        "pattern": "^0x([a-fA-F\\d])+$",
                                        "title": "uint64",
                                        "type": "string"
                                    },
                                    "type": {
                                        "pattern": "^0x([a-fA-F\\d])+$",
                                        "title": "uint64",
                                        "type": "string"
                                    },
                                    "v": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "value": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "yParity": {
                                        "pattern": "^0x([a-fA-F\\d])+$",
                                        "title": "uint64",
                                        "type": "string"
                                    }
                                },
                                "type": "object"
                            }
                        },
                        "type": "object"
                    }
                },
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "txpool_content", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "txpool_content", "params": []}'
	```


=== "Javascript Console"

	``` js
	txpool.content();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TxPoolAPI) Content() map // Content returns the transactions contained within the transaction pool.
[string]map[string]map[string]*RPCTransaction {
	content := map[string]map[string]map[string]*RPCTransaction{"pending": make(map[string]map[string]*RPCTransaction), "queued": make(map[string]map[string]*RPCTransaction)}
	pending, queue := s.b.TxPoolContent()
	curHeader := s.b.CurrentHeader()
	for account, txs := range pending {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
		}
		content["pending"][account.Hex()] = dump
	}
	for account, txs := range queue {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L166" target="_">View on GitHub →</a>
</p>
</details>

---



### txpool_contentFrom

ContentFrom returns the transactions contained within the transaction pool.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
addr <code>common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- title: `keccak`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a Keccak 256 hash POINTER",
        "pattern": "^0x[a-fA-F\\d]{64}$",
        "title": "keccak",
        "type": [
            "string"
        ]
    }
	```





#### Result



mapstringmapstringRPCTransaction <code>map[string]map[string]*RPCTransaction</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- patternProperties: 
				- .*: 
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

						- blockHash: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- blockNumber: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- chainId: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
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

						- hash: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- input: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `dataWord`
							- type: `string`

						- maxFeePerGas: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- maxPriorityFeePerGas: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- nonce: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`
							- type: `string`

						- r: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- s: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- to: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- transactionIndex: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`
							- type: `string`

						- type: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`
							- type: `string`

						- v: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- value: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- yParity: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`
							- type: `string`


					- type: `object`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "patternProperties": {
            ".*": {
                "patternProperties": {
                    ".*": {
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
                            "blockHash": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "blockNumber": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "chainId": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
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
                            "hash": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "input": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "dataWord",
                                "type": "string"
                            },
                            "maxFeePerGas": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "maxPriorityFeePerGas": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "nonce": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "uint64",
                                "type": "string"
                            },
                            "r": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "s": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "to": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "transactionIndex": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "uint64",
                                "type": "string"
                            },
                            "type": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "uint64",
                                "type": "string"
                            },
                            "v": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "value": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "yParity": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "uint64",
                                "type": "string"
                            }
                        },
                        "type": "object"
                    }
                },
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "txpool_contentFrom", "params": [<addr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "txpool_contentFrom", "params": [<addr>]}'
	```


=== "Javascript Console"

	``` js
	txpool.contentFrom(addr);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TxPoolAPI) ContentFrom(addr common.Address) map // ContentFrom returns the transactions contained within the transaction pool.
[string]map[string]*RPCTransaction {
	content := make(map[string]map[string]*RPCTransaction, 2)
	pending, queue := s.b.TxPoolContentFrom(addr)
	curHeader := s.b.CurrentHeader()
	dump := make(map[string]*RPCTransaction, len(pending))
	for _, tx := range pending {
		dump[fmt.Sprintf("%d", tx.Nonce())] = NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
	}
	content["pending"] = dump
	dump = make(map[string]*RPCTransaction, len(queue))
	for _, tx := range queue {
		dump[fmt.Sprintf("%d", tx.Nonce())] = NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig())
	}
	content["queued"] = dump
	return content
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L193" target="_">View on GitHub →</a>
</p>
</details>

---



### txpool_inspect

Inspect retrieves the content of the transaction pool and flattens it into an
easily inspectable list.


#### Params (0)

_None_

#### Result



mapstringmapstringmapstringstring <code>map[string]map[string]map[string]string</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- patternProperties: 
				- .*: 
					- patternProperties: 
						- .*: 
							- type: `string`


					- type: `object`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "patternProperties": {
            ".*": {
                "patternProperties": {
                    ".*": {
                        "patternProperties": {
                            ".*": {
                                "type": "string"
                            }
                        },
                        "type": "object"
                    }
                },
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "txpool_inspect", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "txpool_inspect", "params": []}'
	```


=== "Javascript Console"

	``` js
	txpool.inspect();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TxPoolAPI) Inspect() map // Inspect retrieves the content of the transaction pool and flattens it into an
// easily inspectable list.
[string]map[string]map[string]string {
	content := map[string]map[string]map[string]string{"pending": make(map[string]map[string]string), "queued": make(map[string]map[string]string)}
	pending, queue := s.b.TxPoolContent()
	var format = func(tx *types.Transaction) string {
		if to := tx.To(); to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei", tx.To().Hex(), tx.Value(), tx.Gas(), tx.GasPrice())
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei", tx.Value(), tx.Gas(), tx.GasPrice())
	}
	for account, txs := // Define a formatter to flatten a transaction into a string
	range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	for account, txs := range queue {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L226" target="_">View on GitHub →</a>
</p>
</details>

---



### txpool_status

Status returns the number of pending and queued transaction in the pool.


#### Params (0)

_None_

#### Result



mapstringhexutilUint <code>map[string]hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- description: `Hex representation of a uint`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "patternProperties": {
            ".*": {
                "description": "Hex representation of a uint",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint",
                "type": "string"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "txpool_status", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "txpool_status", "params": []}'
	```


=== "Javascript Console"

	``` js
	txpool.status();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TxPoolAPI) Status() map // Status returns the number of pending and queued transaction in the pool.
[string]hexutil.Uint {
	pending, queue := s.b.Stats()
	return map[string]hexutil.Uint{"pending": hexutil.Uint(pending), "queued": hexutil.Uint(queue)}
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L216" target="_">View on GitHub →</a>
</p>
</details>

---

