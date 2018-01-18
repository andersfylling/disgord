#!/bin/sh

# generate coverage stats
go test -coverprofile=coverage.out

# list coverage profile for each func
go tool cover -func=coverage.out

# open a html file in your web browser for a pretty output
go tool cover -html=coverage.out
