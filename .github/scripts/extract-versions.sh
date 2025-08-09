#!/bin/bash

# Script for extracting tool versions from Taskfile.yml
# with help of yq embedded in GitHub Actions runner

set -euo pipefail

# Path to Taskfile.yml
TASKFILE="Taskfile.yml"

# Check if file exists
if [ ! -f "$TASKFILE" ]; then
  echo "Taskfile.yml not found" >&2
  exit 1
fi

echo "Extracting variables from Taskfile.yml using yq..."

# Extract all .vars entries as step outputs
yq -r '.vars | to_entries[] | "\(.key)=\(.value)"' "$TASKFILE" \
  | while IFS= read -r kv; do 
      echo "$kv" >> "$GITHUB_OUTPUT"
      echo "  $kv"
    done

echo "Processing MODULES..."

# MODULES: prefer .vars.MODULES if present, else scan repo for go.mod
modules_from_vars=$(yq -r '.vars.MODULES // ""' "$TASKFILE")
if [ -n "$modules_from_vars" ] && [ "$modules_from_vars" != "null" ]; then
  echo "MODULES=$modules_from_vars" >> "$GITHUB_OUTPUT"
  echo "  MODULES=$modules_from_vars (from Taskfile.yml)"
else
  # Try git-based discovery first (works reliably on checked-out workspace)
  modules_candidates=$(git ls-files -- ':**/go.mod' 2>/dev/null || true)
  # Fallback to find if git returns nothing (e.g., detached or unusual checkout)
  if [ -z "$modules_candidates" ]; then
    modules_candidates=$(find . -path './.*' -prune -o -type f -name 'go.mod' -print)
  fi
  if [ -z "$modules_candidates" ]; then
    if [ -f go.mod ]; then
      modules_scan='["."]'
    else
      modules_scan='[]'
    fi
  else
    echo "modules_candidates raw:\n$modules_candidates"
    modules_list=$(echo "$modules_candidates" \
      | sed 's|/go\.mod$||' \
      | sort -u \
      | grep -v '^$' \
      | sed 's|^\./$|.|' \
      | sed 's/^/"/; s/$/"/' \
      | paste -sd ',' -)
    modules_scan="[$modules_list]"
    if [ "$modules_scan" = "[]" ] && [ -f go.mod ]; then
      modules_scan='["."]'
    fi
  fi
  echo "MODULES=$modules_scan" >> "$GITHUB_OUTPUT"
  echo "  MODULES=$modules_scan (auto-detected)"
fi

echo "Variables extraction completed."