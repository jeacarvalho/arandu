---
name: privacy-compliance
description: >
  Conformidade de privacidade e anonimização para plataforma de psicologia com dados
  sensíveis de saúde. Use esta skill SEMPRE que o usuário trabalhar com: dados de
  pacientes, chamadas a LLMs externos, prompts de IA, migrations de banco de dados,
  handlers que retornam dados de pacientes, código que toca SessionForm/MedicalRecord/
  Patient, ou qualquer feature que envolva processar informações clínicas. Também
  dispara para: "como anonimizar X?", "posso enviar esse campo para a IA?", "como
  implementar LGPD aqui?", "esse prompt está seguro?", "como fazer audit trail?".
  Esta skill é um BLOQUEADOR — nenhum código ou prompt que toque dados de pacientes
  deve ser gerado sem passar por estas verificações.
---

# Privacidade e Conformidade — Plataforma de Psicologia

Dados de saúde são **dados sensíveis** sob a LGPD (Art. 11). Esta skill define as
regras inegociáveis de privacidade para toda a plataforma.

**Princípio central**: dados que identificam um paciente **nunca saem do sistema**
para serviços externos. A IA trabalha apenas com dados anonimizados, estruturados
e cuidadosamente selecionados.

---

## Classificação de dados

### Tier 1-Plus — Dados de Identidade (proteção máxima — NUNCA processados por IA, NUNCA em logs)
```
patient_context.ethnicity
patient_context.gender_identity
patient_context.sexual_orientation
patient_context.occupation        ← pode revelar identidade indiretamente
patient_context.education_level
```
> Estes dados têm proteção além dos dados de saúde comuns. Discriminação por identidade
> de gênero, orientação sexual ou etnia é risco real. Tratamento só com base legal
> explícita e consentimento destacado. Nunca aparecem em logs, audit trails ou payloads.

### Tier 1 — Identificação Direta (NUNCA processado por IA externa)
```
patient.name
patient.cpf
patient.birth_date
patient.phone
patient.email
patient.address.*
legal_guardian.*  (todos os campos)
psychologist.name
psychologist.crp
psychologist.email
psychologist.phone
appointment.location  (endereço físico)
```

### Tier 2 — Texto Livre com PII implícito (NUNCA enviado para IA externa)
```
session_form.observations      ← pode conter nomes, locais, relações
session_form.chief_complaint   ← idem
session_form.homework          ← idem
intake.clinical_history        ← idem
intake.family_history          ← idem
intake.life_context            ← idem
session_form.risk_indicators   ← CRÍTICO: nunca, sem exceção
```
> Texto livre é impossível de anonimizar de forma confiável com regex ou substituição
> simples. A única proteção segura é a exclusão do payload.

### Tier 3 — Dados Estruturados Anonimizáveis (permitido com anonimização)
```
session_form.themes[]           ← tags controladas
session_form.emotional_state    ← escala numérica + categoria
session_form.techniques_used[]  ← lista controlada
session_form.next_session_focus ← texto curto e estruturado*
intake.therapeutic_goals        ← categorias, não texto livre*
appointment.date                ← data sem hora exata
session count                   ← número de sessões
psychologist.therapeutic_approach ← abordagem terapêutica
```
> * Atenção: se o psicólogo digitou texto livre nesses campos, aplicar verificação
> adicional antes de incluir no payload.

### Tier 4 — Metadados Operacionais (uso interno — nunca enviado externamente)
```
UUIDs de entidades
timestamps de criação/atualização
audit logs
```

---

## Isolamento físico por tenant — garantia arquitetural

O Arandu usa **database-per-tenant**: cada psicólogo tem seu próprio arquivo `.db`.
Isso garante isolamento físico — não há risco de query cross-tenant por SQL injection
ou bug de filtro. É a primeira linha de defesa de privacidade do sistema.

**Implicações práticas:**
- Um bug de autorização pode expor dados de um paciente ao próprio psicólogo errado —
  mas nunca a outro psicólogo (bancos separados)
- Backup, exportação e exclusão de dados são operações por arquivo — triviais e auditáveis
- O Control Plane (banco central) nunca contém dados clínicos — só metadados de acesso

---

## Pipeline de anonimização

Todo payload destinado a LLM externo deve passar por este pipeline:

### Passo 1 — Seleção de campos
```go
// Apenas campos do Tier 3 são candidatos
// Lista explícita — nunca usar "todos os campos exceto..."
type AnonymizedSessionPayload struct {
    Themes          []string `json:"themes"`
    EmotionalState  int      `json:"emotional_state"`   // só escala numérica
    EmotionalLabel  string   `json:"emotional_label"`   // categoria, não texto livre
    TechniquesUsed  []string `json:"techniques_used"`
    NextFocus       string   `json:"next_session_focus,omitempty"`
    SessionSequence int      `json:"session_number"`    // número ordinal, não data
}
```

### Passo 2 — Token opaco rotativo
```go
// Nunca usar UUID real do paciente como identificador em payloads de IA
// Gerar token de sessão de análise — descartado após uso

type AnalysisToken struct {
    Token     string    // gerado com crypto/rand
    PatientID uuid.UUID // mapeamento local — nunca sai do sistema
    ExpiresAt time.Time // token expira em 1h
}

func NewAnalysisToken(patientID uuid.UUID) AnalysisToken {
    b := make([]byte, 32)
    rand.Read(b)
    return AnalysisToken{
        Token:     base64.URLEncoding.EncodeToString(b),
        PatientID: patientID,
        ExpiresAt: time.Now().Add(time.Hour),
    }
}
```

### Passo 3 — Verificação de risco
```go
// Bloqueia envio se houver RiskIndicator nas últimas N sessões
func (s *AIService) canProcessPatient(ctx context.Context, patientID uuid.UUID) error {
    hasRisk, err := s.sessions.HasRecentRiskIndicators(ctx, patientID, 3)
    if err != nil {
        return fmt.Errorf("risk check failed: %w", err)
    }
    if hasRisk {
        return ErrRiskIndicatorPresent
    }
    return nil
}
```

### Passo 4 — Verificação de menores
```go
// Menores têm proteção adicional — requer consentimento explícito do responsável
func (s *AIService) checkMinorConsent(ctx context.Context, patientID uuid.UUID) error {
    patient, err := s.patients.FindByID(ctx, patientID)
    if err != nil {
        return err
    }
    if patient.IsMinor() {
        consent, err := s.consents.FindGuardianAIConsent(ctx, patientID)
        if err != nil || !consent.IsValid() {
            return ErrMinorWithoutGuardianConsent
        }
    }
    return nil
}
```

### Passo 5 — Sanitização de texto
```go
// Para campos next_session_focus e therapeutic_goals (texto curto mas livre)
// Verificar comprimento e ausência de padrões de PII
func sanitizeShortText(text string) (string, error) {
    if len(text) > 200 {
        return "", ErrTextTooLong // texto longo → risco alto de PII
    }
    if containsPIIPatterns(text) {
        return "", ErrPIIDetected
    }
    return text, nil
}

var piiPatterns = []*regexp.Regexp{
    regexp.MustCompile(`\d{3}\.\d{3}\.\d{3}-\d{2}`), // CPF
    regexp.MustCompile(`\d{2}/\d{2}/\d{4}`),           // data de nascimento
    // emails, telefones...
}
```

### Passo 6 — Audit log obrigatório
```go
// Todo envio a LLM externo deve ser logado ANTES do envio
type AICallLog struct {
    ID              uuid.UUID
    PsychologistID  uuid.UUID
    AnalysisToken   string        // token opaco — não o patient ID
    Feature         string        // "pattern_detection" | "journey_suggestion" | etc.
    FieldsSent      []string      // lista de campos incluídos no payload
    ModelUsed       string        // ex: "claude-sonnet-4-6"
    RequestHash     string        // SHA256 do payload — para auditoria sem reexpor dados
    ResponseHash    string        // SHA256 da resposta
    CalledAt        time.Time
    ConsentVersion  string        // versão do TCLE vigente
}
```

---

## Checklist — antes de qualquer chamada a LLM

Use este checklist ao revisar ou escrever código que chama IA:

```
DADOS
[ ] O payload usa apenas campos do Tier 3?
[ ] Nenhum campo do Tier 1 está presente?
[ ] Nenhum campo do Tier 2 (texto livre) está presente?
[ ] risk_indicators está ausente do payload?
[ ] O identificador do paciente é um token opaco (não UUID real)?

CONSENTIMENTO
[ ] Psicólogo optou explicitamente por esta análise?
[ ] Se menor: consentimento do responsável para IA está registrado?
[ ] TCLE vigente cobre uso de IA para este tipo de análise?

RISCO
[ ] Verificação de RiskIndicators recentes foi executada?
[ ] Bloqueio está implementado quando risco presente?

AUDITORIA
[ ] AICallLog é gravado antes do envio?
[ ] Hash do payload é registrado?
[ ] Versão do consentimento é registrada?

OUTPUT
[ ] Resultado é apresentado como sugestão/hipótese?
[ ] Disclaimer obrigatório está presente na UI?
[ ] Output é salvo com referência ao AICallLog?
```

---

## Construção de prompts seguros

### Template de system prompt para features de IA

```
Você é um assistente de apoio clínico para psicólogos. Você analisa dados
estruturados e anonimizados de sessões terapêuticas para identificar padrões
e sugerir reflexões clínicas.

REGRAS INEGOCIÁVEIS:
1. Nunca afirme diagnósticos. Use linguagem de hipótese: "pode indicar",
   "sugere explorar", "padrão observado".
2. Nunca solicite ou infira dados pessoais do paciente. Os dados são
   intencionalmente anonimizados.
3. Sempre inclua um disclaimer indicando que a análise é uma ferramenta
   de apoio e não substitui o julgamento clínico do profissional.
4. Se os dados forem insuficientes para uma análise confiável, diga isso
   explicitamente em vez de especular.
5. Respeite a abordagem terapêutica informada — não sugira técnicas
   incompatíveis com ela sem justificativa.

CONTEXTO:
- Abordagem terapêutica do profissional: {therapeutic_approach}
- Número de sessões analisadas: {session_count}
- Período coberto: {period_weeks} semanas

DADOS (anonimizados e estruturados):
{anonymized_payload_json}

Responda em português brasileiro.
```

### Verificação de prompt antes de usar em produção

```go
func validatePrompt(prompt string, payload AIPayload) error {
    // Verificar se nenhum dado Tier 1 foi interpolado
    for _, pii := range []string{
        payload.PatientName,    // deve ser vazio
        payload.PatientCPF,     // deve ser vazio
        payload.PatientEmail,   // deve ser vazio
    } {
        if pii != "" && strings.Contains(prompt, pii) {
            return ErrPIIInPrompt
        }
    }
    return nil
}
```

---

## Implementação de Audit Trail

```go
// Toda operação sensível deve ser logada — não só IA
type AuditEvent struct {
    ID             uuid.UUID
    ActorID        uuid.UUID  // quem fez a ação (psychologist)
    Action         string     // "view_record" | "ai_analysis" | "export_data"
    ResourceType   string     // "MedicalRecord" | "Session" | "Patient"
    ResourceID     uuid.UUID
    IPAddress      string
    OccurredAt     time.Time
    // Nunca logar o conteúdo — só metadados
}

// Eventos que OBRIGATORIAMENTE geram AuditEvent:
// - Abertura de prontuário
// - Criação/edição de SessionForm
// - Solicitação de análise de IA
// - Exportação de dados
// - Acesso por suporte/admin da plataforma
```

---

## Retenção e exclusão

```go
// Política de retenção — implementar como job periódico
type RetentionPolicy struct {
    // Adultos: 5 anos após último atendimento
    AdultRetentionYears int // = 5

    // Menores: 5 anos após completar 18 anos
    MinorRetentionUntil func(birthDate time.Time) time.Time

    // Dados de conta (não clínicos): 2 anos após encerramento
    AccountDataRetentionYears int // = 2
}

// "Exclusão" de paciente = anonimização, não delete
// O prontuário permanece (obrigação legal), mas dados identificadores são
// substituídos por placeholders: "[DADOS REMOVIDOS - LGPD Art. 18]"
func anonymizePatientIdentity(db *sql.DB, patientID uuid.UUID) error {
    _, err := db.Exec(`
        UPDATE patients SET
            name = '[REMOVIDO]',
            cpf = NULL,
            phone = NULL,
            email = NULL,
            address = NULL,
            updated_at = ?
        WHERE id = ?
    `, time.Now(), patientID)
    return err
}
```

---

## Referências
- `references/lgpd-checklist.md` — Checklist LGPD completo para features
- `references/incident-response.md` — Procedimento em caso de vazamento de dados
