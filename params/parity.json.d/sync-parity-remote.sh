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
	"test-specs/frontier_test"
	"test-specs/homestead_test"
	"test-specs/eip150_test"
	"test-specs/eip161_test"
	"test-specs/eip210_test"
	"test-specs/byzantium_test"
	"test-specs/constantinople_test"
	"test-specs/st_peters_test"
	"test-specs/istanbul_test"
	"test-specs/transition_test"

	"foundation"
	"goerli"
	"rinkeby"
	"ropsten"

	"classic"
	"mordor"
	"kotti"

	"mix"
	"musicoin"
)

for spec_name in "${specs[@]}"; do
    echo "Fetching $spec_name..."
	curl -q https://raw.githubusercontent.com/paritytech/parity-ethereum/master/ethcore/res/ethereum/"$spec_name".json > ./params/parity.json.d/"$(basename ${spec_name})".json
done

curl -q https://api.github.com/repos/paritytech/parity-ethereum/git/refs/heads/master | jq .object.sha | sed 's/"//g' > ./params/parity.json.d/commit.txt
