#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🧪 ARANDU SCRIPTS TEST SUITE"
echo "============================"
echo ""

TOTAL_SUITES=0
SUITES_PASSED=0

run_suite() {
  local name="$1"
  local script="$2"
  
  TOTAL_SUITES=$((TOTAL_SUITES + 1))
  
  if [ -f "$SCRIPT_DIR/$script" ]; then
    echo ""
    bash "$SCRIPT_DIR/$script"
    if [ $? -eq 0 ]; then
      SUITES_PASSED=$((SUITES_PASSED + 1))
    fi
  else
    echo "❌ Suite not found: $script"
  fi
}

run_suite "Checkpoint" "test_checkpoint.sh"
run_suite "Guard" "test_guard.sh"
run_suite "Trace" "test_trace.sh"

echo ""
echo "================================"
echo "📊 Suite Summary"
echo "================================"
echo "  Suites Run:    $TOTAL_SUITES"
echo "  Suites Passed: $SUITES_PASSED"
echo "  Suites Failed: $((TOTAL_SUITES - SUITES_PASSED))"
echo "================================"

if [ $SUITES_PASSED -eq $TOTAL_SUITES ]; then
  echo "✅ ALL TEST SUITES PASSED"
  exit 0
else
  echo "❌ SOME TESTS FAILED"
  exit 1
fi
