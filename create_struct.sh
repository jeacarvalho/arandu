#!/usr/bin/env bash

set -e

echo "🚀 Inicializando estrutura de desenvolvimento do Arandu..."
echo

# Diretórios principais
mkdir -p docs
mkdir -p docs/vision
mkdir -p docs/capabilities
mkdir -p docs/requirements
mkdir -p docs/learnings

mkdir -p scripts

mkdir -p work
mkdir -p work/current_session
mkdir -p work/tasks
mkdir -p work/archive

mkdir -p work/templates

echo "✅ Estrutura de diretórios criada"

# Criar arquivos base se não existirem
create_if_missing() {
    FILE=$1
    CONTENT=$2

    if [ ! -f "$FILE" ]; then
        echo "$CONTENT" > "$FILE"
        echo "📄 Criado: $FILE"
    else
        echo "⚠️  Já existe: $FILE"
    fi
}

echo
echo "📄 Criando arquivos base..."

create_if_missing "docs/dvp.md" "# Documento de Visão do Projeto — Arandu"

create_if_missing "docs/learnings/README.md" "# Aprendizados do Projeto

Este diretório contém aprendizados permanentes coletados durante o desenvolvimento.

Cada arquivo deve referenciar um Requirement (REQ)."

create_if_missing "work/README.md" "# Área de Trabalho

Este diretório contém sessões de desenvolvimento e tarefas temporárias.

Estrutura:

work/
  current_session/
  tasks/
  archive/"

create_if_missing "scripts/README.md" "# Scripts de Automação

Scripts utilizados para gerenciar sessões, tarefas e aprendizados."

echo
echo "📁 Estrutura final criada:"
echo

tree_structure='
docs/
  dvp.md
  vision/
  capabilities/
  requirements/
  learnings/

scripts/

work/
  current_session/
  tasks/
  archive/
  templates/
'

echo "$tree_structure"

echo
echo "🎉 Estrutura do Arandu inicializada com sucesso."
