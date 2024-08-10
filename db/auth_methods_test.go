package db

import (
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

/*
var (
	testDBRoot = "/tmp/cuttle/db"
	testPass   = []byte("testUserP@ssw0rd")
	keyPass    = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABATucg6QV
b74QXyKzG7c6YAAAAAZAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIN8GWe3xMFt/5zSP
xbFK7UlOCB72cCvTec2X1fwAFtYgAAAAoCgV9C/P0QHNfo1edW3BgnBQ1bMOpKVxzUkQ7Q
FIHLIj5vRP4Sv7P6d2u4KnVaCsvIuhVyqductwQskVBSsHPU3HwTPQVZZ0Lu8P3cci7oBc
OiOUXdWp4VAqxTXGkpoTs7Kr/WMavOB2C+/AqgWdOhpICpLxAVk5knuXTK9OvSD34EbC0l
GOO5fZbTGQ1XE1ihvWiIAkUn1XyLaBa3xzOZc=
-----END OPENSSH PRIVATE KEY-----`)
	keyNoPass = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD4p4CaynaubF35hzOcEXg6e/mXM4wlluZBKW9FMg8MegAAAKC8UmL4vFJi
+AAAAAtzc2gtZWQyNTUxOQAAACD4p4CaynaubF35hzOcEXg6e/mXM4wlluZBKW9FMg8Meg
AAAECwBTmJkCxA2UyiNnP5Mh3ampIMnZt+wegxE5jqySmfAvingJrKdq5sXfmHM5wReDp7
+ZczjCWW5kEpb0UyDwx6AAAAGGNlbGRyaWRnZUBDRS1PRkZJQ0UtTUFJTgECAwQF
-----END OPENSSH PRIVATE KEY-----`)
)
*/

func TestAuthMethodsAMDBMigrate(t *testing.T) {
	testDBSetup(t)
	deleteDB(amdb_file)
	_, err := os.Stat(db_folder + "/" + amdb_file)
	if err == nil && !os.IsNotExist(err) {
		log.Fatalf("TestAuthMethodsAMDBMigrate: failed to delete db file: %s", amdb_file)
	}

	require := require.New(t)
	amdb, err := NewSqliteDB(amdb_file)
	require.NoError(err, "TestAuthMethodsAMDBMigrate returned an error: %s", err)

	err = amdb.Open()
	require.NoError(err, "TestAuthMethodsAMDBMigrate returned an error: %s", err)
	defer amdb.Close()

	t.Run("empty file", func(t *testing.T) {
		// Migrate the auth_methods table
		require.NoError(amdbMigrate(amdb))

		// Ensure the auth_methods table was created
		var count int
		row := amdb.QueryRow("SELECT COUNT(*) FROM auth_methods")
		require.NoError(row.Scan(&count))
		require.Equal(0, count)
	})

	t.Run("already exists", func(t *testing.T) {
		require.NoError(amdbMigrate(amdb))
	})
}

/*
func TestAuthMethods(t *testing.T) {
	db := NewTestDB(t)
	authMethods, err := NewAuthMethods(db)
	require.NoError(t, err)

	// Create a new auth method
}
*/
