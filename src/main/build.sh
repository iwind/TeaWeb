#!/usr/bin/env bash

export GOPATH=`pwd`/../../
export CGO_ENABLED=1

go build -o teaweb main.go
echo "OK"