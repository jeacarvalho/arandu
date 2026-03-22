REQ-08-02-01 — Padronização de Logs Estruturados

Identificação

ID: REQ-08-02-01

Capability: CAP-08-02 — Observabilidade e Diagnóstico

Vision: VISION-01 — Registro da Prática Clínica

Status: draft

História do utilizador

Como administrador do sistema, quero que todos os eventos do servidor sejam registrados em formato JSON estruturado, para facilitar a indexação por ferramentas de análise e permitir a filtragem de erros por versão e utilizador.

Contexto

O Arandu opera num modelo multi-tenant complexo. Logs textuais legados (ex: fmt.Println) são impossíveis de analisar em larga escala. A mudança para slog (JSON) transforma o log numa base de dados pesquisável, onde cada linha contém o contexto necessário para resolver o problema sem "adivinhação".

Descrição funcional

O sistema deve implementar um motor de logging centralizado:

Formato Único: JSON em todos os ambientes (Dev, Staging, Prod).

Injeção de Versão: Todo log deve carregar a tag do Git e o hash do commit injetados via ldflags.

Context Aware: O logger deve ser capaz de extrair o tenant_id do contexto da requisição automaticamente.

Níveis de Severidade: Uso rigoroso de DEBUG, INFO, WARN e ERROR.

Dados do Log (Atributos Obrigatórios)

time: ISO8601 Timestamp.

level: Nível do evento.

msg: Descrição humana do evento.

version: Tag do Git (ex: v1.2.3).

commit: Hash curto do Git.

tenant_id: UUID do médico (se disponível).

request_id: Identificador único da transação HTTP.

Interface (Visualização de Admin)

Desenvolvimento: Logs coloridos no terminal via ferramenta de CLI (opcional) ou JSON cru.

Produção: Dashboard Grafana lendo do Loki.

Critérios de Aceitação

CA-01: O sistema não deve utilizar log.Print ou fmt.Print para mensagens de sistema.

CA-02: O campo version não deve ser "empty" se o binário foi compilado via pipeline.

CA-03: Dados clínicos sensíveis (notas de sessões) NUNCA devem ser incluídos nas mensagens de log.

CA-04: Erros (nível ERROR) devem incluir o stack trace ou contexto da função falha.

Nota: Este log é direcionado para infraestrutura e não para auditoria clínica.