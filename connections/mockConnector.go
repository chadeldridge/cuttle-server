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
}

func NewMockHandler(username string) (MockConnector, error) {
	m := MockConnector{}

	err := m.SetUser(username)
	return m, err
}

func (h *MockConnector) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.MockHandler.SetUser: username was empty")
	}

	h.user = username
	return nil
}
