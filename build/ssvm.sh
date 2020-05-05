#!/usr/bin/env bash

set -e

mkdir -p build/_workspace/SSVM/build/tools/ssvm-evmc/
wget --no-check-certificate -O build/_workspace/SSVM/build/tools/ssvm-evmc/libssvmEVMC.so https://github.com/second-state/SSVM/releases/download/0.5.0/libssvmEVMC-linux-x86_64.so

# The follow are steps that might be used to build the SSVM SO from source.
# This library is only used for testing.
# Previous versions of this source worked by this pattern, but
# fail for me with 0.5.0.
# So I opted to use the above download pattern, which is simpler.
#mkdir -p build/_workspace
#[ ! -d build/_workspace/SSVM ] && git clone https://github.com/second-state/SSVM build/_workspace/SSVM || echo "SSVM exists."
#cd build/_workspace/SSVM
#git fetch --tags
## git checkout 0.5.0
#git checkout 0.5.0
#mkdir -p build
#cd build
#cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTS=ON ..
#make
#echo "Built library at: $(pwd)/tools/ssvm-evmc/libssvmEVMC.so"