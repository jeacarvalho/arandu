package sqlite

import (
	"context"
	"database/sql"
	"time"

	"arandu/internal/domain/patient"
)

type MedicationRepository struct {
	db *DB
}

func NewMedicationRepository(db *DB) *MedicationRepository {
	return &MedicationRepository{db: db}
}

func (r *MedicationRepository) Save(ctx context.Context, m *patient.Medication) error {
	query := `INSERT INTO patient_medications (id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, m.ID, m.PatientID, m.Name, m.Dosage, m.Frequency, m.Prescriber, m.Status, m.StartedAt, m.EndedAt, m.CreatedAt, m.UpdatedAt)
	return err
}

func (r *MedicationRepository) FindByID(ctx context.Context, id string) (*patient.Medication, error) {
	query := `SELECT id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at 
			  FROM patient_medications WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var m patient.Medication
	var endedAt sql.NullTime
	err := row.Scan(&m.ID, &m.PatientID, &m.Name, &m.Dosage, &m.Frequency, &m.Prescriber, &m.Status, &m.StartedAt, &endedAt, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if endedAt.Valid {
		m.EndedAt = &endedAt.Time
	}
	return &m, nil
}

func (r *MedicationRepository) FindByPatientID(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	query := `SELECT id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at 
			  FROM patient_medications WHERE patient_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMedications(rows)
}

func (r *MedicationRepository) GetActiveMedications(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	query := `SELECT id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at 
			  FROM patient_medications WHERE patient_id = ? AND status = 'active' ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMedications(rows)
}

func (r *MedicationRepository) GetMedicationsByStatus(ctx context.Context, patientID string, status patient.MedicationStatus) ([]*patient.Medication, error) {
	query := `SELECT id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at 
			  FROM patient_medications WHERE patient_id = ? AND status = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMedications(rows)
}

func (r *MedicationRepository) UpdateStatus(ctx context.Context, id string, status patient.MedicationStatus) error {
	var endedAt *time.Time
	if status == patient.MedicationStatusSuspended || status == patient.MedicationStatusFinished {
		now := time.Now()
		endedAt = &now
	}

	query := `UPDATE patient_medications SET status = ?, ended_at = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, endedAt, time.Now(), id)
	return err
}

func (r *MedicationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM patient_medications WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MedicationRepository) scanMedications(rows *sql.Rows) ([]*patient.Medication, error) {
	var medications []*patient.Medication
	for rows.Next() {
		var m patient.Medication
		var endedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.PatientID, &m.Name, &m.Dosage, &m.Frequency, &m.Prescriber, &m.Status, &m.StartedAt, &endedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if endedAt.Valid {
			m.EndedAt = &endedAt.Time
		}
		medications = append(medications, &m)
	}
	return medications, nil
}

// FindByPatientIDAndTimeframe busca medicações de um paciente dentro de um período
func (r *MedicationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Medication, error) {
	query := `SELECT id, patient_id, name, dosage, frequency, prescriber, status, started_at, ended_at, created_at, updated_at 
			  FROM patient_medications WHERE patient_id = ?`

	var args []interface{}
	args = append(args, patientID)

	if !startTime.IsZero() {
		query += " AND started_at >= ?"
		args = append(args, startTime)
	}

	query += " ORDER BY started_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMedications(rows)
}
