package db

import (
	"fmt"
)

var (
	ErrRecordExists = fmt.Errorf("record exists")
	// Invalid parameters
	ErrInvalidID         = fmt.Errorf("invalid ID")
	ErrInvalidUsername   = fmt.Errorf("invalid username")
	ErrInvalidName       = fmt.Errorf("invalid name")
	ErrInvalidAuthType   = fmt.Errorf("invalid auth type")
	ErrInvalidPassphrase = fmt.Errorf("invalid passphrase")
	// Users
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserExists   = fmt.Errorf("user exists")
	// User Groups
	ErrUserGroupNotFound = fmt.Errorf("user group not found")
	ErrUserGroupExists   = fmt.Errorf("user group exists")
)

/*
var	logger               *core.Logger

var Errors = map[string]error{
	"ErrInvalidID":         fmt.Errorf("invalid ID"),
	"ErrorUserNotFound":    fmt.Errorf("user not found"),
	"ErrorUserExists":      fmt.Errorf("user already exists"),
	"ErrInvalidUsername":   fmt.Errorf("invalid username"),
	"ErrInvalidName":       fmt.Errorf("invalid name"),
	"ErrInvalidAuthType":   fmt.Errorf("invalid auth type"),
	"ErrInvalidPassphrase": fmt.Errorf("invalid passphrase"),
}

var RegisteredErrors = map[string]error{}
*/
