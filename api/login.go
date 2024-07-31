package api

import (
	"net/http"

	"github.com/chadeldridge/cuttle/core"
)

func handleLoginGet(logger *core.Logger, server *HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("login: %s\n", r.Method)
			// renderJSON(w, http.StatusOK, struct{ Message string }{Message: "login"})
			page, err := fetchStatic("static/login.html")
			if err != nil {
				logger.Printf("GET login: %v\n", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			renderHTML(w, http.StatusOK, page)
		})
}

/*
func handleLogin(logger *core.Logger, server *HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("login: %s\n", r.Method)

			data, err := readJSON(r)
			if err != nil {
				logger.Printf("login: %v\n", err)
				writeJSON(w, http.StatusBadRequest, struct{ Message string }{Message: "invalid request"})
				return
			}

			user, err := auth.AuthenticateUser(server.repo, data.Username, data.Password)
			if err != nil {
				logger.Printf("login: %v\n", err)
				writeJSON(w, http.StatusUnauthorized, struct{ Message string }{Message: "invalid credentials"})
				return
			}
		})
}
*/
