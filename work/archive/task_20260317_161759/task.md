Task: Implementação de Infraestrutura FTS5 para Análise Clínica

ID da Tarefa: task_20260317_infra_fts5
Requirement Relacionado: REQ-04-01-01
Stack Técnica: SQLite (FTS5), Go, Migrations SQL.

🎯 Objetivo

Configurar a capacidade de busca textual completa (Full Text Search) no SQLite para permitir análises de frequência de termos em milésimos de segundo. Esta é a fundação para a identificação de padrões e temas recorrentes.

🛠️ Escopo Técnico

1. Camada de Persistência (SQL Migrations)

Criar o ficheiro de migração internal/infrastructure/repository/sqlite/migrations/0002_enable_fts5.up.sql.

Conteúdo da Migração:

Criar tabelas virtuais FTS5: observations_fts e interventions_fts.

Configurar estas tabelas como EXTERNAL CONTENT apontando para as tabelas originais (observations e interventions) para economizar espaço em disco.

Implementar Triggers SQL para sincronização automática:

AFTER INSERT nas tabelas originais -> Insert no FTS.

AFTER DELETE nas tabelas originais -> Delete no FTS.

AFTER UPDATE nas tabelas originais -> Update no FTS.

2. Camada de Repositório (Go)

No internal/infrastructure/repository/sqlite/patient_repository.go (ou um novo analysis_repository.go), implementar o método GetTermFrequency(ctx, patientID, filter).

Lógica da Query:

Utilizar a funcionalidade fts5vocabulary do SQLite para extrair a lista de termos e frequências.

A query deve filtrar pelo patient_id associado às sessões (JOIN necessário entre a tabela virtual e a tabela de sessões).

Retornar uma lista de structs (Term string, Count int).

⚙️ Habilitação Técnica (Guia para o Agente/Desenvolvedor)

Se o sistema retornar erro de "unknown module: fts5", siga estes passos:

Verificação de Suporte:
No terminal do SQLite, execute: PRAGMA compile_options;.
Verifique se ENABLE_FTS5 está na lista. Se não estiver, o binário do SQLite do sistema precisa de atualização.

Go Build Tags:
O driver github.com/mattn/go-sqlite3 requer que o FTS5 seja explicitamente habilitado na compilação do Go.
Execute ou compile o projeto usando:

go run -tags "fts5" cmd/arandu/main.go
# ou para build
go build -tags "fts5" -o arandu cmd/arandu/main.go


CGO Required:
Certifique-se de que CGO_ENABLED=1 está definido no seu ambiente, pois o driver de SQLite padrão do Go é um wrapper para a biblioteca em C.

Alternativa CGO-Free:
Se a compilação C for um problema no ambiente, considere migrar para o driver modernc.org/sqlite, que é puramente Go e geralmente traz o FTS5 habilitado por padrão.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Sincronização (SQL)

Inserir uma nova observação na tabela observations.

Verificar se o termo aparece imediatamente na tabela observations_fts via consulta MATCH.

Apagar a observação e verificar se o índice FTS foi limpo.

B. Teste de Performance (Go)

Executar a função GetTermFrequency contra a base de dados de teste (63.000 sessões).

Critério de Aceitação: A resposta deve ser retornada em menos de 200ms.

🛡️ Checklist de Integridade

[ ] A migração SQL utiliza EXTERNAL CONTENT para evitar duplicação desnecessária de dados?

[ ] Os Triggers cobrem os três estados (INSERT, UPDATE, DELETE)?

[ ] O código Go lida corretamente com o fecho dos cursores SQL?

[ ] O projeto é executado com as build tags -tags "fts5"?

[ ] O scripts/arandu_guard.sh confirma que o banco inicializa com a nova migração?

Instruções de Conclusão

Ao finalizar, execute o script:
./scripts/arandu_conclude_task.sh 20260317_infra_fts5 --success