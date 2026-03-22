package patient

import (
	"testing"
	"time"
)

func TestNewAnamnesis(t *testing.T) {
	t.Run("valid patient ID", func(t *testing.T) {
		a, err := NewAnamnesis("patient-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if a.PatientID != "patient-123" {
			t.Errorf("expected patient_id 'patient-123', got %s", a.PatientID)
		}
		if a.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("empty patient ID", func(t *testing.T) {
		_, err := NewAnamnesis("")
		if err == nil {
			t.Error("expected error for empty patient ID")
		}
	})

	t.Run("patient ID too long", func(t *testing.T) {
		longID := ""
		for i := 0; i < 40; i++ {
			longID += "x"
		}
		_, err := NewAnamnesis(longID)
		if err == nil {
			t.Error("expected error for too long patient ID")
		}
	})
}

func TestAnamnesisUpdateSection(t *testing.T) {
	a, _ := NewAnamnesis("patient-123")

	t.Run("update chief complaint", func(t *testing.T) {
		err := a.UpdateSection("chief_complaint", "Anxiety symptoms")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if a.ChiefComplaint != "Anxiety symptoms" {
			t.Errorf("expected 'Anxiety symptoms', got %s", a.ChiefComplaint)
		}
	})

	t.Run("update personal history", func(t *testing.T) {
		err := a.UpdateSection("personal_history", "Born in São Paulo")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if a.PersonalHistory != "Born in São Paulo" {
			t.Errorf("expected 'Born in São Paulo', got %s", a.PersonalHistory)
		}
	})

	t.Run("update family history", func(t *testing.T) {
		err := a.UpdateSection("family_history", "Mother with depression")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if a.FamilyHistory != "Mother with depression" {
			t.Errorf("expected 'Mother with depression', got %s", a.FamilyHistory)
		}
	})

	t.Run("update mental state exam", func(t *testing.T) {
		err := a.UpdateSection("mental_state_exam", "Alert and oriented")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if a.MentalStateExam != "Alert and oriented" {
			t.Errorf("expected 'Alert and oriented', got %s", a.MentalStateExam)
		}
	})

	t.Run("invalid section", func(t *testing.T) {
		err := a.UpdateSection("invalid", "content")
		if err == nil {
			t.Error("expected error for invalid section")
		}
	})
}

func TestAnamnesisIsEmpty(t *testing.T) {
	t.Run("empty anamnesis", func(t *testing.T) {
		a, _ := NewAnamnesis("patient-123")
		if !a.IsEmpty() {
			t.Error("expected empty anamnesis to return true")
		}
	})

	t.Run("non-empty anamnesis", func(t *testing.T) {
		a, _ := NewAnamnesis("patient-123")
		a.ChiefComplaint = "Some complaint"
		if a.IsEmpty() {
			t.Error("expected non-empty anamnesis to return false")
		}
	})
}

func TestAnamnesisValidate(t *testing.T) {
	t.Run("valid anamnesis", func(t *testing.T) {
		a, _ := NewAnamnesis("patient-123")
		err := a.Validate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("nil anamnesis", func(t *testing.T) {
		var a *Anamnesis
		err := a.Validate()
		if err == nil {
			t.Error("expected error for nil anamnesis")
		}
	})

	t.Run("empty patient ID", func(t *testing.T) {
		a, _ := NewAnamnesis("")
		err := a.Validate()
		if err == nil {
			t.Error("expected error for empty patient ID")
		}
	})
}

func TestAnamnesisUpdatedAtSet(t *testing.T) {
	a, _ := NewAnamnesis("patient-123")
	before := time.Now()

	err := a.UpdateSection("chief_complaint", "Test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if a.UpdatedAt.Before(before) {
		t.Error("UpdatedAt should be updated after UpdateSection")
	}
}
