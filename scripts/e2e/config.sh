#!/bin/bash
# scripts/e2e/config.sh
# Configurações globais do E2E Audit

# CRÍTICO: Definir PROJECT_ROOT
export PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# URLs e Portas
export E2E_BASE_URL="${E2E_BASE_URL:-http://localhost:8080}"
export E2E_PORT="${E2E_PORT:-8080}"

# Diretórios
export E2E_AUDIT_DIR="${E2E_AUDIT_DIR:-tmp/audit_logs}"
export E2E_SCREENSHOT_DIR="${E2E_SCREENSHOT_DIR:-tmp/audit_screenshots}"

# Arquivos de sessão
export E2E_COOKIES_FILE="${E2E_COOKIES_FILE:-tmp/e2e_cookies.txt}"

# Credenciais de teste
export E2E_TEST_EMAIL="${E2E_TEST_EMAIL:-arandu_e2e@test.com}"
export E2E_TEST_PASS="${E2E_TEST_PASS:-test123456}"

# Timeouts
export E2E_SERVER_TIMEOUT="${E2E_SERVER_TIMEOUT:-30}"
export E2E_REQUEST_TIMEOUT="${E2E_REQUEST_TIMEOUT:-10}"

# Cores para logging
export E2E_COLOR_RED='\033[0;31m'
export E2E_COLOR_GREEN='\033[0;32m'
export E2E_COLOR_YELLOW='\033[1;33m'
export E2E_COLOR_BLUE='\033[0;34m'
export E2E_COLOR_NC='\033[0m'

# Contadores globais
export E2E_TOTAL_TESTS=0
export E2E_PASSED_TESTS=0
export E2E_FAILED_TESTS=0