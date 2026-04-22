# Convenções de Rotas - Arandu

**Documento de Referência Arquitetural**  
**Última atualização:** 22/04/2026

---

## Princípios Gerais

### 1. RESTful com Plural Consistente

Todas as rotas de recurso usam **plural**. O mux Go usa prefix matching via trailing slash:

- **Coleções**: `/patients`, `/sessions`
- **Recurso específico**: `/patients/{id}` (não `/patient/{id}`)
- **Sub-recursos**: `/patients/{id}/sessions`

> ⚠️ Não existe `/patient/{id}` (singular). O handler `mux.HandleFunc("/patients/", ...)` captura todas as variantes.

### 2. Hierarquia de Recursos

```
/patients                    # Coleção de pacientes
/patients/{id}               # Recurso paciente específico
/patients/{id}/sessions      # Sub-recursos: sessões do paciente
/session/{id}                # Sessão específica (singular — herança histórica)
/session/{id}/observations   # Sub-recursos: observações da sessão
/session/{id}/interventions  # Sub-recursos: intervenções da sessão
```

---

## Rotas Implementadas

### Dashboard e Sistema

| Rota | Método | Descrição |
|------|--------|-----------|
| `/` | GET | Redireciona para `/dashboard` |
| `/dashboard` | GET | Página inicial do sistema |
| `/static/` | GET | Arquivos CSS, JS, imagens |

### Autenticação

| Rota | Método | Descrição |
|------|--------|-----------|
| `/login` | GET/POST | Login com email/senha |
| `/auth/login` | GET/POST | Alias de `/login` |
| `/auth/google` | GET | Inicia OAuth Google |
| `/auth/google/callback` | GET | Callback OAuth Google |
| `/auth/signup` | GET/POST | Cadastro de novo terapeuta |
| `/logout` | GET | Encerra sessão |

### Pacientes

| Rota | Método | Handler | Descrição |
|------|--------|---------|-----------|
| `/patients` | GET | `PatientHandler.ListPatients` | Listar pacientes |
| `/patients/new` | GET | `PatientHandler.NewPatient` | Formulário novo paciente |
| `/patients/create` | POST | `PatientHandler.CreatePatient` | Criar paciente |
| `/patients/search` | GET | `PatientHandler.Search` | Busca FTS5 de pacientes |
| `/patients/{id}` | GET | `PatientHandler.Show` | Perfil do paciente |
| `/patients/{id}/sessions` | GET | `PatientHandler.ListSessions` | Listar sessões |
| `/patients/{id}/sessions/new` | GET | `SessionHandler.NewSession` | Formulário nova sessão |
| `/patients/{id}/history` | GET | `TimelineHandler.ShowPatientHistory` | Timeline longitudinal |
| `/patients/{id}/history/search` | GET | `TimelineHandler.SearchPatientHistory` | Busca na timeline |
| `/patients/{id}/history/more` | GET | `TimelineHandler.LoadMoreEvents` | Infinite scroll |
| `/patients/{id}/bio-context` | GET | `BiopsychosocialHandler.GetContextPanel` | Painel biopsicossocial |
| `/patients/{id}/medications` | POST | `BiopsychosocialHandler.AddMedication` | Adicionar medicação |
| `/patients/{id}/medications/{mid}/status` | POST | `BiopsychosocialHandler.UpdateMedicationStatus` | Atualizar status |
| `/patients/{id}/vitals` | POST | `BiopsychosocialHandler.RecordVitals` | Registrar sinais vitais |
| `/patients/{id}/anamnesis` | GET | `PatientHandler.ShowAnamnesis` | Exibir anamnese |
| `/patients/{id}/anamnesis/{section}` | PATCH | `PatientHandler.UpdateAnamnesisSection` | Atualizar seção anamnese |
| `/patients/{id}/goals` | POST | `SessionHandler.CreateGoal` | Criar meta terapêutica |
| `/patients/{id}/goals/{gid}/close` | POST | `SessionHandler.CloseGoalWithNote` | Encerrar meta |
| `/patients/{id}/therapeutic-plan-report` | GET | `SessionHandler.TherapeuticPlanReport` | Relatório do plano |
| `/patients/{id}/ai/synthesis` | POST | `AIHandler.GeneratePatientSynthesis` | Síntese por IA |
| `/patients/{id}/themes` | GET | `AnalysisHandler.ShowThemes` | Nuvem de temas |

### Sessões

| Rota | Método | Handler | Descrição |
|------|--------|---------|-----------|
| `/session` | POST | `SessionHandler.CreateSession` | Criar sessão |
| `/session/{id}` | GET | — | **Redirect 301 → `/session/{id}/edit`** |
| `/session/{id}/edit` | GET | `SessionHandler.EditSession` | Registro clínico (design Sábio) |
| `/session/{id}/update` | POST | `SessionHandler.UpdateSession` | Atualizar data/resumo |
| `/session/{id}/observations` | POST | `SessionHandler.CreateObservation` | Criar observação |
| `/session/{id}/interventions` | POST | `SessionHandler.CreateIntervention` | Criar intervenção |
| `/session/{id}/summary` | PATCH | `SessionHandler.PatchSummary` | Auto-save resumo (HTMX) |

### Observações

| Rota | Método | Handler | Descrição |
|------|--------|---------|-----------|
| `/observations/{id}` | GET | `ObservationHandler.GetObservation` | Detalhes |
| `/observations/{id}/edit` | GET | `ObservationHandler.GetObservationEditForm` | Formulário edição inline |
| `/observations/{id}` | PUT | `ObservationHandler.UpdateObservation` | Atualizar |
| `/observations/{id}/classify/edit` | GET | `ClassificationHandler.GetClassificationEdit` | Seletor de tags inline |
| `/observations/{id}/classify` | POST | `ClassificationHandler.ClassifyObservation` | Aplicar tag |
| `/observations/{id}/classify/{tag_id}` | DELETE | `ClassificationHandler.RemoveClassification` | Remover tag |
| `/tags` | GET | `ClassificationHandler.GetTagsByType` | Tags de observações por tipo |

### Intervenções

| Rota | Método | Handler | Descrição |
|------|--------|---------|-----------|
| `/interventions/{id}` | GET | `InterventionHandler.GetIntervention` | Detalhes |
| `/interventions/{id}/edit` | GET | `InterventionHandler.GetInterventionEditForm` | Formulário edição inline |
| `/interventions/{id}` | PUT | `InterventionHandler.UpdateIntervention` | Atualizar |
| `/interventions/{id}/classify/edit` | GET | `InterventionClassificationHandler.GetInterventionClassificationEdit` | Seletor de tags inline |
| `/interventions/{id}/classify` | POST | `InterventionClassificationHandler.ClassifyIntervention` | Aplicar tag |
| `/interventions/{id}/classify/{tag_id}` | DELETE | `InterventionClassificationHandler.RemoveInterventionClassification` | Remover tag |
| `/tags/interventions` | GET | `InterventionClassificationHandler.GetInterventionTagsByType` | Tags de intervenções por tipo |

### Agenda

| Rota | Método | Handler | Descrição |
|------|--------|---------|-----------|
| `/agenda` | GET | `AgendaHandler.View` | Visualização principal |
| `/agenda/day` | GET | `AgendaHandler.DayView` | View diária |
| `/agenda/week` | GET | `AgendaHandler.WeekView` | View semanal |
| `/agenda/month` | GET | `AgendaHandler.MonthView` | View mensal |
| `/agenda/new` | GET | `AgendaHandler.NewForm` | Formulário novo compromisso |
| `/agenda/slots` | GET | `AgendaHandler.GetSlots` | Slots disponíveis |
| `/agenda/appointments` | POST | `AgendaHandler.Create` | Criar compromisso |
| `/agenda/appointments/{id}` | GET | `AgendaHandler.Show` | Detalhe do compromisso |
| `/agenda/appointments/{id}` | PUT | `AgendaHandler.Update` | Atualizar compromisso |
| `/agenda/appointments/{id}` | DELETE | `AgendaHandler.Cancel` | Cancelar compromisso |
| `/agenda/appointments/{id}/reschedule` | POST | `AgendaHandler.Reschedule` | Reagendar |
| `/agenda/appointments/{id}/complete` | POST | `AgendaHandler.Complete` | Marcar como realizada |

---

## Regras de Implementação

### 1. URLs em Templates

```go
// ✅ CORRETO
<a href={ templ.URL("/patients/" + patientID) }>Ver paciente</a>

// ❌ ERRADO — sem templ.URL()
<a href="/patients/{patientID}">Ver paciente</a>

// ❌ ERRADO — singular
<a href={ templ.URL("/patient/" + patientID) }>Ver paciente</a>
```

### 2. Extração de IDs em Handlers

```go
// Use r.PathValue() no Go 1.22+ ou strings.Split para extrair segmentos
patientID := strings.TrimPrefix(r.URL.Path, "/patients/")
patientID = strings.Split(patientID, "/")[0]
```

### 3. Verificação de Método HTTP

```go
switch r.Method {
case http.MethodGet:
    h.show(w, r)
case http.MethodPost:
    h.create(w, r)
default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
```

### 4. HTMX: Fragmento vs Página Completa

```go
if r.Header.Get("HX-Request") == "true" {
    // retorna apenas o componente
    component.Render(r.Context(), w)
    return
}
// retorna com layout completo
layout.Shell(config, component).Render(r.Context(), w)
```

---

## Histórico de Mudanças

| Data | Mudança | Justificativa |
|------|---------|---------------|
| 2026-03-17 | Criação do documento | Documentação obrigatória |
| 2026-04-22 | Adição de todas as rotas implementadas | Sync com main.go real |

---

**Este documento é de leitura obrigatória para todas as implementações que envolvam rotas web no projeto Arandu.**
