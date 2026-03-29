#!/bin/bash
# scripts/e2e/modules/test_interventions.sh
# Módulo: Testes de Intervenções (requer autenticação)
# Status: Completo - Não cortar

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"

test_interventions_module() {
    e2e_log_module "MODULE: Interventions"
    local failed=0
    local out session_id int_id

    # Obter session_id
    if [ -f "$E2E_AUDIT_DIR/test_session_id.txt" ]; then
        session_id=$(cat "$E2E_AUDIT_DIR/test_session_id.txt")
    else
        e2e_log_warn "No session ID found, skipping interventions tests"
        return 0
    fi

    # Criar intervenção para testes
    e2e_log_info "Creating test intervention..."
    out="$E2E_AUDIT_DIR/route_intervention_create.html"
    
    local create_response
    create_response=$(curl -s -D - -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session/$session_id/interventions" \
        -d "content=Test intervention E2E" \
        -d "classification=clinical" \
        -w "\n%{http_code}" 2>&1)
    
    # Salvar response
    echo "$create_response" > "$out"
    
    # Tentar extrair ID de várias formas
    int_id=$(echo "$create_response" | grep -oP 'hx-get="/interventions/[a-f0-9-]{36}' | head -1 | sed 's|hx-get="/interventions/||' || echo "")
    
    if [ -z "$int_id" ]; then
        int_id=$(echo "$create_response" | grep -oP 'id="intervention-[a-f0-9-]{36}' | head -1 | sed 's/id="intervention-//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$int_id" ]; then
        e2e_log_error "Failed to create intervention"
        e2e_test_failed
        return 3
    fi
    
    e2e_log_info "✅ Created intervention: $int_id"
    echo "$int_id" > "$E2E_AUDIT_DIR/test_intervention_id.txt"
    e2e_test_passed

    # GET /interventions/{id}
    out="$E2E_AUDIT_DIR/route_intervention_show.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/interventions/$int_id")
    if [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_log_error "intervention_show returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /interventions/{id}/edit
    out="$E2E_AUDIT_DIR/route_intervention_edit.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/interventions/$int_id/edit")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "intervention_edit" "false" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "intervention_edit returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # PUT /interventions/{id}
    out="$E2E_AUDIT_DIR/route_intervention_update.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X PUT \
        "$E2E_BASE_URL/interventions/$int_id" \
        -d "content=Updated intervention" -d "classification=clinical")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "intervention_update returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Interventions: $((4 - failed))/4 passed"
    return $failed
}