package sqlite

import (
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

func (r *InsightRepository) Save(i *insight.Insight) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()

	query := `INSERT INTO insights (id, content, source, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, i.ID, i.Content, i.Source, i.CreatedAt)
	return err
}

func (r *InsightRepository) FindByID(id string) (*insight.Insight, error) {
	query := `SELECT id, content, source, created_at FROM insights WHERE id = ?`
	row := r.db.QueryRow(query, id)

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

func (r *InsightRepository) FindAll() ([]*insight.Insight, error) {
	query := `SELECT id, content, source, created_at FROM insights ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
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

func (r *InsightRepository) Delete(id string) error {
	query := `DELETE FROM insights WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *InsightRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS insights (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		source TEXT NOT NULL,
		created_at DATETIME NOT NULL
	)
	`
	_, err := r.db.Exec(query)
	return err
}
