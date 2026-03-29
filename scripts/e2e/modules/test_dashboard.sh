#!/bin/bash
# scripts/e2e/modules/test_dashboard.sh
# Módulo: Testes do Dashboard (requer autenticação)
# Versão: 1.0 - Completa com SLP e Screenshot

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/slp_validation.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/screenshot.sh"

test_dashboard_module() {
    e2e_log_module "MODULE: Dashboard"
    local failed=0
    local out

    out="$E2E_AUDIT_DIR/route_root.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_test_failed
        failed=$((failed + 1))
    fi

    out="$E2E_AUDIT_DIR/route_dashboard.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/dashboard")
    
    if [ "$code" = "200" ]; then
        slp_validate_full "$out" "dashboard" "false"
        local slp_result=$?
        if [ $slp_result -gt 0 ]; then
            e2e_log_error "dashboard - $slp_result SLP validation error(s)"
            failed=$((failed + slp_result))
        fi
        
        e2e_capture_screenshot "/dashboard" "$E2E_SCREENSHOT_DIR/dashboard_authenticated.png" 1440 900
        
        e2e_test_passed
    else
        e2e_test_failed
        failed=$((failed + 1))
    fi

    out="$E2E_AUDIT_DIR/route_logout.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/logout")
    if [ "$code" = "302" ] || [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Dashboard: $((3 - failed))/3 passed"
    return $failed
}