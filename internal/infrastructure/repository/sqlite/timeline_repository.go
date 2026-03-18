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

	query := `
		SELECT 
			'session' as type,
			id,
			date as event_date,
			summary as content,
			created_at,
			json_object('session_id', id, 'patient_id', patient_id) as metadata
		FROM sessions 
		WHERE patient_id = ?
		
		UNION ALL
		
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

	rows, err := r.db.QueryContext(ctx, query, patientID, patientID, patientID, limit, offset)
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

		if filterType != nil && eventType != *filterType {
			continue
		}

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
