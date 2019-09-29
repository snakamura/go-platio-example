#!/usr/bin/env bash

set -e

rm -rf mock_platio
mkdir -p mock_platio
~/go/bin/mockgen -source=platio/api.go -destination=mock_platio/api.go

mkdir -p coverage
go test -cover ./... -coverprofile=coverage/coverage.out
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
