package config

import "os"

// Global config from env
var Cfg Config

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

func Init() {
	Cfg = Load()
}
