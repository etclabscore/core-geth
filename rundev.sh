#!/usr/bin/env bash

make geth
./build/bin/geth --dev --dev.period 5 --ipcpath /tmp/geth_dev.ipc --keystore ./ks --miner.gaslimit=8000000 --verbosity 5 console 2> geth.log
