// Página: Dashboard
const DashboardPage = ({ data, onOpenPatient }) => {
  const today = new Date().toLocaleDateString("pt-BR", { weekday: "long", day: "numeric", month: "long" });
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 28 }}>
      {/* Cabeçalho editorial */}
      <div style={{
        display: "grid", gridTemplateColumns: "1fr auto", gap: 40, alignItems: "end",
        paddingBottom: 24, borderBottom: "1px solid var(--line)",
      }}>
        <div>
          <div style={{
            fontSize: 11, letterSpacing: 1.6, textTransform: "uppercase",
            color: "var(--ink-3)", fontWeight: 500, marginBottom: 10,
          }}>{today}</div>
          <h1 className="serif" style={{
            margin: 0, fontSize: 44, fontWeight: 400, letterSpacing: -0.8,
            lineHeight: 1.05, color: "var(--ink)",
          }}>
            Bom dia, <em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>Helena</em>.
          </h1>
          <p style={{ margin: "10px 0 0", color: "var(--ink-3)", fontSize: 15, maxWidth: 560 }}>
            Você tem <strong style={{ color: "var(--ink)", fontWeight: 500 }}>5 sessões</strong> hoje e <strong style={{ color: "var(--ink)", fontWeight: 500 }}>2 anotações</strong> aguardando revisão desde ontem.
          </p>
        </div>
        <div style={{ display: "flex", gap: 8 }}>
          <button style={{
            padding: "10px 16px", borderRadius: 10,
            background: "var(--paper-2)", border: "1px solid var(--line)",
            color: "var(--ink-2)", fontSize: 13, fontWeight: 500,
            display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="calendar" size={15} /> Abrir agenda
          </button>
          <button style={{
            padding: "10px 16px", borderRadius: 10,
            background: "var(--ink)", border: "1px solid var(--ink)",
            color: "var(--paper)", fontSize: 13, fontWeight: 500,
            display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="plus" size={15} /> Nova sessão
          </button>
        </div>
      </div>

      {/* KPIs */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 16 }}>
        {data.kpis.map((k, i) => <KpiCard key={k.id} kpi={k} index={i} />)}
      </div>

      {/* Grid principal: agenda + pacientes + sessões */}
      <div style={{ display: "grid", gridTemplateColumns: "1.1fr 1fr", gap: 20 }}>
        <TodayColumn data={data} />
        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          <PatientsPreview data={data} onOpenPatient={onOpenPatient} />
          <RecentSessionsCard data={data} />
        </div>
      </div>
    </div>
  );
};

const KpiCard = ({ kpi, index }) => {
  const tones = {
    neutral: "var(--ink-3)",
    up: "var(--moss)",
    warn: "var(--clay)",
  };
  return (
    <div style={{
      background: index === 0 ? "var(--ink)" : "var(--paper-2)",
      color: index === 0 ? "var(--paper)" : "var(--ink)",
      border: "1px solid " + (index === 0 ? "var(--ink)" : "var(--line)"),
      borderRadius: "var(--radius)",
      padding: "20px 22px",
      display: "flex", flexDirection: "column", gap: 12,
      position: "relative", overflow: "hidden",
      boxShadow: "var(--shadow-sm)",
    }}>
      <div style={{
        fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase",
        color: index === 0 ? "color-mix(in oklab, var(--paper) 65%, transparent)" : "var(--ink-3)",
        fontWeight: 500,
      }}>{kpi.label}</div>
      <div className="serif" style={{
        fontSize: 42, fontWeight: 400, letterSpacing: -1,
        lineHeight: 1, fontVariantNumeric: "tabular-nums",
      }}>{kpi.value.toLocaleString("pt-BR")}</div>
      <div style={{
        fontSize: 11.5,
        color: index === 0 ? "color-mix(in oklab, var(--paper) 65%, transparent)" : tones[kpi.tone],
        display: "flex", alignItems: "center", gap: 6,
      }}>
        {kpi.tone === "up" && <Icon name="trend" size={12} />}
        {kpi.tone === "warn" && <Icon name="alert" size={12} />}
        {kpi.delta}
      </div>
      {/* Decoração sutil */}
      {index === 0 && <div style={{
        position: "absolute", right: -30, top: -30, width: 120, height: 120,
        borderRadius: "50%", background: "color-mix(in oklab, var(--accent) 40%, transparent)",
        filter: "blur(30px)", pointerEvents: "none",
      }} />}
    </div>
  );
};

const TodayColumn = ({ data }) => (
  <Card eyebrow="Hoje" title="Agenda do dia" subtitle="Sessões agendadas para terça-feira" action={
    <button style={{
      background: "transparent", border: 0, color: "var(--ink-3)",
      fontSize: 12, display: "flex", alignItems: "center", gap: 4,
    }}>Semana inteira <Icon name="chevRight" size={12} /></button>
  }>
    <div style={{ display: "flex", flexDirection: "column" }}>
      {data.todaySchedule.map((s, i) => <ScheduleItem key={i} item={s} last={i === data.todaySchedule.length - 1} />)}
    </div>
  </Card>
);

const ScheduleItem = ({ item, last }) => {
  const done = item.status === "done";
  const next = item.status === "next";
  return (
    <div style={{
      display: "grid", gridTemplateColumns: "64px 1fr auto", gap: 16,
      alignItems: "center",
      padding: "14px 0",
      borderBottom: last ? "none" : "1px dashed var(--line)",
      opacity: done ? 0.55 : 1,
    }}>
      <div className="mono" style={{
        fontSize: 13, color: next ? "var(--accent-deep)" : "var(--ink-3)",
        fontWeight: next ? 600 : 400,
        textDecoration: done ? "line-through" : "none",
      }}>{item.time}</div>
      <div>
        <div style={{ fontSize: 14, fontWeight: 500, color: "var(--ink)" }}>{item.patient}</div>
        <div style={{ fontSize: 12, color: "var(--ink-3)", marginTop: 2 }}>{item.type}</div>
      </div>
      <div>
        {done && <Pill tone="ok" size="xs">Concluída</Pill>}
        {next && <Pill tone="accent" size="xs">Próxima</Pill>}
        {!done && !next && <Pill tone="neutral" size="xs">Agendada</Pill>}
      </div>
    </div>
  );
};

const PatientsPreview = ({ data, onOpenPatient }) => (
  <Card eyebrow="Pacientes" title="Em acompanhamento" action={
    <button style={{
      background: "transparent", border: 0, color: "var(--ink-3)",
      fontSize: 12, display: "flex", alignItems: "center", gap: 4,
    }}>Ver todos (42) <Icon name="chevRight" size={12} /></button>
  }>
    <div style={{ display: "flex", flexDirection: "column" }}>
      {data.patients.slice(0, 5).map((p, i) => (
        <button key={p.id} onClick={() => onOpenPatient(p.id)} style={{
          display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 14,
          alignItems: "center", padding: "12px 0",
          borderBottom: i === 4 ? "none" : "1px dashed var(--line)",
          background: "transparent", border: 0, borderBottomStyle: i === 4 ? "none" : "dashed",
          borderBottomWidth: i === 4 ? 0 : 1, borderBottomColor: "var(--line)",
          textAlign: "left", width: "100%", cursor: "pointer",
        }}>
          <Avatar initials={p.name.split(" ").map(w => w[0]).slice(0,2).join("")} size={34} />
          <div style={{ minWidth: 0 }}>
            <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
              <span style={{ fontSize: 14, fontWeight: 500, color: "var(--ink)" }}>{p.name}</span>
              <Pill tone={p.tagTone} size="xs">{p.tag}</Pill>
            </div>
            <div style={{ fontSize: 11.5, color: "var(--ink-3)", marginTop: 3 }}>
              {p.sessions} sessões · última {p.last} · próxima {p.next}
            </div>
          </div>
          <Icon name="chevRight" size={14} className="" style={{ color: "var(--ink-4)" }} />
        </button>
      ))}
    </div>
  </Card>
);

const RecentSessionsCard = ({ data }) => (
  <Card eyebrow="Reflexão" title="Últimas sessões registradas">
    <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
      {data.recentSessions.map((s, i) => (
        <div key={i} style={{ display: "flex", gap: 14 }}>
          <div style={{
            width: 4, borderRadius: 2,
            background: "linear-gradient(to bottom, var(--accent), color-mix(in oklab, var(--accent) 30%, transparent))",
            flexShrink: 0,
          }} />
          <div style={{ flex: 1 }}>
            <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 4 }}>
              <span style={{ fontSize: 12.5, fontWeight: 500 }}>{s.patient}</span>
              <span style={{ fontSize: 11, color: "var(--ink-4)" }}>·</span>
              <span className="mono" style={{ fontSize: 11, color: "var(--ink-3)" }}>{s.date}</span>
              <div style={{ flex: 1 }} />
              <Pill tone="neutral" size="xs">{s.theme}</Pill>
            </div>
            <p className="serif" style={{
              margin: 0, fontSize: 14, color: "var(--ink-2)",
              fontStyle: "italic", lineHeight: 1.5,
            }}>"{s.preview}"</p>
          </div>
        </div>
      ))}
    </div>
  </Card>
);

window.DashboardPage = DashboardPage;
