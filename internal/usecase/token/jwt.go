package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret string, userID int64, expires time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(expires).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
