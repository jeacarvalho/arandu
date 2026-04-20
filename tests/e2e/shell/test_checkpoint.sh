#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "$SCRIPT_DIR/../../../scripts" && pwd)"
source "$SCRIPT_DIR/test_framework.sh"

echo "🧪 Testing arandu_checkpoint.sh"
echo "==============================="

test_start "Script exists"
if [ -f "$SCRIPTS_DIR/arandu_checkpoint.sh" ]; then
  test_pass
else
  test_fail "Script not found"
  exit 1
fi

test_start "Script is executable"
if [ -x "$SCRIPTS_DIR/arandu_checkpoint.sh" ]; then
  test_pass
else
  test_fail "Not executable"
fi

test_start "Calls arandu_validate_handlers.sh"
if grep -q "arandu_validate_handlers.sh" "$SCRIPTS_DIR/arandu_checkpoint.sh"; then
  test_pass
else
  test_fail "Missing handler validation call"
fi

test_start "Checks go build for handlers"
if grep -q "go build ./internal/web/handlers" "$SCRIPTS_DIR/arandu_checkpoint.sh"; then
  test_pass
else
  test_fail "Missing handlers build check"
fi

test_start "Checks go build for main"
if grep -q "go build ./cmd/arandu" "$SCRIPTS_DIR/arandu_checkpoint.sh"; then
  test_pass
else
  test_fail "Missing main build check"
fi

test_start "Checks HTML inline anti-pattern"
if grep -q "HTML_INLINE_COUNT" "$SCRIPTS_DIR/arandu_checkpoint.sh"; then
  test_pass
else
  test_fail "Missing HTML inline check"
fi

test_start "Validates migrations exist"
if grep -q "migrations/0001_initial_schema.up.sql" "$SCRIPTS_DIR/arandu_checkpoint.sh"; then
  test_pass
else
  test_fail "Missing migration check"
fi

test_summary
