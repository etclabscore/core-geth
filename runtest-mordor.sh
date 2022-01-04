#!/usr/bin/env bash

env AM=on AM_CLIENT=/tmp/multigeth_mordor.ipc AM_KEYSTORE=$(pwd)/.ia/ks AM_ADDRESS=0x877BD459c9B7D8576B44E59E09d076C25946F443 AM_PASSWORDFILE=$(pwd)/.ia/877bdpass.txt go test -timeout 99999s -count 1 -run TestState -v ./tests |& tee mordor.out
