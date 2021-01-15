






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-21T13:33:54-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### debug_accountRange

AccountRange enumerates all accounts in the given block and start point in paging request


__Params (6)__

Parameters must be given _by position_.  


__1:__ 
blockNrOrHash <code>rpc.BlockNumberOrHash</code> 

  + Required: ✓ Yes





__2:__ 
start <code>[]byte</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- description: `Hex representation of a variable length byte array`
	- pattern: `^0x([a-fA-F0-9]?)+$`
	- type: string
	- title: `bytes`


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a variable length byte array",
        "pattern": "^0x([a-fA-F0-9]?)+$",
        "title": "bytes",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
maxResults <code>int</code> 

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




__4:__ 
nocode <code>bool</code> 

  + Required: ✓ Yes





__5:__ 
nostorage <code>bool</code> 

  + Required: ✓ Yes





__6:__ 
incompletes <code>bool</code> 

  + Required: ✓ Yes






__Result__




<code>state.IteratorDump</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accounts: 
			- patternProperties: 
				- .*: 
					- type: `object`
					- additionalProperties: `false`
					- properties: 
						- nonce: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- root: 
							- type: `string`

						- storage: 
							- type: `object`
							- patternProperties: 
								- .*: 
									- type: `string`



						- address: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- balance: 
							- type: `string`

						- code: 
							- type: `string`

						- codeHash: 
							- type: `string`

						- key: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `dataWord`
							- type: `string`




			- type: `object`

		- next: 
			- pattern: `^0x([a-fA-F0-9]?)+$`
			- title: `bytes`
			- type: `string`

		- root: 
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "accounts": {
                "patternProperties": {
                    ".*": {
                        "additionalProperties": false,
                        "properties": {
                            "address": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "balance": {
                                "type": "string"
                            },
                            "code": {
                                "type": "string"
                            },
                            "codeHash": {
                                "type": "string"
                            },
                            "key": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "dataWord",
                                "type": "string"
                            },
                            "nonce": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "root": {
                                "type": "string"
                            },
                            "storage": {
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
                "type": "object"
            },
            "next": {
                "pattern": "^0x([a-fA-F0-9]?)+$",
                "title": "bytes",
                "type": "string"
            },
            "root": {
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_accountRange", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.accountRange(blockNrOrHash,start,maxResults,nocode,nostorage,incompletes);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) AccountRange(blockNrOrHash rpc.BlockNumberOrHash, start [ // AccountRange enumerates all accounts in the given block and start point in paging request
]byte, maxResults int, nocode, nostorage, incompletes bool) (state.IteratorDump, error) {
	var stateDb *state.StateDB
	var err error
	if number, ok := blockNrOrHash.Number(); ok {
		if number == rpc.PendingBlockNumber {
			_, stateDb = api.eth.miner.Pending()
		} else {
			var block *types.Block
			if number == rpc.LatestBlockNumber {
				block = api.eth.blockchain.CurrentBlock()
			} else {
				block = api.eth.blockchain.GetBlockByNumber(uint64(number))
			}
			if block == nil {
				return state.IteratorDump{}, fmt.Errorf("block #%d not found", number)
			}
			stateDb, err = api.eth.BlockChain().StateAt(block.Root())
			if err != nil {
				return state.IteratorDump{}, err
			}
		}
	} else if hash, ok := blockNrOrHash.Hash(); ok {
		block := api.eth.blockchain.GetBlockByHash(hash)
		if block == nil {
			return state.IteratorDump{}, fmt.Errorf("block %s not found", hash.Hex())
		}
		stateDb, err = api.eth.BlockChain().StateAt(block.Root())
		if err != nil {
			return state.IteratorDump{}, err
		}
	} else {
		return state.IteratorDump{}, errors.New("either block number or block hash must be specified")
	}
	if maxResults > AccountRangeMaxResults || maxResults <= 0 {
		maxResults = AccountRangeMaxResults
	}
	return stateDb.IteratorDump(nocode, nostorage, incompletes, start, maxResults), nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L382" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_backtraceAt

BacktraceAt sets the log backtrace location. See package log for details on
the pattern syntax.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
location <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_backtraceAt", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.backtraceAt(location);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) BacktraceAt(location string) error {
	return glogger.BacktraceAt(location)
}// BacktraceAt sets the log backtrace location. See package log for details on
// the pattern syntax.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L68" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_blockProfile

BlockProfile turns on goroutine profiling for nsec seconds and writes profile data to
file. It uses a profile rate of 1 for most accurate information. If a different rate is
desired, set the rate and write the profile manually.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_blockProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.blockProfile(file,nsec);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) BlockProfile(file string, nsec uint) error {
	runtime.SetBlockProfileRate(1)
	time.Sleep(time.Duration(nsec) * time.Second)
	defer runtime.SetBlockProfileRate(0)
	return writeProfile("block", file)
}// BlockProfile turns on goroutine profiling for nsec seconds and writes profile data to
// file. It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L147" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_chaindbCompact

ChaindbCompact flattens the entire key-value database into a single level,
removing all unused slots and merging all keys.


__Params (0)__

_None_

__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_chaindbCompact", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.chaindbCompact();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) ChaindbCompact() error {
	for b := byte(0); b < 255; b++ {
		log.Info("Compacting chain database", "range", fmt.Sprintf("0x%0.2X-0x%0.2X", b, b+1))
		if err := api.b.ChainDb().Compact([ // ChaindbCompact flattens the entire key-value database into a single level,
		// removing all unused slots and merging all keys.
		]byte{b}, []byte{b + 1}); err != nil {
			log.Error("Database compaction failed", "err", err)
			return err
		}
	}
	return nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L2004" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_chaindbProperty

ChaindbProperty returns leveldb properties of the key-value database.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
property <code>string</code> 

  + Required: ✓ Yes






__Result__




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_chaindbProperty", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.chaindbProperty(property);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) ChaindbProperty(property string) (string, error) {
	if property == "" {
		property = "leveldb.stats"
	} else if !strings.HasPrefix(property, "leveldb.") {
		property = "leveldb." + property
	}
	return api.b.ChainDb().Stat(property)
}// ChaindbProperty returns leveldb properties of the key-value database.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1993" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_cpuProfile

CpuProfile turns on CPU profiling for nsec seconds and writes
profile data to file.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_cpuProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.cpuProfile(file,nsec);
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) CpuProfile(file string, nsec uint) error {
	if err := h.StartCPUProfile(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second)
	h.StopCPUProfile()
	return nil
}// CpuProfile turns on CPU profiling for nsec seconds and writes
// profile data to file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L88" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_dumpBlock

DumpBlock retrieves the entire state of the database at a given block.


__Params (1)__

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





__Result__




<code>state.Dump</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accounts: 
			- patternProperties: 
				- .*: 
					- properties: 
						- nonce: 
							- type: `string`
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`

						- root: 
							- type: `string`

						- storage: 
							- patternProperties: 
								- .*: 
									- type: `string`


							- type: `object`

						- address: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- balance: 
							- type: `string`

						- code: 
							- type: `string`

						- codeHash: 
							- type: `string`

						- key: 
							- title: `dataWord`
							- type: `string`
							- pattern: `^0x([a-fA-F\d])+$`


					- type: `object`
					- additionalProperties: `false`


			- type: `object`

		- root: 
			- type: `string`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "accounts": {
                "patternProperties": {
                    ".*": {
                        "additionalProperties": false,
                        "properties": {
                            "address": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "balance": {
                                "type": "string"
                            },
                            "code": {
                                "type": "string"
                            },
                            "codeHash": {
                                "type": "string"
                            },
                            "key": {
                                "pattern": "^0x([a-fA-F\\d])+$",
                                "title": "dataWord",
                                "type": "string"
                            },
                            "nonce": {
                                "pattern": "^0x[a-fA-F0-9]+$",
                                "title": "integer",
                                "type": "string"
                            },
                            "root": {
                                "type": "string"
                            },
                            "storage": {
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
                "type": "object"
            },
            "root": {
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_dumpBlock", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.dumpBlock(blockNr);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error) {
	if blockNr == rpc.PendingBlockNumber {
		_, stateDb := api.eth.miner.Pending()
		return stateDb.RawDump(false, false, true), nil
	}
	var block *types.Block
	if blockNr == rpc.LatestBlockNumber {
		block = api.eth.blockchain.CurrentBlock()
	} else {
		block = api.eth.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return state.Dump{}, fmt.Errorf("block #%d not found", blockNr)
	}
	stateDb, err := api.eth.BlockChain().StateAt(block.Root())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(false, false, true), nil
}// DumpBlock retrieves the entire state of the database at a given block.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L304" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_freeOSMemory

FreeOSMemory forces a garbage collection.


__Params (0)__

_None_

__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_freeOSMemory", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.freeOSMemory();
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) FreeOSMemory() {
	debug.FreeOSMemory()
}// FreeOSMemory forces a garbage collection.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L200" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_gcStats

GcStats returns GC statistics.


__Params (0)__

_None_

__Result__




<code>*debug.GCStats</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- PauseTotal: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- LastGC: 
			- type: `string`
			- format: `date-time`

		- NumGC: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- Pause: 
			- items: 
				- type: `string`
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`

			- type: `array`

		- PauseEnd: 
			- items: 
				- format: `date-time`
				- type: `string`

			- type: `array`

		- PauseQuantiles: 
			- type: `array`
			- items: 
				- type: `string`
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`



	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "LastGC": {
                "format": "date-time",
                "type": "string"
            },
            "NumGC": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Pause": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "type": "array"
            },
            "PauseEnd": {
                "items": {
                    "format": "date-time",
                    "type": "string"
                },
                "type": "array"
            },
            "PauseQuantiles": {
                "items": {
                    "description": "Hex representation of the integer",
                    "pattern": "^0x[a-fA-F0-9]+$",
                    "title": "integer",
                    "type": "string"
                },
                "type": "array"
            },
            "PauseTotal": {
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_gcStats", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.gcStats();
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) GcStats() *debug.GCStats {
	s := new(debug.GCStats)
	debug.ReadGCStats(s)
	return s
}// GcStats returns GC statistics.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L80" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getBadBlocks

GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
and returns them as a JSON list of block-hashes


__Params (0)__

_None_

__Result__



BadBlockArgs <code>[]*BadBlockArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- block: 
					- additionalProperties: `false`
					- properties: 
						- difficulty: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- gasLimit: 
							- type: `string`
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`

						- extraData: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `dataWord`
							- type: `string`

						- hash: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- miner: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- number: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- receiptsRoot: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- timestamp: 
							- title: `uint64`
							- type: `string`
							- pattern: `^0x([a-fA-F\d])+$`

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

						- error: 
							- type: `string`

						- gasUsed: 
							- pattern: `^0x([a-fA-F\d])+$`
							- title: `uint64`
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

						- transactions: 
							- items: 
								- additionalProperties: `true`

							- type: `array`

						- logsBloom: 
							- maxItems: `256`
							- minItems: `256`
							- type: `array`
							- items: 
								- type: `string`
								- description: `Hex representation of the integer`
								- pattern: `^0x[a-fA-F0-9]+$`
								- title: `integer`


						- sha3Uncles: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`


					- type: `object`

				- hash: 
					- type: `string`
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`

				- rlp: 
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
                    "block": {
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
                        "type": "object"
                    },
                    "hash": {
                        "pattern": "^0x[a-fA-F\\d]{64}$",
                        "title": "keccak",
                        "type": "string"
                    },
                    "rlp": {
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_getBadBlocks", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.getBadBlocks();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) GetBadBlocks(ctx context.Context) ([ // GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
]*BadBlockArgs, error) {
	blocks := api.eth.BlockChain().BadBlocks()
	results := make([]*BadBlockArgs, len(blocks))
	var err error
	for i, block := range blocks {
		results[i] = &BadBlockArgs{Hash: block.Hash()}
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			results[i].RLP = err.Error()
		} else {
			results[i].RLP = fmt.Sprintf("0x%x", rlpBytes)
		}
		if results[i].Block, err = ethapi.RPCMarshalBlock(block, true, true); err != nil {
			results[i].Block = &ethapi.RPCMarshalBlockT{Error: err.Error()}
		}
	}
	return results, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L357" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getBlockRlp

GetBlockRlp retrieves the RLP encoded for of a single block.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
number <code>uint64</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: string
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`


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





__Result__




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_getBlockRlp", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.getBlockRlp(number);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) GetBlockRlp(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encoded), nil
}// GetBlockRlp retrieves the RLP encoded for of a single block.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1908" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getModifiedAccountsByHash

GetModifiedAccountsByHash returns all accounts that have changed between the
two blocks specified. A change is defined as a difference in nonce, balance,
code hash, or storage hash.

With one parameter, returns the list of accounts modified in the specified block.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
startHash <code>common.Hash</code> 

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
endHash <code>*common.Hash</code> 

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





__Result__



commonAddress <code>[]common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- title: `keccak`
			- description: `Hex representation of a Keccak 256 hash POINTER`
			- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_getModifiedAccountsByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.getModifiedAccountsByHash(startHash,endHash);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([ // GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.eth.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}
	if endHash == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L513" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getModifiedAccountsByNumber

GetModifiedAccountsByNumber returns all accounts that have changed between the
two blocks specified. A change is defined as a difference in nonce, balance,
code hash, or storage hash.

With one parameter, returns the list of accounts modified in the specified block.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
startNum <code>uint64</code> 

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




__2:__ 
endNum <code>*uint64</code> 

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





__Result__



commonAddress <code>[]common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- title: `keccak`
			- description: `Hex representation of a Keccak 256 hash POINTER`
			- pattern: `^0x[a-fA-F\d]{64}$`
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_getModifiedAccountsByNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.getModifiedAccountsByNumber(startNum,endNum);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([ // GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.eth.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}
	if endNum == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L485" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_goTrace

GoTrace turns on tracing for nsec seconds and writes
trace data to file.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_goTrace", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.goTrace(file,nsec);
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) GoTrace(file string, nsec uint) error {
	if err := h.StartGoTrace(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second)
	h.StopGoTrace()
	return nil
}// GoTrace turns on tracing for nsec seconds and writes
// trace data to file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L135" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_memStats

MemStats returns detailed runtime memory statistics.


__Params (0)__

_None_

__Result__




<code>*runtime.MemStats</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- HeapInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MSpanInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NumForcedGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- PauseEnd: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- TotalAlloc: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- BuckHashSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Lookups: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- OtherSys: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- Frees: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapAlloc: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- LastGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NextGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- PauseNs: 
			- minItems: `256`
			- type: `array`
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`

		- StackInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- BySize: 
			- minItems: `61`
			- type: `array`
			- items: 
				- additionalProperties: `false`
				- properties: 
					- Frees: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- Mallocs: 
						- title: `integer`
						- type: `string`
						- pattern: `^0x[a-fA-F0-9]+$`

					- Size: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- maxItems: `61`

		- Sys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- StackSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapIdle: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- HeapReleased: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NumGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- PauseTotalNs: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- DebugGC: 
			- type: `boolean`

		- HeapSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- EnableGC: 
			- type: `boolean`

		- GCSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapObjects: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MCacheInuse: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- MSpanSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Mallocs: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- GCCPUFraction: 
			- type: `number`

		- MCacheSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Alloc: 
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
            "Alloc": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "BuckHashSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "BySize": {
                "items": {
                    "additionalProperties": false,
                    "properties": {
                        "Frees": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "Mallocs": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        },
                        "Size": {
                            "pattern": "^0x[a-fA-F0-9]+$",
                            "title": "integer",
                            "type": "string"
                        }
                    },
                    "type": "object"
                },
                "maxItems": 61,
                "minItems": 61,
                "type": "array"
            },
            "DebugGC": {
                "type": "boolean"
            },
            "EnableGC": {
                "type": "boolean"
            },
            "Frees": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "GCCPUFraction": {
                "type": "number"
            },
            "GCSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapAlloc": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapIdle": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapInuse": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapObjects": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapReleased": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "HeapSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "LastGC": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Lookups": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "MCacheInuse": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "MCacheSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "MSpanInuse": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "MSpanSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Mallocs": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "NextGC": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "NumForcedGC": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "NumGC": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "OtherSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "PauseEnd": {
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
            "PauseNs": {
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
            "PauseTotalNs": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "StackInuse": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "StackSys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "Sys": {
                "pattern": "^0x[a-fA-F0-9]+$",
                "title": "integer",
                "type": "string"
            },
            "TotalAlloc": {
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_memStats", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.memStats();
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) MemStats() *runtime.MemStats {
	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)
	return s
}// MemStats returns detailed runtime memory statistics.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L73" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_mutexProfile

MutexProfile turns on mutex profiling for nsec seconds and writes profile data to file.
It uses a profile rate of 1 for most accurate information. If a different rate is
desired, set the rate and write the profile manually.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_mutexProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.mutexProfile(file,nsec);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) MutexProfile(file string, nsec uint) error {
	runtime.SetMutexProfileFraction(1)
	time.Sleep(time.Duration(nsec) * time.Second)
	defer runtime.SetMutexProfileFraction(0)
	return writeProfile("mutex", file)
}// MutexProfile turns on mutex profiling for nsec seconds and writes profile data to file.
// It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L168" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_preimage

Preimage is a debug API function that returns the preimage for a sha3 hash, if known.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
hash <code>common.Hash</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`


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





__Result__




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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_preimage", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.preimage(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := rawdb.ReadPreimage(api.eth.ChainDb(), hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L341" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_printBlock

PrintBlock retrieves a block and returns its pretty printed form.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
number <code>uint64</code> 

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





__Result__




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_printBlock", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.printBlock(number);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return spew.Sdump(block), nil
}// PrintBlock retrieves a block and returns its pretty printed form.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1960" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_removePendingTransaction

RemovePendingTransaction removes a transaction from the txpool.
It returns the transaction removed, if any.


__Params (1)__

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





__Result__




<code>*types.Transaction</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "type": [
            "object"
        ]
    }
	```



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_removePendingTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.removePendingTransaction(hash);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) RemovePendingTransaction(hash common.Hash) (*types.Transaction, error) {
	return api.eth.txPool.RemoveTx(hash), nil
}// RemovePendingTransaction removes a transaction from the txpool.
// It returns the transaction removed, if any.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L565" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_seedHash

SeedHash retrieves the seed hash of a block.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
number <code>uint64</code> 

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





__Result__




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_seedHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.seedHash(number);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) SeedHash(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	ecip1099FBlock := api.b.ChainConfig().GetEthashECIP1099Transition()
	epochLength := ethash.CalcEpochLength(number, ecip1099FBlock)
	epoch := ethash.CalcEpoch(number, epochLength)
	return fmt.Sprintf("0x%x", ethash.SeedHash(epoch, epochLength)), nil
}// SeedHash retrieves the seed hash of a block.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1969" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setBlockProfileRate

SetBlockProfileRate sets the rate of goroutine block profile data collection.
rate 0 disables block profiling.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
rate <code>int</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_setBlockProfileRate", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.setBlockProfileRate(rate);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) SetBlockProfileRate(rate int) {
	runtime.SetBlockProfileRate(rate)
}// SetBlockProfileRate sets the rate of goroutine block profile data collection.
// rate 0 disables block profiling.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L156" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setGCPercent

SetGCPercent sets the garbage collection target percentage. It returns the previous
setting. A negative value disables GC.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
v <code>int</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- pattern: `^0x[a-fA-F0-9]+$`
	- type: string
	- title: `integer`
	- description: `Hex representation of the integer`


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





__Result__




<code>int</code> 

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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_setGCPercent", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.setGCPercent(v);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) SetGCPercent(v int) int {
	return debug.SetGCPercent(v)
}// SetGCPercent sets the garbage collection target percentage. It returns the previous
// setting. A negative value disables GC.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L206" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setHead

SetHead rewinds the head of the blockchain to a previous block.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
number <code>hexutil.Uint64</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_setHead", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.setHead(number);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) SetHead(number hexutil.Uint64) {
	api.b.SetHead(uint64(number))
}// SetHead rewinds the head of the blockchain to a previous block.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L2016" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setMutexProfileFraction

SetMutexProfileFraction sets the rate of mutex profiling.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
rate <code>int</code> 

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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_setMutexProfileFraction", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.setMutexProfileFraction(rate);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) SetMutexProfileFraction(rate int) {
	runtime.SetMutexProfileFraction(rate)
}// SetMutexProfileFraction sets the rate of mutex profiling.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L177" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stacks

Stacks returns a printed representation of the stacks of all goroutines.


__Params (0)__

_None_

__Result__




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_stacks", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.stacks();
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) Stacks() string {
	buf := new(bytes.Buffer)
	pprof.Lookup("goroutine").WriteTo(buf, 2)
	return buf.String()
}// Stacks returns a printed representation of the stacks of all goroutines.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L193" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_standardTraceBadBlockToFile

StandardTraceBadBlockToFile dumps the structured logs created during the
execution of EVM against a block pulled from the pool of bad ones to the
local file system and returns a list of files to the caller.


__Params (2)__

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




__2:__ 
config <code>*StdTraceConfig</code> 

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

		- TxHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- overrides: 
			- additionalProperties: `true`

		- DisableStack: 
			- type: `boolean`


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
            "TxHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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





__Result__



string <code>[]string</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: array
	- items: 

			- type: string




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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_standardTraceBadBlockToFile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.standardTraceBadBlockToFile(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) StandardTraceBadBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([ // StandardTraceBadBlockToFile dumps the structured logs created during the
// execution of EVM against a block pulled from the pool of bad ones to the
// local file system and returns a list of files to the caller.
]string, error) {
	blocks := api.eth.blockchain.BadBlocks()
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.standardTraceBlockToFile(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L446" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_standardTraceBlockToFile

StandardTraceBlockToFile dumps the structured logs created during the
execution of EVM to the local file system and returns a list of files
to the caller.


__Params (2)__

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




__2:__ 
config <code>*StdTraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- Debug: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- Limit: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- TxHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- overrides: 
			- additionalProperties: `true`




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
            "TxHash": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
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





__Result__



string <code>[]string</code> 

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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_standardTraceBlockToFile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.standardTraceBlockToFile(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) StandardTraceBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([ // StandardTraceBlockToFile dumps the structured logs created during the
// execution of EVM to the local file system and returns a list of files
// to the caller.
]string, error) {
	block := api.eth.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("block %#x not found", hash)
	}
	return api.standardTraceBlockToFile(ctx, block, config)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L435" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_startCPUProfile

StartCPUProfile turns on CPU profiling, writing to the given file.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_startCPUProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.startCPUProfile(file);
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) StartCPUProfile(file string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.cpuW != nil {
		return errors.New("CPU profiling already in progress")
	}
	f, err := os.Create(expandHome(file))
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return err
	}
	h.cpuW = f
	h.cpuFile = file
	log.Info("CPU profiling started", "dump", h.cpuFile)
	return nil
}// StartCPUProfile turns on CPU profiling, writing to the given file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L98" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_startGoTrace

StartGoTrace turns on tracing, writing to the given file.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_startGoTrace", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.startGoTrace(file);
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) StartGoTrace(file string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.traceW != nil {
		return errors.New("trace already in progress")
	}
	f, err := os.Create(expandHome(file))
	if err != nil {
		return err
	}
	if err := trace.Start(f); err != nil {
		f.Close()
		return err
	}
	h.traceW = f
	h.traceFile = file
	log.Info("Go tracing started", "dump", h.traceFile)
	return nil
}// StartGoTrace turns on tracing, writing to the given file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/trace.go#L30" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stopCPUProfile

StopCPUProfile stops an ongoing CPU profile.


__Params (0)__

_None_

__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_stopCPUProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.stopCPUProfile();
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) StopCPUProfile() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	pprof.StopCPUProfile()
	if h.cpuW == nil {
		return errors.New("CPU profiling not in progress")
	}
	log.Info("Done writing CPU profile", "dump", h.cpuFile)
	h.cpuW.Close()
	h.cpuW = nil
	h.cpuFile = ""
	return nil
}// StopCPUProfile stops an ongoing CPU profile.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L119" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stopGoTrace

StopTrace stops an ongoing trace.


__Params (0)__

_None_

__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_stopGoTrace", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.stopGoTrace();
	```


<details><summary>Source code</summary>
<p>
```go
func (h *HandlerT) StopGoTrace() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	trace.Stop()
	if h.traceW == nil {
		return errors.New("trace not in progress")
	}
	log.Info("Done writing Go trace", "dump", h.traceFile)
	h.traceW.Close()
	h.traceW = nil
	h.traceFile = ""
	return nil
}// StopTrace stops an ongoing trace.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/trace.go#L51" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_storageRangeAt

StorageRangeAt returns the storage at the given block height and transaction index.


__Params (5)__

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
txIndex <code>int</code> 

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
contractAddress <code>common.Address</code> 

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




__4:__ 
keyStart <code>hexutil.Bytes</code> 

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




__5:__ 
maxResult <code>int</code> 

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





__Result__




<code>StorageRangeResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- nextKey: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- storage: 
			- patternProperties: 
				- .*: 
					- additionalProperties: `false`
					- properties: 
						- key: 
							- pattern: `^0x[a-fA-F\d]{64}$`
							- title: `keccak`
							- type: `string`

						- value: 
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
        "additionalProperties": false,
        "properties": {
            "nextKey": {
                "pattern": "^0x[a-fA-F\\d]{64}$",
                "title": "keccak",
                "type": "string"
            },
            "storage": {
                "patternProperties": {
                    ".*": {
                        "additionalProperties": false,
                        "properties": {
                            "key": {
                                "pattern": "^0x[a-fA-F\\d]{64}$",
                                "title": "keccak",
                                "type": "string"
                            },
                            "value": {
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_storageRangeAt", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.storageRangeAt(blockHash,txIndex,contractAddress,keyStart,maxResult);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) StorageRangeAt(blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	block := api.eth.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return StorageRangeResult{}, fmt.Errorf("block %#x not found", blockHash)
	}
	_, _, statedb, err := api.computeTxEnv(block, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}// StorageRangeAt returns the storage at the given block height and transaction index.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L440" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_testSignCliqueBlock

TestSignCliqueBlock fetches the given block number, and attempts to sign it as a clique header with the
given address, returning the address of the recovered signature

This is a temporary method to debug the externalsigner integration,
TODO: Remove this method when the integration is mature


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
address <code>common.Address</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- description: `Hex representation of a Keccak 256 hash POINTER`
	- pattern: `^0x[a-fA-F\d]{64}$`
	- type: string
	- title: `keccak`


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
number <code>uint64</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: string
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`


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





__Result__




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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_testSignCliqueBlock", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.testSignCliqueBlock(address,number);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PublicDebugAPI) TestSignCliqueBlock(ctx context.Context, address common.Address, number uint64) (common.Address, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return common.Address{}, fmt.Errorf("block #%d not found", number)
	}
	header := block.Header()
	header.Extra = make([ // TestSignCliqueBlock fetches the given block number, and attempts to sign it as a clique header with the
	// given address, returning the address of the recovered signature
	//
	// This is a temporary method to debug the externalsigner integration,
	// TODO: Remove this method when the integration is mature
	]byte, 32+65)
	encoded := clique.CliqueRLP(header)
	account := accounts.Account{Address: address}
	wallet, err := api.b.AccountManager().Find(account)
	if err != nil {
		return common.Address{}, err
	}
	signature, err := wallet.SignData(account, accounts.MimetypeClique, encoded)
	if err != nil {
		return common.Address{}, err
	}
	sealHash := clique.SealHash(header).Bytes()
	log.Info("test signing of clique block", "Sealhash", fmt.Sprintf("%x", sealHash), "signature", fmt.Sprintf("%x", signature))
	pubkey, err := crypto.Ecrecover(sealHash, signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	return signer, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1925" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBadBlock

TraceBadBlock returns the structured logs created during the execution of
EVM against a block pulled from the pool of bad ones and returns them as a JSON
object.


__Params (2)__

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




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: object
	- additionalProperties: `false`
	- properties: 
		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- Limit: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- Tracer: 
			- type: `string`

		- Debug: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`




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





__Result__



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- result: 
					- additionalProperties: `true`

				- error: 
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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceBadBlock", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceBadBlock(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceBadBlock(ctx context.Context, hash common.Hash, config *TraceConfig) ([ // TraceBadBlock returns the structured logs created during the execution of
// EVM against a block pulled from the pool of bad ones and returns them as a JSON
// object.
]*txTraceResult, error) {
	blocks := api.eth.blockchain.BadBlocks()
	for _, block := range blocks {
		if block.Hash() == hash {
			return api.traceBlock(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L422" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlock

TraceBlock returns the structured logs created during the execution of EVM
and returns them as a JSON object.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
blob <code>[]byte</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- title: `bytes`
	- description: `Hex representation of a variable length byte array`
	- pattern: `^0x([a-fA-F0-9]?)+$`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of a variable length byte array",
        "pattern": "^0x([a-fA-F0-9]?)+$",
        "title": "bytes",
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
	
	- properties: 
		- Limit: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- Tracer: 
			- type: `string`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- overrides: 
			- additionalProperties: `true`

		- Debug: 
			- type: `boolean`


	- type: object
	- additionalProperties: `false`


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





__Result__



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- properties: 
				- error: 
					- type: `string`

				- result: 
					- additionalProperties: `true`


			- type: object
			- additionalProperties: `false`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceBlock", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceBlock(blob,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceBlock(ctx context.Context, blob [ // TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
]byte, config *TraceConfig) ([]*txTraceResult, error) {
	return traceBlockRLP(ctx, api.eth, blob, config)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L405" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockByHash

TraceBlockByHash returns the structured logs created during the execution of
EVM and returns them as a JSON object.


__Params (2)__

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




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Debug: 
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

		- Tracer: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`

		- DisableMemory: 
			- type: `boolean`

		- Reexec: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- Timeout: 
			- type: `string`


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





__Result__



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- properties: 
				- error: 
					- type: `string`

				- result: 
					- additionalProperties: `true`


			- type: object
			- additionalProperties: `false`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceBlockByHash", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceBlockByHash(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([ // TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	block := api.eth.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("block %#x not found", hash)
	}
	return api.traceBlock(ctx, block, config)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L385" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockByNumber

TraceBlockByNumber returns the structured logs created during the execution of
EVM and returns them as a JSON object.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
number <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- oneOf: 

			- title: `blockNumberTag`
			- description: `The block height description`
			- enum: earliest, latest, pending
			- type: string


			- pattern: `^0x([a-fA-F\d])+$`
			- type: string
			- title: `uint64`
			- description: `Hex representation of a uint64`


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
		- DisableStack: 
			- type: `boolean`

		- Limit: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Tracer: 
			- type: `string`

		- DisableReturnData: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Timeout: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`

		- Debug: 
			- type: `boolean`


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





__Result__



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- items: 

			- properties: 
				- error: 
					- type: `string`

				- result: 
					- additionalProperties: `true`


			- type: object
			- additionalProperties: `false`


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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceBlockByNumber", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceBlockByNumber(number,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([ // TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	return traceBlockByNumber(ctx, api.eth, number, config)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L379" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockFromFile

TraceBlockFromFile returns the structured logs created during the execution of
EVM and returns them as a JSON object.


__Params (2)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Debug: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Reexec: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- Timeout: 
			- type: `string`

		- DisableMemory: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
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





__Result__



txTraceResult <code>[]*txTraceResult</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: array
	- items: 

			- additionalProperties: `false`
			- properties: 
				- error: 
					- type: `string`

				- result: 
					- additionalProperties: `true`


			- type: object




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



__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceBlockFromFile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceBlockFromFile(file,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([ // TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, blob, config)
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L411" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceCall

TraceCall lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
if the given transaction was added on top of the provided block and returns them as a JSON object.
You can provide -2 as a block number to trace on top of the pending block.


__Params (3)__

Parameters must be given _by position_.  


__1:__ 
args <code>ethapi.CallArgs</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- to: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- value: 
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`

		- data: 
			- title: `dataWord`
			- type: `string`
			- pattern: `^0x([a-fA-F\d])+$`

		- from: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
			- type: `string`

		- gas: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- gasPrice: 
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
		- overrides: 
			- additionalProperties: `true`

		- Debug: 
			- type: `boolean`

		- DisableReturnData: 
			- type: `boolean`

		- DisableStorage: 
			- type: `boolean`

		- Reexec: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Timeout: 
			- type: `string`

		- DisableMemory: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Tracer: 
			- type: `string`


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





__Result__



interface <code>interface{}</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceCall", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceCall(args,blockNrOrHash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceCall(ctx context.Context, args ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	statedb, header, err := api.eth.APIBackend.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		var block *types.Block
		if hash, ok := blockNrOrHash.Hash(); ok {
			block = api.eth.blockchain.GetBlockByHash(hash)
		} else if number, ok := blockNrOrHash.Number(); ok {
			block = api.eth.blockchain.GetBlockByNumber(uint64(number))
		}
		if block == nil {
			return nil, fmt.Errorf("block %v not found: %v", blockNrOrHash, err)
		}
		reexec := defaultTraceReexec
		if config != nil && config.Reexec != nil {
			reexec = *config.Reexec
		}
		_, _, statedb, err = api.computeTxEnv(block, 0, reexec)
		if err != nil {
			return nil, err
		}
	}
	msg := args.ToMessage(api.eth.APIBackend.RPCGasCap())
	vmctx := core.NewEVMContext(msg, header, api.eth.blockchain, nil)
	return api.traceTx(ctx, msg, vmctx, statedb, config)
}// TraceCall lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
// Try to retrieve the specified block

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L808" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceTransaction

TraceTransaction returns the structured logs created during the execution of EVM
and returns them as a JSON object.


__Params (2)__

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




__2:__ 
config <code>*TraceConfig</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- properties: 
		- DisableStorage: 
			- type: `boolean`

		- Limit: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Debug: 
			- type: `boolean`

		- DisableMemory: 
			- type: `boolean`

		- DisableStack: 
			- type: `boolean`

		- Tracer: 
			- type: `string`

		- overrides: 
			- additionalProperties: `true`

		- DisableReturnData: 
			- type: `boolean`

		- Reexec: 
			- title: `integer`
			- type: `string`
			- pattern: `^0x[a-fA-F0-9]+$`

		- Timeout: 
			- type: `string`


	- type: object
	- additionalProperties: `false`


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





__Result__



interface <code>interface{}</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_traceTransaction", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.traceTransaction(hash,config);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateDebugAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	return traceTransaction(ctx, api.eth, hash, config)
}// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api_tracer.go#L801" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_verbosity

Verbosity sets the log verbosity ceiling. The verbosity of individual packages
and source files can be raised using Vmodule.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
level <code>int</code> 

  + Required: ✓ Yes

 
=== "Schema"

	``` Schema
	
	- type: string
	- title: `integer`
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`


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





__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_verbosity", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.verbosity(level);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) Verbosity(level int) {
	glogger.Verbosity(log.Lvl(level))
}// Verbosity sets the log verbosity ceiling. The verbosity of individual packages
// and source files can be raised using Vmodule.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L57" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_vmodule

Vmodule sets the log verbosity pattern. See package log for details on the
pattern syntax.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
pattern <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_vmodule", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.vmodule(pattern);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) Vmodule(pattern string) error {
	return glogger.Vmodule(pattern)
}// Vmodule sets the log verbosity pattern. See package log for details on the
// pattern syntax.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L62" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeBlockProfile

WriteBlockProfile writes a goroutine blocking profile to the given file.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_writeBlockProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.writeBlockProfile(file);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) WriteBlockProfile(file string) error {
	return writeProfile("block", file)
}// WriteBlockProfile writes a goroutine blocking profile to the given file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L161" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeMemProfile

WriteMemProfile writes an allocation profile to the given file.
Note that the profiling rate cannot be set through the API,
it must be set on the command line.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_writeMemProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.writeMemProfile(file);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) WriteMemProfile(file string) error {
	return writeProfile("heap", file)
}// WriteMemProfile writes an allocation profile to the given file.
// Note that the profiling rate cannot be set through the API,
// it must be set on the command line.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L188" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeMutexProfile

WriteMutexProfile writes a goroutine blocking profile to the given file.


__Params (1)__

Parameters must be given _by position_.  


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






__Result__

_None_

__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "debug_writeMutexProfile", "params": []}'
	```

=== "Javascript Console"

	``` js
	debug.writeMutexProfile(file);
	```


<details><summary>Source code</summary>
<p>
```go
func (*HandlerT) WriteMutexProfile(file string) error {
	return writeProfile("mutex", file)
}// WriteMutexProfile writes a goroutine blocking profile to the given file.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/debug/api.go#L181" target="_">View on GitHub →</a>
</p>
</details>

---

