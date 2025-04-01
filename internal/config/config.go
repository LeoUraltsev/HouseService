package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	HttpAddress           string
	PostgresURLConnection string
	JWTSecret             []byte
	JWTDuration           time.Duration
}

func New() (*Config, error) {
	var cfg *Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	jwtDuration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		jwtDuration = 1 * time.Hour
	}

	cfg = &Config{
		Env:                   os.Getenv("ENV"),
		HttpAddress:           os.Getenv("HTTP_ADDRESS"),
		PostgresURLConnection: os.Getenv("POSTGRES_CONNECTION"),
		JWTSecret:             []byte(os.Getenv("JWT_SECRET")),
		JWTDuration:           jwtDuration,
	}

	return cfg, nil
}
