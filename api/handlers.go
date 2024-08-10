package api

import (
	"net/http"

	"github.com/chadeldridge/cuttle/core"
	"github.com/chadeldridge/cuttle/router"
)

func handleTest(logger *core.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("test: %s\n", r.Method)
			err := router.RenderJSON(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
			if err != nil {
				logger.Printf("test: %v\n", err)
			}
		})
}
