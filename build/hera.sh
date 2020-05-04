#!/usr/bin/env bash

set -e

mkdir -p build/_workspace/hera/build/src
wget --no-check-certificate -O build/_workspace/hera/hera-0.2.5-linux-x86_64.tar.gz https://github.com/ewasm/hera/releases/download/v0.2.5/hera-0.2.5-linux-x86_64.tar.gz
mkdir build/_workspace/hera/tmp
tar -C build/_workspace/hera/tmp -zxvf build/_workspace/hera/hera-0.2.5-linux-x86_64.tar.gz
mv build/_workspace/hera/tmp/lib/libhera.so build/_workspace/hera/build/src/libhera.so
rm -rf build/_workspace/hera/tmp

#mkdir -p build/_workspace
#[ ! -d build/_workspace/hera ] && git clone https://github.com/ewasm/hera build/_workspace/hera || echo "Hera exists."
#cd build/_workspace/hera
#git submodule update --init
#mkdir -p build
#cd build
#cmake -DBUILD_SHARED_LIBS=ON ..
#cmake --build .
#echo "Built library at: $(pwd)/src/libhera.so"