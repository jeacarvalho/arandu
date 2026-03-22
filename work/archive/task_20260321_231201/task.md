Task: Padronização de Identidade com UUID v4

ID da Tarefa: task_20260321_uuid_standardization

Requirement: REQ-07-03-01 e REQ-07-03-05

Dependência: Nenhuma (é um pré-requisito para o Storage Hashing)

Stack: Go, github.com/google/uuid.

🎯 Objetivo

Substituir a geração de IDs baseada em strings/números aleatórios por UUIDs v4 (RFC 4122) reais em todo o Control Plane. Isso garantirá unicidade global, segurança por obscuridade e compatibilidade com a estratégia de hashing de diretórios.

🛠️ Escopo Técnico

1. Adição de Dependência

Executar go get github.com/google/uuid.

2. Refatoração da Camada de Domínio e Aplicação

Identidade de Usuário: No momento da criação do usuário (Login/OAuth), o ID gerado deve ser uuid.New().String().

Identidade de Tenant: O tenant_id deve seguir o mesmo padrão UUID.

Remover Prefixos: Eliminar os prefixos user- e tenant- armazenados no banco. O tipo de dado no SQLite deve continuar sendo TEXT, mas o conteúdo deve ser apenas o UUID (ex: 550e8400-e29b-41d4-a716-446655440000).

3. Impacto no Sistema de Ficheiros

Os ficheiros em storage/tenants/ devem ser renomeados de clinical_tenant-xxx.db para clinical_{uuid}.db.

O DashboardService e o TenantPool devem ser atualizados para esperar este novo formato.

🧪 Protocolo de Testes "Ironclad"

A. Validação de Formato

Criar um novo utilizador via Google ou Registro.

Abrir o arandu_central.db via terminal.

Verificar: O ID do utilizador e o TenantID devem ter o formato xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx.

B. Teste de Colisão e Hashing

Gerar 10 utilizadores em sequência.

Verificar: Se os IDs são suficientemente aleatórios (essencial para que o Directory Hashing distribua bem os ficheiros entre as pastas aa/, bb/, etc).

🛡️ Checklist de Integridade

[x] O comando go get foi executado e o go.mod atualizado?

[x] Não existem mais prefixos "user-" ou "tenant-" hardcoded no código?

[x] O sistema lida com o "erro" de não encontrar bases de dados antigas (ou existe um script de renomeação)?

[x] O scripts/arandu_guard.sh confirma que o login continua funcional?

---

## Status: ✅ COMPLETADO

**Data:** 2026-03-21

**Resumo das alterações:**
1. ✅ Adicionado `github.com/google/uuid` como dependência
2. ✅ Atualizado `auth_handler.go` para usar `uuid.New().String()` em vez de `fmt.Sprintf("tenant-%d", time.Now().UnixNano())` e `fmt.Sprintf("user-%d", time.Now().UnixNano())`
3. ✅ Atualizado formato de session ID para UUID
4. ✅ Criado script de migração `scripts/arandu_migrate_to_uuid.sh`
5. ✅ Migrados todos os tenants, users e sessions existentes para UUID v4
6. ✅ Renomeados arquivos de banco de dados para o novo formato `clinical_{uuid}.db`
7. ✅ Build passou com sucesso