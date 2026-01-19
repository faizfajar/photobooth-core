package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN     string
	AppPort   string
	JWTSecret string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBDSN:     os.Getenv("DATABASE_URL"),
		AppPort:   os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// VALIDATOR: Langsung hentikan aplikasi jika config krusial kosong
	if cfg.DBDSN == "" || cfg.JWTSecret == "" {
		slog.Error("Konfigurasi KRITIKAL hilang! Cek DATABASE_URL dan JWT_SECRET di .env")
		os.Exit(1)
	}

	return cfg
}
