package sqlite

import (
	"context"
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

func (r *ObservationRepository) Save(ctx context.Context, o *observation.Observation) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	o.CreatedAt = time.Now()

	query := `INSERT INTO observations (id, session_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, NULL)`
	_, err := r.db.ExecContext(ctx, query, o.ID, o.SessionID, o.Content, o.CreatedAt)
	return err
}

func (r *ObservationRepository) FindByID(ctx context.Context, id string) (*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

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

func (r *ObservationRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations WHERE session_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
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

func (r *ObservationRepository) FindAll(ctx context.Context) ([]*observation.Observation, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM observations ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *ObservationRepository) Update(ctx context.Context, o *observation.Observation) error {
	query := `UPDATE observations SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, o.Content, time.Now(), o.ID)
	return err
}

func (r *ObservationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM observations WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
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

// FindByPatientIDAndTimeframe busca observações de um paciente dentro de um período
func (r *ObservationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*observation.Observation, error) {
	query := `
		SELECT o.id, o.session_id, o.content, o.created_at, o.updated_at 
		FROM observations o
		JOIN sessions s ON o.session_id = s.id
		WHERE s.patient_id = ?
	`

	var args []interface{}
	args = append(args, patientID)

	if !startTime.IsZero() {
		query += " AND o.created_at >= ?"
		args = append(args, startTime)
	}

	query += " ORDER BY o.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
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

// InitSchema is deprecated - use migrations instead
func (r *ObservationRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}

// GetTags retrieves all available tags
func (r *ObservationRepository) GetTags(ctx context.Context) ([]observation.Tag, error) {
	query := `SELECT id, name, tag_type, color, sort_order, created_at FROM tags ORDER BY tag_type, sort_order`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []observation.Tag
	for rows.Next() {
		var t observation.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.TagType, &t.Color, &t.SortOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// GetTagsByType retrieves tags filtered by type
func (r *ObservationRepository) GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error) {
	query := `SELECT id, name, tag_type, color, sort_order, created_at FROM tags WHERE tag_type = ? ORDER BY sort_order`
	rows, err := r.db.QueryContext(ctx, query, string(tagType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []observation.Tag
	for rows.Next() {
		var t observation.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.TagType, &t.Color, &t.SortOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// AddTagToObservation adds a tag to an observation with intensity
func (r *ObservationRepository) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error {
	id := uuid.New().String()
	query := `INSERT INTO observation_tags (id, observation_id, tag_id, intensity, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, id, observationID, tagID, intensity, time.Now())
	return err
}

// RemoveTagFromObservation removes a tag from an observation
func (r *ObservationRepository) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error {
	query := `DELETE FROM observation_tags WHERE observation_id = ? AND tag_id = ?`
	_, err := r.db.ExecContext(ctx, query, observationID, tagID)
	return err
}

// GetObservationTags retrieves all tags associated with an observation
func (r *ObservationRepository) GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error) {
	query := `
		SELECT ot.id, ot.observation_id, ot.tag_id, ot.intensity, ot.created_at,
		       t.id, t.name, t.tag_type, t.color
		FROM observation_tags ot
		JOIN tags t ON ot.tag_id = t.id
		WHERE ot.observation_id = ?
		ORDER BY t.tag_type, t.sort_order`

	rows, err := r.db.QueryContext(ctx, query, observationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observationTags []observation.ObservationTag
	for rows.Next() {
		var ot observation.ObservationTag
		ot.Tag = &observation.Tag{}
		if err := rows.Scan(&ot.ID, &ot.ObservationID, &ot.TagID, &ot.Intensity, &ot.CreatedAt,
			&ot.Tag.ID, &ot.Tag.Name, &ot.Tag.TagType, &ot.Tag.Color); err != nil {
			return nil, err
		}
		observationTags = append(observationTags, ot)
	}
	return observationTags, nil
}

// GetTagsSummary returns a summary of all tags usage
func (r *ObservationRepository) GetTagsSummary(ctx context.Context) ([]observation.TagSummary, error) {
	query := `
		SELECT t.tag_type, COUNT(*) as count
		FROM observation_tags ot
		JOIN tags t ON ot.tag_id = t.id
		GROUP BY t.tag_type
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []observation.TagSummary
	for rows.Next() {
		var s observation.TagSummary
		var tagType string
		if err := rows.Scan(&tagType, &s.Count); err != nil {
			return nil, err
		}
		s.TagType = observation.TagType(tagType)
		summaries = append(summaries, s)
	}
	return summaries, nil
}

// GetTagsSummaryByPatient returns tag summary for a specific patient
func (r *ObservationRepository) GetTagsSummaryByPatient(ctx context.Context, patientID string) ([]observation.TagSummary, error) {
	query := `
		SELECT t.tag_type, COUNT(*) as count
		FROM observation_tags ot
		JOIN tags t ON ot.tag_id = t.id
		JOIN observations o ON ot.observation_id = o.id
		JOIN sessions s ON o.session_id = s.id
		WHERE s.patient_id = ?
		GROUP BY t.tag_type
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []observation.TagSummary
	for rows.Next() {
		var s observation.TagSummary
		var tagType string
		if err := rows.Scan(&tagType, &s.Count); err != nil {
			return nil, err
		}
		s.TagType = observation.TagType(tagType)
		summaries = append(summaries, s)
	}
	return summaries, nil
}

// FindByTag retrieves observations that have a specific tag
func (r *ObservationRepository) FindByTag(ctx context.Context, tagID string) ([]*observation.Observation, error) {
	query := `
		SELECT o.id, o.session_id, o.content, o.created_at, o.updated_at
		FROM observations o
		JOIN observation_tags ot ON o.id = ot.observation_id
		WHERE ot.tag_id = ?
		ORDER BY o.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, tagID)
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
