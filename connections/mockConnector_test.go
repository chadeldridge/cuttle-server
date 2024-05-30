package connections

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var username = "bob"

func testNewMockHandler(t *testing.T) MockConnector {
	got, err := NewMockHandler(username)
	require.Nil(t, err, "got error when creating MockHandler")
	require.Equal(t, username, got.user, "user did not match username after MockHandler creation")
	return got
}

func TestMockHandlerNewMockHandler(t *testing.T) {
	testNewMockHandler(t)

	// Test empty username
	got, err := NewMockHandler("")
	require.NotNil(t, err, "did not get error when creating MockHandler with empty username")
	require.Equal(t, "", got.user, "MockHandler.user was not empty after empty username given")
}

func TestMockHandlerSetUser(t *testing.T) {
	newUser := "george"

	got := testNewMockHandler(t)
	err := got.SetUser(newUser)
	require.Nil(t, err, "got error when creating MockHandler")
	require.Equal(t, newUser, got.user, "user did not match username after SetUser")

	// Test empty username
	err = got.SetUser("")
	require.NotNil(t, err, "did not get error when passing SetUser an empty username")
	require.Equal(t, newUser, got.user, "MockHandler.user changed after empty username given")
}

// Test Connector{}

func TestMockHandlerProtocol(t *testing.T) {
	got := testNewMockHandler(t)
	require.Equal(t, MockProtocol, got.Protocol(), "unexpected Protocol returned for MockHandler")
}

func TestMockHandlerUser(t *testing.T) {
	got := testNewMockHandler(t)
	require.Equal(t, username, got.User(), "MockHandler.User() did not match username")
}

func TestMockHandlerDefaultPort(t *testing.T) {
	got := testNewMockHandler(t)
	require.Equal(t, MockDefaultPort, got.DefaultPort(), "MockHandler.DefaultPort() di not match MockDefaultPort")
}

func TestMockHandlerIsEmpty(t *testing.T) {
	got := MockConnector{}
	require.True(t, got.IsEmpty(), "MockHandler was not empty somehow! Seriously, how?")

	// Test Not Emtpy
	got = testNewMockHandler(t)
	require.False(t, got.IsEmpty(), "MockHandler was empty but user should have been set")
}

func TestMockHandlerIsValid(t *testing.T) {
	got := MockConnector{}
	require.False(t, got.IsValid(), "MockHandler was valid somehow! Seriously, how?")

	got = testNewMockHandler(t)
	require.True(t, got.IsValid(), "MockHandler was not valid but user shoud have been set")
}

func TestMockHandlerOpen(t *testing.T) {
	got := testNewMockHandler(t)
	err := got.Open(Server{})
	require.Nil(t, err, "MockHandler.Open() returned an error")
	require.True(t, got.isConnected, "failed to open MockHandler")

	// Test invalid MockHandler
	got = MockConnector{}
	err = got.Open(Server{})
	require.NotNil(t, err, "MockHandler.Open() did not return an error")
	require.False(t, got.isConnected, "MockHandler openned despite invalid state")
}

func TestMockHandlerClose(t *testing.T) {
	got := testNewMockHandler(t)
	err := got.Open(Server{})
	require.Nil(t, err, "MockHandler.Open() returned an error")
	require.True(t, got.isConnected, "failed to open MockHandler")

	got.Close(true)
	require.False(t, got.isConnected, "failed to close MockHandler")
}

func TestMockHandlerRun(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewMockHandler(t)
	cmd := "echo testing"
	exp := "testing"

	server := Server{hostname: "testing.test", Results: &res, Logs: &log}
	err := server.SetHandler(&conn)
	require.Nil(t, err, "MockHandler.SetHandler() returned an error")

	err = conn.Open(server)
	require.Nil(t, err, "MockHandler.Open() returned an error")

	err = conn.Run(server, cmd, exp)
	require.Nil(t, err, "MockHandler.Run() returned an error")
	require.NotEmpty(t, res.String(), "results Buffer was empty")
	require.NotEmpty(t, log.String(), "logs Buffer was empty")

	conn.Close(true)
	require.False(t, conn.isConnected, "failed to close MockHandler")

	// Test while not connected
	err = conn.Run(server, cmd, exp)
	require.NotNil(t, err, "MockHandler.Run() did not return an error")

	// Setup for failure tests
	err = conn.Open(server)
	require.Nil(t, err, "MockHandler.Open() returned an error")

	// Test empty cmd
	err = conn.Run(server, "", exp)
	require.NotNil(t, err, "MockHandler.Run() did not return an error")

	// Test empty exp
	err = conn.Run(server, cmd, "")
	require.NotNil(t, err, "MockHandler.Run() did not return an error")
}

func TestMockHandlerTestConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	conn := testNewMockHandler(t)

	server := Server{hostname: "testing.test", Results: &res, Logs: &log}
	err := server.SetHandler(&conn)
	require.Nil(t, err, "MockHandler.SetHandler() returned an error")

	err = conn.Open(server)
	require.Nil(t, err, "MockHandler.Open() returned an error")

	err = conn.TestConnection(server)
	require.Nil(t, err, "MockHandler.TestConnection() returned an error")
	require.NotEmpty(t, res.String(), "results Buffer was empty")
	require.NotEmpty(t, log.String(), "logs Buffer was empty")
}
