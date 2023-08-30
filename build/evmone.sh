#!/usr/bin/env bash

set -e

TARGET_VERSION=0.9.1
# TARGET_VERSION=0.9.0
# TARGET_VERSION=0.6.0

EXTENSION="so"

_OSTYPE=${OSTYPE}
if [[ ${_OSTYPE} == "linux"* ]]; then
    _OSTYPE="linux"
    EXTENSION="so"
fi
if [[ ${_OSTYPE} == "darwin"* ]]; then
    _OSTYPE="darwin"
    EXTENSION="dylib"
fi

main() {
    mkdir -p build/_workspace
    [ ! -d build/_workspace/evmone ] && git clone --recursive https://github.com/ethereum/evmone build/_workspace/evmone || echo "Evmone exists."
    cd build/_workspace/evmone
    git checkout v${TARGET_VERSION}
    cmake -S . -B build -DEVMONE_TESTING=ON
    cmake --build build --parallel
    echo "Built library at: $(pwd)/build/lib/libevmone.${EXTENSION}"
}
main
