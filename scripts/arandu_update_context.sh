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
# 1. ATUALIZAR/CRIAR CONTEXTO PERMANENTE
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
else
    # Atualizar data de última atualização
    TEMP_FILE=$(mktemp)
    awk -v new_date="**Última atualização:** $(date)" \
        '/^\*\*Última atualização:\*\*/ {print new_date; next} {print}' \
        "$PERMANENT_FILE" > "$TEMP_FILE"
    mv "$TEMP_FILE" "$PERMANENT_FILE"
fi

# ------------------------------------------------
# 2. ADICIONAR TAREFAS RECENTES AO HISTÓRICO
# ------------------------------------------------

if [ -d "work/archive" ]; then
    # Processar todas as sessões arquivadas
    for SESSION_DIR in work/archive/*; do
        if [ -d "$SESSION_DIR" ] && [ -f "$SESSION_DIR/task.md" ]; then
            SESSION_NAME=$(basename "$SESSION_DIR")
            
            # Extrair informações da tarefa
            REQ=$(grep -i "^Requirement:" "$SESSION_DIR/task.md" | head -1 | cut -d: -f2 | xargs || echo "")
            TITLE=$(grep -i "^Task:" "$SESSION_DIR/task.md" | head -1 | cut -d: -f2- | xargs || echo "Tarefa $SESSION_NAME")
            
            # Verificar se esta tarefa já está no histórico e tem formato válido
            if [ -n "$REQ" ] && [[ "$REQ" =~ ^REQ-[0-9]{2}-[0-9]{2}-[0-9]{2} ]]; then
                # Verificar se já existe no histórico (por requirement)
                if ! grep -q "| $REQ |" "$PERMANENT_FILE" 2>/dev/null; then
                    # Adicionar linha ao histórico
                    sed -i "/^| Data | Tarefa | Requirement | Status |$/a | $(date +%Y-%m-%d) | $TITLE | $REQ | ✅ Concluída |" "$PERMANENT_FILE"
                    echo "  ✅ Adicionada tarefa: $TITLE ($REQ)"
                fi
            fi
        fi
    done
fi

# ------------------------------------------------
# 3. CRIAR CONTEXTO DA SESSÃO ATUAL
# ------------------------------------------------

create_session_context() {
    echo "# 🧠 Estado Atual do Projeto Arandu" > "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
    echo "**Gerado em:** $(date)" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"

    # Incluir link para contexto permanente
    echo "> 📚 **Contexto permanente:** [permanent_state.md](../project_context/permanent_state.md)" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"

    # Visões
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

    # Capabilities
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

    # Requirements
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

    # Tarefas recentes
    echo "## 🛠️ Tarefas recentes" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
    if [ -d "work/archive" ]; then
        ls -t work/archive | head -5 | while read SESSION; do
            echo "- $SESSION" >> "$SESSION_FILE"
        done
    else
        echo "Nenhuma sessão arquivada ainda." >> "$SESSION_FILE"
    fi
    echo "" >> "$SESSION_FILE"

    # Aprendizados consolidados
    echo "## 📚 Sistema de Aprendizados" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
    if [ -f "docs/learnings/MASTER_LEARNINGS.md" ]; then
        echo "✅ **Sistema consolidado ativo**" >> "$SESSION_FILE"
        echo "" >> "$SESSION_FILE"
        echo "Arquivos principais:" >> "$SESSION_FILE"
        echo "- MASTER_LEARNINGS.md (aprendizados consolidados)" >> "$SESSION_FILE"
        echo "- ARCHITECTURE_PATTERNS.md (padrões arquiteturais)" >> "$SESSION_FILE"
        echo "- TEMPL_GUIDE.md (guia Templ)" >> "$SESSION_FILE"
        echo "- SQLITE_BEST_PRACTICES.md (práticas SQLite)" >> "$SESSION_FILE"
        echo "" >> "$SESSION_FILE"
        echo "📊 **Estatísticas:**" >> "$SESSION_FILE"
        echo "- Arquivos consolidados: 4 principais" >> "$SESSION_FILE"
        echo "- Arquivos em archive: $(ls docs/learnings/archive/*.md 2>/dev/null | wc -l)" >> "$SESSION_FILE"
    else
        echo "⚠️ Sistema de aprendizados não configurado." >> "$SESSION_FILE"
    fi
    echo "" >> "$SESSION_FILE"

    # Status do sistema
    echo "## 📊 Status do Sistema" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"
    echo "- Visions: $(ls docs/vision 2>/dev/null | wc -l)" >> "$SESSION_FILE"
    echo "- Capabilities: $(ls docs/capabilities 2>/dev/null | wc -l)" >> "$SESSION_FILE"
    echo "- Requirements: $(ls docs/requirements 2>/dev/null | wc -l)" >> "$SESSION_FILE"
    echo "- Learnings: Sistema consolidado (4 arquivos principais + archive)" >> "$SESSION_FILE"
    echo "" >> "$SESSION_FILE"

    # Progresso do projeto
    if [ -f "$PERMANENT_FILE" ]; then
        TOTAL_TASKS=$(grep -c "✅ Concluída" "$PERMANENT_FILE" || echo "0")
        LAST_UPDATE=$(grep "Última atualização:" "$PERMANENT_FILE" | head -1 | cut -d: -f2-)
        echo "## 📈 Progresso do Projeto" >> "$SESSION_FILE"
        echo "" >> "$SESSION_FILE"
        echo "- **Tarefas concluídas:** $TOTAL_TASKS" >> "$SESSION_FILE"
        echo "- **Requirements implementados:** $TOTAL_TASKS" >> "$SESSION_FILE"
        echo "- **Última atualização permanente:** $LAST_UPDATE" >> "$SESSION_FILE"
        echo "" >> "$SESSION_FILE"
    fi
}

create_session_context

echo "✅ Contexto atualizado:"
echo "   - Sessão: $SESSION_FILE"
echo "   - Permanente: $PERMANENT_FILE"