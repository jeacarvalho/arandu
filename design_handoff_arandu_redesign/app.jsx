// App principal

const DEFAULTS = /*EDITMODE-BEGIN*/{
  "theme": "day",
  "palette": "sage",
  "density": "cozy",
  "sidebarStyle": "rail"
}/*EDITMODE-END*/;

function App() {
  const [page, setPage] = React.useState(() => localStorage.getItem("arandu:page") || "dashboard");
  const [sidebarCollapsed, setSidebarCollapsed] = React.useState(false);
  const [llmOpen, setLlmOpen] = React.useState(false);
  const [tweaksOpen, setTweaksOpen] = React.useState(false);
  const [editMode, setEditMode] = React.useState(false);
  const [tweaks, setTweaks] = React.useState(DEFAULTS);

  // persiste página
  React.useEffect(() => { localStorage.setItem("arandu:page", page); }, [page]);

  // aplica tweaks no body
  React.useEffect(() => {
    document.body.dataset.theme = tweaks.theme;
    document.body.dataset.palette = tweaks.palette;
    document.body.dataset.density = tweaks.density;
    if (tweaks.sidebarStyle === "compact") setSidebarCollapsed(true);
  }, [tweaks]);

  // Tweaks edit-mode protocol
  React.useEffect(() => {
    const onMsg = (e) => {
      const d = e.data || {};
      if (d.type === "__activate_edit_mode") { setEditMode(true); setTweaksOpen(true); }
      if (d.type === "__deactivate_edit_mode") { setEditMode(false); setTweaksOpen(false); }
    };
    window.addEventListener("message", onMsg);
    window.parent.postMessage({ type: "__edit_mode_available" }, "*");
    return () => window.removeEventListener("message", onMsg);
  }, []);

  // Atalho ⌘J / Ctrl+J abre LLM
  React.useEffect(() => {
    const onKey = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "j") {
        e.preventDefault(); setLlmOpen(o => !o);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const setTweak = (k, v) => {
    setTweaks(t => {
      const next = { ...t, [k]: v };
      window.parent.postMessage({ type: "__edit_mode_set_keys", edits: { [k]: v } }, "*");
      return next;
    });
  };

  const data = window.ARANDU_DATA;
  const contextLabel = page === "patients" ? "André Barbosa" : page === "session" ? "Sessão 5 · André" : "Clínica geral";

  return (
    <div style={{ display: "flex", minHeight: "100vh", background: "var(--paper)" }} data-screen-label={`01 ${page}`}>
      <Sidebar page={page} setPage={setPage} collapsed={sidebarCollapsed} setCollapsed={setSidebarCollapsed} />
      <div style={{ flex: 1, display: "flex", flexDirection: "column", minWidth: 0 }}>
        <Topbar
          page={page}
          onOpenLLM={() => setLlmOpen(true)}
          onToggleTweaks={() => setTweaksOpen(o => !o)}
        />
        <main style={{
          flex: 1, padding: "32px 40px 60px",
          maxWidth: 1320, width: "100%", margin: "0 auto",
          display: "flex", flexDirection: "column",
        }}>
          {page === "dashboard" && <DashboardPage data={data} onOpenPatient={(id) => setPage("patients")} />}
          {page === "patients"  && <PatientPage data={data} onOpenSession={() => setPage("session")} onBack={() => setPage("dashboard")} />}
          {page === "session"   && <SessionPage data={data} onBack={() => setPage("patients")} />}
          {page === "agenda"    && <AgendaPage data={data} />}
          {page === "notes"     && <NotesPage data={data} />}
          {page === "insights"  && <InsightsPage data={data} />}
        </main>
      </div>

      <LLMPanel data={data} open={llmOpen} onClose={() => setLlmOpen(false)} context={contextLabel} />

      <TweaksPanel tweaks={tweaks} setTweak={setTweak} open={tweaksOpen} onClose={() => setTweaksOpen(false)} />

      {/* FAB para abrir Tweaks rapidamente (sempre visível) */}
      {!tweaksOpen && (
        <button onClick={() => setTweaksOpen(true)} style={{
          position: "fixed", bottom: 20, right: 20,
          width: 44, height: 44, borderRadius: "50%",
          background: "var(--paper)",
          border: "1px solid var(--line-2)",
          boxShadow: "var(--shadow-md)",
          color: "var(--ink-2)",
          display: "flex", alignItems: "center", justifyContent: "center",
          zIndex: 30,
        }} title="Tweaks">
          <Icon name="settings" size={17} />
        </button>
      )}
    </div>
  );
}

const Placeholder = ({ label }) => (
  <div style={{
    minHeight: 400, display: "flex", alignItems: "center", justifyContent: "center",
    flexDirection: "column", gap: 10,
    color: "var(--ink-3)", textAlign: "center",
  }}>
    <div style={{
      width: 56, height: 56, borderRadius: "50%",
      background: "var(--paper-2)", border: "1px dashed var(--line-2)",
      display: "flex", alignItems: "center", justifyContent: "center",
    }}>
      <Icon name="spark" size={24} />
    </div>
    <h2 className="serif" style={{ margin: "12px 0 0", fontSize: 24, fontWeight: 400, color: "var(--ink-2)" }}>{label}</h2>
    <p style={{ margin: 0, fontSize: 13, maxWidth: 320 }}>
      Nesta proposta, priorizamos Dashboard, Perfil do Paciente e Sessão. Esta tela seguirá o mesmo sistema visual.
    </p>
  </div>
);

// Mount: carrega os scripts na ordem correta
const mount = () => {
  const root = ReactDOM.createRoot(document.getElementById("root"));
  root.render(<App />);
};

// Espera todos os globals carregarem
const waitForGlobals = (names, cb) => {
  const check = () => names.every(n => window[n]) ? cb() : setTimeout(check, 20);
  check();
};

waitForGlobals([
  "ARANDU_DATA", "Icon", "Sidebar", "Topbar", "Card", "Pill", "Avatar", "Kbd",
  "DashboardPage", "PatientPage", "SessionPage", "AgendaPage", "NotesPage", "InsightsPage", "LLMPanel", "TweaksPanel",
], mount);
