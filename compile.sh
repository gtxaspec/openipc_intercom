#!/bin/bash

set -x

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o intercom_mipsle main.go
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o intercom_amd64 main.go

##upx

if [[ "$1" == "upx" ]]; then

UPX_BINARY=upx
$UPX_BINARY --ultra-brute intercom_mipsle
$UPX_BINARY --ultra-brute intercom_amd64

fi
