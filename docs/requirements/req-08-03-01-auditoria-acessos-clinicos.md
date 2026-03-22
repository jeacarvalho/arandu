REQ-08-03-01 — Auditoria de Acessos Clínicos

Identificação

ID: REQ-08-03-01

Capability: CAP-08-03 — Auditoria e Conformidade

Vision: VISION-07 — Organização Operacional (Compliance)

Status: draft

História do utilizador

Como psicólogo clínico e proprietário dos dados, quero que todas as consultas a prontuários e exportações de dados sejam registradas de forma imutável, para garantir a segurança jurídica e a conformidade com as normas de sigilo profissional.

Contexto

Diferente dos logs técnicos (REQ-08-02-01), os logs de auditoria são registros de negócio. Eles provam que o médico acessou o paciente X na data Y. Estes dados devem ser persistidos no Banco Central (arandu_central.db) pois sobrevivem à exclusão de um tenant específico e servem para auditorias externas.

Descrição funcional

Implementar o serviço de registro de trilha de auditoria:

Eventos Auditáveis: Login, Acesso a Prontuário, Edição de Sessão, Exportação de PDF, Alteração de Medicação.

Imutabilidade: Os registros de auditoria só permitem inserção (Append Only).

Consolidação Central: Os registros devem ser gravados na tabela audit_logs do banco central.

Dados de Auditoria

timestamp: Data e hora do acesso.

user_id: ID do médico que realizou a ação.

tenant_id: ID do banco clínico acessado.

action: Verbo da ação (ex: VIEW_PATIENT, UPDATE_GOAL).

resource_id: ID do paciente ou registro afetado.

metadata: JSON com informações de IP e User-Agent.

Interface (Admin Hub)

O administrador do Arandu (ou o próprio médico em sua área de segurança) poderá visualizar uma tabela cronológica desses acessos através do Dashboard de Administração.

Critérios de Aceitação

CA-01: Toda leitura de prontuário deve gerar uma entrada na tabela de auditoria.

CA-02: O registro deve ser assíncrono para não travar a UI (goroutine), mas deve ser garantido.

CA-03: A tabela de auditoria deve ser isolada dos bancos clínicos individuais.

CA-04: O médico deve poder ver quem (qual credencial dele) acessou os dados em caso de consultórios com secretárias (futuro).

Nota: Este dado reside no Banco de Dados Central (arandu_central.db).