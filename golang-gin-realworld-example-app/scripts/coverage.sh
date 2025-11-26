#!/bin/bash

# Run tests and generate coverage profile
go test ./... -coverprofile=coverage.out

# Check if coverage.out file was created
if [ -f coverage.out ]; then
    # Generate HTML report
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
else
    echo "No coverage profile found."
fi

# Clean up
rm -f coverage.out