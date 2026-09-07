// Package middleware reúne os middlewares HTTP de segurança e
// observabilidade (docs/plano.md §4, §7): request-id, recover, CORS,
// security headers, rate limit, CSRF e autenticação JWT.
//
// Todos têm a assinatura func(http.Handler) http.Handler e são
// compostos por Chain, na ordem definida em handler/http/router.go.
package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"incidenttrack/internal/pkg/token"
)

// Middleware é a unidade de composição.
type Middleware func(http.Handler) http.Handler

// Chain aplica os middlewares na ordem em que são passados: o primeiro
// da lista é o mais externo (primeiro a ver o request, último a ver a
// response).
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type ctxKey int

const (
	ctxKeyRequestID ctxKey = iota
	ctxKeyClaims
)

// WithRequestID guarda o id no contexto.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}

// RequestIDFrom recupera o id do contexto ("" se ausente).
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyRequestID).(string)
	return id
}

// WithClaims guarda os claims do token autenticado no contexto.
func WithClaims(ctx context.Context, c token.Claims) context.Context {
	return context.WithValue(ctx, ctxKeyClaims, c)
}

// ClaimsFrom recupera os claims; ok=false se a rota não passou pelo auth.
func ClaimsFrom(ctx context.Context) (token.Claims, bool) {
	c, ok := ctx.Value(ctxKeyClaims).(token.Claims)
	return c, ok
}

// UserIDFrom é um atalho para o subject dos claims.
func UserIDFrom(ctx context.Context) string {
	c, ok := ClaimsFrom(ctx)
	if !ok {
		return ""
	}
	return c.Subject
}

// ClientIP extrai o IP do cliente. Só confia em X-Forwarded-For quando
// trustProxy é true (atrás de Nginx no deploy); caso contrário usa o
// RemoteAddr real, evitando spoof do rate limiter (OWASP A07). Exportado
// para que handlers usem exatamente a mesma chave do rate limiter.
func ClientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// primeiro IP da lista é o cliente original
			if i := strings.IndexByte(xff, ','); i >= 0 {
				return strings.TrimSpace(xff[:i])
			}
			return strings.TrimSpace(xff)
		}
		if xr := r.Header.Get("X-Real-IP"); xr != "" {
			return strings.TrimSpace(xr)
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
