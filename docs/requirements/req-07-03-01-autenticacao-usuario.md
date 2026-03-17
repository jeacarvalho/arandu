REQ-07-03-01 — Autenticação de Usuário

Identificação

ID: REQ-07-03-01

Capability: CAP-07-03 Gestão de Acesso e Multi-tenancy

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do usuário

Como psicólogo clínico, quero me autenticar de forma segura no sistema, para acessar meu consultório digital exclusivo e garantir que apenas eu tenha acesso aos dados dos meus pacientes.

Contexto

A autenticação é o gatilho para a orquestração de multi-tenancy. Ao validar as credenciais no Banco Central (Control Plane), o sistema identifica qual arquivo SQLite clínico deve ser "montado" para a sessão do usuário. O login deve ser simples, mas robusto o suficiente para proteger informações sensíveis de saúde.

Descrição funcional

O sistema deve prover uma interface de login para terapeutas cadastrados.

Entrada: E-mail e senha.

Segurança: As senhas devem ser armazenadas utilizando hashes seguros (ex: Argon2 ou bcrypt).

Sessão: Utilização de cookies seguros ou JWT para manter o estado da sessão.

Logout: Encerrar a sessão e garantir o fechamento seguro das conexões com o banco de dados clínico.

Dados de Acesso (Banco Central)

Tabela users

ID: UUID (Chave primária).

Email: String única (Login).

PasswordHash: String (Hash da senha).

TenantID: String única (Utilizada para localizar o arquivo .db clínico).

CreatedAt / LastLogin: Timestamps.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Estética: Tela de login extremamente limpa, sem distrações visuais, com foco total nos campos de entrada.

Feedback: Mensagens de erro genéricas ("Credenciais inválidas") para evitar enumeração de usuários.

Loading: Indicador visual sutil durante a validação para evitar múltiplos cliques.

Fluxo

O usuário acessa a rota /login.

Preenche e-mail e senha.

O sistema consulta o Banco Central.

Sucesso:

Gera token de sessão.

Identifica o TenantID.

Redireciona para o /dashboard.

O middleware de conexão (REQ-07-03-02) abre o banco clínico correspondente.

Falha: Retorna erro visual via HTMX sem recarregar a página.

Rotas Esperadas

GET /login: Renderiza a página de login.

POST /login: Processa a autenticação.

POST /logout: Finaliza a sessão.

Critérios de Aceitação

CA-01: O sistema não deve permitir o acesso a nenhuma rota clínica sem uma sessão ativa.

CA-02: A senha nunca deve ser armazenada em texto plano.

CA-03: O sistema deve suportar redirecionamento automático para a página solicitada após o login bem-sucedido.

CA-04: O logout deve invalidar completamente o token de sessão tanto no cliente quanto no servidor.

CA-05: A interface de login deve ser responsiva e funcional em dispositivos móveis.

Persistência (SQL Central)

-- Busca de usuário para autenticação
SELECT id, password_hash, tenant_id FROM users WHERE email = ? LIMIT 1;

-- Registro de log de acesso
UPDATE users SET last_login = ? WHERE id = ?;


Fora do escopo

Recuperação de senha via e-mail (REQ-07-03-04).

Autenticação de dois fatores (MFA).

Registro de novos usuários (Self-service signup).