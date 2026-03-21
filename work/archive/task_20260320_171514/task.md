Task: Resiliência, Rollback e Hardening (Finalização)

ID da Tarefa: task_20260320_tenant_resilience

Requirement: REQ-07-03-05

Dependência: task_20260320_tenant_provisioning (Concluída)

Stack: Go, SQLite, os, chmod.

Status: CONCLUÍDA

## Resumo das Alterações

### 1. Atomicidade de Provisão (Rollback Manual)
- Implementado mecanismo de cleanup com flag `provisioned`
- Se o registo no arandu_central.db falhar, o ficheiro .db é removido automaticamente
- Teste `TestTenantService_ProvisionNewTenant_RollbackOnCentralFailure` verifica o rollback

### 2. Hardening do File System (Permissões 0600)
- Pasta `storage/tenants/` criada com permissões 0700 (em vez de 0755)
- Ficheiros .db criados com permissões 0600 após criação
- Teste `TestTenantService_ProvisionNewTenant_Permissions` verifica as permissões

### 3. Log de Auditoria de Infraestrutura
- Adicionado `log.Printf("[Provisioning] New Tenant Created: {tenant_id: %s, user_id: %s}", tenantID, userID)`
- Log apenas contém IDs, sem informações sensíveis

### 4. Guard Script Atualizado
- Guard script agora segue redirects com `-L` flag
- Permite verificar rotas protegidas que redirecionam para login

## Ficheiros Modificados

- `internal/application/services/tenant_service.go` - Rollback, permissions, audit logging
- `internal/application/services/tenant_service_test.go` - Novos testes
- `scripts/arandu_guard.sh` - Suporte a redirects

## Testes

✅ `TestTenantService_ProvisionNewTenant` - Provisionamento normal
✅ `TestTenantService_ProvisionNewTenant_RollbackOnCentralFailure` - Rollback automático
✅ `TestTenantService_ProvisionNewTenant_Permissions` - Verificação de permissões 0600/0700
✅ `arandu_guard.sh` - Sistema operacional

## Checklist de Integridade

[x] O processo de limpeza (rollback) foi testado e funciona?
[x] As permissões 0600 estão a ser aplicadas no momento da criação do ficheiro?
[x] O log de auditoria não contém informações sensíveis (como senhas), apenas IDs?
[x] O scripts/arandu_guard.sh confirma que o fluxo de login continua operacional