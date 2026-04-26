// Página: Inteligência Clínica — painel de insights agregados
const InsightsPage = ({ data }) => {
  const ins = data.insights;
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 22 }}>
      {/* Hero editorial */}
      <div style={{
        display: "grid", gridTemplateColumns: "1fr auto", gap: 32, alignItems: "end",
        paddingBottom: 22, borderBottom: "1px solid var(--line)",
      }}>
        <div>
          <div style={{ fontSize: 11, letterSpacing: 1.6, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 10, display: "flex", gap: 10, alignItems: "center" }}>
            <Icon name="sparkles" size={13} style={{ color: "var(--accent-deep)" }} />
            Arandu · inteligência clínica
          </div>
          <h1 className="serif" style={{ margin: 0, fontSize: 40, fontWeight: 400, letterSpacing: -0.8, lineHeight: 1.05, maxWidth: 780 }}>
            {ins.headline.title.split("emergentes")[0]}<em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>emergentes</em> na sua prática
          </h1>
          <p className="serif" style={{ margin: "12px 0 0", fontSize: 17, color: "var(--ink-2)", lineHeight: 1.55, maxWidth: 720 }}>
            {ins.headline.body}
          </p>
        </div>
        <div style={{ textAlign: "right", fontSize: 11, color: "var(--ink-3)", letterSpacing: .5 }}>
          <div>gerado {ins.generatedAt}</div>
          <div style={{ marginTop: 4, color: "var(--ink-4)" }}>{ins.coverage}</div>
        </div>
      </div>

      {/* Grid principal */}
      <div style={{ display: "grid", gridTemplateColumns: "1.5fr 1fr", gap: 20 }}>
        {/* Coluna A */}
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          {/* Temas dominantes */}
          <Card eyebrow="Temas" title="Padrões dominantes na prática" subtitle="Proporção de pacientes impactados e tendência nas últimas 4 semanas">
            <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
              {ins.themes.map((t, i) => <ThemeBar key={i} theme={t} />)}
            </div>
          </Card>

          {/* Cohorts */}
          <Card eyebrow="Cohorts" title="Evolução por agrupamento clínico" subtitle="Tendência de volume de menções nas últimas 12 semanas">
            <div style={{ display: "flex", flexDirection: "column", gap: 18 }}>
              {ins.cohorts.map((c, i) => <CohortRow key={i} cohort={c} />)}
            </div>
          </Card>
        </div>

        {/* Coluna B */}
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          {/* Alertas */}
          <Card eyebrow="Atenção" title="Alertas clínicos" subtitle={`${ins.alerts.length} pontos merecem um olhar`}>
            <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
              {ins.alerts.map((a, i) => <AlertItem key={i} alert={a} />)}
            </div>
          </Card>

          {/* Conversas */}
          <Card eyebrow="Conversas" title="Histórico com Arandu">
            <div style={{ display: "flex", flexDirection: "column", gap: 2 }}>
              {ins.threads.map((t, i) => (
                <button key={t.id} style={{
                  textAlign: "left", background: "transparent", border: 0,
                  padding: "12px 0", borderBottom: i === ins.threads.length - 1 ? "none" : "1px dashed var(--line)",
                  display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 10,
                  alignItems: "center", cursor: "pointer",
                }}>
                  <div style={{
                    width: 26, height: 26, borderRadius: 7,
                    background: "linear-gradient(135deg, var(--accent-deep), var(--accent))",
                    color: "var(--paper)",
                    display: "flex", alignItems: "center", justifyContent: "center",
                  }}>
                    <Icon name="sparkles" size={12} />
                  </div>
                  <div style={{ minWidth: 0 }}>
                    <div className="serif" style={{ fontSize: 14, color: "var(--ink)", letterSpacing: -0.1, fontStyle: "italic", fontWeight: 500, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                      {t.title}
                    </div>
                    <div style={{ fontSize: 11, color: "var(--ink-3)", marginTop: 2 }}>
                      {t.msgs} msgs · {t.updated}
                    </div>
                  </div>
                  <Icon name="chevRight" size={12} style={{ color: "var(--ink-4)" }} />
                </button>
              ))}
            </div>
          </Card>

          {/* Perguntas sugeridas */}
          <Card eyebrow="Perguntar a Arandu" title="Sugestões de investigação">
            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
              {ins.questions.map((q, i) => (
                <button key={i} style={{
                  textAlign: "left", padding: "10px 14px",
                  background: "var(--paper)", border: "1px dashed var(--line-2)",
                  borderRadius: 10, color: "var(--ink-2)",
                  fontFamily: "var(--font-serif)", fontStyle: "italic", fontSize: 13.5,
                  display: "flex", alignItems: "center", gap: 10, cursor: "pointer",
                }}>
                  <Icon name="sparkles" size={12} style={{ color: "var(--accent)", flexShrink: 0 }} />
                  <span style={{ flex: 1 }}>{q}</span>
                  <Icon name="arrowRight" size={12} style={{ color: "var(--ink-4)" }} />
                </button>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

const ThemeBar = ({ theme }) => {
  const colorByTone = {
    accent: "var(--accent)", warn: "var(--clay)", info: "var(--accent-soft)",
    ok: "var(--moss-2)", neutral: "var(--ink-4)",
  };
  const trendUp = theme.trend.startsWith("+");
  const trendDown = theme.trend.startsWith("−");
  return (
    <div>
      <div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 6 }}>
        <span className="serif" style={{ fontSize: 16, fontWeight: 500, color: "var(--ink)", flex: 1, letterSpacing: -0.1 }}>
          {theme.label}
        </span>
        <span style={{ fontSize: 11, color: "var(--ink-3)" }}>
          {theme.patients} {theme.patients === 1 ? "paciente" : "pacientes"}
        </span>
        <span style={{
          fontSize: 11, fontWeight: 500, padding: "2px 8px", borderRadius: 999,
          background: trendUp ? "color-mix(in oklab, var(--clay) 14%, transparent)" : trendDown ? "color-mix(in oklab, var(--moss-2) 14%, transparent)" : "color-mix(in oklab, var(--ink) 6%, transparent)",
          color: trendUp ? "#8B3A24" : trendDown ? "var(--moss)" : "var(--ink-3)",
          minWidth: 52, textAlign: "center",
        }}>{theme.trend}</span>
      </div>
      <div style={{ position: "relative", height: 6, borderRadius: 3, background: "var(--paper-3)", overflow: "hidden" }}>
        <div style={{
          position: "absolute", left: 0, top: 0, bottom: 0,
          width: `${theme.weight * 100}%`,
          background: colorByTone[theme.tone] || "var(--accent)",
          borderRadius: 3,
        }} />
      </div>
    </div>
  );
};

const CohortRow = ({ cohort }) => {
  const max = Math.max(...cohort.trend);
  const points = cohort.trend.map((v, i) => {
    const x = (i / (cohort.trend.length - 1)) * 100;
    const y = 100 - (v / max) * 100;
    return `${x},${y}`;
  }).join(" ");
  const last = cohort.trend[cohort.trend.length - 1];
  const first = cohort.trend[0];
  const delta = last - first;
  return (
    <div style={{ display: "grid", gridTemplateColumns: "160px 1fr auto", gap: 14, alignItems: "center" }}>
      <div>
        <div className="serif" style={{ fontSize: 15, fontWeight: 500, letterSpacing: -0.1 }}>{cohort.name.split(" (")[0]}</div>
        <div style={{ fontSize: 11, color: "var(--ink-3)", marginTop: 2 }}>{cohort.note}</div>
      </div>
      <div style={{ height: 44, width: "100%" }}>
        <svg viewBox="0 0 100 100" preserveAspectRatio="none" style={{ width: "100%", height: "100%" }}>
          <polyline points={points} fill="none" stroke="var(--accent)" strokeWidth="1.5" vectorEffect="non-scaling-stroke" />
          <polyline points={`0,100 ${points} 100,100`} fill="color-mix(in oklab, var(--accent) 10%, transparent)" stroke="none" />
        </svg>
      </div>
      <div style={{ textAlign: "right" }}>
        <div className="serif" style={{ fontSize: 20, fontWeight: 500, lineHeight: 1 }}>{last}</div>
        <div style={{ fontSize: 10.5, color: delta > 0 ? "var(--clay)" : delta < 0 ? "var(--moss)" : "var(--ink-3)", marginTop: 3, fontWeight: 500 }}>
          {delta > 0 ? "+" : ""}{delta} vs início
        </div>
      </div>
    </div>
  );
};

const AlertItem = ({ alert }) => {
  const isAttention = alert.level === "attention";
  return (
    <div style={{
      padding: "12px 14px",
      background: isAttention ? "color-mix(in oklab, var(--clay) 8%, var(--paper))" : "var(--paper)",
      border: "1px solid " + (isAttention ? "color-mix(in oklab, var(--clay) 28%, transparent)" : "var(--line)"),
      borderRadius: 10,
      display: "flex", gap: 12, alignItems: "flex-start",
    }}>
      <div style={{
        width: 26, height: 26, borderRadius: "50%", flexShrink: 0,
        background: isAttention ? "var(--clay)" : "var(--paper-3)",
        color: isAttention ? "var(--paper)" : "var(--ink-3)",
        display: "flex", alignItems: "center", justifyContent: "center",
        marginTop: 2,
      }}>
        <Icon name="alert" size={13} />
      </div>
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 4 }}>
          <span style={{ fontSize: 13, fontWeight: 500 }}>{alert.patient}</span>
          <span style={{ fontSize: 10.5, color: "var(--ink-4)", marginLeft: "auto" }}>{alert.when}</span>
        </div>
        <p className="serif" style={{ margin: 0, fontSize: 13.5, lineHeight: 1.5, color: "var(--ink-2)" }}>{alert.text}</p>
      </div>
    </div>
  );
};

window.InsightsPage = InsightsPage;
