#!/bin/bash

# This script builds for windows/amd64, arm64, arm, and linux/

echo "Please wait while binaries are compiled... (can take 5-30 minutes)"

# indexer
cd ../indexer

GOOS=windows GOARCH=amd64 go build -o ../release/binaries/windows-amd64-indexer.exe
GOOS=windows GOARCH=arm go build -o ../release/binaries/windows-arm-indexer.exe
GOOS=windows GOARCH=arm64 go build -o ../release/binaries/windows-arm64-indexer.exe

GOOS=linux GOARCH=386 go build -o ../release/binaries/linux-386-indexer
GOOS=linux GOARCH=amd64 go build -o ../release/binaries/linux-amd64-indexer
GOOS=linux GOARCH=arm go build -o ../release/binaries/linux-arm-indexer
GOOS=linux GOARCH=arm64 go build -o ../release/binaries/linux-arm64-indexer

GOOS=darwin GOARCH=amd64 go build -o ../release/binaries/darwin-amd64-indexer
GOOS=darwin GOARCH=arm64 go build -o ../release/binaries/darwin-arm64-indexer

GOOS=js GOARCH=wasm go build -o ../release/binaries/js-wasm-indexer

# server
cd ../server

GOOS=windows GOARCH=amd64 go build -o ../release/binaries/windows-amd64-server.exe
GOOS=windows GOARCH=arm go build -o ../release/binaries/windows-arm-server.exe
GOOS=windows GOARCH=arm64 go build -o ../release/binaries/windows-arm64-server.exe

GOOS=linux GOARCH=386 go build -o ../release/binaries/linux-386-server
GOOS=linux GOARCH=amd64 go build -o ../release/binaries/linux-amd64-server
GOOS=linux GOARCH=arm go build -o ../release/binaries/linux-arm-server
GOOS=linux GOARCH=arm64 go build -o ../release/binaries/linux-arm64-server

GOOS=darwin GOARCH=amd64 go build -o ../release/binaries/darwin-amd64-server
GOOS=darwin GOARCH=arm64 go build -o ../release/binaries/darwin-arm64-server

GOOS=js GOARCH=wasm go build -o ../release/binaries/js-wasm-server

cd ../release
