






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### admin_addPeer

AddPeer requests connecting to a remote node, and also maintaining the new
connection at all times, even reconnecting if it is lost.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_addPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_addPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.addPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) AddPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.AddPeer(node)
	return true, nil
}// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L57" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_addTrustedPeer

AddTrustedPeer allows a remote node to always connect, even if slots are full


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_addTrustedPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_addTrustedPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.addTrustedPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) AddTrustedPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.AddTrustedPeer(node)
	return true, nil
}// AddTrustedPeer allows a remote node to always connect, even if slots are full

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L89" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_datadir

Datadir retrieves the current data directory the node is using.


#### Params (0)

_None_

#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_datadir", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_datadir", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.datadir();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) Datadir() string {
	return api.node.DataDir()
}// Datadir retrieves the current data directory the node is using.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L320" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_ecbp1100



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




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_ecbp1100", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_ecbp1100", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	admin.ecbp1100(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Ecbp1100(blockNr rpc.BlockNumber) (bool, error) {
	i := uint64(blockNr.Int64())
	err := api.eth.blockchain.Config().SetECBP1100Transition(&i)
	return api.eth.blockchain.IsArtificialFinalityEnabled() && api.eth.blockchain.Config().IsEnabled(api.eth.blockchain.Config().GetECBP1100Transition, api.eth.blockchain.CurrentBlock().Number), err
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_admin.go#L142" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_exportChain

ExportChain exports the current blockchain into a local file,
or a range of blocks if first and last are non-nil.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
first <code>*uint64</code> 

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
last <code>*uint64</code> 

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




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_exportChain", "params": [<file>, <first>, <last>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_exportChain", "params": [<file>, <first>, <last>]}'
	```


=== "Javascript Console"

	``` js
	admin.exportChain(file,first,last);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first != nil && last == nil {
		head := api.eth.BlockChain().CurrentHeader().Number.Uint64()
		last = &head
	}
	if _, err := os.Stat(file); err == nil {
		return false, errors.New("location would overwrite an existing file")
	}
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false, err
	}
	defer out.Close()
	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}
	if first != nil {
		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
			return false, err
		}
	} else if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_admin.go#L46" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_importChain

ImportChain imports a blockchain from a local file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_importChain", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_importChain", "params": [<file>]}'
	```


=== "Javascript Console"

	``` js
	admin.importChain(file);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) ImportChain(file string) (bool, error) {
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()
	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}
	stream := rlp.NewStream(reader, 0)
	blocks, index := make([ // ImportChain imports a blockchain from a local file.
	]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}
		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_admin.go#L94" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_maxPeers

MaxPeers sets the maximum peer limit for the protocol manager and the p2p server.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
n <code>int</code> 

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




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_maxPeers", "params": [<n>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_maxPeers", "params": [<n>]}'
	```


=== "Javascript Console"

	``` js
	admin.maxPeers(n);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MaxPeers(n int) (bool, error) {
	api.eth.handler.maxPeers = n
	api.eth.p2pServer.MaxPeers = n
	for i := api.eth.handler.peers.len(); i > n; i = api.eth.handler.peers.len() {
		p := api.eth.handler.peers.WorstPeer()
		if p == nil {
			break
		}
		api.eth.handler.removePeer(p.ID())
	}
	return true, nil
}// MaxPeers sets the maximum peer limit for the protocol manager and the p2p server.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_admin.go#L152" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_nodeInfo

NodeInfo retrieves all the information we know about the host node at the
protocol granularity.


#### Params (0)

_None_

#### Result




<code>*p2p.NodeInfo</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- enode: 
			- type: `string`

		- enr: 
			- type: `string`

		- id: 
			- type: `string`

		- ip: 
			- type: `string`

		- listenAddr: 
			- type: `string`

		- name: 
			- type: `string`

		- ports: 
			- additionalProperties: `false`
			- properties: 
				- discovery: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- listener: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`


			- type: `object`

		- protocols: 
			- additionalProperties: `false`
			- properties: 
				- discovery: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- listener: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "enode": {
                "type": "string"
            },
            "enr": {
                "type": "string"
            },
            "id": {
                "type": "string"
            },
            "ip": {
                "type": "string"
            },
            "listenAddr": {
                "type": "string"
            },
            "name": {
                "type": "string"
            },
            "ports": {
                "additionalProperties": false,
                "properties": {
                    "discovery": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "listener": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    }
                },
                "type": "object"
            },
            "protocols": {
                "additionalProperties": false,
                "properties": {
                    "discovery": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "listener": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_nodeInfo", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_nodeInfo", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.nodeInfo();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	return server.NodeInfo(), nil
}// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L310" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_peerEvents

PeerEvents creates an RPC subscription which receives peer events from the
node's p2p.Server


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
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_subscribe", "params": ["peerEvents"]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) PeerEvents(ctx context.Context) (*rpc.Subscription, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		events := make(chan *p2p.PeerEvent)
		sub := server.SubscribeEvents(events)
		defer sub.Unsubscribe()
		for {
			select {
			case event := <-events:
				notifier.Notify(rpcSub.ID, event)
			case <-sub.Err():
				return
			case <-rpcSub.Err():
				return
			case <-notifier.Closed():
				return
			}
		}
	}()
	return rpcSub, nil
}// PeerEvents creates an RPC subscription which receives peer events from the
// node's p2p.Server

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L121" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_peers

Peers retrieves all the information we know about each individual peer at the
protocol granularity.


#### Params (0)

_None_

#### Result



p2pPeerInfo <code>[]*p2p.PeerInfo</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- caps: 
					- items: 
						- type: `string`

					- type: `array`

				- enode: 
					- type: `string`

				- enr: 
					- type: `string`

				- id: 
					- type: `string`

				- name: 
					- type: `string`

				- network: 
					- additionalProperties: `false`
					- properties: 
						- inbound: 
							- type: `boolean`

						- localAddress: 
							- type: `string`

						- remoteAddress: 
							- type: `string`

						- static: 
							- type: `boolean`

						- trusted: 
							- type: `boolean`


					- type: `object`

				- protocols: 
					- additionalProperties: `false`
					- properties: 
						- inbound: 
							- type: `boolean`

						- localAddress: 
							- type: `string`

						- remoteAddress: 
							- type: `string`

						- static: 
							- type: `boolean`

						- trusted: 
							- type: `boolean`


					- type: `object`


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
                    "caps": {
                        "items": {
                            "type": "string"
                        },
                        "type": "array"
                    },
                    "enode": {
                        "type": "string"
                    },
                    "enr": {
                        "type": "string"
                    },
                    "id": {
                        "type": "string"
                    },
                    "name": {
                        "type": "string"
                    },
                    "network": {
                        "additionalProperties": false,
                        "properties": {
                            "inbound": {
                                "type": "boolean"
                            },
                            "localAddress": {
                                "type": "string"
                            },
                            "remoteAddress": {
                                "type": "string"
                            },
                            "static": {
                                "type": "boolean"
                            },
                            "trusted": {
                                "type": "boolean"
                            }
                        },
                        "type": "object"
                    },
                    "protocols": {
                        "additionalProperties": false,
                        "properties": {
                            "inbound": {
                                "type": "boolean"
                            },
                            "localAddress": {
                                "type": "string"
                            },
                            "remoteAddress": {
                                "type": "string"
                            },
                            "static": {
                                "type": "boolean"
                            },
                            "trusted": {
                                "type": "boolean"
                            }
                        },
                        "type": "object"
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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_peers", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_peers", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.peers();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) Peers() ([ // Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
]*p2p.PeerInfo, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	return server.PeersInfo(), nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L300" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_removePeer

RemovePeer disconnects from a remote node if the connection exists


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_removePeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_removePeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.removePeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) RemovePeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.RemovePeer(node)
	return true, nil
}// RemovePeer disconnects from a remote node if the connection exists

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L73" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_removeTrustedPeer

RemoveTrustedPeer removes a remote node from the trusted peer set, but it
does not disconnect it automatically.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_removeTrustedPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_removeTrustedPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.removeTrustedPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) RemoveTrustedPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.RemoveTrustedPeer(node)
	return true, nil
}// RemoveTrustedPeer removes a remote node from the trusted peer set, but it
// does not disconnect it automatically.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L105" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startHTTP

StartHTTP starts the HTTP RPC API server.


#### Params (5)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

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
cors <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes





__5:__ 
vhosts <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startHTTP", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startHTTP", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```


=== "Javascript Console"

	``` js
	admin.startHTTP(host,port,cors,apis,vhosts);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartHTTP(host *string, port *int, cors *string, apis *string, vhosts *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()
	if host == nil {
		h := DefaultHTTPHost
		if api.node.config.HTTPHost != "" {
			h = api.node.config.HTTPHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.config.HTTPPort
	}
	config := httpConfig{CorsAllowedOrigins: api.node.config.HTTPCors, Vhosts: api.node.config.HTTPVirtualHosts, Modules: api.node.config.HTTPModules, rpcEndpointConfig: rpcEndpointConfig{batchItemLimit: api.node.config.BatchRequestLimit, batchResponseSizeLimit: api.node.config.BatchResponseMaxSize}}
	if cors != nil {
		config.CorsAllowedOrigins = nil
		for _, origin := // StartHTTP starts the HTTP RPC API server.
		range strings.Split(*cors, ",") {
			config.CorsAllowedOrigins = append(config.CorsAllowedOrigins, strings.TrimSpace(origin))
		}
	}
	if vhosts != nil {
		config.Vhosts = nil
		for _, vhost := range strings.Split(*host, ",") {
			config.Vhosts = append(config.Vhosts, strings.TrimSpace(vhost))
		}
	}
	if apis != nil {
		config.Modules = nil
		for _, m := range strings.Split(*apis, ",") {
			config.Modules = append(config.Modules, strings.TrimSpace(m))
		}
	}
	if err := api.node.http.setListenAddr(*host, *port); err != nil {
		return false, err
	}
	if err := api.node.http.enableRPC(api.node.rpcAPIs, config); err != nil {
		return false, err
	}
	if err := api.node.http.start(); err != nil {
		return false, err
	}
	return true, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L158" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startRPC

StartRPC starts the HTTP RPC API server.
Deprecated: use StartHTTP instead.


#### Params (5)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

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
cors <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes





__5:__ 
vhosts <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startRPC", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startRPC", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```


=== "Javascript Console"

	``` js
	admin.startRPC(host,port,cors,apis,vhosts);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartRPC(host *string, port *int, cors *string, apis *string, vhosts *string) (bool, error) {
	log.Warn("Deprecation warning", "method", "admin.StartRPC", "use-instead", "admin.StartHTTP")
	return api.StartHTTP(host, port, cors, apis, vhosts)
}// StartRPC starts the HTTP RPC API server.
// Deprecated: use StartHTTP instead.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L217" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startWS

StartWS starts the websocket RPC API server.


#### Params (4)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

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
allowedOrigins <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startWS", "params": [<host>, <port>, <allowedOrigins>, <apis>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startWS", "params": [<host>, <port>, <allowedOrigins>, <apis>]}'
	```


=== "Javascript Console"

	``` js
	admin.startWS(host,port,allowedOrigins,apis);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartWS(host *string, port *int, allowedOrigins *string, apis *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()
	if host == nil {
		h := DefaultWSHost
		if api.node.config.WSHost != "" {
			h = api.node.config.WSHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.config.WSPort
	}
	config := wsConfig{Modules: api.node.config.WSModules, Origins: api.node.config.WSOrigins, rpcEndpointConfig: rpcEndpointConfig{batchItemLimit: api.node.config.BatchRequestLimit, batchResponseSizeLimit: api.node.config.BatchResponseMaxSize}}
	if apis != nil {
		config.Modules = nil
		for _, m := // StartWS starts the websocket RPC API server.
		range strings.Split(*apis, ",") {
			config.Modules = append(config.Modules, strings.TrimSpace(m))
		}
	}
	if allowedOrigins != nil {
		config.Origins = nil
		for _, origin := range strings.Split(*allowedOrigins, ",") {
			config.Origins = append(config.Origins, strings.TrimSpace(origin))
		}
	}
	server := api.node.wsServerForPort(*port, false)
	if err := server.setListenAddr(*host, *port); err != nil {
		return false, err
	}
	openApis, _ := api.node.getAPIs()
	if err := server.enableWS(openApis, config); err != nil {
		return false, err
	}
	if err := server.start(); err != nil {
		return false, err
	}
	api.node.http.log.Info("WebSocket endpoint opened", "url", api.node.WSEndpoint())
	return true, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L236" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopHTTP

StopHTTP shuts down the HTTP server.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopHTTP", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopHTTP", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopHTTP();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopHTTP() (bool, error) {
	api.node.http.stop()
	return true, nil
}// StopHTTP shuts down the HTTP server.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L223" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopRPC

StopRPC shuts down the HTTP server.
Deprecated: use StopHTTP instead.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopRPC", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopRPC", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopRPC();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopRPC() (bool, error) {
	log.Warn("Deprecation warning", "method", "admin.StopRPC", "use-instead", "admin.StopHTTP")
	return api.StopHTTP()
}// StopRPC shuts down the HTTP server.
// Deprecated: use StopHTTP instead.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L230" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopWS

StopWS terminates all WebSocket servers.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopWS", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopWS", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopWS();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopWS() (bool, error) {
	api.node.http.stopWS()
	api.node.ws.stop()
	return true, nil
}// StopWS terminates all WebSocket servers.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/node/api.go#L292" target="_">View on GitHub →</a>
</p>
</details>

---

