#!/bin/bash
# Script seguro de deployment do Arandu
# Garante que só existe uma instância do app rodando na porta 8080
# Compila templates e Tailwind CSS automaticamente

set -e

echo "🛑 Parando todas as instâncias do Arandu..."
pkill -9 -f "arandu" 2>/dev/null || true
sleep 2

# Verificar se a porta 8080 está liberada
PORT_PID=$(lsof -ti:8080 2>/dev/null)
if [ ! -z "$PORT_PID" ]; then
    echo "⚠️  Porta 8080 ainda em uso pelo PID $PORT_PID. Finalizando..."
    kill -9 $PORT_PID 2>/dev/null || true
    sleep 1
fi

# Verificar novamente
if lsof -ti:8080 >/dev/null 2>&1; then
    echo "❌ ERRO: Não foi possível liberar a porta 8080"
    exit 1
fi

echo "✅ Porta 8080 liberada"

# ===== Compilar Tailwind CSS =====
echo "🎨 Verificando Tailwind CSS..."
TAILWIND_INPUT="web/static/css/input.css"
TAILWIND_OUTPUT="web/static/css/tailwind.css"

# Verificar se precisa recompilar Tailwind
need_tailwind=false
if [ ! -f "$TAILWIND_OUTPUT" ]; then
    need_tailwind=true
elif [ "$TAILWIND_INPUT" -nt "$TAILWIND_OUTPUT" ]; then
    need_tailwind=true
fi

if [ "$need_tailwind" = true ]; then
    echo "🔄 Compilando Tailwind CSS..."
    # Verificar se npm está disponível
    if ! command -v npm >/dev/null 2>&1; then
        echo "❌ ERRO: npm não encontrado"
        exit 1
    fi
    npm run tailwind:build
    echo "✅ Tailwind CSS compilado"
else
    echo "✅ Tailwind CSS está atualizado"
fi

# ===== Compilar Templates =====
echo "📝 Verificando templates..."
need_templ=false

# Forçar regeneração para garantir que最新 alterações sejam incluídas
# (Mais seguro que verificar timestamps, que podem falhar)
for file in $(find web/components -name "*.templ" 2>/dev/null); do
    generated="${file%.templ}_templ.go"
    if [ ! -f "$generated" ] || [ "$file" -nt "$generated" ]; then
        need_templ=true
        break
    fi
done

# Também verifica se input.css foi modificado (indica possível mudança nos componentes)
if [ "$need_tailwind" = true ]; then
    need_templ=true
fi

if [ "$need_templ" = true ]; then
    echo "🔄 Recompilando templates..."
    
    # Tentar usar templ do PATH ou do GOPATH
    TEMPL_CMD=""
    if command -v templ >/dev/null 2>&1; then
        TEMPL_CMD="templ"
    elif [ -f "$HOME/go/bin/templ" ]; then
        TEMPL_CMD="$HOME/go/bin/templ"
    elif [ -f "$(go env GOPATH)/bin/templ" ]; then
        TEMPL_CMD="$(go env GOPATH)/bin/templ"
    else
        echo "📦 Instalando templ..."
        go install github.com/a-h/templ/cmd/templ@latest
        TEMPL_CMD="$HOME/go/bin/templ"
    fi
    
    $TEMPL_CMD generate
    if [ $? -ne 0 ]; then
        echo "❌ ERRO: Falha ao compilar templates"
        exit 1
    fi
    echo "✅ Templates recompilados"
else
    echo "✅ Templates estão atualizados"
fi

# ===== Compilar Go =====
echo "🔨 Compilando Arandu..."
go build -o arandu cmd/arandu/main.go
if [ $? -ne 0 ]; then
    echo "❌ ERRO: Falha na compilação"
    exit 1
fi

echo "🚀 Iniciando Arandu..."
./arandu > server.log 2>&1 &
sleep 3

# Verificar se o app está rodando
if curl -s "http://localhost:8080/login" -o /dev/null; then
    echo "✅ Arandu iniciado com sucesso em http://localhost:8080"
else
    echo "❌ ERRO: Arandu não respondeu"
    exit 1
fi
