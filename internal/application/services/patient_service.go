package services

import (
	"context"
	"fmt"
	"strings"

	"arandu/internal/domain/patient"
)

// CreatePatientInput represents the input data for creating a patient
type CreatePatientInput struct {
	Name  string `json:"name"`
	Notes string `json:"notes,omitempty"`
}

// Validate validates the CreatePatientInput
func (input *CreatePatientInput) Validate() error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("patient name cannot be empty or whitespace only")
	}

	if len(input.Name) > 255 {
		return fmt.Errorf("patient name cannot exceed 255 characters")
	}

	if len(input.Notes) > 5000 {
		return fmt.Errorf("patient notes cannot exceed 5000 characters")
	}

	// Additional application-level validation
	// Example: Name cannot contain special characters (except spaces, hyphens, apostrophes)
	for _, r := range input.Name {
		if !isValidNameRune(r) {
			return fmt.Errorf("patient name contains invalid character: %q", r)
		}
	}

	return nil
}

// isValidNameRune checks if a rune is valid in a patient name
func isValidNameRune(r rune) bool {
	// Allow letters (including accented), digits, spaces, hyphens, apostrophes
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == ' ' || r == '-' || r == '\'' ||
		r == 'á' || r == 'é' || r == 'í' || r == 'ó' || r == 'ú' ||
		r == 'Á' || r == 'É' || r == 'Í' || r == 'Ó' || r == 'Ú' ||
		r == 'ã' || r == 'õ' || r == 'â' || r == 'ê' || r == 'î' || r == 'ô' || r == 'û' ||
		r == 'Ã' || r == 'Õ' || r == 'Â' || r == 'Ê' || r == 'Î' || r == 'Ô' || r == 'Û' ||
		r == 'à' || r == 'è' || r == 'ì' || r == 'ò' || r == 'ù' ||
		r == 'À' || r == 'È' || r == 'Ì' || r == 'Ò' || r == 'Ù' ||
		r == 'ç' || r == 'Ç'
}

// Sanitize sanitizes the input data
func (input *CreatePatientInput) Sanitize() {
	input.Name = strings.TrimSpace(input.Name)
	input.Notes = strings.TrimSpace(input.Notes)

	// Normalize multiple spaces to single space
	input.Name = strings.Join(strings.Fields(input.Name), " ")
}

// UpdatePatientInput represents the input data for updating a patient
type UpdatePatientInput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Notes string `json:"notes,omitempty"`
}

// Validate validates the UpdatePatientInput
func (input *UpdatePatientInput) Validate() error {
	if strings.TrimSpace(input.ID) == "" {
		return fmt.Errorf("patient ID cannot be empty")
	}

	if len(input.ID) > 36 {
		return fmt.Errorf("patient ID cannot exceed 36 characters")
	}

	// Reuse CreatePatientInput validation for name and notes
	createInput := CreatePatientInput{
		Name:  input.Name,
		Notes: input.Notes,
	}

	return createInput.Validate()
}

// Sanitize sanitizes the input data
func (input *UpdatePatientInput) Sanitize() {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.Notes = strings.TrimSpace(input.Notes)

	// Normalize multiple spaces to single space
	input.Name = strings.Join(strings.Fields(input.Name), " ")
}

// Application errors
var (
	ErrPatientNotFound      = fmt.Errorf("patient not found")
	ErrInvalidInput         = fmt.Errorf("invalid input")
	ErrPatientAlreadyExists = fmt.Errorf("patient already exists")
	ErrRepository           = fmt.Errorf("repository error")
)

// PatientService implements the application service for patient operations
type PatientService struct {
	repo patient.Repository
}

// NewPatientService creates a new PatientService
func NewPatientService(repo patient.Repository) *PatientService {
	return &PatientService{repo: repo}
}

// CreatePatient creates a new patient with enhanced validation and error handling
// ctx: Context for cancellation and timeout
// input: Validated and sanitized input data
// Returns the created patient or an application-specific error
func (s *PatientService) CreatePatient(ctx context.Context, input CreatePatientInput) (*patient.Patient, error) {
	// Step 1: Sanitize input
	input.Sanitize()

	// Step 2: Validate input at application level
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Step 3: Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Step 4: Check for duplicate patients (optional business rule)
	// This is an example of application-level business logic
	// In a real application, you might want to check for duplicates
	// existingPatients, err := s.repo.FindByName(input.Name)
	// if err != nil {
	//     return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	// }
	// if len(existingPatients) > 0 {
	//     return nil, ErrPatientAlreadyExists
	// }

	// Step 5: Create domain entity
	p, err := patient.NewPatient(input.Name, input.Notes)
	if err != nil {
		// Domain validation error
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Step 6: Persist to repository
	if err := s.repo.Save(ctx, p); err != nil {
		// Repository error
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	// Step 7: Return created patient
	return p, nil
}

// CreatePatientLegacy provides backward compatibility with the old signature
// Deprecated: Use CreatePatient with context and input model instead
func (s *PatientService) CreatePatientLegacy(name, notes string) (*patient.Patient, error) {
	input := CreatePatientInput{
		Name:  name,
		Notes: notes,
	}

	// Use background context for legacy method
	return s.CreatePatient(context.Background(), input)
}

// GetPatientByID retrieves a patient by ID with enhanced error handling
func (s *PatientService) GetPatientByID(ctx context.Context, id string) (*patient.Patient, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate ID
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: patient ID cannot be empty", ErrInvalidInput)
	}

	// Retrieve from repository
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	if p == nil {
		return nil, ErrPatientNotFound
	}

	return p, nil
}

// GetPatientLegacy provides backward compatibility
// Deprecated: Use GetPatientByID with context instead
func (s *PatientService) GetPatientLegacy(id string) (*patient.Patient, error) {
	return s.GetPatientByID(context.Background(), id)
}

// ListPatients retrieves all patients with enhanced error handling
func (s *PatientService) ListPatients(ctx context.Context) ([]*patient.Patient, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Retrieve from repository
	patients, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return patients, nil
}

// ListPatientsLegacy provides backward compatibility
// Deprecated: Use ListPatients with context instead
func (s *PatientService) ListPatientsLegacy() ([]*patient.Patient, error) {
	return s.ListPatients(context.Background())
}

// UpdatePatient updates an existing patient with enhanced validation and error handling
func (s *PatientService) UpdatePatient(ctx context.Context, input UpdatePatientInput) error {
	// Step 1: Sanitize input
	input.Sanitize()

	// Step 2: Validate input at application level
	if err := input.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Step 3: Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Step 4: Retrieve existing patient
	p, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepository, err)
	}

	if p == nil {
		return ErrPatientNotFound
	}

	// Step 5: Update domain entity
	if err := p.Update(input.Name, input.Notes); err != nil {
		// Domain validation error
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Step 6: Persist changes
	if err := s.repo.Update(ctx, p); err != nil {
		return fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return nil
}

// UpdatePatientLegacy provides backward compatibility
// Deprecated: Use UpdatePatient with context and input model instead
func (s *PatientService) UpdatePatientLegacy(id, name, notes string) error {
	input := UpdatePatientInput{
		ID:    id,
		Name:  name,
		Notes: notes,
	}

	return s.UpdatePatient(context.Background(), input)
}

// DeletePatient deletes a patient by ID with enhanced error handling
func (s *PatientService) DeletePatient(ctx context.Context, id string) error {
	// Check context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Validate ID
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: patient ID cannot be empty", ErrInvalidInput)
	}

	// Optional: Check if patient exists before deleting
	// This is an application-level decision - you might want to skip this
	// to avoid an extra database query
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepository, err)
	}

	if p == nil {
		return ErrPatientNotFound
	}

	// Delete from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return nil
}

// DeletePatientLegacy provides backward compatibility
// Deprecated: Use DeletePatient with context instead
func (s *PatientService) DeletePatientLegacy(id string) error {
	return s.DeletePatient(context.Background(), id)
}

// SearchPatientsByName searches for patients by name (case-insensitive, partial match)
func (s *PatientService) SearchPatientsByName(ctx context.Context, name string) ([]*patient.Patient, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate search term
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%w: search name cannot be empty", ErrInvalidInput)
	}

	if len(name) > 100 {
		return nil, fmt.Errorf("%w: search name cannot exceed 100 characters", ErrInvalidInput)
	}

	// Search in repository
	patients, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return patients, nil
}

// SearchPatients searches for patients with pagination
// Returns a slice of patients and the total count
func (s *PatientService) SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Empty query returns empty slice
	query = strings.TrimSpace(query)
	if query == "" {
		return []*patient.Patient{}, nil
	}

	if len(query) > 100 {
		return nil, fmt.Errorf("%w: search query cannot exceed 100 characters", ErrInvalidInput)
	}

	// Default limit
	if limit < 1 || limit > 100 {
		limit = 15
	}

	if offset < 0 {
		offset = 0
	}

	// Search in repository with pagination
	patients, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return patients, nil
}

// GetPatientCount returns the total number of patients
func (s *PatientService) GetPatientCount(ctx context.Context) (int, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	// Get count from repository
	count, err := s.repo.CountAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return count, nil
}

// ListPatientsPaginated retrieves a paginated list of patients
func (s *PatientService) ListPatientsPaginated(ctx context.Context, page, pageSize int) ([]*patient.Patient, int, error) {
	// Check context cancellation
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	// Validate pagination parameters
	if page < 1 {
		return nil, 0, fmt.Errorf("%w: page must be at least 1", ErrInvalidInput)
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, 0, fmt.Errorf("%w: page size must be between 1 and 100", ErrInvalidInput)
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	patients, err := s.repo.FindPaginated(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	// Get total count for pagination metadata
	total, err := s.repo.CountAll(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return patients, total, nil
}

// GetThemeFrequency retrieves the most common terms from a patient's records
func (s *PatientService) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if patientID == "" {
		return nil, fmt.Errorf("%w: patient ID cannot be empty", ErrInvalidInput)
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	themes, err := s.repo.GetThemeFrequency(ctx, patientID, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRepository, err)
	}

	return themes, nil
}
