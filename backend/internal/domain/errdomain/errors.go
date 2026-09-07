// Package errdomain define os erros sentinela do domínio. As camadas
// externas (handlers HTTP) traduzem esses erros para status codes sem
// vazar detalhe interno ao cliente (OWASP A02 / A10; docs/plano.md §7).
package errdomain

import "errors"

var (
	// ErrNotFound: recurso inexistente.
	ErrNotFound = errors.New("recurso não encontrado")

	// ErrForbidden: recurso existe mas não pertence ao solicitante
	// (OWASP A01 — Broken Access Control).
	ErrForbidden = errors.New("acesso negado ao recurso")

	// ErrEmailTaken: e-mail já cadastrado.
	ErrEmailTaken = errors.New("e-mail já cadastrado")

	// ErrInvalidCredentials: e-mail ou senha inválidos. Mensagem única e
	// genérica para não permitir enumeração de usuários (OWASP A07).
	ErrInvalidCredentials = errors.New("credenciais inválidas")

	// ErrValidation: entrada malformada. Normalmente acompanhada de um
	// mapa campo→mensagens produzido pela camada de validação.
	ErrValidation = errors.New("dados inválidos")
)
