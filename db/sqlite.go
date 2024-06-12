package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type SqliteDB struct {
	Name string
	db   *sql.DB

	ctx    context.Context
	cancel func()
}

func NewSqliteDB(name string) *SqliteDB {
	db := &SqliteDB{Name: name}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *SqliteDB) Open() error {
	if db.Name == "" {
		return fmt.Errorf("no database name provided")
	}

	if db.Name != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.Name), 0o700); err != nil {
			return err
		}
	}

	var err error
	db.db, err = sql.Open("sqlite3", db.Name)
	if err != nil {
		return err
	}

	// Enable WAL mode.
	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreing key checks.
	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	return nil
}

func (db *SqliteDB) Migrate(query string) error {
	_, err := db.db.Exec(query)
	return err
}

func (db *SqliteDB) Close() error {
	if db.db == nil {
		return nil
	}

	db.cancel()
	return db.db.Close()
}
