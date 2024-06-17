@echo off

echo "Please wait while binaries are compiled... (can take 5-30 minutes)"

REM This script builds for windows/amd64, arm64, arm and linux/

setlocal

REM indexer
cd ../indexer
set GOOS=windows&&set GOARCH=amd64&&go build -o ../release/binaries/windows-amd64-indexer.exe
set GOOS=windows&&set GOARCH=arm&&go build -o ../release/binaries/windows-arm-indexer.exe
set GOOS=windows&&set GOARCH=arm64&&go build -o ../release/binaries/windows-arm64-indexer.exe

set GOOS=linux&&set GOARCH=386&&go build -o ../release/binaries/linux-386-indexer
set GOOS=linux&&set GOARCH=amd64&&go build -o ../release/binaries/linux-amd64-indexer
set GOOS=linux&&set GOARCH=arm&&go build -o ../release/binaries/linux-arm-indexer
set GOOS=linux&&set GOARCH=arm64&&go build -o ../release/binaries/linux-arm64-indexer

set GOOS=darwin&&set GOARCH=amd64&&go build -o ../release/binaries/darwin-amd64-indexer
set GOOS=darwin&&set GOARCH=arm64&&go build -o ../release/binaries/darwin-arm64-indexer

set GOOS=js&&set GOARCH=wasm&&go build -o ../release/binaries/js-wasm-indexer

REM server
cd ../server
set GOOS=windows&&set GOARCH=amd64&&go build -o ../release/binaries/windows-amd64-server.exe
set GOOS=windows&&set GOARCH=arm&&go build -o ../release/binaries/windows-arm-server.exe
set GOOS=windows&&set GOARCH=arm64&&go build -o ../release/binaries/windows-arm64-server.exe

set GOOS=linux&&set GOARCH=386&&go build -o ../release/binaries/linux-386-server
set GOOS=linux&&set GOARCH=amd64&&go build -o ../release/binaries/linux-amd64-server
set GOOS=linux&&set GOARCH=arm&&go build -o ../release/binaries/linux-arm-server
set GOOS=linux&&set GOARCH=arm64&&go build -o ../release/binaries/linux-arm64-server

set GOOS=darwin&&set GOARCH=amd64&&go build -o ../release/binaries/darwin-amd64-server
set GOOS=darwin&&set GOARCH=arm64&&go build -o ../release/binaries/darwin-arm64-server

set GOOS=js&&set GOARCH=wasm&&go build -o ../release/binaries/js-wasm-server

REM go back to release
cd ../release

endlocal