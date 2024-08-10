package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// libray has to be imported to register the driver.
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	Name string // DB file name.
	DB   *sql.DB

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

func (db *SqliteDB) IsUnique(query string, args ...any) bool {
	row := db.QueryRow(query, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatalf("SqliteDB.AuthMethods.IsUnique: %s", err)
	}

	return count == 0
}

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

func (db *SqliteDB) Close() error {
	if db.DB == nil {
		return nil
	}

	db.cancel()
	return db.DB.Close()
}
