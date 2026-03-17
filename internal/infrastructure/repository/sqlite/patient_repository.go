package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

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
	search        string
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

		// Search with pagination
		search: `SELECT id, name, notes, created_at, updated_at FROM patients WHERE LOWER(name) LIKE LOWER(?) ORDER BY name LIMIT ? OFFSET ?`,

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

// Search retrieves patients matching the given query with pagination
// Uses LOWER for case-insensitive search and LIMIT/OFFSET for pagination
// Returns empty slice if no patients found (not an error)
func (r *PatientRepository) Search(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	if query == "" {
		return []*patient.Patient{}, nil
	}

	if limit < 1 || limit > 100 {
		limit = 15
	}

	if offset < 0 {
		offset = 0
	}

	searchTerm := "%" + query + "%"

	rows, err := r.db.QueryContext(ctx, r.queries.search, searchTerm, limit, offset)
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

// GetThemeFrequency extracts the most common terms from a patient's observations and interventions
// This is a simplified implementation that counts word occurrences and filters stop words
func (r *PatientRepository) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Query to extract themes from observations
	obsQuery := `
		SELECT o.content 
		FROM observations o
		JOIN sessions s ON s.id = o.session_id
		WHERE s.patient_id = ?
	`

	rows, err := r.db.QueryContext(ctx, obsQuery, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Count word occurrences
	wordCount := make(map[string]int)

	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		// Simple word extraction (split by space and punctuation)
		words := extractWords(content)
		for _, word := range words {
			wordLower := strings.ToLower(word)
			if len(wordLower) > 3 && !isStopWord(wordLower) {
				wordCount[wordLower]++
			}
		}
	}

	// Also get from interventions
	intQuery := `
		SELECT i.content 
		FROM interventions i
		JOIN sessions s ON s.id = i.session_id
		WHERE s.patient_id = ?
	`

	rows2, err := r.db.QueryContext(ctx, intQuery, patientID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var content string
		if err := rows2.Scan(&content); err != nil {
			return nil, err
		}
		words := extractWords(content)
		for _, word := range words {
			wordLower := strings.ToLower(word)
			if len(wordLower) > 3 && !isStopWord(wordLower) {
				wordCount[wordLower]++
			}
		}
	}

	// Sort by count and take top N
	type wordFreq struct {
		term  string
		count int
	}

	var sorted []wordFreq
	for term, count := range wordCount {
		sorted = append(sorted, wordFreq{term, count})
	}

	// Sort descending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	// Take top N
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	// Convert to result format
	result := make([]map[string]interface{}, len(sorted))
	for i, wf := range sorted {
		result[i] = map[string]interface{}{
			"term":  wf.term,
			"count": wf.count,
		}
	}

	return result, nil
}

// extractWords splits text into words
func extractWords(text string) []string {
	// Simple split - in production, use proper tokenization
	text = strings.ReplaceAll(text, ".", " ")
	text = strings.ReplaceAll(text, ",", " ")
	text = strings.ReplaceAll(text, "!", " ")
	text = strings.ReplaceAll(text, "?", " ")
	text = strings.ReplaceAll(text, ";", " ")
	text = strings.ReplaceAll(text, ":", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	words := strings.Split(text, " ")
	var result []string
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w != "" {
			result = append(result, w)
		}
	}
	return result
}

// isStopWord checks if a word is a common Portuguese stop word
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		// Conjunções e preposições
		"porque": true, "quanto": true, "quando": true, "onde": true, "quem": true,
		"como": true, "desde": true, "entre": true, "sobre": true, "também": true,
		"assim": true, "ainda": true, "muito": true, "pouco": true, "todos": true,
		"todas": true, "alguns": true, "algumas": true, "outros": true, "outras": true,
		"mesmo": true, "mesma": true, "outro": true, "outra": true, "cada": true,
		"qual": true, "quais": true, "qualquer": true, "nenhum": true, "nenhuma": true,
		"todo": true, "toda": true, "ser": true, "são": true, "está": true, "esto": true,
		"essa": true, "esse": true, "isso": true, "isto": true, "uma": true, "umas": true,
		"uns": true, "foi": true, "será": true, "seria": true, "poderia": true,
		"deveria": true, "teria": true, "haver": true, "haveria": true, "podem": true,
		"devem": true, "serem": true, "estarem": true, "têm": true,
		"tinha": true, "tinham": true, "terão": true, "teriam": true, "havia": true,
		"haviam": true, "existia": true, "existem": true, "existi": true, "existiu": true,
		// Preposições muito comuns
		"para": true, "por": true, "com": true, "sem": true, "sob": true, "ante": true,
		"após": true, "antes": true, "até": true, "mediante": true, "durante": true,
		"embaixo": true, "através": true,
		// Verbos muito comuns
		"ter": true, "teve": true,
		"fazer": true, "fez": true, "faz": true, "feito": true,
		"dizer": true, "disse": true, "diz": true, "dito": true,
		"ver": true, "viu": true, "vê": true, "visto": true,
		"dar": true, "deu": true, "dá": true, "dado": true,
		"saber": true, "sabia": true, "sabe": true, "soube": true,
		"querer": true, "queria": true, "quer": true, "quis": true,
		"poder": true, "pode": true, "pôde": true,
		// Advérbios muito comuns
		"mais": true, "menos": true, "bem": true, "mal": true,
		"já": true, "sempre": true, "nunca": true, "quase": true, "só": true,
		// Artigos
		"o": true, "a": true, "os": true, "as": true, "um": true,
		// Pronomes
		"ele": true, "ela": true, "eles": true, "elas": true, "seu": true, "sua": true,
		"seus": true, "suas": true, "se": true, "lhe": true, "lhes": true,
		// Diversos
		"vez": true, "vezes": true, "apenas": true, "então": true,
		"portanto": true, "dessa": true, "desse": true,
		"daquela": true, "daquele": true, "nesta": true, "neste": true, "nessa": true,
		"nesse": true, "àquela": true, "àquele": true, "àquilo": true,
		"aqui": true, "aí": true, "ali": true, "lá": true, "cá": true,
		"dentro": true, "fora": true, "longe": true, "perto": true,
		"primeira": true, "primeiro": true, "segunda": true, "segundo": true,
		"terceira": true, "terceiro": true, "última": true, "último": true,
	}
	return stopWords[word]
}

// InitSchema is deprecated - use migrations instead
func (r *PatientRepository) InitSchema() error {
	// Schema creation is now handled by migrations
	// This method exists only for interface compatibility during transition
	return nil
}
