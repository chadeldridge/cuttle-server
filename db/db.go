package db

import (
	"fmt"
	"os"
)

const (
	DefaultDBFolder = "db"
)

var db_folder string

var (
	ErrRecordExists      = fmt.Errorf("record exists")
	ErrInvalidID         = fmt.Errorf("invalid ID")
	ErrAliasInUse        = fmt.Errorf("db alias in use")
	ErrInvalidUsername   = fmt.Errorf("invalid username")
	ErrInvalidAuthType   = fmt.Errorf("invalid auth type")
	ErrInvalidName       = fmt.Errorf("invalid name")
	ErrInvalidPassphrase = fmt.Errorf("invalid passphrase")
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

type CuttleDB interface {
	Open() error
	CuttleMigrate() error
	Close() error
	// AddRepo(file, alias string, migrate migrater) error
	// Attach(filename, alias string) error
}

type AuthDB interface {
	Open() error
	AuthMigrate() error
	Close() error
	// AddRepo(file, alias string, migrate migrater) error
	// Attach(filename, alias string) error
	// Users
	UserIsUnique(username string) error
	UserCreate(username, name, password, groups string) (UserData, error)
	UserGet(id int) (UserData, error)
	UserGetByUsername(username string) (UserData, error)
	UserUpdate(user UserData) (UserData, error)
	UserDelete(id int) error
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
