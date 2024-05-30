package connections

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// iSSH contains the implementation of the Connector interface for SSH Connections.

func (c *SSHConnector) IsConnected() bool  { return c.isConnected }
func (c *SSHConnector) IsActive() bool     { return c.hasSession }
func (c *SSHConnector) Protocol() Protocol { return SSHProtocol }
func (c *SSHConnector) User() string       { return c.user }
func (c *SSHConnector) DefaultPort() int   { return SSHDefaultPort }
func (c *SSHConnector) IsEmpty() bool      { return c.user == "" }
func (c *SSHConnector) IsValid() bool      { return c.user != "" && len(c.auth) > 0 }

func (c *SSHConnector) TestConnection(server Server) error {
	hostLocal, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("connections.SSHHandler.TestConnection: error retrieving local hostname: %s", err)
	}

	expect := fmt.Sprintf("cuttle from %s ok", hostLocal)
	return c.run(server, fmt.Sprintf("echo '%s'", expect), expect)
}

func (c *SSHConnector) Run(server Server, cmd string, exp string) error {
	// Replace command variables before s.run()
	return c.run(server, cmd, exp)
}

func (c *SSHConnector) run(server Server, cmd string, exp string) error {
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

func (c *SSHConnector) Open(server Server) error {
	config := &ssh.ClientConfig{
		User:            c.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            c.auth,
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
	// INCOMPLETE: Add a keepalive later.
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
