package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/KyAnhVo/mystock/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	DBQuerier  *db.DBQueryMachine
	dummy_hash []byte
	is_secure  bool
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
	var login_attempt struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&login_attempt)
	if err != nil {
		w.WriteHeader(405)
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
	_, err = auth.DBQuerier.Querier.Exec(
		ctx,
		"INSERT INTO users.session (session_token, user_id) VALUES ($1, $2)",
		session_token,
		user_id,
	)
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
		Secure:   auth.is_secure,
	})
	w.WriteHeader(200)
}

// Logs the user out, and sends a "delete session" cookie.
func (auth *AuthMiddleware) Logout(w http.ResponseWriter, r *http.Request) {
	// if session token is wrong, authenticate error.
	cookie, err := r.Cookie("session-token")
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintln(w, "authentication error: session cookie not included")
		return
	}

	// remove the current session
	session_token := cookie.Value
	ctx := r.Context()
	_, err = auth.DBQuerier.Querier.Exec(
		ctx,
		"DELETE FROM users.session WHERE session_token = $1",
		session_token,
	)
	if err != nil {
		// if deletion fail, send 500 and report so
		w.WriteHeader(500)
		fmt.Fprintln(w, "cannot log the user out")
		return
	}

	// tell the browser to delete this cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    "",
		HttpOnly: true,
		Secure:   auth.is_secure,
		MaxAge:   -1,
	})
	w.WriteHeader(200)
}

// Creates an account for a user.
func (auth *AuthMiddleware) Signup(w http.ResponseWriter, r *http.Request) {
	var signup_attempt struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&signup_attempt)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "create account error: malformed request body")
		return
	}

	// check if any is empty
	if len(signup_attempt.Username) == 0 ||
		len(signup_attempt.Email) == 0 ||
		len(signup_attempt.Password) == 0 {
		w.WriteHeader(400)
		fmt.Fprintln(w, "create account error: must have all fields: username, email, password")
		return
	}

	// hash the password, check if it is long enough
	password_hashed, err := bcrypt.GenerateFromPassword(
		[]byte(signup_attempt.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "create account error: password too long")
		return
	}

	// Check if email is actually an email, and parsable.
	_, err = mail.ParseAddress(signup_attempt.Email)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "create account error: invalid email format")
		return
	}

	// Create uuid, insert the whole thing into users. If
	// there is error, check if it is constraint violation
	// and send error back correspondingly.
	user_id, _ := uuid.NewV7()
	ctx := r.Context()
	_, err = auth.DBQuerier.Querier.Exec(
		ctx,
		"INSERT INTO users.users (id, email, username, password_hashed) VALUES ($1, $2, $3, $4)",
		user_id,
		signup_attempt.Email,
		signup_attempt.Username,
		password_hashed,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// unique constraint violated
			// there is virtually no chance uuid is repeated
			// thus email or name is repeated.
			w.WriteHeader(400)
			fmt.Fprintln(w, "create account error: username or email is used")
			return
		} else {
			w.WriteHeader(500)
			fmt.Fprintln(w, "server error")
			return
		}
	}

	w.WriteHeader(201)
}

// This function acts as a middleware for any request needing authentication.
// If authentication fails, reject the attempt to run the function.
// If authentication successes, runs the function with the user id in it.
//
// Note: function must have a user_id field which is a UUID.
func (auth *AuthMiddleware) Authenticate(
	w http.ResponseWriter,
	r *http.Request,
	fn func(http.ResponseWriter, *http.Request, uuid.UUID),
) {
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
	err = auth.DBQuerier.Querier.QueryRow(
		ctx,
		"SELECT user_id, expires_at FROM users.session WHERE session_token = $1",
		session_token,
	).Scan(&user_id, &expires_at)
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
