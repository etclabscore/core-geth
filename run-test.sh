#!/usr/bin/env bash

for f in $(find ./tests/testdata/GeneralStateTests/stEIP161Specific -type f); do echo "TESTFILE: $f"; ./build/bin/evm --verbosity 2 --dump-always statetest "$f" |& tee -a test-state.out; done

