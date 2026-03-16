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

	query := `INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, NULL)`
	_, err := r.db.Exec(query, o.ID, o.SessionID, o.Content, o.CreatedAt)
	return err
}

func (r *ObservationRepository) FindByID(id string) (*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var o observation.Observation
	var updatedAt sql.NullTime
	err := row.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if updatedAt.Valid {
		o.UpdatedAt = updatedAt.Time
	}
	return &o, nil
}

func (r *ObservationRepository) FindBySessionID(sessionID string) ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations WHERE session_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []*observation.Observation
	for rows.Next() {
		var o observation.Observation
		var updatedAt sql.NullTime
		if err := rows.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			o.UpdatedAt = updatedAt.Time
		}
		observations = append(observations, &o)
	}
	return observations, nil
}

func (r *ObservationRepository) FindAll() ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []*observation.Observation
	for rows.Next() {
		var o observation.Observation
		var updatedAt sql.NullTime
		if err := rows.Scan(&o.ID, &o.SessionID, &o.Content, &o.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			o.UpdatedAt = updatedAt.Time
		}
		observations = append(observations, &o)
	}
	return observations, nil
}

func (r *ObservationRepository) Update(o *observation.Observation) error {
	query := `UPDATE observations SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, o.Content, time.Now(), o.ID)
	return err
}

func (r *ObservationRepository) Delete(id string) error {
	query := `DELETE FROM observations WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// InitSchema is deprecated - use migrations instead
func (r *ObservationRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
