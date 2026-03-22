#!/bin/bash
# Script seguro de deployment do Arandu
# Garante que só existe uma instância do app rodando na porta 8080

echo "🛑 Parando todas as instâncias do Arandu..."
pkill -9 -f "arandu" 2>/dev/null
sleep 2

# Verificar se a porta 8080 está liberada
PORT_PID=$(lsof -ti:8080 2>/dev/null)
if [ ! -z "$PORT_PID" ]; then
    echo "⚠️  Porta 8080 ainda em uso pelo PID $PORT_PID. Finalizando..."
    kill -9 $PORT_PID 2>/dev/null
    sleep 1
fi

# Verificar novamente
if lsof -ti:8080 >/dev/null 2>&1; then
    echo "❌ ERRO: Não foi possível liberar a porta 8080"
    exit 1
fi

echo "✅ Porta 8080 liberada"

# Verificar se precisa recompilar templates
echo "📝 Verificando templates..."
need_templ=false
for file in $(find web/components -name "*.templ" 2>/dev/null); do
    generated="${file%.templ}_templ.go"
    if [ ! -f "$generated" ] || [ "$file" -nt "$generated" ]; then
        need_templ=true
        break
    fi
done

if [ "$need_templ" = true ]; then
    echo "🔄 Recompilando templates..."
    # Tentar usar templ do PATH ou do GOPATH
    if command -v templ >/dev/null 2>&1; then
        templ generate
    elif [ -f "$HOME/go/bin/templ" ]; then
        "$HOME/go/bin/templ" generate
    else
        echo "📦 Instalando templ..."
        go install github.com/a-h/templ/cmd/templ@latest
        "$HOME/go/bin/templ" generate
    fi
    if [ $? -ne 0 ]; then
        echo "❌ ERRO: Falha ao compilar templates"
        exit 1
    fi
    echo "✅ Templates recompilados"
else
    echo "✅ Templates estão atualizados"
fi

echo "🔨 Compilando Arandu..."
go build -o arandu cmd/arandu/main.go
if [ $? -ne 0 ]; then
    echo "❌ ERRO: Falha na compilação"
    exit 1
fi

echo "🚀 Iniciando Arandu..."
./arandu &
sleep 3

# Verificar se o app está rodando
if curl -s "http://localhost:8080/login" -o /dev/null; then
    echo "✅ Arandu iniciado com sucesso em http://localhost:8080"
else
    echo "❌ ERRO: Arandu não respondeu"
    exit 1
fi
