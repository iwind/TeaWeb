#!/usr/bin/env bash

export GOPATH=`pwd`/../../
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go build -o teaweb-linux main.go
echo "OK"