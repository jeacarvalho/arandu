package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	config      *oauth2.Config
	redirectURL string
}

type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewGoogleProvider() (*GoogleProvider, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("Google OAuth credentials not configured")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
	}

	return &GoogleProvider{
		config:      config,
		redirectURL: redirectURL,
	}, nil
}

func (gp *GoogleProvider) GetAuthURL(state string) string {
	return gp.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (gp *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := gp.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (gp *GoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := gp.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := decodeJSON(resp.Body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &GoogleUser{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (gp *GoogleProvider) ValidateState(state, expectedState string) bool {
	if state == "" || expectedState == "" {
		return false
	}
	return state == expectedState
}

type OAuthSession struct {
	State        string
	ExpiresAt    time.Time
	CallbackPath string
}

func NewOAuthSession(state string) *OAuthSession {
	return &OAuthSession{
		State:        state,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
		CallbackPath: "/auth/google/callback",
	}
}

func (s *OAuthSession) IsValid() bool {
	return time.Now().Before(s.ExpiresAt)
}

func decodeJSON(data io.Reader, v interface{}) error {
	return json.NewDecoder(data).Decode(v)
}
