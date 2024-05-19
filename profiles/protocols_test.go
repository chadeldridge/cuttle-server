package profiles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testStringToProtocol(t *testing.T, in string, want Protocol) {
	got := StringToProtocol(in)
	require.Equal(t, want, got, "recieved wrong Protocol back")
}

func TestProtocolsStringToProtocol(t *testing.T) {
	testStringToProtocol(t, "ssh", SSH)
	testStringToProtocol(t, "SSHpwd", SSHPWD)
	testStringToProtocol(t, "TELNET", TELNET)
	testStringToProtocol(t, "", INVALID)
	testStringToProtocol(t, "testing", INVALID)
}
