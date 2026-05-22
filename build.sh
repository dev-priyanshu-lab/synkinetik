#!/usr/bin/env bash
set -euo pipefail

case "$(uname -s)" in
  Darwin)
    if [[ -d /opt/homebrew/include && -d /opt/homebrew/lib ]]; then
      export CGO_CFLAGS="${CGO_CFLAGS:-} -I/opt/homebrew/include"
      export CGO_LDFLAGS="${CGO_LDFLAGS:-} -L/opt/homebrew/lib"
    fi
    ;;
  Linux)
    if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists libpcap; then
      export CGO_CFLAGS="${CGO_CFLAGS:-} $(pkg-config --cflags libpcap)"
      export CGO_LDFLAGS="${CGO_LDFLAGS:-} $(pkg-config --libs libpcap)"
    fi
    ;;
esac

go build -o synkinetik .
