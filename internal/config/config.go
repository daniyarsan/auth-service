package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        int
	DBName        string
	DBUser        string
	DBPassword    string
	DBSSLMode     string
	JWTAccessKey  string
	JWTRefreshKey string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func Load() *Config {
	p, _ := strconv.Atoi(getenv("POSTGRES_PORT", "5432"))
	accessMinutes, _ := strconv.Atoi(getenv("ACCESS_TTL_MIN", "15"))
	refreshDays, _ := strconv.Atoi(getenv("REFRESH_TTL_DAYS", "30"))

	return &Config{
		Port:          getenv("PORT", "8080"),
		DBHost:        getenv("POSTGRES_HOST", "localhost"),
		DBPort:        p,
		DBName:        getenv("POSTGRES_DB", "auth"),
		DBUser:        getenv("POSTGRES_USER", "postgres"),
		DBPassword:    getenv("POSTGRES_PASSWORD", ""),
		DBSSLMode:     getenv("PGSSLMODE", "disable"),
		JWTAccessKey:  getenv("JWT_ACCESS_SECRET", "access-secret"),
		JWTRefreshKey: getenv("JWT_REFRESH_SECRET", "refresh-secret"),
		AccessTTL:     time.Minute * time.Duration(accessMinutes),
		RefreshTTL:    time.Hour * 24 * time.Duration(refreshDays),
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
