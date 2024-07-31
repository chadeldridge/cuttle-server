package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chadeldridge/cuttle/core"
)

func addRoutes(mux *http.ServeMux, server *HTTPServer) error {
	root, err := NewRouterGroup(mux, "/v1")
	if err != nil {
		return err
	}

	root.GET("/test", handleTest(server.logger), AuthMiddleware(server.logger))
	root.GET("/login", handleLoginGet(server.logger, server))
	// mux.Handle("/v1/test", handleTest(server.logger))

	return nil
}

// func renderJSON[T any](w http.ResponseWriter, r *http.Request, status int, obj T) error {
func renderJSON[T any](w http.ResponseWriter, status int, obj T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return fmt.Errorf("encoder: %w", err)
	}

	return nil
}

// data, err := readJSON[MyStructType](r)
func readJSON[T any](req *http.Request) (T, error) {
	var obj T
	if err := json.NewDecoder(req.Body).Decode(&obj); err != nil {
		return obj, fmt.Errorf("decoder: %w", err)
	}

	return obj, nil
}

func renderHTML(w http.ResponseWriter, status int, s string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	_, err := w.Write([]byte(s))
	return err
}

func fetchStatic(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func handleTest(logger *core.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("test: %s\n", r.Method)
			err := renderJSON(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
			if err != nil {
				logger.Printf("test: %v\n", err)
			}
		})
}
