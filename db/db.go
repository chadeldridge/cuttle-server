package db

import (
	"database/sql"
	"os"
)

const (
	DefaultDBFolder = "db"
)

var db_folder string

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
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) error
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
