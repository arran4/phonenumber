#!/bin/bash
set -euo pipefail

EVENT_NAME="${1:-}"
REF="${2:-}"
REF_TYPE="${3:-}"
EVENT_ACTION="${4:-}"
INPUTS_MODE="${5:-}"
EVENT_SCHEDULE="${6:-}"

run_code_checks="false"
run_pr_meta_checks="false"
run_cleanup="false"
run_release="false"
run_post_release="false"
is_monthly="false"
is_nightly="false"

case "${EVENT_NAME}" in
  push)
    run_code_checks="true"
    if [[ "${REF}" == refs/tags/v* || "${REF}" == refs/tags/test-* ]]; then
      run_release="true"
    fi
    ;;
  pull_request)
    if [[ "${EVENT_ACTION}" == "closed" ]]; then
      run_cleanup="true"
    else
      run_pr_meta_checks="true"
      run_code_checks="true"
    fi
    ;;
  workflow_dispatch)
    case "${INPUTS_MODE}" in
      lint-fix)
        run_code_checks="true"
        is_nightly="true"
        ;;
      build)
        run_code_checks="true"
        ;;
      release-*)
        run_code_checks="true"
        run_release="true"
        ;;
      publish-tag)
        if [[ "${REF_TYPE}" != "tag" || (! "${REF}" =~ ^refs/tags/v.* && ! "${REF}" =~ ^refs/tags/test-.*) ]]; then
          echo "publish-tag mode requires an eligible tag context" >&2
          exit 1
        fi
        run_code_checks="true"
        run_release="true"
        ;;
      monthly-maintenance)
        run_code_checks="true"
        is_monthly="true"
        ;;
    esac
    ;;
  schedule)
    run_code_checks="true"
    if [[ "${EVENT_SCHEDULE}" == "17 3 1 * *" ]]; then
      is_monthly="true"
    else
      is_nightly="true"
    fi
    ;;
  release)
    run_post_release="true"
    ;;
esac

echo "run_code_checks=$run_code_checks"
echo "run_pr_meta_checks=$run_pr_meta_checks"
echo "run_cleanup=$run_cleanup"
echo "run_release=$run_release"
echo "run_post_release=$run_post_release"
echo "is_monthly=$is_monthly"
echo "is_nightly=$is_nightly"
