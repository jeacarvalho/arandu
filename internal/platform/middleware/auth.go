package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"
)

const (
	SessionCookieName = "arandu_session"
	TenantIDKey       = "tenant_id"
	TenantDBKey       = "tenant_db"
	UserIDKey         = "user_id"
)

type Session struct {
	ID        string
	UserID    string
	TenantID  string
	ExpiresAt time.Time
}

type AuthMiddleware struct {
	centralDB *sqlite.CentralDB
	pool      *sqlite.TenantPool
}

func NewAuthMiddleware(centralDB *sqlite.CentralDB, pool *sqlite.TenantPool) *AuthMiddleware {
	return &AuthMiddleware{
		centralDB: centralDB,
		pool:      pool,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRoute(r) {
			next.ServeHTTP(w, r)
			return
		}

		session, err := am.getSession(r)
		if err != nil || session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			am.clearSession(w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		db, err := am.pool.GetConnection(session.TenantID)
		if err != nil {
			am.renderMaintenancePage(w)
			return
		}

		ctx := context.WithValue(r.Context(), TenantIDKey, session.TenantID)
		ctx = context.WithValue(ctx, TenantDBKey, db)
		ctx = context.WithValue(ctx, UserIDKey, session.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicRoute(r *http.Request) bool {
	path := r.URL.Path

	publicExactPaths := []string{
		"/",
		"/login",
		"/logout",
		"/test",
		"/auth/login",
		"/auth/google",
		"/auth/google/callback",
	}

	publicPrefixPaths := []string{
		"/static/",
	}

	for _, p := range publicExactPaths {
		if path == p {
			return true
		}
	}

	for _, p := range publicPrefixPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

func (am *AuthMiddleware) getSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, fmt.Errorf("no session cookie: %w", err)
	}

	if cookie.Value == "" {
		return nil, fmt.Errorf("empty session token")
	}

	session, err := am.validateSessionToken(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}

	return session, nil
}

func (am *AuthMiddleware) validateSessionToken(token string) (*Session, error) {
	var session Session
	var expiresAt int64

	query := `
		SELECT id, user_id, tenant_id, expires_at 
		FROM sessions 
		WHERE id = ? AND expires_at > ?`

	err := am.centralDB.QueryRow(query, token, time.Now().Unix()).Scan(
		&session.ID,
		&session.UserID,
		&session.TenantID,
		&expiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	session.ExpiresAt = time.Unix(expiresAt, 0)
	return &session, nil
}

func (am *AuthMiddleware) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (am *AuthMiddleware) renderMaintenancePage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(`<!DOCTYPE html>
<html><head><title>Manutenção</title></head>
<body style="font-family: system-ui; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0;">
<div style="text-align: center;"><h1>Manutenção Temporária</h1><p>Seu consultório está temporariamente indisponível. Tente novamente em alguns minutos.</p></div>
</body></html>`))
}

func GetTenantID(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	if !ok || tenantID == "" {
		return "", fmt.Errorf("tenant ID not found in context")
	}
	return tenantID, nil
}

func GetTenantDB(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(TenantDBKey).(*sql.DB)
	if db == nil || !ok {
		return nil, fmt.Errorf("tenant database not found in context")
	}
	return db, nil
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
