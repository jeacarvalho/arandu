// Ícones em linha fina (1.4px), estilo editorial, inspirados em Lucide mas customizados
const Icon = ({ name, size = 18, stroke = 1.5, className = "", style = {} }) => {
  const common = {
    width: size, height: size, viewBox: "0 0 24 24",
    fill: "none", stroke: "currentColor", strokeWidth: stroke,
    strokeLinecap: "round", strokeLinejoin: "round",
    className, style,
  };
  const paths = {
    dashboard: <><path d="M3 13h8V3H3z"/><path d="M13 21h8V11h-8z"/><path d="M3 21h8v-4H3z"/><path d="M13 7h8V3h-8z"/></>,
    patients:  <><circle cx="9" cy="8" r="3.5"/><path d="M3 20c0-3.3 2.7-6 6-6s6 2.7 6 6"/><circle cx="17" cy="7" r="2.5"/><path d="M15 14c3 0 5 2 5 5"/></>,
    calendar:  <><rect x="3" y="5" width="18" height="16" rx="2"/><path d="M3 9h18M8 3v4M16 3v4"/><circle cx="8" cy="14" r=".8" fill="currentColor" stroke="none"/><circle cx="12" cy="14" r=".8" fill="currentColor" stroke="none"/></>,
    notes:     <><path d="M6 3h9l4 4v14H6z"/><path d="M14 3v5h5"/><path d="M9 13h7M9 17h5"/></>,
    spark:     <><path d="M12 3v4M12 17v4M3 12h4M17 12h4M5.6 5.6l2.8 2.8M15.6 15.6l2.8 2.8M5.6 18.4l2.8-2.8M15.6 8.4l2.8-2.8"/></>,
    search:    <><circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/></>,
    plus:      <><path d="M12 5v14M5 12h14"/></>,
    chevRight: <><path d="m9 6 6 6-6 6"/></>,
    chevDown:  <><path d="m6 9 6 6 6-6"/></>,
    arrowLeft: <><path d="M19 12H5M12 5l-7 7 7 7"/></>,
    arrowRight:<><path d="M5 12h14M12 5l7 7-7 7"/></>,
    close:     <><path d="M18 6 6 18M6 6l12 12"/></>,
    check:     <><path d="M20 6 9 17l-5-5"/></>,
    edit:      <><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 1 1 3 3L7 19l-4 1 1-4z"/></>,
    settings:  <><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.7 1.7 0 0 0-1.8-.3 1.7 1.7 0 0 0-1 1.5V21a2 2 0 1 1-4 0v-.1a1.7 1.7 0 0 0-1.1-1.5 1.7 1.7 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.7 1.7 0 0 0 .3-1.8 1.7 1.7 0 0 0-1.5-1H3a2 2 0 1 1 0-4h.1a1.7 1.7 0 0 0 1.5-1.1 1.7 1.7 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.7 1.7 0 0 0 1.8.3H9a1.7 1.7 0 0 0 1-1.5V3a2 2 0 1 1 4 0v.1a1.7 1.7 0 0 0 1 1.5 1.7 1.7 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.7 1.7 0 0 0-.3 1.8V9a1.7 1.7 0 0 0 1.5 1H21a2 2 0 1 1 0 4h-.1a1.7 1.7 0 0 0-1.5 1z"/></>,
    logout:    <><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><path d="m10 17-5-5 5-5"/><path d="M15 12H5"/></>,
    user:      <><circle cx="12" cy="8" r="4"/><path d="M4 21c0-4.4 3.6-8 8-8s8 3.6 8 8"/></>,
    clock:     <><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></>,
    tag:       <><path d="M20 12 12 20l-8-8V4h8z"/><circle cx="8" cy="8" r="1.2"/></>,
    flag:      <><path d="M4 21V4h13l-2 5 2 5H4"/></>,
    dot:       <><circle cx="12" cy="12" r="3" fill="currentColor"/></>,
    star:      <><path d="m12 3 2.8 5.7 6.2.9-4.5 4.4 1 6.2L12 17.3 6.5 20.2l1-6.2L3 9.6l6.2-.9z"/></>,
    book:      <><path d="M4 4h10a4 4 0 0 1 4 4v12H8a4 4 0 0 1-4-4z"/><path d="M4 16a4 4 0 0 1 4-4h10"/></>,
    brain:     <><path d="M9 3a3 3 0 0 0-3 3v0a3 3 0 0 0-2 5 3 3 0 0 0 0 4 3 3 0 0 0 3 4 3 3 0 0 0 5 0V6a3 3 0 0 0-3-3z"/><path d="M15 3a3 3 0 0 1 3 3v0a3 3 0 0 1 2 5 3 3 0 0 1 0 4 3 3 0 0 1-3 4 3 3 0 0 1-5 0"/><path d="M9 8h2M13 12h2M9 16h2"/></>,
    sun:       <><circle cx="12" cy="12" r="4"/><path d="M12 3v2M12 19v2M3 12h2M19 12h2M5.6 5.6l1.4 1.4M17 17l1.4 1.4M5.6 18.4 7 17M17 7l1.4-1.4"/></>,
    moon:      <><path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z"/></>,
    filter:    <><path d="M3 5h18l-7 9v5l-4 2v-7z"/></>,
    sparkles:  <><path d="M9 3v4M7 5h4M18 9v3M16.5 10.5h3"/><path d="m12 11 1.5 3.5L17 16l-3.5 1.5L12 21l-1.5-3.5L7 16l3.5-1.5z"/></>,
    send:      <><path d="m22 2-20 9 9 2 2 9z"/><path d="m22 2-11 11"/></>,
    pin:       <><path d="M12 17v5"/><path d="M9 3h6v4l3 4H6l3-4z"/><path d="M6 11h12"/></>,
    trend:     <><path d="M3 17 9 11l4 4 8-8"/><path d="M14 7h7v7"/></>,
    alert:     <><path d="M12 3 2 21h20z"/><path d="M12 10v5M12 18v.5"/></>,
    menu:      <><path d="M4 7h16M4 12h16M4 17h16"/></>,
    expand:    <><path d="M4 4h6M4 4v6M20 20h-6M20 20v-6"/></>,
    collapse:  <><path d="M10 4H4v6M14 20h6v-6"/></>,
    mic:       <><rect x="9" y="3" width="6" height="12" rx="3"/><path d="M5 11a7 7 0 0 0 14 0M12 18v3"/></>,
    link:      <><path d="M10 14a5 5 0 0 0 7 0l3-3a5 5 0 0 0-7-7l-1 1"/><path d="M14 10a5 5 0 0 0-7 0l-3 3a5 5 0 0 0 7 7l1-1"/></>,
  };
  return <svg {...common}>{paths[name] || null}</svg>;
};

window.Icon = Icon;
