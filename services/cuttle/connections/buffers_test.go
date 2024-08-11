package connections

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testTime    = "2021/01/01 12:00:00"
	testTimeObj = time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)
)

func TestBuffersNewBuffers(t *testing.T) {
	var results bytes.Buffer
	var logs bytes.Buffer
	require := require.New(t)

	bufs := NewBuffers("testHost", &results, &logs)

	require.Equal("testHost", bufs.Hostname)
	require.NotNil(bufs.Results)
	require.NotNil(bufs.Logs)
}

func TestBuffersClear(t *testing.T) {
	var results bytes.Buffer
	var logs bytes.Buffer
	require := require.New(t)

	bufs := Buffers{Results: &results, Logs: &logs}
	bufs.Results.WriteString("test")
	bufs.Logs.WriteString("test")

	bufs.Clear()

	require.Equal("", bufs.Results.String())
	require.Equal("", bufs.Logs.String())
}

func TestBuffersClearResults(t *testing.T) {
	var results bytes.Buffer
	var logs bytes.Buffer
	require := require.New(t)

	bufs := Buffers{Results: &results, Logs: &logs}
	bufs.Results.WriteString("test")
	bufs.Logs.WriteString("test")

	bufs.ClearResults()

	require.Equal("", bufs.Results.String())
	require.Equal("test", bufs.Logs.String())
}

func TestBuffersClearLogs(t *testing.T) {
	var results bytes.Buffer
	var logs bytes.Buffer
	require := require.New(t)

	bufs := Buffers{Results: &results, Logs: &logs}
	bufs.Results.WriteString("test")
	bufs.Logs.WriteString("test")

	bufs.ClearLogs()

	require.Equal("test", bufs.Results.String())
	require.Equal("", bufs.Logs.String())
}

func TestBuffersPrintResults(t *testing.T) {
	var results bytes.Buffer
	var logs bytes.Buffer
	ok := "ok"
	fail := "failed"
	require := require.New(t)

	bufs := NewBuffers("testHost", &results, &logs)
	bufs.User = testUser

	t.Run("ok", func(t *testing.T) {
		bufs.PrintResults(testTimeObj, ok, nil)
		require.NotEmpty(results)
		exp := fmt.Sprintf("%s: %s...%s\n", testTime, bufs.Hostname, ok)
		require.Equal(exp, results.String())
	})

	t.Run("failed", func(t *testing.T) {
		bufs.Clear()
		testErr := fmt.Errorf("test error")
		bufs.PrintResults(testTimeObj, fail, testErr)
		require.NotEmpty(results)

		exp := fmt.Sprintf("%s: %s...%s: %s\n", testTime, bufs.Hostname, fail, testErr)
		require.Equal(exp, results.String())
	})
}
