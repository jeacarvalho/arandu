# 💾 SQLite Best Practices - Arandu

**Última atualização:** $(date +"%d de %B de %Y")
**Foco:** Banco de dados SQLite, FTS5, Migrations, Performance

> 📚 **Consulte também:** [MASTER_LEARNINGS.md](./MASTER_LEARNINGS.md)

---

## 📋 Índice

1. [Driver e Configuração](#driver-e-configuração)
2. [Migrations](#migrations)
3. [Full-Text Search (FTS5)](#full-text-search-fts5)
4. [Performance e Otimização](#performance-e-otimização)
5. [Transações e Concorrência](#transações-e-concorrência)
6. [Backup e Recuperação](#backup-e-recuperação)
7. [Monitoramento](#monitoramento)
8. [Referências](#referências)

---

## Driver e Configuração

### Driver Moderno (sem CGO)

**Problema:** `database/sql` com driver CGO tem dependências de compilação.

**Solução:** Usar `modernc.org/sqlite` (pure Go):

```go
import (
    "database/sql"
    _ "modernc.org/sqlite" // Driver pure Go
)

func setupDatabase() (*sql.DB, error) {
    // Conexão sem CGO
    db, err := sql.Open("sqlite", "./arandu.db")
    if err != nil {
        return nil, err
    }
    
    // Configurações recomendadas
    db.SetMaxOpenConns(1)          // SQLite é single-writer
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0)       // Conexões persistentes
    
    // Ativar WAL mode para melhor concorrência
    _, err = db.Exec("PRAGMA journal_mode = WAL;")
    if err != nil {
        return nil, err
    }
    
    // Ativar foreign keys
    _, err = db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {
        return nil, err
    }
    
    // Synchronous NORMAL para performance (trade-off segurança)
    _, err = db.Exec("PRAGMA synchronous = NORMAL;")
    if err != nil {
        return nil, err
    }
    
    return db, nil
}
```

### Configurações PRAGMA Recomendadas

```sql
-- Modo WAL para melhor concorrência leitura/escrita
PRAGMA journal_mode = WAL;

-- Foreign keys ativadas
PRAGMA foreign_keys = ON;

-- Synchronous NORMAL (balance entre performance e segurança)
PRAGMA synchronous = NORMAL;

-- Cache size (em páginas, 2000 = ~32MB)
PRAGMA cache_size = 2000;

-- Temp store em memória
PRAGMA temp_store = MEMORY;

-- Busy timeout (30 segundos)
PRAGMA busy_timeout = 30000;
```

---

## Migrations

### Estrutura de Migrations

```
internal/infrastructure/repository/sqlite/migrations/
├── 0001_initial_schema.up.sql
├── 0001_initial_schema.down.sql
├── 0002_patients_table.up.sql
├── 0002_patients_table.down.sql
├── 0003_sessions_table.up.sql
├── 0003_sessions_table.down.sql
├── 0004_observations_interventions.up.sql
├── 0004_observations_interventions.down.sql
├── 0005_biopsychosocial_tables.up.sql
└── 0005_biopsychosocial_tables.down.sql
```

### Padrão de Numeração

- **4 dígitos:** `0001`, `0002`, etc.
- **Descrição:** `_patients_table`, `_fts5_search`, etc.
- **Extensões:** `.up.sql` (aplicar), `.down.sql` (reverter)

### Exemplo de Migration Completa

```sql
-- 0005_biopsychosocial_tables.up.sql

-- Tabela de medicamentos
CREATE TABLE IF NOT EXISTS medications (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT NOT NULL,
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Tabela de sinais vitais
CREATE TABLE IF NOT EXISTS vitals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    recorded_at DATETIME NOT NULL,
    heart_rate INTEGER,
    blood_pressure_systolic INTEGER,
    blood_pressure_diastolic INTEGER,
    temperature REAL,
    respiratory_rate INTEGER,
    oxygen_saturation INTEGER,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_medications_patient_id ON medications(patient_id);
CREATE INDEX IF NOT EXISTS idx_medications_started_at ON medications(started_at);
CREATE INDEX IF NOT EXISTS idx_vitals_patient_id ON vitals(patient_id);
CREATE INDEX IF NOT EXISTS idx_vitals_recorded_at ON vitals(recorded_at);

-- Trigger para updated_at automático
CREATE TRIGGER IF NOT EXISTS update_medications_timestamp 
AFTER UPDATE ON medications
BEGIN
    UPDATE medications SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_vitals_timestamp 
AFTER UPDATE ON vitals
BEGIN
    UPDATE vitals SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

```sql
-- 0005_biopsychosocial_tables.down.sql

-- Reverter na ordem inversa (dependências)
DROP TRIGGER IF EXISTS update_vitals_timestamp;
DROP TRIGGER IF EXISTS update_medications_timestamp;

DROP INDEX IF EXISTS idx_vitals_recorded_at;
DROP INDEX IF EXISTS idx_vitals_patient_id;
DROP INDEX IF EXISTS idx_medications_started_at;
DROP INDEX IF EXISTS idx_medications_patient_id;

DROP TABLE IF EXISTS vitals;
DROP TABLE IF EXISTS medications;
```

### Boas Práticas para Migrations

1. **Sempre incluir `.down.sql`** - Para rollback
2. **Usar `IF NOT EXISTS`/`IF EXISTS`** - Idempotência
3. **Manter ordem de dependências** - Criar tabelas antes de índices/triggers
4. **Testar rollback** - Executar `.down.sql` após `.up.sql`
5. **Documentar mudanças** - Comentários no arquivo

### Sistema de Migração Automática

```go
// internal/infrastructure/repository/sqlite/migrator.go
func RunMigrations(db *sql.DB) error {
    migrations := []string{
        "0001_initial_schema.up.sql",
        "0002_patients_table.up.sql",
        // ...
    }
    
    for _, migration := range migrations {
        content, err := readMigrationFile(migration)
        if err != nil {
            return fmt.Errorf("failed to read migration %s: %w", migration, err)
        }
        
        // Executar em transação
        tx, err := db.Begin()
        if err != nil {
            return fmt.Errorf("failed to begin transaction for %s: %w", migration, err)
        }
        
        _, err = tx.Exec(content)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to execute migration %s: %w", migration, err)
        }
        
        if err := tx.Commit(); err != nil {
            return fmt.Errorf("failed to commit migration %s: %w", migration, err)
        }
        
        log.Printf("✅ Migration applied: %s", migration)
    }
    
    return nil
}
```

---

## Full-Text Search (FTS5)

### Por que FTS5?

- **Busca textual avançada** - Stemming, ranking, highlighting
- **Performance** - Índices otimizados para texto
- **Integração nativa** - Parte do SQLite

### Configuração FTS5

```sql
-- Tabela virtual FTS5 para observações clínicas
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
    content,                    -- Campo a indexar
    content='observations',     -- Tabela externa
    tokenize='porter'           -- Stemming em inglês
);

-- Tabela virtual FTS5 para pacientes (múltiplos campos)
CREATE VIRTUAL TABLE IF NOT EXISTS patients_fts USING fts5(
    name,                       -- Nome do paciente
    notes,                      -- Observações
    content='patients',         -- Tabela externa
    tokenize='porter'
);
```

### Triggers para Sincronização Automática

```sql
-- Trigger para INSERT
CREATE TRIGGER IF NOT EXISTS observations_ai AFTER INSERT ON observations BEGIN
    INSERT INTO observations_fts(rowid, content) VALUES (NEW.id, NEW.content);
END;

-- Trigger para DELETE
CREATE TRIGGER IF NOT EXISTS observations_ad AFTER DELETE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', OLD.id, OLD.content);
END;

-- Trigger para UPDATE
CREATE TRIGGER IF NOT EXISTS observations_au AFTER UPDATE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', OLD.id, OLD.content);
    INSERT INTO observations_fts(rowid, content) VALUES (NEW.id, NEW.content);
END;
```

### Queries FTS5 com Highlighting

```sql
-- Busca simples
SELECT * FROM observations_fts 
WHERE observations_fts MATCH 'ansiedade' 
ORDER BY rank;

-- Busca com highlighting (tags <b>)
SELECT 
    snippet(observations_fts, 0, '<b>', '</b>', '...', 10) as highlighted,
    observations.*
FROM observations_fts
JOIN observations ON observations_fts.rowid = observations.id
WHERE observations_fts MATCH 'ansiedade'
ORDER BY rank
LIMIT 20;
```

### Busca em Múltiplos Campos

```sql
-- Busca em pacientes (nome OU observações)
SELECT 
    patients.*,
    snippet(patients_fts, 0, '<b>', '</b>', '...', 10) as highlighted_name,
    snippet(patients_fts, 1, '<b>', '</b>', '...', 10) as highlighted_notes
FROM patients_fts
JOIN patients ON patients_fts.rowid = patients.id
WHERE patients_fts MATCH 'joão OR terapia'
ORDER BY rank
LIMIT 15;
```

### Tratamento no Go (RawHTML)

```go
type SearchResult struct {
    PatientID   string
    Name        string
    Highlighted templ.HTML // RawHTML para não escapar <b>
}

func (r *Repository) SearchPatients(ctx context.Context, query string) ([]SearchResult, error) {
    rows, err := r.db.QueryContext(ctx, `
        SELECT 
            patients.id,
            patients.name,
            snippet(patients_fts, 0, '<b>', '</b>', '...', 10) as highlighted
        FROM patients_fts
        JOIN patients ON patients_fts.rowid = patients.id
        WHERE patients_fts MATCH ?
        ORDER BY rank
        LIMIT 15
    `, query)
    
    // ...
    
    for rows.Next() {
        var result SearchResult
        var highlighted string
        if err := rows.Scan(&result.PatientID, &result.Name, &highlighted); err != nil {
            return nil, err
        }
        result.Highlighted = templ.HTML(highlighted) // Converter para RawHTML
        results = append(results, result)
    }
    
    return results, nil
}
```

---

## Performance e Otimização

### Índices Estratégicos

```sql
-- Índices para queries comuns
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(name);
CREATE INDEX IF NOT EXISTS idx_sessions_patient_id_date ON sessions(patient_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_observations_session_id ON observations(session_id);
CREATE INDEX IF NOT EXISTS idx_medications_patient_id_started_at ON medications(patient_id, started_at DESC);

-- Índices compostos para queries de timeline
CREATE INDEX IF NOT EXISTS idx_timeline_events_patient_date ON timeline_events(patient_id, event_date DESC);
```

### Análise de Queries (EXPLAIN QUERY PLAN)

```sql
-- Analisar performance de query
EXPLAIN QUERY PLAN
SELECT * FROM sessions 
WHERE patient_id = '123' 
ORDER BY date DESC 
LIMIT 20;

-- Resultado esperado:
-- SEARCH TABLE sessions USING INDEX idx_sessions_patient_id_date (patient_id=?)
```

### Otimização de Queries

**RUIM:** N+1 queries
```go
// RUIM: Query para cada paciente
for _, patient := range patients {
    sessions, _ := repo.GetSessionsByPatientID(patient.ID)
    // ...
}
```

**BOM:** JOIN ou batch query
```go
// BOM: Uma query com JOIN
patientsWithSessions, _ := repo.GetPatientsWithRecentSessions(limit)
```

**BOM:** IN clause para múltiplos IDs
```sql
-- Buscar sessões para múltiplos pacientes
SELECT * FROM sessions 
WHERE patient_id IN (?, ?, ?)
ORDER BY date DESC;
```

### Paginação Eficiente

```sql
-- Paginação com LIMIT/OFFSET (para datasets moderados)
SELECT * FROM observations 
WHERE patient_id = ? 
ORDER BY created_at DESC 
LIMIT 20 OFFSET 40;

-- Paginação com cursor (para datasets grandes)
SELECT * FROM observations 
WHERE patient_id = ? AND created_at < ? 
ORDER BY created_at DESC 
LIMIT 20;
```

### Vacuum e Manutenção

```sql
-- Vacuum periódico (reduz tamanho do arquivo)
VACUUM;

-- Analisar índices
ANALYZE;

-- Otimizar database (SQLite 3.38.0+)
PRAGMA optimize;
```

---

## Transações e Concorrência

### Padrão de Transações

```go
func (r *Repository) CreateWithTransaction(ctx context.Context, patient Patient) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    // Sempre fazer defer rollback (será ignorado se commit successful)
    defer tx.Rollback()
    
    // Executar operações
    _, err = tx.ExecContext(ctx, `
        INSERT INTO patients (id, name, notes) 
        VALUES (?, ?, ?)
    `, patient.ID, patient.Name, patient.Notes)
    if err != nil {
        return err
    }
    
    // Criar entrada na FTS5 (via trigger automático)
    // Trigger será executado dentro da mesma transação
    
    // Commit se tudo ok
    if err := tx.Commit(); err != nil {
        return err
    }
    
    return nil
}
```

### WAL Mode para Concorrência

**Vantagens do WAL mode:**
- Leituras não bloqueiam escritas
- Escritas não bloqueiam leituras
- Melhor performance em multi-thread

**Ativação:**
```sql
PRAGMA journal_mode = WAL;
```

**Checkpoint periódico (opcional):**
```sql
PRAGMA wal_checkpoint(TRUNCATE);
```

### Locking e Timeouts

```go
// Configurar busy timeout
db.SetConnMaxLifetime(0)
db.SetMaxOpenConns(1)  // SQLite é single-writer
db.SetMaxIdleConns(1)

// No SQL
PRAGMA busy_timeout = 30000;  // 30 segundos
```

---

## Backup e Recuperação

### Backup Online (WAL mode)

```bash
# Backup enquanto database está em uso
sqlite3 arandu.db ".backup backup.db"

# Backup com progresso
sqlite3 arandu.db ".backup ? backup.db"
```

### Backup Programático

```go
func BackupDatabase(sourcePath, backupPath string) error {
    sourceDB, err := sql.Open("sqlite", sourcePath)
    if err != nil {
        return err
    }
    defer sourceDB.Close()
    
    backupDB, err := sql.Open("sqlite", backupPath)
    if err != nil {
        return err
    }
    defer backupDB.Close()
    
    // Usar backup API do SQLite
    _, err = sourceDB.Exec(`VACUUM INTO ?`, backupPath)
    return err
}
```

### Restauração

```bash
# Restaurar backup
sqlite3 arandu.db ".restore backup.db"

# Ou copiar arquivo diretamente (se database não está em uso)
cp backup.db arandu.db
```

### Verificação de Integridade

```sql
-- Verificar integridade do database
PRAGMA integrity_check;

-- Verificar foreign keys
PRAGMA foreign_key_check;

-- Verificar quick (mais rápido)
PRAGMA quick_check;
```

---

## Monitoramento

### Estatísticas do Database

```sql
-- Tamanho do database
SELECT page_count * page_size as size_bytes FROM pragma_page_count(), pragma_page_size();

-- Contagem de registros
SELECT 
    (SELECT COUNT(*) FROM patients) as patient_count,
    (SELECT COUNT(*) FROM sessions) as session_count,
    (SELECT COUNT(*) FROM observations) as observation_count;

-- Espaço usado por tabelas
SELECT 
    name, 
    SUM(pgsize) as size_bytes
FROM dbstat 
GROUP BY name 
ORDER BY size_bytes DESC;
```

### Performance Queries

```go
// Middleware para log de queries lentas
func QueryLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        
        if duration > 100*time.Millisecond {
            log.Printf("⚠️ Slow request: %s %s took %v", r.Method, r.URL.Path, duration)
        }
    })
}
```

### Health Check

```go
func HealthCheck(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Verificar conexão com database
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()
        
        var result int
        err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
        
        if err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "unhealthy",
                "error":  err.Error(),
            })
            return
        }
        
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "healthy",
            "database": "connected",
        })
    }
}
```

---

## Referências

### Arquivos de Exemplo

1. **`internal/infrastructure/repository/sqlite/migrations/`** - Todas as migrations
2. **`internal/infrastructure/repository/sqlite/patient_repository.go`** - Repository com queries
3. **`work/learnings/task_20260317_220659.md`** - Infinite scroll e performance

### Documentação SQLite

- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [FTS5 Documentation](https://www.sqlite.org/fts5.html)
- [PRAGMA Statements](https://www.sqlite.org/pragma.html)

### Ferramentas

- **`sqlite3` CLI** - Interface de linha de comando
- **DB Browser for SQLite** - GUI para análise
- **`EXPLAIN QUERY PLAN`** - Análise de performance

### Troubleshooting

1. **Database locked** - Verificar WAL mode e busy_timeout
2. **Slow queries** - Usar EXPLAIN QUERY PLAN, adicionar índices
3. **FTS5 not working** - Verificar triggers de sincronização
4. **Large file size** - Executar VACUUM periodicamente

---

*Baseado em aprendizados de implementações reais no Arandu, incluindo FTS5, migrations e otimizações de performance*