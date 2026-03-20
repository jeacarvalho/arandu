Task: Tela de Login Botânica (Silent UI)

ID da Tarefa: task_20260320_login_ui

Requirement: REQ-07-03-01 e REQ-07-03-04

Dependência: task_20260320_session_middleware

Stack: Go, templ, Tailwind CSS.

🎯 Objetivo

Implementar a interface de login do Arandu seguindo a filosofia de "Tecnologia Silenciosa" e a Identidade Botânica. A tela deve oferecer suporte para entrada via E-mail/Senha e o botão de acesso rápido via Google OAuth.

🛠️ Escopo Técnico

1. Camada Web (Componente templ)

Local: web/components/auth/login.templ.

Estrutura:

Container: Fundo em --arandu-bg (#E1F5EE).

Logotipo: Texto "Arandu" em Source Serif 4 Itálico, cor --arandu-dark.

Formulário de E-mail/Senha:

Inputs no padrão Silent Input (apenas borda inferior, sem bordas laterais).

Labels em Inter (Sans) 10px, uppercase, cor --arandu-primary.

Botão Primário: Fundo --arandu-primary, texto branco, cantos arredondados (rounded-xl).

Botão Google: Estilo minimalista com o ícone oficial do Google e o texto "Entrar com Google".

2. Camada Web (Handlers & Rotas)

GET /auth/login: Renderiza a página de login completa (usando o Layout se necessário ou um layout simplificado de Auth).

POST /auth/login: Handler (stub) que recebe os dados. A lógica de validação de senha foi preparada na Task 1.

🎨 Design System: A Estética do Acolhimento

Fundo: #E1F5EE (Papel de seda).

Sombras: Evitar sombras pesadas. Se necessário, use shadow-sm apenas no card de login.

Tipografia: * Mensagens de boas-vindas: Source Serif 4.

Instruções de formulário: Inter.

Feedback de Erro: Texto em text-red-400 de 12px que aparece via HTMX abaixo do botão de login em caso de falha.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Renderização

Aceder a /auth/login.

Verificar: O fundo verde-água está correto? A fonte Serif é usada no logotipo?

Verificar: O botão do Google está alinhado e visível?

B. Teste de Responsividade

Redimensionar para 375px (Mobile).

Verificar: O formulário ocupa a largura adequada com margens laterais seguras? Os campos de toque têm pelo menos 44px de altura?

C. Teste HTMX

Tentar submeter o formulário vazio.

Verificar: O HTMX dispara o POST e o backend (ainda que em modo stub) retorna uma resposta de validação parcial?

🛡️ Checklist de Integridade

[✓] O componente usa a tipografia Source Serif 4 para elementos narrativos?

[✓] O padrão Silent Input foi respeitado (sem "rings" azuis no foco)?

[✓] O arquivo .templ foi gerado com sucesso via templ generate?

[✓] O scripts/arandu_guard.sh confirma que a nova rota de login está acessível?

---

status: completed

## ✅ Implementação Concluída

### Arquivos Criados:
- `web/components/auth/login.templ` - Componente de login Botânica
- `web/components/auth/login_templ.go` - Arquivo gerado
- `web/components/auth/login_test.go` - Testes do componente
- `internal/web/handlers/auth_handler.go` - Handler de autenticação
- `internal/web/handlers/auth_handler_test.go` - Testes do handler

### Funcionalidades:
- Página de login com designbotânico (fundo #E1F5EE)
- Logotipo em Source Serif 4 Itálico
- Silent Inputs (apenas borda inferior)
- Labels em Inter, uppercase
- Botão Google com ícone oficial
- Rota GET/POST /login e /auth/login
- Formulário com suporte HTMX

### Testes Automatizados:
- LoginData: 2 testes
- AuthHandler: 4 testes
- Total: 6 testes passando

### Verificações:
- /login retorna 200 (página acessível)
- Middleware redireciona rotas protegidas para /login (302)
- Arquivo .templ gerado com sucesso