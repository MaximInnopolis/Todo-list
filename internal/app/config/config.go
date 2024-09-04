package config

import (
	"fmt"
	"os"
)

var defaultHttpPort = ":8080"

type Config struct {
	DbUrl    string
	HttpPort string
}

func New() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не задан")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
	}

	return &Config{
		DbUrl:    dbURL,
		HttpPort: httpPort,
	}, nil
}
