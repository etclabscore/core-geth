






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### eth_accounts

Accounts returns the collection of accounts this node manages.


#### Params (0)

_None_

#### Result



commonAddress <code>[]common.Address</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- description: `Hex representation of a Keccak 256 hash POINTER`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: string


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "description": "Hex representation of a Keccak 256 hash POINTER",
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": [
                    "string"
                ]
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_accounts", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_accounts", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.accounts();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *EthereumAccountAPI) Accounts() [ // Accounts returns the collection of accounts this node manages.
]common.Address {
	return s.am.Accounts()
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L271" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_blockNumber

BlockNumber returns the block number of the chain head.


#### Params (0)

_None_

#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_blockNumber", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.blockNumber();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	return hexutil.Uint64(header.Number.Uint64())
}// BlockNumber returns the block number of the chain head.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L628" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_call

Call executes the given transaction on the state for the given block number.

Additionally, the caller can specify a batch of contract for fields overriding.

Note, this function doesn't make and changes in the state/blockchain and is
useful to execute and retrieve values.


#### Params (4)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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
overrides <code>*StateOverride</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- additionalProperties: `false`
			- properties: 
				- balance: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- code: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `dataWord`
					- type: `string`

				- nonce: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `uint64`
					- type: `string`

				- state: 
					- patternProperties: 
						- .*: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`


					- type: `object`

				- stateDiff: 
					- patternProperties: 
						- .*: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
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
                "additionalProperties": false,
                "properties": {
                    "balance": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "code": {
                        "pattern": "^0x([a-fA-F\\d])+$",
                        "title": "dataWord",
                        "type": "string"
                    },
                    "nonce": {
                        "pattern": "^0x([a-fA-F\\d])+$",
                        "title": "uint64",
                        "type": "string"
                    },
                    "state": {
                        "patternProperties": {
                            ".*": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            }
                        },
                        "type": "object"
                    },
                    "stateDiff": {
                        "patternProperties": {
                            ".*": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
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




__4:__ 
blockOverrides <code>*BlockOverrides</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- BaseFee: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Coinbase: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- Difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- GasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- Number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Random: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- Time: 
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
            "BaseFee": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Coinbase": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "Difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "GasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "Number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Random": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "Time": {
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





#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_call", "params": [<args>, <blockNrOrHash>, <overrides>, <blockOverrides>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_call", "params": [<args>, <blockNrOrHash>, <overrides>, <blockOverrides>]}'
	```


=== "Javascript Console"

	``` js
	eth.call(args,blockNrOrHash,overrides,blockOverrides);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) Call(ctx context.Context, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	result, err := DoCall(ctx, s.b, args, blockNrOrHash, overrides, blockOverrides, s.b.RPCEVMTimeout(), s.b.RPCGasCap())
	if err != nil {
		return nil, err
	}
	if len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}
	return result.Return(), result.Err
}// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1126" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_chainId

ChainId is the EIP-155 replay-protection chain id for the current Ethereum chain config.

Note, this method does not conform to EIP-695 because the configured chain ID is always
returned, regardless of the current head block. We used to return an error when the chain
wasn't synced up to a block where EIP-155 is enabled, but this behavior caused issues
in CL clients.


#### Params (0)

_None_

#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_chainId", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_chainId", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.chainId();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *BlockChainAPI) ChainId() *hexutil.Big {
	return (*hexutil.Big)(api.b.ChainConfig().GetChainID())
}// ChainId is the EIP-155 replay-protection chain id for the current Ethereum chain config.
//
// Note, this method does not conform to EIP-695 because the configured chain ID is always
// returned, regardless of the current head block. We used to return an error when the chain
// wasn't synced up to a block where EIP-155 is enabled, but this behavior caused issues
// in CL clients.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L623" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_coinbase

Coinbase is the address that mining rewards will be sent to (alias for Etherbase).


#### Params (0)

_None_

#### Result




<code>common.Address</code> 

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_coinbase", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_coinbase", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.coinbase();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *EthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}// Coinbase is the address that mining rewards will be sent to (alias for Etherbase).

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L40" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_createAccessList

CreateAccessList creates an EIP-2930 type AccessList for the given transaction.
Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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
blockNrOrHash <code>*rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>*accessListResult</code> 

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

		- error: 
			- type: `string`

		- gasUsed: 
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
            "error": {
                "type": "string"
            },
            "gasUsed": {
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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_createAccessList", "params": [<args>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_createAccessList", "params": [<args>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.createAccessList(args,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	acl, gasUsed, vmerr, err := AccessList(ctx, s.b, bNrOrHash, args)
	if err != nil {
		return nil, err
	}
	result := &accessListResult{Accesslist: &acl, GasUsed: hexutil.Uint64(gasUsed)}
	if vmerr != nil {
		result.Error = vmerr.Error()
	}
	return result, nil
}// CreateAccessList creates an EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1645" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_estimateGas

EstimateGas returns an estimate of the amount of gas needed to execute the
given transaction against the current pending block.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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
blockNrOrHash <code>*rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes





__3:__ 
overrides <code>*StateOverride</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- additionalProperties: `false`
			- properties: 
				- balance: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- code: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `dataWord`
					- type: `string`

				- nonce: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `uint64`
					- type: `string`

				- state: 
					- patternProperties: 
						- .*: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`


					- type: `object`

				- stateDiff: 
					- patternProperties: 
						- .*: 
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
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
                "additionalProperties": false,
                "properties": {
                    "balance": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "code": {
                        "pattern": "^0x([a-fA-F\\d])+$",
                        "title": "dataWord",
                        "type": "string"
                    },
                    "nonce": {
                        "pattern": "^0x([a-fA-F\\d])+$",
                        "title": "uint64",
                        "type": "string"
                    },
                    "state": {
                        "patternProperties": {
                            ".*": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            }
                        },
                        "type": "object"
                    },
                    "stateDiff": {
                        "patternProperties": {
                            ".*": {
                                "description": "Hex representation of a Keccak 256 hash",
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
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





#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_estimateGas", "params": [<args>, <blockNrOrHash>, <overrides>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_estimateGas", "params": [<args>, <blockNrOrHash>, <overrides>]}'
	```


=== "Javascript Console"

	``` js
	eth.estimateGas(args,blockNrOrHash,overrides);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, overrides, s.b.RPCGasCap())
}// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1273" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_etherbase

Etherbase is the address that mining rewards will be sent to.


#### Params (0)

_None_

#### Result




<code>common.Address</code> 

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_etherbase", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_etherbase", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.etherbase();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *EthereumAPI) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}// Etherbase is the address that mining rewards will be sent to.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L35" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_feeHistory

FeeHistory returns the fee market history.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
blockCount <code>math.HexOrDecimal64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
lastBlock <code>rpc.BlockNumber</code> 

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




__3:__ 
rewardPercentiles <code>[]float64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- type: number


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "type": [
                    "number"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```





#### Result




<code>*feeHistoryResult</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- type: `array`

		- gasUsedRatio: 
			- items: 
				- type: `number`

			- type: `array`

		- oldestBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- reward: 
			- items: 
				- items: 
					- description: `Hex representation of the integer`
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- type: `array`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "type": "array"
            },
            "gasUsedRatio": {
                "items": {
                    "type": "number"
                },
                "type": "array"
            },
            "oldestBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "reward": {
                "items": {
                    "items": {
                        "description": "Hex representation of the integer",
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "type": "array"
                },
                "type": "array"
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_feeHistory", "params": [<blockCount>, <lastBlock>, <rewardPercentiles>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_feeHistory", "params": [<blockCount>, <lastBlock>, <rewardPercentiles>]}'
	```


=== "Javascript Console"

	``` js
	eth.feeHistory(blockCount,lastBlock,rewardPercentiles);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *EthereumAPI) FeeHistory(ctx context.Context, blockCount math.HexOrDecimal64, lastBlock rpc.BlockNumber, rewardPercentiles [ // FeeHistory returns the fee market history.
]float64) (*feeHistoryResult, error) {
	oldest, reward, baseFee, gasUsed, err := s.b.FeeHistory(ctx, uint64(blockCount), lastBlock, rewardPercentiles)
	if err != nil {
		return nil, err
	}
	results := &feeHistoryResult{OldestBlock: (*hexutil.Big)(oldest), GasUsedRatio: gasUsed}
	if reward != nil {
		results.Reward = make([][]*hexutil.Big, len(reward))
		for i, w := range reward {
			results.Reward[i] = make([]*hexutil.Big, len(w))
			for j, v := range w {
				results.Reward[i][j] = (*hexutil.Big)(v)
			}
		}
	}
	if baseFee != nil {
		results.BaseFee = make([]*hexutil.Big, len(baseFee))
		for i, v := range baseFee {
			results.BaseFee[i] = (*hexutil.Big)(v)
		}
	}
	return results, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L94" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_fillTransaction

FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
on a given unsigned transaction, and returns it to the caller for further
processing (signing + broadcast).


#### Params (1)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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





#### Result




<code>*SignTransactionResult</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- raw: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- tx: 
			- additionalProperties: `false`
			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "raw": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "tx": {
                "additionalProperties": false,
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_fillTransaction", "params": [<args>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_fillTransaction", "params": [<args>]}'
	```


=== "Javascript Console"

	``` js
	eth.fillTransaction(args);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) FillTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx := args.toTransaction()
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1996" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_gasPrice

GasPrice returns a suggestion for a gas price for legacy transactions.


#### Params (0)

_None_

#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_gasPrice", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_gasPrice", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.gasPrice();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *EthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	if head := s.b.CurrentHeader(); head.BaseFee != nil {
		tipcap.Add(tipcap, head.BaseFee)
	}
	return (*hexutil.Big)(tipcap), err
}// GasPrice returns a suggestion for a gas price for legacy transactions.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L66" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBalance

GetBalance returns the amount of wei for the given address in the state of the
given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
block numbers are also allowed.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
address <code>common.Address</code> 

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




__2:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getBalance", "params": [<address>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getBalance", "params": [<address>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getBalance(address,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L636" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBlockByHash

GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
detail, otherwise only the transaction hash is returned.


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
fullTx <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>*RPCMarshalBlockT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- error: 
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- uncles: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- withdrawals: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- amount: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- index: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- validatorIndex: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- type: `array`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "error": {
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "size": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactions": {
                "items": {
                    "additionalProperties": true
                },
                "type": "array"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "uncles": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "withdrawals": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "amount": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "index": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "validatorIndex": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getBlockByHash", "params": [<hash>, <fullTx>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getBlockByHash", "params": [<hash>, <fullTx>]}'
	```


=== "Javascript Console"

	``` js
	eth.getBlockByHash(hash,fullTx);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (*RPCMarshalBlockT, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		return s.rpcMarshalBlock(ctx, block, true, fullTx)
	}
	return nil, err
}// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L817" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBlockByNumber

GetBlockByNumber returns the requested canonical block.
  - When blockNr is -1 the chain pending block is returned.
  - When blockNr is -2 the chain latest block is returned.
  - When blockNr is -3 the chain finalized block is returned.
  - When blockNr is -4 the chain safe block is returned.
  - When fullTx is true all transactions in the block are returned, otherwise
    only the transaction hash is returned.


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
fullTx <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>*RPCMarshalBlockT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- error: 
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- uncles: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- withdrawals: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- amount: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- index: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- validatorIndex: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- type: `array`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "error": {
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "size": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactions": {
                "items": {
                    "additionalProperties": true
                },
                "type": "array"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "uncles": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "withdrawals": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "amount": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "index": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "validatorIndex": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getBlockByNumber", "params": [<number>, <fullTx>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getBlockByNumber", "params": [<number>, <fullTx>]}'
	```


=== "Javascript Console"

	``` js
	eth.getBlockByNumber(number,fullTx);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (*RPCMarshalBlockT, error) {
	block, err := s.b.BlockByNumber(ctx, number)
	if block != nil && err == nil {
		response, err := s.rpcMarshalBlock(ctx, block, true, fullTx)
		if err == nil && number == rpc.PendingBlockNumber {
			response.setAsPending()
		}
		return response, err
	}
	return nil, err
}// GetBlockByNumber returns the requested canonical block.
//   - When blockNr is -1 the chain pending block is returned.
//   - When blockNr is -2 the chain latest block is returned.
//   - When blockNr is -3 the chain finalized block is returned.
//   - When blockNr is -4 the chain safe block is returned.
//   - When fullTx is true all transactions in the block are returned, otherwise
//     only the transaction hash is returned.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L802" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBlockTransactionCountByHash

GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockHash <code>common.Hash</code> 

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





#### Result




<code>*hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getBlockTransactionCountByHash", "params": [<blockHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getBlockTransactionCountByHash", "params": [<blockHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getBlockTransactionCountByHash(blockHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1757" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBlockTransactionCountByNumber

GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

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





#### Result




<code>*hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getBlockTransactionCountByNumber", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getBlockTransactionCountByNumber", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	eth.getBlockTransactionCountByNumber(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1748" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getCode

GetCode returns the code stored at the given address in the state for the given block number.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
address <code>common.Address</code> 

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




__2:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getCode", "params": [<address>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getCode", "params": [<address>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getCode(address,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}// GetCode returns the code stored at the given address in the state for the given block number.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L876" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getFilterChanges

GetFilterChanges returns the logs for the filter with the given id since
last time it was called. This can be used for polling.

For pending transaction and block filters the result is []common.Hash.
(pending)Log filters return []Log.


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



interface <code>interface{}</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getFilterChanges", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getFilterChanges", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	eth.getFilterChanges(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()
	chainConfig := api.sys.backend.ChainConfig()
	latest := api.sys.backend.CurrentHeader()
	if f, found := api.filters[id]; found {
		if !f.deadline.Stop() {
			<-f.deadline.C
		}
		f.deadline.Reset(api.timeout)
		switch f.typ {
		case BlocksSubscription, SideBlocksSubscription:
			hashes := f.hashes
			f.hashes = nil
			return returnHashes(hashes), nil
		case PendingTransactionsSubscription:
			if f.fullTx {
				txs := make([ // GetFilterChanges returns the logs for the filter with the given id since
				// last time it was called. This can be used for polling.
				//
				// For pending transaction and block filters the result is []common.Hash.
				// (pending)Log filters return []Log.
				]*ethapi.RPCTransaction, 0, len(f.txs))
				for _, tx := range f.txs {
					txs = append(txs, ethapi.NewRPCPendingTransaction(tx, latest, chainConfig))
				}
				f.txs = nil
				return txs, nil
			} else {
				hashes := make([]common.Hash, 0, len(f.txs))
				for _, tx := range f.txs {
					hashes = append(hashes, tx.Hash())
				}
				f.txs = nil
				return hashes, nil
			}
		case LogsSubscription, MinedAndPendingLogsSubscription:
			logs := f.logs
			f.logs = nil
			return returnLogs(logs), nil
		}
	}
	return []interface{}{}, errFilterNotFound
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L480" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getFilterLogs

GetFilterLogs returns the logs for the filter with the given id.
If the filter could not be found an empty array of logs is returned.


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



typesLog <code>[]*types.Log</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- address: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockNumber: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- data: 
					- pattern: `^0x([a-fA-F0-9]?)+$`
					- title: `bytes`
					- type: `string`

				- logIndex: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- removed: 
					- type: `boolean`

				- topics: 
					- items: 
						- description: `Hex representation of a Keccak 256 hash`
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- type: `array`

				- transactionHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- transactionIndex: 
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
                    "address": {
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
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
                    "data": {
                        "pattern": "^0x([a-fA-F0-9]?)+$",
                        "title": "bytes",
                        "type": "string"
                    },
                    "logIndex": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "removed": {
                        "type": "boolean"
                    },
                    "topics": {
                        "items": {
                            "description": "Hex representation of a Keccak 256 hash",
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "type": "array"
                    },
                    "transactionHash": {
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "transactionIndex": {
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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getFilterLogs", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getFilterLogs", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	eth.getFilterLogs(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([ // GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
]*types.Log, error) {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()
	if !found || f.typ != LogsSubscription {
		return nil, errFilterNotFound
	}
	var filter *Filter
	if f.crit.BlockHash != nil {
		filter = api.sys.NewBlockFilter(*f.crit.BlockHash, f.crit.Addresses, f.crit.Topics)
	} else {
		begin := rpc.LatestBlockNumber.Int64()
		if f.crit.FromBlock != nil {
			begin = f.crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if f.crit.ToBlock != nil {
			end = f.crit.ToBlock.Int64()
		}
		filter = api.sys.NewRangeFilter(begin, end, f.crit.Addresses, f.crit.Topics)
	}
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L441" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getHashrate

GetHashrate returns the current hashrate for local CPU miner and remote miner.


#### Params (0)

_None_

#### Result




<code>uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getHashrate", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getHashrate", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.getHashrate();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) GetHashrate() uint64 {
	return uint64(api.ethash.Hashrate())
}// GetHashrate returns the current hashrate for local CPU miner and remote miner.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/consensus/ethash/api.go#L111" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getHeaderByHash

GetHeaderByHash returns the requested header by hash.


#### Params (1)

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





#### Result




<code>*RPCMarshalHeaderT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getHeaderByHash", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getHeaderByHash", "params": [<hash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getHeaderByHash(hash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) *RPCMarshalHeaderT {
	header, _ := s.b.HeaderByHash(ctx, hash)
	if header != nil {
		return s.rpcMarshalHeader(ctx, header)
	}
	return nil
}// GetHeaderByHash returns the requested header by hash.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L787" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getHeaderByNumber

GetHeaderByNumber returns the requested canonical block header.
  - When blockNr is -1 the chain pending header is returned.
  - When blockNr is -2 the chain latest header is returned.
  - When blockNr is -3 the chain finalized header is returned.
  - When blockNr is -4 the chain safe header is returned.


#### Params (1)

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





#### Result




<code>*RPCMarshalHeaderT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getHeaderByNumber", "params": [<number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getHeaderByNumber", "params": [<number>]}'
	```


=== "Javascript Console"

	``` js
	eth.getHeaderByNumber(number);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*RPCMarshalHeaderT, error) {
	header, err := s.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		response := s.rpcMarshalHeader(ctx, header)
		if number == rpc.PendingBlockNumber {
			response.setAsPending()
		}
		return response, err
	}
	return nil, err
}// GetHeaderByNumber returns the requested canonical block header.
//   - When blockNr is -1 the chain pending header is returned.
//   - When blockNr is -2 the chain latest header is returned.
//   - When blockNr is -3 the chain finalized header is returned.
//   - When blockNr is -4 the chain safe header is returned.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L773" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getLogs

GetLogs returns logs matching the given argument that are stored within the state.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
crit <code>FilterCriteria</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Addresses: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash POINTER`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- BlockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- FromBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- ToBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Topics: 
			- items: 
				- items: 
					- description: `Hex representation of a Keccak 256 hash`
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- type: `array`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Addresses": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash POINTER",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "BlockHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "FromBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "ToBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Topics": {
                "items": {
                    "items": {
                        "description": "Hex representation of a Keccak 256 hash",
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "type": "array"
                },
                "type": "array"
            }
        },
        "type": [
            "object"
        ]
    }
	```





#### Result



typesLog <code>[]*types.Log</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- address: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockNumber: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- data: 
					- pattern: `^0x([a-fA-F0-9]?)+$`
					- title: `bytes`
					- type: `string`

				- logIndex: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- removed: 
					- type: `boolean`

				- topics: 
					- items: 
						- description: `Hex representation of a Keccak 256 hash`
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- type: `array`

				- transactionHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- transactionIndex: 
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
                    "address": {
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
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
                    "data": {
                        "pattern": "^0x([a-fA-F0-9]?)+$",
                        "title": "bytes",
                        "type": "string"
                    },
                    "logIndex": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "removed": {
                        "type": "boolean"
                    },
                    "topics": {
                        "items": {
                            "description": "Hex representation of a Keccak 256 hash",
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "type": "array"
                    },
                    "transactionHash": {
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "transactionIndex": {
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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getLogs", "params": [<crit>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getLogs", "params": [<crit>]}'
	```


=== "Javascript Console"

	``` js
	eth.getLogs(crit);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) GetLogs(ctx context.Context, crit FilterCriteria) ([ // GetLogs returns logs matching the given argument that are stored within the state.
]*types.Log, error) {
	var filter *Filter
	if crit.BlockHash != nil {
		filter = api.sys.NewBlockFilter(*crit.BlockHash, crit.Addresses, crit.Topics)
	} else {
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		filter = api.sys.NewRangeFilter(begin, end, crit.Addresses, crit.Topics)
	}
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), err
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L398" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getProof

GetProof returns the Merkle-proof for a given account and optionally some storage keys.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
address <code>common.Address</code> 

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




__2:__ 
storageKeys <code>[]string</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- type: string


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "type": [
                    "string"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```




__3:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>*AccountResult</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accountProof: 
			- items: 
				- type: `string`

			- type: `array`

		- address: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- balance: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- codeHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- storageHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- storageProof: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- key: 
						- type: `string`

					- proof: 
						- items: 
							- type: `string`

						- type: `array`

					- value: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "accountProof": {
                "items": {
                    "type": "string"
                },
                "type": "array"
            },
            "address": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "balance": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "codeHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "storageHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "storageProof": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "key": {
                            "type": "string"
                        },
                        "proof": {
                            "items": {
                                "type": "string"
                            },
                            "type": "array"
                        },
                        "value": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getProof", "params": [<address>, <storageKeys>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getProof", "params": [<address>, <storageKeys>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getProof(address,storageKeys,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys [ // GetProof returns the Merkle-proof for a given account and optionally some storage keys.
]string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	var (
		keys		= make([]common.Hash, len(storageKeys))
		keyLengths	= make([]int, len(storageKeys))
		storageProof	= make([]StorageResult, len(storageKeys))
		storageTrie	state.Trie
		storageHash	= types.EmptyRootHash
		codeHash	= types.EmptyCodeHash
	)
	for i, hexKey := range storageKeys {
		var err error
		keys[i], keyLengths[i], err = decodeHash(hexKey)
		if err != nil {
			return nil, err
		}
	}
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	if storageTrie, err = state.StorageTrie(address); err != nil {
		return nil, err
	}
	if storageTrie != nil {
		storageHash = storageTrie.Hash()
		codeHash = state.GetCodeHash(address)
	}
	for i, key := range keys {
		var outputKey string
		if keyLengths[i] != 32 {
			outputKey = hexutil.EncodeBig(key.Big())
		} else {
			outputKey = hexutil.Encode(key[ // Output key encoding is a bit special: if the input was a 32-byte hash, it is
			// returned as such. Otherwise, we apply the QUANTITY encoding mandated by the
			// JSON-RPC spec for getProof. This behavior exists to preserve backwards
			// compatibility with older client versions.
			:])
		}
		if storageTrie == nil {
			storageProof[i] = StorageResult{outputKey, &hexutil.Big{}, []string{}}
			continue
		}
		var proof proofList
		if err := storageTrie.Prove(crypto.Keccak256(key.Bytes()), &proof); err != nil {
			return nil, err
		}
		value := (*hexutil.Big)(state.GetState(address, key).Big())
		storageProof[i] = StorageResult{outputKey, value, proof}
	}
	accountProof, proofErr := state.GetProof(address)
	if proofErr != nil {
		return nil, proofErr
	}
	return &AccountResult{Address: address, AccountProof: toHexSlice(accountProof), Balance: (*hexutil.Big)(state.GetBalance(address)), CodeHash: codeHash, Nonce: hexutil.Uint64(state.GetNonce(address)), StorageHash: storageHash, StorageProof: storageProof}, state.Error()
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L675" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getRawTransactionByBlockHashAndIndex

GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockHash <code>common.Hash</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getRawTransactionByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getRawTransactionByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getRawTransactionByBlockHashAndIndex(blockHash,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1790" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getRawTransactionByBlockNumberAndIndex

GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getRawTransactionByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getRawTransactionByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getRawTransactionByBlockNumberAndIndex(blockNr,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1782" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getRawTransactionByHash

GetRawTransactionByHash returns the bytes of the transaction for the given hash.


#### Params (1)

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





#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getRawTransactionByHash", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getRawTransactionByHash", "params": [<hash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getRawTransactionByHash(hash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	tx, _, _, _, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			return nil, nil
		}
	}
	return tx.MarshalBinary()
}// GetRawTransactionByHash returns the bytes of the transaction for the given hash.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1840" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getStorageAt

GetStorageAt returns the storage from the state at the given address, key and
block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
numbers are also allowed.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
address <code>common.Address</code> 

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




__2:__ 
hexKey <code>string</code> 

  + Required: ✓ Yes





__3:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getStorageAt", "params": [<address>, <hexKey>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getStorageAt", "params": [<address>, <hexKey>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getStorageAt(address,hexKey,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	key, _, err := decodeHash(hexKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode storage key: %s", err)
	}
	res := state.GetState(address, key)
	return res[ // GetStorageAt returns the storage from the state at the given address, key and
	// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
	// numbers are also allowed.
	:], state.Error()
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L888" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getTransactionByBlockHashAndIndex

GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockHash <code>common.Hash</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>*RPCTransaction</code> 

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
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getTransactionByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getTransactionByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getTransactionByBlockHashAndIndex(blockHash,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index), s.b.ChainConfig())
	}
	return nil
}// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1774" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getTransactionByBlockNumberAndIndex

GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>*RPCTransaction</code> 

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
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getTransactionByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getTransactionByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getTransactionByBlockNumberAndIndex(blockNr,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index), s.b.ChainConfig())
	}
	return nil
}// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1766" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getTransactionByHash

GetTransactionByHash returns the transaction for the given hash


#### Params (1)

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





#### Result




<code>*RPCTransaction</code> 

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
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getTransactionByHash", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getTransactionByHash", "params": [<hash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getTransactionByHash(hash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		header, err := s.b.HeaderByHash(ctx, blockHash)
		if err != nil {
			return nil, err
		}
		return newRPCTransaction(tx, blockHash, blockNumber, header.Time, index, header.BaseFee, s.b.ChainConfig()), nil
	}
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return NewRPCPendingTransaction(tx, s.b.CurrentHeader(), s.b.ChainConfig()), nil
	}
	return nil, nil
}// GetTransactionByHash returns the transaction for the given hash

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1817" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getTransactionCount

GetTransactionCount returns the number of transactions the given address has sent for the given block number


#### Params (2)

Parameters must be given _by position_.


__1:__ 
address <code>common.Address</code> 

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




__2:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>*hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getTransactionCount", "params": [<address>, <blockNrOrHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getTransactionCount", "params": [<address>, <blockNrOrHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getTransactionCount(address,blockNrOrHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		nonce, err := s.b.GetPoolNonce(ctx, address)
		if err != nil {
			return nil, err
		}
		return (*hexutil.Uint64)(&nonce), nil
	}
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	nonce := state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}// GetTransactionCount returns the number of transactions the given address has sent for the given block number

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1798" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getTransactionReceipt

GetTransactionReceipt returns the transaction receipt for the given transaction hash.


#### Params (1)

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





#### Result



mapstringinterface <code>map[string]interface{}</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- additionalProperties: `true`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "patternProperties": {
            ".*": {
                "additionalProperties": true
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getTransactionReceipt", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getTransactionReceipt", "params": [<hash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getTransactionReceipt(hash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map // GetTransactionReceipt returns the transaction receipt for the given transaction hash.
[string]interface{}, error) {
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if tx == nil || err != nil {
		return nil, nil
	}
	header, err := s.b.HeaderByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	receipts, err := s.b.GetReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	if uint64(len(receipts)) <= index {
		return nil, nil
	}
	receipt := receipts[index]
	signer := types.MakeSigner(s.b.ChainConfig(), header.Number, header.Time)
	from, _ := types.Sender(signer, tx)
	fields := map[string]interface{}{"blockHash": blockHash, "blockNumber": hexutil.Uint64(blockNumber), "transactionHash": hash, "transactionIndex": hexutil.Uint64(index), "from": from, "to": tx.To(), "gasUsed": hexutil.Uint64(receipt.GasUsed), "cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed), "contractAddress": nil, "logs": receipt.Logs, "logsBloom": receipt.Bloom, "type": hexutil.Uint(tx.Type()), "effectiveGasPrice": (*hexutil.Big)(receipt.EffectiveGasPrice)}
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = []*types.Log{}
	}
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1857" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getUncleByBlockHashAndIndex

GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockHash <code>common.Hash</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>*RPCMarshalBlockT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- error: 
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- uncles: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- withdrawals: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- amount: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- index: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- validatorIndex: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- type: `array`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "error": {
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "size": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactions": {
                "items": {
                    "additionalProperties": true
                },
                "type": "array"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "uncles": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "withdrawals": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "amount": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "index": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "validatorIndex": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getUncleByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getUncleByBlockHashAndIndex", "params": [<blockHash>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getUncleByBlockHashAndIndex(blockHash,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCMarshalBlockT, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L843" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getUncleByBlockNumberAndIndex

GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

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
index <code>hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>*RPCMarshalBlockT</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- baseFeePerGas: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- error: 
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- uncles: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- withdrawals: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- address: 
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- amount: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- index: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- validatorIndex: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- type: `array`

		- withdrawalsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "baseFeePerGas": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "difficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "error": {
                "type": "string"
            },
            "extraData": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "gasLimit": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "gasUsed": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "hash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "logsBloom": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "maxItems": 256,
                "minItems": 256,
                "type": "array"
            },
            "miner": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "mixHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "nonce": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "number": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "parentHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "receiptsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "sha3Uncles": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "size": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "stateRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "timestamp": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "totalDifficulty": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "transactions": {
                "items": {
                    "additionalProperties": true
                },
                "type": "array"
            },
            "transactionsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "uncles": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "withdrawals": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "address": {
                            "pattern": "^0x[a-fA-F\\d]{64}$",
                            "title": "keccak",
                            "type": "string"
                        },
                        "amount": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "index": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "validatorIndex": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "type": "array"
            },
            "withdrawalsRoot": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getUncleByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getUncleByBlockNumberAndIndex", "params": [<blockNr>, <index>]}'
	```


=== "Javascript Console"

	``` js
	eth.getUncleByBlockNumberAndIndex(blockNr,index);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (*RPCMarshalBlockT, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L827" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getUncleCountByBlockHash

GetUncleCountByBlockHash returns number of uncles in the block for the given block hash


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockHash <code>common.Hash</code> 

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





#### Result




<code>*hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getUncleCountByBlockHash", "params": [<blockHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getUncleCountByBlockHash", "params": [<blockHash>]}'
	```


=== "Javascript Console"

	``` js
	eth.getUncleCountByBlockHash(blockHash);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L867" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getUncleCountByBlockNumber

GetUncleCountByBlockNumber returns number of uncles in the block for the given block number


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

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





#### Result




<code>*hexutil.Uint</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getUncleCountByBlockNumber", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getUncleCountByBlockNumber", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	eth.getUncleCountByBlockNumber(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *BlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L858" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getWork

GetWork returns a work package for external miner.

The work package consists of 3 strings:

	result[0] - 32 bytes hex encoded current block header pow-hash
	result[1] - 32 bytes hex encoded seed hash used for DAG
	result[2] - 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
	result[3] - hex encoded block number


#### Params (0)

_None_

#### Result



num4string <code>[4]string</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- type: string


	- maxItems: `4`
	- minItems: `4`
	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "type": [
                    "string"
                ]
            }
        ],
        "maxItems": 4,
        "minItems": 4,
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_getWork", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_getWork", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.getWork();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) GetWork() ([4]string, error) {
	if api.ethash.remote == nil {
		return [4]string{}, errors.New("not supported")
	}
	var (
		workCh	= make(chan [4]string, 1)
		errc	= make(chan error, 1)
	)
	select {
	case api.ethash.remote.fetchWorkCh <- &sealWork{errc: errc, res: workCh}:
	case <-api.ethash.remote.exitCh:
		return [4]string{}, errEthashStopped
	}
	select {
	case work := <-workCh:
		return work, nil
	case err := <-errc:
		return [4]string{}, err
	}
}// GetWork returns a work package for external miner.
//
// The work package consists of 3 strings:
//
//	result[0] - 32 bytes hex encoded current block header pow-hash
//	result[1] - 32 bytes hex encoded seed hash used for DAG
//	result[2] - 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
//	result[3] - hex encoded block number

```
<a href="https://github.com/etclabscore/core-geth/blob/master/consensus/ethash/api.go#L42" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_hashrate

Hashrate returns the POW hashrate.


#### Params (0)

_None_

#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_hashrate", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_hashrate", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.hashrate();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *EthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().Hashrate())
}// Hashrate returns the POW hashrate.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L45" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_logs

Logs creates a subscription that fires for all new log that match the given filter criteria.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
crit <code>FilterCriteria</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Addresses: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash POINTER`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- BlockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- FromBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- ToBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Topics: 
			- items: 
				- items: 
					- description: `Hex representation of a Keccak 256 hash`
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- type: `array`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Addresses": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash POINTER",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "BlockHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "FromBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "ToBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Topics": {
                "items": {
                    "items": {
                        "description": "Hex representation of a Keccak 256 hash",
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "type": "array"
                },
                "type": "array"
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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": ["logs", <crit>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) Logs(ctx context.Context, crit FilterCriteria) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	var (
		rpcSub		= notifier.CreateSubscription()
		matchedLogs	= make(chan [ // Logs creates a subscription that fires for all new log that match the given filter criteria.
		]*types.Log)
	)
	logsSub, err := api.events.SubscribeLogs(ethereum.FilterQuery(crit), matchedLogs)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case logs := <-matchedLogs:
				for _, log := range logs {
					log := log
					notifier.Notify(rpcSub.ID, &log)
				}
			case <-rpcSub.Err():
				logsSub.Unsubscribe()
				return
			case <-notifier.Closed():
				logsSub.Unsubscribe()
				return
			}
		}
	}()
	return rpcSub, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L313" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_maxPriorityFeePerGas

MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.


#### Params (0)

_None_

#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_maxPriorityFeePerGas", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_maxPriorityFeePerGas", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.maxPriorityFeePerGas();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *EthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(tipcap), err
}// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L78" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_mining

Mining returns an indication if this node is currently mining.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_mining", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_mining", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.mining();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *EthereumAPI) Mining() bool {
	return api.e.IsMining()
}// Mining returns an indication if this node is currently mining.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L51" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newBlockFilter

NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
It is part of the filter package since polling goes with eth_getFilterChanges.


#### Params (0)

_None_

#### Result




<code>rpc.ID</code> 

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_newBlockFilter", "params": []}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewBlockFilter() rpc.ID {
	var (
		headers		= make(chan *types.Header)
		headerSub	= api.events.SubscribeNewHeads(headers)
	)
	api.filtersMu.Lock()
	api.filters[headerSub.ID] = &filter{typ: BlocksSubscription, deadline: time.NewTimer(api.timeout), hashes: make([ // NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
	// It is part of the filter package since polling goes with eth_getFilterChanges.
	]common.Hash, 0), s: headerSub}
	api.filtersMu.Unlock()
	go func() {
		for {
			select {
			case h := <-headers:
				api.filtersMu.Lock()
				if f, found := api.filters[headerSub.ID]; found {
					f.hashes = append(f.hashes, h.Hash())
				}
				api.filtersMu.Unlock()
			case <-headerSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, headerSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()
	return headerSub.ID
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L188" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newFilter

NewFilter creates a new filter and returns the filter id. It can be
used to retrieve logs when the state changes. This method cannot be
used to fetch logs that are already stored in the state.

Default criteria for the from and to block are "latest".
Using "latest" as block number will return logs for mined blocks.
Using "pending" as block number returns logs for not yet mined (pending) blocks.
In case logs are removed (chain reorg) previously returned logs are returned
again but with the removed property set to true.

In case "fromBlock" > "toBlock" an error is returned.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
crit <code>FilterCriteria</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Addresses: 
			- items: 
				- description: `Hex representation of a Keccak 256 hash POINTER`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`
				- type: `string`

			- type: `array`

		- BlockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- FromBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- ToBlock: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Topics: 
			- items: 
				- items: 
					- description: `Hex representation of a Keccak 256 hash`
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- type: `array`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "Addresses": {
                "items": {
                    "description": "Hex representation of a Keccak 256 hash POINTER",
                    "pattern": "^0x[a-fA-F\\d]{64}$",
                    "title": "keccak",
                    "type": "string"
                },
                "type": "array"
            },
            "BlockHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "FromBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "ToBlock": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Topics": {
                "items": {
                    "items": {
                        "description": "Hex representation of a Keccak 256 hash",
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "type": "array"
                },
                "type": "array"
            }
        },
        "type": [
            "object"
        ]
    }
	```





#### Result




<code>rpc.ID</code> 

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_newFilter", "params": [<crit>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewFilter(crit FilterCriteria) (rpc.ID, error) {
	logs := make(chan [ // NewFilter creates a new filter and returns the filter id. It can be
	// used to retrieve logs when the state changes. This method cannot be
	// used to fetch logs that are already stored in the state.
	//
	// Default criteria for the from and to block are "latest".
	// Using "latest" as block number will return logs for mined blocks.
	// Using "pending" as block number returns logs for not yet mined (pending) blocks.
	// In case logs are removed (chain reorg) previously returned logs are returned
	// again but with the removed property set to true.
	//
	// In case "fromBlock" > "toBlock" an error is returned.
	]*types.Log)
	logsSub, err := api.events.SubscribeLogs(ethereum.FilterQuery(crit), logs)
	if err != nil {
		return "", err
	}
	api.filtersMu.Lock()
	api.filters[logsSub.ID] = &filter{typ: LogsSubscription, crit: crit, deadline: time.NewTimer(api.timeout), logs: make([]*types.Log, 0), s: logsSub}
	api.filtersMu.Unlock()
	go func() {
		for {
			select {
			case l := <-logs:
				api.filtersMu.Lock()
				if f, found := api.filters[logsSub.ID]; found {
					f.logs = append(f.logs, l...)
				}
				api.filtersMu.Unlock()
			case <-logsSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, logsSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()
	return logsSub.ID, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L365" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newHeads

NewHeads send a notification each time a new (header) block is appended to the chain.


#### Params (0)

_None_

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": ["newHeads"]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		headers := make(chan *types.Header)
		headersSub := api.events.SubscribeNewHeads(headers)
		for {
			select {
			case h := <-headers:
				notifier.Notify(rpcSub.ID, h)
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()
	return rpcSub, nil
}// NewHeads send a notification each time a new (header) block is appended to the chain.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L253" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newPendingTransactionFilter

NewPendingTransactionFilter creates a filter that fetches pending transactions
as transactions enter the pending state.

It is part of the filter package because this filter can be used through the
`eth_getFilterChanges` polling method that is also used for log filters.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
fullTx <code>*bool</code> 

  + Required: ✓ Yes






#### Result




<code>rpc.ID</code> 

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_newPendingTransactionFilter", "params": [<fullTx>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewPendingTransactionFilter(fullTx *bool) rpc.ID {
	var (
		pendingTxs	= make(chan [ // NewPendingTransactionFilter creates a filter that fetches pending transactions
		// as transactions enter the pending state.
		//
		// It is part of the filter package because this filter can be used through the
		// `eth_getFilterChanges` polling method that is also used for log filters.
		]*types.Transaction)
		pendingTxSub	= api.events.SubscribePendingTxs(pendingTxs)
	)
	api.filtersMu.Lock()
	api.filters[pendingTxSub.ID] = &filter{typ: PendingTransactionsSubscription, fullTx: fullTx != nil && *fullTx, deadline: time.NewTimer(api.timeout), txs: make([]*types.Transaction, 0), s: pendingTxSub}
	api.filtersMu.Unlock()
	go func() {
		for {
			select {
			case pTx := <-pendingTxs:
				api.filtersMu.Lock()
				if f, found := api.filters[pendingTxSub.ID]; found {
					f.txs = append(f.txs, pTx...)
				}
				api.filtersMu.Unlock()
			case <-pendingTxSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, pendingTxSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()
	return pendingTxSub.ID
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L112" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newPendingTransactions

NewPendingTransactions creates a subscription that is triggered each time a
transaction enters the transaction pool. If fullTx is true the full tx is
sent to the client, otherwise the hash is sent.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
fullTx <code>*bool</code> 

  + Required: ✓ Yes






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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": ["newPendingTransactions", <fullTx>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		txs := make(chan [ // NewPendingTransactions creates a subscription that is triggered each time a
		// transaction enters the transaction pool. If fullTx is true the full tx is
		// sent to the client, otherwise the hash is sent.
		]*types.Transaction, 128)
		pendingTxSub := api.events.SubscribePendingTxs(txs)
		chainConfig := api.sys.backend.ChainConfig()
		for {
			select {
			case txs := <-txs:
				latest := api.sys.backend.CurrentHeader()
				for _, tx := range txs {
					if fullTx != nil && *fullTx {
						rpcTx := ethapi.NewRPCPendingTransaction(tx, latest, chainConfig)
						notifier.Notify(rpcSub.ID, rpcTx)
					} else {
						notifier.Notify(rpcSub.ID, tx.Hash())
					}
				}
			case <-rpcSub.Err():
				pendingTxSub.Unsubscribe()
				return
			case <-notifier.Closed():
				pendingTxSub.Unsubscribe()
				return
			}
		}
	}()
	return rpcSub, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L146" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newSideBlockFilter

NewSideBlockFilter creates a filter that fetches blocks that are imported into the chain with a non-canonical status.
It is part of the filter package since polling goes with eth_getFilterChanges.


#### Params (0)

_None_

#### Result




<code>rpc.ID</code> 

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_newSideBlockFilter", "params": []}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewSideBlockFilter() rpc.ID {
	var (
		headers		= make(chan *types.Header)
		headerSub	= api.events.SubscribeNewSideHeads(headers)
	)
	api.filtersMu.Lock()
	api.filters[headerSub.ID] = &filter{typ: SideBlocksSubscription, deadline: time.NewTimer(api.timeout), hashes: make([ // NewSideBlockFilter creates a filter that fetches blocks that are imported into the chain with a non-canonical status.
	// It is part of the filter package since polling goes with eth_getFilterChanges.
	]common.Hash, 0), s: headerSub}
	api.filtersMu.Unlock()
	go func() {
		for {
			select {
			case h := <-headers:
				api.filtersMu.Lock()
				if f, found := api.filters[headerSub.ID]; found {
					f.hashes = append(f.hashes, h.Hash())
				}
				api.filtersMu.Unlock()
			case <-headerSub.Err():
				api.filtersMu.Lock()
				delete(api.filters, headerSub.ID)
				api.filtersMu.Unlock()
				return
			}
		}
	}()
	return headerSub.ID
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L221" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newSideHeads

NewSideHeads send a notification each time a new non-canonical (header) block is written to the database.


#### Params (0)

_None_

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": ["newSideHeads"]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) NewSideHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		headers := make(chan *types.Header)
		headersSub := api.events.SubscribeNewSideHeads(headers)
		for {
			select {
			case h := <-headers:
				notifier.Notify(rpcSub.ID, h)
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()
	return rpcSub, nil
}// NewSideHeads send a notification each time a new non-canonical (header) block is written to the database.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L283" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_pendingTransactions

PendingTransactions returns the transactions that are in the transaction pool
and have a from address that is one of the accounts this node manages.


#### Params (0)

_None_

#### Result



RPCTransaction <code>[]*RPCTransaction</code> 

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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_pendingTransactions", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_pendingTransactions", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.pendingTransactions();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) PendingTransactions() ([ // PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
]*RPCTransaction, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := make(map[common.Address]struct{})
	for _, wallet := range s.b.AccountManager().Wallets() {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = struct{}{}
		}
	}
	curHeader := s.b.CurrentHeader()
	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from, _ := types.Sender(s.signer, tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig()))
		}
	}
	return transactions, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2085" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_resend

Resend accepts an existing transaction and a new gas price and limit. It will remove
the given transaction from the pool and reinsert it with the new gas price and limit.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
sendArgs <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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
gasPrice <code>*hexutil.Big</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
gasLimit <code>*hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>common.Hash</code> 

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_resend", "params": [<sendArgs>, <gasPrice>, <gasLimit>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_resend", "params": [<sendArgs>, <gasPrice>, <gasLimit>]}'
	```


=== "Javascript Console"

	``` js
	eth.resend(sendArgs,gasPrice,gasLimit);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) Resend(ctx context.Context, sendArgs TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.Nonce == nil {
		return common.Hash{}, errors.New("missing transaction nonce in transaction spec")
	}
	if err := sendArgs.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx := sendArgs.toTransaction()
	var price = matchTx.GasPrice()
	if gasPrice != nil {
		price = gasPrice.ToInt()
	}
	var gas = matchTx.Gas()
	if gasLimit != nil {
		gas = uint64(*gasLimit)
	}
	if err := checkTxFee(price, gas, s.b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}
	for _, p := // Resend accepts an existing transaction and a new gas price and limit. It will remove
	// the given transaction from the pool and reinsert it with the new gas price and limit.
	// Before replacing the old transaction, ensure the _new_ transaction fee is reasonable.
	range pending {
		wantSigHash := s.signer.Hash(matchTx)
		pFrom, err := types.Sender(s.signer, p)
		if err == nil && pFrom == sendArgs.from() && s.signer.Hash(p) == wantSigHash {
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.Gas = gasLimit
			}
			signedTx, err := s.sign(sendArgs.from(), sendArgs.toTransaction())
			if err != nil {
				return common.Hash{}, err
			}
			if err = s.b.SendTx(ctx, signedTx); err != nil {
				return common.Hash{}, err
			}
			return signedTx.Hash(), nil
		}
	}
	return common.Hash{}, fmt.Errorf("transaction %#x not found", matchTx.Hash())
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2109" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sendRawTransaction

SendRawTransaction will add the signed transaction to the transaction pool.
The sender is responsible for signing the transaction and using the correct nonce.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
input <code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>common.Hash</code> 

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_sendRawTransaction", "params": [<input>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_sendRawTransaction", "params": [<input>]}'
	```


=== "Javascript Console"

	``` js
	eth.sendRawTransaction(input);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(input); err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, tx)
}// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2012" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sendTransaction

SendTransaction creates a transaction for the given argument, sign it and submit it to the
transaction pool.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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





#### Result




<code>common.Hash</code> 

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_sendTransaction", "params": [<args>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_sendTransaction", "params": [<args>]}'
	```


=== "Javascript Console"

	``` js
	eth.sendTransaction(args);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) SendTransaction(ctx context.Context, args TransactionArgs) (common.Hash, error) {
	account := accounts.Account{Address: args.from()}
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}
	if args.Nonce == nil {
		s.nonceLock.LockAddr(args.from())
		defer s.nonceLock.UnlockAddr(args.from())
	}
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	tx := args.toTransaction()
	signed, err := wallet.SignTx(account, tx, s.b.ChainConfig().GetChainID())
	if err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, signed)
}// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1963" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sign

Sign calculates an ECDSA signature for:
keccak256("\x19Ethereum Signed Message:\n" + len(message) + message).

Note, the produced signature conforms to the secp256k1 curve R, S and V values,
where the V value will be 27 or 28 for legacy reasons.

The account associated with addr must be unlocked.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign


#### Params (2)

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




__2:__ 
data <code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `dataWord`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of some bytes",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "dataWord",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_sign", "params": [<addr>, <data>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_sign", "params": [<addr>, <data>]}'
	```


=== "Javascript Console"

	``` js
	eth.sign(addr,data);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	account := accounts.Account{Address: addr}
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	signature, err := wallet.SignText(account, data)
	if err == nil {
		signature[64] += 27
	}
	return signature, err
}// Sign calculates an ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2029" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_signTransaction

SignTransaction will sign the given transaction with the from account.
The node needs to have the private key of the account corresponding with
the given from address and it needs to be unlocked.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
args <code>TransactionArgs</code> 

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

		- chainId: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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
            "chainId": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
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





#### Result




<code>*SignTransactionResult</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- raw: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- tx: 
			- additionalProperties: `false`
			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "raw": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "dataWord",
                "type": "string"
            },
            "tx": {
                "additionalProperties": false,
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_signTransaction", "params": [<args>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_signTransaction", "params": [<args>]}'
	```


=== "Javascript Console"

	``` js
	eth.signTransaction(args);
	```



<details><summary>Source code</summary>
<p>
```go
func (s *TransactionAPI) SignTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	if args.Gas == nil {
		return nil, errors.New("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
		return nil, errors.New("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, errors.New("nonce not specified")
	}
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := s.sign(args.from(), tx)
	if err != nil {
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2054" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_submitHashrate

SubmitHashrate can be used for remote miners to submit their hash rate.
This enables the node to report the combined hash rate of all miners
which submit work through this node.

It accepts the miner hash rate and an identifier which must be unique
between nodes.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
rate <code>hexutil.Uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
	- title: `uint64`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a uint64",
        "pattern": "^0x([a-fA-F\\d])+$",
        "title": "uint64",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
id <code>common.Hash</code> 

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





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_submitHashrate", "params": [<rate>, <id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_submitHashrate", "params": [<rate>, <id>]}'
	```


=== "Javascript Console"

	``` js
	eth.submitHashrate(rate,id);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) SubmitHashrate(rate hexutil.Uint64, id common.Hash) bool {
	if api.ethash.remote == nil {
		return false
	}
	var done = make(chan struct{}, 1)
	select {
	case api.ethash.remote.submitRateCh <- &hashrate{done: done, rate: uint64(rate), id: id}:
	case <-api.ethash.remote.exitCh:
		return false
	}
	<-done
	return true
}// SubmitHashrate can be used for remote miners to submit their hash rate.
// This enables the node to report the combined hash rate of all miners
// which submit work through this node.
//
// It accepts the miner hash rate and an identifier which must be unique
// between nodes.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/consensus/ethash/api.go#L93" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_submitWork

SubmitWork can be used by external miner to submit their POW solution.
It returns an indication if the work was accepted.
Note either an invalid solution, a stale work a non-existent work will return false.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
nonce <code>types.BlockNonce</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__2:__ 
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




__3:__ 
digest <code>common.Hash</code> 

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





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_submitWork", "params": [<nonce>, <hash>, <digest>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_submitWork", "params": [<nonce>, <hash>, <digest>]}'
	```


=== "Javascript Console"

	``` js
	eth.submitWork(nonce,hash,digest);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) SubmitWork(nonce types.BlockNonce, hash, digest common.Hash) bool {
	if api.ethash.remote == nil {
		return false
	}
	var errc = make(chan error, 1)
	select {
	case api.ethash.remote.submitWorkCh <- &mineResult{nonce: nonce, mixDigest: digest, hash: hash, errc: errc}:
	case <-api.ethash.remote.exitCh:
		return false
	}
	err := <-errc
	return err == nil
}// SubmitWork can be used by external miner to submit their POW solution.
// It returns an indication if the work was accepted.
// Note either an invalid solution, a stale work a non-existent work will return false.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/consensus/ethash/api.go#L67" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_subscribe

Subscribe creates a subscription to an event channel.
Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
subscriptionName <code>RPCEthSubscriptionParamsName</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `Fires a notification each time a new header is appended to the chain, including chain reorganizations.`
			- enum: newHeads
			- type: string


			- description: `Fires a notification each time a new header is appended to the non-canonical (side) chain, including chain reorganizations.`
			- enum: newSideHeads
			- type: string


			- description: `Returns logs that are included in new imported blocks and match the given filter criteria.`
			- enum: logs
			- type: string


			- description: `Returns the hash for all transactions that are added to the pending state and are signed with a key that is available in the node.`
			- enum: newPendingTransactions
			- type: string


			- description: `Indicates when the node starts or stops synchronizing. The result can either be a boolean indicating that the synchronization has started (true), finished (false) or an object with various progress indicators.`
			- enum: syncing
			- type: string


	- title: `subscriptionName`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "Fires a notification each time a new header is appended to the chain, including chain reorganizations.",
                "enum": [
                    "newHeads"
                ],
                "type": [
                    "string"
                ]
            },
            {
                "description": "Fires a notification each time a new header is appended to the non-canonical (side) chain, including chain reorganizations.",
                "enum": [
                    "newSideHeads"
                ],
                "type": [
                    "string"
                ]
            },
            {
                "description": "Returns logs that are included in new imported blocks and match the given filter criteria.",
                "enum": [
                    "logs"
                ],
                "type": [
                    "string"
                ]
            },
            {
                "description": "Returns the hash for all transactions that are added to the pending state and are signed with a key that is available in the node.",
                "enum": [
                    "newPendingTransactions"
                ],
                "type": [
                    "string"
                ]
            },
            {
                "description": "Indicates when the node starts or stops synchronizing. The result can either be a boolean indicating that the synchronization has started (true), finished (false) or an object with various progress indicators.",
                "enum": [
                    "syncing"
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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": [<subscriptionName>, <subscriptionOptions>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (sub *RPCEthSubscription) Subscribe(subscriptionName RPCEthSubscriptionParamsName, subscriptionOptions interface{}) (subscriptionID rpc.ID, err error) {
	return
}// Subscribe creates a subscription to an event channel.
// Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L233" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_syncing

Syncing provides information when this nodes starts synchronising with the Ethereum network and when it's finished.


#### Params (0)

_None_

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_subscribe", "params": ["syncing"]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *DownloaderAPI) Syncing(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		statuses := make(chan interface{})
		sub := api.SubscribeSyncStatus(statuses)
		for {
			select {
			case status := <-statuses:
				notifier.Notify(rpcSub.ID, status)
			case <-rpcSub.Err():
				sub.Unsubscribe()
				return
			case <-notifier.Closed():
				sub.Unsubscribe()
				return
			}
		}
	}()
	return rpcSub, nil
}// Syncing provides information when this nodes starts synchronising with the Ethereum network and when it's finished.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/downloader/api.go#L93" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_syncing

Syncing returns false in case the node is currently not syncing with the network. It can be up-to-date or has not
yet received the latest block headers from its pears. In case it is synchronizing:
- startingBlock: block number this node started to synchronize from
- currentBlock:  block number this node is currently importing
- highestBlock:  block number of the highest block header this node has received from peers
- pulledStates:  number of state entries processed until now
- knownStates:   number of known state entries that still need to be pulled


#### Params (0)

_None_

#### Result



interface <code>interface{}</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_syncing", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_syncing", "params": []}'
	```


=== "Javascript Console"

	``` js
	eth.syncing();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *EthereumAPI) Syncing() (interface{}, error) {
	progress := s.b.SyncProgress()
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	return map // Syncing returns false in case the node is currently not syncing with the network. It can be up-to-date or has not
	// yet received the latest block headers from its pears. In case it is synchronizing:
	// - startingBlock: block number this node started to synchronize from
	// - currentBlock:  block number this node is currently importing
	// - highestBlock:  block number of the highest block header this node has received from peers
	// - pulledStates:  number of state entries processed until now
	// - knownStates:   number of known state entries that still need to be pulled
	[string]interface{}{"startingBlock": hexutil.Uint64(progress.StartingBlock), "currentBlock": hexutil.Uint64(progress.CurrentBlock), "highestBlock": hexutil.Uint64(progress.HighestBlock), "syncedAccounts": hexutil.Uint64(progress.SyncedAccounts), "syncedAccountBytes": hexutil.Uint64(progress.SyncedAccountBytes), "syncedBytecodes": hexutil.Uint64(progress.SyncedBytecodes), "syncedBytecodeBytes": hexutil.Uint64(progress.SyncedBytecodeBytes), "syncedStorage": hexutil.Uint64(progress.SyncedStorage), "syncedStorageBytes": hexutil.Uint64(progress.SyncedStorageBytes), "healedTrienodes": hexutil.Uint64(progress.HealedTrienodes), "healedTrienodeBytes": hexutil.Uint64(progress.HealedTrienodeBytes), "healedBytecodes": hexutil.Uint64(progress.HealedBytecodes), "healedBytecodeBytes": hexutil.Uint64(progress.HealedBytecodeBytes), "healingTrienodes": hexutil.Uint64(progress.HealingTrienodes), "healingBytecode": hexutil.Uint64(progress.HealingBytecode)}, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L128" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_uninstallFilter

UninstallFilter removes the filter with the given filter id.


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




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_uninstallFilter", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_uninstallFilter", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	eth.uninstallFilter(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *FilterAPI) UninstallFilter(id rpc.ID) bool {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	if found {
		delete(api.filters, id)
	}
	api.filtersMu.Unlock()
	if found {
		f.s.Unsubscribe()
	}
	return found
}// UninstallFilter removes the filter with the given filter id.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/filters/api.go#L425" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_unsubscribe

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "eth_unsubscribe", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "eth_unsubscribe", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	eth.unsubscribe(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (sub *RPCEthSubscription) Unsubscribe(id rpc.ID) error {
	return nil
}// Unsubscribe terminates an existing subscription by ID.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L224" target="_">View on GitHub →</a>
</p>
</details>

---

