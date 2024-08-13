package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/chadeldridge/cuttle-server/db"
	"golang.org/x/crypto/bcrypt"
)

const (
	pw_min_length          = 12
	pw_max_length          = 72
	pw_complexity_required = 2
	pw_allowed_sequential  = 3
)

var (
	ErrInvalidName  = fmt.Errorf("invalid name")
	ErrUserNotFound = fmt.Errorf("user not found")
	// Password errors.
	ErrPwEmpty      = fmt.Errorf("password is empty")
	ErrPwTooShort   = fmt.Errorf("password too short, must be at least %d characters", pw_min_length)
	ErrPwTooLong    = fmt.Errorf("password too long, must be at most %d characters", pw_max_length)
	ErrPwSequential = fmt.Errorf("password contains %d or more sequential characters", pw_allowed_sequential+1)
	ErrPwNotComplex = fmt.Errorf(
		"password must contain at least one of the following: special character (non alpha-numeric) or space",
	)

	// Regex patterns for password complexity.
	alpha = regexp.MustCompile(`[a-zA-Z]`)
	dg    = regexp.MustCompile(`[0-9]`)
	sp    = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// ValidatePassword checks a password for complexity and length requirements.
func ValidatePassword(password string) error {
	// Reject empty passwords.
	if password == "" {
		return fmt.Errorf("auth.ValidatePassword: %w", ErrPwEmpty)
	}

	// Reject passwords that are too short.
	if len(password) < pw_min_length {
		return ErrPwTooShort
	}

	// Reject passwords that are too long for bcrypt.
	if len(password) > pw_max_length {
		return ErrPwTooLong
	}

	// INCOMPLETE: Reject common passwords before checking complexity.

	// Reject passwords that contain too many sequential characters.
	if HasSequential(password) {
		return fmt.Errorf("auth.ValidatePassword: %w", ErrPwSequential)
	}

	// If the password is short, check for password complexity.
	if !IsComplex(password) {
		return fmt.Errorf("auth.ValidatePassword: %w", ErrPwNotComplex)
	}

	return nil
}

// HasSequential checks for sequential characters in a password. If the password contains
// more sequential characters than allowed, it will return true.
func HasSequential(password string) bool {
	pwl := len(password)
	end := pwl - 1
	score := 0

	for i, c := range password {
		// If we have more sequential characters than allowed, return true.
		if score == pw_allowed_sequential {
			return true
		}

		// If we are on the last character there is nothing to compare. Return false since
		// we haven't found enough sequential characters yet.
		if i == end {
			return false
		}

		// If there aren't enough characters left to meet the score needed to fail,
		// return false.
		if i+pw_allowed_sequential-score >= pwl {
			return false
		}

		next := string(password[i+1])
		s := string(c)

		// If the character is a special character, check for repeat of the same character.
		if sp.MatchString(s) {
			if s == next {
				score++
				continue
			}

			// If the next character is not the same, reset the score.
			score = 0
			continue
		}

		// If the character is a letter, check for sequential letters.
		if alpha.MatchString(s) {
			// Check to make sure the next character is also a letter.
			if alpha.MatchString(next) {
				if IsCharSequential(s, next) {
					score++
					continue
				}

				// If the next character is not sequential, reset the score.
				score = 0
				continue
			}

			// If the next character is not a letter, reset the score.
			score = 0
			continue
		}

		// At this point we should only be dealing with numbers.
		// If the next character is not a number, reset the score.
		if !dg.MatchString(next) {
			score = 0
			continue
		}

		// Check for sequential numbers.
		if IsIntSequential(int(c), int(password[i+1])) {
			score++
			continue
		}

		// If the next character is not sequential, reset the score.
		score = 0

	}

	return false
}

// IsIntSequential checks if two integers are sequential.
func IsIntSequential(i1 int, i2 int) bool {
	next := int(i2)
	log.Printf("i1: %d, next: %d\n", i1, next)
	if i1 == next || i1+1 == next || i1-1 == next {
		return true
	}

	return false
}

// IsCharSequential checks if two characters are sequential.
func IsCharSequential(char1, char2 string) bool {
	if !alpha.MatchString(char2) {
		return false
	}

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	char1 = strings.ToLower(char1)
	char2 = strings.ToLower(char2)

	if char1 == char2 {
		return true
	}

	for i, c := range alphabet {
		// Go to the next this isn't the character we're looking for.
		if string(c) != char1 {
			continue
		}

		var next, prev string

		// If we reach 'z', then compare the next character to 'a'.
		if i == len(alphabet)-1 {
			next = string(alphabet[0])
		} else {
			next = string(alphabet[i+1])
		}

		// If we reach 'a', then compare the previous character to 'z'.
		if i == 0 {
			prev = string(alphabet[len(alphabet)-1])
		} else {
			prev = string(alphabet[i-1])
		}

		// Return the boolean result of comparing the next character.
		return char2 == next || char2 == prev
	}

	// If we can't find the character in the alphabet, return false. (How though?!)
	return false
}

// IsComplex checks that the password contains at least one special character or whitespace.
func IsComplex(password string) bool { return sp.MatchString(password) }

// HashPassword will hash a password using bcrypt and return the hashed password as a string. It
// will run ValidatePassword on the password before hashing it.
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", fmt.Errorf("auth.HashPassword: %w", err)
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("auth.HashPassword: %w", err)
	}

	return string(p), nil
}

// AuthenticateUser will authenticate a user based on the provided username and password and return
// the user data if successful. You can check error for bcrypt.ErrMismatchedHashAndPassword to
// determine if the password was incorrect.
func AuthenticateUser(authDB db.AuthDB, username string, password string) (User, error) {
	if authDB == nil {
		return User{}, fmt.Errorf("auth.AuthenticateUser: authDb - %w", core.ErrParamEmpty)
	}

	if username == "" {
		return User{}, fmt.Errorf("auth.AuthenticateUser: username - %w", core.ErrParamEmpty)
	}

	if password == "" {
		return User{}, fmt.Errorf("auth.AuthenticateUser: password - %w", core.ErrParamEmpty)
	}

	data, err := authDB.UserGetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("auth.AuthenticateUser: %w", ErrUserNotFound)
		}

		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Hash), []byte(password)); err != nil {
		return User{}, fmt.Errorf("auth.AuthenticateUser: %w", err)
	}

	return NewUserFromUserData(data)
}

// ############################################################################################## //
// ####################################        Users         #################################### //
// ############################################################################################## //

type User struct {
	ID       int64
	Username string
	Name     string
	Groups   []int64
	IsAdmin  bool
	Created  time.Time
	Updated  time.Time
}

func NewUserFromUserData(data db.UserData) (User, error) {
	if data.ID == 0 && data.Username == "" && data.Name == "" {
		return User{}, fmt.Errorf("auth.NewUserFromUserData: data %w", core.ErrParamEmpty)
	}

	GIDs, err := UnmarshGroupIDs([]byte(data.Groups))
	if err != nil {
		return User{}, fmt.Errorf("auth.NewUserFromUserData: %w", err)
	}

	return User{
		ID:       data.ID,
		Username: data.Username,
		Name:     data.Name,
		Groups:   GIDs,
		IsAdmin:  data.IsAdmin,
		Created:  data.Created,
		Updated:  data.Updated,
	}, nil
}

func Signup(authDB db.AuthDB, username, name, password string) (User, error) {
	if authDB == nil {
		return User{}, fmt.Errorf("auth.Signup: authDb - %w", core.ErrParamEmpty)
	}

	password, err := HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("auth.Signup: %w", err)
	}

	data, err := authDB.UserCreate(username, name, password, "[]")
	if err != nil {
		return User{}, fmt.Errorf("auth.Signup: %w", err)
	}
	return NewUserFromUserData(data)
}

func (u *User) HasGroup(id int64) bool {
	for _, group := range u.Groups {
		if group == id {
			return true
		}
	}

	return false
}

// ############################################################################################## //
// ##################################        User Groups        ################################# //
// ############################################################################################## //

type UserGroup struct {
	ID       int64
	Name     string
	Members  []int64
	Profiles map[string]Permissions
}

type UserGroups []UserGroup

func (g UserGroups) HasGroup(id int64) bool {
	for _, group := range g {
		if group.ID == id {
			return true
		}
	}

	return false
}

func (g UserGroups) MarshalIDs() ([]byte, error) {
	var ids []int64
	for _, group := range g {
		ids = append(ids, group.ID)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshGroupIDs(data []byte) ([]int64, error) {
	var ids []int64
	err := json.Unmarshal(data, &ids)

	return ids, err
}
