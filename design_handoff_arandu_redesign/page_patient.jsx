// Página: Perfil do Paciente
const PatientPage = ({ data, onOpenSession, onBack }) => {
  const p = data.patientProfile;
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 24 }}>
      {/* Hero editorial */}
      <div style={{
        display: "grid", gridTemplateColumns: "auto 1fr auto",
        gap: 28, alignItems: "center",
        padding: "8px 0 24px",
        borderBottom: "1px solid var(--line)",
      }}>
        <Avatar initials={p.initials} size={72} />
        <div>
          <div style={{ fontSize: 11, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", marginBottom: 8, fontWeight: 500 }}>
            Paciente #{p.id.toUpperCase()}
          </div>
          <h1 className="serif" style={{ margin: 0, fontSize: 40, fontWeight: 400, letterSpacing: -0.8, lineHeight: 1 }}>
            {p.name}
          </h1>
          <div style={{ display: "flex", flexWrap: "wrap", gap: 14, marginTop: 10, fontSize: 13, color: "var(--ink-3)" }}>
            <span>{p.age} anos · {p.pronouns}</span>
            <span>·</span>
            <span>{p.marker}</span>
            <span>·</span>
            <span>{p.since}</span>
          </div>
        </div>
        <div style={{ display: "flex", gap: 24, alignItems: "center" }}>
          <StatBlock label="Sessões" value={p.sessions} />
          <div style={{ width: 1, height: 40, background: "var(--line)" }} />
          <StatBlock label="Em terapia" value={p.therapyDuration} />
          <div style={{ width: 1, height: 40, background: "var(--line)" }} />
          <StatBlock label="Frequência" value={p.frequency} />
        </div>
      </div>

      {/* Corpo em 2 colunas */}
      <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: 20 }}>
        {/* Coluna principal */}
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          {/* Triagem em destaque — tratada como citação editorial */}
          <Card eyebrow="Notas de triagem" title="Primeira escuta" subtitle="Registrada em 28/01/2026">
            <blockquote style={{
              margin: 0, padding: "4px 0 4px 20px",
              borderLeft: "2px solid var(--accent)",
              fontFamily: "var(--font-serif)",
              fontSize: 18, lineHeight: 1.55,
              color: "var(--ink-2)",
              fontWeight: 400, letterSpacing: -0.1,
            }}>{p.triage}</blockquote>
          </Card>

          {/* Timeline */}
          <Card eyebrow="Percurso" title="Linha do tempo clínica" action={
            <div style={{ display: "flex", gap: 6 }}>
              <button style={filterBtn(true)}>Tudo</button>
              <button style={filterBtn(false)}>Sessões</button>
              <button style={filterBtn(false)}>Notas</button>
            </div>
          }>
            <div style={{ display: "flex", flexDirection: "column" }}>
              {p.timeline.map((t, i) => (
                <TimelineItem key={i} item={t} last={i === p.timeline.length - 1} onClick={() => t.kind === "sessão" && onOpenSession()} />
              ))}
            </div>
          </Card>
        </div>

        {/* Coluna lateral */}
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          {/* Ações rápidas */}
          <Card eyebrow="Atalhos" title="Ações">
            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
              {p.quickActions.map(a => (
                <button key={a.id} style={{
                  display: "flex", alignItems: "center", gap: 10,
                  padding: "11px 14px", borderRadius: 10,
                  background: a.primary ? "var(--ink)" : "var(--paper)",
                  border: "1px solid " + (a.primary ? "var(--ink)" : "var(--line)"),
                  color: a.primary ? "var(--paper)" : "var(--ink-2)",
                  fontSize: 13, fontWeight: 500, textAlign: "left",
                }}>
                  <Icon name={a.primary ? "plus" : a.id === "anamnese" ? "notes" : a.id === "plano" ? "book" : "brain"} size={15} />
                  <span style={{ flex: 1 }}>{a.label}</span>
                  <Icon name="chevRight" size={12} />
                </button>
              ))}
            </div>
          </Card>

          {/* Observações recentes */}
          <Card eyebrow="Percepções" title="Observações recentes" subtitle="Anotações da clínica">
            <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
              {p.recentObservations.map((o, i) => (
                <div key={i} style={{
                  padding: "12px 14px",
                  background: "var(--paper)",
                  border: "1px solid var(--line)",
                  borderRadius: 10,
                }}>
                  <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 8 }}>
                    <Pill tone="accent" size="xs">#{o.tag}</Pill>
                    <span className="mono" style={{ fontSize: 10.5, color: "var(--ink-3)", marginLeft: "auto" }}>{o.date}</span>
                  </div>
                  <p className="serif" style={{
                    margin: 0, fontSize: 13.5, lineHeight: 1.55,
                    color: "var(--ink-2)",
                  }}>{o.text}</p>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

const filterBtn = (active) => ({
  padding: "5px 11px", borderRadius: 20,
  background: active ? "var(--ink)" : "transparent",
  border: "1px solid " + (active ? "var(--ink)" : "var(--line)"),
  color: active ? "var(--paper)" : "var(--ink-3)",
  fontSize: 11.5, fontWeight: 500,
});

const StatBlock = ({ label, value }) => (
  <div style={{ textAlign: "right" }}>
    <div className="serif" style={{ fontSize: 26, fontWeight: 500, letterSpacing: -0.5, lineHeight: 1 }}>{value}</div>
    <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", marginTop: 4, fontWeight: 500 }}>{label}</div>
  </div>
);

const TimelineItem = ({ item, last, onClick }) => {
  const isSession = item.kind === "sessão";
  return (
    <button onClick={onClick} style={{
      display: "grid", gridTemplateColumns: "96px auto 1fr auto",
      gap: 16, alignItems: "flex-start",
      padding: "18px 0",
      borderBottom: last ? "none" : "1px dashed var(--line)",
      background: "transparent", border: 0,
      borderBottomStyle: last ? "none" : "dashed",
      borderBottomWidth: last ? 0 : 1, borderBottomColor: "var(--line)",
      textAlign: "left", width: "100%", cursor: isSession ? "pointer" : "default",
      position: "relative",
    }}>
      <div className="mono" style={{ fontSize: 12, color: "var(--ink-3)", paddingTop: 2 }}>{item.date}</div>
      <div style={{
        width: 28, height: 28, borderRadius: "50%",
        background: isSession ? "var(--accent)" : "var(--paper-2)",
        border: "1px solid " + (isSession ? "var(--accent-deep)" : "var(--line)"),
        color: isSession ? "var(--paper)" : "var(--ink-3)",
        display: "flex", alignItems: "center", justifyContent: "center",
        flexShrink: 0,
      }}>
        <Icon name={isSession ? "notes" : "pin"} size={13} />
      </div>
      <div>
        <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 4 }}>
          <h4 className="serif" style={{ margin: 0, fontSize: 16, fontWeight: 500, color: "var(--ink)" }}>{item.title}</h4>
          <Pill tone={isSession ? "accent" : "neutral"} size="xs">{item.kind}</Pill>
        </div>
        <p style={{ margin: 0, fontSize: 13, color: "var(--ink-3)", lineHeight: 1.5 }}>{item.summary}</p>
      </div>
      {isSession && <Icon name="chevRight" size={14} style={{ color: "var(--ink-4)", marginTop: 4 }} />}
    </button>
  );
};

window.PatientPage = PatientPage;
