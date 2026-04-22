# Referência Regulatória

## CFP — Resolução 11/2018 (Telepsicologia)

Pontos críticos para o sistema:

- Psicólogo deve manter **prontuário eletrônico** de todos os atendimentos online
- Atendimentos online exigem **plataforma segura** (criptografia, controle de acesso)
- Psicólogo é responsável pela **guarda e sigilo** dos dados, mesmo em plataforma terceira
- Menor de idade: telepsicologia requer autorização expressa do responsável legal
- Obrigatório: identificação do psicólogo (nome + CRP) em toda comunicação com paciente
- Vedado: gravação de sessões sem consentimento explícito e documentado

**Implicações no sistema:**
- Todo acesso ao prontuário deve ser logado (audit trail)
- TCLE digital deve ser assinado antes da primeira sessão
- Plataforma não pode acessar conteúdo das sessões sem autorização do psicólogo

---

## LGPD aplicada ao contexto clínico

### Dados de saúde = Dados Sensíveis (Art. 11)
Dados de saúde têm proteção reforçada. Tratamento só permitido com:
- Consentimento específico e destacado do titular
- Finalidade específica (não pode usar dados de saúde para fins não declarados)

### Base legal para tratamento
- **Consentimento** (Art. 7, I): TCLE assinado pelo paciente/responsável
- **Tutela da saúde** (Art. 11, II, f): único profissional de saúde responsável

### Direitos do titular que o sistema deve suportar
| Direito | Art. | Implementação |
|---------|------|---------------|
| Acesso | 18, II | Paciente pode solicitar cópia do prontuário |
| Correção | 18, III | Psicólogo corrige dados cadastrais (não clínicos) |
| Eliminação | 18, VI | Complexo: dados clínicos têm retenção obrigatória |
| Portabilidade | 18, V | Exportação do prontuário em formato legível |
| Revogação do consentimento | 18, IX | Encerra novos atendimentos; mantém histórico |

### DPO e Relatório de Impacto
- Plataforma deve nomear DPO (Encarregado de Dados)
- Elaborar RIPD (Relatório de Impacto à Proteção de Dados) para dados sensíveis de saúde

---

## Prontuário Eletrônico — CFP

- Resolução CFP 001/2009 e atualizações: prontuário é documento legal
- **Retenção mínima**: 5 anos após último atendimento (adultos)
- **Retenção mínima**: 5 anos após o paciente completar 18 anos (menores)
- Prontuário não pode ser deletado — apenas arquivado com acesso restrito
- Alterações devem manter histórico (nunca sobrescrever, sempre versionar)
- Psicólogo responde pelo prontuário mesmo após encerrar uso da plataforma

### Conteúdo mínimo obrigatório do prontuário
1. Identificação do paciente
2. Data de início do atendimento
3. Registro de cada sessão (data, duração, observações)
4. TCLE assinado
5. Identificação do psicólogo responsável (nome + CRP)

---

## Sigilo Profissional — Código de Ética CFP

- Art. 9: psicólogo é obrigado a manter sigilo sobre informações do paciente
- **Exceções ao sigilo** (situações em que o sistema deve alertar o psicólogo):
  - Risco de vida para o paciente ou terceiros
  - Determinação judicial
  - Proteção de menor em situação de risco
- A plataforma **não tem acesso** ao conteúdo clínico — só o psicólogo
- Qualquer análise de IA deve ser **opt-in** explícito do psicólogo, sessão a sessão ou por configuração
