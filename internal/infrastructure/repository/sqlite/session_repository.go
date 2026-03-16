package sqlite

import (
	"context"
	"database/sql"

	"arandu/internal/domain/session"
)

type SessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *session.Session) error {
	query := `INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.PatientID, s.Date, s.Summary, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	query := `SELECT id, patient_id, date, summary, created_at, updated_at FROM sessions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var s session.Session
	err := row.Scan(&s.ID, &s.PatientID, &s.Date, &s.Summary, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) ListByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	query := `SELECT id, patient_id, date, summary, created_at, updated_at FROM sessions WHERE patient_id = ? ORDER BY date DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*session.Session
	for rows.Next() {
		var s session.Session
		if err := rows.Scan(&s.ID, &s.PatientID, &s.Date, &s.Summary, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, s *session.Session) error {
	query := `UPDATE sessions SET date = ?, summary = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, s.Date, s.Summary, s.UpdatedAt, s.ID)
	return err
}

// InitSchema is deprecated - use migrations instead
func (r *SessionRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
