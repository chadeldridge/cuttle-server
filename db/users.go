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

type Users struct {
	DB
}

type UserData struct {
	ID       int
	Username string
	Password string
	Groups   string
}

func NewUsers(db DB) (Users, error) {
	r := Users{
		DB: db,
	}

	if db == nil {
		return r, fmt.Errorf("db.NewUsers: db is nil")
	}

	udb := NewSqliteDB(udb_file)
	err := udb.Open()
	if err != nil {
		return r, fmt.Errorf("db.NewUsers: failed to open db: %w", err)
	}
	defer udb.Close()

	udbMigrate(udb)
	err = r.Attach(udb_file, udb_alias)
	if err != nil {
		return r, fmt.Errorf("db.NewUsers: failed to attach udb: %w", err)
	}

	return r, nil
}

func udbMigrate(db DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + udb_name + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
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

func (r Users) Create(username, pwHash, groups string) error {
	if username == "" {
		return fmt.Errorf("db.Users.Create: %w", ErrInvalidName)
	}

	if !r.IsUnique(username) {
		return fmt.Errorf("db.Users.Create: %w", ErrDuplicateEntry)
	}

	if groups == "" {
		return fmt.Errorf("db.Users.Create: groups is empty")
	}

	query := `INSERT INTO ` + udb_ref + ` (username, password, groups) VALUES (?, ?, ?)`
	if err := r.Exec(query, username, pwHash, groups); err != nil {
		return fmt.Errorf("db.Users.Create: %w", err)
	}

	return nil
}

func (r Users) GetByName(username string) (UserData, error) {
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
	query := `UPDATE ` + udb_ref + ` SET username = ?, password = ?, groups = ? WHERE id = ?`
	if err := r.Exec(query, data.Username, data.Password, data.Groups, data.ID); err != nil {
		return fmt.Errorf("db.Users.Update: %w", err)
	}

	return nil
}

func (r Users) Delete(id int) error {
	query := `DELETE FROM ` + udb_ref + ` WHERE id = ? LIMIT 1`
	if err := r.Exec(query, id); err != nil {
		return fmt.Errorf("db.Users.Delete: %w", err)
	}

	return nil
}
