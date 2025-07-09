#!/bin/sh

# Check if the correct number of arguments is provided.
if [ $# -ne 2 ]; then
    echo "usage: $0 key filename"
    exit 1
fi

file_to_read="$1"
key_to_find="$2"

# Read the file line by line.
while IFS= read -r line; do
    # Skip empty lines.
    [ -z "$line" ] && continue

    # Check if the line starts with the key followed by '='.
    case "$line" in
        "$key_to_find="*)
            # Extract the value after the '=' character.
            value="${line#*=}"
            # Output the value and exit.
            printf "%s" "$value"
            exit 0
            ;;
    esac
done < "$file_to_read"

echo "error: key not found."
exit 2
