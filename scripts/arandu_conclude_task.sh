#!/usr/bin/env bash

# Verificar argumentos
if [ $# -lt 2 ]; then
  echo "Uso: $0 TASK_ID \"descrição do aprendizado\" [--success|--failure]"
  echo "Exemplo: $0 20260313_185952 \"Versão inicial\" --success"
  exit 1
fi

TASK_ID=$1
LEARNING=$2
STATUS=${3:-"--success"}

TASK_DIR="work/tasks/task_${TASK_ID}"

# Verificar se o diretório da tarefa existe
if [ ! -d "$TASK_DIR" ]; then
  echo "❌ Diretório da tarefa não encontrado: $TASK_DIR"
  echo "   A tarefa pode já ter sido arquivada."
  exit 1
fi

# Tentar extrair o requirement do arquivo da tarefa
REQ=""
if [ -f "$TASK_DIR/task.md" ]; then
  REQ=$(grep "Requirement:" "$TASK_DIR/task.md" | awk '{print $2}' | head -1)
fi

# Se não encontrou requirement, usar um nome padrão baseado na data
if [ -z "$REQ" ]; then
  REQ="task_${TASK_ID}"
fi

LEARN_FILE="docs/learnings/${REQ}.md"

mkdir -p docs/learnings

# Adicionar aprendizado ao arquivo
echo "## Aprendizado $(date +'%Y-%m-%d %H:%M:%S')" >> "$LEARN_FILE"
echo "**Tarefa:** $TASK_ID" >> "$LEARN_FILE"
echo "**Status:** $STATUS" >> "$LEARN_FILE"
echo "**Conteúdo:**" >> "$LEARN_FILE"
echo "$LEARNING" >> "$LEARN_FILE"
echo "" >> "$LEARN_FILE"

# Mover tarefa para archive
mv "$TASK_DIR" work/archive/ 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ tarefa arquivada"
else
  echo "⚠️  Aviso: Não foi possível mover a tarefa para o archive"
fi

echo "📚 aprendizado registrado em ${LEARN_FILE}"

# Atualizar contexto após conclusão de task
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Atualizando contexto do projeto..."
  bash scripts/arandu_update_context.sh
fi

