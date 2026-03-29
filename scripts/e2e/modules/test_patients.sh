#!/bin/bash
# scripts/e2e/modules/test_patients.sh
# Módulo: Testes de Pacientes (requer autenticação)
# Versão: 2.0 — Com screenshots e validações SLP

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../core.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/html_validation.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/slp_validation.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../utils/screenshot.sh"

test_patients_module() {
    e2e_log_module "MODULE: Patients"
    local failed=0
    local out patient_id

    # Criar paciente para testes
    e2e_log_info "Creating test patient..."
    out="$E2E_AUDIT_DIR/route_patients_create.html"
    
    local create_response
    create_response=$(curl -s -D - -b "arandu_session=$E2E_COOKIE_VALUE" -X POST "$E2E_BASE_URL/patients/create" \
        -d "name=Test Patient E2E" \
        -d "ethnicity=branca" \
        -d "gender=masculino" \
        -w "\n%{http_code}" 2>&1)
    
    # Salvar response no arquivo
    echo "$create_response" > "$out"
    
    patient_id=$(echo "$create_response" | grep -oP '[a-f0-9-]{36}' | head -1 || echo "")
    
    if [ -z "$patient_id" ]; then
        patient_id=$(echo "$create_response" | grep -oP '/patients/[a-f0-9-]{36}' | head -1 | sed 's|/patients/||' || echo "")
    fi
    
    if [ -z "$patient_id" ]; then
        e2e_log_error "Failed to create patient - cannot proceed with patient tests"
        e2e_log_error "Response: $create_response"
        e2e_test_failed
        return 14
    fi
    
    e2e_log_info "✅ Created patient: $patient_id"
    echo "$patient_id" > "$E2E_AUDIT_DIR/test_patient_id.txt"
    e2e_test_passed

    # GET /patients (COM SCREENSHOT)
    out="$E2E_AUDIT_DIR/route_patients_list.html"
    local code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients")
    if [ "$code" = "200" ]; then
        slp_validate_full "$out" "patients_list" "false" || failed=$((failed + $?))
        
        # ← NOVO: Capturar screenshot da lista de pacientes
        e2e_capture_screenshot "/patients" "$E2E_SCREENSHOT_DIR/patients_list_authenticated.png" 1440 900 || true
        
        e2e_test_passed
    else
        e2e_log_error "patients_list returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/new
    out="$E2E_AUDIT_DIR/route_patients_new.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/new")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "patients_new" "true" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "patients_new returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/search
    out="$E2E_AUDIT_DIR/route_patients_search.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/search?q=test")
    if [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_log_error "patients_search returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/{id} (COM SCREENSHOT)
    out="$E2E_AUDIT_DIR/route_patients_detail.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/$patient_id")
    if [ "$code" = "200" ]; then
        slp_validate_full "$out" "patients_detail" "false" || failed=$((failed + $?))
        
        # ← NOVO: Capturar screenshot do detalhe do paciente
        e2e_capture_screenshot "/patients/$patient_id" "$E2E_SCREENSHOT_DIR/patient_detail_authenticated.png" 1440 900 || true
        
        e2e_test_passed
    else
        e2e_log_error "patients_detail returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/{id}/anamnesis
    out="$E2E_AUDIT_DIR/route_patients_anamnesis.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/$patient_id/anamnesis")
    if [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_log_error "patients_anamnesis returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/{id}/history
    out="$E2E_AUDIT_DIR/route_patients_history.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/$patient_id/history")
    if [ "$code" = "200" ]; then
        e2e_validate_html "$out" "patients_history" "true" || failed=$((failed + 1))
        e2e_test_passed
    else
        e2e_log_error "patients_history returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # GET /patients/{id}/sessions
    out="$E2E_AUDIT_DIR/route_patients_sessions.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" "$E2E_BASE_URL/patients/$patient_id/sessions")
    if [ "$code" = "200" ]; then
        e2e_test_passed
    else
        e2e_log_error "patients_sessions returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # PATCH /patients/{id}/anamnesis/chief_complaint
    out="$E2E_AUDIT_DIR/route_patients_anamnesis_update.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X PATCH \
        "$E2E_BASE_URL/patients/$patient_id/anamnesis/chief_complaint" \
        -d "chief_complaint=test")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "anamnesis_update returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /patients/{id}/medications
    out="$E2E_AUDIT_DIR/route_patients_medications.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/patients/$patient_id/medications" \
        -d "name=Aspirin" -d "dosage=100mg" -d "frequency=daily")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "medications_add returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /patients/{id}/vitals
    out="$E2E_AUDIT_DIR/route_patients_vitals.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/patients/$patient_id/vitals" \
        -d "blood_pressure=120/80" -d "heart_rate=72" -d "temperature=36.5")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "vitals_add returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    # POST /patients/{id}/goals
    out="$E2E_AUDIT_DIR/route_patients_goals.html"
    code=$(curl -s -o "$out" -w "%{http_code}" -b "arandu_session=$E2E_COOKIE_VALUE" -X POST \
        "$E2E_BASE_URL/patients/$patient_id/goals" \
        -d "title=Test Goal" -d "description=Test description")
    if [ "$code" = "200" ] || [ "$code" = "302" ]; then
        e2e_test_passed
    else
        e2e_log_error "goals_create returned $code"
        e2e_test_failed
        failed=$((failed + 1))
    fi

    e2e_log_info "Patients: $((14 - failed))/14 passed"
    return $failed
}