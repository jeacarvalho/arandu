package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"arandu/internal/domain/intervention"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/patient"
	"arandu/internal/infrastructure/ai"
)

type mockAIObservationRepository struct {
	observations []*observation.Observation
	err          error
}

func (m *mockAIObservationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*observation.Observation, error) {
	return m.observations, m.err
}

type mockInterventionRepository struct {
	interventions []*intervention.Intervention
	err           error
}

func (m *mockInterventionRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*intervention.Intervention, error) {
	return m.interventions, m.err
}

type mockPatientVitalsRepository struct {
	vitals []*patient.Vitals
	err    error
}

func (m *mockPatientVitalsRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Vitals, error) {
	return m.vitals, m.err
}

type mockPatientMedicationRepository struct {
	medications []*patient.Medication
	err         error
}

func (m *mockPatientMedicationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Medication, error) {
	return m.medications, m.err
}

type mockCache struct {
	getFunc func(patientID, timeframe string) (*ai.CacheEntry, bool)
	setFunc func(patientID, timeframe, synthesis string)
}

func (m *mockCache) Get(patientID, timeframe string) (*ai.CacheEntry, bool) {
	if m.getFunc != nil {
		return m.getFunc(patientID, timeframe)
	}
	return nil, false
}

func (m *mockCache) Set(patientID, timeframe, synthesis string) {
	if m.setFunc != nil {
		m.setFunc(patientID, timeframe, synthesis)
	}
}

type mockAIGeminiClient struct {
	response string
	err      error
}

func (m *mockAIGeminiClient) GenerateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error) {
	return m.response, m.err
}

func (m *mockAIGeminiClient) Close() error {
	return nil
}

func TestAIService_GeneratePatientSynthesis(t *testing.T) {
	now := time.Now()

	testObservations := []*observation.Observation{
		{
			ID:        "obs1",
			SessionID: "sess1",
			Content:   "Paciente relatou ansiedade elevada durante a sessão",
			CreatedAt: now.AddDate(0, -1, 0),
		},
		{
			ID:        "obs2",
			SessionID: "sess1",
			Content:   "Demonstrou progresso na regulação emocional",
			CreatedAt: now.AddDate(0, -2, 0),
		},
	}

	testInterventions := []*intervention.Intervention{
		{
			ID:        "int1",
			SessionID: "sess1",
			Content:   "Técnica de respiração diafragmática",
			CreatedAt: now.AddDate(0, -1, 0),
		},
	}

	weight := 70.5
	sleepHours := 7.5
	appetiteLevel := 7
	testVitals := []*patient.Vitals{
		{
			ID:               "vital1",
			PatientID:        "patient1",
			Date:             now.AddDate(0, -1, 0),
			SleepHours:       &sleepHours,
			AppetiteLevel:    &appetiteLevel,
			Weight:           &weight,
			PhysicalActivity: 6,
			Notes:            "Paciente relatou sono mais regular",
		},
	}

	endDate := now.AddDate(0, 0, -15)
	testMedications := []*patient.Medication{
		{
			ID:         "med1",
			PatientID:  "patient1",
			Name:       "Sertralina",
			Dosage:     "50mg",
			Frequency:  "1x/dia",
			Prescriber: "Dr. Silva",
			Status:     patient.MedicationStatusActive,
			StartedAt:  now.AddDate(0, -3, 0),
			EndedAt:    &endDate,
		},
	}

	mockGemini := &mockAIGeminiClient{
		response: "Síntese reflexiva gerada com sucesso",
		err:      nil,
	}

	mockObsRepo := &mockAIObservationRepository{
		observations: testObservations,
		err:          nil,
	}

	mockIntervRepo := &mockInterventionRepository{
		interventions: testInterventions,
		err:           nil,
	}

	mockVitalsRepo := &mockPatientVitalsRepository{
		vitals: testVitals,
		err:    nil,
	}

	mockMedsRepo := &mockPatientMedicationRepository{
		medications: testMedications,
		err:         nil,
	}

	mockCache := &mockCache{}
	service := NewAIService(mockGemini, mockCache, mockObsRepo, mockIntervRepo, mockVitalsRepo, mockMedsRepo)

	req := PatientSynthesisRequest{
		PatientID:  "patient1",
		Timeframe:  "3 months",
		MaxRetries: 3,
	}

	ctx := context.Background()
	resp, err := service.GeneratePatientSynthesis(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Synthesis != "Síntese reflexiva gerada com sucesso" {
		t.Errorf("Expected synthesis 'Síntese reflexiva gerada com sucesso', got %s", resp.Synthesis)
	}

	if resp.GeneratedAt.IsZero() {
		t.Error("Expected GeneratedAt to be set")
	}
}

func TestAIService_GeneratePatientSynthesis_EmptyData(t *testing.T) {
	mockGemini := &mockAIGeminiClient{
		response: "",
		err:      nil,
	}

	mockObsRepo := &mockAIObservationRepository{
		observations: []*observation.Observation{},
		err:          nil,
	}

	mockIntervRepo := &mockInterventionRepository{
		interventions: []*intervention.Intervention{},
		err:           nil,
	}

	mockVitalsRepo := &mockPatientVitalsRepository{
		vitals: []*patient.Vitals{},
		err:    nil,
	}

	mockMedsRepo := &mockPatientMedicationRepository{
		medications: []*patient.Medication{},
		err:         nil,
	}

	mockCache := &mockCache{}
	service := NewAIService(mockGemini, mockCache, mockObsRepo, mockIntervRepo, mockVitalsRepo, mockMedsRepo)

	req := PatientSynthesisRequest{
		PatientID: "patient1",
		Timeframe: "3 months",
	}

	ctx := context.Background()
	_, err := service.GeneratePatientSynthesis(ctx, req)

	if err == nil {
		t.Fatal("Expected error for empty data, got nil")
	}

	expectedErr := "dados clínicos insuficientes para análise"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
	}
}

func TestAIService_CalculateStartTime(t *testing.T) {
	service := &AIService{}
	now := time.Now()

	testCases := []struct {
		timeframe string
		expected  time.Time
	}{
		{"1 month", now.AddDate(0, -1, 0)},
		{"3 months", now.AddDate(0, -3, 0)},
		{"6 months", now.AddDate(0, -6, 0)},
		{"1 year", now.AddDate(-1, 0, 0)},
		{"all", time.Time{}},
		{"invalid", now.AddDate(0, -3, 0)}, // default
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			result := service.calculateStartTime(tc.timeframe)

			if tc.timeframe == "all" {
				if !result.IsZero() {
					t.Errorf("Expected zero time for 'all', got %v", result)
				}
				return
			}

			if tc.timeframe == "invalid" {
				expected := now.AddDate(0, -3, 0)
				diff := result.Sub(expected)
				if diff > time.Hour || diff < -time.Hour {
					t.Errorf("Expected ~%v for invalid timeframe, got %v", expected, result)
				}
				return
			}

			diff := result.Sub(tc.expected)
			if diff > time.Hour || diff < -time.Hour {
				t.Errorf("Expected ~%v for timeframe %s, got %v", tc.expected, tc.timeframe, result)
			}
		})
	}
}

func TestAIService_BuildPrompt(t *testing.T) {
	service := &AIService{}
	now := time.Now()

	weight := 70.5
	sleepHours := 7.5
	appetiteLevel := 7
	endDate := now.AddDate(0, 0, -15)

	data := ClinicalData{
		Observations: []*observation.Observation{
			{
				ID:        "obs1",
				SessionID: "sess1",
				Content:   "Paciente relatou ansiedade elevada",
				CreatedAt: now.AddDate(0, -1, 0),
			},
		},
		Interventions: []*intervention.Intervention{
			{
				ID:        "int1",
				SessionID: "sess1",
				Content:   "Técnica de respiração",
				CreatedAt: now.AddDate(0, -1, 0),
			},
		},
		Vitals: []*patient.Vitals{
			{
				ID:               "vital1",
				PatientID:        "patient1",
				Date:             now.AddDate(0, -1, 0),
				SleepHours:       &sleepHours,
				AppetiteLevel:    &appetiteLevel,
				Weight:           &weight,
				PhysicalActivity: 6,
				Notes:            "Sono regular",
			},
		},
		Medications: []*patient.Medication{
			{
				ID:         "med1",
				PatientID:  "patient1",
				Name:       "Sertralina",
				Dosage:     "50mg",
				Frequency:  "1x/dia",
				Prescriber: "Dr. Silva",
				Status:     patient.MedicationStatusActive,
				StartedAt:  now.AddDate(0, -3, 0),
				EndedAt:    &endDate,
			},
		},
	}

	prompt := service.buildPrompt(data)

	expectedSections := []string{
		"## Dados Clínicos para Análise",
		"### Observações Clínicas:",
		"### Intervenções Terapêuticas:",
		"### Sinais Vitais:",
		"### Histórico Farmacológico:",
		"## Instruções para Análise:",
		"1. Temas Dominantes",
		"2. Pontos de Inflexão",
		"3. Correlações Sugeridas",
		"4. Provocação Clínica",
	}

	for _, section := range expectedSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("Expected prompt to contain section: %s", section)
		}
	}

	if !strings.Contains(prompt, "Paciente relatou ansiedade elevada") {
		t.Error("Expected prompt to contain observation content")
	}

	if !strings.Contains(prompt, "Técnica de respiração") {
		t.Error("Expected prompt to contain intervention content")
	}

	if !strings.Contains(prompt, "70.5kg") {
		t.Error("Expected prompt to contain weight")
	}

	if !strings.Contains(prompt, "Sertralina") {
		t.Error("Expected prompt to contain medication name")
	}
}
