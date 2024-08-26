package core

import "fmt"

var (
	ErrParamEmpty     = fmt.Errorf("parameter is empty")
	ErrParamInvalid   = fmt.Errorf("parameter is invalid")
	ErrParamTooLong   = fmt.Errorf("parameter is too long")
	ErrParamBadFormat = fmt.Errorf("parameter format is invalid")
)
