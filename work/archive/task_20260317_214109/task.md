Guia de Implementação Master: SQLite FTS5 no Arandu

Contexto: Após a confirmação de compatibilidade (Teste de Sanidade OK), este prompt guia a implementação definitiva da busca textual de alta performance para o banco de dados clínico, utilizando o padrão de "Conteúdo Externo" para máxima eficiência.

⚡ Passo 0: Ajuste do Driver (Obrigatório)

Para evitar falhas de compilação C e garantir FTS5 nativo:

Instale: go get modernc.org/sqlite

Importe: No arquivo de conexão, substitua o driver anterior por import _ "modernc.org/sqlite".

Conexão: Use o driver name "sqlite" em sql.Open("sqlite", ...).

📂 Passo 1: Migração SQL de "Conteúdo Externo"

Não duplicaremos os dados brutos. O FTS5 indexará as tabelas observations e interventions existentes.

Ficheiro: internal/infrastructure/repository/sqlite/migrations/0002_enable_fts5.up.sql

-- 1. Tabelas Virtuais FTS5 (External Content)
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
    content,
    content='observations',
    content_rowid='id'
);

CREATE VIRTUAL TABLE IF NOT EXISTS interventions_fts USING fts5(
    content,
    content='interventions',
    content_rowid='id'
);

-- 2. Triggers de Sincronização: OBSERVATIONS
-- Sincroniza INSERT
CREATE TRIGGER IF NOT EXISTS observations_ai AFTER INSERT ON observations BEGIN
  INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;

-- Sincroniza DELETE
CREATE TRIGGER IF NOT EXISTS observations_ad AFTER DELETE ON observations BEGIN
  INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', old.id, old.content);
END;

-- Sincroniza UPDATE
CREATE TRIGGER IF NOT EXISTS observations_au AFTER UPDATE ON observations BEGIN
  INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', old.id, old.content);
  INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;

-- 3. Triggers de Sincronização: INTERVENTIONS
CREATE TRIGGER IF NOT EXISTS interventions_ai AFTER INSERT ON interventions BEGIN
  INSERT INTO interventions_fts(rowid, content) VALUES (new.id, new.content);
END;

CREATE TRIGGER IF NOT EXISTS interventions_ad AFTER DELETE ON interventions BEGIN
  INSERT INTO interventions_fts(interventions_fts, rowid, content) VALUES('delete', old.id, old.content);
END;

CREATE TRIGGER IF NOT EXISTS interventions_au AFTER UPDATE ON interventions BEGIN
  INSERT INTO interventions_fts(interventions_fts, rowid, content) VALUES('delete', old.id, old.content);
  INSERT INTO interventions_fts(rowid, content) VALUES (new.id, new.content);
END;

-- 4. Reconstrução do Índice Inicial (Para dados já existentes)
INSERT INTO observations_fts(observations_fts) VALUES('rebuild');
INSERT INTO interventions_fts(interventions_fts) VALUES('rebuild');


🔍 Passo 2: Lógica de Busca e Análise (Go)

Busca Global (MATCH)

Para localizar termos em milissegundos:

SELECT session_id, content 
FROM observations 
WHERE id IN (
    SELECT rowid FROM observations_fts WHERE content MATCH ?
);


Análise de Frequência (Vocabulary)

Para extrair os temas mais citados sem percorrer o texto bruto:

SELECT term, count 
FROM fts5vocabulary('observations_fts', 'col')
WHERE term NOT IN ('de', 'a', 'o', 'que', 'e', 'do', 'da', 'em', 'um', 'para', 'com') -- Stopwords
ORDER BY count DESC 
LIMIT 15;


🛡️ Checklist de Validação para o Agente

[ ] A migração foi aplicada com sucesso via Migrator?

[ ] O driver utilizado é o modernc.org/sqlite?

[ ] Ao inserir uma observação, ela aparece na busca MATCH instantaneamente?

[ ] O comando sqlite3 arandu.db "SELECT * FROM observations_fts LIMIT 1;" retorna dados?

Atenção: Se houver erro de "database is locked", certifique-se de que o SQLite está configurado com _pragma=foreign_keys(1)&_pragma=journal_mode(WAL).