package sqlite

import (
	"context"
	"database/sql"
	"time"

	"arandu/internal/domain/patient"

	"github.com/google/uuid"
)

type GoalRepository struct {
	db *DB
}

func NewGoalRepository(db *DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Save(ctx context.Context, goal *patient.TherapeuticGoal) error {
	query := `INSERT INTO therapeutic_goals (id, patient_id, title, description, status, closure_note, closed_at, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		goal.ID,
		goal.PatientID,
		goal.Title,
		goal.Description,
		string(goal.Status),
		goal.ClosureNote,
		goal.ClosedAt,
		goal.CreatedAt,
		goal.UpdatedAt,
	)
	return err
}

func (r *GoalRepository) Create(ctx context.Context, patientID, title, description string) (*patient.TherapeuticGoal, error) {
	goal := &patient.TherapeuticGoal{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		Title:       title,
		Description: description,
		Status:      patient.GoalStatusInProgress,
	}
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = goal.CreatedAt

	if err := r.Save(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

func (r *GoalRepository) FindByID(ctx context.Context, id string) (*patient.TherapeuticGoal, error) {
	query := `SELECT id, patient_id, title, description, status, closure_note, closed_at, created_at, updated_at 
			  FROM therapeutic_goals WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanGoal(row)
}

func (r *GoalRepository) FindByPatientID(ctx context.Context, patientID string) ([]*patient.TherapeuticGoal, error) {
	query := `SELECT id, patient_id, title, description, status, closure_note, closed_at, created_at, updated_at 
			  FROM therapeutic_goals WHERE patient_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanGoals(rows)
}

func (r *GoalRepository) GetActiveGoals(ctx context.Context, patientID string) ([]*patient.TherapeuticGoal, error) {
	query := `SELECT id, patient_id, title, description, status, closure_note, closed_at, created_at, updated_at 
			  FROM therapeutic_goals WHERE patient_id = ? AND status = 'in_progress' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanGoals(rows)
}

func (r *GoalRepository) GetGoalsByStatus(ctx context.Context, patientID string, status patient.GoalStatus) ([]*patient.TherapeuticGoal, error) {
	query := `SELECT id, patient_id, title, description, status, closure_note, closed_at, created_at, updated_at 
			  FROM therapeutic_goals WHERE patient_id = ? AND status = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, patientID, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanGoals(rows)
}

func (r *GoalRepository) UpdateStatus(ctx context.Context, id string, status patient.GoalStatus) error {
	query := `UPDATE therapeutic_goals SET status = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(status), id)
	return err
}

func (r *GoalRepository) CloseWithNote(ctx context.Context, id string, status patient.GoalStatus, closureNote string) error {
	query := `UPDATE therapeutic_goals SET status = ?, closure_note = ?, closed_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(status), closureNote, id)
	return err
}

func (r *GoalRepository) Update(ctx context.Context, goal *patient.TherapeuticGoal) error {
	query := `UPDATE therapeutic_goals SET title = ?, description = ?, status = ?, closure_note = ?, closed_at = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		goal.Title,
		goal.Description,
		string(goal.Status),
		goal.ClosureNote,
		goal.ClosedAt,
		goal.UpdatedAt,
		goal.ID,
	)
	return err
}

func (r *GoalRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM therapeutic_goals WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *GoalRepository) scanGoal(row *sql.Row) (*patient.TherapeuticGoal, error) {
	var goal patient.TherapeuticGoal
	var description, closureNote sql.NullString
	var closedAt sql.NullTime
	err := row.Scan(
		&goal.ID,
		&goal.PatientID,
		&goal.Title,
		&description,
		&goal.Status,
		&closureNote,
		&closedAt,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		goal.Description = description.String
	}
	if closureNote.Valid {
		goal.ClosureNote = closureNote.String
	}
	if closedAt.Valid {
		goal.ClosedAt = &closedAt.Time
	}
	return &goal, nil
}

func (r *GoalRepository) scanGoals(rows *sql.Rows) ([]*patient.TherapeuticGoal, error) {
	var goals []*patient.TherapeuticGoal
	for rows.Next() {
		var goal patient.TherapeuticGoal
		var description, closureNote sql.NullString
		var closedAt sql.NullTime
		err := rows.Scan(
			&goal.ID,
			&goal.PatientID,
			&goal.Title,
			&description,
			&goal.Status,
			&closureNote,
			&closedAt,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			goal.Description = description.String
		}
		if closureNote.Valid {
			goal.ClosureNote = closureNote.String
		}
		if closedAt.Valid {
			goal.ClosedAt = &closedAt.Time
		}
		goals = append(goals, &goal)
	}
	return goals, rows.Err()
}
