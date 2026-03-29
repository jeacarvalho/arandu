#!/bin/bash
# scripts/e2e/modules/test_observations.sh
# Módulo: Testes de Observações (requer autenticação)
# Status: Completo - Não cortar

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"

test_observations_module() {
    e2e_log_module "MODULE: Observations"
    local failed=0
    local out session_id obs_id

    # Obter session_id
    if [ -f "$E2E_AUDIT_DIR/test_session_id.txt" ]; then
        session_id=$(cat "$E2E_AUDIT_DIR/test_session_id.txt")
    else
        e2e_log_warn "No session ID found, skipping observations tests"
        return 0
    fi

    # Criar observação para testes
    e2e_log_info "Creating test observation..."
    out="$E2E_AUDIT_DIR/route_observation_create.html"
    
    local create_response
    create_response=$(curl -s -D - -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session/$session_id/observations" \
        -d "content=Test observation E2E" \
        -d "classification=clinical" \
        -w "\n%{http_code}" 2>&1)
    
    # Salvar response
    echo "$create_response" > "$out"
    
    # Tentar extrair ID de várias formas
    obs_id=$(echo "$create_response" | grep -oP 'hx-get="/observations/[a-f0-9-]{36}' | head -1 | sed 's|hx-get="/observations/||' || echo "")
    
    if [ -z "$obs_id" ]; then
        obs_id=$(echo "$create_response" | grep -oP 'id="observation-[a-f0-9-]{36}' | head -1 | sed 's/id="observation-//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$obs_id" ]; then
        e2e_log_error "Failed to create observation"
        e2e_test_failed
        return 3
    fi
    
    e2e_log_info "✅ Created observation: $obs_id"
    echo "$obs_id" > "$E2E_AUDIT_DIR/test_observation_id.txt"
    e2e_test_passed

    # GET /observations/{id}
    out="$E2E_AUDIT_DIR/route_observation_show.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/observations/$obs_id")
    if [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_log_error "observation_show returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /observations/{id}/edit
    out="$E2E_AUDIT_DIR/route_observation_edit.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/observations/$obs_id/edit")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "observation_edit" "false" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "observation_edit returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # PUT /observations/{id}
    out="$E2E_AUDIT_DIR/route_observation_update.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X PUT \
        "$E2E_BASE_URL/observations/$obs_id" \
        -d "content=Updated observation" -d "classification=clinical")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "observation_update returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Observations: $((4 - failed))/4 passed"
    return $failed
}