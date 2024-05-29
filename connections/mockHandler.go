package connections

import (
	"errors"
	"time"
)

const (
	MockDefaultPort = 555
	MockProtocol    = MOCK
)

type MockHandler struct {
	user      string
	connected bool
}

// Begin Connector{}
func (h MockHandler) Protocol() Protocol                 { return MockProtocol }
func (h MockHandler) User() string                       { return h.user }
func (h MockHandler) DefaultPort() int                   { return MockDefaultPort }
func (h MockHandler) IsEmpty() bool                      { return h.user == "" }
func (h MockHandler) IsValid() bool                      { return h.user != "" }
func (h *MockHandler) Close()                            { h.connected = false }
func (h MockHandler) TestConnection(server Server) error { return h.Run(server, "echo", "any") }

func (h *MockHandler) Open(server Server) error {
	if !h.IsValid() {
		return errors.New("connections.MockHandler.Open: cannot open, not in a valid state")
	}

	h.connected = true
	return nil
}

func (h MockHandler) Run(server Server, cmd, exp string) error {
	if !h.connected {
		return errors.New("connections.MockHandler.Run: connection has not been openned for MockHandler")
	}

	if cmd == "" {
		return errors.New("connections.MockHandler.Run: cmd was empty")
	}

	if exp == "" {
		return errors.New("connections.MockHandler.Run: exp was empty")
	}

	eventTime := time.Now()
	server.Log(eventTime, "mock ok")
	server.PrintResults(eventTime, "ok", nil)
	return nil
}

// End Connector{}

func NewMockHandler(username string) (MockHandler, error) {
	m := MockHandler{}

	err := m.SetUser(username)
	return m, err
}

func (h *MockHandler) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.MockHandler.SetUser: username was empty")
	}

	h.user = username
	return nil
}
