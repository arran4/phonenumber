#!/bin/bash
set -euo pipefail

EVENT_NAME="${1:-}"
INPUTS_SNAPSHOT_MODE="${2:-false}"
RELEASE_TAG="${3:-}"

windows_upload="true"

if [[ "${EVENT_NAME}" == "workflow_dispatch" && "${INPUTS_SNAPSHOT_MODE}" == "true" ]]; then
  windows_upload="false"
fi

if [[ "${RELEASE_TAG}" == test-* ]]; then
  windows_upload="false"
fi

echo "windows_upload=${windows_upload}"
