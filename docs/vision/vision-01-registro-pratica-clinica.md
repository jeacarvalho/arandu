# VISION-01 — Registro da prática clínica

## Identificação

**ID:** VISION-01
**Nome:** Registro da prática clínica
**Status:** ativo

---

# Propósito

Permitir que profissionais de saúde mental **registrem de forma estruturada sua prática clínica**, preservando a riqueza das sessões terapêuticas e criando uma base consistente para reflexão e aprendizado ao longo do tempo.

O registro clínico é a **fundação operacional do Arandu**.

Sem registros estruturados de pacientes, sessões, observações e intervenções, não é possível construir:

* memória clínica longitudinal
* organização do conhecimento terapêutico
* identificação de padrões clínicos
* assistência reflexiva por IA

Essa vision estabelece, portanto, a **camada fundamental de dados clínicos do sistema**.

---

# Problema

A maioria dos profissionais de saúde mental registra suas sessões de forma fragmentada:

* cadernos pessoais
* documentos dispersos
* notas rápidas em aplicativos genéricos
* memória do terapeuta

Esse modelo gera diversos problemas:

```text
dificuldade de recuperar histórico clínico
perda de observações relevantes
baixa capacidade de identificar padrões entre casos
falta de organização da evolução terapêutica
```

Além disso, ferramentas tradicionais de prontuário costumam focar em **registro administrativo**, e não em **registro cognitivo da prática clínica**.

---

# Visão de solução

O Arandu deve fornecer um ambiente onde o terapeuta possa registrar sua prática clínica de forma:

```text
simples
estruturada
recuperável
reflexiva
```

Esse registro deve preservar elementos fundamentais da prática terapêutica:

```text
paciente
sessões
observações clínicas
intervenções terapêuticas
```

Esses elementos formarão a base para todas as funcionalidades mais avançadas do sistema.

---

# Elementos fundamentais da prática clínica

A prática clínica registrada no sistema será organizada em torno da seguinte estrutura conceitual:

```text
Paciente
  └── Sessões
        ├── Observações
        └── Intervenções
```

### Paciente

Representa a pessoa atendida pelo profissional.

É a **entidade raiz do registro clínico**.

---

### Sessão

Representa um encontro terapêutico entre profissional e paciente.

Cada sessão registra o contexto e os acontecimentos relevantes do atendimento.

---

### Observação

Representa percepções clínicas registradas pelo terapeuta durante ou após a sessão.

Exemplos:

```text
padrões emocionais
relatos significativos
comportamentos observados
hipóteses clínicas iniciais
```

---

### Intervenção

Representa ações terapêuticas realizadas durante a sessão.

Exemplos:

```text
perguntas exploratórias
exercícios terapêuticos
reformulações narrativas
experimentos comportamentais
```

---

# Valor para o profissional

Ao registrar a prática clínica dessa forma, o terapeuta passa a ter:

```text
histórico organizado de pacientes
registro consistente de sessões
memória clínica confiável
material estruturado para reflexão posterior
```

Isso melhora diretamente:

```text
continuidade do tratamento
qualidade da análise clínica
capacidade de aprendizado a partir da própria prática
```

---

# Relação com outras visions

A VISION-01 é **fundacional**.

Ela habilita diretamente as seguintes visions:

```text
VISION-02 Memória clínica longitudinal
VISION-03 Organização do conhecimento clínico
VISION-04 Descoberta de padrões clínicos
VISION-05 Assistência reflexiva por IA
```

Sem registros estruturados de prática clínica, nenhuma dessas visões pode existir.

---

# Capabilities derivadas

A partir desta vision surgem as seguintes capabilities:

```text
CAP-01-00 Gestão de pacientes
CAP-01-01 Registro de sessões
CAP-01-02 Observações clínicas
CAP-01-03 Intervenções terapêuticas
```

Essas capabilities definem **as funcionalidades operacionais do registro clínico**.

---

# Fora do escopo desta vision

Esta vision **não inclui**:

```text
análise automática de padrões
assistência de IA
comparação entre casos
agenda clínica
gestão financeira do consultório
```

Essas funcionalidades pertencem a outras visions do sistema.

---

# Resultado esperado

Quando esta vision estiver implementada, o Arandu permitirá que o profissional:

```text
registre pacientes
registre sessões
registre observações clínicas
registre intervenções terapêuticas
```

Esses registros formarão a **base de conhecimento clínica inicial do sistema**.

---

# Observação estratégica

O objetivo desta vision **não é apenas registrar dados**, mas **preservar a riqueza da prática clínica do terapeuta**.

O Arandu deve tratar esses registros não como formulários burocráticos, mas como **material vivo de reflexão clínica**.

