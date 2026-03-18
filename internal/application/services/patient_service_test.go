package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"arandu/internal/domain/patient"
)

// mockPatientRepository is a mock implementation of patient.Repository for testing
type mockPatientRepository struct {
	saveFunc              func(p *patient.Patient) error
	findByIDFunc          func(id string) (*patient.Patient, error)
	findAllFunc           func() ([]*patient.Patient, error)
	updateFunc            func(p *patient.Patient) error
	deleteFunc            func(id string) error
	findByNameFunc        func(name string) ([]*patient.Patient, error)
	searchFunc            func(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error)
	countAllFunc          func() (int, error)
	findPaginatedFunc     func(limit, offset int) ([]*patient.Patient, error)
	getThemeFrequencyFunc func(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error)
}

func (m *mockPatientRepository) Save(p *patient.Patient) error {
	if m.saveFunc != nil {
		return m.saveFunc(p)
	}
	return nil
}

func (m *mockPatientRepository) FindByID(id string) (*patient.Patient, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, nil
}

func (m *mockPatientRepository) FindAll() ([]*patient.Patient, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return nil, nil
}

func (m *mockPatientRepository) Update(p *patient.Patient) error {
	if m.updateFunc != nil {
		return m.updateFunc(p)
	}
	return nil
}

func (m *mockPatientRepository) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockPatientRepository) FindByName(name string) ([]*patient.Patient, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return nil, nil
}

func (m *mockPatientRepository) Search(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query, limit, offset)
	}
	return nil, nil
}

func (m *mockPatientRepository) CountAll() (int, error) {
	if m.countAllFunc != nil {
		return m.countAllFunc()
	}
	return 0, nil
}

func (m *mockPatientRepository) FindPaginated(limit, offset int) ([]*patient.Patient, error) {
	if m.findPaginatedFunc != nil {
		return m.findPaginatedFunc(limit, offset)
	}
	return nil, nil
}

func (m *mockPatientRepository) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	if m.getThemeFrequencyFunc != nil {
		return m.getThemeFrequencyFunc(ctx, patientID, limit)
	}
	return nil, nil
}

func TestPatientService_CreatePatient(t *testing.T) {
	tests := []struct {
		name      string
		input     CreatePatientInput
		mockSave  func(p *patient.Patient) error
		wantError bool
		errorType error
	}{
		{
			name: "Valid input creates patient",
			input: CreatePatientInput{
				Name:  "John Doe",
				Notes: "Test patient",
			},
			mockSave: func(p *patient.Patient) error {
				// Simulate successful save
				return nil
			},
			wantError: false,
		},
		{
			name: "Empty name returns error",
			input: CreatePatientInput{
				Name:  "",
				Notes: "Test patient",
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Name with only whitespace returns error",
			input: CreatePatientInput{
				Name:  "   ",
				Notes: "Test patient",
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Name too long returns error",
			input: CreatePatientInput{
				Name:  string(make([]byte, 256)), // 256 characters
				Notes: "Test patient",
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Notes too long returns error",
			input: CreatePatientInput{
				Name:  "John Doe",
				Notes: string(make([]byte, 5001)), // 5001 characters
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Invalid characters in name returns error",
			input: CreatePatientInput{
				Name:  "John@Doe", // @ is not allowed
				Notes: "Test patient",
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Repository error returns wrapped error",
			input: CreatePatientInput{
				Name:  "John Doe",
				Notes: "Test patient",
			},
			mockSave: func(p *patient.Patient) error {
				return errors.New("database connection failed")
			},
			wantError: true,
			errorType: ErrRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				saveFunc: tt.mockSave,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			result, err := service.CreatePatient(ctx, tt.input)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check result
				if result == nil {
					t.Error("Expected patient but got nil")
				}

				if result.Name != tt.input.Name {
					t.Errorf("Expected patient name %q, got %q", tt.input.Name, result.Name)
				}

				if result.Notes != tt.input.Notes {
					t.Errorf("Expected patient notes %q, got %q", tt.input.Notes, result.Notes)
				}

				if result.ID == "" {
					t.Error("Expected patient to have ID")
				}
			}
		})
	}
}

func TestPatientService_CreatePatient_ContextCancellation(t *testing.T) {
	// Create mock repository
	mockRepo := &mockPatientRepository{}

	// Create service
	service := NewPatientService(mockRepo)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Call method
	input := CreatePatientInput{
		Name:  "John Doe",
		Notes: "Test patient",
	}

	_, err := service.CreatePatient(ctx, input)

	// Should return context error
	if err == nil {
		t.Error("Expected context cancellation error but got none")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestPatientService_GetPatientByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockFindByID  func(id string) (*patient.Patient, error)
		wantError     bool
		errorType     error
		wantPatientID string
	}{
		{
			name: "Valid ID returns patient",
			id:   "test-id-123",
			mockFindByID: func(id string) (*patient.Patient, error) {
				return &patient.Patient{
					ID:   id,
					Name: "John Doe",
				}, nil
			},
			wantError:     false,
			wantPatientID: "test-id-123",
		},
		{
			name:      "Empty ID returns error",
			id:        "",
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name: "Patient not found returns error",
			id:   "non-existent-id",
			mockFindByID: func(id string) (*patient.Patient, error) {
				return nil, nil // Repository returns nil, nil for not found
			},
			wantError: true,
			errorType: ErrPatientNotFound,
		},
		{
			name: "Repository error returns wrapped error",
			id:   "test-id-123",
			mockFindByID: func(id string) (*patient.Patient, error) {
				return nil, errors.New("database error")
			},
			wantError: true,
			errorType: ErrRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				findByIDFunc: tt.mockFindByID,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			result, err := service.GetPatientByID(ctx, tt.id)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check result
				if result == nil {
					t.Error("Expected patient but got nil")
				}

				if result.ID != tt.wantPatientID {
					t.Errorf("Expected patient ID %q, got %q", tt.wantPatientID, result.ID)
				}
			}
		})
	}
}

func TestPatientService_UpdatePatient(t *testing.T) {
	existingPatient := &patient.Patient{
		ID:   "existing-id",
		Name: "Old Name",
	}

	tests := []struct {
		name         string
		input        UpdatePatientInput
		mockFindByID func(id string) (*patient.Patient, error)
		mockUpdate   func(p *patient.Patient) error
		wantError    bool
		errorType    error
	}{
		{
			name: "Valid update succeeds",
			input: UpdatePatientInput{
				ID:    "existing-id",
				Name:  "New Name",
				Notes: "Updated notes",
			},
			mockFindByID: func(id string) (*patient.Patient, error) {
				if id == "existing-id" {
					return existingPatient, nil
				}
				return nil, nil
			},
			mockUpdate: func(p *patient.Patient) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Patient not found returns error",
			input: UpdatePatientInput{
				ID:   "non-existent-id",
				Name: "New Name",
			},
			mockFindByID: func(id string) (*patient.Patient, error) {
				return nil, nil
			},
			wantError: true,
			errorType: ErrPatientNotFound,
		},
		{
			name: "Invalid name returns error",
			input: UpdatePatientInput{
				ID:   "existing-id",
				Name: "", // Empty name
			},
			wantError: true,
			errorType: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				findByIDFunc: tt.mockFindByID,
				updateFunc:   tt.mockUpdate,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			err := service.UpdatePatient(ctx, tt.input)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPatientService_DeletePatient(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockFindByID func(id string) (*patient.Patient, error)
		mockDelete   func(id string) error
		wantError    bool
		errorType    error
	}{
		{
			name: "Valid delete succeeds",
			id:   "existing-id",
			mockFindByID: func(id string) (*patient.Patient, error) {
				if id == "existing-id" {
					return &patient.Patient{ID: id}, nil
				}
				return nil, nil
			},
			mockDelete: func(id string) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Patient not found returns error",
			id:   "non-existent-id",
			mockFindByID: func(id string) (*patient.Patient, error) {
				return nil, nil
			},
			wantError: true,
			errorType: ErrPatientNotFound,
		},
		{
			name:      "Empty ID returns error",
			id:        "",
			wantError: true,
			errorType: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				findByIDFunc: tt.mockFindByID,
				deleteFunc:   tt.mockDelete,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			err := service.DeletePatient(ctx, tt.id)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPatientService_SearchPatientsByName(t *testing.T) {
	tests := []struct {
		name           string
		searchTerm     string
		mockFindByName func(name string) ([]*patient.Patient, error)
		wantError      bool
		errorType      error
		wantCount      int
	}{
		{
			name:       "Valid search returns patients",
			searchTerm: "John",
			mockFindByName: func(name string) ([]*patient.Patient, error) {
				return []*patient.Patient{
					{ID: "1", Name: "John Doe"},
					{ID: "2", Name: "Johnny Smith"},
				}, nil
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name:       "Empty search term returns error",
			searchTerm: "",
			wantError:  true,
			errorType:  ErrInvalidInput,
		},
		{
			name:       "Search term too long returns error",
			searchTerm: string(make([]byte, 101)), // 101 characters
			wantError:  true,
			errorType:  ErrInvalidInput,
		},
		{
			name:       "Repository error returns wrapped error",
			searchTerm: "John",
			mockFindByName: func(name string) ([]*patient.Patient, error) {
				return nil, errors.New("database error")
			},
			wantError: true,
			errorType: ErrRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				findByNameFunc: tt.mockFindByName,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			results, err := service.SearchPatientsByName(ctx, tt.searchTerm)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check results
				if len(results) != tt.wantCount {
					t.Errorf("Expected %d patients, got %d", tt.wantCount, len(results))
				}
			}
		})
	}
}

func TestPatientService_ListPatientsPaginated(t *testing.T) {
	tests := []struct {
		name              string
		page              int
		pageSize          int
		mockFindPaginated func(limit, offset int) ([]*patient.Patient, error)
		mockCountAll      func() (int, error)
		wantError         bool
		errorType         error
		wantPatientsCount int
		wantTotalCount    int
	}{
		{
			name:     "Valid pagination returns results",
			page:     1,
			pageSize: 10,
			mockFindPaginated: func(limit, offset int) ([]*patient.Patient, error) {
				// Return 5 patients for page 1
				return []*patient.Patient{
					{ID: "1", Name: "Patient 1"},
					{ID: "2", Name: "Patient 2"},
					{ID: "3", Name: "Patient 3"},
					{ID: "4", Name: "Patient 4"},
					{ID: "5", Name: "Patient 5"},
				}, nil
			},
			mockCountAll: func() (int, error) {
				return 25, nil // Total 25 patients
			},
			wantError:         false,
			wantPatientsCount: 5,
			wantTotalCount:    25,
		},
		{
			name:      "Page less than 1 returns error",
			page:      0,
			pageSize:  10,
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name:      "Page size less than 1 returns error",
			page:      1,
			pageSize:  0,
			wantError: true,
			errorType: ErrInvalidInput,
		},
		{
			name:      "Page size greater than 100 returns error",
			page:      1,
			pageSize:  101,
			wantError: true,
			errorType: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &mockPatientRepository{
				findPaginatedFunc: tt.mockFindPaginated,
				countAllFunc:      tt.mockCountAll,
			}

			// Create service
			service := NewPatientService(mockRepo)

			// Create context
			ctx := context.Background()

			// Call method
			patients, total, err := service.ListPatientsPaginated(ctx, tt.page, tt.pageSize)

			// Check error
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				// Check error type
				if tt.errorType != nil {
					if !errors.Is(err, tt.errorType) {
						t.Errorf("Expected error type %v, got %v", tt.errorType, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check results
				if len(patients) != tt.wantPatientsCount {
					t.Errorf("Expected %d patients, got %d", tt.wantPatientsCount, len(patients))
				}

				if total != tt.wantTotalCount {
					t.Errorf("Expected total count %d, got %d", tt.wantTotalCount, total)
				}
			}
		})
	}
}

func TestPatientService_ContextCancellation(t *testing.T) {
	// Create mock repository
	mockRepo := &mockPatientRepository{
		findAllFunc: func() ([]*patient.Patient, error) {
			return []*patient.Patient{{ID: "1", Name: "Test"}}, nil
		},
	}

	// Create service
	service := NewPatientService(mockRepo)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Call method - should return cancellation error
	_, err := service.ListPatients(ctx)

	if err == nil {
		t.Error("Expected cancellation error but got none")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestInputValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     CreatePatientInput
		wantError bool
	}{
		{
			name: "Valid name with accented characters",
			input: CreatePatientInput{
				Name:  "João da Silva",
				Notes: "Patient with accented name",
			},
			wantError: false,
		},
		{
			name: "Valid name with hyphen",
			input: CreatePatientInput{
				Name:  "Anne-Marie",
				Notes: "Patient with hyphenated name",
			},
			wantError: false,
		},
		{
			name: "Valid name with apostrophe",
			input: CreatePatientInput{
				Name:  "O'Connor",
				Notes: "Patient with apostrophe",
			},
			wantError: false,
		},
		{
			name: "Invalid name with special character",
			input: CreatePatientInput{
				Name:  "John@Doe",
				Notes: "Invalid character @",
			},
			wantError: true,
		},
		{
			name: "Invalid name with underscore",
			input: CreatePatientInput{
				Name:  "John_Doe",
				Notes: "Invalid character _",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()

			if tt.wantError {
				if err == nil {
					t.Error("Expected validation error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestInputSanitization(t *testing.T) {
	input := CreatePatientInput{
		Name:  "  John   Doe  ", // Extra spaces
		Notes: "  Test notes  ", // Extra spaces
	}

	// Sanitize
	input.Sanitize()

	// Check results
	if input.Name != "John Doe" {
		t.Errorf("Expected sanitized name 'John Doe', got %q", input.Name)
	}

	if input.Notes != "Test notes" {
		t.Errorf("Expected sanitized notes 'Test notes', got %q", input.Notes)
	}
}

func TestPatientService_CreatePatientLegacy(t *testing.T) {
	repo := &mockPatientRepository{
		saveFunc: func(p *patient.Patient) error {
			return nil
		},
	}
	service := NewPatientService(repo)

	t.Run("Valid input creates patient", func(t *testing.T) {
		p, err := service.CreatePatientLegacy("John Doe", "Test notes")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if p == nil {
			t.Error("Expected patient to be created")
		}
		if p.Name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got %q", p.Name)
		}
		if p.Notes != "Test notes" {
			t.Errorf("Expected notes 'Test notes', got %q", p.Notes)
		}
	})

	t.Run("Invalid name returns error", func(t *testing.T) {
		p, err := service.CreatePatientLegacy("", "Test notes")
		if err == nil {
			t.Error("Expected error for empty name")
		}
		if p != nil {
			t.Error("Expected no patient to be created on error")
		}
	})
}

func TestPatientService_GetPatientLegacy(t *testing.T) {
	expectedPatient := &patient.Patient{
		ID:        "patient-123",
		Name:      "John Doe",
		Notes:     "Test notes",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := &mockPatientRepository{
		findByIDFunc: func(id string) (*patient.Patient, error) {
			if id == "patient-123" {
				return expectedPatient, nil
			}
			return nil, nil
		},
	}
	service := NewPatientService(repo)

	t.Run("Valid ID returns patient", func(t *testing.T) {
		p, err := service.GetPatientLegacy("patient-123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if p != expectedPatient {
			t.Error("Expected patient to be returned")
		}
	})

	t.Run("Empty ID returns error", func(t *testing.T) {
		p, err := service.GetPatientLegacy("")
		if err == nil {
			t.Error("Expected error for empty ID")
		}
		if p != nil {
			t.Error("Expected no patient to be returned on error")
		}
	})

	t.Run("Non-existent patient returns error", func(t *testing.T) {
		p, err := service.GetPatientLegacy("non-existent")
		if err == nil {
			t.Error("Expected error for non-existent patient")
		}
		if p != nil {
			t.Error("Expected no patient to be returned for non-existent ID")
		}
	})
}

func TestPatientService_ListPatientsLegacy(t *testing.T) {
	expectedPatients := []*patient.Patient{
		{
			ID:        "patient-1",
			Name:      "John Doe",
			Notes:     "Notes 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "patient-2",
			Name:      "Jane Smith",
			Notes:     "Notes 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	repo := &mockPatientRepository{
		findAllFunc: func() ([]*patient.Patient, error) {
			return expectedPatients, nil
		},
	}
	service := NewPatientService(repo)

	t.Run("Returns all patients", func(t *testing.T) {
		patients, err := service.ListPatientsLegacy()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(patients) != len(expectedPatients) {
			t.Errorf("Expected %d patients, got %d", len(expectedPatients), len(patients))
		}
	})

	t.Run("Repository error returns wrapped error", func(t *testing.T) {
		expectedErr := errors.New("repository error")
		errorRepo := &mockPatientRepository{
			findAllFunc: func() ([]*patient.Patient, error) {
				return nil, expectedErr
			},
		}
		errorService := NewPatientService(errorRepo)

		patients, err := errorService.ListPatientsLegacy()
		if err == nil {
			t.Error("Expected error from repository")
		}
		if patients != nil {
			t.Error("Expected no patients on repository error")
		}
	})
}

func TestPatientService_UpdatePatientLegacy(t *testing.T) {
	existingPatient := &patient.Patient{
		ID:        "patient-123",
		Name:      "Old Name",
		Notes:     "Old notes",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	var updatedPatient *patient.Patient
	repo := &mockPatientRepository{
		findByIDFunc: func(id string) (*patient.Patient, error) {
			if id == "patient-123" {
				return existingPatient, nil
			}
			return nil, nil
		},
		updateFunc: func(p *patient.Patient) error {
			updatedPatient = p
			return nil
		},
	}
	service := NewPatientService(repo)

	t.Run("Valid update succeeds", func(t *testing.T) {
		err := service.UpdatePatientLegacy("patient-123", "New Name", "New notes")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedPatient == nil {
			t.Error("Expected patient to be updated")
		}
		if updatedPatient.Name != "New Name" {
			t.Errorf("Expected name 'New Name', got %q", updatedPatient.Name)
		}
		if updatedPatient.Notes != "New notes" {
			t.Errorf("Expected notes 'New notes', got %q", updatedPatient.Notes)
		}
		// Note: UpdatePatientLegacy calls UpdatePatient which should update UpdatedAt
		// but we're mocking the repository, so the domain object's UpdatedAt might not be updated
		// in our mock. This is acceptable for unit testing.
	})

	t.Run("Patient not found returns error", func(t *testing.T) {
		err := service.UpdatePatientLegacy("non-existent", "New Name", "New notes")
		if err == nil {
			t.Error("Expected error for non-existent patient")
		}
	})

	t.Run("Invalid name returns error", func(t *testing.T) {
		err := service.UpdatePatientLegacy("patient-123", "", "New notes")
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})
}
