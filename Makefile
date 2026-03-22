# Arandu - Makefile
# Garante que templates sejam recompilados antes do build

.PHONY: all build templ deploy clean test lint

# Default target
all: build

# Instalar templ se necessário e gerar arquivos .go
# Detecta se arquivos .templ foram modificados
templ:
	@echo "📝 Verificando templates..."
	@# Verificar se templ está instalado
	@if [ ! -f "$(HOME)/go/bin/templ" ] && [ ! -f "/usr/local/bin/templ" ]; then \
		echo "📦 Instalando templ..."; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	fi
	@# Detectar se precisa recompilar
	@need_recompile=false; \
	for file in $$(find web/components -name "*.templ" 2>/dev/null); do \
		generated="$${file%.templ}_templ.go"; \
		if [ ! -f "$$generated" ] || [ "$$file" -nt "$$generated" ]; then \
			need_recompile=true; \
			break; \
		fi; \
	done; \
	if [ "$$need_recompile" = true ]; then \
		echo "🔄 Recompilando templates..."; \
		$(HOME)/go/bin/templ generate || templ generate || { echo "❌ Falha ao compilar templates"; exit 1; }; \
		echo "✅ Templates recompilados"; \
	else \
		echo "✅ Templates estão atualizados"; \
	fi

# Build completo (templates + binário)
build: templ
	@echo "🔨 Compilando Arandu..."
	@go build -o arandu cmd/arandu/main.go
	@echo "✅ Build concluído"

# Deploy seguro (matar processos, build, iniciar)
deploy: build
	@echo "🚀 Iniciando deploy..."
	@./scripts/safe_deploy.sh

# Limpar arquivos gerados
clean:
	@echo "🧹 Limpando arquivos gerados..."
	@rm -f arandu
	@find web/components -name "*_templ.go" -delete
	@echo "✅ Limpo"

# Executar testes
test:
	@echo "🧪 Executando testes..."
	@go test ./... -v

# Lint
golint:
	@echo "🔍 Executando lint..."
	@gofmt -d .
	@go vet ./...

# Desenvolvimento (recompilar templates automaticamente em loop)
dev:
	@echo "👁️  Modo desenvolvimento - observando mudanças..."
	@which templ > /dev/null || go install github.com/a-h/templ/cmd/templ@latest
	@$(HOME)/go/bin/templ generate --watch &
	@go run cmd/arandu/main.go

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make build    - Compilar templates e binário"
	@echo "  make templ    - Compilar apenas templates"
	@echo "  make deploy   - Deploy seguro (build + start)"
	@echo "  make clean    - Limpar arquivos gerados"
	@echo "  make test     - Executar testes"
	@echo "  make golint   - Executar lint"
	@echo "  make dev      - Modo desenvolvimento com hot reload"
	@echo "  make help     - Mostrar esta ajuda"
