#!/bin/sh
CGO_ENABLED=0 go build -ldflags "-X main.debug=no -s -w" -o "bin"
