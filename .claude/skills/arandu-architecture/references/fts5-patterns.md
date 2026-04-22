# SQLite FTS5 — Busca Full-Text no Arandu

## Configuração das tabelas FTS5

```sql
-- observations_fts espelha a tabela observations
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts
USING fts5(
    content,                        -- campo indexado
    content='observations',         -- tabela de origem
    content_rowid='id'              -- chave primária
);

-- Triggers para manter índice sincronizado
CREATE TRIGGER observations_ai AFTER INSERT ON observations BEGIN
    INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;

CREATE TRIGGER observations_ad AFTER DELETE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content)
    VALUES ('delete', old.id, old.content);
END;

CREATE TRIGGER observations_au AFTER UPDATE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content)
    VALUES ('delete', old.id, old.content);
    INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;
```

## Queries FTS5

```go
// Busca simples
func (r *observationRepo) Search(ctx context.Context, query string) ([]*domain.Observation, error) {
    db, _ := tenant.TenantDB(ctx)

    rows, err := db.QueryContext(ctx, `
        SELECT o.id, o.session_id, o.content, o.created_at
        FROM observations o
        JOIN observations_fts fts ON o.id = fts.rowid
        WHERE observations_fts MATCH ?
        ORDER BY rank
        LIMIT 50
    `, query)
    // ...
}

// Busca com highlight (marca o termo encontrado)
func (r *observationRepo) SearchWithHighlight(ctx context.Context, query string) ([]SearchResult, error) {
    db, _ := tenant.TenantDB(ctx)

    rows, err := db.QueryContext(ctx, `
        SELECT
            o.id,
            o.session_id,
            highlight(observations_fts, 0, '<mark>', '</mark>') as excerpt,
            rank
        FROM observations o
        JOIN observations_fts fts ON o.id = fts.rowid
        WHERE observations_fts MATCH ?
        ORDER BY rank
        LIMIT 20
    `, query)
    // ...
}
```

## HTMX search com delay

```templ
// Busca ativada 500ms após parar de digitar — evita IO excessivo
templ SearchInput(placeholder string) {
    <input
        type="search"
        name="q"
        placeholder={ placeholder }
        hx-get="/observations/search"
        hx-trigger="input changed delay:500ms, search"
        hx-target="#search-results"
        hx-swap="innerHTML"
        hx-indicator="#search-spinner"
        class="input-silent w-full"
        autocomplete="off"
    />
    <span id="search-spinner" class="htmx-indicator">
        @Spinner("sm")
    </span>
}
```

## Sanitização de query FTS5

FTS5 tem sua própria sintaxe — sanitize antes de executar:

```go
func sanitizeFTSQuery(raw string) string {
    // Remove caracteres especiais do FTS5
    replacer := strings.NewReplacer(
        `"`, ``,
        `*`, ``,
        `^`, ``,
        `(`, ``,
        `)`, ``,
        `-`, ` `,
    )
    query := strings.TrimSpace(replacer.Replace(raw))
    if query == "" {
        return ""
    }
    // Adiciona * para prefix search
    return query + "*"
}
```
