package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"
	appcontext "arandu/internal/platform/context"
)

const (
	SessionCookieName = "arandu_session"
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

		log.Printf("[AUTH DEBUG] Path: %s, Method: %s", r.URL.Path, r.Method)

		session, err := am.getSession(r)
		if err != nil || session == nil {
			log.Printf("[AUTH DEBUG] getSession error: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			log.Printf("[AUTH DEBUG] Session expired: %v", session.ExpiresAt)
			am.clearSession(w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		log.Printf("[AUTH DEBUG] Session valid for tenant: %s", session.TenantID)

		db, err := am.pool.GetConnection(session.TenantID)
		if err != nil {
			log.Printf("[AUTH DEBUG] GetConnection error: %v", err)
			am.renderMaintenancePage(w, err)
			return
		}

		userEmail, err := am.getUserEmail(session.UserID)
		if err != nil {
			am.renderMaintenancePage(w, fmt.Errorf("failed to get user email: %w", err))
			return
		}

		ctx := appcontext.WithTenantID(r.Context(), session.TenantID)
		ctx = appcontext.WithTenantDB(ctx, db)
		ctx = appcontext.WithUserID(ctx, session.UserID)
		ctx = appcontext.WithUserEmail(ctx, userEmail)

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
		"/favicon.ico",
		"/auth/login",
		"/auth/google",
		"/auth/google/callback",
		"/auth/signup",
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

func (am *AuthMiddleware) renderMaintenancePage(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	errMsg := ""
	if err != nil {
		errMsg = fmt.Sprintf("<p style='color: #666; font-size: 0.9em;'>Detalhe técnico: %v</p>", err)
	}
	w.Write([]byte(`<!DOCTYPE html>
<html><head><title>Manuten&#231;&#227;o</title><meta charset="UTF-8"></head>
<body style="font-family: system-ui; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #E1F5EE;">
<div style="text-align: center; padding: 2rem; background: white; border-radius: 1rem; box-shadow: 0 4px 6px rgba(0,0,0,0.1);"><h1 style="color: #1B4D3E; font-family: 'Source Serif 4', serif;">Manuten&#231;&#227;o Tempor&#225;ria</h1><p style="color: #2D6A4F;">Seu consult&#243;rio est&#225; temporariamente indispon&#237;vel.<br>Tente novamente em alguns minutos.</p>` + errMsg + `</div>
</body></html>`))
}

func GetTenantID(ctx context.Context) (string, error) {
	return appcontext.GetTenantID(ctx)
}

func GetTenantDB(ctx context.Context) (*sql.DB, error) {
	return appcontext.GetTenantDB(ctx)
}

func GetUserID(ctx context.Context) (string, error) {
	return appcontext.GetUserID(ctx)
}

func (am *AuthMiddleware) getUserEmail(userID string) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = ?`

	err := am.centralDB.QueryRow(query, userID).Scan(&email)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to query user email: %w", err)
	}

	return email, nil
}
