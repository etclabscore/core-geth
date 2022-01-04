#!/usr/bin/env bash

env AM=on AM_CLIENT=/tmp/geth_dev.ipc AM_KEYSTORE=$(pwd)/.ia/ks AM_ADDRESS=0xf6d7e4e39f35a0b8f61dd0a24a2dc92a3a5e0b01 AM_PASSWORDFILE=$(pwd)/pass.txt go test -timeout=8h -count 1 -run TestState -v ./tests
