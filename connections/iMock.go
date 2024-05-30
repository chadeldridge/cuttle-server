package connections

import (
	"errors"
	"time"
)

// iMock implements the Connector{} interface for MockConnector.

func (h MockConnector) IsConnected() bool                  { return h.isConnected }
func (h MockConnector) IsActive() bool                     { return h.isConnected }
func (h MockConnector) Protocol() Protocol                 { return MockProtocol }
func (h MockConnector) User() string                       { return h.user }
func (h MockConnector) DefaultPort() int                   { return MockDefaultPort }
func (h MockConnector) IsEmpty() bool                      { return h.user == "" }
func (h MockConnector) IsValid() bool                      { return h.user != "" }
func (h MockConnector) TestConnection(server Server) error { return h.Run(server, "echo", "any") }

func (h *MockConnector) Close(force bool) error {
	h.isConnected = false
	return nil
}

func (h *MockConnector) Open(server Server) error {
	if !h.IsValid() {
		return errors.New("connections.MockHandler.Open: cannot open, not in a valid state")
	}

	h.isConnected = true
	return nil
}

func (h MockConnector) Run(server Server, cmd, exp string) error {
	if !h.isConnected {
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
