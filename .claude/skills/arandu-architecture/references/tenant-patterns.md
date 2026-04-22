# Multi-tenancy — Padrões Avançados

## Pool de conexões por tenant

Abrir uma nova conexão SQLite a cada request é caro. Use um pool por tenant:

```go
// internal/infrastructure/tenant/pool.go

type TenantPool struct {
    mu    sync.RWMutex
    pools map[string]*sql.DB  // key: tenant db path
}

var globalPool = &TenantPool{
    pools: make(map[string]*sql.DB),
}

func (p *TenantPool) Get(dbPath string) (*sql.DB, error) {
    p.mu.RLock()
    db, ok := p.pools[dbPath]
    p.mu.RUnlock()
    if ok {
        return db, nil
    }

    // Primeira vez: abre e registra
    p.mu.Lock()
    defer p.mu.Unlock()

    // Double-check após lock
    if db, ok = p.pools[dbPath]; ok {
        return db, nil
    }

    db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_timeout=5000")
    if err != nil {
        return nil, fmt.Errorf("open tenant db %s: %w", dbPath, err)
    }
    db.SetMaxOpenConns(1)  // SQLite: uma conexão de escrita por vez
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0)

    p.pools[dbPath] = db
    return db, nil
}
```

## WAL mode — essencial para concorrência

SQLite no modo WAL (Write-Ahead Log) permite leituras concorrentes sem bloquear escritas:

```sql
-- Aplicar no startup de cada tenant DB
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;  -- Balanço entre segurança e performance
PRAGMA busy_timeout=5000;   -- Aguarda até 5s antes de SQLITE_BUSY
```

## Backup por tenant

Cada tenant tem seu próprio arquivo — backup é trivial:

```go
func BackupTenant(dbPath string, backupDir string) error {
    tenantID := filepath.Base(dbPath)
    backupPath := filepath.Join(backupDir, tenantID+".backup."+time.Now().Format("20060102"))

    // SQLite Online Backup API via VACUUM INTO
    db, _ := sql.Open("sqlite3", dbPath)
    defer db.Close()

    _, err := db.Exec("VACUUM INTO ?", backupPath)
    return err
}
```

## Criação de novo tenant

```go
func CreateTenantDB(controlDB *sql.DB, psychologistID uuid.UUID, dbsDir string) (string, error) {
    dbPath := filepath.Join(dbsDir, psychologistID.String()+".db")

    // 1. Cria o arquivo DB do tenant
    tenantDB, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return "", err
    }
    defer tenantDB.Close()

    // 2. Aplica todas as migrations
    if err := applyMigrations(tenantDB, migrationsFS); err != nil {
        os.Remove(dbPath)  // rollback: remove arquivo se migration falhou
        return "", err
    }

    // 3. Registra no Control Plane
    _, err = controlDB.Exec(
        `INSERT INTO tenants (psychologist_id, db_path, created_at) VALUES (?, ?, ?)`,
        psychologistID, dbPath, time.Now(),
    )
    if err != nil {
        os.Remove(dbPath)
        return "", err
    }

    return dbPath, nil
}
```

## Isolamento de queries — verificação de segurança

Nunca deve existir um JOIN entre o Control Plane e um tenant DB.
Se precisar de dados de múltiplos tenants (ex: analytics da plataforma), use agregação
via leitura sequencial dos DBs — nunca JOIN cross-tenant.

```go
// ✅ Correto: agrega sequencialmente
func CountAllPatients(controlDB *sql.DB) (int, error) {
    paths, _ := listAllTenantDBPaths(controlDB)
    total := 0
    for _, path := range paths {
        db, _ := sql.Open("sqlite3", path)
        var count int
        db.QueryRow("SELECT COUNT(*) FROM patients").Scan(&count)
        total += count
        db.Close()
    }
    return total, nil
}

// ❌ Impossível e proibido: ATTACH DATABASE para cross-tenant query
```
