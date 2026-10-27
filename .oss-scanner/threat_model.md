# CoreGeth threat model

CoreGeth is a full-node client for Ethereum Classic (ETC) and its Mordor testnet, based on go-ethereum v1.13.14 with individual upstream backports since. It validates blocks under proof of work (etchash), relays blocks and transactions over devp2p, serves peers and JSON-RPC, and can mine.

CoreGeth runs on most of the ETC network. A bug that every node on a version shares is therefore a network-level problem, not a single-host one: a consensus divergence splits the network, and a crash that a relayed block or transaction can trigger takes most of the network down at once.

## What we protect, in order

1. **Consensus.** Every CoreGeth node accepts exactly the blocks the ECIPs make valid and rejects the rest, and reaches the same state root, whatever the version, platform, sync mode (full or snap), cache state, or history of restarts and reorgs. A block that one node accepts and another rejects, or that two nodes execute to different states, is the worst outcome.
2. **Network availability.** An unauthenticated peer cannot crash, stall or exhaust nodes, keep them off the canonical chain, or stop them relaying and serving peers. An input that is relayed before it crashes the receiver (a block, a header, a transaction) is worse than one that only affects the peer it is sent to, because it spreads through the network by itself.
3. **Mining.** Untrusted input cannot make the work packages a node hands to miners invalid or stale, or stall their creation.
4. **Keys and funds.** No remote party can read, use or delete keys in the keystore, or make the node sign something the operator did not ask for.
5. **Durable state.** The chain database survives restarts, crashes and reorgs without corruption that a remote party can cause.

## Who is trusted

| Source | Trust | Notes |
|---|---|---|
| devp2p peers | Untrusted | Discovery (`p2p/discover`), RLPx (`p2p/rlpx`), and the `eth` and `snap` protocols (`eth/protocols/`). Any internet host can connect to a node with default settings. Say how many connections an attack needs. |
| Blocks and headers with valid proof of work | Untrusted, costly | Producing them costs hashrate. State the difficulty the attack needs and what it would cost on ETC mainnet. A side chain forking from an old block can be cheap to build; state the fork point. Reorgs by an attacker with majority hashrate are a property of proof of work, not a CoreGeth vulnerability, unless CoreGeth handles them differently from what the ECIPs specify. |
| Transactions that pay fees | Untrusted, costly | State the fees the attack spends. |
| JSON-RPC over HTTP/WebSocket | Untrusted only where exposed | Off by default and bound to localhost when enabled; the default HTTP modules are `net` and `web3`. Issues reachable through `eth`, `net` and `web3`, which operators commonly expose, are in scope. `admin`, `debug`, `personal`, `miner` and `txpool` callers are trusted. |
| IPC, the datadir, config files, CLI flags | Trusted | Anything that needs local access or write access to these is out of scope. |

## How we rate severity

- **Critical**
  - Consensus divergence: accepting a block the ECIPs make invalid, rejecting a valid one, or computing a different state root than the rest of the network, on ETC mainnet or Mordor.
  - Any remote crash, hang or out-of-memory that can be triggered by a block, header or transaction a node relays, because it propagates to the whole network.
  - Creating, destroying or moving ETC outside the consensus rules (for example a wrong block reward or a wrong balance change).
  - Remote code execution.
- **High**
  - A remote crash, hang or out-of-memory caused by a single peer in the default configuration, when the input is not relayed.
  - A peer, with no special position and little traffic, keeping a node off the canonical chain or stalling sync (eclipse, peer-scoring bypass, poisoning the downloader or fetcher).
  - Remote theft or use of keystore keys.
  - Database corruption that a remote party can cause.
- **Medium**
  - Resource exhaustion that needs sustained traffic or many connections.
  - DoS through the `eth`/`net`/`web3` RPC namespaces with small requests.
  - Anything that needs a non-default flag; name the flag.
- **Low**
  - Anything that needs a trusted RPC namespace, local access, or operator misconfiguration.

Give a severity and the reason; do not inflate it.

## What counts as a resource-exhaustion finding

Cost proportional to the input, and per-peer state freed on disconnect, are expected. Use input that the network would accept: a block that is invalid under consensus is dropped before the expensive path, and ETC's gas limit keeps a legitimate block under about 2 MiB. Measure memory and CPU on the build in this image, mark each figure as measured or extrapolated, and say how many peers it needs.

## Components that matter most

These have no upstream project reviewing them:

- `consensus/ethash`: etchash, difficulty calculation, block rewards.
- `params/` and `params/types/`: the chain configurations for ETC and Mordor. A wrong activation block, or a rule enabled on the wrong chain, is a consensus bug.
- Mining (`miner/`) and proof-of-work block propagation (`eth/fetcher`, `eth/handler*.go`), which upstream removed after moving to proof of stake.

Code shared with upstream go-ethereum (`core/`, `core/vm`, `eth/protocols/`, `p2p/`, `trie/`, `rlp/`) is in scope as well. **When a finding is in code identical to upstream go-ethereum, say so**: our policy is to report it upstream too and ship upstream's fix, and other clients may need coordinated disclosure.

## Out of scope

- Ethereum mainnet, Sepolia and Holesky (`--mainnet`, `--sepolia`, `--holesky`). CoreGeth only implements their upgrades through Cancun and cannot follow them since Pectra.
- Proof of stake and the engine API (`consensus/beacon`, `beacon/`, `eth/catalyst`): Ethereum Classic does not use them.
- `BadHashes` in `core/blockchain.go` is kept deliberately. Do not report it or propose removing it.
- `--dev` mode and private chains configured by the operator.
- `cmd/` tools other than `cmd/geth` (`evm`, `devp2p`, `abigen`, `clef`, `ethkey`, `rlpdump`): operators run them by hand on their own input.
- Test code, `tests/`, `docs/`, `build/`, and Go module dependencies unless CoreGeth reaches the bug with attacker input; then report it against CoreGeth.

## How to exercise it

- The checkout is at `/src`, the binaries are in `/src/build/bin`, and test binaries are pre-compiled. Prefer targeted package tests, for example `go test -run TestX ./eth/protocols/eth/`.
- Consensus fixtures are in `tests/testdata` (ethereum/tests) and `tests/testdata-etc` (ETC-specific). `tests/regression/simulated/` imports recorded chains.
- Go fuzz targets exist (`grep -rl 'func Fuzz' --include='*_test.go' .`), including `rlp`, `trie`, `core/vm`, `eth/protocols/snap` and `tests/fuzzers/`.
- Use private chains only (`--dev`, or a genesis file with `--nodiscover`). Never connect to ETC mainnet or Mordor.

## Reports and patches we want

- The scanned revision, the entry point, the flags needed (default or not), and the attacker position: remote peer, block producer, transaction sender, or RPC client.
- A reproducer as a Go test in the affected package, run against the unmodified build, with a passing control run. Code reading alone is not enough.
- One root cause per report. Check the history and tests for an earlier fix, or for a test that asserts the behaviour as intended.
- A minimal patch with a test that fails before it and passes after.
