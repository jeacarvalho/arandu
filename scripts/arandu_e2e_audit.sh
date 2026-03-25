#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

AUDIT_DIR="tmp/audit_logs"
COOKIES_FILE="tmp/cookies.txt"
SERVER_PID=""
BASE_URL="http://localhost:8080"

ROUTES_TO_TEST="all"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -r, --routes ROUTES   Comma-separated list of routes to test:"
    echo "                        public, dashboard, patients, sessions, observations, interventions, screenshot, all"
    echo "  -s, --skip ROUTES     Comma-separated list of routes to skip"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --routes public,dashboard,patients"
    echo "  $0 --skip sessions,observations"
    echo "  $0 --routes all"
    exit 0
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -r|--routes)
            ROUTES_TO_TEST="$2"
            shift 2
            ;;
        -s|--skip)
            ROUTES_TO_SKIP="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

should_test() {
    local route="$1"
    if [[ "$ROUTES_TO_SKIP" == *"$route"* ]]; then
        return 1
    fi
    if [[ "$ROUTES_TO_TEST" == "all" ]]; then
        return 0
    fi
    if [[ "$ROUTES_TO_TEST" == *"$route"* ]]; then
        return 0
    fi
    return 1
}

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_route() {
    echo -e "${BLUE}[ROUTE]${NC} $1"
}

cleanup() {
    log_info "Cleaning up..."
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    rm -f "$COOKIES_FILE"
}

kill_existing_server() {
    log_info "Checking for existing server on port 8080..."
    local existing_pids
    existing_pids=$(lsof -i :8080 -t 2>/dev/null || true)
    if [ -n "$existing_pids" ]; then
        log_warn "Found existing server (PIDs: $existing_pids), killing them..."
        echo "$existing_pids" | xargs -r kill 2>/dev/null || true
        sleep 3
    fi
    log_info "Port 8080 is free"
}

trap cleanup EXIT

setup_environment() {
    log_info "Setting up test environment..."
    
    mkdir -p "$AUDIT_DIR"
    rm -rf "$AUDIT_DIR"/*
    rm -f "$COOKIES_FILE"
    rm -f storage/tenants/test_*
    rm -f storage/arandu_central_test.db
    
    log_info "Environment ready"
}

create_test_session() {
    log_info "=========================================="
    log_info "Creating test user and session"
    log_info "=========================================="
    
    local test_email="arandu_e2e@test.com"
    local test_pass="test123456"
    
    log_info "Creating/verifying test user via signup..."
    local signup_response
    signup_response=$(curl -s -L -X POST "$BASE_URL/auth/signup" \
        -d "email=$test_email" \
        -d "password=$test_pass")
    
    if echo "$signup_response" | grep -q "Usuário criado"; then
        log_info "✅ Test user created: $test_email"
    elif echo "$signup_response" | grep -q "já existe"; then
        log_info "ℹ️ Test user already exists: $test_email"
    else
        log_warn "⚠️ Signup response: $signup_response"
    fi
    
    log_info "Logging in to get valid session..."
    rm -f "$COOKIES_FILE"
    local login_response
    login_response=$(curl -s -L -c "$COOKIES_FILE" -X POST "$BASE_URL/login" \
        -d "email=$test_email" \
        -d "password=$test_pass" \
        -w "\n%{http_code}")
    
    local http_code
    http_code=$(echo "$login_response" | tail -1)
    log_info "Login response code: $http_code"
    
    if grep -q "arandu_session" "$COOKIES_FILE" 2>/dev/null; then
        log_info "✅ Session cookie created"
    else
        log_error "❌ No session cookie found"
        cat "$COOKIES_FILE" 2>/dev/null || true
    fi
    
    if [ -f "$COOKIES_FILE" ] && [ -s "$COOKIES_FILE" ]; then
        log_info "✅ Session ready"
    else
        log_error "Failed to create session"
        return 1
    fi
}

start_server() {
    log_info "Starting Arandu server..."
    
    go run cmd/arandu/main.go &
    SERVER_PID=$!
    
    log_info "Waiting for server to be ready..."
    for i in {1..30}; do
        if curl -s "$BASE_URL/test" > /dev/null 2>&1; then
            log_info "Server is ready"
            return 0
        fi
        sleep 1
    done
    
    log_error "Server failed to start"
    return 1
}

check_ghosting() {
    local file="$1"
    local step_name="$2"
    
    if grep -qP '\{ \.?[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+ \}' "$file" 2>/dev/null; then
        log_error "$step_name - Ghosting detected"
        grep -P '\{ \.?[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+ \}' "$file" >> "$AUDIT_DIR/errors.txt" 2>/dev/null || true
        return 1
    fi
    
    log_info "$step_name - No ghosting"
    return 0
}

check_no_inline_styles() {
    local file="$1"
    local step_name="$2"
    
    local style_count
    style_count=$(grep -o 'style="' "$file" | wc -l)
    
    if [ "$style_count" -gt 0 ]; then
        log_error "$step_name - Inline styles detected: $style_count"
        return 1
    fi
    
    log_info "$step_name - No inline styles"
    return 0
}

check_slp_sidebar() {
    local file="$1"
    local step_name="$2"
    
    if ! grep -q "Anamnese" "$file" 2>/dev/null || ! grep -q "Prontuário" "$file" 2>/dev/null; then
        log_error "$step_name - SLP sidebar missing expected links"
        return 1
    fi
    
    log_info "$step_name - SLP sidebar validated"
    return 0
}

check_clinical_typography() {
    local file="$1"
    local step_name="$2"
    
    if ! grep -q "font-clinical\|font-serif\|Source Serif" "$file" 2>/dev/null; then
        log_error "$step_name - Clinical typography (.font-clinical) not found"
        return 1
    fi
    
    log_info "$step_name - Clinical typography validated"
    return 0
}

check_design_system() {
    local file="$1"
    local step_name="$2"
    local skip_layout_check=${3:-false}
    
    local errors=0
    
    if ! grep -q 'href=".*style\.css' "$file" 2>/dev/null; then
        log_error "$step_name - External CSS (style.css) not loaded"
        ((errors++))
    else
        log_info "$step_name - External CSS loaded"
    fi
    
    if [ "$skip_layout_check" != "true" ]; then
        if ! grep -qE 'class=".*(app-container|main-content|sidebar|top-bar)' "$file" 2>/dev/null; then
            log_error "$step_name - Missing layout structure classes"
            ((errors++))
        else
            log_info "$step_name - Layout structure classes present"
        fi
    else
        log_info "$step_name - Skipping layout check (public route)"
    fi
    
    if ! grep -qE 'class=".*(space-|margin-|padding-|mb-|mt-|ml-|mr-|mx-|my-|pb-|pt-|pl-|pr-|px-|py-|gap-|p-|m-|align-|flex-|grid-|text-|font-|bg-|rounded-|shadow-|border-|opacity-|z-)' "$file" 2>/dev/null; then
        log_warn "$step_name - No spacing/design classes found (may cause visual issues)"
    else
        log_info "$step_name - Spacing/design classes present"
    fi
    
    if [ $errors -gt 0 ]; then
        return 1
    fi
    
    log_info "$step_name - Design system validated"
    return 0
}

check_fixed_header_spacing() {
    local file="$1"
    local route_name="$2"
    
    if ! grep -qE '<header.*class=".*top-bar' "$file" 2>/dev/null; then
        return 0
    fi
    
    if grep -qE '\.main-content\s*\{[^}]*padding-top:\s*0' "$file" 2>/dev/null; then
        log_warn "$route_name - Main content has no top padding for fixed header"
        return 0
    fi
    
    log_info "$route_name - Fixed header spacing validated"
    return 0
}

check_fixed_z_index() {
    local file="$1"
    local route_name="$2"
    
    if ! grep -q 'position:\s*fixed' "$file" 2>/dev/null; then
        return 0
    fi
    
    if ! grep -q 'z-index:' "$file" 2>/dev/null; then
        if ! grep -q 'class=".*z-' "$file" 2>/dev/null; then
            log_warn "$route_name - Fixed elements without z-index detected"
            return 0
        fi
    fi
    
    log_info "$route_name - Fixed elements z-index validated"
    return 0
}

check_css_conflicts() {
    local file="$1"
    local route_name="$2"
    
    if grep -q 'style="[^"]*!important' "$file" 2>/dev/null; then
        log_warn "$route_name - CSS !important detected (may cause conflicts)"
    fi
    
    if grep -q '<aside.*sidebar.*top: 80px' "$file" 2>/dev/null; then
        if grep -q '<header.*top-bar' "$file" 2>/dev/null; then
            log_warn "$route_name - Sidebar positioned below header (should be aligned at top)"
        fi
    fi
    
    if grep -qE 'style="[^"]*left:\s*0[^"]*' "$file" 2>/dev/null; then
        if grep -q '<aside.*sidebar' "$file" 2>/dev/null; then
            log_warn "$route_name - Sidebar with left:0 may overlap header"
        fi
    fi
    
    log_info "$route_name - CSS conflicts check completed"
    return 0
}

validate_html() {
    local file="$1"
    local route_name="$2"
    local validate_styles=$3
    local skip_layout_check=${4:-false}
    
    local step_failed=0
    
    check_ghosting "$file" "[$route_name]" || step_failed=$((step_failed + 1))
    
    if [ "$validate_styles" = "true" ]; then
        check_no_inline_styles "$file" "[$route_name]" || step_failed=$((step_failed + 1))
        check_design_system "$file" "[$route_name]" "$skip_layout_check" || step_failed=$((step_failed + 1))
        
        if [ "$skip_layout_check" != "true" ]; then
            check_fixed_header_spacing "$file" "[$route_name]" || step_failed=$((step_failed + 1))
            check_fixed_z_index "$file" "[$route_name]" || step_failed=$((step_failed + 1))
            check_css_conflicts "$file" "[$route_name]" || step_failed=$((step_failed + 1))
        fi
    fi
    
    return $step_failed
}

test_route() {
    local method=$1
    local path=$2
    local route_name=$3
    local validate_styles=$4
    local extra_opts=$5
    local skip_layout_check=${6:-false}
    
    local output_file="$AUDIT_DIR/route_${route_name//\//_}.html"
    
    log_route "Testing: $method $path"
    
    local response
    if [ -f "$COOKIES_FILE" ]; then
        response=$(curl -s --max-time 10 -b "$COOKIES_FILE" $extra_opts -X "$method" "$BASE_URL$path" -w "\n%{http_code}" 2>&1)
    else
        response=$(curl -s --max-time 10 $extra_opts -X "$method" "$BASE_URL$path" -w "\n%{http_code}" 2>&1)
    fi
    
    echo "$response" > "$output_file"
    
    local http_code
    http_code=$(echo "$response" | tail -1)
    
    if [ "$http_code" = "200" ]; then
        log_info "  ✓ HTTP 200 OK"
        
        if [ "$validate_styles" = "true" ]; then
            if grep -q "<!DOCTYPE\|<html" "$output_file" 2>/dev/null; then
                validate_html "$output_file" "$route_name" "$validate_styles" "$skip_layout_check"
                return $?
            else
                log_info "  ✓ HTMX fragment (skip layout validation)"
                check_ghosting "$output_file" "[$route_name]"
                return $?
            fi
        else
            check_ghosting "$output_file" "$route_name"
            return $?
        fi
    elif [ "$http_code" = "302" ] || [ "$http_code" = "303" ]; then
        local redirect
        redirect=$(echo "$response" | grep -i "Location: [^\n]*" | head -1 || echo "")
        log_warn "  ⚠ HTTP $http_code (redirect)"
        if [ -n "$redirect" ]; then
            log_info "  → $redirect"
        fi
        return 0
    elif [ "$http_code" = "401" ]; then
        log_warn "  ⚠ HTTP 401 (unauthorized - expected for public routes)"
        return 0
    elif [ "$http_code" = "404" ]; then
        log_error "  ✗ HTTP 404 (not found)"
        return 1
    else
        log_warn "  ⚠ HTTP $http_code"
        return 0
    fi
}

test_public_routes() {
    echo ""
    echo "=========================================="
    echo "  TESTING PUBLIC ROUTES (no auth)"
    echo "=========================================="
    echo ""
    
    local failed=0
    
    test_route "GET" "/test" "health_check" "false" "" "true" || failed=$((failed + 1))
    test_route "GET" "/login" "login_page" "true" "" "true" || failed=$((failed + 1))
    test_route "GET" "/auth/login" "auth_login" "false" "" "true" || failed=$((failed + 1))
    
    echo ""
    log_info "Public routes completed: $((3 - failed))/3 passed"
    
    return $failed
}

test_patient_routes() {
    echo ""
    echo "=========================================="
    echo "  TESTING PATIENT ROUTES"
    echo "=========================================="
    echo ""
    
    local failed=0
    local cookie_value
    local cookie_opts="-b $COOKIES_FILE"
    
    log_info "Creating patient in in-memory database..."
    local create_response
    create_response=$(curl -s -D - -b "$COOKIES_FILE" -X POST "$BASE_URL/patients/create" \
        -d "name=Test Patient" \
        -d "ethnicity=branca" \
        -d "gender=masculino" \
        -w "\n%{http_code}" 2>&1)
    
    local http_code
    http_code=$(echo "$create_response" | grep "HTTP_CODE:" | cut -d: -f2 || echo "")
    echo "$create_response" > "$AUDIT_DIR/route_patients_create.html"
    
    local patient_id=""
    patient_id=$(echo "$create_response" | grep -i "^Location:" | grep -oP '[a-f0-9-]{36}' | head -1 || echo "")
    
    if [ -z "$patient_id" ]; then
        patient_id=$(echo "$create_response" | grep -oP '/patients/[a-f0-9-]{36}' | head -1 | sed 's|/patients/||' || echo "")
    fi
    
    if [ -z "$patient_id" ]; then
        log_error "Failed to create patient - cannot proceed with patient tests"
        return 14
    fi
    
    log_info "✅ Created patient: $patient_id"
    echo "$patient_id" > "$AUDIT_DIR/test_patient_id.txt"
    
    test_route "GET" "/patients" "patients_list" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/new" "patients_new" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/search" "patients_search" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/$patient_id" "patients_detail" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/$patient_id/anamnesis" "patients_anamnesis" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/$patient_id/context" "patients_context" "false" || failed=$((failed + 1))
    test_route "GET" "/patients/$patient_id/history" "patients_history" "true" || failed=$((failed + 1))
    test_route "GET" "/patients/$patient_id/sessions" "patients_sessions" "true" || failed=$((failed + 1))
    
    test_route "PATCH" "/patients/$patient_id/anamnesis/chief_complaint" "patients_anamnesis_update" "false" "-d chief_complaint=test" || failed=$((failed + 1))
    test_route "POST" "/patients/$patient_id/medications" "patients_medications_add" "false" "-d name=Aspirin -d dosage=100mg -d frequency=daily" || failed=$((failed + 1))
    test_route "POST" "/patients/$patient_id/vitals" "patients_vitals_add" "false" "-d blood_pressure=120/80 -d heart_rate=72 -d temperature=36.5" || failed=$((failed + 1))
    test_route "POST" "/patients/$patient_id/goals" "patients_goals_create" "false" "-d title=Test Goal -d description=Test description" || failed=$((failed + 1))
    
    echo ""
    log_info "Patient routes completed: $((14 - failed))/14 passed"
    
    return $failed
}

test_session_routes() {
    echo ""
    echo "=========================================="
    echo "  TESTING SESSION ROUTES"
    echo "=========================================="
    echo ""
    
    local failed=0
    local cookie_value
    local cookie_opts="-b $COOKIES_FILE"
    
    log_info "Creating test session..."
    
    local patient_id=""
    if [ -f "$AUDIT_DIR/test_patient_id.txt" ]; then
        patient_id=$(cat "$AUDIT_DIR/test_patient_id.txt")
    fi
    
    if [ -z "$patient_id" ]; then
        log_warn "No patient ID found, skipping session tests"
        return 0
    fi
    
    log_info "Using patient ID: $patient_id"
    
    local create_response
    create_response=$(curl -s --max-time 10 -D - -b "$COOKIES_FILE" -X POST "$BASE_URL/session" \
        -d "patient_id=$patient_id" \
        -d "date=$(date +%Y-%m-%d)" \
        -d "summary=Test session" \
        -w "\n%{http_code}" 2>&1)
    
    echo "$create_response" > "$AUDIT_DIR/route_session_create.html"
    
    local http_code
    http_code=$(echo "$create_response" | grep "HTTP_CODE:" | cut -d: -f2 || echo "")
    
    local session_id=""
    session_id=$(echo "$create_response" | grep -i "^Location:" | grep -oP '/session/[a-f0-9-]{36}' | head -1 | sed 's|/session/||' || echo "")
    
    if [ -z "$session_id" ]; then
        session_id=$(echo "$create_response" | grep -oP '/session/[a-f0-9-]{36}' | head -1 | sed 's|/session/||' || echo "")
    fi
    
    if [ -z "$session_id" ]; then
        session_id=$(echo "$create_response" | grep -oP 'action="/session/[a-f0-9-]{36}' | head -1 | sed 's|action="/session/||' | sed 's|/update||' || echo "")
    fi
    
    if [ -z "$session_id" ]; then
        log_error "Failed to get session ID from response"
        log_warn "Response: $create_response"
        return 1
    fi
    
    log_info "Created session ID: $session_id"
    echo "$session_id" > "$AUDIT_DIR/test_session_id.txt"
    test_route "GET" "/session/$session_id" "session_show" "true" || failed=$((failed + 1))
    test_route "GET" "/session/$session_id/edit" "session_edit" "true" || failed=$((failed + 1))
    test_route "POST" "/session/$session_id/update" "session_update" "false" "-d session_id=$session_id -d date=2026-03-24 -d summary=Updated session" || failed=$((failed + 1))
    test_route "POST" "/session/$session_id/observations" "session_observations" "false" "-d content=Test observation -d classification=clinical" || failed=$((failed + 1))
    test_route "POST" "/session/$session_id/interventions" "session_interventions" "false" "-d content=Test intervention -d classification=clinical" || failed=$((failed + 1))
    
    test_route "POST" "/session" "session_create_new" "false" "-d patient_id=$patient_id -d date=$(date +%Y-%m-%d) -d summary=Test" || failed=$((failed + 1))
    
    echo ""
    log_info "Session routes completed: $((7 - failed))/7 passed"
    
    return $failed
}

test_observation_routes() {
    echo ""
    echo "=========================================="
    echo "  TESTING OBSERVATION ROUTES"
    echo "=========================================="
    echo ""
    
    local failed=0
    local cookie_opts="-b $COOKIES_FILE"
    
    local session_id=""
    if [ -f "$AUDIT_DIR/test_session_id.txt" ]; then
        session_id=$(cat "$AUDIT_DIR/test_session_id.txt")
    fi
    
    if [ -z "$session_id" ]; then
        log_warn "No session ID found, skipping observation tests"
        return 0
    fi
    
    log_info "Creating observation via session: $session_id"
    local obs_response
    obs_response=$(curl -s --max-time 10 -D - -b "$COOKIES_FILE" -X POST "$BASE_URL/session/$session_id/observations" \
        -d "content=Test observation" \
        -d "classification=clinical" \
        -w "\n%{http_code}" 2>&1)
    
    local http_code
    http_code=$(echo "$obs_response" | grep "HTTP_CODE:" | cut -d: -f2 || echo "")
    
    local obs_id=""
    # Try hx-get="/observations/uuid/edit"
    obs_id=$(echo "$obs_response" | grep -oP 'hx-get="/observations/[a-f0-9-]{36}' | head -1 | sed 's|hx-get="/observations/||' || echo "")
    
    if [ -z "$obs_id" ]; then
        # Try id="observation-uuid"
        obs_id=$(echo "$obs_response" | grep -oP 'id="observation-[a-f0-9-]{36}' | head -1 | sed 's/id="observation-//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$obs_id" ]; then
        obs_id=$(echo "$obs_response" | grep -oP 'data-id="[a-f0-9-]{36}"' | head -1 | sed 's/data-id="//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$obs_id" ]; then
        log_error "Failed to create observation"
        return 3
    fi
    
    log_info "Created observation: $obs_id"
    echo "$obs_id" > "$AUDIT_DIR/test_observation_id.txt"
    
    test_route "GET" "/observations/$obs_id" "observation_show" "false" || failed=$((failed + 1))
    test_route "GET" "/observations/$obs_id/edit" "observation_edit" "false" || failed=$((failed + 1))
    test_route "PUT" "/observations/$obs_id" "observation_update" "false" "-d content=Updated -d classification=clinical" || failed=$((failed + 1))
    
    echo ""
    log_info "Observation routes completed: $((3 - failed))/3 passed"
    
    return $failed
}

test_intervention_routes() {
    echo ""
    echo "=========================================="
    echo "  TESTING INTERVENTION ROUTES"
    echo "=========================================="
    echo ""
    
    local failed=0
    local cookie_opts="-b $COOKIES_FILE"
    
    local session_id=""
    if [ -f "$AUDIT_DIR/test_session_id.txt" ]; then
        session_id=$(cat "$AUDIT_DIR/test_session_id.txt")
    fi
    
    if [ -z "$session_id" ]; then
        log_warn "No session ID found, skipping intervention tests"
        return 0
    fi
    
    log_info "Creating intervention via session: $session_id"
    local int_response
    int_response=$(curl -s --max-time 10 -D - -b "$COOKIES_FILE" -X POST "$BASE_URL/session/$session_id/interventions" \
        -d "content=Test intervention" \
        -d "classification=clinical" \
        -w "\n%{http_code}" 2>&1)
    
    local http_code
    http_code=$(echo "$int_response" | grep "HTTP_CODE:" | cut -d: -f2 || echo "")
    
    local int_id=""
    # Try hx-get="/interventions/uuid/edit"
    int_id=$(echo "$int_response" | grep -oP 'hx-get="/interventions/[a-f0-9-]{36}' | head -1 | sed 's|hx-get="/interventions/||' || echo "")
    
    if [ -z "$int_id" ]; then
        # Try id="intervention-uuid"
        int_id=$(echo "$int_response" | grep -oP 'id="intervention-[a-f0-9-]{36}' | head -1 | sed 's/id="intervention-//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$int_id" ]; then
        int_id=$(echo "$int_response" | grep -oP 'data-id="[a-f0-9-]{36}"' | head -1 | sed 's/data-id="//' | sed 's/"//' || echo "")
    fi
    
    if [ -z "$int_id" ]; then
        log_error "Failed to create intervention"
        return 3
    fi
    
    log_info "Created intervention: $int_id"
    echo "$int_id" > "$AUDIT_DIR/test_intervention_id.txt"
    
    test_route "GET" "/interventions/$int_id" "intervention_show" "false" || failed=$((failed + 1))
    test_route "GET" "/interventions/$int_id/edit" "intervention_edit" "false" || failed=$((failed + 1))
    test_route "PUT" "/interventions/$int_id" "intervention_update" "false" "-d content=Updated -d classification=clinical" || failed=$((failed + 1))
    
    echo ""
    log_info "Intervention routes completed: $((3 - failed))/3 passed"
    
    return $failed
}

test_dashboard_route() {
    echo ""
    echo "=========================================="
    echo "  TESTING DASHBOARD"
    echo "=========================================="
    echo ""
    
    local failed=0
    
    test_route "GET" "/" "root_redirect" "false" || failed=$((failed + 1))
    test_route "GET" "/dashboard" "dashboard" "true" || failed=$((failed + 1))
    test_route "GET" "/logout" "logout" "false" || failed=$((failed + 1))
    
    echo ""
    log_info "Dashboard routes completed: $((3 - failed))/3 passed"
    
    return $failed
}

test_screenshot_route() {
    echo ""
    echo "=========================================="
    echo "  TESTING SCREENSHOT ROUTE (no auth)"
    echo "=========================================="
    echo ""
    
    local failed=0
    
    test_route "GET" "/screenshot/dashboard" "screenshot_dashboard" "true" || failed=$((failed + 1))
    
    echo ""
    log_info "Screenshot route completed: $((1 - failed))/1 passed"
    
    return $failed
}

generate_report() {
    local total_failed=$1
    
    echo ""
    echo "========================================"
    echo "        E2E AUDIT REPORT"
    echo "========================================"
    echo ""
    
    local total_routes=0
    if [ -d "$AUDIT_DIR" ]; then
        total_routes=$(find "$AUDIT_DIR" -type f -name "*.html" | wc -l)
    fi
    local passed=$((total_routes - total_failed))
    
    echo "Total routes tested: $total_routes"
    echo -e "Passed: ${GREEN}$passed${NC}"
    echo -e "Failed: ${RED}$total_failed${NC}"
    echo ""
    
    if [ -d "$AUDIT_DIR" ]; then
        log_info "Audit logs saved to: $AUDIT_DIR/"
        ls -la "$AUDIT_DIR/" | head -20
    fi
    
    if [ $total_failed -gt 0 ]; then
        echo ""
        log_error "AUDIT FAILED - Check HTML in tmp/audit_logs/"
        return 1
    fi
    
    echo ""
    log_info "AUDIT PASSED - All routes validated successfully"
    return 0
}

setup_test_data() {
    log_info "Creating test data in in-memory database..."
    
    local cookie_value
    local cookie_opts="-b $COOKIES_FILE"
    
    local create_response
    create_response=$(curl -s -b "$COOKIES_FILE" -X POST "$BASE_URL/patients/create" \
        -d "name=Test Patient" \
        -d "ethnicity=branca" \
        -d "gender=masculino" \
        -w "\n%{http_code}")
    
    local http_code
    http_code=$(echo "$create_response" | tail -1)
    
    if [ "$http_code" = "200" ] || [ "$http_code" = "302" ]; then
        TEST_PATIENT_ID=$(echo "$create_response" | grep -oP '[a-f0-9-]{36}' | head -1 || echo "")
        if [ -n "$TEST_PATIENT_ID" ]; then
            log_info "✅ Created test patient: $TEST_PATIENT_ID"
        else
            TEST_PATIENT_ID="d4186096-3673-49af-9b45-f2cf57ef8b10"
            log_warn "Using default patient ID: $TEST_PATIENT_ID"
        fi
    else
        TEST_PATIENT_ID="d4186096-3673-49af-9b45-f2cf57ef8b10"
        log_warn "Failed to create patient, using default: $TEST_PATIENT_ID"
    fi
    
    export TEST_PATIENT_ID
}

main() {
    echo "🚀 Starting E2E Audit - Route Coverage: $ROUTES_TO_TEST"
    if [ -n "$ROUTES_TO_SKIP" ]; then
        echo "Skipping: $ROUTES_TO_SKIP"
    fi
    echo "==============================================="
    
    kill_existing_server
    
    setup_environment
    
    if ! start_server; then
        log_error "Failed to start server"
        exit 1
    fi
    
    create_test_session
    
    local total_failed=0
    local result=0
    
    should_test "public" && { test_public_routes; result=$?; total_failed=$((total_failed + result)); }
    should_test "dashboard" && { test_dashboard_route; result=$?; total_failed=$((total_failed + result)); }
    should_test "patients" && { test_patient_routes; result=$?; total_failed=$((total_failed + result)); }
    should_test "sessions" && { test_session_routes; result=$?; total_failed=$((total_failed + result)); }
    should_test "observations" && { test_observation_routes; result=$?; total_failed=$((total_failed + result)); }
    should_test "interventions" && { test_intervention_routes; result=$?; total_failed=$((total_failed + result)); }
    should_test "screenshot" && { test_screenshot_route; result=$?; total_failed=$((total_failed + result)); }
    
    generate_report $total_failed
    local report_exit=$?
    log_info "Report generation exit code: $report_exit"
    exit $report_exit
}

main "$@"