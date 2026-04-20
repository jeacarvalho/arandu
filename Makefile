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
test: test-unit test-e2e

test-unit:
	@echo "🧪 Executando unit tests..."
	@bash tests/run_unit.sh

test-e2e:
	@echo "🌐 Executando E2E tests..."
	@bash tests/run_e2e.sh

test-all: test-unit test-e2e

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

# Monitoring stack
monitor-up:
	@echo "🚀 Iniciando stack de monitoramento..."
	@docker-compose -f docker-compose.monitoring.yml up -d
	@echo "⏳ Aguardando serviços iniciarem..."
	@sleep 5
	@echo "✅ Grafana: http://localhost:3000 (admin/arandu2024)"
	@echo "✅ Loki: http://localhost:3100"

monitor-down:
	@echo "🛑 Parando stack de monitoramento..."
	@docker-compose -f docker-compose.monitoring.yml down
	@echo "✅ Stack parado"

monitor-logs:
	@docker-compose -f docker-compose.monitoring.yml logs -f

monitor-status:
	@docker-compose -f docker-compose.monitoring.yml ps

monitor-clean:
	@echo "🧹 Removendo stack e volumes..."
	@docker-compose -f docker-compose.monitoring.yml down -v
	@echo "✅ Stack e volumes removidos"

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo " make build - Compilar templates e binário"
	@echo " make templ - Compilar apenas templates"
	@echo " make deploy - Deploy seguro (build + start)"
	@echo " make clean - Limpar arquivos gerados"
	@echo " make test - Executar testes"
	@echo " make golint - Executar lint"
	@echo " make dev - Modo desenvolvimento com hot reload"
	@echo ""
	@echo " Monitoring:"
	@echo " make monitor-up - Iniciar stack (Loki, Promtail, Grafana)"
	@echo " make monitor-down - Parar stack"
	@echo " make monitor-logs - Ver logs em tempo real"
	@echo " make monitor-status - Status dos containers"
	@echo " make monitor-clean - Remover stack e volumes"
	@echo ""
	@echo " make help - Mostrar esta ajuda"
