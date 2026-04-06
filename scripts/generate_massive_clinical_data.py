#!/usr/bin/env python3
\\"\\"\\"
Gerador Massivo de Dados Clínicos para Arandu
Cria 500+ pacientes, 50.000+ sessões, 100.000+ observações, 50.000+ intervenções
\\"\\"\\"

import uuid
import random
from datetime import datetime, timedelta
from typing import List, Tuple, Dict
import json

# Configurações
CONFIG = {
    'num_patients': 500,
    'min_sessions_per_patient': 100,
    'max_sessions_per_patient': 150,
    'min_observations_per_session': 2,
    'max_observations_per_session': 4,
    'min_interventions_per_session': 1,
    'max_interventions_per_session': 3,
    'start_date': datetime(2023, 1, 1),
    'end_date': datetime(2026, 3, 17),
}

# Nomes brasileiros
FIRST_NAMES_MALE = [
    'João', 'José', 'Antônio', 'Francisco', 'Carlos', 'Paulo', 'Pedro', 'Lucas', 'Luiz', 'Marcos',
    'Luís', 'Gabriel', 'Rafael', 'Daniel', 'Marcelo', 'Bruno', 'Eduardo', 'Felipe', 'Guilherme',
    'Leonardo', 'Rodrigo', 'Fernando', 'Fabio', 'Ricardo', 'Gustavo', 'André', 'Alexandre', 'Marco',
    'Julio', 'Roberto', 'Sérgio', 'Mário', 'César', 'Victor', 'Mateus', 'Davi', 'Samuel', 'Arthur',
    'Heitor', 'Lorenzo', 'Enzo', 'Nicolas', 'Benjamin', 'Henry', 'Noah', 'Murilo', 'Vicente',
    'Otávio', 'Augusto', 'Cauã', 'Kaique', 'Anthony', 'Yuri', 'Breno', 'Thiago', 'Diego',
    'Vitor', 'Renato', 'Igor', 'Douglas', 'Caio', 'Leandro', 'Emerson', 'Wagner', 'Ronaldo',
    'Rogério', 'Maurício', 'Cláudio', 'Jonas', 'Cristiano', 'Adriano', 'Patrik', 'Hugo', 'Alan',
    'Elias', 'Vinícius', 'Jorge', 'William', 'Tiago', 'Danilo', 'Flávio', 'Anderson', 'Nathan',
    'Tomás', 'Ivan', 'Miguel', 'Theo', 'Davi', 'Helena', 'Alice', 'Manuela', 'Valentina',
    'Laura', 'Isabella', 'Heloísa', 'Lara', 'Lívia', 'Lorena', 'Aisha', 'Melissa', 'Sarah',
    'Mirela', 'Mariana', 'Yasmin', 'Isabelly', 'Raissa', 'Julia', 'Sophia', 'Luna', 'Alícia',
    'Cecília', 'Elisa', 'Júlia', 'Isis', 'Antonella', 'Maya', 'Pietra', 'Agatha', 'Nicole',
    'Emanuelly', 'Milena', 'Rebeca', 'Bella', 'Kamilly', 'Stella', 'Catarina', 'Olivia',
    'Ana', 'Beatriz', 'Sofia', 'Gabriela', 'Vitória', 'Clara', 'Fernanda', 'Bianca',
    'Camila', 'Amanda', 'Natália', 'Letícia', 'Luana', 'Débora', 'Cristina', 'Patrícia',
    'Simone', 'Adriana', 'Márcia', 'Silvia', 'Tatiana', 'Monica', 'Eliane', 'Sandra',
    'Rosa', 'Vera', 'Luciana', 'Elaine', 'Daniela', 'Vanessa', 'Jéssica', 'Priscila',
    'Carolina', 'Natalia', 'Tamires', 'Aline', 'Jaqueline', 'Bruna', 'Fernanda', 'Juliana',
    'Karina', 'Sabrina', 'Tatiane', 'Viviane', 'Raquel', 'Tereza', 'Luzia', 'Marlene',
    'Aparecida', 'Francisca', 'Raimunda', 'Antonia', 'Sebastiana', 'Joaquina', 'Benedita'
]

FIRST_NAMES_FEMALE = [
    'Maria', 'Ana', 'Francisca', 'Antônia', 'Adriana', 'Juliana', 'Fernanda', 'Patrícia',
    'Mariana', 'Vanessa', 'Amanda', 'Bruna', 'Beatriz', 'Carolina', 'Daniela', 'Eduarda',
    'Fernanda', 'Gabriela', 'Helena', 'Isabela', 'Julia', 'Karina', 'Larissa', 'Mariana',
    'Natália', 'Olivia', 'Patricia', 'Quezia', 'Raquel', 'Sabrina', 'Tatiana', 'Ursula',
    'Valentina', 'Wanessa', 'Xenia', 'Yasmin', 'Zara', 'Alice', 'Bianca', 'Clara',
    'Diana', 'Elisa', 'Flávia', 'Giovana', 'Heloisa', 'Ingrid', 'Joana', 'Kelly',
    'Lara', 'Manuela', 'Nina', 'Odete', 'Priscila', 'Quiteria', 'Rafaela', 'Sofia',
    'Teresa', 'Umbelina', 'Violeta', 'Walesca', 'Ximena', 'Yara', 'Zilda', 'Aline',
    'Barbara', 'Camila', 'Debora', 'Elaine', 'Fabiana', 'Gisele', 'Helen', 'Iracema',
    'Jéssica', 'Kamila', 'Luana', 'Melissa', 'Nayara', 'Ornela', 'Paloma', 'Quésia',
    'Renata', 'Samara', 'Talita', 'Viviane', 'Wanessa', 'Xuxa', 'Yeda', 'Zelia',
    'Agatha', 'Betina', 'Cintia', 'Dalila', 'Estela', 'Fabíola', 'Geovana', 'Haydée',
    'Ilana', 'Jandira', 'Kellen', 'Lidiane', 'Marlene', 'Nadir', 'Onélia', 'Paula'
]

LAST_NAMES = [
    'Silva', 'Santos', 'Oliveira', 'Souza', 'Pereira', 'Costa', 'Rodrigues', 'Almeida',
    'Nascimento', 'Lima', 'Carvalho', 'Araújo', 'Ferreira', 'Gomes', 'Ribeiro', 'Martins',
    'Barbosa', 'Alves', 'Rocha', 'Cardoso', 'Correia', 'Nunes', 'Mendes', 'Cavalcanti',
    'Dias', 'Teixeira', 'Monteiro', 'Freitas', 'Machado', 'Moreira', 'Andrade', 'Borges',
    'Pinto', 'Vieira', 'Moura', 'Cunha', 'Sales', 'Campos', 'Reis', 'Duarte',
    'Bezerra', 'Moraes', 'Castro', 'Sousa', 'Gonçalves', 'Melo', 'Ramos', 'Barros',
    'Macedo', 'Pinheiro', 'Azevedo', 'Farias', 'Braga', 'Neves', 'Vasconcelos', 'Leite',
    'Batista', 'Peixoto', 'Guimarães', 'Miranda', 'Dantas', 'Tavares', 'Aguiar', 'Amaral',
    'Assis', 'Bittencourt', 'Brito', 'Cabral', 'Caldas', 'Câmara', 'Carneiro', 'Coelho',
    'Coutinho', 'Cruz', 'Cunha', 'Domingues', 'Esteves', 'Figueiredo', 'Franco', 'Freire',
    'Furtado', 'Garcia', 'Henrique', 'Lopes', 'Marques', 'Matos', 'Medeiros', 'Mesquita',
    'Montenegro', 'Neto', 'Nóbrega', 'Paiva', 'Pimentel', 'Porto', 'Queiroz', 'Rangel',
    'Rego', 'Resende', 'Sampaio', 'Siqueira', 'Torres', 'Veloso', 'Viana', 'Xavier',
    'Zaragoza', 'Abreu', 'Albuquerque', 'Almada', 'Alvarenga', 'Alvarez', 'Amorim',
    'Antunes', 'Ávila', 'Azambuja', 'Baptista', 'Barreto', 'Barros', 'Beltrão', 'Bernardes'
]

# Contextos clínicos variados
CLINICAL_CONTEXTS = [
    {
        'tipo': 'Transtorno de Ansiedade Generalizada',
        'sintomas': ['preocupação excessiva', 'tensão muscular', 'insônia', 'irritabilidade', 'fadiga'],
        'abordagem': 'TCC',
        'evolucao': ['piora inicial', 'estabilização', 'melhora gradual', 'alta']
    },
    {
        'tipo': 'Depressão Moderada',
        'sintomas': ['humor deprimido', 'anhedonia', 'falta de energia', 'alterações do sono', 'culpa'],
        'abordagem': 'TCC',
        'evolucao': ['piora', 'mesma', 'melhora lenta', 'remissão parcial']
    },
    {
        'tipo': 'Transtorno do Pânico',
        'sintomas': ['crises de pânico', 'fear de locais públicos', 'sintomas físicos', 'anticipação ansiosa'],
        'abordagem': 'TCC',
        'evolucao': ['frequente', 'espaçada', 'rara', 'remissão']
    },
    {
        'tipo': 'Fobia Social',
        'sintomas': ['evitação social', 'medo de julgamento', 'ansiedade em público', 'isamento'],
        'abordagem': 'TCC',
        'evolucao': ['severa', 'moderada', 'leve', 'controle']
    },
    {
        'tipo': 'Transtorno Obsessivo-Compulsivo',
        'sintomas': ['obsessões', 'compulsões', 'ansiedade de contaminação', 'verificação excessiva'],
        'abordagem': 'TCC',
        'evolucao': ['incapacitante', 'interferente', 'gerenciável', 'minimizado']
    },
    {
        'tipo': 'TEPT (Trauma)',
        'sintomas': ['flashbacks', 'pesadelos', 'hipervigilância', 'evitação', 'humor alterado'],
        'abordagem': 'TCC',
        'evolucao': ['agudo', 'crônico', 'melhora', 'processado']
    },
    {
        'tipo': 'Luto Complicado',
        'sintomas': ['dor intensa', 'identidade perdida', 'culpa', 'isolamento', 'falta de sentido'],
        'abordagem': 'Psicanalítica',
        'evolucao': ['recente', 'persistent', 'elaboração', 'integração']
    },
    {
        'tipo': 'Crise de Meia-Idade',
        'sintomas': ['questionamento existencial', 'insatisfação', 'ansiedade de morte', 'mudanças radicais'],
        'abordagem': 'Psicanalítica',
        'evolucao': ['turbulento', 'exploratório', 'reestruturação', 'integração']
    },
    {
        'tipo': 'Transtorno Bipolar II',
        'sintomas': ['hipomania', 'depressão', 'ciclos', 'impulsividade', 'instabilidade'],
        'abordagem': 'Integrativa',
        'evolucao': ['descontrolado', 'identificado', 'estabilizado', 'manejado']
    },
    {
        'tipo': 'Borderline',
        'sintomas': ['instabilidade emocional', 'medo de abandono', 'impulsividade', 'identidade difusa'],
        'abordagem': 'DBT',
        'evolucao': ['crise', 'instável', 'flutuante', 'regulado']
    },
    {
        'tipo': 'Dependência Emocional',
        'sintomas': ['dificuldade em estar só', 'medo de abandono', 'autoestima baixa', 'relações tóxicas'],
        'abordagem': 'TCC',
        'evolucao': ['severa', 'consciente', 'trabalhada', 'autônoma']
    },
    {
        'tipo': 'Burnout Profissional',
        'sintomas': ['exaustão', 'cinismo', 'ineficácia', 'despersonalização', 'falta de sentido'],
        'abordagem': 'TCC',
        'evolucao': ['agudo', 'reconhecido', 'mudanças', 'recuperação']
    },
    {
        'tipo': 'Adoção e Identidade',
        'sintomas': ['busca por origens', 'identidade fragmentada', 'abandono', 'não pertencimento'],
        'abordagem': 'Psicanalítica',
        'evolucao': ['questionamento', 'exploração', 'integração', 'aceitação']
    },
    {
        'tipo': 'Separacao e Divorcio',
        'sintomas': ['luto do relacionamento', 'medo de solidão', 'raiva', 'questões práticas'],
        'abordagem': 'TCC',
        'evolucao': ['recente', 'conflituoso', 'processamento', 'reorganização']
    },
    {
        'tipo': 'Cuidador Primario Esgotado',
        'sintomas': ['fadiga', 'ressentimento', 'culpa', 'perda de identidade', 'isolamento'],
        'abordagem': 'TCC',
        'evolucao': ['sobrecarregado', 'consciente', 'limites', 'equilíbrio']
    }
]

# Templates de observações por contexto
OBSERVATION_TEMPLATES = {
    'Transtorno de Ansiedade Generalizada': [
        \\"Paciente apresenta tensão muscular visível, especialmente em ombros e mandíbula. Postura rígida, movimentos controlados. Fala rápida com respiração curta.\\",
        \\"Preocupações circulares evidentes: mesmo tema retorna múltiplas vezes. Dificuldade em conter o fluxo de pensamentos catastróficos.\\",
        \\"Sinais de hipervigilância: observa constantemente ambiente, reage a estímulos mínimos. Estado de alerta elevado presente.\\",
        \\"Insônia relatada: dificuldade em desligar mente à noite. Rotina de sono desregulada afeta humor e energia.\\",
        \\"Progresso notável: primeira sessão com postura relaxada. Ombros caídos naturalmente, respiração diafragmática observada.\\",
        \\"Uso efetivo de técnicas de relaxamento: demonstrou exercício de grounding quando ansiedade elevou. Internalização de coping.\\",
        \\"Recaída situacional: evento estressor externo reativou padrão antigo. Porém, recuperação mais rápida que em episódios anteriores.\\",
        \\"Insight emergente: reconhece padrão de antecipar problemas. 'Já percebo que estou catastrofizando'. Consciência metacognitiva.\\"
    ],
    'Depressão Moderada': [
        \\"Psicomotricidade retardada: movimentos lentos, voz baixa, postura encolhida. Contacto visual mínimo, olhar fixo no chão.\\",
\\"Anedonia evidente: nenhuma atividade gera prazer relatado. Descreve dias monótonos, \\\"tudo igual, tudo cinza\\\".\\",
\\"Autocrítica severa: linguagem depreciativa sobre si mesma. \\\"Sou um fracasso\\\", \\\"Não sirvo para nada\\\". Cognições disfuncionais.\\",
        \\"Isolamento social: evita contatos, cancela compromissos. Vínculos afetivos mantidos superficialmente.\\",
        \\"Primeiros sinais de melhora: trouxe assunto externo não solicitado. Interesse diminuto mas presente em algo além do sofrimento.\\",
        \\"Ativação comportamental: retomou atividade física. Relata cansaço mas também sensação de realização. Ciclo virtuoso iniciando.\\",
        \\"Humor labil: oscila entre tristeza profunda e irritabilidade. Períodos de choro alternam com momentos de raiva reprimida.\\",
\\"Cognições mais flexíveis: consegue questionar pensamentos automáticos negativos. \\\"Talvez eu esteja sendo dura demais comigo\\\".\\",
        \\"Reconexão com valores: identificou o que realmente importa. Discussão sobre sentido de vida além do humor deprimido.\\"
    ],
    'Transtorno do Pânico': [
        \\"Crise de pânico relatada na semana: taquicardia, sudorese, medo de morrer. Foi ao pronto-socorro, exames normais.\\",
        \\"Evitação de locais: deixou de frequentar shopping, cinema, transporte público. Espaço vital diminuído progressivamente.\\",
        \\"Hipervigilância corporal: monitora constantemente batimentos, respiração. Interpreta sensações normais como ameaça.\\",
        \\"Antecipação ansiosa: medo de ter medo. Ansiedade secundária sobre possíveis crises. Ciclo de feedback estabelecido.\\",
        \\"Exposição bem-sucedida: enfrentou situação temida. Ansiedade elevou mas diminuiu. Habituação ocorrida.\\",
\\"Reestruturação cognitiva: questiona interpretações catastróficas. \\\"E se for só ansiedade e não infarto?\\\".\\",
        \\"Diário de pensamentos: registros detalhados mostram padrão. Antecedentes, pensamentos, consequências mapeados.\\",
        \\"Medicação estabilizadora: relata melhora significativa com ISRS. Crises menos intensas e frequentes.\\",
        \\"Recuperação de autonomia: retomou atividades evitadas. Confiança progressiva no próprio corpo.\\"
    ],
    'Fobia Social': [
        \\"Ansiedade social severa: contato visual fugidio, voz baixa, postura retrátil. Linguagem corporal de invisibilidade.\\",
        \\"Evitação sistemática: recusa convites, sai cedo de eventos. Comportamentos de segurança minimizam exposição social.\\",
\\"Medo de julgamento: preocupação obsessiva com avaliação alheia. \\\"Vão perceber que sou esquisita\\\".\\",
        \\"Isolamento compensatório: redes sociais substituem interação presencial. Vida social virtual mas não real.\\",
        \\"Exposição gradual: interação planejada bem-sucedida. Ansiedade inicial alta, diminuição natural com tempo.\\",
        \\"Feedback positivo: interação social foi melhor que esperado. Reforço externo contradiz expectativas negativas.\\",
        \\"Assertividade emergente: consegue dizer não sem culpa excessiva. Limites saudáveis estabelecidos.\\",
\\"Autoimagem em transformação: menos comparação social, mais aceitação. \\\"Sou tímida e tudo bem\\\".\\",
        \\"Rede social construída: novos contatos, amizades desenvolvidas. Pertencimento social estabelecido.\\"
    ],
    'Transtorno Obsessivo-Compulsivo': [
        \\"Rituais compulsivos evidentes: repetições de ações, verificações excessivas. Perda de tempo significativa.\\",
        \\"Obsessões intrusivas: pensamentos indesejados invadem constantemente. Esforço cognitivo para neutralizá-los.\\",
        \\"Ansiedade de contaminação: medo intenso de sujeira, germes. Lavagens excessivas, pele ressecada.\\",
        \\"Dano evitado: verificações repetidas de portas, fogão, tomadas. Incapacidade de confiar na memória.\\",
        \\"Prevenção de resposta: resistiu à compulsão. Ansiedade elevada inicialmente, habituação ao longo da semana.\\",
        \\"Exposição intencional: confronto com tema temido. Iniciação de comportamento alvo da ansiedade.\\",
        \\"Insight sobre o ciclo: reconhece que compulsões aliviam temporariamente mas reforçam o TOC.\\",
        \\"Redução de rituos: tempo gasto em compulsões diminuiu. Flexibilidade comportamental aumentada.\\",
        \\"Aceitação da incerteza: desenvolvendo tolerância a dúvidas. \\"Posso viver sem ter certeza absoluta\\".\\"
    ],
    'TEPT (Trauma)': [
        \\"Flashbacks relatados: memória intrusiva do trauma. Dissociação, sensação de reviver evento. Desorientação temporal.\\",
        \\"Pesadelos frequentes: sono perturbado, medo de dormir. Insônia secundária ao trauma.\\",
        \\"Hipervigilância excessiva: estado de alerta constante. Sinais de perigo procurados ativamente.\\",
        \\"Evitação de gatilhos: lugares, pessoas, situações remetem ao trauma. Espaço vital restrito.\\",
        \\"Processamento do trauma: narração coesa do evento. Integrando experiência na história de vida.\\",
        \\"Técnica de grounding: uso de ancoragem no presente. Redução de dissociação durante gatilhos.\\",
\\"Reestruturação de culpa: questionamento de autoacusações. \\\"Não foi minha culpa, eu fiz o que podia\\\".\\",
        \\"Sono restaurado: pesadelos diminuíram. Tecnicas de higiene do sono efetivas.\\",
        \\"PTSD em remissão: sintomas mínimos, funcionamento restaurado. Trauma integrado, não superado.\\"
    ],
    'Luto Complicado': [
        \\"Luto persistente: 18 meses após perda, sofrimento intenso. Identificação simbiótica com falecido.\\",
\\"Melancolia evidente: perda de interesse em vida sem pessoa. \\\"Não existe mais razão para viver\\\".\\",
        \\"Culpa de sobrevivente: questiona por que ficou vivo. Autocrítica sobre coisas ditas/não ditas.\\",
        \\"Isolamento social: afastamento de amigos, família. Dificuldade em ver outros seguindo vida.\\",
        \\"Rituais de despedida: primeiro aniversário sem pessoa. Celebração diferente, mas presente.\\",
        \\"Reconexão com vida: pequenos prazeres emergindo. Jardim, caminhadas, momentos de paz.\\",
        \\"Integração da perda: falecido lembrado sem dor esmagadora. Memórias trazem afeto misturado.\\",
\\"Nova identidade: \\\"Sou viúva mas também sou eu mesma\\\". Autonomia reconstruída.\\",
        \\"Projeto futuro: caderno de receitas para netos. Legado familiar preservado, futuro construído.\\"
    ],
    'Crise de Meia-Idade': [
\\"Questionamento existencial: \\\"É isso? É só isso?\\\". Insatisfação com conquistas aparentes.\\",
        \\"Medo de morte: ansiedade de finitude. Tempo limitado percebido, pressão para mudanças.\\",
        \\"Impulsividade compensatória: mudanças radicais consideradas. Término, mudança de carreira, aventuras.\\",
        \\"Comparação social: avaliação de vida em relação a pares. Sentimento de atraso ou escolhas erradas.\\",
        \\"Exploração de valores: identificação do que realmente importa. Desapego de sucessos externos.\\",
        \\"Reavaliação de relacionamentos: vínculos autênticos versus sociais. Qualidade sobre quantidade.\\",
        \\"Transição profissional: mudança de carreira planejada. Propósito sobre status.\\",
\\"Aceitação da finitude: desenvolvimento de sabedoria. \\\"Não dá para agradar a todos, escolho o que importa\\\".\\",
        \\"Integridade alcançada: alinhamento entre valores e ações. Vida autêntica emergindo.\\"
    ],
    'Transtorno Bipolar II': [
        \\"Episódio hipomaníaco: humor elevado, energia excessiva, redução do sono. Grandiosidade presente.\\",
        \\"Fala pressada: ideias saltitantes, dificuldade em seguir linha de raciocínio. Produtividade aumentada.\\",
        \\"Impulsividade: gastos excessivos, projetos múltiplos. Julgamento prejudicado na euforia.\\",
        \\"Queda depressiva: humor deprimido após euforia. Culpa pelos excessos, fadiga intensa.\\",
        \\"Estabilização: humor equilibrado, sono regular. Rotina estruturada mantida.\\",
        \\"Diário de humor: padrões identificados. Previsibilidade dos ciclos aumentando.\\",
        \\"Psicoeducação: compreensão do transtorno. Aceitação do diagnóstico, não mais resistência.\\",
        \\"Medicamento estabilizador: lítio iniciado. Efeitos colaterais gerenciáveis, benefícios claros.\\",
        \\"Estabilidade prolongada: primeiro trimestre sem episódio severo. Qualidade de vida restaurada.\\"
    ],
    'Borderline': [
        \\"Instabilidade emocional: oscilações rápidas e intensas. Regulação afetiva comprometida.\\",
        \\"Medo de abandono intenso: reações desproporcionais a percepção de rejeição. Vínculos instáveis.\\",
        \\"Impulsividade autodestrutiva: comportamentos de risco, automutilação histórica. Coping disfuncional.\\",
\\"Identidade difusa: \\\"Não sei quem sou\\\". Interesses, valores, relacionamentos em constante mudança.\\",
        \\"Crise de choro: desregulação emocional. Uso de técnicas de distress tolerance.\\",
        \\"Mindfulness praticado: momentos de consciência plena. Atenção ao presente, não reatividade.\\",
        \\"Habilidades de interpessoal: validação própria, não apenas externa. Relacionamento mais estável.\\",
        \\"Tolerância à frustração: pequenos contratempos gerenciados sem crise. Resiliência emergente.\\",
        \\"Identidade consolidada: valores claros, metas definidas. \\"Sei quem sou e para onde vou\\".\\"
    ],
    'Dependência Emocional': [
        \\"Dificuldade extrema em estar só: ansiedade de separação. Telefonemas constantes ao parceiro.\\",
\\"Autoestima baseada no outro: \\\"Se ele me ama, sou válida\\\". Valorização condicional.\\",
        \\"Medo de abandono: terror de ser deixada. Submissão para manter relacionamento.\\",
        \\"Relações tóxicas repetidas: mesmo padrão com parceiros diferentes. Vício emocional.\\",
        \\"Primeira vez sozinha: viagem sozinha, resistida mas realizada. Autonomia construída.\\",
        \\"Autoestima própria: identificação de qualidades independentes do relacionamento.\\",
\\"Limites saudáveis: consegue dizer não. Não mais \\\"pessoa agradável\\\" compulsiva.\\",
        \\"Relacionamento consigo mesma: autocuidado, prazeres individuais. Individualidade cultivada.\\",
        \\"Parceria saudável: relacionamento atual baseado em escolha, não necessidade. Liberdade no amor.\\"
    ],
    'Burnout Profissional': [
        \\"Exaustão severa: energia depletada, esgotamento físico e emocional. Dificuldade em recuperar.\\",
        \\"Cinismo profissional: atitudes negativas sobre trabalho, colegas, pacientes. Desumanização.\\",
        \\"Ineficácia percebida: sensação de não fazer diferença. Competência questionada.\\",
        \\"Despersonalização: tratamento mecânico de tarefas. Desconexão do sentido do trabalho.\\",
        \\"Licença médica: pausa forçada, reconhecimento do problema. Espaço para recuperação.\\",
        \\"Reavaliação de prioridades: o que realmente importa na vida? Trabalho não é tudo.\\",
        \\"Mudanças estruturais: limites no trabalho estabelecidos. Não mais disponível 24/7.\\",
        \\"Novo projeto profissional: redução de cargo para qualidade de vida. Propósito realinhado.\\",
        \\"Recuperação completa: energia restaurada, sentido recuperado. Trabalho saudável possível.\\"
    ],
    'Adoção e Identidade': [
        \\"Busca por origens: necessidade de saber história. Identidade fragmentada sem raízes.\\",
\\"Não pertencimento: \\\"Nunca me senti realmente filho\\\". Diferença física, temperamental.\\",
        \\"Culpa de buscar: medo de magoar pais adotivos. Lealdade conflitante.\\",
        \\"Raiva reprimida: por que fui abandonado? Sentimento de descartabilidade.\\",
        \\"Contato com família biológica: encontro realizado. Misto de emoções, mas integrador.\\",
\\"Narrativa coesa: história de vida reconstruída. \\\"Sou filho de dois mundos\\\".\\",
        \\"Aceitação das duas famílias: amor multiplicado, não dividido. Gratidão e curiosidade.\\",
\\"Identidade integrada: \\\"Sou adotado e sou eu mesmo\\\". Ambas verdades coexistem.\\",
        \\"Projeto de origem: árvore genealógica completa. Pertencimento aos dois lados.\\"
    ],
    'Separacao e Divorcio': [
        \\"Luto do relacionamento: término recente, dor intensa. Vida construída juntos desmontada.\\",
        \\"Medo de solidão: medo de estar só, reconstruir vida. Ansiedade sobre futuro incerto.\\",
        \\"Raiva e traição: sentimento de injustiça. Histórias conflitantes sobre o término.\\",
        \\"Questões práticas: divisão de bens, guarda de filhos. Conflitos além do emocional.\\",
        \\"Processamento do fim: narrativa coesa da história. Aceitação das partes de cada um.\\",
        \\"Novo apartamento: espaço próprio construído. Individualidade territorial estabelecida.\\",
        \\"Rede de apoio: amigos, família próxima. Solidão escolhida versus imposta.\\",
        \\"Novos interesses: hobbies, atividades individuais. Vida própria reconstruída.\\",
        \\"Ciclo fechado: término processado, aprendizados integrados. Aberto ao novo.\\"
    ],
    'Cuidador Primario Esgotado': [
        \\"Fadiga crônica: anos cuidando de familiar doente. Energia física e emocional esgotada.\\",
\\"Ressentimento: \\\"Minha vida parou para cuidar dele\\\". Culpa pelo ressentimento.\\",
\\"Perda de identidade: \\\"Sou só cuidadora, esqueci quem sou\\\". Individualidade anulada.\\",
        \\"Isolamento social: amigos afastados, vida social inexistente. Solidão no cuidado.\\",
        \\"Pedido de ajuda: aceitou cuidador profissional 2x/semana. Primeira pausa em anos.\\",
        \\"Reconexão consigo: saídas sozinha, hobbies retomados. Individualidade resgatada.\\",
        \\"Limites saudáveis: não mais culpa por não estar presente 24/7. Autocuidado permitido.\\",
        \\"Grupo de apoio: encontro com outros cuidadores. Normalização da experiência.\\",
        \\"Equilíbrio alcançado: cuidado sem sacrifício total. Qualidade sobre quantidade de presença.\\"
    ]
}

# Templates de intervenções por contexto
INTERVENTION_TEMPLATES = {
    'Transtorno de Ansiedade Generalizada': [
        \\"Psicoeducação sobre mecanismo da ansiedade: explicação do ciclo de feedback entre pensamentos catastróficos, tensão corporal e comportamentos de segurança.\\",
        \\"Técnica de respiração diafragmática 4-7-8: demonstração em sessão, prática supervisionada. Prescrição de exercícios diários.\\",
        \\"Exposição gradual in vivo: construção de hierarquia de ansiedade, planejamento de experimento comportamental para próxima semana.\\",
        \\"Reestruturação cognitiva: identificação de pensamentos automáticos, questionamento socrático de evidências, geração de alternativas mais realistas.\\",
        \\"Mindfulness de atenção plena: exercício de consciência do momento presente, desidentificação dos pensamentos.\\",
        \\"Relaxamento progressivo de Jacobson: técnica de tensão e relaxamento muscular sistemático. Gravação fornecida.\\",
        \\"Técnica de postponement da preocupação: agendamento de horário específico para preocupações, redução da frequência intrusiva.\\",
        \\"Exposição interoceptiva: exercícios que induzem sensações físicas similares às da ansiedade, habituação às sensações.\\",
        \\"Comunicação assertiva: treinamento de expressão de necessidades, limites e desejos de forma respeitosa mas firme.\\",
        \\"Plano de manutenção: revisão de estratégias aprendidas, identificação de sinais de alerta, protocolo de retorno.\\"
    ],
    'Depressão Moderada': [
        \\"Psicoeducação sobre depressão: normalização da condição, explicação dos sintomas físicos e psicológicos, prognóstico positivo com tratamento.\\",
        \\"Ativação comportamental: programação de atividades prazerosas e de realização, comportamento oposto às vontades depressivas.\\",
        \\"Registro de pensamentos: diário de situações, pensamentos automáticos, emoções e comportamentos. Identificação de padrões.\\",
        \\"Reestruturação cognitiva: questionamento de cognições disfuncionais, busca de evidências contrárias, reformulação de crenças.\\",
        \\"Técnica de reestruturação em duas colunas: evidências a favor e contra pensamentos negativos. Análise crítica e equilibrada.\\",
        \\"Comportamento oposto às vontades: ação apesar da falta de vontade, confiando que humor seguirá comportamento e não vice-versa.\\",
        \\"Atividades de prazer e domínio: identificação de atividades que geram prazer ou sensação de competência. Programação intencional.\\",
        \\"Terapia comportamental de ativação: análise funcional da depressão, intervenções nos gatilhos ambientais e nas consequências.\\",
        \\"Trabalho com ruminação: técnicas de distanciamento cognitivo, mindfulness, redirecionamento da atenção para o presente.\\",
        \\"Prevenção de recaída: identificação de sinais de alerta precoce, estratégias de coping, rede de apoio, plano de ação.\\"
    ],
    'Transtorno do Pânico': [
        \\"Explicação do ciclo do pânico: interpretação catastrófica de sensações corporais normais -> ansiedade -> sensações intensificadas -> ciclo de feedback.\\",
        \\"Psicoeducação sobre sintomas físicos: normalização de taquicardia, sudorese, falta de ar como reações de luta-fuga benignas.\\",
        \\"Exposição interoceptiva em sessão: hiperventilação controlada, provocação intencional de sintomas para habituação e desastramento.\\",
        \\"Exposição in vivo gradual: construção de hierarquia, confronto sistemático com situações evitadas, sessões de exposição acompanhadas.\\",
\\"Reestruturação cognitiva do medo do medo: \\\"O que mais temo?\\\" \\\"E se acontecesse?\\\" \\\"Como lidaria?\\\" Desastramento do pior cenário.\\",
        \\"Técnica de respiração controlada: respiração diafragmática lenta para prevenir hiperventilação e reduzir sintomas físicos.\\",
        \\"Diário de auto-observação: registro de crises, antecedentes, pensamentos, sensações, comportamentos e consequências. Análise de padrões.\\",
        \\"Prevenção de resposta: resistência às seguranças (anxiolíticos, escape), enfrentamento da ansiedade até habituação.\\",
\\"Técnica de acatastrofização: \\\"Qual a probabilidade real?\\\" \\\"Como me senti da última vez?\\\" \\\"Quais são minhas alternativas?\\\"\\",
        \\"Plano de manutenção: revisão de conquistas, sinais de alerta, estratégias de enfrentamento, protocolo de retorno.\\"
    ],
    'Fobia Social': [
        \\"Psicoeducação sobre ciclo da ansiedade social: expectativa negativa -> ansiedade -> comportamento de segurança -> resultado pobre -> confirmação da expectativa.\\",
        \\"Exposição gradual social: hierarquia de situações temidas, desde cumprimentar conhecido até falar em público.\\",
        \\"Role-play de interação social: ensaio de situações temidas em sessão, feedback comportamental, repetição até confiança.\\",
\\"Reestruturação cognitiva do medo de julgamento: \\\"Como sei o que pensam?\\\" \\\"E se fossem gentis?\\\" \\\"Qual o custo de tentar?\\\"\\",
        \\"Experimento de atenção: foco externo em vez de introspecção excessiva. Observação do ambiente, não de si mesmo.\\",
        \\"Comportamentos de segurança: identificação e eliminação gradual de comportamentos que mantêm ansiedade (evitar olhar, falar baixo).\\",
        \\"Técnica de exposição a rejeição: pedidos intencionais para ser recusado, desmistificação do não como catástrofe.\\",
        \\"Assertividade social: expressão de opiniões, pedidos, limites de forma respeitosa. Treinamento de comunicação.\\",
        \\"Mindfulness social: atenção plena durante interações, aceitação da ansiedade sem luta, foco na conexão.\\",
        \\"Exposição máxima: evento social planejado e enfrentado. Celebração do sucesso, processamento da experiência.\\"
    ],
    'Transtorno Obsessivo-Compulsivo': [
        \\"Psicoeducação sobre TOC: explicação do ciclo obsessão-ansiedade-compulsão-alívio temporário-reforço do ciclo.\\",
        \\"Prevenção de resposta: resistência à realização da compulsão, tolerância à ansiedade até habituação natural.\\",
        \\"Exposição e prevenção de resposta (ERP): confronto gradual com temas temidos sem realizar rituais de alívio.\\",
        \\"Exposição imaginária: narração detalhada de cenas temidas escritas, gravação em áudio, escuta repetida.\\",
\\"Reestruturação cognitiva da importância dos pensamentos: \\\"Pensar não é fazer\\\", \\\"Todos têm pensamentos intrusivos\\\".\\",
        \\"Técnica de aceitação da incerteza: desenvolvimento de tolerância à dúvida sem necessidade de certeza absoluta.\\",
        \\"Deliberação imposível: adiamento intencional da decisão, treinamento em lidar com incerteza prolongada.\\",
        \\"Redução de verificações: limitação de número de vezes, gradualmente diminuído. Aceitação do risco residual.\\",
        \\"Técnica de normalização: observação de comportamentos de outras pessoas, comportamento similar em intensidade normal.\\",
        \\"Plano de manutenção: ERP doméstico, sinais de recaída, estratégias de enfrentamento, sessões de reforço.\\"
    ],
    'TEPT (Trauma)': [
        \\"Psicoeducação sobre TEPT: explicação neurobiológica, normalização dos sintomas, prognóstico de recuperação.\\",
        \\"Técnica de grounding: ancoragem no presente usando sentidos (5-4-3-2-1), redução de flashbacks e dissociação.\\",
        \\"Estabilização e segurança: construção de recursos internos, janela de tolerância, técnicas de autorregulação.\\",
        \\"Processamento do trauma: narração da história de forma segura, integração da experiência na identidade.\\",
        \\"Técnica de cadeira vazia: diálogo simbólico com figura do trauma ou aspectos internos, resolução de assuntos pendentes.\\",
\\"Reestruturação cognitiva da culpa: questionamento de responsabilidade, \\\"O que você diria a um amigo nessa situação?\\\"\\",
        \\"Exposição imaginária controlada: revisitação imaginária do trauma com recursos terapêuticos, reprocessamento.\\",
        \\"Integração da história de vida: colocação do trauma no contexto da trajetória, sentido do sofrimento.\\",
        \\"Técnica de reescrita da narrativa: história de vitória em vez de vítima, resiliência, crescimento pós-traumático.\\",
        \\"Plano de futuro: metas, sonhos, projeto de vida. Trauma é parte da história, não define todo futuro.\\"
    ],
    'Luto Complicado': [
        \\"Psicoeducação sobre luto normal versus complicado: diferenciação, normalização da dor, prognóstico de elaboração.\\",
        \\"Escuta empática não-diretiva: permissão para a dor, presença sem pressa, validação do sofrimento.\\",
        \\"Técnica de cadeira vazia: diálogo com falecido, despedida simbólica, expressão de não-ditos.\\",
        \\"Rituais de despedida: criação de cerimônias pessoais para marcar transições, aniversários, datas especiais.\\",
        \\"Ativação comportamental: retorno gradual a atividades, reconstrução de prazeres e sentido na vida sem a pessoa.\\",
        \\"Trabalho com culpa e arrependimentos: reestruturação cognitiva, perdão a si mesmo, aceitação da imperfeição.\\",
\\"Reconstrução de identidade: \\\"Quem sou eu agora?\\\" Exploração de papéis, interesses, valores independentes.\\",
        \\"Integração da perda: lembranças trazem afeto sem dor esmagadora, falecido internalizado de forma saudável.\\",
        \\"Projeto de legado: continuidade dos valores, projetos em homenagem, perpetuação da memória.\\",
        \\"Fechamento do ciclo: revisão da trajetória de luto, reconhecimento da elaboração, possibilidade de término.\\"
    ],
    'Crise de Meia-Idade': [
        \\"Escuta do questionamento existencial: validação da busca por sentido, normalização da crise como oportunidade.\\",
        \\"Exercício de valores: identificação do que realmente importa, comparação com valores vividos, alinhamento.\\",
        \\"Análise de transição de vida: avaliação de conquistas, insatisfações, metas não realizadas, novas possibilidades.\\",
        \\"Reavaliação de relacionamentos: quais vínculos alimentam, quais esgotam? Decisões sobre investimento emocional.\\",
\\"Exploração de novas identidades: \\\"Quem eu posso ser além do que sou?\\\" Experimentação segura de novos papéis.\\",
        \\"Trabalho com medo da morte: normalização da finitude, urgência de viver autenticamente, legado desejado.\\",
        \\"Reestruturação de sucesso: definição própria versus social, qualidade versus status, significado versus aparência.\\",
        \\"Planejamento de mudanças: transição profissional, geográfica, de estilo de vida. Análise de custo-benefício.\\",
        \\"Integração de sombra: aceitação de aspectos reprimidos, potencialidades negligenciadas, totalidade do self.\\",
        \\"Criação de legado: contribuições desejadas, impacto pretendido, sentido da existência além do individual.\\"
    ],
    'Transtorno Bipolar II': [
        \\"Psicoeducação sobre bipolaridade: explicação de ciclos, gatilhos, manejo. Aceitação do diagnóstico.\\",
        \\"Monitoramento de humor: diário diário com escala 0-10, identificação de padrões e sinais de alerta.\\",
        \\"Higiene do sono: regularização de horários, higiene do sono, sono como estabilizador de humor.\\",
        \\"Estabilização de rotina: horários regulares de refeições, atividades, descanso. Previsibilidade organizacional.\\",
        \\"Identificação de gatilhos: estresse, sono irregular, substâncias, épocas do ano. Prevenção de episódios.\\",
        \\"Intervenção precoce: plano de ação para sinais de mania ou depressão leve. Contenção antes da crise.\\",
        \\"Técnica de contenção na euforia: limites externos quando insight comprometido, contato com familiares.\\",
        \\"Ativação comportamental na depressão: comportamento oposto às vontades, ação precede motivação.\\",
        \\"Psicoeducação familiar: explicação para familiares, sinais de alerta, como ajudar, limites saudáveis.\\",
        \\"Plano de manutenção: revisão de estratégias, adesão ao tratamento, contato com equipe, prevenção de recaídas.\\"
    ],
    'Borderline': [
        \\"Psicoeducação sobre borderline: validação da intensidade emocional, normalização da instabilidade, prognóstico de regulação.\\",
        \\"Treinamento de mindfulness: habilidades de consciência plena, observação sem julgamento, atenção ao momento presente.\\",
        \\"Tolerância à angústia: técnicas de distracao, relaxamento, self-soothing. Sobrevivência à crise sem piorar.\\",
        \\"Regulação emocional: identificação e nomeação de emoções, entendimento de função, modulação de intensidade.\\",
        \\"Eficácia interpessoal: solicitação de mudanças, recusa de pedidos, resolução de conflitos, validação própria e alheia.\\",
        \\"Análise de comportamento alvo: identificação de gatilhos, consequências, função. Mudança de contingências.\\",
        \\"Técnica de aceitação radical: aceitação da realidade como é, não como deveria ser. Redução do sofrimento.\\",
        \\"Validação emocional: reconhecimento de sentimentos como compreensíveis dados as circunstâncias. Não aprovação, mas entendimento.\\",
        \\"Interrupção de ciclo de crise: identificação de cadeia de eventos, pontos de intervenção, estratégias de desaceleração.\\",
        \\"Construção de vida digna de ser vivida: valores, metas, relacionamentos saudáveis. Significado além da sobrevivência.\\"
    ],
    'Dependência Emocional': [
        \\"Psicoeducação sobre dependência: explicação do ciclo de medo de abandono, apego ansioso, comportamentos de aproximação excessiva.\\",
        \\"Análise de história de apego: padrões relacionais iniciais, internalização de modelos, repetição compulsiva.\\",
        \\"Técnica de tolerância à solidão: exposição gradual a estar consigo mesmo, desenvolvimento de companhia interna.\\",
        \\"Reconstrução de autoestima: identificação de qualidades próprias, independência da validação externa.\\",
        \\"Estabelecimento de limites: prática de dizer não, espaço individual, autonomia no relacionamento.\\",
        \\"Atividades individuais: desenvolvimento de hobbies, interesses, amizades independentes do parceiro.\\",
        \\"Técnica de postergamento do contato: resistir à urgência de ligar/mensagem, tolerar a ansiedade de separação.\\",
        \\"Análise de relacionamentos passados: identificação de padrões, escolha de parceiros, repetição de dinâmicas.\\",
        \\"Comunicação de necessidades: expressão de inseguranças de forma madura, pedido de apoio sem demanda.\\",
        \\"Integração de autonomia: \\"Posso estar em relacionamento e ser completo\\". Amor interdependente, não codependente.\\"
    ],
    'Burnout Profissional': [
        \\"Avaliação de burnout: exaustão, cinismo, ineficácia. Identificação de fatores de risco individuais e organizacionais.\\",
        \\"Psicoeducação sobre recuperação: necessidade de pausa, reversibilidade do processo, prognóstico positivo.\\",
        \\"Análise de demandas e recursos: balanço entre o que é exigido e o que é disponível. Desequilíbrio identificado.\\",
        \\"Estabelecimento de limites: não mais disponibilidade ilimitada, horários definidos, recusa de demandas excessivas.\\",
        \\"Reavaliação de valores profissionais: sentido do trabalho, propósito, contribuição. Reconexão com o porquê.\\",
        \\"Técnica de desconexão: transição entre trabalho e vida pessoal, ritual de desligamento, presente no não-trabalho.\\",
        \\"Atividades de recuperação: lazer, descanso, exercício, socialização. Recarga de energia deliberada.\\",
        \\"Comunicação assertiva no trabalho: expressão de limites, negociação de demandas, pedido de apoio.\\",
        \\"Mudança estrutural se necessário: rediscussão de cargo, mudança de área, saída se insustentável.\\",
        \\"Plano de prevenção: sinais de alerta, estratégias de coping, equilíbrio vida-trabalho sustentável.\\"
    ],
    'Adoção e Identidade': [
        \\"Escuta da história de adoção: validação da experiência, permissão para sentimentos contraditórios, sem julgamento.\\",
        \\"Normalização do não pertencimento: compreensível dada a história, não patológico, parte da experiência.\\",
        \\"Exploração de busca por origens: pros e contras, preparo para possíveis resultados, apoio no processo.\\",
        \\"Técnica de diálogo com duas famílias: integração do amor adotivo e da curiosidade biológica. Ambos válidos.\\",
        \\"Reconstrução de narrativa de vida: história coesa incluindo adoção, construção de sentido, identidade integrada.\\",
        \\"Trabalho com abandono: ressignificação, não pessoal, circunstâncias dos pais biológicos, valor próprio.\\",
        \\"Visita ou contato com origens: preparação emocional, apoio durante, processamento depois.\\",
\\"Integração de identidades: \\\"Filho de ambos\\\". Aceitação de complexidade, pertencimento múltiplo.\\",
        \\"Construção de árvore genealógica: inclusão de ambas famílias, reconhecimento de todas as raízes.\\",
        \\"Sentido do pertencimento: escolha de quem é família, vínculos afetivos, comunidades de identidade.\\"
    ],
    'Separacao e Divorcio': [
\\"Escuta da história de término: validação da dor, permissão para raiva e tristeza, sem pressa para \\\"superar\\\".\\",
        \\"Análise do ciclo de luto: negação, raiva, barganha, depressão, aceitação. Normalização das fases.\\",
        \\"Técnica de cadeira vazia com ex-parceiro: diálogo simbólico, despedida, expressão de não-ditos.\\",
        \\"Processamento de culpa e raiva: reestruturação cognitiva, distribuição equilibrada de responsabilidade.\\",
        \\"Apoio nas questões práticas: organização de prioridades, decisões racionais apesar da emoção.\\",
        \\"Construção de nova vida: identificação de valores, metas individuais, projeto de futuro solo.\\",
        \\"Reconstrução de rede social: fortalecimento de amizades, família, novos contatos. Combate à solidão.\\",
        \\"Desenvolvimento de autonomia: tarefas antes divididas, independência prática e emocional.\\",
        \\"Análise de padrões relacionais: escolhas passadas, aprendizados, sabedoria para futuros relacionamentos.\\",
        \\"Integração da experiência: história completa do relacionamento, bons momentos reconhecidos, lições aprendidas.\\"
    ],
    'Cuidador Primario Esgotado': [
        \\"Validação do esgotamento: reconhecimento da sobrecarga, permissão para limites, não é egoísmo.\\",
        \\"Análise da situação: demandas da pessoa cuidada, recursos disponíveis, equilíbrio atual.\\",
        \\"Psicoeducação sobre sobrecarga: normal da situação, não fraqueza pessoal, prevenção de saúde.\\",
        \\"Estabelecimento de limites: não é possível fazer tudo, escolha de prioridades, delegação possível.\\",
        \\"Busca de apoio externo: cuidadores profissionais, grupos de apoio, familiares, serviços comunitários.\\",
        \\"Resgate de identidade própria: interesses individuais, vida além do cuidado, quem é além disso.\\",
        \\"Técnica de respiradores: pausas durante o dia, momentos de cuidado consigo mesmo.\\",
        \\"Comunicação com familiares: divisão de responsabilidades, expressão de necessidades, pedido de ajuda.\\",
        \\"Planejamento de cuidado sustentável: rotina viável, alternância de cuidadores, prevenção de crises.\\",
        \\"Qualidade versus quantidade: presença significativa em vez de exaustiva. Cuidado que alimenta ambos.\\"
    ]
}


def generate_patient(patient_num: int) -> Tuple[str, str, str, str]:
    \\"\\"\\"Gera dados de um paciente\\"\\"\\"
    patient_id = f\\"p{patient_num:04d}\\"
    
    # Decide gênero
    is_male = random.choice([True, False])
    if is_male:
        first_name = random.choice(FIRST_NAMES_MALE)
    else:
        first_name = random.choice(FIRST_NAMES_FEMALE)
    
    last_name = random.choice(LAST_NAMES)
    name = f\\"{first_name} {last_name}\\"
    
    # Contexto clínico
    context = random.choice(CLINICAL_CONTEXTS)
    
    # Idade baseada no contexto
    if context['tipo'] in ['Adoção e Identidade']:
        age = random.randint(18, 35)
    elif context['tipo'] in ['Adolescente com Ansiedade Social']:
        age = random.randint(14, 17)
    elif context['tipo'] in ['Luto Complicado', 'Crise de Meia-Idade', 'Cuidador Primario Esgotado']:
        age = random.randint(45, 70)
    elif context['tipo'] in ['Transtorno Bipolar II']:
        age = random.randint(20, 40)
    else:
        age = random.randint(25, 55)
    
    # Notas clínicas
    notes = f\\"Paciente de {age} anos. {context['tipo']}. Queixas principais: {', '.join(random.sample(context['sintomas'], 3))}. Início do tratamento: busca por {random.choice(['melhora da qualidade de vida', 'redução de sintomas', 'autoconhecimento', 'apoio em momento difícil'])}. Abordagem: {context['abordagem']}.\\"
    
    # Data de criação (aleatória entre 2023 e 2024)
    created_at = CONFIG['start_date'] + timedelta(days=random.randint(0, 365))
    
    return patient_id, name, notes, created_at.strftime('%Y-%m-%d %H:%M:%S')


def generate_session(session_num: int, patient_id: str, patient_context: Dict, session_idx: int, total_sessions: int) -> Tuple[str, str, str, str, str]:
    \\"\\"\\"Gera dados de uma sessão\\"\\"\\"
    session_id = f\\"s{patient_id}-{session_num:03d}\\"
    
    # Data distribuída ao longo do tratamento
    start = CONFIG['start_date']
    end = min(start + timedelta(days=session_idx * 7), CONFIG['end_date'])
    session_date = start + timedelta(days=session_idx * 7 + random.randint(-2, 2))
    
    # Determina evolução baseada na posição
    progress = session_idx / total_sessions
    if progress < 0.2:
        evol_stage = 0  # Início
    elif progress < 0.4:
        evol_stage = 1  # Processo
    elif progress < 0.7:
        evol_stage = 2  # Melhora
    else:
        evol_stage = 3  # Final
    
    # Gera resumo
    evol_desc = patient_context['evolucao'][evol_stage]
    summary_templates = [
        f\\"Sessão {session_idx + 1}. Paciente em fase {evol_desc} do tratamento. {random.choice(['Abordagem de temas centrais', 'Exploração de dinâmicas relacionais', 'Trabalho com padrões disfuncionais', 'Processamento de emoções difíceis'])}. {random.choice(['Resistência inicial presente', 'Participação ativa', 'Insight emergente', 'Trabalho terapêutico intenso'])}.\\",
        f\\"Sessão {session_idx + 1}. {random.choice(['Avanços significativos', 'Estagnação momentânea', 'Crise superada', 'Momento de transição'])} no processo. {random.choice(['Técnicas específicas aplicadas', 'Escuta empática predominante', 'Confronto suave realizado', 'Reforço de progressos'])}. Estado emocional: {evol_desc}.\\",
        f\\"Sessão {session_idx + 1}. {random.choice(['Tema novo emergiu', 'Retorno a questões antigas', 'Consolidação de aprendizados', 'Preparação para encerramento'])}. {random.choice(['Trabalho emocional intenso', 'Reflexão cognitiva predominante', 'Integração de experiências', 'Planejamento de mudanças'])}. Evolução: {evol_desc}.\\",
    ]
    summary = random.choice(summary_templates)
    
    created_at = session_date.strftime('%Y-%m-%d %H:%M:%S')
    updated_at = created_at
    
    return session_id, patient_id, session_date.strftime('%Y-%m-%d %H:%M:%S'), summary, created_at, updated_at


def generate_observation(session_id: str, patient_context: Dict, obs_num: int) -> Tuple[str, str, str, str]:
    \\"\\"\\"Gera uma observação clínica\\"\\"\\"
    obs_id = f\\"obs{session_id}-{obs_num}\\"
    
    # Escolhe template apropriado para o contexto
    templates = OBSERVATION_TEMPLATES.get(patient_context['tipo'], OBSERVATION_TEMPLATES['Transtorno de Ansiedade Generalizada'])
    content = random.choice(templates)
    
    # Timestamp da observação (5 minutos após início da sessão)
    created_at = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    updated_at = created_at
    
    return obs_id, session_id, content, created_at, updated_at


def generate_intervention(session_id: str, patient_context: Dict, int_num: int) -> Tuple[str, str, str, str]:
    \\"\\"\\"Gera uma intervenção terapêutica\\"\\"\\"
    int_id = f\\"int{session_id}-{int_num}\\"
    
    # Escolhe template apropriado para o contexto
    templates = INTERVENTION_TEMPLATES.get(patient_context['tipo'], INTERVENTION_TEMPLATES['Transtorno de Ansiedade Generalizada'])
    content = random.choice(templates)
    
    # Timestamp da intervenção (10 minutos após início da sessão)
    created_at = datetime.now().strftime('%Y-%m-%d %H:%M:%00')
    updated_at = created_at
    
    return int_id, session_id, content, created_at, updated_at


def generate_sql():
    \\"\\"\\"Gera o script SQL completo\\"\\"\\"
    
    sql_lines = [
        \\"-- Massa de Dados Clínica Massiva para Arandu\\",
        \\"-- Gerado em: {}\\".format(datetime.now().strftime('%Y-%m-%d %H:%M:%S')),
        \\"-- Configuração: {} pacientes, {}-{} sessões/paciente, {}-{} observações/sessão, {}-{} intervenções/sessão\\".format(
            CONFIG['num_patients'],
            CONFIG['min_sessions_per_patient'],
            CONFIG['max_sessions_per_patient'],
            CONFIG['min_observations_per_session'],
            CONFIG['max_observations_per_session'],
            CONFIG['min_interventions_per_session'],
            CONFIG['max_interventions_per_session']
        ),
        \\"\\",
        \\"-- Limpar dados existentes\\",
        \\"DELETE FROM interventions;\\",
        \\"DELETE FROM observations;\\",
        \\"DELETE FROM sessions;\\",
        \\"DELETE FROM patients;\\",
        \\"DELETE FROM insights;\\",
        \\"\\",
        \\"BEGIN TRANSACTION;\\",
        \\"\\",
    ]
    
    total_patients = 0
    total_sessions = 0
    total_observations = 0
    total_interventions = 0
    
    # Gera pacientes
    print(f\\"Gerando {CONFIG['num_patients']} pacientes...\\")
    patient_values = []
    
    for i in range(1, CONFIG['num_patients'] + 1):
        patient_id, name, notes, created_at = generate_patient(i)
        safe_notes = notes.replace(\\"'\\", \\"''\\")
        patient_values.append(f\\"('{patient_id}', '{name}', '{safe_notes}', '{created_at}', '{created_at}')\\")
        total_patients += 1
    
    # Insere pacientes em batches
    batch_size = 100
    for i in range(0, len(patient_values), batch_size):
        batch = patient_values[i:i+batch_size]
        sql_lines.append(\\"INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES\\")
        sql_lines.append(\\",\n\\".join(batch) + \\";\\")
        sql_lines.append(\\"\\")
    
    # Gera sessões, observações e intervenções para cada paciente
    print(\\"Gerando sessões, observações e intervenções...\\")
    
    session_values = []
    observation_values = []
    intervention_values = []
    
    for patient_num in range(1, CONFIG['num_patients'] + 1):
        patient_id = f\\"p{patient_num:04d}\\"
        patient_context = random.choice(CLINICAL_CONTEXTS)
        
        num_sessions = random.randint(CONFIG['min_sessions_per_patient'], CONFIG['max_sessions_per_patient'])
        
        for session_idx in range(1, num_sessions + 1):
            session_id, _, session_date, summary, created_at, updated_at = generate_session(
                session_idx, patient_id, patient_context, session_idx, num_sessions
            )
            
            safe_summary = summary.replace(\\"'\\", \\"''\\")
            session_values.append(
                f\\"('{session_id}', '{patient_id}', '{session_date}', '{safe_summary}', '{created_at}', '{updated_at}')\\"
            )
            total_sessions += 1
            
            # Gera observações
            num_obs = random.randint(CONFIG['min_observations_per_session'], CONFIG['max_observations_per_session'])
            for obs_idx in range(1, num_obs + 1):
                obs_id, _, content, obs_created, obs_updated = generate_observation(session_id, patient_context, obs_idx)
                safe_content = content.replace(\\"'\\", \\"''\\")
                observation_values.append(f\\"('{obs_id}', '{session_id}', '{safe_content}', '{obs_created}', '{obs_updated}')\\")
                total_observations += 1
            
            # Gera intervenções
            num_int = random.randint(CONFIG['min_interventions_per_session'], CONFIG['max_interventions_per_session'])
            for int_idx in range(1, num_int + 1):
                int_id, _, int_content, int_created, int_updated = generate_intervention(session_id, patient_context, int_idx)
                safe_int_content = int_content.replace(\\"'\\", \\"''\\")
                intervention_values.append(f\\"('{int_id}', '{session_id}', '{safe_int_content}', '{int_created}', '{int_updated}')\\")
                total_interventions += 1
        
        if patient_num % 50 == 0:
            print(f\\"  Processados {patient_num}/{CONFIG['num_patients']} pacientes...\\")
    
    # Insere sessões em batches
    print(\\"Inserindo sessões no SQL...\\")
    for i in range(0, len(session_values), batch_size):
        batch = session_values[i:i+batch_size]
        sql_lines.append(\\"INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES\\")
        sql_lines.append(\\",\n\\".join(batch) + \\";\\")
        sql_lines.append(\\"\\")
    
    # Insere observações em batches
    print(\\"Inserindo observações no SQL...\\")
    for i in range(0, len(observation_values), batch_size):
        batch = observation_values[i:i+batch_size]
        sql_lines.append(\\"INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES\\")
        sql_lines.append(\\",\n\\".join(batch) + \\";\\")
        sql_lines.append(\\"\\")
    
    # Insere intervenções em batches
    print(\\"Inserindo intervenções no SQL...\\")
    for i in range(0, len(intervention_values), batch_size):
        batch = intervention_values[i:i+batch_size]
        sql_lines.append(\\"INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES\\")
        sql_lines.append(\\",\n\\".join(batch) + \\";\\")
        sql_lines.append(\\"\\")
    
    # Finaliza
    sql_lines.append(\\"COMMIT;\\")
    sql_lines.append(\\"\\")
    sql_lines.append(\\"-- ============================================\\")
    sql_lines.append(\\"-- ESTATÍSTICAS DO SEED\\")
    sql_lines.append(\\"-- ============================================\\")
    sql_lines.append(f\\"-- Total de pacientes: {total_patients}\\")
    sql_lines.append(f\\"-- Total de sessões: {total_sessions}\\")
    sql_lines.append(f\\"-- Total de observações: {total_observations}\\")
    sql_lines.append(f\\"-- Total de intervenções: {total_interventions}\\")
    sql_lines.append(f\\"-- Média de sessões por paciente: {total_sessions // total_patients}\\")
    sql_lines.append(\\"-- ============================================\\")
    
    print(f\\"\nResumo:\\")
    print(f\\"  Pacientes: {total_patients}\\")
    print(f\\"  Sessões: {total_sessions}\\")
    print(f\\"  Observações: {total_observations}\\")
    print(f\\"  Intervenções: {total_interventions}\\")
    
    return \\"\n\\".join(sql_lines)


if __name__ == \\"__main__\\":
    print(\\"=\\" * 60)
    print(\\"GERADOR DE MASSA DE DADOS CLÍNICA - ARANDU\\")
    print(\\"=\\" * 60)
    print()
    
    sql_content = generate_sql()
    
    # Salva arquivo
    output_file = \\"internal/infrastructure/repository/sqlite/seeds/seed_massive_clinical_data.sql\\"
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print(f\\"\n✅ Arquivo gerado: {output_file}\\")
    print(f\\"📊 Tamanho: {len(sql_content) / 1024 / 1024:.2f} MB\\")
    print()
    print(\\"Para executar:\\")
    print(f\\"  sqlite3 arandu.db < {output_file}\\")
