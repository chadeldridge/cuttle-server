package api

import (
	"github.com/chadeldridge/cuttle/router"
)

func AddRoutes(server *router.HTTPServer) error {
	mwLogger := router.LoggerMiddleware(server.Logger)
	mwAuth := router.AuthMiddleware(server.Logger, server.Users)
	root, err := router.NewRouterGroup(server.Mux, "/api")
	if err != nil {
		return err
	}

	server.Logger.Debug("adding /app/v1 routes")
	v1 := root.Group("/v1", nil)

	// v1.GET("/test", handleTest(server.logger), AuthMiddleware(server.logger))
	v1.GET("/metrics", router.HandleMetrics(server.Logger), mwLogger)
	v1.GET("/test", handleTest(server.Logger), mwLogger, mwAuth)
	// v1.GET("/login", handleLoginGet(server.logger, server), mwLogger)

	return nil
}
