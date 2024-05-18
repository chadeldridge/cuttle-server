package connections

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/chadeldridge/cuttle/helpers"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	auth     []ssh.AuthMethod
	IP       net.IP
	Port     string
	user     string
	key      ssh.Signer
	password string
	Results  *helpers.Queue
	Logs     *helpers.Queue
}

// NewSSH
func NewSSH(port int, params map[string]interface{}) (SSH, error) {
	var sshKey []byte
	var passphrase []byte

	s := SSH{}

	if port < 0 || port > 65535 {
		return s, fmt.Errorf("port must be between 0 and 65535, 0 will default to 22")
	}

	if port == 0 {
		port = 22
	}

	s.Port = strconv.Itoa(port)

	if len(params) < 1 {
		return s, nil
	}

	for key, val := range params {
		switch key {
		case "ip":
			err := s.parseIP(val)
			if err != nil {
				return s, err
			}
		case "key":
			k, ok := val.([]byte)
			if !ok {
				return s, fmt.Errorf("expected key of type([]byte), got type(%T)", val)
			}

			sshKey = k
		case "passphrase":
			p, err := parsePassphrase(val)
			if err != nil {
				return s, err
			}

			passphrase = p
		case "username":
			user, ok := val.(string)
			if !ok {
				return s, fmt.Errorf("expected username of type(string), got type(%T)", val)
			}

			// Add some username validation here
			err := s.SetUser(user)
			if err != nil {
				return s, err
			}
		case "password":
			err := s.parsePassword(val)
			if err != nil {
				return s, err
			}
		default:
			return s, fmt.Errorf("unrecognized key in params: %s", key)
		}
	}

	if len(sshKey) > 0 {
		if passphrase == nil {
			err := s.ParseKey(sshKey)
			if err != nil {
				return s, err
			}
		}

		err := s.ParseKeyWithPassphrase(sshKey, passphrase)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func parsePassphrase(passphrase interface{}) ([]byte, error) {
	switch pass := passphrase.(type) {
	case string:
		return []byte(pass), nil
	case []byte:
		return pass, nil
	default:
		return nil, fmt.Errorf("expected passphrase of type([]byte) or type(string), got type(%T)", pass)
	}
}

func (s *SSH) parsePassword(password interface{}) error {
	var p string

	switch pass := password.(type) {
	case string:
		p = pass
	case []byte:
		p = string(pass)
	default:
		return fmt.Errorf("expected password of type([]byte) or type(string), got type(%T)", pass)
	}

	s.SetPassword(p)
	return nil
}

func (s *SSH) parseIP(ip interface{}) error {
	switch ip := ip.(type) {
	case string:
		return s.SetIP(ip)
	case net.IP:
		s.IP = ip
		return nil
	default:
		return fmt.Errorf("expected ip of type(net.IP) or type(string), got type(%T)", ip)
	}
}

func (s *SSH) SetIP(ip string) error {
	addr := net.ParseIP(ip)
	if addr == nil {
		return errors.New("provided string is not a valid IP")
	}

	s.IP = addr
	return nil
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
	s.Logs.Append(fmt.Sprintf("%s@%s:~ %s", s.user, s.IP, txt))
}
