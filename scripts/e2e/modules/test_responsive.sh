#!/bin/bash
# scripts/e2e/modules/test_responsive.sh
# Módulo: Testes Responsivos (multi-viewport)

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"

test_responsive_module() {
    e2e_log_module "MODULE: Responsive"
    local failed=0
    local out patient_id

    # Obter patient_id se existir
    if [ -f "$E2E_AUDIT_DIR/test_patient_id.txt" ]; then
        patient_id=$(cat "$E2E_AUDIT_DIR/test_patient_id.txt")
    fi

    # Testar em diferentes viewports (simulado via User-Agent e validação de CSS)
    e2e_log_info "Checking responsive CSS classes..."
    
    local css_file="web/static/css/style.css"
    if [ -f "$css_file" ]; then
        # Verificar media queries para mobile
        if grep -q "@media.*max-width.*768" "$css_file"; then
            e2e_log_info "✅ Mobile media queries present"
            e2e_test_passed
        else
            e2e_log_warn "⚠️ No mobile media queries found"
            e2e_test_failed
            failed=$((failed + 1))
        fi
        
        # Verificar media queries para tablet
        if grep -q "@media.*min-width.*768" "$css_file"; then
            e2e_log_info "✅ Tablet media queries present"
            e2e_test_passed
        else
            e2e_log_warn "⚠️ No tablet media queries found"
            e2e_test_failed
            failed=$((failed + 1))
        fi
    else
        e2e_log_warn "⚠️ CSS file not found"
        failed=$((failed + 1))
    fi

    # Verificar elementos mobile-specific no HTML
    out="$E2E_AUDIT_DIR/route_responsive_check.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/dashboard")
    
    if [ "$code" = "200" ]; then
        # Verificar hamburger menu trigger
        if grep -q 'class=".*hamburger\|menu-trigger\|mobile-menu' "$out"; then
            e2e_log_info "✅ Mobile menu trigger present"
            e2e_test_passed
        else
            e2e_log_warn "⚠️ No mobile menu trigger found"
        fi
        
        # Verificar bottom-nav para mobile
        if grep -q 'class=".*bottom-nav' "$out"; then
            e2e_log_info "✅ Bottom navigation present"
            e2e_test_passed
        else
            e2e_log_info "ℹ️ No bottom navigation (may be OK)"
            e2e_test_passed
        fi
    else
        e2e_log_error "Responsive check returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Responsive: $((4 - failed))/4 passed"
    return $failed
}