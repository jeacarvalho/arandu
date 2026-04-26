// Componentes reutilizáveis base (Card, Pill, Kbd, etc) + Sidebar + Topbar

const Pill = ({ children, tone = "neutral", size = "sm", className = "", ...rest }) => {
  const tones = {
    neutral: { bg: "color-mix(in oklab, var(--ink) 6%, transparent)", fg: "var(--ink-2)", bd: "var(--line)" },
    warn:    { bg: "color-mix(in oklab, var(--clay) 14%, transparent)", fg: "#8B3A24", bd: "color-mix(in oklab, var(--clay) 30%, transparent)" },
    ok:      { bg: "color-mix(in oklab, var(--moss-2) 14%, transparent)", fg: "var(--moss)", bd: "color-mix(in oklab, var(--moss-2) 30%, transparent)" },
    info:    { bg: "color-mix(in oklab, var(--accent-soft) 14%, transparent)", fg: "var(--accent-deep)", bd: "color-mix(in oklab, var(--accent-soft) 30%, transparent)" },
    danger:  { bg: "color-mix(in oklab, var(--danger) 14%, transparent)", fg: "var(--danger)", bd: "color-mix(in oklab, var(--danger) 30%, transparent)" },
    accent:  { bg: "color-mix(in oklab, var(--accent) 12%, transparent)", fg: "var(--accent-deep)", bd: "color-mix(in oklab, var(--accent) 28%, transparent)" },
  };
  const t = tones[tone] || tones.neutral;
  return (
    <span className={className} style={{
      display: "inline-flex", alignItems: "center", gap: 6,
      padding: size === "xs" ? "2px 8px" : "3px 10px",
      fontSize: size === "xs" ? 10.5 : 11.5,
      fontWeight: 500,
      letterSpacing: .2,
      textTransform: "uppercase",
      background: t.bg, color: t.fg,
      border: `1px solid ${t.bd}`,
      borderRadius: 999,
      lineHeight: 1,
      ...rest.style,
    }}>{children}</span>
  );
};

const Card = ({ children, title, subtitle, action, eyebrow, style = {}, bodyStyle = {}, noPad = false }) => (
  <section style={{
    background: "var(--paper-2)",
    border: "1px solid var(--line)",
    borderRadius: "var(--radius)",
    boxShadow: "var(--shadow-sm)",
    overflow: "hidden",
    ...style,
  }}>
    {(title || subtitle || action || eyebrow) && (
      <header style={{
        padding: "16px 20px 12px",
        display: "flex", alignItems: "flex-start", justifyContent: "space-between", gap: 16,
        borderBottom: "1px solid var(--line)",
      }}>
        <div>
          {eyebrow && <div style={{
            fontSize: 10.5, letterSpacing: 1.4, textTransform: "uppercase",
            color: "var(--ink-3)", fontWeight: 500, marginBottom: 4,
          }}>{eyebrow}</div>}
          {title && <h3 className="serif" style={{
            margin: 0, fontSize: 19, fontWeight: 500, color: "var(--ink)",
            letterSpacing: -0.2, lineHeight: 1.25,
          }}>{title}</h3>}
          {subtitle && <p style={{ margin: "4px 0 0", color: "var(--ink-3)", fontSize: 13 }}>{subtitle}</p>}
        </div>
        {action}
      </header>
    )}
    <div style={{ padding: noPad ? 0 : "18px 20px", ...bodyStyle }}>{children}</div>
  </section>
);

const Kbd = ({ children }) => (
  <kbd style={{
    fontFamily: "var(--font-mono)",
    fontSize: 10.5,
    padding: "2px 6px",
    border: "1px solid var(--line-2)",
    borderBottomWidth: 2,
    borderRadius: 5,
    background: "var(--paper)",
    color: "var(--ink-3)",
  }}>{children}</kbd>
);

const SidebarItem = ({ icon, label, active, onClick, badge, collapsed }) => (
  <button onClick={onClick} title={collapsed ? label : undefined} style={{
    width: "100%",
    display: "flex", alignItems: "center", gap: 12,
    padding: collapsed ? "10px" : "9px 12px",
    justifyContent: collapsed ? "center" : "flex-start",
    borderRadius: 10,
    background: active ? "color-mix(in oklab, var(--accent) 14%, transparent)" : "transparent",
    border: "1px solid " + (active ? "color-mix(in oklab, var(--accent) 22%, transparent)" : "transparent"),
    color: active ? "var(--accent-deep)" : "var(--ink-2)",
    fontSize: 13.5, fontWeight: active ? 500 : 400,
    textAlign: "left",
    position: "relative",
    transition: "background .15s ease, color .15s ease",
  }}>
    <Icon name={icon} size={18} stroke={active ? 1.8 : 1.5} />
    {!collapsed && <span style={{ flex: 1 }}>{label}</span>}
    {!collapsed && badge && <Pill tone="accent" size="xs">{badge}</Pill>}
  </button>
);

const Sidebar = ({ page, setPage, collapsed, setCollapsed }) => {
  const items = [
    { id: "dashboard", icon: "dashboard", label: "Dashboard" },
    { id: "patients",  icon: "patients",  label: "Pacientes", badge: "42" },
    { id: "agenda",    icon: "calendar",  label: "Agenda" },
    { id: "notes",     icon: "notes",     label: "Prontuários" },
    { id: "insights",  icon: "brain",     label: "Inteligência", badge: "3" },
  ];
  const w = collapsed ? 68 : 232;
  return (
    <aside style={{
      width: w, flexShrink: 0,
      background: "var(--paper)",
      borderRight: "1px solid var(--line)",
      display: "flex", flexDirection: "column",
      position: "sticky", top: 0, height: "100vh",
      transition: "width .2s ease",
    }}>
      {/* Logo / Brand */}
      <div style={{
        padding: collapsed ? "22px 12px" : "22px 20px",
        borderBottom: "1px solid var(--line)",
        display: "flex", alignItems: "center", gap: 10,
      }}>
        <BrandMark size={30} />
        {!collapsed && <div>
          <div className="serif" style={{ fontSize: 20, letterSpacing: -0.3, fontWeight: 500, lineHeight: 1 }}>Arandu</div>
          <div style={{ fontSize: 10.5, color: "var(--ink-3)", letterSpacing: 1.2, textTransform: "uppercase", marginTop: 3 }}>Clínico</div>
        </div>}
      </div>

      {/* Nav */}
      <nav style={{ padding: 14, display: "flex", flexDirection: "column", gap: 4, flex: 1 }}>
        {!collapsed && <div style={{
          fontSize: 10, letterSpacing: 1.4, textTransform: "uppercase",
          color: "var(--ink-4)", padding: "6px 12px 8px", fontWeight: 500,
        }}>Menu</div>}
        {items.map(it => (
          <SidebarItem key={it.id} {...it}
            active={page === it.id}
            onClick={() => setPage(it.id)}
            collapsed={collapsed}
          />
        ))}

        {!collapsed && <div style={{
          fontSize: 10, letterSpacing: 1.4, textTransform: "uppercase",
          color: "var(--ink-4)", padding: "20px 12px 8px", fontWeight: 500,
        }}>Atalhos</div>}
        <SidebarItem icon="plus" label="Nova sessão" onClick={() => {}} collapsed={collapsed} />
        <SidebarItem icon="user" label="Novo paciente" onClick={() => {}} collapsed={collapsed} />
      </nav>

      {/* Rodapé: clínico + colapsar */}
      <div style={{ padding: 14, borderTop: "1px solid var(--line)", display: "flex", flexDirection: "column", gap: 10 }}>
        <div style={{
          display: "flex", alignItems: "center", gap: 10,
          padding: collapsed ? 6 : "8px 10px",
          borderRadius: 10,
          justifyContent: collapsed ? "center" : "flex-start",
        }}>
          <Avatar initials="HM" size={collapsed ? 30 : 32} />
          {!collapsed && <div style={{ minWidth: 0, flex: 1 }}>
            <div style={{ fontSize: 12.5, fontWeight: 500, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>Dra. Helena Moraes</div>
            <div style={{ fontSize: 10.5, color: "var(--ink-3)", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>CRP 06/98432</div>
          </div>}
        </div>
        <button onClick={() => setCollapsed(c => !c)} style={{
          background: "transparent", border: "1px dashed var(--line-2)",
          color: "var(--ink-3)", borderRadius: 8,
          padding: 8, display: "flex", alignItems: "center", justifyContent: "center", gap: 8,
          fontSize: 11.5,
        }}>
          <Icon name={collapsed ? "expand" : "collapse"} size={14} />
          {!collapsed && <span>Recolher</span>}
        </button>
      </div>
    </aside>
  );
};

const BrandMark = ({ size = 30 }) => (
  <div style={{
    width: size, height: size, borderRadius: size/4,
    background: "linear-gradient(135deg, var(--accent), var(--accent-deep))",
    display: "flex", alignItems: "center", justifyContent: "center",
    color: "var(--paper)", boxShadow: "inset 0 -2px 4px rgba(0,0,0,.18), 0 1px 0 rgba(255,255,255,.15)",
    position: "relative", overflow: "hidden", flexShrink: 0,
  }}>
    <svg width={size*0.62} height={size*0.62} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
      {/* Glifo: folha + olho (sabedoria + cuidado) */}
      <path d="M12 3c-4 3-6 6-6 10 0 4 2 7 6 8 4-1 6-4 6-8 0-4-2-7-6-10z" />
      <path d="M12 7v12" opacity=".55" />
      <circle cx="12" cy="13" r="1.4" fill="currentColor" stroke="none" />
    </svg>
  </div>
);

const Avatar = ({ initials, size = 36, tone = "accent" }) => {
  const bg = tone === "accent"
    ? "linear-gradient(135deg, color-mix(in oklab, var(--accent) 80%, var(--paper)), var(--accent))"
    : "linear-gradient(135deg, var(--sage), var(--moss-2))";
  return (
    <div style={{
      width: size, height: size, borderRadius: "50%",
      background: bg,
      color: "var(--paper)",
      display: "flex", alignItems: "center", justifyContent: "center",
      fontFamily: "var(--font-serif)",
      fontWeight: 500,
      fontSize: size * 0.38,
      letterSpacing: .2,
      flexShrink: 0,
      boxShadow: "inset 0 -1px 2px rgba(0,0,0,.15)",
    }}>{initials}</div>
  );
};

const Topbar = ({ onOpenLLM, onToggleTweaks, page }) => {
  const breadcrumb = {
    dashboard: ["Hoje"],
    patients:  ["Pacientes", "André Barbosa"],
    session:   ["Pacientes", "André Barbosa", "Sessão 5"],
    agenda:    ["Agenda"],
    notes:     ["Prontuários"],
    insights:  ["Inteligência Clínica"],
  }[page] || ["Dashboard"];
  return (
    <header style={{
      height: 64, flexShrink: 0,
      borderBottom: "1px solid var(--line)",
      background: "color-mix(in oklab, var(--paper) 88%, transparent)",
      backdropFilter: "blur(8px)",
      display: "flex", alignItems: "center",
      padding: "0 28px", gap: 20,
      position: "sticky", top: 0, zIndex: 20,
    }}>
      {/* Breadcrumb */}
      <nav style={{ display: "flex", alignItems: "center", gap: 8, color: "var(--ink-3)", fontSize: 13 }}>
        {breadcrumb.map((b, i) => (
          <React.Fragment key={i}>
            {i > 0 && <Icon name="chevRight" size={12} />}
            <span style={{
              color: i === breadcrumb.length - 1 ? "var(--ink)" : "var(--ink-3)",
              fontFamily: i === breadcrumb.length - 1 ? "var(--font-serif)" : "var(--font-ui)",
              fontSize: i === breadcrumb.length - 1 ? 15 : 13,
              fontWeight: i === breadcrumb.length - 1 ? 500 : 400,
            }}>{b}</span>
          </React.Fragment>
        ))}
      </nav>

      <div style={{ flex: 1 }} />

      {/* Search */}
      <div style={{
        display: "flex", alignItems: "center", gap: 10,
        width: 340,
        padding: "7px 12px",
        background: "var(--paper-2)",
        border: "1px solid var(--line)",
        borderRadius: 10,
        color: "var(--ink-3)",
      }}>
        <Icon name="search" size={15} />
        <input placeholder="Buscar paciente, sessão ou observação…" style={{
          flex: 1, background: "transparent", border: 0, outline: "none",
          color: "var(--ink)", fontSize: 13, fontFamily: "inherit",
        }} />
        <Kbd>⌘K</Kbd>
      </div>

      {/* LLM Toggle */}
      <button onClick={onOpenLLM} style={{
        display: "flex", alignItems: "center", gap: 8,
        padding: "8px 14px",
        borderRadius: 10,
        background: "linear-gradient(135deg, var(--accent-deep), var(--accent))",
        color: "var(--paper)",
        border: "1px solid color-mix(in oklab, var(--accent-deep) 60%, #000)",
        boxShadow: "0 1px 0 rgba(255,255,255,.2) inset, 0 4px 12px -4px color-mix(in oklab, var(--accent) 70%, transparent)",
        fontSize: 13, fontWeight: 500,
      }}>
        <Icon name="sparkles" size={15} />
        <span>Arandu</span>
        <Kbd>⌘J</Kbd>
      </button>

      <button onClick={onToggleTweaks} title="Configurar" style={{
        background: "transparent", border: "1px solid var(--line)",
        padding: 9, borderRadius: 10, color: "var(--ink-2)",
      }}>
        <Icon name="settings" size={16} />
      </button>
    </header>
  );
};

window.Pill = Pill;
window.Card = Card;
window.Kbd = Kbd;
window.Sidebar = Sidebar;
window.Topbar = Topbar;
window.BrandMark = BrandMark;
window.Avatar = Avatar;
