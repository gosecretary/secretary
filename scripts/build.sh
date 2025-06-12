#!/bin/bash -e

set -e

BIN_DIR="./bin"

echo "Building Secretary..."

# Ensure bin directory exists
mkdir -p "$BIN_DIR"

# Tidy up dependencies
go mod tidy

# Build the server
GOOS=linux GOARCH=amd64 go build -o "$BIN_DIR/secretary" ./cmd/server/main.go

# Optionally, build a Darwin binary for Mac users
go build -o "$BIN_DIR/secretary-darwin" ./cmd/server/main.go

echo "Build complete! Executables: $BIN_DIR/secretary, $BIN_DIR/secretary-darwin"
