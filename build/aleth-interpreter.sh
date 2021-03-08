#!/bin/sh

set -e

mkdir -p build/_workspace/aleth
wget -O build/_workspace/aleth/aleth-1.8.0-linux-x86_64.tar.gz https://github.com/ethereum/aleth/releases/download/v1.8.0/aleth-1.8.0-linux-x86_64.tar.gz
tar xzvf build/_workspace/aleth/aleth-1.8.0-linux-x86_64.tar.gz -C build/_workspace/aleth/
