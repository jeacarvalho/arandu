#!/usr/bin/env bash

set -e

BASE_DIR="docs/vision"

mkdir -p "$BASE_DIR"

create_file() {
  FILE="$1"
  ID="$2"
  TITLE="$3"

  if [ -f "$FILE" ]; then
    echo "⚠️  Já existe: $FILE"
    return
  fi

  cat <<EOF > "$FILE"
---
id: $ID
parent: dvp.md
status: draft
---

# $ID — $TITLE

## Descrição

TODO: Expandir a visão.

## Valor para o usuário

TODO: Descrever valor para o terapeuta.

## Capacidades relacionadas

TODO

## Observações

Arquivo originado do DVP.
EOF

  echo "✅ Criado: $FILE"
}

create_file "$BASE_DIR/vision-01-registro-pratica-clinica.md" "VISION-01" "Registro estruturado da prática clínica"
create_file "$BASE_DIR/vision-02-memoria-clinica-longitudinal.md" "VISION-02" "Memória clínica longitudinal"
create_file "$BASE_DIR/vision-03-organizacao-conhecimento-clinico.md" "VISION-03" "Organização do conhecimento clínico"
create_file "$BASE_DIR/vision-04-descoberta-padroes-clinicos.md" "VISION-04" "Descoberta de padrões clínicos"
create_file "$BASE_DIR/vision-05-assistencia-reflexiva-ia.md" "VISION-05" "Assistência reflexiva com IA"
create_file "$BASE_DIR/vision-06-comparacao-casos-clinicos.md" "VISION-06" "Comparação entre casos clínicos"
create_file "$BASE_DIR/vision-07-organizacao-operacional-consultorio.md" "VISION-07" "Organização operacional do consultório"
create_file "$BASE_DIR/vision-08-base-clinica-evolutiva.md" "VISION-08" "Base clínica evolutiva"
create_file "$BASE_DIR/vision-09-inteligencia-clinica-ampliada.md" "VISION-09" "Inteligência clínica ampliada"
create_file "$BASE_DIR/vision-10-aprendizado-clinico-coletivo.md" "VISION-10" "Aprendizado clínico coletivo"

echo
echo "🎉 Arquivos de visão criados em $BASE_DIR"
