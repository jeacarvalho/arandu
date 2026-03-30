#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "👁️  Validação Visual Completa"
echo "=============================="
echo ""

PASSED=0
FAILED=0

echo "1️⃣  Validando CSS e Layout..."
echo "------------------------------"
if bash "$SCRIPT_DIR/arandu_visual_check.sh"; then
  echo "   ✅ CSS/Layout OK"
  PASSED=$((PASSED + 1))
else
  echo "   ❌ CSS/Layout com problemas"
  FAILED=$((FAILED + 1))
fi
echo ""

echo "2️⃣  Verificando Screenshots..."
echo "-----------------------------"
if [ -d "screenshots" ]; then
  CURRENT_COUNT=$(find screenshots/current -name "*.png" 2>/dev/null | wc -l)
  BASELINE_COUNT=$(find screenshots/baseline -name "*.png" 2>/dev/null | wc -l)
  echo "   📸 Screenshots atuais: $CURRENT_COUNT"
  echo "   📋 Baseline: $BASELINE_COUNT"
  if [ "$CURRENT_COUNT" -gt 0 ]; then
    echo "   ✅ Screenshots encontrados"
    PASSED=$((PASSED + 1))
  else
    echo "   ⚠️  Nenhum screenshot encontrado"
    echo "   💡 Execute: ./scripts/arandu_screenshot.sh"
  fi
else
  echo "   ⚠️  Diretório screenshots não existe"
  echo "   💡 Execute: ./scripts/arandu_screenshot.sh"
fi
echo ""

echo "3️⃣  Checklist Visual Manual"
echo "---------------------------"
echo "   Responda as perguntas:"
echo ""
read -p "   [1] Testou em desktop (1920px)? (s/n): " DESKTOP
read -p "   [2] Testou em mobile (375px)? (s/n): " MOBILE
read -p "   [3] Testou em tablet (768px)? (s/n): " TABLET
read -p "   [4] Verificou que não há elementos sobrepostos? (s/n): " OVERLAP
read -p "   [5] Verificou que scroll funciona sem corte? (s/n): " SCROLL
read -p "   [6] Verificou contraste de cores? (s/n): " CONTRAST
echo ""

MANUAL_OK=true
if [[ "$DESKTOP" != "s" && "$DESKTOP" != "S" ]]; then
  echo "   ❌ Teste desktop obrigatório"
  MANUAL_OK=false
fi
if [[ "$MOBILE" != "s" && "$MOBILE" != "S" ]]; then
  echo "   ❌ Teste mobile obrigatório"
  MANUAL_OK=false
fi

if [ "$MANUAL_OK" = true ]; then
  echo "   ✅ Checklist manual aprovado"
  PASSED=$((PASSED + 1))
else
  echo "   ❌ Checklist manual incompleto"
  FAILED=$((FAILED + 1))
fi
echo ""

echo "4️⃣  Verificando Densidade de Layout..."
echo "--------------------------------------"
if bash "$SCRIPT_DIR/analise_densidade_layout.sh" 2>/dev/null | grep -q "WARNING"; then
  echo "   ⚠️  Alguns arquivos com alta densidade de layout"
  echo "   💡 Revise: ./scripts/analise_densidade_layout.sh"
else
  echo "   ✅ Densidade de layout OK"
  PASSED=$((PASSED + 1))
fi
echo ""

echo "================================"
echo "📊 Visual Test Summary"
echo "================================"
echo "Passed: $PASSED"
echo "Failed: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo "✅ VALIDAÇÃO VISUAL COMPLETA PASSOU"
  exit 0
else
  echo "❌ VALIDAÇÃO VISUAL COMPLETA FALHOU"
  echo ""
  echo "Próximos passos:"
  echo "1. Execute: ./scripts/arandu_screenshot.sh"
  echo "2. Revise manualmente os screenshots em screenshots/current/"
  echo "3. Execute: ./scripts/arandu_visual_check.sh para detalhes"
  exit 1
fi
