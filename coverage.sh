#!/bin/sh

# generate coverage stats
go test -coverprofile=coverage.txt

# list coverage profile for each func
go tool cover -coverprofile=coverage.txt -covermode=atomic

# open a html file in your web browser for a pretty output
#go tool cover -html=coverage.out
