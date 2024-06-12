package connections

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func testNewMockConnector() MockConnector {
	return MockConnector{user: testUser}
}

func TestMockConnectorNewMockConnector(t *testing.T) {
	require := require.New(t)
	t.Run("new MockConnector", func(t *testing.T) {
		conn, err := NewMockConnector(testUser)
		require.NoError(err, "conn error when creating MockConnector")
		require.Equal(testUser, conn.user, "user did not match username after MockConnector creation")
	})

	t.Run("empty username", func(t *testing.T) {
		conn, err := NewMockConnector("")
		require.Error(err, "did not get error when creating MockConnector with empty username")
		require.Equal("", conn.user, "MockConnector.user was not empty after empty username given")
	})
}

func TestMockConnectorSetUser(t *testing.T) {
	require := require.New(t)
	newUser := "george"
	conn := testNewMockConnector()

	t.Run("good username", func(t *testing.T) {
		err := conn.SetUser(newUser)
		require.NoError(err, "conn error when creating MockConnector")
		require.Equal(newUser, conn.user, "user did not match username after SetUser")
	})

	t.Run("empty username", func(t *testing.T) {
		err := conn.SetUser("")
		require.Error(err, "did not get error when passing SetUser an empty username")
		require.Equal(newUser, conn.user, "MockConnector.user changed after empty username given")
	})
}

func TestMockConnectorErrOnConnectionOpen(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewMockConnector()
	server := Server{Hostname: testHost, Results: &res, Logs: &log}

	t.Run("default", func(t *testing.T) {
		err := conn.Open(server)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})

	t.Run("true", func(t *testing.T) {
		conn.ErrOnConnectionOpen(true)
		err := conn.Open(server)
		require.Error(err, "MockConnector.Open() did not return an error")
	})

	t.Run("false", func(t *testing.T) {
		conn.isConnected = false
		conn.ErrOnConnectionOpen(false)
		err := conn.Open(server)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})
}

func TestMockConnectorErrOnConnectionClose(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()

	t.Run("default", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Close(false)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})

	t.Run("true", func(t *testing.T) {
		conn.isConnected = true
		conn.ErrOnConnectionClose(true)
		err := conn.Close(false)
		require.Error(err, "MockConnector.Open() did not return an error")
	})

	t.Run("false", func(t *testing.T) {
		conn.isConnected = true
		conn.ErrOnConnectionClose(false)
		err := conn.Close(false)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})
}

func TestMockConnectorErrOnSessionOpen(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewMockConnector()
	server := Server{Hostname: testHost, Results: &res, Logs: &log}

	conn.isConnected = true

	t.Run("default", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
		conn.hasSession = false
	})

	t.Run("true", func(t *testing.T) {
		conn.ErrOnSessionOpen(true)
		err := conn.OpenSession(server)
		require.Error(err, "MockConnector.Open() did not return an error")
	})

	t.Run("false", func(t *testing.T) {
		conn.ErrOnSessionOpen(false)
		err := conn.OpenSession(server)
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
		conn.hasSession = false
	})
}

func TestMockConnectorErrOnSessionClose(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()
	conn.isConnected = true

	t.Run("default", func(t *testing.T) {
		conn.hasSession = true
		err := conn.CloseSession()
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})

	t.Run("true", func(t *testing.T) {
		conn.hasSession = true
		conn.ErrOnSessionClose(true)
		err := conn.CloseSession()
		require.Error(err, "MockConnector.Open() did not return an error")
	})

	t.Run("false", func(t *testing.T) {
		conn.hasSession = true
		conn.ErrOnSessionClose(false)
		err := conn.CloseSession()
		require.NoError(err, "MockConnector.Open() returned an error: %s", err)
	})
}

func TestMockConnectorOpenSession(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewMockConnector()
	server := Server{Hostname: testHost, Results: &res, Logs: &log}

	t.Run("not connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.Error(err, "did not get error with no connection")
		require.False(conn.hasSession, "MockConnector.hasSession was true")
	})

	t.Run("open error", func(t *testing.T) {
		conn.isConnected = true
		conn.sessOpenErr = true

		err := conn.OpenSession(server)
		require.Error(err, "did not get error with no connection")
		require.False(conn.hasSession, "MockConnector.hasSession was true")

		conn.sessOpenErr = false
	})

	t.Run("open success", func(t *testing.T) {
		conn.isConnected = true
		err := conn.OpenSession(server)
		require.NoError(err, "MockConnector.OpenSession() returned an error")
		require.True(conn.isConnected, "MockConnector.isConnected was false")
		conn.hasSession = false
	})
}

func TestMockConnectorCloseSession(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()

	t.Run("no session", func(t *testing.T) {
		conn.isConnected = true
		conn.hasSession = false
		err := conn.CloseSession()
		require.NoError(err, "Connection.CloseSession() returned an error: %s", err)
	})

	t.Run("err on close", func(t *testing.T) {
		conn.isConnected = true
		conn.hasSession = true
		conn.sessCloseErr = true

		err := conn.CloseSession()
		require.Error(err, "Connection.CloseSession() did not return an error")
		require.True(conn.hasSession, "MockConnector.hasSession was false")

		conn.sessCloseErr = false
	})

	t.Run("good close", func(t *testing.T) {
		conn.isConnected = true
		conn.hasSession = true
		err := conn.CloseSession()
		require.NoError(err, "Connection.CloseSession() returned an error: %s", err)
		require.False(conn.hasSession, "MockConnector.hasSession was true")
	})
}

//					//
//	Test Connector Interface	//
//					//

func TestMockConnectorIsConnected(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()

	t.Run("default", func(t *testing.T) {
		ok := conn.IsConnected()
		require.False(ok, "connector was connected")
	})

	t.Run("true", func(t *testing.T) {
		conn.isConnected = true
		ok := conn.IsConnected()
		require.True(ok, "connector was not connected")
	})

	t.Run("false", func(t *testing.T) {
		conn.isConnected = false
		ok := conn.IsConnected()
		require.False(ok, "connector was connected")
	})
}

func TestMockConnectorIsActive(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()

	t.Run("default", func(t *testing.T) {
		ok := conn.IsActive()
		require.False(ok, "connector was connected")
	})

	t.Run("true", func(t *testing.T) {
		conn.hasSession = true
		ok := conn.IsActive()
		require.True(ok, "connector was not connected")
	})

	t.Run("false", func(t *testing.T) {
		conn.hasSession = false
		ok := conn.IsActive()
		require.False(ok, "connector was connected")
	})
}

func TestMockConnectorProtocol(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()
	require.Equal(MockProtocol, conn.Protocol(), "unexpected Protocol returned for MockConnector")
}

func TestMockConnectorUser(t *testing.T) {
	require := require.New(t)

	t.Run("default", func(t *testing.T) {
		conn := MockConnector{}
		require.Empty(conn.GetUser(), "MockConnector.User() did not return empty username")
	})

	t.Run("user set", func(t *testing.T) {
		conn := testNewMockConnector()
		require.Equal(testUser, conn.GetUser(), "MockConnector.User() did not match username")
	})
}

func TestMockConnectorDefaultPort(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()
	require.Equal(MockDefaultPort, conn.DefaultPort(), "MockConnector.DefaultPort() di not match MockDefaultPort")
}

func TestMockConnectorIsEmpty(t *testing.T) {
	require := require.New(t)

	t.Run("empty", func(t *testing.T) {
		conn := MockConnector{}
		require.True(conn.IsEmpty(), "MockConnector was not empty somehow! Seriously, how?")
	})

	t.Run("not empty", func(t *testing.T) {
		conn := testNewMockConnector()
		require.False(conn.IsEmpty(), "MockConnector was empty")
	})
}

func TestMockConnectorIsValid(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		conn := testNewMockConnector()
		require.True(conn.IsValid(), "MockConnector was not valid but user shoud have been set")
	})

	t.Run("invalid", func(t *testing.T) {
		conn := MockConnector{}
		require.False(conn.IsValid(), "MockConnector was valid somehow! Seriously, how?")
	})
}

func TestMockConnectorValidate(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		conn := testNewMockConnector()
		err := conn.Validate()
		require.NoError(err, "MockConnector.Validate() returned an error: %s", err)
	})

	t.Run("invalid", func(t *testing.T) {
		conn := MockConnector{}
		err := conn.Validate()
		require.Error(err, "MockConnector.Validate() did not return an error")
	})
}

func TestMockConnectorOpen(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		conn := testNewMockConnector()
		err := conn.Open(Server{})
		require.NoError(err, "MockConnector.Open() returned an error")
		require.True(conn.isConnected, "failed to open MockConnector")
	})

	t.Run("invalid", func(t *testing.T) {
		conn := MockConnector{}
		err := conn.Open(Server{})
		require.Error(err, "MockConnector.Open() did not return an error")
		require.False(conn.isConnected, "MockConnector openned despite invalid state")
	})
}

func TestMockConnectorClose(t *testing.T) {
	require := require.New(t)
	conn := testNewMockConnector()

	t.Run("no connection", func(t *testing.T) {
		err := conn.Close(true)
		require.Error(err, "MockConnector.Close() did not returned an error")
		require.False(conn.isConnected, "failed to close MockConnector")
	})

	t.Run("open connection", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Close(true)
		require.NoError(err, "MockConnector.Close() returned an error: %s", err)
		require.False(conn.isConnected, "MockConnector.isConnected was true")
	})
}

func TestMockConnectorRun(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewMockConnector()
	server := Server{Hostname: testHost, Results: &res, Logs: &log}
	cmd := "echo testing"
	exp := "testing"

	server.Connector = &conn

	t.Run("not connected", func(t *testing.T) {
		conn.isConnected = false
		err := conn.Run(server, cmd, exp)
		require.Error(err, "MockConnector.Run() did not return an error")
	})

	t.Run("good connector", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Run(server, cmd, exp)
		require.NoError(err, "MockConnector.Run() returned an error: %s", err)
		require.NotEmpty(res.String(), "results Buffer was empty")
		require.NotEmpty(log.String(), "logs Buffer was empty")
	})

	t.Run("empty cmd", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Run(server, "", exp)
		require.Error(err, "MockConnector.Run() did not return an error")
	})

	t.Run("empty exp", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Run(server, cmd, "")
		require.Error(err, "MockConnector.Run() did not return an error")
	})

	t.Run("bad cmd", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Run(server, "blahIsNotACommand -with args", exp)
		require.Error(err, "MockConnector.Run() did not return an error")
	})

	t.Run("bad exp", func(t *testing.T) {
		conn.isConnected = true
		err := conn.Run(server, cmd, "this won't match")
		require.NoError(err, "MockConnector.Run() did not return an error: %s", err)
	})
}

func TestMockConnectorTestConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := testNewMockConnector()
	server := Server{Hostname: testHost, Results: &res, Logs: &log}

	server.Connector = &conn
	conn.isConnected = true

	err := conn.TestConnection(server)
	require.NoError(err, "MockConnector.TestConnection() returned an error: %s", err)
	require.NotEmpty(res.String(), "res Buffer was empty")
	require.NotEmpty(log.String(), "log Buffer was empty")
}
