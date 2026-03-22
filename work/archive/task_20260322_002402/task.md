Task: Implementação do Serviço de Auditoria (Control Plane)

ID da Tarefa: task_20260321_infra_audit_service

Requirement: REQ-08-01-02 (Auditoria e Conformidade)

Dependência: task_20260321_infra_logger (Task 1)

Stack: Go, SQLite (Central), context.

🎯 Objetivo

Implementar o serviço de registo de auditoria imutável no arandu_central.db. Este serviço deve capturar acções críticas de soberania de dados, garantindo que o administrador saiba sempre quem, quando, onde (tenant) e o quê (recurso) foi acedido.

🛠️ Escopo Técnico

1. Camada de Infraestrutura (Migration Central)

Arquivo: internal/infrastructure/repository/sqlite/migrations_central/0002_add_audit_logs.up.sql.

Schema SOTA:

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL,
    tenant_id TEXT NOT NULL,      -- Identifica o banco clínico de origem
    action TEXT NOT NULL,         -- ex: 'ACCESS_PATIENT', 'EXPORT_DATA', 'LOGIN'
    resource_id TEXT,             -- ex: UUID do paciente ou da sessão
    ip_address TEXT,
    user_agent TEXT
);
CREATE INDEX IF NOT EXISTS idx_audit_tenant ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id);


2. Camada de Aplicação (internal/application/services/audit_service.go)

Criar o AuditService com o método Log(ctx, action, resourceID).

Automatização: O serviço deve extrair automaticamente user_id e tenant_id do contexto (injetados pelo middleware de auth).

Async: O registo de auditoria deve ser disparado numa goroutine (Fire and Forget) para não adicionar latência à experiência do terapeuta, mas deve ter um mecanismo de recuperação em caso de falha de escrita no banco central.

3. Integração nos Serviços Clínicos

Injetar o AuditService no PatientService.

Ação: No método GetPatientByID, disparar a auditoria: s.audit.Log(ctx, "ACCESS_PATIENT", patientID).

🎨 Design System (Monitoração)

Os logs de auditoria são dados frios. No terminal, o logger da Task 1 deve emitir uma linha INFO sempre que um log de auditoria for persistido com sucesso no banco central.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Rastreabilidade

Logar como Médico A (Tenant 1).

Aceder ao prontuário do Paciente X.

Verificar: No banco arandu_central.db, deve existir uma nova linha com:

user_id do Médico A.

tenant_id do Médico A.

action = 'ACCESS_PATIENT'.

resource_id = UUID do Paciente X.

B. Teste de Performance

Simular 100 acessos rápidos.

Verificar: Se o sistema continua responsivo (indicando que a auditoria está a correr de forma assíncrona corretamente).

🛡️ Checklist de Integridade

[x] O tenant_id está a ser gravado corretamente em todas as acções clínicas?

[x] A tabela de auditoria é protegida (sem rota de DELETE ou UPDATE exposta)?

[x] O scripts/arandu_guard.sh confirma a criação da tabela no banco central?

---

## Status: ✅ COMPLETADO

**Data:** 2026-03-22

**Resumo das alterações:**

1. ✅ Migration 0002_add_audit_logs adicionada ao central_db.go:
   - Tabela `audit_logs` com índices para tenant, user, timestamp e action

2. ✅ AuditService criado em `internal/application/services/audit_service.go`:
   - Log assíncrono (fire-and-forget) via channel
   - Extrai user_id e tenant_id do contexto
   - Worker goroutine para persistência
   - Método Close() para shutdown graceful

3. ✅ PatientService integrado com AuditService:
   - `NewPatientServiceWithAudit()` com suporte a auditoria
   - CREATE_PATIENT logado após criação
   - ACCESS_PATIENT logado após acesso

4. ✅ Testes criados:
   - TestAuditService_Log
   - TestAuditService_AsyncBehavior
   - TestAuditService_Close
   - TestAuditService_LogWithoutContext
   - TestAuditService_GetLogsByTenant

5. ✅ Verificação em produção:
   - Migration 0002 aplicada com sucesso
   - Audit logs sendo criados ao criar pacientes