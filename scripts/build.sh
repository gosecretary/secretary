#!/bin/bash -e

echo "Building Secretary..."

# Tidy up dependencies
go mod tidy

# Build the server
go build -o ./bin/secretary ./cmd/server/main.go

echo "Build complete! Executable: ./bin/secretary"
