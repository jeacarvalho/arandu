#!/bin/bash
# scripts/debug_login.sh

echo "🔍 Debug: Testando autenticação Arandu"
echo "======================================="

BASE_URL="http://localhost:8080"
COOKIES="/tmp/debug_cookies.txt"
EMAIL="arandu_e2e@test.com"
PASS="test123456"

rm -f "$COOKIES"

echo "1. Tentando POST /login..."
curl -v -L -c "$COOKIES" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -X POST "$BASE_URL/login" \
    -d "email=$EMAIL" \
    -d "password=$PASS" \
    2>&1 | grep -E "Set-Cookie|Location|< HTTP|arandu_session"

echo ""
echo "2. Conteúdo do arquivo de cookies:"
cat "$COOKIES" 2>/dev/null || echo "(arquivo não criado)"

echo ""
echo "3. Testando acesso ao dashboard com cookie..."
if grep -q "arandu_session" "$COOKIES" 2>/dev/null; then
    curl -s -b "$COOKIES" -o /tmp/dashboard.html "$BASE_URL/dashboard"
    if grep -q "Arandu\|dashboard\|Pacientes" /tmp/dashboard.html 2>/dev/null; then
        echo "✅ Dashboard acessado com sucesso!"
    else
        echo "❌ Dashboard não carregou conteúdo esperado"
        head -20 /tmp/dashboard.html
    fi
else
    echo "❌ Cookie de sessão não encontrado"
fi