#!/bin/bash

# Hardcoded list of file paths
files=(
    "./handlers/"
    "./services/"
    "./utils/"
    "./middleware/"
)

# Loop through the hardcoded list of files
for file in "${files[@]}"; do
    echo "Testing $file"
    go test $file
done
