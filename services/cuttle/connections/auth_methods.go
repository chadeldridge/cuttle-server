package connections

import (
	"encoding/pem"
	"fmt"

	"github.com/chadeldridge/cuttle/db"
	"golang.org/x/crypto/ssh"
)

var ErrInvalidAuthType = fmt.Errorf("invalid auth type")

type AuthMethod struct {
	ID       int
	Name     string
	AuthType string
	Proto    Protocol
	Data     []byte
}

func NewAuthMethod(name string) AuthMethod {
	return AuthMethod{Name: name}
}

func ParseAuthMethod(data db.AuthMethodData) (AuthMethod, error) {
	a := AuthMethod{ID: data.ID}
	switch data.AuthType {
	case "ssh_password":
		a.SSHPassword(data.Name, []byte(data.Data))
		return a, nil
	case "ssh_key":
		a.SSHKey(data.Name, []byte(data.Data))
		return a, nil
	default:
		return a, fmt.Errorf("connections.ParseAuthMethod: auth_type not supported: %s", data.AuthType)
	}
}

func (a *AuthMethod) SSHPassword(name string, password []byte) {
	a.AuthType = "ssh_password"
	a.Proto = SSH
	a.Data = password
}

func (a *AuthMethod) SSHKey(name string, key []byte) {
	a.AuthType = "ssh_key"
	a.Proto = SSH
	a.Data = key
}

func (a AuthMethod) ToSSHAuthMethod(passphrase []byte) (ssh.AuthMethod, error) {
	switch a.AuthType {
	case "password":
		// INCOMPLETE: Implement password decryption. Passwords should be encrypted in the database so we have to decrypt them here.
		return ssh.Password(string(a.Data)), nil
	case "ssh_key":
		key, err := ssh.ParsePrivateKeyWithPassphrase(a.Data, passphrase)
		if err != nil {
			return nil, fmt.Errorf("connections.AuthMethod.ToSSHAuthMethod: %w", err)
		}

		return ssh.PublicKeys(key), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

func (a AuthMethod) ToAuthMethodData(passphrase []byte) (db.AuthMethodData, error) {
	var key string

	switch a.AuthType {
	case "password":
		// INCOMPLETE: Implement password encryption. Passwords should be encrypted in the database so we have to encrypt them here.
		key = string(a.Data)
	case "ssh_key":
		data, err := ssh.MarshalPrivateKeyWithPassphrase(a.Data, "", passphrase)
		if err != nil {
			return db.AuthMethodData{}, fmt.Errorf("db.AuthMethods.Create: failed to marshal key: %w", err)
		}
		key = string(pem.EncodeToMemory(data))
	default:
		return db.AuthMethodData{}, ErrInvalidAuthType
	}

	return db.AuthMethodData{
		ID:       a.ID,
		Name:     a.Name,
		AuthType: a.AuthType,
		Data:     key,
	}, nil
}
