package db

import (
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/golang-jwt/jwt/v4"
)

const (
	JWT_DEFAULT_PATH                = "/"
	JWT_DEFAULT_SESSION_EXPIRES     = time.Hour * 3
	JWT_DEFAULT_MAX_SESSION_EXPIRES = time.Hour * 6
	JWT_COOKIE_NAME                 = "session_token"
)

var (
	sessionExpires    = JWT_DEFAULT_SESSION_EXPIRES
	sessionMaxExpires = JWT_DEFAULT_MAX_SESSION_EXPIRES
)

var (
	ErrExpiredCookie  = fmt.Errorf("cookie has expired")
	ErrSessionExpired = fmt.Errorf("session has expired")
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func CreateJWT(userID int64, username, name, secret string, isAdmin bool) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("CreateJWT: secret - %w", core.ErrParamEmpty)
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Name:     name,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionExpires)),
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

func RefreshJWT(claims *Claims, secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("CreateJWT: secret - %w", core.ErrParamEmpty)
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(sessionExpires))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("RefreshJWT: %w", err)
	}

	return signedToken, nil
}

func ParseJWT(tokenString, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("ParseJWT: token - %w", core.ErrParamEmpty)
	}

	if secret == "" {
		return nil, fmt.Errorf("ParseJWT: secret - %w", core.ErrParamEmpty)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
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

	if !token.Valid {
		// Check for expiration.
		if !claims.VerifyExpiresAt(time.Now(), true) {
			return nil, fmt.Errorf("ParseJWT: %w", ErrSessionExpired)
		}

		// Check for max duration.
		if claims.VerifyIssuedAt(time.Now().Add(-sessionMaxExpires), true) {
			return nil, fmt.Errorf("ParseJWT: %w", ErrSessionExpired)
		}

		return nil, fmt.Errorf("ParseJWT: %w", jwt.ErrSignatureInvalid)
	}

	return claims, nil
}
