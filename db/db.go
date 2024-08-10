package db

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	DefaultDBFolder = "db"
)

var db_folder string

var (
	ErrDuplicateEntry = fmt.Errorf("UNIQUE constraint failed")
	ErrRecordNotFound = fmt.Errorf("record not found")
	ErrInvalidID      = fmt.Errorf("invalid ID")
	ErrAliasInUse     = fmt.Errorf("db alias in use")
)

func init() {
	db_folder = GenDBFolder()
}

func GenDBFolder() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return currentDir + "/" + DefaultDBFolder
}

type DB interface {
	Open() error
	AddRepo(file, alias string, migrate migrater) error
	Attach(filename, alias string) error
	IsUnique(query string, args ...any) bool
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) error
	Close() error
}

// SetDBRoot sets the root directory for the database. If this is not set, the default behavior is to
// create a directory called "db" in the current working directory.
// Example:
//
//	db.InitDB("/path/to/db")
//
// Expected Behavior:
//
//	db.InitDB("/tmp/db")
//	db.NewSqliteDB("mydb.db")
//	db.Open()
//
//	`ls /tmp/db/`
//	mydb.db
func SetDBRoot(rootDir string) error {
	if rootDir == "" {
		return fmt.Errorf("db.SetDBRoot: rootDir is empty")
	}

	if _, err := os.Stat(rootDir); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("db.SetDBRoot: failed to stat rootDir: %w", err)
		}

		if err := os.MkdirAll(rootDir, 0o755); err != nil {
			return fmt.Errorf("db.SetDBRoot: failed to create rootDir: %w", err)
		}
	}

	db_folder = rootDir
	return nil
}
