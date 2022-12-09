#!/usr/bin/env bash
set -euo pipefail

die() {
	echo "ERR: $*"
	exit 1
}

if [ -n "$(git ls-files '*.go' | grep -v -E '.pb.gw.go' | grep -v -e '^docs/' | xargs gofumpt -l 2>/dev/null)" ]; then
	git ls-files '*.go' | grep -v -E '.pb.gw.go' | grep -v -e '^docs/' | xargs gofumpt -d 2>/dev/null
	die "Go formatting errors"
fi

go mod verify
