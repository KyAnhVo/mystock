package main

import (
	"fmt"
	"net/http"

	"github.com/KyAnhVo/mystock/config"
	"github.com/KyAnhVo/mystock/internal/auth"
	"github.com/KyAnhVo/mystock/internal/db"
	"github.com/google/uuid"
)

func main() {
	config.GetCfg()
	database, err := db.Init()
	if err != nil {
		fmt.Println("fail to connect to DB")
		return
	}
	auth := auth.Init(database)

	// Authentication
	http.HandleFunc("POST /api/auth/login", auth.Login)
	http.HandleFunc("POST /api/auth/logout", auth.Logout)
	http.HandleFunc("POST /api/auth/signup", auth.Signup)
	http.HandleFunc("GET /api/auth/test", func(w http.ResponseWriter, r *http.Request) {
		auth.Authenticate(w, r, func(w http.ResponseWriter, r *http.Request, user_id uuid.UUID) {
			ctx := r.Context()
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
		})
	})

	// Finally, run it.
	http.ListenAndServe(":8080", nil)
}
