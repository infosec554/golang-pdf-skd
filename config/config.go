package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	ServiceName string
	LoggerLevel string

	GotenbergURL string

	AppHost string
	AppPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults")
	}

	cfg := &Config{}

	cfg.AppHost = cast.ToString(getOrReturnDefault("APP_HOST", "localhost"))
	cfg.AppPort = cast.ToString(getOrReturnDefault("APP_PORT", ":8080"))

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "golang-pdf-sdk"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))

	cfg.GotenbergURL = cast.ToString(getOrReturnDefault("GOTENBERG_URL", "http://localhost:3000"))

	return cfg
}

func NewWithURL(gotenbergURL string) *Config {
	return &Config{
		ServiceName:  "golang-pdf-sdk",
		LoggerLevel:  "info",
		GotenbergURL: gotenbergURL,
		AppHost:      "localhost",
		AppPort:      ":8080",
	}
}

func getOrReturnDefault(key string, defaultValue any) any {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
