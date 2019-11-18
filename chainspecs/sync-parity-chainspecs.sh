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

mkdir -p ./originals
mv ./*.json ./originals/

cp ./originals/*.json ./

echo "> Fixing hex encoding rules to play nicely with our libs."
set -x
echo  "Trim all addresses and hashes of 0x- prefixes"
sed -E -i 's/0x([a-fA-F0-9]{16})/\1/g' ./*.json
#sed -E -i 's/0x([a-fA-F0-9]{40,})/\1/g' ./*.json

echo  "Remove all leading-0 hex quantities"
sed -E -i 's/0x0([a-fA-F1-9])/0x\1/g' ./*.json

echo "Replace all 0x00 -> 0x0"
sed -E -i 's/0x00/0x0/g' ./*.json

echo "Replace all stringy single-digit decimals with hex encoding."
sed -E -i 's/\"([0-9]{1})\"/\"0x\1\"/g' ./*.json

echo  "Replace all addresses and hashes with 0x- prefixes."
sed -E -i 's/(\b[a-fA-F0-9]{16})/0x\1/g' ./*.json

echo "Replace default max code size decimal -> hex encoding."
sed -E -i 's/24576/"0x6000"/g' ./*.json