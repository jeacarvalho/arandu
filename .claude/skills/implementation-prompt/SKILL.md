---
name: implementation-prompt
description: >
  Gera Prompts de Implementação estruturados para agentes codificadores (OpenCode,
  Gemini CLI, DeepSeek, modelos locais via Ollama) executarem tarefas no projeto Arandu.
  Use esta skill SEMPRE que o usuário pedir para "criar um prompt para o agente", "gerar
  o briefing de implementação", "preparar a tarefa para o codificador", "delegar essa
  feature para o OpenCode", ou qualquer variação de "quero que um modelo mais barato
  implemente isso". Também dispara quando o usuário descreve uma feature e pergunta
  "como eu passo isso para o agente?". O output desta skill é um documento Markdown
  completo, pronto para ser colado no agente codificador — não uma conversa, um artefato.
---

# Implementation Prompt — Skill de Delegação para Agentes Codificadores

Esta skill transforma intenções de feature em **Prompts de Implementação** — documentos
precisos que um agente codificador (OpenCode + DeepSeek, Gemini CLI, modelo local) pode
executar sem raciocinar sobre arquitetura, domínio ou privacidade.

**Princípio central**: o Claude resolve o *quê* e o *porquê*. O agente codificador
executa o *como*. A qualidade do prompt é o que separa execução boa de execução cara.

---

## Quando usar esta skill

- Usuário descreve uma feature e quer delegar a implementação
- Usuário tem um requisito (REQ-XXX) pronto e quer o prompt de execução
- Usuário quer paralelizar trabalho entre Claude Code e agente mais barato
- Qualquer tarefa de codificação que seja *execução guiada*, não *design*

## Quando NÃO usar — manter no Claude Code

- Modelagem de novo Aggregate ou Bounded Context (requer raciocínio DDD)
- Features que tocam privacidade de forma não trivial (requer julgamento)
- Refatorações cross-cutting que afetam múltiplas camadas
- Qualquer decisão arquitetural nova

---

## Processo de geração

### Passo 1 — Classificar a tarefa

Antes de gerar, classifique:

| Tipo | Características | Delegável? |
|------|----------------|------------|
| **CRUD simples** | Uma entidade, operações padrão, sem regra de negócio complexa | ✅ Sim |
| **Componente UI** | Templ + HTMX, sem lógica de domínio | ✅ Sim |
| **Migration** | Schema SQL descrito, sem ambiguidade | ✅ Sim |
| **Feature com regra de negócio** | Invariantes de domínio envolvidas | ⚠️ Apenas com spec detalhada |
| **Feature com IA/privacidade** | Anonimização, payload para LLM | ❌ Manter no Claude Code |
| **Refatoração arquitetural** | Múltiplos pacotes, decisões de design | ❌ Manter no Claude Code |

### Passo 2 — Coletar contexto necessário

Antes de escrever o prompt, confirme com o usuário:
- Qual o requisito ou descrição da feature?
- Já existe código relacionado que o agente deve seguir como padrão?
- Há alguma restrição específica além das padrões do Arandu?
- Qual o critério de "pronto"? (testes, compilação, validação visual)

### Passo 3 — Gerar o Prompt de Implementação

Use o template abaixo. Preencha **todas** as seções — nunca deixe seção vaga.
Um prompt incompleto é mais perigoso que nenhum prompt.

---

## Template do Prompt de Implementação

```markdown
# ARANDU — Prompt de Implementação
## [TÍTULO DA FEATURE — verbo + substantivo do domínio]

---

### 🎯 Objetivo
[1-3 frases descrevendo o que deve existir ao final. Use a Ubiquitous Language do domínio.]

---

### 🏗️ Contexto do sistema (leia antes de escrever qualquer código)

**Stack**: Go 1.22+ · Templ · HTMX · Tailwind CSS · SQLite (database-per-tenant)

**Estrutura de pastas** (não desvie):
```
internal/domain/{entidade}/         ← domínio puro, zero deps externas
internal/application/services/      ← application services
internal/infrastructure/repository/sqlite/  ← repositories + migrations SQL
web/handlers/                       ← HTTP handlers
web/components/{entidade}/          ← componentes .templ
```

**Multi-tenancy** (crítico): cada psicólogo tem seu próprio arquivo `.db`.
O handler DEVE obter a conexão assim — sem exceção:
```go
db, err := tenant.TenantDB(r.Context())
if err != nil {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
}
```

**HTMX**: handler verifica `HX-Request` header para retornar fragmento ou página completa:
```go
if r.Header.Get("HX-Request") == "true" {
    // retorna fragmento
    components.MeuComponente(vm).Render(r.Context(), w)
    return
}
// retorna página completa
templates.Layout(pages.MinhaPagina(vm)).Render(r.Context(), w)
```

**URLs em Templ**: sempre `templ.URL("/rota/" + id)` — nunca string literal.

---

### 📦 Domínio

**Entidade(s) envolvida(s)**: [Nome exato da Ubiquitous Language]

**Campos**:
| Campo | Tipo Go | Tipo SQL | Regra |
|-------|---------|----------|-------|
| id | uuid.UUID | TEXT PRIMARY KEY | gerado no construtor |
| [campo] | [tipo] | [tipo SQL] | [regra ou "obrigatório"] |

**Invariantes** (regras que o domínio deve garantir):
- [regra 1: ex: "content não pode ser vazio"]
- [regra 2: ex: "session_id deve existir"]

**Construtor**:
```go
// New[Entidade] — cria nova instância, valida invariantes
func New[Entidade]([params]) (*[Entidade], error) { ... }

// Reconstitute[Entidade] — rebuild do banco, sem re-validar
func Reconstitute[Entidade]([params]) *[Entidade] { ... }
```

---

### 🛣️ Rotas

| Método | Rota | Descrição | Retorna |
|--------|------|-----------|---------|
| [METHOD] | [/rota/{id}] | [descrição] | [fragmento HTMX \| página completa] |

---

### 📁 Arquivos a criar/modificar

**Criar:**
- `internal/domain/[entidade]/[entidade].go` — struct, construtor, Reconstitute, métodos de domínio
- `internal/infrastructure/repository/sqlite/[entidade]_repository.go` — Save, FindByID, [outros]
- `internal/infrastructure/repository/sqlite/migrations/[NNNN]_[descricao].up.sql` — schema SQL
- `web/handlers/[entidade]_handler.go` — handlers HTTP
- `web/components/[entidade]/[componente].templ` — UI

**Modificar:**
- `cmd/arandu/main.go` — registrar novas rotas
- [outros arquivos se necessário]

---

### 🔒 Restrições de privacidade

[Selecione as aplicáveis e remova as demais]

- [ ] **Sem restrições especiais** — entidade não contém dados de paciente
- [ ] **Tier 2 — texto livre**: campos `[listar]` NUNCA saem do sistema para IA externa
- [ ] **Tier 1 — PII direto**: campos `[listar]` NUNCA em logs ou audit trail
- [ ] **Tier 1-Plus**: `patient_context.*` — proteção máxima, nunca em nenhum payload externo
- [ ] **RiskIndicator presente**: bloquear chamada de IA se registrado nas últimas 3 sessões

---

### ✅ Critérios de aceite

O código está pronto quando:

**Compilação**
- [ ] `templ generate` executa sem erros
- [ ] `go build ./...` compila sem erros

**Testes**
- [ ] `go test ./internal/domain/[entidade]/...` passa
- [ ] `go test ./internal/infrastructure/repository/sqlite/...` passa (SQLite in-memory)

**Comportamento**
- [ ] [CA01: Dado X, quando Y, então Z]
- [ ] [CA02: Dado X inválido, quando Y, então sistema retorna erro/status correto]
- [ ] [CA03 negativo: caso de borda ou rejeição esperada]

**Integridade do sistema**
- [ ] `./scripts/arandu_guard.sh` passa sem erros
- [ ] `./scripts/arandu_validate_handlers.sh` passa

---

### 🚫 NÃO faça

- Não crie schema SQL em código Go — use arquivo `.sql` em `migrations/`
- Não use `html/template` — use apenas `.templ`
- Não passe domain struct diretamente ao template — use ViewModel
- Não acesse tenant DB globalmente — sempre extraia do `context.Context`
- Não invente rotas — siga a tabela acima
- Não use `float64` para valores monetários — use `int64` (centavos)
- [restrições específicas da feature, se houver]

---

### 📎 Padrão de referência

[Se houver código existente similar no projeto que o agente deve seguir como modelo,
cite aqui o arquivo e o que deve ser replicado. Ex:]

Siga o padrão de `web/handlers/session_handler.go` para estrutura do handler.
Siga o padrão de `internal/infrastructure/repository/sqlite/session_repository.go`
para o repository — especialmente o uso de `ReconstituteSession`.
```

---

## Exemplo preenchido — Criar Observation

```markdown
# ARANDU — Prompt de Implementação
## Registrar Observation em uma Session

---

### 🎯 Objetivo
Permitir que o Psicólogo adicione Observations de texto livre a uma Session existente.
Observations são indexadas via FTS5 para busca instantânea. Uma Session pode ter
múltiplas Observations.

---

### 🏗️ Contexto do sistema (leia antes de escrever qualquer código)

**Stack**: Go 1.22+ · Templ · HTMX · Tailwind CSS · SQLite (database-per-tenant)

**Estrutura de pastas** (não desvie):
```
internal/domain/observation/
internal/infrastructure/repository/sqlite/
web/handlers/
web/components/session/
```

**Multi-tenancy**: handler DEVE obter conexão via `tenant.TenantDB(r.Context())`.

**HTMX**: POST /session/{id}/observations retorna fragmento com a nova Observation
inserida no topo da lista — não recarrega a página.

---

### 📦 Domínio

**Entidade**: `Observation`

**Campos**:
| Campo | Tipo Go | Tipo SQL | Regra |
|-------|---------|----------|-------|
| id | uuid.UUID | TEXT PRIMARY KEY | gerado no construtor |
| session_id | uuid.UUID | TEXT NOT NULL | FK para sessions |
| content | string | TEXT NOT NULL | obrigatório, mín 3 chars |
| created_at | time.Time | DATETIME NOT NULL | gerado no construtor |
| updated_at | time.Time | DATETIME NOT NULL | gerado no construtor |

**Invariantes**:
- content não pode ser vazio ou menor que 3 caracteres
- session_id deve ser um UUID válido

**Construtores**:
```go
func NewObservation(sessionID uuid.UUID, content string) (*Observation, error)
func ReconstituteObservation(id, sessionID uuid.UUID, content string, createdAt, updatedAt time.Time) *Observation
```

---

### 🛣️ Rotas

| Método | Rota | Descrição | Retorna |
|--------|------|-----------|---------|
| POST | /session/{id}/observations | Cria nova Observation | Fragmento HTMX: novo item na lista |
| GET | /observations/{id} | Detalhes da Observation | Fragmento ou página |
| GET | /observations/{id}/edit | Formulário de edição | Fragmento HTMX |
| PUT | /observations/{id} | Atualiza Observation | Fragmento HTMX atualizado |

---

### 📁 Arquivos a criar/modificar

**Criar:**
- `internal/domain/observation/observation.go`
- `internal/infrastructure/repository/sqlite/observation_repository.go` — Save, FindByID, FindBySession
- `internal/infrastructure/repository/sqlite/migrations/0004_add_observations_fts.up.sql`
- `web/handlers/observation_handler.go`
- `web/components/session/observation_item.templ`
- `web/components/session/observation_form.templ`

**Modificar:**
- `cmd/arandu/main.go` — registrar rotas de observations

---

### 🔒 Restrições de privacidade

- [x] **Tier 2 — texto livre**: campo `content` NUNCA sai do sistema para IA externa.
  Não incluir em nenhum payload de análise. Não logar o conteúdo — apenas metadados.

---

### ✅ Critérios de aceite

**Compilação**
- [ ] `templ generate` sem erros
- [ ] `go build ./...` sem erros

**Testes**
- [ ] `go test ./internal/domain/observation/...` — NewObservation valida content vazio
- [ ] `go test ./internal/infrastructure/repository/sqlite/...` — Save e FindBySession com SQLite in-memory

**Comportamento**
- [ ] CA01: POST /session/{id}/observations com content válido → 200 + fragmento com nova Observation
- [ ] CA02: POST com content vazio → 422 Unprocessable Entity
- [ ] CA03: POST com session_id inexistente → 404 Not Found
- [ ] CA04: Observation aparece na lista sem reload da página (HTMX swap)

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa

---

### 🚫 NÃO faça

- Não crie a tabela FTS5 em Go — use o arquivo `.sql` de migration
- Não inclua `content` em nenhum log ou payload externo
- Não retorne a página completa no POST — apenas o fragmento HTMX da nova Observation
- Não use html/template — apenas .templ

---

### 📎 Padrão de referência

Siga `internal/infrastructure/repository/sqlite/session_repository.go` para estrutura
do repository e uso de `ReconstituteSession`.
Siga `web/handlers/session_handler.go` para padrão de handler com verificação HX-Request.
```

---

## Checklist de qualidade do prompt

Antes de entregar o prompt ao agente, verifique:

```
COMPLETUDE
[ ] Objetivo está em 1-3 frases claras com termos do domínio?
[ ] Todos os campos da entidade estão listados com tipos Go E SQL?
[ ] Invariantes de domínio estão explícitas?
[ ] Todos os arquivos a criar estão listados?
[ ] Rotas seguem a convenção singular/plural do Arandu?

CONTEXTO TÉCNICO
[ ] Multi-tenancy está documentado com o trecho de código exato?
[ ] Padrão HTMX (fragmento vs página) está especificado por rota?
[ ] Há um arquivo de referência para o agente seguir como modelo?

PRIVACIDADE
[ ] Campos Tier 1, Tier 2 ou Tier 1-Plus foram identificados?
[ ] Restrição de não logar / não enviar para IA está explícita?

CRITÉRIOS DE ACEITE
[ ] Há pelo menos um CA positivo, um negativo e um de compilação?
[ ] Os testes estão especificados (não só "escreva testes")?
[ ] `arandu_guard.sh` está nos critérios?

NÃO FAZER
[ ] As proibições específicas da feature estão listadas?
```

---

## Níveis de detalhe por complexidade

**Tarefa simples** (CRUD, componente isolado)
→ Use o template completo — agente não deve tomar nenhuma decisão de design.

**Tarefa média** (feature com regra de negócio)
→ Template completo + adicione seção `## Lógica de negócio detalhada` com pseudocódigo
dos casos de borda.

**Tarefa complexa** (múltiplos aggregates, feature de IA)
→ Não delegar. Manter no Claude Code.
