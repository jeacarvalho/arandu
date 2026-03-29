#!/bin/bash
# scripts/validate_modules.sh

echo "🔍 Validando estrutura de módulos E2E..."

MODULES_DIR="scripts/e2e/modules"

for module in "$MODULES_DIR"/test_*.sh; do
    [ -f "$module" ] || continue
    
    module_name=$(basename "$module")
    expected_function="${module_name%.sh}_module"
    
    echo "📦 $module_name → Esperado: $expected_function"
    
    # Verificar se a função existe
    if grep -q "^$expected_function()" "$module"; then
        echo "   ✅ Função encontrada"
    else
        echo "   ❌ Função NÃO encontrada!"create_test_session() {
    log_info "=========================================="
    log_info "Creating test user and session"
    log_info "=========================================="
    
    local test_email="${E2E_TEST_EMAIL:-arandu_e2e@test.com}"
    local test_pass="${E2E_TEST_PASS:-test123456}"

    # Signup (pode falhar se usuário já existe - OK)
    curl -s -L -X POST "$E2E_BASE_URL/auth/signup" \
        -d "email=$test_email" \
        -d "password=$test_pass" > /dev/null 2>&1 || true

    # DEBUG: Listar possíveis endpoints de login
    log_info "Testing login endpoints..."
    
    rm -f "$E2E_COOKIES_FILE"
    
    # Tentar múltiplos endpoints COM HEADERS corretos
    local endpoints=("/login" "/auth/login" "/api/login")
    local login_success=false
    
    for endpoint in "${endpoints[@]}"; do
        log_info "  Trying: POST $endpoint"
        
        # CRÍTICO: Usar -c para salvar cookies, -L para seguir redirects
        local response
        response=$(curl -s -v -L -c "$E2E_COOKIES_FILE" \
            -H "Content-Type: application/x-www-form-urlencoded" \
            -X POST "$E2E_BASE_URL$endpoint" \
            -d "email=$test_email" \
            -d "password=$test_pass" \
            2>&1)
        
        # Verificar se cookie foi salvo
        if [ -f "$E2E_COOKIES_FILE" ] && [ -s "$E2E_COOKIES_FILE" ]; then
            if grep -q "arandu_session" "$E2E_COOKIES_FILE" 2>/dev/null; then
                log_info "✅ Login successful via $endpoint"
                login_success=true
                
                # Extrair cookie para uso em Playwright/scripts externos
                export E2E_SESSION_COOKIE=$(grep "arandu_session" "$E2E_COOKIES_FILE" | awk '{print $7}')
                log_info "✅ Session cookie: ${E2E_SESSION_COOKIE:0:20}..."
                break
            fi
        fi
    done
    
    if [ "$login_success" = false ]; then
        log_error "❌ Failed to login via any endpoint"
        log_error "Cookies file content:"
        cat "$E2E_COOKIES_FILE" 2>/dev/null || echo "(empty or not created)"
        log_error "Last response:"
        echo "$response" | tail -20
        return 1
    fi

    log_info "✅ Session ready"
}
        echo "   💡 Funções disponíveis no arquivo:"
        grep -E "^[a-z_]+\(\)" "$module" | sed 's/^/      /'
    fi
done

echo ""
echo "✅ Validação concluída"