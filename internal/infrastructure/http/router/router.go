package router

import (
	middleware "github.com/daniyarsan/auth-service/internal/infrastructure/http"
	"github.com/gin-gonic/gin"

	"github.com/daniyarsan/auth-service/internal/infrastructure/http/handler"
	jwtinfra "github.com/daniyarsan/auth-service/internal/infrastructure/jwt"
	authuc "github.com/daniyarsan/auth-service/internal/usecase/auth"
)

func Setup(r *gin.Engine, svc *authuc.Service, jwtSvc *jwtinfra.Service) {
	h := handler.NewAuthHandler(svc)
	authMW := middleware.JWT(jwtSvc)

	api := r.Group("/api")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.POST("/refresh", h.Refresh)

		protected := api.Group("/")
		protected.Use(authMW)
		{
			protected.GET("/me", h.Me)
			protected.POST("/logout", h.Logout)
		}
	}
}
