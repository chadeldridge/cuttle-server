package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testUser1 = UserData{
	Username: "testUser1",
	Name:     "Bob",
	Password: "myGr3atP@ssword",
	Groups:   "{}",
}

func TestUsersUserData(t *testing.T) {
	require := require.New(t)
	user := UserData{
		Username: testUser1.Username,
		Name:     testUser1.Name,
		Password: testUser1.Password,
		Groups:   testUser1.Groups,
	}

	require.Equal(testUser1.Username, user.Username)
	require.Equal(testUser1.Name, user.Name)
	require.Equal(testUser1.Password, user.Password)
	require.Equal(testUser1.Groups, user.Groups)
}

func TestUsersNewUsers(t *testing.T) {
	require := require.New(t)
	db := testDBSetup(t)
	defer db.Close()

	users, err := NewUsers(db)
	require.NoError(err, "NewUsers returned an error: %s", err)
	require.NotNil(users)
}
