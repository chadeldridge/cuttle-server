package tests

import (
	"fmt"

	"github.com/chadeldridge/cuttle-server/services/cuttle/connections"
)

var ErrTestFailed = fmt.Errorf("failed")

type Test struct {
	Name        string
	MustSucceed bool
	Tester
}

type Tester interface {
	Run(server connections.Server, args ...TestArg) error
}
