package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// RequestIDHeader é o cabeçalho de correlação exposto na resposta.
const RequestIDHeader = "X-Request-ID"

// RequestID injeta um identificador único por requisição no contexto e
// no cabeçalho de resposta, para correlacionar logs sem expor dados
// internos. Um X-Request-ID vindo do cliente é ignorado (poderia ser
// usado para forjar rastros de log).
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := uuid.NewString()
			w.Header().Set(RequestIDHeader, id)
			next.ServeHTTP(w, r.WithContext(WithRequestID(r.Context(), id)))
		})
	}
}
