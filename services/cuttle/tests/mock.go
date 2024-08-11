package tests

import "github.com/chadeldridge/cuttle-server/services/cuttle/connections"

type MockTest struct {
	fail bool
}

func NewMockTest(fail bool) Test {
	return Test{
		Name:        "Mock Test",
		MustSucceed: true,
		Tester: &MockTest{
			fail: fail,
		},
	}
}

func (t *MockTest) Run(server connections.Server, args ...TestArg) error {
	if t.fail {
		return ErrTestFailed
	}

	return nil
}
