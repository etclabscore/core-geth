#!/usr/bin/env bash

set -e

specs=(
	"frontier_test"
	"homestead_test"
	"eip150_test"
	"eip161_test"
	"eip210_test"
	"byzantium_test"
	"constantinople_test"
	"st_peters_test"
	"istanbul_test"
	"transition_test"

	"foundation"
	"goerli"
	"rinkeby"
	"ropsten"

	"classic"
	"morden"
	"mordor"
	"kotti"

	"mix"
	"musicoin"
)

if [[ $# -ne 0 ]]; then
	specs=("$*")
fi

for spec_name in "${specs[@]}"; do
    echo "Fetching $spec_name..."
	curl -q -O https://raw.githubusercontent.com/paritytech/parity-ethereum/master/ethcore/res/ethereum/"$spec_name".json
done
