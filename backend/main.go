package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/KyAnhVo/mystock/internal/handler"
)

func main() {
	cfg := config.GetCfg()
	database, err := db.Init()
	if err != nil {
		fmt.Println("fail to connect to DB")
		return
	}
	logger, err := createLoggingHandler()
	if err != nil {
		fmt.Println("logger creation error")
		os.Exit(1)
	}
	auth := handler.NewAuthMiddleware(database, logger)
	stockHandler := handler.NewStockHandler(database)

	// Authentication
	http.HandleFunc("POST /api/auth/login", auth.Login)
	http.HandleFunc("POST /api/auth/logout", auth.Logout)
	http.HandleFunc("POST /api/auth/signup", auth.Signup)

	// Simple ticker functionalities
	http.HandleFunc("GET /api/ticker/{ticker}", stockHandler.OverviewTicker)

	// Finally, run it.
	http.ListenAndServe(cfg.Port, nil)
}

func createLoggingHandler() (*slog.Logger, error) {
	handlerOption := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	logFile, err := os.OpenFile("./backend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	multi := slog.NewMultiHandler(
		slog.NewJSONHandler(logFile, handlerOption),
		slog.NewJSONHandler(os.Stdout, handlerOption),
	)
	logger := slog.New(multi)

	return logger, nil
}
