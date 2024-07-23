package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRunServer(out io.Writer, args []string, env map[string]string) error {
	// Capture the interrupt signal to gracefully shutdown the server.
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error)

	go func() {
		err := run(ctx, out, args, env)
		ch <- err
	}()
	cancel()

	err := <-ch
	return err
}

func TestServerRun(t *testing.T) {
	var buf bytes.Buffer
	require := require.New(t)

	t.Run("version", func(t *testing.T) {
		ctx := context.Background()
		env := map[string]string{}
		args := []string{"app", "--version"}

		err := run(ctx, &buf, args, env)
		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(buf.String(), "Cuttle v", "run() did not return the version")
		buf.Reset()
	})

	t.Run("help", func(t *testing.T) {
		ctx := context.Background()
		env := map[string]string{}
		args := []string{"app", "--help"}

		err := run(ctx, &buf, args, env)
		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(buf.String(), "Usage:", "run() did not return the version")
		buf.Reset()
	})

	t.Run("missing config file", func(t *testing.T) {
		ctx := context.Background()
		env := map[string]string{}
		args := []string{"app", "--config-file", "/tmp/noconfig.yaml"}

		err := run(ctx, &buf, args, env)
		require.Error(err, "run() did not return an error")
		require.Contains(
			err.Error(),
			"stat /tmp/noconfig.yaml: no such file or directory",
			"run() did not return the correct error",
		)
		buf.Reset()
	})

	t.Run("debug", func(t *testing.T) {
		env := map[string]string{}
		args := []string{"app", "-v"}
		err := testRunServer(&buf, args, env)

		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(buf.String(), "Config: ", "run() did not return the config")
		buf.Reset()
	})

	t.Run("defaults", func(t *testing.T) {
		env := map[string]string{}
		args := []string{"app"}
		err := testRunServer(&buf, args, env)

		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(
			buf.String(),
			"http server listening on 0.0.0.0:8080",
			"run() output did not contain the expected string",
		)
		buf.Reset()
	})

	t.Run("with flags", func(t *testing.T) {
		env := map[string]string{}
		args := []string{"app", "-v", "--host", "127.0.0.1", "--port", "9090"}
		err := testRunServer(&buf, args, env)

		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(
			buf.String(),
			"http server listening on 127.0.0.1:9090",
			"run() output did not contain the expected string",
		)
		buf.Reset()
	})

	t.Run("all interfaces", func(t *testing.T) {
		env := map[string]string{}
		args := []string{"app", "-v", "--host", "0.0.0.0", "--port", "9090"}
		err := testRunServer(&buf, args, env)

		require.NoError(err, "run() returned an error: %s", err)
		require.Contains(
			buf.String(),
			"http server listening on 0.0.0.0:9090",
			"run() output did not contain the expected string",
		)
		buf.Reset()
	})

	/*
		t.Run("invalid host", func(t *testing.T) {
			env := map[string]string{}
			args := []string{"app", "--host", "invalid"}
			err := testRunServer(&buf, args, env)

			log.Printf("buffer: %s", buf.String())
			require.Error(err, "run() did not return an error")
			log.Printf("testing err: %s\n", err)
			buf.Reset()
		})
	*/
}

func TestServerGetEnv(t *testing.T) {
	require := require.New(t)

	err := os.Setenv("CUTTLE_TEST_ENV", "test")
	require.Nil(err, "os.Setenv() returned an error: %s", err)

	env := getEnv()
	require.Greater(len(env), 0, "getEnv() did not return any environment variables")
	require.Equal("test", env["CUTTLE_TEST_ENV"], "getEnv() did not return the correct value")
}
