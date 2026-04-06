Implementação: REQ-03-02-01 — Classificar Intervenções Terapêuticas

📝 O que construir
Sistema de classificação de intervenções terapêuticas usando tags predefinidas organizadas em 6 tipos (Cognitiva, Comportamental, Emocional, Psicoeducação, Exploração Narrativa, Corporal). Similar ao REQ-03-01-01 mas adaptado para intervenções.

🧩 Componentes necessários
1. Tag Selector Modal (col-span-2)
2. Tag List Display (col-span-2)
3. Classification Summary Panel (col-span-1)
4. Intervention List with Tags (col-span-3)

📋 Regras específicas
- Tags organizadas em 6 tipos com cores específicas
- Cada tag pode ter intensidade 1-5 (opcional)
- HTMX para adição/remoção sem reload
- Persistência em tabelas intervention_tags e intervention_classifications
- Busca por tipo de tag
- Estatísticas de uso por tipo

⚠️ PADRÕES OBRIGATÓRIOS (NÃO IGNORAR)
[MANter todos os padrões do template]

🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

6. 🆕 Tailwind Config (Tokens a Criar/Usar)
// tailwind.config.js - Adicionar/Verificar:
theme: {
  extend: {
    colors: {
      'intervention-cognitive': '#7C3AED',
      'intervention-behavioral': '#1D9E75',
      'intervention-emotional': '#0F6E56',
      'intervention-psychoeducation': '#F59E0B',
      'intervention-narrative': '#3B82F6',
      'intervention-body': '#DC2626',
    },
    spacing: {
      'tag-selector': '400px',
      'classification-panel': '280px',
    },
    fontFamily: {
      'clinical': ['Source Serif 4', 'serif'],
    },
  },
}

7. 🆕 Estrutura de Arquivos
Criar:
- internal/domain/intervention/classification.go
- internal/web/handlers/classification_handler.go
- internal/infrastructure/repository/sqlite/intervention_tag_repository.go
- web/components/intervention/tag_selector.templ
- web/components/intervention/tag_badge.templ
- web/components/intervention/classification_panel.templ
- internal/infrastructure/repository/sqlite/migrations/0012_add_intervention_tags.up.sql
- internal/infrastructure/repository/sqlite/migrations/0012_add_intervention_tags.down.sql

8. 🆕 Seed Data
Criar script de seed para popular 30+ tags predefinidas distribuídas nos 6 tipos

9. 🆕 Testes
- Testes de unitários para service de classificação
- Testes de integração para repository
- Testes E2E para fluxo HTMX de adição/remoção
```