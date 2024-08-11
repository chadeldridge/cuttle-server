package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// libray has to be imported to register the driver.
	"github.com/chadeldridge/cuttle/core"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Auth Tables.
	sqlite_tb_users       = "users"
	sqlite_tb_user_groups = "user_groups"
)

type SqliteDB struct {
	Name string // DB file name.
	*sql.DB

	ctx    context.Context
	cancel func()
}

// NewSqliteDB creates a new Sqlite3 DB instance.
func NewSqliteDB(filename string) (*SqliteDB, error) {
	if filename == "" {
		return nil, fmt.Errorf("db.NewSqliteDB: filename is empty")
	}

	db := &SqliteDB{Name: filename}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db, nil
}

// Open opens the database. It attempts to create the path to the database file if it does not
// exist, opens the database file, and enables WAL mode and foreign key checks.
func (db *SqliteDB) Open() error {
	if db.Name == "" {
		return fmt.Errorf("no database name provided")
	}

	var err error
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=ON", db_folder+"/"+db.Name)
	db.DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("SqliteDB.Open: %w", err)
	}

	err = db.DB.Ping()
	if err != nil {
		return fmt.Errorf("SqliteDB.Open: failed to ping db: %w", err)
	}

	// TODO: Implement zstd compression. https://phiresky.github.io/blog/2022/sqlite-zstd/

	return nil
}

func (db *SqliteDB) CuttleMigrate() error {
	return nil
}

func (db *SqliteDB) AuthMigrate() error {
	if err := db.UsersMigrate(); err != nil {
		return fmt.Errorf("db.AuthMigrate: failed to migrate %s: %w", sqlite_tb_users, err)
	}

	return nil
}

// IsUnique returns nil if no records exist in the table that match the where clause. If a record
// exists, it returns an ErrRecordExists error.
func (db *SqliteDB) IsUnique(table string, where string, args ...any) error {
	if table == "" {
		return fmt.Errorf("SqliteDB.IsUnique: table is empty")
	}

	if where == "" {
		return fmt.Errorf("SqliteDB.IsUnique: where is empty")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s;", table, where)
	row := db.QueryRow(query, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("SqliteDB.IsUnique: %s", err)
	}

	if count > 0 {
		return fmt.Errorf("SqliteDB.IsUnique: %w", ErrRecordExists)
	}

	return nil
}

func (db *SqliteDB) Close() error {
	if db.DB == nil {
		return nil
	}

	db.cancel()
	return db.DB.Close()
}

/*
// Attach the filename database tot he current database with the given alias.
func (db *SqliteDB) Attach(filename, alias string) error {
	if filename == "" {
		return fmt.Errorf("SqliteDB.Attach: filename is empty")
	}

	if alias == "" {
		return fmt.Errorf("SqliteDB.Attach: alias is empty")
	}

	_, err := db.DB.ExecContext(db.ctx, "ATTACH DATABASE ? AS ?", filename, alias)
	return err
}

func (db *SqliteDB) IsAttached(alias string) bool {
	rows, err := db.Query("PRAGMA database_list")
	if err != nil {
		log.Fatalf("SqliteDB.IsAttached: %s", err)
		return false
	}

	for rows.Next() {
		var id int
		var name, file string

		err := rows.Scan(&id, &name, &file)
		if err != nil {
			log.Fatalf("SqliteDB.IsAttached: %s", err)
			return false
		}

		fmt.Printf("id: %d, name: %s, file: %s, alias: %s\n", id, name, file, alias)
		if name == alias {
			return true
		}
	}

	return false
}

type migrater func(DB) error

func (db *SqliteDB) AddRepo(file, alias string, migrate migrater) error {
	if db.DB == nil {
		return fmt.Errorf("SqliteDB.AddRepo: db.DB is nil")
	}

	if db.IsAttached(alias) {
		return fmt.Errorf("SqliteDB.AddRepo: %w", ErrAliasInUse)
	}

	repo, err := NewSqliteDB(file)
	if err != nil {
		return fmt.Errorf("SqliteDB.AddRepo: %w", err)
	}

	err = repo.Open()
	if err != nil {
		return fmt.Errorf("SqliteDB.AddRepo: failed to open repo db: %w", err)
	}
	defer repo.Close()

	err = migrate(repo)
	if err != nil {
		return fmt.Errorf("SqliteDB.AddRepo: failed to migrate repo: %w", err)
	}

	// Attach the repo database to the main database so we can perform joins.
	// Access tables in the attached repo with "alias.table_name".
	err = db.Attach(file, alias)
	if err != nil {
		return fmt.Errorf("SqliteDB.AddRepo: failed to attach repo: %w", err)
	}

	return nil
}
*/

func (db *SqliteDB) Query(query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(db.ctx, query, args...)
}

func (db *SqliteDB) QueryRow(query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(db.ctx, query, args...)
}

func (db *SqliteDB) Exec(query string, args ...any) error {
	_, err := db.DB.ExecContext(db.ctx, query, args...)
	return err
}

// ############################################################################################# //
// ###################################        Users         #################################### //
// ############################################################################################# //

// UserData represents a user in the database.
type UserData struct {
	ID       int
	Username string // Username to login with.
	Name     string // Name to show in app.
	Password string // Hashed password.
	Groups   string // JSON string of group IDs. Empty should be "[]".
	Created  time.Time
	Updated  time.Time
}

func (db *SqliteDB) UsersMigrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + sqlite_tb_users + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(32) NOT NULL,
		password VARCHAR(32) NOT NULL,
		groups TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("db.UserMigrate: %w", err)
	}

	return nil
}

func (db *SqliteDB) UserIsUnique(username string) error {
	return db.IsUnique(sqlite_tb_users, "username = ?", username)
}

func (db *SqliteDB) UserCreate(username, name, pwHash, groups string) (UserData, error) {
	if username == "" {
		return UserData{}, fmt.Errorf("SqliteDB.UserCreate: %w", ErrInvalidUsername)
	}

	if name == "" {
		return UserData{}, fmt.Errorf("SqliteDB.UserCreate: %w", ErrInvalidName)
	}

	// Password hash should be a 32 byte hex string.
	if err := core.ValidatePasswordHash(pwHash); err != nil {
		return UserData{}, fmt.Errorf("SqliteDB.UserCreate: %w", err)
	}

	if groups == "" {
		groups = "{}"
	}

	if err := db.UserIsUnique(username); err != nil {
		return UserData{}, fmt.Errorf("SqliteDB.UserCreate: %w", err)
	}

	query := `INSERT INTO ` + sqlite_tb_users + ` (username, name, password, groups) VALUES (?, ?, ?, ?)`
	if err := db.Exec(query, username, name, pwHash, groups); err != nil {
		return UserData{}, fmt.Errorf("SqliteDB.UserCreate: %w", err)
	}

	return db.UserGetByUsername(username)
}

func (db *SqliteDB) UserGet(id int) (UserData, error) {
	query := `SELECT * FROM ` + sqlite_tb_users + ` WHERE id = ?`
	row := db.QueryRow(query, id)

	var data UserData
	err := row.Scan(&data.ID, &data.Username, &data.Name, &data.Password, &data.Groups, &data.Created, &data.Updated)
	if err != nil {
		return data, fmt.Errorf("SqliteDB.UserGet: %w", err)
	}

	return data, nil
}

func (db *SqliteDB) UserGetByUsername(username string) (UserData, error) {
	if username == "" {
		return UserData{}, fmt.Errorf("SqliteDB.UserGetByUsername: %w", ErrInvalidUsername)
	}

	query := `SELECT * FROM ` + sqlite_tb_users + ` WHERE username = ?`
	row := db.QueryRow(query, username)

	var data UserData
	err := row.Scan(&data.ID, &data.Username, &data.Name, &data.Password, &data.Groups, &data.Created, &data.Updated)
	if err != nil {
		return data, fmt.Errorf("SqliteDB.UserGetByUsername: %w", err)
	}

	return data, nil
}

func (db *SqliteDB) UserUpdate(user UserData) (UserData, error) {
	if user.ID == 0 {
		return UserData{}, fmt.Errorf("SqliteDB.UserUpdate: %w", ErrInvalidID)
	}

	user.Updated = time.Now()
	query := `UPDATE ` + sqlite_tb_users + ` SET username = ?, name = ?, password = ?, groups = ?, updated_at = ? WHERE id = ?`
	if err := db.Exec(query, user.Username, user.Name, user.Password, user.Groups, user.Updated, user.ID); err != nil {
		return UserData{}, fmt.Errorf("SqliteDB.UserUpdate: %w", err)
	}

	return db.UserGet(user.ID)
}

func (db *SqliteDB) UserDelete(id int) error {
	if id == 0 {
		return fmt.Errorf("db.UserUpdate: %w", ErrInvalidID)
	}

	if _, err := db.UserGet(id); err != nil {
		return fmt.Errorf("db.UserDelete: %w", err)
	}

	query := `DELETE FROM ` + sqlite_tb_users + ` WHERE id = ?`
	if err := db.Exec(query, id); err != nil {
		return fmt.Errorf("db.UserDelete: %w", err)
	}

	return nil
}
