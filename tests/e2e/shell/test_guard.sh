#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "$SCRIPT_DIR/../../../scripts" && pwd)"
source "$SCRIPT_DIR/test_framework.sh"

echo "🧪 Testing arandu_guard.sh"
echo "========================"

test_start "Script exists"
if [ -f "$SCRIPTS_DIR/arandu_guard.sh" ]; then
  test_pass
else
  test_fail "Script not found"
  exit 1
fi

test_start "Script is executable"
if [ -x "$SCRIPTS_DIR/arandu_guard.sh" ]; then
  test_pass
else
  test_fail "Not executable"
fi

test_start "Checks templ generation"
if grep -q "templ generate" "$SCRIPTS_DIR/arandu_guard.sh"; then
  test_pass
else
  test_fail "Missing templ check"
fi

test_start "Verifies _templ.go files exist"
if grep -q "_templ.go" "$SCRIPTS_DIR/arandu_guard.sh"; then
  test_pass
else
  test_fail "Missing _templ.go verification"
fi

test_start "Checks routes with curl"
if grep -q "curl" "$SCRIPTS_DIR/arandu_guard.sh"; then
  test_pass
else
  test_fail "Missing route check"
fi

test_start "Handles missing central.db"
if grep -q "CENTRAL_DB" "$SCRIPTS_DIR/arandu_guard.sh"; then
  test_pass
else
  test_fail "Missing database check"
fi

test_summary

if [ $TESTS_FAILED -eq 0 ]; then
  echo "✅ ALL TESTS PASSED"
  exit 0
else
  echo "❌ SOME TESTS FAILED"
  exit 1
fi