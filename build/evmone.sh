#!/usr/bin/env bash

set -e

if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/ethereum/evmone\" documentation for instructions to build in other environments."
    exit 1
fi

mkdir -p build/_workspace/evmone
[[ -f build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz ]] && exit 0
wget -O build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz https://github.com/ethereum/evmone/releases/download/v0.5.0/evmone-0.5.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz -C build/_workspace/evmone/
