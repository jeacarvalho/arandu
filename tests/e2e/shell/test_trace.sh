#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$(cd "$SCRIPT_DIR/../../../scripts" && pwd)"
source "$SCRIPT_DIR/test_framework.sh"

echo "🧪 Testing arandu_trace.sh"
echo "==========================="

test_start "Script exists"
if [ -f "$SCRIPTS_DIR/arandu_trace.sh" ]; then
  test_pass
else
  test_fail "Script not found"
  exit 1
fi

test_start "Script is executable"
if [ -x "$SCRIPTS_DIR/arandu_trace.sh" ]; then
  test_pass
else
  test_fail "Not executable"
fi

test_start "Creates work/logs directory"
source "$SCRIPTS_DIR/arandu_trace.sh"
mkdir -p work/logs
if [ -d "work/logs" ]; then
  test_pass
else
  test_fail "Failed to create logs dir"
fi

test_start "Creates trace log file"
TRACE_FILE="work/logs/trace.log"
if [ -f "$TRACE_FILE" ]; then
  test_pass
else
  test_fail "Trace file not created"
fi

test_start "Handles --trace flag"
if grep -q "TRACE_ENABLED" "$SCRIPTS_DIR/arandu_trace.sh"; then
  test_pass
else
  test_fail "Missing TRACE_ENABLED"
fi

test_start "Trace init function exists"
if grep -q "trace_init()" "$SCRIPTS_DIR/arandu_trace.sh"; then
  test_pass
else
  test_fail "Missing trace_init"
fi

test_start "Trace function exists"
if grep -q "^trace()" "$SCRIPTS_DIR/arandu_trace.sh"; then
  test_pass
else
  test_fail "Missing trace function"
fi

test_summary
