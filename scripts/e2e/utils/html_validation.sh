#!/bin/bash
# scripts/e2e/utils/html_validation.sh

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"

e2e_check_ghosting() {
    local file="$1" step="$2"
    grep -qP '\{ \.?[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+ \}' "$file" 2>/dev/null && { e2e_log_error "$step - Ghosting"; return 1; }
    return 0
}

e2e_check_no_inline_styles() {
    local file="$1" step="$2"
    local count=$(grep -o 'style="' "$file" | wc -l)
    [ "$count" -gt 0 ] && { e2e_log_error "$step - Inline styles: $count"; return 1; }
    return 0
}

e2e_validate_html() {
    local file="$1" step="$2" validate_layout="${3:-true}" errors=0
    e2e_check_ghosting "$file" "$step" || errors=$((errors+1))
    e2e_check_no_inline_styles "$file" "$step" || errors=$((errors+1))
    return $errors
}