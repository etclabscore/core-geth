#!/usr/bin/env bash

set -e

TARGET_VERSION=0.9.1
# TARGET_VERSION=0.9.0
# TARGET_VERSION=0.6.0


if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/ethereum/evmone\" documentation for instructions to build in other environments."
    exit 1
fi

download()
{
    wget -O build/_workspace/evmone/evmone-${TARGET_VERSION}-linux-x86_64.tar.gz https://github.com/ethereum/evmone/releases/download/v${TARGET_VERSION}/evmone-${TARGET_VERSION}-linux-x86_64.tar.gz
}

unpack()
{
    tar xzvf build/_workspace/evmone/evmone-${TARGET_VERSION}-linux-x86_64.tar.gz -C build/_workspace/evmone/
}

mkdir -p build/_workspace/evmone
[[ -f build/_workspace/evmone/evmone-${TARGET_VERSION}-linux-x86_64.tar.gz ]] || download
unpack

