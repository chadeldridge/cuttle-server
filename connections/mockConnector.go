package connections

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	MockDefaultPort = 555
	MockProtocol    = MOCK
)

// MockConnector holds the minimal information needed for creating a mock Connector interface.
type MockConnector struct {
	user         string
	isConnected  bool
	hasSession   bool
	connOpenErr  bool
	sessOpenErr  bool
	connCloseErr bool
	sessCloseErr bool
}

// NewMockConnector creates a MockConnector to simulate connecting to a server.
func NewMockConnector(username string) (MockConnector, error) {
	m := MockConnector{}

	err := m.SetUser(username)
	return m, err
}

// SetUser sets the username to be used for the connection credentials.
func (c *MockConnector) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.MockHandler.SetUser: username was empty")
	}

	c.user = username
	return nil
}

func (c *MockConnector) ErrOnConnectionOpen(do bool)  { c.connOpenErr = do }
func (c *MockConnector) ErrOnConnectionClose(do bool) { c.connCloseErr = do }
func (c *MockConnector) ErrOnSessionOpen(do bool)     { c.sessOpenErr = do }
func (c *MockConnector) ErrOnSessionClose(do bool)    { c.sessCloseErr = do }

// OpenSession creates a new single command session.
func (c *MockConnector) OpenSession(server Server) error {
	// log.Print(" - Creating session...")
	if !c.isConnected {
		return ErrNotConnected
	}

	if c.sessOpenErr {
		return errors.New("error opening session because you asked me to")
	}

	c.hasSession = true
	return nil
}

// CloseSession closes an open session.
func (c *MockConnector) CloseSession() error {
	if c.sessCloseErr {
		return errors.New("error closing session because you asked me to")
	}

	c.hasSession = false
	return nil
}

//						//
//	Connector Interface Implementation	//
//						//

func (c MockConnector) IsConnected() bool  { return c.isConnected }
func (c MockConnector) IsActive() bool     { return c.hasSession }
func (c MockConnector) Protocol() Protocol { return MockProtocol }
func (c MockConnector) GetUser() string    { return c.user }
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

	if c.connOpenErr {
		return errors.New("error opening connection because you asked me to")
	}

	c.isConnected = true
	return nil
}

func (c MockConnector) TestConnection(server Server) error {
	expect := "cuttle ok"
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

	if c.hasSession && !force {
		return ErrSessionActive
	}

	if c.connCloseErr {
		return errors.New("error closing connection because you asked me to")
	}

	c.isConnected = false
	return nil
}
