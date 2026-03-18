package patient

import (
	"testing"
	"time"
)

func TestNewVitals(t *testing.T) {
	patientID := "patient-123"
	date := time.Now().AddDate(0, 0, -1)
	sleepHours := 8.0
	appetiteLevel := 7
	weight := 70.5
	physicalActivity := 3
	notes := "Paciente dormiu bem"

	vitals, err := NewVitals(patientID, date, &sleepHours, &appetiteLevel, &weight, physicalActivity, notes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if vitals.ID == "" {
		t.Error("Expected vitals ID to be set")
	}
	if vitals.PatientID != patientID {
		t.Errorf("Expected patient ID %s, got %s", patientID, vitals.PatientID)
	}
	if vitals.Date != date {
		t.Errorf("Expected date %v, got %v", date, vitals.Date)
	}
	if *vitals.SleepHours != sleepHours {
		t.Errorf("Expected sleep hours %f, got %f", sleepHours, *vitals.SleepHours)
	}
	if *vitals.AppetiteLevel != appetiteLevel {
		t.Errorf("Expected appetite level %d, got %d", appetiteLevel, *vitals.AppetiteLevel)
	}
	if *vitals.Weight != weight {
		t.Errorf("Expected weight %f, got %f", weight, *vitals.Weight)
	}
	if vitals.PhysicalActivity != physicalActivity {
		t.Errorf("Expected physical activity %d, got %d", physicalActivity, vitals.PhysicalActivity)
	}
	if vitals.Notes != notes {
		t.Errorf("Expected notes %s, got %s", notes, vitals.Notes)
	}
}

func TestNewVitalsValidation(t *testing.T) {
	lowAppetite := 0
	highAppetite := 11
	lowSleep := -1.0
	highSleep := 25.0

	tests := []struct {
		name          string
		patientID     string
		date          time.Time
		sleepHours    *float64
		appetiteLevel *int
		expectErr     bool
		errMsg        string
	}{
		{
			name:      "empty patient ID",
			patientID: "",
			date:      time.Now(),
			expectErr: true,
			errMsg:    "patient ID cannot be empty",
		},
		{
			name:      "future date",
			patientID: "patient-123",
			date:      time.Now().AddDate(0, 0, 1),
			expectErr: true,
			errMsg:    "vitals date cannot be in the future",
		},
		{
			name:          "low appetite level",
			patientID:     "patient-123",
			date:          time.Now(),
			appetiteLevel: &lowAppetite,
			expectErr:     true,
			errMsg:        "appetite level must be between 1 and 10",
		},
		{
			name:          "high appetite level",
			patientID:     "patient-123",
			date:          time.Now(),
			appetiteLevel: &highAppetite,
			expectErr:     true,
			errMsg:        "appetite level must be between 1 and 10",
		},
		{
			name:       "low sleep hours",
			patientID:  "patient-123",
			date:       time.Now(),
			sleepHours: &lowSleep,
			expectErr:  true,
			errMsg:     "sleep hours must be between 0 and 24",
		},
		{
			name:       "high sleep hours",
			patientID:  "patient-123",
			date:       time.Now(),
			sleepHours: &highSleep,
			expectErr:  true,
			errMsg:     "sleep hours must be between 0 and 24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVitals(tt.patientID, tt.date, tt.sleepHours, tt.appetiteLevel, nil, 0, "")
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error %s, got nil", tt.errMsg)
				} else if err.Error() != tt.errMsg {
					t.Errorf("Expected error %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestNewVitalsNilValues(t *testing.T) {
	vitals, err := NewVitals("patient-123", time.Now(), nil, nil, nil, 0, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if vitals.SleepHours != nil {
		t.Error("Expected sleep hours to be nil")
	}
	if vitals.AppetiteLevel != nil {
		t.Error("Expected appetite level to be nil")
	}
	if vitals.Weight != nil {
		t.Error("Expected weight to be nil")
	}
}

func TestVitalsUpdate(t *testing.T) {
	sleepHours := 8.0
	appetiteLevel := 7
	vitals, _ := NewVitals("patient-123", time.Now(), &sleepHours, &appetiteLevel, nil, 0, "")

	newSleepHours := 6.0
	newAppetiteLevel := 5
	newWeight := 69.5
	newPhysicalActivity := 4
	newNotes := "Atualizado"

	err := vitals.Update(&newSleepHours, &newAppetiteLevel, &newWeight, newPhysicalActivity, newNotes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if *vitals.SleepHours != newSleepHours {
		t.Errorf("Expected sleep hours %f, got %f", newSleepHours, *vitals.SleepHours)
	}
	if *vitals.AppetiteLevel != newAppetiteLevel {
		t.Errorf("Expected appetite level %d, got %d", newAppetiteLevel, *vitals.AppetiteLevel)
	}
	if *vitals.Weight != newWeight {
		t.Errorf("Expected weight %f, got %f", newWeight, *vitals.Weight)
	}
	if vitals.PhysicalActivity != newPhysicalActivity {
		t.Errorf("Expected physical activity %d, got %d", newPhysicalActivity, vitals.PhysicalActivity)
	}
	if vitals.Notes != newNotes {
		t.Errorf("Expected notes %s, got %s", newNotes, vitals.Notes)
	}
}

func TestVitalsUpdateInvalidAppetite(t *testing.T) {
	vitals, _ := NewVitals("patient-123", time.Now(), nil, nil, nil, 0, "")

	invalidAppetite := 15
	err := vitals.Update(nil, &invalidAppetite, nil, 0, "")
	if err == nil {
		t.Error("Expected error for invalid appetite level")
	}
}

func TestVitalsUpdateInvalidSleep(t *testing.T) {
	vitals, _ := NewVitals("patient-123", time.Now(), nil, nil, nil, 0, "")

	invalidSleep := 30.0
	err := vitals.Update(&invalidSleep, nil, nil, 0, "")
	if err == nil {
		t.Error("Expected error for invalid sleep hours")
	}
}
