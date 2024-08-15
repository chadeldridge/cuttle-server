package connections

import (
	"bytes"
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
	// AuthMethods []AuthMethod     // A list of AuthMethods to be used for authentication.
	User string // The username to login to the server with.
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
		return fmt.Errorf("connections.SSHConnector.SetName: username was empty")
	}

	// INCOMPLETE: Add name validation here
	c.Name = name
	return nil
}

// SetUser sets the User to be used for connection credentials.
func (c *SSHConnector) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("connections.SSHConnector.SetUser: username was empty")
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
		return fmt.Errorf("connections.SSHConnector.AddKeyAuth: key was nil")
	}

	c.Auth = append(c.Auth, ssh.PublicKeys(key))
	return nil
}

// ParseKey parses the private key into a key signer and sends it to SSHConnector.AddKeyAuth().
func (c *SSHConnector) ParseKey(privateKey []byte) error {
	if privateKey == nil || len(privateKey) < 1 {
		return fmt.Errorf("connections.SSHConnector.ParseKey: privateKey was empty")
	}

	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	return c.AddKeyAuth(key)
}

// ParseKeyWithPassphrase parses a passhphrase protected private key into a key signer
// and sends it to SSHConnector.SetKey().
func (c *SSHConnector) ParseKeyWithPassphrase(privateKey, passphrase []byte) error {
	if privateKey == nil || len(privateKey) < 1 {
		return fmt.Errorf("connections.SSHConnector.ParseKeyWithPassphrase: privateKey was empty")
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
	if err != nil {
		return err
	}

	return c.AddKeyAuth(key)
}

/*
// GenerateAuth creates the ssh.AuthMethod slice from the AuthMethods slice.
func (c *SSHConnector) GenerateAuth() error {
	if len(c.AuthMethods) < 1 {
		return fmt.Errorf("connections.SSHConnector.GenerateSAuth: no auth methods available")
	}

	var errs error
	for _, a := range c.AuthMethods {
		if a.Proto != SSH {
			continue
		}

		am, err := a.ToSSHAuthMethod()
		if err != nil {
			errs = fmt.Errorf("%w\n%s", errs, err)
			continue
		}

		c.Auth = append(c.Auth, am)
	}

	return errs
}
*/

// OpenSession creates a new single command session.
func (c *SSHConnector) OpenSession(bufs Buffers) error {
	// log.Print(" - Creating session...")
	if !c.isConnected {
		return ErrNotConnected
	}

	sess, err := c.NewSession()
	if err != nil {
		bufs.Log(time.Now(), err.Error())
		bufs.PrintResults(time.Now(), "error", err)
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
		return fmt.Errorf("connections.SSHConnector.CloseSession: no session avaiable")
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

// Open creates a connection to the server. addr is the server address to connect to in the format
// of "hostname:port" or "ip:port".
func (c *SSHConnector) Open(addr string, bufs Buffers) error {
	config := &ssh.ClientConfig{
		User:            c.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // INCOMPLETE: Find a way to add host key checking.
		Auth:            c.Auth,
		// TODO: Add a timeout to the connection.
	}

	// log.Print("Dialing server...")
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		bufs.Log(time.Now(), err.Error())
		bufs.PrintResults(time.Now(), "error", err)
		return err
	}

	c.isConnected = true
	c.Client = client
	// log.Print("done.")
	// INCOMPLETE: Add a keepalive later. Make sure keepalive is cancelled when the connection is closed.
	return nil
}

func (c *SSHConnector) TestConnection(bufs Buffers) error {
	expect := "cuttle ok"
	return c.run(bufs, fmt.Sprintf("echo '%s'", expect), expect)
}

func (c *SSHConnector) Run(bufs Buffers, cmd string, exp string) error {
	return c.run(bufs, cmd, exp)
}

func (c *SSHConnector) run(bufs Buffers, cmd string, exp string) error {
	if cmd == "" {
		return ErrEmtpyCmd
	}

	if exp == "" {
		return ErrEmtpyExp
	}

	err := c.OpenSession(bufs)
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
		bufs.Log(eventTime, err.Error())
		bufs.PrintResults(eventTime, "error", err)
		return err
	}
	// log.Print("done.")

	// Log the full output of the command
	bufs.Log(eventTime, b.String())

	// Match results to the expected results and print
	ok := foundExpect(b.Bytes(), exp)
	if !ok {
		bufs.PrintResults(eventTime, "failed", nil)
		return nil
	}

	bufs.PrintResults(eventTime, "ok", nil)
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
