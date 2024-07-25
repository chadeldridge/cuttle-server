package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/chadeldridge/cuttle/core"
)

type HTTPServer struct {
	logger  *core.Logger
	config  *core.Config
	Handler http.Handler
}

func NewHTTPServer(logger *core.Logger, config *core.Config) HTTPServer {
	return HTTPServer{logger: logger, config: config}
}

func (s *HTTPServer) Build() error {
	mux := http.NewServeMux()

	// Add routes.
	addRoutes(mux, s)

	s.Handler = mux
	// Add middleware.
	// server.Handler = someMiddleware(server)
	return nil
}

func (s *HTTPServer) Start(ctx context.Context, timeoutSec int) error {
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(s.config.APIHost, s.config.APIPort),
		Handler: s.Handler,
	}

	// Start the server.
	srvErr := make(chan error)
	go func() {
		s.logger.Printf("http server listening on %s\n", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				s.logger.Printf("server closed")
				close(srvErr)
			} else {
				s.logger.Debugf("http server error: %v\n", err)
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
			s.logger.Debugf("http server shutdown error: %v\n", err)
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
