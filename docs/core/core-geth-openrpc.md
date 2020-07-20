# Core-geth and Open-RPC

Since ~~September 2019~~ Core-geth has implemented support for the Open-RPC service description endpoint
`rpc_discover`. Querying this endpoint returns the user a service definition document defined and
intended to be standardized at github.com/etclabscore/ethereum-json-rpc-specification.

The definition there hopes to establish and codify the common ground between existing Ethereum Protocol provider
implemenations (eg. ethereum/go-ethereum, etclabscore/core-geth, openethereum/openethereum), and in doing so,
build a target specification for new providers to match when building their JSON-RPC APIs. 

The document details about 45 methods and their associated paramaters and results, defined as JSON Schemas.

With a little weekend-hacking, it's very possible to build some super-convenient tools like the one I use below, `ethrpc`:

```
> ethrpc | head -20
> This is an auto-generated CLI interface for an Open-RPC compliant API.
>
> Open-RPC Version: 1.0.10
>
> Run 'ethrpc completion --help' to learn about auto-auto-completion! It's easy!
>
> Usage:
>   ethrpc [command]
>
>   Available Commands:
>     completion                                 Generates bash completion scripts
>       eth_blockNumber                            Returns the number of most recent block.
>         eth_call                                   Executes a new message call (locally) immediately without creating a transaction on the block chain.
>           eth_chainId                                Returns the currently configured chain id
>             eth_coinbase                               Returns the client coinbase address.
>               eth_estimateGas                            Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.
>                 eth_gasPrice                               Returns the current price per gas in wei
>                   eth_getBalance                             Returns Ether balance of a given or account or contract
>                     eth_getBlockByHash                         Gets a block for a given hash
>                       eth_getBlockByNumber                       Gets a block for a given number salad
>
```
