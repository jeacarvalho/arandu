#!/usr/bin/env bash

TASK_TITLE="${1:-}"
REQ_ID="${2:-}"

if [ -z "$TASK_TITLE" ]; then
  echo "🎯 Cycle de Tarefa Completo"
  echo "==========================="
  echo ""
  echo "Uso: $0 \"Título da tarefa\" [REQ-ID]"
  echo ""
  echo "Exemplo 1: $0 \"Implementar feature X\""
  echo "Exemplo 2: $0 \"Criar componente Y\" req-01-01-01"
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🎯 Starting Task Cycle"
echo "====================="
echo ""
echo "Tarefa: $TASK_TITLE"
[ -n "$REQ_ID" ] && echo "Requirement: $REQ_ID"
echo ""

TASK_ID=$(date +%Y%m%d_%H%M%S)

echo "1️⃣  Criando tarefa..."
bash "$SCRIPT_DIR/arandu_create_task.sh" "$TASK_TITLE" $REQ_ID
echo ""

echo "📝 INSTRUÇÕES:"
echo "=============="
echo "1. Edite o arquivo de tarefa criado"
echo "2. Remova a seção 'Instruções para o Agente'"
echo "3. Atualize o status para 'PRONTO_PARA_IMPLEMENTACAO'"
echo ""
read -p "Pressione ENTER quando a tarefa estiver pronta..."

echo ""
echo "4️⃣  Processando tarefa..."
bash "$SCRIPT_DIR/arandu_process_task.sh" "$TASK_ID"
echo ""

echo "5️⃣  Validando integridade..."
bash "$SCRIPT_DIR/arandu_checkpoint.sh"
[ $? -ne 0 ] && echo "⚠️  checkpoint com problemas" && exit 1
echo ""

echo "6️⃣  CHECKLIST VISUAL OBRIGATÓRIO"
echo "================================"
echo "Responda:"
read -p "Testou em desktop (1920px)? (s/n): " DESKTOP
read -p "Testou em mobile (375px)? (s/n): " MOBILE
read -p "Verificou que não há elementos sobrepostos? (s/n): " OVERLAP
read -p "Verificou scroll sem corte de conteúdo? (s/n): " SCROLL
read -p "Gerou e revisou screenshots? (s/n): " SCREEN

if [[ "$DESKTOP" != "s" && "$DESKTOP" != "S" ]]; then
  echo "❌ Teste desktop obrigatório"
  exit 1
fi
if [[ "$MOBILE" != "s" && "$MOBILE" != "S" ]]; then
  echo "❌ Teste mobile obrigatório"
  exit 1
fi

echo ""
echo "7️⃣  Concluindo tarefa..."
bash "$SCRIPT_DIR/arandu_conclude_task.sh" "$TASK_ID" --success
echo ""

echo "================================"
echo "✅ TASK CYCLE COMPLETO"
echo "================================"
