package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"incidenttrack/internal/pkg/httpx"
)

// maxBodyBytes limita o corpo de qualquer request JSON (defesa simples
// contra payloads absurdos; OWASP A05).
const maxBodyBytes = 1 << 20 // 1 MiB

// decodeJSON lê exatamente um objeto JSON do corpo para dst. Rejeita
// Content-Type inválido, campos desconhecidos, corpo vazio, JSON extra e
// payloads acima do limite. Em erro, já escreve a resposta 400/413/415 e
// devolve false.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		if mt := strings.TrimSpace(strings.SplitN(ct, ";", 2)[0]); mt != "application/json" {
			httpx.WriteError(w, http.StatusUnsupportedMediaType,
				"Content-Type deve ser application/json.")
			return false
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxErr):
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "Corpo da requisição muito grande.")
		case errors.Is(err, io.EOF):
			httpx.WriteError(w, http.StatusBadRequest, "Corpo da requisição vazio.")
		default:
			// Não repassa a mensagem crua do parser ao cliente.
			httpx.WriteError(w, http.StatusBadRequest, "JSON inválido.")
		}
		return false
	}

	// Garante que não há um segundo valor JSON depois do objeto.
	if dec.More() {
		httpx.WriteError(w, http.StatusBadRequest, "JSON inválido: conteúdo extra.")
		return false
	}
	return true
}
