#!/bin/sh

# --- Link Checker for Local HTML Files ---
# This script recursively finds all .html files in the specified directory (or current directory)
# and runs 'linkchecker' on each file, displaying only the failing links ([NG] results).

# Name of the linkchecker executable (ensure this is in your PATH)
LINK_CHECKER_BIN="linkchecker"
OPTION="-no-internal"

# Directory to scan (default to current directory)
TARGET_DIR="${1:-.}"

echo "--- Starting local link check in: ${TARGET_DIR} ---"
echo "Searching recursively for .html files..."
echo ""

# Check if linkchecker executable is available
if ! command -v "$LINK_CHECKER_BIN" &> /dev/null
then
    echo "Error: '$LINK_CHECKER_BIN' command not found."
    echo "Please ensure linkchecker is installed and available in your PATH (e.g., \$GOPATH/bin)."
    exit 1
fi

# Use find to recursively locate .html files and process them
find "$TARGET_DIR" -type f -name "*.html" | while read -r FILE_PATH
do
    echo "File: $FILE_PATH"
    $LINK_CHECKER_BIN -u "$FILE_PATH" "$OPTION" | grep '\[NG\]'
done

echo "--- Local link check finished ---"

exit 0