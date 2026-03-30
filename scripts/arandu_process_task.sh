#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/arandu_trace.sh" "$@"

trace "PROCESS_TASK" "$TASK_ID" "Starting"

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 TASK_ID [--execute]"
  echo "Exemplo: $0 20260313_185952"
  echo "Exemplo: $0 20260313_185952 --execute"
  exit 1
fi

TASK_ID=$1
EXECUTE=${2:-""}

TASK_DIR="work/tasks/task_${TASK_ID}"

if [ ! -d "$TASK_DIR" ]; then
  echo "❌ Tarefa não encontrada: $TASK_ID"
  echo "   Verifique se o diretório existe: $TASK_DIR"
  exit 1
fi

echo "📋 Tarefa: ${TASK_ID}"
echo "📁 Diretório: ${TASK_DIR}"
echo ""

# Verificar se o arquivo task.md existe
if [ ! -f "$TASK_DIR/task.md" ]; then
  echo "❌ Arquivo task.md não encontrado na tarefa"
  echo "   Crie o arquivo $TASK_DIR/task.md com os detalhes da implementação primeiro."
  exit 1
fi

# Verificar se o task.md tem conteúdo além do template básico
TASK_CONTENT=$(cat "$TASK_DIR/task.md")
BASIC_TEMPLATE_LINES=$(cat "$TASK_DIR/task.md" | wc -l)

# Se o arquivo tem apenas o template básico (13 linhas ou menos), pedir para preencher
if [ "$BASIC_TEMPLATE_LINES" -le 13 ]; then
  echo "⚠️  Arquivo task.md contém apenas o template básico."
  echo ""
  echo "📝 Por favor, preencha o arquivo $TASK_DIR/task.md com:"
  echo "   1. Detalhes específicos da implementação"
  echo "   2. Passos a serem seguidos"
  echo "   3. Requisitos técnicos"
  echo "   4. Restrições de design"
  echo ""
  echo "📋 Conteúdo atual do task.md:"
  echo "---"
  cat "$TASK_DIR/task.md"
  echo "---"
  echo ""
  echo "❌ Não é possível processar a tarefa até que task.md seja preenchido."
  exit 1
fi

# Mostrar conteúdo da tarefa
echo "✅ Tarefa pronta para processamento:"
echo "---"
cat "$TASK_DIR/task.md"
echo "---"

echo ""
echo "⚠️  Antes de implementar, leia o requirement correspondente"

# Se o parâmetro --execute foi fornecido
if [ "$EXECUTE" = "--execute" ]; then
  echo ""
  echo "🚀 Modo execute ativado"
  echo "   Esta funcionalidade executaria a tarefa automaticamente."
  echo "   (Implementação futura)"
fi
