






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### net_listening

Listening returns an indication if the node is listening for network connections.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "net_listening", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "net_listening", "params": []}'
	```


=== "Javascript Console"

	``` js
	net.listening();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *NetAPI) Listening() bool {
	return true
}// Listening returns an indication if the node is listening for network connections.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2313" target="_">View on GitHub →</a>
</p>
</details>

---



### net_peerCount

PeerCount returns the number of connected peers


#### Params (0)

_None_

#### Result




<code>hexutil.Uint</code> 

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "net_peerCount", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "net_peerCount", "params": []}'
	```


=== "Javascript Console"

	``` js
	net.peerCount();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *NetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(s.net.PeerCount())
}// PeerCount returns the number of connected peers

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2317" target="_">View on GitHub →</a>
</p>
</details>

---



### net_version

Version returns the current ethereum protocol version.


#### Params (0)

_None_

#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "net_version", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "net_version", "params": []}'
	```


=== "Javascript Console"

	``` js
	net.version();
	```



<details><summary>Source code</summary>
<p>
```go
func (s *NetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}// Version returns the current ethereum protocol version.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/internal/ethapi/api.go#L2322" target="_">View on GitHub →</a>
</p>
</details>

---

