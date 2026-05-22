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
	config.GetCfg()
	database, err := db.Init()
	if err != nil {
		fmt.Println("fail to connect to DB")
		return
	}
	logger, err := createHandler()
	if err != nil {
		fmt.Println("logger creation error")
		os.Exit(1)
	}
	auth := handler.NewAuthMiddleware(database, logger)

	// Authentication
	http.HandleFunc("POST /api/auth/login", auth.Login)
	http.HandleFunc("POST /api/auth/logout", auth.Logout)
	http.HandleFunc("POST /api/auth/signup", auth.Signup)
	http.HandleFunc(
		"GET /api/auth/test",
		auth.Authenticate(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user_id, _ := handler.RetrieveUserID(ctx)
			var username string
			err := database.Querier.QueryRow(
				ctx,
				"SELECT username FROM users.users WHERE id = $1",
				user_id,
			).Scan(&username)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, "server error")
				return
			}
			w.WriteHeader(200)
			fmt.Fprintf(w, "Hello, %v\n", username)
		}),
	)

	// Finally, run it.
	http.ListenAndServe(":8080", nil)
}

func createHandler() (*slog.Logger, error) {
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
