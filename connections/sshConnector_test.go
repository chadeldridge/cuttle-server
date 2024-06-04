package connections

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

var (
	testHost = "localhost"
	testUser = "bob"
	testPass = "testUserP@ssw0rd"

	keyFile         = "../testHelpers/testServer_ed25519_no_pass"
	keyFileWithPass = "../testHelpers/testServer_ed25519_pass"
)

func testNewSSHConnector(t *testing.T) SSHConnector {
	got, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("got error when creating SSHConnector: %s", err)
	}

	got.AddPasswordAuth(testPass)
	return got
}

func TestSSHConnectorNewSSHConnector(t *testing.T) {
	got, err := NewSSHConnector(testUser)
	t.Run("new SSHConnector", func(t *testing.T) {
		require.Nil(t, err, "got error when creating SSHConnector")
		require.Equal(t, testUser, got.user, "user did not match username after SSHConnector creation")
	})

	// Test empty username
	got, err = NewSSHConnector("")
	t.Run("empty username", func(t *testing.T) {
		require.NotNil(t, err, "did not get error when creating SSHConnector with empty username")
		require.Empty(t, got.user, "SSHConnector.user was not empty after empty username given")
	})
}

func TestSSHConnectorSetUser(t *testing.T) {
	newUser := "george"

	got := testNewSSHConnector(t)
	t.Run("good username", func(t *testing.T) {
		err := got.SetUser(newUser)
		require.Nil(t, err, "got error when creating SSHConnector")
		require.Equal(t, newUser, got.user, "user did not match username after SetUser")
	})

	// Test empty username
	t.Run("empty username", func(t *testing.T) {
		err := got.SetUser("")
		require.NotNil(t, err, "did not get error when passing SetUser an empty username")
		require.Equal(t, newUser, got.user, "SSHConnector.user changed after empty username given")
	})
}

func TestSSHConnectorAddPasswordAuth(t *testing.T) {}

func TestSSHConnectorAddKeyAuth(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	t.Run("new SSHConnector", func(t *testing.T) {
		require.Nil(t, err, "got error when creating SSHConnector")
		require.Equal(t, testUser, conn.user, "user did not match username after SSHConnector creation")
	})

	raw, err := os.ReadFile(keyFile)
	require.Nil(t, err, "got err reading keyFile", err)

	key, err := ssh.ParsePrivateKey(raw)
	require.Nil(t, err, "got err parsing private key", err)

	conn.AddKeyAuth(key)
	require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
}

func TestSSHConnectorParseKey(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	t.Run("new SSHConnector", func(t *testing.T) {
		require.Nil(t, err, "got error when creating SSHConnector")
		require.Equal(t, testUser, conn.user, "user did not match username after SSHConnector creation")
	})

	raw, err := os.ReadFile(keyFile)
	require.Nil(t, err, "got err reading keyFile", err)

	err = conn.ParseKey(raw)
	require.Nil(t, err, "got err from ParseKey", err)
	require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
}

func TestSSHConnectorParseKeyWithPassphrase(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	t.Run("new SSHConnector", func(t *testing.T) {
		require.Nil(t, err, "got error when creating SSHConnector")
		require.Equal(t, testUser, conn.user, "user did not match username after SSHConnector creation")
	})

	raw, err := os.ReadFile(keyFileWithPass)
	require.Nil(t, err, "got err reading keyFile", err)

	err = conn.ParseKeyWithPassphrase(raw, testPass)
	require.Nil(t, err, "got err from ParseKeyWithPassphrase", err)
	require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
}

func TestSSHConnectorOpenSession(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)
	server := Server{hostname: testHost, Results: &res, Logs: &log}
	err := server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	err = conn.Open(server)
	require.Nil(t, err, "SSHConnector.Open() returned an error")
	require.True(t, conn.isConnected, "failed to open SSHConnector")

	t.Run("connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.Nil(t, err, "SSHConnector.OpenSession() returned an error", err)
		require.True(t, conn.hasSession, "SSHConnector.hasSession was false")
	})

	err = conn.CloseSession()
	require.Nil(t, err, "SSHConnector.CloseSession() returned an error")
	require.False(t, conn.hasSession, "failed to close SSHConnector Session")

	err = conn.Close(true)
	require.Nil(t, err, "SSHConnector.CloseSession() returned an error", err)
	require.False(t, conn.isConnected, "failed to close SSHConnector")

	conn = SSHConnector{}
	t.Run("not connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.NotNil(t, err, "SSHConnector.OpenSession() did not return an error")
		require.False(t, conn.hasSession, "SSHConnector Session openned despite not being connected")
	})
}

func TestSSHConnectorCloseSession(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	conn := testNewSSHConnector(t)
	server := Server{hostname: testHost, Results: &res, Logs: &logs}
	err := server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	err = conn.Open(server)
	require.Nil(t, err, "SSHConnector.Open() returned an error")
	require.True(t, conn.isConnected, "failed to open SSHConnector")

	err = conn.OpenSession(server)
	require.Nil(t, err, "SSHConnector.OpenSession() returned an error", err)
	require.True(t, conn.hasSession, "SSHConnector.hasSession was false")

	t.Run("open session", func(t *testing.T) {
		err = conn.CloseSession()
		require.Nil(t, err, "SSHConnector.CloseSession() returned an error")
		require.False(t, conn.hasSession, "failed to close SSHConnector Session")
	})

	t.Run("closed session", func(t *testing.T) {
		err = conn.CloseSession()
		require.NotNil(t, err, "SSHConnector.CloseSession() did not returned an error")
		require.False(t, conn.hasSession, "failed to close SSHConnector Session")
	})

	err = conn.Close(true)
	require.Nil(t, err, "SSHConnector.CloseSession() returned an error", err)
	require.False(t, conn.isConnected, "failed to close SSHConnector")
}

// Test Connector{}
func TestSSHConnectorIsConnected(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)
	server := Server{hostname: testHost, Results: &res, Logs: &log}
	err := server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	t.Run("setup connector", func(t *testing.T) {
		err := conn.Open(server)
		require.Nil(t, err, "SSHConnector.Open() returned an error")
		require.True(t, conn.isConnected, "failed to open SSHConnector")
	})

	t.Run("check for connection", func(t *testing.T) {
		ok := conn.IsConnected()
		require.True(t, ok, "connector was not connected")
	})

	conn.Close(true)
	require.False(t, conn.isConnected, "failed to close SSHConnector")

	t.Run("check for disconnected", func(t *testing.T) {
		ok := conn.IsConnected()
		require.False(t, ok, "connector was connected")
	})
}

func TestSSHConnectorIsActive(t *testing.T) {
	conn := testNewSSHConnector(t)

	t.Run("has session", func(t *testing.T) {
		conn.hasSession = true
		ok := conn.IsActive()
		require.True(t, ok, "connector did not have an active session")
	})

	t.Run("no session", func(t *testing.T) {
		conn.hasSession = false
		ok := conn.IsActive()
		require.False(t, ok, "connector had an active session")
	})
}

func TestSSHConnectorProtocol(t *testing.T) {
	got := testNewSSHConnector(t)
	require.Equal(t, SSHProtocol, got.Protocol(), "unexpected Protocol returned for SSHConnector")
}

func TestSSHConnectorUser(t *testing.T) {
	got := testNewSSHConnector(t)
	require.Equal(t, fakeUser, got.User(), "SSHConnector.User() did not match username")
}

func TestSSHConnectorDefaultPort(t *testing.T) {
	got := testNewSSHConnector(t)
	require.Equal(t, SSHDefaultPort, got.DefaultPort(), "SSHConnector.DefaultPort() di not match SSHDefaultPort")
}

func TestSSHConnectorIsEmpty(t *testing.T) {
	got := SSHConnector{}
	t.Run("empty", func(t *testing.T) {
		require.True(t, got.IsEmpty(), "SSHConnector was not empty somehow! Seriously, how?")
	})

	// Test Not Emtpy
	t.Run("not empty", func(t *testing.T) {
		got = testNewSSHConnector(t)
		require.False(t, got.IsEmpty(), "SSHConnector was empty but user should have been set")
	})
}

func TestSSHConnectorIsValid(t *testing.T) {
	conn := testNewSSHConnector(t)

	t.Run("valid", func(t *testing.T) {
		require.True(t, conn.IsValid(), "SSHConnector was not valid but user shoud have been set")
	})

	conn = SSHConnector{}
	t.Run("invalid", func(t *testing.T) {
		require.False(t, conn.IsValid(), "SSHConnector was valid somehow! Seriously, how?")
	})
}

func TestSSHConnectorValidate(t *testing.T) {
	got := testNewSSHConnector(t)
	t.Run("valid", func(t *testing.T) {
		err := got.Validate()
		require.Nil(t, err, "err while validating valid SSHConnector")
	})

	got = SSHConnector{}
	t.Run("invalid", func(t *testing.T) {
		err := got.Validate()
		require.NotNil(t, err, "no err while validating invalid SSHConnector")
	})
}

func TestSSHConnectorOpen(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)
	server := Server{hostname: testHost, Results: &res, Logs: &log}
	err := server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	t.Run("valid", func(t *testing.T) {
		err := conn.Open(server)
		require.Nil(t, err, "SSHConnector.Open() returned an error")
		require.True(t, conn.isConnected, "failed to open SSHConnector")
	})

	conn.Close(true)
	require.False(t, conn.isConnected, "failed to close SSHConnector")

	conn = SSHConnector{}
	t.Run("invalid", func(t *testing.T) {
		err := conn.Open(Server{})
		require.NotNil(t, err, "SSHConnector.Open() did not return an error")
		require.False(t, conn.isConnected, "SSHConnector openned despite invalid state")
	})
}

func TestSSHConnectorClose(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)
	server := Server{hostname: testHost, Results: &res, Logs: &log}
	err := server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	t.Run("close open connection", func(t *testing.T) {
		err := conn.Open(server)
		require.Nil(t, err, "SSHConnector.Open() returned an error")
		require.True(t, conn.isConnected, "failed to open SSHConnector")

		err = conn.Close(true)
		require.Nil(t, err, "SSHConnector.Close() returned an error")
		require.False(t, conn.isConnected, "failed to close SSHConnector")
	})

	t.Run("close closed connection", func(t *testing.T) {
		err := conn.Close(true)
		require.NotNil(t, err, "SSHConnector.Close() did not returned an error")
		require.False(t, conn.isConnected, "failed to close SSHConnector")
	})
}

func TestSSHConnectorRun(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)

	cmd := "echo testing | grep testing"
	exp := "testing"

	server := Server{hostname: testHost, Results: &res, Logs: &log}
	// This also verifies that SSHConnector properly implements the Connector interface.
	t.Run("good connector", func(t *testing.T) {
		err := server.SetConnector(&conn)
		require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

		err = conn.Open(server)
		require.Nil(t, err, "SSHConnector.Open() returned an error")

		err = conn.Run(server, cmd, exp)
		require.Nil(t, err, "SSHConnector.Run() returned an error")
		require.NotEmpty(t, res.String(), "results Buffer was empty")
		require.NotEmpty(t, log.String(), "logs Buffer was empty")

		conn.Close(true)
		require.False(t, conn.isConnected, "failed to close SSHConnector")
	})

	t.Run("not connected", func(t *testing.T) {
		err := conn.Run(server, cmd, exp)
		require.NotNil(t, err, "SSHConnector.Run() did not return an error")
	})

	// Setup for failure tests
	err := conn.Open(server)
	require.Nil(t, err, "SSHConnector.Open() returned an error")

	t.Run("empty cmd", func(t *testing.T) {
		err := conn.Run(server, "", exp)
		require.NotNil(t, err, "SSHConnector.Run() did not return an error")
	})

	t.Run("empty exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "")
		require.NotNil(t, err, "SSHConnector.Run() did not return an error")
	})

	t.Run("bad cmd", func(t *testing.T) {
		err := conn.Run(server, "blahIsNotACommand -with args", exp)
		require.NotNil(t, err, "SSHConnector.Run() did not return an error")
	})

	t.Run("bad exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "this won't match")
		require.Nil(t, err, "SSHConnector.Run() return an error")
		require.Contains(t, getLastLine(server.Results), "failed", "SSHConnector.Run() did not fail match")
	})

	conn.Close(true)
	require.False(t, conn.isConnected, "failed to close SSHConnector")
}

func TestSSHConnectorTestConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewSSHConnector(t)

	server := Server{hostname: testHost, Results: &res, Logs: &log}
	t.Run("connection setup", func(t *testing.T) {
		err := server.SetConnector(&conn)
		require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

		err = conn.Open(server)
		require.Nil(t, err, "SSHConnector.Open() returned an error")
	})

	t.Run("test connection", func(t *testing.T) {
		err := conn.TestConnection(server)
		require.Nil(t, err, "SSHConnector.TestConnection() returned an error")
		require.NotEmpty(t, res.String(), "results Buffer was empty")
		require.NotEmpty(t, log.String(), "logs Buffer was empty")
	})

	conn.Close(true)
	require.False(t, conn.isConnected, "failed to close SSHConnector")
}

func getLastLine(buf *bytes.Buffer) string {
	var b []string
	s := bufio.NewScanner(buf)
	for s.Scan() {
		b = append(b, s.Text())
	}

	return b[len(b)-1]
}
