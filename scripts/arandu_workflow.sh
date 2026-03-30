#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🎯 Workflow Completo de Validação"
echo "================================="
echo ""

STEP=1

echo "$STEP. Executando checkpoint..."
bash "$SCRIPT_DIR/arandu_checkpoint.sh"
[ $? -ne 0 ] && echo "❌Checkpoint falhou. Corrija os erros." && exit 1
echo ""

STEP=$((STEP + 1))
echo "$STEP. Validando handlers..."
bash "$SCRIPT_DIR/arandu_validate_handlers.sh"
[ $? -ne 0 ] && echo "❌Handlers com problemas." && exit 1
echo ""

STEP=$((STEP + 1))
echo "$STEP. Verificando guard..."
bash "$SCRIPT_DIR/arandu_guard.sh"
[ $? -ne 0 ] && echo "❌Guard falhou." && exit 1
echo ""

STEP=$((STEP + 1))
echo "$STEP. Validando CSS..."
bash "$SCRIPT_DIR/verify_css.sh"
echo ""

STEP=$((STEP + 1))
echo "$STEP. Verificando setup E2E..."
bash "$SCRIPT_DIR/verify_e2e_setup.sh"
echo ""

echo ""
echo "================================"
echo "✅ WORKFLOW COMPLETO PASSOU"
echo "================================"
echo ""
echo "Próximos passos:"
echo "1. Execute testes E2E: ./scripts/e2e/modules/test_*.sh"
echo "2. Faça visual check: ./scripts/arandu_visual_check.sh"
echo "3. Conclua a tarefa: ./scripts/arandu_conclude_task.sh [TASK_ID]"
