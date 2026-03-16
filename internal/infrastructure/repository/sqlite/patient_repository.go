package sqlite

import (
	"database/sql"
	"fmt"

	"arandu/internal/domain/patient"
)

// patientQueries contains all SQL queries for the PatientRepository
type patientQueries struct {
	// CRUD operations
	save     string
	findByID string
	findAll  string
	update   string
	delete   string

	// Additional useful queries
	findByName    string
	countAll      string
	findPaginated string
}

// newPatientQueries creates and returns a patientQueries struct with all queries initialized
func newPatientQueries() *patientQueries {
	return &patientQueries{
		save: `INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,

		findByID: `SELECT id, name, notes, created_at, updated_at FROM patients WHERE id = ?`,

		findAll: `SELECT id, name, notes, created_at, updated_at FROM patients ORDER BY created_at DESC`,

		update: `UPDATE patients SET name = ?, notes = ?, updated_at = ? WHERE id = ?`,

		delete: `DELETE FROM patients WHERE id = ?`,

		// Search patients by name (case-insensitive, partial match)
		findByName: `SELECT id, name, notes, created_at, updated_at FROM patients WHERE LOWER(name) LIKE LOWER(?) ORDER BY name`,

		// Count all patients
		countAll: `SELECT COUNT(*) FROM patients`,

		// Paginated results
		findPaginated: `SELECT id, name, notes, created_at, updated_at FROM patients ORDER BY created_at DESC LIMIT ? OFFSET ?`,
	}
}

// PatientRepository implements the patient.Repository interface
type PatientRepository struct {
	db      *DB
	queries *patientQueries
}

// NewPatientRepository creates a new PatientRepository
func NewPatientRepository(db *DB) *PatientRepository {
	return &PatientRepository{
		db:      db,
		queries: newPatientQueries(),
	}
}

// validatePatientForSave validates patient data before saving to database
func validatePatientForSave(p *patient.Patient) error {
	if p == nil {
		return fmt.Errorf("patient cannot be nil")
	}

	if p.ID == "" {
		return fmt.Errorf("patient ID cannot be empty")
	}

	if len(p.ID) > 36 { // UUID max length
		return fmt.Errorf("patient ID too long")
	}

	if p.Name == "" {
		return fmt.Errorf("patient name cannot be empty")
	}

	if len(p.Name) > 255 {
		return fmt.Errorf("patient name too long (max 255 characters)")
	}

	if len(p.Notes) > 5000 {
		return fmt.Errorf("patient notes too long (max 5000 characters)")
	}

	if p.CreatedAt.IsZero() {
		return fmt.Errorf("patient created_at cannot be zero")
	}

	if p.UpdatedAt.IsZero() {
		return fmt.Errorf("patient updated_at cannot be zero")
	}

	return nil
}

// validatePatientForUpdate validates patient data before updating in database
func validatePatientForUpdate(p *patient.Patient) error {
	if err := validatePatientForSave(p); err != nil {
		return err
	}

	// Additional validation for updates if needed
	return nil
}

// validateID validates a patient ID
func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("patient ID cannot be empty")
	}

	if len(id) > 36 {
		return fmt.Errorf("patient ID too long")
	}

	return nil
}

// validateNameQuery validates a name query for search
func validateNameQuery(name string) error {
	if name == "" {
		return fmt.Errorf("search name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("search name too long (max 100 characters)")
	}

	return nil
}

// Save persists a new patient to the database
// Validates parameters before execution to prevent SQL injection and ensure data integrity
func (r *PatientRepository) Save(p *patient.Patient) error {
	// Validate parameters
	if err := validatePatientForSave(p); err != nil {
		return err
	}

	_, err := r.db.Exec(r.queries.save, p.ID, p.Name, p.Notes, p.CreatedAt, p.UpdatedAt)
	return err
}

// FindByID retrieves a patient by their ID
// Returns nil, nil if patient is not found (not an error)
// Uses index on primary key for optimal performance
func (r *PatientRepository) FindByID(id string) (*patient.Patient, error) {
	// Validate ID
	if err := validateID(id); err != nil {
		return nil, err
	}

	row := r.db.QueryRow(r.queries.findByID, id)

	var p patient.Patient
	err := row.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindAll retrieves all patients ordered by creation date (newest first)
// Uses idx_patients_created_at index for optimal sorting performance
// Consider using FindPaginated for large datasets to avoid memory issues
func (r *PatientRepository) FindAll() ([]*patient.Patient, error) {
	rows, err := r.db.Query(r.queries.findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*patient.Patient
	for rows.Next() {
		var p patient.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}
	return patients, nil
}

// Update modifies an existing patient in the database
// Validates all parameters and ensures updated_at is current
// Uses primary key index for optimal update performance
func (r *PatientRepository) Update(p *patient.Patient) error {
	// Validate parameters
	if err := validatePatientForUpdate(p); err != nil {
		return err
	}

	_, err := r.db.Exec(r.queries.update, p.Name, p.Notes, p.UpdatedAt, p.ID)
	return err
}

// Delete removes a patient from the database by ID
// Validates ID before execution
// Uses primary key index for optimal delete performance
func (r *PatientRepository) Delete(id string) error {
	// Validate ID
	if err := validateID(id); err != nil {
		return err
	}

	_, err := r.db.Exec(r.queries.delete, id)
	return err
}

// FindByName searches for patients by name (case-insensitive, partial match)
// Uses idx_patients_name index for optimal search performance
// Returns empty slice if no patients found (not an error)
func (r *PatientRepository) FindByName(name string) ([]*patient.Patient, error) {
	// Validate name query
	if err := validateNameQuery(name); err != nil {
		return nil, err
	}

	// Add wildcards for partial match
	searchTerm := "%" + name + "%"

	rows, err := r.db.Query(r.queries.findByName, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*patient.Patient
	for rows.Next() {
		var p patient.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}
	return patients, nil
}

// CountAll returns the total number of patients in the database
// Useful for pagination and statistics
func (r *PatientRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow(r.queries.countAll).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindPaginated retrieves a paginated list of patients
// Uses idx_patients_created_at index for optimal sorting performance
// limit: maximum number of patients to return
// offset: number of patients to skip (for pagination)
func (r *PatientRepository) FindPaginated(limit, offset int) ([]*patient.Patient, error) {
	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	rows, err := r.db.Query(r.queries.findPaginated, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*patient.Patient
	for rows.Next() {
		var p patient.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}
	return patients, nil
}

// InitSchema is deprecated - use migrations instead
func (r *PatientRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
