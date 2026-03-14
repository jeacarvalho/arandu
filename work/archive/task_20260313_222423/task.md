# TASK 20260313_222423

Requirement: req-01-00-01

Title: Criar migration SQLite para Patient

## Objetivo

Implementar um sistema de migrations para o SQLite, especificamente para a tabela `patients`, seguindo boas práticas de versionamento de banco de dados.

## Contexto

Atualmente, o schema da tabela `patients` é criado através do método `InitSchema()` no repositório. Esta abordagem tem limitações:

1. **Sem versionamento:** Não há histórico de mudanças no schema
2. **Sem rollback:** Não é possível reverter mudanças
3. **Sem migrações incrementais:** Apenas "create if not exists"
4. **Acoplamento:** Lógica de schema no repositório

Precisamos implementar um sistema de migrations que:
- Versiona mudanças no banco de dados
- Permite migrações para frente (up) e para trás (down)
- É independente do código da aplicação
- Pode ser executado em diferentes ambientes

## Análise do Estado Atual

**Arquivo atual:** `internal/infrastructure/repository/sqlite/schema.go`
```go
func InitSchema(db *sql.DB) error {
	// Cria todas as tabelas necessárias
	// ...
}
```

**Problemas identificados:**
1. Schema fixo, sem versionamento
2. Sem suporte a migrações incrementais
3. Dificuldade para evoluir o schema ao longo do tempo
4. Não lida com dados de migração (seed data)

## Tarefas Específicas

### 1. Definir Estrutura de Migrations
Criar diretório e estrutura para migrations:
```
internal/infrastructure/repository/sqlite/migrations/
├── 0001_create_patients_table.up.sql
├── 0001_create_patients_table.down.sql
└── migration.go (gerenciador)
```

### 2. Implementar Migration Manager
Criar um gerenciador de migrations que:
- Rastreia migrations aplicadas (tabela `schema_migrations`)
- Executa migrations na ordem correta
- Suporta rollback (down migrations)
- Valida integridade das migrations

### 3. Criar Migration para Tabela Patients
Criar os arquivos SQL para:
- **Up migration:** Criação da tabela `patients`
- **Down migration:** Remoção da tabela `patients`

### 4. Atualizar Inicialização do Banco
Modificar a inicialização do banco para:
1. Criar tabela `schema_migrations` se não existir
2. Executar migrations pendentes
3. Remover chamada a `InitSchema()` do repositório

### 5. Adicionar Seed Data (Opcional)
Considerar adicionar migration para dados iniciais:
- Pacientes de exemplo para desenvolvimento
- Dados de configuração do sistema

## Requisitos Técnicos

### Estrutura de Migration
Cada migration deve ter:
- **ID único:** Timestamp ou sequência numérica
- **Nome descritivo:** Ex: `create_patients_table`
- **Arquivo .up.sql:** Comandos SQL para aplicar a migration
- **Arquivo .down.sql:** Comandos SQL para reverter a migration

### Tabela schema_migrations
```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Migration Manager Interface
```go
type MigrationManager interface {
    Migrate() error
    Rollback(version string) error
    CurrentVersion() (string, error)
    PendingMigrations() ([]string, error)
}
```

## Restrições de Design

1. **Independência:** O sistema de migrations deve funcionar independentemente do repositório
2. **Simplicidade:** Manter o sistema simples, sem dependências complexas
3. **SQL puro:** Migrations em SQL, não em Go (mais portável)
4. **Transações:** Cada migration deve ser executada em transação
5. **Idempotência:** Executar múltiplas vezes não deve causar erros

## Critérios de Aceitação

✅ Sistema de migrations implementado com estrutura de diretórios  
✅ Migration manager com suporte a up/down  
✅ Tabela `schema_migrations` para rastreamento  
✅ Migration para tabela `patients` criada  
✅ Inicialização do banco atualizada para usar migrations  
✅ `InitSchema()` removido do repositório Patient  
✅ Testes para verificar aplicação de migrations  
✅ Rollback funcional para a migration  
✅ Projeto compila e testes passam  

## Passos de Implementação

1. Criar estrutura de diretórios para migrations
2. Implementar migration manager em Go
3. Criar migration SQL para tabela patients
4. Atualizar inicialização do banco de dados
5. Atualizar repositório Patient para não usar InitSchema()
6. Criar testes para o sistema de migrations
7. Validar fluxo completo: migrate → rollback → migrate
8. Documentar uso do sistema de migrations

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md`
- `internal/infrastructure/repository/sqlite/schema.go` (atual)
- `internal/infrastructure/repository/sqlite/patient_repository.go`
- `work/tasks/task_20260313_221358/implementation.md` (tarefa anterior do repository)
