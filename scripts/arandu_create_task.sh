#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/arandu_trace.sh" "$@"

trace "CREATE_TASK" "$TITLE" "REQ: $REQ_ID"

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 \"titulo da tarefa\" [REQ-ID]"
  echo "Exemplo 1: $0 \"Criando as fundações\""
  echo "Exemplo 2: $0 \"Implementar feature X\" req-01-01-01"
  exit 1
fi

# Se apenas um argumento foi fornecido, é o título
if [ $# -eq 1 ]; then
  TITLE=$1
  REQ_ID=""
else
  TITLE=$1
  REQ_ID=$2
fi

TASK_ID=$(date +%Y%m%d_%H%M%S)
TASK_DIR="work/tasks/task_${TASK_ID}"

mkdir -p "$TASK_DIR"

# Criar arquivo da tarefa
if [ -n "$REQ_ID" ]; then
  cat > "$TASK_DIR/task.md" <<EOF
# TASK ${TASK_ID}

Requirement: ${REQ_ID}

Title: ${TITLE}

## Objetivo

Implementar requirement ${REQ_ID}

## Referências

docs/requirements/${REQ_ID}.md
EOF
else
  cat > "$TASK_DIR/task.md" <<EOF
# TASK ${TASK_ID}

Title: ${TITLE}

## Status
**AGUARDANDO_DETALHES_DO_USUARIO** - NÃO inicie trabalho até que o usuário edite este arquivo

## Objetivo

${TITLE}

## Descrição

Tarefa criada sem requirement específico.

## Instruções para o Agente

**CRÍTICO: NÃO leia, edite ou execute qualquer ação nesta tarefa até que o usuário tenha editado este arquivo com os detalhes completos.**

1. Aguarde o usuário fornecer detalhes da tarefa neste arquivo
2. Quando o usuário editar este arquivo com os detalhes completos (removendo esta seção de instruções), inicie a implementação
3. Siga o padrão de referenciar requirements quando aplicável

**Verificação obrigatória antes de iniciar:**
- Esta seção "Instruções para o Agente" deve ter sido removida/replaceada pelo usuário
- O arquivo deve conter uma descrição detalhada da tarefa fornecida pelo usuário
- O status deve ter sido atualizado para "PRONTO_PARA_IMPLEMENTACAO" ou similar

## Checklist de Integridade (OBRIGATÓRIO)
- [ ] O componente usa .templ e herda de Layout?
- [ ] A tipografia Source Serif 4 foi aplicada ao conteúdo clínico?
- [ ] Executei 'templ generate' e o código Go compilou?
- [ ] Testei a rota atual e as rotas vizinhas (Regressão)?
- [ ] O banco de dados foi atualizado via migration .up.sql?

EOF
fi

trace "TASK_CREATED" "$TASK_ID" "Title: $TITLE"
echo "✅ tarefa criada: ${TASK_ID}"
echo "📁 Tarefa em: $TASK_DIR/task.md"
[ "$TRACE_ENABLED" = true ] && echo "📍 Trace: $TRACE_FILE"
echo ""
echo "📝 Instruções CRÍTICAS:"
echo "1. Edite $TASK_DIR/task.md para detalhar a tarefa COMPLETAMENTE"
echo "2. REMOVA a seção 'Instruções para o Agente' do arquivo task.md"
echo "3. Atualize o status para 'PRONTO_PARA_IMPLEMENTACAO'"
echo "4. O agente NÃO iniciará até que você tenha feito essas alterações"
