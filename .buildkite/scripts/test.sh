#!/usr/bin/env bash
set -euo pipefail

echo "--- Prepare enviroment"
source .buildkite/scripts/pre-install-command.sh
add_bin_path
with_go_junit_report

echo "--- Run the tests"
export OUT_FILE="build/test-report.out"
mkdir -p build
set +e
go test -v -race ./... 2>&1 | tee ${OUT_FILE}
status=$?
go-junit-report > "build/junit-${GO_VERSION}.xml" < ${OUT_FILE}

exit ${status}
