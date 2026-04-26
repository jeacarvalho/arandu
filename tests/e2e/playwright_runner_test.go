package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/appointment"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/session"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/middleware"
	"arandu/internal/web"
	"arandu/internal/web/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const playwrightEmail = "playwright_e2e@arandu.internal"
const playwrightPassword = "playwright_test_2026"

type pwClassificationMock struct{}

func (m *pwClassificationMock) GetTags(ctx context.Context) ([]observation.Tag, error)             { return nil, nil }
func (m *pwClassificationMock) GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error) { return nil, nil }
func (m *pwClassificationMock) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error { return nil }
func (m *pwClassificationMock) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error   { return nil }
func (m *pwClassificationMock) GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error) { return nil, nil }
func (m *pwClassificationMock) GetObservation(ctx context.Context, id string) (*observation.Observation, error) { return nil, nil }

type pwAgendaMock struct{}

func (m *pwAgendaMock) CreateAppointment(ctx context.Context, patientID, patientName string, date time.Time, startTime, endTime string, duration int, apptType appointment.AppointmentType, notes string) (*appointment.Appointment, error) {
	id := uuid.New().String()
	return &appointment.Appointment{
		ID:          id,
		PatientID:  patientID,
		PatientName: patientName,
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
		Duration:  duration,
		Type:       apptType,
		Status:    appointment.AppointmentStatusScheduled,
		Notes:      notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *pwAgendaMock) GetAppointment(ctx context.Context, id string) (*appointment.Appointment, error) {
	return nil, nil
}

func (m *pwAgendaMock) UpdateAppointment(ctx context.Context, id string, patientID, patientName string, date time.Time, startTime, endTime string, duration int, notes string) error {
	return nil
}

func (m *pwAgendaMock) CancelAppointment(ctx context.Context, id string) error {
	return nil
}

func (m *pwAgendaMock) ConfirmAppointment(ctx context.Context, id string) error {
	return nil
}

func (m *pwAgendaMock) MarkNoShow(ctx context.Context, id string) error {
	return nil
}

func (m *pwAgendaMock) CompleteAppointment(ctx context.Context, id string, sessionID string) error {
	return nil
}

func (m *pwAgendaMock) GetDayView(ctx context.Context, date time.Time) (*services.DayView, error) {
	return &services.DayView{}, nil
}

func (m *pwAgendaMock) GetWeekView(ctx context.Context, date time.Time) (*services.WeekView, error) {
	return &services.WeekView{}, nil
}

func (m *pwAgendaMock) GetMonthView(ctx context.Context, year int, month int) (*services.MonthView, error) {
	return &services.MonthView{}, nil
}

func (m *pwAgendaMock) GetAvailableSlots(ctx context.Context, date time.Time) ([]appointment.TimeSlot, error) {
	return nil, nil
}

func (m *pwAgendaMock) CheckConflicts(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

func (m *pwAgendaMock) GetUpcomingAppointments(ctx context.Context, limit int) ([]*appointment.Appointment, error) {
	return nil, nil
}

func (m *pwAgendaMock) GetPatientAppointments(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

type pwAgendaMockFull struct {
	pwAgendaMock
}

func (m *pwAgendaMockFull) GetAppointment(ctx context.Context, id string) (*appointment.Appointment, error) {
	return nil, nil
}

func (m *pwAgendaMockFull) GetDayView(ctx context.Context, date time.Time) (*services.DayView, error) {
	return &services.DayView{}, nil
}

func (m *pwAgendaMockFull) GetWeekView(ctx context.Context, date time.Time) (*services.WeekView, error) {
	return &services.WeekView{}, nil
}

func (m *pwAgendaMockFull) GetMonthView(ctx context.Context, year int, month int) (*services.MonthView, error) {
	return &services.MonthView{}, nil
}

func (m *pwAgendaMockFull) GetAvailableSlots(ctx context.Context, date time.Time) ([]appointment.TimeSlot, error) {
	return nil, nil
}

func (m *pwAgendaMockFull) CheckConflicts(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

func (m *pwAgendaMockFull) GetUpcomingAppointments(ctx context.Context, limit int) ([]*appointment.Appointment, error) {
	return nil, nil
}

type pwAgendaSessionServiceMock struct{}

func (m *pwAgendaSessionServiceMock) CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error) {
	id := uuid.New().String()
	return &session.Session{
		ID:        id,
		PatientID: patientID,
		Date:     date,
		Summary:  summary,
	}, nil
}

type PlaywrightTestSuite struct {
	t             *testing.T
	router        http.Handler
	centralDB     *sqlite.CentralDB
	db            *sqlite.DB
	tenantPool    *sqlite.TenantPool
	tenantID      string
	sessionCookie *http.Cookie
	patientID     string
}

func setupPlaywrightEnvironment(t *testing.T) *PlaywrightTestSuite {
	tmpDir, err := os.MkdirTemp("", "playwright-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	centralDB, err := sqlite.NewCentralDB(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create central DB: %v", err)
	}

	if err := centralDB.Migrate(nil); err != nil {
		centralDB.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	tenantPool := sqlite.NewTenantPool(tmpDir, nil)

	tenantID := uuid.New().String()
	tenantDBPath := filepath.Join(tmpDir, fmt.Sprintf("arandu_%s.db", tenantID))
	tenantDB, err := sqlite.NewDB(tenantDBPath)
	if err != nil {
		centralDB.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create tenant DB: %v", err)
	}

	centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
		tenantID, tenantDBPath)

	if err := tenantDB.Migrate(); err != nil {
		centralDB.Close()
		tenantDB.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to migrate tenant DB: %v", err)
	}

	return &PlaywrightTestSuite{
		t:          t,
		centralDB:  centralDB,
		db:         tenantDB,
		tenantPool: tenantPool,
		tenantID:   tenantID,
	}
}

func setupRouterPlaywright(s *PlaywrightTestSuite) {
	mux := http.NewServeMux()

	patientRepo := sqlite.NewContextAwarePatientRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	sessionRepo := sqlite.NewContextAwareSessionRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	observationRepo := sqlite.NewContextAwareObservationRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	interventionRepo := sqlite.NewContextAwareInterventionRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	insightRepo := sqlite.NewContextAwareInsightRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	medicationRepo := sqlite.NewContextAwareMedicationRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	vitalsRepo := sqlite.NewContextAwareVitalsRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	goalRepo := sqlite.NewContextAwareGoalRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))
	timelineRepo := sqlite.NewContextAwareTimelineRepository(sqlite.NewContextAwareRepositoryFactory(s.db, s.tenantPool))

	biopsychosocialService := services.NewBiopsychosocialService(medicationRepo, vitalsRepo)

	patientService := services.NewPatientService(patientRepo)
	sessionService := services.NewSessionService(sessionRepo)
	observationService := services.NewObservationService(observationRepo)
	interventionService := services.NewInterventionService(interventionRepo)
	insightService := services.NewInsightService(insightRepo)

	sessionServiceAdapter := web.NewSessionServiceAdapter(sessionService)
	insightServiceAdapter := web.NewInsightServiceAdapter(insightService)
	patientServiceAdapter := web.NewPatientServiceAdapter(patientService)
	observationServiceAdapter := web.NewObservationServiceAdapter(observationService)
	interventionServiceAdapter := web.NewInterventionServiceAdapter(interventionService)
	goalServiceAdapter := web.NewGoalServiceAdapter(goalRepo)

	timelineService := services.NewTimelineServiceContext(timelineRepo, patientServiceAdapter)
	timelineServiceAdapter := web.NewTimelineServiceAdapter(timelineService)

	var classificationServiceAdapter handlers.ClassificationServiceInterface = &pwClassificationMock{}
	var agendaServiceAdapter handlers.AgendaServicePort = &pwAgendaMock{}
	var dashboardAgendaAdapter handlers.DashboardAgendaService = &pwAgendaMock{}
	var pwAgendaSessionMock handlers.AgendaSessionServiceInterface = &pwAgendaSessionServiceMock{}

	biopsychosocialServiceAdapterImpl := handlers.BiopsychosocialServiceFuncs{
		GetMedicationsFunc: func(ctx context.Context, patientID string) ([]interface{}, error) {
			meds, err := biopsychosocialService.GetMedications(ctx, patientID)
			if err != nil {
				return nil, err
			}
			result := make([]interface{}, len(meds))
			for i, m := range meds {
				result[i] = m
			}
			return result, nil
		},
		GetLatestVitalsFunc: func(ctx context.Context, patientID string) (interface{}, error) {
			return biopsychosocialService.GetLatestVitals(ctx, patientID)
		},
		GetAverageVitalsFunc: func(ctx context.Context, patientID string, days int) (interface{}, error) {
			return biopsychosocialService.GetAverageVitals(ctx, patientID, days)
		},
	}

	patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, insightServiceAdapter, biopsychosocialServiceAdapterImpl, timelineServiceAdapter, web.NewAnamnesisServiceAdapter(patientRepo), agendaServiceAdapter)
	sessionHandler := handlers.NewSessionHandler(sessionServiceAdapter, patientServiceAdapter, observationServiceAdapter, interventionServiceAdapter, goalServiceAdapter, classificationServiceAdapter, nil)
	observationHandler := handlers.NewObservationHandler(observationServiceAdapter)
	interventionHandler := handlers.NewInterventionHandler(interventionServiceAdapter)
	dashboardHandler := handlers.NewDashboardHandler(patientServiceAdapter, sessionServiceAdapter, dashboardAgendaAdapter)
	timelineHandler := handlers.NewTimelineHandler(timelineServiceAdapter, patientServiceAdapter)
	biopsychosocialHandler := handlers.NewBiopsychosocialHandler(biopsychosocialService)
	authHandler := handlers.NewAuthHandler(s.centralDB)
	agendaHandler := handlers.NewAgendaHandler(&pwAgendaMockFull{}, patientServiceAdapter, pwAgendaSessionMock)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/dashboard", dashboardHandler.Show)
	mux.HandleFunc("/patients", patientHandler.ListPatients)
	mux.HandleFunc("/patients/new", patientHandler.NewPatient)
	mux.HandleFunc("/patients/search", patientHandler.Search)
	mux.HandleFunc("/patients/create", patientHandler.CreatePatient)

	mux.HandleFunc("/session/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/session/" && r.Method == "GET" {
			sessionHandler.NewSession(w, r)
		} else if r.URL.Path == "/session/" && r.Method == "POST" {
			sessionHandler.CreateSession(w, r)
		} else if r.URL.Path == "/session" && r.Method == "POST" {
			sessionHandler.CreateSession(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/session/") {
			id := r.URL.Path[9:] // strip "/session/"
			if strings.HasSuffix(id, "/edit") {
				sessionHandler.EditSession(w, r)
			} else if strings.HasSuffix(id, "/update") {
				sessionHandler.UpdateSession(w, r)
			} else if strings.HasSuffix(id, "/observations") {
				sessionHandler.CreateObservation(w, r)
			} else if strings.HasSuffix(id, "/interventions") {
				sessionHandler.CreateIntervention(w, r)
			} else {
				sessionHandler.Show(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/observations/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 13 && r.URL.Path[len(r.URL.Path)-4:] == "/edit" {
			observationHandler.GetObservationEditForm(w, r)
		} else if r.Method == "GET" {
			observationHandler.GetObservation(w, r)
		} else if r.Method == "PUT" {
			observationHandler.UpdateObservation(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/interventions/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 14 && r.URL.Path[len(r.URL.Path)-4:] == "/edit" {
			interventionHandler.GetInterventionEditForm(w, r)
		} else if r.Method == "GET" {
			interventionHandler.GetIntervention(w, r)
		} else if r.Method == "PUT" {
			interventionHandler.UpdateIntervention(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/patients/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 9 && path[len(path)-6:] == "/new" {
			sessionHandler.NewSession(w, r)
		} else if len(path) > 8 && path[len(path)-11:] == "/sessions/new" {
			sessionHandler.NewSession(w, r)
		} else if len(path) > 8 && path[len(path)-10:] == "/load-more" {
			timelineHandler.LoadMoreEvents(w, r)
		} else if len(path) > 7 && path[len(path)-7:] == "/search" {
			timelineHandler.SearchPatientHistory(w, r)
		} else if len(path) > 7 && path[len(path)-7:] == "/history" {
			timelineHandler.ShowPatientHistory(w, r)
		} else if len(path) > 5 && path[len(path)-8:] == "/sessions" {
			patientHandler.ListSessions(w, r)
		} else if len(path) > 5 && path[len(path)-7:] == "/context" {
			biopsychosocialHandler.GetContextPanel(w, r)
		} else if len(path) > 5 && path[len(path)-12:] == "/medications" && r.Method == "POST" {
			biopsychosocialHandler.AddMedication(w, r)
		} else if len(path) > 5 && path[len(path)-6:] == "/vitals" && r.Method == "POST" {
			biopsychosocialHandler.RecordVitals(w, r)
		} else {
			patientHandler.Show(w, r)
		}
	})

	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     middleware.SessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	mux.HandleFunc("/agenda", agendaHandler.View)
	mux.HandleFunc("/agenda/day", agendaHandler.DayView)
	mux.HandleFunc("/agenda/week", agendaHandler.WeekView)
	mux.HandleFunc("/agenda/month", agendaHandler.MonthView)
	mux.HandleFunc("/agenda/new", agendaHandler.NewForm)
	mux.HandleFunc("/agenda/slots", agendaHandler.GetSlots)

	mux.HandleFunc("/agenda/appointments/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[18:]
		if len(path) > 0 {
			id := path
			if r.Method == "GET" {
				agendaHandler.Show(w, r)
			} else if r.Method == "PUT" {
				agendaHandler.Update(w, r)
			} else if r.Method == "DELETE" {
				agendaHandler.Cancel(w, r)
			} else if r.Method == "POST" && id == "confirm" {
				agendaHandler.Confirm(w, r)
			} else if r.Method == "POST" && id == "noshow" {
				agendaHandler.NoShow(w, r)
			} else if r.Method == "POST" && id == "complete" {
				agendaHandler.Complete(w, r)
			} else {
				http.NotFound(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/agenda/appointments", agendaHandler.Create)

	searchHandler := handlers.NewSearchHandler(timelineServiceAdapter, patientServiceAdapter)
	mux.HandleFunc("/search", searchHandler.Search)

	authMiddleware := middleware.NewAuthMiddleware(s.centralDB, s.tenantPool)
	protectedHandler := authMiddleware.Middleware(mux)

	s.router = protectedHandler
}

func (s *PlaywrightTestSuite) createTestUserAndSession() error {
	userID := uuid.New().String()
	sessionID := uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(playwrightPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.centralDB.Exec(`
		INSERT INTO users (id, email, password_hash, tenant_id, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
	`, userID, playwrightEmail, string(hash), s.tenantID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	_, err = s.centralDB.Exec(`
		INSERT INTO sessions (id, user_id, tenant_id, created_at, expires_at)
		VALUES (?, ?, ?, datetime('now'), ?)
	`, sessionID, userID, s.tenantID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	s.sessionCookie = &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}

	return nil
}

func (s *PlaywrightTestSuite) teardown() {
	s.centralDB.Close()
	s.tenantPool.CloseAll()
	s.db.Close()
}

func TestPlaywrightClinicalWorkflow(t *testing.T) {
	suite := setupPlaywrightEnvironment(t)
	defer suite.teardown()
	setupRouterPlaywright(suite)

	if err := suite.createTestUserAndSession(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	server := httptest.NewServer(suite.router)
	defer server.Close()

	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not available — skipping Playwright tests")
	}

	cmd := exec.Command("npx", "playwright", "test", "--project=chromium", "--reporter=list", "tests/e2e/playwright/")
	cmd.Env = append(os.Environ(),
		"PLAYWRIGHT_BASE_URL="+server.URL,
		"E2E_EMAIL="+playwrightEmail,
		"E2E_PASSWORD="+playwrightPassword,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/home/s015533607/Documentos/desenv/arandu"

	if err := cmd.Run(); err != nil {
		t.Fatalf("Playwright tests failed: %v", err)
	}
}