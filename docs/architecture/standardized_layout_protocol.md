# 🖼️ Standardized Layout Protocol (SLP) — Arandu SOTA

Este documento define a anatomia obrigatória de todas as telas do sistema para garantir consistência visual, ergonomia e imersão clínica.

---

## 🏗️ 1. A Anatomia das 3 Zonas

Todas as páginas devem ser renderizadas através do componente central de layout, dividindo o ecrã em três áreas distintas:

### A. Top Bar (Navegação Universal)
* **Esquerda:** Logo "Arandu" + Gatilho de Menu (Hambúrguer no Mobile).
* **Centro:** Barra de busca de pacientes (`.silent-search`). A busca deve estar sempre presente para permitir saltos rápidos entre prontuários.
* **Direita:** Identificação do Utilizador (Avatar + Nome/Iniciais).

### B. Sidebar (Menu Sensível ao Contexto)
A barra lateral não é estática; ela reflete o **Domínio Atual**.
* **Contexto Geral (Dashboard):** Links para Pacientes, Agenda, Configurações.
* **Contexto Paciente:** Links para Dados Cadastrais, Anamnese, Linha do Tempo, Plano Terapêutico.
* **Contexto Sessão:** Links para Observações, Intervenções, Vitais daquela sessão específica.

### C. Main Canvas (Área de Trabalho)
* Onde o dado de domínio é processado.
* **Regra de Ouro:** Fundo `--arandu-bg` (#E1F5EE), com o conteúdo principal dentro de `.clinical-card` ou sobre o "papel digital".

---

## 📱 2. Comportamento Mobile (Responsividade)

| Zona | Comportamento Mobile |
| :--- | :--- |
| **Top Bar** | Mantém Logo e Avatar. A Busca pode recolher para um ícone de lupa. |
| **Sidebar** | Transforma-se num **Drawer Lateral** (esquerda) acionado pelo hambúrguer. |
| **Main Canvas** | Ocupa 100% da largura. Margens reduzidas de 32px para 16px. |

---

## 🛡️ 3. Mecanismos de Checagem (Anti-Quebra)

Para garantir que novas implementações não fujam do padrão:

1.  **Contrato do Template:** O componente `Layout` deve receber obrigatoriamente um parâmetro `contextID` (string) e um `sidebarComponent` (templ.Component).
2.  **Verificação Manual (Checklist):**
    * [ ] O menu lateral mudou ao entrar num paciente?
    * [ ] A busca na Top Bar continua funcional?
    * [ ] Em 375px, a sidebar está escondida sob o hambúrguer?
3.  **Verificação Automatizada (`arandu_guard.sh`):**
    * O script deve verificar a existência da classe `.top-bar-user-id` e `.contextual-sidebar` no HTML retornado pelas rotas principais.

---

## 📝 4. Matriz de Menus por Contexto

| Rota | Contexto | Ações Sugeridas na Sidebar |
| :--- | :--- | :--- |
| `/dashboard` | `global` | Pacientes, Relatórios Gerais, Perfil. |
| `/patients/{id}/*` | `patient` | Anamnese, Prontuário, Medicamentos, Metas. |
| `/sessions/{id}` | `session` | Nova Observação, Nova Intervenção, Vitais, Fechar Sessão. |