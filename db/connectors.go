package db

import "fmt"

type Connectors struct {
	DB
}

func NewConnectors(db DB) (Connectors, error) {
	r := Connectors{db}
	if db == nil {
		return r, fmt.Errorf("db.NewConnectors: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r Connectors) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS connectors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		user VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Connectors.Exec: %w", err)
	}

	// Create the table for the Connector to AuthMethod relation.
	query = `
	CREATE TABLE IF NOT EXISTS connector_authmethods (
		connector_id INTEGER NOT NULL,
		auth_method_id INTEGER NOT NULL,
		FOREIGN KEY (connector_id) REFERENCES connectors(id) ON DELETE CASCADE,
		FOREIGN KEY (auth_method_id) REFERENCES auth_methods(id) ON DELETE CASCADE
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Connectors.Exec: %w", err)
	}

	return nil
}
