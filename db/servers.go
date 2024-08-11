package db

/*
import "fmt"

type Servers struct {
	DB
}

func NewServers(db DB) (Servers, error) {
	r := Servers{db}
	if db == nil {
		return r, fmt.Errorf("db.NewServer: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r Servers) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		hostname VARCHAR(255) NOT NULL,
		ip VARCHAR(15),
		port INTEGER NOT NULL DEFAULT 0,
		useIP BOOLEAN NOT NULL DEFAULT FALSE,
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Servers.Exec: %w", err)
	}

	// Create the table for the Server to Connector relation.
	query = `
	CREATE TABLE IF NOT EXISTS server_connectors (
		server_id INTEGER NOT NULL,
		connector_id INTEGER NOT NULL,
		FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
		FOREIGN KEY (connector_id) REFERENCES connectors(id) ON DELETE CASCADE
	)`

	if err := r.Exec(query); err != nil {
		return fmt.Errorf("db.Servers.Exec: %w", err)
	}

	return nil
}
*/
