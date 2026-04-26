#!/bin/bash
set -e

echo "🧪 Unit Tests (Go)"
echo "=================="

cd "$(dirname "$0")/.."

# Check for -p flag (parallel execution)
if [ "$1" = "-p" ]; then
    PARALLEL_FLAG=1
    shift
fi

# Check for -cover flag
COVER=""
if [ "$1" = "-cover" ]; then
    COVER="-cover"
fi

# Get all test packages
UNIT_PACKAGES=$(go list ./... | grep -v -E '(tests/e2e|scripts|cmd|migrations)')

if [ -n "$PARALLEL_FLAG" ]; then
    # Get number of CPU cores
    NUM_CORES=$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)
    echo "📊 Executando testes em paralelo ($$NUM_CORES núcleos)..."
    if [ "$COVER" = "-cover" ]; then
        go test $UNIT_PACKAGES -p $$NUM_CORES -coverprofile=coverage.out -covermode=atomic
        go tool cover -func=coverage.out | tail -1
    else
        go test $UNIT_PACKAGES -p $$NUM_CORES
    fi
else
    if [ "$COVER" = "-cover" ]; then
        echo "📊 Com coverage..."
        go test $UNIT_PACKAGES -coverprofile=coverage.out -covermode=atomic
        go tool cover -func=coverage.out | tail -1
    else
        go test $UNIT_PACKAGES
    fi
fi

echo "✅ Unit tests concluídos"