package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/web/handlers"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, *sqlite.DB) {
	tmpfile, err := os.CreateTemp("", "testdb-e2e-*.db")
	require.NoError(t, err)

	db, err := sqlite.NewDB(tmpfile.Name())
	require.NoError(t, err)

	// Run migrations
	migrationsDir := "../../internal/infrastructure/repository/sqlite/migrations"
	err = db.Migrate(migrationsDir)
	if err != nil {
		// Fallback for different CWD
		migrationsDir = "../internal/infrastructure/repository/sqlite/migrations"
		err = db.Migrate(migrationsDir)
		require.NoError(t, err)
	}

	err = db.InitSchema()
	require.NoError(t, err)

	patientRepo := sqlite.NewPatientRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)
	observationRepo := sqlite.NewObservationRepository(db)
	interventionRepo := sqlite.NewInterventionRepository(db)
	insightRepo := sqlite.NewInsightRepository(db)

	patientService := services.NewPatientService(patientRepo)
	sessionService := services.NewSessionService(sessionRepo)
	observationService := services.NewObservationService(observationRepo)
	interventionService := services.NewInterventionService(interventionRepo)
	insightService := services.NewInsightService(insightRepo)

	templatePath := "../../web/templates"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		templatePath = "../web/templates"
	}

	h := handlers.NewHandler(patientService, sessionService, observationService, interventionService, insightService, templatePath)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Dashboard)
	mux.HandleFunc("/dashboard", h.Dashboard)
	mux.HandleFunc("/patients", h.Patients)
	mux.HandleFunc("/patients/new", h.NewPatient)
	mux.HandleFunc("/session/", h.Session)
	mux.HandleFunc("/sessions", h.CreateSession)
	mux.HandleFunc("/sessions/edit/", h.EditSession)
	mux.HandleFunc("/sessions/update/", h.UpdateSession)
	mux.HandleFunc("/patient/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sessions/new") {
			h.NewSession(w, r)
		} else {
			h.Patient(w, r)
		}
	})

	staticPath := "../../web/static"
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		staticPath = "../web/static"
	}

	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := httptest.NewServer(mux)

	return server, db
}

func TestEditSessionE2E(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()

	// 1. Preparation: Create a patient and a session
	ctx := context.Background()
	patientRepo := sqlite.NewPatientRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)

	p, err := patient.NewPatient("E2E Patient", "Initial notes")
	require.NoError(t, err)
	require.NoError(t, patientRepo.Save(p))

	sess := session.NewSession(p.ID, time.Now(), "Original session summary")
	require.NoError(t, sessionRepo.Create(ctx, sess))

	// 2. Playwright Test
	pw, err := playwright.Run()
	require.NoError(t, err)
	defer func() { require.NoError(t, pw.Stop()) }()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, browser.Close()) }()

	page, err := browser.NewPage()
	require.NoError(t, err)

	// 3. Navigation
	editURL := fmt.Sprintf("%s/sessions/edit/%s", server.URL, sess.ID)
	_, err = page.Goto(editURL)
	require.NoError(t, err)

	// Verify original summary is present
	originalSummary, err := page.InputValue("#summary")
	require.NoError(t, err)
	require.Equal(t, "Original session summary", originalSummary)

	// Alter data
	updatedDate := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	updatedSummary := "This is the updated summary from the E2E test."
	require.NoError(t, page.Fill("#date", updatedDate))
	require.NoError(t, page.Fill("#summary", updatedSummary))

	// Submit form
	require.NoError(t, page.Click("button[type=submit]"))

	// 5. Verification (UI)
	// Wait for redirection to the patient page
	require.NoError(t, page.WaitForURL(fmt.Sprintf("%s/patient/%s", server.URL, p.ID), playwright.PageWaitForURLOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}))

	// Check if the updated summary is visible on the patient page
	pageContent, err := page.Content()
	require.NoError(t, err)
	require.Contains(t, pageContent, updatedSummary, "The patient page should display the updated session summary")

	// Also check for the updated date, formatted as it appears on the page
	formattedUpdatedDate := time.Now().Add(-48 * time.Hour).Format("02/01/2006")
	require.Contains(t, pageContent, formattedUpdatedDate, "The patient page should display the updated session date")

	// 6. Verification (DB)
	updatedSess, err := sessionRepo.GetByID(ctx, sess.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedSess)
	require.Equal(t, updatedSummary, updatedSess.Summary)

	// Compare dates (ignoring time part)
	require.Equal(t, updatedDate, updatedSess.Date.Format("2006-01-02"))
	require.True(t, updatedSess.UpdatedAt.After(sess.UpdatedAt), "UpdatedAt should have been updated")
}
