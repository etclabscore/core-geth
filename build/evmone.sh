#!/usr/bin/env bash

set -e

TARGET_VERSION=0.9.1
# TARGET_VERSION=0.9.0
# TARGET_VERSION=0.6.0

_OSTYPE=${OSTYPE}
if [[ ${_OSTYPE} == "linux"* ]]; then
    _OSTYPE="linux"
fi
if [[ ${_OSTYPE} == "darwin"* ]]; then
    _OSTYPE="darwin"
fi


download()
{
    echo "Downloading artifact for '${_OSTYPE}'"
    wget -O build/_workspace/evmone/evmone-${TARGET_VERSION}-${_OSTYPE}-x86_64.tar.gz https://github.com/ethereum/evmone/releases/download/v${TARGET_VERSION}/evmone-${TARGET_VERSION}-${_OSTYPE}-x86_64.tar.gz
}

unpack()
{
    rm -rf build/_workspace/evmone/bin
    rm -rf build/_workspace/evmone/include
    rm -rf build/_workspace/evmone/lib

    tar xzvf build/_workspace/evmone/evmone-${TARGET_VERSION}-${_OSTYPE}-x86_64.tar.gz -C build/_workspace/evmone/
}

mkdir -p build/_workspace/evmone
[[ -f build/_workspace/evmone/evmone-${TARGET_VERSION}-${_OSTYPE}-x86_64.tar.gz ]] || download
unpack

