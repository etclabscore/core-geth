#!/usr/bin/env bash
set -e

go build -o ./build/bin/devgeth ./cmd/geth

# Clean up any bad leftovers.
pkill -f ./build/bin/devgeth || true
rm -fr /tmp/gethddd

onexit(){
	pkill -f ./build/bin/devgeth
	rm -fr /tmp/gethddd
}
trap onexit EXIT

./build/bin/geth --port 30313 --datadir=/tmp/gethddd --nodiscover --maxpeers=0 --rpc --rpcapi=admin,debug,eth,ethash,miner,net,personal,rpc,txpool,web3 >/tmp/geth.log 2>&1 &
disown
gethpid=$!
echo "Geth PID: ${gethpid}"

# Wait for geth to start up.
echo "Waiting 3 for geth to startup..."
sleep 3

# Save a copy of the generated openrpc document by HTTP RPC query.
http --json POST http://localhost:8545 id:=$(date +%s) method='rpc_describe' params:='[]' | jj -p 'result' | tee "$(pwd)/.develop/spec.json"

# Update our gist before validation; the script will exit if the validator fails.
gist -u 4da4c08765679dac1899543002d1f545 "$(pwd)/.develop/spec.json" >/dev/null 2>&1 &

# Run the OpenRPC document validator.
# Not sure the install here actually works. FIXME.
# command -v openrpc-generator-cli >/dev/null 2>&1 || { npm install -g @etclabscore/openrpc-generator-cli; }
openrpc-validator-cli "$(pwd)/.develop/spec.json"

kill $gethpid

