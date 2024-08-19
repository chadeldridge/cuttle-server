package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/stretchr/testify/require"
)

/*
func testDBSetup(t *testing.T) DB {
	return testSqliteDBSetup(t)
}
*/

func TestSqliteDBNewSqliteDB(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		db, err := NewSqliteDB(TestCuttleDBName)
		require.NoError(err, "NewSqliteDB returned an error: %s", err)
		require.NotNil(db)
		require.Equal(TestCuttleDBName, db.Name)
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

	if _, err := db.Exec(query); err != nil {
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

func TestSqliteDBOpen(t *testing.T) {
	require := require.New(t)
	SetDBRoot(TestDBRoot)
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
		db, err := NewSqliteDB(TestCuttleDBName)
		require.NoError(err, "NewSqliteDB returned an error: %s", err)

		fmt.Printf("db: %s\n", db_folder+"/"+db.Name)
		err = db.Open()
		require.NoError(err, "Open returned an error: %s", err)
		require.FileExists(db_folder+"/"+db.Name, "Open did not create the database file")
		db.DB.Close()
		DeleteDB(TestCuttleDBName)
	})
}

func TestSqliteDBClose(t *testing.T) {
	require := require.New(t)
	db := TestSqliteCuttleDBSetup(t)
	defer DeleteDB(TestCuttleDBName)

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

func TestSqliteDBIsUnique(t *testing.T) {
	require := require.New(t)
	db := TestSqliteCuttleDBSetup(t)
	// Setup the test table.
	testDBMigrater(db)
	defer db.Close()
	defer DeleteDB(TestCuttleDBName)

	t.Run("unique", func(t *testing.T) {
		err := db.IsUnique("test", "name = ?", "testRecord_1")
		require.NoError(err, "IsUnique returned an error: %s", err)
	})

	// Insert a row.
	_, err := db.Exec("INSERT INTO test (name) VALUES ('testRecord_1')")
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
		Hash:     "102650912390a29378e092378b29834f",
		Groups:   "[]",
	}

	testUser2 = UserData{
		Username: "user2",
		Name:     "Jan",
		Hash:     "102650912390a29378e092378b29834f",
		Groups:   `[1, 45]`,
	}
)

func TestSqliteDBUserMigrate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	createQuery := `CREATE TABLE ` + sqlite_tb_users + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(32) NOT NULL,
		password VARCHAR(32) NOT NULL,
		groups TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	indexQuery := `CREATE UNIQUE INDEX idx_users_username ON ` + sqlite_tb_users + ` (username)`

	t.Run("table", func(t *testing.T) {
		err := db.AuthMigrate()
		require.NoError(err, "AuthMigrate returned an error: %s", err)
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = '" + sqlite_tb_users + "'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(createQuery, schema)
	})

	t.Run("index", func(t *testing.T) {
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = 'idx_users_username'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(indexQuery, schema)
	})
}

func TestSqliteDBUserCreate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("empty username", func(t *testing.T) {
		_, err := db.UserCreate("", testUser1.Name, testUser1.Hash, testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "UserCreate returned an unexpected error")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := db.UserCreate(testUser1.Username, "", testUser1.Hash, testUser1.Groups)
		require.Error(err, "UserCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "UserCreate returned an unexpected error")
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
		_, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, "")
		require.NoError(err, "UserCreate returned an error: %s", err)
	})

	t.Run("valid", func(t *testing.T) {
		got, err := db.UserCreate(testUser2.Username, testUser2.Name, testUser2.Hash, testUser2.Groups)
		require.NoError(err, "UserCreate returned an error: %s", err)
		require.Equal(testUser2.Username, got.Username)
		require.Equal(testUser2.Name, got.Name)
		require.Equal(testUser2.Hash, got.Hash)
		require.Equal(testUser2.Groups, got.Groups)
		require.NotZero(got.Created)
		require.NotZero(got.Updated)
	})

	t.Run("duplicate", func(t *testing.T) {
		_, err := db.UserCreate(testUser2.Username, testUser2.Name, testUser2.Hash, testUser2.Groups)
		require.Error(err, "UserCreate returned an error: %s", err)
		require.ErrorIs(err, ErrUserExists, "UserCreate did not return the expected error")
	})
}

func TestSqliteDBUserIsUnique(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("unique", func(t *testing.T) {
		err := db.UserIsUnique(testUser1.Username)
		require.NoError(err, "IsUnique returned an error: %s", err)
	})

	// Insert a row.
	_, err = db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("not unique", func(t *testing.T) {
		err := db.UserIsUnique(testUser1.Username)
		require.Error(err, "IsUnique did not return an error")
		require.ErrorIs(err, ErrUserExists, "IsUnique did not return the expected error")
	})
}

func TestSqliteDBUserGet(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	want, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("valid", func(t *testing.T) {
		data, err := db.UserGet(1)
		require.NoError(err, "UserGet returned an error: %s", err)
		require.Equal(want.Username, data.Username)
		require.Equal(want.Name, data.Name)
		require.Equal(want.Hash, data.Hash)
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

func TestSqliteDBUserGetByUsername(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	_, err = db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, testUser1.Groups)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("empty username", func(t *testing.T) {
		data, err := db.UserGetByUsername("")
		require.Error(err, "UserGetByUsername did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "UserGetByUsername did not return the expected error")
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
		require.Equal(testUser1.Hash, data.Hash)
		require.Equal(testUser1.Groups, data.Groups)
	})
}

func TestSqliteDBUserUpdate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, testUser1.Groups)
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
		data.Hash = testUser2.Hash
		data.Groups = testUser2.Groups

		updated, err := db.UserUpdate(data)
		require.NoError(err, "UserUpdate returned an error: %s", err)
		require.Equal(data.Username, updated.Username)
		require.Equal(data.Name, updated.Name)
		require.Equal(data.Hash, updated.Hash)
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

func TestSqliteDBUserDelete(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserCreate(testUser1.Username, testUser1.Name, testUser1.Hash, testUser1.Groups)
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

// ############################################################################################## //
// ##################################        User Groups        ################################# //
// ############################################################################################## //

var (
	testUserGroup1 = UserGroupData{
		Name:     "Test UserGroup 1",
		Members:  "[]",
		Profiles: "{}",
	}

	testUserGroup2 = UserGroupData{
		Name:     "Test UserGroup 2",
		Members:  "[1,5,28,349]",
		Profiles: `{"Web Servers": {"POST": false, "GET": true, "PUT": false, "DELETE": false}, "DB Servers": {"POST": false, "GET": true, "PUT": true, "DELETE": false}}`,
	}
)

func TestSqliteDBUserGroupMigrate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	createQuery := `CREATE TABLE ` + sqlite_tb_user_groups + ` (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255) NOT NULL UNIQUE,
	members TEXT NOT NULL,
	profiles TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
	indexQuery := `CREATE INDEX idx_user_groups_name ON ` + sqlite_tb_user_groups + ` (name)`

	t.Run("table", func(t *testing.T) {
		err := UserGroupsMigrate(db)
		require.NoError(err, "UserGroupMigrate returned an error: %s", err)
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = '" + sqlite_tb_user_groups + "'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(createQuery, schema)
	})

	t.Run("index", func(t *testing.T) {
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = 'idx_user_groups_name'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(indexQuery, schema)
	})
}

func TestSqliteDBUserGroupCreate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("empty name", func(t *testing.T) {
		_, err := db.UserGroupCreate("", testUserGroup1.Members, testUserGroup1.Profiles)
		require.Error(err, "UserGroupCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "UserGroupCreate did not return the expected error")
	})

	t.Run("empty members", func(t *testing.T) {
		got, err := db.UserGroupCreate(testUserGroup1.Name, "", testUserGroup1.Profiles)
		require.NoError(err, "UserGroupCreate returned an error: %s", err)
		require.Equal(testUserGroup1.Name, got.Name)
		require.Equal(testUserGroup1.Members, got.Members)
		require.Equal(testUserGroup1.Profiles, got.Profiles)
		require.NotZero(got.Created)
		require.NotZero(got.Updated)

		err = db.UserGroupDelete(got.ID)
		require.NoError(err, "UserGroupDelete returned an error: %s", err)
	})

	t.Run("empty profiles", func(t *testing.T) {
		got, err := db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, "")
		require.NoError(err, "UserGroupCreate returned an error: %s", err)
		require.Equal(testUserGroup1.Name, got.Name)
		require.Equal(testUserGroup1.Members, got.Members)
		require.Equal(testUserGroup1.Profiles, got.Profiles)
		require.NotZero(got.Created)
		require.NotZero(got.Updated)

		err = db.UserGroupDelete(got.ID)
		require.NoError(err, "UserGroupDelete returned an error: %s", err)
	})

	t.Run("all values", func(t *testing.T) {
		got, err := db.UserGroupCreate(testUserGroup2.Name, testUserGroup2.Members, testUserGroup2.Profiles)
		require.NoError(err, "UserGroupCreate returned an error: %s", err)
		require.Equal(testUserGroup2.Name, got.Name)
		require.Equal(testUserGroup2.Members, got.Members)
		require.Equal(testUserGroup2.Profiles, got.Profiles)
		require.NotZero(got.Created)
		require.NotZero(got.Updated)

		err = db.UserGroupDelete(got.ID)
		require.NoError(err, "UserGroupDelete returned an error: %s", err)
	})

	t.Run("duplicate", func(t *testing.T) {
		got, err := db.UserGroupCreate(testUserGroup2.Name, testUserGroup2.Members, testUserGroup2.Profiles)
		require.NoError(err, "UserGroupCreate returned an error: %s", err)

		_, err = db.UserGroupCreate(testUserGroup2.Name, testUserGroup2.Members, testUserGroup2.Profiles)
		require.Error(err, "UserGroupCreate did not return an error")
		require.ErrorIs(err, ErrUserGroupExists, "UserGroupCreate did not return the expected error")

		err = db.UserGroupDelete(got.ID)
		require.NoError(err, "UserGroupDelete returned an error: %s", err)
	})
}

func TestSqliteDBUserGroupIsUnique(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("unique", func(t *testing.T) {
		err := db.UserGroupIsUnique(testUserGroup1.Name)
		require.NoError(err, "IsUnique returned an error: %s", err)
	})

	// Insert a row.
	_, err = db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("not unique", func(t *testing.T) {
		err := db.UserGroupIsUnique(testUserGroup1.Name)
		require.Error(err, "IsUnique did not return an error")
		require.ErrorIs(err, ErrUserGroupExists, "IsUnique did not return the expected error")
	})
}

func TestSqliteDBUserGroupGet(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	want, err := db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("valid", func(t *testing.T) {
		data, err := db.UserGroupGet(want.ID)
		require.NoError(err, "UserGroupGet returned an error: %s", err)
		require.Equal(want.Name, data.Name)
		require.Equal(want.Members, data.Members)
		require.Equal(want.Profiles, data.Profiles)
		require.Equal(want.Created, data.Created)
		require.Equal(want.Updated, data.Updated)
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := db.UserGroupGet(9999)
		require.Error(err, "UserGroupGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGroupGet did not return the expected error")
		require.Equal(UserGroupData{}, data)
	})
}

func TestSqliteDBUserGroupGetByName(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	_, err = db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("empty name", func(t *testing.T) {
		data, err := db.UserGroupGetByName("")
		require.Error(err, "UserGroupGetByName did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "UserGroupGetByName did not return the expected error")
		require.Equal(UserGroupData{}, data)
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := db.UserGroupGetByName("not_a_user_group")
		require.Error(err, "UserGroupGetByName did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGroupGetByName did not return the expected error")
		require.Equal(UserGroupData{}, data)
	})

	t.Run("valid", func(t *testing.T) {
		data, err := db.UserGroupGetByName(testUserGroup1.Name)
		require.NoError(err, "UserGroupGet returned an error: %s", err)
		require.Equal(testUserGroup1.Name, data.Name)
		require.Equal(testUserGroup1.Members, data.Members)
		require.Equal(testUserGroup1.Profiles, data.Profiles)
	})
}

func TestSqliteDBUserGroupGetGroups(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("empty gids", func(t *testing.T) {
		groups, err := db.UserGroupGetGroups([]int64{})
		require.NoError(err, "UserGroupGetGroups returned an error: %s", err)
		require.Empty(groups)
	})

	t.Run("not found", func(t *testing.T) {
		groups, err := db.UserGroupGetGroups([]int64{1})
		require.NoError(err, "UserGroupGetGroups returned an error: %s", err)
		require.Empty(groups)
	})

	// Insert a row.
	want1, err := db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("one", func(t *testing.T) {
		data, err := db.UserGroupGetGroups([]int64{want1.ID})
		require.NoError(err, "UserGroupGetGroups returned an error: %s", err)
		require.Len(data, 1)
		require.Equal(want1.Name, data[0].Name)
		require.Equal(want1.Members, data[0].Members)
		require.Equal(want1.Profiles, data[0].Profiles)
	})

	// Insert a second row.
	want2, err := db.UserGroupCreate(testUserGroup2.Name, testUserGroup2.Members, testUserGroup2.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("two", func(t *testing.T) {
		data, err := db.UserGroupGetGroups([]int64{want1.ID, want2.ID})
		require.NoError(err, "UserGroupGetGroups returned an error: %s", err)
		require.Len(data, 2)
		require.Equal(want1.Name, data[0].Name)
		require.Equal(want1.Members, data[0].Members)
		require.Equal(want1.Profiles, data[0].Profiles)
		require.Equal(want2.Name, data[1].Name)
		require.Equal(want2.Members, data[1].Members)
		require.Equal(want2.Profiles, data[1].Profiles)
	})
}

func TestSqliteDBUserGroupUpdate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("invalid id", func(t *testing.T) {
		_, err := db.UserGroupUpdate(UserGroupData{ID: 0})
		require.Error(err, "UserGroupUpdate did not return an error")
		require.ErrorIs(err, ErrInvalidID, "UserGroupUpdate did not return the expected error")
	})

	createdAt := data.Created
	updatedAt := data.Updated

	t.Run("valid", func(t *testing.T) {
		data.Name = testUserGroup2.Name
		data.Members = testUserGroup2.Members
		data.Profiles = testUserGroup2.Profiles

		updated, err := db.UserGroupUpdate(data)
		require.NoError(err, "UserGroupUpdate returned an error: %s", err)
		require.Equal(data.Name, updated.Name)
		require.Equal(data.Members, updated.Members)
		require.Equal(data.Profiles, updated.Profiles)
		require.Equal(createdAt, updated.Created)
		require.Greater(updated.Updated, updatedAt)
	})

	t.Run("not found", func(t *testing.T) {
		data.ID = 9999
		_, err := db.UserGroupUpdate(data)
		require.Error(err, "UserGroupUpdate did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGroupUpdate did not return the expected error")
	})
}

func TestSqliteDBUserGroupDelete(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	// Setup the test tables.
	err := db.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	// Insert a row.
	data, err := db.UserGroupCreate(testUserGroup1.Name, testUserGroup1.Members, testUserGroup1.Profiles)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("invalid id", func(t *testing.T) {
		err := db.UserGroupDelete(0)
		require.Error(err, "UserGroupDelete did not return an error")
		require.ErrorIs(err, ErrInvalidID, "UserGroupDelete did not return the expected error")
	})

	t.Run("valid", func(t *testing.T) {
		err := db.UserGroupDelete(data.ID)
		require.NoError(err, "UserGroupDelete returned an error: %s", err)

		_, err = db.UserGroupGet(data.ID)
		require.Error(err, "UserGroupGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGroupGet did not return the expected error")
	})

	t.Run("not found", func(t *testing.T) {
		err := db.UserGroupDelete(9999)
		require.Error(err, "UserGroupDelete did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "UserGroupDelete did not return the expected error")
	})
}

// ############################################################################################## //
// ##################################        Tokens        ###################################### //
// ############################################################################################## //

func TestSqliteDBTokenMigrate(t *testing.T) {
	require := require.New(t)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	createQuery := `CREATE TABLE ` + sqlite_tb_tokens + ` (
		bearer VARCHAR(73) NOT NULL UNIQUE,
		jwt TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	indexQuery := `CREATE INDEX idx_tokens_bearer ON ` + sqlite_tb_tokens + ` (bearer)`

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	t.Run("table", func(t *testing.T) {
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = '" + sqlite_tb_tokens + "'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(createQuery, schema)
	})

	t.Run("index", func(t *testing.T) {
		row, err := db.QueryRow("SELECT sql FROM sqlite_schema WHERE name = 'idx_tokens_bearer'")
		require.NoError(err, "QueryRow returned an error: %s", err)

		var schema string
		err = row.Scan(&schema)
		require.NoError(err, "Scan returned an error: %s", err)
		require.Equal(indexQuery, schema)
	})
}

func TestSqliteDBTokenCreate(t *testing.T) {
	require := require.New(t)
	SetAuthSecret(test_secret)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	t.Run("empty userID", func(t *testing.T) {
		_, err := db.TokenCreate(0, want.username, want.name, false)
		require.Error(err, "TokenCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenCreate did not return the expected error")
	})

	t.Run("empty username", func(t *testing.T) {
		_, err := db.TokenCreate(want.userID, "", want.name, false)
		require.Error(err, "TokenCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenCreate did not return the expected error")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := db.TokenCreate(want.userID, "", want.name, false)
		require.Error(err, "TokenCreate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenCreate did not return the expected error")
	})

	t.Run("valid", func(t *testing.T) {
		bearer, err := db.TokenCreate(want.userID, want.username, want.name, false)
		require.NoError(err, "TokenCreate returned an error: %s", err)
		require.Len(bearer, 72)
	})
}

func TestSqliteDBTokenGet(t *testing.T) {
	require := require.New(t)
	SetAuthSecret(test_secret)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	bearer, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	t.Run("empty bearer", func(t *testing.T) {
		_, err := db.TokenGet("")
		require.Error(err, "TokenGet did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenGet did not return the expected error")
	})

	t.Run("invalid bearer", func(t *testing.T) {
		_, err := db.TokenGet("not_a_bearer")
		require.Error(err, "TokenGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "TokenGet did not return the expected error")
	})

	t.Run("valid", func(t *testing.T) {
		token, err := db.TokenGet(bearer)
		require.NoError(err, "TokenGet returned an error: %s", err)
		require.Equal(want.userID, token.UserID)
		require.Equal(want.username, token.Username)
		require.Equal(want.name, token.Name)
		require.False(token.IsAdmin)
	})
}

func TestSqliteDBTokenUpdate(t *testing.T) {
	require := require.New(t)
	SetAuthSecret(test_secret)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	bearer, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	got, err := db.TokenGet(bearer)
	require.NoError(err, "TokenGet returned an error: %s", err)

	got.UserID = 5678
	got.Username = "bobs"
	got.Name = "Bob Smith"
	got.IsAdmin = true

	t.Run("empty bearer", func(t *testing.T) {
		err := db.TokenUpdate("", got)
		require.Error(err, "TokenUpdate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenUpdate did not return the expected error")
	})

	t.Run("nil claims", func(t *testing.T) {
		err := db.TokenUpdate(bearer, nil)
		require.Error(err, "TokenUpdate did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenUpdate did not return the expected error")
	})

	t.Run("empty claims", func(t *testing.T) {
		err := db.TokenUpdate(bearer, &Claims{})
		require.NoError(err, "TokenUpdate did not return an error")
	})

	t.Run("valid", func(t *testing.T) {
		err := db.TokenUpdate(bearer, got)
		require.NoError(err, "TokenUpdate returned an error: %s", err)

		token, err := db.TokenGet(bearer)
		require.NoError(err, "TokenGet returned an error: %s", err)
		require.Equal(got.UserID, token.UserID)
		require.Equal(got.Username, token.Username)
		require.Equal(got.Name, token.Name)
		require.Equal(got.IsAdmin, token.IsAdmin)
	})
}

func TestSqliteDBTokenDelete(t *testing.T) {
	require := require.New(t)
	SetAuthSecret(test_secret)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	bearer, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	t.Run("empty bearer", func(t *testing.T) {
		err := db.TokenDelete("")
		require.Error(err, "TokenDelete did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "TokenDelete did not return the expected error")
	})

	t.Run("invalid bearer", func(t *testing.T) {
		err := db.TokenDelete("not_a_bearer")
		require.NoError(err, "TokenDelete did not return an error")
	})

	t.Run("valid", func(t *testing.T) {
		err := db.TokenDelete(bearer)
		require.NoError(err, "TokenDelete returned an error: %s", err)

		_, err = db.TokenGet(bearer)
		require.Error(err, "TokenGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "TokenGet did not return the expected error")
	})
}

func TestSqliteDBTokenClean(t *testing.T) {
	require := require.New(t)
	SetAuthSecret(test_secret)
	db := TestSqliteAuthDBSetup(t)
	defer db.Close()
	defer DeleteDB(TestAuthDBName)

	err := TokensMigrate(db)
	require.NoError(err, "TokensMigrate returned an error: %s", err)

	bearer1, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	bearer2, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	bearer3, err := db.TokenCreate(want.userID, want.username, want.name, false)
	require.NoError(err, "TokenCreate returned an error: %s", err)

	_, err = db.Exec(
		"UPDATE "+sqlite_tb_tokens+" SET expires_at = ? where bearer = ?",
		time.Now().Add(-(sessionMaxExpires + 1*time.Minute)),
		bearer1,
	)
	require.NoError(err, "Exec returned an error: %s", err)

	_, err = db.Exec(
		"UPDATE "+sqlite_tb_tokens+" SET created_at = ? where bearer = ?",
		time.Now().Add(-(sessionMaxExpires + 1*time.Minute)),
		bearer2,
	)
	require.NoError(err, "Exec returned an error: %s", err)

	t.Run("valid", func(t *testing.T) {
		err := db.TokenClean()
		require.NoError(err, "TokenClean returned an error: %s", err)

		_, err = db.TokenGet(bearer1)
		require.Error(err, "TokenGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "TokenGet did not return the expected error")

		_, err = db.TokenGet(bearer2)
		require.Error(err, "TokenGet did not return an error")
		require.ErrorIs(err, sql.ErrNoRows, "TokenGet did not return the expected error")

		got, err := db.TokenGet(bearer3)
		require.NoError(err, "TokenGet returned an error: %s", err)
		require.NotNil(got)
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

	DeleteDB(db_file)
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

	DeleteDB(db_file)
}
*/
