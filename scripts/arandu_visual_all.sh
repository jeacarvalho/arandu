#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "👁️  Validação Visual Automatizada"
echo "=================================="
echo ""

PASSED=0
FAILED=0

echo "1️⃣  Validando CSS e Layout..."
echo "------------------------------"
if bash "$SCRIPT_DIR/arandu_visual_check.sh" 2>&1 | grep -q "APROVADA\|PASSED"; then
  echo "   ✅ CSS/Layout OK"
  PASSED=$((PASSED + 1))
else
  echo "   ❌ CSS/Layout com problemas"
  FAILED=$((FAILED + 1))
fi
echo ""

echo "2️⃣  Capturando Screenshots..."
echo "---------------------------"
SCREENSHOT_RESULT=0
if bash "$SCRIPT_DIR/arandu_screenshot.sh" 2>&1; then
  CURRENT_COUNT=$(find screenshots/current -name "*.png" 2>/dev/null | wc -l)
  if [ "$CURRENT_COUNT" -gt 0 ]; then
    echo "   ✅ $CURRENT_COUNT screenshots capturados"
    PASSED=$((PASSED + 1))
  else
    echo "   ⚠️  Nenhum screenshot capturado"
    SCREENSHOT_RESULT=1
  fi
else
  echo "   ⚠️  Falha ao capturar screenshots (servidor pode não estar rodando)"
  SCREENSHOT_RESULT=1
fi
echo ""

echo "3️⃣  Verificando Densidade de Layout..."
echo "--------------------------------------"
if bash "$SCRIPT_DIR/analise_densidade_layout.sh" 2>/dev/null | grep -q "WARNING"; then
  echo "   ⚠️  Alta densidade detectada - revise manualmente"
else
  echo "   ✅ Densidade de layout OK"
  PASSED=$((PASSED + 1))
fi
echo ""

echo "4️⃣  Verificando Responsividade..."
echo "--------------------------------"
RESPONSIVE_CHECK=0
for size in "375,667" "768,1024" "1440,900"; do
  WIDTH=$(echo $size | cut -d, -f1)
  HEIGHT=$(echo $size | cut -d, -f2)
  
  if [ "$WIDTH" -eq 375 ]; then
    LABEL="Mobile"
  elif [ "$WIDTH" -eq 768 ]; then
    LABEL="Tablet"
  else
    LABEL="Desktop"
  fi
  
  if grep -q "@media.*max-width.*$WIDTH" web/static/css/*.css 2>/dev/null; then
    echo "   ✅ $LABEL ($WIDTH px) - media query encontrada"
    PASSED=$((PASSED + 1))
  else
    echo "   ⚠️  $LABEL ($WIDTH px) - sem media query específica"
    RESPONSIVE_CHECK=1
  fi
done
echo ""

echo "================================"
echo "📊 Visual Test Summary"
echo "================================"
echo "Passed: $PASSED"
echo "Failed: $FAILED"
echo ""

if [ $FAILED -eq 0 ] && [ $SCREENSHOT_RESULT -eq 0 ]; then
  echo "✅ VALIDAÇÃO VISUAL AUTOMATIZADA PASSOU"
  exit 0
else
  echo "❌ VALIDAÇÃO VISUAL AUTOMATIZADA FALHOU"
  if [ $SCREENSHOT_RESULT -ne 0 ]; then
    echo ""
    echo "Dica: Inicie o servidor antes:"
    echo "  ./scripts/safe_deploy.sh"
    echo "  # ou"
    echo "  go run cmd/arandu/main.go"
  fi
  exit 1
fi
