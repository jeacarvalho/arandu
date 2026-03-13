package sqlite

import (
	"database/sql"
	"time"

	"arandu/internal/domain/session"
	"github.com/google/uuid"
)

type SessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(s *session.Session) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	query := `INSERT INTO sessions (id, patient_id, date, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.PatientID, s.Date, s.Notes, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SessionRepository) FindByID(id string) (*session.Session, error) {
	query := `SELECT id, patient_id, date, notes, created_at, updated_at FROM sessions WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var s session.Session
	err := row.Scan(&s.ID, &s.PatientID, &s.Date, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) FindByPatientID(patientID string) ([]*session.Session, error) {
	query := `SELECT id, patient_id, date, notes, created_at, updated_at FROM sessions WHERE patient_id = ? ORDER BY date DESC`
	rows, err := r.db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		if err := rows.Scan(&s.ID, &s.PatientID, &s.Date, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (r *SessionRepository) FindAll() ([]*session.Session, error) {
	query := `SELECT id, patient_id, date, notes, created_at, updated_at FROM sessions ORDER BY date DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		if err := rows.Scan(&s.ID, &s.PatientID, &s.Date, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (r *SessionRepository) Update(s *session.Session) error {
	s.UpdatedAt = time.Now()
	query := `UPDATE sessions SET patient_id = ?, date = ?, notes = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, s.PatientID, s.Date, s.Notes, s.UpdatedAt, s.ID)
	return err
}

func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SessionRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		patient_id TEXT NOT NULL,
		date DATETIME NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
	)
	`
	_, err := r.db.Exec(query)
	return err
}
