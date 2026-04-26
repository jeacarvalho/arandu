#!/bin/bash
# Debug parallel test runner
set -x

cd "$(dirname "$0")/.."

echo "🧪 Parallel Test Runner"
echo "===================="

NPROC=$(nproc 2>/dev/null || echo "4")
echo "📊 Usando $NPROC núcleos..."
echo ""

# Unit tests
echo ""
echo "🧪 Unit tests..."
UNIT_START=$(date +%s.%N)
go test -p $NPROC $(go list ./... | grep -v -E '(tests/e2e|scripts|cmd|migrations)') 2>&1 | tail -1
UNIT_END=$(date +%s.%N)
UNIT_ELAPSED=$(echo "$UNIT_END - $UNIT_START" | bc)
echo "📊 Unit: ${UNIT_ELAPSED}s"

# Check server 
echo ""
echo "🌐 Checking server..."
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/dashboard | grep -q "200\|302\|303"; then
    echo "✅ Server running"
    SERVER_TIME=0
else
    echo "🚀 Starting server..."
    START_SERVER=$(date +%s.%N)
    ./scripts/safe_deploy.sh >/dev/null 2>&1 
    sleep 3
    END_SERVER=$(date +%s.%N)
    SERVER_TIME=$(echo "$END_SERVER - $START_SERVER" | bc)
    echo "📊 Server: ${SERVER_TIME}s"
fi

# E2E parallel - LAUNCH ALL AT ONCE
echo ""
echo "🌐 Running E2E parallel..."

E2E_START=$(date +%s.%N)

# Job 1: Go E2E (runs in background with &)
(echo ">>> Starting Go E2E..." 
 START=$(date +%s.%N)
 go test -v ./tests/e2e/... 2>&1 | tail -3
 END=$(date +%s.%N)
 echo ">>> Go E2E done: $(echo "$END - $START" | bc)s") &
pid1=$!
echo "PID1: $pid1"

# Job 2: JS Tests  
(echo ">>> Starting JS Tests..."
 START=$(date +%s.%N)
 bash tests/e2e/js/run_js_tests.sh 2>&1 | tail -3
 END=$(date +%s.%N)
 echo ">>> JS Tests done: $(echo "$END - $START" | bc)s") &
pid2=$!
echo "PID2: $pid2"

# Job 3: Shell Tests
(echo ">>> Starting Shell Tests..."
 START=$(date +%s.%N)
 bash tests/e2e/shell/run_shell_tests.sh 2>&1 | tail -3
 END=$(date +%s.%N)
 echo ">>> Shell done: $(echo "$END - $START" | bc)s") &
pid3=$!
echo "PID3: $pid3"

# Job 4: Modules
(echo ">>> Starting Modules..."
 START=$(date +%s.%N)
 for script in scripts/e2e/modules/*.sh; do
     [ -f "$script" ] && bash "$script" >/dev/null 2>&1
 done
 END=$(date +%s.%N)
 echo ">>> Modules done: $(echo "$END - $START" | bc)s") &
pid4=$!
echo "PID4: $pid4"

echo "Waiting for all jobs..."

# Wait for ALL at once (not sequentially)
wait $pid1 $pid2 $pid3 $pid4

E2E_END=$(date +%s.%N)
E2E_ELAPSED=$(echo "$E2E_END - $E2E_START" | bc)
echo ""
echo "📊 E2E total: ${E2E_ELAPSED}s"

# Summary
TOTAL=$(echo "$UNIT_ELAPSED + $SERVER_TIME + $E2E_ELAPSED" | bc)
echo ""
echo "⏱️  TOTAL: ${TOTAL}s"