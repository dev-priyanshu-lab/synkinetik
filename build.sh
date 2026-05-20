#!/usr/bin/env bash
set -euo pipefail

CGO_CFLAGS="-I/opt/homebrew/include" \
CGO_LDFLAGS="-L/opt/homebrew/lib" \
go build -o synkinetik .
