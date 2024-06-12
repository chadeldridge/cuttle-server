package connections

import (
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

func testNewSSHConnector() SSHConnector {
	return SSHConnector{User: testUser, Auth: []ssh.AuthMethod{ssh.Password(testPass)}}
}

func TestSSHConnectorNewSSHConnector(t *testing.T) {
	require := require.New(t)

	t.Run("good username", func(t *testing.T) {
		conn, err := NewSSHConnector("my connector", testUser)
		require.NoError(err, "NewSSHConnector() returned an error: %s", err)
		require.Equal("my connector", conn.Name, "SSHConnector.Name did not match")
		require.Equal(testUser, conn.User, "SSHConnector.User did not match")
	})

	t.Run("empty name", func(t *testing.T) {
		conn, err := NewSSHConnector("", testUser)
		require.Error(err, "NewSSHConnector() did not return an error")
		require.Empty(conn.Name, "SSHConnector.Name was not empty")
		require.Empty(conn.User, "SSHConnector.User was not empty")
	})

	t.Run("empty username", func(t *testing.T) {
		conn, err := NewSSHConnector("my connector", "")
		require.Error(err, "NewSSHConnector() did not return an error")
		require.Equal("my connector", conn.Name, "SSHConnector.Name did not match")
		require.Empty(conn.User, "SSHConnector.User was not empty")
	})
}

func TestSSHConnectorSetName(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()
	newName := "george"

	t.Run("good", func(t *testing.T) {
		err := conn.SetName(newName)
		require.NoError(err, "SSHConnector.SetName() returned an error: %s", err)
		require.Equal(newName, conn.Name, "user did not match expected username")
	})

	// Test empty username
	t.Run("empty", func(t *testing.T) {
		err := conn.SetName("")
		require.Error(err, "SSHConnector.SetName() did not return an error")
		require.Equal(newName, conn.Name, "SSHConnector.user changed after empty username given")
	})
}

func TestSSHConnectorSetUser(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()
	newUser := "george"

	t.Run("good", func(t *testing.T) {
		err := conn.SetUser(newUser)
		require.NoError(err, "SSHConnector.SetUser() returned an error: %s", err)
		require.Equal(newUser, conn.User, "user did not match expected username")
	})

	// Test empty username
	t.Run("empty", func(t *testing.T) {
		err := conn.SetUser("")
		require.Error(err, "SSHConnector.SetUser() did not return an error")
		require.Equal(newUser, conn.User, "SSHConnector.user changed after empty username given")
	})
}

func TestSSHConnectorAddPasswordAuth(t *testing.T) {
	require := require.New(t)
	conn := SSHConnector{User: testUser}

	conn.AddPasswordAuth(testPass)
	require.Len(conn.Auth, 1, "missing AuthMethod after AddKeyAuth")
}

func TestSSHConnectorAddKeyAuth(t *testing.T) {
	require := require.New(t)
	conn := SSHConnector{User: testUser}

	key, err := ssh.ParsePrivateKey(keyNoPass)
	require.NoError(err, "ssh.ParsePrivateKey() returned an error: %s", err)

	t.Run("good key", func(t *testing.T) {
		err = conn.AddKeyAuth(key)
		require.NoError(err, "SSHConnector.AddKeyAuth() returned an error: %s", err)
		require.Len(conn.Auth, 1, "missing AuthMethod after AddKeyAuth")
	})

	// Reset auth so we get the right count.
	conn.Auth = []ssh.AuthMethod{}
	t.Run("nil key", func(t *testing.T) {
		err = conn.AddKeyAuth(nil)
		require.Error(err, "SSHConnector.AddKeyAuth() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})
}

func TestSSHConnectorParseKey(t *testing.T) {
	require := require.New(t)
	conn := SSHConnector{User: testUser}

	t.Run("good key", func(t *testing.T) {
		err := conn.ParseKey(keyNoPass)
		require.NoError(err, "SSHConnector.ParseKey() returned an error: %s", err)
		require.Len(conn.Auth, 1, "missing AuthMethod after AddKeyAuth")
	})

	t.Run("bad key", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKey([]byte("not a real key"))
		require.Error(err, "SSHConnector.ParseKey() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})

	t.Run("empty key", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKey([]byte{})
		require.Error(err, "SSHConnector.ParseKey() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})

	t.Run("nil key", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKey(nil)
		require.Error(err, "SSHConnector.ParseKey() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})
}

func TestSSHConnectorParseKeyWithPassphrase(t *testing.T) {
	require := require.New(t)
	conn := SSHConnector{User: testUser}

	/*
		raw, err := os.ReadFile(keyFileWithPass)
		require.NoError(err, "conn err reading keyFile", err)
	*/

	t.Run("good passphrase", func(t *testing.T) {
		err := conn.ParseKeyWithPassphrase(keyPass, testPass)
		require.NoError(err, "SSHConnector.ParseKeyWithPassphrase() returned an error: %s", err)
		require.Len(conn.Auth, 1, "missing AuthMethod after AddKeyAuth")
	})

	t.Run("bad passphrase", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKeyWithPassphrase(keyPass, "")
		require.Error(err, "SSHConnector.ParseKeyWithPassphrase() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})

	t.Run("empty key", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKeyWithPassphrase([]byte{}, testPass)
		require.Error(err, "SSHConnector.ParseKeyWithPassphrase() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})

	t.Run("nil key", func(t *testing.T) {
		conn.Auth = []ssh.AuthMethod{}
		err := conn.ParseKeyWithPassphrase(nil, testPass)
		require.Error(err, "SSHConnector.ParseKeyWithPassphrase() did not return an error")
		require.Len(conn.Auth, 0, "AuthMethod found")
	})
}

func TestSSHConnectorOpenSession(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &log}

	err := conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
	require.True(conn.isConnected, "failed to open SSHConnector")

	t.Run("connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.NoError(err, "SSHConnector.OpenSession() returned an error: %s", err)
		require.True(conn.hasSession, "SSHConnector.hasSession was false")
	})

	err = conn.CloseSession()
	require.NoError(err, "SSHConnector.CloseSession() returned an error")
	require.False(conn.hasSession, "failed to close SSHConnector Session")

	err = conn.Close(true)
	require.NoError(err, "SSHConnector.CloseSession() returned an error", err)
	require.False(conn.isConnected, "failed to close SSHConnector")

	conn = SSHConnector{}
	t.Run("not connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.Error(err, "SSHConnector.OpenSession() did not return an error")
		require.False(conn.hasSession, "SSHConnector Session openned despite not being connected")
	})

	// TODO: Find a way to get ssh.Client.NewSession() to return an error.
}

func TestSSHConnectorCloseSession(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &logs}

	err := conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
	require.True(conn.isConnected, "failed to open SSHConnector")

	t.Run("no session", func(t *testing.T) {
		err = conn.CloseSession()
		require.NoError(err, "SSHConnector.CloseSession() returned an error: %s", err)
		require.False(conn.hasSession, "failed to close SSHConnector Session")
	})

	err = conn.OpenSession(server)
	require.NoError(err, "SSHConnector.OpenSession() returned an error: %s", err)
	require.True(conn.hasSession, "SSHConnector.hasSession was false")

	t.Run("open session", func(t *testing.T) {
		err = conn.CloseSession()
		require.NoError(err, "SSHConnector.CloseSession() returned an error")
		require.False(conn.hasSession, "failed to close SSHConnector Session")
	})

	t.Run("closed session", func(t *testing.T) {
		err = conn.CloseSession()
		require.Error(err, "SSHConnector.CloseSession() did not return an error")
		require.False(conn.hasSession, "failed to close SSHConnector Session")
	})

	t.Run("nil session", func(t *testing.T) {
		conn.hasSession = true
		conn.Session = nil
		err = conn.CloseSession()
		require.Error(err, "SSHConnector.CloseSession() did not return an error")
		require.False(conn.hasSession, "failed to close SSHConnector Session")
	})

	err = conn.Close(true)
	require.NoError(err, "SSHConnector.CloseSession() returned an error", err)
	require.False(conn.isConnected, "failed to close SSHConnector")
}

func TestSSHConnectorFoundExpect(t *testing.T) {
	require := require.New(t)
	expect := "my test data"
	data := []byte(expect)

	t.Run("matching data", func(t *testing.T) {
		matched := foundExpect(data, expect)
		require.True(matched, "SSHConnector.foundExpect() did not match matching data")
	})

	t.Run("non-matching data", func(t *testing.T) {
		matched := foundExpect(data, "does not compute")
		require.False(matched, "SSHConnector.foundExpect() matched non-matching data")
	})

	/*
		// TODO: Find a way to get regexp.MatchString to return an error.
		t.Run("non-matching data", func(t *testing.T) {
			matched := foundExpect(data, "{(si|sa|za|ja|to)}")
			require.False(matched, "SSHConnector.foundExpect() matched non-matching data")
		})
	*/
}

//							//
//	Test Connector Interface Implementation		//
//							//

func TestSSHConnectorIsConnected(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()

	t.Run("connected", func(t *testing.T) {
		conn.isConnected = true
		require.True(conn.IsConnected(), "SSHConnector.IsConnected() returned false")
	})

	t.Run("not connected", func(t *testing.T) {
		conn.isConnected = false
		require.False(conn.IsConnected(), "SSHConnector.IsConnected() returned true")
	})
}

func TestSSHConnectorIsActive(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()

	t.Run("has session", func(t *testing.T) {
		conn.hasSession = true
		require.True(conn.IsActive(), "SSHConnector.IsActive() returned false")
	})

	t.Run("no session", func(t *testing.T) {
		conn.hasSession = false
		require.False(conn.IsActive(), "SSHConnector.IsActive() returned true")
	})
}

func TestSSHConnectorProtocol(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()
	require.Equal(SSHProtocol, conn.Protocol(), "SSHConnector.Protocol() did not match SSHProtocol")
}

func TestSSHConnectorUser(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()
	require.Equal(testUser, conn.GetUser(), "SSHConnector.User() did not match testUser")
}

func TestSSHConnectorDefaultPort(t *testing.T) {
	require := require.New(t)
	conn := testNewSSHConnector()
	require.Equal(SSHDefaultPort, conn.DefaultPort(), "SSHConnector.DefaultPort() did not match SSHDefaultPort")
}

func TestSSHConnectorIsEmpty(t *testing.T) {
	require := require.New(t)

	t.Run("empty", func(t *testing.T) {
		conn := SSHConnector{}
		require.True(conn.IsEmpty(), "SSHConnector was not empty somehow! Seriously, how?")
	})

	// Test Not Emtpy
	t.Run("not empty", func(t *testing.T) {
		conn := testNewSSHConnector()
		require.False(conn.IsEmpty(), "SSHConnector was empty but user should have been set")
	})
}

func TestSSHConnectorIsValid(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		conn := testNewSSHConnector()
		require.True(conn.IsValid(), "SSHConnector was not valid")
	})

	t.Run("invalid", func(t *testing.T) {
		conn := SSHConnector{}
		require.False(conn.IsValid(), "SSHConnector was valid somehow! Seriously, how?")
	})
}

func TestSSHConnectorValidate(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		conn := testNewSSHConnector()
		err := conn.Validate()
		require.NoError(err, "SSHConnector.Validate() returned an error: %s", err)
	})

	t.Run("no user", func(t *testing.T) {
		conn := SSHConnector{}
		require.Error(conn.Validate(), "SSHConnector.Validate() did not return an error")
	})

	t.Run("no auth", func(t *testing.T) {
		conn := SSHConnector{User: testUser}
		require.Error(conn.Validate(), "SSHConnector.Validate() did not return an error")
	})
}

func TestSSHConnectorOpen(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &log}

	t.Run("valid", func(t *testing.T) {
		err := conn.Open(server)
		require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
		require.True(conn.isConnected, "failed to open SSHConnector")
	})

	err := conn.Close(true)
	require.NoError(err, "SSHConnector.Close() returned an error: %s", err)
	require.False(conn.isConnected, "failed to close SSHConnector")

	t.Run("invalid server", func(t *testing.T) {
		require.Error(conn.Open(Server{}), "SSHConnector.Open() did not return an error")
		require.False(conn.isConnected, "SSHConnector openned despite invalid state")
	})

	t.Run("invalid connector", func(t *testing.T) {
		conn := SSHConnector{}
		require.Error(conn.Open(server), "SSHConnector.Open() did not return an error")
		require.False(conn.isConnected, "SSHConnector openned despite invalid state")
	})

	t.Run("dial err", func(t *testing.T) {
		conn := testNewSSHConnector()
		server := Server{Hostname: "not.likely", Connector: &conn, Results: &res, Logs: &log}
		require.Error(conn.Open(server), "SSHConnector.Open() did not return an error")
		require.False(conn.isConnected, "SSHConnector openned despite invalid state")
	})
}

func TestSSHConnectorTestConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &log}

	err := conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
	defer conn.Close(true)

	t.Run("test connection", func(t *testing.T) {
		err := conn.TestConnection(server)
		require.NoError(err, "SSHConnector.TestConnection() returned an error: %s", err)
		require.NotEmpty(res.String(), "results Buffer was empty")
		require.NotEmpty(log.String(), "logs Buffer was empty")
	})
}

func TestSSHConnectorRun(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &log}
	cmd := "echo testing | grep testing"
	exp := "testing"

	err := conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)

	// This also verifies that SSHConnector properly implements the Connector interface.
	t.Run("good connector", func(t *testing.T) {
		err = conn.Run(server, cmd, exp)
		require.NoError(err, "SSHConnector.Run() returned an error: %s", err)
		require.NotEmpty(res.String(), "results Buffer was empty")
		require.NotEmpty(log.String(), "logs Buffer was empty")
	})

	t.Run("empty cmd", func(t *testing.T) {
		err := conn.Run(server, "", exp)
		require.Error(err, "SSHConnector.Run() did not return an error")
	})

	t.Run("empty exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "")
		require.Error(err, "SSHConnector.Run() did not return an error")
	})

	t.Run("bad cmd", func(t *testing.T) {
		err := conn.Run(server, "blahIsNotACommand -with args", exp)
		require.Error(err, "SSHConnector.Run() did not return an error")
	})

	t.Run("bad exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "this won't match")
		require.NoError(err, "SSHConnector.Run() return an error")
		require.Contains(GetLastBufferLine(server.Results), "failed", "SSHConnector.Run() did not fail match")
	})

	conn.Close(true)
	require.False(conn.isConnected, "failed to close SSHConnector")

	t.Run("not connected", func(t *testing.T) {
		err := conn.Run(server, cmd, exp)
		require.Error(err, "SSHConnector.Run() did not return an error")
	})
}

func TestSSHConnectorClose(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewSSHConnector()
	server := Server{Hostname: testHost, Connector: &conn, Results: &res, Logs: &log}

	err := conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
	require.True(conn.isConnected, "failed to open SSHConnector")

	t.Run("has session", func(t *testing.T) {
		conn.hasSession = true
		err := conn.Close(false)
		require.Error(err, "SSHConnector.Close() did not return an error")
		require.True(conn.isConnected, "SSHConnector.Close() closed a connection with an open session")
	})

	t.Run("force close", func(t *testing.T) {
		conn.hasSession = true
		err := conn.Close(true)
		require.NoError(err, "SSHConnector.Close() returned an error: %s", err)
		require.False(conn.isConnected, "failed to close SSHConnector")
	})

	err = conn.Open(server)
	require.NoError(err, "SSHConnector.Open() returned an error: %s", err)
	require.True(conn.isConnected, "failed to open SSHConnector")

	t.Run("open connection", func(t *testing.T) {
		conn.hasSession = false
		err = conn.Close(false)
		require.NoError(err, "SSHConnector.Close() returned an error: %s", err)
		require.False(conn.isConnected, "failed to close SSHConnector")
	})

	t.Run("closed connection", func(t *testing.T) {
		conn.hasSession = false
		err := conn.Close(false)
		require.Error(err, "SSHConnector.Close() did not return an error")
		require.False(conn.isConnected, "failed to close SSHConnector")
	})
}
