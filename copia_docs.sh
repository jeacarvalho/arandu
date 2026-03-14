#!/usr/bin/env bash

SRC_DIR="$1"
OUTPUT_FILE="$2"

# Validação de argumentos
if [ -z "$SRC_DIR" ] || [ -z "$OUTPUT_FILE" ]; then
  echo "Uso: $0 <diretorio_origem> <arquivo_saida>"
  exit 1
fi

# Converte para caminho absoluto para evitar erros de referência
SRC_DIR=$(realpath "$SRC_DIR")

if [ ! -d "$SRC_DIR" ]; then
  echo "Erro: O diretório de origem '$SRC_DIR' não existe."
  exit 1
fi

echo "Diretório origem: $SRC_DIR"
echo "Arquivo saída: $OUTPUT_FILE"
echo "--------------------------------"

# Criar lista temporária de arquivos .md recursivamente
# O 'sort -V' organiza numericamente (ex: 1, 2, 10 em vez de 1, 10, 2)
TMP_LIST=$(mktemp)
find "$SRC_DIR" -type f -name "*.md" | sort -V > "$TMP_LIST"

# Verifica se encontrou arquivos
if [ ! -s "$TMP_LIST" ]; then
  echo "Aviso: Nenhum arquivo .md encontrado em $SRC_DIR"
  rm "$TMP_LIST"
  exit 1
fi

# Limpa/Cria o arquivo de saída
: > "$OUTPUT_FILE"

{
    echo "# MERGED DOCUMENTATION"
    echo "Generated on: $(date)"
    echo ""
    echo "## FILE INDEX"
    echo ""

    # 1. Gerar o índice primeiro
    while read -r file; do
        rel="${file#$SRC_DIR/}"
        echo "- $rel"
    done < "$TMP_LIST"

    echo ""
    echo "---"
    echo ""

    # 2. Concatenar o conteúdo
    while read -r file; do
        rel="${file#$SRC_DIR/}"
        
        echo "Processando: $rel" >&2

        echo "===== BEGIN FILE: $rel ====="
        echo ""
        # tr -d remove caracteres nulos que podem corromper o output
        tr -d '\000' < "$file"
        echo ""
        echo "===== END FILE: $rel ====="
        echo ""
    done < "$TMP_LIST"

} >> "$OUTPUT_FILE"

rm "$TMP_LIST"

echo "--------------------------------"
echo "Merge concluído: $(wc -l < "$OUTPUT_FILE") linhas geradas."
