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
	logger.Debugf("signup - name: %s, username: %s, password: %s, confirmPassword: %s\n", n, u, p, c)

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

	// TODO: Do something with the returned user?
	_, err := auth.Signup(authDB, u, n, p)
	if err != nil {
		logger.Printf("%s %s: %s\n", r.Method, r.RequestURI, err)
		err := components.LoginError(err.Error()).Render(r.Context(), w)
		if err != nil {
			handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
		}

		return
	}

	w.Header().Set("HX-Redirect", "/login.html")
}

func handleLogin(server *router.HTTPServer) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// handle login
				handleLoginPost(server.Logger, server.TokenCache, server.AuthDB, w, r)
				return
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

func handleLoginPost(
	logger *core.Logger,
	tokenCache *db.TokenCache,
	authDB db.AuthDB,
	w http.ResponseWriter,
	r *http.Request,
) {
	r.ParseForm()
	u := r.FormValue("username")
	p := r.FormValue("password")
	redirect := r.FormValue("redirect")

	logger.Debugf("login - username: %s, password: %s, redirect: %s\n", u, p, redirect)
	// handle login
	if u == "" {
		logger.Debug("login - missing username: %s\n", u)
		returnError(logger, w, r, "Missing username!")
		return
	}
	if p == "" {
		logger.Debug("login - missing password: %s\n", p)
		returnError(logger, w, r, "Missing password!")
		return
	}

	logger.Debug("login - authenticating user")
	user, err := auth.AuthenticateUser(authDB, u, p)
	if err != nil {
		logger.Printf("(%s) %s - login failed: %s\n", router.ClientIP(r), u, err)
		returnError(logger, w, r, err.Error())
		return
	}

	logger.Debug("login - creating bearer token")
	bearer, err := tokenCache.NewBearerToken(user.ID, user.Username, user.Name, user.IsAdmin)
	if err != nil {
		logger.Printf("handleLoginPost: (%s) failed to create bearer token: %s\n", u, err)
		returnError(logger, w, r, "internal server error")
		return
	}

	logger.Debug("login - creating session cookie")
	// Create a new session cookie.
	cookie, err := router.NewSessionCookie(bearer)
	if err != nil {
		logger.Printf("handleLoginPost: (%s) failed to create session_cookie: %s\n", u, err)
		returnError(logger, w, r, "internal server error")
		return
	}

	logger.Debug("login - writing session cookie")
	err = cookie.Write(w)
	if err != nil {
		logger.Printf("handleLoginPost: (%s) failed to write session_cookie: %s\n", u, err)
		returnError(logger, w, r, "internal server error")
		return
	}

	if redirect == "" {
		redirect = "/index.html"
	}

	logger.Printf("(%s) %s (redirect->%s)- login successful\n", router.ClientIP(r), u, redirect)
	w.Header().Set("HX-Redirect", redirect)
}

func returnError(logger *core.Logger, w http.ResponseWriter, r *http.Request, errMsg string) {
	err := components.LoginError(errMsg).Render(r.Context(), w)
	if err != nil {
		handleError(logger, w, r, http.StatusInternalServerError, "internal server error", nil)
	}
}
