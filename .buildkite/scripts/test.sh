#!/usr/bin/env bash
set -euo pipefail

echo "--- Prepare enviroment"
source .buildkite/scripts/pre-install-command.sh
add_bin_path
with_go_junit_report

echo "--- Run the tests"
export OUT_FILE="build/test-report-${GO_VERSION}"
mkdir -p build
set +e
go test -v -race ./... > "${OUT_FILE}"
status=$?
set -e

# Buildkite collapse logs under --- symbols
# need to change --- to anything else or switch off collapsing (note: not available at the moment of this commit)
awk '{gsub("---", "----"); print }' ${OUT_FILE}

exit ${status}
