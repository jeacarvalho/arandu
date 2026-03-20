 Plano de Sprint: Autenticação e Isolamento de Dados

Este plano divide a implementação da segurança e multi-tenancy em micro-tarefas atómicas para garantir a precisão do desenvolvimento.

📋 Lista de Requisitos

[ ] REQ-07-03-01: Autenticação E-mail/Senha (Central DB).

[ ] REQ-07-03-04: Integração Google OAuth.

[ ] REQ-07-03-05: Provisão Automática de Tenant e Criação de DB Físico.

🛠️ Sequência de Execução (Tasks Atómicas)

Task 1: O "Control Plane" (Infra)

Ação: Criar o banco central arandu_central.db e as tabelas de users e tenants.

Sucesso: Binário inicializa e cria o banco central vazio mas estruturado.

Task 2: Middleware de Sessão (Core)

Ação: Implementar o middleware que lê o cookie e injeta o tenant_id no contexto.

Sucesso: Rotas protegidas retornam 401 se o cookie não estiver presente.

Task 3: Tela de Login Botânica (UI)

Ação: Implementar a UI de Login (E-mail/Senha + Botão Google) no padrão Silent UI.

Sucesso: Interface renderiza perfeitamente no mobile e desktop.

Task 4: Google OAuth Integration (Feature)

Ação: Configurar as rotas de callback do Google e validação de JWT/Cookie.

Sucesso: Utilizador é redirecionado para o Google e volta para o app com sucesso.

Task 5: Onboarding e Provisão (Flow)

Ação: Lógica de criação física do arquivo .db para novos utilizadores.

Sucesso: Ao logar pela primeira vez, um novo arquivo SQLite surge na pasta storage/tenants/.