package api

/*
import "net/http"

func newMiddleware(server HTTPServer) func(http.Handler) http.Handler {
	// Add middleware.
	// server.Handler = someMiddleware(server)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Do something that fails.
				//if !something {
				//	// Return early.
				//	http.NotFound(w, r)
				//	return
				//}

				// Allow original handler to run.
				h.ServeHTTP(w, r)
			})
	}
}
*/
