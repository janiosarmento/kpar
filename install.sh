#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="$HOME/.local/bin"
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$PROJECT_DIR/kpar"

echo "Building kpar..."
go build -o "$BINARY" "$PROJECT_DIR/cmd/kpar/"

mkdir -p "$INSTALL_DIR"
ln -sf "$BINARY" "$INSTALL_DIR/kpar"

echo "Installed: $INSTALL_DIR/kpar -> $BINARY"
