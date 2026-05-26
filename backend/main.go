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
	authHandler := handler.NewAuthMiddleware(database, logger)
	stockHandler := handler.NewStockHandler(database, logger)
	corsHandler := handler.NewCORSMiddleware()

	// Preflight
	http.HandleFunc("OPTIONS /", corsHandler.Middleware(func(w http.ResponseWriter, r *http.Request) {}))

	// Authentication
	http.HandleFunc("POST /api/auth/login", corsHandler.Middleware(authHandler.Login))
	http.HandleFunc("POST /api/auth/logout", corsHandler.Middleware(authHandler.Logout))
	http.HandleFunc("POST /api/auth/signup", corsHandler.Middleware(authHandler.Signup))

	// Simple ticker functionalities
	http.HandleFunc("GET /api/ticker/{ticker}", corsHandler.Middleware(stockHandler.OverviewTicker))

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
