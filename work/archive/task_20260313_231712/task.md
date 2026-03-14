# TASK 20260313_231712

Requirement: req-01-00-01

Title: Patient: Implementar UI/UX

## Objetivo

Implementar a interface de usuário (UI/UX) para a funcionalidade de criação e gerenciamento de pacientes, integrando com o backend completo já implementado.

## Contexto

Nas tarefas anteriores, implementamos completamente o backend para gerenciamento de pacientes:
- ✅ **Domínio:** Entidade Patient com validação
- ✅ **Repository:** PatientRepository com queries otimizadas
- ✅ **Migrations:** Sistema de versionamento de banco
- ✅ **Application Service:** PatientService com input models e validação
- ✅ **Handlers HTTP:** Endpoints básicos já existem

Agora precisamos implementar a interface de usuário que permita:
1. **Listar pacientes:** Visualizar todos os pacientes cadastrados
2. **Criar paciente:** Formulário para cadastro de novo paciente
3. **Visualizar paciente:** Detalhes de um paciente específico
4. **Buscar pacientes:** Funcionalidade de busca por nome
5. **Navegação:** Interface intuitiva e responsiva

## Análise do Estado Atual

**Arquitetura atual do frontend:**
- **Tecnologia:** Go templates (HTML) + HTMX para interatividade
- **Estrutura:** `web/templates/` para templates, `web/static/` para assets
- **Endpoints existentes:** `/patients` (GET), `/patient/{id}` (GET)

**Problemas identificados:**
1. **Templates básicos:** HTML simples sem estilização adequada
2. **Sem formulários:** Não há interface para criação de pacientes
3. **Sem feedback visual:** Não mostra mensagens de sucesso/erro
4. **Sem responsividade:** Não otimizado para mobile
5. **Sem validação client-side:** Apenas validação no backend
6. **Sem experiência de usuário:** Fluxo não é intuitivo

## Tarefas Específicas

### 1. Criar Sistema de Layout/Templates
Implementar sistema de templates com:
- Layout base com header, navigation, footer
- Sistema de partials/templates reutilizáveis
- Suporte a mensagens flash (success, error, warning)
- Sistema de breadcrumbs para navegação

### 2. Implementar Listagem de Pacientes
Criar página `/patients` com:
- Tabela/listagem de pacientes
- Paginação (usando `ListPatientsPaginated`)
- Busca por nome (integrando `SearchPatientsByName`)
- Ordenação por diferentes campos
- Ações: Visualizar, Editar (futuro), Excluir

### 3. Implementar Formulário de Criação
Criar página `/patients/new` com:
- Formulário com campos: Nome (obrigatório), Notas (opcional)
- Validação client-side (JavaScript/HTMX)
- Feedback visual durante submit
- Redirecionamento após criação bem-sucedida
- Mensagens de erro específicas

### 4. Implementar Página de Detalhes
Melhorar página `/patient/{id}` com:
- Informações completas do paciente
- Histórico de sessões (futuro - req-01-01-01)
- Ações: Editar, Excluir, Voltar para lista
- Layout organizado e legível

### 5. Adicionar Estilização (CSS)
Implementar sistema de estilos:
- Framework CSS leve ou custom
- Design responsivo (mobile-first)
- Componentes reutilizáveis (buttons, forms, cards)
- Temas (light/dark mode opcional)

### 6. Adicionar Interatividade (HTMX/JavaScript)
Implementar interações:
- Submit de formulários sem page reload (HTMX)
- Validação em tempo real
- Confirmação para exclusões
- Loading states durante operações
- Notificações/toasts para feedback

### 7. Implementar Navegação e UX
Criar experiência de usuário completa:
- Menu de navegação intuitivo
- Breadcrumbs para contexto
- Estados vazios (empty states)
- Mensagens de ajuda/guidance
- Acessibilidade básica (ARIA labels, keyboard navigation)

## Requisitos Técnicos

### Stack Tecnológica
- **Backend:** Go + HTML templates
- **Frontend:** HTMX + CSS (Tailwind/Bootstrap ou custom)
- **Interatividade:** HTMX para AJAX, JavaScript mínimo para validação
- **Estilos:** CSS com sistema de design consistente

### Estrutura de Templates
```
web/templates/
├── layouts/
│   └── base.html      # Layout principal
├── partials/
│   ├── header.html    # Cabeçalho
│   ├── footer.html    # Rodapé
│   ├── nav.html       # Navegação
│   └── messages.html  # Mensagens flash
├── patients/
│   ├── list.html      # Listagem
│   ├── new.html       # Criação
│   ├── show.html      # Detalhes
│   └── _form.html     # Formulário partial
└── components/
    ├── table.html     # Componente de tabela
    ├── card.html      # Componente de card
    └── pagination.html # Componente de paginação
```

### Design System
**Cores (exemplo):**
- Primária: Azul profissional (#2563eb)
- Secundária: Verde sucesso (#10b981)
- Perigo: Vermelho erro (#ef4444)
- Neutro: Cinzas (#6b7280, #9ca3af, #d1d5db)

**Tipografia:**
- Fontes: System fonts (sans-serif)
- Hierarquia: h1-h6 com tamanhos consistentes
- Legibilidade: Contraste adequado, line-height

**Componentes:**
- Botões: Primário, Secundário, Perigo
- Formulários: Labels, inputs, validation states
- Cards: Para informações agrupadas
- Tabelas: Para listagens
- Modals: Para confirmações

## Restrições de Design

1. **Mobile-first:** Design responsivo desde o início
2. **Acessibilidade:** HTML semântico, ARIA labels, keyboard nav
3. **Performance:** CSS/JS otimizados, lazy loading quando apropriado
4. **Manutenibilidade:** CSS organizado, componentes reutilizáveis
5. **Consistência:** Design system aplicado uniformemente
6. **UX intuitiva:** Fluxos claros, feedback imediato
7. **Error handling:** Mensagens claras, recovery fácil

## Critérios de Aceitação

✅ Sistema de templates com layout base implementado  
✅ Página de listagem de pacientes com tabela e paginação  
✅ Formulário de criação de paciente com validação client-side  
✅ Página de detalhes do paciente com informações completas  
✅ Sistema de estilos responsivo e consistente  
✅ Interatividade com HTMX (form submit, delete confirmation)  
✅ Navegação intuitiva com menu e breadcrumbs  
✅ Mensagens de feedback (success, error, validation)  
✅ Acessibilidade básica implementada  
✅ Projeto compila e funcionalidades trabalham end-to-end  
✅ Backend integration completa com handlers existentes  

## Passos de Implementação

1. Analisar estrutura atual de templates e handlers
2. Criar sistema de layout base com partials
3. Implementar página de listagem de pacientes
4. Implementar formulário de criação
5. Implementar página de detalhes
6. Adicionar sistema de estilos (CSS)
7. Adicionar interatividade com HTMX
8. Implementar navegação e UX improvements
9. Testar fluxo completo end-to-end
10. Validar responsividade e acessibilidade

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md` (requirement original)
- `web/handlers/handler.go` (handlers existentes)
- `web/templates/` (estrutura atual de templates)
- `internal/application/services/patient_service.go` (service com métodos)
- `work/tasks/task_20260313_224928/implementation.md` (tarefa anterior do service)
