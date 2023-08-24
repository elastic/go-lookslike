#!/bin/bash

set -euo pipefail

echo "--- Pre install"
source .buildkite/scripts/pre-install-command.sh
add_bin_path
with_go_junit_report

# Create Junit report for junit annotation plugin
build_folder="/build"
mkdir $build_folder
buildkite-agent artifact download "build/test-report-*" "${build_folder}" --step test-matrix
find ./build -name "test-report-*" -exec sh -c 'f=$1; go-junit-report < ${f} >> ${f}.xml' shell {}  \;
