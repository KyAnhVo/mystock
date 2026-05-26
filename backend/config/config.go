package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Global config from env
var cfg *Config

type Config struct {
	// Program, connection configs
	Port          string
	DBConn        string
	AllowedOrigin string

	// Stock configs
	StockApiKey    string
	StockApiHeader string
}

func GetCfg() *Config {
	if cfg == nil {
		initCfg()
	}
	return cfg
}

func initCfg() {
	godotenv.Load()
	cfg = &Config{
		Port:          os.Getenv("PORT"),
		DBConn:        os.Getenv("DB_URL"),
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),

		StockApiKey:    os.Getenv("STOCK_API_KEY"),
		StockApiHeader: os.Getenv("STOCK_API_HEADER"),
	}
}
