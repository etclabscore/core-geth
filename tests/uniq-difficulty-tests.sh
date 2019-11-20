#!/usr/bin/env bash

f=tests/testdata/BasicTests/generated_difficulty/all_difficulty_tests.json
t=$(mktemp)
cat $f | sort | uniq > $t
mv $t $f
