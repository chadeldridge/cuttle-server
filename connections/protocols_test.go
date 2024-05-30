package connections

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testStringToProtocol(t *testing.T, in string, want Protocol) {
	got := StringToProtocol(in)
	require.Equal(t, want, got, "connections.testStringToProtocol: returned Protocol did not match")
}

func TestProtocolsStringToProtocol(t *testing.T) {
	testStringToProtocol(t, "ssh", SSH)
	testStringToProtocol(t, "Rdp", RDP)
	testStringToProtocol(t, "TELNET", TELNET)
	testStringToProtocol(t, "", INVALID)
	testStringToProtocol(t, "testing", INVALID)
}
