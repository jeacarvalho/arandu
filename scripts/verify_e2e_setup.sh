#!/bin/bash
# scripts/verify_e2e_setup.sh
echo "🔍 Verificando setup do E2E Audit..."

echo "1. Verificando utils..."
[ -f "scripts/e2e/utils/slp_validation.sh" ] && echo "✅ slp_validation.sh" || echo "❌ slp_validation.sh MISSING"
[ -f "scripts/e2e/utils/screenshot.sh" ] && echo "✅ screenshot.sh" || echo "❌ screenshot.sh MISSING"
[ -f "scripts/e2e/utils/html_validation.sh" ] && echo "✅ html_validation.sh" || echo "❌ html_validation.sh MISSING"

echo ""
echo "2. Verificando módulos..."
[ -f "scripts/e2e/modules/test_dashboard.sh" ] && echo "✅ test_dashboard.sh" || echo "❌ test_dashboard.sh MISSING"

echo ""
echo "3. Verificando Playwright..."
if node -e "require('playwright')" 2>/dev/null; then
    echo "✅ Playwright instalado"
else
    echo "❌ Playwright NÃO instalado"
    echo "   Execute: npm install playwright"
fi

echo ""
echo "4. Verificando fontes nos módulos..."
if grep -q "source.*slp_validation.sh" scripts/e2e/modules/test_dashboard.sh; then
    echo "✅ test_dashboard.sh source slp_validation.sh"
else
    echo "❌ test_dashboard.sh NÃO source slp_validation.sh"
fi

echo ""
echo "✅ Verificação concluída"