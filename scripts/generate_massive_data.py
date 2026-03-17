#!/usr/bin/env python3
"""
Gerador Massivo de Dados Clínicos para Arandu
Cria 500 pacientes, 50.000+ sessões, 100.000+ observações, 50.000+ intervenções
"""

import random
from datetime import datetime, timedelta

# Configurações
NUM_PATIENTS = 500
MIN_SESSIONS = 100
MAX_SESSIONS = 150
MIN_OBS = 2
MAX_OBS = 4
MIN_INT = 1
MAX_INT = 3

# Dados
FIRST_NAMES = [
    "João",
    "Maria",
    "José",
    "Ana",
    "Antônio",
    "Francisca",
    "Carlos",
    "Adriana",
    "Paulo",
    "Juliana",
    "Lucas",
    "Fernanda",
    "Pedro",
    "Patrícia",
    "Marcos",
    "Mariana",
    "Gabriel",
    "Vanessa",
    "Rafael",
    "Amanda",
    "Bruno",
    "Bruna",
    "Eduardo",
    "Beatriz",
    "Felipe",
    "Carolina",
    "Gustavo",
    "Daniela",
    "André",
    "Larissa",
]

LAST_NAMES = [
    "Silva",
    "Santos",
    "Oliveira",
    "Souza",
    "Pereira",
    "Costa",
    "Rodrigues",
    "Almeida",
    "Lima",
    "Carvalho",
    "Araújo",
    "Ferreira",
    "Gomes",
    "Ribeiro",
    "Martins",
    "Barbosa",
    "Alves",
    "Rocha",
    "Cardoso",
    "Correia",
]

CLINICAL_TYPES = [
    "Transtorno de Ansiedade Generalizada",
    "Depressão Moderada",
    "Transtorno do Pânico",
    "Fobia Social",
    "Transtorno Obsessivo-Compulsivo",
    "TEPT",
    "Luto Complicado",
    "Crise de Meia-Idade",
    "Transtorno Bipolar",
    "Borderline",
    "Dependência Emocional",
    "Burnout Profissional",
    "Adoção e Identidade",
    "Separação e Divórcio",
    "Cuidador Primário Esgotado",
]

OBS_TEMPLATES = [
    "Paciente apresenta tensão muscular visível, postura rígida e movimentos controlados. Sinais de ansiedade evidentes.",
    "Preocupações circulares evidentes, mesmo tema retorna múltiplas vezes. Dificuldade em conter pensamentos.",
    "Sinais de hipervigilância presentes, reage a estímulos mínimos. Estado de alerta elevado observado.",
    "Insônia relatada afetando humor e energia. Rotina de sono desregulada compromete funcionamento diário.",
    "Progresso notável observado, postura relaxada e respiração diafragmática presentes pela primeira vez.",
    "Uso efetivo de técnicas de relaxamento demonstrado. Internalização de coping confirmada na prática.",
    "Recaída situacional identificada, porém recuperação mais rápida que episódios anteriores. Resiliência em desenvolvimento.",
    "Insight emergente sobre padrões automáticos. Consciência metacognitiva demonstrada na sessão.",
    "Psicomotricidade retardada evidente, movimentos lentos e voz baixa. Postura encolhida e contato visual mínimo.",
    "Anedonia presente, nenhuma atividade gera prazer relatado. Descreve dias monótonos sem motivação.",
    "Autocrítica severa identificada, linguagem depreciativa sobre si mesmo. Cognições disfuncionais presentes.",
    "Isolamento social relatado, evita contatos e cancela compromissos. Vínculos afetivos mantidos superficialmente.",
    "Primeiros sinais de melhora, trouxe assunto externo não solicitado. Interesse diminuto mas presente.",
    "Ativação comportamental em curso, retomou atividade física. Ciclo virtuoso de humor e comportamento iniciando.",
    "Humor labil observado, oscilação entre tristeza profunda e irritabilidade. Choro e raiva alternam-se.",
    "Cognições mais flexíveis demonstradas, consegue questionar pensamentos automáticos negativos.",
    "Reconexão com valores identificada, discussão sobre sentido de vida além do humor deprimido.",
    "Crise de pânico relatada na semana, sintomas físicos intensos. Correu ao pronto-socorro, exames normais.",
    "Evitação de locais evidente, espaço vital diminuído progressivamente. Restrição de atividades diárias.",
    "Hipervigilância corporal presente, monitora constantemente batimentos e respiração.",
    "Antecipação ansiosa identificada, medo de ter medo. Ciclo de feedback estabelecido e reforçado.",
    "Exposição bem-sucedida realizada, enfrentou situação temida. Ansiedade elevou mas diminuiu com tempo.",
    "Reestruturação cognitiva em andamento, questiona interpretações catastróficas de forma mais realista.",
    "Diário de pensamentos mantido com registros detalhados. Padrões de antecedentes e consequências mapeados.",
    "Recuperação de autonomia em progresso, retomou atividades anteriormente evitadas com confiança.",
    "Ansiedade social severa evidente, contato visual fugidio e postura retrátil. Comportamentos de segurança presentes.",
    "Evitação sistemática identificada, recusa convites e comportamentos minimizam exposição social.",
    "Medo de julgamento presente, preocupação obsessiva com avaliação alheia afeta interações.",
    "Isolamento compensatório observado, redes sociais substituem interação presencial.",
    "Exposição gradual realizada com sucesso, interação planejada bem executada. Ansiedade inicial alta mas habituável.",
    "Feedback positivo recebido, interação social foi melhor que esperado. Reforço externo contradiz expectativas.",
    "Assertividade emergente observada, consegue dizer não sem culpa excessiva. Limites saudáveis estabelecidos.",
    "Rituais compulsivos evidentes, repetições de ações e verificações excessivas causando perda de tempo.",
    "Obsessões intrusivas relatadas, pensamentos indesejados invadem constantemente o fluxo mental.",
    "Ansiedade de contaminação intensa, medo de sujeira e germes levando a comportamentos excessivos.",
    "Prevenção de resposta realizada, resistiu à compulsão e tolerou ansiedade até habituação.",
    "Exposição intencional em andamento, confronto com tema temido sem comportamentos de segurança.",
    "Insight sobre o ciclo do TOC demonstrado, reconhece que compulsões reforçam o transtorno.",
    "Flashbacks relatados com memórias intrusivas do trauma. Dissociação e sensação de reviver evento.",
    "Pesadelos frequentes perturbando sono e causando medo de dormir. Insônia secundária ao trauma.",
    "Hipervigilância excessiva com estado de alerta constante. Sinais de perigo procurados ativamente.",
    "Processamento do trauma em curso, narração coesa do evento e integração na história de vida.",
    "Técnica de grounding utilizada, ancoragem no presente reduzindo flashbacks e dissociação.",
    "Reestruturação de culpa em andamento, questionamento de autoacusações excessivas.",
    "Luto persistente identificado, identificação simbiótica com falecido após período prolongado.",
    "Melancolia evidente com perda de interesse em vida sem a pessoa amada. Sentimento de falta de propósito.",
    "Culpa de sobrevivente presente, questiona por que ficou vivo e critica ações passadas.",
    "Rituais de despedida realizados, celebração de datas significativas de forma diferente mas presente.",
    "Reconexão com vida emergindo, pequenos prazeres retornando como jardim e caminhadas.",
    "Questionamento existencial intenso sobre sentido da vida e insatisfação com conquistas aparentes.",
    "Medo de morte e ansiedade de finitude presentes. Pressão percebida para realizar mudanças.",
    "Impulsividade compensatória observada, mudanças radicais consideradas sem planejamento adequado.",
    "Exploração de valores em curso, identificação do que realmente importa além de sucessos externos.",
    "Reavaliação de relacionamentos identificada, busca por vínculos autênticos em vez de sociais.",
    "Episódio hipomaníaco presente, humor elevado com energia excessiva e redução do sono.",
    "Fala pressada com ideias saltitantes e dificuldade em seguir linha racional.",
    "Impulsividade com gastos excessivos e múltiplos projetos simultâneos. Julgamento prejudicado.",
    "Queda depressiva após euforia, humor deprimido com culpa pelos excessos e fadiga intensa.",
    "Estabilização de humor alcançada, rotina estruturada sendo mantida com regularidade.",
    "Instabilidade emocional severa, oscilações rápidas e intensas afetando funcionamento.",
    "Medo de abandono intenso com reações desproporcionais à percepção de rejeição.",
    "Impulsividade autodestrutiva presente, comportamentos de risco e coping disfuncional.",
    "Identidade difusa demonstrada, interesses e valores em constante mudança.",
    "Dificuldade extrema em estar só, ansiedade de separação levando a telefonemas constantes.",
    "Autoestima baseada no outro identificada, validação condicional afetando autoconcepção.",
    "Relações tóxicas repetidas, mesmo padrão com parceiros diferentes evidenciando ciclo.",
    "Primeira experiência de autonomia realizada, viagem sozinha apesar da resistência.",
    "Exaustão severa presente, energia física e emocional completamente esgotada.",
    "Cinismo profissional evidente, atitudes negativas sobre trabalho e colegas.",
    "Ineficácia percebida com sensação de não fazer diferença ou ter competência.",
    "Licença médica em curso, pausa forçada para recuperação do esgotamento.",
    "Reavaliação de prioridades identificada, questionamento sobre o que realmente importa na vida.",
    "Busca por origens presente, necessidade de saber história afetando identidade.",
    "Não pertencimento relatado, dificuldade em se sentir realmente parte da família.",
    "Luto do relacionamento intenso após término, vida construída jun sendo desmontada.",
    "Medo de solidão presente com ansiedade sobre reconstruir vida sozinho.",
    "Raiva e traição identificadas, sentimento de injustiça sobre o término.",
    "Esgotamento de cuidador presente, anos de sobrecarga afetando saúde física e mental.",
    "Ressentimento identificado, sensação de que vida parou para cuidar do outro.",
    "Perda de identidade demonstrada, esquecimento de quem era além do papel de cuidador.",
]

INT_TEMPLATES = [
    "Psicoeducação sobre mecanismo da ansiedade e ciclo de feedback entre pensamentos e comportamentos.",
    "Técnica de respiração diafragmática 4-7-8 demonstrada e prescrita para prática diária.",
    "Exposição gradual in vivo planejada com construção de hierarquia de ansiedade.",
    "Reestruturação cognitiva realizada com identificação de pensamentos automáticos e questionamento socrático.",
    "Mindfulness de atenção plena exercitado para desidentificação dos pensamentos.",
    "Relaxamento progressivo de Jacobson ensinado com gravação fornecida para casa.",
    "Técnica de postponement da preocupação aplicada para reduzir frequência intrusiva.",
    "Exposição interoceptiva realizada em sessão para habituação às sensações físicas.",
    "Comunicação assertiva treinada com prática de expressão de necessidades e limites.",
    "Plano de manutenção revisado com estratégias aprendidas e sinais de alerta identificados.",
    "Ativação comportamental programada com atividades prazerosas e de realização.",
    "Registro de pensamentos orientado com diário de situações e emoções.",
    "Comportamento oposto às vontades prescrito, confiando que humor seguirá comportamento.",
    "Terapia comportamental de ativação com análise funcional e intervenções nos gatilhos.",
    "Explicação do ciclo do pânico e normalização dos sintomas físicos como reações benignas.",
    "Exposição interoceptiva com hiperventilação controlada para habituação.",
    "Exposição in vivo gradual realizada com acompanhamento terapêutico.",
    "Reestruturação cognitiva do medo do medo com desastramento do pior cenário.",
    "Diário de auto-observação prescrito para registro de crises e análise de padrões.",
    "Técnica de acatastrofização aplicada questionando probabilidades reais.",
    "Psicoeducação sobre ciclo da ansiedade social e expectativas negativas.",
    "Exposição gradual social planejada desde cumprimentos até falar em público.",
    "Role-play de interação social ensaiado em sessão com feedback comportamental.",
    "Experimento de atenção externa em vez de introspecção excessiva.",
    "Comportamentos de segurança identificados e eliminados gradualmente.",
    "Psicoeducação sobre TOC e ciclo obsessão-compulsão-alívio temporário.",
    "Prevenção de resposta aplicada com resistência à realização da compulsão.",
    "Exposição e prevenção de resposta realizadas com confronto gradual.",
    "Exposição imaginária com narração detalhada de cenas temidas.",
    "Técnica de aceitação da incerteza desenvolvida para tolerar dúvidas.",
    "Psicoeducação sobre TEPT e normalização dos sintomas.",
    "Técnica de grounding aplicada para ancoragem no presente.",
    "Estabilização e segurança construídas com janela de tolerância.",
    "Processamento do trauma com narração coesa da história.",
    "Técnica de cadeira vazia utilizada para diálogo simbólico.",
    "Reestruturação cognitiva da culpa com questionamento de responsabilidade.",
    "Escuta empática não-diretiva para validação da dor do luto.",
    "Rituais de despedida criados para marcar transições significativas.",
    "Ativação comportamental para retorno gradual às atividades.",
    "Trabalho com culpa e arrependimentos para perdão próprio.",
    "Escuta do questionamento existencial e validação da busca por sentido.",
    "Exercício de valores para identificação do que realmente importa.",
    "Análise de transição de vida com avaliação de conquistas e metas.",
    "Reavaliação de relacionamentos e escolha de vínculos alimentares.",
    "Psicoeducação sobre bipolaridade e monitoramento de humor.",
    "Higiene do sono regularizada como estabilizador de humor.",
    "Estabilização de rotina com horários regulares estabelecidos.",
    "Identificação de gatilhos para prevenção de episódios.",
    "Intervenção precoce aplicada para sinais de alteração de humor.",
    "Psicoeducação sobre borderline e validação da intensidade emocional.",
    "Treinamento de mindfulness para consciência plena.",
    "Tolerância à angústia desenvolvida com técnicas de distress.",
    "Regulação emocional treinada com identificação e nomeação.",
    "Eficácia interpessoal desenvolvida para assertividade.",
    "Análise de história de apego para compreensão de padrões.",
    "Técnica de tolerância à solidão com exposição gradual.",
    "Reconstrução de autoestima independente da validação externa.",
    "Estabelecimento de limites saudáveis no relacionamento.",
    "Avaliação de burnout com identificação de fatores de risco.",
    "Psicoeducação sobre recuperação e necessidade de pausa.",
    "Análise de demandas e recursos para equilíbrio.",
    "Estabelecimento de limites de disponibilidade.",
    "Reavaliação de valores profissionais e sentido do trabalho.",
    "Escuta da história de adoção e normalização da experiência.",
    "Exploração de busca por origens com preparo para resultados.",
    "Reconstrução de narrativa de vida com história coesa.",
    "Trabalho com abandono para ressignificação não pessoal.",
    "Escuta da história de término e validação da dor.",
    "Análise do ciclo de luto com normalização das fases.",
    "Técnica de cadeira vazia para diálogo com ex-parceiro.",
    "Processamento de culpa e raiva com distribuição equilibrada.",
    "Validação do esgotamento do cuidador e permissão para limites.",
    "Análise da situação de sobrecarga e recursos disponíveis.",
    "Estabelecimento de limites realistas no cuidado.",
    "Busca de apoio externo com cuidadores profissionais.",
    "Resgate de identidade própria além do papel de cuidador.",
]


def escape_sql(text):
    """Escapa aspas simples para SQL"""
    return text.replace("'", "''")


def generate_patient(patient_num):
    """Gera um paciente"""
    patient_id = f"p{patient_num:04d}"
    first_name = random.choice(FIRST_NAMES)
    last_name = random.choice(LAST_NAMES)
    name = f"{first_name} {last_name}"
    clinical_type = random.choice(CLINICAL_TYPES)
    age = random.randint(20, 70)
    notes = f"Paciente de {age} anos. {clinical_type}. Início do tratamento buscando melhora da qualidade de vida."
    created_at = datetime(2023, 1, 1) + timedelta(days=random.randint(0, 365))
    return (
        patient_id,
        name,
        notes,
        created_at.strftime("%Y-%m-%d %H:%M:%S"),
        clinical_type,
    )


def generate_session(session_num, patient_id, session_idx):
    """Gera uma sessão"""
    session_id = f"s{patient_id[1:]}-{session_num:03d}"
    session_date = datetime(2023, 1, 1) + timedelta(
        days=session_idx * 7 + random.randint(-2, 2)
    )
    evol_desc = random.choice(["inicial", "em processo", "de melhora", "consolidação"])
    summary = f"Sessão {session_num}. Paciente em fase {evol_desc} do tratamento. Trabalho terapêutico em andamento."
    created_at = session_date.strftime("%Y-%m-%d %H:%M:%S")
    return (
        session_id,
        patient_id,
        session_date.strftime("%Y-%m-%d %H:%M:%S"),
        summary,
        created_at,
        created_at,
    )


def generate_observation(session_id, obs_num):
    """Gera uma observação"""
    obs_id = f"obs{session_id[1:]}-{obs_num}"
    content = random.choice(OBS_TEMPLATES)
    created_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    return obs_id, session_id, content, created_at, created_at


def generate_intervention(session_id, int_num):
    """Gera uma intervenção"""
    int_id = f"int{session_id[1:]}-{int_num}"
    content = random.choice(INT_TEMPLATES)
    created_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    return int_id, session_id, content, created_at, created_at


def main():
    print("=" * 60)
    print("GERADOR DE MASSA DE DADOS CLÍNICA - ARANDU")
    print("=" * 60)
    print()

    output_file = (
        "internal/infrastructure/repository/sqlite/seeds/seed_massive_clinical_data.sql"
    )

    with open(output_file, "w", encoding="utf-8") as f:
        # Header
        f.write("-- Massa de Dados Clínica Massiva para Arandu\n")
        f.write(f"-- Gerado em: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        f.write(f"-- Configuração: {NUM_PATIENTS} pacientes\n")
        f.write(f"--                {MIN_SESSIONS}-{MAX_SESSIONS} sessões/paciente\n")
        f.write(f"--                {MIN_OBS}-{MAX_OBS} observações/sessão\n")
        f.write(f"--                {MIN_INT}-{MAX_INT} intervenções/sessão\n")
        f.write("\n")
        f.write("DELETE FROM interventions;\n")
        f.write("DELETE FROM observations;\n")
        f.write("DELETE FROM sessions;\n")
        f.write("DELETE FROM patients;\n")
        f.write("DELETE FROM insights;\n")
        f.write("\n")
        f.write("BEGIN TRANSACTION;\n\n")

        total_patients = 0
        total_sessions = 0
        total_obs = 0
        total_int = 0

        # Gera pacientes em batches
        print(f"Gerando {NUM_PATIENTS} pacientes...")
        batch_size = 50
        patient_batch = []

        for i in range(1, NUM_PATIENTS + 1):
            patient_id, name, notes, created_at, clinical_type = generate_patient(i)
            patient_batch.append((patient_id, name, notes, created_at))
            total_patients += 1

            if len(patient_batch) >= batch_size:
                # Insere pacientes
                f.write(
                    "INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES\n"
                )
                values = []
                for p in patient_batch:
                    safe_notes = escape_sql(p[2])
                    values.append(
                        f"('{p[0]}', '{p[1]}', '{safe_notes}', '{p[3]}', '{p[3]}')"
                    )
                f.write(",\n".join(values) + ";\n\n")

                # Gera sessões, observações e intervenções para este batch
                for patient_data in patient_batch:
                    patient_id, _, _, _ = patient_data
                    num_sessions = random.randint(MIN_SESSIONS, MAX_SESSIONS)

                    session_batch = []
                    obs_batch = []
                    int_batch = []

                    for session_idx in range(1, num_sessions + 1):
                        session_id, pid, s_date, summary, created, updated = (
                            generate_session(session_idx, patient_id, session_idx)
                        )
                        safe_summary = escape_sql(summary)
                        session_batch.append(
                            f"('{session_id}', '{pid}', '{s_date}', '{safe_summary}', '{created}', '{updated}')"
                        )
                        total_sessions += 1

                        # Gera observações
                        num_obs = random.randint(MIN_OBS, MAX_OBS)
                        for obs_idx in range(1, num_obs + 1):
                            obs_id, sid, content, o_created, o_updated = (
                                generate_observation(session_id, obs_idx)
                            )
                            safe_content = escape_sql(content)
                            obs_batch.append(
                                f"('{obs_id}', '{sid}', '{safe_content}', '{o_created}', '{o_updated}')"
                            )
                            total_obs += 1

                        # Gera intervenções
                        num_int = random.randint(MIN_INT, MAX_INT)
                        for int_idx in range(1, num_int + 1):
                            int_id, sid, content, i_created, i_updated = (
                                generate_intervention(session_id, int_idx)
                            )
                            safe_content = escape_sql(content)
                            int_batch.append(
                                f"('{int_id}', '{sid}', '{safe_content}', '{i_created}', '{i_updated}')"
                            )
                            total_int += 1

                    # Insere dados deste paciente
                    if session_batch:
                        f.write(
                            "INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES\n"
                        )
                        f.write(
                            ",\n".join(session_batch[:500]) + ";\n\n"
                        )  # Limita batch para não estourar memória
                        if len(session_batch) > 500:
                            f.write(
                                "INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES\n"
                            )
                            f.write(",\n".join(session_batch[500:]) + ";\n\n")

                    if obs_batch:
                        for j in range(0, len(obs_batch), 500):
                            f.write(
                                "INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES\n"
                            )
                            f.write(",\n".join(obs_batch[j : j + 500]) + ";\n\n")

                    if int_batch:
                        for j in range(0, len(int_batch), 500):
                            f.write(
                                "INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES\n"
                            )
                            f.write(",\n".join(int_batch[j : j + 500]) + ";\n\n")

                patient_batch = []

                if i % 100 == 0:
                    print(f"  Processados {i}/{NUM_PATIENTS} pacientes...")
                    print(
                        f"    Total acumulado: {total_sessions} sessões, {total_obs} observações, {total_int} intervenções"
                    )

        # Insere restante
        if patient_batch:
            f.write(
                "INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES\n"
            )
            values = []
            for p in patient_batch:
                safe_notes = escape_sql(p[2])
                values.append(
                    f"('{p[0]}', '{p[1]}', '{safe_notes}', '{p[3]}', '{p[3]}')"
                )
            f.write(",\n".join(values) + ";\n\n")

            for patient_data in patient_batch:
                patient_id, _, _, _ = patient_data
                num_sessions = random.randint(MIN_SESSIONS, MAX_SESSIONS)

                session_batch = []
                obs_batch = []
                int_batch = []

                for session_idx in range(1, num_sessions + 1):
                    session_id, pid, s_date, summary, created, updated = (
                        generate_session(session_idx, patient_id, session_idx)
                    )
                    safe_summary = escape_sql(summary)
                    session_batch.append(
                        f"('{session_id}', '{pid}', '{s_date}', '{safe_summary}', '{created}', '{updated}')"
                    )
                    total_sessions += 1

                    num_obs = random.randint(MIN_OBS, MAX_OBS)
                    for obs_idx in range(1, num_obs + 1):
                        obs_id, sid, content, o_created, o_updated = (
                            generate_observation(session_id, obs_idx)
                        )
                        safe_content = escape_sql(content)
                        obs_batch.append(
                            f"('{obs_id}', '{sid}', '{safe_content}', '{o_created}', '{o_updated}')"
                        )
                        total_obs += 1

                    num_int = random.randint(MIN_INT, MAX_INT)
                    for int_idx in range(1, num_int + 1):
                        int_id, sid, content, i_created, i_updated = (
                            generate_intervention(session_id, int_idx)
                        )
                        safe_content = escape_sql(content)
                        int_batch.append(
                            f"('{int_id}', '{sid}', '{safe_content}', '{i_created}', '{i_updated}')"
                        )
                        total_int += 1

                if session_batch:
                    f.write(
                        "INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES\n"
                    )
                    f.write(",\n".join(session_batch[:500]) + ";\n\n")
                    if len(session_batch) > 500:
                        f.write(
                            "INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES\n"
                        )
                        f.write(",\n".join(session_batch[500:]) + ";\n\n")

                if obs_batch:
                    for j in range(0, len(obs_batch), 500):
                        f.write(
                            "INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES\n"
                        )
                        f.write(",\n".join(obs_batch[j : j + 500]) + ";\n\n")

                if int_batch:
                    for j in range(0, len(int_batch), 500):
                        f.write(
                            "INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES\n"
                        )
                        f.write(",\n".join(int_batch[j : j + 500]) + ";\n\n")

        # Finaliza
        f.write("COMMIT;\n\n")
        f.write("-- ============================================\n")
        f.write("-- ESTATÍSTICAS DO SEED\n")
        f.write("-- ============================================\n")
        f.write(f"-- Total de pacientes: {total_patients}\n")
        f.write(f"-- Total de sessões: {total_sessions}\n")
        f.write(f"-- Total de observações: {total_obs}\n")
        f.write(f"-- Total de intervenções: {total_int}\n")
        f.write(
            f"-- Média de sessões por paciente: {total_sessions // total_patients}\n"
        )
        f.write("-- ============================================\n")

    print()
    print("=" * 60)
    print("RESUMO DA GERAÇÃO")
    print("=" * 60)
    print(f"  Pacientes: {total_patients}")
    print(f"  Sessões: {total_sessions}")
    print(f"  Observações: {total_obs}")
    print(f"  Intervenções: {total_int}")
    print()

    import os

    file_size = os.path.getsize(output_file)
    print(f"✅ Arquivo gerado: {output_file}")
    print(f"📊 Tamanho: {file_size / 1024 / 1024:.2f} MB")
    print()
    print("Para executar:")
    print(f"  sqlite3 arandu.db < {output_file}")


if __name__ == "__main__":
    main()
