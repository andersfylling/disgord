#!/bin/bash

# https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
set -euox pipefail

DIR=$(pwd)
echo "$DIR"

for d in docs/examples/*/; do
  cd "./$d"
  go fmt .
  go build .
  cd "$DIR"
done
