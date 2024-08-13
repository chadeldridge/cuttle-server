package auth

import (
	"testing"

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
		require.False(IsIntSequential(int("1"[0]), "9"[0]), "IsSequential returned true")
	})

	t.Run("true", func(t *testing.T) {
		require.True(IsIntSequential(int("0"[0]), `1`[0]), "IsSequential returned false")
	})

	t.Run("true reverse", func(t *testing.T) {
		require.True(IsIntSequential(int("1"[0]), "0"[0]), "IsSequential returned false")
	})

	t.Run("true repeated", func(t *testing.T) {
		require.True(IsIntSequential(int("1"[0]), "1"[0]), "IsSequential returned false")
	})
}

func TestAuthenticationIsCharSequential(t *testing.T) {
	require := require.New(t)

	t.Run("false", func(t *testing.T) {
		require.False(IsCharSequential("a", "p"), "IsSequential returned true")
	})

	t.Run("false with int", func(t *testing.T) {
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
