REQ-08-01-02 — Serviço de Auditoria Centralizada (Control Plane)

Identificação

ID: REQ-08-01-02

Capability: CAP-08-03 — Auditoria e Conformidade

Vision: VISION-07 — Organização Operacional (Compliance)

Status: draft

História do utilizador

Como administrador do sistema Arandu, quero que todas as ações críticas realizadas pelos utilizadores e pela infraestrutura sejam registradas de forma imutável no Banco Central, para garantir a rastreabilidade total, auditoria de segurança e conformidade com normas de proteção de dados (LGPD).

Contexto

Diferente dos logs técnicos (orientados a depuração), os registros de auditoria são evidências de negócio. Eles residem no Control Plane (arandu_central.db) para garantir que, mesmo que um banco de dados clínico (Tenant) seja removido ou corrompido, a trilha de quem acedeu a esses dados permaneça preservada para fins legais e de suporte.

Descrição funcional

O sistema deve implementar um motor de auditoria persistente com as seguintes características:

Imutabilidade: Os registros de auditoria são "append-only". Não deve existir lógica no sistema que permita a edição ou exclusão de entradas na tabela audit_logs.

Centralização: Todos os eventos, independentemente do Tenant de origem, devem ser gravados no banco de dados central.

Execução Assíncrona: Para não impactar a percepção de performance do terapeuta (Tecnologia Silenciosa), o disparo do log de auditoria deve ocorrer em segundo plano (goroutines).

Isolamento de PHI: O serviço de auditoria NUNCA deve registrar informações de saúde protegidas (conteúdo de notas, diagnósticos). Deve registrar apenas IDs de recursos e verbos de ação.

Estrutura do Registro (Audit Schema)

id: UUID v4 único do registro.

timestamp: Momento exato da ação (UTC).

user_id: Identificador do profissional que realizou a ação.

tenant_id: Identificador do banco de dados clínico afetado.

action: Verbo normalizado (ex: LOGIN, VIEW_PATIENT, EXPORT_REPORT, UPDATE_MEDICATION).

resource_id: ID opcional do objeto afetado (ex: UUID do paciente).

metadata: JSON contendo ip_address e user_agent para análise forense.

Eventos Mandatórios para Auditoria

Acessos: Login, Logout, falha de autenticação.

Privacidade: Visualização de prontuário, abertura de ficha de paciente.

Modificação: Alteração de histórico farmacológico ou vitais.

Exportação: Download de backups, geração de PDFs de relatório.

Critérios de Aceitação

CA-01: O sistema deve extrair automaticamente o user_id e tenant_id do contexto da requisição para o log.

CA-02: Uma falha na escrita do log de auditoria não deve interromper a ação do utilizador, mas deve emitir um ERROR crítico no log de sistema.

CA-03: A consulta aos logs de auditoria deve ser restrita apenas a utilizadores com privilégio de administrador do sistema.

CA-04: O campo metadata deve ser preenchido corretamente com as informações de rede do cliente.

Nota de Soberania: Este dado reside exclusivamente no Banco de Dados Central (arandu_central.db).