#!/usr/bin/env bash
#
# suite_compressor is intended for use with the tests/ suite.
# It strips the filename from a benchmark test result line.
# This generalization causes the benchstat util to treat each file as a rerun of a single
# test, which yields a more generalized statistic including almost believable p values.
#
# Eg.
#
# BenchmarkVM/vmArithmeticTest/addmod1_overflow2.json-12             35870             36298 ns/op                24.27 mgas/s       21434 B/op     311 allocs/op
# BenchmarkVM/vmArithmeticTest/addmod1_overflow3.json-12             29626             37401 ns/op               102.9 mgas/s        22371 B/op     315 allocs/op
# BenchmarkVM/vmArithmeticTest/addmod1_overflow4.json-12             30820             37519 ns/op               103.2 mgas/s        22382 B/op     315 allocs/op
#
# becomes:
#
# BenchmarkVM/vmArithmeticTest             35870             36298 ns/op                24.27 mgas/s       21434 B/op     311 allocs/op
# BenchmarkVM/vmArithmeticTest             29626             37401 ns/op               102.9 mgas/s        22371 B/op     315 allocs/op
# BenchmarkVM/vmArithmeticTest             30820             37519 ns/op               103.2 mgas/s        22382 B/op     315 allocs/op
#
while read -r line
do
    if grep -q '/op' <<< "$line"
    then echo "$(dirname $(echo ${line} | cut -d' ' -f1)) $(echo $line | cut -d' ' -f2-)"
    fi
done
