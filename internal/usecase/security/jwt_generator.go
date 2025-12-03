package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	Secret string
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{Secret: secret}
}

func (j *JWTGenerator) Generate(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(j.Secret))
}

func (j *JWTGenerator) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
