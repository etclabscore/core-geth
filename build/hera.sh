#!/usr/bin/env bash

set -e

# rm -rf build/_workspace/hera/0.3.2
# mkdir -p build/_workspace/hera/0.3.2
# wget -P build/_workspace/hera/ https://github.com/ewasm/hera/releases/download/v0.3.2/hera-0.3.2-linux-x86_64.tar.gz
# tar -C build/_workspace/hera/0.3.2 -xzf build/_workspace/hera/hera-0.3.2-linux-x86_64.tar.gz

if [[ "$OSTYPE" != "linux"* ]]; then
    echo "This script is only currently configured to work on Linux. Please see \"https://github.com/ewasm/hera#building-hera\" documentation for instructions to build in other environments."
    exit 1
fi

main() {
    mkdir -p build/_workspace
    [ ! -d build/_workspace/hera ] && git clone https://github.com/ewasm/hera build/_workspace/hera || echo "Hera exists."
    cd build/_workspace/hera
    git checkout v0.3.2
    git submodule update --init
    mkdir -p build
    cd build
    cmake -DBUILD_SHARED_LIBS=ON ..
    cmake --build .
    echo "Built library at: $(pwd)/src/libhera.so"
}
main
