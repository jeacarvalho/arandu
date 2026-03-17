REQ-07-03-03 — Gestão de Migrações Multi-tenant

Identificação

ID: REQ-07-03-03

Capability: CAP-07-03 Gestão de Acesso e Multi-tenancy

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do utilizador

Como administrador do sistema, quero que todas as bases de dados clínicas individuais sejam actualizadas automaticamente quando houver mudanças no esquema do sistema, para que todas as funcionalidades novas funcionem correctamente para todos os terapeutas sem intervenção manual.

Contexto

No modelo de Database-per-Tenant, o Arandu lida com centenas ou milhares de ficheiros SQLite independentes. Quando uma nova funcionalidade exige uma nova tabela (ex: interventions), o sistema deve garantir que todos esses ficheiros recebam a migração. Este requisito define o mecanismo de "Broadcast" de esquema.

Descrição funcional

O sistema deve implementar um motor de migração capaz de operar em lote (batch) e sob demanda.

Migração no Startup: Ao iniciar, o servidor deve ser capaz de percorrer todos os tenants registados no Banco Central e aplicar as migrações pendentes.

Migração Just-in-Time: Como medida de segurança extra, ao abrir uma conexão (REQ-07-03-02), o sistema deve verificar e aplicar migrações específicas para aquele tenant.

Versionamento Individual: Cada base de dados clínica deve manter a sua própria tabela schema_migrations para controlo local.

Atomicidade: Cada migração deve ser executada dentro de uma transação SQLite para evitar estados inconsistentes em caso de erro.

Lógica Técnica (SOTA)

1. Descoberta de Tenants

O sistema deve consultar o Banco Central para obter a lista de caminhos de ficheiros .db activos.

2. O Loop de Migração (Batch)

func RunGlobalMigrations(ctx context.Context) error {
    tenants, _ := centralDB.GetAllTenants(ctx)
    for _, t := range tenants {
        db, _ := sql.Open("sqlite3", t.Path)
        migrator := sqlite.NewMigrator(db)
        if err := migrator.Migrate(ctx); err != nil {
            log.Printf("Falha ao migrar tenant %s: %v", t.ID, err)
            // Estratégia: Logar erro mas continuar para os outros
        }
        db.Close()
    }
    return nil
}


3. Utilização do go:embed

O Migrator (já definido no projeto) deve utilizar os ficheiros .up.sql embutidos no binário para garantir que todos os bancos recebam exactamente a mesma versão do código SQL.

Interface (Padrão Arandu SOTA)

Dashboard de Admin (Interno): Uma visão simplificada que mostra o estado das migrações (ex: "150/150 bancos actualizados").

Feedback de Erro: Se um banco falhar a migração, o utilizador correspondente deve ser impedido de aceder e ver uma mensagem de "Manutenção Programada" para proteger a integridade dos dados.

Fluxo

O Administrador faz o deploy de uma nova versão do binário Arandu (com novos ficheiros .up.sql).

O sistema inicia e executa o RunGlobalMigrations.

Para cada ficheiro .db encontrado:

Abre a conexão.

Compara a tabela schema_migrations local com os ficheiros disponíveis.

Executa os scripts SQL em falta.

Regista o sucesso e fecha a conexão.

O sistema entra em modo operacional.

Critérios de Aceitação

CA-01: O sistema deve aplicar migrações com sucesso em múltiplos ficheiros SQLite sem corrupção de dados.

CA-02: Se uma migração falhar num banco específico, o erro não deve impedir a migração dos outros bancos.

CA-03: O sistema deve suportar migrações que adicionam colunas, tabelas ou índices.

CA-04: O tempo de migração global no startup deve ser optimizado para não causar timeouts de infraestrutura (execução assíncrona se necessário).

CA-05: O Banco Central (Control Plane) também deve possuir o seu próprio fluxo de migração isolado.

Persistência

Tabela de Controlo (em cada DB): schema_migrations
Localização: internal/infrastructure/repository/sqlite/migrations/*.sql

Fora do escopo

Rollback automático de migrações (requer backups complexos pré-migração).

Migrações de dados complexas (ETL) que exijam lógica externa ao SQL puro.

Migrações via UI para o utilizador final.