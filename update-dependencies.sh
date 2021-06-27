#!/bin/bash

# https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
set -euox pipefail

DIR=$(pwd)
echo "$DIR"

go get -u
go mod tidy

for d in docs/examples/*/; do
  cd "./$d"
  go get -u
  go mod tidy
  cd "$DIR"
done
