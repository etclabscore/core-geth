






| Entity | Version |
| --- | --- |
| Source | <code>1.11.22-unstable/generated-at:2021-01-23T04:50:40-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### miner_getHashrate

GetHashrate returns the current hashrate of the miner.


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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_getHashrate", "params": []}'
	```

=== "Javascript Console"

	``` js
	miner.getHashrate();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) GetHashrate() uint64 {
	return api.e.miner.HashRate()
}// GetHashrate returns the current hashrate of the miner.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L153" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_setEtherbase

SetEtherbase sets the etherbase of the miner


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
etherbase <code>common.Address</code> 

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





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_setEtherbase", "params": [<etherbase>]}'
	```

=== "Javascript Console"

	``` js
	miner.setEtherbase(etherbase);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) SetEtherbase(etherbase common.Address) bool {
	api.e.SetEtherbase(etherbase)
	return true
}// SetEtherbase sets the etherbase of the miner

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L142" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_setExtra

SetExtra sets the extra data string that is included when this miner mines a block.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
extra <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_setExtra", "params": [<extra>]}'
	```

=== "Javascript Console"

	``` js
	miner.setExtra(extra);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([ // SetExtra sets the extra data string that is included when this miner mines a block.
	]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}
```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L124" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_setGasPrice

SetGasPrice sets the minimum accepted gas price for the miner.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
gasPrice <code>hexutil.Big</code> 

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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_setGasPrice", "params": [<gasPrice>]}'
	```

=== "Javascript Console"

	``` js
	miner.setGasPrice(gasPrice);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()
	api.e.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}// SetGasPrice sets the minimum accepted gas price for the miner.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L132" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_setRecommitInterval

SetRecommitInterval updates the interval for miner sealing work recommitting.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
interval <code>int</code> 

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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_setRecommitInterval", "params": [<interval>]}'
	```

=== "Javascript Console"

	``` js
	miner.setRecommitInterval(interval);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) SetRecommitInterval(interval int) {
	api.e.Miner().SetRecommitInterval(time.Duration(interval) * time.Millisecond)
}// SetRecommitInterval updates the interval for miner sealing work recommitting.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L148" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_start

Start starts the miner with the given number of threads. If threads is nil,
the number of workers started is equal to the number of logical CPUs that are
usable by this process. If mining is already running, this method adjust the
number of threads allowed to use and updates the minimum price required by the
transaction pool.


#### Params (1)

Parameters must be given _by position_.  


__1:__ 
threads <code>*int</code> 

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

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_start", "params": [<threads>]}'
	```

=== "Javascript Console"

	``` js
	miner.start(threads);
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) Start(threads *int) error {
	if threads == nil {
		return api.e.StartMining(runtime.NumCPU())
	}
	return api.e.StartMining(*threads)
}// Start starts the miner with the given number of threads. If threads is nil,
// the number of workers started is equal to the number of logical CPUs that are
// usable by this process. If mining is already running, this method adjust the
// number of threads allowed to use and updates the minimum price required by the
// transaction pool.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L110" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_stop

Stop terminates the miner, both at the consensus engine level as well as at
the block creation level.


#### Params (0)

_None_

#### Result

_None_

#### Client Method Invocation Examples

=== "Shell"

	``` shell
	curl -X POST http://localhost:8545 --data '{"jsonrpc": "2.0", id": 42, "method": "miner_stop", "params": []}'
	```

=== "Javascript Console"

	``` js
	miner.stop();
	```


<details><summary>Source code</summary>
<p>
```go
func (api *PrivateMinerAPI) Stop() {
	api.e.StopMining()
}// Stop terminates the miner, both at the consensus engine level as well as at
// the block creation level.

```
<a href="https://github.com/ethereum/go-ethereum/blob/master/eth/api.go#L119" target="_">View on GitHub →</a>
</p>
</details>

---

