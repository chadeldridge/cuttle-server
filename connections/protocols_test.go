package connections

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testStringToProtocol(t *testing.T, proto string, want Protocol) {
	got := StringToProtocol(proto)
	testName := proto
	if proto == "" {
		testName = "empty"
	}

	t.Run(fmt.Sprintf("%s protocol", testName), func(t *testing.T) {
		require.Equal(t, want, got, "connections.testStringToProtocol: returned Protocol did not match")
	})
}

func TestProtocolsStringToProtocol(t *testing.T) {
	testStringToProtocol(t, "ssh", SSH)
	testStringToProtocol(t, "Rdp", RDP)
	testStringToProtocol(t, "TELNET", TELNET)
	testStringToProtocol(t, "", INVALID)
	testStringToProtocol(t, "bad", INVALID)
}
