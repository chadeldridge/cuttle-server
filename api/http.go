package api

import (
	"net/http"

	"github.com/chadeldridge/cuttle/core"
)

type HTTPServer struct {
	logger  *core.Logger
	config  *core.Config
	Handler http.Handler
}

func NewHTTPServer(logger *core.Logger, config *core.Config) http.Handler {
	server := HTTPServer{logger: logger, config: config}
	mux := http.NewServeMux()

	// Add routes.
	addRoutes(mux, server)

	server.Handler = mux
	// Add middleware.
	// server.Handler = someMiddleware(server)
	return server.Handler
}
