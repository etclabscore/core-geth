# Core-Geth and Open-RPC

Since January 2020 Core-Geth has included support for the [Open-RPC](https://spec.open-rpc.org/) service description endpoint
`rpc_discover`. Querying this endpoint returns the user a service definition document detailed at [etclabscore/ethereum-json-rpc-specification](https://github.com/etclabscore/ethereum-json-rpc-specification).

This specification hopes to establish and codify the common ground between Ethereum Protocol provider
implementations (eg. ethereum/go-ethereum, etclabscore/Core-Geth, openethereum/openethereum), and in doing so,
establish a target specification for new providers to match when crafting compliant JSON-RPC APIs.

The document details 44 methods and their associated paramaters and results, defined as JSON Schemas. I find this document best viewed at the Open-RPC Playground, [here](https://playground.open-rpc.org/?schemaUrl=https://raw.githubusercontent.com/etclabscore/ethereum-json-rpc-specification/master/openrpc.json&uiSchema%5BappBar%5D%5Bui:input%5D=false). 

Let's take Core-Geth and open-rpc for a spin.

First, we'll start an ephemeral Core-Geth `geth` instance with `--http` JSON-RPC enabled. Then we'll `curl` the `rpc_discover` method and take a look at the response.

```sh
> ./build/bin/geth --dev --http console 2>/dev/null
```

```sh
> curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc": "2.0","method":"rpc_discover","params":[],"id":71}' http://localhost:8545 > openrpc.json
```

```sh
> cat openrpc.json | jq -r '.result' | head -20
{
  "openrpc": "1.0.0",
  "info": {
    "description": "This API lets you interact with an EVM-based client via JSON-RPC",
    "license": {
      "name": "Apache 2.0",
      "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "title": "Ethereum JSON-RPC",
    "version": "1.0.10"
  },
  "servers": [],
  "methods": [
    {
      "description": "Returns the version of the current client",
      "name": "web3_clientVersion",
      "params": [],
      "result": {
        "description": "client version",
        "name": "clientVersion",
```

```sh
> cat openrpc.json | jq '.result.methods[].name'
"web3_clientVersion"
"web3_sha3"
"net_listening"
"net_peerCount"
"net_version"
"eth_blockNumber"
"eth_call"
"eth_chainId"
"eth_coinbase"
"eth_estimateGas"
"eth_gasPrice"
"eth_getBalance"
"eth_getBlockByHash"
"eth_getBlockByNumber"
"eth_getBlockTransactionCountByHash"
"eth_getBlockTransactionCountByNumber"
"eth_getCode"
"eth_getFilterChanges"
"eth_getFilterLogs"
"eth_getRawTransactionByHash"
"eth_getRawTransactionByBlockHashAndIndex"
"eth_getRawTransactionByBlockNumberAndIndex"
"eth_getLogs"
"eth_getStorageAt"
"eth_getTransactionByBlockHashAndIndex"
"eth_getTransactionByBlockNumberAndIndex"
"eth_getTransactionByHash"
"eth_getTransactionCount"
"eth_getTransactionReceipt"
"eth_getUncleByBlockHashAndIndex"
"eth_getUncleByBlockNumberAndIndex"
"eth_getUncleCountByBlockHash"
"eth_getUncleCountByBlockNumber"
"eth_getProof"
"eth_hashrate"
"eth_mining"
"eth_newBlockFilter"
"eth_newFilter"
"eth_newPendingTransactionFilter"
"eth_pendingTransactions"
"eth_protocolVersion"
"eth_sendRawTransaction"
"eth_syncing"
"eth_uninstallFilter"
```

With this information and a little weekend hacking, it's possible to build some convenient tools like the one I made and use below, `ethrpc`:

```
> ethrpc | head -20
This is an auto-generated CLI interface for an Open-RPC compliant API.

Open-RPC Version: 1.0.10

Run 'ethrpc completion --help' to learn about auto-auto-completion! It's easy!

Usage:
  ethrpc [command]

  Available Commands:
    completion                                 Generates bash completion scripts
    eth_blockNumber                            Returns the number of most recent block.
    eth_call                                   Executes a new message call (locally) immediately without creating a transaction on the block chain.
    eth_chainId                                Returns the currently configured chain id
    eth_coinbase                               Returns the client coinbase address.
    eth_estimateGas                            Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.
    eth_gasPrice                               Returns the current price per gas in wei
    eth_getBalance                             Returns Ether balance of a given or account or contract
    eth_getBlockByHash                         Gets a block for a given hash
    eth_getBlockByNumber                       Gets a block for a given number salad

  [...]
```

It's not very fancy, but was a good exercise in working with the API definition document, and saves me a few quote-laden `curl`s. It has auto-complete and parameter validation, both of which I get essentially for free via, again, the Open-RPC API definition.


---

## On the Future of Core-Geth and Open-RPC

As it's currently implemented, Core-Geth provides the Open-RPC API definition by handing the consumer a static document. This is the simplest possible way to get the document from source to query, and -- given an existing and adhered-to specification for this document (mentioned above) -- it works.

But serving a static document as a definition of an API has limits and risks, too. 
- The document is static; if your API is able to toggle endpoint available (as Core-Geth can with `--http.api`), then it's possible for the document to include endpoints which are not actually available; likewise, Core-Geth has about 100 total methods, with only a little less than half of them documented in the service description (eg. `debug_` and `admin_` -prefixed endpoints, since these are not specificed to be standard across Ethereum Protocol Provider JSON-RPC APIs). 
- Another potential limitation/risk is correctness; since the document had to have-been hand-rolled, there's a possibility of non-systemic bugs in the definition. Tools like validators and control-test suites can serve to minimize their probability, but as with anything human-determined, we shouldn't be 100% confident of correctness. 
- Finally, although we don't expect APIs to change (often, or at all) once established, it does sometimes happen, and when it does, both the code and the corresponding documentation must be modified cooperatively since they're essentially independent mechanisms.

In service of addressing these limits Core-Geth intends to supersede the static-document paradigm with a dynamically generated service description provided by the API itself via both build and runtime introspection.

Work on this project at Core-Geth can be tracked at [this PR](https://github.com/etclabscore/Core-Geth/pull/137), which makes use of a new library [etclabscore/go-openrpc-reflect](https://github.com/etclabscore/go-openrpc-reflect).

