package connections

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

const (
	SSHDefaultPort = 22
	SSHProtocol    = SSH
)

type SSHHandler struct {
	IsConnected bool             // Track if we have an active connection to the server.
	HasSession  bool             // Indicates there's an active session so we don't close the connection on it.
	auth        []ssh.AuthMethod // Each auth method will be tried in turn until one works or all fail.
	user        string           // The username to login to the server with.
	*ssh.Client
	*ssh.Session
}

// NewSSH creates an SSH struct and sets the Server, Results, and Logs fields. Results and Logs
// can be set to nil if you don't want to ignore them.
func NewSSH(username string) (SSHHandler, error) {
	s := SSHHandler{}

	err := s.SetUser(username)
	return s, err
}

// SetUser sets the username to be used for connection credentials.
func (h *SSHHandler) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.SSHHandler.SetUser: username was empty")
	}

	//								//
	// Add username validation here //
	//								//
	h.user = username
	return nil
}

// SetKey sets the key signer and appends it as an auth method.
func (h *SSHHandler) SetKey(key ssh.Signer) {
	h.auth = append(h.auth, ssh.PublicKeys(key))
}

// ParseKey parses the private key into a key signer and sends it to SSHHandler.SetKey()
func (h *SSHHandler) ParseKey(privateKey []byte) error {
	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	h.SetKey(key)
	return nil
}

// ParseKeyWithPassphrase parses a passhphrase protected private key into a key signer
// and sends it to SSHHandler.SetKey().
func (h *SSHHandler) ParseKeyWithPassphrase(privateKey, passphrase []byte) error {
	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, passphrase)
	if err != nil {
		return err
	}

	h.SetKey(key)
	return nil
}

// SetPassword sets the password field and appends it as an auth method.
func (h *SSHHandler) SetPassword(password string) {
	h.auth = append(h.auth, ssh.Password(password))
}
