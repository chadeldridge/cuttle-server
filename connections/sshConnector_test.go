package connections

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

const (
	testHost = "localhost"
	testUser = "bob"
	testPass = "testUserP@ssw0rd"
)

var (
	keyPass = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABATucg6QV
b74QXyKzG7c6YAAAAAZAAAAAEAAAAzAAAAC3NzaC1lZDI1NTE5AAAAIN8GWe3xMFt/5zSP
xbFK7UlOCB72cCvTec2X1fwAFtYgAAAAoCgV9C/P0QHNfo1edW3BgnBQ1bMOpKVxzUkQ7Q
FIHLIj5vRP4Sv7P6d2u4KnVaCsvIuhVyqductwQskVBSsHPU3HwTPQVZZ0Lu8P3cci7oBc
OiOUXdWp4VAqxTXGkpoTs7Kr/WMavOB2C+/AqgWdOhpICpLxAVk5knuXTK9OvSD34EbC0l
GOO5fZbTGQ1XE1ihvWiIAkUn1XyLaBa3xzOZc=
-----END OPENSSH PRIVATE KEY-----`)
	keyNoPass = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD4p4CaynaubF35hzOcEXg6e/mXM4wlluZBKW9FMg8MegAAAKC8UmL4vFJi
+AAAAAtzc2gtZWQyNTUxOQAAACD4p4CaynaubF35hzOcEXg6e/mXM4wlluZBKW9FMg8Meg
AAAECwBTmJkCxA2UyiNnP5Mh3ampIMnZt+wegxE5jqySmfAvingJrKdq5sXfmHM5wReDp7
+ZczjCWW5kEpb0UyDwx6AAAAGGNlbGRyaWRnZUBDRS1PRkZJQ0UtTUFJTgECAwQF
-----END OPENSSH PRIVATE KEY-----`)
)

func testNewSSHConnector(t *testing.T) SSHConnector {
	conn, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("conn error when creating SSHConnector: %s", err)
	}

	conn.AddPasswordAuth(testPass)
	return conn
}

func TestSSHConnectorNewSSHConnector(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	t.Run("new SSHConnector", func(t *testing.T) {
		require.Nil(t, err, "conn error when creating SSHConnector")
		require.Equal(t, testUser, conn.user, "user did not match username after SSHConnector creation")
	})

	// Test empty username
	conn, err = NewSSHConnector("")
	t.Run("empty username", func(t *testing.T) {
		require.NotNil(t, err, "did not get error when creating SSHConnector with empty username")
		require.Empty(t, conn.user, "SSHConnector.user was not empty after empty username given")
	})
}

func TestSSHConnectorSetUser(t *testing.T) {
	newUser := "george"

	conn := testNewSSHConnector(t)
	t.Run("good username", func(t *testing.T) {
		err := conn.SetUser(newUser)
		require.Nil(t, err, "conn error when creating SSHConnector")
		require.Equal(t, newUser, conn.user, "user did not match username after SetUser")
	})

	// Test empty username
	t.Run("empty username", func(t *testing.T) {
		err := conn.SetUser("")
		require.NotNil(t, err, "did not get error when passing SetUser an empty username")
		require.Equal(t, newUser, conn.user, "SSHConnector.user changed after empty username given")
	})
}

func TestSSHConnectorAddPasswordAuth(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("got error when creating SSHConnector: %s", err)
	}

	conn.AddPasswordAuth(testPass)
	require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
}

func TestSSHConnectorAddKeyAuth(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("got error when creating SSHConnector: %s", err)
	}

	key, err := ssh.ParsePrivateKey(keyNoPass)
	require.Nil(t, err, "conn err parsing private key", err)

	t.Run("good key", func(t *testing.T) {
		err = conn.AddKeyAuth(key)
		require.Nil(t, err, "SSHConnector.AddKeyAuth() returned an error")
		require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
	})

	// Reset auth so we get the right count.
	conn.auth = make([]ssh.AuthMethod, 0)
	t.Run("nil key", func(t *testing.T) {
		err = conn.AddKeyAuth(nil)
		require.NotNil(t, err, "SSHConnector.AddKeyAuth() did not return an error")
		require.Equal(t, 0, len(conn.auth), "AuthMethod found")
	})
}

func TestSSHConnectorParseKey(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("got error when creating SSHConnector: %s", err)
	}

	t.Run("good key", func(t *testing.T) {
		err = conn.ParseKey(keyNoPass)
		require.Nil(t, err, "conn err from ParseKey", err)
		require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
	})

	conn.auth = make([]ssh.AuthMethod, 0)
	t.Run("bad key", func(t *testing.T) {
		err := conn.ParseKey([]byte("not a real key"))
		require.NotNil(t, err, "SSHConnector.ParseKey() did not return an error")
		require.Equal(t, 0, len(conn.auth), "AuthMethod found")
	})

	// Reset auth so we get the right count.
	conn.auth = make([]ssh.AuthMethod, 0)
	t.Run("empty key", func(t *testing.T) {
		err := conn.ParseKey([]byte{})
		require.NotNil(t, err, "SSHConnector.ParseKey() did not return an error")
		require.Equal(t, 0, len(conn.auth), "AuthMethod found")
	})
}

func TestSSHConnectorParseKeyWithPassphrase(t *testing.T) {
	conn, err := NewSSHConnector(testUser)
	if err != nil {
		t.Fatalf("got error when creating SSHConnector: %s", err)
	}

	/*
		raw, err := os.ReadFile(keyFileWithPass)
		require.Nil(t, err, "conn err reading keyFile", err)
	*/

	t.Run("good passphrase", func(t *testing.T) {
		err := conn.ParseKeyWithPassphrase(keyPass, testPass)
		require.Nil(t, err, "SSHConnector.ParseKeyWithPassphrase() returned an error")
		require.Equal(t, 1, len(conn.auth), "missing AuthMethod after AddKeyAuth")
	})

	// Reset auth so we get the right count.
	conn.auth = make([]ssh.AuthMethod, 0)
	t.Run("bad passphrase", func(t *testing.T) {
		err := conn.ParseKeyWithPassphrase(keyPass, "")
		require.NotNil(t, err, "SSHConnector.ParseKeyWithPassphrase() did not return an error")
		require.Equal(t, 0, len(conn.auth), "AuthMethod found")
	})

	// Reset auth so we get the right count.
	conn.auth = make([]ssh.AuthMethod, 0)
	t.Run("empty key", func(t *testing.T) {
		err := conn.ParseKeyWithPassphrase([]byte{}, testPass)
		require.NotNil(t, err, "SSHConnector.ParseKeyWithPassphrase() did not return an error")
		require.Equal(t, 0, len(conn.auth), "AuthMethod found")
	})
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

	t.Run("no session", func(t *testing.T) {
		err = conn.CloseSession()
		require.Nil(t, err, "SSHConnector.CloseSession() returned an error")
		require.False(t, conn.hasSession, "failed to close SSHConnector Session")
	})

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

	t.Run("nil session", func(t *testing.T) {
		conn.hasSession = true
		conn.Session = nil
		err = conn.CloseSession()
		require.NotNil(t, err, "SSHConnector.CloseSession() did not returned an error")
		require.False(t, conn.hasSession, "failed to close SSHConnector Session")
	})

	err = conn.Close(true)
	require.Nil(t, err, "SSHConnector.CloseSession() returned an error", err)
	require.False(t, conn.isConnected, "failed to close SSHConnector")
}

func TestSSHConnectorFoundExpect(t *testing.T) {
	expect := "my test data"
	data := []byte(expect)

	t.Run("matching data", func(t *testing.T) {
		matched := foundExpect(data, expect)
		require.True(t, matched, "SSHConnector.foundExpect() did not match matching data")
	})

	t.Run("non-matching data", func(t *testing.T) {
		matched := foundExpect(data, "does not compute")
		require.False(t, matched, "SSHConnector.foundExpect() matched non-matching data")
	})

	t.Run("non-matching data", func(t *testing.T) {
		matched := foundExpect(data, "{(si|sa|za|ja|to)}")
		require.False(t, matched, "SSHConnector.foundExpect() matched non-matching data")
	})
}

//							//
//	Test Connector Interface Implementation		//
//							//

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
	conn := testNewSSHConnector(t)
	require.Equal(t, SSHProtocol, conn.Protocol(), "unexpected Protocol returned for SSHConnector")
}

func TestSSHConnectorUser(t *testing.T) {
	conn := testNewSSHConnector(t)
	require.Equal(t, fakeUser, conn.User(), "SSHConnector.User() did not match username")
}

func TestSSHConnectorDefaultPort(t *testing.T) {
	conn := testNewSSHConnector(t)
	require.Equal(t, SSHDefaultPort, conn.DefaultPort(), "SSHConnector.DefaultPort() di not match SSHDefaultPort")
}

func TestSSHConnectorIsEmpty(t *testing.T) {
	conn := SSHConnector{}
	t.Run("empty", func(t *testing.T) {
		require.True(t, conn.IsEmpty(), "SSHConnector was not empty somehow! Seriously, how?")
	})

	// Test Not Emtpy
	t.Run("not empty", func(t *testing.T) {
		conn = testNewSSHConnector(t)
		require.False(t, conn.IsEmpty(), "SSHConnector was empty but user should have been set")
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
	conn := testNewSSHConnector(t)
	t.Run("valid", func(t *testing.T) {
		err := conn.Validate()
		require.Nil(t, err, "err while validating valid SSHConnector")
	})

	conn = SSHConnector{}
	t.Run("no user", func(t *testing.T) {
		err := conn.Validate()
		require.NotNil(t, err, "no err while validating invalid SSHConnector")
	})

	conn.SetUser(testUser)
	t.Run("no user", func(t *testing.T) {
		err := conn.Validate()
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

	conn = testNewSSHConnector(t)
	server = Server{hostname: "not.likely", Results: &res, Logs: &log}
	err = server.SetConnector(&conn)
	require.Nil(t, err, "SSHConnector.SetConnector() returned an error")

	t.Run("dial err", func(t *testing.T) {
		err := conn.Open(server)
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

	err = conn.Open(server)
	require.Nil(t, err, "SSHConnector.Open() returned an error")
	require.True(t, conn.isConnected, "failed to open SSHConnector")

	t.Run("has session", func(t *testing.T) {
		conn.hasSession = true
		err := conn.Close(false)
		require.NotNil(t, err, "SSHConnector.Close() did not returned an error")
		require.True(t, conn.isConnected, "SSHConnector.Close() closed a connection with an open session")
	})

	t.Run("has session forced", func(t *testing.T) {
		conn.hasSession = true
		err := conn.Close(true)
		require.Nil(t, err, "SSHConnector.Close() did not returned an error")
		require.False(t, conn.isConnected, "failed to close SSHConnector")
	})

	// Setup for next test.
	err = conn.Open(server)
	require.Nil(t, err, "SSHConnector.Open() returned an error")
	require.True(t, conn.isConnected, "failed to open SSHConnector")

	t.Run("open connection", func(t *testing.T) {
		err = conn.Close(true)
		require.Nil(t, err, "SSHConnector.Close() returned an error")
		require.False(t, conn.isConnected, "failed to close SSHConnector")
	})

	t.Run("closed connection", func(t *testing.T) {
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
