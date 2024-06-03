package connections

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// iMock implements the Connector{} interface for MockConnector.

func (c MockConnector) IsConnected() bool  { return c.isConnected }
func (c MockConnector) IsActive() bool     { return c.hasSession }
func (c MockConnector) Protocol() Protocol { return MockProtocol }
func (c MockConnector) User() string       { return c.user }
func (c MockConnector) DefaultPort() int   { return MockDefaultPort }
func (c MockConnector) IsEmpty() bool      { return c.user == "" }
func (c MockConnector) IsValid() bool      { return c.user != "" }

func (c MockConnector) Validate() error {
	if c.user == "" {
		return ErrInvalidEmtpyUser
	}

	return nil
}

func (c *MockConnector) Open(server Server) error {
	if err := c.Validate(); err != nil {
		return err
	}

	c.isConnected = true
	return nil
}

func (c MockConnector) TestConnection(server Server) error {
	localhost, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("connections.SSHConnector.TestConnection: error retrieving local hostname: %s", err)
	}

	expect := fmt.Sprintf("cuttle from %s ok", localhost)
	return c.Run(server, fmt.Sprintf("echo '%s'", expect), expect)
}

func (c MockConnector) Run(server Server, cmd, exp string) error {
	if !c.isConnected {
		return ErrNotConnected
	}

	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ErrEmtpyCmd
	}

	if exp == "" {
		return ErrEmtpyExp
	}

	// We have to split cmd into the command name and args for exec to work. This adds
	// complication but we do not have a choice.
	parts := strings.SplitN(cmd, " ", 2)
	eventTime := time.Now()
	out, err := exec.Command(parts[0], parts[1]).Output()
	if err != nil {
		server.Log(eventTime, err.Error())
		server.PrintResults(eventTime, "error", err)
		return err
	}

	// Log the full output of the command.
	server.Log(eventTime, string(out))

	ok := foundExpect(out, exp)
	if !ok {
		server.PrintResults(eventTime, "failed", nil)
	}

	server.PrintResults(eventTime, "ok", nil)
	return nil
}

func (c *MockConnector) Close(force bool) error {
	if !c.isConnected {
		return ErrNotConnected
	}

	c.isConnected = false
	return nil
}
