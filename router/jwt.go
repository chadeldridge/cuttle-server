package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JWT_DEFAULT_PATH             = "/"
	JWT_DEFAULT_SESSION_DURATION = time.Hour * 3
	JWT_COOKIE_NAME              = "session_token"
)

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

type JWTCookie struct {
	Value   string
	Path    string
	Expires time.Time
}

func CreateJWT(username, userid, name, secret string, expires time.Time) (string, error) {
	claims := Claims{
		Username: username,
		UserID:   userid,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("NewJWT: %w", err)
	}

	return signedToken, nil
}

func ParseJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("ParseJWT: unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ParseJWT: %w", err)
	}

	return claims, nil
}

func NewJWTCookie(username, userid, name, secret string) (JWTCookie, error) {
	expires := time.Now().Add(JWT_DEFAULT_SESSION_DURATION)
	j, err := CreateJWT(username, userid, name, secret, expires)
	if err != nil {
		return JWTCookie{}, fmt.Errorf("NewJWTCookie: %w", err)
	}

	return JWTCookie{
		Value:   j,
		Path:    JWT_DEFAULT_PATH,
		Expires: expires,
	}, nil
}

func NewJWTCookieFromCookie(cookie *http.Cookie) JWTCookie {
	return JWTCookie{
		Value:   cookie.Value,
		Path:    cookie.Path,
		Expires: cookie.Expires,
	}
}

func JWTGetCookie(r *http.Request) (JWTCookie, error) {
	cookie, err := r.Cookie(JWT_COOKIE_NAME)
	if err != nil {
		return JWTCookie{}, err
	}

	j := NewJWTCookieFromCookie(cookie)
	return j, nil
}

func (j JWTCookie) Validate() error {
	if j.Value == "" {
		return fmt.Errorf("JWTCookie.Validate: value is empty")
	}

	if j.Path == "" {
		return fmt.Errorf("JWTCookie.Validate: path is empty")
	}

	if j.Expires.IsZero() {
		return fmt.Errorf("JWTCookie.Validate: expires is empty")
	}

	if j.Expires.Before(time.Now()) {
		return fmt.Errorf("JWTCookie.Validate: expires is in the past")
	}

	return nil
}

func (j JWTCookie) JWTRefreshCookie(secret string) (JWTCookie, error) {
	err := j.Validate()
	if err != nil {
		return JWTCookie{}, fmt.Errorf("JWTRefreshCookie: %w", err)
	}

	claims, err := ParseJWT(j.Value, secret)
	if err != nil {
		return JWTCookie{}, fmt.Errorf("JWTRefreshCookie: %w", err)
	}

	expires := time.Now().Add(JWT_DEFAULT_SESSION_DURATION)
	newToken, err := CreateJWT(claims.Username, claims.UserID, claims.Name, secret, expires)
	if err != nil {
		return JWTCookie{}, fmt.Errorf("JWTRefreshCookie: %w", err)
	}

	return JWTCookie{
		Value:   newToken,
		Path:    j.Path,
		Expires: expires,
	}, nil
}

func (j JWTCookie) JWTDeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    JWT_COOKIE_NAME,
		Value:   "",
		Expires: time.Now().Add(-JWT_DEFAULT_SESSION_DURATION),
	})
}

func (j JWTCookie) JWTWriteCookie(w http.ResponseWriter) error {
	err := j.Validate()
	if err != nil {
		return fmt.Errorf("JWTWriteCookie: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     JWT_COOKIE_NAME,
		Value:    j.Value,
		Path:     j.Path,
		Expires:  j.Expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
