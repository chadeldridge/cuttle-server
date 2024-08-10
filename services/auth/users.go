package auth

import (
	"fmt"

	"github.com/chadeldridge/cuttle/db"
)

type ID int

type User struct {
	ID
	Username string
	Name     string
	Groups   []ID
}

func Signup(username, name, password string, udb *db.Users) error {
	// INCOMPLETE: implement password strength checks
	password, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("auth.Signup: %w", err)
	}

	err = udb.Create(username, name, password, "{}")
	if err != nil {
		return fmt.Errorf("auth.Signup: %w", err)
	}
	return nil
}
