






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-21T17:27:32-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### web3_clientVersion

ClientVersion returns the node name


#### Params (0)

_None_

#### Result




<code>string</code> 

  + Required: ✓ Yes




__Client Method Invocation Examples__

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "web3_clientVersion", "params": []}'
	```

=== "Javascript Console"

	``` js
	web3.clientVersion();
	```


<details><summary>Source code</summary>
<p>
```go
func (s *publicWeb3API) ClientVersion() string {
	return s.stack.Server().Name
}// ClientVersion returns the node name

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/node/api.go#L320" target="_">View on GitHub →</a>
</p>
</details>

---



### web3_sha3

Sha3 applies the ethereum sha3 implementation on the input.
It assumes the input is hex encoded.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
input <code>hexutil.Bytes</code> 

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
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "web3_sha3", "params": []}'
	```

=== "Javascript Console"

	``` js
	web3.sha3(input);
	```


<details><summary>Source code</summary>
<p>
```go
func (s *publicWeb3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}// Sha3 applies the ethereum sha3 implementation on the input.
// It assumes the input is hex encoded.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/node/api.go#L326" target="_">View on GitHub →</a>
</p>
</details>

---

