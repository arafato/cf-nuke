#!/bin/bash
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/cf-nuke-linux-amd64 main.go

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/cf-nuke-windows-amd64.exe main.go

echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o bin/cf-nuke-macos-arm64 main.go

echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o bin/cf-nuke-macos-amd64 main.go

echo "Build complete."
