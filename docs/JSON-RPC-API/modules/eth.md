






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-21T17:27:32-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### eth_accounts

Accounts returns the collection of accounts this node manages


#### Params (0)

_None_

#### Result



commonAddress <code>[]common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: array
	- items: 

			- title: `keccak`
			- description: `Hex representation of a Keccak 256 hash POINTER`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- type: string




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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_accounts", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.accounts();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicAccountAPI) Accounts() [ // Accounts returns the collection of accounts this node manages
]common.Address {
	return s.am.Accounts()
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L190" target="_">View on GitHub →</a>
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
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_blockNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.blockNumber();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	return hexutil.Uint64(header.Number.Uint64())
}// BlockNumber returns the block number of the chain head.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L545" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_call

Call executes the given transaction on the state for the given block number.

Additionally, the caller can specify a batch of contract for fields overriding.

Note, this function doesn't make and changes in the state/blockchain and is
useful to execute and retrieve values.


#### Params (3)

Parameters must be given _by position_.  


__1:__ 
args <code>CallArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- value: 
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

		- to: 
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
overrides <code>*map[common.Address]account</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- patternProperties: 
		- .*: 
			- properties: 
				- nonce: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `uint64`
					- type: `string`

				- state: 
					- patternProperties: 
						- .*: 
							- title: `keccak`
							- type: `string`
							- description: `Hex representation of a Keccak 256 hash`
							- pattern: `^0x[a-fA-F\d]{64}$`


					- type: `object`

				- stateDiff: 
					- patternProperties: 
						- .*: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`
							- description: `Hex representation of a Keccak 256 hash`


					- type: `object`

				- balance: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- code: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `dataWord`
					- type: `string`


			- type: `object`
			- additionalProperties: `false`


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




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: string
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_call", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.call(args,blockNrOrHash,overrides);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) Call(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *map // Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
[common.Address]account) (hexutil.Bytes, error) {
	var accounts map[common.Address]account
	if overrides != nil {
		accounts = *overrides
	}
	result, err := DoCall(ctx, s.b, args, blockNrOrHash, accounts, vm.Config{}, 5*time.Second, s.b.RPCGasCap())
	if err != nil {
		return nil, err
	}
	if len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}
	return result.Return(), result.Err
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L928" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_chainId

ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.


#### Params (0)

_None_

#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_chainId", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.chainId();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicEthereumAPI) ChainId() hexutil.Uint64 {
	chainID := new(big.Int)
	if config := api.e.blockchain.Config(); config.IsEnabled(config.GetEIP155Transition, api.e.blockchain.CurrentBlock().Number()) {
		chainID = config.GetChainID()
	}
	return (hexutil.Uint64)(chainID.Uint64())
}// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L70" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_chainId

ChainId returns the chainID value for transaction replay protection.


#### Params (0)

_None_

#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_chainId", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.chainId();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) ChainId() *hexutil.Big {
	return (*hexutil.Big)(s.b.ChainConfig().GetChainID())
}// ChainId returns the chainID value for transaction replay protection.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L540" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_coinbase

Coinbase is the address that mining rewards will be send to (alias for Etherbase)


#### Params (0)

_None_

#### Result




<code>common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_coinbase", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.coinbase();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicEthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}// Coinbase is the address that mining rewards will be send to (alias for Etherbase)

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L60" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_estimateGas

EstimateGas returns an estimate of the amount of gas needed to execute the
given transaction against the current pending block.


#### Params (2)

Parameters must be given _by position_.  


__1:__ 
args <code>CallArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- value: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- data: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- from: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- gas: 
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
blockNrOrHash <code>*rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_estimateGas", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.estimateGas(args,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) EstimateGas(ctx context.Context, args CallArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, s.b.RPCGasCap())
}// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1055" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_etherbase

Etherbase is the address that mining rewards will be send to


#### Params (0)

_None_

#### Result




<code>common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_etherbase", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.etherbase();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicEthereumAPI) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}// Etherbase is the address that mining rewards will be send to

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L55" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_fillTransaction

FillTransaction fills the defaults (nonce, gas, gasPrice) on a given unsigned transaction,
and returns it to the caller for further processing (signing + broadcast)


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
args <code>SendTxArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: object
	- additionalProperties: `false`
	- properties: 
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

		- data: 
			- title: `dataWord`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- from: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- input: 
			- title: `dataWord`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`




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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_fillTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.fillTransaction(args);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) FillTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx := args.toTransaction()
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}// FillTransaction fills the defaults (nonce, gas, gasPrice) on a given unsigned transaction,
// and returns it to the caller for further processing (signing + broadcast)

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1725" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_gasPrice

GasPrice returns a suggestion for a gas price.


#### Params (0)

_None_

#### Result




<code>*hexutil.Big</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_gasPrice", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.gasPrice();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicEthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	price, err := s.b.SuggestPrice(ctx)
	return (*hexutil.Big)(price), err
}// GasPrice returns a suggestion for a gas price.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L63" target="_">View on GitHub →</a>
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
	
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`


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
	
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getBalance", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getBalance(address,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L553" target="_">View on GitHub →</a>
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
	
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`


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
		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
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

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- receiptsRoot: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- sha3Uncles: 
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

		- error: 
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

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

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- stateRoot: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
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
            }
        },
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getBlockByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getBlockByHash(hash,fullTx);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (*RPCMarshalBlockT, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		return s.rpcMarshalBlock(ctx, block, true, fullTx)
	}
	return nil, err
}// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L672" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getBlockByNumber

GetBlockByNumber returns the requested canonical block.
* When blockNr is -1 the chain head is returned.
* When blockNr is -2 the pending chain head is returned.
* When fullTx is true all transactions in the block are returned, otherwise
  only the transaction hash is returned.


#### Params (2)

Parameters must be given _by position_.  


__1:__ 
number <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- enum: earliest, latest, pending
			- type: string
			- title: `blockNumberTag`
			- description: `The block height description`


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- type: string
			- title: `uint64`




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
		- miner: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- parentHash: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- timestamp: 
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- uncles: 
			- items: 
				- title: `keccak`
				- type: `string`
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`

			- type: `array`

		- error: 
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
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

		- gasLimit: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- logsBloom: 
			- items: 
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`
				- description: `Hex representation of the integer`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- totalDifficulty: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
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
            }
        },
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getBlockByNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getBlockByNumber(number,fullTx);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (*RPCMarshalBlockT, error) {
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
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
// * When fullTx is true all transactions in the block are returned, otherwise
//   only the transaction hash is returned.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L657" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getBlockTransactionCountByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getBlockTransactionCountByHash(blockHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1421" target="_">View on GitHub →</a>
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
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- title: `blockNumberTag`
			- description: `The block height description`
			- enum: earliest, latest, pending
			- type: string


			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- type: string




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
	- type: string
	- title: `uint`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getBlockTransactionCountByNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getBlockTransactionCountByNumber(blockNr);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1412" target="_">View on GitHub →</a>
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
	
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`


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
	
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getCode", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getCode(address,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}// GetCode returns the code stored at the given address in the state for the given block number.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L731" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getFilterChanges

GetFilterChanges returns the logs for the filter with the given id since
last time it was called. This can be used for polling.

For pending transaction and block filters the result is []common.Hash.
(pending)Log filters return []Log.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
id <code>rpc.ID</code> 

  + Required: ✓ Yes






#### Result



interface <code>interface{}</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getFilterChanges", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getFilterChanges(id);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()
	if f, found := api.filters[id]; found {
		if !f.deadline.Stop() {
			<-f.deadline.C
		}
		f.deadline.Reset(deadline)
		switch f.typ {
		case PendingTransactionsSubscription, BlocksSubscription, SideBlocksSubscription:
			hashes := f.hashes
			f.hashes = nil
			return returnHashes(hashes), nil
		case LogsSubscription, MinedAndPendingLogsSubscription:
			logs := f.logs
			f.logs = nil
			return returnLogs(logs), nil
		}
	}
	return [ // GetFilterChanges returns the logs for the filter with the given id since
	// last time it was called. This can be used for polling.
	//
	// For pending transaction and block filters the result is []common.Hash.
	// (pending)Log filters return []Log.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges
	]interface{}{}, fmt.Errorf("filter not found")
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L477" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getFilterLogs

GetFilterLogs returns the logs for the filter with the given id.
If the filter could not be found an empty array of logs is returned.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
id <code>rpc.ID</code> 

  + Required: ✓ Yes






#### Result



typesLog <code>[]*types.Log</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- data: 
					- pattern: `^0x([a-fA-F0-9]?)+$`
					- title: `bytes`
					- type: `string`

				- transactionHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockNumber: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- blockHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- logIndex: 
					- title: `integer`
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`

				- removed: 
					- type: `boolean`

				- topics: 
					- items: 
						- description: `Hex representation of a Keccak 256 hash`
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- type: `array`

				- transactionIndex: 
					- title: `integer`
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`

				- address: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getFilterLogs", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getFilterLogs(id);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([ // GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs
]*types.Log, error) {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()
	if !found || f.typ != LogsSubscription {
		return nil, fmt.Errorf("filter not found")
	}
	var filter *Filter
	if f.crit.BlockHash != nil {
		filter = NewBlockFilter(api.backend, *f.crit.BlockHash, f.crit.Addresses, f.crit.Topics)
	} else {
		begin := rpc.LatestBlockNumber.Int64()
		if f.crit.FromBlock != nil {
			begin = f.crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if f.crit.ToBlock != nil {
			end = f.crit.ToBlock.Int64()
		}
		filter = NewRangeFilter(api.backend, begin, end, f.crit.Addresses, f.crit.Topics)
	}
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L436" target="_">View on GitHub →</a>
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
	
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getHashrate", "params": []}'
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/api.go#L110" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- extraData: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- gasLimit: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- miner: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- timestamp: 
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- logsBloom: 
			- type: `array`
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- sha3Uncles: 
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
            "transactionsRoot": {
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getHeaderByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getHeaderByHash(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) *RPCMarshalHeaderT {
	header, _ := s.b.HeaderByHash(ctx, hash)
	if header != nil {
		return s.rpcMarshalHeader(ctx, header)
	}
	return nil
}// GetHeaderByHash returns the requested header by hash.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L644" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getHeaderByNumber

GetHeaderByNumber returns the requested canonical block header.
* When blockNr is -1 the chain head is returned.
* When blockNr is -2 the pending chain head is returned.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
number <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- enum: earliest, latest, pending
			- type: string
			- title: `blockNumberTag`
			- description: `The block height description`


			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- type: string




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
	
	- properties: 
		- gasLimit: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- extraData: 
			- title: `dataWord`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- receiptsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- items: 
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`
				- description: `Hex representation of the integer`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- nonce: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- stateRoot: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- totalDifficulty: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- difficulty: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- timestamp: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`


	- type: object
	- additionalProperties: `false`


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
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
            "transactionsRoot": {
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getHeaderByNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getHeaderByNumber(number);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*RPCMarshalHeaderT, error) {
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
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L630" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_getLogs

GetLogs returns logs matching the given argument that are stored within the state.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
crit <code>FilterCriteria</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- properties: 
		- Addresses: 
			- items: 
				- type: `string`
				- description: `Hex representation of a Keccak 256 hash POINTER`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`

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
	- additionalProperties: `false`


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
				- topics: 
					- items: 
						- description: `Hex representation of a Keccak 256 hash`
						- pattern: `^0x[a-fA-F\d]{64}$`
						- title: `keccak`
						- type: `string`

					- type: `array`

				- blockHash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- blockNumber: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- logIndex: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- removed: 
					- type: `boolean`

				- transactionHash: 
					- type: `string`
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`

				- transactionIndex: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- address: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- data: 
					- pattern: `^0x([a-fA-F0-9]?)+$`
					- title: `bytes`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getLogs", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getLogs(crit);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) GetLogs(ctx context.Context, crit FilterCriteria) ([ // GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
]*types.Log, error) {
	var filter *Filter
	if crit.BlockHash != nil {
		filter = NewBlockFilter(api.backend, *crit.BlockHash, crit.Addresses, crit.Topics)
	} else {
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		filter = NewRangeFilter(api.backend, begin, end, crit.Addresses, crit.Topics)
	}
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), err
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L389" target="_">View on GitHub →</a>
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
	
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`


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
	
	- properties: 
		- accountProof: 
			- type: `array`
			- items: 
				- type: `string`


		- address: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- balance: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

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
				- type: `object`
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



			- type: `array`


	- type: object
	- additionalProperties: `false`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getProof", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getProof(address,storageKeys,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys [ // GetProof returns the Merkle-proof for a given account and optionally some storage keys.
]string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	storageTrie := state.StorageTrie(address)
	storageHash := types.EmptyRootHash
	codeHash := state.GetCodeHash(address)
	storageProof := make([]StorageResult, len(storageKeys))
	if storageTrie != nil {
		storageHash = storageTrie.Hash()
	} else {
		codeHash = crypto.Keccak256Hash(nil)
	}
	for i, key := range storageKeys {
		if storageTrie != nil {
			proof, storageError := state.GetStorageProof(address, common.HexToHash(key))
			if storageError != nil {
				return nil, storageError
			}
			storageProof[i] = StorageResult{key, (*hexutil.Big)(state.GetState(address, common.HexToHash(key)).Big()), toHexSlice(proof)}
		} else {
			storageProof[i] = StorageResult{key, &hexutil.Big{}, []string{}}
		}
	}
	accountProof, proofErr := state.GetProof(address)
	if proofErr != nil {
		return nil, proofErr
	}
	return &AccountResult{Address: address, AccountProof: toHexSlice(accountProof), Balance: (*hexutil.Big)(state.GetBalance(address)), CodeHash: codeHash, Nonce: hexutil.Uint64(state.GetNonce(address)), StorageHash: storageHash, StorageProof: storageProof}, state.Error()
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L578" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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
	
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getRawTransactionByBlockHashAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getRawTransactionByBlockHashAndIndex(blockHash,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1454" target="_">View on GitHub →</a>
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
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- type: string
			- title: `blockNumberTag`
			- description: `The block height description`
			- enum: earliest, latest, pending


			- type: string
			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`




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
	- type: string
	- title: `uint`


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
	
	- type: string
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getRawTransactionByBlockNumberAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getRawTransactionByBlockNumberAndIndex(blockNr,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1446" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getRawTransactionByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getRawTransactionByHash(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	tx, _, _, _, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			return nil, nil
		}
	}
	return rlp.EncodeToBytes(tx)
}// GetRawTransactionByHash returns the bytes of the transaction for the given hash.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1500" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
key <code>string</code> 

  + Required: ✓ Yes





__3:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes






#### Result




<code>hexutil.Bytes</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getStorageAt", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getStorageAt(address,key,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[ // GetStorageAt returns the storage from the state at the given address, key and
	// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
	// numbers are also allowed.
	:], state.Error()
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L743" target="_">View on GitHub →</a>
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
	
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`


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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- from: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- gasPrice: 
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

		- s: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- value: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- blockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- blockNumber: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- gas: 
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- v: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- r: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`




	```

=== "Raw"

	``` Raw
	{
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
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getTransactionByBlockHashAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getTransactionByBlockHashAndIndex(blockHash,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1438" target="_">View on GitHub →</a>
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
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- type: string
			- title: `blockNumberTag`
			- description: `The block height description`
			- enum: earliest, latest, pending


			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- type: string




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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- r: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- from: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`

		- hash: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- blockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- s: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- transactionIndex: 
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

		- blockNumber: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`




	```

=== "Raw"

	``` Raw
	{
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
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getTransactionByBlockNumberAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getTransactionByBlockNumberAndIndex(blockNr,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1430" target="_">View on GitHub →</a>
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
	
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`


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
		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- blockHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- transactionIndex: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- v: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
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

		- value: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- gas: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- blockNumber: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`


	- type: object


	```

=== "Raw"

	``` Raw
	{
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
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getTransactionByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getTransactionByHash(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return newRPCTransaction(tx, blockHash, blockNumber, index), nil
	}
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx), nil
	}
	return nil, nil
}// GetTransactionByHash returns the transaction for the given hash

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1481" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getTransactionCount", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getTransactionCount(address,blockNrOrHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1462" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getTransactionReceipt", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getTransactionReceipt(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map // GetTransactionReceipt returns the transaction receipt for the given transaction hash.
[string]interface{}, error) {
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, nil
	}
	receipts, err := s.b.GetReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	if len(receipts) <= int(index) {
		return nil, nil
	}
	receipt := receipts[index]
	var signer types.Signer = types.FrontierSigner{}
	if tx.Protected() {
		signer = types.NewEIP155Signer(tx.ChainId())
	}
	from, _ := types.Sender(signer, tx)
	fields := map[string]interface{}{"blockHash": blockHash, "blockNumber": hexutil.Uint64(blockNumber), "transactionHash": hash, "transactionIndex": hexutil.Uint64(index), "from": from, "to": tx.To(), "gasUsed": hexutil.Uint64(receipt.GasUsed), "cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed), "contractAddress": nil, "logs": receipt.Logs, "logsBloom": receipt.Bloom}
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1517" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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
		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- logsBloom: 
			- maxItems: `256`
			- minItems: `256`
			- type: `array`
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`


		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
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

		- error: 
			- type: `string`

		- mixHash: 
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

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- miner: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- receiptsRoot: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- sha3Uncles: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
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

		- parentHash: 
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
            }
        },
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getUncleByBlockHashAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getUncleByBlockHashAndIndex(blockHash,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCMarshalBlockT, error) {
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L698" target="_">View on GitHub →</a>
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
	
	- title: `blockNumberIdentifier`
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- type: string
			- title: `blockNumberTag`


			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- type: string




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
	
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`
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
		- miner: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- size: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- totalDifficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- uncles: 
			- items: 
				- type: `string`
				- description: `Hex representation of a Keccak 256 hash`
				- pattern: `^0x[a-fA-F\d]{64}$`
				- title: `keccak`

			- type: `array`

		- difficulty: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- number: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- timestamp: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
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

		- logsBloom: 
			- items: 
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`
				- description: `Hex representation of the integer`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- mixHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- stateRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- transactionsRoot: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gasUsed: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- hash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- nonce: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- parentHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- receiptsRoot: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- sha3Uncles: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- transactions: 
			- items: 
				- additionalProperties: `true`

			- type: `array`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
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
            }
        },
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getUncleByBlockNumberAndIndex", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getUncleByBlockNumberAndIndex(blockNr,index);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (*RPCMarshalBlockT, error) {
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L682" target="_">View on GitHub →</a>
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- pattern: `^0x([a-fA-F\d])+$`
	- type: string
	- title: `uint`
	- description: `Hex representation of a uint`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getUncleCountByBlockHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getUncleCountByBlockHash(blockHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L722" target="_">View on GitHub →</a>
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

			- enum: earliest, latest, pending
			- type: string
			- title: `blockNumberTag`
			- description: `The block height description`


			- title: `uint64`
			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
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
	
	- type: string
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getUncleCountByBlockNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.getUncleCountByBlockNumber(blockNr);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicBlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L713" target="_">View on GitHub →</a>
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_getWork", "params": []}'
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
//   result[0] - 32 bytes hex encoded current block header pow-hash
//   result[1] - 32 bytes hex encoded seed hash used for DAG
//   result[2] - 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
//   result[3] - hex encoded block number

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/api.go#L41" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_hashrate

Hashrate returns the POW hashrate


#### Params (0)

_None_

#### Result




<code>hexutil.Uint64</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_hashrate", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.hashrate();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().HashRate())
}// Hashrate returns the POW hashrate

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L65" target="_">View on GitHub →</a>
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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_mining", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.mining();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicMinerAPI) Mining() bool {
	return api.e.IsMining()
}// Mining returns an indication if this node is currently mining.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L91" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newBlockFilter

NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
It is part of the filter package since polling goes with eth_getFilterChanges.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter


#### Params (0)

_None_

#### Result




<code>rpc.ID</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_newBlockFilter", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.newBlockFilter();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) NewBlockFilter() rpc.ID {
	var (
		headers		= make(chan *types.Header)
		headerSub	= api.events.SubscribeNewHeads(headers)
	)
	api.filtersMu.Lock()
	api.filters[headerSub.ID] = &filter{typ: BlocksSubscription, deadline: time.NewTimer(deadline), hashes: make([ // NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
	// It is part of the filter package since polling goes with eth_getFilterChanges.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L175" target="_">View on GitHub →</a>
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

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter


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
				- title: `keccak`
				- type: `string`
				- description: `Hex representation of a Keccak 256 hash POINTER`
				- pattern: `^0x[a-fA-F\d]{64}$`

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
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_newFilter", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.newFilter(crit);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) NewFilter(crit FilterCriteria) (rpc.ID, error) {
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
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter
	]*types.Log)
	logsSub, err := api.events.SubscribeLogs(ethereum.FilterQuery(crit), logs)
	if err != nil {
		return rpc.ID(""), err
	}
	api.filtersMu.Lock()
	api.filters[logsSub.ID] = &filter{typ: LogsSubscription, crit: crit, deadline: time.NewTimer(deadline), logs: make([]*types.Log, 0), s: logsSub}
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L354" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_newPendingTransactionFilter

NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
as transactions enter the pending state.

It is part of the filter package because this filter can be used through the
`eth_getFilterChanges` polling method that is also used for log filters.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter


#### Params (0)

_None_

#### Result




<code>rpc.ID</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_newPendingTransactionFilter", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.newPendingTransactionFilter();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) NewPendingTransactionFilter() rpc.ID {
	var (
		pendingTxs	= make(chan [ // NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
		// as transactions enter the pending state.
		//
		// It is part of the filter package because this filter can be used through the
		// `eth_getFilterChanges` polling method that is also used for log filters.
		//
		// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter
		]common.Hash)
		pendingTxSub	= api.events.SubscribePendingTxs(pendingTxs)
	)
	api.filtersMu.Lock()
	api.filters[pendingTxSub.ID] = &filter{typ: PendingTransactionsSubscription, deadline: time.NewTimer(deadline), hashes: make([]common.Hash, 0), s: pendingTxSub}
	api.filtersMu.Unlock()
	go func() {
		for {
			select {
			case ph := <-pendingTxs:
				api.filtersMu.Lock()
				if f, found := api.filters[pendingTxSub.ID]; found {
					f.hashes = append(f.hashes, ph...)
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L105" target="_">View on GitHub →</a>
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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_newSideBlockFilter", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.newSideBlockFilter();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) NewSideBlockFilter() rpc.ID {
	var (
		headers		= make(chan *types.Header)
		headerSub	= api.events.SubscribeNewSideHeads(headers)
	)
	api.filtersMu.Lock()
	api.filters[headerSub.ID] = &filter{typ: SideBlocksSubscription, deadline: time.NewTimer(deadline), hashes: make([ // NewSideBlockFilter creates a filter that fetches blocks that are imported into the chain with a non-canonical status.
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L208" target="_">View on GitHub →</a>
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
				- v: 
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`

				- blockNumber: 
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`

				- from: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- s: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- value: 
					- title: `integer`
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`

				- blockHash: 
					- title: `keccak`
					- type: `string`
					- pattern: `^0x[a-fA-F\d]{64}$`

				- gasPrice: 
					- title: `integer`
					- type: `string`
					- pattern: `^0x[a-fA-F0-9]+$`

				- r: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- to: 
					- title: `keccak`
					- type: `string`
					- pattern: `^0x[a-fA-F\d]{64}$`

				- transactionIndex: 
					- title: `uint64`
					- type: `string`
					- pattern: `^0x([a-fA-F\d])+$`

				- gas: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `uint64`
					- type: `string`

				- hash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

				- nonce: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `uint64`
					- type: `string`

				- input: 
					- pattern: `^0x([a-fA-F\d])+$`
					- title: `dataWord`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_pendingTransactions", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.pendingTransactions();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) PendingTransactions() ([ // PendingTransactions returns the transactions that are in the transaction pool
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
	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		var signer types.Signer = types.HomesteadSigner{}
		if tx.Protected() {
			signer = types.NewEIP155Signer(tx.ChainId())
		}
		from, _ := types.Sender(signer, tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, newRPCPendingTransaction(tx))
		}
	}
	return transactions, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1813" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_protocolVersion

ProtocolVersion returns the current Ethereum protocol version this node supports


#### Params (0)

_None_

#### Result




<code>hexutil.Uint</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: string
	- title: `uint`
	- description: `Hex representation of a uint`
	- pattern: `^0x([a-fA-F\d])+$`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_protocolVersion", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.protocolVersion();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicEthereumAPI) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(s.b.ProtocolVersion())
}// ProtocolVersion returns the current Ethereum protocol version this node supports

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L69" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_resend

Resend accepts an existing transaction and a new gas price and limit. It will remove
the given transaction from the pool and reinsert it with the new gas price and limit.


#### Params (3)

Parameters must be given _by position_.  


__1:__ 
sendArgs <code>SendTxArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- nonce: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- to: 
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`

		- value: 
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
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`


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
	
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
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
	
	- pattern: `^0x([a-fA-F\d])+$`
	- type: string
	- title: `uint64`
	- description: `Hex representation of a uint64`


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
	
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_resend", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.resend(sendArgs,gasPrice,gasLimit);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
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
		var signer types.Signer = types.HomesteadSigner{}
		if p.Protected() {
			signer = types.NewEIP155Signer(p.ChainId())
		}
		wantSigHash := signer.Hash(matchTx)
		if pFrom, err := types.Sender(signer, p); err == nil && pFrom == sendArgs.From && signer.Hash(p) == wantSigHash {
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.Gas = gasLimit
			}
			signedTx, err := s.sign(sendArgs.From, sendArgs.toTransaction())
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1840" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sendRawTransaction

SendRawTransaction will add the signed transaction to the transaction pool.
The sender is responsible for signing the transaction and using the correct nonce.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
encodedTx <code>hexutil.Bytes</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_sendRawTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.sendRawTransaction(encodedTx);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, tx)
}// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1741" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sendTransaction

SendTransaction creates a transaction for the given argument, sign it and submit it to the
transaction pool.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
args <code>SendTxArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gas: 
			- title: `uint64`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- gasPrice: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
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

		- data: 
			- title: `dataWord`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`




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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_sendTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.sendTransaction(args);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	account := accounts.Account{Address: args.From}
	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}
	if args.Nonce == nil {
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1693" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_sign

Sign calculates an ECDSA signature for:
keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).

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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- type: string
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`


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
	
	- type: string
	- title: `dataWord`
	- description: `Hex representation of some bytes`
	- pattern: `^0x([a-fA-F\d])+$`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_sign", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.sign(addr,data);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
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
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1758" target="_">View on GitHub →</a>
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
args <code>SendTxArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- input: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
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

		- data: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- from: 
			- title: `keccak`
			- type: `string`
			- pattern: `^0x[a-fA-F\d]{64}$`


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
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- raw: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `dataWord`
			- type: `string`

		- tx: 
			- additionalProperties: `false`
			- type: `object`




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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_signTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.signTransaction(args);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicTransactionPoolAPI) SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil {
		return nil, fmt.Errorf("gasPrice not specified")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	if err := checkTxFee(args.GasPrice.ToInt(), uint64(*args.Gas), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	tx, err := s.sign(args.From, args.toTransaction())
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, tx}, nil
}// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1783" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_submitHashRate

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
	
	- title: `uint64`
	- description: `Hex representation of a uint64`
	- pattern: `^0x([a-fA-F\d])+$`
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_submitHashRate", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.submitHashRate(rate,id);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *API) SubmitHashRate(rate hexutil.Uint64, id common.Hash) bool {
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/api.go#L92" target="_">View on GitHub →</a>
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
	- type: string
	- title: `integer`


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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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
	
	- title: `keccak`
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_submitWork", "params": []}'
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
<a href="https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/api.go#L66" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_subscribe

Subscribe creates a subscription to an event channel.
Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.


#### Params (2)

Parameters must be given _by position_.  


__1:__ 
subscriptionName <code>RPCSubscriptionParamsName</code> 

  + Required: ✓ Yes





__2:__ 
subscriptionOptions <code>interface{}</code> 

  + Required: No






#### Result



subscriptionID <code>rpc.ID</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_subscribe", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.subscribe(subscriptionName,subscriptionOptions);
	```


<details><summary>Source code</summary>
<p>
```go
func (sub *RPCSubscription) Subscribe(subscriptionName RPCSubscriptionParamsName, subscriptionOptions interface{}) (subscriptionID rpc.ID, err error) {
	return
}// Subscribe creates a subscription to an event channel.
// Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/node/openrpc.go#L211" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_syncing

Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
yet received the latest block headers from its pears. In case it is synchronizing:
- startingBlock: block number this node started to synchronise from
- currentBlock:  block number this node is currently importing
- highestBlock:  block number of the highest block header this node has received from peers
- pulledStates:  number of state entries processed until now
- knownStates:   number of known state entries that still need to be pulled


#### Params (0)

_None_

#### Result



interface <code>interface{}</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_syncing", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.syncing();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicEthereumAPI) Syncing() (interface{}, error) {
	progress := s.b.Downloader().Progress()
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	return map // Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
	// yet received the latest block headers from its pears. In case it is synchronizing:
	// - startingBlock: block number this node started to synchronise from
	// - currentBlock:  block number this node is currently importing
	// - highestBlock:  block number of the highest block header this node has received from peers
	// - pulledStates:  number of state entries processed until now
	// - knownStates:   number of known state entries that still need to be pulled
	[string]interface{}{"startingBlock": hexutil.Uint64(progress.StartingBlock), "currentBlock": hexutil.Uint64(progress.CurrentBlock), "highestBlock": hexutil.Uint64(progress.HighestBlock), "pulledStates": hexutil.Uint64(progress.PulledStates), "knownStates": hexutil.Uint64(progress.KnownStates)}, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L80" target="_">View on GitHub →</a>
</p>
</details>

---



### eth_uninstallFilter

UninstallFilter removes the filter with the given filter id.

https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
id <code>rpc.ID</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_uninstallFilter", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.uninstallFilter(id);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicFilterAPI) UninstallFilter(id rpc.ID) bool {
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
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/filters/api.go#L418" target="_">View on GitHub →</a>
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






#### Result

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "eth_unsubscribe", "params": []}'
	```

=== "Javascript Console"

	``` js
	eth.unsubscribe(id);
	```


<details><summary>Source code</summary>
<p>
```go
func (sub *RPCSubscription) Unsubscribe(id rpc.ID) error {
	return nil
}// Unsubscribe terminates an existing subscription by ID.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/node/openrpc.go#L202" target="_">View on GitHub →</a>
</p>
</details>

---

