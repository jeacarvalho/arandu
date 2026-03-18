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

// SearchFTS busca observações usando FTS5
func (r *ObservationRepository) SearchFTS(query string, limit int) ([]*observation.Observation, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	sqlQuery := `
		SELECT o.id, o.session_id, o.content, o.created_at, o.updated_at 
		FROM observations o
		WHERE o.id IN (
			SELECT rowid FROM observations_fts WHERE content MATCH ?
		)
		ORDER BY o.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(sqlQuery, query, limit)
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

// GetTopTerms retorna os termos mais frequentes nas observações
func (r *ObservationRepository) GetTopTerms(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := `
		SELECT term, count 
		FROM fts5vocabulary('observations_fts', 'col')
		WHERE term NOT IN ('de', 'a', 'o', 'que', 'e', 'do', 'da', 'em', 'um', 'para', 'com', 'os', 'as', 'no', 'na', 'mais', 'como', 'mas', 'ao', 'aos', 'à', 'sua', 'seu', 'das', 'dos')
		ORDER BY count DESC 
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var term string
		var count int
		if err := rows.Scan(&term, &count); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"term":  term,
			"count": count,
		})
	}
	return result, nil
}

// InitSchema is deprecated - use migrations instead
func (r *ObservationRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
