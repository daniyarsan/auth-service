package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/daniyarsan/auth-service/internal/config"
	dbinfra "github.com/daniyarsan/auth-service/internal/infrastructure/db"
	httpRouter "github.com/daniyarsan/auth-service/internal/infrastructure/http/router"
	jwtinfra "github.com/daniyarsan/auth-service/internal/infrastructure/jwt"
	reprepo "github.com/daniyarsan/auth-service/internal/infrastructure/repository"
	authuc "github.com/daniyarsan/auth-service/internal/usecase/auth"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	conn, err := dbinfra.New(dbinfra.Settings{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBName,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer conn.Pool.Close()

	// repos
	userRepo := reprepo.NewUserRepository(conn.Pool)
	tokenRepo := reprepo.NewRefreshTokenRepository(conn.Pool)

	// jwt service
	jwtSvc := jwtinfra.New(cfg.JWTAccessKey, cfg.JWTRefreshKey, cfg.AccessTTL, cfg.RefreshTTL)

	// usecase service
	authSvc := authuc.NewService(userRepo, tokenRepo, jwtSvc)

	// http
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	httpRouter.Setup(r, authSvc, jwtSvc)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
