






| Entity | Version |
| --- | --- |
| Source | <code>1.11.23-unstable/generated-at:2021-04-30T19:24:24+03:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### debug_accountRange

AccountRange enumerates all accounts in the given block and start point in paging request


#### Params (6)

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
	- title: `bytes`
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




__3:__ 
maxResults <code>int</code> 

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




__4:__ 
nocode <code>bool</code> 

  + Required: ✓ Yes





__5:__ 
nostorage <code>bool</code> 

  + Required: ✓ Yes





__6:__ 
incompletes <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>state.IteratorDump</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accounts: 
			- patternProperties: 
				- .*: 
					- additionalProperties: `false`
					- properties: 
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

						- nonce: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- root: 
							- type: `string`

						- storage: 
							- patternProperties: 
								- .*: 
									- type: `string`


							- type: `object`


					- type: `object`


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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_accountRange", "params": [<blockNrOrHash>, <start>, <maxResults>, <nocode>, <nostorage>, <incompletes>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_accountRange", "params": [<blockNrOrHash>, <start>, <maxResults>, <nocode>, <nostorage>, <incompletes>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L389" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_backtraceAt

BacktraceAt sets the log backtrace location. See package log for details on
the pattern syntax.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
location <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_backtraceAt", "params": [<location>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_backtraceAt", "params": [<location>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L68" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_blockProfile

BlockProfile turns on goroutine profiling for nsec seconds and writes profile data to
file. It uses a profile rate of 1 for most accurate information. If a different rate is
desired, set the rate and write the profile manually.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_blockProfile", "params": [<file>, <nsec>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_blockProfile", "params": [<file>, <nsec>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L147" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_chaindbCompact

ChaindbCompact flattens the entire key-value database into a single level,
removing all unused slots and merging all keys.


#### Params (0)

_None_

#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_chaindbCompact", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_chaindbCompact", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2069" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_chaindbProperty

ChaindbProperty returns leveldb properties of the key-value database.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
property <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_chaindbProperty", "params": [<property>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_chaindbProperty", "params": [<property>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2058" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_cpuProfile

CpuProfile turns on CPU profiling for nsec seconds and writes
profile data to file.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_cpuProfile", "params": [<file>, <nsec>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_cpuProfile", "params": [<file>, <nsec>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L88" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_dumpBlock

DumpBlock retrieves the entire state of the database at a given block.


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




<code>state.Dump</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- accounts: 
			- patternProperties: 
				- .*: 
					- additionalProperties: `false`
					- properties: 
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

						- nonce: 
							- pattern: `^0x[a-fA-F0-9]+$`
							- title: `integer`
							- type: `string`

						- root: 
							- type: `string`

						- storage: 
							- patternProperties: 
								- .*: 
									- type: `string`


							- type: `object`


					- type: `object`


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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_dumpBlock", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_dumpBlock", "params": [<blockNr>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L304" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_freeOSMemory

FreeOSMemory forces a garbage collection.


#### Params (0)

_None_

#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_freeOSMemory", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_freeOSMemory", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L200" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_gcStats

GcStats returns GC statistics.


#### Params (0)

_None_

#### Result




<code>*debug.GCStats</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- LastGC: 
			- format: `date-time`
			- type: `string`

		- NumGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Pause: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- type: `array`

		- PauseEnd: 
			- items: 
				- format: `date-time`
				- type: `string`

			- type: `array`

		- PauseQuantiles: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- type: `array`

		- PauseTotal: 
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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_gcStats", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_gcStats", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L80" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getBadBlocks

GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
and returns them as a JSON list of block-hashes


#### Params (0)

_None_

#### Result



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


					- type: `object`

				- hash: 
					- pattern: `^0x[a-fA-F\d]{64}$`
					- title: `keccak`
					- type: `string`

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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_getBadBlocks", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_getBadBlocks", "params": []}'
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
	var (
		err	error
		blocks	= rawdb.ReadAllBadBlocks(api.eth.chainDb)
		results	= make([]*BadBlockArgs, 0, len(blocks))
	)
	for _, block := range blocks {
		var (
			blockRlp	string
			blockJSON	*ethapi.RPCMarshalBlockT
		)
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			blockRlp = err.Error()
		} else {
			blockRlp = fmt.Sprintf("0x%x", rlpBytes)
		}
		if blockJSON, err = ethapi.RPCMarshalBlock(block, true, true); err != nil {
			blockJSON = &ethapi.RPCMarshalBlockT{Error: err.Error()}
		}
		results = append(results, &BadBlockArgs{Hash: block.Hash(), RLP: blockRlp, Block: blockJSON})
	}
	return results, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L357" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getBlockRlp

GetBlockRlp retrieves the RLP encoded for of a single block.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
number <code>uint64</code> 

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





#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_getBlockRlp", "params": [<number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_getBlockRlp", "params": [<number>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1973" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getModifiedAccountsByHash

GetModifiedAccountsByHash returns all accounts that have changed between the
two blocks specified. A change is defined as a difference in nonce, balance,
code hash, or storage hash.

With one parameter, returns the list of accounts modified in the specified block.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
startHash <code>common.Hash</code> 

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
endHash <code>*common.Hash</code> 

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_getModifiedAccountsByHash", "params": [<startHash>, <endHash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_getModifiedAccountsByHash", "params": [<startHash>, <endHash>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L521" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_getModifiedAccountsByNumber

GetModifiedAccountsByNumber returns all accounts that have changed between the
two blocks specified. A change is defined as a difference in nonce, balance,
code hash, or storage hash.

With one parameter, returns the list of accounts modified in the specified block.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
startNum <code>uint64</code> 

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
endNum <code>*uint64</code> 

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_getModifiedAccountsByNumber", "params": [<startNum>, <endNum>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_getModifiedAccountsByNumber", "params": [<startNum>, <endNum>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L493" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_goTrace

GoTrace turns on tracing for nsec seconds and writes
trace data to file.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_goTrace", "params": [<file>, <nsec>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_goTrace", "params": [<file>, <nsec>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L135" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_memStats

MemStats returns detailed runtime memory statistics.


#### Params (0)

_None_

#### Result




<code>*runtime.MemStats</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- Alloc: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- BuckHashSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- BySize: 
			- items: 
				- additionalProperties: `false`
				- properties: 
					- Frees: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- Mallocs: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`

					- Size: 
						- pattern: `^0x[a-fA-F0-9]+$`
						- title: `integer`
						- type: `string`


				- type: `object`

			- maxItems: `61`
			- minItems: `61`
			- type: `array`

		- DebugGC: 
			- type: `boolean`

		- EnableGC: 
			- type: `boolean`

		- Frees: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- GCCPUFraction: 
			- type: `number`

		- GCSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapAlloc: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapIdle: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapObjects: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapReleased: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- HeapSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- LastGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Lookups: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MCacheInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MCacheSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MSpanInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- MSpanSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Mallocs: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NextGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NumForcedGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- NumGC: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- OtherSys: 
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

		- PauseNs: 
			- items: 
				- description: `Hex representation of the integer`
				- pattern: `^0x[a-fA-F0-9]+$`
				- title: `integer`
				- type: `string`

			- maxItems: `256`
			- minItems: `256`
			- type: `array`

		- PauseTotalNs: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- StackInuse: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- StackSys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- Sys: 
			- pattern: `^0x[a-fA-F0-9]+$`
			- title: `integer`
			- type: `string`

		- TotalAlloc: 
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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_memStats", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_memStats", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L73" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_mutexProfile

MutexProfile turns on mutex profiling for nsec seconds and writes profile data to file.
It uses a profile rate of 1 for most accurate information. If a different rate is
desired, set the rate and write the profile manually.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
nsec <code>uint</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_mutexProfile", "params": [<file>, <nsec>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_mutexProfile", "params": [<file>, <nsec>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L168" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_preimage

Preimage is a debug API function that returns the preimage for a sha3 hash, if known.


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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_preimage", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_preimage", "params": [<hash>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L341" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_printBlock

PrintBlock retrieves a block and returns its pretty printed form.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
number <code>uint64</code> 

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





#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_printBlock", "params": [<number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_printBlock", "params": [<number>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2025" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_removePendingTransaction

RemovePendingTransaction removes a transaction from the txpool.
It returns the transaction removed, if any.


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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_removePendingTransaction", "params": [<hash>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_removePendingTransaction", "params": [<hash>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L573" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_seedHash

SeedHash retrieves the seed hash of a block.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
number <code>uint64</code> 

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





#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_seedHash", "params": [<number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_seedHash", "params": [<number>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2034" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setBlockProfileRate

SetBlockProfileRate sets the rate of goroutine block profile data collection.
rate 0 disables block profiling.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
rate <code>int</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_setBlockProfileRate", "params": [<rate>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_setBlockProfileRate", "params": [<rate>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L156" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setGCPercent

SetGCPercent sets the garbage collection target percentage. It returns the previous
setting. A negative value disables GC.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
v <code>int</code> 

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





#### Result




<code>int</code> 

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_setGCPercent", "params": [<v>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_setGCPercent", "params": [<v>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L206" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setHead

SetHead rewinds the head of the blockchain to a previous block.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
number <code>hexutil.Uint64</code> 

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

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_setHead", "params": [<number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_setHead", "params": [<number>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2081" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_setMutexProfileFraction

SetMutexProfileFraction sets the rate of mutex profiling.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
rate <code>int</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_setMutexProfileFraction", "params": [<rate>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_setMutexProfileFraction", "params": [<rate>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L177" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stacks

Stacks returns a printed representation of the stacks of all goroutines.


#### Params (0)

_None_

#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_stacks", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_stacks", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L193" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_standardTraceBadBlockToFile

StandardTraceBadBlockToFile dumps the structured logs created during the
execution of EVM against a block pulled from the pool of bad ones to the
local file system and returns a list of files to the caller.


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

		- TxHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
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





#### Result



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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_standardTraceBadBlockToFile", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_standardTraceBadBlockToFile", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.standardTraceBadBlockToFile(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) StandardTraceBadBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([ // StandardTraceBadBlockToFile dumps the structured logs created during the
// execution of EVM against a block pulled from the pool of bad ones to the
// local file system and returns a list of files to the caller.
]string, error) {
	for _, block := range rawdb.ReadAllBadBlocks(api.backend.ChainDb()) {
		if block.Hash() == hash {
			return api.standardTraceBlockToFile(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L448" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_standardTraceBlockToFile

StandardTraceBlockToFile dumps the structured logs created during the
execution of EVM to the local file system and returns a list of files
to the caller.


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

		- TxHash: 
			- pattern: `^0x[a-fA-F\d]{64}$`
			- title: `keccak`
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





#### Result



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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_standardTraceBlockToFile", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_standardTraceBlockToFile", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.standardTraceBlockToFile(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) StandardTraceBlockToFile(ctx context.Context, hash common.Hash, config *StdTraceConfig) ([ // StandardTraceBlockToFile dumps the structured logs created during the
// execution of EVM to the local file system and returns a list of files
// to the caller.
]string, error) {
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return api.standardTraceBlockToFile(ctx, block, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L437" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_startCPUProfile

StartCPUProfile turns on CPU profiling, writing to the given file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_startCPUProfile", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_startCPUProfile", "params": [<file>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L98" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_startGoTrace

StartGoTrace turns on tracing, writing to the given file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_startGoTrace", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_startGoTrace", "params": [<file>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/trace.go#L30" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stopCPUProfile

StopCPUProfile stops an ongoing CPU profile.


#### Params (0)

_None_

#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_stopCPUProfile", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_stopCPUProfile", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L119" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_stopGoTrace

StopTrace stops an ongoing trace.


#### Params (0)

_None_

#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_stopGoTrace", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_stopGoTrace", "params": []}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/trace.go#L51" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_storageRangeAt

StorageRangeAt returns the storage at the given block height and transaction index.


#### Params (5)

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
txIndex <code>int</code> 

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
contractAddress <code>common.Address</code> 

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




__4:__ 
keyStart <code>hexutil.Bytes</code> 

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




__5:__ 
maxResult <code>int</code> 

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





#### Result




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



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_storageRangeAt", "params": [<blockHash>, <txIndex>, <contractAddress>, <keyStart>, <maxResult>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_storageRangeAt", "params": [<blockHash>, <txIndex>, <contractAddress>, <keyStart>, <maxResult>]}'
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
	_, _, statedb, release, err := api.eth.stateAtTransaction(block, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	defer release()
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}// StorageRangeAt returns the storage at the given block height and transaction index.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api.go#L447" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_subscribe

Subscribe creates a subscription to an event channel.
Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
subscriptionName <code>RPCDebugSubscriptionParamsName</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `Returns transaction traces within a range of blocks.`
			- enum: traceChain
			- type: string


	- title: `subscriptionName`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "Returns transaction traces within a range of blocks.",
                "enum": [
                    "traceChain"
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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_subscribe", "params": [<subscriptionName>, <subscriptionOptions>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (sub *RPCDebugSubscription) Subscribe(subscriptionName RPCDebugSubscriptionParamsName, subscriptionOptions interface{}) (subscriptionID rpc.ID, err error) {
	return
}// Subscribe creates a subscription to an event channel.
// Subscriptions are not available over HTTP; they are only available over WS, IPC, and Process connections.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L252" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_testSignCliqueBlock

TestSignCliqueBlock fetches the given block number, and attempts to sign it as a clique header with the
given address, returning the address of the recovered signature

This is a temporary method to debug the externalsigner integration,
TODO: Remove this method when the integration is mature


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
number <code>uint64</code> 

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_testSignCliqueBlock", "params": [<address>, <number>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_testSignCliqueBlock", "params": [<address>, <number>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L1990" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBadBlock

TraceBadBlock returns the structured logs created during the execution of
EVM against a block pulled from the pool of bad ones and returns them as a JSON
object.


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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceBadBlock", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceBadBlock", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceBadBlock(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceBadBlock(ctx context.Context, hash common.Hash, config *TraceConfig) ([ // TraceBadBlock returns the structured logs created during the execution of
// EVM against a block pulled from the pool of bad ones and returns them as a JSON
// object.
]*txTraceResult, error) {
	for _, block := range rawdb.ReadAllBadBlocks(api.backend.ChainDb()) {
		if block.Hash() == hash {
			return api.traceBlock(ctx, block, config)
		}
	}
	return nil, fmt.Errorf("bad block %#x not found", hash)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L425" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlock

TraceBlock returns the structured logs created during the execution of EVM
and returns them as a JSON object.


#### Params (2)

Parameters must be given _by position_.


__1:__ 
blob <code>[]byte</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of a variable length byte array`
	- pattern: `^0x([a-fA-F0-9]?)+$`
	- title: `bytes`
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceBlock", "params": [<blob>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceBlock", "params": [<blob>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceBlock(blob,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceBlock(ctx context.Context, blob [ // TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
]byte, config *TraceConfig) ([]*txTraceResult, error) {
	block := new(types.Block)
	if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
		return nil, fmt.Errorf("could not decode block: %v", err)
	}
	return api.traceBlock(ctx, block, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L404" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockByHash

TraceBlockByHash returns the structured logs created during the execution of
EVM and returns them as a JSON object.


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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceBlockByHash", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceBlockByHash", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceBlockByHash(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([ // TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	block, err := api.blockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L394" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockByNumber

TraceBlockByNumber returns the structured logs created during the execution of
EVM and returns them as a JSON object.


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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceBlockByNumber", "params": [<number>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceBlockByNumber", "params": [<number>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceBlockByNumber(number,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([ // TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	block, err := api.blockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return api.traceBlock(ctx, block, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L384" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceBlockFromFile

TraceBlockFromFile returns the structured logs created during the execution of
EVM and returns them as a JSON object.


#### Params (2)

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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceBlockFromFile", "params": [<file>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceBlockFromFile", "params": [<file>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceBlockFromFile(file,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([ // TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
]*txTraceResult, error) {
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, blob, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L414" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceCall

TraceCall lets you trace a given eth_call. It collects the structured logs
created during the execution of EVM if the given transaction was added on
top of the provided block and returns them as a JSON object.
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceCall", "params": [<args>, <blockNrOrHash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceCall", "params": [<args>, <blockNrOrHash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceCall(args,blockNrOrHash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceCall(ctx context.Context, args ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	var (
		err	error
		block	*types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = api.blockByHash(ctx, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		block, err = api.blockByNumber(ctx, number)
	}
	if err != nil {
		return nil, err
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, release, err := api.backend.StateAtBlock(ctx, block, reexec)
	if err != nil {
		return nil, err
	}
	defer release()
	msg := args.ToMessage(api.backend.RPCGasCap())
	vmctx := core.NewEVMBlockContext(block.Header(), api.chainContext(ctx), nil)
	originalCanTransfer := vmctx.CanTransfer
	originalTransfer := vmctx.Transfer
	gasCost := new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.GasPrice())
	totalCost := new(big.Int).Add(gasCost, msg.Value())
	hasFromSufficientBalanceForValueAndGasCost := vmctx.CanTransfer(statedb, msg.From(), totalCost)
	hasFromSufficientBalanceForGasCost := vmctx.CanTransfer(statedb, msg.From(), gasCost)
	vmctx.CanTransfer = func(db vm.StateDB, sender common.Address, amount *big.Int) bool {
		if msg.From() == sender {
			return true
		}
		res := originalCanTransfer(db, sender, amount)
		return res
	}
	vmctx.Transfer = func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
		toAmount := new(big.Int).Set(amount)
		senderBalance := db.GetBalance(sender)
		if senderBalance.Cmp(toAmount) < 0 {
			toAmount.Set(big.NewInt(0))
		}
		originalTransfer(db, sender, recipient, toAmount)
	}
	taskExtraContext := map // TraceCall lets you trace a given eth_call. It collects the structured logs
	// created during the execution of EVM if the given transaction was added on
	// top of the provided block and returns them as a JSON object.
	// You can provide -2 as a block number to trace on top of the pending block.
	// Try to retrieve the specified block
	[string]interface{}{"hasFromSufficientBalanceForValueAndGasCost": hasFromSufficientBalanceForValueAndGasCost, "hasFromSufficientBalanceForGasCost": hasFromSufficientBalanceForGasCost}
	return api.traceTx(ctx, msg, vmctx, statedb, taskExtraContext, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L725" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceCallMany

TraceCallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceCallMany", "params": [<txs>, <blockNrOrHash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceCallMany", "params": [<txs>, <blockNrOrHash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceCallMany(txs,blockNrOrHash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceCallMany(ctx context.Context, txs [ // TraceCallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
]ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (interface{}, error) {
	var (
		err	error
		block	*types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = api.blockByHash(ctx, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		block, err = api.blockByNumber(ctx, number)
	}
	if err != nil {
		return nil, err
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, release, err := api.backend.StateAtBlock(ctx, block, reexec)
	if err != nil {
		return nil, err
	}
	defer release()
	var results = make([ // Try to retrieve the specified block
	]interface{}, len(txs))
	for idx, args := range txs {
		msg := args.ToMessage(api.backend.RPCGasCap())
		vmctx := core.NewEVMBlockContext(block.Header(), api.chainContext(ctx), nil)
		originalCanTransfer := vmctx.CanTransfer
		originalTransfer := vmctx.Transfer
		gasCost := new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.GasPrice())
		totalCost := new(big.Int).Add(gasCost, msg.Value())
		hasFromSufficientBalanceForValueAndGasCost := vmctx.CanTransfer(statedb, msg.From(), totalCost)
		hasFromSufficientBalanceForGasCost := vmctx.CanTransfer(statedb, msg.From(), gasCost)
		vmctx.CanTransfer = func(db vm.StateDB, sender common.Address, amount *big.Int) bool {
			if msg.From() == sender {
				return true
			}
			res := originalCanTransfer(db, sender, amount)
			return res
		}
		vmctx.Transfer = func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
			toAmount := new(big.Int).Set(amount)
			senderBalance := db.GetBalance(sender)
			if senderBalance.Cmp(toAmount) < 0 {
				toAmount.Set(big.NewInt(0))
			}
			originalTransfer(db, sender, recipient, toAmount)
		}
		taskExtraContext := map[string]interface{}{"hasFromSufficientBalanceForValueAndGasCost": hasFromSufficientBalanceForValueAndGasCost, "hasFromSufficientBalanceForGasCost": hasFromSufficientBalanceForGasCost}
		res, err := api.traceTx(ctx, msg, vmctx, statedb, taskExtraContext, config)
		if err != nil {
			results[idx] = &txTraceResult{Error: err.Error()}
			continue
		}
		res, err = decorateResponse(res, config)
		if err != nil {
			return nil, fmt.Errorf("failed to decorate response for transaction at index %d with error %v", idx, err)
		}
		results[idx] = res
	}
	return results, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L798" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceChain

TraceChain returns the structured logs created during the execution of EVM
between two blocks (excluding start) and returns them as a JSON object.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
start <code>rpc.BlockNumber</code> 

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
end <code>rpc.BlockNumber</code> 

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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_subscribe", "params": ["traceChain", <start>, <end>, <config>]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	from, err := api.blockByNumber(ctx, start)
	if err != nil {
		return nil, err
	}
	to, err := api.blockByNumber(ctx, end)
	if err != nil {
		return nil, err
	}
	if from.Number().Cmp(to.Number()) >= 0 {
		return nil, fmt.Errorf("end block (#%d) needs to come after start block (#%d)", end, start)
	}
	return api.traceChain(ctx, from, to, config)
}// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L207" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_traceTransaction

TraceTransaction returns the structured logs created during the execution of EVM
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_traceTransaction", "params": [<hash>, <config>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_traceTransaction", "params": [<hash>, <config>]}'
	```


=== "Javascript Console"

	``` js
	debug.traceTransaction(hash,config);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *API) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	tx, blockHash, blockNumber, index, err := api.backend.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if blockNumber == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	block, err := api.blockByNumberAndHash(ctx, rpc.BlockNumber(blockNumber), blockHash)
	if err != nil {
		return nil, err
	}
	msg, vmctx, statedb, release, err := api.backend.StateAtTransaction(ctx, block, int(index), reexec)
	if err != nil {
		return nil, err
	}
	defer release()
	taskExtraContext := map // TraceTransaction returns the structured logs created during the execution of EVM
	// and returns them as a JSON object.
	[string]interface{}{"blockNumber": block.NumberU64(), "blockHash": blockHash.Hex(), "transactionHash": tx.Hash().Hex(), "transactionPosition": index}
	return api.traceTx(ctx, msg, vmctx, statedb, taskExtraContext, config)
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/tracers/api.go#L688" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_unsubscribe

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_unsubscribe", "params": [<id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_unsubscribe", "params": [<id>]}'
	```


=== "Javascript Console"

	``` js
	debug.unsubscribe(id);
	```



<details><summary>Source code</summary>
<p>
```go
func (sub *RPCDebugSubscription) Unsubscribe(id rpc.ID) error {
	return nil
}// Unsubscribe terminates an existing subscription by ID.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/openrpc.go#L243" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_verbosity

Verbosity sets the log verbosity ceiling. The verbosity of individual packages
and source files can be raised using Vmodule.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
level <code>int</code> 

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





#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_verbosity", "params": [<level>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_verbosity", "params": [<level>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L57" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_vmodule

Vmodule sets the log verbosity pattern. See package log for details on the
pattern syntax.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
pattern <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_vmodule", "params": [<pattern>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_vmodule", "params": [<pattern>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L62" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeBlockProfile

WriteBlockProfile writes a goroutine blocking profile to the given file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_writeBlockProfile", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_writeBlockProfile", "params": [<file>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L161" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeMemProfile

WriteMemProfile writes an allocation profile to the given file.
Note that the profiling rate cannot be set through the API,
it must be set on the command line.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_writeMemProfile", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_writeMemProfile", "params": [<file>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L188" target="_">View on GitHub →</a>
</p>
</details>

---



### debug_writeMutexProfile

WriteMutexProfile writes a goroutine blocking profile to the given file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result

_None_

#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "debug_writeMutexProfile", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "debug_writeMutexProfile", "params": [<file>]}'
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
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/debug/api.go#L181" target="_">View on GitHub →</a>
</p>
</details>

---

