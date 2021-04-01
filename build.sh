#!/bin/sh
go build -ldflags "-X main.debug=false -s -w" -tags netgo -o "bin"
upx --brute bin
