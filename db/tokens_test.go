package db

import (
	"testing"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/stretchr/testify/require"
)

const (
	test_secret = "928hqg0pg9a8342892hgp9834hgf98hgp894q498yt2hhgvbncj892948y27"
)

var want = struct {
	userID   int64
	username string
	name     string
	isAdmin  bool
}{
	userID:   int64(1823),
	username: "test",
	name:     "Test User 1",
	isAdmin:  false,
}

func TestCreateJWT(t *testing.T) {
	require := require.New(t)

	t.Run("empty secret", func(t *testing.T) {
		_, err := CreateJWT(want.userID, want.username, want.name, "", want.isAdmin)
		require.Error(err, "CreateJWT did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "CreateJWT did not return the correct error")
	})

	t.Run("valid secret", func(t *testing.T) {
		token, err := CreateJWT(want.userID, want.username, want.name, test_secret, want.isAdmin)
		require.NoError(err, "CreateJWT returned an error")
		require.NotEmpty(token, "CreateJWT returned an empty token")
	})
}

func TestRefreshJWT(t *testing.T) {
	require := require.New(t)
	claims := &Claims{
		UserID:   want.userID,
		Username: want.username,
		Name:     want.name,
		IsAdmin:  want.isAdmin,
	}

	t.Run("empty secret", func(t *testing.T) {
		_, err := RefreshJWT(claims, "")
		require.Error(err, "RefreshJWT did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "RefreshJWT did not return the correct error")
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := RefreshJWT(&Claims{}, test_secret)
		require.NoError(err, "RefreshJWT did not return an error")
	})

	t.Run("valid token", func(t *testing.T) {
		token, err := RefreshJWT(claims, test_secret)
		require.NoError(err, "RefreshJWT returned an error")
		require.NotEmpty(token, "RefreshJWT returned an empty token")
	})
}

func TestParseJWT(t *testing.T) {
	require := require.New(t)

	token, err := CreateJWT(want.userID, want.username, want.name, test_secret, want.isAdmin)
	require.NoError(err, "CreateJWT returned an error")
	require.NotEmpty(token, "CreateJWT returned an empty token")

	t.Run("empty token", func(t *testing.T) {
		_, err := ParseJWT("", test_secret)
		require.Error(err, "ParseJWT did not return an error")
	})

	t.Run("empty secret", func(t *testing.T) {
		_, err := ParseJWT(token, "")
		require.Error(err, "ParseJWT did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "ParseJWT did not return the correct error")
	})

	t.Run("valid token", func(t *testing.T) {
		claims, err := ParseJWT(token, test_secret)
		require.NoError(err, "ParseJWT returned an error")
		require.NotNil(claims, "ParseJWT returned nil claims")
		require.Equal(want.userID, claims.UserID, "ParseJWT returned incorrect userID")
		require.Equal(want.username, claims.Username, "ParseJWT returned incorrect username")
		require.Equal(want.name, claims.Name, "ParseJWT returned incorrect name")
		require.Equal(want.isAdmin, claims.IsAdmin, "ParseJWT returned incorrect isAdmin")
	})
}
