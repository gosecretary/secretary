#!/bin/bash -e

set -e

BIN_DIR="./bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

echo "Building Secretary v${VERSION}..."

# Ensure bin directory exists
mkdir -p "$BIN_DIR"

# Tidy up dependencies
go mod tidy

# Build the server
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=${VERSION}'" -o "$BIN_DIR/secretary" ./cmd/secretary/main.go

# Optionally, build a Darwin binary for Mac users
go build -ldflags="-X 'main.Version=${VERSION}'" -o "$BIN_DIR/secretary-darwin" ./cmd/secretary/main.go

echo "Build complete! Executables:"
echo "  - $BIN_DIR/secretary (Linux)"
echo "  - $BIN_DIR/secretary-darwin (macOS)"
