#!/usr/bin/env bash

mkdir -p build/_workspace
[ ! -d build/_workspace/hera ] && git clone https://github.com/ewasm/hera build/_workspace/hera || echo "Hera exists."
cd build/_workspace/hera
git submodule update --init
mkdir build
cd build
cmake -DBUILD_SHARED_LIBS=ON ..
cmake --build .
echo "Built library at: $(pwd)/src/libhera.so"