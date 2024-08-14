package core

import "fmt"

var (
	ErrParamEmpty     = fmt.Errorf("parameter is empty")
	ErrParamTooLong   = fmt.Errorf("parameter is too long")
	ErrParamBadFormat = fmt.Errorf("parameter format is invalid")
)
