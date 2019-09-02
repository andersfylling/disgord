#!/bin/bash

go test -bench=BenchmarkCRS -count 10 >results-native.txt

benchstat results-interface.txt results-native.txt >perf-diff.txt