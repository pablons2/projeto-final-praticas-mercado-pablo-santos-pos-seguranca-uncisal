package httpapi

import (
	"net/http"
	"time"
)

// CookieConfig controla os atributos do cookie de sessão (docs/plano.md §5).
type CookieConfig struct {
	Name   string
	Secure bool
	Domain string
	TTL    time.Duration
}

// setSessionCookie grava o JWT num cookie HttpOnly. O JavaScript do
// navegador nunca tem acesso a ele — mesmo um XSS no frontend não
// consegue ler o token (OWASP A07; docs/plano.md §5).
func (c CookieConfig) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    token,
		Path:     "/",
		Domain:   c.Domain,
		MaxAge:   int(c.TTL.Seconds()),
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// clearSessionCookie expira o cookie no cliente (logout).
func (c CookieConfig) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    "",
		Path:     "/",
		Domain:   c.Domain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSiteStrictMode,
	})
}
