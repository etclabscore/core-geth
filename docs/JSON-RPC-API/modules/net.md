






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-21T17:27:32-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### net_listening

Listening returns an indication if the node is listening for network connections.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "net_listening", "params": []}'
	```

=== "Javascript Console"

	``` js
	net.listening();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicNetAPI) Listening() bool {
	return true
}// Listening returns an indication if the node is listening for network connections.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L2033" target="_">View on GitHub →</a>
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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "net_peerCount", "params": []}'
	```

=== "Javascript Console"

	``` js
	net.peerCount();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicNetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(s.net.PeerCount())
}// PeerCount returns the number of connected peers

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L2037" target="_">View on GitHub →</a>
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




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "net_version", "params": []}'
	```

=== "Javascript Console"

	``` js
	net.version();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}// Version returns the current ethereum protocol version.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L2042" target="_">View on GitHub →</a>
</p>
</details>

---

