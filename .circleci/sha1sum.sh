#!/bin/bash
find . -type f -name "*.go" -print0 | sort -z | xargs -0 sha1sum | sha1sum