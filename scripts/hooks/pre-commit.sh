#!/bin/sh

STAGED_GO_FILES=$(git diff --cached --name-only --filter=d -- '*.go')

if [ -z "$STAGED_GO_FILES" ]; then
  echo "no go files updated, skipping format."
else
  echo "Formatting staged Go files..."
  for file in $STAGED_GO_FILES; do
    if [ -f "$file" ]; then
      go fmt "$file"
      git add "$file"
    fi
  done
fi

echo "Running golangci-lint"
golangci-lint run ./...
