#!/bin/bash
set -euo pipefail

# Run all Go tests with coverage and verbose output
echo "Running all Go tests..."
go test -v -cover ./...

# Print summary
echo "All tests completed successfully." 