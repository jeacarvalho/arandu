# Convenções de Rotas - Arandu

**Documento de Referência Arquitetural**

---

## Princípios Gerais

### 1. RESTful com Singular/Plural Consistente

Seguimos convenções REST com distinção clara entre coleções (plural) e recursos específicos (singular):

- **Coleções**: Usam plural (`/patients`, `/sessions`)
- **Recursos específicos**: Usam singular (`/patient/{id}`, `/session/{id}`)
- **Sub-recursos**: Usam plural após recurso singular (`/patient/{id}/sessions`)

### 2. Hierarquia de Recursos

```
/patients                    # Coleção de pacientes (plural)
/patient/{id}               # Recurso paciente específico (singular)
/patient/{id}/sessions      # Sub-recursos: sessões do paciente (plural)
/session/{id}               # Recurso sessão específica (singular)
/session/{id}/observations  # Sub-recursos: observações da sessão (plural)
/session/{id}/interventions # Sub-recursos: intervenções da sessão (plural)
```

---

## Convenções Específicas

### Pacientes

| Rota | Método | Descrição | Exemplo |
|------|--------|-----------|---------|
| `/patients` | GET | Listar todos pacientes | `GET /patients` |
| `/patients` | POST | Criar novo paciente | `POST /patients` |
| `/patient/{id}` | GET | Detalhes do paciente | `GET /patient/abc123` |
| `/patient/{id}/sessions` | GET | Listar sessões do paciente | `GET /patient/abc123/sessions` |
| `/patient/{id}/sessions/new` | GET | Formulário nova sessão | `GET /patient/abc123/sessions/new` |
| `/patient/{id}/history` | GET | Timeline do paciente | `GET /patient/abc123/history` |

### Sessões

| Rota | Método | Descrição | Exemplo |
|------|--------|-----------|---------|
| `/session` | POST | Criar nova sessão | `POST /session` |
| `/session/{id}` | GET | Detalhes da sessão | `GET /session/xyz789` |
| `/session/{id}/edit` | GET | Formulário editar sessão | `GET /session/xyz789/edit` |
| `/session/{id}/update` | POST | Atualizar sessão | `POST /session/xyz789/update` |
| `/session/{id}/observations` | POST | Criar observação | `POST /session/xyz789/observations` |
| `/session/{id}/interventions` | POST | Criar intervenção | `POST /session/xyz789/interventions` |

### Observações e Intervenções

| Rota | Método | Descrição | Exemplo |
|------|--------|-----------|---------|
| `/observations/{id}` | GET | Detalhes da observação | `GET /observations/obs123` |
| `/observations/{id}/edit` | GET | Formulário editar observação | `GET /observations/obs123/edit` |
| `/observations/{id}` | PUT | Atualizar observação | `PUT /observations/obs123` |
| `/interventions/{id}` | GET | Detalhes da intervenção | `GET /interventions/int456` |
| `/interventions/{id}/edit` | GET | Formulário editar intervenção | `GET /interventions/int456/edit` |
| `/interventions/{id}` | PUT | Atualizar intervenção | `PUT /interventions/int456` |

---

## Regras de Implementação

### 1. URLs em Templates

Sempre use `templ.URL()` para gerar URLs em templates:

```go
// ✅ CORRETO
<a href={ templ.URL("/patient/" + patientID) }>Ver paciente</a>

// ❌ ERRADO  
<a href="/patient/{patientID}">Ver paciente</a>
```

### 2. Extração de IDs em Handlers

Use funções auxiliares consistentes para extrair IDs:

```go
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

### 3. Verificação de Método HTTP

Sempre verifique o método HTTP no handler:

```go
func (h *Handler) HandlePatient(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        h.showPatient(w, r)
    case "POST":
        h.createPatient(w, r)
    case "PUT":
        h.updatePatient(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

---

## Exceções e Casos Especiais

### 1. Dashboard

| Rota | Método | Descrição |
|------|--------|-----------|
| `/` | GET | Redireciona para dashboard |
| `/dashboard` | GET | Página inicial do sistema |

### 2. Arquivos Estáticos

| Rota | Descrição |
|------|-----------|
| `/static/` | Arquivos CSS, JS, imagens |

### 3. Rotas Especiais

| Rota | Método | Descrição |
|------|--------|-----------|
| `/patients/new` | GET | Formulário novo paciente |
| `/patient/create` | POST | Criar paciente (legado - migrar para POST /patients) |

---

## Validação e Manutenção

### Scripts de Verificação

Use os scripts do projeto para validar consistência:

```bash
./scripts/arandu_validate_handlers.sh  # Valida handlers
./scripts/arandu_checkpoint.sh         # Checkpoint arquitetural
./scripts/arandu_guard.sh              # Verifica integridade
```

### Atualização de Contexto

Após mudanças significativas nas rotas, atualize o contexto:

```bash
./scripts/arandu_update_context.sh
```

---

## Histórico de Mudanças

| Data | Mudança | Justificativa |
|------|---------|---------------|
| 2026-03-17 | Padronização singular/plural | Consistência RESTful |
| 2026-03-17 | Criação deste documento | Documentação obrigatória |

---

**Este documento é de leitura obrigatória para todas as implementações que envolvam rotas web no projeto Arandu.**