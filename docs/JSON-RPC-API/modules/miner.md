






| Entity | Version |
| --- | --- |
| Source | <code>1.12.14-unstable/generated-at:2023-09-04T08:02:34-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### miner_setEtherbase

SetEtherbase sets the etherbase of the miner.


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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_setEtherbase", "params": [<etherbase>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_setEtherbase", "params": [<etherbase>]}'
	```


=== "Javascript Console"

	``` js
	miner.setEtherbase(etherbase);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) SetEtherbase(etherbase common.Address) bool {
	api.e.SetEtherbase(etherbase)
	return true
}// SetEtherbase sets the etherbase of the miner.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L81" target="_">View on GitHub →</a>
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_setExtra", "params": [<extra>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_setExtra", "params": [<extra>]}'
	```


=== "Javascript Console"

	``` js
	miner.setExtra(extra);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([ // SetExtra sets the extra data string that is included when this miner mines a block.
	]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}
```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L57" target="_">View on GitHub →</a>
</p>
</details>

---



### miner_setGasLimit

SetGasLimit sets the gaslimit to target towards during mining.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
gasLimit <code>hexutil.Uint64</code> 

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




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_setGasLimit", "params": [<gasLimit>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_setGasLimit", "params": [<gasLimit>]}'
	```


=== "Javascript Console"

	``` js
	miner.setGasLimit(gasLimit);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) SetGasLimit(gasLimit hexutil.Uint64) bool {
	api.e.Miner().SetGasCeil(uint64(gasLimit))
	return true
}// SetGasLimit sets the gaslimit to target towards during mining.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L75" target="_">View on GitHub →</a>
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_setGasPrice", "params": [<gasPrice>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_setGasPrice", "params": [<gasPrice>]}'
	```


=== "Javascript Console"

	``` js
	miner.setGasPrice(gasPrice);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()
	api.e.txPool.SetGasTip((*big.Int)(&gasPrice))
	return true
}// SetGasPrice sets the minimum accepted gas price for the miner.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L65" target="_">View on GitHub →</a>
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_setRecommitInterval", "params": [<interval>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_setRecommitInterval", "params": [<interval>]}'
	```


=== "Javascript Console"

	``` js
	miner.setRecommitInterval(interval);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) SetRecommitInterval(interval int) {
	api.e.Miner().SetRecommitInterval(time.Duration(interval) * time.Millisecond)
}// SetRecommitInterval updates the interval for miner sealing work recommitting.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L87" target="_">View on GitHub →</a>
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_start", "params": [<threads>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_start", "params": [<threads>]}'
	```


=== "Javascript Console"

	``` js
	miner.start(threads);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) Start(threads *int) error {
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
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L43" target="_">View on GitHub →</a>
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


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "miner_stop", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "miner_stop", "params": []}'
	```


=== "Javascript Console"

	``` js
	miner.stop();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *MinerAPI) Stop() {
	api.e.StopMining()
}// Stop terminates the miner, both at the consensus engine level as well as at
// the block creation level.

```
<a href="https://github.com/etclabscore/core-geth/blob/master/eth/api_miner.go#L52" target="_">View on GitHub →</a>
</p>
</details>

---

