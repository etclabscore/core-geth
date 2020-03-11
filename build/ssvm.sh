#!/usr/bin/env bash

set -e

mkdir -p build/_workspace
[ ! -d build/_workspace/SSVM ] && git clone https://github.com/second-state/SSVM build/_workspace/SSVM || echo "SSVM exists."
cd build/_workspace/SSVM
git fetch --tags
git checkout 0.4.0
mkdir -p build
cd build
cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTS=ON ..
make
echo "Built library at: $(pwd)/tools/ssvm-evmc/libssvmEVMC.so"