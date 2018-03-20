#!/bin/bash

export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
export GOOS=linux
export GOARCH=arm
export GOARM=5

go get -u github.com/golang/dep/cmd/dep

# dep init ./
dep ensure --update ./
go build -o onwire-prom-exporter main.go