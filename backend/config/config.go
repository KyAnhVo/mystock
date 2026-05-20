package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Global config from env
var cfg *Config

type Config struct {
	// Program, connection configs
	Port   string
	DBConn string

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
		Port:   os.Getenv("PORT"),
		DBConn: os.Getenv("DB_URL"),

		StockApiKey:    os.Getenv("STOCK_API_KEY"),
		StockApiHeader: os.Getenv("STOCK_API_HEADER"),
	}
}
