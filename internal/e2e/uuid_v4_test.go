package e2e

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/storage"
	"github.com/google/uuid"
)

func TestUUIDv4_FullAuthFlow(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")
	tenantsPath := filepath.Join(storagePath, "tenants")

	if err := os.MkdirAll(tenantsPath, 0755); err != nil {
		t.Fatalf("Failed to create storage dir: %v", err)
	}

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	tenantPool := sqlite.NewTenantPool(storagePath, nil)
	pathResolver := storage.NewPathResolver(storagePath)

	tenantID := uuid.New().String()
	dbPath := pathResolver.ResolveTenantPath(tenantID)
	userID := uuid.New().String()
	email := "e2e-" + tenantID[:8] + "@test.com"
	sessionID := uuid.New().String()

	t.Logf("=== E2E Test: UUID v4 Full Auth Flow ===")
	t.Logf("Tenant ID: %s", tenantID)
	t.Logf("User ID: %s", userID)
	t.Logf("Email: %s", email)

	t.Run("CreateTenant", func(t *testing.T) {
		_, err := centralDB.Exec(
			`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
			tenantID, dbPath,
		)
		if err != nil {
			t.Fatalf("Failed to create tenant: %v", err)
		}
		t.Logf("✅ Tenant created in central DB")
	})

	t.Run("CreateUser", func(t *testing.T) {
		_, err := centralDB.Exec(
			`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
			userID, email, tenantID,
		)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		t.Logf("✅ User created in central DB")
	})

	t.Run("CreateSession", func(t *testing.T) {
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
		_, err := centralDB.Exec(
			`INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
			sessionID, userID, tenantID, expiresAt,
		)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		t.Logf("✅ Session created in central DB")
	})

	t.Run("TenantPoolCreatesDBFile", func(t *testing.T) {
		if _, err := os.Stat(dbPath); err == nil {
			t.Fatalf("DB file should NOT exist before TenantPool.GetConnection: %s", dbPath)
		}

		tenantDB, err := tenantPool.GetConnection(tenantID)
		if err != nil {
			t.Fatalf("TenantPool.GetConnection failed: %v", err)
		}
		defer tenantPool.CloseConnection(tenantID)

		if _, err := os.Stat(dbPath); err != nil {
			t.Fatalf("DB file should exist after TenantPool.GetConnection: %v", err)
		}
		t.Logf("✅ TenantPool created DB file at hashed path: %s", dbPath)

		if err := tenantDB.Ping(); err != nil {
			t.Fatalf("TenantDB ping failed: %v", err)
		}
		t.Logf("✅ TenantDB is accessible")
	})

	t.Run("InsertPatientInTenantDB", func(t *testing.T) {
		tenantDB, err := tenantPool.GetConnection(tenantID)
		if err != nil {
			t.Fatalf("TenantPool.GetConnection failed: %v", err)
		}
		defer tenantPool.CloseConnection(tenantID)

		patientID := uuid.New().String()
		patientName := "Paciente E2E Teste"
		patientNotes := "Notas de teste E2E"

		_, err = tenantDB.Exec(
			`INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))`,
			patientID, patientName, patientNotes,
		)
		if err != nil {
			t.Fatalf("Failed to insert patient: %v", err)
		}
		t.Logf("✅ Patient created: %s", patientID)

		var name, notes string
		err = tenantDB.QueryRow(`SELECT name, notes FROM patients WHERE id = ?`, patientID).Scan(&name, &notes)
		if err != nil {
			t.Fatalf("Failed to query patient: %v", err)
		}

		if name != patientName {
			t.Errorf("Expected name %q, got %q", patientName, name)
		}
		if notes != patientNotes {
			t.Errorf("Expected notes %q, got %q", patientNotes, notes)
		}
		t.Logf("✅ Patient data verified: %s - %s", name, notes)
	})

	t.Run("VerifyUUIDFormat", func(t *testing.T) {
		var fetchedTenantID, fetchedUserID, fetchedSessionID string
		err := centralDB.QueryRow(`SELECT id FROM tenants WHERE id = ?`, tenantID).Scan(&fetchedTenantID)
		if err != nil {
			t.Fatalf("Failed to query tenant: %v", err)
		}

		err = centralDB.QueryRow(`SELECT id FROM users WHERE id = ?`, userID).Scan(&fetchedUserID)
		if err != nil {
			t.Fatalf("Failed to query user: %v", err)
		}

		err = centralDB.QueryRow(`SELECT id FROM sessions WHERE id = ?`, sessionID).Scan(&fetchedSessionID)
		if err != nil {
			t.Fatalf("Failed to query session: %v", err)
		}

		if !isValidUUID(fetchedTenantID) {
			t.Errorf("Tenant ID is not valid UUID v4: %s", fetchedTenantID)
		}
		if !isValidUUID(fetchedUserID) {
			t.Errorf("User ID is not valid UUID v4: %s", fetchedUserID)
		}
		if !isValidUUID(fetchedSessionID) {
			t.Errorf("Session ID is not valid UUID v4: %s", fetchedSessionID)
		}

		t.Logf("✅ All IDs are valid UUID v4 format")
	})

	t.Run("Cleanup", func(t *testing.T) {
		tenantPool.CloseAll()
		os.RemoveAll(tempDir)
		t.Logf("✅ Cleanup complete")
	})
}

func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}

func TestUUIDv4_AuthHandlerSimulatedFlow(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")
	tenantsPath := filepath.Join(storagePath, "tenants")

	os.MkdirAll(tenantsPath, 0755)

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	centralDB.Migrate(nil)
	tenantPool := sqlite.NewTenantPool(storagePath, nil)

	ctx := context.Background()

	simulateOAuthUserCreation := func(email string) (string, string, error) {
		var existingUserID, existingTenantID string
		err := centralDB.QueryRow(
			`SELECT id, tenant_id FROM users WHERE email = ? LIMIT 1`,
			email,
		).Scan(&existingUserID, &existingTenantID)

		if err == nil {
			return existingUserID, existingTenantID, nil
		}
		if err != sql.ErrNoRows {
			return "", "", err
		}

		tenantID := uuid.New().String()
		dbPath := filepath.Join(tenantsPath, "clinical_"+tenantID+".db")

		_, err = centralDB.Exec(
			`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
			tenantID, dbPath,
		)
		if err != nil {
			return "", "", err
		}

		userID := uuid.New().String()
		_, err = centralDB.Exec(
			`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
			userID, email, tenantID,
		)
		if err != nil {
			return "", "", err
		}

		return userID, tenantID, nil
	}

	t.Run("SimulatedOAuthLogin_NewUser", func(t *testing.T) {
		email := "new-user-" + uuid.New().String()[:8] + "@test.com"

		userID, tenantID, err := simulateOAuthUserCreation(email)
		if err != nil {
			t.Fatalf("simulateOAuthUserCreation failed: %v", err)
		}

		t.Logf("Created user: %s, tenant: %s", userID, tenantID)

		tenantDB, err := tenantPool.GetConnection(tenantID)
		if err != nil {
			t.Fatalf("TenantPool.GetConnection failed: %v", err)
		}
		defer tenantPool.CloseConnection(tenantID)

		patientID := uuid.New().String()
		_, err = tenantDB.ExecContext(ctx,
			`INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))`,
			patientID, "Test Patient", "Test Notes",
		)
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}

		var count int
		err = tenantDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM patients`).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count patients: %v", err)
		}

		if count != 1 {
			t.Errorf("Expected 1 patient, got %d", count)
		}

		t.Logf("✅ OAuth flow simulation successful: patient created in tenant DB")
	})

	t.Run("SimulatedOAuthLogin_ExistingUser", func(t *testing.T) {
		email := "existing-user-" + uuid.New().String()[:8] + "@test.com"

		_, _, err := simulateOAuthUserCreation(email)
		if err != nil {
			t.Fatalf("First login failed: %v", err)
		}

		existingUserID, existingTenantID, err := simulateOAuthUserCreation(email)
		if err != nil {
			t.Fatalf("Second login failed: %v", err)
		}

		tenantDB, err := tenantPool.GetConnection(existingTenantID)
		if err != nil {
			t.Fatalf("TenantPool.GetConnection failed: %v", err)
		}
		defer tenantPool.CloseConnection(existingTenantID)

		var count int
		err = tenantDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM patients`).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count patients: %v", err)
		}

		t.Logf("✅ Existing user login successful: %s, %d patients", existingUserID, count)
	})

	tenantPool.CloseAll()
	os.RemoveAll(tempDir)
}
