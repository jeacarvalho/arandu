#!/usr/bin/env bash

set -e

BASE_DIR="docs/requirements"

mkdir -p "$BASE_DIR"

create_requirement() {
  FILE="$1"
  ID="$2"
  VISION="$3"
  CAP="$4"
  TITLE="$5"

  if [ -f "$FILE" ]; then
    echo "⚠️  Já existe: $FILE"
    return
  fi

  cat <<EOF > "$FILE"
---
id: $ID
vision: VISION-$VISION
capability: CAP-$VISION-$CAP
status: draft
---

# $ID — $TITLE

## Visão associada

VISION-$VISION

## Capability associada

CAP-$VISION-$CAP

## Descrição

TODO: Descrever o comportamento esperado do sistema.

## Critérios de aceitação

- TODO

## Observações

Requirement derivado de CAP-$VISION-$CAP.
EOF

  echo "✅ Criado: $FILE"
}

# CAP-01-01 Registro de sessões
create_requirement "$BASE_DIR/req-01-01-01-criar-sessao.md" "REQ-01-01-01" "01" "01" "Criar sessão clínica"
create_requirement "$BASE_DIR/req-01-01-02-editar-sessao.md" "REQ-01-01-02" "01" "01" "Editar sessão clínica"
create_requirement "$BASE_DIR/req-01-01-03-listar-sessoes.md" "REQ-01-01-03" "01" "01" "Listar sessões do paciente"

# CAP-01-02 Observações clínicas
create_requirement "$BASE_DIR/req-01-02-01-adicionar-observacao.md" "REQ-01-02-01" "01" "02" "Adicionar observação clínica"
create_requirement "$BASE_DIR/req-01-02-02-editar-observacao.md" "REQ-01-02-02" "01" "02" "Editar observação clínica"

# CAP-01-03 Intervenções terapêuticas
create_requirement "$BASE_DIR/req-01-03-01-registrar-intervencao.md" "REQ-01-03-01" "01" "03" "Registrar intervenção terapêutica"

# CAP-02-01 Histórico do paciente
create_requirement "$BASE_DIR/req-02-01-01-visualizar-historico.md" "REQ-02-01-01" "02" "01" "Visualizar histórico clínico do paciente"

# CAP-02-02 Linha do tempo clínica
create_requirement "$BASE_DIR/req-02-02-01-linha-tempo.md" "REQ-02-02-01" "02" "02" "Visualizar linha do tempo clínica"

# CAP-03-01 Organização de observações
create_requirement "$BASE_DIR/req-03-01-01-classificar-observacao.md" "REQ-03-01-01" "03" "01" "Classificar observações clínicas"

# CAP-03-02 Organização de intervenções
create_requirement "$BASE_DIR/req-03-02-01-classificar-intervencao.md" "REQ-03-02-01" "03" "02" "Classificar intervenções terapêuticas"

# CAP-04-01 Padrões clínicos
create_requirement "$BASE_DIR/req-04-01-01-detectar-padroes.md" "REQ-04-01-01" "04" "01" "Detectar padrões clínicos"

# CAP-05-01 Assistente reflexivo
create_requirement "$BASE_DIR/req-05-01-01-consulta-ia.md" "REQ-05-01-01" "05" "01" "Consultar assistente reflexivo com IA"

# CAP-06-01 Comparação de casos
create_requirement "$BASE_DIR/req-06-01-01-comparar-casos.md" "REQ-06-01-01" "06" "01" "Comparar casos clínicos"

# CAP-07-01 Agenda
create_requirement "$BASE_DIR/req-07-01-01-gerenciar-agenda.md" "REQ-07-01-01" "07" "01" "Gerenciar agenda clínica"

# CAP-07-02 Atendimentos
create_requirement "$BASE_DIR/req-07-02-01-registrar-atendimento.md" "REQ-07-02-01" "07" "02" "Registrar atendimento clínico"

# CAP-08-01 Evolução da base clínica
create_requirement "$BASE_DIR/req-08-01-01-evolucao-base.md" "REQ-08-01-01" "08" "01" "Evolução da base clínica"

# CAP-09-01 Inteligência clínica
create_requirement "$BASE_DIR/req-09-01-01-analise-ia.md" "REQ-09-01-01" "09" "01" "Análise clínica assistida por IA"

# CAP-10-01 Aprendizado coletivo
create_requirement "$BASE_DIR/req-10-01-01-base-anonimizada.md" "REQ-10-01-01" "10" "01" "Gerar base clínica anonimizada"

echo
echo "🎉 Requirements criados em $BASE_DIR"
