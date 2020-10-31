#!/bin/sh

# https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
set -euox pipefail

DIR=$(pwd)
echo "$DIR"

for d in docs/examples/*/; do
  cd "./$d"
  go build .
  cd "$DIR"
done