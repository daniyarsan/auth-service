package middleware

import (
	"net/http"
	"strings"

	"github.com/daniyarsan/auth-service/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
)

func JWT(jwtSvc *jwt.Service) gin.HandlerFunc {
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
		token := parts[1]
		uid, err := jwtSvc.ParseAccess(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// save user id in context
		c.Set("user_id", uid)
		c.Next()
	}
}
