package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Global config from env
var cfg Config

// True if Cfg is loaded
var cfgLoaded bool

type Config struct {
	// Program, connection configs
	Port   string
	DBConn string

	// Stock configs
	StockApiKey    string
	StockApiHeader string
}

func Load() Config {
	return Config{
		Port:   os.Getenv("PORT"),
		DBConn: os.Getenv("DB_URL"),

		StockApiKey:    os.Getenv("STOCK_API_KEY"),
		StockApiHeader: os.Getenv("STOCK_API_HEADER"),
	}
}

func GetCfg() *Config {
	if !cfgLoaded {
		initCfg()
	}
	return &cfg
}

func initCfg() {
	godotenv.Load()
	cfg = Load()
	cfgLoaded = true
}
