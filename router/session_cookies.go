package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/chadeldridge/cuttle-server/db"
)

const SESSION_COOKIE_NAME = "session_token"

var ErrExpiredCookie = fmt.Errorf("cookie has expired")

type SessionCookie struct {
	Value   string // Bearer token used to lookup JWT from AuthCache.
	Path    string
	Expires time.Time // Cookie expiration time.
}

func NewSessionCookie(bearer string) (SessionCookie, error) {
	if bearer == "" {
		return SessionCookie{}, fmt.Errorf("NewSessionCookie: bearer - %s", core.ErrParamEmpty)
	}

	expires := time.Now().Add(db.JWT_DEFAULT_SESSION_EXPIRES)
	return SessionCookie{
		Value:   bearer,
		Path:    db.JWT_DEFAULT_PATH,
		Expires: expires,
	}, nil
}

func NewSessionCookieFromCookie(cookie *http.Cookie) (SessionCookie, error) {
	if cookie == nil {
		return SessionCookie{}, fmt.Errorf("NewSessionCookieFromCookie: cookie %s", core.ErrParamEmpty)
	}

	s := SessionCookie{
		Value: cookie.Value,
	}

	return s, s.Validate()
}

func GetSessionCookie(r *http.Request) (SessionCookie, error) {
	cookie, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return SessionCookie{}, err
	}

	// Create, validate, and return a new SessionCookie.
	return NewSessionCookieFromCookie(cookie)
}

func (s SessionCookie) Validate() error {
	if s.Value == "" {
		return fmt.Errorf("SessionCookie.Validate: value - %s", core.ErrParamEmpty)
	}

	return nil
}

func (s SessionCookie) Write(w http.ResponseWriter) error {
	err := s.Validate()
	if err != nil {
		return fmt.Errorf("SessionCookie.Write: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    s.Value,
		Path:     s.Path,
		Expires:  s.Expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (s SessionCookie) Delete(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
