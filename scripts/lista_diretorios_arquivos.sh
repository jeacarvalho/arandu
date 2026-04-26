#!/bin/bash

# --- CONFIGURAÇÕES ---
# Cores para o terminal
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # Sem cor

echo -e "${BLUE}=== AUDITORIA DE ESTRUTURA ARANDU ===${NC}"
echo -e "Data: $(date)"
echo ""

# 1. VISUALIZAÇÃO DA ÁRVORE (Ignorando lixo e pastas ocultas)
echo -e "${YELLOW}--- Mapeamento de Directórios e Ficheiros ---${NC}"
if command -v tree >/dev/null 2>&1; then
    # Se o 'tree' estiver instalado, usa-o para melhor visualização
    tree -I ".git|node_modules|tmp|vendor" -F
else
    # Fallback com 'find' formatado
    find . -maxdepth 4 -not -path '*/.*' \
           -not -path './node_modules*' \
           -not -path './vendor*' \
           -not -path './tmp*' | sed -e 's/[^-][^\/]*\// |/g' -e 's/|\([^ ]\)/|-\1/'
fi

echo ""

# 2. DETECTOR DE DUPLICIDADES DE NOMES
# Procura ficheiros com o mesmo nome em locais diferentes
echo -e "${YELLOW}--- Detector de Possíveis Duplicidades (Mesmo nome em pastas diferentes) ---${NC}"
find . -type f -not -path '*/.*' -not -path './tmp/*' -not -path './vendor/*' | \
    sed 's|.*/||' | sort | uniq -d | while read -r filename; do
    echo -e "${RED}⚠ Duplicado detectado: ${filename}${NC}"
    find . -name "$filename" -not -path './tmp/*'
done

echo ""

# 3. RESUMO DE CONTEÚDO
echo -e "${YELLOW}--- Resumo de Quantidades ---${NC}"
echo -e "Total de Ficheiros Go:    $(find . -name "*.go" | wc -l)"
echo -e "Total de Templates Templ: $(find . -name "*.templ" | wc -l)"
echo -e "Total de Requisitos (.md):$(find . -name "req-*.md" | wc -l)"
echo -e "Total de Migrações SQL:   $(find . -name "*.sql" | wc -l)"

echo ""
echo -e "${BLUE}=====================================${NC}"
