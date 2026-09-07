#!/bin/bash
set -euo pipefail

MODE="${1:-}"

if [[ "$MODE" == "release-test" || "$MODE" == "release-rc" || "$MODE" == "release-alpha" ]]; then
  echo "snapshot_mode=true"
elif [[ "$MODE" == "release-major" || "$MODE" == "release-minor" || "$MODE" == "release-patch" ]]; then
  echo "snapshot_mode=false"
else
  # Default/fallback
  echo "snapshot_mode=false"
fi
