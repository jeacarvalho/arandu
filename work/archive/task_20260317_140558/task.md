Prompt: Geração de Massa de Dados Clínica (Narrativa SOTA)

Objetivo: Gerar um script SQL (seed_data.sql) para popular o banco de dados SQLite do Arandu com casos clínicos realistas, complexos e longitudinais.

🎭 Contexto da IA

Você é um psicólogo clínico sênior com 20 anos de experiência e um especialista em bancos de dados SQL. Sua tarefa é criar a história de 5 pacientes distintos, cada um com uma trajetória terapêutica de pelo menos 2 anos (mínimo de 15 sessões por paciente).

Archetipos de Pacientes Sugeridos:

O Executivo em Burnout: Crises de pânico, alta produtividade mascarando depressão.

O Idoso em Luto: Perda de cônjuge, solidão, busca por novo sentido vital.

A Adolescente com Ansiedade Social: Dificuldade escolar, uso excessivo de telas, questões de autoimagem.

O Casal em Crise (focar em um dos membros): Conflitos de comunicação, traição, reconstrução de confiança.

O Artista com Bloqueio Criativo: Transtorno bipolar leve (ciclotimia), flutuações de humor e produtividade.

🛠️ Requisitos Técnicos do Script SQL

O script deve seguir rigorosamente o schema do Arandu:

UUIDs: Use UUIDs válidos para todos os campos id.

Foreign Keys: Respeite a hierarquia: Patient -> Session -> (Observation & Intervention).

Datas: Distribua as sessões ao longo dos últimos 24 meses. Use timestamps SQLite válidos (YYYY-MM-DD HH:MM:SS).

Coerência: Uma Intervention deve fazer sentido clínico em relação às Observations daquela mesma sessão ou sessões anteriores.

Estrutura das Tabelas:

patients: (id, name, notes, created_at, updated_at)

sessions: (id, patient_id, date, summary, created_at, updated_at)

observations: (id, session_id, content, created_at, updated_at)

interventions: (id, session_id, content, created_at, updated_at)

✍️ Qualidade do Conteúdo (Silent UI Ready)

Linguagem: O conteúdo clínico deve ser escrito em Português (Brasil), com terminologia técnica adequada (Psicanálise, TCC ou Fenomenologia - escolha uma abordagem por paciente).

Narrativa: As observações não devem ser apenas fatos, mas percepções: "Paciente evitou contato visual ao mencionar a figura paterna, indicando uma resistência já observada na sessão 04".

Intervenções: Devem ser ações: "Realizada técnica de cadeira vazia para processamento do luto não resolvido".

Evolução: O "Summary" das sessões deve mostrar a evolução (ou retrocesso) do caso ao longo do tempo.

📤 Saída Esperada

Um arquivo SQL único que possa ser executado diretamente no SQLite:

DELETE FROM em todas as tabelas (para resetar o ambiente de teste).

INSERT INTO patients.

INSERT INTO sessions.

INSERT INTO observations.

INSERT INTO interventions.

Instrução Final: Gere dados que, ao serem lidos na fonte Source Serif 4, pareçam um prontuário médico real e respeitado.