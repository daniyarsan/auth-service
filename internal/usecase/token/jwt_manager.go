package token

import (
	"context"
	"errors"
	"time"

	"github.com/goccy/go-yaml/token"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJwtManager(secret string, accessTTL, refreshTTL time.Duration) token.TokenManager {
	return &JwtManager{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

type claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func (m *JwtManager) Generate(ctx context.Context, userID int64) (token.Tokens, error) {
	now := time.Now().UTC()

	accessExp := now.Add(m.accessTTL)
	accessClaims := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := at.SignedString([]byte(m.secret))
	if err != nil {
		return token.Tokens{}, err
	}

	// Refresh token (simple; in prod, store it or use rotating refresh tokens)
	refreshExp := now.Add(m.refreshTTL)
	refreshClaims := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := rt.SignedString([]byte(m.secret))
	if err != nil {
		return token.Tokens{}, err
	}

	return token.Tokens{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
	}, nil
}

func (m *JwtManager) ParseAccessToken(ctx context.Context, tokenStr string) (int64, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := parsed.Claims.(*claims); ok && parsed.Valid {
		return claims.UserID, nil
	}
	return 0, errors.New("invalid token")
}
