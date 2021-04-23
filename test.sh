#!/bin/sh
go clean -testcache
go test ./internal/board/
go test ./internal/game/
go test ./internal/rest/
go test ./internal/rest/auth/
