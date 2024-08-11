package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testCuttleDBName = "test_cuttle.db"
	testAuthDBName   = "test_auth.db"
)

/*
func testDBSetup(t *testing.T) DB {
	return testSqliteDBSetup(t)
}
*/

func testSqliteCuttleDBSetup(t *testing.T) *SqliteDB {
	require := require.New(t)
	SetDBRoot(testDBRoot)

	db, err := NewSqliteDB(testCuttleDBName)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)

	err = db.Open()
	require.NoError(err, "testDBSetup returned an error: %s", err)
	return db
}

func testSqliteAuthDBSetup(t *testing.T) *SqliteDB {
	require := require.New(t)
	SetDBRoot(testDBRoot)

	db, err := NewSqliteDB(testAuthDBName)
	require.NoError(err, "NewSqliteDB returned an error: %s", err)

	err = db.Open()
	require.NoError(err, "testDBSetup returned an error: %s", err)
	return db
}

/*
func deleteDBDir() {
	if db_folder == "" || db_folder == "/" {
		log.Printf("deleteDBDir: db_folder is dangerous: %s", db_folder)
		return
	}

	err := os.RemoveAll(db_folder)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("deleteDBDir: %s", err)
	}
}
*/

func deleteDB(filename string) {
	if filename == "" || filename == "/" {
		log.Printf("deleteDBDir: db_folder is dangerous: %s", db_folder)
		return
	}

	deleteFile(filename)
	deleteFile(filename + "-shm")
	deleteFile(filename + "-wal")
}

func deleteFile(filename string) {
	err := os.Remove(db_folder + "/" + filename)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
		log.Fatalf("deleteDB: %s", err)
	}
}

func TestSqliteNewSqliteDB(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		db, err := NewSqliteDB(testCuttleDBName)
		require.NoError(err, "NewSqliteDB returned an error: %s", err)
		require.NotNil(db)
		require.Equal(testCuttleDBName, db.Name)
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

func testDBMigrater(db *SqliteDB) error {
	query := `
	CREATE TABLE IF NOT EXISTS test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(32) NOT NULL
	)`

	if err := db.Exec(query); err != nil {
		return fmt.Errorf("testDBMigrater: %w", err)
	}

	return nil
}

/*
func testBadMigrater(db *SqliteDB) error {
	if db == nil {
		return fmt.Errorf("testBadMigrate: db is nil")
	}

	return fmt.Errorf("testBadMigrater: failed to migrate")
}

func testGetAll(db *SqliteDB) {
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		log.Fatalf("Query: %s", err)
	}

	found := false
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("Scan: %s", err)
		}

		found = true
		log.Printf("id: %d, name: %s", id, name)
	}

	if !found {
		log.Println("    ----    No rows found    ----")
	}
}
*/

func TestSqliteOpen(t *testing.T) {
	require := require.New(t)
	SetDBRoot(testDBRoot)
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
		db, err := NewSqliteDB(testCuttleDBName)
		require.NoError(err, "NewSqliteDB returned an error: %s", err)

		fmt.Printf("db: %s\n", db_folder+"/"+db.Name)
		err = db.Open()
		require.NoError(err, "Open returned an error: %s", err)
		require.FileExists(db_folder+"/"+db.Name, "Open did not create the database file")
		db.DB.Close()
		deleteDB(testCuttleDBName)
	})
}

func TestSqliteClose(t *testing.T) {
	require := require.New(t)
	db := testSqliteCuttleDBSetup(t)
	defer deleteDB(testCuttleDBName)

	t.Run("valid", func(t *testing.T) {
		err := db.Close()
		require.NoError(err, "Close returned an error: %s", err)
	})

	t.Run("nil db.DB", func(t *testing.T) {
		db.DB = nil
		err := db.Close()
		require.NoError(err, "Close returned an error: %s", err)
	})
}

func TestSqliteIsUnique(t *testing.T) {
	require := require.New(t)
	db := testSqliteCuttleDBSetup(t)
	// Setup the test table.
	testDBMigrater(db)
	defer db.Close()
	defer deleteDB(testCuttleDBName)

	t.Run("unique", func(t *testing.T) {
		err := db.IsUnique("test", "name = ?", "testRecord_1")
		require.NoError(err, "IsUnique returned an error: %s", err)
	})

	// Insert a row.
	err := db.Exec("INSERT INTO test (name) VALUES ('testRecord_1')")
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("not unique", func(t *testing.T) {
		err := db.IsUnique("test", "name = ?", "testRecord_1")
		require.Error(err, "IsUnique did not return an error")
		require.ErrorIs(err, ErrRecordExists, "IsUnique did not return the expected error")
	})
}

// ############################################################################################# //
// ###################################        Users         #################################### //
// ############################################################################################# //

var (
	testUser1 = UserData{
		Username: "user1",
		Name:     "Bob",
		Password: "102650912390a29378e092378b29834f",
		Groups:   "{}",
	}

	testUser2 = UserData{
		Username: "user2",
		Name:     "Jan",
		Password: "102650912390a29378e092378b29834f",
		Groups:   `[1, 45]`,
	}
)

func TestSqliteUserCreate(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("empty username", func(t *testing.T) {
		_, err := db.UserCreate("", testUser1.Name, testUser1.Password, testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.ErrorIs(err, ErrInvalidUsername, "UserCreate returned an unexpected error")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := db.UserCreate(testUser1.Username, "", testUser1.Password, testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.ErrorIs(err, ErrInvalidName, "UserCreate returned an unexpected error")
	})

	t.Run("empty pwHash", func(t *testing.T) {
		_, err := db.UserCreate(testUser1.Username, testUser1.Name, "", testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.Equal(
			"SqliteDB.UserCreate: core.ValidatePasswordHash: hash was empty",
			err.Error(),
			"UserCreate returned an unexpected error",
		)
	})

	t.Run("short pwHash", func(t *testing.T) {
		_, err := db.UserCreate(testUser1.Username, testUser1.Name, "102650912390a29", testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.Equal(
			"SqliteDB.UserCreate: core.ValidatePasswordHash: incorrect hash length: 15",
			err.Error(),
			"UserCreate returned an unexpected error",
		)
	})

	t.Run("non-hex pwHash", func(t *testing.T) {
		_, err := db.UserCreate(
			testUser1.Username,
			testUser1.Name,
			"102650912390a29378-092378b29834f",
			testUser1.Groups,
		)
		require.Error(err, "UserCreate did not return an error")
		require.Equal(
			"SqliteDB.UserCreate: core.ValidatePasswordHash: hash is not a hex string",
			err.Error(),
			"UserCreate returned an unexpected error",
		)
	})

	t.Run("empty groups", func(t *testing.T) {
		_, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, "")
		require.NoError(err, "UserCreate returned an error: %s", err)
	})

	t.Run("valid", func(t *testing.T) {
		got, err := db.UserCreate(testUser2.Username, testUser2.Name, testUser2.Password, testUser2.Groups)
		require.NoError(err, "UserCreate returned an error: %s", err)
		require.Equal(testUser2.Username, got.Username)
		require.Equal(testUser2.Name, got.Name)
		require.Equal(testUser2.Password, got.Password)
		require.Equal(testUser2.Groups, got.Groups)
		require.NotZero(got.Created)
		require.NotZero(got.Updated)
	})
}

func TestSqliteUserIsUnique(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("unique", func(t *testing.T) {
		err := db.UserIsUnique(testUser1.Username)
		require.NoError(err, "IsUnique returned an error: %s", err)
	})

	// Insert a row.
	_, err = db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("not unique", func(t *testing.T) {
		err := db.UserIsUnique(testUser1.Username)
		require.Error(err, "IsUnique did not return an error")
		require.ErrorIs(err, ErrRecordExists, "IsUnique did not return the expected error")
	})
}

func TestSqliteUserGet(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	want, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("valid", func(t *testing.T) {
		data, err := db.UserGet(1)
		require.NoError(err, "UserGet returned an error: %s", err)
		require.Equal(want.Username, data.Username)
		require.Equal(want.Name, data.Name)
		require.Equal(want.Password, data.Password)
		require.Equal(want.Groups, data.Groups)
		require.Equal(want.Created, data.Created)
		require.Equal(want.Updated, data.Updated)
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := db.UserGet(9999)
		require.Error(err, "UserGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGet did not return the expected error")
		require.Equal(UserData{}, data)
	})
}

func TestSqliteUserGetByUsername(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	_, err = db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("empty username", func(t *testing.T) {
		data, err := db.UserGetByUsername("")
		require.Error(err, "UserGetByUsername did not return an error")
		require.ErrorIs(err, ErrInvalidUsername, "UserGetByUsername did not return the expected error")
		require.Equal(UserData{}, data)
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := db.UserGetByUsername("not_a_user")
		require.Error(err, "UserGetByUsername did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGetByUsername did not return the expected error")
		require.Equal(UserData{}, data)
	})

	t.Run("valid", func(t *testing.T) {
		data, err := db.UserGetByUsername(testUser1.Username)
		require.NoError(err, "UserGet returned an error: %s", err)
		require.Equal(testUser1.Username, data.Username)
		require.Equal(testUser1.Name, data.Name)
		require.Equal(testUser1.Password, data.Password)
		require.Equal(testUser1.Groups, data.Groups)
	})
}

func TestSqliteUserUpdate(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("invalid ID", func(t *testing.T) {
		_, err := db.UserUpdate(UserData{ID: 0})
		require.Error(err, "UserUpdate did not return an error")
		require.ErrorIs(err, ErrInvalidID, "UserUpdate did not return the expected error")
	})

	createdAt := data.Created
	updatedAt := data.Updated

	t.Run("valid", func(t *testing.T) {
		data.Username = testUser2.Username
		data.Name = testUser2.Name
		data.Password = testUser2.Password
		data.Groups = testUser2.Groups

		updated, err := db.UserUpdate(data)
		require.NoError(err, "UserUpdate returned an error: %s", err)
		require.Equal(data.Username, updated.Username)
		require.Equal(data.Name, updated.Name)
		require.Equal(data.Password, updated.Password)
		require.Equal(data.Groups, updated.Groups)
		require.Equal(createdAt, updated.Created)
		require.Greater(updated.Updated, updatedAt)
	})

	t.Run("invalid id", func(t *testing.T) {
		data.ID = 9999
		_, err := db.UserUpdate(data)
		require.Error(err, "UserUpdate did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserUpdate did not return the expected error")
	})
}

func TestSqliteUserDelete(t *testing.T) {
	require := require.New(t)
	db := testSqliteAuthDBSetup(t)
	defer db.Close()
	defer deleteDB(testAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Password, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("invalid ID", func(t *testing.T) {
		err := db.UserDelete(0)
		require.Error(err, "UserDelete did not return an error")
		require.ErrorIs(err, ErrInvalidID, "UserDelete did not return the expected error")
	})

	t.Run("valid", func(t *testing.T) {
		err := db.UserDelete(data.ID)
		require.NoError(err, "UserDelete returned an error: %s", err)

		_, err = db.UserGet(data.ID)
		require.Error(err, "UserGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGet did not return the expected error")
	})

	t.Run("invalid", func(t *testing.T) {
		err := db.UserDelete(9999)
		require.Error(err, "UserDelete did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserDelete did not return the expected error")
	})
}

/*
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
*/
