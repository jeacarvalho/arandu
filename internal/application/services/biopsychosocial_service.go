package services

import (
	"time"

	"arandu/internal/domain/patient"
	"arandu/internal/infrastructure/repository/sqlite"
)

type BiopsychosocialService struct {
	medicationRepo *sqlite.MedicationRepository
	vitalsRepo     *sqlite.VitalsRepository
}

func NewBiopsychosocialService(medicationRepo *sqlite.MedicationRepository, vitalsRepo *sqlite.VitalsRepository) *BiopsychosocialService {
	return &BiopsychosocialService{
		medicationRepo: medicationRepo,
		vitalsRepo:     vitalsRepo,
	}
}

func (s *BiopsychosocialService) AddMedication(patientID, name, dosage, frequency, prescriber string, startedAt time.Time) (*patient.Medication, error) {
	medication, err := patient.NewMedication(patientID, name, dosage, frequency, prescriber, startedAt)
	if err != nil {
		return nil, err
	}
	if err := s.medicationRepo.Save(medication); err != nil {
		return nil, err
	}
	return medication, nil
}

func (s *BiopsychosocialService) GetMedications(patientID string) ([]*patient.Medication, error) {
	return s.medicationRepo.FindByPatientID(patientID)
}

func (s *BiopsychosocialService) GetActiveMedications(patientID string) ([]*patient.Medication, error) {
	return s.medicationRepo.GetActiveMedications(patientID)
}

func (s *BiopsychosocialService) UpdateMedicationStatus(medicationID string, status patient.MedicationStatus) (*patient.Medication, error) {
	if err := s.medicationRepo.UpdateStatus(medicationID, status); err != nil {
		return nil, err
	}
	return s.medicationRepo.FindByID(medicationID)
}

func (s *BiopsychosocialService) SuspendMedication(medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(medicationID, patient.MedicationStatusSuspended)
}

func (s *BiopsychosocialService) ActivateMedication(medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(medicationID, patient.MedicationStatusActive)
}

func (s *BiopsychosocialService) FinishMedication(medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(medicationID, patient.MedicationStatusFinished)
}

func (s *BiopsychosocialService) RecordVitals(patientID string, date time.Time, sleepHours *float64, appetiteLevel *int, weight *float64, physicalActivity int, notes string) (*patient.Vitals, error) {
	vitals, err := patient.NewVitals(patientID, date, sleepHours, appetiteLevel, weight, physicalActivity, notes)
	if err != nil {
		return nil, err
	}
	if err := s.vitalsRepo.Save(vitals); err != nil {
		return nil, err
	}
	return vitals, nil
}

func (s *BiopsychosocialService) GetLatestVitals(patientID string) (*patient.Vitals, error) {
	return s.vitalsRepo.GetLatestVitals(patientID)
}

func (s *BiopsychosocialService) GetVitalsHistory(patientID string, limit int) ([]*patient.Vitals, error) {
	return s.vitalsRepo.FindByPatientID(patientID, limit)
}

func (s *BiopsychosocialService) GetAverageVitals(patientID string, days int) (*sqlite.VitalsAverage, error) {
	return s.vitalsRepo.GetAverageVitals(patientID, days)
}

type BiopsychosocialContext struct {
	PatientID         string
	ActiveMedications []*patient.Medication
	AllMedications    []*patient.Medication
	LatestVitals      *patient.Vitals
	VitalsAverage     *sqlite.VitalsAverage
}

func (s *BiopsychosocialService) GetContext(patientID string) (*BiopsychosocialContext, error) {
	activeMeds, err := s.GetActiveMedications(patientID)
	if err != nil {
		return nil, err
	}
	allMeds, err := s.GetMedications(patientID)
	if err != nil {
		return nil, err
	}
	latestVitals, err := s.GetLatestVitals(patientID)
	if err != nil {
		return nil, err
	}
	avgVitals, err := s.GetAverageVitals(patientID, 30)
	if err != nil {
		return nil, err
	}

	return &BiopsychosocialContext{
		PatientID:         patientID,
		ActiveMedications: activeMeds,
		AllMedications:    allMeds,
		LatestVitals:      latestVitals,
		VitalsAverage:     avgVitals,
	}, nil
}
