package db

import (
	"fmt"
	"log"
)

const (
	amdb_file  = "auth_methods.db"
	amdb_name  = "auth_methods"
	amdb_alias = "am"
	amdb_ref   = amdb_alias + "." + amdb_name
)

var (
	ErrInvalidAuthType   = fmt.Errorf("invalid auth type")
	ErrInvalidName       = fmt.Errorf("invalid name")
	ErrInvalidPassphrase = fmt.Errorf("invalid passphrase")
)

type AuthMethods struct {
	DB
}

type AuthMethodData struct {
	ID       int
	Name     string
	AuthType string
	Data     string
}

// NewAuthMethods attaches the auth_methods database to the current database. It first opens the
// auth_methods database, migrates the database, closes the instance, and then attaches it to the
// current database.
func NewAuthMethods(db DB) (AuthMethods, error) {
	r := AuthMethods{
		DB: db,
	}

	if db == nil {
		return r, fmt.Errorf("db.NewAuthMethods: db is nil")
	}

	amdb := NewSqliteDB(amdb_file)
	err := amdb.Open()
	if err != nil {
		return r, fmt.Errorf("db.NewAuthMethods: failed to open db: %w", err)
	}
	defer amdb.Close()

	amdbMigrate(amdb)
	// Attach the auth_methods database to the current database so we can perform joins.
	err = r.Attach(amdb_file, amdb_alias)
	if err != nil {
		return r, fmt.Errorf("db.NewAuthMethods: failed to attach amdb: %w", err)
	}

	return r, nil
}

// Migrate creates the auth_methods table if it does not exist.
func amdbMigrate(db DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + amdb_name + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		auth_type VARCHAR(32) NOT NULL,
		data TEXT NOT NULL
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("db.AuthMethods.Migrate: %w", err)
	}

	return nil
}

func (r AuthMethods) IsUnique(name string) bool {
	query := `SELECT COUNT(*) FROM ? WHERE name = ?`
	row := r.QueryRow(query, amdb_ref, name)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatalf("db.AuthMethods.IsUnique: %s", err)
	}

	return count == 0
}

func (r AuthMethods) Create(data AuthMethodData) error {
	if data.Name == "" {
		return ErrInvalidName
	}

	if !r.IsUnique(data.Name) {
		return ErrDuplicateEntry
	}

	query := `INSERT INTO ? (name, data) VALUES (?, ?)`
	if err := r.Exec(query, amdb_ref, data.Name, data.Data); err != nil {
		return fmt.Errorf("db.AuthMethods.Create: %w", err)
	}

	return nil
}

func (r AuthMethods) Read(id int) (AuthMethodData, error) {
	var data AuthMethodData
	query := `SELECT name, data FROM ? WHERE id = ?`

	row := r.QueryRow(query, amdb_ref, id)
	err := row.Scan(&data.ID, &data.Name, &data.AuthType, &data.Data)
	if err != nil {
		return data, fmt.Errorf("db.AuthMethods.Read: %w", err)
	}

	return data, nil
}

func (r AuthMethods) Update(data AuthMethodData) error {
	query := `UPDATE ? SET name = ?, auth_type = ?, data = ? WHERE id = ?`
	if err := r.Exec(query, amdb_ref, data.Name, data.AuthType, data.Data, data.ID); err != nil {
		return fmt.Errorf("db.AuthMethods.Update: %w", err)
	}

	return nil
}

func (r AuthMethods) Delete(id int) error {
	query := `DELETE FROM ? WHERE id = ?`
	if err := r.Exec(query, amdb_ref, id); err != nil {
		return fmt.Errorf("db.AuthMethods.Delete: %w", err)
	}

	return nil
}
