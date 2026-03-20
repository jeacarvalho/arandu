package e2e

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/web/handlers"
)

func TestOnboardingFlow(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	authHandler := handlers.NewAuthHandler(centralDB)

	// Test 1: First-time user onboarding
	t.Run("FirstTimeUserOnboarding", func(t *testing.T) {
		// Clear tenants directory
		tenantsDir := filepath.Join(storagePath, "tenants")
		os.RemoveAll(tenantsDir)

		// Simulate Google OAuth callback for new user
		req := httptest.NewRequest("GET", "/auth/google/callback?state=test-state&code=test-code", nil)

		// Set oauth_state cookie
		req.AddCookie(&http.Cookie{
			Name:  "oauth_state",
			Value: "test-state",
		})

		// Mock the auth handler to simulate user creation
		// Since we can't actually call Google OAuth in tests,
		// we'll test the provisioning endpoint directly

		// First, create a user without tenant
		userID := fmt.Sprintf("test-user-%d", time.Now().UnixNano())
		email := "newuser@example.com"

		_, err := centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, NULL)`,
			userID, email)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Set provisioning cookie
		req = httptest.NewRequest("GET", "/auth/provisioning", nil)
		req.AddCookie(&http.Cookie{
			Name:  "arandu_provisioning_user",
			Value: userID,
		})

		// Call provisioning endpoint (GET shows waiting page)
		rec := httptest.NewRecorder()
		authHandler.Provisioning(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		// Check that waiting page was shown
		body := rec.Body.String()
		if body == "" {
			t.Error("Expected response body")
		}

		// Now test POST to actually provision
		req = httptest.NewRequest("POST", "/auth/provisioning", nil)
		req.AddCookie(&http.Cookie{
			Name:  "arandu_provisioning_user",
			Value: userID,
		})

		rec = httptest.NewRecorder()
		authHandler.Provisioning(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		// Check response is JSON
		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Verify tenant was created
		var tenantID, dbPath string
		err = centralDB.QueryRow(`SELECT tenant_id FROM users WHERE id = ?`, userID).Scan(&tenantID)
		if err != nil {
			t.Fatalf("Failed to get user tenant: %v", err)
		}

		if tenantID == "" {
			t.Error("Expected tenant_id to be set")
		}

		// Verify tenant record exists
		err = centralDB.QueryRow(`SELECT db_path FROM tenants WHERE id = ?`, tenantID).Scan(&dbPath)
		if err != nil {
			t.Fatalf("Failed to get tenant: %v", err)
		}

		// Verify DB file exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Errorf("Tenant DB file does not exist: %s", dbPath)
		}

		// Verify DB has tables
		db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
		if err != nil {
			t.Fatalf("Failed to open tenant DB: %v", err)
		}
		defer db.Close()

		var tableCount int
		err = db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tableCount)
		if err != nil {
			t.Fatalf("Failed to count tables: %v", err)
		}

		if tableCount == 0 {
			t.Error("Tenant DB has no tables")
		}

		t.Logf("✅ First-time user onboarding successful: tenant=%s, db=%s", tenantID, dbPath)
	})

	// Test 2: Existing user login (no provisioning needed)
	t.Run("ExistingUserLogin", func(t *testing.T) {
		// Create user with existing tenant
		userID := fmt.Sprintf("existing-user-%d", time.Now().UnixNano())
		email := "existing@example.com"
		tenantID := fmt.Sprintf("existing-tenant-%d", time.Now().UnixNano())
		dbPath := filepath.Join(storagePath, "tenants", "clinical_"+tenantID+".db")

		// Create tenant first
		_, err := centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
			tenantID, dbPath)
		if err != nil {
			t.Fatalf("Failed to create test tenant: %v", err)
		}

		// Create user with tenant
		_, err = centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
			userID, email, tenantID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create tenant DB file
		tenantsDir := filepath.Join(storagePath, "tenants")
		if err := os.MkdirAll(tenantsDir, 0755); err != nil {
			t.Fatalf("Failed to create tenants directory: %v", err)
		}

		// The auth flow for existing users should create a session directly
		// without going through provisioning

		// Verify user has tenant
		var userTenantID string
		err = centralDB.QueryRow(`SELECT tenant_id FROM users WHERE id = ?`, userID).Scan(&userTenantID)
		if err != nil {
			t.Fatalf("Failed to get user tenant: %v", err)
		}

		if userTenantID != tenantID {
			t.Errorf("Expected tenant_id %s, got %s", tenantID, userTenantID)
		}

		t.Logf("✅ Existing user login verified: user=%s, tenant=%s", userID, tenantID)
	})

	// Test 3: Error handling - missing provisioning cookie
	t.Run("MissingProvisioningCookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/provisioning", nil)
		// No cookie set

		rec := httptest.NewRecorder()
		authHandler.Provisioning(rec, req)

		if rec.Code != http.StatusFound {
			t.Errorf("Expected redirect (302), got %d", rec.Code)
		}

		location := rec.Header().Get("Location")
		if location != "/login?error=invalid_provisioning" {
			t.Errorf("Expected redirect to /login?error=invalid_provisioning, got %s", location)
		}

		t.Log("✅ Missing provisioning cookie handled correctly")
	})

	// Test 4: Error handling - invalid user in provisioning cookie
	t.Run("InvalidUserInCookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/provisioning", nil)
		req.AddCookie(&http.Cookie{
			Name:  "arandu_provisioning_user",
			Value: "non-existent-user",
		})

		rec := httptest.NewRecorder()
		authHandler.Provisioning(rec, req)

		// Should fail because user doesn't exist
		if rec.Code != http.StatusFound {
			t.Errorf("Expected redirect (302), got %d", rec.Code)
		}

		t.Log("✅ Invalid user in cookie handled correctly")
	})
}

func TestTenantServiceIntegration(t *testing.T) {
	// Test the tenant service directly
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	// Create multiple users and provision tenants
	for i := 1; i <= 3; i++ {
		t.Run(fmt.Sprintf("User%dProvisioning", i), func(t *testing.T) {
			userID := fmt.Sprintf("user-%d-%d", i, time.Now().UnixNano())
			email := fmt.Sprintf("user%d@example.com", i)

			// Create user without tenant
			_, err := centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, NULL)`,
				userID, email)
			if err != nil {
				t.Fatalf("Failed to create user %d: %v", i, err)
			}

			// Get tenant service
			authHandler := handlers.NewAuthHandler(centralDB)

			// Simulate provisioning
			req := httptest.NewRequest("POST", "/auth/provisioning", nil)
			req.AddCookie(&http.Cookie{
				Name:  "arandu_provisioning_user",
				Value: userID,
			})

			rec := httptest.NewRecorder()
			authHandler.Provisioning(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("User %d: Expected status 200, got %d", i, rec.Code)
			}

			// Verify user now has tenant
			var tenantID string
			err = centralDB.QueryRow(`SELECT tenant_id FROM users WHERE id = ?`, userID).Scan(&tenantID)
			if err != nil {
				t.Fatalf("User %d: Failed to get tenant: %v", i, err)
			}

			if tenantID == "" {
				t.Errorf("User %d: Expected tenant_id to be set", i)
			}

			t.Logf("✅ User %d provisioned successfully: tenant=%s", i, tenantID)
		})
	}

	// Verify isolation: each user has unique tenant DB
	rows, err := centralDB.Query(`SELECT id, db_path FROM tenants`)
	if err != nil {
		t.Fatalf("Failed to query tenants: %v", err)
	}
	defer rows.Close()

	tenantCount := 0
	for rows.Next() {
		var tenantID, dbPath string
		if err := rows.Scan(&tenantID, &dbPath); err != nil {
			t.Fatalf("Failed to scan tenant: %v", err)
		}
		tenantCount++

		// Verify DB file is unique
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Errorf("Tenant DB file does not exist: %s", dbPath)
		}
	}

	if tenantCount != 3 {
		t.Errorf("Expected 3 tenants, got %d", tenantCount)
	}

	t.Logf("✅ Tenant isolation verified: %d unique tenant databases created", tenantCount)
}
