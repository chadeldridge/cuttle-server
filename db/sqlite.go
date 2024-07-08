package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var ErrDuplicateEntry = fmt.Errorf("duplicate entry")

type SqliteDB struct {
	Name string // DB file name.
	DB   *sql.DB

	ctx    context.Context
	cancel func()
}

// NewSqliteDB creates a new Sqlite3 DB instance.
func NewSqliteDB(filename string) *SqliteDB {
	db := &SqliteDB{Name: filename}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

// Open opens the database. It attempts to create the path to the database file if it does not
// exist, opens the database file, and enables WAL mode and foreign key checks.
func (db *SqliteDB) Open() error {
	if db.Name == "" {
		return fmt.Errorf("no database name provided")
	}

	// if db.Name != ":memory:" {
	if err := os.MkdirAll(filepath.Dir(db_folder+"/"+db.Name), 0o700); err != nil {
		return err
	}
	//}

	var err error
	db.DB, err = sql.Open("sqlite3", db_folder+"/"+db.Name)
	if err != nil {
		log.Fatalf("SqliteDB.Open: %s", err)
		return err
	}

	// Enable WAL mode.
	if _, err := db.DB.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreign key checks.
	if _, err := db.DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	// TODO: Implement zstd compression. https://phiresky.github.io/blog/2022/sqlite-zstd/

	return nil
}

// Attach the filename database tot he current database with the given alias.
func (db *SqliteDB) Attach(filename, alias string) error {
	_, err := db.DB.ExecContext(db.ctx, "ATTACH DATABASE ? AS ?", filename, alias)
	return err
}

func (db *SqliteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(db.ctx, query, args...)
}

func (db *SqliteDB) Exec(query string, args ...interface{}) error {
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
