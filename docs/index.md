---
hide:
  - navigation # Hide navigation
title: CoreGeth
---

# CoreGeth

> A Go client for Ethereum Classic, based on [go-ethereum](https://github.com/ethereum/go-ethereum).

[![Go Reference](https://pkg.go.dev/badge/github.com/etclabscore/core-geth.svg)](https://pkg.go.dev/github.com/etclabscore/core-geth)
[![Go Report Card](https://goreportcard.com/badge/github.com/etclabscore/core-geth)](https://goreportcard.com/report/github.com/etclabscore/core-geth)
[![Tests](https://github.com/etclabscore/core-geth/actions/workflows/test-linux.yml/badge.svg?branch=master)](https://github.com/etclabscore/core-geth/actions/workflows/test-linux.yml)
[![Latest release](https://img.shields.io/github/v/release/etclabscore/core-geth)](https://github.com/etclabscore/core-geth/releases/latest)
[![Release date](https://img.shields.io/github/release-date/etclabscore/core-geth)](https://github.com/etclabscore/core-geth/releases/latest)
[![Last commit](https://img.shields.io/github/last-commit/etclabscore/core-geth)](https://github.com/etclabscore/core-geth/commits/master)

CoreGeth runs Ethereum Classic mainnet, the Mordor testnet and private proof-of-work chains. It keeps the proof-of-work consensus and mining that Ethereum Classic runs on, which upstream go-ethereum removed after Ethereum moved to proof of stake.

It stays as close to upstream go-ethereum as Ethereum Classic allows. RPC, JS and CLI APIs, data storage locations and schemas, and node protocols are kept compatible with upstream, so tooling written for geth works with CoreGeth. Bug reports, fixes and proposals that apply to upstream code belong upstream. The [README](https://github.com/etclabscore/core-geth#relationship-with-upstream-go-ethereum) describes how changes are merged.

## Maintenance

CoreGeth is maintained in [etclabscore/core-geth](https://github.com/etclabscore/core-geth) and has been released from it since March 2020. It has implemented every Ethereum Classic upgrade from Phoenix (2020) to Spiral (2024) as specified through the [ECIP process](https://ecips.ethereumclassic.org). It is the client most of the network runs today, as [etcnodes.org](https://etcnodes.org) shows.

If development ever moves to another repository, the project's README will say so.

## Getting a release

CoreGeth releases are published only on the repository's [releases page](https://github.com/etclabscore/core-geth/releases), with container images at [etclabscore/core-geth](https://hub.docker.com/r/etclabscore/core-geth) on Docker Hub. To check that a binary is a CoreGeth release:

- It comes from `github.com/etclabscore/core-geth/releases` or the `etclabscore/core-geth` Docker Hub image, and nowhere else.
- Its version is in the `v1.12.x` line. Version numbers are assigned by this project and say nothing about where a binary came from, since anyone can tag a fork with a higher number.
- Every archive comes with a `.sha256` file built by the repository's release workflow. Container images have no such file; they are identified by their tag on the `etclabscore` Docker Hub account, and each tag matches a release on the releases page.
- The project's only website is this documentation, at [etclabscore.github.io/core-geth](https://etclabscore.github.io/core-geth). Any other domain using the CoreGeth name is unrelated to this project.

To install and run it, see [Installation](getting-started/installation.md) and [Running the CLI](getting-started/run-cli.md).

## Networks

Networks selectable with a flag in the `geth` program.

| Network                     | Flag         | Consensus     | Status    |
| ---                         | ---          | ---           | ---       |
| Ethereum Classic (ETC)      | `--classic`  | PoW (etchash) | supported |
| Mordor (ETC testnet)        | `--mordor`   | PoW (etchash) | supported |
| MintMe.com Coin (MINTME)    | `--mintme`   | PoW           | supported |
| Private chains              | genesis file | PoW / PoA     | supported |
| Ephemeral development chain | `--dev`      | PoA           | supported |

## Security

Please do not report vulnerabilities in public issues. Use GitHub's private [vulnerability report](https://github.com/etclabscore/core-geth/security/advisories/new); see [SECURITY.md](https://github.com/etclabscore/core-geth/blob/master/SECURITY.md) for details.

## Forks and other clients

Please be cautious about using other repositories that carry the CoreGeth name. They are not maintained or reviewed by this project, and a release published there is not a CoreGeth release. Forks are welcome. We ask that they use a different name so operators can tell which client they are running.

CoreGeth is one of several Ethereum Classic clients. Others include:

- [getc](https://github.com/go-ethereum-classic/go-ethereum-classic), a Go client based on a newer go-ethereum.
- [Nethermind](https://github.com/NethermindEth/nethermind) with the [ETC plugin](https://github.com/nethermindEtc/nethermind-etc-plugin), a .NET client.
- [Besu](https://github.com/besu-eth/besu) with the [ETC plugin](https://github.com/besu-etc/besu-etc-plugin), a Java client.

Which one to run is the operator's choice. Running more than one client makes the network safer.
