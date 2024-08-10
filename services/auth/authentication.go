package auth

import (
	"errors"
	"fmt"

	"github.com/chadeldridge/cuttle/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidName     = fmt.Errorf("invalid name")
	ErrInvalidPassword = fmt.Errorf("invalid password")
	ErrPasswordNoMatch = fmt.Errorf("incorrect password")
	ErrUserNotFound    = fmt.Errorf("user not found")
)

func HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("auth.HashPassword: %w", err)
	}

	return string(p), nil
}

func IsMatch(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func AuthenticateUser(repo *db.Users, username string, password string) (User, error) {
	data, err := repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return User{}, fmt.Errorf("auth.AuthenticateUser: %w", ErrUserNotFound)
		}

		return User{}, err
	}

	if !IsMatch(data.Password, password) {
		return User{}, fmt.Errorf("auth.AuthenticateUser: %w", ErrPasswordNoMatch)
	}

	GIDs, err := UnmarshGroupIDs([]byte(data.Groups))
	if err != nil {
		return User{}, fmt.Errorf("auth.AuthenticateUser: %w", err)
	}

	return User{
		ID:       ID(data.ID),
		Username: data.Username,
		Name:     data.Name,
		Groups:   GIDs,
	}, nil
}
