---
hide:
  - navigation # Hide navigation
title: CoreGeth
---

# CoreGeth: An Ethereum Protocol Provider

[![OpenRPC](https://img.shields.io/static/v1.svg?label=OpenRPC&message=1.14.0&color=blue)](#openrpc-discovery)
[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://godoc.org/github.com/etclabscore/core-geth)
[![Go Report Card](https://goreportcard.com/badge/github.com/etclabscore/core-geth)](https://goreportcard.com/report/github.com/etclabscore/core-geth)
[![Travis](https://travis-ci.org/etclabscore/core-geth.svg?branch=master)](https://travis-ci.org/etclabscore/core-geth)
[![Gitter](https://badges.gitter.im/core-geth/community.svg)](https://gitter.im/core-geth/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

> An [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) downstream effort to make the Ethereum Protocol accessible and extensible for a diverse ecosystem.

Priority is given to reducing opinions around chain configuration, IP-based feature implementations, and API predictability.
Upstream development from [ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) is merged to this repository regularly,
 usually at every upstream tagged release. Every effort is made to maintain seamless compatibility with upstream source, including compatible RPC, JS, and CLI
 APIs, data storage locations and schemas, and, of course, interoperable node protocols. Applicable bug reports, bug fixes, features, and proposals should be
 made upstream whenever possible.

## Network/provider comparison

Networks supported by the respective core-geth packaged `geth` program.

| Ticker | Consensus         | Network                               | core-geth                                                | ethereum/go-ethereum |
| ---    | ---               | ---                                   | ---                                                      | ---                  |
| ETC    | :zap:             | Ethereum Classic                      | :heavy_check_mark:                                       |                      |
| ETH    | :zap:             | Ethereum (Foundation)                 | :heavy_check_mark:                                       | :heavy_check_mark:   |
| -      | :zap: :handshake: | Private chains                        | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :zap:             | Mordor (Geth+Parity ETH PoW Testnet)  | :heavy_check_mark:                                       |                      |
|        | :zap:             | Morden (Geth+Parity ETH PoW Testnet)  |                                                          |                      |
|        | :zap:             | Ropsten (Geth+Parity ETH PoW Testnet) | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Rinkeby (Geth-only ETH PoA Testnet)   | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Goerli (Geth+Parity ETH PoA Testnet)  | :heavy_check_mark:                                       | :heavy_check_mark:   |
|        | :handshake:       | Kovan (Parity-only ETH PoA Testnet)   |                                                          |                      |
|        |                   | Tobalaba (EWF Testnet)                |                                                          |                      |
|        |                   | Ephemeral development PoA network     | :heavy_check_mark:                                       | :heavy_check_mark:   |

- :zap: = __Proof of Work__
- :handshake: = __Proof of Authority__

---

> <a name="ellaism-footnote">1</a>: This is originally an [Ellaism
Project](https://github.com/ellaism). However, A [recent hard
fork](https://github.com/ellaism/specs/blob/master/specs/2018-0003-wasm-hardfork.md)
makes Ellaism not feasible to support with go-ethereum any more. Existing
Ellaism users are asked to switch to
[Parity](https://github.com/paritytech/parity).

> <a name="configuration-capable">2</a>: Network not supported by default, but network configuration is possible. Make a PR!
