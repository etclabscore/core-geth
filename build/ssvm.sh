#!/usr/bin/env bash

set -e

if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/second-state/ssvm-evmc#notice\" documentation for instructions to build in other environments."
    exit 1
fi

mkdir -p build/_workspace/SSVM/build/tools/ssvm-evmc/
[[ -f build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so ]] && exit 0
wget -O build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so \
    https://github.com/second-state/ssvm-evmc/releases/download/evmc7-0.1.1/libssvm-evmc.so
