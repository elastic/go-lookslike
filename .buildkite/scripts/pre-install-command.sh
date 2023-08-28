#!/bin/bash
set -euo pipefail

add_bin_path(){
    mkdir -p "${WORKSPACE}/bin"
    export PATH="${WORKSPACE}/bin:${PATH}"
}

with_go_junit_report() {
    go get -v -u github.com/jstemmer/go-junit-report
}

WORKSPACE=${WORKSPACE:-"$(pwd)"}
