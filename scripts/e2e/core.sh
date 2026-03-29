#!/bin/bash
# scripts/e2e/core.sh
# Funções compartilhadas

source "$(dirname "${BASH_SOURCE[0]}")/config.sh"

e2e_log_info() { echo -e "${E2E_COLOR_GREEN}[INFO]${E2E_COLOR_NC} $1"; }
e2e_log_warn() { echo -e "${E2E_COLOR_YELLOW}[WARN]${E2E_COLOR_NC} $1"; }
e2e_log_error() { echo -e "${E2E_COLOR_RED}[ERROR]${E2E_COLOR_NC} $1"; }
e2e_log_module() { echo -e "${E2E_COLOR_BLUE}[MODULE]${E2E_COLOR_NC} $1"; }

e2e_kill_existing_server() {
    e2e_log_info "Checking for existing server on port $E2E_PORT..."
    local pids=$(lsof -i :$E2E_PORT -t 2>/dev/null || true)
    [ -n "$pids" ] && echo "$pids" | xargs -r kill 2>/dev/null && sleep 3
}

e2e_start_server() {
    e2e_log_info "Starting Arandu server..."
    go run cmd/arandu/main.go > "$E2E_AUDIT_DIR/server.log" 2>&1 &
    export E2E_SERVER_PID=$!
    for i in $(seq 1 $E2E_SERVER_TIMEOUT); do
        curl -s "$E2E_BASE_URL/test" > /dev/null 2>&1 && { e2e_log_info "Server ready"; return 0; }
        sleep 1
    done
    e2e_log_error "Server failed to start"; return 1
}

e2e_cleanup() {
    [ -n "$E2E_SERVER_PID" ] && kill -0 "$E2E_SERVER_PID" 2>/dev/null && kill "$E2E_SERVER_PID" 2>/dev/null
    rm -f "$E2E_COOKIES_FILE"
}

e2e_setup_environment() {
    mkdir -p "$E2E_AUDIT_DIR" "$E2E_SCREENSHOT_DIR"
    rm -rf "$E2E_AUDIT_DIR"/* "$E2E_SCREENSHOT_DIR"/* "$E2E_COOKIES_FILE"
}

e2e_create_test_session() {
    e2e_log_info "Creating test session..."
    curl -s -L -X POST "$E2E_BASE_URL/auth/signup" -d "email=$E2E_TEST_EMAIL" -d "password=$E2E_TEST_PASS" > /dev/null 2>&1 || true
    rm -f "$E2E_COOKIES_FILE"
    curl -s -L -c "$E2E_COOKIES_FILE" -X POST "$E2E_BASE_URL/login" -d "email=$E2E_TEST_EMAIL" -d "password=$E2E_TEST_PASS" > /dev/null 2>&1
    if grep -q "arandu_session" "$E2E_COOKIES_FILE" 2>/dev/null; then
        export E2E_COOKIE_VALUE=$(grep "arandu_session" "$E2E_COOKIES_FILE" | awk '{print $NF}')
        e2e_log_info "✅ Session cookie extracted"
    else
        e2e_log_error "❌ Failed to extract session cookie"; return 1
    fi
}

e2e_test_passed() { export E2E_TOTAL_TESTS=$((E2E_TOTAL_TESTS+1)); export E2E_PASSED_TESTS=$((E2E_PASSED_TESTS+1)); }
e2e_test_failed() { export E2E_TOTAL_TESTS=$((E2E_TOTAL_TESTS+1)); export E2E_FAILED_TESTS=$((E2E_FAILED_TESTS+1)); }