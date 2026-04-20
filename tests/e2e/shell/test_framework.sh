#!/usr/bin/env bash

TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

test_start() {
  local name="$1"
  TESTS_RUN=$((TESTS_RUN + 1))
  echo -n "  🔄 $name... "
}

test_pass() {
  TESTS_PASSED=$((TESTS_PASSED + 1))
  echo "✅"
}

test_fail() {
  local msg="${1:-}"
  TESTS_FAILED=$((TESTS_FAILED + 1))
  echo "❌ $msg"
}

test_summary() {
  echo ""
  echo "================================"
  echo "📊 Test Summary"
  echo "================================"
  echo "  Total:  $TESTS_RUN"
  echo "  Passed: $TESTS_PASSED"
  echo "  Failed: $TESTS_FAILED"
  echo "================================"
  
  if [ $TESTS_FAILED -gt 0 ]; then
    echo "❌ TESTS FAILED"
    return 1
  else
    echo "✅ ALL TESTS PASSED"
    return 0
  fi
}
