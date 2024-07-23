package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/chadeldridge/cuttle/core"
)

func Start(ctx context.Context, httpServer *http.Server, logger *core.Logger, timeoutSec int) error {
	// Start the server.
	srvErr := make(chan error)
	go func() {
		logger.Printf("http server listening on %s\n", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				logger.Printf("server closed")
				close(srvErr)
			} else {
				// logger.Printf("http server error: %v\n", err)
				srvErr <- err
			}
		}
	}()

	// Create a wait group to handle a graceful shutdown.
	var wg sync.WaitGroup
	wg.Add(1)
	wgErr := make(chan error)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(
			shutdownCtx,
			time.Duration(timeoutSec)*time.Second,
		)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			// logger.Printf("http server shutdown error: %v\n", err)
			wgErr <- fmt.Errorf("http server shutdown error: %w", err)
		}
	}()
	wg.Wait()

	select {
	case err := <-srvErr:
		if err != nil {
			return err
		}
	case err := <-wgErr:
		if err != nil {
			return err
		}
	}

	return nil
}
