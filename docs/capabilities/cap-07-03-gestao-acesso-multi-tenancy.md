CAP-07-03 — Gestão de Acesso e Multi-tenancy

Identificação

ID: CAP-07-03

Título: Gestão de Acesso e Multi-tenancy

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

Descrição

Prover uma infraestrutura de acesso seguro onde múltiplos terapeutas podem utilizar a mesma instância da aplicação web, garantindo o isolamento físico e absoluto dos seus dados clínicos através de uma arquitetura de banco de dados por usuário.

Objetivos

Autenticação Segura: Validar a identidade do profissional antes de conceder acesso ao ambiente clínico.

Isolamento de Dados (Soberania): Garantir que um terapeuta nunca tenha acesso, mesmo que acidental, aos dados de outro profissional.

Portabilidade: Permitir que o banco de dados individual (SQLite) possa ser exportado ou migrado de forma independente.

Escalabilidade Administrativa: Facilitar a manutenção e atualização de múltiplos bancos de dados simultaneamente.

Contexto Técnico (SOTA)

Diferente de sistemas tradicionais que usam uma tabela única com filtros de user_id, esta capability implementa o modelo Database-per-tenant:

Control Plane: Um banco central gerencia as credenciais e o mapeamento de qual arquivo .db pertence a qual usuário.

Data Plane: Milhares de arquivos SQLite individuais que são "montados" dinamicamente na sessão do usuário.

Requisitos Relacionados

REQ-07-03-01: Autenticação de Usuário (Login/Logout).

REQ-07-03-02: Orquestração de Conexão Dinâmica (Middleware de Seleção de DB).

REQ-07-03-03: Gestão de Migrações em Lote (Updates de Schema Multi-tenant).

Valor de Negócio

Para o terapeuta, esta capacidade representa a Soberania do Dado. Em caso de saída da plataforma ou auditoria, o profissional possui um arquivo físico exclusivo com seu histórico clínico, reforçando a confiança e a conformidade com normas de sigilo profissional e proteção de dados.

Critérios de Sucesso

Um usuário logado deve visualizar apenas os seus pacientes e sessões.

A performance de abertura de conexão com o banco individual não deve impactar a experiência do HTMX (latência mínima).

Falhas em um banco de dados individual não devem afetar a disponibilidade do sistema para outros usuários.