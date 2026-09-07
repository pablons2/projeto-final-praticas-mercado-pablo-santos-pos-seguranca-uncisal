package middleware

import (
	"net/http"

	"incidenttrack/internal/pkg/httpx"
)

// CSRF protege as rotas de mutação numa API baseada em cookie de sessão
// (docs/plano.md §5). Complementa SameSite=Strict com duas checagens
// que um formulário cross-site forjado não consegue satisfazer:
//
//  1. presença do header customizado X-Requested-With — que dispara
//     preflight CORS e não pode ser enviado por <form>/<img> cross-site;
//  2. quando há header Origin, ele precisa bater com a origem do frontend.
//
// Métodos seguros (GET/HEAD/OPTIONS) passam direto.
func CSRF(allowedOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(w, r)
				return
			}

			if r.Header.Get("X-Requested-With") == "" {
				httpx.WriteError(w, http.StatusForbidden,
					"Requisição bloqueada por proteção CSRF.")
				return
			}

			if origin := r.Header.Get("Origin"); origin != "" && origin != allowedOrigin {
				httpx.WriteError(w, http.StatusForbidden,
					"Origem não autorizada.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
