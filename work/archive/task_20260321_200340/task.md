# TASK 20260321_200340

Title: Integrar timeline dinâmica na página de detalhes do paciente

## Status
**PRONTO_PARA_IMPLEMENTACAO**

## Objetivo

Substituir a seção "Linha do Tempo" estática/mock na página de detalhes do paciente por uma timeline dinâmica que busca os 5 eventos mais recentes do banco de dados.

## Descrição Detalhada

A página de detalhes do paciente (`/patients/{id}`) atualmente exibe uma seção "Linha do Tempo" com dados estáticos/mock (datas fixas: "15 Jan 2024", "22 Jan 2024", etc.). Esta tarefa visa integrar a timeline real do sistema.

### Requisitos Funcionais

1. **Buscar eventos do banco**: No handler `PatientHandler.Show`, buscar os 5 eventos mais recentes da timeline para o paciente atual usando o `TimelineService`

2. **Exibir eventos reais**: No template `detail.templ`, substituir os dados estáticos por eventos reais da timeline
   - Mostrar: tipo de evento (sessão, observação, intervenção), data, conteúdo/resumo
   - Manter o padrão visual atual (timeline-container, timeline-event, etc.)
   - Se não houver eventos, mostrar mensagem "Nenhum evento registrado"

3. **Adicionar link para histórico completo**: Incluir botão/link "Ver histórico completo" que redireciona para `/patients/{id}/history`

4. **Tratamento de erros**: Se houver erro ao buscar eventos, mostrar mensagem amigável e manter o resto da página funcional

### Requisitos Técnicos

1. **Injetar TimelineService**: Adicionar `timelineService TimelineService` ao `PatientHandler` e atualizar `NewPatientHandler` e `main.go`

2. **Criar view model**: Criar `TimelineEventViewModel` no pacote `patient` com campos necessários para o template

3. **Mapear eventos**: Converter `timeline.TimelineEvent` para `TimelineEventViewModel` com formatação de data amigável

4. **Tipografia**: Manter uso de `var(--font-clinical)` (Source Serif 4) para o conteúdo clínico

### Alterações Necessárias

**Arquivos Go:**
- `internal/web/handlers/patient_handler.go` - Adicionar TimelineService e buscar eventos no método Show
- `cmd/arandu/main.go` - Injetar TimelineService no PatientHandler

**Templates:**
- `web/components/patient/detail.templ` - Atualizar seção da timeline para usar eventos reais

### Critérios de Aceite

- [ ] A timeline mostra os 5 eventos mais recentes do paciente
- [ ] Eventos são ordenados cronologicamente (mais recentes primeiro)
- [ ] Cada evento mostra: ícone por tipo, data formatada, conteúdo/resumo
- [ ] Botão "Ver histórico completo" leva a `/patients/{id}/history`
- [ ] Se não houver eventos, exibe mensagem "Nenhum evento registrado"
- [ ] Página funciona normalmente mesmo se houver erro na timeline
- [ ] Design mantém consistência visual com o resto da aplicação
- [ ] Nenhuma regressão nas outras funcionalidades da página

## Checklist de Integridade (OBRIGATÓRIO)
- [ ] O componente usa .templ e herda de Layout?
- [ ] A tipografia Source Serif 4 foi aplicada ao conteúdo clínico?
- [ ] Executei 'templ generate' e o código Go compilou?
- [ ] Testei a rota atual e as rotas vizinhas (Regressão)?
- [ ] O banco de dados foi atualizado via migration .up.sql?
