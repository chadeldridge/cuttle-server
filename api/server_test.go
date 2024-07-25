package api

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"testing"

	"github.com/chadeldridge/cuttle/core"
	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	require := require.New(t)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var out bytes.Buffer
	logger := core.NewLogger(&out, "cuttle: ", log.LstdFlags, false)

	// Setup the configuration.
	core.SetTester(core.MockTester)
	core.SetReader(core.MockReader)
	core.MockWriteFile("/tmp/cuttle.yaml", core.MockTestConfig, true, nil)
	core.MockWriteFile("/tmp/cuttle_cert.pem", core.MockTestCert, true, nil)
	core.MockWriteFile("/tmp/cuttle_key.pem", core.MockTestKey, true, nil)

	config, err := core.NewConfig(map[string]string{"config_file": "/tmp/cuttle.yaml"}, []string{}, map[string]string{})
	require.NoError(err, "NewConfig() returned an error: %s", err)

	// Update logger with config value.
	logger.DebugMode = config.Debug

	// Setup the HTTP server.
	srv := NewHTTPServer(logger, config)
	err = srv.Build()
	require.NoError(err, "Build() returned an error: %s", err)

	t.Run("Start", func(t *testing.T) {
		// Test the Start function.
		err := testRunServer(ctx, &srv, 5)
		require.NoError(err, "Start() returned an error: %s", err)
	})
}

func testRunServer(ctx context.Context, srv *HTTPServer, timeout int) error {
	// Capture the interrupt signal to gracefully shutdown the server.
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error)

	go func() {
		err := srv.Start(ctx, timeout)
		ch <- err
	}()
	cancel()

	err := <-ch
	return err
}
