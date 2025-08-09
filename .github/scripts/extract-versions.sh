#!/bin/bash

# Script for extracting tool versions from Taskfile.yml
# with help of yq embedded in GitHub Actions runner

set -euo pipefail

TASKFILE="Taskfile.yml"

if [ ! -f "$TASKFILE" ]; then
  echo "Taskfile.yml not found" >&2
  exit 1
fi

echo "Extracting variables from Taskfile.yml..."

# Extract all .vars entries as step outputs
yq -r '.vars | to_entries[] | "\(.key)=\(.value)"' "$TASKFILE" \
  | while IFS= read -r kv; do 
      echo "$kv" >> "$GITHUB_OUTPUT"
      echo "  $kv"
    done

# Handle MODULES: prefer .vars.MODULES, else auto-detect
modules_from_vars=$(yq -r '.vars.MODULES // ""' "$TASKFILE")
if [ -n "$modules_from_vars" ] && [ "$modules_from_vars" != "null" ]; then
  echo "MODULES=$modules_from_vars" >> "$GITHUB_OUTPUT"
  echo "  MODULES=$modules_from_vars (from vars)"
else
  # Simple auto-detection: if go.mod exists in root, return ["."]
  if [ -f "go.mod" ]; then
    echo 'MODULES=["."]' >> "$GITHUB_OUTPUT"
    echo '  MODULES=["."] (auto-detected single module)'
  else
    echo 'MODULES=[]' >> "$GITHUB_OUTPUT"
    echo '  MODULES=[] (no modules found)'
  fi
fi

echo "Done."