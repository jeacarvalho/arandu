Task: Refatoração de Armazenamento (Directory Hashing)

ID da Tarefa: task_20260321_storage_hashing

Requirement: REQ-07-03-02 (Evolução de Infraestrutura)

Stack: Go, os, path/filepath.

🎯 Objetivo

Implementar a estratégia de Directory Hashing para a gestão de arquivos de banco de dados. Isso evita o gargalo de performance do Linux quando milhares de arquivos residem na mesma pasta e prepara o Arandu para escala maciça.

🛠️ Escopo Técnico

1. Camada de Infraestrutura (internal/platform/storage/path_resolver.go)

Implementar um serviço de resolução de caminhos:

Função ResolveTenantPath(tenantID string) string:

Recebe: a1b2c3d4...

Retorna: storage/tenants/a1/b2/clinical_a1b2c3d4.db

Função EnsureTenantDir(tenantID string) error:

Garante a criação recursiva das pastas (os.MkdirAll) com permissões 0700.

2. Atualização do Connection Manager / Tenant Pool

Substituir a construção manual de strings de caminho pela chamada ao PathResolver.

Garantir que, ao criar um novo tenant (Task 5), o diretório hashed seja criado antes da inicialização do SQLite.

3. Script de Migração de Dados (Legado)

Criar um pequeno script/utilitário em Go para mover bancos de dados existentes na raiz de storage/tenants/ para a nova estrutura de pastas.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Resolução

Passar um UUID fixo e verificar se o caminho retornado segue o padrão aa/bb/filename.db.

B. Teste de Criação Física

Simular o onboarding de um novo médico.

Verificar: Se as subpastas foram criadas corretamente no disco.

Verificar: Se o SQLite consegue abrir e escrever no arquivo dentro da nova estrutura.

🛡️ Checklist de Integridade

[x] O sistema lida corretamente com a transição (encontra arquivos antigos ou move-os)?

[x] As permissões de diretório continuam restritivas (0700)?

[x] O scripts/arandu_guard.sh confirma que o login e acesso ao dashboard continuam funcionando?

---

## Status: ✅ COMPLETADO

**Data:** 2026-03-22

**Resumo das alterações:**

1. ✅ Criado `internal/platform/storage/path_resolver.go` com:
   - `ResolveTenantPath(tenantID)` - retorna caminho hashed (ex: `storage/tenants/fe/9b/clinical_xxx.db`)
   - `EnsureTenantDir(tenantID)` - cria diretórios com permissões 0700
   - `TenantDirExists()`, `TenantDBExists()` - verificações

2. ✅ Atualizado `internal/infrastructure/repository/sqlite/tenant_pool.go`:
   - Usa PathResolver para construir caminhos
   - Cria estrutura hashed automaticamente

3. ✅ Atualizado `internal/application/services/tenant_service.go`:
   - Usa PathResolver para ProvisionNewTenant
   - Atualizado ValidateTenantDB e GetTenantDBPath

4. ✅ Criado script de migração `scripts/arandu_migrate_to_hashed_storage.sh`:
   - Migra bancos legados para estrutura hashed
   - Atualiza db_path no banco central

5. ✅ Migração executada:
   - Todos os bancos existentes movidos para `aa/bb/` estrutura

6. ✅ Testes E2E atualizados para novo formato de caminhos