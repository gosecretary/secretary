#!/bin/bash

set -e

# Default to development mode
DEV_MODE="--dev"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --no-dev)
      DEV_MODE=""
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--no-dev]"
      exit 1
      ;;
  esac
done

# Load environment variables from .env if it exists
if [ -f .env ]; then
  echo "Loading environment variables from .env..."
  export $(grep -v '^#' .env | xargs)
fi

# Build the project
./scripts/build.sh

# Run the server
echo "Starting Secretary server..."
if [ -n "$DEV_MODE" ]; then
  echo "Running in development mode with admin user..."
fi
./bin/secretary-darwin server $DEV_MODE