package cuttle

import (
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
	IP       net.IP
	Protocol
	Port string
}

var (
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
		IP:       net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 50, 105}),
		Protocol: SSH,
		Port:     "22",
	}

	badWant = serverWants{
		Name:     goodInputs.Name, // Change this when exploit validation is added for Name
		Hostname: "",
		IP:       net.IP([]byte{}),
		Protocol: INVALID,
		Port:     "",
	}
)

func testNewServer(t *testing.T, input serverInputs) Server {
	got, err := NewServer(input.Name, input.Protocol)
	require.Nil(t, err, "error recieved for NewServer", err)
	require.Equal(t, goodWant.Name, got.Name())
	require.Equal(t, goodWant.Protocol, got.Protocol())
	return got
}

func TestServersNewServer(t *testing.T) {
	// Good Protol
	testNewServer(t, goodInputs)

	// Bad Protocol
	got, err := NewServer(badInputs.Name, badInputs.Protocol)
	require.NotNil(t, err, "did not recieve error using NewServer with bad Protocol", err, badInputs.Protocol)
	require.Equal(t, badWant.Name, got.Name())
	require.Equal(t, badWant.Protocol, got.Protocol())
}

func TestServersSetHostname(t *testing.T) {
	// Good Hostname
	got := testNewServer(t, goodInputs)
	err := got.SetHostname(goodInputs.Hostname)
	require.Nil(t, err, "error recieved when setting a good hostname", err, goodInputs.Hostname)
	require.Equal(t, goodWant.Hostname, got.Hostname(), "hostname did not match")
	// IP should be empty because we should not be resolving a non-IP hostname here.
	require.Empty(t, got.IP(), "got.IP() was not empty when setting a no-IP hostname", got.IP())

	// Good IP Hostname
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetHostname(goodInputs.IP)
	require.Nil(t, err, "error recieved when setting a good IP hostname", err, goodInputs.IP)
	require.Equal(t, goodInputs.IP, got.Hostname(), "ip hostname did not match")
	require.Equal(t, goodWant.IP, got.IP(), "got.IP() did not match expected ip")

	// Empty Hostname
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetHostname("")
	require.NotNil(t, err, "did not recieve error when setting an empty hostname")
	require.Equal(t, badWant.Hostname, got.Hostname(), "hostname was not empty when setting an empty hostname")
	require.Empty(t, got.IP(), "got.IP() was not empty when setting an empty hostname", got.IP())

	// Bad Hostname
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetHostname(badInputs.Hostname)
	require.NotNil(t, err, "did not recieve error when setting a bad hostname", badInputs.Hostname)
	require.Empty(t, got.IP(), "got.IP() was not empty when setting a no-IP hostname", got.IP())

	// Bad IP Hostname
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetHostname(badInputs.IP)
	require.NotNil(t, err, "did not recieve error when setting a bad IP hostname", badInputs.IP)
	require.Equal(t, badWant.Hostname, got.Hostname(), "hostname not set for bad IP")
	require.Empty(t, got.IP(), "got.IP() was not empty when setting a bad IP hostname", got.IP())
}

func TestServersSetIP(t *testing.T) {
	// Good IP
	got := testNewServer(t, goodInputs)
	err := got.SetIP(goodInputs.IP)
	require.Nil(t, err, "error recieved when setting good IP", goodInputs.IP)
	require.Equal(t, goodWant.IP, got.IP(), "got.IP() did not match expected ip")

	// Empty IP
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetIP("")
	require.NotNil(t, err, "did not recieve error when setting an empty IP")
	require.Empty(t, got.IP(), "got.IP() was not empty when setting an empty hostname", got.IP())

	// Bad IP
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetIP(badInputs.IP)
	require.NotNil(t, err, "did not recieve error when setting a bad IP", badInputs.IP)
	require.Empty(t, got.IP(), "got.IP() was not empty when setting a bad IP hostname", got.IP())
}

func TestServersSetPort(t *testing.T) {
	// Good Port
	got := testNewServer(t, goodInputs)
	err := got.SetPort(goodInputs.Port)
	require.Nil(t, err, "error recieved when setting good Port", goodInputs.Port)
	require.Equal(t, goodWant.Port, got.Port(), "got.Port() did not match expected port")

	// Bad Port
	got, _ = NewServer(goodInputs.Name, goodInputs.Protocol)
	err = got.SetPort(badInputs.Port)
	require.NotNil(t, err, "did not recieve error when setting a bad Port", badInputs.Port)
	require.Empty(t, got.Port(), "got.IP() was not empty when setting a bad IP hostname", got.IP())
}
