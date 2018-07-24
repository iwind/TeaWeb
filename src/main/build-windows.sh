#!/usr/bin/env bash

export GOPATH=`pwd`/../../
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=386

go build -o ${GOPATH}/dist/bin/teaweb-32.exe main.go

export GOARCH=amd64
go build -o ${GOPATH}/dist/bin/teaweb-64.exe main.go
