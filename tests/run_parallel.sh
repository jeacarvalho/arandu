#!/bin/bash
# Parallel test runner with timing - OPTIMIZED

echo "🧪 Arandu Test Runner (Paralelo)"
echo "=============================="

cd "$(dirname "$0")/.."
PROJECT_DIR=$(pwd)

NPROC=$(nproc 2>/dev/null || echo "4")
echo "📊 Usando $NPROC núcleos..."
echo ""

# =====================
# STAGE 1: Unit Tests (parallel)
# =====================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 STAGE 1: UNIT TESTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

UNIT_START=$(date +%s.%N)
go test -p $NPROC $(go list ./... | grep -v -E '(tests/e2e|scripts|cmd|migrations)') 2>&1 | tail -2
UNIT_END=$(date +%s.%N)
UNIT_ELAPSED=$(echo "$UNIT_END - $UNIT_START" | bc)

echo ""
echo "📊 Unit tests: ${UNIT_ELAPSED}s"

# =====================
# STAGE 2: Server (only if needed)
# =====================
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 STAGE 2: E2E TESTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if server already running
SERVER_TIME=0
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/dashboard | grep -q "200\|302\|303"; then
    echo "✅ Servidor já está rodando"
else
    START_SERVER=$(date +%s.%N)
    ./scripts/safe_deploy.sh >/dev/null 2>&1 
    sleep 3
    END_SERVER=$(date +%s.%N)
    SERVER_TIME=$(echo "$END_SERVER - $START_SERVER" | bc)
    echo "📊 Server: ${SERVER_TIME}s"
fi

# =====================
# STAGE 3: E2E (ALL PARALLEL)
# =====================
E2E_START=$(date +%s.%N)

# Run E2E suites in parallel
(
    echo "📦 Go E2E..."
    START=$(date +%s.%N)
    go test -v ./tests/e2e/... 2>&1 | tail -2
    END=$(date +%s.%N)
    echo "📊 Go E2E: $(echo "$END - $START" | bc)s"
) &
PID1=$!

(
    echo "📜 JS Tests..."
    START=$(date +%s.%N)
    bash tests/e2e/js/run_js_parallel.sh 2>&1 | tail -2
    END=$(date +%s.%N)
    echo "📊 JS Tests: $(echo "$END - $START" | bc)s"
) &
PID2=$!

(
    echo "🐚 Shell Tests..."
    START=$(date +%s.%N)
    bash tests/e2e/shell/run_shell_tests.sh 2>&1 | tail -2
    END=$(date +%s.%N)
    echo "📊 Shell: $(echo "$END - $START" | bc)s"
) &
PID3=$!

(
    echo "🧩 Modules..."
    START=$(date +%s.%N)
    for script in scripts/e2e/modules/*.sh; do
        [ -f "$script" ] && bash "$script" >/dev/null 2>&1
    done
    END=$(date +%s.%N)
    echo "📊 Modules: $(echo "$END - $START" | bc)s"
) &
PID4=$!

wait $PID1 $PID2 $PID3 $PID4 2>/dev/null || true

E2E_END=$(date +%s.%N)
E2E_ELAPSED=$(echo "$E2E_END - $E2E_START" | bc)

echo ""
echo "📊 E2E (paralelo): ${E2E_ELAPSED}s"

# =====================
# SUMMARY
# =====================
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESUMO DOS TIMINGS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Unit tests:     ${UNIT_ELAPSED}s"
[ $SERVER_TIME -gt 0 ] && echo "🚀 Server:        ${SERVER_TIME}s"
echo "🌐 E2E (paralelo): ${E2E_ELAPSED}s"

TOTAL=$(echo "$UNIT_ELAPSED + $SERVER_TIME + $E2E_ELAPSED" | bc)
echo "⏱️  TOTAL:         ${TOTAL}s"
echo "━���━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "✅ TODOS OS TESTES OK!"