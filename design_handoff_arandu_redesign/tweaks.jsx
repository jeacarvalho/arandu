// Painel de Tweaks — flutuante no canto
const TweaksPanel = ({ tweaks, setTweak, open, onClose }) => {
  if (!open) return null;
  return (
    <div style={{
      position: "fixed", bottom: 20, right: 20,
      width: 300,
      background: "var(--paper)",
      border: "1px solid var(--line-2)",
      borderRadius: 14,
      boxShadow: "var(--shadow-lg)",
      zIndex: 50,
      overflow: "hidden",
    }}>
      <header style={{
        padding: "12px 16px",
        borderBottom: "1px solid var(--line)",
        display: "flex", alignItems: "center", gap: 10,
        background: "var(--paper-2)",
      }}>
        <Icon name="settings" size={14} />
        <div className="serif" style={{ fontSize: 15, fontWeight: 500, flex: 1 }}>Tweaks</div>
        <button onClick={onClose} style={{ background: "transparent", border: 0, color: "var(--ink-3)", display: "flex" }}>
          <Icon name="close" size={13} />
        </button>
      </header>
      <div style={{ padding: 16, display: "flex", flexDirection: "column", gap: 16 }}>
        <TweakRow label="Paleta">
          <Seg value={tweaks.palette} onChange={v => setTweak("palette", v)} options={[
            { v: "sage", l: "Sábio" },
            { v: "moss", l: "Musgo" },
            { v: "clay", l: "Barro" },
          ]} />
        </TweakRow>
        <TweakRow label="Modo">
          <Seg value={tweaks.theme} onChange={v => setTweak("theme", v)} options={[
            { v: "day", l: "Dia" },
            { v: "night", l: "Noite" },
          ]} />
        </TweakRow>
        <TweakRow label="Densidade">
          <Seg value={tweaks.density} onChange={v => setTweak("density", v)} options={[
            { v: "compact", l: "Denso" },
            { v: "cozy", l: "Equilibrado" },
            { v: "roomy", l: "Amplo" },
          ]} />
        </TweakRow>
        <TweakRow label="Sidebar">
          <Seg value={tweaks.sidebarStyle} onChange={v => setTweak("sidebarStyle", v)} options={[
            { v: "rail", l: "Completa" },
            { v: "compact", l: "Compacta" },
          ]} />
        </TweakRow>
      </div>
    </div>
  );
};

const TweakRow = ({ label, children }) => (
  <div>
    <div style={{ fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase", color: "var(--ink-3)", fontWeight: 500, marginBottom: 6 }}>{label}</div>
    {children}
  </div>
);

const Seg = ({ value, onChange, options }) => (
  <div style={{
    display: "flex", padding: 3, background: "var(--paper-2)",
    border: "1px solid var(--line)", borderRadius: 8,
  }}>
    {options.map(o => (
      <button key={o.v} onClick={() => onChange(o.v)} style={{
        flex: 1, padding: "5px 8px",
        background: value === o.v ? "var(--ink)" : "transparent",
        color: value === o.v ? "var(--paper)" : "var(--ink-3)",
        border: 0, borderRadius: 6,
        fontSize: 11.5, fontWeight: 500,
      }}>{o.l}</button>
    ))}
  </div>
);

window.TweaksPanel = TweaksPanel;
