#!/bin/bash

# Change to the project root directory
#cd /path/to/your/project

# Generate documentation for all packages and subpackages
for package in $(go list ./...); do
    go doc "$package" >> documentation.txt
done

# Check if documentation.txt is empty and delete it if it is
if [ -s "documentation.txt" ]; then
    echo "Documentation was generated and saved in documentation.txt"
else
    echo "Documentation is empty. Removing documentation.txt."
    rm "documentation.txt"
fi