package connections

import "errors"

var (
	// Invalid Connector Errors
	ErrInvalidEmtpyUser    = errors.New("user is empty")
	ErrInvalidNoAuthMethod = errors.New("no AuthMethod set")

	// Client Errors
	ErrNotConnected = errors.New("not connected")

	// Session Errors
	ErrSessionActive = errors.New("cannot close connection, session active")

	// Run Errors
	ErrEmtpyCmd = errors.New("cmd is empty")
	ErrEmtpyExp = errors.New("exp is empty")
)

type Connector interface {
	// IsConnected returns true if there is a connection established to the server.
	IsConnected() bool
	// IsActive returns true if the Connector is being used such as if there's an open session
	// with the SSHConnector.
	IsActive() bool
	// Protocol returns the Protocol enum for this Connector type.
	Protocol() Protocol
	// GetUser returns the username used by this Connector.
	GetUser() string
	// DefaultPort returns the default port number used by this Connector type.
	DefaultPort() int
	// IsEmpty checks that fields populated by New contain data.
	IsEmpty() bool
	// IsValid returns true if no errors are returned by Connector.Validate(). If IsValid
	// returns true then you should be able to create a connection using this Connector.
	IsValid() bool
	// Validate checks each Connector field required for connecting to an endpoint are returns
	// an error if anything is missing or incoorect.
	Validate() error
	// TestConnection creates a connection to the server and performs a minimal command test
	// such as a basic echo for ssh. Logs and Results are handled the same way as with
	// Connector.Run().
	TestConnection(bufs Buffers) error
	// Run executes the given cmd(command) against the server, if exp(expect) != "" performs a
	// match of expect against the output of the command. The output of command is sent to
	// Server.Log() and the expect is sent to Server.PrintResults(). Results will either be
	// "ok" or "failed" with the error.
	// Example:
	// Connector.Run(server, "echo 'we did it'", "we did it")
	// Logs Buffer
	// 2024/05/30 12:15:42 debian@test.home:~ we did it
	// Results Buffer
	// 2024/05/30 12:15:42: test.home...ok
	Run(bufs Buffers, cmd string, exp string) error
	// Open creates a connection to the server.
	Open(addr string, bufs Buffers) error
	// Close ends the connecton to the server. Setting force to true will close the connection
	// even if there is an active session.
	Close(force bool) error
}
