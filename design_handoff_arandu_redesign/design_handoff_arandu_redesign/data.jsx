// Dados mockados para a proposta
const ARANDU_DATA = {
  clinician: {
    name: "Dra. Helena Moraes",
    role: "Psicóloga Clínica · CRP 06/98432",
    initials: "HM",
  },
  kpis: [
    { id: "sessions", label: "Sessões registradas", value: 618, delta: "+24 esta semana", tone: "neutral" },
    { id: "patients", label: "Pacientes ativos", value: 42, delta: "3 novos no mês", tone: "up" },
    { id: "today",    label: "Hoje",              value: 5,  delta: "próxima às 14h", tone: "neutral" },
    { id: "pending",  label: "Anotações pendentes", value: 2, delta: "de ontem",     tone: "warn" },
  ],
  todaySchedule: [
    { time: "09:00", patient: "Amanda Rocha",  type: "Retorno", status: "done" },
    { time: "10:30", patient: "André Barbosa", type: "Acompanhamento", status: "done" },
    { time: "14:00", patient: "Carolina Costa", type: "Primeira consulta", status: "next" },
    { time: "15:30", patient: "Felipe Lima",   type: "Retorno", status: "upcoming" },
    { time: "17:00", patient: "Gustavo Pereira", type: "Acompanhamento", status: "upcoming" },
  ],
  patients: [
    { id: "p0446", name: "André Barbosa", since: "dez/2023", tag: "Burnout", tagTone: "warn", last: "há 2 dias", next: "qua, 10:30", sessions: 5, risk: "moderado" },
    { id: "p0312", name: "Amanda Rocha", since: "mar/2024", tag: "Ansiedade", tagTone: "neutral", last: "hoje", next: "ter, 09:00", sessions: 14, risk: "baixo" },
    { id: "p0188", name: "Carolina Costa", since: "abr/2026", tag: "Triagem", tagTone: "info", last: "—", next: "hoje, 14:00", sessions: 0, risk: "a avaliar" },
    { id: "p0077", name: "Felipe Lima", since: "ago/2025", tag: "Luto", tagTone: "warn", last: "há 7 dias", next: "hoje, 15:30", sessions: 11, risk: "moderado" },
    { id: "p0054", name: "Gustavo Pereira", since: "jan/2025", tag: "TOC", tagTone: "neutral", last: "há 3 dias", next: "hoje, 17:00", sessions: 22, risk: "baixo" },
    { id: "p0029", name: "Carlos Santos", since: "nov/2024", tag: "Depressão", tagTone: "danger", last: "há 1 dia", next: "sex, 11:00", sessions: 31, risk: "atenção" },
  ],
  recentSessions: [
    { id: "s0446-148", patient: "André Barbosa", date: "02/04/2026", preview: "Questionamento existencial intenso sobre sentido da vida…", theme: "Sentido" },
    { id: "s0312-091", patient: "Amanda Rocha",  date: "02/04/2026", preview: "Padrões de pensamento automático revisitados.", theme: "Cognição" },
    { id: "s0077-044", patient: "Felipe Lima",   date: "01/04/2026", preview: "Conversa sobre aniversário de perda. Afeto congruente.", theme: "Luto" },
    { id: "s0054-203", patient: "Gustavo Pereira", date: "31/03/2026", preview: "Redução de rituais na semana. Satisfação relatada.", theme: "TOC" },
  ],
  patientProfile: {
    id: "p0446",
    name: "André Barbosa",
    initials: "AB",
    age: 57,
    pronouns: "ele/dele",
    marker: "Masculino · branco · Ensino superior",
    since: "Desde dez/2023",
    sessions: 5,
    therapyDuration: "2,3 anos",
    frequency: "Quinzenal",
    triage: "Paciente de 57 anos, engenheiro sênior. Apresenta Burnout profissional com sintomas somáticos associados (insônia, fadiga, irritabilidade). Início do tratamento buscando melhora da qualidade de vida e reencontro com sentido profissional.",
    recentObservations: [
      { date: "02/04/2026", tag: "luto", text: "Luto persistente identificado; identificação simbiótica com figura paterna dificulta separação psíquica." },
      { date: "26/03/2026", tag: "humor", text: "Episódio hipomaníaco presente, humor elevado com projetos simultâneos. Investigar história de oscilações." },
      { date: "19/03/2026", tag: "social", text: "Isolamento social relatado, evita contatos e cancela compromissos. Padrão recente pós-evento de trabalho." },
      { date: "12/03/2026", tag: "sentido", text: "Questiona escolhas profissionais dos últimos 20 anos. Traz frase: ‘para que serviu tudo isso?’" },
    ],
    quickActions: [
      { id: "anamnese", label: "Anamnese clínica" },
      { id: "plano", label: "Plano terapêutico" },
      { id: "temas", label: "Análise de temas" },
      { id: "nova", label: "Nova sessão", primary: true },
    ],
    timeline: [
      { date: "02/04/2026", kind: "sessão", title: "Sessão 5 · Sentido e existência", summary: "Exploração de questionamentos existenciais e insatisfação com conquistas." },
      { date: "19/03/2026", kind: "nota", title: "Observação entre sessões", summary: "Paciente relatou por mensagem melhora no sono após ajuste de rotina." },
      { date: "12/03/2026", kind: "sessão", title: "Sessão 4 · Histórico profissional", summary: "Retomada de escolhas dos últimos 20 anos. Identificação com figura paterna." },
      { date: "26/02/2026", kind: "sessão", title: "Sessão 3 · Sintomas somáticos", summary: "Mapeamento de insônia, fadiga e gatilhos corporais." },
      { date: "12/02/2026", kind: "sessão", title: "Sessão 2 · Burnout", summary: "Construção de narrativa de esgotamento e expectativa." },
      { date: "29/01/2026", kind: "sessão", title: "Sessão 1 · Acolhimento inicial", summary: "Primeiro contato; queixa principal e contrato terapêutico." },
    ],
  },
  session: {
    id: "s0446-148",
    patient: "André Barbosa",
    initials: "AB",
    number: 5,
    date: "18/10/2025",
    time: "18:00 · 50min",
    observations: [
      { id: "o1", tag: "Observação", timestamp: "02/04/2026 às 20:18", text: "Questionamento existencial intenso sobre sentido da vida e insatisfação com conquistas aparentes." },
      { id: "o2", tag: "Observação", timestamp: "02/04/2026 às 20:18", text: "Insight emergente sobre padrões automáticos. Consciência metacognitiva demonstrada na sessão." },
      { id: "o3", tag: "Observação", timestamp: "02/04/2026 às 20:18", text: "Impulsividade com gastos excessivos e múltiplos projetos simultâneos. Julgamento prejudicado." },
    ],
    interventions: [
      { id: "i1", tag: "Intervenção", timestamp: "02/04/2026 às 20:18", text: "Estabelecimento de limites de disponibilidade entre papéis profissional e pessoal." },
      { id: "i2", tag: "Intervenção", timestamp: "02/04/2026 às 20:18", text: "Exposição in vivo gradual realizada com acompanhamento terapêutico." },
      { id: "i3", tag: "Intervenção", timestamp: "02/04/2026 às 20:18", text: "Escuta do questionamento existencial e validação da busca por sentido." },
    ],
    summary: "Sessão marcada por abertura reflexiva sobre sentido profissional e emergência de material hipomaníaco. Recomenda-se monitorar oscilações de humor entre sessões.",
  },
  llm: {
    threads: [
      {
        id: "t1",
        title: "Padrões cruzados entre sessões de André Barbosa",
        updated: "há 3h",
        messages: [
          { role: "user", text: "O que vem se repetindo nas últimas 5 sessões do André?" },
          { role: "assistant", text: "Três temas recorrentes emergem:\n\n**1. Identificação paterna** — mencionada em 4 das 5 sessões, geralmente em contextos de dúvida profissional.\n\n**2. Oscilação de humor** — picos de entusiasmo (13/02, 26/03) seguidos de retração (19/03, 02/04). O intervalo é curto; vale investigar ritmo hipomaníaco.\n\n**3. Questionamento existencial** — intensificado desde a sessão 3. Sempre aparece após menção ao trabalho.", citations: [
            { session: "S1", date: "29/01" },
            { session: "S3", date: "26/02" },
            { session: "S4", date: "12/03" },
            { session: "S5", date: "02/04" },
          ]},
        ],
      },
    ],
    suggestions: [
      "Resumir evolução dos últimos 3 meses",
      "Comparar temas com outros pacientes em Burnout",
      "Sugerir hipóteses diagnósticas diferenciais",
      "Gerar rascunho de plano terapêutico",
    ],
    crossInsights: [
      { label: "Tema dominante", value: "Sentido / existência", weight: 0.82 },
      { label: "Emergente", value: "Oscilação de humor", weight: 0.64 },
      { label: "Em recuo", value: "Sintomas somáticos", weight: 0.31 },
    ],
  },
};

window.ARANDU_DATA = ARANDU_DATA;
