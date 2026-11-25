#!/usr/bin/env bash
# Thin wrapper — all the real logic is in new_service.py (text manipulation
# is far less error-prone in Python than in bash 3.2, which is what ships
# as /bin/bash on macOS). See document/SCAFFOLD_CHECKLIST.md for the flags.
set -eu

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PYTHON_BIN="python3"
if ! command -v "$PYTHON_BIN" >/dev/null 2>&1; then
  PYTHON_BIN="python"
fi

exec "$PYTHON_BIN" "$DIR/new_service.py" "$@"
