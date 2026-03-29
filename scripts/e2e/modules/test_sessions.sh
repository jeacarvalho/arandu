#!/bin/bash
# scripts/e2e/modules/test_sessions.sh
# Módulo: Testes de Sessões (requer autenticação)
# Status: Completo - Não cortar

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"

test_sessions_module() {
    e2e_log_module "MODULE: Sessions"
    local failed=0
    local out patient_id session_id

    # Obter patient_id
    if [ -f "$E2E_AUDIT_DIR/test_patient_id.txt" ]; then
        patient_id=$(cat "$E2E_AUDIT_DIR/test_patient_id.txt")
    else
        e2e_log_warn "No patient ID found, creating one..."
        local create_response
        create_response=$(curl -s -b "arandu_session=$E2E_COOKIE_VALUE" -X POST "$E2E_BASE_URL/patients/create" \
            -d "name=Session Test Patient" -w "\n%{http_code}")
        patient_id=$(echo "$create_response" | grep -oP '[a-f0-9-]{36}' | head -1 || echo "")
        echo "$patient_id" > "$E2E_AUDIT_DIR/test_patient_id.txt"
    fi

    if [ -z "$patient_id" ]; then
        e2e_log_error "Failed to get patient ID"
        return 1
    fi

    # Criar sessão para testes
    e2e_log_info "Creating test session..."
    out="$E2E_AUDIT_DIR/route_session_create.html"
    
    local create_response
    create_response=$(curl -s -D - -b "arandu_session=$E2E_COOKIE_VALUE" -X POST "$E2E_BASE_URL/session" \
        -d "patient_id=$patient_id" \
        -d "date=$(date +%Y-%m-%d)" \
        -d "summary=Test session E2E" \
        -w "\n%{http_code}" 2>&1)
    
    # Salvar response
    echo "$create_response" > "$out"
    
    session_id=$(echo "$create_response" | grep -oP '/session/[a-f0-9-]{36}' | head -1 | sed 's|/session/||' || echo "")
    
    if [ -z "$session_id" ]; then
        e2e_log_error "Failed to create session"
        e2e_test_failed
        return 1
    fi
    
    e2e_log_info "✅ Created session: $session_id"
    echo "$session_id" > "$E2E_AUDIT_DIR/test_session_id.txt"
    e2e_test_passed

    # GET /session/{id}
    out="$E2E_AUDIT_DIR/route_session_show.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/session/$session_id")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "session_show" "true" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "session_show returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /session/{id}/edit
    out="$E2E_AUDIT_DIR/route_session_edit.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/session/$session_id/edit")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "session_edit" "true" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "session_edit returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /session/{id}/update
    out="$E2E_AUDIT_DIR/route_session_update.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session/$session_id/update" \
        -d "session_id=$session_id" -d "date=2026-03-24" -d "summary=Updated session")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "session_update returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /session/{id}/observations
    out="$E2E_AUDIT_DIR/route_session_observations.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session/$session_id/observations" \
        -d "content=Test observation" -d "classification=clinical")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "session_observations returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /session/{id}/interventions
    out="$E2E_AUDIT_DIR/route_session_interventions.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session/$session_id/interventions" \
        -d "content=Test intervention" -d "classification=clinical")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "session_interventions returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /session (nova sessão)
    out="$E2E_AUDIT_DIR/route_session_create_new.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/session" \
        -d "patient_id=$patient_id" -d "date=$(date +%Y-%m-%d)" -d "summary=New test session")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "session_create_new returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Sessions: $((7 - failed))/7 passed"
    return $failed
}