#!/bin/bash
# scripts/arandu_e2e_audit.sh — Entry Point Modular
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/e2e/config.sh"
source "$SCRIPT_DIR/e2e/core.sh"
source "$SCRIPT_DIR/e2e/report.sh"
trap e2e_cleanup EXIT

main() {
    e2e_kill_existing_server; e2e_setup_environment
    e2e_start_server || exit 1
    e2e_create_test_session || exit 1
    
    local total_failed=0
    while IFS='=' read -r module enabled; do
        [[ "$module" =~ ^#.*$ || -z "$module" ]] && continue
        [ "$enabled" != "enabled" ] && continue
        local path="$SCRIPT_DIR/e2e/modules/$module"
        [ ! -f "$path" ] && continue
        source "$path"
        local func="${module%.sh}_module"
        declare -F "$func" >/dev/null && { $func; total_failed=$((total_failed + $?)); }
    done < "$SCRIPT_DIR/e2e/registry.conf"
    
    e2e_generate_report; 
    # No final, antes de exit
    echo ""
    echo "📊 ARQUIVOS GERADOS:"
    echo "===================="
    echo "HTMLs: $(find tmp/audit_logs -name '*.html' 2>/dev/null | wc -l)"
    echo "Screenshots: $(find tmp/audit_screenshots -name '*.png' 2>/dev/null | wc -l)"
    echo ""
    ls -la tmp/audit_logs/*.html 2>/dev/null | head -10
    exit $?
}
main "$@"