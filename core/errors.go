package core

import "fmt"

var (
	ErrParamEmpty    = fmt.Errorf("parameter is empty")
	ErrInvalidFormat = fmt.Errorf("invalid format")
)
