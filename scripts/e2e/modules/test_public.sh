#!/bin/bash
# scripts/e2e/modules/test_public.sh

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"

test_public_module() {
    e2e_log_module "MODULE: Public Routes"
    local failed=0 out
    
    # /test
    out="$E2E_AUDIT_DIR/route_test.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" "$E2E_BASE_URL/test")
    [ "$code" = "200" ] && { e2e_test_passed; } || { e2e_test_failed; failed=$((failed+1)); }
    
    # /login
    out="$E2E_AUDIT_DIR/route_login.html"
    code=$(curl -s -o "$out" -w "%{http_code}" "$E2E_BASE_URL/login")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "login" "false" || failed=$((failed+1))
        e2e_test_passed
    else
        e2e_test_failed; failed=$((failed+1))
    fi
    
    e2e_log_info "Public: $((3-failed))/3 passed"; return $failed
}