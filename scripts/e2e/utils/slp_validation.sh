#!/bin/bash
# scripts/e2e/utils/slp_validation.sh
# Validações do Standardized Layout Protocol (SLP)
# Versão: 1.0 - Completa

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"

# Verificar se sidebar tem links esperados (Anamnese, Prontuário)
slp_check_sidebar() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -q "Anamnese" "$file" 2>/dev/null || ! grep -q "Prontuário" "$file" 2>/dev/null; then
        e2e_log_error "$test_name - SLP sidebar missing expected links"
        return 1
    fi
    
    e2e_log_info "$test_name - SLP sidebar validated"
    return 0
}

# Verificar tipografia clínica (Source Serif 4)
slp_check_clinical_typography() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -q "font-clinical\|font-serif\|Source Serif" "$file" 2>/dev/null; then
        e2e_log_error "$test_name - Clinical typography (.font-clinical) not found"
        return 1
    fi
    
    e2e_log_info "$test_name - Clinical typography validated"
    return 0
}

# Verificar estrutura de layout (app-container, main-content, sidebar, top-bar)
slp_check_layout_structure() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -qE 'class=".*(app-container|main-content|sidebar|top-bar)' "$file" 2>/dev/null; then
        e2e_log_error "$test_name - Missing layout structure classes"
        return 1
    fi
    
    e2e_log_info "$test_name - Layout structure classes present"
    return 0
}

# Verificar spacing/design classes
slp_check_spacing_classes() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -qE 'class=".*(space-|margin-|padding-|mb-|mt-|ml-|mr-|mx-|my-|pb-|pt-|pl-|pr-|px-|py-|gap-|p-|m-|align-|flex-|grid-|text-|font-|bg-|rounded-|shadow-|border-|opacity-|z-)' "$file" 2>/dev/null; then
        e2e_log_warn "$test_name - No spacing/design classes found (may cause visual issues)"
        return 0
    fi
    
    e2e_log_info "$test_name - Spacing/design classes present"
    return 0
}

# Verificar padding-top para header fixo
slp_check_fixed_header_spacing() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -qE '<header.*class=".*top-bar' "$file" 2>/dev/null; then
        return 0
    fi
    
    if grep -qE '\.main-content\s*\{[^}]*padding-top:\s*0' "$file" 2>/dev/null; then
        e2e_log_warn "$test_name - Main content has no top padding for fixed header"
        return 0
    fi
    
    e2e_log_info "$test_name - Fixed header spacing validated"
    return 0
}

# Verificar z-index em elementos fixed
slp_check_fixed_z_index() {
    local file="$1"
    local test_name="$2"
    
    if ! grep -q 'position:\s*fixed' "$file" 2>/dev/null; then
        return 0
    fi
    
    if ! grep -q 'z-index:' "$file" 2>/dev/null; then
        if ! grep -q 'class=".*z-' "$file" 2>/dev/null; then
            e2e_log_warn "$test_name - Fixed elements without z-index detected"
            return 0
        fi
    fi
    
    e2e_log_info "$test_name - Fixed elements z-index validated"
    return 0
}

# Verificar CSS conflicts (!important, sidebar overlap)
slp_check_css_conflicts() {
    local file="$1"
    local test_name="$2"
    
    if grep -q 'style="[^"]*!important' "$file" 2>/dev/null; then
        e2e_log_warn "$test_name - CSS !important detected (may cause conflicts)"
    fi
    
    if grep -q '<aside.*sidebar.*top: 80px' "$file" 2>/dev/null; then
        if grep -q '<header.*top-bar' "$file" 2>/dev/null; then
            e2e_log_warn "$test_name - Sidebar positioned below header (should be aligned at top)"
        fi
    fi
    
    if grep -qE 'style="[^"]*left:\s*0[^"]*' "$file" 2>/dev/null; then
        if grep -q '<aside.*sidebar' "$file" 2>/dev/null; then
            e2e_log_warn "$test_name - Sidebar with left:0 may overlap header"
        fi
    fi
    
    e2e_log_info "$test_name - CSS conflicts check completed"
    return 0
}

# Validação completa SLP (chama todas as validações)
slp_validate_full() {
    local file="$1"
    local test_name="$2"
    local skip_layout_check="${3:-false}"
    local errors=0
    local warnings=0
    
    # Validações críticas (contam como erros)
    slp_check_sidebar "$file" "$test_name" || errors=$((errors + 1))
    slp_check_clinical_typography "$file" "$test_name" || errors=$((errors + 1))
    slp_check_layout_structure "$file" "$test_name" || errors=$((errors + 1))
    
    # Validações não-críticas (apenas warnings)
    if [ "$skip_layout_check" != "true" ]; then
        slp_check_spacing_classes "$file" "$test_name" || warnings=$((warnings + 1))
        slp_check_fixed_header_spacing "$file" "$test_name" || warnings=$((warnings + 1))
        slp_check_fixed_z_index "$file" "$test_name" || warnings=$((warnings + 1))
        slp_check_css_conflicts "$file" "$test_name" || warnings=$((warnings + 1))
    fi
    
    # Reportar warnings
    if [ $warnings -gt 0 ]; then
        e2e_log_warn "$test_name - $warnings warning(s) found"
    fi
    
    # Retornar APENAS erros críticos
    return $errors
}

# Verificar se main-content tem padding-top adequado para header fixo
slp_check_main_content_overflow() {
    local file="$1"
    local test_name="$2"
    
    # Verificar se há header fixo
    if grep -qE '<header.*class=".*top-bar.*fixed' "$file" 2>/dev/null; then
        # Verificar se main-content tem padding-top ou margin-top
        if ! grep -qE '\.main-content\s*\{[^}]*padding-top:\s*[6-9][0-9]px' "$file" 2>/dev/null && \
           ! grep -qE 'class=".*main-content.*pt-[123456789]' "$file" 2>/dev/null; then
            e2e_log_warn "$test_name - Main content may overlap fixed header (check padding-top)"
            return 0  # Warning, não erro crítico
        fi
    fi
    
    # Verificar z-index da sidebar vs main-content
    if grep -q '<aside.*sidebar' "$file" 2>/dev/null; then
        if ! grep -qE 'class=".*sidebar.*z-\[?[1-9]' "$file" 2>/dev/null; then
            e2e_log_warn "$test_name - Sidebar may need explicit z-index to stay above content"
            return 0
        fi
    fi
    
    e2e_log_info "$test_name - Layout overflow check passed"
    return 0
}

# Verificar densidade de conteúdo no main canvas
slp_check_layout_density() {
    local file="$1"
    local test_name="$2"
    local min_density="${3:-40}"  # Mínimo 40% de densidade por padrão
    
    # Analisar estrutura HTML para estimar uso de espaço
    local total_elements=$(grep -c '<div\|<section\|<article' "$file" 2>/dev/null || echo "0")
    local empty_containers=$(grep -cE '<div[^>]*>\s*</div>|<div[^>]*>\s*<br\s*/?>\s*</div>' "$file" 2>/dev/null || echo "0")
    
    # Verificar se há grid/flexbox sendo usado
    local has_layout=$(grep -cE 'class=".*grid|class=".*flex|class=".*grid-cols|class=".*flex-col' "$file" 2>/dev/null || echo "0")
    
    # Verificar cards/containers vazios ou subutilizados
    local large_empty_spaces=$(grep -cE 'class=".*h-96|h-screen|min-h-screen' "$file" 2>/dev/null || echo "0")
    
    if [ "$total_elements" -gt 0 ]; then
        local waste_ratio=$((empty_containers * 100 / total_elements))
        
        if [ "$waste_ratio" -gt 30 ]; then
            e2e_log_warn "$test_name - Alto índice de containers vazios ($waste_ratio% dos elementos)"
            e2e_log_warn "$test_name - Considere usar grid/flexbox para melhor aproveitamento"
            return 1
        fi
    fi
    
    if [ "$has_layout" -eq 0 ] && [ "$total_elements" -gt 10 ]; then
        e2e_log_warn "$test_name - Múltiplos elementos sem layout grid/flex (pode desperdiçar espaço)"
        return 0
    fi
    
    e2e_log_info "$test_name - Layout density acceptable"
    return 0
}

# Verificar se está usando sistema de grid adequadamente
slp_check_grid_usage() {
    local file="$1"
    local test_name="$2"
    
    # Verificar se há múltiplos cards/elementos que poderiam estar em grid
    local card_count=$(grep -cE 'class=".*card|class=".*rounded-lg.*p-' "$file" 2>/dev/null || echo "0")
    local grid_usage=$(grep -cE 'class=".*grid|class=".*grid-cols' "$file" 2>/dev/null || echo "0")
    
    if [ "$card_count" -ge 2 ] && [ "$grid_usage" -eq 0 ]; then
        e2e_log_warn "$test_name - $card_count cards detectados sem uso de grid (espaço pode estar sendo desperdiçado)"
        e2e_log_info "$test_name - Considere usar: grid grid-cols-2 gap-4 ou similar"
        return 1
    fi
    
    e2e_log_info "$test_name - Grid usage appropriate"
    return 0
}