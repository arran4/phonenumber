#!/bin/bash
set -euo pipefail

EVENT_NAME="${1:-}"
INPUTS_SNAPSHOT_MODE="${2:-false}"
RELEASE_TAG="${3:-}"
REF="${4:-}"

is_snapshot="false"

if [[ "${EVENT_NAME}" == "workflow_dispatch" && "${INPUTS_SNAPSHOT_MODE}" == "true" ]]; then
  is_snapshot="true"
fi

if [[ "${RELEASE_TAG}" == test-* || "${REF}" == refs/tags/test-* ]]; then
  is_snapshot="true"
fi

echo "is_snapshot=${is_snapshot}"
