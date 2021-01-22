






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-22T08:53:19-06:00</code> |
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
								- input: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `dataWord`
									- type: `string`

								- r: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- s: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- value: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- from: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- hash: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- blockNumber: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- v: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- gasPrice: 
									- pattern: `^0x[a-fA-F0-9]+$`
									- title: `integer`
									- type: `string`

								- to: 
									- title: `keccak`
									- type: `string`
									- pattern: `^0x[a-fA-F\d]{64}$`

								- nonce: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`

								- transactionIndex: 
									- pattern: `^0x([a-fA-F\d])+$`
									- title: `uint64`
									- type: `string`

								- blockHash: 
									- pattern: `^0x[a-fA-F\d]{64}$`
									- title: `keccak`
									- type: `string`

								- gas: 
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
                                    "v": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
                                        "type": "string"
                                    },
                                    "value": {
                                        "pattern": "^0x[a-fA-F0-9]+$",
                                        "title": "integer",
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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "txpool_content", "params": []}'
	```

=== "Javascript Console"

	``` js
	txpool.content();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTxPoolAPI) Content() map // Content returns the transactions contained within the transaction pool.
[string]map[string]map[string]*RPCTransaction {
	content := map[string]map[string]map[string]*RPCTransaction{"pending": make(map[string]map[string]*RPCTransaction), "queued": make(map[string]map[string]*RPCTransaction)}
	pending, queue := s.b.TxPoolContent()
	for account, txs := range pending {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	for account, txs := range queue {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L108" target="_">View on GitHub →</a>
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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "txpool_inspect", "params": []}'
	```

=== "Javascript Console"

	``` js
	txpool.inspect();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTxPoolAPI) Inspect() map // Inspect retrieves the content of the transaction pool and flattens it into an
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L145" target="_">View on GitHub →</a>
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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "txpool_status", "params": []}'
	```

=== "Javascript Console"

	``` js
	txpool.status();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTxPoolAPI) Status() map // Status returns the number of pending and queued transaction in the pool.
[string]hexutil.Uint {
	pending, queue := s.b.Stats()
	return map[string]hexutil.Uint{"pending": hexutil.Uint(pending), "queued": hexutil.Uint(queue)}
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L135" target="_">View on GitHub →</a>
</p>
</details>

---

