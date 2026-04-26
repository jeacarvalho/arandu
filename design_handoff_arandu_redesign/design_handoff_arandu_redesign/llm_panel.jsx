// Painel lateral: Inteligência Clínica (LLM) — drawer
const LLMPanel = ({ data, open, onClose, context }) => {
  const [input, setInput] = React.useState("");
  const thread = data.llm.threads[0];

  return (
    <>
      {open && <div onClick={onClose} style={{
        position: "fixed", inset: 0,
        background: "color-mix(in oklab, var(--ink) 40%, transparent)",
        backdropFilter: "blur(2px)",
        zIndex: 40,
      }} />}
      <aside style={{
        position: "fixed", top: 0, right: 0, height: "100vh",
        width: 460, maxWidth: "90vw",
        background: "var(--paper)",
        borderLeft: "1px solid var(--line)",
        boxShadow: "-24px 0 60px -24px rgba(31,26,21,.25)",
        display: "flex", flexDirection: "column",
        zIndex: 41,
        transform: open ? "translateX(0)" : "translateX(100%)",
        transition: "transform .28s cubic-bezier(.4,0,.2,1)",
      }}>
        {/* Header */}
        <header style={{
          padding: "18px 22px 14px",
          borderBottom: "1px solid var(--line)",
          display: "flex", alignItems: "center", gap: 12,
        }}>
          <div style={{
            width: 36, height: 36, borderRadius: 10,
            background: "linear-gradient(135deg, var(--accent-deep), var(--accent))",
            color: "var(--paper)",
            display: "flex", alignItems: "center", justifyContent: "center",
            boxShadow: "inset 0 -2px 3px rgba(0,0,0,.2)",
          }}>
            <Icon name="sparkles" size={17} />
          </div>
          <div style={{ flex: 1 }}>
            <div className="serif" style={{ fontSize: 18, fontWeight: 500, letterSpacing: -0.2 }}>Arandu</div>
            <div style={{ fontSize: 11.5, color: "var(--ink-3)" }}>
              Inteligência clínica · analisando <strong style={{ color: "var(--ink-2)", fontWeight: 500 }}>{context}</strong>
            </div>
          </div>
          <button onClick={onClose} style={{
            background: "transparent", border: "1px solid var(--line)",
            padding: 7, borderRadius: 8, color: "var(--ink-2)",
          }}>
            <Icon name="close" size={14} />
          </button>
        </header>

        {/* Insights cruzados */}
        <div style={{
          padding: "16px 22px",
          borderBottom: "1px dashed var(--line-2)",
          background: "var(--paper-2)",
        }}>
          <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 10 }}>
            Padrões cruzados
          </div>
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {data.llm.crossInsights.map((c, i) => (
              <div key={i} style={{
                display: "grid", gridTemplateColumns: "110px 1fr auto",
                gap: 12, alignItems: "center",
                fontSize: 12.5,
              }}>
                <span style={{ color: "var(--ink-3)", fontSize: 11, textTransform: "uppercase", letterSpacing: .6 }}>{c.label}</span>
                <div style={{ position: "relative", height: 4, borderRadius: 2, background: "var(--line)", overflow: "hidden" }}>
                  <div style={{
                    position: "absolute", left: 0, top: 0, bottom: 0,
                    width: `${c.weight * 100}%`,
                    background: i === 0 ? "var(--accent)" : i === 1 ? "var(--gold)" : "var(--sage)",
                    borderRadius: 2,
                  }} />
                </div>
                <span className="serif" style={{ fontSize: 13, color: "var(--ink)", textAlign: "right", minWidth: 120 }}>{c.value}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Thread */}
        <div style={{ flex: 1, overflow: "auto", padding: "18px 22px" }}>
          <div style={{ fontSize: 11, color: "var(--ink-3)", marginBottom: 12, letterSpacing: .5 }}>
            <span className="serif" style={{ fontSize: 14, fontStyle: "italic", color: "var(--ink-2)", fontWeight: 500 }}>{thread.title}</span>
            <span style={{ marginLeft: 8 }}>· atualizado {thread.updated}</span>
          </div>

          {thread.messages.map((m, i) => <Message key={i} m={m} />)}
        </div>

        {/* Sugestões rápidas */}
        <div style={{ padding: "10px 22px", borderTop: "1px dashed var(--line-2)" }}>
          <div style={{ display: "flex", flexWrap: "wrap", gap: 6 }}>
            {data.llm.suggestions.map((s, i) => (
              <button key={i} style={{
                fontSize: 11.5, padding: "6px 12px",
                background: "var(--paper-2)", border: "1px solid var(--line)",
                borderRadius: 20, color: "var(--ink-2)",
                fontFamily: "var(--font-serif)", fontStyle: "italic",
              }}>{s}</button>
            ))}
          </div>
        </div>

        {/* Input */}
        <div style={{
          padding: "14px 22px 18px",
          borderTop: "1px solid var(--line)",
          background: "var(--paper-2)",
        }}>
          <div style={{
            display: "flex", alignItems: "flex-end", gap: 10,
            padding: "10px 12px",
            background: "var(--paper)",
            border: "1px solid var(--line-2)",
            borderRadius: 12,
          }}>
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Pergunte sobre padrões, hipóteses, comparações…"
              rows={1}
              style={{
                flex: 1, resize: "none",
                background: "transparent", border: 0, outline: "none",
                fontSize: 13.5, lineHeight: 1.5, color: "var(--ink)",
                fontFamily: "inherit",
                minHeight: 22, maxHeight: 120,
              }}
            />
            <button style={{
              padding: 8, borderRadius: 8,
              background: "var(--ink)", color: "var(--paper)", border: 0,
              display: "flex",
            }}>
              <Icon name="send" size={14} />
            </button>
          </div>
          <div style={{ fontSize: 10.5, color: "var(--ink-4)", marginTop: 8, display: "flex", alignItems: "center", gap: 8 }}>
            <Icon name="alert" size={11} />
            Respostas de IA são auxiliares e não substituem julgamento clínico.
          </div>
        </div>
      </aside>
    </>
  );
};

const Message = ({ m }) => {
  if (m.role === "user") {
    return (
      <div style={{
        marginBottom: 16,
        padding: "10px 14px",
        background: "var(--paper-2)",
        border: "1px solid var(--line)",
        borderRadius: "12px 12px 12px 4px",
        fontSize: 13.5,
        color: "var(--ink)",
        maxWidth: "85%",
      }}>{m.text}</div>
    );
  }
  return (
    <div style={{ marginBottom: 18, display: "flex", gap: 12 }}>
      <div style={{
        width: 26, height: 26, borderRadius: 7, flexShrink: 0,
        background: "linear-gradient(135deg, var(--accent-deep), var(--accent))",
        color: "var(--paper)",
        display: "flex", alignItems: "center", justifyContent: "center",
        marginTop: 2,
      }}>
        <Icon name="sparkles" size={13} />
      </div>
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{
          fontSize: 14, lineHeight: 1.6, color: "var(--ink)",
          whiteSpace: "pre-wrap",
        }} className="serif" dangerouslySetInnerHTML={{
          __html: m.text
            .replace(/\*\*(.+?)\*\*/g, '<strong style="font-weight:600;color:var(--accent-deep)">$1</strong>')
        }} />
        {m.citations && (
          <div style={{ display: "flex", flexWrap: "wrap", gap: 6, marginTop: 12 }}>
            {m.citations.map((c, i) => (
              <span key={i} style={{
                fontSize: 10.5, padding: "3px 8px",
                background: "var(--paper-2)",
                border: "1px solid var(--line)",
                borderRadius: 20,
                color: "var(--ink-3)",
                display: "inline-flex", alignItems: "center", gap: 5,
                fontFamily: "var(--font-mono)",
              }}>
                <Icon name="link" size={10} />
                {c.session} · {c.date}
              </span>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

window.LLMPanel = LLMPanel;
