package sqlite

import (
	"database/sql"
	"time"

	"arandu/internal/domain/observation"
	"github.com/google/uuid"
)

type ObservationRepository struct {
	db *DB
}

func NewObservationRepository(db *DB) *ObservationRepository {
	return &ObservationRepository{db: db}
}

func (r *ObservationRepository) Save(o *observation.Observation) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	o.CreatedAt = time.Now()

	query := `INSERT INTO observations (id, session_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, o.ID, o.SessionID, o.Content, o.CreatedAt)
	return err
}

func (r *ObservationRepository) FindByID(id string) (*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at FROM observations WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var o observation.Observation
	err := row.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *ObservationRepository) FindBySessionID(sessionID string) ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at FROM observations WHERE session_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []*observation.Observation
	for rows.Next() {
		var o observation.Observation
		if err := rows.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, &o)
	}
	return observations, nil
}

func (r *ObservationRepository) FindAll() ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at FROM observations ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []*observation.Observation
	for rows.Next() {
		var o observation.Observation
		if err := rows.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, &o)
	}
	return observations, nil
}

func (r *ObservationRepository) Update(o *observation.Observation) error {
	query := `UPDATE observations SET content = ? WHERE id = ?`
	_, err := r.db.Exec(query, o.Content, o.ID)
	return err
}

func (r *ObservationRepository) Delete(id string) error {
	query := `DELETE FROM observations WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ObservationRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS observations (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	)
	`
	_, err := r.db.Exec(query)
	return err
}
