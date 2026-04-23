#!/usr/bin/env bash
set -euo pipefail

OLD="apigo"

if [ $# -ne 1 ]; then
    echo "Usage: $0 <new-name>"
    exit 1
fi

NEW="$1"

if [ "$OLD" = "$NEW" ]; then
    echo "New name is the same as current name, nothing to do."
    exit 0
fi

# Replace module name in all Go and mod files
find . -type f \( -name '*.go' -o -name 'go.mod' \) -exec sed -i '' "s|${OLD}|${NEW}|g" {} +

echo "Renamed module from '${OLD}' to '${NEW}'"
