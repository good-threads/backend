#!/bin/bash

# Check if correct number of arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <filename> <token>"
    exit 1
fi

filename=$1
token=$2

# Check if the file exists
if [ ! -f "$filename" ]; then
    echo "Error: File '$filename' not found."
    exit 1
fi

# Use grep to filter lines that do not contain the specified token
grep -v "$token" "$filename" > "$filename.tmp"

# Replace the original file with the temporary file
mv "$filename.tmp" "$filename"

echo "Lines containing '$token' have been deleted from $filename."