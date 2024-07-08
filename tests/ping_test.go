package tests

import (
	"bytes"
	"testing"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/stretchr/testify/require"
)

const testHost = "test.home"

var (
	testResults bytes.Buffer
	testLogs    bytes.Buffer
)

func testServerSetup(t *testing.T) connections.Server {
	require := require.New(t)
	testResults.Reset()
	testLogs.Reset()

	s, err := connections.NewServer(testHost, 0, &testResults, &testLogs)
	require.NoError(err, "connections.NewServer() returned an error: %s", err)

	return s
}

func TestPingPing(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)

	t.Run("success", func(t *testing.T) {
		err := Ping(server, 1)
		require.NoError(err, "Ping() returned an error: %s", err)
		require.NotEmpty(testLogs, "testLogs.Len() is empty")
	})

	t.Run("count zero", func(t *testing.T) {
		testLogs.Reset()
		err := Ping(server, 0)
		require.NoError(err, "Ping() returned an error: %s", err)
		require.Greater(testLogs.Len(), 1, "testLogs.Len() is empty")
	})

	t.Run("invalid server", func(t *testing.T) {
		testLogs.Reset()
		s, err := connections.NewServer("invalid", 0, &testResults, &testLogs)
		require.NoError(err, "connections.NewServer() returned an error: %s", err)
		err = Ping(s, 1)
		require.Error(err, "Ping() did not return an error")
		require.Equal("", testLogs.String(), "testLogs.String() is not empty ~:%s:~", testLogs.String())
	})

	/*
		// Currently causes pro-bing to panic. Not sure why.
		t.Run("fail", func(t *testing.T) {
			testLogs.Reset()
			s, err := connections.NewServer("192.168.199.199", 0, &testResults, &testLogs)
			require.NoError(err, "connections.NewServer() returned an error: %s", err)
			res, err := Ping(s, 1)
			require.NoError(err, "Ping() did not return an error")
			require.False(res, "Ping() did not return false")
			require.Greater(testLogs.Len(), 1, "testLogs.Len() is empty")
		})
	*/
}
