#!/bin/bash
# scripts/e2e/report.sh

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

e2e_generate_report() {
    echo ""
    echo "=========================================="
    echo "        E2E AUDIT REPORT"
    echo "=========================================="
    echo ""
    echo "📁 HTMLs gerados:"
    find "$E2E_AUDIT_DIR" -name "*.html" -exec basename {} \; 2>/dev/null | sort
    echo ""
    echo "📊 Resumo por módulo:"
    echo "  - Public:      $PUBLIC_PASSED/$PUBLIC_TOTAL"
    echo "  - Dashboard:   $DASHBOARD_PASSED/$DASHBOARD_TOTAL"
    echo "  - Patients:    $PATIENTS_PASSED/$PATIENTS_TOTAL"
    # ... etc
    echo ""
    echo "✅ Total: $E2E_PASSED_TESTS/$E2E_TOTAL_TESTS"
}