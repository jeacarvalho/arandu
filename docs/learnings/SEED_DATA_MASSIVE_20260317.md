# Aprendizado: Geração de Massa de Dados Clínica Massiva

**Data:** 2026-03-17  
**Tarefa:** task_20260317_140558  
**Autor:** Agente Arandu

## Resumo

Criado script Python gerador e SQL com massa de dados clínica realista para testes profundos do sistema Arandu.

## Dados Gerados

- **500 pacientes** com nomes brasileiros realistas
- **62.849 sessões** distribuídas ao longo de ~2 anos (100-150 por paciente)
- **188.532 observações clínicas** (2-4 por sessão)
- **125.583 intervenções terapêuticas** (1-3 por sessão)

## Arquivos Criados

1. **Gerador Python:** `scripts/generate_massive_data.py`
2. **Script SQL:** `internal/infrastructure/repository/sqlite/seeds/seed_massive_clinical_data.sql` (61.48 MB)

## Características dos Dados

- **15 contextos clínicos** diferentes (TCC, Psicanálise, DBT, etc.)
- **15 archetypes de pacientes**: ansiedade, depressão, pânico, fobia social, TOC, TEPT, luto, crise de meia-idade, bipolar, borderline, dependência emocional, burnout, adoção, separação, cuidador
- **Templates de observações**: 85+ frases clínicas realistas
- **Templates de intervenções**: 75+ técnicas terapêuticas

## Uso

```bash
# Gerar novos dados
python3 scripts/generate_massive_data.py

# Executar no banco
sqlite3 arandu.db < internal/infrastructure/repository/sqlite/seeds/seed_massive_clinical_data.sql
```

## Objetivo do Teste

Esta massa de dados foi criada para:
1. **Testar performance** do SQLite com base grande
2. **Testar funcionalidades de IA** para gerar insights clínicos
3. **Validação de UI/UX** com dados realistas
4. **Testes de busca e filtragem** em grandes volumes
