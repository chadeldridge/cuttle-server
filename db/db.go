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
)

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	db_folder = currentDir + "/" + DefaultDBFolder
}

type DB interface {
	Open() error
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
func SetDBRoot(rootDir string) {
	db_folder = rootDir
}
