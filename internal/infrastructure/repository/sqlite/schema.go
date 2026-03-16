package sqlite

// InitSchema is deprecated - use Migrate() instead
// Kept for backward compatibility during transition
func (db *DB) InitSchema() error {
	// Simply call Migrate() which now handles schema creation
	return db.Migrate()
}
