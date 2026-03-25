package context

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoTenantDB  = errors.New("no tenant database in context")
	ErrNoTenantID  = errors.New("no tenant ID in context")
	ErrNoUserID    = errors.New("no user ID in context")
	ErrNoRequestID = errors.New("no request ID in context")
	ErrNoUserEmail = errors.New("no user email in context")
)

type contextKey string

const (
	tenantDBKey  contextKey = "tenant_db"
	tenantIDKey  contextKey = "tenant_id"
	userIDKey    contextKey = "user_id"
	requestIDKey contextKey = "request_id"
	userEmailKey contextKey = "user_email"
)

func WithTenantDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, tenantDBKey, db)
}

func GetTenantDB(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(tenantDBKey).(*sql.DB)
	if !ok || db == nil {
		return nil, ErrNoTenantDB
	}
	return db, nil
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func GetTenantID(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)
	if !ok || tenantID == "" {
		return "", ErrNoTenantID
	}
	return tenantID, nil
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", ErrNoUserID
	}
	return userID, nil
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok || requestID == "" {
		return "", ErrNoRequestID
	}
	return requestID, nil
}

func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailKey, email)
}

func GetUserEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(userEmailKey).(string)
	if !ok || email == "" {
		return "", ErrNoUserEmail
	}
	return email, nil
}
