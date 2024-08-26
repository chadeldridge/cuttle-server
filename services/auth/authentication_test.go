package auth

import (
	"testing"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/chadeldridge/cuttle-server/db"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticationValidatePassword(t *testing.T) {
	require := require.New(t)

	t.Run("empty password", func(t *testing.T) {
		password := ""
		err := ValidatePassword(password)
		require.Error(err, "ValidatePassword returned an error: %s", err)
		require.ErrorIs(err, ErrPwEmpty, "ValidatePassword did not return the correct error")
	})

	t.Run("too short", func(t *testing.T) {
		password := "password"
		err := ValidatePassword(password)
		require.Error(err, "ValidatePassword did not return an error")
		require.ErrorIs(err, ErrPwTooShort, "ValidatePassword did not return the correct error")
	})

	t.Run("too long", func(t *testing.T) {
		// 73 characters
		password := `wMtHN.Yu5yr&4(ZsfeF?k6{"mzh;,Lq)*aWC]A!D@38SXT^Bcgc59bd;sm#NAuBLHWw[%e>2a`
		err := ValidatePassword(password)
		require.Error(err, "ValidatePassword did not return an error")
		require.ErrorIs(err, ErrPwTooLong, "ValidatePassword did not return the correct error")
	})

	t.Run("sequential", func(t *testing.T) {
		password := "a809wep[04hew398pabcd"
		err := ValidatePassword(password)
		require.Error(err, "ValidatePassword did not return an error")
		require.ErrorIs(err, ErrPwSequential, "ValidatePassword did not return the correct error")
	})

	t.Run("valid no spaces", func(t *testing.T) {
		password := "MyT0tallyC0mpl3xP@ssw0rd!"
		err := ValidatePassword(password)
		require.NoError(err, "ValidatePassword returned an error: %s", err)
	})

	t.Run("valid with spaces", func(t *testing.T) {
		password := "My T0tally C0mpl3x Passw0rd"
		err := ValidatePassword(password)
		require.NoError(err, "ValidatePassword returned an error: %s", err)
	})

	t.Run("max length", func(t *testing.T) {
		// 72 characters is the maximum length for a bcrypt password.
		password := `wMtHN.Yu5yr&4(ZsfeF?k6{"mzh;,Lq)*aWC]A!D@38SXT^Bcgc59bd;sm#NAuBLHWw[%e>2`
		err := ValidatePassword(password)
		require.NoError(err, "HashPassword returned an error: %s", err)
	})
}

func TestAuthenticationHasSequential(t *testing.T) {
	require := require.New(t)

	t.Run("none", func(t *testing.T) {
		password := "a809wep[04hew398p"
		require.False(HasSequential(password), "HasSequential returned true")
	})

	t.Run("abc", func(t *testing.T) {
		password := "a809wep[04hew398pabc"
		require.False(HasSequential(password), "HasSequential returned true")
	})

	t.Run("abcd", func(t *testing.T) {
		password := "a809wep[04hew398pabcd"
		require.True(HasSequential(password), "HasSequential returned true")
	})

	t.Run("aaaa", func(t *testing.T) {
		password := "a809wep[04hew398pabcd"
		require.True(HasSequential(password), "HasSequential returned true")
	})

	t.Run("1234", func(t *testing.T) {
		password := "a809w1234ep[04hew398p"
		require.True(HasSequential(password), "HasSequential returned true")
	})

	t.Run("1212", func(t *testing.T) {
		password := "a809w1234ep[04hew398p"
		require.True(HasSequential(password), "HasSequential returned true")
	})

	t.Run("5555", func(t *testing.T) {
		password := "a809w5555ep[04hew398p"
		require.True(HasSequential(password), "HasSequential returned true")
	})

	t.Run("sequential /", func(t *testing.T) {
		password := "a809wep[04he////w398p"
		require.True(HasSequential(password), "HasSequential returned true")
	})
}

func TestAuthenticationIsIntSequential(t *testing.T) {
	require := require.New(t)

	t.Run("false", func(t *testing.T) {
		require.False(IsIntSequential(int("1"[0]), int("9"[0])), "IsSequential returned true")
	})

	t.Run("true", func(t *testing.T) {
		require.True(IsIntSequential(int("0"[0]), int("1"[0])), "IsSequential returned false")
	})

	t.Run("true reverse", func(t *testing.T) {
		require.True(IsIntSequential(int("1"[0]), int("0"[0])), "IsSequential returned false")
	})

	t.Run("true repeated", func(t *testing.T) {
		require.True(IsIntSequential(int("1"[0]), int("1"[0])), "IsSequential returned false")
	})
}

func TestAuthenticationIsCharSequential(t *testing.T) {
	require := require.New(t)

	t.Run("false", func(t *testing.T) {
		require.False(IsCharSequential("a", "p"), "IsSequential returned true")
	})

	t.Run("false char1 int", func(t *testing.T) {
		require.False(IsCharSequential("1", "p"), "IsSequential returned true")
	})

	t.Run("false char2 int", func(t *testing.T) {
		require.False(IsCharSequential("a", "2"), "IsSequential returned true")
	})

	t.Run("true", func(t *testing.T) {
		require.True(IsCharSequential("a", "b"), "IsSequential returned false")
	})

	t.Run("true reverse", func(t *testing.T) {
		require.True(IsCharSequential("b", "a"), "IsSequential returned false")
	})

	t.Run("true repeated", func(t *testing.T) {
		require.True(IsCharSequential("a", "a"), "IsSequential returned false")
	})
}

func TestAuthenticationIsComplex(t *testing.T) {
	require := require.New(t)

	t.Run("false", func(t *testing.T) {
		password := "MyT0tallyNotSpecialPassword"
		require.False(IsComplex(password), "IsComplex returned true")
	})

	t.Run("true with special char", func(t *testing.T) {
		password := "MyT0tallyC0mpl3xP@ssw0rd!"
		require.True(IsComplex(password), "IsComplex returned false")
	})

	t.Run("true with spaces", func(t *testing.T) {
		password := "My T0tally C0mpl3x Passw0rd"
		require.True(IsComplex(password), "IsComplex returned false")
	})
}

func TestAuthenticationHashPassword(t *testing.T) {
	require := require.New(t)

	t.Run("empty password", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)
		require.Error(err, "HashPassword did not return an error")
		require.ErrorIs(err, ErrPwEmpty, "HashPassword did not return the correct error")
		require.Empty(hash, "HashPassword returned a hash")
	})

	t.Run("valid password", func(t *testing.T) {
		password := "MyT0tallyC0mpl3xP@ssw0rd!"
		hash, err := HashPassword(password)
		require.NoError(err, "HashPassword returned an error: %s", err)
		require.NotEmpty(hash, "HashPassword did not return a hash")

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		require.NoError(err, "bcrypt.CompareHashAndPassword did not return the expected hash")
	})

	t.Run("valid password with spaces", func(t *testing.T) {
		password := "My T0tally C0mpl3x Passw0rd"
		hash, err := HashPassword(password)
		require.NoError(err, "HashPassword returned an error: %s", err)
		require.NotEmpty(hash, "HashPassword did not return a hash")

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		require.NoError(err, "bcrypt.CompareHashAndPassword did not return the expected hash")
	})
}

func TestAuthenticationAuthenticateUser(t *testing.T) {
	require := require.New(t)
	want := struct{ username, name, password, groups string }{
		username: "testUser1",
		name:     "Test User 1",
		password: "My T0tally C0mpl3x Passw0rd",
		groups:   "[]",
	}
	authDB := db.TestSqliteAuthDBSetup(t)
	defer authDB.Close()
	defer db.DeleteDB(db.TestAuthDBName)

	// Setup the test tables.
	err := authDB.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("nil authDB", func(t *testing.T) {
		u, err := AuthenticateUser(nil, want.username, "")
		require.Error(err, "AuthenticateUser did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "HashPassword did not return the correct error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	t.Run("empty username", func(t *testing.T) {
		u, err := AuthenticateUser(authDB, "", want.password)
		require.Error(err, "AuthenticateUser did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "HashPassword did not return the correct error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	t.Run("empty password", func(t *testing.T) {
		u, err := AuthenticateUser(authDB, want.username, "")
		require.Error(err, "AuthenticateUser did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "HashPassword did not return the correct error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	t.Run("user not found", func(t *testing.T) {
		u, err := AuthenticateUser(authDB, want.username, want.password)
		require.Error(err, "AuthenticateUser did not return an error")
		require.ErrorIs(err, ErrUserNotFound, "HashPassword did not return the correct error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	t.Run("empty authDB", func(t *testing.T) {
		u, err := AuthenticateUser(&db.SqliteDB{}, want.username, want.password)
		require.Error(err, "AuthenticateUser did not return an error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	hash, err := HashPassword(want.password)
	require.NoError(err, "HashPassword returned an error: %s", err)

	user, err := authDB.UserCreate(want.username, want.name, hash, want.groups)
	require.NoError(err, "UserCreate returned an error: %s", err)
	require.NotEmpty(user, "UserCreate did not return a user")

	t.Run("wrong password", func(t *testing.T) {
		u, err := AuthenticateUser(authDB, want.username, "wrongPassword")
		require.Error(err, "AuthenticateUser did not return an error")
		require.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword, "HashPassword did not return the correct error")
		require.Empty(u, "AuthenticateUser returned an non-empty user")
	})

	t.Run("valid password", func(t *testing.T) {
		u, err := AuthenticateUser(authDB, want.username, want.password)
		require.NoError(err, "AuthenticateUser returned an error: %s", err)
		require.Equal(user.ID, u.ID, "AuthenticateUser returned the wrong user")
		require.Equal(user.Username, u.Username, "AuthenticateUser returned the wrong user")
		require.Equal(user.Name, u.Name, "AuthenticateUser returned the wrong user")
		require.Empty(u.Groups, "AuthenticateUser returned the wrong user")
		require.False(u.IsAdmin, "AuthenticateUser returned the wrong user")
		require.Equal(user.Created, u.Created, "AuthenticateUser returned the wrong user")
		require.Equal(user.Updated, u.Updated, "AuthenticateUser returned the wrong user")
	})
}

func TestAuthenticationNewUserFromUserData(t *testing.T) {
	require := require.New(t)
	want := User{
		ID:       1,
		Username: "testUser1",
		Name:     "Test User 1",
		Groups:   make([]int64, 0),
		IsAdmin:  false,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	ud := db.UserData{
		ID:       1,
		Username: want.Username,
		Name:     want.Name,
		Hash:     "hash",
		Groups:   "[]",
		IsAdmin:  false,
		Created:  want.Created,
		Updated:  want.Updated,
	}

	t.Run("empty data", func(t *testing.T) {
		u, err := NewUserFromUserData(db.UserData{})
		require.Error(err, "NewUserFromUserData did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "NewUserFromUserData did not return the correct error")
		require.Empty(u, "NewUserFromUserData returned an non-empty user")
	})

	t.Run("empty groups", func(t *testing.T) {
		badUD := ud
		badUD.Groups = "{}"
		u, err := NewUserFromUserData(badUD)
		require.Error(err, "NewUserFromUserData did not return an error")
		require.Empty(u, "NewUserFromUserData returned an non-empty user")
	})

	t.Run("valid", func(t *testing.T) {
		u, err := NewUserFromUserData(ud)
		require.NoError(err, "NewUserFromUserData returned an error: %s", err)
		require.Equal(want.ID, u.ID, "NewUserFromUserData returned the wrong user")
		require.Equal(want.Username, u.Username, "NewUserFromUserData returned the wrong user")
		require.Equal(want.Name, u.Name, "NewUserFromUserData returned the wrong user")
		require.Equal(want.Groups, u.Groups, "NewUserFromUserData returned the wrong user")
		require.Equal(want.IsAdmin, u.IsAdmin, "NewUserFromUserData returned the wrong user")
		require.Equal(want.Created, u.Created, "NewUserFromUserData returned the wrong user")
		require.Equal(want.Updated, u.Updated, "NewUserFromUserData returned the wrong user")
	})
}

func TestAuthenticationSignup(t *testing.T) {
	require := require.New(t)
	want := struct{ username, name, password, groups string }{
		username: "testUser1",
		name:     "Test User 1",
		password: "My T0tally C0mpl3x Passw0rd",
		groups:   "[]",
	}

	authDB := db.TestSqliteAuthDBSetup(t)
	defer authDB.Close()
	defer db.DeleteDB(db.TestAuthDBName)

	// Setup the test tables.
	err := authDB.AuthMigrate()
	require.NoError(err, "AuthMigrate returned an error: %s", err)

	t.Run("nil authDB", func(t *testing.T) {
		u, err := Signup(nil, want.username, want.name, want.password)
		require.Error(err, "Signup did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "Signup did not return the correct error")
		require.Empty(u, "Signup returned an non-empty user")
	})

	t.Run("bad password", func(t *testing.T) {
		u, err := Signup(authDB, want.username, want.name, "")
		require.Error(err, "Signup did not return an error")
		require.ErrorIs(err, ErrPwEmpty, "Signup did not return the correct error")
		require.Empty(u, "Signup returned an non-empty user")
	})

	t.Run("empty username", func(t *testing.T) {
		u, err := Signup(authDB, "", want.name, want.password)
		require.Error(err, "Signup did not return an error")
		require.ErrorIs(err, core.ErrParamEmpty, "Signup did not return the correct error")
		require.Empty(u, "Signup returned an non-empty user")
	})

	hash, err := HashPassword(want.password)
	require.NoError(err, "HashPassword returned an error: %s", err)

	user, err := authDB.UserCreate("testUser2", "Test User 2", hash, want.groups)
	require.NoError(err, "UserCreate returned an error: %s", err)
	require.NotEmpty(user, "UserCreate did not return a user")

	t.Run("user exists", func(t *testing.T) {
		u, err := Signup(authDB, "testUser2", "Test User 2", want.password)
		require.Error(err, "Signup did not return an error")
		require.ErrorIs(err, db.ErrExists, "Signup did not return the correct error")
		require.Empty(u, "Signup returned an non-empty user")
	})

	t.Run("valid", func(t *testing.T) {
		u, err := Signup(authDB, want.username, want.name, want.password)
		require.NoError(err, "Signup returned an error: %s", err)
		require.Equal(want.username, u.Username, "Signup returned the wrong user")
		require.Equal(want.name, u.Name, "Signup returned the wrong user")
		require.Empty(u.Groups, "Signup returned the wrong user")
		require.False(u.IsAdmin, "Signup returned the wrong user")
		require.NotZero(u.Created, "Signup returned the wrong user")
		require.NotZero(u.Updated, "Signup returned the wrong user")

		data, err := authDB.UserGet(int64(u.ID))
		require.NoError(err, "UserGet returned an error: %s", err)
		err = bcrypt.CompareHashAndPassword([]byte(data.Hash), []byte(want.password))
		require.NoError(err, "bcrypt.CompareHashAndPassword did not return the expected hash")
	})
}

func TestAuthenticationUserHasGroup(t *testing.T) {
	require := require.New(t)
	// We only need a list of groups for this test.

	t.Run("empty groups", func(t *testing.T) {
		user := User{Groups: []int64{}}
		require.False(user.HasGroup(1), "HasGroup returned true")
	})

	t.Run("valid", func(t *testing.T) {
		user := User{Groups: []int64{1, 2, 3}}
		require.True(user.HasGroup(2), "HasGroup returned false")
	})
}

/*
func TestAuthenticationUserGroupHasGroup(t *testing.T) {
	require := require.New(t)

	t.Run("empty groups", func(t *testing.T) {
		groups := make([]int64, 0)
		require.False(HasGroup(groups, 1), "HasGroup returned true")
	})

	t.Run("group not found", func(t *testing.T) {
		groups := []int64{1, 2, 3}
		require.False(HasGroup(groups, 4), "HasGroup returned true")
	})

	t.Run("group found", func(t *testing.T) {
		groups := []int64{1, 2, 3}
		require.True(HasGroup(groups, 2), "HasGroup returned false")
	})
}
*/
