package db

import "fmt"

type Profiles struct {
	DB
}

func NewProfiles(db DB) (Profiles, error) {
	r := Profiles{db}
	if db == nil {
		return r, fmt.Errorf("db.NewProfiles: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r Profiles) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name varchar(255) NOT NULL,
	)`

	if err := r.Migrate(query); err != nil {
		return fmt.Errorf("db.Profiles.Migrate: %w", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS profile_tiles (
		profile_id INTEGER NOT NULL,
		tile_id INTEGER NOT NULL,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
		FOREIGN KEY (tile_id) REFERENCES tiles(id) ON DELETE CASCADE
	)`

	if err := r.Migrate(query); err != nil {
		return fmt.Errorf("db.Profiles.Migrate: %w", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS profile_groups (
		profile_id INTEGER NOT NULL,
		group_id INTEGER NOT NULL,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
	)`

	if err := r.Migrate(query); err != nil {
		return fmt.Errorf("db.Profiles.Migrate: %w", err)
	}

	return nil
}
