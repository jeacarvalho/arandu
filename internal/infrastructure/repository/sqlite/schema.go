package sqlite

func (db *DB) InitSchema() error {
	repos := []interface{ InitSchema() error }{
		NewPatientRepository(db),
		NewSessionRepository(db),
		NewObservationRepository(db),
		NewInterventionRepository(db),
		NewInsightRepository(db),
	}

	for _, repo := range repos {
		if err := repo.InitSchema(); err != nil {
			return err
		}
	}
	return nil
}
