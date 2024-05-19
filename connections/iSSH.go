package connections

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// Implementation of Handler interface for SSH Connections

func (s SSH) TestConnection() error {
	hostLocal, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error retrieving local hostname: %s", err)
	}

	expect := fmt.Sprintf("cuttle from %s ok", hostLocal)
	res, err := s.run(fmt.Sprintf("echo '%s'", expect), expect)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.Results, "%s...%s", s.Hostname(), res)
	return nil
}

func (s SSH) Run(cmd string, expect string) (string, error) {
	// Replace command variables before s.run()
	return s.run(cmd, expect)
}

func (s *SSH) run(cmd string, expect string) (string, error) {
	c := &ssh.ClientConfig{
		User:            s.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            s.auth,
	}

	host := s.Server.Hostname()
	if i := s.Server.IP(); i != "<nil>" {
		host = i
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, s.Port()), c)
	if err != nil {
		return "", err
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	var b bytes.Buffer
	sess.Stdout = &b

	err = sess.Run(cmd)
	if err != nil {
		return "", err
	}

	s.Log(b.String())

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
