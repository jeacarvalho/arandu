#!/bin/bash
set -e

echo "🧪 Unit Tests (Go)"
echo "=================="

cd "$(dirname "$0")/.."

UNIT_PACKAGES=$(go list ./... | grep -v -E '(tests/e2e|scripts|cmd|migrations)')

if [ "$1" = "-cover" ]; then
    echo "📊 Com coverage..."
    go test $UNIT_PACKAGES -coverprofile=coverage.out -covermode=atomic
    go tool cover -func=coverage.out | tail -1
else
    go test $UNIT_PACKAGES
fi

echo "✅ Unit tests concluídos"