#!/bin/bash
# ABOUTME: Script to run the BASIC interpreter using go run
# ABOUTME: Pass any arguments (like .bas files) to the interpreter

set -euo pipefail

# Use workspace-local Go caches to satisfy sandbox write rules
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CACHE_DIR="$ROOT_DIR/.cache"
mkdir -p "$CACHE_DIR/go-build" "$CACHE_DIR/mod" "$CACHE_DIR/tmp"

export GOCACHE="$CACHE_DIR/go-build"
export GOMODCACHE="$CACHE_DIR/mod"
export GOTMPDIR="$CACHE_DIR/tmp"

go run -modcacherw ./cmd/basic "$@"
