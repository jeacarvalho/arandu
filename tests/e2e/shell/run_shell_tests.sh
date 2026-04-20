#!/bin/bash
set -e

echo "🐚 Shell E2E Tests"
echo "================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

MODE="${1:-all}"

run_all() {
    echo "🧪 Executando todos os testes shell..."
    bash "$SCRIPT_DIR/test_all.sh"
}

run_checkpoint() {
    echo "📍 Checkpoint..."
    bash "$SCRIPT_DIR/test_checkpoint.sh"
}

run_guard() {
    echo "🛡️ Guard..."
    bash "$SCRIPT_DIR/test_guard.sh"
}

run_trace() {
    echo "📝 Trace..."
    bash "$SCRIPT_DIR/test_trace.sh"
}

run_framework() {
    echo "🔧 Framework..."
    bash "$SCRIPT_DIR/test_framework.sh"
}

run_screenshot() {
    echo "📸 Screenshot..."
    bash "$SCRIPT_DIR/test_screenshot_manual.sh"
}

case "$MODE" in
    all)
        run_all
        ;;
    checkpoint)
        run_checkpoint
        ;;
    guard)
        run_guard
        ;;
    trace)
        run_trace
        ;;
    framework)
        run_framework
        ;;
    screenshot)
        run_screenshot
        ;;
    *)
        echo "Uso: $0 [all|checkpoint|guard|trace|framework|screenshot]"
        ;;
esac

echo "✅ Shell tests concluídos"