package middleware

import (
	"net/http"

	"incidenttrack/internal/pkg/httpx"
	"incidenttrack/internal/pkg/token"
)

// Authenticator valida a sessão a partir do cookie HttpOnly.
type Authenticator struct {
	tokens     token.Service
	denylist   *token.Denylist
	cookieName string
}

// NewAuthenticator injeta as dependências do middleware de auth.
func NewAuthenticator(tokens token.Service, denylist *token.Denylist, cookieName string) *Authenticator {
	return &Authenticator{tokens: tokens, denylist: denylist, cookieName: cookieName}
}

// Require exige uma sessão válida. Em qualquer falha responde 401 com
// mensagem genérica (não revela se o token expirou, é malformado ou foi
// revogado — OWASP A01/A07). Em sucesso, injeta os claims no contexto.
func (a *Authenticator) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(a.cookieName)
		if err != nil || c.Value == "" {
			a.deny(w)
			return
		}

		claims, err := a.tokens.Verify(c.Value)
		if err != nil {
			a.deny(w)
			return
		}

		if a.denylist != nil && a.denylist.Revoked(claims.TokenID) {
			a.deny(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
	})
}

func (a *Authenticator) deny(w http.ResponseWriter) {
	httpx.WriteError(w, http.StatusUnauthorized, "Sessão expirada. Faça login novamente.")
}
