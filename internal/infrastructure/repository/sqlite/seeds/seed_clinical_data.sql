-- Seed Data: Massa de Dados Clínica para Testes
-- Gerado em: 2026-03-17
-- Contém: 5 pacientes, 78 sessões, observações e intervenções realistas
-- Período: 24 meses de histórico terapêutico

-- Limpar dados existentes
DELETE FROM interventions;
DELETE FROM observations;
DELETE FROM sessions;
DELETE FROM patients;
DELETE FROM insights;

-- ============================================
-- PACIENTE 1: Executivo em Burnout
-- Abordagem: TCC
-- Período: 24 meses
-- ============================================
INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES
('p1-exec-burnout', 'Ricardo Almeida', 'Executivo de 45 anos, diretor comercial. Início do tratamento após primeira crise de pânico no trânsito. Alta funcionalidade mascarando quadro depressivo.', '2024-01-15 09:00:00', '2024-01-15 09:00:00');

-- Sessões do Paciente 1
INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES
('s1-001', 'p1-exec-burnout', '2024-01-15 14:00:00', 'Primeira sessão. Paciente relata crise de pânico ocorrida há 3 semanas no trânsito. Sintomas: taquicardia, sudorese excessiva, sensação de morte iminente. Atendeu pronto-socorro, exames cardiológicos normais. Nega história de crises anteriores. Histórico de trabalho excessivo (60+ horas/semana).', '2024-01-15 14:00:00', '2024-01-15 14:00:00'),
('s1-002', 'p1-exec-burnout', '2024-01-22 14:00:00', 'Segunda sessão. Relata ansiedade generalizada, dificuldade de concentração e insônia. Apresenta-se hipervigilante, monitorando sinais corporais constantemente. Evita autopistas desde a crise. Distorções cognitivas identificadas: catastrofização, leitura mental.', '2024-01-22 14:00:00', '2024-01-22 14:00:00'),
('s1-003', 'p1-exec-burnout', '2024-01-29 14:00:00', 'Terceira sessão. Início da psicoeducação sobre transtorno de pânico. Paciente demonstra resistência em reconhecer componente psicológico, busca causas orgânicas. Histórico familiar: pai com infarto aos 50 (sobreviveu). Medo de "morrer como o pai".', '2024-01-29 14:00:00', '2024-01-29 14:00:00'),
('s1-004', 'p1-exec-burnout', '2024-02-05 14:00:00', 'Quarta sessão. Apresenta diário de pensamentos. Identificada tendência à catastrofização em situações profissionais. Discussão sobre pressão no trabalho e medo de falha. Lágrimas durante sessão - primeira vulnerabilidade demonstrada.', '2024-02-05 14:00:00', '2024-02-05 14:00:00'),
('s1-005', 'p1-exec-burnout', '2024-02-12 14:00:00', 'Quinta sessão. Início do registro de pensamentos automáticos. Paciente relata evitar eventos sociais por medo de crises. Isolamento progressivo notado. Discussão sobre associação entre produtividade e valor pessoal.', '2024-02-12 14:00:00', '2024-02-12 14:00:00'),
('s1-006', 'p1-exec-burnout', '2024-02-19 14:00:00', 'Sexta sessão. Primeira tentativa de exposição gradual: conduzir curta distância. Relata ansiedade moderada (6/10) mas concluiu trajeto. Reforço positivo. Discussão sobre pensamentos "O que as pessoas vão pensar".', '2024-02-19 14:00:00', '2024-02-19 14:00:00'),
('s1-007', 'p1-exec-burnout', '2024-02-26 14:00:00', 'Sétima sessão. Crise de pânico relatada ontem, após reunião estressante. Análise da situação: gatilho foi pensamento "Vou falhar". Discussão sobre estresse acumulado e necessidade de estabelecer limites no trabalho.', '2024-02-26 14:00:00', '2024-02-26 14:00:00'),
('s1-008', 'p1-exec-burnout', '2024-03-05 14:00:00', 'Oitava sessão. Paciente apresenta-se mais reflexivo. Discussão sobre infância: pais exigentes, mensagem implícita de que amor era condicional ao sucesso. Insight parcial: "Sempre precisei provar meu valor".', '2024-03-05 14:00:00', '2024-03-05 14:00:00'),
('s1-009', 'p1-exec-burnout', '2024-03-12 14:00:00', 'Nona sessão. Exposição a autopista com terapeuta acompanhando. Ansiedade inicial 7/10, decrescente ao longo do trajeto. Discussão sobre catastrofização de sintomas corporais normais.', '2024-03-12 14:00:00', '2024-03-12 14:00:00'),
('s1-010', 'p1-exec-burnout', '2024-03-19 14:00:00', 'Décima sessão. Paciente relata delegar mais tarefas no trabalho. Redução de ansiedade geral (escala 4/10). Discute medo de ser visto como "fraco" por diminuir ritmo. Cognição central identificada.', '2024-03-19 14:00:00', '2024-03-19 14:00:00'),
('s1-011', 'p1-exec-burnout', '2024-04-02 14:00:00', 'Décima primeira sessão. Conflito com esposa: ela celebrou mudanças, ele interpretou como pressão. Discussão sobre mudança de comportamento versus mudança de identidade. Reestruturação cognitiva da crença central.', '2024-04-02 14:00:00', '2024-04-02 14:00:00'),
('s1-012', 'p1-exec-burnout', '2024-04-16 14:00:00', 'Décima segunda sessão. Viajou a negócios de avião (exposição planejada). Relata ansiedade gerenciável, usou técnicas de respiração. Primeira viagem sem medo debilitante em meses. Humor melhorado.', '2024-04-16 14:00:00', '2024-04-16 14:00:00'),
('s1-013', 'p1-exec-burnout', '2024-04-30 14:00:00', 'Décima terceira sessão. Reflexão sobre 4 meses de terapia. Reconhece padrão de busca excessiva de aprovação externa. Discussão sobre autocompaixão e autoestima condicional. Início de diário de gratidão.', '2024-04-30 14:00:00', '2024-04-30 14:00:00'),
('s1-014', 'p1-exec-burnout', '2024-05-14 14:00:00', 'Décima quarta sessão. Recaída: medo de perder o emprego após reestruturação na empresa. Ansiedade elevada novamente. Retorno às técnicas TCC. Discussão sobre estabilidade versus mudança: recaídas fazem parte.', '2024-05-14 14:00:00', '2024-05-14 14:00:00'),
('s1-015', 'p1-exec-burnout', '2024-05-28 14:00:00', 'Décima quinta sessão. Paciente processa notícia de que mantém emprego. Reflete sobre quanta energia desperdiçou em antecipar catástrofe. Discute plano de manutenção. Redução para sessões quinzenais.', '2024-05-28 14:00:00', '2024-05-28 14:00:00'),
('s1-016', 'p1-exec-burnout', '2024-06-11 14:00:00', 'Décima sexta sessão. Continuação quinzenal. Paciente mantém práticas: limites no trabalho, exercícios, técnicas de respiração. Relata crise suave há 2 semanas, gerenciada com sucesso. Evolução mantida.', '2024-06-11 14:00:00', '2024-06-11 14:00:00'),
('s1-017', 'p1-exec-burnout', '2024-07-09 14:00:00', 'Décima sétima sessão. Intervalo de férias. Paciente viajou com família, primeira vez em anos que relaxou verdadeiramente. Reflexão profunda sobre valores de vida. Possível início de transição profissional.', '2024-07-09 14:00:00', '2024-07-09 14:00:00'),
('s1-018', 'p1-exec-burnout', '2024-08-06 14:00:00', 'Décima oitava sessão. Tomou decisão de pedir diminuição de cargo (redução salarial) para qualidade de vida. Ansiedade anticipatória mas firmeza na escolha. Discussão sobre identidade sem status profissional.', '2024-08-06 14:00:00', '2024-08-06 14:00:00'),
('s1-019', 'p1-exec-burnout', '2024-09-03 14:00:00', 'Décima nona sessão. Assumiu novo cargo menos estressante. Adaptação em curso. Relata que, pela primeira vez, tem energia para hobbies. Discussão sobre definição de sucesso pessoal versus social.', '2024-09-03 14:00:00', '2024-09-03 14:00:00'),
('s1-020', 'p1-exec-burnout', '2024-10-01 14:00:00', 'Vigésima sessão. Avaliação de 9 meses. Sem crises de pânico há 4 meses. Ansiedade gerenciável. Relacionamento familiar melhorado. Paciente solicita espaçamento para sessões mensais.', '2024-10-01 14:00:00', '2024-10-01 14:00:00'),
('s1-021', 'p1-exec-burnout', '2024-11-05 14:00:00', 'Vigésima primeira sessão. Mensal. Mantém estabilidade. Relata conflito leve com colega, gerenciado assertivamente sem ansiedade excessiva. Demonstra ferramentas internalizadas.', '2024-11-05 14:00:00', '2024-11-05 14:00:00'),
('s1-022', 'p1-exec-burnout', '2024-12-03 14:00:00', 'Vigésima segunda sessão. Mensal. Final de ano. Paciente faz balanço positivo do ano. Reconhece que transformação exigiu coragem. Planeja alta para próximo trimestre.', '2024-12-03 14:00:00', '2024-12-03 14:00:00'),
('s1-023', 'p1-exec-burnout', '2025-01-14 14:00:00', 'Vigésima terceira sessão. Retorno após férias. Ano novo, novos desafios. Paciente mantém equilíbrio. Discussão sobre sustentabilidade da mudança.', '2025-01-14 14:00:00', '2025-01-14 14:00:00'),
('s1-024', 'p1-exec-burnout', '2025-02-11 14:00:00', 'Vigésima quarta sessão. Sessão de encerramento do primeiro ciclo. Paciente emocionado, agradece processo. Revisão de aprendizados. Plano de alta com possibilidade de retorno se necessário.', '2025-02-11 14:00:00', '2025-02-11 14:00:00');

-- Observações do Paciente 1
INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES
('obs1-001', 's1-001', 'Paciente apresenta postura rígida, evita contato visual prolongado. Fala rápida, pressa evidente. Linguagem corporal tensa: ombros elevados, mandíbula contraída. Sinais de hipervigilância: observa constantemente saída da sala.', '2024-01-15 14:50:00', '2024-01-15 14:50:00'),
('obs1-002', 's1-001', 'Dificuldade em nomear emoções: descreve sintomas físicos mas não afetivos. Quando perguntado sobre medo, responde "não sou de me abalar". Defesas intelectualizadoras presentes.', '2024-01-15 14:50:00', '2024-01-15 14:50:00'),
('obs1-003', 's1-002', 'Hoje mais falante, porém superficial. Lista sintomas como itens de agenda. Resistência em explorar significado afetivo. Quando aborda medo, minimiza: "Todo mundo passa por isso".', '2024-01-22 14:50:00', '2024-01-22 14:50:00'),
('obs1-004', 's1-002', 'Evita palavras como "ansiedade" ou "medo", substitui por "preocupação" ou "estresse". Observada primeira lágrima quando menciona noites em claro. Rápida recuperação e mudança de assunto.', '2024-01-22 14:50:00', '2024-01-22 14:50:00'),
('obs1-005', 's1-004', 'Momento significativo: lágrimas ao dizer "não sei quem sou sem o trabalho". Primeira vulnerabilidade genuína. Contato visual melhorado. Mãos relaxaram sobre joelhos.', '2024-02-05 14:50:00', '2024-02-05 14:50:00'),
('obs1-006', 's1-008', 'Insight parcial emergente: associa crise de pânico ao pai. "Meu pai infartou porque trabalhava demais, e eu estou fazendo o mesmo". Consciência do padrão repetitivo.', '2024-03-05 14:50:00', '2024-03-05 14:50:00'),
('obs1-007', 's1-012', 'Retorno de viagem: animação diferente do habitual. Contato visual espontâneo. Sorrisos genuínos. Primeira vez que traz tema positivo sem solicitação.', '2024-04-16 14:50:00', '2024-04-16 14:50:00'),
('obs1-008', 's1-018', 'Conflito interno evidente: satisfação com decisão de vida versus medo de julgamento. Oscila entre orgulho e vergonha. "Meus colegas vão achar que desisti".', '2024-08-06 14:50:00', '2024-08-06 14:50:00'),
('obs1-009', 's1-024', 'Última sessão: postura aberta, contato visual frequente e genuíno. Expressões faciais suaves. Despedida calorosa. Disponibilidade para retorno se necessário.', '2025-02-11 14:50:00', '2025-02-11 14:50:00');

-- Intervenções do Paciente 1
INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES
('int1-001', 's1-001', 'Psicoeducação sobre transtorno de pánico: explicação do ciclo de feedback entre interpretação catastrófica de sensações corporais e ansiedade crescente. Entrega de folder informativo.', '2024-01-15 14:50:00', '2024-01-15 14:50:00'),
('int1-002', 's1-002', 'Exercício de respiração diafragmática 4-7-8. Prática em sessão. Prescrição: 2x ao dia e durante episódios de ansiedade. Orientação sobre hyperventilação.', '2024-01-22 14:50:00', '2024-01-22 14:50:00'),
('int1-003', 's1-003', 'Técnica de registro de pensamentos automáticos. Modelo aplicado: situação -> pensamento -> emoção -> comportamento. Exemplo prático com episódio de trânsito.', '2024-01-29 14:50:00', '2024-01-29 14:50:00'),
('int1-004', 's1-006', 'Exposição gradual in vivo: plano de hierarquia de ansiedade. Primeiro item: dirigir 2 km em ruas secundárias. Acompanhamento terapêutico.', '2024-02-19 14:50:00', '2024-02-19 14:50:00'),
('int1-005', 's1-007', 'Reestruturação cognitiva da catastrofização. Questionamento socrático: "Qual evidência real de que vai falhar?" "O que aconteceria de pior?" "Como lidaria?"', '2024-02-26 14:50:00', '2024-02-26 14:50:00'),
('int1-006', 's1-009', 'Exposição interoceptiva: exercício de hiperventilação controlada em sessão para desfazer associação entre sintomas físicos e perigo. Habituação demonstrada.', '2024-03-12 14:50:00', '2024-03-12 14:50:00'),
('int1-007', 's1-011', 'Técnica de assertividade: ensaiada em sessão conversa com esposa sobre mudanças. Role-play. Feedback sobre comunicação não-violenta.', '2024-04-02 14:50:00', '2024-04-02 14:50:00'),
('int1-008', 's1-017', 'Exercício de valores: identificação de 3 valores fundamentais (família, saúde, autenticidade) e avaliação de alinhamento entre valores e comportamentos atuais.', '2024-07-09 14:50:00', '2024-07-09 14:50:00'),
('int1-009', 's1-024', 'Plano de manutenção: revisão de estratégias TCC aprendidas. Agenda de sessões de reforço trimestrais. Lista de sinais de alerta para retorno.', '2025-02-11 14:50:00', '2025-02-11 14:50:00');

-- ============================================
-- PACIENTE 2: Idosa em Luto
-- Abordagem: Psicanalítica
-- Período: 24 meses
-- ============================================
INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES
('p2-luto-idoso', 'Dona Helena Costa', 'Viúva de 72 anos, casada 45 anos. Início de terapia 8 meses após morte do marido. Queixas: solidão intensa, falta de sentido para vida, dificuldade em cuidar da casa sozinha.', '2024-02-01 10:00:00', '2024-02-01 10:00:00');

-- Sessões do Paciente 2 (primeiras 15 de 20)
INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES
('s2-001', 'p2-luto-idoso', '2024-02-01 10:00:00', 'Primeira sessão. Paciente traz foto do marido. Chora abundantemente. Relata que "não existe mais motivo para levantar da cama". Identificação simbiótica com falecido evidente. Nega estar deprimida, diz apenas "sentir falta".', '2024-02-01 10:00:00', '2024-02-01 10:00:00'),
('s2-002', 'p2-luto-idoso', '2024-02-08 10:00:00', 'Segunda sessão. Discute rotina diária: acorda às 10h, fica na cama até 14h. Refeições irregulares. Filhos preocupados mas distantes (mora em outro estado). Sentimento de abandono.', '2024-02-08 10:00:00', '2024-02-08 10:00:00'),
('s2-003', 'p2-luto-idoso', '2024-02-15 10:00:00', 'Terceira sessão. Associações livres sobre infância. Pai ausente, mãe dominadora. Casamento como "única escolha que fiz sozinha". Pergunta: "Quem sou eu sem ele?"', '2024-02-15 10:00:00', '2024-02-15 10:00:00'),
('s2-004', 'p2-luto-idoso', '2024-02-22 10:00:00', 'Quarta sessão. Traço depressivo mais evidente. Disfórica, psicomotora retardada. Refere dores inexplicáveis. Nega ideação suicida mas "não vê problema se não acordasse".', '2024-02-22 10:00:00', '2024-02-22 10:00:00'),
('s2-005', 'p2-luto-idoso', '2024-02-29 10:00:00', 'Quinta sessão. Primeira associação não relacionada ao marido: fala do jardim. Interesse despertado. Sugestão terapêutica de voltar a cuidar das plantas.', '2024-02-29 10:00:00', '2024-02-29 10:00:00'),
('s2-006', 'p2-luto-idoso', '2024-03-07 10:00:00', 'Sexta sessão. Retorna animada: reativou jardim, encontrou vizinha. Primeiro momento de afeto não-triste. Discussão sobre culpa em sentir prazer.', '2024-03-07 10:00:00', '2024-03-07 10:00:00'),
('s2-007', 'p2-luto-idoso', '2024-03-14 10:00:00', 'Sétima sessão. Oscilação: volte a falar do marido com intensidade. Culpa por "esquecer" dele. Interpretação sobre luto vs melancolia. Primeira compreensão intelectual.', '2024-03-14 10:00:00', '2024-03-14 10:00:00'),
('s2-008', 'p2-luto-idoso', '2024-03-21 10:00:00', 'Oitava sessão. Aniversário de casamento (seriam 46 anos). Dias difíceis, mas compareceu. Velório simbólico: trouxe flores do jardim para fotos. Ritual de despedida.', '2024-03-21 10:00:00', '2024-03-21 10:00:00'),
('s2-009', 'p2-luto-idoso', '2024-04-04 10:00:00', 'Nona sessão. Duas semanas sem sessão (feriado). Relata sonhos com marido. Associação: sonhos como forma de continuidade do vínculo. Angústia de separação diminuída.', '2024-04-04 10:00:00', '2024-04-04 10:00:00'),
('s2-010', 'p2-luto-idoso', '2024-04-11 10:00:00', 'Décima sessão. Revela que nunca dormiu sozinha antes: sempre com família, depois casada. Medo da escuridão infantilizado. Associação ao silêncio da casa vazia.', '2024-04-11 10:00:00', '2024-04-11 10:00:00'),
('s2-011', 'p2-luto-idoso', '2024-04-18 10:00:00', 'Décima primeira sessão. Comprou abajur de cabeceira. Primeira compra "só dela" em décadas. Discussão sobre autonomia e identidade feminina reprimida.', '2024-04-18 10:00:00', '2024-04-18 10:00:00'),
('s2-012', 'p2-luto-idoso', '2024-04-25 10:00:00', 'Décima segunda sessão. Conflito com filha: filha quer que mude de cidade, ela recusa. Afirmação de limites inédita. "Minha vida é aqui, com minhas memórias".', '2024-04-25 10:00:00', '2024-04-25 10:00:00'),
('s2-013', 'p2-luto-idoso', '2024-05-02 10:00:00', 'Décima terceira sessão. Revela talento oculto: poesia. Escreveu durante luto. Compartilha poemas. Emoção positiva, vergonha misturada com orgulho.', '2024-05-02 10:00:00', '2024-05-02 10:00:00'),
('s2-014', 'p2-luto-idoso', '2024-05-09 10:00:00', 'Décima quarta sessão. Participação em grupo de luto da igreja. Primeiro contato social estruturado. Relata tanto conforto quanto ressentimento de outras viúvas mais jovens.', '2024-05-09 10:00:00', '2024-05-09 10:00:00'),
('s2-015', 'p2-luto-idoso', '2024-05-16 10:00:00', 'Décima quinta sessão. Avaliação trimestral: já acorda cedo, alimentação regular, atividades sociais. Ainda triste, mas "tristeza normal, não aquela paralisia". Melhora significativa.', '2024-05-16 10:00:00', '2024-05-16 10:00:00'),
('s2-016', 'p2-luto-idoso', '2024-06-13 10:00:00', 'Décima sexta sessão. Reencontro com amiga de infância por acaso. Revela vida social anterior esquecida. "Eu tinha amigas, tinha vida... fiquei só dele". Insight sobre fusão.', '2024-06-13 10:00:00', '2024-06-13 10:00:00'),
('s2-017', 'p2-luto-idoso', '2024-07-11 10:00:00', 'Décima sétima sessão. Filha visita, elogio mútuo. Relacionamento reparado. Discussão sobre luto completo: preservar memória sem anular presente.', '2024-07-11 10:00:00', '2024-07-11 10:00:00'),
('s2-018', 'p2-luto-idoso', '2024-08-08 10:00:00', 'Décima oitava sessão. Início de projeto: caderno de receitas da família para netos. Primeiro projeto futuro. Integração do passado ao futuro.', '2024-08-08 10:00:00', '2024-08-08 10:00:00'),
('s2-019', 'p2-luto-idoso', '2024-09-12 10:00:00', 'Décima nona sessão. Espaçamento para quinzenal. Paciente confortável com autonomia. Sessões mais leves, risonhas. Matização do luto.', '2024-09-12 10:00:00', '2024-09-12 10:00:00'),
('s2-020', 'p2-luto-idoso', '2024-10-10 10:00:00', 'Vigésima sessão. Sessão de encerramento do ciclo inicial. Gratidão expressa. Reconhecimento do processo. Plano de retorno se necessário.', '2024-10-10 10:00:00', '2024-10-10 10:00:00');

-- Observações do Paciente 2
INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES
('obs2-001', 's2-001', 'Luto patológico suspeito: passados 8 meses, paciente ainda sem aceitação. Choros prolongados, catatonia psicomotora. Dificuldade em falar em tempo presente.', '2024-02-01 10:50:00', '2024-02-01 10:50:00'),
('obs2-002', 's2-003', 'Associação espontânea sobre infância reprimida: "Minha mãe dizia que casamento era obrigação". Primeira conexão entre casamento e falta de autonomia.', '2024-02-15 10:50:00', '2024-02-15 10:50:00'),
('obs2-003', 's2-006', 'Mudança qualitativa no afeto: primeiro sorriso genuíno ao falar do jardim. Emoções se ampliando além da tristeza. Sinais de elaboração.', '2024-03-07 10:50:00', '2024-03-07 10:50:00'),
('obs2-004', 's2-010', 'Regressão evidente: medo do escuro, dificuldade em dormir sozinha, busca por figura parental. Luto reativando questões de separação infantis não resolvidas.', '2024-04-11 10:50:00', '2024-04-11 10:50:00'),
('obs2-005', 's2-012', 'Afirmação de limites: resistência em ceder aos desejos da filha. Primeira vez que colocou suas necessidades em primeiro plano. Autonomia emergente.', '2024-04-25 10:50:00', '2024-04-25 10:50:00'),
('obs2-006', 's2-016', 'Insight sobre fusão marital: "Eu me cancelei para ser só esposa dele". Reconhecimento do custo da identidade perdida no casamento.', '2024-06-13 10:50:00', '2024-06-13 10:50:00'),
('obs2-007', 's2-020', 'Transformação observada: postura mais ereta, voz mais firme, contato visual sustentado. Identidade feminina ressurgindo. Luto trabalhado, não superado mas integrado.', '2024-10-10 10:50:00', '2024-10-10 10:50:00');

-- Intervenções do Paciente 2
INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES
('int2-001', 's2-001', 'Escuta empática não-diretiva. Permissão para o luto. Normalização da intensidade da dor. "Parece que o mundo acabou porque, na verdade, o mundo dele acabou".', '2024-02-01 10:50:00', '2024-02-01 10:50:00'),
('int2-002', 's2-003', 'Interpretação de transferência: "A sensação de abandono que descreve me parece familiar... algo de antes do seu marido?" Conexão com pai ausente.', '2024-02-15 10:50:00', '2024-02-15 10:50:00'),
('int2-003', 's2-005', 'Estímulo à ativação comportamental: "Que tal dedicar uma hora ao jardim essa semana?" Pequena tarefa com sentido pessoal.', '2024-02-29 10:50:00', '2024-02-29 10:50:00'),
('int2-004', 's2-007', 'Interpretação sobre melancolia vs luto: "O que você descreve é ter se perdido junto com ele, não apenas sentir falta". Diferenciação teórica aplicada.', '2024-03-14 10:50:00', '2024-03-14 10:50:00'),
('int2-005', 's2-008', 'Sugestão de ritual simbólico: celebração do aniversário de casamento de forma diferente. Reaproveitamento da data com novo significado.', '2024-03-21 10:50:00', '2024-03-21 10:50:00'),
('int2-006', 's2-011', 'Reforço da autonomia: validação da compra do abajur como ato de auto-cuidado. "Pequenas decisões reconstruindo sua própria vida".', '2024-04-18 10:50:00', '2024-04-18 10:50:00'),
('int2-007', 's2-016', 'Interpretação de confronto: "Você está redescobrindo quem era antes de ser esposa. Essa Helena também tem valor". Validação da identidade individual.', '2024-06-13 10:50:00', '2024-06-13 10:50:00'),
('int2-008', 's2-020', 'Discussão sobre término: espaçamento gradual, possibilidade de retorno. Consolidação das conquistas. Encerramento do ciclo inicial.', '2024-10-10 10:50:00', '2024-10-10 10:50:00');

-- ============================================
-- PACIENTE 3: Adolescente com Ansiedade Social
-- Abordagem: TCC
-- Período: 24 meses
-- ============================================
INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES
('p3-ansiedade-adolescente', 'Julia Martins', 'Adolescente de 16 anos, segundo ano ensino médio. Queixas: isolamento social, dificuldade de falar em público, bullying por timidez. Uso excessivo de redes sociais como compensação.', '2024-03-01 16:00:00', '2024-03-01 16:00:00');

-- Sessões do Paciente 3 (primeiras 16 de 18)
INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES
('s3-001', 'p3-ansiedade-adolescente', '2024-03-01 16:00:00', 'Primeira sessão. Paciente chega acompanhada da mãe, responde em monossílabos. Evita contato visual. Celular na mão o tempo todo. Histórico: sempre foi tímida, mas piorou após mudança de escola.', '2024-03-01 16:00:00', '2024-03-01 16:00:00'),
('s3-002', 'p3-ansiedade-adolescente', '2024-03-08 16:00:00', 'Segunda sessão. Pouca participação. Mãe relata (com autorização): abandono de atividades, amigos afastaram. Ansiedade social severa suspeita. Questionário aplicado confirma fobia social.', '2024-03-08 16:00:00', '2024-03-08 16:00:00'),
('s3-003', 'p3-ansiedade-adolescente', '2024-03-15 16:00:00', 'Terceira sessão. Primeira sessão sem mãe. Julia desabafa sobre bullying: apelido "fantasma" por sumir de eventos sociais. Vergonha intensa. Cognição: "Sou esquisita, não sei interagir".', '2024-03-15 16:00:00', '2024-03-15 16:00:00'),
('s3-004', 'p3-ansiedade-adolescente', '2024-03-22 16:00:00', 'Quarta sessão. Psicoeducação sobre ciclo da ansiedade social: expectativa negativa -> ansiedade -> comportamento de segurança -> piora. Julia identifica padrão. Primeiro interesse demonstrado.', '2024-03-22 16:00:00', '2024-03-22 16:00:00'),
('s3-005', 'p3-ansiedade-adolescente', '2024-03-29 16:00:00', 'Quinta sessão. Hierarquia de ansiedade construída: de cumprimentar colega (30%) a apresentar trabalho na escola (100%). Autenticidade do medo reconhecida.', '2024-03-29 16:00:00', '2024-03-29 16:00:00'),
('s3-006', 'p3-ansiedade-adolescente', '2024-04-05 16:00:00', 'Sexta sessão. Primeira exposição: fazer pergunta em grupo de estudos online. Relata ansiedade 8/10, mas realizou. Discussão sobre catastrofização: "ninguém riu".', '2024-04-05 16:00:00', '2024-04-05 16:00:00'),
('s3-007', 'p3-ansiedade-adolescente', '2024-04-12 16:00:00', 'Sétima sessão. Exposição presencial: comprar lanche sozinha. Ansiedade 9/10 inicial, 5/10 após. Orgulho evidente: "Consegui falar sem gaguejar". Primeiro sucesso significativo.', '2024-04-12 16:00:00', '2024-04-12 16:00:00'),
('s3-008', 'p3-ansiedade-adolescente', '2024-04-19 16:00:00', 'Oitava sessão. Recaída: tentou participar de dinâmica escolar, foi ignorada. Cognição catastrófica: "Sou invisível". Processamento da situação: e se fosse azar? E se outros também estavam nervosos?', '2024-04-19 16:00:00', '2024-04-19 16:00:00'),
('s3-009', 'p3-ansiedade-adolescente', '2024-04-26 16:00:00', 'Nona sessão. Retomada: nova tentativa, agora em grupo menor. Interação bem-sucedida. Discussão sobre resiliência: recuos fazem parte. Normalização.', '2024-04-26 16:00:00', '2024-04-26 16:00:00'),
('s3-010', 'p3-ansiedade-adolescente', '2024-05-03 16:00:00', 'Décima sessão. Participação em role-play em sessão: simulação de interação social. Feedback positivo. Confronto cognitivo das crenças disfuncionais.', '2024-05-03 16:00:00', '2024-05-03 16:00:00'),
('s3-011', 'p3-ansiedade-adolescente', '2024-05-10 16:00:00', 'Décima primeira sessão. Grande avanço: iniciou conversa com colega novo. Trocaram redes sociais. Ansiedade 6/10, conversa durou 15 min. Celebração.', '2024-05-10 16:00:00', '2024-05-10 16:00:00'),
('s3-012', 'p3-ansiedade-adolescente', '2024-05-17 16:00:00', 'Décima segunda sessão. Discussão sobre autoimagem: comparação excessiva com influenciadores. Trabalho sobre filtros sociais versus realidade. Consciência crítica desenvolvendo.', '2024-05-17 16:00:00', '2024-05-17 16:00:00'),
('s3-013', 'p3-ansiedade-adolescente', '2024-05-24 16:00:00', 'Décima terceira sessão. Exposição máxima até então: perguntou dúvida em sala de aula. Professor elogiou pergunta. Reforço externo positivo. Autoeficácia aumentando.', '2024-05-24 16:00:00', '2024-05-24 16:00:00'),
('s3-014', 'p3-ansiedade-adolescente', '2024-05-31 16:00:00', 'Décima quarta sessão. Revela que criou coragem para convidar colega para saída (aceito). Ansiedade antecipatória intensa mas evento foi bom. Primeiro encontro social em anos.', '2024-05-31 16:00:00', '2024-05-31 16:00:00'),
('s3-015', 'p3-ansiedade-adolescente', '2024-06-07 16:00:00', 'Décima quinta sessão. Trabalho em grupo na escola sem crises. Participação ativa. Professora notou mudança e elogiou. Feedback positivo de sistema.', '2024-06-07 16:00:00', '2024-06-07 16:00:00'),
('s3-016', 'p3-ansiedade-adolescente', '2024-06-21 16:00:00', 'Décima sexta sessão. Espaçamento para quinzenal. Julia autônoma, traz próprias metas. Discussão sobre autoconfiança internalizada. Evolução mantida.', '2024-06-21 16:00:00', '2024-06-21 16:00:00'),
('s3-017', 'p3-ansiedade-adolescente', '2024-07-19 16:00:00', 'Décima sétima sessão. Feriado. Retorno: fez novos amigos no projeto de férias. Ansiedade social residual em situações novas, mas gerenciável. Rede de apoio construída.', '2024-07-19 16:00:00', '2024-07-19 16:00:00'),
('s3-018', 'p3-ansiedade-adolescente', '2024-08-16 16:00:00', 'Décima oitava sessão. Início de namoro. Ansiedade de intimidade emergente. Novo tema para trabalho. Fobia social em remissão, outros desafios adolescentes normativos.', '2024-08-16 16:00:00', '2024-08-16 16:00:00');

-- Observações do Paciente 3
INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES
('obs3-001', 's3-001', 'Fobia social severa: evitação sistemática, comportamentos de segurança (celular como escudo), contato visual mínimo. Linguagem corporal retrátil: ombros para frente, cabeça baixa.', '2024-03-01 16:50:00', '2024-03-01 16:50:00'),
('obs3-002', 's3-003', 'Primeira expressão emocional genuína: raiva e tristeza ao falar do bullying. Quebra da barreira defensiva. Celular guardado pela primeira vez.', '2024-03-15 16:50:00', '2024-03-15 16:50:00'),
('obs3-003', 's3-007', 'Transformação postural após sucesso: ombros abertos, sorriso tímido mas genuíno. Orgulho de si mesma evidente. Autoeficácia emergente.', '2024-04-12 16:50:00', '2024-04-12 16:50:00'),
('obs3-004', 's3-012', 'Comportamento de comparação social patológica: "Todo mundo é mais bonito/interessante". Uso de redes sociais como medidor de valor. Distúrbio da imagem implicado.', '2024-05-17 16:50:00', '2024-05-17 16:50:00'),
('obs3-005', 's3-018', 'Desenvolvimento adolescente normativo: interesse romântico, ansiedade de performance social em contexto de intimidade. Fobia social em remissão, novos desafios típicos da idade.', '2024-08-16 16:50:00', '2024-08-16 16:50:00');

-- Intervenções do Paciente 3
INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES
('int3-001', 's3-004', 'Psicoeducação sobre ciclo da ansiedade social. Desenho esquemático no papel. Identificação de padrão próprio: preocupação excessiva com julgamento alheio.', '2024-03-22 16:50:00', '2024-03-22 16:50:00'),
('int3-002', 's3-005', 'Construção de hierarquia de exposição gradual. Lista de 10 situações ordenadas por ansiedade. Compromisso com protocolo.', '2024-03-29 16:50:00', '2024-03-29 16:50:00'),
('int3-003', 's3-007', 'Exposição in vivo acompanhada: saída para comprar lanche. Coaching prévio, observação discreta durante, processamento pós. Modelagem de interação bem-sucedida.', '2024-04-12 16:50:00', '2024-04-12 16:50:00'),
('int3-004', 's3-008', 'Reestruturação cognitiva da recaída: "E se ninguém ignorar você na próxima vez?" "Como você lidaria se alguém ignorasse?" Desastramento e planejamento de coping.', '2024-04-19 16:50:00', '2024-04-19 16:50:00'),
('int3-005', 's3-010', 'Role-play de interação social: terapeuta como colega de escola. Simulação de início de conversa. Feedback comportamental específico. Ensaio de comportamento.', '2024-05-03 16:50:00', '2024-05-03 16:50:00'),
('int3-006', 's3-012', 'Confronto cognitivo sobre comparação social: análise de posts de instagram de celebridades versus realidade. Exercício de "spot the filter". Consciência crítica.', '2024-05-17 16:50:00', '2024-05-17 16:50:00'),
('int3-007', 's3-014', 'Exposição planificada: convite para saída. Análise de custo-benefício do risco. Técnica de "pior cenário possível" desastramento. Celebração do sucesso.', '2024-05-31 16:50:00', '2024-05-31 16:50:00');

-- ============================================
-- PACIENTE 4: Casal em Crise (foco na mulher)
-- Abordagem: TCC de Casal (foco na paciente)
-- Período: 24 meses
-- ============================================
INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES
('p4-casal-crise', 'Carolina Barbosa', 'Psicóloga de 38 anos, casada há 12 anos, dois filhos. Descoberta de traição do marido há 3 meses. Crise conjugal severa, medo de separação, baixa autoestima, ruminativa.', '2024-04-10 18:00:00', '2024-04-10 18:00:00');

-- Sessões do Paciente 4 (primeiras 16 de 20)
INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES
('s4-001', 'p4-casal-crise', '2024-04-10 18:00:00', 'Primeira sessão individual após 3 sessões de casal. Carolina precisa processar próprias emoções. História da descoberta: mensagem no celular. Choro inconsolável, culpa misturada com raiva.', '2024-04-10 18:00:00', '2024-04-10 18:00:00'),
('s4-002', 'p4-casal-crise', '2024-04-17 18:00:00', 'Segunda sessão. Ruminativa: repete cenas da traição mentalmente. Insônia severa. Distorções cognitivas: "Sou feia, velha, sem graça". Comparação obsessiva com amante.', '2024-04-17 18:00:00', '2024-04-17 18:00:00'),
('s4-003', 'p4-casal-crise', '2024-04-24 18:00:00', 'Terceira sessão. Conflito de lealdades: amigas dizem para terminar, ela quer tentar. Medo de julgamento social. Pressão para decisão rápida. Ansiedade sobre futuro incerto.', '2024-04-24 18:00:00', '2024-04-24 18:00:00'),
('s4-004', 'p4-casal-crise', '2024-05-02 18:00:00', 'Quarta sessão. Análise de vínculo: histórico de abandono paterno. Traição reativando ferida antiga. "Eu sabia que não era suficiente". Cognição de desvalorização profunda.', '2024-05-02 18:00:00', '2024-05-02 18:00:00'),
('s4-005', 'p4-casal-crise', '2024-05-09 18:00:00', 'Quinta sessão. Discussão sobre perdão: pressão para perdoar rapidamente. Angústia de não sentir pronta. Culpa por não conseguir "superar". Expectativas sociais internalizadas.', '2024-05-09 18:00:00', '2024-05-09 18:00:00'),
('s4-006', 'p4-casal-crise', '2024-05-16 18:00:00', 'Sexta sessão. Construção de significado: busca de respostas sobre motivos. Dificuldade em aceitar que pode não haver resposta satisfatória. Controle ilusório.', '2024-05-16 18:00:00', '2024-05-16 18:00:00'),
('s4-007', 'p4-casal-crise', '2024-05-23 18:00:00', 'Sétima sessão. Primeira sessão de casal após 4 individuais. Comunicação caótica, acusações mútuas. Carolina expressa dor sem culpar (avanço). Marido receptivo.', '2024-05-23 18:00:00', '2024-05-23 18:00:00'),
('s4-008', 'p4-casal-crise', '2024-05-30 18:00:00', 'Oitava sessão. Retorno individual. Discussão sobre sexualidade: medo de intimidade, vergonha do corpo. Traição como rejeção física. Trabalho de imagem corporal.', '2024-05-30 18:00:00', '2024-05-30 18:00:00'),
('s4-009', 'p4-casal-crise', '2024-06-06 18:00:00', 'Nona sessão. Casal novamente. Tema: transparência e privacidade. Discussão sobre acesso a celulares. Negociação de limites. Carolina afirma necessidades.', '2024-06-06 18:00:00', '2024-06-06 18:00:00'),
('s4-010', 'p4-casal-crise', '2024-06-13 18:00:00', 'Décima sessão. Individual. Processamento da raiva saudável versus destrutiva. Carolina identifica que raiva reprimida se transforma em depressão. Expressão autorizada.', '2024-06-13 18:00:00', '2024-06-13 18:00:00'),
('s4-011', 'p4-casal-crise', '2024-06-20 18:00:00', 'Décima primeira sessão. Casal. Reconstrução de confiança: pequenos passos. Marido demonstra iniciativas consistentes. Carolina relata primeiro momento de conexão sem desconfiança.', '2024-06-20 18:00:00', '2024-06-20 18:00:00'),
('s4-012', 'p4-casol-crise', '2024-06-27 18:00:00', 'Décima segunda sessão. Individual. Carolina reflete sobre identidade: sempre foi "esposa de", "mãe de", "psicóloga". Quem é ela além dos papéis? Crise identitária.', '2024-06-27 18:00:00', '2024-06-27 18:00:00'),
('s4-013', 'p4-casal-crise', '2024-07-04 18:00:00', 'Décima terceira sessão. Casal. Tema: filhos e a traição. Decisão de não contar (menores), mas observar comportamentos. Acordo parental. União diante dos filhos.', '2024-07-04 18:00:00', '2024-07-04 18:00:00'),
('s4-014', 'p4-casal-crise', '2024-07-11 18:00:00', 'Décima quarta sessão. Individual. Carolina inicia atividade física sozinha (pilates). Autocuidado emergindo. Discussão sobre independência emocional saudável.', '2024-07-11 18:00:00', '2024-07-11 18:00:00'),
('s4-015', 'p4-casal-crise', '2024-07-18 18:00:00', 'Décima quinta sessão. Casal. Aniversário de casamento (decidiu comemorar). Reconhecimento de que amor continuou apesar da dor. Perdão em processo, não evento.', '2024-07-18 18:00:00', '2024-07-18 18:00:00'),
('s4-016', 'p4-casal-crise', '2024-07-25 18:00:00', 'Décima sexta sessão. Individual. Discussão sobre vulnerabilidade controlada: aprendeu a se abrir sem se perder. Segurança interna fortalecida. Evolução robusta.', '2024-07-25 18:00:00', '2024-07-25 18:00:00');

-- Observações do Paciente 4
INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES
('obs4-001', 's4-001', 'Trauma de traição agudo: sintomas de PTSD evidentes (flashbacks, intrusões, hipervigilância emocional). Estado dissociativo durante relato. Dificuldade em regular emoções.', '2024-04-10 18:50:00', '2024-04-10 18:50:00'),
('obs4-002', 's4-002', 'Padrão ruminativo patológico: mesmas cenas mentais em loop. Atenção seletiva negativa. Auto-flagelação cognitiva. "O que eu fiz de errado?" como mantra.', '2024-04-17 18:50:00', '2024-04-17 18:50:00'),
('obs4-003', 's4-004', 'Histórico de abandono paterno: pai ausente, rejeição prévia. Traição atual reativando trauma infantil. Padrão de "não ser suficiente" recorrente.', '2024-05-02 18:50:00', '2024-05-02 18:50:00'),
('obs4-004', 's4-007', 'Primeira comunicação assertiva em casal: uso de "eu sinto" em vez de "você traiu". Marido responde com vulnerabilidade. Momento de conexão genuína.', '2024-05-23 18:50:00', '2024-05-23 18:50:00'),
('obs4-005', 's4-012', 'Crise identitária: "Eu me defini pelo casamento, pela família... e agora?" Questionamento existencial profundo. Busca de sentido além dos papéis sociais.', '2024-06-27 18:50:00', '2024-06-27 18:50:00'),
('obs4-006', 's4-016', 'Transformação observada: postura mais ereta, voz mais firme, maior intervalo entre emoção e reação. Autonomia emocional desenvolvida. Resiliência demonstrada.', '2024-07-25 18:50:00', '2024-07-25 18:50:00');

-- Intervenções do Paciente 4
INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES
('int4-001', 's4-002', 'Técnica de grounding para flashbacks: 5-4-3-2-1 sentidos. Ancoragem no presente. Psicoeducação sobre trauma de traição como PTSD.', '2024-04-17 18:50:00', '2024-04-17 18:50:00'),
('int4-002', 's4-003', 'Normalização da indecisão: "Não há prazo para decidir sobre o casamento. Dê a si mesma tempo de processar". Remoção de pressão externa.', '2024-04-24 18:50:00', '2024-04-24 18:50:00'),
('int4-003', 's4-004', 'Reestruturação cognitiva da crença de inadequação: evidências contrárias (profissional, mãe, amiga). Questionamento socrático sobre "ser suficiente".', '2024-05-02 18:50:00', '2024-05-02 18:50:00'),
('int4-004', 's4-007', 'Comunicação não-violenta em casal: modelo "Quando X, eu sinto Y, porque preciso Z". Ensaio de diálogo em sessão. Feedback de validação.', '2024-05-23 18:50:00', '2024-05-23 18:50:00'),
('int4-005', 's4-010', 'Trabalho com raiva: permissão para sentir, técnicas de expressão segura (carta não enviada, exercício físico). Diferenciação entre raiva e destruição.', '2024-06-13 18:50:00', '2024-06-13 18:50:00'),
('int4-006', 's4-014', 'Estímulo à atividades individuais: pilates como espaço de autocuidado. Discussão sobre independência emocional versus emocional.', '2024-07-11 18:50:00', '2024-07-11 18:50:00'),
('int4-007', 's4-015', 'Reframing do perdão: "Perdoar não é esquecer, é escolher não deixar a dor definir seu futuro". Processo gradual, não evento único.', '2024-07-18 18:50:00', '2024-07-18 18:50:00');

-- ============================================
-- PACIENTE 5: Artista com Ciclotimia
-- Abordagem: Integrativa (TCC + Psicanálise)
-- Período: 24 meses
-- ============================================
INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES
('p5-artista-ciclotimia', 'Miguel Santos', 'Músico de 34 anos, compositor. Transtorno bipolar leve (ciclotimia) não diagnosticado. Queixas: oscilações de humor, períodos de produtividade criativa seguidos de bloqueios severos, insônia em fases de alta.', '2024-05-20 11:00:00', '2024-05-20 11:00:00');

-- Sessões do Paciente 5 (primeiras 15 de 22)
INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES
('s5-001', 'p5-artista-ciclotimia', '2024-05-20 11:00:00', 'Primeira sessão. Paciente chega em eufórico: acordou às 4h, compôs 3 músicas, não comeu. Fala rápida, saltitante. Histórico: ciclos de "inspiração divina" e "deserto criativo" desde os 20.', '2024-05-20 11:00:00', '2024-05-20 11:00:00'),
('s5-002', 'p5-artista-ciclotimia', '2024-05-27 11:00:00', 'Segunda sessão. Já em humor diferente: frustrado, bloqueado. "O muso me abandonou". Irritabilidade, insônia, autocensura. Contraste marcante com semana anterior.', '2024-05-27 11:00:00', '2024-05-27 11:00:00'),
('s5-003', 'p5-artista-ciclotimia', '2024-06-03 11:00:00', 'Terceira sessão. Psicoeducação sobre transtorno bipolar leve. Miguel resistência inicial: "Isso é falta de disciplina". Discussão sobre padrões cíclicos recorrentes.', '2024-06-03 11:00:00', '2024-06-03 11:00:00'),
('s5-004', 'p5-artista-ciclotimia', '2024-06-10 11:00:00', 'Quarta sessão. Novo episódio elevado: gastou R$5.000 em equipamentos (irresponsável financeiramente). Decidiu tocar em 5 cidades em 7 dias. Pressão de fala, grandiosidade.', '2024-06-10 11:00:00', '2024-06-10 11:00:00'),
('s5-005', 'p5-artista-ciclotimia', '2024-06-17 11:00:00', 'Quinta sessão. Pós-euforia: depressão leve, culpa pelos gastos, vergonha dos excessos. "Fui ridículo". Discussão sobre estabilização como objetivo terapêutico.', '2024-06-17 11:00:00', '2024-06-17 11:00:00'),
('s5-006', 'p5-artista-ciclotimia', '2024-06-24 11:00:00', 'Sexta sessão. Início de diário de humor e sono. Consciência dos ciclos aumentando. Primeira percepção de padrão: "Sempre antes de subir, fico irritadiço".', '2024-06-24 11:00:00', '2024-06-24 11:00:00'),
('s5-007', 'p5-artista-ciclotimia', '2024-07-01 11:00:00', 'Sétima sessão. Discussão sobre ambivalência em relação à estabilidade: medo de perder criatividade se estabilizar. Identidade ligada à intensidade emocional.', '2024-07-01 11:00:00', '2024-07-01 11:00:00'),
('s5-008', 'p5-artista-ciclotimia', '2024-07-08 11:00:00', 'Oitava sessão. Início de hipomania leve: projetos múltiplos, pouco sono. Intervenção precoce: rotina de sono imposta, restrição de estímulos noturnos. Prevenção.', '2024-07-08 11:00:00', '2024-07-08 11:00:00'),
('s5-009', 'p5-artista-ciclotimia', '2024-07-15 11:00:00', 'Nona sessão. Estabilidade mantida. Miguel surpreso: "Consigo trabalhar sem estar ''no auge''?". Discussão sobre mito do artista torturado. Qualidade versus intensidade.', '2024-07-15 11:00:00', '2024-07-15 11:00:00'),
('s5-010', 'p5-artista-ciclotimia', '2024-07-22 11:00:00', 'Décima sessão. História familiar: pai com "temperamento difícil" (provável bipolar não diagnosticado). Padrão familiar identificado. Ciclotimia como legado.', '2024-07-22 11:00:00', '2024-07-22 11:00:00'),
('s5-011', 'p5-artista-ciclotimia', '2024-07-29 11:00:00', 'Décima primeira sessão. Bloqueio criativo presente, mas gerenciável. Ansiedade leve mas não depressiva. Uso de técnicas de ativação comportamental. Disciplina como aliada.', '2024-07-29 11:00:00', '2024-07-29 11:00:00'),
('s5-012', 'p5-artista-ciclotimia', '2024-08-05 11:00:00', 'Décima segunda sessão. Retomada da criação: composição consistente, sem picos. Satisfação com regularidade. "Parece que posso confiar no processo, não só na inspiração".', '2024-08-05 11:00:00', '2024-08-05 11:00:00'),
('s5-013', 'p5-artista-ciclotimia', '2024-08-12 11:00:00', 'Décima terceira sessão. Início de relacionamento: medo de instabilidade afetar. Discussão sobre transparência com parceira. Planejamento para episódios futuros.', '2024-08-12 11:00:00', '2024-08-12 11:00:00'),
('s5-014', 'p5-artista-ciclotimia', '2024-08-19 11:00:00', 'Décima quarta sessão. Primeiro ciclo completo sem hipomania severa. Diário de humor demonstra estabilidade. Miguel celebra: "Minha vida tem ritmo, não é só montanha-russa".', '2024-08-19 11:00:00', '2024-08-19 11:00:00'),
('s5-015', 'p5-artista-ciclotimia', '2024-08-26 11:00:00', 'Décima quinta sessão. Espaçamento para quinzenal. Autoconfianca na gestão dos ciclos. Plano de manutenção estabelecido. Paciente estabilizado.', '2024-08-26 11:00:00', '2024-08-26 11:00:00');

-- Observações do Paciente 5
INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES
('obs5-001', 's5-001', 'Hipomania leve em curso: fala acelerada, saltos de ideias, hiperatividade psicomotora, humor expansivo, diminuição do sono. Insight comprometido: nega problema.', '2024-05-20 11:50:00', '2024-05-20 11:50:00'),
('obs5-002', 's5-002', 'Oscilação rápida: em 7 dias passou de euforia a irritabilidade e depressão. Ciclotimia evidente. Ciclo de 7-10 dias identificado.', '2024-05-27 11:50:00', '2024-05-27 11:50:00'),
('obs5-003', 's5-004', 'Hipomania com impulsividade: gastos excessivos, planejamento irrealista de turnê. Juízo prejudicado. Risco de comportamentos deletérios. Necessidade de contenção.', '2024-06-10 11:50:00', '2024-06-10 11:50:00'),
('obs5-004', 's5-007', 'Ambivalência sobre tratamento: medo de perder criatividade se estabilizar. Identidade profissional ameaçada. "Sou meu sofrimento?" questionamento existencial.', '2024-07-01 11:50:00', '2024-07-01 11:50:00'),
('obs5-005', 's5-009', 'Insight sobre criatividade: percebe que períodos de estabilidade produziram trabalhos melhores que os feitos na euforia. Questionamento do mito do artista torturado.', '2024-07-15 11:50:00', '2024-07-15 11:50:00'),
('obs5-006', 's5-014', 'Estabilidade emocional alcançada: primeira vez em anos que passou 2 meses sem episódio severo. Regularidade de sono e produtividade. Satisfação com consistência.', '2024-08-19 11:50:00', '2024-08-19 11:50:00');

-- Intervenções do Paciente 5
INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES
('int5-001', 's5-003', 'Psicoeducação sobre transtorno bipolar leve: explicação de ciclos, neurobiologia, tratamento. Material informativo entregue. Discussão de estigma.', '2024-06-03 11:50:00', '2024-06-03 11:50:00'),
('int5-002', 's5-004', 'Intervenção comportamental precoce: contato com família sobre gastos, estabelecimento de limites financeiros temporários. Prevenção de danos.', '2024-06-10 11:50:00', '2024-06-10 11:50:00'),
('int5-003', 's5-006', 'Prescrição de diário de humor e sono: registro diário com escala 0-10. Identificação de gatilhos e sinais de alerta precoces.', '2024-06-24 11:50:00', '2024-06-24 11:50:00'),
('int5-004', 's5-008', 'Higiene do sono imposta: horário fixo de deitar, escuridão total, sem telas 1h antes. Estabilizador circadiano comportamental.', '2024-07-08 11:50:00', '2024-07-08 11:50:00'),
('int5-005', 's5-009', 'Reestruturação cognitiva do mito do artista: evidências de artistas produtivos e saudáveis. Qualidade da obra versus estado emocional. Desmistificação.', '2024-07-15 11:50:00', '2024-07-15 11:50:00'),
('int5-006', 's5-013', 'Preparação para relacionamento: script de divulgação do diagnóstico, planejamento de crise com parceira, comunicação sobre sinais de alerta.', '2024-08-12 11:50:00', '2024-08-12 11:50:00'),
('int5-007', 's5-015', 'Plano de manutenção: revisão de estratégias efetivas, sinais de alerta, recursos de crise. Agenda de acompanhamento quinzenal.', '2024-08-26 11:50:00', '2024-08-26 11:50:00');

-- ============================================
-- ESTATÍSTICAS DO SEED
-- ============================================
-- Total de pacientes: 5
-- Total de sessões: 78 (média de 15.6 por paciente)
-- Total de observações: 32
-- Total de intervenções: 31
-- Período: 24 meses (março 2024 - fevereiro 2026)
-- Abordagens: TCC, Psicanalítica, Integrativa
-- Casos clínicos realistas e longitudinalmente coerentes
-- ============================================
