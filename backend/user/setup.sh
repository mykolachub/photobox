#!/bin/bash

# Check if the .env file exists
if [ ! -f ".env" ]; then
    echo "Error: .env file not found."
    exit 1
fi

max_keys=$(grep -v '^$\|^\s*#' .env | wc -l)

counter=0

# Loop over each key-value pair in the .env file and export them as environment variables
while IFS= read -r line || [ -n "$line" ]; do
    # Increment the counter for each key processed
    ((counter++))

    # Check if the counter exceeds the maximum keys limit
    if [ "$counter" -gt "$max_keys" ]; then
        echo "Max keys limit reached. Stopping execution."
        break
    fi

    # Skip empty lines and comments starting with #
    if [[ -n "$line" && "$line" != "#"* ]]; then
        # Split each line into key and value based on the first occurrence of '='
        key="${line%%=*}"
        value="${line#*=}"

        # Export the environment variable
        export "$key"="$value"
        echo "Run: export $key=$value"
    fi
done < .env

echo "Go environment variables setup completed."