REQ-01-04-01 — Registro de Histórico Farmacológico e Sinais Vitais

Identificação

ID: REQ-01-04-01

Capability: CAP-01-04 — Contexto Biopsicossocial e Farmacológico

Vision: VISION-09 — Inteligência Clínica Ampliada

Status: draft

História do utilizador

Como psicólogo clínico, quero registrar e acompanhar a medicação atual e indicadores fisiológicos (sono, apetite, peso) do paciente, para compreender como fatores biológicos e farmacológicos influenciam o processo psíquico e o humor.

Contexto

Na clínica moderna (SOTA), o psicólogo não ignora o corpo. Medicamentos psiquiátricos alteram a cognição, e a qualidade do sono é um biomarcador crítico de recaída. Este requisito permite que o Arandu monitore esses dados de forma estruturada, preparando o terreno para que a IA identifique correlações (ex: "A ansiedade aumenta quando o paciente relata menos de 6h de sono").

Descrição funcional

O sistema deve prover seções específicas dentro do perfil do paciente para:

1. Histórico Farmacológico

Cadastro de Medicação: Nome do fármaco, dosagem, frequência e médico prescritor.

Status: Indicar se o uso é "Ativo", "Suspenso" ou "Finalizado".

Reconciliação: Durante a sessão, o terapeuta deve poder marcar rapidamente se houve mudança na medicação.

2. Sinais Vitais e Hábitos (Time-series)

Sono: Registro de horas ou qualidade percebida (1-10).

Apetite/Peso: Monitoramento de mudanças bruscas.

Atividade Física: Frequência semanal.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Visibilidade Progressiva: Esses dados não devem poluir a tela principal. Eles residem em um "Painel Biopsicossocial" lateral ou em uma aba dedicada.

Edição Rápida: Uso de micro-formulários HTMX que salvam sem recarregar a página.

Visualização: Medicamentos ativos devem ter um destaque sutil (cor verde Arandu pálida); suspensos devem aparecer riscados ou em cinza claro.

Lógica Técnica (SQL)

Novas tabelas para suportar os dados estruturados:

-- Medicamentos
CREATE TABLE IF NOT EXISTS patient_medications (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT,
    frequency TEXT,
    status TEXT DEFAULT 'active', -- active, suspended, finished
    started_at DATETIME,
    ended_at DATETIME,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Sinais Vitais (Série Temporal)
CREATE TABLE IF NOT EXISTS patient_vitals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATETIME NOT NULL,
    sleep_hours REAL,
    appetite_level INTEGER, -- 1 a 10
    weight REAL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);


Fluxo HTMX

O terapeuta abre o painel "Contexto Biológico".

POST: /patients/{id}/medications para adicionar novo remédio.

GET: /patients/{id}/medications para listar via fragmento.

PUT: /patients/{id}/medications/{med_id}/status para suspender/ativar rapidamente.

Critérios de Aceitação

CA-01: O sistema deve permitir listar todos os medicamentos ativos do paciente.

CA-02: Mudanças de status de medicação devem ser registradas com data automática.

CA-03: O registro de sinais vitais deve permitir a visualização de uma média recente.

CA-04: A interface deve respeitar a tipografia Inter para dados técnicos e Source Serif para observações farmacológicas.

CA-05: Os dados devem ser isolados por tenant_id (multi-tenancy).

Fora do escopo

Gráficos de linha complexos para o peso (será CAP-09-01).

Integração com APIs de farmácias.

Prescrição digital (o Arandu é apenas registro).