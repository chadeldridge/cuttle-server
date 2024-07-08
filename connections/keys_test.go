package connections

import (
	"encoding/pem"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestKeysStripPasswordFromKey(t *testing.T) {
	require := require.New(t)
	pemBlock, err := ssh.ParseRawPrivateKeyWithPassphrase(keyPass, testPass)
	require.Nil(err, "ssh.ParseRawPrivateKeyWithPassphrase() returned an error: %w", err)
	require.NotNil(pemBlock, "ssh.ParseRawPrivateKeyWithPassphrase() returned nil")

	ePEMBlock, err := ssh.MarshalPrivateKeyWithPassphrase(pemBlock, "", []byte(testPass))
	require.Nil(err, "ssh.MarshalPrivateKeyWithPassphrase() returned an error: %w", err)
	writeable := pem.EncodeToMemory(ePEMBlock)
	log.Println("Encrypted key:\n", string(writeable))
}
