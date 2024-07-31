package auth

import "github.com/chadeldridge/cuttle/db"

type ID int

type User struct {
	ID
	Username string
	Groups   []ID
}

func GetUser(repo *db.Users, username string) (db.UserData, error) {
	return repo.GetByName(username)
}
