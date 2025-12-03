package main

import (
	"log"

	"github.com/daniyarsan/auth-service/internal/config"
	"github.com/daniyarsan/auth-service/internal/usecase/token/refresh_token"
	"github.com/daniyarsan/auth-service/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := cfg.PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&user.User{}, &refresh_token.RefreshToken{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
