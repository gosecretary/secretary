#!/bin/bash

set -e

# Load environment variables from .env if it exists
if [ -f .env ]; then
  echo "Loading environment variables from .env..."
  export $(grep -v '^#' .env | xargs)
fi

# Build the project
./scripts/build.sh

# Run the server
./bin/secretary 