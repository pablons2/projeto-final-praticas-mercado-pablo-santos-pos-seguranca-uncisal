package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// CORS libera apenas a origem exata do frontend, com credenciais
// (cookies) habilitadas. Nunca reflete origem arbitrária nem usa "*"
// (OWASP A02 — Security Misconfiguration; docs/plano.md §7).
//
// allowedOrigin deve ser algo como "https://app.exemplo.com".
func CORS(allowedOrigin string) Middleware {
	const allowedMethods = "GET, POST, PUT, DELETE, OPTIONS"
	const allowedHeaders = "Content-Type, Accept, X-Requested-With"
	maxAge := strconv.Itoa(int((12 * time.Hour).Seconds()))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Vary sempre, para caches não misturarem respostas por origem.
			w.Header().Add("Vary", "Origin")

			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
				w.Header().Set("Access-Control-Expose-Headers", RequestIDHeader)
				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			// Preflight: encerra aqui.
			if r.Method == http.MethodOptions {
				// 204 mesmo para origem não permitida — sem corpo, sem detalhe.
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
