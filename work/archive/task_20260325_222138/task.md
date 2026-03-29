# TASK 20260325_222138

**Requirement:** REQ-01-00-02, REQ-02-01-01
**Title:** Correção de Tipografia Clínica e Validações SLP
**Status:** PRONTO_PARA_IMPLEMENTACAO

---

## 🎯 Objetivo

Corrigir validações SLP falhando no E2E Audit:
1. Adicionar `.font-clinical` a todo conteúdo clínico
2. Corrigir rota `/patients/{id}/history` retornando 302
3. Opcional: Ajustar session_update para retornar 200 com HTMX

---

## 📊 Problemas Identificados (E2E Audit)

### Problema 1: Clinical Typography Missing
```
[ERROR] dashboard - Clinical typography (.font-clinical) not found
[ERROR] patients_list - Clinical typography (.font-clinical) not found
[ERROR] patients_detail - Clinical typography (.font-clinical) not found
```

**O que o script valida:**
```bash
grep -q "font-clinical\|font-serif\|Source Serif" "$file"
```

**Solução:** Adicionar classe `.font-clinical` a:
- Nomes de pacientes (quando exibidos como conteúdo clínico)
- Observações clínicas
- Resumos de sessões
- Notas de pacientes
- Qualquer texto narrativo clínico

**CSS necessário (web/static/css/style.css):**
```css
.font-clinical {
    font-family: 'Source Serif 4', serif;
    font-size: 1.125rem;
    line-height: 1.75;
    color: #1F2937;
}
```

---

### Problema 2: Patients History 302
```
[ERROR] patients_history returned 302
```

**Debug:**
```bash
# Verificar se rota existe
grep -r "history" internal/web/handlers/ --include="*.go"

# Verificar se está registrada no main.go
grep -r "/patients/.*/history" cmd/arandu/main.go
```

**Solução:** Garantir que rota existe e está acessível com sessão válida

---

## ✅ Critérios de Aceitação

- [ ] E2E Audit passa sem erros de tipografia clínica
- [ ] `/patients/{id}/history` retorna 200 (não 302)
- [ ] Screenshots continuam sendo gerados
- [ ] `./scripts/arandu_e2e_audit.sh` retorna exit code 0

---

## 🧪 Validação Pós-Implementação

```bash
# 1. Rodar E2E Audit
./scripts/arandu_e2e_audit.sh

# 2. Verificar validações de tipografia
grep -c "font-clinical" tmp/audit_logs/route_*.html

# 3. Verificar screenshots
ls -la tmp/audit_screenshots/

# 4. Verificar exit code
echo $?  # Deve ser 0
```

---

## 📚 Referências

- `docs/design-system.md` — Tipografia Source Serif 4
- `docs/architecture/standardized_layout_protocol.md` — SLP
- `scripts/arandu_e2e_audit.sh` — Validações SLP
```

---
