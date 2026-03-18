#!/usr/bin/env bash

set -e

echo "🧠 Atualizando contexto permanente do projeto Arandu..."
echo

CONTEXT_DIR="docs/project_context"
SESSION_CONTEXT_DIR="docs/agent_context"
PERMANENT_FILE="${CONTEXT_DIR}/permanent_state.md"
SESSION_FILE="${SESSION_CONTEXT_DIR}/project_state.md"

mkdir -p "$CONTEXT_DIR"
mkdir -p "$SESSION_CONTEXT_DIR"

# ------------------------------------------------
# 1. ATUALIZAR CONTEXTO PERMANENTE (se existir)
# ------------------------------------------------

if [ ! -f "$PERMANENT_FILE" ]; then
    echo "# 🧠 Estado Permanente do Projeto Arandu" > "$PERMANENT_FILE"
    echo "" >> "$PERMANENT_FILE"
    echo "**Criado em:** $(date)" >> "$PERMANENT_FILE"
    echo "**Última atualização:** $(date)" >> "$PERMANENT_FILE"
    echo "" >> "$PERMANENT_FILE"
    
    echo "## 📋 Histórico de Implementações" >> "$PERMANENT_FILE"
    echo "" >> "$PERMANENT_FILE"
    echo "| Data | Tarefa | Requirement | Status |" >> "$PERMANENT_FILE"
    echo "|------|--------|-------------|--------|" >> "$PERMANENT_FILE"
fi

# Atualizar data de última atualização no arquivo permanente
if [ -f "$PERMANENT_FILE" ]; then
    sed -i "s/^\*\*Última atualização:\*\* .*$/**Última atualização:** $(date)/" "$PERMANENT_FILE"
fi

# Verificar se há tarefas recentes para adicionar ao histórico
if [ -d "work/archive" ]; then
    LATEST_SESSION=$(ls -t work/archive | head -1)
    if [ -n "$LATEST_SESSION" ] && [ -f "work/archive/$LATEST_SESSION/task.md" ]; then
        # Extrair informações da tarefa
        REQ=$(grep -i "^Requirement:" "work/archive/$LATEST_SESSION/task.md" | cut -d: -f2 | xargs || echo "")
        TITLE=$(grep -i "^Task:" "work/archive/$LATEST_SESSION/task.md" | head -1 | cut -d: -f2- | xargs || echo "Tarefa $LATEST_SESSION")
        
        # Verificar se esta tarefa já está no histórico (por session ID)
        if ! grep -q "task_$LATEST_SESSION" "$PERMANENT_FILE" 2>/dev/null; then
            # Adicionar ao histórico
            sed -i "/^| Data | Tarefa | Requirement | Status |$/a | $(date +%Y-%m-%d) | $TITLE | $REQ | ✅ Concluída |" "$PERMANENT_FILE"
            echo "  ✅ Adicionada tarefa ao histórico permanente: $TITLE ($REQ)"
        fi
    fi
fi

# ------------------------------------------------
# 2. CRIAR CONTEXTO DA SESSÃO ATUAL
# ------------------------------------------------

echo "# 🧠 Estado Atual do Projeto Arandu" > "$SESSION_FILE"
echo "" >> "$SESSION_FILE"
echo "**Gerado em:** $(date)" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

# Incluir link para contexto permanente
echo "> 📚 **Contexto permanente:** [permanent_state.md](../project_context/permanent_state.md)" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# VISÃO DO SISTEMA
# ------------------------------------------------

echo "## 🌳 Estrutura de Visões" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

if [ -d "docs/vision" ]; then
  for V in docs/vision/*.md; do
    [ -f "$V" ] || continue
    echo "- $(basename "$V")" >> "$SESSION_FILE"
  done
else
  echo "Nenhuma visão encontrada." >> "$SESSION_FILE"
fi

echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# CAPABILITIES
# ------------------------------------------------

echo "## ⚙️ Capabilities" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

if [ -d "docs/capabilities" ]; then
  for C in docs/capabilities/*.md; do
    [ -f "$C" ] || continue
    echo "- $(basename "$C")" >> "$SESSION_FILE"
  done
else
  echo "Nenhuma capability encontrada." >> "$SESSION_FILE"
fi

echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# REQUIREMENTS
# ------------------------------------------------

echo "## 📋 Requirements" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

REQ_TOTAL=0

if [ -d "docs/requirements" ]; then
  for R in docs/requirements/*.md; do
    [ -f "$R" ] || continue
    echo "- $(basename "$R")" >> "$SESSION_FILE"
    REQ_TOTAL=$((REQ_TOTAL+1))
  done
fi

echo "" >> "$SESSION_FILE"
echo "**Total de requirements:** $REQ_TOTAL" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# TAREFAS RECENTES
# ------------------------------------------------

echo "## 🛠️ Tarefas recentes" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

if [ -d "work/archive" ]; then
  ls -t work/archive | head -5 | while read SESSION
  do
    echo "- $SESSION" >> "$SESSION_FILE"
  done
else
  echo "Nenhuma sessão arquivada ainda." >> "$SESSION_FILE"
fi

echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# APRENDIZADOS RECENTES
# ------------------------------------------------

echo "## 📚 Aprendizados recentes" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

if [ -d "docs/learnings" ]; then
  ls -t docs/learnings | head -5 | while read FILE
  do
    echo "- $FILE" >> "$SESSION_FILE"
  done
else
  echo "Nenhum aprendizado registrado." >> "$SESSION_FILE"
fi

echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# STATUS DO SISTEMA
# ------------------------------------------------

echo "## 📊 Status do Sistema" >> "$SESSION_FILE"
echo "" >> "$SESSION_FILE"

echo "- Visions: $(ls docs/vision 2>/dev/null | wc -l)" >> "$SESSION_FILE"
echo "- Capabilities: $(ls docs/capabilities 2>/dev/null | wc -l)" >> "$SESSION_FILE"
echo "- Requirements: $(ls docs/requirements 2>/dev/null | wc -l)" >> "$SESSION_FILE"
echo "- Learnings: $(ls docs/learnings 2>/dev/null | wc -l)" >> "$SESSION_FILE"

echo "" >> "$SESSION_FILE"

# ------------------------------------------------
# RESUMO DO PROGRESSO (do contexto permanente)
# ------------------------------------------------

if [ -f "$PERMANENT_FILE" ]; then
    TOTAL_TASKS=$(grep -c "✅ Concluída" "$PERMANENT_FILE" || echo "0")
    echo "## 📈 Progresso do Projeto" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
    echo "- **Tarefas concluídas:** $TOTAL_TASKS" >> "$SESSION_FILE"
    echo "- **Requirements implementados:** $((TOTAL_TASKS))" >> "$SESSION_FILE"
    echo "- **Última atualização permanente:** $(grep "Última atualização:" "$PERMANENT_FILE" | head -1 | cut -d: -f2-)" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
fi

echo "✅ Contexto atualizado:"
echo "   - Sessão: $SESSION_FILE"
echo "   - Permanente: $PERMANENT_FILE"