package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/chadeldridge/cuttle-server/core"
	"github.com/chadeldridge/cuttle-server/db"
	"github.com/chadeldridge/cuttle-server/router"
	"github.com/chadeldridge/cuttle-server/services/auth"
	"github.com/chadeldridge/cuttle-server/web/components"
)

type ErrorHandler func(error) templ.Component

func handleError(
	logger *core.Logger,
	w http.ResponseWriter,
	r *http.Request,
	status int,
	pageMsg string,
	errMsg error,
) {
	logger.Printf("%s %s: %s\n", r.Method, r.RequestURI, errMsg)
	err := components.ErrorPage(status, pageMsg, errMsg).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(server *router.HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := components.Page("Cuttle", components.Index()).Render(r.Context(), w)
			if err != nil {
				handleError(server.Logger, w, r, http.StatusInternalServerError, "internal server error", nil)
			}
		})
}

func handleLogin(server *router.HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// handle login
				handleLoginPost(server.Logger, w, r)
			}

			handleLoginGet(server.Logger, w, r)
		})
}

func handleLoginGet(logger *core.Logger, w http.ResponseWriter, r *http.Request) {
	redirect := "/"
	if ref := r.Referer(); ref != "" {
		redirect = ref
	}
	logger.Debugf("login referer: %s\n", redirect)

	err := components.Page("Cuttle Login", components.Login(redirect)).Render(r.Context(), w)
	if err != nil {
		handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
	}
}

func handleLoginPost(logger *core.Logger, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("username")
	p := r.FormValue("password")
	redirect := r.FormValue("redirect")

	// handle login
	if u != "" && p != "" {
		err := components.Login(redirect).Render(r.Context(), w)
		if err != nil {
			handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
		}
	}
}

func handleSignup(server *router.HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// handle signup
				handleSignupPost(server.Logger, server.AuthDB, w, r)
			}

			handleSignupGet(server.Logger, w, r)
		})
}

func handleSignupGet(logger *core.Logger, w http.ResponseWriter, r *http.Request) {
	err := components.Page("Cuttle Signup", components.Signup()).Render(r.Context(), w)
	if err != nil {
		handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
	}
}

func handleSignupPost(logger *core.Logger, authDB db.AuthDB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	n := r.FormValue("name")
	u := r.FormValue("username")
	p := r.FormValue("password")
	c := r.FormValue("confirmPassword")

	// handle signup
	if p == "" || c == "" {
		// Return password required error
		err := router.RenderHTML(w, http.StatusBadRequest, "Signup Error: Password required.")
		if err != nil {
			handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	if p != c {
		// Return password mismatch error
		err := router.RenderHTML(w, http.StatusBadRequest, "Signup Error: Passwords did not match.")
		if err != nil {
			handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	// TODO: Do something with the returned user data
	_, err := auth.Signup(authDB, u, n, p)
	if err != nil {
		handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
