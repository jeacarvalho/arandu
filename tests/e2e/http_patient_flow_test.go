package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/appointment"
	"arandu/internal/domain/observation"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/logger"
	"arandu/internal/platform/middleware"
	"arandu/internal/web"
	"arandu/internal/web/handlers"

	"github.com/google/uuid"
)

type classificationMock2 struct{}

func (m *classificationMock2) GetTags(ctx context.Context) ([]observation.Tag, error)             { return nil, nil }
func (m *classificationMock2) GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error) { return nil, nil }
func (m *classificationMock2) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error { return nil }
func (m *classificationMock2) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error   { return nil }
func (m *classificationMock2) GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error) { return nil, nil }
func (m *classificationMock2) GetObservation(ctx context.Context, id string) (*observation.Observation, error) { return nil, nil }

type agendaMock2 struct{}

func (m *agendaMock2) GetDayView(ctx context.Context, date time.Time) (*services.DayView, error) {
	return &services.DayView{}, nil
}

func (m *agendaMock2) GetPatientAppointments(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	return nil, nil
}

var _ handlers.DashboardAgendaService = &agendaMock2{}
var _ handlers.AgendaServicePort = &agendaMock2{}

func setupRouter(db *sqlite.DB, centralDB *sqlite.CentralDB, tenantPool *sqlite.TenantPool) http.Handler {
	mux := http.NewServeMux()

	patientRepo := sqlite.NewContextAwarePatientRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	sessionRepo := sqlite.NewContextAwareSessionRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	observationRepo := sqlite.NewContextAwareObservationRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	interventionRepo := sqlite.NewContextAwareInterventionRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	insightRepo := sqlite.NewContextAwareInsightRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	medicationRepo := sqlite.NewContextAwareMedicationRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	vitalsRepo := sqlite.NewContextAwareVitalsRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	goalRepo := sqlite.NewContextAwareGoalRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))
	timelineRepo := sqlite.NewContextAwareTimelineRepository(sqlite.NewContextAwareRepositoryFactory(db, tenantPool))

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

	// Classification service mock
	var classificationServiceAdapter handlers.ClassificationServiceInterface = &classificationMock2{}

	// DashboardAgendaService mock
	var agendaServiceAdapter handlers.AgendaServicePort = &agendaMock2{}
	var dashboardAgendaAdapter handlers.DashboardAgendaService = &agendaMock2{}

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/dashboard", dashboardHandler.Show)
	mux.HandleFunc("/patients", patientHandler.ListPatients)
	mux.HandleFunc("/patients/new", patientHandler.NewPatient)
	mux.HandleFunc("/patients/search", patientHandler.Search)
	mux.HandleFunc("/patients/create", patientHandler.CreatePatient)

	mux.HandleFunc("/session/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
			sessionHandler.EditSession(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/update") && r.Method == "POST" {
			sessionHandler.UpdateSession(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/observations") && r.Method == "POST" {
			sessionHandler.CreateObservation(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/interventions") && r.Method == "POST" {
			sessionHandler.CreateIntervention(w, r)
		} else if r.Method == "GET" {
			sessionHandler.Show(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/session", sessionHandler.CreateSession)

	mux.HandleFunc("/observations/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
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
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
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
		logger.InfoContext(r.Context(), "Route /patients/ called",
			logger.String("path", r.URL.Path),
			logger.String("method", r.Method),
		)
		if strings.HasSuffix(r.URL.Path, "/sessions/new") {
			sessionHandler.NewSession(w, r)
		} else if strings.Contains(r.URL.Path, "/history/load-more") && r.Method == "GET" {
			timelineHandler.LoadMoreEvents(w, r)
		} else if strings.Contains(r.URL.Path, "/history/search") && r.Method == "GET" {
			timelineHandler.SearchPatientHistory(w, r)
		} else if strings.Contains(r.URL.Path, "/history") && r.Method == "GET" {
			timelineHandler.ShowPatientHistory(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/sessions") && r.Method == "GET" {
			patientHandler.ListSessions(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/context") && r.Method == "GET" {
			biopsychosocialHandler.GetContextPanel(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/medications") && r.Method == "POST" {
			biopsychosocialHandler.AddMedication(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/vitals") && r.Method == "POST" {
			biopsychosocialHandler.RecordVitals(w, r)
		} else if strings.Contains(r.URL.Path, "/medications/") && r.Method == "PUT" {
			biopsychosocialHandler.UpdateMedicationStatus(w, r)
		} else if strings.Contains(r.URL.Path, "/analysis/synthesis") && r.Method == "POST" {
			handlers.NewAIHandler(nil).GeneratePatientSynthesis(w, r)
		} else if strings.Contains(r.URL.Path, "/plan/report") && r.Method == "GET" {
			sessionHandler.TherapeuticPlanReport(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/goals") && r.Method == "POST" {
			sessionHandler.CreateGoal(w, r)
		} else if strings.Contains(r.URL.Path, "/goals/") && strings.HasSuffix(r.URL.Path, "/close") && r.Method == "POST" {
			sessionHandler.CloseGoalWithNote(w, r)
		} else {
			logger.InfoContext(r.Context(), "Routing to patientHandler.Show")
			patientHandler.Show(w, r)
		}
	})

	authMiddleware := middleware.NewAuthMiddleware(centralDB, tenantPool)
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		mux.ServeHTTP(w, r)
	})

	protected := authMiddleware.Middleware(corsHandler)
	return protected
}

func TestHTTP_PatientCreationFlow_UsingRealRoutes(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")
	tenantsPath := filepath.Join(storagePath, "tenants")
	os.MkdirAll(tenantsPath, 0755)

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	tenantPool := sqlite.NewTenantPool(storagePath, nil)

	dbPath := filepath.Join(storagePath, "arandu.db")
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create DB: %v", err)
	}
	defer db.Close()

	db.Migrate()

	router := setupRouter(db, centralDB, tenantPool)
	server := httptest.NewServer(router)
	defer server.Close()

	tenantID := uuid.New().String()
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	email := "e2e-" + tenantID[:8] + "@test.com"
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
		tenantID, filepath.Join(tenantsPath, "clinical_"+tenantID+".db"))
	centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
		userID, email, tenantID)
	centralDB.Exec(`INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, userID, tenantID, expiresAt)

	cookie := &http.Cookie{Name: "arandu_session", Value: sessionID}

	t.Run("Step1_Dashboard", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", nil)
		req.AddCookie(cookie)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Dashboard status %d, want %d", resp.Code, http.StatusOK)
		}
	})

	t.Run("Step2_NewPatientForm", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/new", nil)
		req.AddCookie(cookie)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("NewPatientForm status %d, want %d", resp.Code, http.StatusOK)
		}
	})

	t.Run("Step3_CreatePatient", func(t *testing.T) {
		patientName := "Paciente E2E " + uuid.New().String()[:8]
		body := strings.NewReader(fmt.Sprintf("name=%s&notes=Teste E2E", patientName))

		req := httptest.NewRequest("POST", "/patients/create", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		t.Logf("Create response: status=%d, location=%s", resp.Code, resp.Header().Get("Location"))

		if resp.Code != http.StatusSeeOther {
			t.Errorf("CreatePatient status %d, want %d. Body: %s", resp.Code, http.StatusSeeOther, resp.Body.String())
			return
		}

		location := resp.Header().Get("Location")
		if !strings.HasPrefix(location, "/patients/") {
			t.Errorf("Redirect location %q, want /patients/<id>", location)
			return
		}

		t.Logf("✅ Patient created, redirect to: %s", location)
	})

	t.Run("Step4_VerifyPatientInDB", func(t *testing.T) {
		tenantDB, err := tenantPool.GetConnection(tenantID)
		if err != nil {
			t.Fatalf("Failed to get tenant connection: %v", err)
		}

		var count int
		err = tenantDB.QueryRow(`SELECT COUNT(*) FROM patients`).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count patients: %v", err)
		}

		if count == 0 {
			t.Error("No patients found in tenant DB")
		}

		t.Logf("✅ %d patient(s) in tenant DB", count)
	})

	t.Run("Step5_GetPatientDetail", func(t *testing.T) {
		tenantDB, _ := tenantPool.GetConnection(tenantID)
		var patientID, patientName string
		tenantDB.QueryRow(`SELECT id, name FROM patients ORDER BY created_at DESC LIMIT 1`).Scan(&patientID, &patientName)

		req := httptest.NewRequest("GET", "/patients/"+patientID, nil)
		req.AddCookie(cookie)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("GetPatient status %d, want %d", resp.Code, http.StatusOK)
			return
		}

		if !strings.Contains(resp.Body.String(), patientName) {
			t.Errorf("Response does not contain patient name %q", patientName)
		}

		t.Logf("✅ Patient detail page shows: %s", patientName)
	})

	t.Run("Step6_VerifyUUIDFormat", func(t *testing.T) {
		var tid, uid, sid string
		centralDB.QueryRow(`SELECT id FROM tenants WHERE id=?`, tenantID).Scan(&tid)
		centralDB.QueryRow(`SELECT id FROM users WHERE id=?`, userID).Scan(&uid)
		centralDB.QueryRow(`SELECT id FROM sessions WHERE id=?`, sessionID).Scan(&sid)

		if !isValidUUID(tid) || !isValidUUID(uid) || !isValidUUID(sid) {
			t.Errorf("IDs are not valid UUID v4: tenant=%s, user=%s, session=%s", tid, uid, sid)
		}

		t.Logf("✅ All IDs are UUID v4: tenant=%s", tid[:8])
	})

	tenantPool.CloseAll()
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
