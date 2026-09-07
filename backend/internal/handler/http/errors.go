// Package httpapi contém os controllers HTTP: traduzem request↔usecase
// e não carregam regra de negócio (SRP; docs/plano.md §4).
package httpapi

import (
	"errors"
	"log/slog"
	"net/http"

	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/httpx"
)

// writeDomainError mapeia um erro de domínio para o status HTTP
// correspondente, sempre com mensagem genérica ao cliente. Erros
// inesperados viram 500 e são registrados com o request_id — o detalhe
// nunca chega ao cliente (OWASP A02/A10).
func writeDomainError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	switch {
	case errors.Is(err, errdomain.ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "Recurso não encontrado.")
	case errors.Is(err, errdomain.ErrForbidden):
		// 404 em vez de 403 para não confirmar a existência de um recurso
		// de outro usuário (evita enumeração de IDs; OWASP A01).
		httpx.WriteError(w, http.StatusNotFound, "Recurso não encontrado.")
	case errors.Is(err, errdomain.ErrEmailTaken):
		httpx.WriteError(w, http.StatusConflict, "E-mail já cadastrado.")
	case errors.Is(err, errdomain.ErrInvalidCredentials):
		httpx.WriteError(w, http.StatusUnauthorized, "Credenciais inválidas.")
	case errors.Is(err, errdomain.ErrValidation):
		httpx.WriteError(w, http.StatusUnprocessableEntity, "Dados inválidos.")
	default:
		log.Error("erro não tratado no handler",
			"err", err,
			"method", r.Method,
			"path", r.URL.Path,
			"request_id", middleware.RequestIDFrom(r.Context()),
		)
		httpx.WriteError(w, http.StatusInternalServerError,
			"Algo deu errado no servidor. Tente novamente.")
	}
}
