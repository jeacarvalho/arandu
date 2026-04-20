#!/bin/bash
set -e

echo "🌐 E2E Tests"
echo "============"

cd "$(dirname "$0")/.."

MODE="${1:-all}"

check_playwright() {
    if command -v playwright &> /dev/null || [ -d "node_modules/playwright" ]; then
        return 0
    fi
    return 1
}

run_js_tests() {
    echo "📜 Executando JS tests..."
    bash tests/e2e/js/run_js_tests.sh "$1"
}

run_go_e2e() {
    echo "📦 Executando testes Go E2E..."
    go test -v ./tests/e2e/...
}

run_playwright() {
    echo "🎭 Executando Playwright..."
    if [ -f "scripts/arandu_e2e_all.sh" ]; then
        bash scripts/arandu_e2e_all.sh
    else
        echo "❌ Playwright script não encontrado"
        exit 1
    fi
}

run_modules() {
    echo "🧩 Executando módulos E2E..."
    for script in scripts/e2e/modules/*.sh; do
        if [ -f "$script" ]; then
            echo "▶ $(basename "$script")"
            bash "$script" || true
        fi
    done
}

run_shell_tests() {
    echo "🐚 Executando shell tests..."
    local mode="${1:-all}"
    bash tests/e2e/shell/run_shell_tests.sh "$mode"
}

case "$MODE" in
    go)
        run_go_e2e
        ;;
    js)
        run_js_tests "$2"
        ;;
    shell)
        run_shell_tests "$2"
        ;;
    pw|playwright)
        if check_playwright; then
            run_playwright
        else
            echo "⚠️ Playwright não instalado, pulando..."
        fi
        ;;
    modules)
        run_modules
        ;;
    all|*)
        run_go_e2e
        if check_playwright; then
            run_playwright
        fi
        run_js_tests "$2"
        run_shell_tests "$2"
        run_modules
        ;;
esac

echo "✅ E2E tests concluídos"