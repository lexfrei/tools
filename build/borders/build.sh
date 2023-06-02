#!/bin/bash

# Build for Windows (amd64)
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o borders-windows-amd64.exe main.go
upx --lzma --best borders-windows-amd64.exe

# Build for Mac (amd64)
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o borders-mac-amd64 main.go
upx --lzma --best borders-mac-amd64

# Build for Mac (arm64)
env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o borders-mac-arm64 main.go
upx --lzma --best borders-mac-arm64