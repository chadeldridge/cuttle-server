package db

import (
	"log"
	"testing"

	"github.com/chadeldridge/cuttle-server/test_helpers"
	"github.com/stretchr/testify/require"
)

var (
	TestDBRoot       = "/tmp/cuttle/db"
	TestCuttleDBName = "test_cuttle.db"
	TestAuthDBName   = "test_auth.db"
)

func DeleteDB(filename string) {
	if filename == "" || filename == "/" {
		log.Printf("deleteDBDir: testDBRoot is dangerous: '%s'", TestDBRoot)
		return
	}

	test_helpers.DeleteFile(TestDBRoot + "/" + filename)
	test_helpers.DeleteFile(TestDBRoot + "/" + filename + "-shm")
	test_helpers.DeleteFile(TestDBRoot + "/" + filename + "-wal")
}

// TestSqliteCuttleDBSetup creates a new SqliteDB instance for testing. You will still need to run CuttleMigrate.
func TestSqliteCuttleDBSetup(t *testing.T) *SqliteDB {
	require := require.New(t)
	SetDBRoot(TestDBRoot)

	db, err := NewSqliteDB(TestCuttleDBName)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)

	err = db.Open()
	require.NoError(err, "testDBSetup returned an error: %s", err)
	return db
}

// TestSqliteAuthDBSetup creates a new SqliteDB instance for testing. You will still need to run AuthMigrate.
func TestSqliteAuthDBSetup(t *testing.T) *SqliteDB {
	require := require.New(t)
	SetDBRoot(TestDBRoot)

	db, err := NewSqliteDB(TestAuthDBName)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)

	err = db.Open()
	require.NoError(err, "testDBSetup returned an error: %s", err)
	return db
}
