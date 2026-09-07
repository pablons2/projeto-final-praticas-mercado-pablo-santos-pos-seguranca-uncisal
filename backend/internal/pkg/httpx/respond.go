// Package httpx centraliza a escrita de respostas JSON e o formato de
// erro da API. Manter isso num único lugar garante que nenhuma camada
// vaze stack trace ou detalhe interno ao cliente (OWASP A02/A10;
// docs/plano.md §7). O shape do erro espelha frontend/src/lib/http.ts:
//
//	{ "error": "mensagem amigável", "fields": { "campo": ["..."] } }
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorBody é o corpo padrão de erro.
type ErrorBody struct {
	Error  string              `json:"error"`
	Fields map[string][]string `json:"fields,omitempty"`
}

// WriteJSON serializa v como JSON com o status informado.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// resposta já parcialmente escrita; só registra.
		slog.Error("httpx: falha ao serializar resposta", "err", err)
	}
}

// WriteError responde com uma mensagem genérica e segura.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorBody{Error: msg})
}

// WriteValidationError responde 422 com o mapa de erros por campo.
func WriteValidationError(w http.ResponseWriter, msg string, fields map[string][]string) {
	WriteJSON(w, http.StatusUnprocessableEntity, ErrorBody{Error: msg, Fields: fields})
}

// NoContent responde 204 sem corpo.
func NoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }
