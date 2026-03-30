#!/usr/bin/env bash

TRACE_FILE="work/logs/trace.log"
TRACE_ENABLED=false

if [[ "$*" == *"--trace"* ]]; then
  TRACE_ENABLED=true
fi

trace_init() {
  mkdir -p work/logs
  if [ ! -f "$TRACE_FILE" ]; then
    echo "# Arandu Trace Log" > "$TRACE_FILE"
    echo "# Generated: $(date)" >> "$TRACE_FILE"
    echo "" >> "$TRACE_FILE"
  fi
}

trace() {
  local action="$1"
  local target="$2"
  local details="${3:-}"
  
  if [ "$TRACE_ENABLED" = true ]; then
    echo "[$(date +%Y-%m-%d\ %H:%M:%S)] [$action] $target $details" >> "$TRACE_FILE"
  fi
}

trace_summary() {
  if [ "$TRACE_ENABLED" = true ]; then
    echo "" >> "$TRACE_FILE"
    echo "--- Session Summary ---" >> "$TRACE_FILE"
    echo "Script: $0" >> "$TRACE_FILE"
    echo "Args: $*" >> "$TRACE_FILE"
    echo "Exit: $?" >> "$TRACE_FILE"
    echo "----------------------" >> "$TRACE_FILE"
    echo "" >> "$TRACE_FILE"
  fi
}

trace_init
