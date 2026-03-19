package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"arandu/internal/domain/intervention"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/patient"
	"arandu/internal/infrastructure/ai"
)

type ObservationRepository interface {
	FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*observation.Observation, error)
}

type InterventionRepository interface {
	FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*intervention.Intervention, error)
}

type PatientVitalsRepository interface {
	FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Vitals, error)
}

type PatientMedicationRepository interface {
	FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Medication, error)
}

type GeminiClient interface {
	GenerateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error)
	Close() error
}

type Cache interface {
	Get(patientID, timeframe string) (*ai.CacheEntry, bool)
	Set(patientID, timeframe, synthesis string)
}

type AIService struct {
	geminiClient GeminiClient
	cache        Cache
	obsRepo      ObservationRepository
	intervRepo   InterventionRepository
	vitalsRepo   PatientVitalsRepository
	medsRepo     PatientMedicationRepository
}

func NewAIService(
	geminiClient GeminiClient,
	cache Cache,
	obsRepo ObservationRepository,
	intervRepo InterventionRepository,
	vitalsRepo PatientVitalsRepository,
	medsRepo PatientMedicationRepository,
) *AIService {
	return &AIService{
		geminiClient: geminiClient,
		cache:        cache,
		obsRepo:      obsRepo,
		intervRepo:   intervRepo,
		vitalsRepo:   vitalsRepo,
		medsRepo:     medsRepo,
	}
}

type PatientSynthesisRequest struct {
	PatientID  string
	Timeframe  string
	MaxRetries int
}

type PatientSynthesisResponse struct {
	Synthesis   string
	GeneratedAt time.Time
}

func (s *AIService) GeneratePatientSynthesis(ctx context.Context, req PatientSynthesisRequest) (*PatientSynthesisResponse, error) {
	patientID := req.PatientID
	if patientID == "" {
		return nil, fmt.Errorf("patientID não pode ser vazio")
	}

	timeframe := req.Timeframe
	if timeframe == "" {
		timeframe = "3 months"
	}

	maxRetries := req.MaxRetries
	if maxRetries == 0 {
		maxRetries = 5
	}

	// Check cache first
	if s.cache != nil {
		if entry, found := s.cache.Get(patientID, timeframe); found {
			return &PatientSynthesisResponse{
				Synthesis:   entry.Synthesis,
				GeneratedAt: entry.GeneratedAt,
			}, nil
		}
	}

	clinicalData, err := s.collectClinicalData(ctx, patientID, timeframe)
	if err != nil {
		return nil, fmt.Errorf("falha ao coletar dados clínicos: %w", err)
	}

	if clinicalData.IsEmpty() {
		return nil, fmt.Errorf("dados clínicos insuficientes para análise")
	}

	prompt := s.buildPrompt(clinicalData)

	synthesis, err := s.geminiClient.GenerateWithRetry(ctx, prompt, maxRetries)
	if err != nil {
		return nil, fmt.Errorf("falha na geração da síntese: %w", err)
	}

	// Store in cache
	if s.cache != nil {
		s.cache.Set(patientID, timeframe, synthesis)
	}

	return &PatientSynthesisResponse{
		Synthesis:   synthesis,
		GeneratedAt: time.Now(),
	}, nil
}

type ClinicalData struct {
	Observations  []*observation.Observation
	Interventions []*intervention.Intervention
	Vitals        []*patient.Vitals
	Medications   []*patient.Medication
}

func (cd ClinicalData) IsEmpty() bool {
	return len(cd.Observations) == 0 &&
		len(cd.Interventions) == 0 &&
		len(cd.Vitals) == 0 &&
		len(cd.Medications) == 0
}

func (s *AIService) collectClinicalData(ctx context.Context, patientID, timeframe string) (ClinicalData, error) {
	var data ClinicalData
	var err error

	startTime := s.calculateStartTime(timeframe)

	data.Observations, err = s.obsRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
	if err != nil {
		return data, fmt.Errorf("erro ao buscar observações: %w", err)
	}

	data.Interventions, err = s.intervRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
	if err != nil {
		return data, fmt.Errorf("erro ao buscar intervenções: %w", err)
	}

	data.Vitals, err = s.vitalsRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
	if err != nil {
		return data, fmt.Errorf("erro ao buscar sinais vitais: %w", err)
	}

	data.Medications, err = s.medsRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
	if err != nil {
		return data, fmt.Errorf("erro ao buscar medicações: %w", err)
	}

	return data, nil
}

func (s *AIService) calculateStartTime(timeframe string) time.Time {
	now := time.Now()

	switch strings.ToLower(timeframe) {
	case "1 month":
		return now.AddDate(0, -1, 0)
	case "3 months":
		return now.AddDate(0, -3, 0)
	case "6 months":
		return now.AddDate(0, -6, 0)
	case "1 year":
		return now.AddDate(-1, 0, 0)
	case "all":
		return time.Time{}
	default:
		return now.AddDate(0, -3, 0)
	}
}

func (s *AIService) buildPrompt(data ClinicalData) string {
	var prompt strings.Builder

	prompt.WriteString("## Dados Clínicos para Análise\n\n")

	if len(data.Observations) > 0 {
		prompt.WriteString("### Observações Clínicas:\n")
		for i, obs := range data.Observations {
			if i < 10 {
				prompt.WriteString(fmt.Sprintf("- %s (Data: %s)\n",
					obs.Content,
					obs.CreatedAt.Format("02/01/2006")))
			}
		}
		if len(data.Observations) > 10 {
			prompt.WriteString(fmt.Sprintf("... e mais %d observações\n", len(data.Observations)-10))
		}
		prompt.WriteString("\n")
	}

	if len(data.Interventions) > 0 {
		prompt.WriteString("### Intervenções Terapêuticas:\n")
		for i, interv := range data.Interventions {
			if i < 10 {
				prompt.WriteString(fmt.Sprintf("- %s (Data: %s)\n",
					interv.Content,
					interv.CreatedAt.Format("02/01/2006")))
			}
		}
		if len(data.Interventions) > 10 {
			prompt.WriteString(fmt.Sprintf("... e mais %d intervenções\n", len(data.Interventions)-10))
		}
		prompt.WriteString("\n")
	}

	if len(data.Vitals) > 0 {
		prompt.WriteString("### Sinais Vitais:\n")
		for i, vital := range data.Vitals {
			if i < 5 {
				sleepHours := ""
				if vital.SleepHours != nil {
					sleepHours = fmt.Sprintf("%.1f horas", *vital.SleepHours)
				}
				weight := ""
				if vital.Weight != nil {
					weight = fmt.Sprintf("%.1fkg", *vital.Weight)
				}
				prompt.WriteString(fmt.Sprintf("- Peso: %s, Sono: %s, Atividade Física: %d/10 (Data: %s)\n",
					weight,
					sleepHours,
					vital.PhysicalActivity,
					vital.Date.Format("02/01/2006")))
			}
		}
		prompt.WriteString("\n")
	}

	if len(data.Medications) > 0 {
		prompt.WriteString("### Histórico Farmacológico:\n")
		for i, med := range data.Medications {
			if i < 5 {
				endDate := ""
				if med.EndedAt != nil {
					endDate = med.EndedAt.Format("02/01/2006")
				}
				prompt.WriteString(fmt.Sprintf("- %s: %s %s (Início: %s, Fim: %s, Status: %s)\n",
					med.Name,
					med.Dosage,
					med.Frequency,
					med.StartedAt.Format("02/01/2006"),
					endDate,
					string(med.Status)))
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("## Instruções para Análise:\n")
	prompt.WriteString("Analise os dados acima e forneça uma síntese reflexiva estruturada em:\n")
	prompt.WriteString("1. Temas Dominantes\n")
	prompt.WriteString("2. Pontos de Inflexão\n")
	prompt.WriteString("3. Correlações Sugeridas\n")
	prompt.WriteString("4. Provocação Clínica\n")
	prompt.WriteString("\nSeja técnico, reflexivo e evite diagnósticos fechados.")

	return prompt.String()
}
