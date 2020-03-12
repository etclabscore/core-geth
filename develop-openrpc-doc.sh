#!/usr/bin/env bash

set -e

[[ -f openrpc.json ]] || { wget --no-check-certificate https://raw.githubusercontent.com/etclabscore/ethereum-json-rpc-specification/master/openrpc.json; }

if ! pgrep geth > /dev/null 2>&1; then
(
make geth && ./build/bin/geth --datadir=/tmp/gethddd --nodiscover --maxpeers=0 --rpc --rpcapi=admin,debug,eth,ethash,miner,net,personal,rpc,txpool,web3
)&
fi

# Logs are truncated with each run, they only show the latest state.
http --json POST http://localhost:8545 id:=$(date +%s) method=rpc_setOpenRPCDiscoverDocument params:='["./openrpc.json"]' |& tee openrpc_set.log
grep -q error openrpc_set.log && exit 1 
http --json POST http://localhost:8545 id:=$(date +%s) method=rpc_describeOpenRPC params:='[]' | jj -p > openrpc_describe.log

# Developer can then inspect the logs, eg.

# cat openrpc_describe.log | jj -p | head -40

