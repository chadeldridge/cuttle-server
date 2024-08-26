package db

import (
	"fmt"
)

var (
	ErrRecordExists = fmt.Errorf("record exists")
	// Invalid parameters
	ErrInvalidID         = fmt.Errorf("invalid ID")
	ErrInvalidAuthType   = fmt.Errorf("invalid auth type")
	ErrInvalidPassphrase = fmt.Errorf("invalid passphrase")
	// Record errors
	ErrNotFound = fmt.Errorf("not found")
	ErrExists   = fmt.Errorf("already exists")
	// Tokens
	ErrTokenExpired = fmt.Errorf("token expired")
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
