package connections

import (
	"bytes"
	"errors"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	SSHDefaultPort = 22
	SSHProtocol    = SSH
)

// SSHConnector impletments the Connector interface for SSH connectivity.
type SSHConnector struct {
	isConnected bool             // Track if we have an active connection to the server.
	hasSession  bool             // Indicates there's an active session so we don't close the connection on it.
	auth        []ssh.AuthMethod // Each auth method will be tried in turn until one works or all fail.
	user        string           // The username to login to the server with.
	*ssh.Client
	*ssh.Session
}

// NewSSH creates an SSHHandler struct to be used to connect via SSH to a server.
func NewSSH(username string) (SSHConnector, error) {
	s := SSHConnector{}

	err := s.SetUser(username)
	return s, err
}

// SetUser sets the username to be used for connection credentials.
func (c *SSHConnector) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.SSHHandler.SetUser: username was empty")
	}

	// INCOMPLETE: Add username validation here
	c.user = username
	return nil
}

// AddKeyAuth adds an AuthMethod using the ssh private key.
func (c *SSHConnector) AddKeyAuth(key ssh.Signer) {
	c.auth = append(c.auth, ssh.PublicKeys(key))
}

// ParseKey parses the private key into a key signer and sends it to SSHHandler.AddKeyAuth().
func (c *SSHConnector) ParseKey(privateKey []byte) error {
	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	c.AddKeyAuth(key)
	return nil
}

// ParseKeyWithPassphrase parses a passhphrase protected private key into a key signer
// and sends it to SSHHandler.SetKey().
func (c *SSHConnector) ParseKeyWithPassphrase(privateKey, passphrase []byte) error {
	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, passphrase)
	if err != nil {
		return err
	}

	c.AddKeyAuth(key)
	return nil
}

// AddPasswordAuth adds an AuthMethod using a password.
func (c *SSHConnector) AddPasswordAuth(password string) {
	c.auth = append(c.auth, ssh.Password(password))
}

// OpenSession creates a new single command session.
func (c *SSHConnector) OpenSession(server Server) error {
	// log.Print(" - Creating session...")
	sess, err := c.NewSession()
	if err != nil {
		server.Log(time.Now(), err.Error())
		server.PrintResults(time.Now(), "error", err)
		return err
	}

	c.hasSession = true
	c.Session = sess
	return nil
}

// CloseSession closes an open session.
func (c *SSHConnector) CloseSession() error {
	c.hasSession = false
	if c.Session != nil {
		return c.Session.Close()
	}

	return errors.New("connections.SSHConnector.CloseSession: no session avaiable")
}

// foundExpect returns true if expect matches anywhere in the byte array.
func foundExpect(data []byte, expect string) bool {
	m := bytes.Index(data, []byte(expect))
	return m > -1
}