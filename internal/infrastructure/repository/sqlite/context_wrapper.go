package sqlite

import (
	"context"
	"database/sql"
	"time"

	"arandu/internal/domain/appointment"
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

func (r *ContextAwarePatientRepository) ListForDashboard(ctx context.Context, limit int) ([]*patient.DashboardSummary, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.ListForDashboard(ctx, limit)
}

func (r *ContextAwarePatientRepository) GetAnamnesis(ctx context.Context, patientID string) (*patient.Anamnesis, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.GetAnamnesis(ctx, patientID)
}

func (r *ContextAwarePatientRepository) SaveAnamnesis(ctx context.Context, anamnesis *patient.Anamnesis) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewPatientRepository(&DB{db})
	return tempRepo.SaveAnamnesis(ctx, anamnesis)
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

func (r *ContextAwareObservationRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*observation.Observation, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
}

func (r *ContextAwareObservationRepository) GetTags(ctx context.Context) ([]observation.Tag, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.GetTags(ctx)
}

func (r *ContextAwareObservationRepository) GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.GetTagsByType(ctx, tagType)
}

func (r *ContextAwareObservationRepository) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.AddTagToObservation(ctx, observationID, tagID, intensity)
}

func (r *ContextAwareObservationRepository) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.RemoveTagFromObservation(ctx, observationID, tagID)
}

func (r *ContextAwareObservationRepository) GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.GetObservationTags(ctx, observationID)
}

func (r *ContextAwareObservationRepository) GetTagsSummary(ctx context.Context) ([]observation.TagSummary, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.GetTagsSummary(ctx)
}

func (r *ContextAwareObservationRepository) GetTagsSummaryByPatient(ctx context.Context, patientID string) ([]observation.TagSummary, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.GetTagsSummaryByPatient(ctx, patientID)
}

func (r *ContextAwareObservationRepository) FindByTag(ctx context.Context, tagID string) ([]*observation.Observation, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewObservationRepository(&DB{db})
	return tempRepo.FindByTag(ctx, tagID)
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

func (r *ContextAwareInterventionRepository) FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*intervention.Intervention, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.FindByPatientIDAndTimeframe(ctx, patientID, startTime)
}

// AddTagToIntervention adiciona uma tag a uma intervenção
func (r *ContextAwareInterventionRepository) AddTagToIntervention(ctx context.Context, interventionID, tagID string, intensity int) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.AddTagToIntervention(ctx, interventionID, tagID, intensity)
}

// RemoveTagFromIntervention remove uma tag de uma intervenção
func (r *ContextAwareInterventionRepository) RemoveTagFromIntervention(ctx context.Context, interventionID, tagID string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.RemoveTagFromIntervention(ctx, interventionID, tagID)
}

// GetInterventionTags retorna todas as tags de uma intervenção
func (r *ContextAwareInterventionRepository) GetInterventionTags(ctx context.Context, interventionID string) ([]*intervention.InterventionClassification, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetInterventionTags(ctx, interventionID)
}

// GetAllInterventionTags retorna todas as tags predefinidas
func (r *ContextAwareInterventionRepository) GetAllInterventionTags(ctx context.Context) ([]*intervention.Tag, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetAllInterventionTags(ctx)
}

// GetAllTagsByType retorna tags por tipo
func (r *ContextAwareInterventionRepository) GetAllTagsByType(ctx context.Context, tagType intervention.TagType) ([]*intervention.Tag, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetAllTagsByType(ctx, tagType)
}

// GetIntervention retorna uma intervenção pelo ID
func (r *ContextAwareInterventionRepository) GetIntervention(ctx context.Context, id string) (*intervention.Intervention, error) {
	return r.FindByID(ctx, id)
}

// GetTopInterventionTags retorna as tags mais utilizadas
func (r *ContextAwareInterventionRepository) GetTopInterventionTags(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetTopInterventionTags(ctx, limit)
}

// FindInterventionsByTagID retorna intervenções por tag
func (r *ContextAwareInterventionRepository) FindInterventionsByTagID(ctx context.Context, tagID string) ([]*intervention.Intervention, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.FindInterventionsByTagID(ctx, tagID)
}

// GetTagCountByType retorna a contagem de tags por tipo
func (r *ContextAwareInterventionRepository) GetTagCountByType(ctx context.Context, interventionID string) (map[intervention.TagType]int, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetTagCountByType(ctx, interventionID)
}

// GetInterventionTagsGroupedByType retorna as tags de uma intervenção agrupadas por tipo
func (r *ContextAwareInterventionRepository) GetInterventionTagsGroupedByType(ctx context.Context, interventionID string) (map[intervention.TagType][]*intervention.InterventionClassification, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewInterventionRepository(&DB{db})
	return tempRepo.GetInterventionTagsGroupedByType(ctx, interventionID)
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

func (r *ContextAwareTimelineRepository) SearchGlobal(ctx context.Context, query string, limit int) ([]*timeline.SearchResult, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewTimelineRepository(&DB{db})
	return tempRepo.SearchGlobal(ctx, query, limit)
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

// ContextAwareAppointmentRepository wraps AppointmentRepository with multi-tenant support
type ContextAwareAppointmentRepository struct {
	factory *ContextAwareRepositoryFactory
}

// NewContextAwareAppointmentRepository creates a new context-aware appointment repository
func NewContextAwareAppointmentRepository(factory *ContextAwareRepositoryFactory) *ContextAwareAppointmentRepository {
	return &ContextAwareAppointmentRepository{factory: factory}
}

func (r *ContextAwareAppointmentRepository) Save(ctx context.Context, appt *appointment.Appointment) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.Save(ctx, appt)
}

func (r *ContextAwareAppointmentRepository) FindByID(ctx context.Context, id string) (*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindByID(ctx, id)
}

func (r *ContextAwareAppointmentRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindByDateRange(ctx, startDate, endDate)
}

func (r *ContextAwareAppointmentRepository) FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindByPatient(ctx, patientID)
}

func (r *ContextAwareAppointmentRepository) FindByPatientAndDateRange(ctx context.Context, patientID string, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindByPatientAndDateRange(ctx, patientID, startDate, endDate)
}

func (r *ContextAwareAppointmentRepository) FindByDate(ctx context.Context, date time.Time) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindByDate(ctx, date)
}

func (r *ContextAwareAppointmentRepository) FindOverlapping(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindOverlapping(ctx, date, startTime, endTime, excludeID)
}

func (r *ContextAwareAppointmentRepository) Update(ctx context.Context, appt *appointment.Appointment) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.Update(ctx, appt)
}

func (r *ContextAwareAppointmentRepository) Delete(ctx context.Context, id string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.Delete(ctx, id)
}

func (r *ContextAwareAppointmentRepository) CountByDate(ctx context.Context, date time.Time) (int, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return 0, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.CountByDate(ctx, date)
}

func (r *ContextAwareAppointmentRepository) FindUpcoming(ctx context.Context, fromDate time.Time, limit int) ([]*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindUpcoming(ctx, fromDate, limit)
}

func (r *ContextAwareAppointmentRepository) FindBySessionID(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return nil, err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.FindBySessionID(ctx, sessionID)
}

func (r *ContextAwareAppointmentRepository) UpdatePatientName(ctx context.Context, patientID, patientName string) error {
	db, err := r.factory.getDB(ctx)
	if err != nil || db == nil {
		return err
	}
	tempRepo := NewAppointmentRepository(&DB{db})
	return tempRepo.UpdatePatientName(ctx, patientID, patientName)
}
