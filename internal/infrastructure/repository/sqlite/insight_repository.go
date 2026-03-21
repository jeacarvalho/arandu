package sqlite

import (
	"context"
	"database/sql"
	"time"

	"arandu/internal/domain/insight"
	"github.com/google/uuid"
)

type InsightRepository struct {
	db *DB
}

func NewInsightRepository(db *DB) *InsightRepository {
	return &InsightRepository{db: db}
}

func (r *InsightRepository) Save(ctx context.Context, i *insight.Insight) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()

	query := `INSERT INTO insights (id, content, source, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, i.ID, i.Content, i.Source, i.CreatedAt)
	return err
}

func (r *InsightRepository) FindByID(ctx context.Context, id string) (*insight.Insight, error) {
	query := `SELECT id, content, source, created_at FROM insights WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var i insight.Insight
	err := row.Scan(&i.ID, &i.Content, &i.Source, &i.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *InsightRepository) FindAll(ctx context.Context) ([]*insight.Insight, error) {
	query := `SELECT id, content, source, created_at FROM insights ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insights []*insight.Insight
	for rows.Next() {
		var i insight.Insight
		if err := rows.Scan(&i.ID, &i.Content, &i.Source, &i.CreatedAt); err != nil {
			return nil, err
		}
		insights = append(insights, &i)
	}
	return insights, nil
}

func (r *InsightRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM insights WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// InitSchema is deprecated - use migrations instead
func (r *InsightRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
