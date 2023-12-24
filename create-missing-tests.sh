#!/bin/bash

# Find all .go files excluding _test.go files
go_files=$(find . -type f -name "*.go" ! -name "*_test.go" ! -name "*_mock.go")

# Loop through each .go file and generate _test.go files
for file in $go_files; do
    # Extract the directory name
    dir=$(dirname "$file")

    # Create the corresponding _test.go file
    test_file="${dir}/$(basename "$file" .go)_test.go"

    # Skip if the _test.go file already exists
    if [ -e "$test_file" ]; then
        echo "Skipping existing test file: $test_file"
        continue
    fi

    # Create the _test.go file with the package declaration
    echo "package $(basename "$dir")" > "$test_file"
    echo "Created $test_file"
done

echo "Done."