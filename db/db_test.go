package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testDBRoot = "/tmp/cuttle/db"

func TestDBSetDBRoot(t *testing.T) {
	require := require.New(t)

	t.Run("empty rootDir", func(t *testing.T) {
		err := SetDBRoot("")
		require.Error(err, "SetDBRoot did not return an error")
		require.Equal("db.SetDBRoot: rootDir is empty", err.Error(), "SetDBRoot did not return the expected error")
		require.Equal(GenDBFolder(), db_folder)
	})

	t.Run("valid", func(t *testing.T) {
		err := SetDBRoot(testDBRoot)
		require.NoError(err, "SetDBRoot returned an error: %s", err)
		require.Equal(testDBRoot, db_folder)
	})
}
