package sqlite

import (
	"context"
	"database/sql"
	"time"

	"arandu/internal/domain/insight"
	"arandu/internal/domain/intervention"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
	"arandu/internal/domain/timeline"
	appcontext "arandu/internal/platform/context"
)

type ContextAwareRepositoryFactory struct {
	baseDB *DB
	pool   *TenantPool
}

func NewContextAwareRepositoryFactory(baseDB *DB, pool *TenantPool) *ContextAwareRepositoryFactory {
	return &ContextAwareRepositoryFactory{
		baseDB: baseDB,
		pool:   pool,
	}
}

func (f *ContextAwareRepositoryFactory) getDB(ctx context.Context) (*sql.DB, error) {
	if db, err := appcontext.GetTenantDB(ctx); err == nil && db != nil {
		return db, nil
	}
	if f.baseDB != nil {
		return f.baseDB.DB, nil
	}
	return nil, nil
}

type ContextAwarePatientRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwarePatientRepository(factory *ContextAwareRepositoryFactory) *ContextAwarePatientRepository {
	return &ContextAwarePatientRepository{factory: factory}
}

func (r *ContextAwarePatientRepository) Save(ctx context.Context, p *patient.Patient) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.Save(ctx, p)
}

func (r *ContextAwarePatientRepository) FindByID(ctx context.Context, id string) (*patient.Patient, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwarePatientRepository) FindAll(ctx context.Context) ([]*patient.Patient, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.FindAll(ctx)
}

func (r *ContextAwarePatientRepository) Update(ctx context.Context, p *patient.Patient) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.Update(ctx, p)
}

func (r *ContextAwarePatientRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

func (r *ContextAwarePatientRepository) FindByName(ctx context.Context, name string) ([]*patient.Patient, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.FindByName(ctx, name)
}

func (r *ContextAwarePatientRepository) Search(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.Search(ctx, query, limit, offset)
}

func (r *ContextAwarePatientRepository) CountAll(ctx context.Context) (int, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return 0, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.CountAll(ctx)
}

func (r *ContextAwarePatientRepository) FindPaginated(ctx context.Context, limit, offset int) ([]*patient.Patient, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.FindPaginated(ctx, limit, offset)
}

func (r *ContextAwarePatientRepository) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.GetThemeFrequency(ctx, patientID, limit)
}

type ContextAwareSessionRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareSessionRepository(factory *ContextAwareRepositoryFactory) *ContextAwareSessionRepository {
	return &ContextAwareSessionRepository{factory: factory}
}

func (r *ContextAwareSessionRepository) Create(ctx context.Context, sess *session.Session) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewSessionRepository(&DB{db})
	return tempRepo.Create(ctx, sess)
}

func (r *ContextAwareSessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewSessionRepository(&DB{db})
	return tempRepo.GetByID(ctx, id)
}

func (r *ContextAwareSessionRepository) ListByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewSessionRepository(&DB{db})
	return tempRepo.ListByPatient(ctx, patientID)
}

func (r *ContextAwareSessionRepository) Update(ctx context.Context, sess *session.Session) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewSessionRepository(&DB{db})
	return tempRepo.Update(ctx, sess)
}

type ContextAwareObservationRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareObservationRepository(factory *ContextAwareRepositoryFactory) *ContextAwareObservationRepository {
	return &ContextAwareObservationRepository{factory: factory}
}

func (r *ContextAwareObservationRepository) Save(ctx context.Context, o *observation.Observation) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.Save(ctx, o)
}

func (r *ContextAwareObservationRepository) FindByID(ctx context.Context, id string) (*observation.Observation, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareObservationRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*observation.Observation, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.FindBySessionID(ctx, sessionID)
}

func (r *ContextAwareObservationRepository) FindAll(ctx context.Context) ([]*observation.Observation, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.FindAll(ctx)
}

func (r *ContextAwareObservationRepository) Update(ctx context.Context, o *observation.Observation) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.Update(ctx, o)
}

func (r *ContextAwareObservationRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

type ContextAwareInterventionRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareInterventionRepository(factory *ContextAwareRepositoryFactory) *ContextAwareInterventionRepository {
	return &ContextAwareInterventionRepository{factory: factory}
}

func (r *ContextAwareInterventionRepository) Save(ctx context.Context, i *intervention.Intervention) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.Save(ctx, i)
}

func (r *ContextAwareInterventionRepository) FindByID(ctx context.Context, id string) (*intervention.Intervention, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareInterventionRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*intervention.Intervention, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.FindBySessionID(ctx, sessionID)
}

func (r *ContextAwareInterventionRepository) FindAll(ctx context.Context) ([]*intervention.Intervention, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.FindAll(ctx)
}

func (r *ContextAwareInterventionRepository) Update(ctx context.Context, i *intervention.Intervention) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.Update(ctx, i)
}

func (r *ContextAwareInterventionRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

type ContextAwareInsightRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareInsightRepository(factory *ContextAwareRepositoryFactory) *ContextAwareInsightRepository {
	return &ContextAwareInsightRepository{factory: factory}
}

func (r *ContextAwareInsightRepository) Save(ctx context.Context, i *insight.Insight) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInsightRepository(&DB{db})
	return tempRepo.Save(ctx, i)
}

func (r *ContextAwareInsightRepository) FindByID(ctx context.Context, id string) (*insight.Insight, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInsightRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareInsightRepository) FindAll(ctx context.Context) ([]*insight.Insight, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInsightRepository(&DB{db})
	return tempRepo.FindAll(ctx)
}

func (r *ContextAwareInsightRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInsightRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

type ContextAwareTimelineRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareTimelineRepository(factory *ContextAwareRepositoryFactory) *ContextAwareTimelineRepository {
	return &ContextAwareTimelineRepository{factory: factory}
}

func (r *ContextAwareTimelineRepository) GetTimelineByPatientID(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return timeline.Timeline{}, err
	}
	tempRepo := NewTimelineRepository(&DB{db})
	return tempRepo.GetTimelineByPatientID(ctx, patientID, filterType, limit, offset)
}

func (r *ContextAwareTimelineRepository) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewTimelineRepository(&DB{db})
	return tempRepo.SearchInHistory(ctx, patientID, query)
}

type ContextAwareMedicationRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareMedicationRepository(factory *ContextAwareRepositoryFactory) *ContextAwareMedicationRepository {
	return &ContextAwareMedicationRepository{factory: factory}
}

func (r *ContextAwareMedicationRepository) Save(ctx context.Context, m *patient.Medication) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.Save(ctx, m)
}

func (r *ContextAwareMedicationRepository) FindByID(ctx context.Context, id string) (*patient.Medication, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareMedicationRepository) FindByPatientID(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.FindByPatientID(ctx, patientID)
}

func (r *ContextAwareMedicationRepository) GetActiveMedications(ctx context.Context, patientID string) ([]*patient.Medication, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.GetActiveMedications(ctx, patientID)
}

func (r *ContextAwareMedicationRepository) GetMedicationsByStatus(ctx context.Context, patientID string, status patient.MedicationStatus) ([]*patient.Medication, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.GetMedicationsByStatus(ctx, patientID, status)
}

func (r *ContextAwareMedicationRepository) UpdateStatus(ctx context.Context, id string, status patient.MedicationStatus) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.UpdateStatus(ctx, id, status)
}

func (r *ContextAwareMedicationRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

func (r *ContextAwareMedicationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Medication, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewMedicationRepository(&DB{db})
	return tempRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
}

type ContextAwareVitalsRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareVitalsRepository(factory *ContextAwareRepositoryFactory) *ContextAwareVitalsRepository {
	return &ContextAwareVitalsRepository{factory: factory}
}

func (r *ContextAwareVitalsRepository) Save(ctx context.Context, v *patient.Vitals) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.Save(ctx, v)
}

func (r *ContextAwareVitalsRepository) FindByID(ctx context.Context, id string) (*patient.Vitals, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareVitalsRepository) FindByPatientID(ctx context.Context, patientID string, limit int) ([]*patient.Vitals, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.FindByPatientID(ctx, patientID, limit)
}

func (r *ContextAwareVitalsRepository) GetLatestVitals(ctx context.Context, patientID string) (*patient.Vitals, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.GetLatestVitals(ctx, patientID)
}

func (r *ContextAwareVitalsRepository) GetAverageVitals(ctx context.Context, patientID string, days int) (*patient.VitalsAverage, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.GetAverageVitals(ctx, patientID, days)
}

func (r *ContextAwareVitalsRepository) Update(ctx context.Context, v *patient.Vitals) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.Update(ctx, v)
}

func (r *ContextAwareVitalsRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

func (r *ContextAwareVitalsRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*patient.Vitals, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewVitalsRepository(&DB{db})
	return tempRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
}

type ContextAwareGoalRepository struct {
	factory *ContextAwareRepositoryFactory
}

func NewContextAwareGoalRepository(factory *ContextAwareRepositoryFactory) *ContextAwareGoalRepository {
	return &ContextAwareGoalRepository{factory: factory}
}

func (r *ContextAwareGoalRepository) Create(ctx context.Context, patientID, title, description string) (*patient.TherapeuticGoal, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.Create(ctx, patientID, title, description)
}

func (r *ContextAwareGoalRepository) GetActiveGoals(ctx context.Context, patientID string) ([]*patient.TherapeuticGoal, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.GetActiveGoals(ctx, patientID)
}

func (r *ContextAwareGoalRepository) Save(ctx context.Context, g *patient.TherapeuticGoal) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.Save(ctx, g)
}

func (r *ContextAwareGoalRepository) FindByID(ctx context.Context, id string) (*patient.TherapeuticGoal, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareGoalRepository) FindByPatientID(ctx context.Context, patientID string) ([]*patient.TherapeuticGoal, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.FindByPatientID(ctx, patientID)
}

func (r *ContextAwareGoalRepository) GetGoalsByStatus(ctx context.Context, patientID string, status patient.GoalStatus) ([]*patient.TherapeuticGoal, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.GetGoalsByStatus(ctx, patientID, status)
}

func (r *ContextAwareGoalRepository) UpdateStatus(ctx context.Context, id string, status patient.GoalStatus) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.UpdateStatus(ctx, id, status)
}

func (r *ContextAwareGoalRepository) Update(ctx context.Context, g *patient.TherapeuticGoal) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.Update(ctx, g)
}

func (r *ContextAwareGoalRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

func (r *ContextAwareGoalRepository) CloseWithNote(ctx context.Context, id string, status patient.GoalStatus, closureNote string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewGoalRepository(&DB{db})
	return tempRepo.CloseWithNote(ctx, id, status, closureNote)
}
