#!/usr/bin/env bash
# ⚠️  ONE-TIME SETUP SCRIPT - Run once only
# Initializes design system structure (may have already been run)
set -euo pipefail

echo "🎨 Inicializando identidade visual do Arandu..."

BASE="design"
DOCS="docs"

mkdir -p $BASE/css
mkdir -p $BASE/components
mkdir -p $BASE/wireframes
mkdir -p $DOCS

create_if_missing() {
FILE=$1
CONTENT=$2

if [ ! -f "$FILE" ]; then
  echo "$CONTENT" > "$FILE"
  echo "📄 Criado $FILE"
else
  echo "⚠️  $FILE já existe"
fi
}

# ------------------------------------------------------------------
# DESIGN SYSTEM
# ------------------------------------------------------------------

create_if_missing "$DOCS/design-system.md" "# Design System — Arandu

## Filosofia

Arandu é uma ferramenta de reflexão clínica.

Princípios:

- clareza > estética
- calma > impacto visual
- conteúdo > interface

---

## Paleta

Primária

#1E3A5F

Secundária

#3A7D6B

Insight / IA

#D4A84F

Fundo

#F7F8FA

Texto

#1F2937

---

## Tipografia

Interface

Inter

Conteúdo clínico

Source Serif

---

## Estrutura

Layout principal

Pacientes | Sessão | Insights

---

## Tipos de bloco

Observação  
Hipótese  
Intervenção  
Insight IA

---

## Stack UI

Go  
HTMX  
TailwindCSS  
AlpineJS (mínimo)

---

## Filosofia visual

Tecnologia silenciosa.

A interface nunca deve competir com o pensamento do terapeuta.
"

# ------------------------------------------------------------------
# TAILWIND CONFIG
# ------------------------------------------------------------------

create_if_missing "$BASE/tailwind.config.js" "module.exports = {
theme: {
extend: {
colors: {
arandu: {
primary: '#1E3A5F',
secondary: '#3A7D6B',
insight: '#D4A84F',
background: '#F7F8FA',
text: '#1F2937'
}
}
}
}
}
"

# ------------------------------------------------------------------
# BASE CSS
# ------------------------------------------------------------------

create_if_missing "$BASE/css/base.css" "body {
font-family: Inter, system-ui, sans-serif;
background: #F7F8FA;
color: #1F2937;
margin: 0;
padding: 0;
}

.card {
background: white;
border: 1px solid #E5E7EB;
border-radius: 10px;
padding: 16px;
}

.button-primary {
background: #1E3A5F;
color: white;
border-radius: 8px;
padding: 8px 14px;
}

.button-secondary {
background: transparent;
border: 1px solid #D1D5DB;
border-radius: 8px;
padding: 8px 14px;
}

.block-observation {
border-left: 4px solid #3A7D6B;
padding-left: 12px;
}

.block-hypothesis {
border-left: 4px solid #1E3A5F;
padding-left: 12px;
}

.block-intervention {
border-left: 4px solid #6B7280;
padding-left: 12px;
}

.block-insight {
border-left: 4px solid #D4A84F;
background: #FFFBEB;
padding-left: 12px;
}
"

# ------------------------------------------------------------------
# COMPONENTES
# ------------------------------------------------------------------

create_if_missing "$BASE/components/card.md" "# Card

Elemento básico de interface.

CSS:

background: white  
border: 1px solid #E5E7EB  
border-radius: 10px  
padding: 16px
"

create_if_missing "$BASE/components/buttons.md" "# Botões

Primário

- fundo azul
- texto branco

Secundário

- borda leve
- fundo transparente
"

create_if_missing "$BASE/components/editor-clinico.md" "# Editor Clínico

Elemento central do Arandu.

Estrutura

Observações  
Hipóteses  
Intervenções  
Insights IA

Deve parecer um caderno clínico digital.
"

# ------------------------------------------------------------------
# WIREFRAMES
# ------------------------------------------------------------------

create_if_missing "$BASE/wireframes/main_layout.md" "# Layout Principal

Pacientes | Sessão | Insights

┌───────────────┬───────────────────────────┬─────────────┐
│ Pacientes     │ Sessão atual              │ Insights IA │
│               │                           │             │
│ lista         │ editor clínico            │ sugestões   │
│               │                           │ padrões     │
└───────────────┴───────────────────────────┴─────────────┘
"

create_if_missing "$BASE/wireframes/session_editor.md" "# Editor de Sessão

Sessão — Paciente

Observações

Paciente relata ansiedade antes de reuniões.

Hipóteses

Possível associação com experiências escolares.

Intervenções

Exploração narrativa da memória.

Insights IA

Possível padrão de avaliação social.
"

# ------------------------------------------------------------------
# README
# ------------------------------------------------------------------

create_if_missing "$BASE/README.md" "# Identidade Visual do Arandu

Este diretório contém o design system do projeto.

Estrutura

css/  
componentes/  
wireframes/  

Documentação completa em

docs/design-system.md
"

echo
echo "🎉 Identidade visual inicial criada com sucesso."
echo "📁 Diretório: design/"


