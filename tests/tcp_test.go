package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTCPPortOpen(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)

	t.Run("success", func(t *testing.T) {
		require.NoError(PortOpen(server, 22), "PortOpen() returned an error")
	})

	t.Run("with timeout", func(t *testing.T) {
		err := PortOpen(server, 22, map[string]any{"timeout": time.Second * 3})
		require.NoError(err, "PortOpen() returned an error: %s", err)
	})

	t.Run("timeout 0", func(t *testing.T) {
		err := PortOpen(server, 22, map[string]any{"timeout": time.Duration(0)})
		require.NoError(err, "PortOpen() returned an error: %s", err)
	})

	t.Run("random arg", func(t *testing.T) {
		err := PortOpen(server, 22, map[string]any{"bla": 1})
		require.NoError(err, "PortOpen() returned an error: %s", err)
	})

	t.Run("invalid port", func(t *testing.T) {
		err := PortOpen(server, 2222)
		require.Error(err, "PortOpen() did not return an error")
	})
}
