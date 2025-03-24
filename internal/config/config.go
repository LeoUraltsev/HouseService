package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	HttpAddress           string
	PostgresURLConnection string
}

func New() (*Config, error) {
	var cfg *Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	cfg = &Config{
		Env:                   os.Getenv("ENV"),
		HttpAddress:           os.Getenv("HTTP_ADDRESS"),
		PostgresURLConnection: os.Getenv("POSTGRES_CONNECTION"),
	}

	return cfg, nil
}
