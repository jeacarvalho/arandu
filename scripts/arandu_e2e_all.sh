#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🧪 E2E Full Test Suite"
echo "====================="
echo ""

echo "Verificando setup..."
bash "$SCRIPT_DIR/verify_e2e_setup.sh"
[ $? -ne 0 ] && echo "❌ Setup E2E inválido" && exit 1
echo ""

MODULES=(
  "e2e/modules/test_dashboard.sh"
  "e2e/modules/test_patients.sh"
  "e2e/modules/test_sessions.sh"
  "e2e/modules/test_interventions.sh"
  "e2e/modules/test_observations.sh"
  "e2e/modules/test_responsive.sh"
  "e2e/modules/test_public.sh"
)

PASSED=0
FAILED=0

for module in "${MODULES[@]}"; do
  echo "▶️  Executando: $module"
  if bash "$SCRIPT_DIR/$module"; then
    PASSED=$((PASSED + 1))
    echo "   ✅ PASSED"
  else
    FAILED=$((FAILED + 1))
    echo "   ❌ FAILED"
  fi
  echo ""
done

echo "================================"
echo "📊 E2E Test Summary"
echo "================================"
echo "Modules Passed: $PASSED"
echo "Modules Failed: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo "✅ ALL E2E TESTS PASSED"
  exit 0
else
  echo "❌ SOME E2E TESTS FAILED"
  exit 1
fi
