package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
	appcontext "arandu/internal/platform/context"
)

func TestContextAwareRepositoryFactory_getDB_WithContext(t *testing.T) {
	baseDB := &DB{DB: &sql.DB{}}
	factory := NewContextAwareRepositoryFactory(baseDB, nil)

	tenantDB := &sql.DB{}
	ctx := appcontext.WithTenantDB(context.Background(), tenantDB)

	db, err := factory.getDB(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db != tenantDB {
		t.Error("expected tenant DB from context")
	}
}

func TestContextAwareRepositoryFactory_getDB_Fallback(t *testing.T) {
	baseDB := &DB{DB: &sql.DB{}}
	factory := NewContextAwareRepositoryFactory(baseDB, nil)

	ctx := context.Background()

	db, err := factory.getDB(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db != baseDB.DB {
		t.Error("expected fallback to baseDB")
	}
}

func TestContextAwareRepositoryFactory_getDB_NoDB(t *testing.T) {
	factory := NewContextAwareRepositoryFactory(nil, nil)

	ctx := context.Background()

	db, err := factory.getDB(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db != nil {
		t.Error("expected nil DB when no fallback")
	}
}

func TestContextAwarePatientRepository_Save(t *testing.T) {
	baseDB, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	defer baseDB.Close()

	if err := baseDB.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	factory := NewContextAwareRepositoryFactory(baseDB, nil)
	repo := NewContextAwarePatientRepository(factory)

	ctx := context.Background()

	patient := &patient.Patient{
		ID:        "test-patient-1",
		Name:      "Test Patient",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Save(ctx, patient)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, "test-patient-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected patient to be found")
	}
	if retrieved.Name != "Test Patient" {
		t.Errorf("expected 'Test Patient', got '%s'", retrieved.Name)
	}
}

func TestContextAwarePatientRepository_Save_WithTenantContext(t *testing.T) {
	baseDB, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	defer baseDB.Close()

	if err := baseDB.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	factory := NewContextAwareRepositoryFactory(baseDB, nil)
	repo := NewContextAwarePatientRepository(factory)

	tenantDB := baseDB
	ctx := appcontext.WithTenantDB(context.Background(), tenantDB.DB)

	patient := &patient.Patient{
		ID:        "test-patient-2",
		Name:      "Tenant Patient",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Save(ctx, patient)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, "test-patient-2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected patient to be found")
	}
	if retrieved.Name != "Tenant Patient" {
		t.Errorf("expected 'Tenant Patient', got '%s'", retrieved.Name)
	}
}

func TestContextAwareSessionRepository_Create(t *testing.T) {
	baseDB, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	defer baseDB.Close()

	if err := baseDB.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	factory := NewContextAwareRepositoryFactory(baseDB, nil)
	repo := NewContextAwareSessionRepository(factory)

	ctx := context.Background()
	now := time.Now()

	sess := &session.Session{
		ID:        "test-session-1",
		PatientID: "patient-1",
		Date:      now,
		Summary:   "Test session",
	}

	err = repo.Create(ctx, sess)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, "test-session-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected session to be found")
	}
	if retrieved.Summary != "Test session" {
		t.Errorf("expected 'Test session', got '%s'", retrieved.Summary)
	}
}

func TestContextAwareMedicationRepository_Save(t *testing.T) {
	baseDB, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	defer baseDB.Close()

	if err := baseDB.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	factory := NewContextAwareRepositoryFactory(baseDB, nil)
	repo := NewContextAwareMedicationRepository(factory)

	ctx := context.Background()

	med := &patient.Medication{
		ID:        "med-1",
		PatientID: "patient-1",
		Name:      "Test Drug",
		Dosage:    "100mg",
		Status:    patient.MedicationStatusActive,
	}

	err = repo.Save(ctx, med)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, "med-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected medication to be found")
	}
	if retrieved.Name != "Test Drug" {
		t.Errorf("expected 'Test Drug', got '%s'", retrieved.Name)
	}
}

func TestContextAwareVitalsRepository_Save(t *testing.T) {
	baseDB, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	defer baseDB.Close()

	if err := baseDB.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	factory := NewContextAwareRepositoryFactory(baseDB, nil)
	repo := NewContextAwareVitalsRepository(factory)

	ctx := context.Background()

	vitals := &patient.Vitals{
		ID:        "vitals-1",
		PatientID: "patient-1",
		Notes:     "Test vitals",
	}

	err = repo.Save(ctx, vitals)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	retrieved, err := repo.FindByID(ctx, "vitals-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved == nil {
		t.Fatal("expected vitals to be found")
	}
	if retrieved.Notes != "Test vitals" {
		t.Errorf("expected 'Test vitals', got '%s'", retrieved.Notes)
	}
}
