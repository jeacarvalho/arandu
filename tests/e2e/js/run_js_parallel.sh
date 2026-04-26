#!/bin/bash
# JS E2E Tests - PARALLEL VERSION

echo "📜 JS E2E Tests (Paralelo)"
echo "========================="

cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"

# Run each test suite in parallel!
(
    echo ">>> baseline..."
    node "$SCRIPT_DIR/test_baseline.js" 2>&1 | tail -3
) &
PID1=$!

(
    echo ">>> dashboard..."
    node "$SCRIPT_DIR/test_dashboard_migration.js" 2>&1 | tail -3
) &
PID2=$!

(
    echo ">>> patients..."
    node "$SCRIPT_DIR/test_patients_list_migration.js" 2>&1 | tail -3
) &
PID3=$!

(
    echo ">>> fase1..."
    node "$SCRIPT_DIR/test_fase1_complete.js" 2>&1 | tail -3
) &
PID4=$!

(
    echo ">>> comparison..."
    node "$SCRIPT_DIR/test_new_patient_comparison.js" 2>&1 | tail -3
) &
PID5=$!

(
    echo ">>> global..."
    node "$SCRIPT_DIR/test_global_final.js" 2>&1 | tail -3
) &
PID6=$!

echo "▶ Executando 6 testes JS em paralelo..."

wait $PID1 $PID2 $PID3 $PID4 $PID5 $PID6

echo ""
echo "✅ JS tests concluídos"