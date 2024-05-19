package connections

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/chadeldridge/cuttle/profiles"
	"golang.org/x/crypto/ssh"
)

const (
	SSHDefaultPort = 22
)

type SSH struct {
	auth     []ssh.AuthMethod
	user     string
	key      ssh.Signer
	password string
	Results  *bytes.Buffer
	Logs     *bytes.Buffer
	profiles.Server
}

// NewSSH creates an SSH struct and sets the Server, Results, and Logs fields. Results and Logs
// can be set to nil if you don't want to ignore them.
func NewSSH(server profiles.Server, results, logs *bytes.Buffer) (SSH, error) {
	s := SSH{Results: results, Logs: logs}
	if server.IsEmpty() {
		return s, errors.New("empty server profile")
	}

	if server.Port() == "0" {
		server.SetPort(SSHDefaultPort)
	}

	s.Server = server
	return s, nil
}

func (s *SSH) SetUser(username string) error {
	// Add username validation here

	s.user = username
	return nil
}

func (s *SSH) SetKey(key ssh.Signer) {
	s.key = key
	s.auth = append(s.auth, ssh.PublicKeys(key))
}

func (s *SSH) ParseKey(privateKey []byte) error {
	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	s.SetKey(key)
	return nil
}

func (s *SSH) ParseKeyWithPassphrase(privateKey, passphrase []byte) error {
	key, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, passphrase)
	if err != nil {
		return err
	}

	s.SetKey(key)
	return nil
}

func (s *SSH) SetPassword(password string) {
	s.password = password
	s.auth = append(s.auth, ssh.Password(password))
}

func (s *SSH) Log(txt string) {
	fmt.Fprintf(s.Logs, "%s@%s:~ %s", s.user, s.Server.Hostname(), txt)
}
