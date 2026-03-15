package session

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	ListByPatient(ctx context.Context, patientID string) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
}
