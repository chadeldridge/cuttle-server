package connections

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	SSHDefaultPort = 22
	SSHProtocol    = SSH
)

// SSHConnector impletments the Connector interface for SSH connectivity.
type SSHConnector struct {
	Name        string           // A unique name for the connector to make it easier to add to a server.
	isConnected bool             // Track if we have an active connection to the server.
	hasSession  bool             // Indicates there's an active session so we don't close the connection on it.
	Auth        []ssh.AuthMethod // Each auth method will be tried in turn until one works or all fail.
	User        string           // The username to login to the server with.
	*ssh.Client
	*ssh.Session
}

// NewSSHConnector creates an SSHConnector struct to be used to connect via SSH to a server.
func NewSSHConnector(name, username string) (SSHConnector, error) {
	s := SSHConnector{}

	if err := s.SetName(name); err != nil {
		return s, err
	}

	return s, s.SetUser(username)
}

// SetName sets a unique Name to make it easier to add to a server.
func (c *SSHConnector) SetName(name string) error {
	if name == "" {
		return errors.New("connections.SSHConnector.SetName: username was empty")
	}

	// INCOMPLETE: Add name validation here
	c.Name = name
	return nil
}

// SetUser sets the User to be used for connection credentials.
func (c *SSHConnector) SetUser(username string) error {
	if username == "" {
		return errors.New("connections.SSHConnector.SetUser: username was empty")
	}

	// INCOMPLETE: Add username validation here
	c.User = username
	return nil
}

// AddPasswordAuth adds an AuthMethod using a password.
func (c *SSHConnector) AddPasswordAuth(password string) {
	c.Auth = append(c.Auth, ssh.Password(password))
}

// AddKeyAuth adds an AuthMethod using the ssh private key.
func (c *SSHConnector) AddKeyAuth(key ssh.Signer) error {
	if key == nil {
		return errors.New("connections.SSHConnector.AddKeyAuth: key was nil")
	}

	c.Auth = append(c.Auth, ssh.PublicKeys(key))
	return nil
}

// ParseKey parses the private key into a key signer and sends it to SSHConnector.AddKeyAuth().
func (c *SSHConnector) ParseKey(privateKey []byte) error {
	if privateKey == nil || len(privateKey) < 1 {
		return errors.New("connections.SSHConnector.ParseKey: privateKey was empty")
	}

	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	return c.AddKeyAuth(key)
}

// ParseKeyWithPassphrase parses a passhphrase protected private key into a key signer
// and sends it to SSHConnector.SetKey().
func (c *SSHConnector) ParseKeyWithPassphrase(privateKey []byte, passphrase string) error {
	if privateKey == nil || len(privateKey) < 1 {
		return errors.New("connections.SSHConnector.ParseKeyWithPassphrase: privateKey was empty")
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
	if err != nil {
		return err
	}

	return c.AddKeyAuth(key)
}

// OpenSession creates a new single command session.
func (c *SSHConnector) OpenSession(server Server) error {
	// log.Print(" - Creating session...")
	if !c.isConnected {
		return ErrNotConnected
	}

	sess, err := c.NewSession()
	if err != nil {
		server.Log(time.Now(), err.Error())
		server.PrintResults(time.Now(), "error", err)
		return err
	}

	c.hasSession = true
	c.Session = sess
	// log.Print("done.\n")
	return nil
}

// CloseSession closes an open session.
func (c *SSHConnector) CloseSession() error {
	// If hasSession is false and there's not Session ref then we have nothing to do.
	if !c.hasSession && c.Session == nil {
		return nil
	}

	c.hasSession = false
	if c.Session == nil {
		return errors.New("connections.SSHConnector.CloseSession: no session avaiable")
	}

	return c.Session.Close()
}

// foundExpect returns true if expect matches anywhere in the byte array.
func foundExpect(data []byte, expect string) bool {
	matched, err := regexp.MatchString(expect, string(data))
	if err != nil {
		log.Printf("connections.SSHConnector.foundExpect: %s", err)
	}

	return matched
	// m := bytes.Index(data, []byte(expect))
	// return m > -1
}

//						//
//	Connector Interface Implementation	//
//						//

func (c *SSHConnector) IsConnected() bool  { return c.isConnected }
func (c *SSHConnector) IsActive() bool     { return c.hasSession }
func (c *SSHConnector) Protocol() Protocol { return SSHProtocol }
func (c *SSHConnector) GetUser() string    { return c.User }
func (c *SSHConnector) DefaultPort() int   { return SSHDefaultPort }
func (c *SSHConnector) IsEmpty() bool      { return c.User == "" }
func (c *SSHConnector) IsValid() bool      { err := c.Validate(); return err == nil }

func (c SSHConnector) Validate() error {
	if c.User == "" {
		return ErrInvalidEmtpyUser
	}

	if len(c.Auth) < 1 {
		return ErrInvalidNoAuthMethod
	}

	return nil
}

func (c *SSHConnector) Open(server Server) error {
	if err := server.Validate(); err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // INCOMPLETE: Find a way to add host key checking.
		Auth:            c.Auth,
		// TODO: Add a timeout to the connection.
	}

	// log.Print("Dialing server...")
	client, err := ssh.Dial("tcp", server.GetAddr(), config)
	if err != nil {
		server.Log(time.Now(), err.Error())
		server.PrintResults(time.Now(), "error", err)
		return err
	}

	c.isConnected = true
	c.Client = client
	// log.Print("done.")
	// INCOMPLETE: Add a keepalive later. Make sure keepalive is cancelled when the connection is closed.
	return nil
}

func (c *SSHConnector) TestConnection(server Server) error {
	expect := "cuttle ok"
	return c.run(server, fmt.Sprintf("echo '%s'", expect), expect)
}

func (c *SSHConnector) Run(server Server, cmd string, exp string) error {
	return c.run(server, cmd, exp)
}

func (c *SSHConnector) run(server Server, cmd string, exp string) error {
	if cmd == "" {
		return ErrEmtpyCmd
	}

	if exp == "" {
		return ErrEmtpyExp
	}

	err := c.OpenSession(server)
	if err != nil {
		return err
	}

	// We have to close the session each time or it will block further command execution.
	defer c.CloseSession()

	// Set ssh.Session.Stdout so we capture the output
	var b bytes.Buffer
	c.Session.Stdout = &b
	eventTime := time.Now()

	// log.Print("   - Running cmd...")
	err = c.Session.Run(cmd)
	if err != nil {
		server.Log(eventTime, err.Error())
		server.PrintResults(eventTime, "error", err)
		return err
	}
	// log.Print("done.")

	// Log the full output of the command
	server.Log(eventTime, b.String())

	// Match results to the expected results and print
	ok := foundExpect(b.Bytes(), exp)
	if !ok {
		server.PrintResults(eventTime, "failed", nil)
		return nil
	}

	server.PrintResults(eventTime, "ok", nil)
	return nil
}

func (c *SSHConnector) Close(force bool) error {
	if c.hasSession {
		// If we don't want to foce close the connection return an error.
		if !force {
			return ErrSessionActive
		}

		// Otherwise force the session closed.
		c.CloseSession()
	}

	c.isConnected = false
	return c.Client.Close()
}
