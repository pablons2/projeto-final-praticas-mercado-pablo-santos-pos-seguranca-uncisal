package httpapi

import (
	"log/slog"
	"net/http"

	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/httpx"
)

// RouterDeps são todas as dependências necessárias para montar as rotas.
type RouterDeps struct {
	Auth      *AuthHandler
	Incidents *IncidentHandler

	Authenticator *middleware.Authenticator
	APILimiter    *middleware.Limiter
	LoginLimiter  *middleware.Limiter

	AllowedOrigin string
	TrustProxy    bool
	IsProduction  bool
	Log           *slog.Logger
}

// NewRouter constrói o http.Handler completo: mux de rotas + cadeia
// global de middlewares. A ordem da cadeia global (de fora para dentro):
//
//	RequestID → Recover → RequestLogger → SecurityHeaders → CORS → rotas
//
// Assim tanto o Recover quanto o logger enxergam o request_id, e o
// preflight CORS é respondido antes de qualquer roteamento.
func NewRouter(d RouterDeps) http.Handler {
	mux := http.NewServeMux()

	authRequired := d.Authenticator.Require
	csrf := middleware.CSRF(d.AllowedOrigin)
	apiRL := middleware.RateLimit(d.APILimiter, d.TrustProxy, "Muitas requisições. Tente novamente mais tarde.")
	loginRL := middleware.RateLimit(d.LoginLimiter, d.TrustProxy, "Muitas tentativas de login. Tente novamente mais tarde.")

	// Health check — público, sem rate limit (usado por monitor/deploy).
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// --- Autenticação ---
	mux.Handle("POST /api/auth/register",
		middleware.Chain(http.HandlerFunc(d.Auth.Register), apiRL, csrf))
	mux.Handle("POST /api/auth/login",
		middleware.Chain(http.HandlerFunc(d.Auth.Login), apiRL, loginRL, csrf))
	mux.Handle("POST /api/auth/logout",
		middleware.Chain(http.HandlerFunc(d.Auth.Logout), apiRL, csrf, authRequired))
	mux.Handle("GET /api/auth/me",
		middleware.Chain(http.HandlerFunc(d.Auth.Me), apiRL, authRequired))

	// --- Incidentes (todas exigem sessão) ---
	mux.Handle("GET /api/incidents",
		middleware.Chain(http.HandlerFunc(d.Incidents.List), apiRL, authRequired))
	mux.Handle("POST /api/incidents",
		middleware.Chain(http.HandlerFunc(d.Incidents.Create), apiRL, csrf, authRequired))
	mux.Handle("GET /api/incidents/{id}",
		middleware.Chain(http.HandlerFunc(d.Incidents.Get), apiRL, authRequired))
	mux.Handle("PUT /api/incidents/{id}",
		middleware.Chain(http.HandlerFunc(d.Incidents.Update), apiRL, csrf, authRequired))
	mux.Handle("DELETE /api/incidents/{id}",
		middleware.Chain(http.HandlerFunc(d.Incidents.Delete), apiRL, csrf, authRequired))

	// Catch-all: resposta 404 em JSON, sem revelar rotas existentes.
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteError(w, http.StatusNotFound, "Recurso não encontrado.")
	})

	return middleware.Chain(mux,
		middleware.RequestID(),
		middleware.Recover(d.Log),
		middleware.RequestLogger(d.Log),
		middleware.SecurityHeaders(d.IsProduction),
		middleware.CORS(d.AllowedOrigin),
	)
}
