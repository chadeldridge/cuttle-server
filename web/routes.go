package web

import (
	"net/http"

	"github.com/chadeldridge/cuttle-server/router"
)

var BodyInternalServerError = "internal server error"

func AddRoutes(server *router.HTTPServer) error {
	// Initialize middleware
	mwLogger := router.LoggerMiddleware(server.Logger)
	mwAuth := router.AuthMiddleware(server.Logger, server.CuttleDB)

	// Create a new router group
	root, err := router.NewRouterGroup(server.Mux, "/", mwLogger)
	if err != nil {
		return err
	}

	server.Logger.Debug("adding htmx routes")
	server.Mux.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))),
	)
	root.ANY("/login.html", handleLogin(server))
	root.ANY("/signup.html", handleSignup(server))
	root.GET("/index.html", handleIndex(server), mwAuth)
	// root.GET("/login", handleLoginGet(server.logger, server))

	return nil
}
