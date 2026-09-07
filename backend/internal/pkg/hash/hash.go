// Package hash encapsula o hashing de senha com bcrypt (OWASP A07 —
// Authentication Failures; docs/plano.md §5). A senha em texto plano
// nunca sai deste pacote e nunca é registrada em log.
package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hasher é a porta usada pelos usecases de autenticação (DIP/ISP).
type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// ErrMismatch é retornado quando a senha não confere.
var ErrMismatch = errors.New("hash: senha não confere")

// Bcrypt implementa Hasher usando golang.org/x/crypto/bcrypt.
type Bcrypt struct{ cost int }

// NewBcrypt cria um Hasher bcrypt. cost < bcrypt.MinCost cai para o custo
// padrão da biblioteca.
func NewBcrypt(cost int) *Bcrypt {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &Bcrypt{cost: cost}
}

// Hash devolve o hash bcrypt da senha.
func (b *Bcrypt) Hash(password string) (string, error) {
	// bcrypt trunca em 72 bytes; rejeitamos explicitamente para não dar
	// falsa sensação de força a senhas longas coladas.
	if len(password) > 72 {
		return "", errors.New("hash: senha excede 72 bytes")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// Compare confere a senha contra o hash. Retorna ErrMismatch em caso de
// divergência e o erro original para hashes corrompidos.
func (b *Bcrypt) Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatch
	}
	return err
}
