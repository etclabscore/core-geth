#!/usr/bin/env bash

# This script downloads JSON chain configurations from github.com/parity-tech/parity
# master branch.
# It uses a whitelist approach to avoid including irrelevant and/or unsupported
# configurations.
# Must be run from project root.

set -e

if [[ ! $(pwd) =~ go-ethereum$ ]]; then
    echo "Must be run from go-ethereum project root"
    exit 1
fi

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
	curl -q https://raw.githubusercontent.com/paritytech/parity-ethereum/master/ethcore/res/ethereum/"$spec_name".json > ./params/parity.json.d/"$spec_name".json
done

curl -q https://api.github.com/repos/paritytech/parity-ethereum/git/refs/heads/master | jq .object.sha | sed 's/"//g' > ./params/parity.json.d/commit.txt
