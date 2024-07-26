package api

import (
	"net/http"

	"github.com/chadeldridge/cuttle/core"
)

type Middleware func(http.Handler) http.Handler

// Middleware should call the next handler on success.
// doAuth := authMiddleware(logger *core.Logger, db *db.Users) // returns func(http.Handler) http.Handler
// mux.Handle("/v1/test", doAuth(handleTest(server.logger)))
//
//func newMiddleware(server HTTPServer) Middleware {
//	// Add middleware.
//	// server.Handler = someMiddleware(server)
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(
//			func(w http.ResponseWriter, r *http.Request) {
//				// Do something that fails.
//				//if !something {
//				//	// Return early.
//				//	http.NotFound(w, r)
//				//	return
//				//}
//
//				// Allow original handler to run.
//				next.ServeHTTP(w, r)
//			})
//	}
//}

func AuthMiddleware(logger *core.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "" {
					logger.Println("no auth header")
					err := encode(w, http.StatusUnauthorized, struct {
						Error   string
						Message string
					}{
						Error:   "unauthorized",
						Message: "you need to login",
					})
					if err != nil {
						logger.Printf("AuthMiddlware: response encode failed: %v\n", err)
					}
					return
				}
				// Do something that fails.
				//if !something {
				//	// Return early.
				//	http.NotFound(w, r)
				//	return
				//}

				// Allow original handler to run.
				next.ServeHTTP(w, r)
			})
	}
}
