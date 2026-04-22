---
name: arandu-architecture
description: >
  Documento mestre de arquitetura do sistema Arandu — plataforma de inteligência clínica
  para psicólogos. Use esta skill SEMPRE que trabalhar no projeto Arandu, independente
  da tarefa: criar handler, migration, componente templ, rota, serviço, teste, ou qualquer
  artefato do sistema. Também dispara para: "como estruturar X no Arandu?", "onde fica
  esse arquivo?", "como funciona o multi-tenancy?", "qual a rota para X?", "como injetar
  o db do tenant?", "como rodar os scripts?". Esta skill deve ser lida ANTES de qualquer
  outra skill quando o contexto é o projeto Arandu — ela define a estrutura real do
  projeto que sobrepõe sugestões genéricas de outras skills.
---

# Arandu — Arquitetura do Sistema

Documento mestre. Leia antes de qualquer tarefa no projeto Arandu.
Qualquer conflito entre esta skill e outra skill genérica: **esta skill prevalece**.

---

## O que é o Arandu

Sistema de **inteligência clínica** para profissionais de saúde mental. Construído com
DDD + Clean Architecture + renderização type-safe. Projetado para resiliência com grandes
volumes de dados clínicos ("Big Data Clínico") e **privacidade absoluta** por isolamento
físico de dados.

Filosofia: **"Tecnologia Silenciosa"** — a UI não compete com o conteúdo clínico.
O sistema desaparece; o psicólogo e o paciente permanecem.

---

## Estrutura real do projeto

```
arandu/
├── cmd/arandu/
│   └── main.go                          # Ponto de entrada
├── internal/
│   ├── domain/                          # Domínio puro — zero deps externas
│   │   ├── patient/
│   │   │   ├── patient.go               # Aggregate Root Patient
│   │   │   ├── context.go               # PatientContext (biopsicossocial)
│   │   │   ├── medication.go            # Medication
│   │   │   └── vitals.go                # Vitals
│   │   ├── session/                     # Session + SessionForm
│   │   ├── observation/                 # Observation (FTS5)
│   │   ├── intervention/                # Intervention
│   │   └── timeline/                    # Agregadores longitudinais (read models)
│   ├── application/
│   │   └── services/                    # Application Services (orquestração)
│   │       ├── patient_service.go
│   │       ├── session_service.go
│   │       └── biopsychosocial_service.go
│   └── infrastructure/
│       └── repository/
│           └── sqlite/                  # Implementações de repository
│               ├── migrations/          # SQL puro — NUNCA schema em Go hardcoded
│               │   ├── 0001_initial_schema.up.sql
│               │   ├── 0002_enable_fts5.up.sql
│               │   ├── 0003_add_biopsychosocial_tables.up.sql
│               │   └── migrator.go
│               ├── patient_repository.go
│               ├── session_repository.go
│               ├── medication_repository.go
│               └── vitals_repository.go
├── web/                                 # Camada web — FORA de internal/
│   ├── handlers/                        # HTTP handlers (não em internal/infra)
│   │   ├── patient_handler.go
│   │   ├── session_handler.go
│   │   └── biopsychosocial_handler.go
│   ├── components/                      # Componentes .templ por entidade
│   │   ├── patient/
│   │   │   ├── biopsychosocial_panel.templ
│   │   │   ├── medication_list.templ
│   │   │   └── vitals_widget.templ
│   │   └── session/
│   └── static/                         # CSS (Tailwind build), fontes
├── scripts/                            # Automação e guard
│   ├── arandu_guard.sh                 # Verifica integridade do sistema
│   ├── arandu_validate_handlers.sh     # Valida handlers
│   ├── arandu_checkpoint.sh            # Checkpoint arquitetural
│   └── arandu_update_context.sh        # Atualiza contexto após mudanças de rota
└── docs/                               # Requisitos, visões, estratégias
```

> ⚠️ **Atenção**: handlers ficam em `web/handlers/`, não em `internal/infra/http/`.
> Componentes ficam em `web/components/`, não em `internal/ui/`.
> Esta estrutura difere das sugestões genéricas da skill `ddd-go`.

---

## Multi-tenancy — o pilar central

### Modelo: Database-per-tenant

Cada psicólogo tem um arquivo SQLite físico separado. Não há tabela `tenant_id` —
o isolamento é **físico**, não lógico.

```
Control Plane:  arandu_control.db        ← usuários, tenants, paths dos DBs
Data Plane:     {clinical_uuid}.db       ← dados clínicos de um psicólogo
                {clinical_uuid2}.db      ← dados clínicos de outro psicólogo
```

### Fluxo de autenticação e injeção de contexto

```
1. Psicólogo faz login
2. Control Plane autentica e resolve: psychologist_id → db_file_path
3. db_file_path é injetado no contexto do request (context.Context)
4. Cada handler extrai a conexão correta do contexto
5. Todos os repositories da request usam essa conexão
```

### Implementação do contexto de tenant

```go
// internal/infrastructure/tenant/context.go

type contextKey string
const tenantDBKey contextKey = "tenant_db"

// Injetado pelo middleware de autenticação
func WithTenantDB(ctx context.Context, db *sql.DB) context.Context {
    return context.WithValue(ctx, tenantDBKey, db)
}

// Usado por handlers e repositories
func TenantDB(ctx context.Context) (*sql.DB, error) {
    db, ok := ctx.Value(tenantDBKey).(*sql.DB)
    if !ok || db == nil {
        return nil, ErrNoTenantDB
    }
    return db, nil
}
```

### Middleware de tenant

```go
// web/handlers/middleware.go

func TenantMiddleware(controlDB *sql.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Extrai session/token do request
            psychologistID := extractPsychologistID(r)

            // 2. Resolve o path do db no Control Plane
            dbPath, err := resolveTenantDB(controlDB, psychologistID)
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            // 3. Abre (ou reutiliza do pool) a conexão com o db do tenant
            tenantDB, err := openTenantDB(dbPath)
            if err != nil {
                http.Error(w, "internal error", http.StatusInternalServerError)
                return
            }

            // 4. Injeta no contexto
            ctx := tenant.WithTenantDB(r.Context(), tenantDB)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Migrations em todos os tenants

```go
// internal/infrastructure/repository/sqlite/migrations/migrator.go
// Roda no startup — aplica migrations em TODOS os tenant DBs

func RunAllTenantMigrations(controlDB *sql.DB, migrationsFS embed.FS) error {
    tenantPaths, err := listAllTenantDBPaths(controlDB)
    if err != nil {
        return err
    }
    for _, path := range tenantPaths {
        db, err := sql.Open("sqlite3", path)
        if err != nil {
            return fmt.Errorf("open tenant %s: %w", path, err)
        }
        if err := applyMigrations(db, migrationsFS); err != nil {
            return fmt.Errorf("migrate tenant %s: %w", path, err)
        }
        db.Close()
    }
    return nil
}
```

> **Regra**: NUNCA crie schema via código Go hardcoded.
> Sempre use arquivos `.sql` em `migrations/` com prefixo numérico.

---

## Convenções de rotas

### Princípio singular/plural

| Padrão | Uso | Exemplo |
|--------|-----|---------|
| Plural | Coleção | `GET /patients` |
| Singular | Recurso específico | `GET /patient/{id}` |
| Singular + plural | Sub-recursos | `GET /patient/{id}/sessions` |

### Tabela completa de rotas

```
# Pacientes
GET  /patients                     → listar todos
POST /patients                     → criar novo
GET  /patients/new                 → formulário de criação
GET  /patient/{id}                 → detalhes
GET  /patient/{id}/sessions        → sessões do paciente
GET  /patient/{id}/sessions/new    → formulário nova sessão
GET  /patient/{id}/history         → timeline longitudinal

# Sessões
POST /session                      → criar nova sessão
GET  /session/{id}                 → detalhes
GET  /session/{id}/edit            → formulário edição
POST /session/{id}/update          → atualizar
POST /session/{id}/observations    → criar observação
POST /session/{id}/interventions   → criar intervenção

# Observações
GET  /observations/{id}            → detalhes
GET  /observations/{id}/edit       → formulário edição
PUT  /observations/{id}            → atualizar

# Intervenções
GET  /interventions/{id}           → detalhes
GET  /interventions/{id}/edit      → formulário edição
PUT  /interventions/{id}           → atualizar

# Sistema
GET  /                             → redirect para /dashboard
GET  /dashboard                    → página inicial
GET  /static/                      → arquivos estáticos
```

### Extração de IDs em handlers

```go
// Sempre use funções auxiliares — nunca parse manual inline
func extractPatientID(path string) string {
    parts := strings.Split(path, "/")
    for i, part := range parts {
        if part == "patient" && i+1 < len(parts) {
            return parts[i+1]
        }
    }
    return ""
}
```

### URLs em templates Templ

```go
// ✅ SEMPRE use templ.URL()
<a href={ templ.URL("/patient/" + patientID) }>Ver paciente</a>

// ❌ NUNCA string literal com interpolação
<a href="/patient/{patientID}">Ver paciente</a>
```

---

## Silent UI — sistema de design

O Arandu usa uma estética de **"anti-fadiga"** — a UI desaparece para que o conteúdo
clínico seja o protagonista.

### Tokens de design

```css
/* Fundo global */
background: #F7F8FA;  /* Cinza papel — não branco puro */

/* Tipografia dual */
--font-ui:       'Inter', sans-serif;          /* Toda a interface */
--font-clinical: 'Source Serif 4', serif;      /* Conteúdo clínico */

/* Silent inputs — sem bordas pesadas */
.input-silent {
  @apply bg-transparent border-0 border-b border-gray-200
         focus:border-gray-400 focus:ring-0
         text-gray-900 placeholder-gray-300;
}
```

### Regra tipográfica

| Contexto | Fonte | Tailwind |
|----------|-------|---------|
| Labels, botões, navegação, metadados | Inter | `font-sans` |
| Observações clínicas, notas, prontuário | Source Serif 4 | `font-serif text-xl leading-relaxed` |
| Código, IDs, timestamps | Mono | `font-mono text-sm` |

### Hierarquia visual

```templ
// Conteúdo clínico — sempre Serif
<p class="font-serif text-xl leading-relaxed text-gray-800">
    { observation.Content }
</p>

// Metadados — sempre Sans
<span class="font-sans text-xs text-gray-400">
    { observation.CreatedAt.Format("02/01/2006 15:04") }
</span>
```

---

## Protocolo de implementação (checklist por tarefa)

### Ao criar um novo handler

```
[ ] Arquivo em web/handlers/ (não em internal/)
[ ] Extrai tenant DB do contexto: tenant.TenantDB(r.Context())
[ ] Verifica HX-Request para servir fragmento vs página completa
[ ] Usa ViewModel — nunca passa domain struct direto ao template
[ ] Verifica método HTTP (GET/POST/PUT) explicitamente
[ ] Rota segue convenção singular/plural documentada
```

### Ao criar uma nova migration

```
[ ] Arquivo SQL em internal/infrastructure/repository/sqlite/migrations/
[ ] Nome com prefixo numérico sequencial: NNNN_descricao.up.sql
[ ] NUNCA schema hardcoded em Go
[ ] Testa em DB limpo antes de commitar
[ ] Verifica se migrator.go aplica a nova migration nos tenants existentes
```

### Ao criar um componente Templ

```
[ ] Arquivo em web/components/{entidade}/
[ ] Executa templ generate antes de testar
[ ] Conteúdo clínico usa font-serif
[ ] URLs via templ.URL()
[ ] Página completa herda de templates.Layout()
```

### Após mudanças de rota

```bash
./scripts/arandu_update_context.sh   # atualiza contexto
./scripts/arandu_validate_handlers.sh # valida handlers
./scripts/arandu_guard.sh             # verifica integridade geral
```

---

## Fases de evolução

| Fase | Status | Descrição |
|------|--------|-----------|
| 1 | ✅ Concluída | Consolidação Web e DDD |
| 2 | 🔄 Em andamento | Multi-tenancy e Gestão de Acesso |
| 3 | 🔄 Concluindo | Escalabilidade, FTS5 e Navegação em Larga Escala |
| 4 | 🚧 Em implementação | Contexto Biopsicossocial e Identidade SOTA |
| 5 | 📋 Planejada | Inteligência Assistida (Vision-05) |

---

## Referências
- `references/tenant-patterns.md` — padrões avançados de multi-tenancy: pool de conexões, failover, backup por tenant
- `references/fts5-patterns.md` — busca full-text com SQLite FTS5: indexação, queries, performance
