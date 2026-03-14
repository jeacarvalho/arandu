package handlers

import (
	"time"
)

// DashboardData contém todos os dados necessários para o dashboard clínico
type DashboardData struct {
	ActivePatients   []ActivePatient
	RecentSessions   []RecentSession
	AIInsights       []AIInsight
	EmergingPatterns []EmergingPattern
	TotalSessions    int
}

// ActivePatient representa um paciente ativo para o dashboard
type ActivePatient struct {
	ID           string
	Name         string
	LastSession  time.Time
	SessionCount int
	CreatedAt    time.Time
	Notes        string
}

// RecentSession representa uma sessão recente para o dashboard
type RecentSession struct {
	ID            string
	PatientID     string
	PatientName   string
	Date          time.Time
	Summary       string
	SessionNumber int
}

// AIInsight representa um insight gerado por IA
type AIInsight struct {
	ID         string
	Title      string
	Content    string
	Confidence float64 // 0.0 a 1.0
	CreatedAt  time.Time
}

// EmergingPattern representa um padrão emergente detectado
type EmergingPattern struct {
	ID            string
	Theme         string
	Description   string
	PatientCount  int
	FirstDetected time.Time
}

// MockDashboardData retorna dados mock para o dashboard
func MockDashboardData() DashboardData {
	now := time.Now()

	return DashboardData{
		TotalSessions: 25,
		ActivePatients: []ActivePatient{
			{
				ID:           "1",
				Name:         "Maria Silva",
				LastSession:  now.Add(-24 * time.Hour), // ontem
				SessionCount: 12,
				CreatedAt:    now.Add(-30 * 24 * time.Hour), // 30 dias atrás
				Notes:        "Ansiedade generalizada, foco em desempenho profissional",
			},
			{
				ID:           "2",
				Name:         "João Pereira",
				LastSession:  now.Add(-48 * time.Hour), // 2 dias atrás
				SessionCount: 8,
				CreatedAt:    now.Add(-45 * 24 * time.Hour), // 45 dias atrás
				Notes:        "Dificuldades com limites pessoais e assertividade",
			},
			{
				ID:           "3",
				Name:         "Ana Rodrigues",
				LastSession:  now.Add(-72 * time.Hour), // 3 dias atrás
				SessionCount: 5,
				CreatedAt:    now.Add(-20 * 24 * time.Hour), // 20 dias atrás
				Notes:        "Processamento de luto recente, apoio no ajuste emocional",
			},
		},
		RecentSessions: []RecentSession{
			{
				ID:            "s1",
				PatientID:     "1",
				PatientName:   "Maria S.",
				Date:          now.Add(-24 * time.Hour),
				Summary:       "Ansiedade relacionada à avaliação de desempenho no trabalho",
				SessionNumber: 12,
			},
			{
				ID:            "s2",
				PatientID:     "2",
				PatientName:   "João P.",
				Date:          now.Add(-48 * time.Hour),
				Summary:       "Dificuldades em estabelecer limites pessoais",
				SessionNumber: 8,
			},
			{
				ID:            "s3",
				PatientID:     "3",
				PatientName:   "Ana R.",
				Date:          now.Add(-72 * time.Hour),
				Summary:       "Processamento de luto recente",
				SessionNumber: 5,
			},
		},
		AIInsights: []AIInsight{
			{
				ID:         "i1",
				Title:      "Padrão detectado",
				Content:    "Ansiedade relacionada à avaliação de desempenho aparece em múltiplos pacientes.",
				Confidence: 0.85,
				CreatedAt:  now.Add(-2 * time.Hour),
			},
			{
				ID:         "i2",
				Title:      "Possível correlação",
				Content:    "Pacientes com dificuldades de limite mostram padrões similares de estresse.",
				Confidence: 0.72,
				CreatedAt:  now.Add(-5 * time.Hour),
			},
		},
		EmergingPatterns: []EmergingPattern{
			{
				ID:            "p1",
				Theme:         "Ansiedade social",
				Description:   "Medo de avaliação social em contextos profissionais",
				PatientCount:  4,
				FirstDetected: now.Add(-7 * 24 * time.Hour), // 1 semana atrás
			},
			{
				ID:            "p2",
				Theme:         "Conflito com autoridade",
				Description:   "Dificuldade em lidar com figuras de autoridade",
				PatientCount:  3,
				FirstDetected: now.Add(-14 * 24 * time.Hour), // 2 semanas atrás
			},
		},
	}
}
