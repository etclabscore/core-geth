#!/usr/bin/env bash

set -e

main() {
    mkdir -p build/_workspace
    [ ! -d build/_workspace/hera ] && git clone https://github.com/ewasm/hera build/_workspace/hera || echo "Hera exists."
    cd build/_workspace/hera
    git checkout v0.2.5 # Last EVMCv6-compatible tag (v0.3.0 bumps to v7).
    git submodule update --init
    mkdir -p build
    cd build
    cmake -DBUILD_SHARED_LIBS=ON ..
    cmake --build .
    echo "Built library at: $(pwd)/src/libhera.so"
}
main