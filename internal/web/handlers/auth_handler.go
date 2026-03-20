package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/repository/sqlite"
	authComponents "arandu/web/components/auth"
)

type AuthHandler struct {
	centralDB     *sqlite.CentralDB
	tenantService *services.TenantService
}

func NewAuthHandler(centralDB *sqlite.CentralDB) *AuthHandler {
	tenantService := services.NewTenantService(centralDB.DB, "storage")

	return &AuthHandler{
		centralDB:     centralDB,
		tenantService: tenantService,
	}
}

func (h *AuthHandler) Provisioning(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Mostrar página de aguarde
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		authComponents.Provisioning().Render(r.Context(), w)
		return
	}

	if r.Method == http.MethodPost {
		// Processar provisionamento
		cookie, err := r.Cookie("arandu_provisioning_user")
		if err != nil {
			http.Redirect(w, r, "/login?error=invalid_provisioning", http.StatusFound)
			return
		}

		userID := cookie.Value

		// Provisionar tenant
		ctx := r.Context()
		tenantID, err := h.tenantService.ProvisionNewTenant(ctx, userID)
		if err != nil {
			log.Printf("Failed to provision tenant: %v", err)
			http.Redirect(w, r, "/login?error=provisioning_failed", http.StatusFound)
			return
		}

		// Criar sessão
		sessionID, err := h.createSession(userID, tenantID)
		if err != nil {
			log.Printf("Failed to create session after provisioning: %v", err)
			http.Redirect(w, r, "/login?error=session_creation_failed", http.StatusFound)
			return
		}

		// Limpar cookie de provisionamento
		http.SetCookie(w, &http.Cookie{
			Name:   "arandu_provisioning_user",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		// Definir cookie de sessão
		http.SetCookie(w, &http.Cookie{
			Name:     "arandu_session",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400 * 7,
		})

		// Responder com sucesso para redirecionamento via JavaScript
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "redirect": "/dashboard"}`))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *AuthHandler) createSession(userID, tenantID string) (string, error) {
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	_, err := h.centralDB.Exec(
		`INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, tenantID, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/auth/provisioning":
		h.Provisioning(w, r)
	default:
		http.NotFound(w, r)
	}
}
