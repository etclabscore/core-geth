#!/bin/sh

set -e

mkdir -p build/_workspace/evmone
wget -O build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz https://github.com/ethereum/evmone/releases/download/v0.5.0/evmone-0.5.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/evmone/evmone-0.5.0-linux-x86_64.tar.gz -C build/_workspace/evmone/
