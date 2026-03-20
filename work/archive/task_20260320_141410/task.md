Task: Middleware de Sessão e Orquestração de Tenant

ID da Tarefa: task_20260320_session_middleware

Requirement: REQ-07-03-02

Dependência: task_20260320_control_plane_infra (Task 1)

Stack: Go, context, http.Cookie.

🎯 Objetivo

pré requisito: leia o work/plano_atual.md para entender todo o contexto onde a tarefa está incluída.

Implementar o middleware de autenticação e orquestração que identifica o utilizador através da sessão, localiza o seu tenant_id no Banco Central e injeta a conexão SQLite clínica correta no contexto da requisição.

🛠️ Escopo Técnico

1. Camada de Plataforma (Middleware)

Local: internal/platform/middleware/auth.go.

Lógica:

Extrair o cookie arandu_session.

Validar o token contra a tabela sessions (ou users) no Control Plane.

Recuperar o tenant_id associado.

Injetar o tenant_id no contexto: ctx := context.WithValue(r.Context(), "tenant_id", tenantID).

2. Gerenciador de Conexões (Tenant Pool)

Local: internal/infrastructure/repository/sqlite/tenant_pool.go.

Funcionalidade:

Implementar um Map thread-safe (sync.RWMutex) para armazenar conexões abertas: map[tenantID]*sql.DB.

Método GetConnection(tenantID):

Se a conexão existe no map e está ativa (Ping()), retorna-a.

Se não, abre o arquivo storage/tenants/clinical_{tenantID}.db, aplica as pragmas SOTA (WAL mode, foreign_keys), adiciona ao map e retorna.

Integrar o Migrator clínico para rodar automaticamente na primeira abertura de cada banco.

3. Injeção de Dependência no Contexto

O middleware deve chamar o tenantPool.GetConnection(tenantID) e injetar o *sql.DB resultante no contexto sob a chave tenant_db.

🎨 Design de Erros (Tecnologia Silenciosa)

Não autorizado: Se a sessão for inválida em rotas protegidas, redirecionar via http.Redirect para /login.

Erro de Banco: Se o banco do tenant estiver inacessível, retornar uma página de "Manutenção Temporária" minimalista, sem expor detalhes técnicos.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Middleware

Tentar acessar /dashboard sem cookie -> Esperado: Redirecionamento para /login.

Tentar acessar com um cookie de um utilizador inexistente -> Esperado: 401 Unauthorized.

B. Teste de Performance (Pool)

Realizar 50 requisições HTMX rápidas para o mesmo tenant.

Verificar: O número de arquivos abertos no sistema operacional (lsof) não deve crescer linearmente; a conexão deve ser reutilizada.

🛡️ Checklist de Integridade

[✓] O tenant_id é extraído de forma segura (sem possibilidade de injeção)?

[✓] O tenantPool utiliza RWMutex para evitar race conditions em acessos concorrentes?

[✓] O contexto da requisição é limpo adequadamente?

[✓] O scripts/arandu_guard.sh confirma que as rotas protegidas estão bloqueadas para anónimos?

---

status: completed

## ✅ Implementação Concluída

### Arquivos Criados:
- `internal/infrastructure/repository/sqlite/tenant_pool.go` - Pool de conexões para tenants
- `internal/platform/middleware/auth.go` - Middleware de autenticação
- `internal/platform/middleware/auth_test.go` - Testes do middleware
- `internal/infrastructure/repository/sqlite/tenant_pool_test.go` - Testes do pool

### Alterações:
- `internal/infrastructure/repository/sqlite/central_db.go` - Adicionada tabela sessions
- `cmd/arandu/main.go` - Integracao do TenantPool e AuthMiddleware
- Adicionado import do package middleware

### Funcionalidades Implementadas:
- TenantPool com RWMutex para thread-safety
- GetConnection com cache de conexões
- Auto-migration na primeira abertura do banco
- Middleware extrai cookie de sessão
- Redirecionamento para /login em rotas protegidas sem sessão
- Página de manutenção em caso de erro de banco
- Tabela sessions adicionada ao banco central

### Testes Automatizados:
- TenantPool: 8 testes (New, GetActiveCount, IsConnected, CloseConnection, CloseAll, GetConnection)
- AuthMiddleware: 9 testes (isPublicRoute, GetTenantID, GetTenantDB, GetUserID, PublicRoute, NoSessionCookie, ExpiredSession)
- Total: 17 testes passando