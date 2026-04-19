package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"arandu/internal/domain/patient"
)

func setupDashboardTestDB(t *testing.T) (*DB, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "dashboard-test-*.db")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	f.Close()
	db, err := NewDB(f.Name())
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return db, func() {
		db.Close()
		os.Remove(f.Name())
	}
}

func TestPatientRepository_ListForDashboard_EmptyDB(t *testing.T) {
	db, cleanup := setupDashboardTestDB(t)
	defer cleanup()
	repo := NewPatientRepository(db)

	result, err := repo.ListForDashboard(context.Background(), 20)
	if err != nil {
		t.Fatalf("ListForDashboard error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 rows, got %d", len(result))
	}
}

func TestPatientRepository_ListForDashboard_ReturnsPatientFields(t *testing.T) {
	db, cleanup := setupDashboardTestDB(t)
	defer cleanup()
	repo := NewPatientRepository(db)
	ctx := context.Background()

	p, _ := patient.NewPatient("Beatriz Souza", "f", "b", "Médica", "s", "")
	p.Tag = "ANSIEDADE"
	if err := repo.Save(ctx, p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	result, err := repo.ListForDashboard(ctx, 20)
	if err != nil {
		t.Fatalf("ListForDashboard: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}

	row := result[0]
	if row.ID != p.ID {
		t.Errorf("ID: got %q, want %q", row.ID, p.ID)
	}
	if row.Name != "Beatriz Souza" {
		t.Errorf("Name: got %q, want %q", row.Name, "Beatriz Souza")
	}
	if row.Tag != "ANSIEDADE" {
		t.Errorf("Tag: got %q, want %q", row.Tag, "ANSIEDADE")
	}
}

func TestPatientRepository_ListForDashboard_SessionCount(t *testing.T) {
	db, cleanup := setupDashboardTestDB(t)
	defer cleanup()
	repo := NewPatientRepository(db)
	sessRepo := NewSessionRepository(db)
	ctx := context.Background()

	p, _ := patient.NewPatient("Carlos Melo", "m", "", "", "", "")
	repo.Save(ctx, p)

	// insert 3 sessions
	for i := 0; i < 3; i++ {
		_, err := db.ExecContext(ctx,
			`INSERT INTO sessions (id, patient_id, date, summary, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			newID(), p.ID, time.Now().AddDate(0, 0, -i).Format("2006-01-02T15:04:05Z"), "s", time.Now(), time.Now(),
		)
		if err != nil {
			t.Fatalf("insert session: %v", err)
		}
	}
	_ = sessRepo

	result, err := repo.ListForDashboard(ctx, 20)
	if err != nil {
		t.Fatalf("ListForDashboard: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least 1 row")
	}
	if result[0].SessionCount != 3 {
		t.Errorf("SessionCount: got %d, want 3", result[0].SessionCount)
	}
	if result[0].LastSessionDate == nil {
		t.Error("LastSessionDate should not be nil when sessions exist")
	}
}

func TestPatientRepository_ListForDashboard_NextAppointment(t *testing.T) {
	db, cleanup := setupDashboardTestDB(t)
	defer cleanup()
	repo := NewPatientRepository(db)
	ctx := context.Background()

	p, _ := patient.NewPatient("Diana Faria", "f", "", "", "", "")
	repo.Save(ctx, p)

	futureDate := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	_, err := db.ExecContext(ctx,
		`INSERT INTO appointments (id, patient_id, patient_name, date, start_time, end_time, duration, type, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newID(), p.ID, p.Name, futureDate, "10:30", "11:20", 50, "session", "scheduled", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("insert appointment: %v", err)
	}

	result, err := repo.ListForDashboard(ctx, 20)
	if err != nil {
		t.Fatalf("ListForDashboard: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected 1 row")
	}
	if result[0].NextApptDate != futureDate {
		t.Errorf("NextApptDate: got %q, want %q", result[0].NextApptDate, futureDate)
	}
	if result[0].NextApptTime != "10:30" {
		t.Errorf("NextApptTime: got %q, want %q", result[0].NextApptTime, "10:30")
	}
}

func TestPatientRepository_ListForDashboard_LimitRespected(t *testing.T) {
	db, cleanup := setupDashboardTestDB(t)
	defer cleanup()
	repo := NewPatientRepository(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		p, _ := patient.NewPatient("Paciente "+string(rune('A'+i)), "f", "", "", "", "")
		repo.Save(ctx, p)
	}

	result, err := repo.ListForDashboard(ctx, 3)
	if err != nil {
		t.Fatalf("ListForDashboard: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 rows (limit), got %d", len(result))
	}
}

// newID generates a simple unique string for test fixtures.
func newID() string {
	return "test-" + time.Now().Format("150405.000000000")
}
