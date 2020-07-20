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

It's not that fancy, but was a good exercise in working with the API definition document, and saves me a few quote-laden `curl`s. It has auto-complete and parameter validation, both of which I get essentially for free via, again, the Open-RPC API definition.


---

## On the Future of Core-geth and Open-RPC.

As it's currently implemented, core-geth provides the Open-RPC API definition by handing the consumer a static document. This is the simplest possible way to get the document from source to query, and -- given an existing and adhered-to specification for this document (mentioned above) -- it works.

But serving a static document as a definition of an API has limits and risks, too. The document is static; if your API is able to toggle endpoint available (as core-geth can with `--http.api`), then it's possible for the document to include endpoints which are not actually available; likewise, core-geth has about 100 total methods, with only a little less than half of them documented in the service description (eg. `debug_` and `admin_` -prefixed endpoints, since these are not specificed to be standard across Ethereum Protocol Provider JSON-RPC APIs). Another potential limitation/risk is correctness; since the document had to have-been hand-rolled, there's a possibility of non-systemic bugs in the definition. Tools like validators and control-test suites can serve to minimize their probability, but as with anything human-determined, we shouldn't be 100% confident of correctness. Finally, although we don't often expect APIs to change once established (as core-geth's is), it does sometimes happen, and when it does, both the code and the corresponding documentation must be modified cooperatively since they're essentially independent mechanisms.

In service of addressing these limits core-geth intends to supersede the static-document paradigm with a dynamically generated service description provided by the API itself via both build and runtime introspection. 


