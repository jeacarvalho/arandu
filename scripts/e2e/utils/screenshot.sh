#!/bin/bash
# scripts/e2e/utils/screenshot.sh
# Utilitário para captura de screenshots com autenticação
# Versão: 2.1 — Correção PROJECT_ROOT + cleanup

source "$(dirname "${BASH_SOURCE[0]}")/../config.sh"

e2e_capture_screenshot() {
    local route="$1"
    local output_file="$2"
    local width="${3:-1440}"
    local height="${4:-900}"
    
    # Verificar cookie
    if [ -z "$E2E_COOKIE_VALUE" ]; then
        e2e_log_warn "⚠️ No session cookie for screenshot"
        return 1
    fi
    
    # Verificar Node.js
    if ! command -v node &>/dev/null; then
        e2e_log_warn "⚠️ Node.js not installed, skipping screenshot"
        return 0
    fi
    
    # Verificar Playwright NO CONTEXTO DO PROJETO
    if ! (cd "$PROJECT_ROOT" && node -e "require('playwright')" &>/dev/null); then
        e2e_log_warn "⚠️ Playwright not installed in project context"
        e2e_log_info "💡 Install: cd $PROJECT_ROOT && npm install playwright && npx playwright install chromium"
        return 0
    fi
    
    e2e_log_info "📸 Capturing: $route (${width}x${height})"
    
    # CRÍTICO: Criar script NO DIRETÓRIO DO PROJETO (não em /tmp/)
    local tmp_script="${PROJECT_ROOT}/tmp_screenshot_$$.js"
    
    # Escapar caracteres especiais no cookie
    local escaped_cookie
    escaped_cookie=$(printf '%s\n' "$E2E_COOKIE_VALUE" | sed "s/'/\\\\'/g")
    
    cat > "$tmp_script" <<EOF
const { chromium } = require('playwright');

(async () => {
    const browser = await chromium.launch();
    const context = await browser.newContext({
        viewport: { width: $width, height: $height }
    });
    
    await context.addCookies([{
        name: 'arandu_session',
        value: '$escaped_cookie',
        domain: 'localhost',
        path: '/',
    }]);
    
    const page = await context.newPage();
    await page.goto('$E2E_BASE_URL$route', { waitUntil: 'networkidle' });
    await page.screenshot({ path: '$output_file', fullPage: true });
    
    await browser.close();
})();
EOF
    
    # CRÍTICO: Executar node NO DIRETÓRIO DO PROJETO (subshell)
    (cd "$PROJECT_ROOT" && node "$tmp_script" 2>&1)
    local node_exit=$?
    
    # CRÍTICO: CLEANUP — Remover script temporário
    rm -f "$tmp_script"
    
    if [ $node_exit -eq 0 ] && [ -f "$output_file" ]; then
        e2e_log_info "✅ Saved: $output_file"
        return 0
    else
        e2e_log_warn "⚠️ Screenshot failed: $route"
        return 1
    fi
}
