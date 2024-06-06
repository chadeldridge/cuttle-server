package connections

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var fakeUser = "bob"

func testNewMockConnector(t *testing.T) MockConnector {
	got, err := NewMockConnector(fakeUser)
	t.Run("new MockConnector", func(t *testing.T) {
		require.Nil(t, err, "got error when creating MockConnector")
		require.Equal(t, fakeUser, got.user, "user did not match username after MockConnector creation")
	})
	return got
}

func TestMockConnectorNewMockConnector(t *testing.T) {
	testNewMockConnector(t)

	// Test empty username
	got, err := NewMockConnector("")
	t.Run("empty username", func(t *testing.T) {
		require.NotNil(t, err, "did not get error when creating MockConnector with empty username")
		require.Equal(t, "", got.user, "MockConnector.user was not empty after empty username given")
	})
}

func TestMockConnectorSetUser(t *testing.T) {
	newUser := "george"

	got := testNewMockConnector(t)
	t.Run("good username", func(t *testing.T) {
		err := got.SetUser(newUser)
		require.Nil(t, err, "got error when creating MockConnector")
		require.Equal(t, newUser, got.user, "user did not match username after SetUser")
	})

	// Test empty username
	t.Run("empty username", func(t *testing.T) {
		err := got.SetUser("")
		require.NotNil(t, err, "did not get error when passing SetUser an empty username")
		require.Equal(t, newUser, got.user, "MockConnector.user changed after empty username given")
	})
}

func TestMockConnectorOpenSession(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn := testNewMockConnector(t)

	t.Run("not connected", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.NotNil(t, err, "did not get error with no connection")
		require.False(t, conn.hasSession, "MockConnector.hasSession was true")
	})

	err = conn.Open(server)
	require.Nil(t, err, "Connector.Open() returned an error: ", err)
	defer conn.Close(false)

	t.Run("good session", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.Nil(t, err, "MockConnector.OpenSession() returned an error")
		require.True(t, conn.isConnected, "MockConnector.isConnected was false")
		conn.CloseSession()
	})
}

func TestMockConnectorCloseSession(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn := testNewMockConnector(t)
	conn.Open(server)
	defer conn.Close(false)

	t.Run("err on close", func(t *testing.T) {
		err := conn.OpenSession(server)
		require.Nil(t, err, "Connection.OpenSession() returned an error")
		conn.ErrOnSessionClose(true)

		err = conn.CloseSession()
		require.NotNil(t, err, "Connection.CloseSession() did not return an error")
		require.True(t, conn.hasSession, "MockConnector.hasSession was false")
		conn.ErrOnSessionClose(false)
	})

	t.Run("good close", func(t *testing.T) {
		err = conn.CloseSession()
		require.Nil(t, err, "Connection.CloseSession() returned an error")
		require.False(t, conn.hasSession, "MockConnector.hasSession was true")
	})
}

// Test Connector{}
func TestMockConnectorIsConnected(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("setup connector", func(t *testing.T) {
		err := got.Open(Server{})
		require.Nil(t, err, "MockConnector.Open() returned an error")
		require.True(t, got.isConnected, "failed to open MockConnector")
	})

	t.Run("check for connection", func(t *testing.T) {
		ok := got.IsConnected()
		require.True(t, ok, "connector was not connected")
	})
}

func TestMockConnectorIsActive(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("setup connector", func(t *testing.T) {
		err := got.Open(Server{})
		require.Nil(t, err, "MockConnector.Open() returned an error")
		require.True(t, got.isConnected, "failed to open MockConnector")
	})

	t.Run("has session", func(t *testing.T) {
		got.hasSession = true
		ok := got.IsActive()
		require.True(t, ok, "connector did not have an active session")
	})

	t.Run("no session", func(t *testing.T) {
		got.hasSession = false
		ok := got.IsActive()
		require.False(t, ok, "connector had an active session")
	})
}

func TestMockConnectorProtocol(t *testing.T) {
	got := testNewMockConnector(t)
	require.Equal(t, MockProtocol, got.Protocol(), "unexpected Protocol returned for MockConnector")
}

func TestMockConnectorUser(t *testing.T) {
	got := testNewMockConnector(t)
	require.Equal(t, fakeUser, got.User(), "MockConnector.User() did not match username")
}

func TestMockConnectorDefaultPort(t *testing.T) {
	got := testNewMockConnector(t)
	require.Equal(t, MockDefaultPort, got.DefaultPort(), "MockConnector.DefaultPort() di not match MockDefaultPort")
}

func TestMockConnectorIsEmpty(t *testing.T) {
	got := MockConnector{}
	t.Run("empty", func(t *testing.T) {
		require.True(t, got.IsEmpty(), "MockConnector was not empty somehow! Seriously, how?")
	})

	// Test Not Emtpy
	t.Run("not empty", func(t *testing.T) {
		got = testNewMockConnector(t)
		require.False(t, got.IsEmpty(), "MockConnector was empty but user should have been set")
	})
}

func TestMockConnectorIsValid(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("valid", func(t *testing.T) {
		require.True(t, got.IsValid(), "MockConnector was not valid but user shoud have been set")
	})

	got = MockConnector{}
	t.Run("invalid", func(t *testing.T) {
		require.False(t, got.IsValid(), "MockConnector was valid somehow! Seriously, how?")
	})
}

func TestMockConnectorValidate(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("valid", func(t *testing.T) {
		err := got.Validate()
		require.Nil(t, err, "err while validating valid MockConnector")
	})

	got = MockConnector{}
	t.Run("invalid", func(t *testing.T) {
		err := got.Validate()
		require.NotNil(t, err, "no err while validating invalid MockConnector")
	})
}

func TestMockConnectorOpen(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("valid", func(t *testing.T) {
		err := got.Open(Server{})
		require.Nil(t, err, "MockConnector.Open() returned an error")
		require.True(t, got.isConnected, "failed to open MockConnector")
	})

	got = MockConnector{}
	t.Run("invalid", func(t *testing.T) {
		err := got.Open(Server{})
		require.NotNil(t, err, "MockConnector.Open() did not return an error")
		require.False(t, got.isConnected, "MockConnector openned despite invalid state")
	})
}

func TestMockConnectorClose(t *testing.T) {
	got := testNewMockConnector(t)
	t.Run("close open connection", func(t *testing.T) {
		err := got.Open(Server{})
		require.Nil(t, err, "MockConnector.Open() returned an error")
		require.True(t, got.isConnected, "failed to open MockConnector")

		err = got.Close(true)
		require.Nil(t, err, "MockConnector.Close() returned an error")
		require.False(t, got.isConnected, "failed to close MockConnector")
	})

	t.Run("close closed connection", func(t *testing.T) {
		err := got.Close(true)
		require.NotNil(t, err, "MockConnector.Close() did not returned an error")
		require.False(t, got.isConnected, "failed to close MockConnector")
	})
}

func TestMockConnectorRun(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewMockConnector(t)
	cmd := "echo testing"
	exp := "testing"

	server := Server{hostname: "testing.test", Results: &res, Logs: &log}
	// This also verifies that MockConnector properly implements the Connector interface.
	t.Run("good connector", func(t *testing.T) {
		err := server.SetConnector(&conn)
		require.Nil(t, err, "MockConnector.SetConnector() returned an error")

		err = conn.Open(server)
		require.Nil(t, err, "MockConnector.Open() returned an error")

		err = conn.Run(server, cmd, exp)
		require.Nil(t, err, "MockConnector.Run() returned an error")
		require.NotEmpty(t, res.String(), "results Buffer was empty")
		require.NotEmpty(t, log.String(), "logs Buffer was empty")

		conn.Close(true)
		require.False(t, conn.isConnected, "failed to close MockConnector")
	})

	t.Run("not connected", func(t *testing.T) {
		err := conn.Run(server, cmd, exp)
		require.NotNil(t, err, "MockConnector.Run() did not return an error")
	})

	// Setup for failure tests
	err := conn.Open(server)
	require.Nil(t, err, "MockConnector.Open() returned an error")

	t.Run("empty cmd", func(t *testing.T) {
		err := conn.Run(server, "", exp)
		require.NotNil(t, err, "MockConnector.Run() did not return an error")
	})

	t.Run("empty exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "")
		require.NotNil(t, err, "MockConnector.Run() did not return an error")
	})

	t.Run("bad cmd", func(t *testing.T) {
		err := conn.Run(server, "blahIsNotACommand -with args", exp)
		require.NotNil(t, err, "MockConnector.Run() did not return an error")
	})

	t.Run("bad exp", func(t *testing.T) {
		err := conn.Run(server, cmd, "this won't match")
		require.Nil(t, err, "MockConnector.Run() did not return an error")
	})
}

func TestMockConnectorTestConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewMockConnector(t)

	server := Server{hostname: "testing.test", Results: &res, Logs: &log}
	t.Run("connection setup", func(t *testing.T) {
		err := server.SetConnector(&conn)
		require.Nil(t, err, "MockConnector.SetConnector() returned an error")

		err = conn.Open(server)
		require.Nil(t, err, "MockConnector.Open() returned an error")
	})

	t.Run("test connection", func(t *testing.T) {
		err := conn.TestConnection(server)
		require.Nil(t, err, "MockConnector.TestConnection() returned an error")
		require.NotEmpty(t, res.String(), "results Buffer was empty")
		require.NotEmpty(t, log.String(), "logs Buffer was empty")
	})
}
