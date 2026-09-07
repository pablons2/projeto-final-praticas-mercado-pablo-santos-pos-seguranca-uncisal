package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"incidenttrack/internal/pkg/httpx"
)

// Recover captura qualquer panic na cadeia, registra o stack trace
// apenas no log do servidor e devolve um 500 genérico ao cliente —
// nunca o stack trace (OWASP A10 — Mishandling of Exceptional
// Conditions; docs/plano.md §7).
func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// http.ErrAbortHandler é o mecanismo idiomático de abortar
					// uma resposta; propaga sem logar como erro.
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					log.Error("panic recuperado",
						"err", rec,
						"method", r.Method,
						"path", r.URL.Path,
						"request_id", RequestIDFrom(r.Context()),
						"stack", string(debug.Stack()),
					)
					httpx.WriteError(w, http.StatusInternalServerError,
						"Algo deu errado no servidor. Tente novamente.")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
