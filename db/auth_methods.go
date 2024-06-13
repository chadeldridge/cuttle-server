package db

import "fmt"

type AuthMethods struct {
	DB
}

func NewAuthMethods(db DB) (AuthMethods, error) {
	r := AuthMethods{db}
	if db == nil {
		return r, fmt.Errorf("db.NewAuthMethods: db is nil")
	}

	r.migrate()
	return r, nil
}

func (r AuthMethods) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS auth_methods (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		data TEXT NOT NULL
	)`

	if err := r.Migrate(query); err != nil {
		return fmt.Errorf("db.AuthMethods.Migrate: %w", err)
	}

	return nil
}

func (r AuthMethods) Create() error {
	/*
		query := `INSERT INTO auth_methods (name, data) VALUES (?, ?)`
			if err := r.Exec(query, name, data); err != nil {
				return fmt.Errorf("db.AuthMethods.Create: %w", err)
			}
	*/

	return nil
}
