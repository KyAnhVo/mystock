package handler

import (
	"net/http"

	"github.com/KyAnhVo/mystock/config"
)

type CORSMiddleware struct {
	allowedOrigin string
}

func NewCORSMiddleware() *CORSMiddleware {
	cfg := config.GetCfg()
	return &CORSMiddleware{
		allowedOrigin: cfg.AllowedOrigin,
	}
}

func (c *CORSMiddleware) Middleware(
	fn func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		fn(w, r)
	}
}
