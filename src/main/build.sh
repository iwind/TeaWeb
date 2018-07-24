#!/usr/bin/env bash

export GOPATH=`pwd`/../../
export CGO_ENABLED=0

go build -o ${GOPATH}/dist/bin/teaweb main.go
echo "OK"