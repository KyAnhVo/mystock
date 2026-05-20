package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KyAnhVo/mystock/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	DBQuerier  *db.DBQueryMachine
	dummy_hash []byte
}

var authMiddleware *AuthMiddleware

func Init(db *db.DBQueryMachine) *AuthMiddleware {
	if authMiddleware == nil {
		dummy_hash, _ := bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)
		authMiddleware = &AuthMiddleware{
			DBQuerier:  db,
			dummy_hash: dummy_hash,
		}
	}

	return authMiddleware
}

// Logs a user in. Only accepts either email or username, not both.
// If authenticate fails, send 401 if auth is wrong or 500 if server error.
// If authenticate successes, send 200 and a cookie with session token.
func (auth *AuthMiddleware) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintln(w, "authentication error: wrong method type")
		return
	}

	var login_attempt struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&login_attempt)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintln(w, "authentication error: malformed request body")
		return
	}

	if len(login_attempt.Email) == 0 && len(login_attempt.Username) == 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "authentication error: missing email and username")
		return
	}
	if len(login_attempt.Email) != 0 && len(login_attempt.Username) != 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "authentication error: only specify one of username or email")
		return
	}

	var hashed_password []byte
	var user_id uuid.UUID

	// Query user login, and check if login is correct
	ctx := r.Context()
	var query string
	var id_str string
	if len(login_attempt.Email) > 0 {
		query = "SELECT password_hashed, id FROM users.users WHERE email = $1"
		id_str = login_attempt.Email
	} else {
		query = "SELECT password_hashed, id FROM users.users WHERE username = $1"
		id_str = login_attempt.Username
	}
	err = auth.DBQuerier.Querier.QueryRow(ctx, query, id_str).Scan(&hashed_password, &user_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// dummy bcrypt run so that attacker cannot find correct email/username
			_ = bcrypt.CompareHashAndPassword(auth.dummy_hash, []byte(login_attempt.Password))
			w.WriteHeader(401)
			fmt.Fprintln(w, "authentication error: username or email or password is incorrect")
		} else {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
		}
		return
	}
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(login_attempt.Password))
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintln(w, "authentication error: username or email or password is incorrect")
		return
	}

	// if login is correct, we create a session.
	bytes := make([]byte, 32)
	rand.Read(bytes)
	session_token := hex.EncodeToString(bytes)
	_, err = auth.DBQuerier.Querier.Exec(ctx, "INSERT INTO users.session (session_token, user_id) VALUES ($1, $2)", session_token, user_id)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "server error: cannot create session")
		return
	}

	// then sends the user back the auth is complete with a session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    session_token,
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(200)
}

func (auth *AuthMiddleware) Logout(w http.ResponseWriter, r *http.Request) {

}

func (auth *AuthMiddleware) Signup(w http.ResponseWriter, r *http.Request) {

}

// This function acts as a middleware for any request needing authentication.
// If authentication fails, reject the attempt to run the function.
// If authentication successes, runs the function with the user id in it.
//
// Note: function must have a user_id field which is a UUID.
func (auth *AuthMiddleware) Authenticate(w http.ResponseWriter, r *http.Request, fn func(http.ResponseWriter, *http.Request, uuid.UUID)) {
	cookie, err := r.Cookie("session-token")
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintln(w, "authentication error: session cookie not included")
		return
	}
	session_token := cookie.Value

	var expires_at time.Time
	var user_id uuid.UUID
	ctx := r.Context()
	err = auth.DBQuerier.Querier.QueryRow(ctx, "SELECT user_id, expires_at FROM users.session WHERE session_token = $1", session_token).Scan(&user_id, &expires_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(401)
			fmt.Fprintln(w, "authentication error: no session found")
		} else {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
		}

		return
	}
	if expires_at.Before(time.Now()) {
		w.WriteHeader(401)
		fmt.Fprintln(w, "authentication error: session expired")
		return
	}

	fn(w, r, user_id)
}
