---
title: Build from Source
---

## Hardware Requirements

Minimum:

* CPU with 2+ cores
* 4GB RAM
* 500GB free storage space to sync the Mainnet
* 8 MBit/sec download Internet service

Recommended:

* Fast CPU with 4+ cores
* 16GB+ RAM
* High Performance SSD with at least 500GB free space
* 25+ MBit/sec download Internet service

## Dependencies

- Make sure your system has __Go 1.16+__ installed. https://golang.org/doc/install
- Make sure your system has a C compiler installed. For example, with Linux Ubuntu:

```shell
$ sudo apt-get install -y build-essential
```

## Source

Once the dependencies have been installed, it's time to clone and build the source:

```shell
$ git clone https://github.com/etclabscore/core-geth.git
$ cd core-geth
$ make all
$ ./build/bin/geth --help
```

## Build docker image

You can build a local docker image directly from the source:

```shell
$ git clone https://github.com/etclabscore/core-geth.git
$ cd core-geth
$ docker build -t=core-geth .
```

Or with all tools:

```shell
$ docker build -t core-geth-alltools -f Dockerfile.alltools .
```
