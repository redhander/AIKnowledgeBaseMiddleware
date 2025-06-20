#!/bin/bash

# Bash script for development with hot reload
# Usage: ./scripts/dev.sh

echo "Starting AI Knowledge Base Middleware in development mode..."

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "Air not found. Installing air..."
    go install github.com/air-verse/air@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install air. Please install it manually."
        exit 1
    fi
    echo "Air installed successfully!"
fi

# Check if tmp directory exists, create if not
if [ ! -d "tmp" ]; then
    mkdir -p tmp
    echo "Created tmp directory"
fi

# Start air for hot reload
echo "Starting hot reload server..."
air 