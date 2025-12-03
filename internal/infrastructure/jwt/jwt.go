package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Service struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type Claims struct {
	Sub int64 `json:"sub"`
	jwt.RegisteredClaims
}

func New(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *Service) Generate(userID int64) (string, string, error) {
	now := time.Now()

	accessClaims := Claims{
		Sub: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessTok := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessTok.SignedString(s.accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		Sub: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshTok := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshTok.SignedString(s.refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

func (s *Service) parse(tokenStr string, secret []byte) (int64, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}
	if !tok.Valid {
		return 0, ErrInvalidToken
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return claims.Sub, nil
}

func (s *Service) ParseAccess(tokenStr string) (int64, error) {
	return s.parse(tokenStr, s.accessSecret)
}
func (s *Service) ParseRefresh(tokenStr string) (int64, error) {
	return s.parse(tokenStr, s.refreshSecret)
}

func (s *Service) AccessTTL() time.Duration  { return s.accessTTL }
func (s *Service) RefreshTTL() time.Duration { return s.refreshTTL }
