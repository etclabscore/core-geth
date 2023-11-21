---
title: Mordor Testnet Guide
---

!!! tip "Mordor Testnet"
    Mordor is a PoW Ethereum Classic testnet. A testnet allows developers to perform specific tests. Developers may want to test protocol changes, test a smart contract, or interact with the network in anyway that does not require real EthClassic (ETC)—just don’t test on mainnet, mainnet is for production.

## Summary:

+ Install Core-geth https://etclabscore.github.io/core-geth/getting-started/installation/
+ Create an account on --mordor
+ Run --mordor with --mine enabled
+ Create a Script to mine Mordor METC

## Install Core-geth

https://etclabscore.github.io/core-geth/getting-started/installation/

## Mordor Testnet Mining Guide

You can visit the Core-geth documentation for more installation options. I’m using Ubuntu 22.04 LTS.

If you just want to download and run `geth --mordor` or any of the other tools here, this is the quickest and simplest way.

Binary archives are published at https://github.com/etclabscore/core-geth/releases. Find the latest one for your OS, download it, (check the SHA sum), unarchive it, and run!

When running Core-geth use `--mordor` flag for Ethereum Classic testnet.

```shell
$ wget https://github.com/etclabscore/core-geth/releases/download/v1.12.16/core-geth-linux-v1.12.16.zip # Update to the most current release version
$ sudo unzip core-geth-linux-v1.12.16.zip -d /bin/ # Update to the most current release version
$ geth --help # Lists available options
$ geth --mordor # Runs Ethereum Classic's testnet Mordor
```

## Account Creation

You'll need an account with an address (0x...) to receive your mETC mining rewards. Here is how you make an address and keystore file with core-geth. You'll be able to import the keystore file into wallets like MetaMask. Backup this file. You'll mine your mordor rewards to it.

```shell
$ geth --mordor account new # Creates a new account with a public address and keystore file
$ geth --mordor account list
$ geth --mordor # Runs Ethereum Classic's testnet Mordor
```

You’ll notice listing the account will print the keystore file location.For example:keystore:///home/USER/.ethereum/mordor/keystore/UTC...

## Run Mordor with Mining Enable

```shell
$ geth --mordor --mine --minerthreads 1 --miner.etherbase 0x_INSERT_YOUR_ADDRESS_HERE_3a087
```

Check Mordor Balance on Blockscout

So, you’re running a Mordor node and mining testnet mETC. Woohoo! An easy way to double check you’re actually growing a Mordor testnet balance is on [Blockscout](https://etc-mordor.blockscout.com). Just search the account address you created earlier.

## Add your Mordor Account to a Wallet?

You can use your keystore file to import your wallet into a wallet application such as MetaMask. In MetaMask

* Add the Ethereum Classic mainnet to your MetaMask by visiting https://chainlist.org/chain/61 and clicking the "Add to MetaMask" button.
* Add the Mordor testnet to your MetaMask by visiting https://chainlist.org/chain/63 and clicking the "Add to MetaMask" button.
* Under your account profile select import account > select type (JSON) > upload your keystore file. You may need to enter the account password.

## Mordor Mining Script

One way to avoid typing or copy and pasting this same text block is creating a shell script file.

```shell
$ geth --mordor --mine --minerthreads 1 --miner.etherbase 0x_INSERT_YOUR_ADDRESS_HERE_3a087
```

Enter the following in a new terminal window (ctrl + alt + t):

```shell
touch start-mordor.sh && echo "geth --mordor --mine --minerthreads 1 --miner.etherbase 0x_INSERT_YOUR_ADDRESS_HERE_3a087" >start-mordor.sh && chmod +x start-mordor.sh
```

touch start-mordor.sh to create the file && echo “the contents” into the shell script file && add chmod executable+x permissions to the file.

Enter the following in a new terminal window (ctrl + alt + t):

```shell
 ./start-mordor.sh
```

Great job! You are mining on Ethereum Classic's testnet.

## Donate unused mETC to a community faucet.

https://faucet.mordortest.net

A faucet is a developer tool that gives users testnet tokens to use when testing smart contracts or interacting with DApps on test networks. https://faucet.mordortest.net gives Mordor testnet ETC to test smart contracts before pushing them to production on the Ethereum Classic mainnet. Faucets like this allow network users and developers to interact with the Mordor network without the prerequisite of mining mETC via a client node. To donate mETC to this public faucet, please send your mETC to address `0x51Cb0EA27f03e56d84E9EB1879F131393a6769bA`.
