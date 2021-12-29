#!/bin/bash

# https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
set -euox pipefail

DIR=$(pwd)
echo "$DIR"

for d in examples/*/; do
  cd "./$d"
  go fmt .
  go build -o example_program_binary .
  rm example_program_binary
  cd "$DIR"
done
