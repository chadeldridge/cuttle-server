package connections

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// iSSH contains the implementation of the Connector interface for SSH Connections.

func (h *SSHConnector) IsConnected() bool  { return h.isConnected }
func (h *SSHConnector) IsActive() bool     { return h.hasSession }
func (h *SSHConnector) Protocol() Protocol { return SSHProtocol }
func (h *SSHConnector) User() string       { return h.user }
func (h *SSHConnector) DefaultPort() int   { return SSHDefaultPort }
func (h *SSHConnector) IsEmpty() bool      { return h.user == "" }
func (h *SSHConnector) IsValid() bool      { return h.user != "" && len(h.auth) > 0 }

func (h *SSHConnector) TestConnection(server Server) error {
	hostLocal, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("connections.SSHHandler.TestConnection: error retrieving local hostname: %s", err)
	}

	expect := fmt.Sprintf("cuttle from %s ok", hostLocal)
	return h.run(server, fmt.Sprintf("echo '%s'", expect), expect)
}

func (h *SSHConnector) Run(server Server, cmd string, exp string) error {
	// Replace command variables before s.run()
	return h.run(server, cmd, exp)
}

func (h *SSHConnector) run(server Server, cmd string, exp string) error {
	err := h.OpenSession(server)
	if err != nil {
		return err
	}

	// We have to close the session each time or it will block further command execution.
	defer h.CloseSession()

	// Set ssh.Session.Stdout so we capture the output
	var b bytes.Buffer
	h.Session.Stdout = &b
	eventTime := time.Now()

	// log.Print("   - Running cmd...")
	err = h.Session.Run(cmd)
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

func (h *SSHConnector) Open(server Server) error {
	config := &ssh.ClientConfig{
		User:            h.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            h.auth,
	}

	// log.Print("Dialing server...")
	client, err := ssh.Dial("tcp", server.GetAddr(), config)
	if err != nil {
		server.Log(time.Now(), err.Error())
		server.PrintResults(time.Now(), "error", err)
		return err
	}

	h.isConnected = true
	h.Client = client
	// log.Print("done.")
	return nil
}

func (h *SSHConnector) Close(force bool) error {
	if h.hasSession {
		// If we don't want to foce close the connection return an error.
		if !force {
			return ErrSessionActive
		}

		// Otherwise force the session closed.
		h.CloseSession()
	}

	h.isConnected = false
	return h.Client.Close()
}
