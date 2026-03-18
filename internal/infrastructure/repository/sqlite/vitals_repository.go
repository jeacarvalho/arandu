package sqlite

import (
	"database/sql"
	"time"

	"arandu/internal/domain/patient"
)

type VitalsRepository struct {
	db *DB
}

func NewVitalsRepository(db *DB) *VitalsRepository {
	return &VitalsRepository{db: db}
}

func (r *VitalsRepository) Save(v *patient.Vitals) error {
	query := `INSERT INTO patient_vitals (id, patient_id, date, sleep_hours, appetite_level, weight, physical_activity, notes, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, v.ID, v.PatientID, v.Date, v.SleepHours, v.AppetiteLevel, v.Weight, v.PhysicalActivity, v.Notes, v.CreatedAt, v.UpdatedAt)
	return err
}

func (r *VitalsRepository) FindByID(id string) (*patient.Vitals, error) {
	query := `SELECT id, patient_id, date, sleep_hours, appetite_level, weight, physical_activity, notes, created_at, updated_at 
			  FROM patient_vitals WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var v patient.Vitals
	var sleepHours, weight sql.NullFloat64
	var appetiteLevel sql.NullInt64
	var notes sql.NullString
	err := row.Scan(&v.ID, &v.PatientID, &v.Date, &sleepHours, &appetiteLevel, &weight, &v.PhysicalActivity, &notes, &v.CreatedAt, &v.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if sleepHours.Valid {
		v.SleepHours = &sleepHours.Float64
	}
	if appetiteLevel.Valid {
		lvl := int(appetiteLevel.Int64)
		v.AppetiteLevel = &lvl
	}
	if weight.Valid {
		v.Weight = &weight.Float64
	}
	if notes.Valid {
		v.Notes = notes.String
	}
	return &v, nil
}

func (r *VitalsRepository) FindByPatientID(patientID string, limit int) ([]*patient.Vitals, error) {
	if limit <= 0 {
		limit = 30
	}
	query := `SELECT id, patient_id, date, sleep_hours, appetite_level, weight, physical_activity, notes, created_at, updated_at 
			  FROM patient_vitals WHERE patient_id = ? ORDER BY date DESC LIMIT ?`
	rows, err := r.db.Query(query, patientID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanVitals(rows)
}

func (r *VitalsRepository) GetLatestVitals(patientID string) (*patient.Vitals, error) {
	query := `SELECT id, patient_id, date, sleep_hours, appetite_level, weight, physical_activity, notes, created_at, updated_at 
			  FROM patient_vitals WHERE patient_id = ? ORDER BY date DESC LIMIT 1`
	row := r.db.QueryRow(query, patientID)

	var v patient.Vitals
	var sleepHours, weight sql.NullFloat64
	var appetiteLevel sql.NullInt64
	var notes sql.NullString
	err := row.Scan(&v.ID, &v.PatientID, &v.Date, &sleepHours, &appetiteLevel, &weight, &v.PhysicalActivity, &notes, &v.CreatedAt, &v.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if sleepHours.Valid {
		v.SleepHours = &sleepHours.Float64
	}
	if appetiteLevel.Valid {
		lvl := int(appetiteLevel.Int64)
		v.AppetiteLevel = &lvl
	}
	if weight.Valid {
		v.Weight = &weight.Float64
	}
	if notes.Valid {
		v.Notes = notes.String
	}
	return &v, nil
}

func (r *VitalsRepository) GetAverageVitals(patientID string, days int) (*VitalsAverage, error) {
	if days <= 0 {
		days = 30
	}
	fromDate := time.Now().AddDate(0, 0, -days)

	query := `SELECT 
				AVG(sleep_hours) as avg_sleep,
				AVG(appetite_level) as avg_appetite,
				AVG(weight) as avg_weight,
				AVG(physical_activity) as avg_activity,
				COUNT(*) as count
			  FROM patient_vitals 
			  WHERE patient_id = ? AND date >= ?`

	row := r.db.QueryRow(query, patientID, fromDate)

	var avg VitalsAverage
	var avgSleep, avgAppetite, avgWeight, avgActivity sql.NullFloat64
	err := row.Scan(&avgSleep, &avgAppetite, &avgWeight, &avgActivity, &avg.Count)
	if err == sql.ErrNoRows {
		return &avg, nil
	}
	if err != nil {
		return nil, err
	}
	if avgSleep.Valid {
		avg.AverageSleepHours = &avgSleep.Float64
	}
	if avgAppetite.Valid {
		avg.AverageAppetiteLevel = &avgAppetite.Float64
	}
	if avgWeight.Valid {
		avg.AverageWeight = &avgWeight.Float64
	}
	if avgActivity.Valid {
		avg.AveragePhysicalActivity = &avgActivity.Float64
	}
	return &avg, nil
}

func (r *VitalsRepository) Update(v *patient.Vitals) error {
	query := `UPDATE patient_vitals SET sleep_hours = ?, appetite_level = ?, weight = ?, physical_activity = ?, notes = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, v.SleepHours, v.AppetiteLevel, v.Weight, v.PhysicalActivity, v.Notes, time.Now(), v.ID)
	return err
}

func (r *VitalsRepository) Delete(id string) error {
	query := `DELETE FROM patient_vitals WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *VitalsRepository) scanVitals(rows *sql.Rows) ([]*patient.Vitals, error) {
	var vitals []*patient.Vitals
	for rows.Next() {
		var v patient.Vitals
		var sleepHours, weight sql.NullFloat64
		var appetiteLevel sql.NullInt64
		var notes sql.NullString
		if err := rows.Scan(&v.ID, &v.PatientID, &v.Date, &sleepHours, &appetiteLevel, &weight, &v.PhysicalActivity, &notes, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		if sleepHours.Valid {
			v.SleepHours = &sleepHours.Float64
		}
		if appetiteLevel.Valid {
			lvl := int(appetiteLevel.Int64)
			v.AppetiteLevel = &lvl
		}
		if weight.Valid {
			v.Weight = &weight.Float64
		}
		if notes.Valid {
			v.Notes = notes.String
		}
		vitals = append(vitals, &v)
	}
	return vitals, nil
}

type VitalsAverage struct {
	AverageSleepHours       *float64
	AverageAppetiteLevel    *float64
	AverageWeight           *float64
	AveragePhysicalActivity *float64
	Count                   int
}
