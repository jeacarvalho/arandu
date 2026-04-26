// Página: Prontuários — biblioteca editorial de registros clínicos
const NotesPage = ({ data, onOpenPatient }) => {
  const r = data.records;
  const [focused, setFocused] = React.useState(r.focused.id);
  const f = r.focused;
  const [tab, setTab] = React.useState("evolucao");

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 22 }}>
      {/* Hero */}
      <div style={{
        display: "grid", gridTemplateColumns: "1fr auto", gap: 32, alignItems: "end",
        paddingBottom: 22, borderBottom: "1px solid var(--line)",
      }}>
        <div>
          <div style={{ fontSize: 11, letterSpacing: 1.6, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 10 }}>Biblioteca clínica</div>
          <h1 className="serif" style={{ margin: 0, fontSize: 40, fontWeight: 400, letterSpacing: -0.8, lineHeight: 1 }}>
            Prontuários<em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>.</em>
          </h1>
          <p style={{ margin: "8px 0 0", color: "var(--ink-3)", fontSize: 14 }}>
            {r.records.length} registros · {r.records.reduce((s, x) => s + x.pages, 0)} páginas indexadas
          </p>
        </div>
        <div style={{ display: "flex", gap: 8 }}>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "var(--paper-2)", border: "1px solid var(--line)",
            color: "var(--ink-2)", fontSize: 13, display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="filter" size={14} /> Filtros
          </button>
          <button style={{
            padding: "9px 14px", borderRadius: 10,
            background: "var(--ink)", color: "var(--paper)", border: "1px solid var(--ink)",
            fontSize: 13, fontWeight: 500, display: "flex", alignItems: "center", gap: 8,
          }}>
            <Icon name="plus" size={14} /> Novo prontuário
          </button>
        </div>
      </div>

      {/* Filter chips */}
      <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
        <button style={chip(true)}>Tudo</button>
        {r.filterTags.map(t => <button key={t} style={chip(false)}>{t}</button>)}
      </div>

      {/* Split view: lista + visualizador */}
      <div style={{ display: "grid", gridTemplateColumns: "minmax(0, 1fr) minmax(0, 1.3fr)", gap: 18 }}>
        {/* Lista */}
        <section style={{
          background: "var(--paper-2)",
          border: "1px solid var(--line)",
          borderRadius: "var(--radius)",
          overflow: "hidden",
          boxShadow: "var(--shadow-sm)",
        }}>
          <header style={{
            padding: "14px 18px",
            borderBottom: "1px solid var(--line)",
            display: "flex", alignItems: "center", gap: 10,
          }}>
            <Icon name="search" size={14} style={{ color: "var(--ink-3)" }} />
            <input placeholder="Buscar por nome, tag, trecho…" style={{
              flex: 1, background: "transparent", border: 0, outline: "none",
              color: "var(--ink)", fontSize: 13, fontFamily: "inherit",
            }} />
            <span style={{ fontSize: 11, color: "var(--ink-4)" }}>{r.records.length}</span>
          </header>
          <div style={{ maxHeight: 640, overflow: "auto" }}>
            {r.records.map((rec, i) => (
              <button key={rec.id} onClick={() => setFocused(rec.id)} style={{
                width: "100%", textAlign: "left",
                display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 12,
                alignItems: "center", padding: "14px 18px",
                borderBottom: i === r.records.length - 1 ? "none" : "1px dashed var(--line)",
                background: focused === rec.id ? "color-mix(in oklab, var(--accent) 8%, transparent)" : "transparent",
                border: 0, cursor: "pointer", position: "relative",
              }}>
                {focused === rec.id && <span style={{ position: "absolute", left: 0, top: 0, bottom: 0, width: 3, background: "var(--accent)" }} />}
                <Avatar initials={rec.initials} size={36} tone={rec.risk === "Atenção" ? "moss" : "accent"} />
                <div style={{ minWidth: 0 }}>
                  <div style={{ display: "flex", alignItems: "center", gap: 6, marginBottom: 2 }}>
                    <span style={{ fontSize: 14, fontWeight: 500 }}>{rec.patient}</span>
                    {rec.pinned && <Icon name="pin" size={11} style={{ color: "var(--accent)" }} />}
                  </div>
                  <div style={{ fontSize: 11.5, color: "var(--ink-3)", display: "flex", gap: 8, alignItems: "center" }}>
                    <span className="mono">{rec.id}</span>
                    <span>·</span>
                    <span>{rec.pages}p</span>
                    <span>·</span>
                    <span>atualiz. {rec.lastUpdate}</span>
                  </div>
                  <div style={{ marginTop: 6, display: "flex", gap: 6, flexWrap: "wrap" }}>
                    {rec.tags.map(t => <Pill key={t} tone="neutral" size="xs">{t}</Pill>)}
                    <Pill tone={rec.risk === "Atenção" ? "warn" : rec.risk === "Baixo" ? "ok" : "neutral"} size="xs">{rec.status}</Pill>
                  </div>
                </div>
                <Icon name="chevRight" size={14} style={{ color: "var(--ink-4)" }} />
              </button>
            ))}
          </div>
        </section>

        {/* Visualizador editorial */}
        <section style={{
          background: "var(--paper-2)",
          border: "1px solid var(--line)",
          borderRadius: "var(--radius)",
          boxShadow: "var(--shadow-sm)",
          display: "flex", flexDirection: "column",
          overflow: "hidden",
        }}>
          <header style={{
            padding: "18px 22px",
            borderBottom: "1px solid var(--line)",
            display: "grid", gridTemplateColumns: "1fr auto", gap: 16, alignItems: "center",
          }}>
            <div>
              <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 4 }}>
                <span className="mono">{f.id}</span> · André Barbosa
              </div>
              <h2 className="serif" style={{ margin: 0, fontSize: 24, fontWeight: 500, letterSpacing: -0.3 }}>
                Prontuário <em style={{ fontStyle: "italic", color: "var(--accent-deep)" }}>completo</em>
              </h2>
            </div>
            <div style={{ display: "flex", gap: 6 }}>
              <button style={iconBtn()}><Icon name="sparkles" size={14} /></button>
              <button style={iconBtn()}><Icon name="edit" size={14} /></button>
              <button style={iconBtn()}><Icon name="book" size={14} /></button>
            </div>
          </header>

          {/* Tabs */}
          <nav style={{
            display: "flex", borderBottom: "1px solid var(--line)",
            background: "color-mix(in oklab, var(--paper) 60%, var(--paper-2))",
          }}>
            {f.sections.map(s => (
              <button key={s.key} onClick={() => setTab(s.key)} style={{
                padding: "12px 18px",
                background: "transparent",
                border: 0, borderBottom: "2px solid " + (tab === s.key ? "var(--accent)" : "transparent"),
                color: tab === s.key ? "var(--ink)" : "var(--ink-3)",
                fontSize: 13, fontWeight: tab === s.key ? 500 : 400,
                display: "flex", alignItems: "center", gap: 8,
              }}>
                {s.label}
                <span style={{ fontSize: 10, color: "var(--ink-4)" }}>{s.pages}p</span>
              </button>
            ))}
          </nav>

          {/* Página editorial */}
          <div style={{
            flex: 1, padding: "28px 36px 40px",
            overflow: "auto",
            background: "var(--paper)",
            fontFamily: "var(--font-serif)",
            maxHeight: 640,
          }}>
            <div style={{ maxWidth: 620, margin: "0 auto" }}>
              <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 8 }}>
                Evolução · última entrada em 02/04/2026
              </div>
              <h3 className="serif" style={{ margin: 0, fontSize: 28, fontWeight: 400, letterSpacing: -0.5, lineHeight: 1.15 }}>
                Sessão 05 — sentido, existência e padrões emergentes
              </h3>
              <div className="dashed-divider" style={{ margin: "18px 0 22px" }} />

              <p className="serif" style={{ fontSize: 16, lineHeight: 1.65, color: "var(--ink-2)", marginTop: 0 }}>
                {f.lastEntry}
              </p>

              <blockquote style={{
                margin: "22px 0", padding: "4px 0 4px 20px",
                borderLeft: "2px solid var(--accent)",
                fontFamily: "var(--font-serif)", fontStyle: "italic",
                fontSize: 17, lineHeight: 1.55, color: "var(--ink-2)",
              }}>
                "Para que serviu tudo isso?" — frase trazida pelo paciente ao final da sessão, reiterando o questionamento existencial das últimas semanas.
              </blockquote>

              <h4 className="serif" style={{ fontSize: 18, fontWeight: 500, marginTop: 28, marginBottom: 10 }}>Observações</h4>
              <ul className="serif" style={{ fontSize: 15, lineHeight: 1.6, color: "var(--ink-2)", paddingLeft: 20, margin: 0 }}>
                <li>Questionamento existencial intenso sobre sentido da vida e insatisfação com conquistas aparentes.</li>
                <li>Insight emergente sobre padrões automáticos; consciência metacognitiva demonstrada.</li>
                <li>Impulsividade relatada: gastos excessivos e múltiplos projetos simultâneos.</li>
              </ul>

              <h4 className="serif" style={{ fontSize: 18, fontWeight: 500, marginTop: 28, marginBottom: 10 }}>Intervenções</h4>
              <ul className="serif" style={{ fontSize: 15, lineHeight: 1.6, color: "var(--ink-2)", paddingLeft: 20, margin: 0 }}>
                <li>Estabelecimento de limites de disponibilidade entre papéis profissional e pessoal.</li>
                <li>Escuta do questionamento existencial e validação da busca por sentido.</li>
              </ul>

              <div style={{
                marginTop: 32, padding: "14px 18px",
                background: "color-mix(in oklab, var(--accent) 6%, transparent)",
                border: "1px solid color-mix(in oklab, var(--accent) 20%, transparent)",
                borderRadius: 10,
                display: "flex", gap: 12, alignItems: "flex-start",
              }}>
                <Icon name="sparkles" size={16} style={{ color: "var(--accent-deep)", marginTop: 2 }} />
                <div>
                  <div style={{ fontSize: 11, letterSpacing: 1.2, textTransform: "uppercase", color: "var(--accent-deep)", fontWeight: 500, marginBottom: 4, fontFamily: "var(--font-ui)" }}>
                    Leitura cruzada de Arandu
                  </div>
                  <p style={{ margin: 0, fontSize: 14, color: "var(--ink-2)", lineHeight: 1.55 }}>
                    Este tema aparece em 4 das últimas 5 sessões. Em pacientes com perfil semelhante, questionamentos existenciais costumam preceder momentos de reorganização profissional.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <footer style={{
            padding: "12px 22px",
            borderTop: "1px solid var(--line)",
            display: "flex", alignItems: "center", gap: 12,
            background: "color-mix(in oklab, var(--paper) 60%, var(--paper-2))",
          }}>
            <span className="mono" style={{ fontSize: 11, color: "var(--ink-3)" }}>pág. 12 de 12</span>
            <div style={{ flex: 1 }} />
            <button style={footerBtn()}>Exportar PDF</button>
            <button style={{ ...footerBtn(), background: "var(--ink)", color: "var(--paper)", borderColor: "var(--ink)" }}>Adicionar entrada</button>
          </footer>
        </section>
      </div>
    </div>
  );
};

const chip = (active) => ({
  padding: "6px 14px", borderRadius: 999, fontSize: 12, fontWeight: 500,
  background: active ? "var(--ink)" : "transparent",
  border: "1px solid " + (active ? "var(--ink)" : "var(--line)"),
  color: active ? "var(--paper)" : "var(--ink-3)",
});
const iconBtn = () => ({
  background: "transparent", border: "1px solid var(--line)",
  padding: 8, borderRadius: 8, color: "var(--ink-2)", display: "flex",
});
const footerBtn = () => ({
  padding: "7px 12px", borderRadius: 8, fontSize: 12, fontWeight: 500,
  background: "var(--paper)", border: "1px solid var(--line)", color: "var(--ink-2)",
});

window.NotesPage = NotesPage;
