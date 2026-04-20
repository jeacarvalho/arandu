#!/bin/bash
set -e

echo "📜 JS E2E Tests"
echo "==============="

cd "$(dirname "$0")"
SCRIPT_DIR="$(pwd)"

MODE="${1:-all}"

run_baseline() {
    echo "📊 Baseline..."
    node "$SCRIPT_DIR/test_baseline.js"
}

run_dashboard_migration() {
    echo "📊 Dashboard Migration..."
    node "$SCRIPT_DIR/test_dashboard_migration.js"
}

run_patients_list_migration() {
    echo "📊 Patients List Migration..."
    node "$SCRIPT_DIR/test_patients_list_migration.js"
}

run_fase1() {
    echo "📊 Fase 1 Complete..."
    node "$SCRIPT_DIR/test_fase1_complete.js"
}

run_comparison() {
    echo "📊 New Patient Comparison..."
    node "$SCRIPT_DIR/test_new_patient_comparison.js"
}

run_global_final() {
    echo "📊 Global Final..."
    node "$SCRIPT_DIR/test_global_final.js"
}

case "$MODE" in
    baseline)
        run_baseline
        ;;
    dashboard)
        run_dashboard_migration
        ;;
    patients)
        run_patients_list_migration
        ;;
    fase1)
        run_fase1
        ;;
    comparison)
        run_comparison
        ;;
    final)
        run_global_final
        ;;
    all)
        run_baseline
        run_dashboard_migration
        run_patients_list_migration
        run_fase1
        run_comparison
        run_global_final
        ;;
    *)
        echo "Uso: $0 [baseline|dashboard|patients|fase1|comparison|final|all]"
        ;;
esac

echo "✅ JS tests concluídos"