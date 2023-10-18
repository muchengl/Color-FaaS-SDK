#!/usr/bin/env bash
export PATH="$PATH:$(go env GOPATH)/bin"

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o helloworld_arm64 helloworld.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o helloworld_amd64 helloworld.go
go build -o helloworld_raw helloworld.go