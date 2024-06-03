package connections

import (
	"bytes"
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
	results bytes.Buffer
	logs    bytes.Buffer

	goodInputs = serverInputs{
		Name:     "test.home Test Server",
		Hostname: "test.home",
		IP:       "192.168.50.105",
		Protocol: "ssh",
		Port:     22,
	}

	badInputs = serverInputs{
		Name:     "test.home Test Server",
		Hostname: "89ey*(#@F*)89023r",
		IP:       "192.168.501.105",
		Protocol: "blah",
		Port:     -1,
	}

	goodWant = serverWants{
		Name:     goodInputs.Name,
		Hostname: goodInputs.Hostname,
		IP:       "192.168.50.105",
		Protocol: SSH,
		Port:     22,
	}

	badWant = serverWants{
		Name:     "", // Change this when exploit validation is added for Name
		Hostname: "",
		IP:       "<nil>",
		Protocol: INVALID,
		Port:     0,
	}
)

func testNewServer(t *testing.T, input serverInputs) Server {
	got, err := NewServer(input.Hostname, input.Port, &results, &logs)
	t.Run("new server", func(t *testing.T) {
		got.SetName(input.Name)

		require.Nil(t, err, "error recieved for NewServer", err)
		require.Equal(t, goodWant.Hostname, got.Hostname())
		require.Equal(t, goodWant.Name, got.Name())
		require.Equal(t, goodWant.Port, got.Port())
	})

	return got
}

func TestServersNewServer(t *testing.T) {
	// Good Inputs
	testNewServer(t, goodInputs)

	got, err := NewServer(badInputs.Hostname, goodInputs.Port, &results, &logs)
	t.Run("bad hostname", func(t *testing.T) {
		require.NotNil(t, err, "did not receive error when creating Server with bad hostname", err)
		require.Equal(t, badWant.Hostname, got.Hostname())
		require.Equal(t, badWant.Name, got.Name())
		require.Equal(t, badWant.Port, got.Port())
	})

	t.Run("bad port", func(t *testing.T) {
		got, err = NewServer(goodInputs.Hostname, badInputs.Port, &results, &logs)
		require.NotNil(t, err, "did not receive error when creating Server with bad port", err)
		require.Equal(t, goodWant.Hostname, got.Hostname())
		require.Equal(t, goodWant.Hostname, got.Name())
		require.Equal(t, badWant.Port, got.Port())
	})
}

func TestServersSetHostname(t *testing.T) {
	got := testNewServer(t, goodInputs)
	t.Run("good hostname", func(t *testing.T) {
		err := got.SetHostname(goodInputs.Hostname)
		require.Nil(t, err, "error recieved when setting a good hostname", err, goodInputs.Hostname)
		require.Equal(t, goodWant.Hostname, got.Hostname(), "hostname did not match")
		// IP should be empty because we should not be resolving a non-IP hostname here.
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting a good hostname")
	})

	t.Run("good ip hostname", func(t *testing.T) {
		got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
		err := got.SetHostname(goodInputs.IP)
		require.Nil(t, err, "error recieved when setting a good IP hostname", err, goodInputs.IP)
		require.Equal(t, goodInputs.IP, got.Hostname(), "ip hostname did not match")
		require.Equal(t, goodWant.IP, got.IP(), "got.IP() did not match expected ip")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("empty hostname", func(t *testing.T) {
		err := got.SetHostname("")
		require.NotNil(t, err, "did not recieve error when setting an empty hostname")
		require.Equal(t, badWant.Hostname, got.Hostname(), "hostname was not empty when setting an empty hostname")
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting an empty hostname")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("bad hostname", func(t *testing.T) {
		err := got.SetHostname(badInputs.Hostname)
		require.NotNil(t, err, "did not recieve error when setting a bad hostname", badInputs.Hostname)
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting a bad hostname")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("bad ip hostname", func(t *testing.T) {
		err := got.SetHostname(badInputs.IP)
		require.NotNil(t, err, "did not recieve error when setting a bad IP hostname", badInputs.IP)
		require.Equal(t, badWant.Hostname, got.Hostname(), "hostname not set for bad IP")
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting a bad IP hostname")
	})
}

func TestServersSetIP(t *testing.T) {
	got := testNewServer(t, goodInputs)
	t.Run("good ip", func(t *testing.T) {
		err := got.SetIP(goodInputs.IP)
		require.Nil(t, err, "error recieved when setting good IP", goodInputs.IP)
		require.Equal(t, goodWant.IP, got.IP(), "got.IP() did not match expected ip")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("empty ip", func(t *testing.T) {
		err := got.SetIP("")
		require.NotNil(t, err, "did not recieve error when setting an empty IP")
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting a empty IP")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("bad ip", func(t *testing.T) {
		err := got.SetIP(badInputs.IP)
		require.NotNil(t, err, "did not recieve error when setting a bad IP", badInputs.IP)
		require.Equal(t, badWant.IP, got.IP(), "got.IP() was not <nil> when setting a bad IP")
	})
}

func TestServersSetPort(t *testing.T) {
	got := testNewServer(t, goodInputs)
	t.Run("good port", func(t *testing.T) {
		err := got.SetPort(goodInputs.Port)
		require.Nil(t, err, "error recieved when setting good Port", goodInputs.Port)
		require.Equal(t, goodWant.Port, got.Port(), "got.Port() did not match expected port")
	})

	got, _ = NewServer(goodInputs.Name, goodInputs.Port, &results, &logs)
	t.Run("bad port", func(t *testing.T) {
		err := got.SetPort(badInputs.Port)
		require.NotNil(t, err, "did not recieve error when setting a bad Port", badInputs.Port)
		require.Equal(t, badWant.Port, got.Port(), "got.Port() was not 0 when setting a bad port", got.Port())
	})
}
