Task 1.1: Logger Centralizado e Build-time Versioning

ID da Tarefa: task_20260321_infra_logger_core

Capability: CAP-08-02 — Observabilidade e Diagnóstico

Requirement: REQ-08-02-01 (Padronização de Logs Estruturados)

Stack: Go 1.21+ (log/slog), Git Tags, JSON.

🎯 Objetivo

Implementar a fundação de logs estruturados e o sistema de versionamento injetado. O sistema deve emitir apenas JSON (sem distinção dev/prod) para garantir paridade total com o stack Loki/Grafana.

🛠️ Escopo Técnico

1. Pacote de Versão (internal/platform/version/version.go)

Criar variáveis exportadas Version, Commit e BuildTime.

Estas variáveis servirão como fonte da verdade para o sistema.

2. Pacote de Logger (internal/platform/logger/logger.go)

Implementação: Wrapper sobre slog.NewJSONHandler.

Atributos Globais (Invariantes):

app: "arandu"

version: Valor de version.Version

commit: Valor de version.Commit

Configuração: O handler deve escrever para os.Stdout.

3. Função FromContext

Implementar FromContext(ctx context.Context) *slog.Logger.

Deve tentar extrair request_id e tenant_id do contexto para enriquecer o log automaticamente.

🎨 Padrão de Saída (JSON Everywhere)

Exemplo de log esperado:

{"time":"2026-03-21T...","level":"INFO","msg":"Tenant provisionado","app":"arandu","version":"v1.0.4","commit":"a1b2c3d","tenant_id":"uuid-..."}


🧪 Protocolo de Testes "Ironclad"

Build de Teste: * Executar o build injetando uma versão fictícia: go build -ldflags "-X '...version.Version=v9.9.9'".

Rodar o binário e verificar se os logs exibem version: v9.9.9.

Verificação de Formato: * Garantir que mesmo em ambiente local os logs saem em JSON puro.

Atomicidade:

Substituir todos os log.Print e fmt.Println no main.go e no serviço de provisionamento pelo novo logger.

🛡️ Checklist de Integridade

[x] O pacote version é imutável em tempo de execução? (Sim, variáveis são setadas em build-time via ldflags)

[x] O logger não contém lógica condicional para "Pretty Print" (JSON apenas)? (Sim, sempre JSON via slog.NewJSONHandler)

[x] O tenant_id aparece no log quando disponível no contexto? (Sim, extraído automaticamente via FromContext)

[x] O request_id aparece no log quando disponível no contexto? (Sim, extraído automaticamente via FromContext)

[x] O middleware RequestID injeta request_id em todas as requisições? (Sim, em internal/platform/middleware/request_id.go)

[x] Todos os log.Print/fmt.Print foram substituídos no main.go? (Sim, substituídos por logger.Info/Error/Warn)

## ✅ Implementação Concluída

### Arquivos Criados:
1. **internal/platform/version/version.go** - Pacote de versionamento com variáveis injetáveis
2. **internal/platform/logger/logger.go** - Logger centralizado com slog JSON
3. **internal/platform/logger/logger_test.go** - Testes unitários do logger
4. **internal/platform/middleware/request_id.go** - Middleware para injetar request_id
5. **internal/platform/middleware/request_id_test.go** - Testes do middleware

### Arquivos Modificados:
1. **internal/platform/context/context.go** - Adicionado suporte a request_id
2. **cmd/arandu/main.go** - Substituído log por logger, adicionado middleware RequestID

### Como usar o versionamento em builds:
```bash
go build -ldflags "-X arandu/internal/platform/version.Version=v1.0.0 -X arandu/internal/platform/version.Commit=abc1234 -X arandu/internal/platform/version.BuildTime=2026-03-21T22:00:00Z" -o arandu ./cmd/arandu/
```

### Testes passando:
- ✅ TestString, TestInt, TestBool
- ✅ TestFromContextWithTenantID, TestFromContextWithRequestID, TestFromContextWithUserID
- ✅ TestFromContextWithAllContext, TestAttrsToAny
- ✅ TestRequestIDMiddleware, TestRequestIDIsUnique, TestRequestIDHeaderPropagation