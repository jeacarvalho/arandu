package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/auth"
	"arandu/internal/infrastructure/repository/sqlite"
	authComponents "arandu/web/components/auth"
)

type AuthHandler struct {
	centralDB     *sqlite.CentralDB
	googleAuth    *auth.GoogleProvider
	tenantService *services.TenantService
}

func NewAuthHandler(centralDB *sqlite.CentralDB) *AuthHandler {
	googleAuth, err := auth.NewGoogleProvider()
	if err != nil {
		log.Printf("⚠️ Google OAuth não configurado: %v", err)
	} else {
		log.Printf("✅ Google OAuth configurado com sucesso")
	}

	tenantService := services.NewTenantService(centralDB.DB, "storage")

	return &AuthHandler{
		centralDB:     centralDB,
		googleAuth:    googleAuth,
		tenantService: tenantService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := authComponents.LoginData{
			Error: "",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		authComponents.Login(data).Render(r.Context(), w)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			data := authComponents.LoginData{
				Error: "Credenciais inválidas",
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Login(data).Render(r.Context(), w)
			return
		}

		w.Write([]byte("Login functionality coming soon"))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.googleAuth == nil {
		log.Printf("⚠️ GoogleLogin: googleAuth é nil")
		http.Redirect(w, r, "/login?error=google_not_configured", http.StatusFound)
		return
	}

	state, err := auth.GenerateState()
	if err != nil {
		log.Printf("⚠️ GoogleLogin: erro ao gerar state: %v", err)
		http.Redirect(w, r, "/login?error=state_generation_failed", http.StatusFound)
		return
	}

	log.Printf("🔐 GoogleLogin: state=%s", state)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // Allow cross-site for OAuth redirect
		Expires:  time.Now().Add(10 * time.Minute),
	})

	authURL := h.googleAuth.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	errorParam := r.URL.Query().Get("error")

	log.Printf("🔙 GoogleCallback: state=%s, code=%s, error=%s", state, code, errorParam)

	if errorParam == "access_denied" {
		log.Printf("🔙 GoogleCallback: usuário cancelou")
		http.Redirect(w, r, "/login?error=cancelled", http.StatusFound)
		return
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie == nil || cookie.Value == "" {
		http.Redirect(w, r, "/login?error=invalid_state", http.StatusFound)
		return
	}

	if !h.googleAuth.ValidateState(state, cookie.Value) {
		http.Redirect(w, r, "/login?error=invalid_state", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if code == "" {
		http.Redirect(w, r, "/login?error=missing_code", http.StatusFound)
		return
	}

	token, err := h.googleAuth.ExchangeCode(r.Context(), code)
	if err != nil {
		log.Printf("Google OAuth code exchange failed: %v", err)
		http.Redirect(w, r, "/login?error=exchange_failed", http.StatusFound)
		return
	}

	userInfo, err := h.googleAuth.GetUserInfo(r.Context(), token)
	if err != nil {
		log.Printf("Google OAuth get user info failed: %v", err)
		http.Redirect(w, r, "/login?error=user_info_failed", http.StatusFound)
		return
	}

	userID, tenantID, err := h.findOrCreateUser(userInfo.Email)
	if err != nil {
		log.Printf("Failed to find or create user: %v", err)
		http.Redirect(w, r, "/login?error=user_creation_failed", http.StatusFound)
		return
	}

	// Se o usuário não tem tenant, provisionar um novo
	if tenantID == "" {
		log.Printf("🔧 Primeiro acesso do usuário %s - provisionando tenant", userInfo.Email)

		// Redirecionar para página de aguarde
		http.SetCookie(w, &http.Cookie{
			Name:     "arandu_provisioning_user",
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   300, // 5 minutos
		})

		http.Redirect(w, r, "/auth/provisioning", http.StatusFound)
		return
	}

	// Usuário existente - criar sessão e redirecionar
	sessionID, err := h.createSession(userID, tenantID)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		http.Redirect(w, r, "/login?error=session_creation_failed", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "arandu_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
	})

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (h *AuthHandler) findOrCreateUser(email string) (string, string, error) {
	var userID, tenantID string

	query := `SELECT id, tenant_id FROM users WHERE email = ? LIMIT 1`
	err := h.centralDB.QueryRow(query, email).Scan(&userID, &tenantID)

	if err == sql.ErrNoRows {
		// Novo usuário via Google OAuth - criar usuário primeiro
		log.Printf("Novo usuário Google: %s - criando usuário", email)

		userID = fmt.Sprintf("user-%d", time.Now().UnixNano())

		// Criar usuário sem tenant_id inicialmente
		_, err = h.centralDB.Exec(
			`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, NULL)`,
			userID, email,
		)
		if err != nil {
			return "", "", fmt.Errorf("failed to create user: %w", err)
		}

		log.Printf("✅ Criado usuário=%s (sem tenant ainda)", email)
		return userID, "", nil
	} else if err != nil {
		return "", "", fmt.Errorf("failed to query user: %w", err)
	}

	return userID, tenantID, nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "arandu_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
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
	case "/login":
		h.Login(w, r)
	case "/auth/login":
		h.Login(w, r)
	case "/auth/google":
		h.GoogleLogin(w, r)
	case "/auth/google/callback":
		h.GoogleCallback(w, r)
	case "/auth/provisioning":
		h.Provisioning(w, r)
	case "/logout":
		h.Logout(w, r)
	default:
		http.NotFound(w, r)
	}
}
