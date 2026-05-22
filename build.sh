#!/usr/bin/env bash
set -euo pipefail

case "$(uname -s)" in
  Darwin)
    if [[ -d /opt/homebrew/include && -d /opt/homebrew/lib ]]; then
      export CGO_CFLAGS="${CGO_CFLAGS:-} -I/opt/homebrew/include"
      export CGO_LDFLAGS="${CGO_LDFLAGS:-} -L/opt/homebrew/lib"
    fi

    if ! command -v brew >/dev/null 2>&1; then
      echo "Homebrew is required on macOS to install dependencies. Install Homebrew: https://brew.sh/"
      exit 1
    fi

    brew install libpcap redis
    ;;
  Linux)
    if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists libpcap; then
      export CGO_CFLAGS="${CGO_CFLAGS:-} $(pkg-config --cflags libpcap)"
      export CGO_LDFLAGS="${CGO_LDFLAGS:-} $(pkg-config --libs libpcap)"
    fi

    if command -v apt-get >/dev/null 2>&1; then
      sudo apt-get update
      sudo apt-get install -y build-essential pkg-config libpcap-dev redis-server
    elif command -v dnf >/dev/null 2>&1; then
      sudo dnf install -y gcc pkgconf-pkg-config libpcap-devel redis
    else
      echo "Unsupported Linux package manager; install libpcap development headers and Redis manually."
    fi
    ;;
  *)
    echo "Unsupported OS: $(uname -s)"
    exit 1
    ;;
esac

if [[ -z "${SKIP_GO_BUILD:-}" ]]; then
  go build -o synkinetik .
fi
