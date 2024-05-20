package connections

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Implementation of Handler interface for SSH Connections

func (h *SSHHandler) Protocol() Protocol { return SSHProtocol }
func (h *SSHHandler) User() string       { return h.user }
func (h *SSHHandler) DefaultPort() int   { return SSHDefaultPort }
func (h *SSHHandler) IsEmpty() bool      { return h.user == "" }

// IsValid determines if the SSHHandler object is in a useable state. The user and
// at least 1 auth method must be set for the SSHHandler to be considered valid.
func (h *SSHHandler) IsValid() bool { return h.user == "" || len(h.auth) < 1 }

// TestConnection connects to the server and attempts to verify a command can be run.
func (h *SSHHandler) TestConnection(server Server) error {
	hostLocal, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error retrieving local hostname: %s", err)
	}

	expect := fmt.Sprintf("cuttle from %s ok", hostLocal)
	res, err := h.run(server, fmt.Sprintf("echo '%s'", expect), expect)
	if err != nil {
		fmt.Fprintf(server.Results, "%s...failed: %s", server.Hostname(), err)
		return err
	}

	fmt.Fprintf(server.Results, "%s...%s", server.Hostname(), res)
	return nil
}

// Run executes a command against the server and compares the return to the expect string.
func (h *SSHHandler) Run(server Server, cmd string, expect string) (string, error) {
	// Replace command variables before s.run()
	return h.run(server, cmd, expect)
}

func (h *SSHHandler) run(server Server, cmd string, expect string) (string, error) {
	c := &ssh.ClientConfig{
		User:            h.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            h.auth,
	}

	// log.Print("Dialing server...")
	client, err := ssh.Dial("tcp", server.GetAddr(), c)
	if err != nil {
		h.Log(server, err.Error())
		return "", err
	}
	defer client.Close()
	// log.Print("done.")

	// log.Print(" - Creating session...")
	sess, err := client.NewSession()
	if err != nil {
		h.Log(server, err.Error())
		return "", err
	}
	defer sess.Close()
	// log.Print("done.")

	var b bytes.Buffer
	sess.Stdout = &b

	// log.Print("   - Running cmd...")
	err = sess.Run(cmd)
	if err != nil {
		h.Log(server, err.Error())
		return "", err
	}
	// log.Print("done.")

	h.Log(server, b.String())

	ok := foundExpect(b.Bytes(), expect)
	if ok {
		return "ok", nil
	}

	return "failed", nil
}

func foundExpect(data []byte, expect string) bool {
	m := bytes.Index(data, []byte(expect))
	return m > -1
}
