package patient

import "time"

type BiopsychosocialPanelViewModel struct {
	PatientID     string
	Medications   []MedicationListItemViewModel
	LatestVitals  *VitalsItemViewModel
	VitalsAverage *VitalsAverageItemViewModel
}

type MedicationListViewModel struct {
	PatientID   string
	Medications []MedicationListItemViewModel
	Error       string
}

type MedicationListItemViewModel struct {
	ID          string
	Name        string
	Dosage      string
	Frequency   string
	Prescriber  string
	Status      string
	StatusKey   string
	StatusLabel string
	StartedAt   string
	IsActive    bool
	IsSuspended bool
	IsFinished  bool
}

type VitalsWidgetViewModel struct {
	PatientID     string
	LatestVitals  *VitalsItemViewModel
	VitalsAverage *VitalsAverageItemViewModel
	Error         string
}

type VitalsItemViewModel struct {
	ID               string
	Date             string
	SleepHours       string
	AppetiteLevel    string
	Weight           string
	PhysicalActivity string
	Notes            string
	HasData          bool
}

type VitalsAverageItemViewModel struct {
	AvgSleepHours       string
	AvgAppetiteLevel    string
	AvgWeight           string
	AvgPhysicalActivity string
	RecordCount         int
	HasData             bool
}

func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}
