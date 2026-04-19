package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"arandu/internal/infrastructure/auth"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/env"
	authComponents "arandu/web/components/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	centralDB  *sqlite.CentralDB
	googleAuth *auth.GoogleProvider
}

func NewAuthHandler(centralDB *sqlite.CentralDB) *AuthHandler {
	googleAuth, err := auth.NewGoogleProvider()
	if err != nil {
		log.Printf("⚠️ Google OAuth não configurado: %v", err)
	} else {
		log.Printf("✅ Google OAuth configurado com sucesso")
	}

	return &AuthHandler{
		centralDB:  centralDB,
		googleAuth: googleAuth,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Log para debug
	log.Printf("[AuthHandler.Login] Method: %s, Path: %s", r.Method, r.URL.Path)

	// Login page should always be a full page request, not an HTMX fragment
	// If accessed via HTMX, redirect to the full login page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	isDev := env.IsDev()

	if r.Method == http.MethodGet {
		data := authComponents.LoginData{
			Error: "",
			IsDev: isDev,
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
				IsDev: isDev,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Login(data).Render(r.Context(), w)
			return
		}

		sessionID, err := h.authenticateUser(email, password)
		if err != nil {
			data := authComponents.LoginData{
				Error: err.Error(),
				IsDev: isDev,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Login(data).Render(r.Context(), w)
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

	sessionID, err := h.findOrCreateUser(userInfo.Email)
	if err != nil {
		log.Printf("Failed to find or create user: %v", err)
		http.Redirect(w, r, "/login?error=user_creation_failed", http.StatusFound)
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

func (h *AuthHandler) authenticateUser(email, password string) (string, error) {
	var userID, tenantID, passwordHash string

	query := `SELECT id, tenant_id, password_hash FROM users WHERE email = ? LIMIT 1`
	err := h.centralDB.QueryRow(query, email).Scan(&userID, &tenantID, &passwordHash)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("usuário não encontrado")
	} else if err != nil {
		return "", fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if passwordHash == "" {
		return "", fmt.Errorf("usuário não possui senha cadastrada")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("senha incorreta")
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	_, err = h.centralDB.Exec(
		`INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, tenantID, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	log.Printf("✅ Login bem-sucedido para usuário=%s", email)
	return sessionID, nil
}

func (h *AuthHandler) findOrCreateUser(email string) (string, error) {
	var userID, tenantID string
	var expiresAt int64

	query := `SELECT id, tenant_id FROM users WHERE email = ? LIMIT 1`
	err := h.centralDB.QueryRow(query, email).Scan(&userID, &tenantID)

	if err == sql.ErrNoRows {
		// Novo usuário via Google OAuth - criar tenant e usuário automaticamente
		log.Printf("Novo usuário Google: %s - criando tenant e usuário", email)

		tenantID = uuid.New().String()
		dbPath := fmt.Sprintf("storage/tenants/clinical_%s.db", tenantID)

		// Criar tenant
		_, err = h.centralDB.Exec(
			`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
			tenantID, dbPath,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create tenant: %w", err)
		}

		// Criar usuário (sem senha pois é OAuth)
		userID = uuid.New().String()
		_, err = h.centralDB.Exec(
			`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
			userID, email, tenantID,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}

		log.Printf("✅ Criado tenant=%s para usuário=%s", tenantID, email)
	} else if err != nil {
		return "", fmt.Errorf("failed to query user: %w", err)
	}

	sessionID := uuid.New().String()
	expiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()

	_, err = h.centralDB.Exec(
		`INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, tenantID, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
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

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Signup] Method: %s", r.Method)

	if r.Method == http.MethodGet {
		data := authComponents.SignupData{
			Error:   "",
			Email:   "",
			Tenant:  "",
			Success: false,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		authComponents.Signup(data).Render(r.Context(), w)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		log.Printf("[Signup] Email: %s, Password length: %d", email, len(password))

		if email == "" || password == "" {
			data := authComponents.SignupData{
				Error:   "Email e senha são obrigatórios",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}

		var existingID string
		err := h.centralDB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
		if err == nil {
			log.Printf("[Signup] User already exists: %s", email)
			data := authComponents.SignupData{
				Error:   "Usuário já existe",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}
		if err != sql.ErrNoRows {
			log.Printf("[Signup] Error checking user: %v", err)
			data := authComponents.SignupData{
				Error:   "Erro ao verificar usuário",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}

		log.Printf("[Signup] Creating new user: %s", email)
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			data := authComponents.SignupData{
				Error:   "Erro ao criar hash da senha",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}

		tenantID := uuid.New().String()
		dbPath := fmt.Sprintf("storage/tenants/clinical_%s.db", tenantID)

		_, err = h.centralDB.Exec(
			`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
			tenantID, dbPath,
		)
		if err != nil {
			data := authComponents.SignupData{
				Error:   "Erro ao criar tenant",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}

		userID := uuid.New().String()
		_, err = h.centralDB.Exec(
			`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, ?, ?)`,
			userID, email, string(passwordHash), tenantID,
		)
		if err != nil {
			data := authComponents.SignupData{
				Error:   "Erro ao criar usuário",
				Email:   email,
				Success: false,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			authComponents.Signup(data).Render(r.Context(), w)
			return
		}

		log.Printf("✅ Usuário criado: %s com tenant: %s", email, tenantID)

		data := authComponents.SignupData{
			Email:   email,
			Tenant:  tenantID,
			Success: true,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		authComponents.Signup(data).Render(r.Context(), w)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	case "/logout":
		h.Logout(w, r)
	case "/auth/signup":
		h.Signup(w, r)
	default:
		http.NotFound(w, r)
	}
}
