package sqlite

import (
	"context"
	"time"

	"arandu/internal/domain/timeline"
)

type TimelineRepository struct {
	db *DB
}

func NewTimelineRepository(db *DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

func (r *TimelineRepository) GetTimelineByPatientID(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	var events timeline.Timeline

	// Construir query dinamicamente baseada no filtro
	var query string
	var args []interface{}

	if filterType == nil {
		// Sem filtro - buscar observações e intervenções (não inclui sessões)
		query = `
			SELECT 
				'observation' as type,
				o.id,
				o.created_at as event_date,
				o.content,
				o.created_at,
				json_object('observation_id', o.id, 'session_id', o.session_id) as metadata
			FROM observations o
			INNER JOIN sessions s ON o.session_id = s.id
			WHERE s.patient_id = ?
			
			UNION ALL
			
			SELECT 
				'intervention' as type,
				i.id,
				i.created_at as event_date,
				i.content,
				i.created_at,
				json_object('intervention_id', i.id, 'session_id', i.session_id) as metadata
			FROM interventions i
			INNER JOIN sessions s ON i.session_id = s.id
			WHERE s.patient_id = ?
			
			ORDER BY event_date DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{patientID, patientID, limit, offset}
	} else {
		// Com filtro - buscar apenas um tipo
		switch *filterType {
		case timeline.EventTypeObservation:
			query = `
				SELECT 
					'observation' as type,
					o.id,
					o.created_at as event_date,
					o.content,
					o.created_at,
					json_object('observation_id', o.id, 'session_id', o.session_id) as metadata
				FROM observations o
				INNER JOIN sessions s ON o.session_id = s.id
				WHERE s.patient_id = ?
				ORDER BY o.created_at DESC
				LIMIT ? OFFSET ?
			`
			args = []interface{}{patientID, limit, offset}

		case timeline.EventTypeIntervention:
			query = `
				SELECT 
					'intervention' as type,
					i.id,
					i.created_at as event_date,
					i.content,
					i.created_at,
					json_object('intervention_id', i.id, 'session_id', i.session_id) as metadata
				FROM interventions i
				INNER JOIN sessions s ON i.session_id = s.id
				WHERE s.patient_id = ?
				ORDER BY i.created_at DESC
				LIMIT ? OFFSET ?
			`
			args = []interface{}{patientID, limit, offset}

		default:
			// Fallback para sem filtro (apenas observações e intervenções)
			query = `
				SELECT 
					'observation' as type,
					o.id,
					o.created_at as event_date,
					o.content,
					o.created_at,
					json_object('observation_id', o.id, 'session_id', o.session_id) as metadata
				FROM observations o
				INNER JOIN sessions s ON o.session_id = s.id
				WHERE s.patient_id = ?
				
				UNION ALL
				
				SELECT 
					'intervention' as type,
					i.id,
					i.created_at as event_date,
					i.content,
					i.created_at,
					json_object('intervention_id', i.id, 'session_id', i.session_id) as metadata
				FROM interventions i
				INNER JOIN sessions s ON i.session_id = s.id
				WHERE s.patient_id = ?
				
				ORDER BY event_date DESC
				LIMIT ? OFFSET ?
			`
			args = []interface{}{patientID, patientID, limit, offset}
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventTypeStr, id, content string
		var eventDate, createdAt time.Time
		var metadataJSON string

		err := rows.Scan(&eventTypeStr, &id, &eventDate, &content, &createdAt, &metadataJSON)
		if err != nil {
			return nil, err
		}

		eventType := timeline.EventType(eventTypeStr)

		metadata := make(map[string]string)
		if metadataJSON != "" {
			metadata = parseMetadataJSON(metadataJSON)
		}

		event := &timeline.TimelineEvent{
			ID:        id,
			Type:      eventType,
			Date:      eventDate,
			Content:   content,
			Metadata:  metadata,
			CreatedAt: createdAt,
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *TimelineRepository) GetTimelineByPatientIDWithFilter(ctx context.Context, patientID string, filterType timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	return r.GetTimelineByPatientID(ctx, patientID, &filterType, limit, offset)
}

func (r *TimelineRepository) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	if query == "" {
		return nil, nil
	}

	searchQuery := `
		SELECT 
			'observation' as type,
			o.id,
			o.created_at as date,
			o.content,
			snippet(observations_fts, 1, '<b>', '</b>', '...', 64) as snippet,
			s.id as session_id,
			s.patient_id
		FROM observations_fts fts
		JOIN observations o ON fts.source_id = o.id
		JOIN sessions s ON o.session_id = s.id
		WHERE fts.content MATCH ? AND s.patient_id = ?
		
		UNION ALL
		
		SELECT 
			'intervention' as type,
			i.id,
			i.created_at as date,
			i.content,
			snippet(interventions_fts, 1, '<b>', '</b>', '...', 64) as snippet,
			s.id as session_id,
			s.patient_id
		FROM interventions_fts fts
		JOIN interventions i ON fts.source_id = i.id
		JOIN sessions s ON i.session_id = s.id
		WHERE fts.content MATCH ? AND s.patient_id = ?
		
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, searchQuery, query, patientID, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*timeline.SearchResult
	for rows.Next() {
		var result timeline.SearchResult
		var eventTypeStr string

		err := rows.Scan(&eventTypeStr, &result.ID, &result.Date, &result.Content, &result.Snippet, &result.SessionID, &result.PatientID)
		if err != nil {
			return nil, err
		}

		result.Type = timeline.EventType(eventTypeStr)
		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func parseMetadataJSON(jsonStr string) map[string]string {
	result := make(map[string]string)

	if jsonStr == "" || jsonStr == "{}" {
		return result
	}

	trimmed := jsonStr[1 : len(jsonStr)-1]
	pairs := splitJSONPairs(trimmed)

	for _, pair := range pairs {
		key, value := splitKeyValue(pair)
		if key != "" && value != "" {
			result[key] = value
		}
	}

	return result
}

func splitJSONPairs(str string) []string {
	var pairs []string
	var current string
	var inQuotes bool

	for i := 0; i < len(str); i++ {
		ch := str[i]

		if ch == '"' {
			inQuotes = !inQuotes
		}

		if ch == ',' && !inQuotes {
			pairs = append(pairs, current)
			current = ""
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		pairs = append(pairs, current)
	}

	return pairs
}

func splitKeyValue(pair string) (string, string) {
	parts := make([]string, 0)
	var current string
	var inQuotes bool

	for i := 0; i < len(pair); i++ {
		ch := pair[i]

		if ch == '"' {
			inQuotes = !inQuotes
		}

		if ch == ':' && !inQuotes {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	if len(parts) != 2 {
		return "", ""
	}

	key := trimQuotes(parts[0])
	value := trimQuotes(parts[1])

	return key, value
}

func trimQuotes(str string) string {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}
