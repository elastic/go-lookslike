#!/bin/bash
set -euo pipefail

if [ -z "$(gofmt -d .)" ]; then
    true
else
    gofmt -d . && false
fi
