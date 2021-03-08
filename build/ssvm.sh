#!/bin/sh

set -e

mkdir -p build/_workspace/SSVM/build/tools/ssvm-evmc/
wget -O build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so \
    https://github.com/second-state/ssvm-evmc/releases/download/evmc7-0.1.1/libssvm-evmc.so
