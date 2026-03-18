Task: Implementação do REQ-01-04-01 — Histórico Farmacológico e Sinais Vitais

ID da Tarefa: task_20260318_biopsychosocial_context

Status: ✅ CONCLUÍDO

Requisito Relacionado: REQ-01-04-01

Stack Técnica: Go, templ, HTMX, SQLite.

🎯 Objetivo

Implementar a capacidade de registar e monitorizar o contexto biopsicossocial do paciente, focando no histórico de medicação e em indicadores fisiológicos (sinais vitais). A interface deve ser implementada como um "Painel de Contexto" lateral, seguindo a estética de Tecnologia Silenciosa.

Comece executando todos os testes de regressão do sistema. 

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/patient)

Entidade Medication: ID, PatientID, Name, Dosage, Frequency, Status (Active/Suspended/Finished), StartedAt, EndedAt.

Entidade Vitals: ID, PatientID, Date, SleepHours, AppetiteLevel, Weight.

Validação: Garantir que as datas de início de medicação não sejam futuras e que os níveis de apetite estejam entre 1 e 10.

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)
REQ-01-04-01
Migração SQL: Criar o ficheiro 0003_add_biopsychosocial_tables.up.sql com as tabelas patient_medications e patient_vitals.

Repositório:

Implementar SaveMedication, UpdateMedicationStatus e GetActiveMedications.

Implementar SaveVitals e GetLatestVitals.

3. Camada de Aplicação (internal/application/services)

Criar BiopsychosocialService para gerir a lógica de reconciliação medicamentosa e registo de vitais.

4. Camada Web (Componentes templ)

BiopsychosocialPanel: Contentor lateral que organiza as secções.

MedicationWidget:

Lista de medicamentos ativos.

Botão rápido para suspender (hx-put) que risca o nome do medicamento instantaneamente.

Formulário minimalista ("Silent Input") para adicionar nova medicação.

VitalsWidget:

Campos de entrada rápida para Sono e Peso.

Exibição dos últimos valores registados.

5. Camada Web (Handlers)

GET /patients/{id}/context: Retorna o fragmento do painel completo.

POST /patients/{id}/medications: Adiciona e retorna a lista atualizada.

PUT /patients/{id}/medications/{med_id}/status: Atualiza status e retorna o item formatado.

🎨 Design System (Silent UI)

Tipografia:

Labels técnicos: Inter (Sans) 12px.

Nomes de medicamentos: Source Serif 4 18px.

Visual:

Medicamentos ativos: Fundo bg-green-50/50 sutil.

Medicamentos suspensos: Opacidade reduzida e line-through.

Interação: O painel deve abrir/fechar lateralmente sem deslocar o conteúdo principal do prontuário (uso de fixed ou sticky lateral).

🧪 Protocolo de Testes "Ironclad"

A. Teste de Persistência

Validar que ao suspender um medicamento, o campo ended_at é preenchido automaticamente com a data atual.

Validar que os sinais vitais são filtrados corretamente por patient_id.

B. Teste E2E (Playwright)

Abrir o perfil de um paciente.

Expandir o painel de "Contexto Biológico".

Adicionar "Sertralina, 50mg, Manhã".

Validar: O item aparece na lista com fonte Serif.

Clicar em "Suspender".

Validar: O item fica cinzento e riscado sem recarga de página.

Registar "8h de sono".

Validar: O valor é atualizado no resumo de vitais.

🛡️ Checklist de Integridade

[x] As novas tabelas possuem ON DELETE CASCADE para o paciente?
[x] O componente BiopsychosocialPanel é responsivo (esconde em mobile ou vira modal)?
[x] Foi utilizado o padrão RawHTML para qualquer renderização dinâmica necessária?
[x] O scripts/arandu_guard.sh confirma a integridade das novas rotas?
[x] Executei templ generate?

## 📋 Notas de Implementação

### Arquivos Criados

1. **Domínio:**
   - `internal/domain/patient/medication.go` - Entidade Medication
   - `internal/domain/patient/vitals.go` - Entidade Vitals
   - `internal/domain/patient/medication_test.go` - Testes Medication
   - `internal/domain/patient/vitals_test.go` - Testes Vitals

2. **Infraestrutura:**
   - `internal/infrastructure/repository/sqlite/migrations/0005_biopsychosocial_tables.up.sql`
   - `internal/infrastructure/repository/sqlite/migrations/0005_biopsychosocial_tables.down.sql`
   - `internal/infrastructure/repository/sqlite/medication_repository.go`
   - `internal/infrastructure/repository/sqlite/vitals_repository.go`

3. **Aplicação:**
   - `internal/application/services/biopsychosocial_service.go`

4. **Web:**
   - `internal/web/handlers/biopsychosocial_handler.go`
   - `web/components/patient/biopsychosocial_panel.templ`
   - `web/components/patient/medication_list.templ`
   - `web/components/patient/vitals_widget.templ`
   - `web/components/patient/biopsychosocial_viewmodel.go`

5. **Integração:**
   - `cmd/arandu/main.go` - Rotas e inicialização

### Rotas Implementadas

- `GET /patients/{id}/context` - Painel de contexto biopsicossocial
- `POST /patients/{id}/medications` - Adicionar medicação
- `PUT /patients/{id}/medications/{med_id}/status?status=suspend|activate` - Atualizar status
- `POST /patients/{id}/vitals` - Registrar sinais vitais

### Testes

- Todos os testes de domínio passam (14 testes)
- Migração 0005 aplicada com sucesso
- Build completo sem erros

