package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTCPNewTCPPortOpen(t *testing.T) {
	require := require.New(t)

	t.Run("default timeout", func(t *testing.T) {
		test := NewTCPPortOpen("TCP Test", true, 22)
		require.Equal("TCP Test", test.Name, "NewTCPTest() Name is not 'TCP Test'")
		require.True(test.MustSucceed, "NewTCPTest() MustSucceed is not true")
		require.Equal("port_open", test.Tester.(*TCPTest).testType, "NewTCPTest() testType is not 'port_open'")
		require.Equal("22", test.Tester.(*TCPTest).port, "NewTCPTest() port is not '22'")
		require.Equal(TCPDefaultTimeout, test.Tester.(*TCPTest).timeout, "NewTCPTest() timeout is not 3s")
	})

	t.Run("int timeout", func(t *testing.T) {
		test := NewTCPPortOpen("TCP Test", true, 22, TestArg{Key: "timeout", Value: 4})
		require.Equal("TCP Test", test.Name, "NewTCPTest() Name is not 'TCP Test'")
		require.True(test.MustSucceed, "NewTCPTest() MustSucceed is not true")
		require.Equal("port_open", test.Tester.(*TCPTest).testType, "NewTCPTest() testType is not 'port_open'")
		require.Equal("22", test.Tester.(*TCPTest).port, "NewTCPTest() port is not '22'")
		require.Equal(time.Second*4, test.Tester.(*TCPTest).timeout, "NewTCPTest() timeout is not 3s")
	})

	t.Run("time.Duration timeout", func(t *testing.T) {
		timeout := time.Second * 4
		test := NewTCPPortOpen("TCP Test", true, 22, TestArg{Key: "timeout", Value: timeout})
		require.Equal("TCP Test", test.Name, "NewTCPTest() Name is not 'TCP Test'")
		require.True(test.MustSucceed, "NewTCPTest() MustSucceed is not true")
		require.Equal("port_open", test.Tester.(*TCPTest).testType, "NewTCPTest() testType is not 'port_open'")
		require.Equal("22", test.Tester.(*TCPTest).port, "NewTCPTest() port is not '22'")
		require.Equal(timeout, test.Tester.(*TCPTest).timeout, "NewTCPTest() timeout is not 3s")
	})
}

func TestTCPRun(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)

	t.Run("port_open", func(t *testing.T) {
		test := TCPTest{
			testType: "port_open",
			port:     "22",
			timeout:  TCPDefaultTimeout,
		}

		require.NoError(test.Run(server), "TCPTest.Run() returned an error")
	})

	t.Run("invalid test type", func(t *testing.T) {
		test := TCPTest{
			testType: "invalid",
		}

		err := test.Run(server)
		require.Error(err, "TCPTest.Run() did not return an error")
		require.Equal(ErrInvalidTestType, err, "TCPTest.Run() did not return ErrInvalidTestType")
	})
}

func TestTCPPortOpen(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)

	test := TCPTest{
		testType: "port_open",
		port:     "22",
		timeout:  TCPDefaultTimeout,
	}

	t.Run("success", func(t *testing.T) {
		require.NoError(PortOpen(test, server), "PortOpen() returned an error")
	})

	t.Run("with timeout", func(t *testing.T) {
		test.timeout = time.Second * 3
		err := PortOpen(test, server)
		require.NoError(err, "PortOpen() returned an error: %s", err)
	})

	t.Run("timeout 0", func(t *testing.T) {
		test.timeout = 0
		err := PortOpen(test, server)
		require.NoError(err, "PortOpen() returned an error: %s", err)
	})

	t.Run("invalid port", func(t *testing.T) {
		test.port = "2222"
		err := PortOpen(test, server)
		require.Error(err, "PortOpen() did not return an error")
	})
}
