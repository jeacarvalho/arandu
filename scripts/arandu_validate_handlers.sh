#!/usr/bin/env bash

set -e

echo "🔍 Validando handlers do sistema Arandu"
echo "========================================"

ERRORS=0
WARNINGS=0

# Função para verificar padrões problemáticos
check_pattern() {
    local pattern="$1"
    local description="$2"
    local severity="$3" # "ERROR" ou "WARNING"
    
    echo
    echo "🔎 Verificando: $description"
    
    local count=$(grep -r "$pattern" internal/web/handlers/ --include="*.go" | wc -l)
    
    if [ "$count" -gt 0 ]; then
        if [ "$severity" = "ERROR" ]; then
            echo "❌ $severity: Encontrado $count ocorrência(s) de '$pattern'"
            grep -r "$pattern" internal/web/handlers/ --include="*.go" --color=always
            ERRORS=$((ERRORS + count))
        else
            echo "⚠️  $severity: Encontrado $count ocorrência(s) de '$pattern'"
            grep -r "$pattern" internal/web/handlers/ --include="*.go" --color=always
            WARNINGS=$((WARNINGS + count))
        fi
    else
        echo "✅ OK: Nenhuma ocorrência encontrada"
    fi
}

# Função para verificar se componentes Templ existem
check_templ_component() {
    local component="$1"
    local description="$2"
    
    echo
    echo "🔎 Verificando: $description"
    
    if find web/components -name "*.templ" -exec grep -l "templ $component" {} \; | grep -q .; then
        echo "✅ OK: Componente '$component' encontrado"
    else
        echo "❌ ERROR: Componente '$component' não encontrado"
        ERRORS=$((ERRORS + 1))
    fi
}

echo
echo "📋 VALIDAÇÃO DE HANDLERS WEB"
echo "============================"

# 1. Verificar violações críticas (ERRORS)
check_pattern "ExecuteTemplate" "Uso de ExecuteTemplate (deprecated)" "ERROR"
check_pattern "html/template" "Import de html/template (deve usar templ)" "ERROR"
check_pattern "\.html\"" "Referência a arquivos .html (deve usar .templ)" "ERROR"
check_pattern "fmt\.Sprintf.*%s.*>" "HTML inline com fmt.Sprintf" "ERROR"
check_pattern "w\.Write.*<.*>" "HTML inline direto no ResponseWriter" "ERROR"

# 2. Verificar padrões de warning
check_pattern "TemplateRenderer" "Interface TemplateRenderer (legado)" "WARNING"
check_pattern "DummyRenderer" "Uso de DummyRenderer" "WARNING"
check_pattern "// TODO.*templ" "TODO para migração templ" "WARNING"

# 3. Verificar componentes Templ obrigatórios
echo
echo "📦 VERIFICAÇÃO DE COMPONENTES TEMPL"
echo "==================================="

check_templ_component "NewPatientForm" "Componente NewPatientForm"
check_templ_component "NewSessionForm" "Componente NewSessionForm"
check_templ_component "EditSessionForm" "Componente EditSessionForm"
check_templ_component "PatientList" "Componente PatientList"
check_templ_component "PatientDetail" "Componente PatientDetail"
check_templ_component "SessionDetailView" "Componente SessionDetailView"

# 4. Verificar imports corretos
echo
echo "📦 VERIFICAÇÃO DE IMPORTS"
echo "=========================="

echo "🔎 Verificando imports de templ"
if grep -r "github.com/a-h/templ" internal/web/handlers/ --include="*.go" | grep -q .; then
    echo "✅ OK: Import de templ encontrado"
else
    echo "⚠️  WARNING: Import de templ não encontrado nos handlers"
    WARNINGS=$((WARNINGS + 1))
fi

echo "🔎 Verificando imports de layoutComponents"
if grep -r "layoutComponents" internal/web/handlers/ --include="*.go" | grep -q .; then
    echo "✅ OK: Import de layoutComponents encontrado"
else
    echo "⚠️  WARNING: Import de layoutComponents não encontrado"
    WARNINGS=$((WARNINGS + 1))
fi

# 5. Verificar build
echo
echo "🔨 VERIFICAÇÃO DE BUILD"
echo "======================="

echo "🔎 Tentando build dos handlers..."
if go build ./internal/web/handlers/...; then
    echo "✅ OK: Build dos handlers bem-sucedido"
else
    echo "❌ ERROR: Falha no build dos handlers"
    ERRORS=$((ERRORS + 1))
fi

echo
echo "📊 RESUMO DA VALIDAÇÃO"
echo "======================"
echo "✅ Erros críticos: $ERRORS"
echo "⚠️  Avisos: $WARNINGS"
echo

if [ "$ERRORS" -gt 0 ]; then
    echo "❌ VALIDAÇÃO FALHOU: $ERRORS erro(s) crítico(s) encontrado(s)"
    echo "   Corrija os erros antes de prosseguir."
    exit 1
elif [ "$WARNINGS" -gt 0 ]; then
    echo "⚠️  VALIDAÇÃO COM AVISOS: $WARNINGS aviso(s)"
    echo "   Considere corrigir os avisos para melhorar a qualidade do código."
    exit 0
else
    echo "🎉 VALIDAÇÃO BEM-SUCEDIDA!"
    echo "   Todos os handlers estão seguindo os padrões arquiteturais."
    exit 0
fi