#!/bin/bash
echo "🔍 Verificando CSS..."

# 1. Verificar se CSS foi modificado recentemente
if find web/static/css -name "*.css" -mmin -10 | grep -q .; then
    echo "⚠️  CSS modificado recentemente"
    echo "💡 Execute: make build && ./safe_deploy.sh"
fi

# 2. Verificar se templates foram gerados
if find web/components -name "*.templ" -newer web/components/*_templ.go 2>/dev/null | grep -q .; then
    echo "❌ Templates .templ modificados mas não gerados"
    echo "💡 Execute: ~/go/bin/templ generate"
    exit 1
fi

# 3. Testar CSS via curl
CSS_VERSION=$(curl -s http://localhost:8080/ | grep -o 'style.css?v=[^"]*' | head -1)
if [ -n "$CSS_VERSION" ]; then
    echo "✅ CSS versionado: $CSS_VERSION"
else
    echo "⚠️  CSS sem versionamento detectado"
fi