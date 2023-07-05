---
title: Command Line Interface (CLI)
---

## Running `geth`

!!! tip "Use for Ethereum mainnet"

    While `core-geth` is mainly used for the Ethereum Classic network, you can use it for Ethereum mainnet and other [supported networks](/#networkprovider-comparison) as well.

### Fast node on an Ethereum Classic network

By far the most common scenario is people wanting to simply interact with the Ethereum
network: create accounts; transfer funds; deploy and interact with contracts. For this
particular use-case the user doesn't care about years-old historical data, so we can
fast-sync quickly to the current state of the network. To do so:

```
$ geth [|--classic|--testnet|--rinkeby|--goerli|--mordor] console
```

This command will:

 * Start `geth` in fast sync mode (default, can be changed with the `--syncmode` flag),
   causing it to download more data in exchange for avoiding processing the entire history
   of the Ethereum network, which is very CPU intensive.
 * Start up `geth`'s built-in interactive [JavaScript console](https://geth.ethereum.org/docs/interface/javascript-console),
   (via the trailing `console` subcommand) through which you can invoke all official [`web3` methods](https://web3js.readthedocs.io/en/v1.2.9/)
   as well as `geth`'s own [management APIs](https://geth.ethereum.org/docs/rpc/server).
   This tool is optional and if you leave it out you can always attach to an already running
   `geth` instance with `geth attach`.

### A Full node on the Mordor test network

Transitioning towards developers, if you'd like to play around with creating Ethereum
contracts, you almost certainly would like to do that without any real money involved until
you get the hang of the entire system. In other words, instead of attaching to the main
network, you want to join the **mordor test** network with your node, which is fully equivalent to
the main network, but with play-Ether only.

```shell
$ geth --mordor console
```

The `console` subcommand has the exact same meaning as above and they are equally
useful on the testnet too. Please see above for their explanations if you've skipped here.

Specifying the `--mordor` flag, however, will reconfigure your `geth` instance a bit:

 * Instead of using the default data directory (`~/.ethereum` on Linux for example), `geth`
   will nest itself one level deeper into a `mordor` subfolder (`~/.ethereum/mordor` on
   Linux). Note, on OSX and Linux this also means that attaching to a running testnet node
   requires the use of a custom endpoint since `geth attach` will try to attach to a
   production node endpoint by default. E.g.
   `geth attach <datadir>/mordor/geth.ipc`. Windows users are not affected by
   this.
 * Instead of connecting the main Ethereum network, the client will connect to the mordor's test
   network, which uses different P2P bootnodes, different network IDs and genesis states.

!!! Note

    Although there are some internal protective measures to prevent transactions from
    crossing over between the classic network and test network, you should make sure to always
    use separate accounts for play-money and real-money. Unless you manually move
    accounts, `geth` will by default correctly separate the two networks and will not make any
    accounts available between them.*

### Configuration

As an alternative to passing the numerous flags to the `geth` binary, you can also pass a
configuration file via:

```shell
$ geth --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to
export your existing configuration:

```shell
$ geth --your-favourite-flags dumpconfig
```

!!! Note

    This works only with `geth` v1.6.0 and above.*

## Command-line Options

```
$ geth --help
NAME:
   geth - the ETC Core Go-Ethereum command line interface

   Copyright 2013-2019 The go-ethereum Authors

USAGE:
   geth [options] command [command options] [arguments...]

VERSION:
   1.11.21-unstable

COMMANDS:
   account                            Manage accounts
   attach                             Start an interactive JavaScript environment (connect to node)
   console                            Start an interactive JavaScript environment
   copydb                             Create a local chain from a target chaindata folder
   dump                               Dump a specific block from storage
   dumpconfig                         Show configuration values
   dumpgenesis                        Dumps genesis block JSON configuration to stdout
   export                             Export blockchain into file
   export-preimages                   Export the preimage database into an RLP stream
   import                             Import a blockchain file
   import-preimages                   Import the preimage database from an RLP stream
   init                               Bootstrap and initialize a new genesis block
   inspect                            Inspect the storage size for each type of data in the database
   js                                 Execute the specified JavaScript files
   license                            Display license information
   makecache                          Generate ethash verification cache (for testing)
   makedag                            Generate ethash mining DAG (for testing)
   removedb                           Remove blockchain and state databases
   show-deprecated-flags              Show flags that have been deprecated
   version                            Print version numbers
   version-check                      Checks (online) whether the current version suffers from any known security vulnerabilities
   wallet                             Manage Ethereum presale wallets
   help, h                            Shows a list of commands or help for one command

ETHEREUM OPTIONS:
  --config value                      TOML configuration file
  --datadir value                     Data directory for the databases and keystore (default: "/Users/ziogaschr/Library/Ethereum")
  --datadir.ancient value             Data directory for ancient chain segments (default = inside chaindata)
  --ancient.rpc value                 Connect to a remote freezer via RPC. Value must an HTTP(S), WS(S), unix socket, or 'stdio' URL. Incompatible with --datadir.ancient
  --keystore value                    Directory for the keystore (default = inside the datadir)
  --nousb                             Disables monitoring for and managing USB hardware wallets
  --pcscdpath value                   Path to the smartcard daemon (pcscd) socket file
  --networkid value                   Explicitly set network id (integer)(For testnets: use --ropsten, --rinkeby, --goerli, --mordor, --yolov2 instead) (default: 1)
  --classic                           Ethereum Classic network: pre-configured Ethereum Classic mainnet
  --mordor                            Mordor network: Ethereum Classic's cross-client proof-of-work test network
  --rinkeby                           Rinkeby network: pre-configured proof-of-authority test network
  --goerli                            Görli network: pre-configured proof-of-authority test network
  --rinkeby                           Rinkeby network: pre-configured proof-of-authority test network
  --yolov2                            YOLOv2 network: pre-configured proof-of-authority shortlived test network.
  --mintme                            MintMe.com Coin mainnet: pre-configured MintMe.com Coin mainnet
  --ropsten                           Ropsten network: pre-configured proof-of-work test network
  --syncmode value                    Blockchain sync mode ("fast", "full", or "light") (default: fast)
  --exitwhensynced                    Exits after block synchronisation completes
  --gcmode value                      Blockchain garbage collection mode ("full", "archive") (default: "full")
  --txlookuplimit value               Number of recent blocks to maintain transactions index by-hash for (default = index all blocks) (default: 0)
  --ethstats value                    Reporting URL of a ethstats service (nodename:secret@host:port)
  --identity value                    Custom node name
  --lightkdf                          Reduce key-derivation RAM & CPU usage at some expense of KDF strength
  --whitelist value                   Comma separated block number-to-hash mappings to enforce (<number>=<hash>)
  --ecbp1100 value                    Configure ECBP-1100 (MESS) block activation number (default: 18446744073709551615)

LIGHT CLIENT OPTIONS:
  --light.serve value                 Maximum percentage of time allowed for serving LES requests (multi-threaded processing allows values over 100) (default: 0)
  --light.ingress value               Incoming bandwidth limit for serving light clients (kilobytes/sec, 0 = unlimited) (default: 0)
  --light.egress value                Outgoing bandwidth limit for serving light clients (kilobytes/sec, 0 = unlimited) (default: 0)
  --light.maxpeers value              Maximum number of light clients to serve, or light servers to attach to (default: 100)
  --ulc.servers value                 List of trusted ultra-light servers
  --ulc.fraction value                Minimum % of trusted ultra-light servers required to announce a new head (default: 75)
  --ulc.onlyannounce                  Ultra light server sends announcements only
  --light.nopruning                   Disable ancient light chain data pruning

DEVELOPER CHAIN OPTIONS:
  --dev                               Ephemeral proof-of-authority network with a pre-funded developer account, mining enabled
  --dev.period value                  Block period for proof-of-authority network to use in developer mode (0 = mine only if transaction pending) (default: 0)
  --dev.pow                           Ephemeral proof-of-work network with a pre-funded developer account, mining enabled

ETHASH OPTIONS:
  --ethash.cachedir value             Directory to store the ethash verification caches (default = inside the datadir)
  --ethash.cachesinmem value          Number of recent ethash caches to keep in memory (16MB each) (default: 2)
  --ethash.cachesondisk value         Number of recent ethash caches to keep on disk (16MB each) (default: 3)
  --ethash.cacheslockmmap             Lock memory maps of recent ethash caches
  --ethash.dagdir value               Directory to store the ethash mining DAGs (default: "/Users/ziogaschr/Library/Ethash")
  --ethash.dagsinmem value            Number of recent ethash mining DAGs to keep in memory (1+GB each) (default: 1)
  --ethash.dagsondisk value           Number of recent ethash mining DAGs to keep on disk (1+GB each) (default: 2)
  --ethash.dagslockmmap               Lock memory maps for recent ethash mining DAGs

TRANSACTION POOL OPTIONS:
  --txpool.locals value               Comma separated accounts to treat as locals (no flush, priority inclusion)
  --txpool.nolocals                   Disables price exemptions for locally submitted transactions
  --txpool.journal value              Disk journal for local transaction to survive node restarts (default: "transactions.rlp")
  --txpool.rejournal value            Time interval to regenerate the local transaction journal (default: 1h0m0s)
  --txpool.pricelimit value           Minimum gas price limit to enforce for acceptance into the pool (default: 1)
  --txpool.pricebump value            Price bump percentage to replace an already existing transaction (default: 10)
  --txpool.accountslots value         Minimum number of executable transaction slots guaranteed per account (default: 16)
  --txpool.globalslots value          Maximum number of executable transaction slots for all accounts (default: 4096)
  --txpool.accountqueue value         Maximum number of non-executable transaction slots permitted per account (default: 64)
  --txpool.globalqueue value          Maximum number of non-executable transaction slots for all accounts (default: 1024)
  --txpool.lifetime value             Maximum amount of time non-executable transaction are queued (default: 3h0m0s)

PERFORMANCE TUNING OPTIONS:
  --cache value                       Megabytes of memory allocated to internal caching (default = 4096 mainnet full node, 128 light mode) (default: 1024)
  --cache.database value              Percentage of cache memory allowance to use for database io (default: 50)
  --cache.trie value                  Percentage of cache memory allowance to use for trie caching (default = 15% full mode, 30% archive mode) (default: 15)
  --cache.trie.journal value          Disk journal directory for trie cache to survive node restarts (default: "triecache")
  --cache.trie.rejournal value        Time interval to regenerate the trie cache journal (default: 1h0m0s)
  --cache.gc value                    Percentage of cache memory allowance to use for trie pruning (default = 25% full mode, 0% archive mode) (default: 25)
  --cache.snapshot value              Percentage of cache memory allowance to use for snapshot caching (default = 10% full mode, 20% archive mode) (default: 10)
  --cache.noprefetch                  Disable heuristic state prefetch during block import (less CPU and disk IO, more time waiting for data)
  --cache.preimages                   Enable recording the SHA3/keccak preimages of trie keys (default: true)

ACCOUNT OPTIONS:
  --unlock value                      Comma separated list of accounts to unlock
  --password value                    Password file to use for non-interactive password input
  --signer value                      External signer (url or path to ipc file)
  --allow-insecure-unlock             Allow insecure account unlocking when account-related RPCs are exposed by http

API AND CONSOLE OPTIONS:
  --ipcdisable                        Disable the IPC-RPC server
  --ipcpath value                     Filename for IPC socket/pipe within the datadir (explicit paths escape it)
  --http                              Enable the HTTP-RPC server
  --http.addr value                   HTTP-RPC server listening interface (default: "localhost")
  --http.port value                   HTTP-RPC server listening port (default: 8545)
  --http.api value                    API's offered over the HTTP-RPC interface
  --http.corsdomain value             Comma separated list of domains from which to accept cross origin requests (browser enforced)
  --http.vhosts value                 Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
  --ws                                Enable the WS-RPC server
  --ws.addr value                     WS-RPC server listening interface (default: "localhost")
  --ws.port value                     WS-RPC server listening port (default: 8546)
  --ws.api value                      API's offered over the WS-RPC interface
  --ws.origins value                  Origins from which to accept websockets requests
  --graphql                           Enable GraphQL on the HTTP-RPC server. Note that GraphQL can only be started if an HTTP server is started as well.
  --graphql.corsdomain value          Comma separated list of domains from which to accept cross origin requests (browser enforced)
  --graphql.vhosts value              Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
  --rpc.gascap value                  Sets a cap on gas that can be used in eth_call/estimateGas (0=infinite) (default: 25000000)
  --rpc.txfeecap value                Sets a cap on transaction fee (in ether) that can be sent via the RPC APIs (0 = no cap) (default: 1)
  --jspath loadScript                 JavaScript root path for loadScript (default: ".")
  --exec value                        Execute JavaScript statement
  --preload value                     Comma separated list of JavaScript files to preload into the console

NETWORKING OPTIONS:
  --bootnodes value                   Comma separated enode URLs for P2P discovery bootstrap
  --bootnodesv4 value                 Comma separated enode URLs for P2P v4 discovery bootstrap (light server, full nodes) (deprecated, use --bootnodes)
  --bootnodesv5 value                 Comma separated enode URLs for P2P v5 discovery bootstrap (light server, light nodes) (deprecated, use --bootnodes)
  --discovery.dns value               Sets DNS discovery entry points (use "" to disable DNS)
  --eth.protocols value               Sets the Ethereum Protocol versions (65|64|63) (default = 65,64,63 first is primary)
  --port value                        Network listening port (default: 30303)
  --maxpeers value                    Maximum number of network peers (network disabled if set to 0) (default: 50)
  --maxpendpeers value                Maximum number of pending connection attempts (defaults used if set to 0) (default: 0)
  --nat value                         NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>) (default: "any")
  --nodiscover                        Disables the peer discovery mechanism (manual peer addition)
  --v5disc                            Enables the experimental RLPx V5 (Topic Discovery) mechanism
  --netrestrict value                 Restricts network communication to the given IP networks (CIDR masks)
  --nodekey value                     P2P node key file
  --nodekeyhex value                  P2P node key as hex (for testing)

MINER OPTIONS:
  --mine                              Enable mining
  --miner.threads value               Number of CPU threads to use for mining (default: 0)
  --miner.notify value                Comma separated HTTP URL list to notify of new work packages
  --miner.gasprice value              Minimum gas price for mining a transaction (default: 1000000000)
  --miner.gastarget value             Target gas floor for mined blocks (default: 8000000)
  --miner.gaslimit value              Target gas ceiling for mined blocks (default: 8000000)
  --miner.etherbase value             Public address for block mining rewards (default = first account) (default: "0")
  --miner.extradata value             Block extra data set by the miner (default = client version)
  --miner.recommit value              Time interval to recreate the block being mined (default: 3s)
  --miner.noverify                    Disable remote sealing verification

GAS PRICE ORACLE OPTIONS:
  --gpo.blocks value                  Number of recent blocks to check for gas prices (default: 20)
  --gpo.percentile value              Suggested gas price is the given percentile of a set of recent transaction gas prices (default: 60)
  --gpo.maxprice value                Maximum gas price will be recommended by gpo (default: 500000000000)

VIRTUAL MACHINE OPTIONS:
  --vmdebug                           Record information useful for VM and contract debugging
  --vm.evm value                      External EVM configuration (default = built-in interpreter)
  --vm.ewasm value                    External ewasm configuration (default = built-in interpreter)

LOGGING AND DEBUGGING OPTIONS:
  --fakepow                           Disables proof-of-work verification
  --fakepow.poisson                   Disables proof-of-work verification and adds mining delay (Poisson) based on --miner.threads
  --nocompaction                      Disables db compaction after import
  --verbosity value                   Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail (default: 3)
  --vmodule value                     Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. eth/*=5,p2p=4)
  --backtrace value                   Request a stack trace at a specific logging statement (e.g. "block.go:271")
  --debug                             Prepends log messages with call-site location (file and line number)
  --pprof                             Enable the pprof HTTP server
  --pprof.addr value                  pprof HTTP server listening interface (default: "127.0.0.1")
  --pprof.port value                  pprof HTTP server listening port (default: 6060)
  --pprof.memprofilerate value        Turn on memory profiling with the given rate (default: 524288)
  --pprof.blockprofilerate value      Turn on block profiling with the given rate (default: 0)
  --pprof.cpuprofile value            Write CPU profile to the given file
  --trace value                       Write execution trace to the given file

METRICS AND STATS OPTIONS:
  --metrics                           Enable metrics collection and reporting
  --metrics.expensive                 Enable expensive metrics collection and reporting
  --metrics.addr value                Enable stand-alone metrics HTTP server listening interface (default: "127.0.0.1")
  --metrics.port value                Metrics HTTP server listening port (default: 6060)
  --metrics.influxdb                  Enable metrics export/push to an external InfluxDB database
  --metrics.influxdb.endpoint value   InfluxDB API endpoint to report metrics to (default: "http://localhost:8086")
  --metrics.influxdb.database value   InfluxDB database name to push reported metrics to (default: "geth")
  --metrics.influxdb.username value   Username to authorize access to the database (default: "test")
  --metrics.influxdb.password value   Password to authorize access to the database (default: "test")
  --metrics.influxdb.tags value       Comma-separated InfluxDB tags (key/values) attached to all measurements (default: "host=localhost")

WHISPER (deprecated) OPTIONS:
  --shh                               Enable Whisper
  --shh.maxmessagesize value          Max message size accepted (default: 1048576)
  --shh.pow value                     Minimum POW accepted (default: 0.2)
  --shh.restrict-light                Restrict connection between two whisper light clients

ALIASED (deprecated) OPTIONS:
  --rpc                               Enable the HTTP-RPC server (deprecated, use --http)
  --rpcaddr value                     HTTP-RPC server listening interface (deprecated, use --http.addr) (default: "localhost")
  --rpcport value                     HTTP-RPC server listening port (deprecated, use --http.port) (default: 8545)
  --rpccorsdomain value               Comma separated list of domains from which to accept cross origin requests (browser enforced) (deprecated, use --http.corsdomain)
  --rpcvhosts value                   Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (deprecated, use --http.vhosts) (default: "localhost")
  --rpcapi value                      API's offered over the HTTP-RPC interface (deprecated, use --http.api)
  --wsaddr value                      WS-RPC server listening interface (deprecated, use --ws.addr) (default: "localhost")
  --wsport value                      WS-RPC server listening port (deprecated, use --ws.port) (default: 8546)
  --wsorigins value                   Origins from which to accept websockets requests (deprecated, use --ws.origins)
  --wsapi value                       API's offered over the WS-RPC interface (deprecated, use --ws.api)
  --gpoblocks value                   Number of recent blocks to check for gas prices (deprecated, use --gpo.blocks) (default: 20)
  --gpopercentile value               Suggested gas price is the given percentile of a set of recent transaction gas prices (deprecated, use --gpo.percentile) (default: 60)
  --graphql.addr value                GraphQL server listening interface (deprecated, graphql can only be enabled on the HTTP-RPC server endpoint, use --graphql)
  --graphql.port value                GraphQL server listening port (deprecated, graphql can only be enabled on the HTTP-RPC server endpoint, use --graphql) (default: 8545)
  --pprofport value                   pprof HTTP server listening port (deprecated, use --pprof.port) (default: 6060)
  --pprofaddr value                   pprof HTTP server listening interface (deprecated, use --pprof.addr) (default: "127.0.0.1")
  --memprofilerate value              Turn on memory profiling with the given rate (deprecated, use --pprof.memprofilerate) (default: 524288)
  --blockprofilerate value            Turn on block profiling with the given rate (deprecated, use --pprof.blockprofilerate) (default: 0)
  --cpuprofile value                  Write CPU profile to the given file (deprecated, use --pprof.cpuprofile)

MISC OPTIONS:
  --snapshot                          Enables snapshot-database mode -- experimental work in progress feature
  --help, -h                          show help


COPYRIGHT:
   Copyright 2013-2020 The core-geth and go-ethereum Authors
```
