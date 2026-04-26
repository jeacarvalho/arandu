#!/bin/bash
# Parallel E2E test runner with timing

cd "$(dirname "$0")/.."

echo ""
echo "🌐 E2E Tests (Paralelo)"
echo "==================="

check_playwright() {
    if command -v playwright &> /dev/null || [ -d "node_modules/playwright" ]; then
        return 0
    fi
    return 1
}

# Track results
RESULTS=()

# Parallel execution using background jobs
run_parallel() {
    echo "🚀 Executando testes em paralelo..."
    
    # Start all test suites in background
    (
        echo "📦 Go E2E..."
        START=$(date +%s.%N)
        go test -v ./tests/e2e/... 2>&1 | tail -5
        END=$(date +%s.%N)
        echo "📊 Go E2E: $(echo "$END - $START" | bc)s"
    ) &
    PID_GO=$!

    if check_playwright; then
        (
            echo "🎭 Playwright..."
            START=$(date +%s.%N)
            bash scripts/arandu_e2e_all.sh 2>&1 | tail -10
            END=$(date +%s.%N)
            echo "📊 Playwright: $(echo "$END - $START" | bc)s"
        ) &
        PID_PW=$!
    fi

    (
        echo "📜 JS Tests..."
        START=$(date +%s.%N)
        bash tests/e2e/js/run_js_tests.sh 2>&1 | tail -10
        END=$(date +%s.%N)
        echo "📊 JS Tests: $(echo "$END - $START" | bc)s"
    ) &
    PID_JS=$!

    (
        echo "🐚 Shell Tests..."
        START=$(date +%s.%N)
        bash tests/e2e/shell/run_shell_tests.sh 2>&1 | tail -10
        END=$(date +%s.%N)
        echo "📊 Shell: $(echo "$END - $START" | bc)s"
    ) &
    PID_SHELL=$!

    (
        echo "🧩 Módulos..."
        START=$(date +%s.%N)
        for script in scripts/e2e/modules/*.sh; do
            [ -f "$script" ] && bash "$script" >/dev/null 2>&1
        done
        END=$(date +%s.%N)
        echo "📊 Módulos: $(echo "$END - $START" | bc)s"
    ) &
    PID_MOD=$!

    # Wait for all
    wait $PID_GO 2>/dev/null || true
    wait $PID_PW 2>/dev/null || true  
    wait $PID_JS 2>/dev/null || true
    wait $PID_SHELL 2>/dev/null || true
    wait $PID_MOD 2>/dev/null || true
}

# Sequential fallback (less parallelism)
run_sequential() {
    echo "📦 Go E2E..."
    START=$(date +%s.%N)
    go test -v ./tests/e2e/... 2>&1 | tail -3
    END=$(date +%s.%N)
    echo "📊 Go E2E: $(echo "$END - $START" | bc)s"

    if check_playwright; then
        echo "🎭 Playwright..."
        START=$(date +%s.%N)
        bash scripts/arandu_e2e_all.sh 2>&1 | tail -5
        END=$(date +%s.%N)
        echo "📊 Playwright: $(echo "$END - $START" | bc)s"
    fi

    echo "📜 JS Tests..."
    START=$(date +%s.%N)
    bash tests/e2e/js/run_js_tests.sh 2>&1 | tail -3
    END=$(date +%s.%N)
    echo "📊 JS Tests: $(echo "$END - $START" | bc)s"

    echo "🐚 Shell Tests..."
    START=$(date +%s.%N)
    bash tests/e2e/shell/run_shell_tests.sh 2>&1 | tail -3
    END=$(date +%s.%N)
    echo "📊 Shell: $(echo "$END - $START" | bc)s"

    echo "🧩 Módulos..."
    START=$(date +%s.%N)
    for script in scripts/e2e/modules/*.sh; do
        [ -f "$script" ] && bash "$script" >/dev/null 2>&1
    done
    END=$(date +%s.%N)
    echo "📊 Módulos: $(echo "$END - $START" | bc)s"
}

# Check if we can run in parallel
# Note: Some tests may conflict if they use same DB/port, so use sequential as default
MODE="${1:-sequential}"
if [ "$MODE" = "parallel" ]; then
    run_parallel
else
    run_sequential
fi

echo "✅ E2E tests concluídos"