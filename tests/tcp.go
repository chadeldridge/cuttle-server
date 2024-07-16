package tests

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/chadeldridge/cuttle/connections"
)

const TCPDefaultTimeout = time.Second * 3

var ErrInvalidTestType = fmt.Errorf("invalid test type")

// TCPTest is a struct that holds the parameters for TCP tests.
type TCPTest struct {
	testType string
	port     string
	timeout  time.Duration
}

// NewTCPPortHalfOpen creates a new Test for tcp port open with the given parameters.
// name: The name of the test.
// mustSucceed: If false, the Tile will continue with the test stack if this test fails.
//
// These TestArg will be evaluated:
// "timeout": (int, int64, time.Duration) int/int64 will be converted into time.Second * int.
func NewTCPPortHalfOpen(name string, mustSucceed bool, port int, args ...TestArg) Test {
	return Test{
		Name:        name,
		MustSucceed: mustSucceed,
		Tester: &TCPTest{
			testType: "port_half_open",
			port:     strconv.Itoa(port),
			timeout:  getTCPTimeout(args),
		},
	}
}

// NewTCPPortOpen creates a new Test for tcp port open with the given parameters.
// name: The name of the test.
// mustSucceed: If false, the Tile will continue with the test stack if this test fails.
//
// These TestArg will be evaluated:
// "timeout": (int, int64, time.Duration) int/int64 will be converted into time.Second * int.
func NewTCPPortOpen(name string, mustSucceed bool, port int, args ...TestArg) Test {
	return Test{
		Name:        name,
		MustSucceed: mustSucceed,
		Tester: &TCPTest{
			testType: "port_open",
			port:     strconv.Itoa(port),
			timeout:  getTCPTimeout(args),
		},
	}
}

func getTCPTimeout(args []TestArg) time.Duration { return GetTimeout(args, TCPDefaultTimeout) }

// Run evaluates TCPTest.testType and runs the appropriate test, passing along server and args.
func (t TCPTest) Run(server connections.Server, args ...TestArg) error {
	switch t.testType {
	case "port_half_open":
		return PortHalfOpen(t, server, args...)
	case "port_open":
		return PortOpen(t, server, args...)
	default:
		return ErrInvalidTestType
	}
}

// PortHalfOpen performs a simple tcp port open test against the server and ignores close errors.
func PortHalfOpen(t TCPTest, server connections.Server, args ...TestArg) error {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(server.GetHostAddr(), t.port),
		t.timeout,
	)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

// PortOpen performs a simple tcp port open test against the server and makes sure it successfully
// closes the connection.
func PortOpen(t TCPTest, server connections.Server, args ...TestArg) error {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(server.GetHostAddr(), t.port),
		t.timeout,
	)
	if err != nil {
		return err
	}

	return conn.Close()
}
