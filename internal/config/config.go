package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	NATSURL     string
	HTTPAddr    string
	NATSSubject string
	LogLevel    string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DatabaseURL: mustEnv("DATABASE_URL"),
		NATSURL:     getenv("NATS_URL", "nats://localhost:4222"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		NATSSubject: getenv("NATS_SUBJECT", "updates"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing required env: " + key)
	}
	return v
}
