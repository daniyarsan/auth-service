package middleware

import (
	"net/http"
	"strings"

	jwtinfra "github.com/daniyarsan/auth-service/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
)

func JWT(jwtSvc *jwtinfra.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		uid, err := jwtSvc.ParseAccess(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", uid)
		c.Next()
	}
}
