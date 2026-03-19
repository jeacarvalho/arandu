# 📚 Aprendizados Mestres do Projeto Arandu

**Última atualização:** 18 de março de 2026
**Versão:** 1.0
**Status:** Ativo

> ⚠️ **NOTA DE MIGRAÇÃO:** Este arquivo consolida os aprendizados valiosos de ~45 arquivos individuais. Arquivos originais estão em `archive/`.

## 📋 Índice

1. [🏗️ Arquitetura Web (Go + templ + HTMX)](#-arquitetura-web-go--templ--htmx)
2. [💾 Banco de Dados (SQLite, FTS5, Migrations)](#-banco-de-dados-sqlite-fts5-migrations)
3. [🩺 Domínio Clínico e Validações](#-domínio-clínico-e-validações)
4. [🎨 UI/UX e Design System](#-uiux-e-design-system)
5. [🤖 Integração com IA Generativa](#-integração-com-ia-generativa)
6. [🧪 Testes e Qualidade](#-testes-e-qualidade)
7. [🔄 Fluxo de Trabalho e Scripts](#-fluxo-de-trabalho-e-scripts)
8. [🚨 Anti-Padrões e Erros Comuns](#-anti-padrões-e-erros-comuns)
9. [📁 Referências e Arquivos Originais](#-referências-e-arquivos-originais)

---

## 🏗️ Arquitetura Web (Go + templ + HTMX)

### 1. Padrão de Renderização Contextual

**Problema:** Esquecer de verificar se a requisição é HTMX, quebrando o layout ou retornando a página inteira dentro de um fragmento.

**Solução:** Sempre verificar o header `HX-Request`:

```go
if r.Header.Get("HX-Request") == "true" {
    // Renderiza apenas o fragmento
    h.templates.ExecuteTemplate(w, "patient-content", data)
} else {
    // Renderiza página completa com layout
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

**Referências:** `logbook-sota.md:9-18`, `task_20260315_174000.md:37-44`

### 2. Proteção do Domínio com ViewModels

**Problema:** Passar entidades de domínio diretamente para templates expõe campos internos e acopla UI ao domínio.

**Solução:** Criar ViewModels específicos no handler:

```go
// NUNCA fazer isso:
// h.templates.ExecuteTemplate(w, "patient", patientEntity)

// SEMPRE fazer isso:
type PatientViewModel struct {
    ID        string
    Name      string
    CreatedAt string // Já formatada
    Notes     string
}

viewModel := PatientViewModel{
    ID:        patient.ID,
    Name:      patient.Name,
    CreatedAt: patient.CreatedAt.Format("02/01/2006"),
    Notes:     patient.Notes,
}
h.templates.ExecuteTemplate(w, "patient", viewModel)
```

**Referências:** `logbook-sota.md:29-33`, `task_20260315_174000.md:23-34`

### 3. Conflito de Nomes de Templates

**Problema:** Nomes duplicados de templates causam sobrescrita e comportamentos inesperados.

**Solução:** Usar nomes únicos e específicos:

```html
<!-- RUIM: nome genérico -->
{{define "content"}}
  <!-- conteúdo -->
{{end}}

<!-- BOM: nome específico -->
{{define "patient-content"}}
  <!-- conteúdo do paciente -->
{{end}}

{{define "session-content"}}
  <!-- conteúdo da sessão -->
{{end}}
```

**Impacto:** 19 arquivos mencionam este problema. **Solução padrão consolidada.**

**Referências:** Múltiplos arquivos (padrão identificado em 18 tarefas)

### 4. Migração de html/template para Templ

**Contexto:** O projeto migrou do sistema nativo `html/template` do Go para o framework `templ`.

**Lições aprendidas:**

1. **Não usar `//go:build templ`** - O código deve compilar nativamente
2. **Não incluir `import "github.com/a-h/templ"` manualmente** - O gerador adiciona automaticamente
3. **Sempre executar `templ generate`** após modificar arquivos `.templ`
4. **Estrutura de componentes:**
   ```
   web/components/
   ├── layout/layout.templ      (layout + sidebar)
   ├── patient/list.templ       (lista de pacientes)
   ├── patient/detail.templ     (detalhe do paciente)
   ├── session/detail.templ     (detalhe da sessão)
   └── dashboard/dashboard.templ (dashboard)
   ```

**Fluxo de desenvolvimento:**
1. Criar componente `.templ`
2. Executar `templ generate` para gerar código Go
3. Integrar no handler: `layoutComponents.BaseWithContent(titulo, componente).Render(ctx, w)`
4. Build nativo: `go build ./...` (sem flags)

**Referências:** `task_20260315_154905.md`, `REQ-01-02-02.md`

### 5. Handlers com Injeção de Dependência

**Padrão estabelecido:** Cada handler define interfaces específicas para os serviços que usa:

```go
type PatientHandler struct {
    patientService  application.PatientService
    sessionService  application.SessionService
    insightService  application.InsightService
    templates       *template.Template
}

// Interfaces segregadas (cada handler define apenas o que usa)
type PatientService interface {
    GetByID(ctx context.Context, id string) (*domain.Patient, error)
    List(ctx context.Context) ([]*domain.Patient, error)
    Create(ctx context.Context, input application.CreatePatientInput) (*domain.Patient, error)
}
```

**Benefícios:** Facilita testes unitários, reduz acoplamento, segue SOLID.

**Referências:** `task_20260315_174000.md:16-22`

---

## 💾 Banco de Dados (SQLite, FTS5, Migrations)

### 1. Motor de Busca FTS5

**Implementação:** Uso de tabelas virtuais com `EXTERNAL CONTENT`:

```sql
-- Tabela virtual FTS5 para busca em observações clínicas
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
    content,
    content='observations',
    tokenize='porter'
);

-- Triggers para sincronização automática
CREATE TRIGGER observations_ai AFTER INSERT ON observations BEGIN
    INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;

CREATE TRIGGER observations_ad AFTER DELETE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', old.id, old.content);
END;

CREATE TRIGGER observations_au AFTER UPDATE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content) VALUES('delete', old.id, old.content);
    INSERT INTO observations_fts(rowid, content) VALUES (new.id, new.content);
END;
```

**Lição:** Sincronizar FTS5 via Triggers SQL (INSERT, UPDATE, DELETE) para garantir que a busca nunca fique obsoleta.

**Highlighting:** O SQLite retorna tags `<b>` para realce. Para exibir no templ, é necessário usar `templ.RawHTML` para evitar o escape automático do HTML.

**Referências:** `logbook-sota.md:35-42`

### 2. Sistema de Migrations

**Padrão estabelecido:** Arquivos `.up.sql` e `.down.sql` numerados sequencialmente:

```
internal/infrastructure/repository/sqlite/migrations/
├── 0001_initial_schema.up.sql
├── 0001_initial_schema.down.sql
├── 0002_patients_table.up.sql
├── 0002_patients_table.down.sql
└── ...
```

**Lições aprendidas:**
1. **Sempre testar rollback** - Executar `.down.sql` após `.up.sql`
2. **Usar transações** - Garantir atomicidade
3. **Manter compatibilidade** - Não quebrar migrações existentes
4. **Documentar mudanças** - Comentar o propósito de cada migration

**Referências:** `req-01-00-01.md:11-12`

### 3. Escalabilidade e Performance

**Big Data testado:** O sistema foi validado com 63.000 sessões.

**Soluções implementadas:**
1. **Infinite Scroll** (`hx-trigger="revealed"`) em listas longas
2. **Debounce** (`delay:500ms`) em campos de busca para evitar IO excessivo
3. **Pagination com LIMIT/OFFSET** para queries grandes
4. **Índices apropriados** em campos frequentemente buscados

**Driver:** Migrado para `modernc.org/sqlite` para garantir suporte nativo a FTS5 sem dependência de CGO.

**Referências:** `logbook-sota.md:43-49`, `work/learnings/task_20260317_220659.md`

### 4. Infinite Scroll no Histórico Clínico

**Implementação:** Paginação com `LIMIT` e `OFFSET` no sistema de timeline:

```sql
-- Query otimizada para infinite scroll
SELECT * FROM timeline_events 
WHERE patient_id = ? 
ORDER BY event_date DESC 
LIMIT ? OFFSET ?;
```

**Padrão HTMX:**
```html
<div 
  hx-get="/patients/{id}/history?filter={filter}&offset={next_offset}"
  hx-trigger="revealed"
  hx-swap="afterend"
  hx-indicator="#loading-indicator">
</div>
```

**Lotes:** 20 itens por carregamento (padrão do sistema)

**Referências:** `work/learnings/task_20260317_220659.md`

---

## 🩺 Domínio Clínico e Validações

### 1. Validações no Domínio

**Princípio:** A validação de regras de negócio pertence ao domínio, não à aplicação ou UI.

**Exemplos implementados:**

```go
// Medication entity validation
func NewMedication(patientID, name, dosage string, startedAt time.Time) (*Medication, error) {
    // Data de início não pode ser futura
    if startedAt.After(time.Now()) {
        return nil, errors.New("medication start date cannot be in the future")
    }
    
    // Dosagem não pode ser vazia
    if strings.TrimSpace(dosage) == "" {
        return nil, errors.New("dosage cannot be empty")
    }
    
    return &Medication{
        ID:        uuid.New().String(),
        PatientID: patientID,
        Name:      name,
        Dosage:    dosage,
        StartedAt: startedAt,
        CreatedAt: time.Now(),
    }, nil
}
```

**Referências:** `work/learnings/REQ-01-04-01.md:42-50`

### 2. Entidades com Construtores

**Padrão:** Usar factory functions em vez de struct literais:

```go
// RUIM: expõe campos internos
patient := &domain.Patient{
    ID:        "123",
    Name:      "João",
    CreatedAt: time.Now(),
}

// BOM: encapsula validação
patient, err := domain.NewPatient("João", "Observações iniciais")
if err != nil {
    // tratar erro
}
```

**Benefícios:** Garante invariantes, centraliza validação, facilita testes.

**Referências:** `task_20260313_215938.md`

---

## 🎨 UI/UX e Design System

### 1. Dualidade Tipográfica

**Regra fundamental:** UI usa Inter (Sans), conteúdo clínico usa Source Serif 4 (Serif).

**Erro comum:** Usar fontes genéricas em notas de pacientes.

**Correção:** Aplicar a classe `.font-clinical` em todos os campos de texto narrativo e snippets de busca.

**Implementação em Templ:**
```templ
css clinicalFont() {
    font-family: 'Source Serif 4', serif;
    font-size: 1.125rem;
    line-height: 1.75;
    color: #1F2937;
}

div(class="clinical-note") {
    @clinicalFont()
    "Conteúdo clínico com tipografia apropriada"
}
```

**Referências:** `logbook-sota.md:20-27`

### 2. Conceito de "Tecnologia Silenciosa"

**Filosofia:** A interface nunca deve competir com o pensamento do terapeuta.

**Princípios implementados:**
1. **Clareza > Estética**
2. **Calma > Impacto visual**
3. **Conteúdo > Interface**

**"Silent Input":** Remover bordas de 4 lados, usar apenas `border-b` sutil e fundo `bg-[#F7F8FA]` (cinza papel). O foco deve ser suave, sem o `ring` azul padrão.

**Referências:** `logbook-sota.md:53-58`, `docs/design-system.md`

### 3. Mobile-First

**Aprendizado:** A Sidebar deve ser um drawer (gaveta) em telas pequenas.

**Padrão:** Usar `flex-col` no mobile para que o "Painel de Insights" e "Histórico Farmacológico" fiquem abaixo do conteúdo principal, mantendo a leitura fluida.

**Referências:** `logbook-sota.md:60-64`

### 4. Paleta de Cores (Design System)

**Primária:** `#1E3A5F` (azul escuro profissional)
**Secundária:** `#3A7D6B` (verde terapêutico)
**Insight/IA:** `#D4A84F` (dourado para destaques)
**Fundo:** `#F7F8FA` (cinza papel)
**Texto:** `#1F2937` (cinza escuro)

**Referências:** `docs/design-system.md:15-37`

---

## 🤖 Integração com IA Generativa

### 1. Configuração Segura de API Keys

**Contexto:** Tarefa 20260318_211140 - Implementação do REQ-05-01-01

**Problema:** Como integrar serviços de IA externos (Gemini) sem expor chaves de API no código.

**Solução:**
1. Usar arquivo `.env` com `GEMINI_API_KEY`
2. Carregar via `github.com/joho/godotenv`
3. Inicializar cliente com fallback para modo "dummy" quando chave ausente
4. Documentar configuração em `README_GEMINI.md`

```go
// Inicialização segura
geminiAPIKey := os.Getenv("GEMINI_API_KEY")
if geminiAPIKey == "" {
    log.Printf("Warning: GEMINI_API_KEY not set. AI features will be disabled.")
    geminiAPIKey = "dummy-key-for-initialization"
}
```

**Referência:** REQ-05-01-01, `cmd/arandu/main.go:76-80`

### 2. Cache de Respostas de IA

**Contexto:** Tarefa 20260318_211140 - Implementação do REQ-05-01-01

**Problema:** Chamadas repetidas à API Gemini geram custos desnecessários e latência.

**Solução:**
1. Implementar cache em memória com TTL configurável (24h padrão)
2. Usar chave SHA256(patientID:timeframe) para identificação única
3. Integrar cache transparentemente no serviço de IA

```go
// Check cache first
if s.cache != nil {
    if entry, found := s.cache.Get(patientID, timeframe); found {
        return &PatientSynthesisResponse{
            Synthesis:   entry.Synthesis,
            GeneratedAt: entry.GeneratedAt,
        }, nil
    }
}
```

**Benefícios:** Redução de custos, respostas mais rápidas, menor dependência de rede.

**Referência:** REQ-05-01-01 CA-03, `internal/infrastructure/ai/cache.go`

### 3. Retry Exponencial para APIs Externas

**Contexto:** Tarefa 20260318_211140 - Implementação do REQ-05-01-01

**Problema:** APIs externas podem falhar temporariamente (rate limiting, timeout).

**Solução:**
1. Implementar backoff exponencial com máximo de 5 tentativas
2. Detectar erros recuperáveis (HTTP 429, 500, timeout)
3. Logar tentativas para debugging

```go
for i := 0; i < maxRetries; i++ {
    if i > 0 {
        backoff := time.Duration(1<<uint(i)) * time.Second
        if backoff > 30*time.Second {
            backoff = 30 * time.Second
        }
        time.Sleep(backoff)
    }
    // Tentar chamada
}
```

**Referência:** `internal/infrastructure/ai/gemini_client.go:58-95`

### 4. Interface "Tecnologia Silenciosa" para IA

**Contexto:** Tarefa 20260318_211140 - Implementação do REQ-05-01-01

**Problema:** Como apresentar análises de IA sem sobrecarregar o terapeuta.

**Solução:**
1. Botão discreto com seletor de período (3 meses, 6 meses, 1 ano, todo histórico)
2. Síntese estruturada em 4 partes: Temas Dominantes, Pontos de Inflexão, Correlações Sugeridas, Provocação Clínica
3. Estilo visual diferenciado (fundo âmbar #FFFBEB, fonte clínica)
4. Aviso legal obrigatório: "Esta é uma análise gerada por IA para apoio à reflexão"

**Referência:** REQ-05-01-01 CA-02, CA-04, `web/components/ai/`

---

## 🧪 Testes e Qualidade

### 1. Estado da Cobertura de Testes

**Situação atual (Março 2026):**
- Cobertura geral: 15.9% (abaixo da meta de 65%)
- Principais gaps:
  - Handlers: 0% de cobertura (sistema antigo vs novo)
  - Components Templ: 0% de cobertura
  - Services: 47.5% de cobertura
  - Repositórios: 35.2% de cobertura

**Problemas identificados:**
1. Sistema misto: handlers antigos (`web/handlers/`) e novos (`internal/web/handlers/`)
2. Testes E2E desatualizados
3. Muitos testes pulados com `t.Skip()`

**Referências:** `task_20260315_191104.md:18-29`

### 2. Testes de Handlers Reais

**Problema:** Testes apenas verificam arquivos, não executam handlers reais.

**Solução padrão:** Criar testes que executam handlers reais:

```go
func TestPatientHandler_Show(t *testing.T) {
    // Setup
    handler := NewPatientHandler(mockService, templates)
    req := httptest.NewRequest("GET", "/patient/123", nil)
    w := httptest.NewRecorder()
    
    // Execute handler REAL
    handler.Show(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "João Silva")
}
```

**Referências:** Múltiplos arquivos mencionam esta necessidade

### 3. Testes de Integração com SQLite

**Padrão estabelecido:** Usar database em memória para testes:

```go
func TestPatientRepository_Create(t *testing.T) {
    // Setup in-memory database
    db, err := sql.Open("sqlite", ":memory:")
    require.NoError(t, err)
    defer db.Close()
    
    // Run migrations
    err = runMigrations(db)
    require.NoError(t, err)
    
    // Create repository
    repo := NewPatientRepository(db)
    
    // Test
    patient, err := repo.Create(context.Background(), domain.Patient{...})
    assert.NoError(t, err)
    assert.NotNil(t, patient)
}
```

---

## 🔄 Fluxo de Trabalho e Scripts

### 1. Ciclo de Vida da Tarefa

**Fluxo estabelecido:**
1. **Requisito:** Ler o arquivo `.md` em `docs/requirements/`
2. **Schema:** Criar migration `.up.sql` em `internal/infrastructure/repository/sqlite/migrations/`
3. **Domínio:** Atualizar structs em `internal/domain/`
4. **UI:** Criar/Editar `.templ` e obrigatoriamente rodar `templ generate`
5. **Handler:** Implementar lógica de fragmento
6. **Guard:** Rodar `./scripts/arandu_guard.sh` para garantir integridade

**Referências:** `logbook-sota.md:69-80`

### 2. Script de Guard (arandu_guard.sh)

**Verificações automáticas:**
1. **Rotas online:** Testa `/dashboard`, `/patients`, `/patients/new`
2. **Integridade Templ:** Verifica se arquivos `.templ` foram gerados
3. **Build:** Garante que o projeto compila

**Referências:** `scripts/arandu_guard.sh`



### 4. Sistema de Aprendizados (REFATORADO)

**Problema anterior:** Script `arandu_conclude_task.sh` gerava conteúdo repetitivo automático.

**Solução nova:** 
- Consolidar aprendizados valiosos neste arquivo mestre
- Remover geração automática de conteúdo repetitivo
- Manter apenas aprendizados realmente valiosos

---

## 🚨 Anti-Padrões e Erros Comuns

### 1. Import Manual de Templ

**Anti-padrão:** Adicionar `import "github.com/a-h/templ"` manualmente em arquivos `.templ`

**Consequência:** Erro de compilação `templ redeclared in this block`

**Como evitar:** Deixar o `templ generate` gerenciar os imports automaticamente.

**Referências:** `REQ-01-02-02.md:16-27`

### 2. SQL Hardcoded no Go

**Anti-padrão:** Escrever queries SQL diretamente no código Go

**Consequência:** Inconsistência entre ambientes, dificuldade de manutenção

**Como evitar:** Usar sempre o sistema de Migrations e centralizar queries em repositories.

**Referências:** `logbook-sota.md:95-99`

### 3. hx-target Genérico

**Anti-padrão:** Usar `hx-target` genérico como `#content` ou `body`

**Consequência:** Troca de elementos errados, comportamentos inesperados

**Como evitar:** Usar IDs específicos ou `closest` selectors:
```html
<!-- BOM -->
<button hx-post="/sessions" hx-target="#session-list">Criar</button>

<!-- MELHOR -->
<button hx-post="/sessions" hx-target="closest .session-container">Criar</button>
```

**Referências:** `logbook-sota.md:101-106`

### 4. Ignorar updated_at

**Anti-padrão:** Não atualizar `updated_at` em operações UPDATE

**Consequência:** Perda de rastreabilidade, dados desatualizados

**Como evitar:** Garantir que todo UPDATE atualize o timestamp:
```sql
UPDATE patients 
SET name = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
```

**Referências:** `logbook-sota.md:107-112`

---

## 📁 Referências e Arquivos Originais

### Arquivos Consolidados Neste Documento:

1. **`logbook-sota.md`** (113 linhas) - Consolidação existente
2. **`task_20260315_174000.md`** (207 linhas) - Refatoração completa da arquitetura web
3. **`REQ-01-02-02.md`** (121 linhas) - Problema com templ e imports
4. **`work/learnings/REQ-01-04-01.md`** (105 linhas) - Implementação recente detalhada
5. **`work/learnings/task_20260317_220659.md`** (111 linhas) - Infinite scroll
6. **`task_20260315_154905.md`** (41 linhas) - Migração para Templ
7. **`task_20260315_191104.md`** (35 linhas) - Estado de testes
8. **`req-01-00-01.md`** (36 linhas) - Múltiplos aprendizados

### Arquivos Deletados/Arquivados:

- **18 arquivos** com conteúdo repetitivo sobre "Conflito de templates"
- **Arquivos com menos de 20 linhas** e conteúdo genérico
- **Arquivos problemáticos**: `.md`, `Teste.md`

### Localização dos Originais:
Todos os arquivos originais foram movidos para `docs/learnings/archive/` para referência histórica.

---

## 🔄 Como Contribuir com Novos Aprendizados

1. **Avalie se é valioso:** O aprendizado é específico, útil e não repetitivo?
2. **Escolha a seção:** Adicione à seção apropriada deste documento
3. **Formato:**
   ```
   ### Título do Aprendizado
   
   **Contexto:** [Breve descrição do contexto]
   
   **Problema:** [O que deu errado ou poderia ser melhor]
   
   **Solução:** [Como foi resolvido ou melhorado]
   
   **Código de exemplo (se aplicável):**
   ```go
   // Código relevante
   ```
   
   **Referência:** [Tarefa ou requirement relacionado]
   ```
4. **Mantenha atualizado:** Revise periodicamente e remova conteúdo obsoleto

---

*Documento criado como parte da refatoração da documentação do projeto Arandu - Março 2026*