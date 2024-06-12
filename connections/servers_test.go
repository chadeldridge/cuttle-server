package connections

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

type serverInputs struct {
	Name     string
	Hostname string
	IP       string
	Protocol string
	Port     int
}

type serverWants struct {
	Name     string
	Hostname string
	IP       string
	Protocol
	Port int
}

var (
	testServerInputs = map[string]serverInputs{
		"good": {
			Name:     "test.home Test Server",
			Hostname: "test.home",
			IP:       "10.0.0.1",
			Protocol: "ssh",
			Port:     22,
		},
		"bad": {
			Name:     "test.home Test Server",
			Hostname: "89ey*(#@F*)89023r",
			IP:       "192.168.501.105",
			Protocol: "blah",
			Port:     -1,
		},
	}

	testServerWants = map[string]serverWants{
		"good": {
			Name:     "test.home Test Server",
			Hostname: "test.home",
			IP:       "10.0.0.1",
			Protocol: SSH,
			Port:     22,
		},
		"bad": {
			Name:     "", // Change this when exploit validation is added for Name
			Hostname: "",
			IP:       "<nil>",
			Protocol: INVALID,
			Port:     0,
		},
	}
)

func testNewServer(inputName string) Server {
	var res bytes.Buffer
	var log bytes.Buffer

	return Server{
		Name:     testServerInputs[inputName].Name,
		Hostname: testServerInputs[inputName].Hostname,
		Port:     testServerInputs[inputName].Port,
		Results:  &res,
		Logs:     &log,
	}
}

func TestServersNewServer(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)

	t.Run("good hostname", func(t *testing.T) {
		server, err := NewServer(
			testServerInputs["good"].Hostname,
			testServerInputs["good"].Port,
			&res,
			&log,
		)
		require.NoError(err, "NewServer() returned an error: %s", err)
		require.Equal(testServerWants["good"].Hostname, server.Hostname, "hostname did not match")
		require.Equal(testServerWants["good"].Hostname, server.Name, "name did not match")
		require.Equal(testServerWants["good"].Port, server.Port, "port did not match")
	})

	t.Run("bad hostname", func(t *testing.T) {
		server, err := NewServer(
			testServerInputs["bad"].Hostname,
			testServerInputs["good"].Port,
			&res,
			&log,
		)
		require.Error(err, "NewServer() did not return an error")
		require.Equal(testServerWants["bad"].Hostname, server.Hostname)
		require.Equal(testServerWants["bad"].Name, server.Name)
		require.Equal(testServerWants["bad"].Port, server.Port)
	})

	t.Run("bad port", func(t *testing.T) {
		server, err := NewServer(
			testServerInputs["good"].Hostname,
			testServerInputs["bad"].Port,
			&res,
			&log,
		)
		require.Error(err, "NewServer() did not return an error")
		require.Equal(testServerWants["good"].Hostname, server.Hostname)
		require.Equal(testServerWants["good"].Hostname, server.Name)
		require.Equal(testServerWants["bad"].Port, server.Port)
	})
}

func TestServersGetIP(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	require.Equal("<nil>", server.GetIP(), "Server.ip did not match <nil> when Server.ip not set")

	server.IP = net.IPv4(10, 0, 0, 1)
	require.Equal(testServerWants["good"].IP, server.GetIP(), "Server.ip did not match expected ip")
}

func TestServersIsEmpty(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")

	t.Run("not empty", func(t *testing.T) {
		require.False(server.IsEmpty(), "Server.IsEmpty() returned true")
	})

	t.Run("empty", func(t *testing.T) {
		server = Server{}
		require.True(server.IsEmpty(), "Server.IsEmpty() returned false")
	})
}

func TestServersIsValid(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	// Should not be valid because we have not added a Connector yet.
	require.False(server.IsValid(), "Server.IsValid() returned true")

	server.Connector = &MockConnector{user: testUser}
	require.True(server.IsValid(), "Server.IsValid() returned false")
}

func TestServersValidate(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	server.Connector = &MockConnector{user: testUser}

	t.Run("valid", func(t *testing.T) {
		err := server.Validate()
		require.NoError(err, "Server.Validate() returned an error: %s", err)
	})

	t.Run("empty hostname", func(t *testing.T) {
		s := server
		s.Hostname = ""
		require.Error(s.Validate(), "Server.Validate() did not return an error")
	})

	t.Run("nil Results", func(t *testing.T) {
		s := server
		s.Results = nil
		require.Error(s.Validate(), "Server.Validate() did not return an error")
	})

	t.Run("nil Logs", func(t *testing.T) {
		s := server
		s.Logs = nil
		require.Error(s.Validate(), "Server.Validate() did not return an error")
	})

	t.Run("nil connector", func(t *testing.T) {
		s := server
		s.Connector = nil
		require.Error(s.Validate(), "Server.Validate() did not return an error")
	})

	t.Run("empty connector", func(t *testing.T) {
		s := server
		s.Hostname = ""
		require.Error(s.Validate(), "Server.Validate() did not return an error")
	})
}

func TestServersRun(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	conn := MockConnector{user: testUser}
	server.Connector = &conn
	conn.isConnected = true

	exp := "my test message"
	err := server.Run(fmt.Sprintf("echo '%s'", exp), exp)
	require.NoError(err, "Server.Run() returned an error: %s", err)
}

func TestServersTestConnection(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	conn := MockConnector{user: testUser}
	server.Connector = &conn

	t.Run("open error", func(t *testing.T) {
		conn.connOpenErr = true
		err := server.TestConnection()
		require.Error(err, "Server.Run() did not return an error")
		conn.connOpenErr = false
	})

	t.Run("connected", func(t *testing.T) {
		err := server.TestConnection()
		require.NoError(err, "Server.Run() returned an error: %s", err)
	})
}

func TestServersGetAddr(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")

	t.Run("hostname", func(t *testing.T) {
		addr := server.GetAddr()
		require.Equal(fmt.Sprintf("%s:%d", server.Hostname, server.Port), addr,
			"Server.GetAddr() output did not match expected value")
	})

	t.Run("ip", func(t *testing.T) {
		server.IP = net.IPv4(10, 0, 0, 1)
		server.UseIP = true
		addr := server.GetAddr()
		require.Equal(fmt.Sprintf("%s:%d", server.GetIP(), server.Port), addr,
			"Server.GetAddr() output did not match expected value")
	})
}

func TestServersSetUseIP(t *testing.T) {
	require := require.New(t)
	server := testNewServer("good")
	server.IP = net.IPv4(10, 0, 0, 1)

	t.Run("true", func(t *testing.T) {
		server.SetUseIP(true)
		require.True(server.UseIP, "Server.useIP did not return true")
	})

	t.Run("false", func(t *testing.T) {
		server.SetUseIP(false)
		require.False(server.UseIP, "Server.useIP did not return false")
	})
}

func TestServersSetName(t *testing.T) {
	require := require.New(t)

	t.Run("good name", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetName(testServerInputs["good"].Name)
		require.NoError(
			err,
			"Server.SetName(%s) returned an error: %s", testServerInputs["good"].Name, err,
		)
		require.Equal(testServerWants["good"].Name, server.Name, "name did not match")
		// ip should be empty because we are not be setting an IP hostname here.
		require.Equal(
			testServerWants["bad"].IP, server.GetIP(),
			"server.IP() was not <nil> when setting a good hostname",
		)
	})

	// INCOMPLETE: Add test for name validation after it is implemented.
}

func TestServersSetHostname(t *testing.T) {
	require := require.New(t)

	t.Run("good hostname", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetHostname(testServerInputs["good"].Hostname)
		require.NoError(
			err,
			"Server.SetHostname(%s) returned an error: %s", testServerInputs["good"].Hostname, err,
		)
		require.Equal(testServerWants["good"].Hostname, server.Hostname, "hostname did not match")
		// ip should be empty because we are not be setting an IP hostname here.
		require.Equal(
			testServerWants["bad"].IP, server.GetIP(),
			"server.IP() was not <nil> when setting a good hostname",
		)
	})

	t.Run("good ip hostname", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetHostname(testServerInputs["good"].IP)
		require.NoError(err, "error recieved when setting a good IP hostname", err, testServerInputs["good"].IP)
		require.Equal(testServerInputs["good"].IP, server.Hostname, "ip hostname did not match")
		require.Equal(testServerWants["good"].IP, server.GetIP(), "server.IP() did not match expected ip")
	})

	t.Run("empty hostname", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetHostname("")
		require.Error(err, "did not recieve error when setting an empty hostname")
		require.Equal(
			testServerWants["good"].Hostname,
			server.Hostname,
			"hostname was not empty when setting an empty hostname",
		)
		require.Equal(
			testServerWants["bad"].IP,
			server.GetIP(),
			"server.IP() was not <nil> when setting an empty hostname",
		)
	})

	t.Run("bad hostname", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetHostname(testServerInputs["bad"].Hostname)
		require.Error(
			err,
			"Server.SetHostname(%s) did not return an error", testServerInputs["bad"].Hostname,
		)
		require.Equal(testServerWants["good"].Hostname, server.Hostname, "hostname not set for bad IP")
		require.Equal(
			testServerWants["bad"].IP,
			server.GetIP(),
			"server.IP() was not %s when setting a bad hostname", testServerWants["bad"].IP,
		)
	})

	t.Run("bad ip hostname", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetHostname(testServerInputs["bad"].IP)
		require.Error(
			err,
			"Server.SetHostname(%s) did not return an error", testServerInputs["bad"].IP,
		)
		require.Equal(testServerWants["good"].Hostname, server.Hostname, "hostname not set for bad IP")
		require.Equal(
			testServerWants["bad"].IP,
			server.GetIP(),
			"server.IP() was not <nil> when setting a bad IP hostname",
		)
	})
}

func TestServersSetIP(t *testing.T) {
	require := require.New(t)

	t.Run("good ip", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetIP(testServerInputs["good"].IP)
		require.NoError(err, "Server.SetIP(%s) returned an error: %s", testServerInputs["good"].IP, err)
		require.Equal(testServerWants["good"].IP, server.GetIP(), "server.IP() did not match expected ip")
	})

	t.Run("empty server", func(t *testing.T) {
		server := Server{}
		err := server.SetIP(testServerInputs["good"].IP)
		require.NoError(err, "Server.SetIP(%s) returned an error: %s", testServerInputs["good"].IP, err)
		require.Equal(testServerWants["good"].IP, server.GetIP(), "server.IP() did not match expected ip")
		require.Equal(testServerWants["good"].IP, server.Hostname, "server.hostname did not match expected ip")
		require.Equal(testServerWants["good"].IP, server.Name, "server.name did not match expected ip")
	})

	t.Run("empty ip", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetIP("")
		require.Error(err, "Server.SetIP() did not return an error when setting an empty IP")
		require.Equal(
			testServerWants["bad"].IP,
			server.GetIP(),
			"server.IP() was not <nil> when setting a empty IP",
		)
	})

	t.Run("bad ip", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetIP(testServerInputs["bad"].IP)
		require.Error(err, "Server.SetIP(%s) did not return an error", testServerInputs["bad"].IP)
		require.Equal(testServerWants["bad"].IP, server.GetIP(), "server.IP() was not <nil> when setting a bad IP")
	})
}

func TestServersSetPort(t *testing.T) {
	require := require.New(t)
	t.Run("good port", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetPort(testServerInputs["good"].Port)
		require.NoError(err, "Server.SetPort(%s) returned an error: %s", testServerInputs["good"].Port, err)
		require.Equal(testServerWants["good"].Port, server.Port, "server.port did not match expected port")
	})

	t.Run("bad port", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetPort(testServerInputs["bad"].Port)
		require.Error(err, "Server.SetPort(%s) did not return an error", testServerInputs["bad"].Port)
	})
}

func TestServersSetConnector(t *testing.T) {
	require := require.New(t)
	t.Run("full connector", func(t *testing.T) {
		server := testNewServer("good")
		conn := MockConnector{user: testUser}
		err := server.SetConnector(&conn)
		require.NoError(err, "Server.SetConnector() returned an error: %s", err)
		require.Equal(testUser, server.Connector.GetUser(), "Server.Connector.User() did not match expected user")
	})

	t.Run("empty connector", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetConnector(&MockConnector{})
		require.NoError(err, "Server.SetConnector() returned an error: %s", err)
		require.Equal("", server.Connector.GetUser(), "Server.Connector.User() did not match expected user")
	})

	t.Run("nil connector", func(t *testing.T) {
		server := testNewServer("good")
		err := server.SetConnector(nil)
		require.Error(err, "Server.SetConnector() did not return an error: %s", err)
	})
}
