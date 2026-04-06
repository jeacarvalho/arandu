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

func (r *InterventionRepository) Save(ctx context.Context, i *intervention.Intervention) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()

	query := `INSERT INTO interventions (id, session_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, i.ID, i.SessionID, i.Content, i.CreatedAt, i.UpdatedAt)
	return err
}

func (r *InterventionRepository) FindByID(ctx context.Context, id string) (*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

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

func (r *InterventionRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions WHERE session_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
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

func (r *InterventionRepository) FindAll(ctx context.Context) ([]*intervention.Intervention, error) {
	query := `SELECT id, session_id, content, created_at, updated_at FROM interventions ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *InterventionRepository) Update(ctx context.Context, i *intervention.Intervention) error {
	i.UpdatedAt = time.Now()
	query := `UPDATE interventions SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, i.Content, i.UpdatedAt, i.ID)
	return err
}

func (r *InterventionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM interventions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
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

// ============================================================================
// INTERVENTION CLASSIFICATION METHODS
// ============================================================================

// AddTagToIntervention adiciona uma tag a uma intervenção
func (r *InterventionRepository) AddTagToIntervention(ctx context.Context, interventionID, tagID string, intensity int) error {
	id := uuid.New().String()
	now := time.Now()

	query := `INSERT INTO intervention_classifications (id, intervention_id, tag_id, intensity, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, id, interventionID, tagID, intensity, now)
	return err
}

// RemoveTagFromIntervention remove uma tag de uma intervenção
func (r *InterventionRepository) RemoveTagFromIntervention(ctx context.Context, interventionID, tagID string) error {
	query := `DELETE FROM intervention_classifications WHERE intervention_id = ? AND tag_id = ?`
	_, err := r.db.ExecContext(ctx, query, interventionID, tagID)
	return err
}

// GetInterventionTags retorna todas as tags de uma intervenção
func (r *InterventionRepository) GetInterventionTags(ctx context.Context, interventionID string) ([]*intervention.InterventionClassification, error) {
	query := `
		SELECT ic.id, ic.intervention_id, ic.tag_id, ic.intensity, ic.created_at,
		       t.id, t.name, t.tag_type, t.color, t.icon
		FROM intervention_classifications ic
		JOIN intervention_tags t ON ic.tag_id = t.id
		WHERE ic.intervention_id = ?
		ORDER BY t.sort_order ASC, t.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, interventionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classifications []*intervention.InterventionClassification
	for rows.Next() {
		var ic intervention.InterventionClassification
		var tag intervention.Tag
		err := rows.Scan(
			&ic.ID, &ic.InterventionID, &ic.TagID, &ic.Intensity, &ic.CreatedAt,
			&tag.ID, &tag.Name, &tag.TagType, &tag.Color, &tag.Icon,
		)
		if err != nil {
			return nil, err
		}
		ic.Tag = &tag
		classifications = append(classifications, &ic)
	}
	return classifications, nil
}

// GetAllInterventionTags retorna todas as tags predefinidas
func (r *InterventionRepository) GetAllInterventionTags(ctx context.Context) ([]*intervention.Tag, error) {
	query := `
		SELECT id, name, tag_type, color, icon, sort_order, created_at
		FROM intervention_tags
		ORDER BY sort_order ASC, name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*intervention.Tag
	for rows.Next() {
		var tag intervention.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.TagType, &tag.Color, &tag.Icon, &tag.SortOrder, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

// GetAllTagsByType retorna todas as tags predefinidas de um tipo específico
func (r *InterventionRepository) GetAllTagsByType(ctx context.Context, tagType intervention.TagType) ([]*intervention.Tag, error) {
	query := `
		SELECT id, name, tag_type, color, icon, sort_order, created_at
		FROM intervention_tags
		WHERE tag_type = ?
		ORDER BY sort_order ASC, name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, string(tagType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*intervention.Tag
	for rows.Next() {
		var tag intervention.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.TagType, &tag.Color, &tag.Icon, &tag.SortOrder, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

// GetInterventionTagsGroupedByType retorna as tags de uma intervenção agrupadas por tipo
func (r *InterventionRepository) GetInterventionTagsGroupedByType(ctx context.Context, interventionID string) (map[intervention.TagType][]*intervention.InterventionClassification, error) {
	query := `
		SELECT ic.id, ic.intervention_id, ic.tag_id, ic.intensity, ic.created_at,
		       t.id, t.name, t.tag_type, t.color, t.icon
		FROM intervention_classifications ic
		JOIN intervention_tags t ON ic.tag_id = t.id
		WHERE ic.intervention_id = ?
		ORDER BY t.tag_type ASC, t.sort_order ASC, t.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, interventionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[intervention.TagType][]*intervention.InterventionClassification)
	for rows.Next() {
		var ic intervention.InterventionClassification
		var tag intervention.Tag
		err := rows.Scan(
			&ic.ID, &ic.InterventionID, &ic.TagID, &ic.Intensity, &ic.CreatedAt,
			&tag.ID, &tag.Name, &tag.TagType, &tag.Color, &tag.Icon,
		)
		if err != nil {
			return nil, err
		}
		ic.Tag = &tag
		grouped[tag.TagType] = append(grouped[tag.TagType], &ic)
	}
	return grouped, nil
}

// FindInterventionsByTagID retorna todas as intervenções com uma determinada tag
func (r *InterventionRepository) FindInterventionsByTagID(ctx context.Context, tagID string) ([]*intervention.Intervention, error) {
	query := `
		SELECT i.id, i.session_id, i.content, i.created_at, i.updated_at
		FROM interventions i
		JOIN intervention_classifications ic ON i.id = ic.intervention_id
		WHERE ic.tag_id = ?
		ORDER BY i.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, tagID)
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

// GetTagCountByType retorna a contagem de tags por tipo para uma intervenção
func (r *InterventionRepository) GetTagCountByType(ctx context.Context, interventionID string) (map[intervention.TagType]int, error) {
	query := `
		SELECT t.tag_type, COUNT(*) as count
		FROM intervention_classifications ic
		JOIN intervention_tags t ON ic.tag_id = t.id
		WHERE ic.intervention_id = ?
		GROUP BY t.tag_type
	`
	rows, err := r.db.QueryContext(ctx, query, interventionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[intervention.TagType]int)
	for rows.Next() {
		var tagType intervention.TagType
		var count int
		if err := rows.Scan(&tagType, &count); err != nil {
			return nil, err
		}
		counts[tagType] = count
	}
	return counts, nil
}

// GetTopInterventionTags retorna as tags mais utilizadas em intervenções
func (r *InterventionRepository) GetTopInterventionTags(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := `
		SELECT t.id, t.name, t.tag_type, t.color, COUNT(ic.id) as usage_count
		FROM intervention_tags t
		LEFT JOIN intervention_classifications ic ON t.id = ic.tag_id
		GROUP BY t.id, t.name, t.tag_type, t.color
		ORDER BY usage_count DESC
		LIMIT ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, name, color string
		var tagType intervention.TagType
		var count int
		if err := rows.Scan(&id, &name, &tagType, &color, &count); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":    id,
			"name":  name,
			"type":  tagType,
			"color": color,
			"count": count,
		})
	}
	return results, nil
}
