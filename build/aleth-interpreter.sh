#!/usr/bin/env bash

set -e


if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/ethereum/aleth\" documentation for instructions to build in other environments."
    exit 1
fi

mkdir -p build/_workspace/aleth
[[ -f build/_workspace/aleth/aleth-1.8.0-linux-x86_64.tar.gz ]] && exit 0
wget -O build/_workspace/aleth/aleth-1.8.0-linux-x86_64.tar.gz https://github.com/ethereum/aleth/releases/download/v1.8.0/aleth-1.8.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/aleth/aleth-1.8.0-linux-x86_64.tar.gz -C build/_workspace/aleth/
