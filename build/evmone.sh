#!/bin/sh

set -e

mkdir -p build/_workspace/evmone/build/tools/ssvm-evmc/
wget -O build/_workspace/evmone/evmone-0.2.0-linux-x86_64.tar.gz https://github.com/ethereum/evmone/releases/download/v0.2.0/evmone-0.2.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/evmone/evmone-0.2.0-linux-x86_64.tar.gz -C build/_workspace/evmone/