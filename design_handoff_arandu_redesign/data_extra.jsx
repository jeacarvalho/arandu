// Dados mockados para as 3 novas telas — carregado APÓS data.jsx
// Anexa ao ARANDU_DATA existente

(function() {
  const D = window.ARANDU_DATA;

  D.agenda = {
    weekLabel: "13 – 19 de abril · 2026",
    weekDays: [
      { dow: "seg", d: 13 },
      { dow: "ter", d: 14 },
      { dow: "qua", d: 15 },
      { dow: "qui", d: 16 },
      { dow: "sex", d: 17 },
      { dow: "sáb", d: 18 },
      { dow: "dom", d: 19, today: true },
    ],
    // formato: col 0..6, start em 30min steps desde 08:00
    events: [
      { day: 0, start: 2,  dur: 2, patient: "Amanda Rocha",    type: "Retorno",            tone: "accent", status: "done" },
      { day: 0, start: 7,  dur: 2, patient: "Felipe Lima",     type: "Acompanhamento",     tone: "neutral", status: "done" },
      { day: 0, start: 12, dur: 2, patient: "Carlos Santos",   type: "Retorno",            tone: "danger", status: "done" },

      { day: 1, start: 4,  dur: 2, patient: "Gustavo Pereira", type: "Acompanhamento",     tone: "neutral", status: "done" },
      { day: 1, start: 10, dur: 2, patient: "André Barbosa",   type: "Retorno",            tone: "warn", status: "done" },

      { day: 2, start: 1,  dur: 2, patient: "Amanda Rocha",    type: "Retorno",            tone: "accent", status: "done" },
      { day: 2, start: 6,  dur: 2, patient: "Felipe Lima",     type: "Acompanhamento",     tone: "neutral", status: "done" },
      { day: 2, start: 12, dur: 3, patient: "Supervisão",      type: "Grupo de estudos",   tone: "moss", status: "done" },

      { day: 3, start: 3,  dur: 2, patient: "Carolina Costa",  type: "Primeira consulta",  tone: "info", status: "done" },
      { day: 3, start: 9,  dur: 2, patient: "Gustavo Pereira", type: "Acompanhamento",     tone: "neutral", status: "done" },

      { day: 4, start: 2,  dur: 2, patient: "Carlos Santos",   type: "Retorno",            tone: "danger", status: "done" },
      { day: 4, start: 8,  dur: 2, patient: "André Barbosa",   type: "Retorno",            tone: "warn", status: "done" },

      { day: 5, start: 4,  dur: 3, patient: "Escrita clínica", type: "Bloco pessoal",      tone: "ghost", status: "blocked" },

      // Hoje (domingo, índice 6) — na prática no design mostramos o dia de hoje destacado
      { day: 6, start: 2,  dur: 2, patient: "Amanda Rocha",     type: "Retorno",            tone: "accent", status: "done" },
      { day: 6, start: 5,  dur: 2, patient: "André Barbosa",    type: "Acompanhamento",     tone: "warn", status: "done" },
      { day: 6, start: 12, dur: 2, patient: "Carolina Costa",   type: "Primeira consulta",  tone: "info", status: "next" },
      { day: 6, start: 15, dur: 2, patient: "Felipe Lima",      type: "Retorno",            tone: "neutral", status: "upcoming" },
      { day: 6, start: 18, dur: 2, patient: "Gustavo Pereira",  type: "Acompanhamento",     tone: "neutral", status: "upcoming" },
    ],
    requests: [
      { patient: "Amanda Rocha", when: "Qua 09:00 → Qui 11:00", reason: "Conflito de agenda" },
      { patient: "Lucas Meirelles", when: "Nova 1ª consulta", reason: "Indicação de colega" },
    ],
    metrics: [
      { label: "Ocupação da semana", value: "78%" },
      { label: "Horas clínicas",     value: "14h" },
      { label: "Cancelamentos",      value: "1" },
    ],
  };

  D.records = {
    filterTags: ["Ansiedade", "Burnout", "Luto", "TOC", "Depressão", "Triagem"],
    records: [
      { id: "PR-0446", patient: "André Barbosa",    initials: "AB", tags: ["Burnout"],        status: "Em acompanhamento", lastUpdate: "02/04/2026", sessions: 5,  pages: 12, risk: "Moderado",   since: "dez/2023", pinned: true },
      { id: "PR-0312", patient: "Amanda Rocha",     initials: "AR", tags: ["Ansiedade"],      status: "Em acompanhamento", lastUpdate: "02/04/2026", sessions: 14, pages: 28, risk: "Baixo",      since: "mar/2024" },
      { id: "PR-0188", patient: "Carolina Costa",   initials: "CC", tags: ["Triagem"],        status: "Em avaliação",      lastUpdate: "hoje",       sessions: 0,  pages: 2,  risk: "A avaliar",  since: "abr/2026" },
      { id: "PR-0077", patient: "Felipe Lima",      initials: "FL", tags: ["Luto"],           status: "Em acompanhamento", lastUpdate: "28/03/2026", sessions: 11, pages: 22, risk: "Moderado",   since: "ago/2025", pinned: true },
      { id: "PR-0054", patient: "Gustavo Pereira",  initials: "GP", tags: ["TOC"],            status: "Em acompanhamento", lastUpdate: "31/03/2026", sessions: 22, pages: 44, risk: "Baixo",      since: "jan/2025" },
      { id: "PR-0029", patient: "Carlos Santos",    initials: "CS", tags: ["Depressão"],      status: "Atenção",           lastUpdate: "01/04/2026", sessions: 31, pages: 61, risk: "Atenção",    since: "nov/2024" },
      { id: "PR-0411", patient: "Beatriz Antunes",  initials: "BA", tags: ["Ansiedade"],      status: "Alta",              lastUpdate: "12/01/2026", sessions: 18, pages: 36, risk: "Alta",       since: "mai/2024" },
      { id: "PR-0398", patient: "Henrique Mourão",  initials: "HM", tags: ["Burnout","Luto"], status: "Em acompanhamento", lastUpdate: "30/03/2026", sessions: 7,  pages: 14, risk: "Moderado",   since: "out/2025" },
    ],
    focused: {
      id: "PR-0446",
      patient: "André Barbosa",
      summary: "Engenheiro, 57 anos. Burnout com sintomas somáticos. Em acompanhamento quinzenal desde jan/2026.",
      sections: [
        { key: "anamnese", label: "Anamnese", updated: "28/01/2026", pages: 3 },
        { key: "hipoteses", label: "Hipóteses clínicas", updated: "12/02/2026", pages: 2 },
        { key: "plano", label: "Plano terapêutico", updated: "26/02/2026", pages: 2 },
        { key: "evolucao", label: "Evolução", updated: "02/04/2026", pages: 5 },
      ],
      lastEntry: "Sessão marcada por abertura reflexiva sobre sentido profissional e emergência de material hipomaníaco. Recomenda-se monitorar oscilações de humor entre sessões.",
    },
  };

  D.insights = {
    generatedAt: "hoje, 07:12",
    coverage: "Últimos 90 dias · 42 pacientes · 186 sessões",
    headline: {
      title: "Três padrões emergentes na sua prática clínica",
      body: "Nas últimas semanas, três temas cruzam múltiplos pacientes. Um deles — oscilação de humor em perfis de burnout — pode merecer um olhar especial.",
    },
    themes: [
      { label: "Sentido e propósito",      weight: 0.82, patients: 9,  trend: "+14%", tone: "accent" },
      { label: "Oscilação de humor",       weight: 0.68, patients: 6,  trend: "+22%", tone: "warn" },
      { label: "Luto persistente",         weight: 0.54, patients: 4,  trend: "+4%",  tone: "info" },
      { label: "Sintomas somáticos",       weight: 0.41, patients: 7,  trend: "−8%",  tone: "ok" },
      { label: "Isolamento social",        weight: 0.38, patients: 5,  trend: "estável", tone: "neutral" },
      { label: "Perfeccionismo",           weight: 0.29, patients: 3,  trend: "+2%",  tone: "neutral" },
    ],
    cohorts: [
      { name: "Burnout (n=8)",    trend: [3,4,4,5,7,8,6,7,9,8,9,10], note: "Picos no retorno das férias." },
      { name: "Ansiedade (n=12)", trend: [6,6,7,5,6,7,7,6,5,6,6,5],  note: "Redução gradual." },
      { name: "Luto (n=4)",       trend: [2,2,3,3,3,4,4,4,4,4,4,4],  note: "Estabilidade no vínculo." },
    ],
    alerts: [
      { level: "attention", patient: "Carlos Santos",   text: "3 sessões consecutivas com afeto disfórico e menções a desesperança. Considerar revisão de plano.", when: "há 2 dias" },
      { level: "watch",     patient: "André Barbosa",   text: "Material hipomaníaco emergente em 2 sessões. Sugerido monitorar ritmo.", when: "há 3 dias" },
      { level: "watch",     patient: "Felipe Lima",     text: "Aniversário de perda próximo — risco de intensificação do luto.", when: "há 1 semana" },
    ],
    questions: [
      "Quais pacientes mencionaram 'sentido' em múltiplas sessões?",
      "Comparar evolução de ansiedade entre cohort feminino e masculino",
      "Sugerir leituras recentes para Burnout com sintomas somáticos",
      "Que intervenções aparecem antes de melhoras em TOC?",
    ],
    threads: [
      { id: "t1", title: "Padrões cruzados entre sessões de André Barbosa", updated: "há 3h", msgs: 8 },
      { id: "t2", title: "Comparação de evolução: Amanda × Beatriz", updated: "ontem", msgs: 14 },
      { id: "t3", title: "Hipóteses diferenciais para Henrique Mourão", updated: "há 2 dias", msgs: 6 },
      { id: "t4", title: "Revisão de literatura: luto persistente", updated: "há 5 dias", msgs: 11 },
    ],
  };
})();
