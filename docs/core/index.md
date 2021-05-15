---
title: About
---

# Core-Geth

CoreGeth is sponsored by and maintained with the leadership of [ETC Labs](https://etclabs.org) with an obvious core intention of stewarding
the Ethereum Classic opinion that the reversion of transactions in inconvenient situations shouldn't be permissible. 

But the spirit of the project intends to reach beyond Ethereum and Ethereum Classic, and indeed to reimagine an EVM node software that 
approaches the EVM-based protocols as technology that can -- and should -- be generalizable.

While CoreGeth inherits from and exposes complete feature parity with Ethereum Foundation's <sup>:registered:</sup> [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum),
there are quite few things that make CoreGeth unique.

## Additional Features

CoreGeth maintainers are [regular](https://github.com/ethereum/go-ethereum/pulls?q=author%3Ameowsbits) [contributors](https://github.com/ethereum/go-ethereum/pulls?q=author%3Aziogaschr+) [upstream](https://github.com/ethereum/go-ethereum/pulls?q=author%3Aiquidus+), but not all CoreGeth features are practicable or accepted there. The following categories document features specific to CoreGeth that ethereum/go-ethereum can't, or won't, implement. 

### Extended RPC API

#### Comprehensive RPC API Service Discovery

CoreGeth features a synthetic build/+runtime service discovery API, allowing you to get a [well structured](https://open-rpc.org/)
description of _all_ available methods, their parameters, and results.

!!! tip "RPC Service Documentation"
    For complete documentation of the available JSON RPC APIs, please see the [JSON RPC API page](/core-geth/JSON-RPC-API/modules/eth).

#### Additional methods and options

- Available `trace_block` and `trace_transaction` RPC API congruent to the OpenEthereum API (including a 1000x performance improvement vs. go-ethereum's `trace_transaction` in some cases).
    + _TODO:_ Talk more about this! And examples!
- Added `debug_removePendingTransaction` API method ([#203](https://github.com/etclabscore/core-geth/pull/203/files))
- Comprehensive service discovery with OpenRPC through method `rpc.discover`.

### EVMCv7 Support

- EVMCv7 support allows use with external EVMs (including EWASM).
- See [Running Geth with an External VM](./core/evmc) for more information.

### Remote Store for Ancient Chaindata

- Remote freezer, store your `ancient` data on Amazon S3 or Storj.
    - _TODO_: Talk more about this, provide examples.

### Extended CLI

- `--eth.protocols` configures `eth/x` protocol prioritization, eg. `65,64,63`.

### Developer Features: Tools

- A developer mode `--dev.pow` able to mock Proof-of-Work block schemas and production at representative Poisson intervals.
  + `--dev.poisson` configures Poisson intervals for block emission
- Chain configuration acceptance of OpenEthereum and go-ethereum chain configuration files (and the extensibility to support _any_ chain configuration schema).
- At the code level, a 1:1 EIP/ECIP specification to implementation pattern; disentangling Ethereum Foundation :registered: hard fork opinions from code. This yields more readable code, more precise naming and conceptual representation, more testable code, and a massive step toward Ethereum as a generalizeable technology.
- `copydb` will default to a sane fallback value if no parameter is passed for the second `<ancient/path>` argument.
- The `faucet` command supports an `--attach` option allowing the program to reference an already-running node instance
  (assuming it has an available RPC API) instead of restricting the faucet to a dedicated light client. Likewise, a `--syncmode=[full|fast|light]` option is provided for networks where _LES_ support may be lacking.

### Risk Management

- [Public](https://ci.etccore.in/blue/organizations/jenkins/core-geth-regression/activity) chaindata [regression testing](https://github.com/etclabscore/core-geth/blob/master/Jenkinsfile) run at each push to `master`.

### Extended _Improvement Proposal_ Support (EIP, ECIP, *IP)

- Myriad additional ECIP support:
  + ECBP1100 (aka MESS, an "artificial finality" gadget)
  + ECIP1099 (DAG growth limit)
  + ECIP1014 (defuse difficulty bomb), etc. :wink:

- Out-of-the-box support for Ethereum Classic.
  Chain configs are selected as `./build/bin/geth --<chain>`. For a list of supported networks and their CLI options, use `./build/bin/geth --help`.

## Divergent Design

How CoreGeth is built differently than ethereum/go-ethereum.

### Developer Features: Code

One of CoreGeth's most significant divergences from ethereum/go-ethereum at the code level is a reimagining (read: _massive overhaul_) of the `ChainConfig` data type and its methods.

At _ethereum/go-ethereum_ the `ChainConfig` makes protocol-facing feature activation decisions as follows:

```go
blockNumber := big.NewInt(0)
config := params.MainnetChainConfig
if config.IsByzantium(blockNumber) {
	// do a special thing for post-Byzantium chains
}
```

This, for the uninitiated developer, raises some questions:
- What's Byzantium?
- Which of the nine distinct Byzantium upgrades is this implementing?
- Does feature `Byzantium.X` depend on also having `Byzantium.Y` activated?

The developers of ethereum/go-ethereum have made this architectural decision because ethereum/go-ethereum is _only designed
and intended_ to support one chain: _Ethereum_. From this perspective, configurability presents a risk rather than a desirable feature.

While a hardcoded feature-group pattern (ie _hardfork upgrades_) in some ways mitigates a risk of "movable parts," and undesirable or unintended feature interactions,
it also presents a massive hurdle for extensibility.

!!! seealso "A metaphor"
    Consider the metaphor of the wiring the electrical circuits of a house. 
    
    With ethereum/go-ethereum's design, the television, the kitchen lights, the garbage disposal, and the garage door are all controlled by the same switch. If you want to watch TV,
    you also have to have the kitchen lights on, the garbage disposal running, and the garage door open.
    
    For an electrician whose _only concern_ is meeting the arbitrary specifications of an eccentric customer who _demands_ that their
    house work in this very strange way (forever), hard-installing these devices on the same circuit makes sense. The electrician commits to only
    serving one customer with this house, and the customer commits to their wiring preference.
    
    But, for anyone _else_ looking at this house, the design is absurd. Another home-owner may want to use a TV
     and a garage door in their own designs, but maybe don't want them to be codependent. Building the feature of a garbage disposal as being inextricable from a TV -- 
    from the perspective of a technologist (or consumer products designer, or anyone interested in these technologies as generalizeable things, rather than details of an eccentric house) -- 
    this arbitrary feature-bundling is patently absurd. 
    
    This is an Ethereum-as-technology perspective versus an Ethereum-as-network perspective, and reimagining a home where you can have the kitchen lights on
    without also turning the TV on is one of the things CoreGeth does.  

This same code as above, in CoreGeth, would look as follows:

```go
blockNumber := big.NewInt(0)
config := params.MainnetChainConfig
if config.IsEnabled(config.EIP658Transition, blockNumber) {
	// do a special thing for post-EIP658 chains
}
```

!!! example "Interface Reference"
    The complete interface pattern for supported feature methods
    can be found here: https://github.com/etclabscore/core-geth/blob/master/params/types/ctypes/configurator_iface.go

The implicit feature-group `Byzantium` is deconstructed into its composite features, using EIPs and ECIP specifications as conceptual delineations as well as naming patterns.

This makes the implementation of Improvement Proposal specifications referencable and readily testable. You can look up the implementation of EIP658 and see directly how it modifies transaction encoding, without having to disambiguate its implementation from state trie cleaning, gas repricing, opcode additions, block reward changes, or difficulty adjustments.
  You can test block reward modifications without _also_ having to test difficulty adjustments (... and state root differences, and ...).

!!! hint "Configuration Data Types"
    Not only does CoreGeth's interface pattern provide descriptive, articulate code; it also
    allows for the support of _arbitrary configuration data types_; CoreGeth supports configuration via
    ethereum/go-ethereum's `genesis` data structure (eg. `geth dumpgenesis genesis.json`) as well as OpenEthereum's JSON configuration schema.
    Extending support for any other configuration schema is likewise possible.

As should be obvious by now, this also allows selective feature adoption for configurations that don't want to bundle changes exactly like the Ethereum Foundation has. 
  For example, without this decomposition, Ethereum Classic would have had to accept and (re)implement the Difficulty Bomb _and_ reduce block rewards in order to adopt a change to the RLP encoding of transaction receipts change :exploding_head:


## Limitations

Things ethereum/go-ethereum can or will do that CoreGeth won't, or doesn't by default.

- A huge and diverse number of default pipeline-delivered build targets.
  This is a defaults and configuration sanity challenge for CoreGeth. We're vastly outnumbered by ethereum/go-ethereum maintainers
  and contributors, and ensuring proper delivery of a whole bunch of diverse artifacts is beyond our capacity.
  With that said, just because CoreGeth doesn't provide artifacts for a given architecture or OS doesn't mean it can't. 
  If ethereum/go-ethereum can build and package for it, then with some elbow grease, CoreGeth can too.
- The `puppeth` CLI program has been [removed](https://github.com/etclabscore/core-geth/pull/270). This is a "wizard"-style interactive program that helps beginners
  configure chain and network settings.
- Trim absolute file paths during build. As of a [somewhat-recent](TODO) Go version, `go build` provides a `-trim` flag
  which reduces the size of the binaries and anonymizes the build environment. This was removed because stripping file paths
  caused automatic service discovery features to break (they depend, in part, on source file path availability for build-time AST and runtime reflection). 
