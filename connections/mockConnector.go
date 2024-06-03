package connections

import "errors"

const (
	MockDefaultPort = 555
	MockProtocol    = MOCK
)

// MockConnector holds the minimal information needed for creating a mock Connector interface.
type MockConnector struct {
	user        string
	isConnected bool
	hasSession  bool
}

// NewMockConnector creates a MockConnector to simulate connecting to a server.
func NewMockConnector(username string) (MockConnector, error) {
	m := MockConnector{}

	err := m.SetUser(username)
	return m, err
}

// SetUser sets the username to be used for the connection credentials.
func (h *MockConnector) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.MockHandler.SetUser: username was empty")
	}

	h.user = username
	return nil
}
