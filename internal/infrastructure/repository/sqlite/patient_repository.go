package sqlite

import (
	"database/sql"
	"time"

	"arandu/internal/domain/patient"
	"github.com/google/uuid"
)

type PatientRepository struct {
	db *DB
}

func NewPatientRepository(db *DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Save(p *patient.Patient) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	query := `INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, p.ID, p.Name, p.Notes, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PatientRepository) FindByID(id string) (*patient.Patient, error) {
	query := `SELECT id, name, notes, created_at, updated_at FROM patients WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var p patient.Patient
	err := row.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PatientRepository) FindAll() ([]*patient.Patient, error) {
	query := `SELECT id, name, notes, created_at, updated_at FROM patients ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*patient.Patient
	for rows.Next() {
		var p patient.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}
	return patients, nil
}

func (r *PatientRepository) Update(p *patient.Patient) error {
	p.UpdatedAt = time.Now()
	query := `UPDATE patients SET name = ?, notes = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, p.Name, p.Notes, p.UpdatedAt, p.ID)
	return err
}

func (r *PatientRepository) Delete(id string) error {
	query := `DELETE FROM patients WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PatientRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS patients (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`
	_, err := r.db.Exec(query)
	return err
}
