REQ-07-03-04 — Integração Google OAuth

Identificação

ID: REQ-07-03-04

Capability: CAP-07-03 — Gestão de Acesso e Multi-tenancy

Status: draft

História do utilizador

Como psicólogo clínico, quero entrar no sistema usando a minha conta Google, para facilitar o meu acesso sem precisar de memorizar mais uma senha, mantendo a segurança da minha conta principal.

Descrição funcional

O sistema deve integrar-se com o Google Identity Services (OAuth 2.0).

Botão de Login: Exibir "Entrar com Google" com o ícone oficial.

Mapeamento: O e-mail retornado pelo Google deve ser a chave primária de busca no Banco Central.

Novo Utilizador: Se o e-mail não existir, o sistema cria um novo Tenant automaticamente.

Critérios de Aceitação

CA-01: O redirecionamento deve usar o parâmetro state para evitar CSRF.

CA-02: O sistema deve solicitar apenas o escopo openid, email e profile.

CA-03: Em caso de sucesso, o utilizador deve ser autenticado e redirecionado para o Dashboard.

CA-04: O layout da página de login deve permanecer responsivo.

Persistência

Os dados de tokens de acesso não devem ser persistidos no banco clínico, apenas a referência do e-mail no Banco Central.