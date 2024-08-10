package db

import "fmt"

const (
	ugdb_file  = "user_groups.db"
	ugdb_name  = "user_groups"
	ugdb_alias = "ug"
	ugdb_ref   = ugdb_alias + "." + ugdb_name
)

type UserGroups struct {
	DB
}

type UserGroupData struct {
	ID      int
	Name    string
	Members string
}

func NewUserGroups(db DB) (UserGroups, error) {
	r := UserGroups{
		DB: db,
	}

	if db == nil {
		return r, fmt.Errorf("db.NewUserGroups: db is nil")
	}

	ugdb, err := NewSqliteDB(ugdb_file)
	if err != nil {
		return r, fmt.Errorf("db.NewUserGroups: failed to create db: %w", err)
	}

	err = ugdb.Open()
	if err != nil {
		return r, fmt.Errorf("db.NewUserGroups: failed to open db: %w", err)
	}
	defer ugdb.Close()

	ugdbMigrate(ugdb)
	err = r.Attach(ugdb_file, ugdb_alias)
	if err != nil {
		return r, fmt.Errorf("db.NewUserGroups: failed to attach ugdb: %w", err)
	}

	return r, nil
}

func ugdbMigrate(db DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + ugdb_name + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		members TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("db.ugdbMigrate: %w", err)
	}

	return nil
}

func (r UserGroups) IsUnique(name string) bool {
	return r.DB.IsUnique(`SELECT COUNT(*) FROM `+ugdb_ref+` WHERE name = ?`, name)
}

func (r UserGroups) Create(name, members string) error {
	if name == "" {
		return ErrInvalidName
	}

	if !r.IsUnique(name) {
		return ErrDuplicateEntry
	}

	query := `INSERT INTO ` + ugdb_ref + ` (name, members) VALUES (?, ?)`
	return r.Exec(query, name, members)
}

func (r UserGroups) Get(name string) (UserGroupData, error) {
	query := `SELECT id, name, members FROM ` + ugdb_ref + ` WHERE name = ?`
	row := r.QueryRow(query, name)

	var data UserGroupData
	err := row.Scan(&data.ID, &data.Name, &data.Members)
	if err != nil {
		return data, fmt.Errorf("db.UserGroups.Get: %w", err)
	}

	return data, nil
}

func (r UserGroups) GetByID(id int) (UserGroupData, error) {
	query := `SELECT name, members FROM ` + ugdb_ref + ` WHERE id = ?`
	row := r.QueryRow(query, id)

	var data UserGroupData
	err := row.Scan(&data.Name, &data.Members)
	if err != nil {
		return data, fmt.Errorf("db.UserGroups.GetByID: %w", err)
	}

	return data, nil
}

func (r UserGroups) GetGroups(gids []int) ([]UserGroupData, error) {
	query := `SELECT id, name, members FROM ` + ugdb_ref + ` WHERE id IN (?)`
	rows, err := r.Query(query, gids)
	if err != nil {
		return nil, fmt.Errorf("db.UserGroups.GetGroups: %w", err)
	}
	defer rows.Close()

	var groups []UserGroupData
	for rows.Next() {
		var data UserGroupData
		err := rows.Scan(&data.ID, &data.Name, &data.Members)
		if err != nil {
			return nil, fmt.Errorf("db.UserGroups.GetGroups: %w", err)
		}

		groups = append(groups, data)
	}

	return groups, nil
}

func (r UserGroups) Update(data UserGroupData) error {
	query := `UPDATE ` + ugdb_ref + ` SET name = ?, members = ? WHERE id = ?`
	if err := r.Exec(query, data.Name, data.Members, data.ID); err != nil {
		return fmt.Errorf("db.UserGroups.Update: %w", err)
	}

	return nil
}

func (r UserGroups) Delete(name string) error {
	query := `DELETE FROM ` + ugdb_ref + ` WHERE name = ?`
	if err := r.Exec(query, name); err != nil {
		return fmt.Errorf("db.UserGroups.Delete: %w", err)
	}

	return nil
}
