package db

import (
	"fmt"
)

const (
	udb_file  = "users.db"
	udb_name  = "users"
	udb_alias = "u"
	udb_ref   = udb_alias + "." + udb_name
)

var ErrInvalidUsername = fmt.Errorf("invalid username")

// Users holds the main database the Users repo is attached to.
type Users struct {
	DB
}

// UserData represents a user in the database.
type UserData struct {
	ID       int
	Username string
	Name     string // Name to show in app.
	Password string
	Groups   string // JSON string of group IDs.
}

// NewUsers attaches the users database to the current database. It first opens the users database,
// creating and migrating it if it does not exist.
func NewUsers(db DB) (Users, error) {
	r := Users{
		DB: db,
	}

	err := db.AddRepo(udb_file, udb_alias, udbMigrate)
	return r, err
}

func udbMigrate(db DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + udb_name + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(32) NOT NULL UNIQUE,
		password VARCHAR(32) NOT NULL,
		groups TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("db.Users.Migrate: %w", err)
	}

	return nil
}

func (r Users) IsUnique(name string) bool {
	return r.DB.IsUnique(`SELECT COUNT(*) FROM `+udb_ref+` WHERE username = ?`, name)
}

func (r Users) Create(username, name, pwHash, groups string) error {
	if username == "" {
		return fmt.Errorf("db.Users.Create: %w", ErrInvalidName)
	}

	if username == "" {
		return fmt.Errorf("db.Users.Create: %w", ErrInvalidUsername)
	}

	if !r.IsUnique(username) {
		return fmt.Errorf("db.Users.Create: %w", ErrDuplicateEntry)
	}

	if groups == "" {
		return fmt.Errorf("db.Users.Create: groups is empty")
	}

	query := `INSERT INTO ` + udb_ref + ` (username, name, password, groups) VALUES (?, ?, ?, ?)`
	if err := r.Exec(query, name, username, pwHash, groups); err != nil {
		return fmt.Errorf("db.Users.Create: %w", err)
	}

	return nil
}

func (r Users) GetByUsername(username string) (UserData, error) {
	query := `SELECT * FROM ` + udb_ref + ` WHERE username = ?`
	row := r.QueryRow(query, username)

	var data UserData
	err := row.Scan(&data.ID, &data.Username, &data.Password, &data.Groups)
	if err != nil {
		return data, fmt.Errorf("db.Users.Get: %w", err)
	}

	return data, nil
}

func (r Users) Get(id int) (UserData, error) {
	query := `SELECT * FROM ` + udb_ref + ` WHERE id = ?`
	row := r.QueryRow(query, id)

	var data UserData
	err := row.Scan(&data.ID, &data.Username, &data.Password, &data.Groups)
	if err != nil {
		return data, fmt.Errorf("db.Users.GetByID: %w", err)
	}

	return data, nil
}

func (r Users) Update(data UserData) error {
	if data.ID == 0 {
		return fmt.Errorf("db.Users.Update: %w", ErrInvalidID)
	}

	query := `UPDATE ` + udb_ref + ` SET username = ?, name = ?, password = ?, groups = ? WHERE id = ? LIMIT 1`
	if err := r.Exec(query, data.Username, data.Name, data.Password, data.Groups, data.ID); err != nil {
		return fmt.Errorf("db.Users.Update: %w", err)
	}

	return nil
}

func (r Users) Delete(id int) error {
	if id == 0 {
		return fmt.Errorf("db.Users.Update: %w", ErrInvalidID)
	}

	query := `DELETE FROM ` + udb_ref + ` WHERE id = ? LIMIT 1`
	if err := r.Exec(query, id); err != nil {
		return fmt.Errorf("db.Users.Delete: %w", err)
	}

	return nil
}
