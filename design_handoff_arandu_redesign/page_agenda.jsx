// Página: Agenda — calendário semanal editorial
const AgendaPage = ({ data, onOpenPatient }) => {
  const a = data.agenda;
  const [view, setView] = React.useState("week"); // "day" | "week" | "month"
  // 22 slots de 30 min: 08:00 → 19:00
  const slots = Array.from({ length: 22 }, (_, i) => {
    const h = 8 + Math.floor(i / 2);
    const m = i % 2 === 0 ? "00" : "30";
    return `${String(h).padStart(2, "0")}:${m}`;
  });

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 22 }}>
      {/* Hero */}
      <div style={{
        display: "grid", gridTemplateColumns: "1fr auto", gap: 32, alignItems: "end",
        paddingBottom: 22, borderBottom: "1px solid var(--line)",
      }}>
        <div>
          <div style={{ fontSize: 11, letterSpacing: 1.6, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 10 }}>
            Agenda · {view === "day" ? "dia" : view === "month" ? "mês" : "semana"}
          </div>
          <h1 className="serif" style={{ margin: 0, fontSize: 40, fontWeight: 400, letterSpacing: -0.8, lineHeight: 1 }}>
            {view === "day" && <>Domingo, <em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>19 abril</em></>}
            {view === "week" && <>{a.weekLabel.split("·")[0]}<em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}> · {a.weekLabel.split("·")[1].trim()}</em></>}
            {view === "month" && <>Abril<em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}> · 2026</em></>}
          </h1>
          <div style={{ marginTop: 10, display: "flex", gap: 14, alignItems: "center", fontSize: 13, color: "var(--ink-3)" }}>
            {a.metrics.map((m, i) => (
              <React.Fragment key={i}>
                {i > 0 && <span>·</span>}
                <span><strong style={{ color: "var(--ink-2)", fontWeight: 500 }}>{m.value}</strong> {m.label.toLowerCase()}</span>
              </React.Fragment>
            ))}
          </div>
        </div>
        <div style={{ display: "flex", gap: 8, alignItems: "center" }}>
          <div style={{ display: "flex", padding: 3, background: "var(--paper-2)", border: "1px solid var(--line)", borderRadius: 10 }}>
            {[{v:"day",l:"Dia"},{v:"week",l:"Semana"},{v:"month",l:"Mês"}].map(o => (
              <button key={o.v} onClick={() => setView(o.v)} style={{
                padding: "6px 12px", border: 0, borderRadius: 7,
                background: view === o.v ? "var(--ink)" : "transparent",
                color: view === o.v ? "var(--paper)" : "var(--ink-3)",
                fontSize: 12, fontWeight: 500, cursor: "pointer",
              }}>{o.l}</button>
            ))}
          </div>
          <button style={navBtn()}><Icon name="arrowLeft" size={14} /></button>
          <button style={{ ...navBtn(), padding: "8px 14px", fontSize: 12.5, color: "var(--ink-2)" }}>Hoje</button>
          <button style={navBtn()}><Icon name="arrowRight" size={14} /></button>
          <button style={{
            marginLeft: 8, padding: "8px 14px", borderRadius: 10,
            background: "var(--ink)", color: "var(--paper)", border: "1px solid var(--ink)",
            fontSize: 13, fontWeight: 500, display: "flex", alignItems: "center", gap: 8,
          }}><Icon name="plus" size={14} /> Novo agendamento</button>
        </div>
      </div>

      {/* Grid principal: calendário + lateral */}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 320px", gap: 20 }}>
        {view === "week" && <WeekView a={a} slots={slots} />}
        {view === "day" && <DayView a={a} slots={slots} />}
        {view === "month" && <MonthView a={a} />}

        {/* Lateral */}
        <div style={{ display: "flex", flexDirection: "column", gap: 18 }}>
          <Card eyebrow="Solicitações" title="Em aberto" subtitle={`${a.requests.length} aguardando resposta`}>
            <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
              {a.requests.map((r, i) => (
                <div key={i} style={{
                  padding: 12, background: "var(--paper)", border: "1px solid var(--line)",
                  borderRadius: 10,
                }}>
                  <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 6 }}>
                    <Avatar initials={r.patient.split(" ").map(w => w[0]).slice(0,2).join("")} size={26} />
                    <span style={{ fontSize: 13, fontWeight: 500 }}>{r.patient}</span>
                  </div>
                  <div style={{ fontSize: 12, color: "var(--ink-3)", marginBottom: 3 }}>{r.when}</div>
                  <div className="serif" style={{ fontSize: 13, fontStyle: "italic", color: "var(--ink-2)", marginBottom: 10 }}>
                    "{r.reason}"
                  </div>
                  <div style={{ display: "flex", gap: 6 }}>
                    <button style={miniBtn(true)}>Aceitar</button>
                    <button style={miniBtn(false)}>Propor outro horário</button>
                  </div>
                </div>
              ))}
            </div>
          </Card>

          <Card eyebrow="Disponibilidade" title="Horários livres esta semana">
            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
              {[
                { day: "Ter", time: "11:00", avail: 1 },
                { day: "Qua", time: "14:00", avail: 1 },
                { day: "Qui", time: "16:30", avail: 1 },
                { day: "Sex", time: "10:00", avail: 2 },
              ].map((s, i) => (
                <div key={i} style={{
                  display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 10,
                  alignItems: "center", padding: "9px 12px",
                  background: "var(--paper)", border: "1px dashed var(--line-2)",
                  borderRadius: 10,
                }}>
                  <span className="mono" style={{ fontSize: 12, color: "var(--ink-3)" }}>{s.day}</span>
                  <span className="serif" style={{ fontSize: 15, fontWeight: 500 }}>{s.time}</span>
                  <span style={{ fontSize: 10.5, color: "var(--ink-4)" }}>livre</span>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

// === VIEW: Semana (original) ===
const WeekView = ({ a, slots }) => (
  <section style={{
    background: "var(--paper-2)", border: "1px solid var(--line)",
    borderRadius: "var(--radius)", overflow: "hidden", boxShadow: "var(--shadow-sm)",
  }}>
    <div style={{
      display: "grid", gridTemplateColumns: "64px repeat(7, 1fr)",
      borderBottom: "1px solid var(--line)", background: "var(--paper)",
    }}>
      <div />
      {a.weekDays.map((d, i) => (
        <div key={i} style={{
          padding: "14px 10px", textAlign: "center",
          borderLeft: "1px solid var(--line)",
          background: d.today ? "color-mix(in oklab, var(--accent) 10%, transparent)" : "transparent",
        }}>
          <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500 }}>{d.dow}</div>
          <div className="serif" style={{
            fontSize: 22, fontWeight: 500, marginTop: 3, letterSpacing: -0.3,
            color: d.today ? "var(--accent-deep)" : "var(--ink)",
          }}>{d.d}</div>
        </div>
      ))}
    </div>

    <div style={{ position: "relative", display: "grid", gridTemplateColumns: "64px repeat(7, 1fr)" }}>
      <div>
        {slots.map((s, i) => (
          <div key={i} style={{
            height: 32, borderBottom: i % 2 === 1 ? "1px solid var(--line)" : "none",
            display: "flex", alignItems: "flex-start", justifyContent: "flex-end",
            padding: "2px 10px 0 0",
          }}>
            {i % 2 === 0 && <span className="mono" style={{ fontSize: 10.5, color: "var(--ink-4)" }}>{s}</span>}
          </div>
        ))}
      </div>
      {a.weekDays.map((d, di) => (
        <div key={di} style={{
          position: "relative", borderLeft: "1px solid var(--line)",
          background: d.today ? "color-mix(in oklab, var(--accent) 4%, transparent)" : "transparent",
        }}>
          {slots.map((_, i) => (
            <div key={i} style={{ height: 32, borderBottom: i % 2 === 1 ? "1px solid var(--line)" : "none" }} />
          ))}
          {a.events.filter(e => e.day === di).map((e, ei) => (
            <AgendaEvent key={ei} event={e} />
          ))}
        </div>
      ))}
    </div>

    <div style={{ padding: "12px 16px", borderTop: "1px solid var(--line)", display: "flex", gap: 16, flexWrap: "wrap", alignItems: "center", fontSize: 11, color: "var(--ink-3)" }}>
      <LegendDot tone="accent" label="Retorno" />
      <LegendDot tone="info" label="1ª consulta" />
      <LegendDot tone="warn" label="Atenção" />
      <LegendDot tone="danger" label="Risco" />
      <LegendDot tone="moss" label="Supervisão" />
      <LegendDot tone="ghost" label="Bloco pessoal" />
    </div>
  </section>
);

// === VIEW: Dia — foco editorial em 1 dia ===
const DayView = ({ a, slots }) => {
  const todayIdx = a.weekDays.findIndex(d => d.today);
  const dayEvents = a.events.filter(e => e.day === todayIdx).sort((x, y) => x.start - y.start);
  const toneColors = {
    accent: "#6B4E3D", info: "#3C5C7A", warn: "#B8842A",
    danger: "#A0463A", moss: "#4A5D4F", ghost: "var(--ink-4)", neutral: "var(--ink-4)",
  };
  const fmtTime = (s) => {
    const h = 8 + Math.floor(s / 2);
    const m = s % 2 === 0 ? "00" : "30";
    return `${String(h).padStart(2, "0")}:${m}`;
  };
  return (
    <section style={{
      background: "var(--paper-2)", border: "1px solid var(--line)",
      borderRadius: "var(--radius)", overflow: "hidden", boxShadow: "var(--shadow-sm)",
    }}>
      {/* Duas colunas: timeline slim + stack de cards */}
      <div style={{ display: "grid", gridTemplateColumns: "96px 1fr", minHeight: 640 }}>
        {/* Régua de horários */}
        <div style={{ borderRight: "1px solid var(--line)", background: "var(--paper)", padding: "20px 0" }}>
          {slots.filter((_, i) => i % 2 === 0).map((s, i) => (
            <div key={i} style={{
              height: 64, padding: "0 14px",
              display: "flex", alignItems: "flex-start", justifyContent: "flex-end",
              borderBottom: "1px dashed var(--line)",
            }}>
              <span className="mono" style={{ fontSize: 11, color: "var(--ink-4)", letterSpacing: .3 }}>{s}</span>
            </div>
          ))}
        </div>
        {/* Cards editoriais */}
        <div style={{ padding: "20px 24px", display: "flex", flexDirection: "column", gap: 10 }}>
          {dayEvents.map((e, i) => {
            const c = toneColors[e.tone] || "var(--ink-4)";
            const done = e.status === "done";
            const next = e.status === "next";
            return (
              <div key={i} style={{
                display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 16,
                alignItems: "center",
                padding: "14px 18px",
                background: "var(--paper)",
                border: `1px solid ${next ? c : "var(--line)"}`,
                borderLeft: `4px solid ${c}`,
                borderRadius: 10,
                opacity: done ? 0.55 : 1,
                boxShadow: next ? `0 0 0 3px color-mix(in oklab, ${c} 20%, transparent)` : "none",
              }}>
                <div style={{ textAlign: "right" }}>
                  <div className="mono" style={{ fontSize: 13, fontWeight: 500, color: c, textDecoration: done ? "line-through" : "none" }}>
                    {fmtTime(e.start)}
                  </div>
                  <div style={{ fontSize: 10.5, color: "var(--ink-4)", marginTop: 2 }}>{e.dur * 30}min</div>
                </div>
                <div style={{ minWidth: 0 }}>
                  <div className="serif" style={{ fontSize: 17, fontWeight: 500, color: "var(--ink)", letterSpacing: -0.2 }}>{e.patient}</div>
                  <div style={{ fontSize: 12, color: "var(--ink-3)", marginTop: 3 }}>{e.type}</div>
                </div>
                <div>
                  {done && <Pill tone="ok" size="xs">Concluída</Pill>}
                  {next && <Pill tone="accent" size="xs">Próxima</Pill>}
                  {!done && !next && <Pill tone="neutral" size="xs">Agendada</Pill>}
                </div>
              </div>
            );
          })}
          {dayEvents.length === 0 && (
            <div style={{ padding: 40, textAlign: "center", color: "var(--ink-3)", fontFamily: "var(--font-serif)", fontStyle: "italic", fontSize: 15 }}>
              Nenhum compromisso para este dia.
            </div>
          )}
        </div>
      </div>
    </section>
  );
};

// === VIEW: Mês — grade compacta com indicadores de volume ===
const MonthView = ({ a }) => {
  // abril 2026 começa numa quarta (4). Cria 6 semanas × 7 dias.
  const daysInMonth = 30;
  const firstWeekday = 3; // 0=dom
  const cells = [];
  for (let i = 0; i < 42; i++) {
    const dayNum = i - firstWeekday + 1;
    const inMonth = dayNum >= 1 && dayNum <= daysInMonth;
    cells.push({ dayNum, inMonth });
  }
  // distribui eventos: gera hash simples a partir dos eventos semanais para simular volume por dia
  const eventCount = (dayNum) => {
    if (!dayNum || dayNum < 1) return 0;
    const seed = (dayNum * 7 + 3) % 11;
    return seed < 2 ? 0 : seed < 5 ? seed - 1 : seed < 9 ? 4 + (dayNum % 3) : 2;
  };
  const tonesForDay = (dayNum) => {
    const n = eventCount(dayNum);
    const palette = ["#6B4E3D", "#3C5C7A", "#B8842A", "#A0463A", "#4A5D4F"];
    return Array.from({ length: Math.min(n, 4) }, (_, i) => palette[(dayNum + i) % palette.length]);
  };
  const todayNum = 19;
  const dows = ["dom", "seg", "ter", "qua", "qui", "sex", "sáb"];
  return (
    <section style={{
      background: "var(--paper-2)", border: "1px solid var(--line)",
      borderRadius: "var(--radius)", overflow: "hidden", boxShadow: "var(--shadow-sm)",
    }}>
      {/* Header dias da semana */}
      <div style={{
        display: "grid", gridTemplateColumns: "repeat(7, 1fr)",
        borderBottom: "1px solid var(--line)", background: "var(--paper)",
      }}>
        {dows.map((d, i) => (
          <div key={i} style={{
            padding: "12px 14px", textAlign: "left",
            borderLeft: i > 0 ? "1px solid var(--line)" : "none",
            fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase",
            color: "var(--ink-3)", fontWeight: 500,
          }}>{d}</div>
        ))}
      </div>
      {/* Grid de 6 semanas */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(7, 1fr)", gridTemplateRows: "repeat(6, minmax(100px, 1fr))" }}>
        {cells.map((c, i) => {
          const isToday = c.inMonth && c.dayNum === todayNum;
          const tones = c.inMonth ? tonesForDay(c.dayNum) : [];
          const n = c.inMonth ? eventCount(c.dayNum) : 0;
          const row = Math.floor(i / 7), col = i % 7;
          return (
            <div key={i} style={{
              padding: "10px 12px",
              borderTop: row > 0 ? "1px solid var(--line)" : "none",
              borderLeft: col > 0 ? "1px solid var(--line)" : "none",
              background: isToday ? "color-mix(in oklab, var(--accent) 8%, transparent)" : "transparent",
              opacity: c.inMonth ? 1 : 0.35,
              display: "flex", flexDirection: "column", gap: 8,
              minHeight: 100,
            }}>
              <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
                <span className="serif" style={{
                  fontSize: 18, fontWeight: 500, letterSpacing: -0.3,
                  color: isToday ? "var(--accent-deep)" : c.inMonth ? "var(--ink)" : "var(--ink-4)",
                  minWidth: 26,
                  display: "inline-flex", alignItems: "center", justifyContent: "center",
                  width: isToday ? 26 : "auto", height: isToday ? 26 : "auto",
                  borderRadius: isToday ? "50%" : 0,
                  background: isToday ? "var(--accent)" : "transparent",
                  color: isToday ? "var(--paper)" : undefined,
                }}>{c.inMonth ? c.dayNum : ""}</span>
                {n > 0 && <span style={{ fontSize: 10, color: "var(--ink-4)", marginLeft: "auto" }}>{n} {n === 1 ? "sessão" : "sessões"}</span>}
              </div>
              {/* Pontos de eventos */}
              {c.inMonth && n > 0 && (
                <div style={{ display: "flex", flexDirection: "column", gap: 3 }}>
                  {tones.map((t, ti) => (
                    <div key={ti} style={{
                      display: "flex", alignItems: "center", gap: 6,
                      fontSize: 10.5, color: "var(--ink-3)",
                    }}>
                      <span style={{ width: 5, height: 5, borderRadius: "50%", background: t, flexShrink: 0 }} />
                      <span style={{ whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                        {ti === 0 ? "09:00 Amanda" : ti === 1 ? "11:00 André" : ti === 2 ? "14:00 Carol" : "16:30 Felipe"}
                      </span>
                    </div>
                  ))}
                  {n > 4 && <span style={{ fontSize: 10, color: "var(--ink-4)" }}>+{n - 4} mais</span>}
                </div>
              )}
            </div>
          );
        })}
      </div>
      {/* Legenda */}
      <div style={{ padding: "12px 16px", borderTop: "1px solid var(--line)", display: "flex", gap: 16, flexWrap: "wrap", alignItems: "center", fontSize: 11, color: "var(--ink-3)" }}>
        <LegendDot tone="accent" label="Retorno" />
        <LegendDot tone="info" label="1ª consulta" />
        <LegendDot tone="warn" label="Atenção" />
        <LegendDot tone="danger" label="Risco" />
        <LegendDot tone="moss" label="Supervisão" />
      </div>
    </section>
  );
};


const AgendaEvent = ({ event }) => {
  // Paleta de categorias — cada tipo com matiz distinto (marrom, azul, ocre, vermelho, verde, cinza)
  const toneStyles = {
    accent: { bg: "#EADFCB",                   bd: "#6B4E3D", fg: "#3E2A1E" }, // Retorno — marrom terroso
    info:   { bg: "#DDE6EE",                   bd: "#3C5C7A", fg: "#1F3A55" }, // 1ª consulta — azul sóbrio
    warn:   { bg: "#F3E3C4",                   bd: "#B8842A", fg: "#6B4A10" }, // Atenção — âmbar/ocre
    danger: { bg: "#EDCFC7",                   bd: "#A0463A", fg: "#6F241A" }, // Risco — vermelho profundo
    moss:   { bg: "#D9E2D3",                   bd: "#4A5D4F", fg: "#263326" }, // Supervisão — verde musgo
    ghost:  { bg: "repeating-linear-gradient(135deg, var(--paper-3), var(--paper-3) 4px, var(--paper-2) 4px, var(--paper-2) 8px)",
              bd: "var(--ink-4)", fg: "var(--ink-3)" },                        // Bloco pessoal — cinza hachurado
    neutral:{ bg: "var(--paper)",              bd: "var(--ink-4)", fg: "var(--ink-2)" },
  };
  const t = toneStyles[event.tone] || toneStyles.neutral;
  const top = event.start * 32;
  const h = event.dur * 32;
  const done = event.status === "done";
  const next = event.status === "next";
  return (
    <div style={{
      position: "absolute", left: 4, right: 4,
      top, height: h - 2,
      background: t.bg,
      borderLeft: `3px solid ${t.bd}`,
      borderRadius: 6,
      padding: "5px 8px",
      fontSize: 11,
      color: t.fg,
      overflow: "hidden",
      display: "flex", flexDirection: "column", gap: 2,
      opacity: done ? 0.5 : 1,
      boxShadow: next ? `0 0 0 2px color-mix(in oklab, ${t.bd} 30%, transparent)` : "none",
    }}>
      <div style={{ fontSize: 12, fontWeight: 500, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>{event.patient}</div>
      <div style={{ fontSize: 10, opacity: 0.75, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>{event.type}</div>
    </div>
  );
};

const LegendDot = ({ tone, label }) => {
  const colors = {
    accent: "#6B4E3D",
    info:   "#3C5C7A",
    warn:   "#B8842A",
    danger: "#A0463A",
    moss:   "#4A5D4F",
    ghost:  "var(--ink-4)",
  };
  return (
    <span style={{ display: "inline-flex", alignItems: "center", gap: 6 }}>
      <span style={{ width: 10, height: 10, borderRadius: 2, background: colors[tone] }} />
      {label}
    </span>
  );
};

const navBtn = () => ({
  background: "var(--paper-2)", border: "1px solid var(--line)",
  padding: 8, borderRadius: 10, color: "var(--ink-2)",
  display: "flex", alignItems: "center", justifyContent: "center",
});
const miniBtn = (primary) => ({
  padding: "5px 10px", borderRadius: 7, fontSize: 11.5,
  background: primary ? "var(--ink)" : "transparent",
  color: primary ? "var(--paper)" : "var(--ink-2)",
  border: "1px solid " + (primary ? "var(--ink)" : "var(--line)"),
  fontWeight: 500,
});

window.AgendaPage = AgendaPage;
