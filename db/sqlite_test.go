package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func testDBSetup(t *testing.T) DB {
	return testSqliteDBSetup(t)
}

func testSqliteDBSetup(t *testing.T) *SqliteDB {
	require := require.New(t)
	db_folder = testDBRoot

	db, err := NewSqliteDB("cuttle.db")
	require.NoError(err, "NewSqliteDB returned an error: %s", err)

	err = db.Open()
	require.NoError(err, "testDBSetup returned an error: %s", err)
	return db
}

func deleteDBDir() {
	err := os.RemoveAll(db_folder)
	if err != nil {
		log.Fatalf("deleteDBDir: %s", err)
	}
}

func deleteDB(filename string) {
	if _, err := os.Stat(db_folder + "/" + filename); os.IsNotExist(err) {
		return
	}

	err := os.Remove(db_folder + "/" + filename)
	if err != nil {
		log.Println(err)
		log.Fatalf("deleteDB: %s", err)
	}
}

func TestSqliteNewSqliteDB(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		db, err := NewSqliteDB("cuttle.db")
		require.NoError(err, "NewSqliteDB returned an error: %s", err)
		require.NotNil(db)
		require.Equal("cuttle.db", db.Name)
		require.NotNil(db.ctx)
		require.NotNil(db.cancel)
	})

	t.Run("empty filename", func(t *testing.T) {
		db, err := NewSqliteDB("")
		require.Error(err, "NewSqliteDB did not return an error")
		require.Equal("db.NewSqliteDB: filename is empty", err.Error(), "NewSqliteDB returned an unexpected error")
		require.Nil(db)
	})
}

func TestSqliteOpen(t *testing.T) {
	require := require.New(t)
	db_folder = testDBRoot
	var db *SqliteDB

	t.Run("empty db.Name", func(t *testing.T) {
		db = &SqliteDB{Name: ""}
		db.ctx, db.cancel = context.WithCancel(context.Background())

		err := db.Open()
		require.Error(err, "Open did not return an error")
		require.Equal("no database name provided", err.Error(), "Open returned an unexpected error")
	})

	t.Run("invalid location", func(t *testing.T) {
		db = &SqliteDB{Name: "not_here/cuttle.db"}
		db.ctx, db.cancel = context.WithCancel(context.Background())

		err := db.Open()
		require.Error(err, "Open did not return an error")
		require.Equal(
			"SqliteDB.Open: failed to ping db: unable to open database file: no such file or directory",
			err.Error(),
			"Open returned an unexpected error: %s",
			err,
		)
	})

	t.Run("valid", func(t *testing.T) {
		db, err := NewSqliteDB("cuttle.db")
		require.NoError(err, "NewSqliteDB returned an error: %s", err)

		err = db.Open()
		require.NoError(err, "Open returned an error: %s", err)
		require.FileExists(db_folder+"/"+db.Name, "Open did not create the database file")
		db.Close()
	})
}

func TestSqliteAttach(t *testing.T) {
	require := require.New(t)
	db_folder = testDBRoot
	db := testSqliteDBSetup(t)
	defer db.Close()

	db_file := "tmp.db"
	db_alias := "t"

	tmpDB, err := NewSqliteDB(db_file)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)
	err = tmpDB.Open()
	require.NoError(err, "Open returned an error: %s", err)
	tmpDB.Close()

	t.Run("empty filename", func(t *testing.T) {
		err := db.Attach("", db_alias)
		require.Error(err, "Attach did not return an error")
		require.Equal("SqliteDB.Attach: filename is empty", err.Error(), "Attach returned an unexpected error")
	})

	t.Run("empty alias", func(t *testing.T) {
		err := db.Attach(db_file, "")
		require.Error(err, "Attach did not return an error")
		require.Equal("SqliteDB.Attach: alias is empty", err.Error(), "Attach returned an unexpected error")
	})

	t.Run("valid", func(t *testing.T) {
		err := db.Attach(db_file, db_alias)
		require.NoError(err, "Attach returned an error: %s", err)

		rows, err := db.Query("PRAGMA database_list")
		require.NoError(err, "Query returned an error: %s", err)

		var found bool
		for rows.Next() {
			var id int
			var name, file string

			err := rows.Scan(&id, &name, &file)
			require.NoError(err, "Scan returned an error: %s", err)

			if name == db_alias {
				found = true
			}
		}

		require.True(found, "Attach did not attach the database")
	})
}

func TestSqliteIsAttached(t *testing.T) {
	require := require.New(t)
	db_folder = testDBRoot
	db := testSqliteDBSetup(t)
	defer db.Close()

	db_file := "tmp.db"
	db_alias := "t"

	tmpDB, err := NewSqliteDB(db_file)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)
	err = tmpDB.Open()
	require.NoError(err, "Open returned an error: %s", err)
	tmpDB.Close()

	t.Run("false", func(t *testing.T) {
		require.False(db.IsAttached(db_alias), "IsAttached returned true")
	})

	err = db.Attach(db_file, db_alias)
	require.NoError(err, "Attach returned an error: %s", err)

	t.Run("true", func(t *testing.T) {
		require.True(db.IsAttached(db_alias), "IsAttached returned false")
	})

	deleteDB(db_file)
}

func testDBMigrater(db DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("testDBMigrater: %w", err)
	}

	return nil
}

func testBadMigrater(db DB) error {
	if db == nil {
		return fmt.Errorf("testBadMigrate: db is nil")
	}

	return fmt.Errorf("testBadMigrater: failed to migrate")
}

func TestSqliteAddRepo(t *testing.T) {
	require := require.New(t)
	db_folder = testDBRoot
	db := testSqliteDBSetup(t)
	defer db.Close()

	db_file := "tmp.db"
	db_alias := "t"

	t.Run("nil db", func(t *testing.T) {
		badDB := &SqliteDB{}
		err := badDB.AddRepo(db_file, db_alias, testDBMigrater)
		require.Error(err, "AddRepo did not return an error")
		require.Equal("SqliteDB.AddRepo: db.DB is nil", err.Error(), "AddRepo returned an unexpected error")
	})

	t.Run("bad migrater", func(t *testing.T) {
		err := db.AddRepo(db_file, db_alias, testBadMigrater)
		require.Error(err, "AddRepo did not return an error")
		require.Equal(
			"SqliteDB.AddRepo: failed to migrate repo: testBadMigrater: failed to migrate",
			err.Error(),
			"AddRepo returned an unexpected error",
		)
	})

	err := db.AddRepo(db_file, db_alias, testDBMigrater)
	t.Run("valid", func(t *testing.T) {
		require.NoError(err, "AddRepo returned an error: %s", err)
		require.FileExists(db_folder+"/"+db_file, "AddRepo did not create the database file")
		require.True(db.IsAttached(db_alias), "AddRepo did not attach the database")
	})

	row := db.QueryRow("SELECT COUNT(*) FROM t.test")
	var count int
	err = row.Scan(&count)
	require.NoError(err, "Scan returned an error: %s", err)
	log.Printf("count: %d", count)
	require.True(db.IsAttached(db_alias), "AddRepo did not attach the database")
	t.Run("duplicate alias", func(t *testing.T) {
		err := db.AddRepo(db_file, db_alias, testDBMigrater)
		require.Error(err, "AddRepo did not return an error")
		require.Equal(
			"SqliteDB.AddRepo: failed to attach repo: SqliteDB.Attach: alias is already in use",
			err.Error(),
			"AddRepo returned an unexpected error",
		)
	})

	deleteDB(db_file)
}
