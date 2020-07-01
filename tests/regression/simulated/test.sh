#!/usr/bin/env bash

set -e

main(){
    datadir="$(mktemp -d)"
    trap "rm -rf $datadir" EXIT
    ./build/bin/geth --datadir="$datadir" init "$1"
	./build/bin/geth --datadir="$datadir" import "$2"
}

main $*