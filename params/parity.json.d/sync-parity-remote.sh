#!/usr/bin/env bash

# This script downloads JSON chain configurations from github.com/parity-tech/parity
# master branch.
# It uses a whitelist approach to avoid including irrelevant and/or unsupported
# configurations.

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

for spec_name in "${specs[@]}"; do
    echo "Fetching $spec_name..."
	curl -q -O https://raw.githubusercontent.com/paritytech/parity-ethereum/master/ethcore/res/ethereum/"$spec_name".json
done
