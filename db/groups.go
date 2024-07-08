package db

import "fmt"

type Groups struct {
	DB
}

func NewGroups(db DB) (Groups, error) {
	r := Groups{db}
	if db == nil {
		return r, fmt.Errorf("db.NewGroups: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r Groups) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name varchar(255) NOT NULL
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Groups.Exec: %w", err)
	}

	// Create the table for the Group to Server relation.
	query = `
	CREATE TABLE IF NOT EXISTS group_servers (
		group_id INTEGER NOT NULL,
		server_id INTEGER NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Groups.Exec: %w", err)
	}

	return nil
}
