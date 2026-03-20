Task: Integração Google OAuth (Feature)

ID da Tarefa: task_20260320_google_oauth

Requirement: REQ-07-03-04

Dependência: task_20260320_login_ui (Task 3) e task_20260320_session_middleware (Task 2)

Stack: Go, markbates/goth (recomendado) ou golang.org/x/oauth2, HTMX.

🎯 Objetivo

Implementar o fluxo completo de autenticação via Google OAuth 2.0. O sistema deve ser capaz de redirecionar o utilizador para o Google, processar o callback, validar a identidade e estabelecer a sessão no Arandu.

🛠️ Escopo Técnico

1. Configuração de Provedor (Infrastructure)

Local: internal/infrastructure/auth/google_provider.go.

Ação: Configurar o provedor OAuth usando variáveis de ambiente (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL).

Segurança: Implementar a validação de state para prevenir ataques CSRF.

2. Camada Web (Handlers)

GET /auth/google: Inicia o fluxo de autenticação (Redirecionamento para o Google).

GET /auth/google/callback:

Receber o código do Google.

Trocar pelo token de acesso e obter as informações do utilizador (E-mail e Nome).

Consultar o Banco Central (arandu_central.db) para verificar se o e-mail já existe.

Se existir: Criar a sessão e redirecionar para /dashboard.

Se NÃO existir: Encaminhar para a lógica de Provisão de Tenant (REQ-07-03-05) — que será a Task 5.

3. Camada Web (Navegação)

Garantir que o botão "Entrar com Google" na auth_page.templ não usa HTMX para o redirecionamento inicial (deve ser um link <a> direto ou window.location para evitar problemas de CORS no frame do Google).

🎨 Design System e UX

Feedback: Durante o redirecionamento, o utilizador deve ver uma mensagem simples de "A redirecionar para o Google..." se houver latência perceptível.

Erro: Caso o utilizador cancele o login no Google, ele deve ser devolvido à página de login com um aviso sutil: "Acesso via Google cancelado."

🧪 Protocolo de Testes "Ironclad"

A. Teste de Fluxo (Manual/Mock)

Clicar no botão Google na tela de login.

Verificar: O URL do browser muda para accounts.google.com.

Simular o retorno com um e-mail de teste.

Verificar: O sistema cria o cookie arandu_session e redireciona para o dashboard.

B. Teste de Segurança

Tentar aceder ao callback diretamente sem os parâmetros do Google.

Verificar: O sistema retorna 400 Bad Request ou redireciona para login com erro de segurança.

🛡️ Checklist de Integridade

[✓] As chaves do Google estão a ser lidas de variáveis de ambiente (sem hardcode)?

[✓] O e-mail retornado é validado contra o Banco Central?

[✓] O scripts/arandu_guard.sh confirma que as rotas /auth/google/* estão abertas?

[✓] O sistema lida corretamente com utilizadores que já têm conta via e-mail mas decidem entrar com Google?

---

status: completed

## ✅ Implementação Concluída

### Arquivos Criados:
- `internal/infrastructure/auth/google_provider.go` - Provedor OAuth
- `internal/infrastructure/auth/google_provider_test.go` - Testes do provedor

### Alterações:
- `internal/web/handlers/auth_handler.go` - Adicionados handlers OAuth (GoogleLogin, GoogleCallback, Logout)
- `internal/platform/middleware/auth.go` - Adicionadas rotas OAuth como públicas
- `cmd/arandu/main.go` - Rotas OAuth configuradas
- `go.mod` - Adicionada dependência golang.org/x/oauth2

### Funcionalidades:
- GET /auth/google - Redireciona para OAuth do Google
- GET /auth/google/callback - Processa callback e cria sessão
- Validação de state para prevenir CSRF
- Cookie oauth_state com expiração
- Tratamento de erros (cancelado, código inválido, etc.)
- Botão Google usa link direto (não HTMX)

### Testes:
- GoogleProvider: 3 testes (GenerateState, OAuthSession IsValid/Expired)
- AuthHandler: 6 testes (Login GET/POST, EmptyCredentials, InvalidMethod, GoogleLogin, ServeHTTP)
- Total: 9 testes passando