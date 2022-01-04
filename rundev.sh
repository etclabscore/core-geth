#!/usr/bin/env bash

make getj
./build/bin/geth --dev --dev.period 1 --miner.etherbase f6d7e4e39f35a0b8f61dd0a24a2dc92a3a5e0b01 --ipcpath /tmp/geth_dev.ipc --keystore $(pwd)/.ia/ks --miner.gaslimit=8000000 --verbosity 5 --http --http.corsdomain='*' console 2> geth.log
#./build/bin/geth --dev --dev.period 5 --ipcpath /tmp/geth_dev.ipc --keystore ./ks --miner.gaslimit=8000000 --verbosity 5 console 2> geth.log
