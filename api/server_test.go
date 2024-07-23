package api

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
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
	config, err := core.NewConfig(map[string]string{}, []string{}, map[string]string{})
	require.NoError(err, "NewConfig() returned an error: %s", err)

	// Update logger with config value.
	logger.DebugMode = config.Debug

	// Setup the HTTP server.
	srv := NewHTTPServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: srv,
	}

	t.Run("Start", func(t *testing.T) {
		// Test the Start function.
		err := testRunServer(ctx, httpServer, logger, 5)
		require.NoError(err, "Start() returned an error: %s", err)
	})
}

func testRunServer(ctx context.Context, httpServer *http.Server, logger *core.Logger, timeout int) error {
	// Capture the interrupt signal to gracefully shutdown the server.
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error)

	go func() {
		err := Start(ctx, httpServer, logger, timeout)
		ch <- err
	}()
	cancel()

	err := <-ch
	return err
}
