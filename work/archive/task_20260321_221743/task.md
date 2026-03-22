Task 1.2: Middleware de Telemetria e Rastreabilidade

ID da Tarefa: task_20260321_infra_telemetry_middleware

Requirement: REQ-08-01-01

Dependência: task_20260321_infra_logger_core (Task 1.1)

Stack: Go, chi (ou router padrão), Context.

🎯 Objetivo

Implementar o middleware global que gera um request_id único para cada interação e loga automaticamente a performance e o desfecho de todas as requisições HTTP.

🛠️ Escopo Técnico

1. Middleware de RequestID

Gerar um UUID v4 (ou um hash curto) para cada requisição.

Injetar no context sob uma chave tipada (evitar colisões).

Adicionar o header X-Request-ID na resposta HTTP para debug do lado do cliente.

2. Middleware de Log de Tráfego

Captura: Capturar Início, Fim, Status Code e Latência (Duration).

Logging: Utilizar o logger da Task 1.1 para emitir uma linha de log INFO ao final de cada request.

Campos: method, path, status, duration_ms, ip, request_id.

3. Integração com TenantID

O middleware deve cooperar com o AuthMiddleware existente para garantir que, se o utilizador estiver logado, o tenant_id também seja incluído nos logs de telemetria.

🧪 Protocolo de Testes "Ironclad"

Rastreio de Fluxo: * Fazer uma requisição. Pegar o X-Request-ID no header do browser.

Procurar no terminal por esse ID. Deve existir exatamente uma linha de log com todos os detalhes da requisição.

Métrica de Latência: * Verificar se o campo duration_ms está a reportar valores realistas.

Filtragem de Sensibilidade:

Garantir que se a URL contiver algo como /auth/google/callback?code=..., o parâmetro sensível seja omitido ou o log focado apenas no path.

🛡️ Checklist de Integridade

[x] O request_id é propagado para todas as camadas através do context.Context? (Sim, via middleware.RequestIDMiddleware)

[x] O middleware de log é o primeiro (ou um dos primeiros) na cadeia do router? (Sim, aplicado após RequestID e Auth)

[x] O request_id aparece nos logs de telemetria? (Sim, extraído automaticamente pelo logger.FromContext)

[x] A filtragem de dados sensíveis está implementada? (Sim, campos code, token, password, secret, api_key, access_token, refresh_token são redacted)

[x] Os testes unitários cobrem o middleware? (Sim, testes em telemetry_test.go)

## ✅ Implementação Concluída

### Arquivos Criados:
1. **internal/platform/middleware/telemetry.go** - Middleware de telemetria HTTP
2. **internal/platform/middleware/telemetry_test.go** - Testes unitários

### Arquivos Modificados:
1. **cmd/arandu/main.go** - Integração do middleware na chain (após RequestID, antes de Recovery)

### Funcionalidades:
- **Captura de métricas**: method, path, status, duration_ms, ip, request_id
- **Sanitização de URLs**: remove parâmetros sensíveis (code, token, password, etc.)
- **Skip paths**: /static/ é ignorado para não poluir logs
- **IP do cliente**: extrai de X-Forwarded-For, X-Real-Ip ou RemoteAddr
- **Integração com logger**: usa o logger centralizado com request_id do contexto

### Testes passando:
- ✅ TestTelemetryMiddleware
- ✅ TestTelemetryMiddlewareShouldSkip  
- ✅ TestSanitizePath
- ✅ TestGetClientIP
- ✅ TestResponseWriterWrapper