#!/usr/bin/env bash

set -e

BASE_DIR="docs/capabilities"

mkdir -p "$BASE_DIR"

create_capability() {
  FILE="$1"
  ID="$2"
  VISION="$3"
  TITLE="$4"

  if [ -f "$FILE" ]; then
    echo "⚠️  Já existe: $FILE"
    return
  fi

  cat <<EOF > "$FILE"
---
id: $ID
vision: VISION-$VISION
status: draft
---

# $ID — $TITLE

## Visão associada

VISION-$VISION

## Descrição

TODO: Descrever esta capability.

## Requisitos relacionados

TODO

## Observações

Capability derivada de VISION-$VISION.
EOF

  echo "✅ Criado: $FILE"
}

# VISION 01 — Registro da prática clínica
create_capability "$BASE_DIR/cap-01-01-registro-sessoes.md" "CAP-01-01" "01" "Registro de sessões clínicas"
create_capability "$BASE_DIR/cap-01-02-observacoes-clinicas.md" "CAP-01-02" "01" "Registro de observações clínicas"
create_capability "$BASE_DIR/cap-01-03-intervencoes-terapeuticas.md" "CAP-01-03" "01" "Registro de intervenções terapêuticas"

# VISION 02 — Memória clínica
create_capability "$BASE_DIR/cap-02-01-historico-paciente.md" "CAP-02-01" "02" "Histórico clínico do paciente"
create_capability "$BASE_DIR/cap-02-02-linha-tempo-clinica.md" "CAP-02-02" "02" "Linha do tempo clínica do paciente"

# VISION 03 — Organização do conhecimento
create_capability "$BASE_DIR/cap-03-01-organizacao-observacoes.md" "CAP-03-01" "03" "Organização de observações clínicas"
create_capability "$BASE_DIR/cap-03-02-organizacao-intervencoes.md" "CAP-03-02" "03" "Organização de intervenções terapêuticas"

# VISION 04 — Padrões clínicos
create_capability "$BASE_DIR/cap-04-01-identificacao-padroes.md" "CAP-04-01" "04" "Identificação de padrões clínicos"

# VISION 05 — Assistência reflexiva
create_capability "$BASE_DIR/cap-05-01-assistente-reflexivo.md" "CAP-05-01" "05" "Assistente reflexivo com IA"

# VISION 06 — Comparação de casos
create_capability "$BASE_DIR/cap-06-01-comparacao-casos.md" "CAP-06-01" "06" "Comparação entre casos clínicos"

# VISION 07 — Gestão do consultório
create_capability "$BASE_DIR/cap-07-01-gestao-agenda.md" "CAP-07-01" "07" "Gestão de agenda clínica"
create_capability "$BASE_DIR/cap-07-02-gestao-atendimentos.md" "CAP-07-02" "07" "Gestão de atendimentos"

# VISION 08 — Evolução da base clínica
create_capability "$BASE_DIR/cap-08-01-evolucao-base-clinica.md" "CAP-08-01" "08" "Evolução da base clínica"

# VISION 09 — Inteligência clínica
create_capability "$BASE_DIR/cap-09-01-analise-clinica.md" "CAP-09-01" "09" "Análise clínica assistida por IA"

# VISION 10 — Aprendizado coletivo
create_capability "$BASE_DIR/cap-10-01-base-clinica-coletiva.md" "CAP-10-01" "10" "Base clínica anonimizada coletiva"

echo
echo "🎉 Capabilities criadas em $BASE_DIR"
