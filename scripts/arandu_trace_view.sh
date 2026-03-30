#!/usr/bin/env bash

TRACE_FILE="work/logs/trace.log"

if [ ! -f "$TRACE_FILE" ]; then
  echo "📭 Nenhum trace encontrado."
  echo "   Execute scripts com --trace para ativar o rastreamento."
  exit 0
fi

echo "📍 Trace Log: $TRACE_FILE"
echo "========================="
echo ""

if [ "$1" == "--watch" ]; then
  tail -f "$TRACE_FILE"
else
  cat "$TRACE_FILE"
fi
