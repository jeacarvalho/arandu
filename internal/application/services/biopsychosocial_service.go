package services

import (
	"context"
	"time"

	"arandu/internal/domain/patient"
)

type BiopsychosocialService struct {
	medicationRepo patient.MedicationRepository
	vitalsRepo     patient.VitalsRepository
}

func NewBiopsychosocialService(medicationRepo patient.MedicationRepository, vitalsRepo patient.VitalsRepository) *BiopsychosocialService {
	return &BiopsychosocialService{
		medicationRepo: medicationRepo,
		vitalsRepo:     vitalsRepo,
	}
}

func (s *BiopsychosocialService) AddMedication(ctx context.Context, patientID, name, dosage, frequency, prescriber string, startedAt time.Time) (*patient.Medication, error) {
	medication, err := patient.NewMedication(patientID, name, dosage, frequency, prescriber, startedAt)
	if err != nil {
		return nil, err
	}
	if err := s.medicationRepo.Save(ctx, medication); err != nil {
		return nil, err
	}
	return medication, nil
}

func (s *BiopsychosocialService) GetMedications(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	return s.medicationRepo.FindByPatientID(ctx, patientID)
}

func (s *BiopsychosocialService) GetActiveMedications(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	return s.medicationRepo.GetActiveMedications(ctx, patientID)
}

func (s *BiopsychosocialService) UpdateMedicationStatus(ctx context.Context, medicationID string, status patient.MedicationStatus) (*patient.Medication, error) {
	if err := s.medicationRepo.UpdateStatus(ctx, medicationID, status); err != nil {
		return nil, err
	}
	return s.medicationRepo.FindByID(ctx, medicationID)
}

func (s *BiopsychosocialService) SuspendMedication(ctx context.Context, medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(ctx, medicationID, patient.MedicationStatusSuspended)
}

func (s *BiopsychosocialService) ActivateMedication(ctx context.Context, medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(ctx, medicationID, patient.MedicationStatusActive)
}

func (s *BiopsychosocialService) FinishMedication(ctx context.Context, medicationID string) (*patient.Medication, error) {
	return s.UpdateMedicationStatus(ctx, medicationID, patient.MedicationStatusFinished)
}

func (s *BiopsychosocialService) RecordVitals(ctx context.Context, patientID string, date time.Time, sleepHours *float64, appetiteLevel *int, weight *float64, physicalActivity int, notes string) (*patient.Vitals, error) {
	vitals, err := patient.NewVitals(patientID, date, sleepHours, appetiteLevel, weight, physicalActivity, notes)
	if err != nil {
		return nil, err
	}
	if err := s.vitalsRepo.Save(ctx, vitals); err != nil {
		return nil, err
	}
	return vitals, nil
}

func (s *BiopsychosocialService) GetLatestVitals(ctx context.Context, patientID string) (*patient.Vitals, error) {
	return s.vitalsRepo.GetLatestVitals(ctx, patientID)
}

func (s *BiopsychosocialService) GetVitalsHistory(ctx context.Context, patientID string, limit int) ([]*patient.Vitals, error) {
	return s.vitalsRepo.FindByPatientID(ctx, patientID, limit)
}

func (s *BiopsychosocialService) GetAverageVitals(ctx context.Context, patientID string, days int) (*patient.VitalsAverage, error) {
	return s.vitalsRepo.GetAverageVitals(ctx, patientID, days)
}

type BiopsychosocialContext struct {
	PatientID         string
	ActiveMedications []*patient.Medication
	AllMedications    []*patient.Medication
	LatestVitals      *patient.Vitals
	VitalsAverage     *patient.VitalsAverage
}

func (s *BiopsychosocialService) GetContext(ctx context.Context, patientID string) (*BiopsychosocialContext, error) {
	activeMeds, err := s.GetActiveMedications(ctx, patientID)
	if err != nil {
		return nil, err
	}
	allMeds, err := s.GetMedications(ctx, patientID)
	if err != nil {
		return nil, err
	}
	latestVitals, err := s.GetLatestVitals(ctx, patientID)
	if err != nil {
		return nil, err
	}
	avgVitals, err := s.GetAverageVitals(ctx, patientID, 30)
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
