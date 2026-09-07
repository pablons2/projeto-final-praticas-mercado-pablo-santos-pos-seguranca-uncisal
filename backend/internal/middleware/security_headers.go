package middleware

import "net/http"

// SecurityHeaders adiciona cabeçalhos de segurança a toda resposta
// (OWASP A02 — Security Misconfiguration; docs/plano.md §7).
//
// Esta é uma API JSON pura (não serve HTML), então a CSP é a mais
// restritiva possível: nada pode ser carregado ou embutido a partir de
// uma resposta desta origem.
func SecurityHeaders(isProd bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Content-Security-Policy",
				"default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "no-referrer")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-site")
			h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), browsing-topics=()")
			h.Set("Cache-Control", "no-store")

			// HSTS só faz sentido sob HTTPS real (produção atrás do Nginx +
			// Certbot). Em dev sobre HTTP o header é ignorado pelo browser,
			// mas evitamos anunciá-lo para não confundir.
			if isProd {
				h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			}

			next.ServeHTTP(w, r)
		})
	}
}
