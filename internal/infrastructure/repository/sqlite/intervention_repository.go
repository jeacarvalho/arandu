package sqlite

import (
	"database/sql"
	"time"

	"arandu/internal/domain/intervention"
	"github.com/google/uuid"
)

type InterventionRepository struct {
	db *DB
}

func NewInterventionRepository(db *DB) *InterventionRepository {
	return &InterventionRepository{db: db}
}

func (r *InterventionRepository) Save(i *intervention.Intervention) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()

	query := `INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, i.ID, i.SessionID, i.Content, i.CreatedAt, i.UpdatedAt)
	return err
}

func (r *InterventionRepository) FindByID(id string) (*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var i intervention.Intervention
	err := row.Scan(&i.ID, &i.SessionID, &i.Content, &i.CreatedAt, &i.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *InterventionRepository) FindBySessionID(sessionID string) ([]*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions WHERE session_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interventions []*intervention.Intervention
	for rows.Next() {
		var i intervention.Intervention
		if err := rows.Scan(&i.ID, &i.SessionID, &i.Content, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		interventions = append(interventions, &i)
	}
	return interventions, nil
}

func (r *InterventionRepository) FindAll() ([]*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interventions []*intervention.Intervention
	for rows.Next() {
		var i intervention.Intervention
		if err := rows.Scan(&i.ID, &i.SessionID, &i.Content, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		interventions = append(interventions, &i)
	}
	return interventions, nil
}

func (r *InterventionRepository) Update(i *intervention.Intervention) error {
	i.UpdatedAt = time.Now()
	query := `UPDATE interventions SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, i.Content, i.UpdatedAt, i.ID)
	return err
}

func (r *InterventionRepository) Delete(id string) error {
	query := `DELETE FROM interventions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *InterventionRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS interventions (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	)
	`
	_, err := r.db.Exec(query)
	return err
}
