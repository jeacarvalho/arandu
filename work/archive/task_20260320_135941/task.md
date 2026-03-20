Task: Implementação da Infraestrutura do Control Plane

ID da Tarefa: task_20260320_control_plane_infra

Requirement: REQ-07-03-01 e REQ-07-03-05

Stack: Go, SQLite, go:embed.

🎯 Objetivo

pré-requisito: leia o arquivo work/plano_atual.md para entender o contexto da tarefa, parte de um escopo maior.

Configurar o "Cérebro" do Arandu: o Control Plane. Este banco de dados centralizado será responsável por gerir as credenciais dos utilizadores e o mapeamento para os seus respectivos ficheiros SQLite clínicos (Data Plane).

🛠️ Escopo Técnico

1. Camada de Infraestrutura (Central DB)

Novo Banco: storage/arandu_central.db.

Migrations Central: Criar diretório internal/infrastructure/repository/sqlite/migrations_central/.

Ficheiro: 0001_initial_central.up.sql:

CREATE TABLE IF NOT EXISTS tenants (
    id TEXT PRIMARY KEY,          -- UUID do Tenant
    db_path TEXT NOT NULL,        -- Caminho físico: storage/tenants/clinical_{id}.db
    status TEXT DEFAULT 'active', -- active, suspended
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,          -- UUID do Utilizador
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,           -- Nulo se for apenas OAuth
    tenant_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);



2. Camada de Domínio (internal/domain/shared)

Criar as structs User e Tenant que serão partilhadas entre módulos.

Implementar o método ValidatePassword(hash, password) utilizando a biblioteca bcrypt.

3. Orquestração de Inicialização (main.go)

Implementar a lógica para inicializar o Banco Central no startup do servidor.

Garantir que o diretório storage/tenants/ é criado automaticamente se não existir.

Configurar o Migrator para aplicar as migrações centrais separadamente das migrações clínicas (Data Plane).

🧪 Protocolo de Verificação SOTA

A. Teste de Inicialização

Apagar a pasta storage/ (em ambiente de desenvolvimento).

Iniciar o servidor.

Verificar: O ficheiro arandu_central.db foi criado e contém as tabelas users e tenants.

Verificar: O diretório storage/tenants/ foi criado.

B. Teste de Integridade de Schema

Validar via terminal sqlite3 que as chaves estrangeiras (FOREIGN KEY) entre users e tenants estão activas.

status: completed

## ✅ Implementação Concluída

### Arquivos Criados:
- `internal/infrastructure/repository/sqlite/migrations_central/0001_initial_central.up.sql` - Migration central
- `internal/infrastructure/repository/sqlite/migrations_central/0001_initial_central.down.sql` - Rollback
- `internal/domain/shared/tenant.go` - Struct Tenant com validação
- `internal/domain/shared/user.go` - Struct User com bcrypt
- `internal/infrastructure/repository/sqlite/central_db.go` - Gerenciador do banco central

### Alterações:
- `cmd/arandu/main.go` - Inicialização do Control Plane no startup
- `go.mod` - Adicionada dependência golang.org/x/crypto/bcrypt

### Verificações Realizadas:
✅ Banco central (arandu_central.db) criado com tabelas users e tenants
✅ Diretório storage/tenants/ criado automaticamente
✅ Chaves estrangeiras ativas entre users e tenants
✅ Pragma WAL configurado no banco central
✅ Script guard.sh executado com sucesso (rotas online, templ OK)

### Testes Automatizados Criados:
- `internal/domain/shared/tenant_test.go` - 4 testes (NewTenant, IsActive, Suspend, Activate)
- `internal/domain/shared/user_test.go` - 6 testes (NewUser, ValidatePassword, HashPassword, UpdateLastLogin, HasPassword)
- `internal/infrastructure/repository/sqlite/central_db_test.go` - 6 testes (NewCentralDB, Migrate, Idempotent, ForeignKeys, Isolation, WAL)

✅ Todos os testes passando (16 testes total)
