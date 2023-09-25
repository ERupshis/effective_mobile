#!/bin/bash

# Change to the project root directory (if necessary)
# cd /path/to/your/project

# Generate documentation for all packages and subpackages
for package in $(go list ./...); do
    # Get the package name without slashes
    package_name=$(echo "$package" | tr / -)

    # Remove the prefix 'github.com-erupshis-effective_mobile-' if it exists
    package_name="${package_name#github.com-erupshis-effective_mobile-}"

    # Create a directory for the package's documentation
    mkdir -p "documentation/$package_name"

    # Generate the documentation content using go doc with the -all flag and redirect to the Markdown file
    go doc -all "$package" > "documentation/$package_name/documentation.md"
done

# Check if any Markdown documentation files were generated
if ls documentation/*/*.md 1>/dev/null 2>&1; then
    echo "Documentation was generated and saved in MD files."
else
    echo "No documentation generated."
fi
