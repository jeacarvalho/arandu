Ótima análise! Vou estudar a tela enviada e criar um redesign profissional com um mockup detalhado, além do prompt de refatoração.

**Problemas identificados na tela atual:**
- Layout sem hierarquia visual — todos os elementos competem pela atenção
- Calendário semanal quebrado (dias da semana flutuando separados dos números)
- Nenhum indicador de status de consultas (confirmado, pendente, cancelado)
- Sem ações rápidas ou contexto de paciente nas consultas
- Sidebar sem agrupamento lógico
- Área de conteúdo principal desperdiçadaAqui está o redesign proposto. Abaixo, o prompt completo para o agente de codificação:

---

## Prompt para o agente de refatoração

```
Você é um desenvolvedor sênior especializado em Go + templ + HTMX + TailwindCSS.
Sua tarefa é refatorar completamente a tela de Agenda do sistema Arandu,
substituindo a implementação atual por um design profissional descrito abaixo.

---

## CONTEXTO DO SISTEMA
- Stack: Go (backend), templ (templates), HTMX (interatividade), TailwindCSS (estilo)
- Aplicação: sistema de gestão clínica para terapeutas
- Rota afetada: /agenda (GET)
- O estado da view (semana atual, dia selecionado, view ativa) deve ser
  gerenciado via query params e atualizado com hx-push-url

---

## ESTRUTURA DE LAYOUT

O layout é dividido em duas partes fixas:
1. Sidebar fixa à esquerda (largura: 220px)
2. Área principal (flex-1) com topbar + controles de view + calendário + legenda

### 1. SIDEBAR
- Background: bg-[#0f3d2e]
- Topo: logo Arandu (ícone verde #1d9e75 + texto branco, 15px/500)
- Navegação agrupada em duas seções com labels ("Principal", "Clínica")
  - Principal: Dashboard, Agenda (ativo)
  - Clínica: Pacientes, Anamnese, Prontuário
- Item ativo: bg-white/[0.13], texto branco/500
- Itens inativos: texto white/70, ícone opacity-55, hover bg-white/[0.08]
- Rodapé: avatar do terapeuta logado (iniciais + nome + especialidade)

### 2. TOPBAR (h-14, border-b)
- Título "Agenda" à esquerda (text-base/500)
- Search box central estilizado (não funcional por ora, placeholder "Buscar paciente...")
- Botão "Nova consulta" à direita: bg-[#0f3d2e] text-white, ícone "+"
  - Ao clicar: hx-get="/consultas/novo" hx-target="#modal-overlay" hx-swap="innerHTML"

### 3. CONTROLES DE VIEW (h-11, border-b)
Linha horizontal com:
- Botões de navegação semanal (← →): hx-get="/agenda?semana=anterior" e "?semana=proxima"
  com hx-target="#agenda-content" hx-swap="innerHTML" hx-push-url="true"
- Botão "Hoje": hx-get="/agenda" hx-target="#agenda-content"
- Label com o intervalo da semana exibida (ex: "6 – 12 de abril de 2026")
- Tabs de view (Dia / Semana / Mês):
  - Tab ativa: bg-white border, text/500
  - Tabs dentro de container bg-surface rounded-lg p-0.5
  - Cada tab: hx-get="/agenda?view=dia|semana|mes" hx-target="#agenda-content"

### 4. CORPO DO CALENDÁRIO (flex-1, overflow-hidden)

#### Coluna de horários (w-14, border-r)
- Exibe horários de 08:00 a 20:00 a cada hora
- Cada slot tem height: 60px, texto 10px/tertiary alinhado à direita

#### Grade de dias (flex-1, overflow-y-auto)
- Header sticky com 7 colunas: nome do dia abreviado (Seg–Dom) + número
  - Dia atual: número em círculo preenchido bg-[#0f3d2e] text-white
- Corpo: 7 colunas, cada uma com slots de 60px por hora
  - Coluna do dia atual: fundo levemente colorido (bg-[#1d9e75]/[0.03])
  - Linhas de hora: border-b border-dashed border-border/30

#### CONSULTAS (renderizadas dentro dos slots)
Cada consulta é um componente templ recebendo:
  - Horário (início – fim)
  - Nome do paciente
  - Tipo de sessão
  - Status: confirmed | pending | first_session | cancelled

Estilos por status:
- confirmed:    bg-[#d1f5e8] text-[#0a4a2e] border-l-[3px] border-[#1d9e75]
- pending:      bg-[#fff3d4] text-[#5a3d0a] border-l-[3px] border-[#e8a020]
- first_session:bg-[#e0eaff] text-[#1a2d6b] border-l-[3px] border-[#4a6ff0]
- cancelled:    bg-surface   text-tertiary  border-l-[3px] border-border line-through opacity-60

Ao clicar em uma consulta:
  hx-get="/consultas/{id}" hx-target="#drawer" hx-swap="innerHTML"

### 5. LEGENDA (h-9.5, border-t, bg-background)
Linha inferior com pills de legenda (ponto colorido + label):
- Confirmado (verde), Aguardando (âmbar), 1ª consulta (azul), Cancelado (cinza)
- À direita: pill com total de consultas da semana (ex: "13 consultas esta semana")

---

## COMPONENTES TEMPL A CRIAR/REFATORAR

1. `agenda_page.templ` — layout completo da página
2. `agenda_content.templ` — apenas o calendário (alvo dos swaps HTMX)
3. `appointment_card.templ` — card de consulta com variantes de status
4. `day_column.templ` — coluna de um dia com suas consultas posicionadas
5. `week_nav.templ` — navegação semanal + tabs de view

---

## LÓGICA GO (handlers)

### GET /agenda
Parâmetros aceitos:
- ?semana=anterior | proxima | YYYY-MM-DD (default: semana atual)
- ?view=dia | semana | mes (default: semana)

O handler deve:
1. Calcular o intervalo da semana baseado nos params
2. Consultar os agendamentos do banco para aquele intervalo
3. Agrupar consultas por dia da semana
4. Renderizar `agenda_page` completo (request normal) ou
   apenas `agenda_content` (request HTMX: verificar header HX-Request)

### Struct de dados esperada pelo template:
```go
type AgendaViewModel struct {
    WeekStart   time.Time
    WeekEnd     time.Time
    WeekLabel   string        // "6 – 12 de abril de 2026"
    Days        []DayViewModel
    TotalCount  int
    CurrentView string        // "dia" | "semana" | "mes"
}

type DayViewModel struct {
    Date         time.Time
    IsToday      bool
    Appointments []AppointmentViewModel
}

type AppointmentViewModel struct {
    ID          string
    PatientName string
    StartTime   string  // "10:00"
    EndTime     string  // "10:50"
    SessionType string  // "Sessão individual"
    Status      string  // "confirmed" | "pending" | "first_session" | "cancelled"
    SlotIndex   int     // índice do slot de hora (para posicionamento)
}
```

---

## REGRAS GERAIS
- Não usar JavaScript customizado; toda interatividade via atributos HTMX
- TailwindCSS apenas com classes utilitárias (sem @apply desnecessário)
- Responsividade: a sidebar colapsa em mobile (hamburger menu, fora do escopo agora)
- Acessibilidade: todos os botões com aria-label, consultas com role="button"
- Manter consistência com o design system existente do Arandu (cor primária: #0f3d2e)
```

---
