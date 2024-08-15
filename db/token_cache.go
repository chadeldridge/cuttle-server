package db

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle-server/core"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

type TokenCache struct {
	*Cache
	secret string
}

// Create a new TokenCache instance. The 'max' expire duration must be greater than or equal to the
// 'expire' time. If 'expire' or 'max' are 0, they will be set to the defaults of 3 and 6 hours.
func NewTokenCache(secret string, expire, max time.Duration) (*TokenCache, error) {
	if secret == "" {
		return nil, fmt.Errorf("NewTokenCache: secret - %s", core.ErrParamEmpty)
	}

	if expire > max {
		expire = max
	}

	if expire != 0 {
		sessionExpires = expire
	}

	if max != 0 {
		sessionMaxExpires = max
	}

	return &TokenCache{
		Cache:  NewCache(sessionExpires, time.Minute*5),
		secret: secret,
	}, nil
}

func (c *TokenCache) NewBearerToken(userID int64, username, name string, isAdmin bool) (string, error) {
	bearer := hex.EncodeToString([]byte(uuid.New().String()))

	token, err := CreateJWT(userID, username, name, c.secret, isAdmin)
	if err != nil {
		return "", fmt.Errorf("NewBearerToken: %w", err)
	}

	err = c.CacheToken(bearer, token)
	if err != nil {
		return "", fmt.Errorf("NewBearerToken: %w", err)
	}

	return bearer, nil
}

// GetClaims retrieves and parses the JWT from cache using the token as the key.
func (c *TokenCache) GetClaims(bearer string) (*Claims, error) {
	if bearer == "" {
		return nil, fmt.Errorf("GetClaims: token %s", core.ErrParamEmpty)
	}

	v, err := c.GetToken(bearer)
	if err != nil {
		return nil, fmt.Errorf("GetClaims: %w", err)
	}

	claims, err := ParseJWT(v, c.secret)
	if err != nil {
		return nil, fmt.Errorf("GetClaims: %w", err)
	}

	token, err := RefreshJWT(claims, c.secret)
	if err != nil {
		return nil, fmt.Errorf("GetClaims: %w", err)
	}

	err = c.CacheToken(bearer, token)
	if err != nil {
		return nil, fmt.Errorf("GetClaims: %w", err)
	}

	return claims, nil
}

func (c *TokenCache) CacheToken(bearer, token string) error {
	if c == nil {
		return fmt.Errorf("CacheBearerToken: cache %s", core.ErrParamEmpty)
	}

	if bearer == "" {
		return fmt.Errorf("CacheBearerToken: token %s", core.ErrParamEmpty)
	}

	if token == "" {
		return fmt.Errorf("CacheBearerToken: authToken %s", core.ErrParamEmpty)
	}

	c.Set(bearer, token)
	return nil
}

func (c *TokenCache) GetToken(bearer string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("GetBearerToken: cache %s", core.ErrParamEmpty)
	}

	if bearer == "" {
		return "", fmt.Errorf("GetBearerToken: token %s", core.ErrParamEmpty)
	}

	if v, ok := c.Get(bearer); ok {
		return v.(string), nil
	}

	return "", fmt.Errorf("GetBearerToken: %s", ErrKeyNotFound)
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func CreateJWT(userID int64, username, name, secret string, isAdmin bool) (string, error) {
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
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(sessionExpires))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("RefreshJWT: %w", err)
	}

	return signedToken, nil
}

func ParseJWT(tokenString, secret string) (*Claims, error) {
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
