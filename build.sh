#!/bin/sh
go build -ldflags "-X main.debug=no -s -w" -tags netgo -o "bin"
upx --brute bin
