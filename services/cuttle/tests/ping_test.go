package tests

import (
	"bytes"
	"testing"
	"time"

	"github.com/chadeldridge/cuttle/connections"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/stretchr/testify/require"
)

const testHost = "test.home"

func testServerSetup(t *testing.T) connections.Server {
	require := require.New(t)
	var testResults bytes.Buffer
	var testLogs bytes.Buffer

	s, err := connections.NewServer(testHost, 0, &testResults, &testLogs)
	require.NoError(err, "connections.NewServer() returned an error: %s", err)

	return s
}

func TestPingNewPingTest(t *testing.T) {
	require := require.New(t)
	perc := float32(1)

	t.Run("success", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, perc)
		require.Equal("Ping Test", test.Name, "NewPingTest() Name is not 'Ping Test'")
		require.True(test.MustSucceed, "NewPingTest() MustSucceed is not true")
		require.Equal(perc, test.Tester.(*PingTest).successPercent, "NewPingTest() successPercent is not 1")
	})

	t.Run("neg successPerc", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, -1.5)
		require.Equal(float32(0), test.Tester.(*PingTest).successPercent, "NewPingTest() successPerc is not 0")
	})

	t.Run(">1 successPerc", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, 1.5)
		require.Equal(float32(1), test.Tester.(*PingTest).successPercent, "NewPingTest() successPerc is not 1")
	})

	t.Run("set count", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, perc, TestArg{Key: "count", Value: 4})
		require.Equal(4, test.Tester.(*PingTest).count, "NewPingTest() count is not 4")
	})

	t.Run("count zero", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, perc, TestArg{Key: "count", Value: 0})
		require.Equal(PingDefaultCount, test.Tester.(*PingTest).count, "NewPingTest() count is not 1")
	})

	t.Run("set timeout", func(t *testing.T) {
		timeout := time.Second * 4
		test := NewPingTest("Ping Test", true, perc, TestArg{Key: "timeout", Value: timeout})
		require.Equal(timeout, test.Tester.(*PingTest).timeout, "NewPingTest() timeout did not match")
	})

	t.Run("timeout zero", func(t *testing.T) {
		test := NewPingTest("Ping Test", true, perc, TestArg{Key: "timeout", Value: 0})
		require.Equal(PingDefaultTimeout, test.Tester.(*PingTest).timeout, "NewPingTest() timeout did not match")
	})
}

func TestPingRunPinger(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)
	p := PingTest{
		successPercent: 1,
		count:          1,
		timeout:        time.Second * 3,
	}

	pinger, err := probing.NewPinger(server.GetHostAddr())
	require.NoError(err, "probing.NewPinger() returned an error: %s", err)
	pinger.Count = p.count

	t.Run("success", func(t *testing.T) {
		err := p.runPinger(pinger, server.Buffers, false)
		require.NoError(err, "runPinger() returned an error: %s", err)
	})
}

func TestPingRun(t *testing.T) {
	require := require.New(t)
	test := PingTest{
		successPercent: 1,
		count:          1,
		timeout:        time.Second * 3,
	}

	t.Run("quiet", func(t *testing.T) {
		server := testServerSetup(t)
		err := test.Run(server, Quiet())
		require.NoError(err, "Ping() returned an error: %s", err)
		require.Empty(server.Buffers.Logs, "testLogs.Len() is empty")
	})

	t.Run("success", func(t *testing.T) {
		server := testServerSetup(t)
		err := test.Run(server)
		require.NoError(err, "Ping() returned an error: %s", err)
		require.NotEmpty(server.Buffers.Logs, "server.Buffers.Logs.Len() is empty")
	})

	t.Run("invalid server", func(t *testing.T) {
		server := testServerSetup(t)
		server.SetHostname("invalid")
		err := test.Run(server, Quiet())
		require.Error(err, "Ping() did not return an error")
		require.Empty(server.Buffers.Logs, "server.Buffers.Logs is not empty")
	})

	t.Run("fail perc 0", func(t *testing.T) {
		server := testServerSetup(t)
		test.successPercent = 0
		err := test.Run(server)
		require.Error(err, "Ping() did not return an error")
		require.NotEmpty(server.Buffers.Logs, "server.Buffers.Logs is empty")
	})

	// INCOMPLETE: Fix this test. Currently causes pro-bing to panic. Not sure why.
	t.Run("fail", func(t *testing.T) {
		server := testServerSetup(t)
		server.SetHostname("192.168.199.199")
		test.successPercent = 1
		test.timeout = time.Second * 1
		err := test.Run(server)
		require.Error(err, "Ping() did not return an error")
		require.Equal(ErrTestFailed, err, "Ping() did not return ErrTestFailed")
		require.Greater(server.Buffers.Logs.Len(), 1, "server.Buffers.Logs.Len() is empty")
	})
}

func TestPingIntegration(t *testing.T) {
	server := testServerSetup(t)
	test := NewPingTest("Ping Test", true, 1, Quiet())
	err := test.Run(server)

	if test.MustSucceed && err != nil {
		t.Errorf("Ping() returned an error: %s", err)
	}
}
