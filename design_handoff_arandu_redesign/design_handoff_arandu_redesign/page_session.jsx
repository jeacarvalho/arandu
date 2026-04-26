// Página: Sessão (notas clínicas lado a lado)
const SessionPage = ({ data, onBack }) => {
  const s = data.session;
  const [obsInput, setObsInput] = React.useState("");
  const [intInput, setIntInput] = React.useState("");

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 22 }}>
      {/* Cabeçalho da sessão */}
      <div style={{
        display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 20, alignItems: "center",
        paddingBottom: 22, borderBottom: "1px solid var(--line)",
      }}>
        <button onClick={onBack} style={{
          background: "var(--paper-2)", border: "1px solid var(--line)",
          padding: "8px 10px", borderRadius: 10, color: "var(--ink-2)",
          display: "flex", alignItems: "center", gap: 6, fontSize: 12.5,
        }}>
          <Icon name="arrowLeft" size={14} /> Voltar ao paciente
        </button>
        <div>
          <div style={{ fontSize: 11, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 6 }}>
            Sessão {String(s.number).padStart(2, "0")} · {s.patient}
          </div>
          <h1 className="serif" style={{ margin: 0, fontSize: 32, fontWeight: 400, letterSpacing: -0.6, lineHeight: 1.1 }}>
            <em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>Sentido</em>, existência e padrões emergentes
          </h1>
          <div style={{ marginTop: 8, fontSize: 13, color: "var(--ink-3)", display: "flex", gap: 14, alignItems: "center" }}>
            <span className="mono">{s.date}</span>
            <span>·</span>
            <span>{s.time}</span>
            <span>·</span>
            <Pill tone="ok" size="xs">Em rascunho</Pill>
          </div>
        </div>
        <div style={{ display: "flex", gap: 8 }}>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "var(--paper-2)", border: "1px solid var(--line)",
            color: "var(--ink-2)", fontSize: 13,
            display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="mic" size={14} /> Ditar
          </button>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "var(--ink)", border: "1px solid var(--ink)",
            color: "var(--paper)", fontSize: 13, fontWeight: 500,
            display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="check" size={14} /> Finalizar sessão
          </button>
        </div>
      </div>

      {/* Duas colunas: Observações + Intervenções */}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 18 }}>
        <NotesColumn
          kind="observação"
          eyebrow="Escuta"
          title="Observações clínicas"
          subtitle="O que foi percebido durante a sessão"
          tone="accent"
          items={s.observations}
          value={obsInput}
          setValue={setObsInput}
          placeholder="Digite sua percepção clínica aqui…"
        />
        <NotesColumn
          kind="intervenção"
          eyebrow="Ação"
          title="Intervenções terapêuticas"
          subtitle="Técnicas e intervenções realizadas"
          tone="ok"
          items={s.interventions}
          value={intInput}
          setValue={setIntInput}
          placeholder="Descreva a técnica ou intervenção realizada…"
        />
      </div>

      {/* Síntese */}
      <Card eyebrow="Síntese" title="Resumo da sessão" subtitle="Texto contínuo para prontuário">
        <p className="serif" style={{
          margin: 0, fontSize: 17, lineHeight: 1.65, color: "var(--ink-2)",
          letterSpacing: -0.1,
        }}>{s.summary}</p>
        <div style={{ display: "flex", gap: 12, marginTop: 18 }}>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "linear-gradient(135deg, var(--accent-deep), var(--accent))",
            border: "1px solid var(--accent-deep)", color: "var(--paper)",
            display: "flex", alignItems: "center", gap: 8,
            fontSize: 12.5, fontWeight: 500,
          }}>
            <Icon name="sparkles" size={14} /> Gerar síntese com Arandu
          </button>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "transparent", border: "1px solid var(--line)",
            color: "var(--ink-2)", fontSize: 12.5,
            display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="edit" size={14} /> Editar manualmente
          </button>
        </div>
      </Card>
    </div>
  );
};

const NotesColumn = ({ kind, eyebrow, title, subtitle, tone, items, value, setValue, placeholder }) => {
  return (
    <section style={{
      background: "var(--paper-2)",
      border: "1px solid var(--line)",
      borderRadius: "var(--radius)",
      display: "flex", flexDirection: "column",
      overflow: "hidden",
      boxShadow: "var(--shadow-sm)",
    }}>
      <header style={{
        padding: "16px 20px 14px",
        borderBottom: "1px solid var(--line)",
      }}>
        <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500 }}>{eyebrow}</div>
        <h3 className="serif" style={{ margin: "3px 0 0", fontSize: 19, fontWeight: 500, letterSpacing: -0.2 }}>{title}</h3>
        <p style={{ margin: "3px 0 0", fontSize: 12.5, color: "var(--ink-3)" }}>{subtitle}</p>
      </header>

      <div style={{ padding: "14px 20px", display: "flex", flexDirection: "column", gap: 10, flex: 1 }}>
        {items.map((item, i) => (
          <div key={item.id} style={{
            padding: "12px 14px",
            background: "var(--paper)",
            border: "1px solid var(--line)",
            borderRadius: 10,
            display: "flex", flexDirection: "column", gap: 8,
            position: "relative",
          }}>
            <div className="serif" style={{ fontSize: 14.5, lineHeight: 1.55, color: "var(--ink)", letterSpacing: -0.05 }}>
              {item.text}
            </div>
            <div style={{ display: "flex", alignItems: "center", gap: 8, fontSize: 11, color: "var(--ink-3)" }}>
              <Icon name="clock" size={11} />
              <span className="mono">{item.timestamp}</span>
              <div style={{ flex: 1 }} />
              <Pill tone={tone} size="xs">{item.tag}</Pill>
              <button style={{ background: "transparent", border: 0, color: "var(--ink-4)", padding: 2, display: "flex" }}>
                <Icon name="edit" size={12} />
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Input */}
      <div style={{
        padding: "14px 20px 18px",
        borderTop: "1px dashed var(--line-2)",
        background: "color-mix(in oklab, var(--paper) 50%, var(--paper-2))",
      }}>
        <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 8, display: "flex", alignItems: "center", gap: 6 }}>
          <Icon name="plus" size={12} /> Nova {kind}
        </div>
        <textarea
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder={placeholder}
          rows={2}
          className="serif"
          style={{
            width: "100%", resize: "none",
            background: "transparent", border: 0, outline: "none",
            fontSize: 14.5, lineHeight: 1.55, color: "var(--ink)",
            fontFamily: "var(--font-serif)",
            fontStyle: value ? "normal" : "italic",
          }}
        />
        <div style={{ display: "flex", alignItems: "center", gap: 8, marginTop: 8 }}>
          <Pill tone="neutral" size="xs">#{kind}</Pill>
          <div style={{ flex: 1 }} />
          <span style={{ fontSize: 11, color: "var(--ink-4)" }}><Kbd>⌘↵</Kbd> registrar</span>
        </div>
      </div>
    </section>
  );
};

window.SessionPage = SessionPage;
