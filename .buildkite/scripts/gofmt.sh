#!/bin/bash
set -euo pipefail

if [ "${GO_VERSION:-}" = '1.12' ]; then

    if [ -z "$(gofmt -d .)" ]; then
        true
    else
        gofmt -d . && false
    fi

fi
