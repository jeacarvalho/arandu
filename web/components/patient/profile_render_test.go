package patient

import (
	"bytes"
	"testing"
	"time"

	domainPatient "arandu/internal/domain/patient"
)

func newTestPatient() *domainPatient.Patient {
	return &domainPatient.Patient{
		ID:        "patient-test-123",
		Name:      "Ana Souza",
		Gender:    "f",
		Ethnicity: "b",
		Occupation: "Professora",
		Education: "s",
		Notes:     "Paciente de teste",
		CreatedAt: time.Now().Add(-180 * 24 * time.Hour), // 6 meses atrás
		UpdatedAt: time.Now(),
	}
}

// TestPatientProfileView_NoFakeStatusTag verifica que o perfil não exibe
// o status "Em tratamento" hardcoded, que era falso para todos os pacientes
// independente do estado real.
func TestPatientProfileView_NoFakeStatusTag(t *testing.T) {
	p := newTestPatient()

	var buf bytes.Buffer
	err := PatientProfileView(p, nil, nil, nil).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("Failed to render PatientProfileView: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte("Em tratamento")) {
		t.Errorf("PatientProfileView não deve conter status hardcoded 'Em tratamento'.\nHTML: %s", html)
	}
}

// TestPatientProfileView_NoFakeDiagnosisTag verifica que o perfil não exibe
// o placeholder "TAG" como diagnóstico, que nunca foi substituído por dado real.
func TestPatientProfileView_NoFakeDiagnosisTag(t *testing.T) {
	p := newTestPatient()

	var buf bytes.Buffer
	err := PatientProfileView(p, nil, nil, nil).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("Failed to render PatientProfileView: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`class="patient-tag-v2 tag-diagnosis"`)) {
		t.Errorf("PatientProfileView não deve conter tag de diagnóstico hardcoded.\nHTML: %s", html)
	}
}

// TestPatientProfileView_NoMisleadingAgeFromRegistration verifica que o perfil
// não exibe o tempo desde o cadastro com o label "anos", que induzia leitura como
// idade biológica do paciente — dado que o sistema não possui.
func TestPatientProfileView_NoMisleadingAgeFromRegistration(t *testing.T) {
	p := newTestPatient()

	var buf bytes.Buffer
	err := PatientProfileView(p, nil, nil, nil).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("Failed to render PatientProfileView: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`tag-age`)) {
		t.Errorf("PatientProfileView não deve conter tag de idade calculada a partir da data de cadastro.\nHTML: %s", html)
	}
}

// TestPatientProfileView_RendersPatientName verifica que o nome do paciente
// está presente no HTML renderizado.
func TestPatientProfileView_RendersPatientName(t *testing.T) {
	p := newTestPatient()

	var buf bytes.Buffer
	err := PatientProfileView(p, nil, nil, nil).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("Failed to render PatientProfileView: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Ana Souza")) {
		t.Errorf("PatientProfileView deve conter o nome do paciente 'Ana Souza'")
	}
}

// TestPatientProfileView_RendersTherapyDuration verifica que o tempo em terapia
// (calculado a partir da data de cadastro) está presente e é o único indicador
// temporal exibido.
func TestPatientProfileView_RendersTherapyDuration(t *testing.T) {
	p := newTestPatient()

	var buf bytes.Buffer
	err := PatientProfileView(p, nil, nil, nil).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("Failed to render PatientProfileView: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("em terapia")) {
		t.Errorf("PatientProfileView deve exibir o tempo em terapia com label 'em terapia'")
	}
}
