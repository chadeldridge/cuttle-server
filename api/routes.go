package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chadeldridge/cuttle/core"
)

func addRoutes(mux *http.ServeMux, server HTTPServer) {
	mux.Handle("/v1/test", handleTest(server.logger))
}

// func encode[T any](w http.ResponseWriter, r *http.Request, status int, obj T) error {
func encode[T any](w http.ResponseWriter, status int, obj T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return fmt.Errorf("encoder: %w", err)
	}

	return nil
}

// data, err := decode[MyStructType](r)
func decode[T any](r *http.Request) (T, error) {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		return obj, fmt.Errorf("decoder: %w", err)
	}

	return obj, nil
}

func handleTest(logger *core.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("test: %s\n", r.Method)
			err := encode(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
			if err != nil {
				logger.Printf("test: %v\n", err)
			}
		})
}
