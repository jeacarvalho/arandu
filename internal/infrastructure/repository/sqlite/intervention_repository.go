package sqlite

import (
	"context"
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

// SearchFTS busca intervenções usando FTS5
func (r *InterventionRepository) SearchFTS(query string, limit int) ([]*intervention.Intervention, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	sqlQuery := `
		SELECT i.id, i.session_id, i.content, i.created_at, i.updated_at 
		FROM interventions i
		WHERE i.id IN (
			SELECT rowid FROM interventions_fts WHERE content MATCH ?
		)
		ORDER BY i.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(sqlQuery, query, limit)
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

// GetTopTerms retorna os termos mais frequentes nas intervenções
func (r *InterventionRepository) GetTopTerms(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := `
		SELECT term, count 
		FROM fts5vocabulary('interventions_fts', 'col')
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

// FindByPatientIDAndTimeframe busca intervenções de um paciente dentro de um período
func (r *InterventionRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*intervention.Intervention, error) {
	query := `
		SELECT i.id, i.session_id, i.content, i.created_at, i.updated_at 
		FROM interventions i
		JOIN sessions s ON i.session_id = s.id
		WHERE s.patient_id = ?
	`

	var args []interface{}
	args = append(args, patientID)

	if !startTime.IsZero() {
		query += " AND i.created_at >= ?"
		args = append(args, startTime)
	}

	query += " ORDER BY i.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
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

// InitSchema is deprecated - use migrations instead
func (r *InterventionRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
