package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func main() {
	if _, err := os.Stat("storage/arandu_central.db"); os.IsNotExist(err) {
		log.Fatal("Central database not found. Run the server first.")
	}

	db, err := sql.Open("sqlite", "storage/arandu_central.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tenantID := "test-tenant-001"
	dbPath := "storage/tenants/clinical_" + tenantID + ".db"

	_, err = db.Exec(`INSERT OR REPLACE INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`, tenantID, dbPath)
	if err != nil {
		log.Fatalf("Failed to insert tenant: %v", err)
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	userID := "test-user-001"

	_, err = db.Exec(`INSERT OR REPLACE INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, ?, ?)`,
		userID, "test@arandu.com", string(passwordHash), tenantID)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}

	sessionID := "test-session-001"
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	_, err = db.Exec(`INSERT OR REPLACE INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, tenantID, expiresAt)
	if err != nil {
		log.Fatalf("Failed to insert session: %v", err)
	}

	fmt.Printf("Created user: test@arandu.com / test123\n")
	fmt.Printf("Session ID: %s\n", sessionID)
	fmt.Printf("Tenant ID: %s\n", tenantID)
}
