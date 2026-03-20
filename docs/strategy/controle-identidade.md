Estratégia: Autenticação, Google OAuth e Identidade

Este documento define como o Arandu gere a identidade do profissional e a orquestração de multi-tenancy (um banco por utilizador).

🏗️ 1. Arquitetura de Identidade (Control Plane)

A autenticação é o porteiro do sistema. Ela reside num banco centralizado, enquanto os dados clínicos residem em "shards" (ficheiros SQLite individuais).

Identidade Híbrida: O sistema deve suportar Login Tradicional (E-mail/Senha) e Google OAuth simultaneamente, vinculando ambas ao mesmo tenant_id.

Segurança de Sessão: Uso de cookies seguros e criptografados (HttpOnly, Secure, SameSite=Lax).

Provisão de Tenant: Na primeira vez que um utilizador entra via Google, o sistema deve:

Criar o registo no Banco Central.

Gerar um tenant_id (UUID).

Criar o ficheiro storage/tenants/clinical_{tenant_id}.db.

Executar as migrations iniciais.

🌐 2. Integração Google OAuth

Utilizaremos a biblioteca markbates/goth ou o pacote oficial golang.org/x/oauth2.

Escopo Mínimo: Apenas email e profile.

Fluxo:

Utilizador clica em "Entrar com Google".

Redirecionamento para o consentimento do Google.

Callback recebe o e-mail.

O sistema verifica se o e-mail existe no Banco Central.

Se não existir, inicia o fluxo de Onboarding Silencioso (Criação de Tenant).

🛠️ 3. Middleware de Orquestração de Dados

O Middleware deve ser "burro" e rápido:

Extrai o tenant_id da sessão.

Injeta o caminho do banco no contexto da requisição.

O Repositório abre a conexão sob demanda (Lazy Loading) ou usa um Pool de Conexões (LRU) para evitar excesso de ficheiros abertos no sistema operacional.

Nota de Soberania: Este documento é o SSoT para qualquer implementação de acesso. Alterações aqui exigem revisão de segurança.