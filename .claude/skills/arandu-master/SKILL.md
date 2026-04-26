---
name: arandu-master
description: >
  Protocolo de comportamento do agente master no projeto Arandu. Define quando delegar
  exploração de codebase para agentes baratos (Explore) em vez de consumir tokens caros
  do master em leitura de arquivos. Use esta skill SEMPRE que o master iniciar uma sessão,
  selecionar a próxima feature, escrever um prompt de implementação ou avaliar entrega do
  coder. Regra de ouro: tokens do master = decisões + prompts. Exploração = Explore agent.
---

# Arandu — Protocolo do Agente Master

O master (Claude Sonnet) é o agente caro. Cada `Read`, `grep` ou exploração de codebase
feita diretamente consome tokens que deveriam ser reservados para raciocínio arquitetural,
escrita de prompts e julgamento de qualidade.

**Regra de ouro:** nunca use `Read` ou `Bash(grep)` quando a pergunta for "o que existe
no codebase". Isso é trabalho para o agente `Explore`.

---

## O que o master faz com seus tokens

| Atividade | Master | Explore |
|-----------|--------|---------|
| Decisão de qual feature tem maior valor | ✅ | ❌ |
| Modelagem DDD / raciocínio arquitetural | ✅ | ❌ |
| Escrita do `task.md` (prompt do coder) | ✅ | ❌ |
| Avaliação de qualidade da entrega | ✅ | ❌ |
| Raciocínio de privacidade / segurança | ✅ | ❌ |
| Mapear estado dos requisitos | ❌ | ✅ |
| Ler arquivos de referência para o prompt | ❌ | ✅ |
| Verificar se arquivo/rota/função existe | ❌ | ✅ |
| Comparar entrega do coder com spec | ❌ | ✅ |

**Exceção:** leituras cirúrgicas de ≤ 2 arquivos cujo caminho exato já é conhecido
podem ser feitas diretamente — o overhead de spawnar Explore não compensa.

---

## Protocolo 1 — Seleção de próxima feature

**Nunca leia os arquivos de requisito diretamente.** Use o Explore com este prompt:

```
Exploração de contexto — próxima feature do Arandu.

Leia todos os arquivos em docs/requirements/*.md e retorne uma tabela compacta:
| ID | Título (curto) | Status | Itens [ ] pendentes (conta) |

Para o status, grep por: "^[Ss]tatus:" ou "| \*\*Status\*\*".
Para itens pendentes, conta as linhas com "- [ ]" em cada arquivo.

Não analise — apenas mapeie. Máximo 300 palavras.
```

Recebido o mapa → master aplica julgamento de valor clínico e recomenda.

---

## Protocolo 2 — Escrita do prompt de implementação

Antes de escrever o `task.md`, delegar coleta de contexto técnico:

```
Coleta de contexto técnico para prompt de implementação — [FEATURE].

Preciso dos seguintes dados para escrever o prompt do coder. Leia os arquivos abaixo
e retorne APENAS os fatos relevantes (sem análise, sem recomendação):

1. [arquivo.go]: interface exata do repositório / assinatura do construtor
2. [arquivo_test.go]: padrão de mock usado (struct com campos *Func? ou interface?)
3. [handler.go]: como extrai tenant DB do context (trecho exato de código)
4. [types.go]: campos do ViewModel relevante

Máximo 400 palavras. Formato: bloco de código por item.
```

Master recebe o contexto compacto → invoca `implementation-prompt` → escreve task.md.

---

## Protocolo 3 — Verificação pontual de existência

Quando precisar saber "isso já existe?":

```
Verificação pontual no Arandu:
- Existe o arquivo [caminho]? Se sim, qual a assinatura de [função]?
- A rota [METHOD /path] está registrada em cmd/arandu/main.go?
- O campo [Campo] existe na struct [Struct] em [arquivo]?

Responda em ≤ 5 linhas por item. Sem contexto adicional.
```

---

## Protocolo 4 — Avaliação de entrega do coder

Antes de avaliar, delegar leitura dos arquivos criados:

```
Avaliação de conformidade — task [ID].

Leia os seguintes arquivos criados pelo coder:
- [arquivo1]
- [arquivo2]

Para cada função/comportamento listado abaixo, classifique:
✅ conforme | ⚠️ desvio (descreva em 1 linha) | ❌ ausente

Especificações:
[cole os itens do checklist da task]

Máximo 400 palavras. Apenas fatos — sem recomendação.
```

Master recebe o relatório → aplica julgamento final → decide se conclui ou solicita correção.

---

## Protocolo 5 — Início de sessão

Ao retomar trabalho após break, antes de qualquer análise:

```
Snapshot de estado — sessão Arandu.

1. Existe alguma task em work/tasks/ com status PRONTO_PARA_IMPLEMENTACAO ou IN_PROGRESS?
   Se sim, retorne o título e o checklist pendente.
2. Há arquivos modificados (git status) que não estão num commit?
   Se sim, liste-os.
3. O servidor está rodando? (verifique arandu.pid ou ps aux | grep arandu)

Máximo 150 palavras.
```

---

## Anti-padrões a evitar

- **Nunca** ler 5+ arquivos seguidos com `Read` para "entender o contexto" — use Explore
- **Nunca** fazer `grep` iterativo em múltiplos diretórios — formule uma única query para Explore
- **Nunca** invocar skills e depois repetir a exploração que a skill já define — confie no Explore
- **Nunca** escrever o task.md antes de ter o contexto técnico (Protocolo 2) — prompts incompletos custam mais caro na reescrita pelo coder

---

## Custo esperado por atividade

| Atividade | Sem skill (tokens master) | Com skill (tokens master) |
|-----------|--------------------------|--------------------------|
| Seleção de feature | ~3.000 (lê 10 req files) | ~400 (recebe tabela) |
| Escrita de prompt | ~2.000 (lê 6 arquivos) | ~600 (recebe contexto compacto) |
| Avaliação de entrega | ~2.500 (lê 4 novos arquivos) | ~500 (recebe relatório) |
| Início de sessão | ~1.500 (lê docs, git status) | ~200 (snapshot compacto) |
