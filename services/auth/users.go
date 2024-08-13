package auth

import (
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle-server/db"
)

type ID int

type User struct {
	ID
	Username string
	Name     string
	Groups   []ID
	IsAdmin  bool
	Created  time.Time
	Updated  time.Time
}

func Signup(username, name, password string, authDB db.AuthDB) (db.UserData, error) {
	// INCOMPLETE: implement password strength checks
	password, err := HashPassword(password)
	if err != nil {
		return db.UserData{}, fmt.Errorf("auth.Signup: %w", err)
	}

	user, err := authDB.UserCreate(username, name, password, "{}")
	if err != nil {
		return user, fmt.Errorf("auth.Signup: %w", err)
	}
	return user, nil
}
