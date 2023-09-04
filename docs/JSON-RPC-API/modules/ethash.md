






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### ethash_getHashrate

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "ethash_getHashrate", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "ethash_getHashrate", "params": []}'
	```


=== "Javascript Console"

	``` js
	ethash.getHashrate();
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



### ethash_getWork

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "ethash_getWork", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "ethash_getWork", "params": []}'
	```


=== "Javascript Console"

	``` js
	ethash.getWork();
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



### ethash_submitHashrate

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "ethash_submitHashrate", "params": [<rate>, <id>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "ethash_submitHashrate", "params": [<rate>, <id>]}'
	```


=== "Javascript Console"

	``` js
	ethash.submitHashrate(rate,id);
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



### ethash_submitWork

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
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "ethash_submitWork", "params": [<nonce>, <hash>, <digest>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "ethash_submitWork", "params": [<nonce>, <hash>, <digest>]}'
	```


=== "Javascript Console"

	``` js
	ethash.submitWork(nonce,hash,digest);
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

